package secret

import (
	"fmt"
	"text/tabwriter"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/cli"
	"github.com/spf13/cobra"
)

func newListCommand(dockerCli *client.DockerCli) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Aliases: []string{"list"},
		Short:   "List secrets",
		Long:    listDescription,
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(dockerCli)
		},
	}

	return cmd
}

func runList(dockerCli *client.DockerCli) error {
	client := dockerCli.Client()
	secrets, err := client.SecretList(context.Background())
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(dockerCli.Out(), 20, 1, 3, ' ', 0)
	fmt.Fprintf(w, "ID \tNAME \tMOUNTPOINT")
	fmt.Fprintf(w, "\n")

	for _, s := range secrets {
		fmt.Fprintf(w, "%s\t%s\t%s\n", s.ID, s.Name, s.Mountpoint)
	}
	w.Flush()
	return nil
}

var listDescription = `

TODO

`
