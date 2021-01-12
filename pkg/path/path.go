package path

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"

	yamltojson "github.com/ghodss/yaml"
	"github.com/spf13/cast"
	"github.com/xh3b4sd/tracer"
	yaml "gopkg.in/yaml.v2"
)

const (
	escapedSeparatorPlaceholder = "%%PLACEHOLDER%%"
)

var (
	placeholderExpression = regexp.MustCompile(escapedSeparatorPlaceholder)
)

type Config struct {
	Bytes     []byte
	Separator string
}

type Path struct {
	isJSON                     bool
	jsonBytes                  []byte
	jsonStructure              interface{}
	escapedSeparatorExpression *regexp.Regexp
	separatorExpression        *regexp.Regexp

	separator string
}

func New(config Config) (*Path, error) {
	if config.Bytes == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Bytes must not be empty", config)
	}
	if config.Separator == "" {
		config.Separator = "."
	}

	var err error

	var isJSON bool
	var jsonBytes []byte
	var jsonStructure interface{}
	{
		jsonBytes, isJSON, err = toJSON(config.Bytes)
		if err != nil {
			return nil, tracer.Mask(err)
		}

		err := json.Unmarshal(jsonBytes, &jsonStructure)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	p := &Path{
		isJSON:                     isJSON,
		jsonBytes:                  jsonBytes,
		jsonStructure:              jsonStructure,
		escapedSeparatorExpression: regexp.MustCompile(fmt.Sprintf(`\\%s`, config.Separator)),
		separatorExpression:        regexp.MustCompile(fmt.Sprintf(`\%s`, config.Separator)),

		separator: config.Separator,
	}

	return p, nil
}

// All returns all paths found in the configured JSON structure.
func (p *Path) All() ([]string, error) {
	paths, err := p.allFromInterface(p.jsonStructure)
	if err != nil {
		return nil, tracer.Mask(err)
	}

	sort.Strings(paths)

	return paths, nil
}

// Get returns the value found under the given path, if any.
func (p *Path) Get(path string) (interface{}, error) {
	value, err := p.getFromInterface(p.escapeKey(path), p.jsonStructure)
	if err != nil {
		return nil, tracer.Mask(err)
	}

	return value, nil
}

