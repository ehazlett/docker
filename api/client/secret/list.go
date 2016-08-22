package secret

import (
	"sort"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/client"
	"github.com/docker/docker/api/client/formatter"
	"github.com/docker/docker/cli"
	"github.com/docker/engine-api/types"
	"github.com/spf13/cobra"
)

type bySecretName []*types.Secret

func (r bySecretName) Len() int      { return len(r) }
func (r bySecretName) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r bySecretName) Less(i, j int) bool {
	return r[i].Name < r[j].Name
}

type listOptions struct {
	quiet  bool
	format string
	filter []string
}

func newListCommand(dockerCli *client.DockerCli) *cobra.Command {
	var opts listOptions

	cmd := &cobra.Command{
		Use:     "ls [OPTIONS]",
		Aliases: []string{"list"},
		Short:   "List secrets",
		Long:    listDescription,
		Args:    cli.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(dockerCli, opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVarP(&opts.quiet, "quiet", "q", false, "Only display secret names")
	flags.StringVar(&opts.format, "format", "", "Pretty-print secrets using a Go template")

	return cmd
}

func runList(dockerCli *client.DockerCli, opts listOptions) error {
	client := dockerCli.Client()

	secrets, err := client.SecretList(context.Background())
	if err != nil {
		return err
	}

	f := opts.format
	if len(f) == 0 {
		if len(dockerCli.ConfigFile().SecretsFormat) > 0 && !opts.quiet {
			f = dockerCli.ConfigFile().SecretsFormat
		} else {
			f = "table"
		}
	}

	sort.Sort(bySecretName(secrets.Secrets))

	secretCtx := formatter.SecretContext{
		Context: formatter.Context{
			Output: dockerCli.Out(),
			Format: f,
			Quiet:  opts.quiet,
		},
		Secrets: secrets.Secrets,
	}

	secretCtx.Write()

	return nil
}

var listDescription = `

TODO

`
