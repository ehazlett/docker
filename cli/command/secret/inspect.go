package secret

import (
	"fmt"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

func newSecretInspectCommand(dockerCli *command.DockerCli) *cobra.Command {
	return &cobra.Command{
		Use:   "inspect [name]",
		Short: "Inspect a secret",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSecretInspect(dockerCli)
		},
	}
}

func runSecretInspect(dockerCli *command.DockerCli) error {
	fmt.Fprintln(dockerCli.Out(), "TODO: secret inspect")
	//client := dockerCli.Client()
	//ctx := context.Background()

	//if err := client.SwarmLeave(ctx, opts.force); err != nil {
	//	return err
	//}

	//fmt.Fprintln(dockerCli.Out(), "Node left the swarm.")
	return nil
}
