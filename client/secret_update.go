package client

import (
	"github.com/docker/docker/api/types/swarm"
	"golang.org/x/net/context"
)

// SecretUpdate updates a new Secret.
func (cli *Client) SecretUpdate(ctx context.Context, secret swarm.SecretSpec) error {
	var headers map[string][]string

	if _, err := cli.post(ctx, "/secrets/update", nil, secret, headers); err != nil {
		return err
	}

	return nil
}
