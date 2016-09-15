package secret

import (
	"fmt"

	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/spf13/cobra"
)

type updateOptions struct{}

func newSecretUpdateCommand(dockerCli *command.DockerCli) *cobra.Command {
	return &cobra.Command{
		Use:   "update [name]",
		Short: "Update a secret",
		Args:  cli.RequiresMinArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := updateOptions{}
			return runSecretUpdate(dockerCli, opts)
		},
	}
}

func runSecretUpdate(dockerCli *command.DockerCli, opts updateOptions) error {
	fmt.Fprintln(dockerCli.Out(), "TODO: secret update")
	//client := dockerCli.Client()
	//ctx := context.Background()

	//if err := client.SwarmLeave(ctx, opts.force); err != nil {
	//	return err
	//}

	//fmt.Fprintln(dockerCli.Out(), "Node left the swarm.")
	return nil
}
