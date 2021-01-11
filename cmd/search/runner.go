package search

import (
	"context"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/dsm/pkg/parser"
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

			Source: r.flag.Source,
		}

		p, err = parser.New(c)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		err = p.Search()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
