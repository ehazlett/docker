package secret

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/server/httputils"
	enginetypes "github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

func (sr *secretRouter) createSecret(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	var s enginetypes.Secret
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		return err
	}
	return sr.backend.CreateSecret(s)
}

func (sr *secretRouter) updateSecret(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	id := vars["id"]

	var s enginetypes.Secret
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		return err
	}
	return sr.backend.UpdateSecret(id, &s)
}

func (sr *secretRouter) listSecrets(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	l, err := sr.backend.ListSecrets()
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, l)
}

func (sr *secretRouter) inspectSecret(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	result, err := sr.backend.InspectSecret(vars["id"])
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, result)
}

func (sr *secretRouter) removeSecret(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	id := vars["id"]
	return sr.backend.RemoveSecret(id)
}
