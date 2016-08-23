package secret

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/cli"
	"github.com/docker/engine-api/types"
	"github.com/spf13/cobra"
)

type createOptions struct {
	name       string
	mountpoint string
}

func newCreateCommand(dockerCli *client.DockerCli) *cobra.Command {
	opts := createOptions{}

	cmd := &cobra.Command{
		Use:   "create [OPTIONS] NAME MOUNTPOINT",
		Short: "Create a secret",
		Long:  createDescription,
		Args:  cli.RequiresMinArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.name = args[0]
			opts.mountpoint = args[1]
			return runCreate(dockerCli, opts)
		},
	}

	return cmd
}

func runCreate(dockerCli *client.DockerCli, opts createOptions) error {
	client := dockerCli.Client()
	r := bufio.NewReader(os.Stdin)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	secretReq := types.SecretCreateRequest{
		Name:       opts.name,
		Mountpoint: opts.mountpoint,
		Data:       string(data),
	}

	secret, err := client.SecretCreate(context.Background(), secretReq)
	if err != nil {
		return err
	}

	fmt.Fprintf(dockerCli.Out(), "%s\n", secret.Name)
	return nil
}

var createDescription = `
TODO
`
