package parser

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/xh3b4sd/dsm/pkg/path"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	FileSystem afero.Fs

	Name     string
	Resource string
	Source   string
}

type Parser struct {
	fileSystem afero.Fs

	name     string
	resource string
	source   string
}

func New(config Config) (*Parser, error) {
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

	p := &Parser{
		fileSystem: config.FileSystem,

		name:     config.Name,
		resource: config.Resource,
		source:   config.Source,
	}

	return p, nil
}

func (p *Parser) Search() (map[string][]byte, error) {
	files, err := p.files(".yaml")
	if err != nil {
		return nil, tracer.Mask(err)
	}

	filtered := map[string][]byte{}
	for f, c := range files {
		var newPath *path.Path
		{
			c := path.Config{
				Bytes: c,
			}

			newPath, err = path.New(c)
			if err != nil {
				return nil, tracer.Mask(err)
			}
		}

		{
			v, err := newPath.Get("metadata.name")
			if err != nil {
				return nil, tracer.Mask(err)
			}

			if v != p.name {
				continue
			}
		}

		{
			v, err := newPath.Get("kind")
			if err != nil {
				return nil, tracer.Mask(err)
			}

			if v != p.resource {
				continue
			}
		}

		filtered[f] = c
	}

	return filtered, nil
}

func (p *Parser) files(exts ...string) (map[string][]byte, error) {
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

			f := filepath.Join(filepath.Dir(r), i.Name())

			b, err := afero.ReadFile(p.fileSystem, f)
			if err != nil {
				return tracer.Mask(err)
			}

			files[f] = b

			return nil
		}

		err := afero.Walk(p.fileSystem, p.source, walkFunc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return files, nil
}
