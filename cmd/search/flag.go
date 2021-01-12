package search

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"
)

type flag struct {
	Key      string
	Name     string
	Resource string
	Source   string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.Key, "key", "k", "", "JSON path key to work with.")
	cmd.Flags().StringVarP(&f.Name, "name", "n", "", "Metadata name of the resources to work with.")
	cmd.Flags().StringVarP(&f.Resource, "resource", "r", "", "Resource kind to work with.")
	cmd.Flags().StringVarP(&f.Source, "source", "s", ".", "Source directory to traverse.")
}

func (f *flag) Validate() error {
	{
		if f.Key == "" {
			return tracer.Maskf(invalidFlagError, "-k/--key must not be empty")
		}
	}

	{
		if f.Name == "" {
			return tracer.Maskf(invalidFlagError, "-n/--name must not be empty")
		}
	}

	{
		if f.Resource == "" {
			return tracer.Maskf(invalidFlagError, "-r/--resource must not be empty")
		}
	}

	{
		if f.Source == "" {
			return tracer.Maskf(invalidFlagError, "-s/--source must not be empty")
		}
	}

	return nil
}
