package update

import (
	"context"
	"io/ioutil"

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

	err := r.run(ctx, cmd, args)
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

	var m map[string][]byte
	{
		m, err = s.Search()
		if err != nil {
			return tracer.Mask(err)
		}
	}

	for p, b := range m {
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

		err := newPath.Set(r.flag.Key, r.flag.Value)
		if err != nil {
			return tracer.Mask(err)
		}

		v, err := newPath.OutputBytes()
		if err != nil {
			return tracer.Mask(err)
		}

		err = ioutil.WriteFile(p, v, 0600)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	return nil
}
