package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/docker/engine-api/types/secret"

	"golang.org/x/net/context"
)

// SecretInspect returns the information about a specific secret
func (cli *Client) SecretInspect(ctx context.Context, secretID string) (secret.Secret, error) {
	secret, _, err := cli.SecretInspectWithRaw(ctx, secretID)
	return secret, err
}

// SecretInspectWithRaw returns the information about a specific secret in the docker host and its raw representation
func (cli *Client) SecretInspectWithRaw(ctx context.Context, secretID string) (secret.Secret, []byte, error) {
	var secret secret.Secret
	resp, err := cli.get(ctx, "/secrets/"+secretID, nil, nil)
	if err != nil {
		if resp.statusCode == http.StatusNotFound {
			return secret, nil, secretNotFoundError{secretID}
		}
		return secret, nil, err
	}
	defer ensureReaderClosed(resp)

	body, err := ioutil.ReadAll(resp.body)
	if err != nil {
		return secret, nil, err
	}
	rdr := bytes.NewReader(body)
	err = json.NewDecoder(rdr).Decode(&secret)
	return secret, body, err
}
