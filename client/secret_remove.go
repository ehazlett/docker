package client

import (
	"net/url"

	"golang.org/x/net/context"
)

// SecretRemove removes a Secret.
func (cli *Client) SecretRemove(ctx context.Context, name, version string) error {
	v := url.Values{}
	if version != "" {
		v.Add("ver", version)
	}
	resp, err := cli.delete(ctx, "/secrets/"+name, v, nil)
	ensureReaderClosed(resp)
	return err
}
