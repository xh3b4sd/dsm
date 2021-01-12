package update

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	name  = "update"
	short = "Update values within YAML or JSON data structures."
	long  = `Update values within YAML or JSON data structures. Consider the following HelmRelease CR
defining a Docker image tag in its spec

    apiVersion: "helm.toolkit.fluxcd.io/v2beta1"
    kind: "HelmRelease"
    metadata:
      name: "apiserver"
    spec:
      values:
        image:
          tag: "8469445410f8a74d72af0cf430ed8dd44fb6b8fa"

The following example shows how to modify the image tag of the YAML file.

    dsm update -r HelmRelease -n apiserver -k spec.values.image.tag -v <new-sha>
`
)

type Config struct {
	Logger logger.Interface
}

func New(config Config) (*cobra.Command, error) {
	if config.Logger == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	var c *cobra.Command
	{
		f := &flag{}

		r := &runner{
			flag:   f,
			logger: config.Logger,
		}

		c = &cobra.Command{
			Use:   name,
			Short: short,
			Long:  long,
			RunE:  r.Run,
		}

		f.Init(c)
	}

	return c, nil
}
