package search

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/dsm/pkg/parser"
	"github.com/xh3b4sd/dsm/pkg/path"
)

type runner struct {
	flag   *flag
	logger logger.Interface
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return tracer.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error

	var p *parser.Parser
	{
		c := parser.Config{
			FileSystem: afero.NewOsFs(),

			Name:     r.flag.Name,
			Resource: r.flag.Resource,
			Source:   r.flag.Source,
		}

		p, err = parser.New(c)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	var m map[string][]byte
	{
		m, err = p.Search()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	for _, c := range m {
		var newPath *path.Path
		{
			c := path.Config{
				Bytes: c,
			}

			newPath, err = path.New(c)
			if err != nil {
				return tracer.Mask(err)
			}
		}

		v, err := newPath.Get(r.flag.Key)
		if err != nil {
			return tracer.Mask(err)
		}

		fmt.Printf("%s\n", v)
	}

	return nil
}
