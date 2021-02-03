package verify

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	name  = "verify"
	short = "Verify the consistency of values within YAML or JSON data structures."
	long  = `Verify the consistency of values within YAML or JSON data structures.
Consider multiple of the following HelmRelease CRs defining a Docker image
tag in its spec.

    apiVersion: "helm.toolkit.fluxcd.io/v2beta1"
    kind: "HelmRelease"
    metadata:
      name: "apiserver"
    spec:
      values:
        image:
          tag: "8469445410f8a74d72af0cf430ed8dd44fb6b8fa"

The following example shows how to validate the image tag of all the YAML
files. In case the validation succeeds the command exits with an exit code 0
and without output. If the validation fails the command exits with an exit
code 1 and with the output of a stack trace.

    $ dsm verify -r HelmRelease -n apiserver -k spec.values.image.tag
    program panic

    {
        "kind": "invalidValueError",
        "stck": [
            "/Users/xh3b4sd/project/xh3b4sd/dsm/cmd/verify/runner.go:87",
            "/Users/xh3b4sd/project/xh3b4sd/dsm/cmd/verify/runner.go:30",
            "/Users/xh3b4sd/project/xh3b4sd/dsm/main.go:47"
        ],
        "type": "*tracer.Error"
    }

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