func (p *Path) OutputBytes() ([]byte, error) {
	b := p.jsonBytes
	if !p.isJSON {
		var err error
		b, err = yamltojson.JSONToYAML(b)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return b, nil
}

// Set changes the value of the given path.
func (p *Path) Set(path string, value interface{}) error {
	var err error

	p.jsonStructure, err = p.setFromInterface(p.escapeKey(path), value, p.jsonStructure)
	if err != nil {
		return tracer.Mask(err)
	}

	b, err := json.MarshalIndent(p.jsonStructure, "", "  ")
	if err != nil {
		return tracer.Mask(err)
	}
	p.jsonBytes = b

	return nil
}

func (p *Path) Validate(paths []string) error {
	all, err := p.All()
	if err != nil {
		return tracer.Mask(err)
	}

	var trimmedAll []string
	for _, service := range all {
		pv := strings.Split(service, ".")
		trimmedAll = append(trimmedAll, pv[len(pv)-1])
	}

	for _, p := range paths {
		fields := trimmedAll
		if strings.Index(p, ".") != -1 { // nolint:gosimple
			fields = all
		}
		if containsString(fields, p) {
			continue
		}

		return tracer.Maskf(notFoundError, "path '%s'", p)
	}

	return nil
}

func (p *Path) allFromInterface(value interface{}) ([]string, error) {
	// process map
	{
		stringMap, err := cast.ToStringMapE(value)
		if err != nil {
			// fall through
		} else {
			var paths []string

			for k, v := range stringMap {
				var l []string
				if reflect.TypeOf(v).String() != "string" {
					l, err = p.allFromInterface(v)
					if err != nil {
						return nil, tracer.Mask(err)
					}
				}

				k := p.separatorExpression.ReplaceAllString(k, fmt.Sprintf(`\%s`, p.separator))

				if l != nil { // nolint:gosimple
					for _, v := range l {
						paths = append(paths, fmt.Sprintf("%s%s%s", k, p.separator, v))
					}
				} else {
					paths = append(paths, k)
				}
			}

			return paths, nil
		}
	}

	// process slice
	{
		slice, err := cast.ToSliceE(value)
		if err != nil {
			// fall through
		} else {
			var paths []string

			for i, v := range slice {
				l, err := p.allFromInterface(v)
				if err != nil {
					return nil, tracer.Mask(err)
				}

				for _, v := range l {
					paths = append(paths, fmt.Sprintf("[%d]%s%s", i, p.separator, v))
				}
			}

			return paths, nil
		}
	}

	// process string
	{
		str, err := cast.ToStringE(value)
		if err != nil {
			// fall through
		} else if str == "" {
			// fall through
		} else {
			jsonBytes, _, err := toJSON([]byte(str))
			if err != nil {
				// fall through
			} else {
				var jsonStructure interface{}
				err := json.Unmarshal(jsonBytes, &jsonStructure)
				if err != nil {
					return nil, tracer.Mask(err)
				}

				l, err := p.allFromInterface(jsonStructure)
				if err != nil {
					return nil, tracer.Mask(err)
				}

				return l, nil
			}
		}
	}

	return nil, nil
}

func (p *Path) escapeKey(key string) string {
	return p.escapedSeparatorExpression.ReplaceAllString(key, escapedSeparatorPlaceholder)
}

func (p *Path) getFromInterface(path string, jsonStructure interface{}) (interface{}, error) {
	split := strings.Split(path, p.separator)
	key := p.unescapeKey(split[0])

	// process map
	{
		stringMap, err := cast.ToStringMapE(jsonStructure)
		if err != nil {
			// fall through
		} else {
			value, ok := stringMap[key]
			if ok {
				if len(split) == 1 {
					return value, nil
				} else {
					recPath := strings.Join(split[1:], p.separator)

					v, err := p.getFromInterface(recPath, value)
					if err != nil {
						return nil, tracer.Mask(err)
					}

					return v, nil
				}
			} else {
				return nil, tracer.Maskf(notFoundError, "key '%s'", path)
			}
		}
	}

	// process slice
	{
		slice, err := cast.ToSliceE(jsonStructure)
		if err != nil {
			// fall through
		} else {
			index, err := indexFromKey(key)
			if err != nil {
				return nil, tracer.Mask(err)
			}

			if index >= len(slice) {
				return nil, tracer.Maskf(notFoundError, "key '%s'", key)
			}
			recPath := strings.Join(split[1:], p.separator)

			v, err := p.getFromInterface(recPath, slice[index])
			if err != nil {
				return nil, tracer.Mask(err)
			}

			return v, nil
		}
	}

	// process string
	{
		str, err := cast.ToStringE(jsonStructure)
		if err != nil {
			// fall through
		} else {
			jsonBytes, _, err := toJSON([]byte(str))
			if err != nil {
				// fall through
			} else {
				var jsonStructure interface{}
				err := json.Unmarshal(jsonBytes, &jsonStructure)
				if err != nil {
					return nil, tracer.Mask(err)
				}

				v, err := p.getFromInterface(path, jsonStructure)
				if err != nil {
					return nil, tracer.Mask(err)
				}

				return v, nil
			}
		}
	}

	return nil, nil
}

func (p *Path) setFromInterface(path string, value interface{}, jsonStructure interface{}) (interface{}, error) {
	split := strings.Split(path, p.separator)
	key := p.unescapeKey(split[0])

	// Create new element when the existing jsonStructure doesn't exist.
	if jsonStructure == nil {
		m := make(map[string]interface{})

		// Just recurse when there are more components left in path with
		// missing elements.
		if len(split) > 1 {
			var err error
			recPath := strings.Join(split[1:], p.separator)
			value, err = p.setFromInterface(recPath, value, nil)
			if err != nil {
				return nil, tracer.Mask(err)
			}
		}

		m[key] = value

		return m, nil
	}

	// process map
	{
		_, ok := jsonStructure.(string)
		if ok {
			// Fall through in case our received JSON structure is actually a string.
			// cast.ToStringMapE was working as expected until
			// https://github.com/spf13/cast/pull/59, so we have to make sure we do
			// not call cast.ToStringMapE only if we do not have an actual string,
			// because cast.ToStringMapE would now accept the string instead of
			// returning an error like it did before.
		} else {
			stringMap, err := cast.ToStringMapE(jsonStructure)
			if err != nil {
				// fall through
			} else {
				if len(split) == 1 {
					stringMap[path] = value
					return stringMap, nil
				} else {
					recPath := strings.Join(split[1:], p.separator)

					modified, err := p.setFromInterface(recPath, value, stringMap[key])
					if err != nil {
						return nil, tracer.Mask(err)
					}
					stringMap[key] = modified

					return stringMap, nil
				}
			}
		}
	}

	// process slice
	{
		slice, err := cast.ToSliceE(jsonStructure)
		if err != nil {
			// fall through
		} else {
			index, err := indexFromKey(key)
			if err != nil {
				return nil, tracer.Mask(err)
			}

			if index >= len(slice) {
				return nil, tracer.Maskf(notFoundError, "key '%s'", key)
			}
			recPath := strings.Join(split[1:], p.separator)

			modified, err := p.setFromInterface(recPath, value, slice[index])
			if err != nil {
				return nil, tracer.Mask(err)
			}
			slice[index] = modified

			return slice, nil
		}
	}

	// process string
	{
		str, err := cast.ToStringE(jsonStructure)
		if err != nil {
			// fall through
		} else {
			jsonBytes, isJSON, err := toJSON([]byte(str))
			if err != nil {
				// fall through
			} else {
				var jsonStructure interface{}
				err := json.Unmarshal(jsonBytes, &jsonStructure)
				if err != nil {
					return nil, tracer.Mask(err)
				}

				modified, err := p.setFromInterface(path, value, jsonStructure)
				if err != nil {
					return nil, tracer.Mask(err)
				}

				var b []byte
				if !isJSON {
					b, err = yamltojson.Marshal(modified)
					if err != nil {
						return nil, tracer.Mask(err)
					}
				} else {
					b, err = json.MarshalIndent(modified, "", "  ")
					if err != nil {
						return nil, tracer.Mask(err)
					}
				}

				return string(b), nil
			}
		}
	}

	return nil, nil
}

func (p *Path) unescapeKey(key string) string {
	return placeholderExpression.ReplaceAllString(key, p.separator)
}

func containsString(list []string, item string) bool {
	for _, l := range list {
		if l == item {
			return true
		}
	}

	return false
}

func indexFromKey(key string) (int, error) {
	re := regexp.MustCompile(`\[[0-9]+\]`)
	ok := re.MatchString(key)
	if !ok {
		return 0, tracer.Maskf(invalidFormatError, key)
	}

	p := key[1 : len(key)-1]
	i, err := strconv.Atoi(p)
	if err != nil {
		return 0, tracer.Mask(err)
	}

	return i, nil
}

func isJSON(b []byte) bool {
	var l []interface{}
	isList := json.Unmarshal(b, &l) == nil

	var m map[string]interface{}
	isObject := json.Unmarshal(b, &m) == nil

	return isObject || isList
}

func isYAMLList(b []byte) bool {
	var l []interface{}
	return yaml.Unmarshal(b, &l) == nil && bytes.HasPrefix(b, []byte("-"))
}

func isYAMLObject(b []byte) bool {
	var m map[interface{}]interface{}
	return yaml.Unmarshal(b, &m) == nil && !bytes.HasPrefix(b, []byte("-"))
}

func toJSON(b []byte) ([]byte, bool, error) {
	if isJSON(b) {
		return b, true, nil
	}

	isYAMLList := isYAMLList(b)
	isYAMLObject := isYAMLObject(b)

	var jsonBytes []byte
	if isYAMLList && !isYAMLObject {
		var jsonList []interface{}
		err := yamltojson.Unmarshal(b, &jsonList)
		if err != nil {
			return nil, false, tracer.Mask(err)
		}

		jsonBytes, err = json.Marshal(jsonList)
		if err != nil {
			return nil, false, tracer.Mask(err)
		}

		return jsonBytes, false, nil
	}

	if !isYAMLList && isYAMLObject {
		var jsonMap map[string]interface{}
		err := yamltojson.Unmarshal(b, &jsonMap)
		if err != nil {
			return nil, false, tracer.Mask(err)
		}

		jsonBytes, err = json.Marshal(jsonMap)
		if err != nil {
			return nil, false, tracer.Mask(err)
		}

		return jsonBytes, false, nil
	}

	return nil, false, tracer.Mask(invalidFormatError)
}
