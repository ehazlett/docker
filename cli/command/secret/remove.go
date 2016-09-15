package secret

import (
	"fmt"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

func newSecretRemoveCommand(dockerCli *command.DockerCli) *cobra.Command {
	return &cobra.Command{
		Use:   "remove [name]",
		Short: "Remove a secret",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecretRemove(dockerCli)
		},
	}
}

func runSecretRemove(dockerCli *command.DockerCli) error {
	fmt.Fprintln(dockerCli.Out(), "TODO: secret remove")
	//client := dockerCli.Client()
	//ctx := context.Background()

	//if err := client.SwarmLeave(ctx, opts.force); err != nil {
	//	return err
	//}

	//fmt.Fprintln(dockerCli.Out(), "Node left the swarm.")
	return nil
}
