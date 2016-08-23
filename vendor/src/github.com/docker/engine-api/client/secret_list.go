package client

import (
	"encoding/json"
	"net/url"

	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

// SecretList returns the secrets configured in the docker host.
func (cli *Client) SecretList(ctx context.Context) (types.SecretsListResponse, error) {
	var secrets types.SecretsListResponse
	query := url.Values{}

	resp, err := cli.get(ctx, "/secrets", query, nil)
	if err != nil {
		return secrets, err
	}

	err = json.NewDecoder(resp.body).Decode(&secrets)
	ensureReaderClosed(resp)
	return secrets, err
}
