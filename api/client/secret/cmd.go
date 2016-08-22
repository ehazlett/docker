package secret

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/cli"
)

// NewSecretCommand returns a cobra command for `secret` subcommands
func NewSecretCommand(dockerCli *client.DockerCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret COMMAND",
		Short: "Manage Docker secrets",
		Long:  secretDescription,
		Args:  cli.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(dockerCli.Err(), "\n"+cmd.UsageString())
		},
	}
	cmd.AddCommand(
		newCreateCommand(dockerCli),
		newInspectCommand(dockerCli),
		newListCommand(dockerCli),
		newRemoveCommand(dockerCli),
	)
	return cmd
}

var secretDescription = `
TODO

To see help for a subcommand, use:

    docker secret CMD help

For full details on using docker secret visit Docker's online documentation.

`
