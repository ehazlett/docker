package secret

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/opts"
	"github.com/docker/engine-api/types"
	"github.com/spf13/cobra"
)

type createOptions struct {
	name       string
	driver     string
	driverOpts opts.MapOpts
	labels     []string
}

func newCreateCommand(dockerCli *client.DockerCli) *cobra.Command {
	opts := createOptions{
		driverOpts: *opts.NewMapOpts(nil, nil),
	}

	cmd := &cobra.Command{
		Use:   "create [OPTIONS]",
		Short: "Create a secret",
		Long:  createDescription,
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(dockerCli, opts)
		},
	}
	flags := cmd.Flags()
	flags.StringVar(&opts.name, "name", "", "Specify secret name")

	return cmd
}

func runCreate(dockerCli *client.DockerCli, opts createOptions) error {
	client := dockerCli.Client()

	volReq := types.SecretCreateRequest{
		Name: opts.name,
		Data: "TODO",
	}

	secret, err := client.CreateSecret(context.Background(), volReq)
	if err != nil {
		return err
	}

	fmt.Fprintf(dockerCli.Out(), "%s\n", secret.Name)
	return nil
}

var createDescription = `
TODO
`
