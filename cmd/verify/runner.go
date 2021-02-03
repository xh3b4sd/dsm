package verify

import (
	"context"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/dsm/pkg/path"
	"github.com/xh3b4sd/dsm/pkg/searcher"
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

	var s *searcher.Searcher
	{
		c := searcher.Config{
			FileSystem: afero.NewOsFs(),

			Name:     r.flag.Name,
			Resource: r.flag.Resource,
			Source:   r.flag.Source,
		}

		s, err = searcher.New(c)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	var f map[string][]byte
	{
		f, err = s.Search()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	if len(f) == 0 {
		return tracer.Mask(notFoundError)
	}

	var x interface{}
	for _, b := range f {
		var newPath *path.Path
		{
			c := path.Config{
				Bytes: b,
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

		if x == nil {
			x = v
		}

		if x != v {
			return tracer.Mask(invalidValueError)
		}
	}

	return nil
}
