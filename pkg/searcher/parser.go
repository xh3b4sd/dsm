package searcher

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/dsm/pkg/path"
)

type Config struct {
	FileSystem afero.Fs

	Name     string
	Resource string
	Source   string
}

type Searcher struct {
	fileSystem afero.Fs

	name     string
	resource string
	source   string
}

func New(config Config) (*Searcher, error) {
	if config.FileSystem == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.FileSystem must not be empty", config)
	}

	if config.Name == "" {
		return nil, tracer.Maskf(invalidConfigError, "%T.Name must not be empty", config)
	}
	if config.Resource == "" {
		return nil, tracer.Maskf(invalidConfigError, "%T.Resource must not be empty", config)
	}
	if config.Source == "" {
		return nil, tracer.Maskf(invalidConfigError, "%T.Source must not be empty", config)
	}

	s := &Searcher{
		fileSystem: config.FileSystem,

		name:     config.Name,
		resource: config.Resource,
		source:   config.Source,
	}

	return s, nil
}

func (s *Searcher) Search() (map[string][]byte, error) {
	files, err := s.files(".yaml")
	if err != nil {
		return nil, tracer.Mask(err)
	}

	filtered := map[string][]byte{}
	for p, b := range files {
		var newPath *path.Path
		{
			c := path.Config{
				Bytes: b,
			}

			newPath, err = path.New(c)
			if err != nil {
				return nil, tracer.Mask(err)
			}
		}

		{
			v, err := newPath.Get("metadata.name")
			if path.IsNotFound(err) {
				continue
			} else if err != nil {
				return nil, tracer.Mask(err)
			}

			if v != s.name {
				continue
			}
		}

		{
			v, err := newPath.Get("kind")
			if err != nil {
				return nil, tracer.Mask(err)
			}

			if v != s.resource {
				continue
			}
		}

		filtered[p] = b
	}

	return filtered, nil
}

func (s *Searcher) files(exts ...string) (map[string][]byte, error) {
	files := map[string][]byte{}
	{
		walkFunc := func(r string, i os.FileInfo, err error) error {
			if err != nil {
				return tracer.Mask(err)
			}

			if i.IsDir() && i.Name() == ".git" {
				return filepath.SkipDir
			}

			if i.IsDir() && i.Name() == ".github" {
				return filepath.SkipDir
			}

			// We do not want to track directories. We are interested in
			// directories containing specific files.
			if i.IsDir() {
				return nil
			}

			// We do not want to track files with the wrong extension. We are
			// interested in protocol buffer files having the ".proto"
			// extension.
			for _, e := range exts {
				if filepath.Ext(i.Name()) != e {
					return nil
				}
			}

			p := filepath.Join(filepath.Dir(r), i.Name())

			b, err := afero.ReadFile(s.fileSystem, p)
			if err != nil {
				return tracer.Mask(err)
			}

			for _, s := range bytes.Split(b, []byte("---")) {
				files[p] = s
			}

			return nil
		}

		err := afero.Walk(s.fileSystem, s.source, walkFunc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return files, nil
}
