package secret

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

type updateOptions struct {
	name string
}

func newSecretUpdateCommand(dockerCli *command.DockerCli) *cobra.Command {
	return &cobra.Command{
		Use:   "update [name]",
		Short: "Update a secret using stdin as content",
		Args:  cli.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := updateOptions{
				name: args[0],
			}
			return runSecretUpdate(dockerCli, opts)
		},
	}
}

func runSecretUpdate(dockerCli *command.DockerCli, opts updateOptions) error {
	client := dockerCli.Client()
	ctx := context.Background()

	secretData, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("Error reading content from STDIN: %v", err)
	}

	spec := swarm.SecretSpec{
		Annotations: swarm.Annotations{
			Name: opts.name,
		},
		Type: swarm.ContainerSecret,
		Data: secretData,
	}

	if err := client.SecretUpdate(ctx, spec); err != nil {
		return err
	}

	fmt.Fprintln(dockerCli.Out(), opts.name)
	return nil
}
