package secret

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

type removeOptions struct {
	names []string
}

func newSecretRemoveCommand(dockerCli *command.DockerCli) *cobra.Command {
	return &cobra.Command{
		Use:   "rm [name]",
		Short: "Remove a secret",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := removeOptions{
				names: args,
			}
			return runSecretRemove(dockerCli, opts)
		},
	}
}

func runSecretRemove(dockerCli *command.DockerCli, opts removeOptions) error {
	client := dockerCli.Client()
	ctx := context.Background()

	for _, name := range opts.names {
		var err error
		tokens := strings.Split(name, "@")
		switch {
		case len(tokens) == 1:
			err = client.SecretRemove(ctx, name, "")
		case len(tokens) == 2 && tokens[1] != "":
			err = client.SecretRemove(ctx, tokens[0], tokens[1])
		default:
			return fmt.Errorf("invalid secret name: %s", name)
		}
		if err != nil {
			return err
		}
		fmt.Fprintln(dockerCli.Out(), name)
	}

	return nil
}
