package secret

import (
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/cli"
	"github.com/docker/docker/cli/command"
	"github.com/docker/docker/opts"
	"github.com/spf13/cobra"
)

type listOptions struct {
	filter opts.FilterOpt
	quiet  bool
}

func newSecretListCommand(dockerCli *command.DockerCli) *cobra.Command {
	opts := listOptions{}

	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List secrets",
		Args:  cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := listOptions{}
			return runSecretList(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display IDs")
	flags.VarP(&opts.filter, "filter", "f", "Filter output based on conditions provided")

	return cmd
}

func runSecretList(dockerCli *command.DockerCli, opts listOptions) error {
	client := dockerCli.Client()
	ctx := context.Background()

	secrets, err := client.SecretList(ctx, types.SecretListOptions{Filter: opts.filter.Value()})
	if err != nil {
		return err
	}

	// TODO (ejh): quiet
	w := tabwriter.NewWriter(dockerCli.Out(), 20, 1, 3, ' ', 0)
	fmt.Fprintf(w, "NAME\tCREATED\tUPDATED\tVERSIONS")
	fmt.Fprintf(w, "\n")

	// TODO (ejh): pretty timestamps
	for _, s := range secrets {
		fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", s.Name, s.Meta.CreatedAt, s.Meta.UpdatedAt, len(s.SecretData))
	}
	w.Flush()

	return nil
}
