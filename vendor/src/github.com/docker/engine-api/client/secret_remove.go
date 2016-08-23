package client

import (
	"net/url"

	"golang.org/x/net/context"
)

// SecretRemove removes a secret from the docker host.
func (cli *Client) SecretRemove(ctx context.Context, secretID string) error {
	query := url.Values{}
	resp, err := cli.delete(ctx, "/secrets/"+secretID, query, nil)
	ensureReaderClosed(resp)
	return err
}
