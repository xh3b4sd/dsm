package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/dsm/cmd/completion"
	"github.com/xh3b4sd/dsm/cmd/search"
	"github.com/xh3b4sd/dsm/cmd/update"
	"github.com/xh3b4sd/dsm/cmd/verify"
	"github.com/xh3b4sd/dsm/cmd/version"
	"github.com/xh3b4sd/dsm/pkg/project"
)

var (
	name  = project.Name()
	short = project.Description()
	long  = project.Description()
)

type Config struct {
	Logger logger.Interface
}

func New(config Config) (*cobra.Command, error) {
	if config.Logger == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	var err error

	var completionCmd *cobra.Command
	{
		c := completion.Config{
			Logger: config.Logger,
		}

		completionCmd, err = completion.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var searchCmd *cobra.Command
	{
		c := search.Config{
			Logger: config.Logger,
		}

		searchCmd, err = search.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var verifyCmd *cobra.Command
	{
		c := verify.Config{
			Logger: config.Logger,
		}

		verifyCmd, err = verify.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var versionCmd *cobra.Command
	{
		c := version.Config{
			Logger: config.Logger,
		}

		versionCmd, err = version.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var updateCmd *cobra.Command
	{
		c := update.Config{
			Logger: config.Logger,
		}

		updateCmd, err = update.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var c *cobra.Command
	{
		r := &runner{
			logger: config.Logger,
		}

		c = &cobra.Command{
			Use:   name,
			Short: short,
			Long:  long,
			RunE:  r.Run,
			// We slience errors because we do not want to see spf13/cobra printing.
			// The errors returned by the commands will be propagated to the main.go
			// anyway, where we have custom error printing for the command line
			// tool.
			SilenceErrors: true,
			SilenceUsage:  true,
		}

		c.AddCommand(completionCmd)
		c.AddCommand(searchCmd)
		c.AddCommand(verifyCmd)
		c.AddCommand(updateCmd)
		c.AddCommand(versionCmd)
	}

	return c, nil
}
