package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/xh3b4sd/tracer"
)

type Config struct {
	FileSystem afero.Fs

	Source string
}

type Parser struct {
	fileSystem afero.Fs

	source string
}

func New(config Config) (*Parser, error) {
	if config.FileSystem == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.FileSystem must not be empty", config)
	}

	if config.Source == "" {
		return nil, tracer.Maskf(invalidConfigError, "%T.Source must not be empty", config)
	}

	p := &Parser{
		fileSystem: config.FileSystem,

		source: config.Source,
	}

	return p, nil
}

func (p *Parser) Search() error {
	files, err := p.files(".proto")
	if err != nil {
		return tracer.Mask(err)
	}

	for p, c := range files {
		fmt.Printf("%#v\n", p)
		fmt.Printf("%#v\n", c)
		fmt.Printf("\n")
	}

	return nil
}

func (p *Parser) files(exts ...string) (map[string][]string, error) {
	files := map[string][]string{}
	{
		walkFunc := func(p string, i os.FileInfo, err error) error {
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

			files[filepath.Dir(p)] = append(files[filepath.Dir(p)], filepath.Join(filepath.Dir(p), i.Name()))

			return nil
		}

		err := afero.Walk(p.fileSystem, p.source, walkFunc)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	return files, nil
}
