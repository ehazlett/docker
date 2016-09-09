package secretsapi

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/api/sorting"
	"github.com/docker/swarmkit/identity"
	"github.com/docker/swarmkit/manager/state/store"

	"golang.org/x/net/context"
)

// MaxSecretSize is the maximum size of the data for any secret
const MaxSecretSize int = 500 * 1024

var (
	errNotImplemented  = errors.New("not implemented")
	errInvalidArgument = errors.New("invalid argument")
	isValidName        = regexp.MustCompile(`^[a-zA-Z0-9](?:[-_.]*[A-Za-z0-9]+)*$`)
)

// Server is the Secrets API gRPC server.
type Server struct {
	memstore *store.MemoryStore
}

// NewServer creates a Secrets API server.
func NewServer(memstore *store.MemoryStore) *Server {
	return &Server{
		memstore: memstore,
	}
}

// assumes spec is not nil
func secretDataFromSecretSpec(spec *api.SecretSpec) *api.SecretData {
	checksumBytes := sha256.Sum256(spec.Data)
	return &api.SecretData{
		ID:         identity.NewID(),
		Spec:       *spec,
		SecretSize: int32(len(spec.Data)),
		Digest:     "sha256:" + hex.EncodeToString(checksumBytes[:]),
	}
}

// zeros out the data on the given secret
func cleanSecret(secret *api.Secret) {
	for _, secretData := range secret.SecretData {
		secretData.Spec.Data = nil
	}
}

// CreateSecret creates and return a Secret based on the provided SecretSpec.
// - Returns `InvalidArgument` if the SecretSpec is malformed.
// - Returns `AlreadyExists` if the Secret's name conflicts.
// - Returns an error if the creation fails.
func (s *Server) CreateSecret(ctx context.Context, request *api.CreateSecretRequest) (*api.CreateSecretResponse, error) {
	if err := validateSecretSpec(request.Spec); err != nil {
		return nil, err
	}

	// creates a secret object and try to insert it into the store - the store will handle name conflicts
	secretData := secretDataFromSecretSpec(request.Spec.Copy())
	storedSecret := api.Secret{
		ID: identity.NewID(),
		SecretData: map[string]*api.SecretData{
			secretData.ID: secretData,
		},
		Name:          secretData.Spec.Annotations.Name,
		LatestVersion: secretData.ID,
	}

	createSecretFunc := func(tx store.Tx) error {
		return store.CreateSecret(tx, &storedSecret)
	}

	if err := s.memstore.Update(createSecretFunc); err != nil {
		if err == store.ErrNameConflict {
			return nil, grpc.Errorf(codes.AlreadyExists, "secret %s already exists", secretData.Spec.Annotations.Name)
		}
		return nil, err
	}

	cleanSecret(&storedSecret)
	return &api.CreateSecretResponse{
		Secret: &storedSecret,
	}, nil
}

// GetSecret returns a secret with the given name
// - Returns `NotFound` if the Secret with the given name is not found.
// - Returns `InvalidArgument` if the name is empty.
// - Returns an error if getting fails.
func (s *Server) GetSecret(ctx context.Context, request *api.GetSecretRequest) (*api.GetSecretResponse, error) {
	if request.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "errInvalidArgument.Error()")
	}

	var (
		secrets []*api.Secret
		err     error
	)

	s.memstore.View(func(tx store.ReadTx) {
		secrets, err = store.FindSecrets(tx, store.ByName(request.Name))
	})

	if err != nil {
		return nil, err
	}

	switch len(secrets) {
	case 0:
		return nil, grpc.Errorf(codes.NotFound, "secret %s not found", request.Name)
	case 1:
		cleanSecret(secrets[0])
		return &api.GetSecretResponse{Secret: secrets[0]}, nil
	default:
		return nil, fmt.Errorf("more than one secret with name %s", request.Name)
	}
}

// UpdateSecret adds a SecretSpec to a Secret as a new version.
// - Returns `NotFound` if the Secret with the given name is not found.
// - Returns `InvalidArgument` if the ServiceSpec is malformed.
// - Returns an error if the update fails.
func (s *Server) UpdateSecret(ctx context.Context, request *api.UpdateSecretRequest) (*api.UpdateSecretResponse, error) {
	if err := validateSecretSpec(request.Spec); err != nil {
		return nil, err
	}

	var secret *api.Secret

	// get the secret by name
	updateSecretFunc := func(tx store.Tx) error {
		secrets, err := store.FindSecrets(tx, store.ByName(request.Spec.Annotations.Name))
		if err != nil {
			return err
		}

		if len(secrets) == 0 {
			return store.ErrNotExist
		}
		if len(secrets) > 1 {
			return fmt.Errorf("more than one secret with name %s", request.Spec.Annotations.Name)
		}

		secret = secrets[0]

		// convert to secret data and add a new version
		secretData := secretDataFromSecretSpec(request.Spec.Copy())
		secret.SecretData[secretData.ID] = secretData
		secret.LatestVersion = secretData.ID
		return store.UpdateSecret(tx, secret)
	}

	err := s.memstore.Update(updateSecretFunc)
	switch {
	case err == store.ErrNotExist:
		return nil, grpc.Errorf(codes.NotFound, "secret %s not found", request.Spec.Annotations.Name)
	case err == store.ErrExist:
		return nil, grpc.Errorf(codes.AlreadyExists, "this update does not change the latest version of secret %s", request.Spec.Annotations.Name)
	case err != nil:
		return nil, err
	default:
		cleanSecret(secret)
		return &api.UpdateSecretResponse{Secret: secret}, nil
	}
}

// RemoveSecret removes a Secret referenced by name or a version of the Secret referenced by name and version.
// - Returns `InvalidArgument` if name is not provided.
// - Returns `NotFound` if the Secret is not found.
// - Returns an error if the deletion fails.
func (s *Server) RemoveSecret(ctx context.Context, request *api.RemoveSecretRequest) (*api.RemoveSecretResponse, error) {
	if request.Name == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "errInvalidArgument.Error()")
	}
	// get the secret by name
	deleteSecretFunc := func(tx store.Tx) error {
		secrets, err := store.FindSecrets(tx, store.ByName(request.Name))
		if err != nil {
			return err
		}

		if len(secrets) == 0 {
			return store.ErrNotExist
		}
		if len(secrets) > 1 {
			return fmt.Errorf("more than one secret with name %s", request.Name)
		}

		if request.Version != "" {
			if _, ok := secrets[0].SecretData[request.Version]; !ok {
				// this secret version doesn't exist
				return store.ErrNotExist
			}

			if len(secrets[0].SecretData) > 1 {
				// delete secret version, but only if there is more than one version.
				// If there are no more versions, just delete the whole secret.
				delete(secrets[0].SecretData, request.Version)
				// update the LatestVersion if that was the one that was removed
				if request.Version == secrets[0].LatestVersion {
					sorted := sorting.GetSortedSecretVersions(secrets[0])
					secrets[0].LatestVersion = sorted[0].ID
				}

				return store.UpdateSecret(tx, secrets[0])
			}
		}

		return store.DeleteSecret(tx, secrets[0].ID)
	}

	err := s.memstore.Update(deleteSecretFunc)
	switch {
	case err == store.ErrNotExist && request.Version == "":
		return nil, grpc.Errorf(codes.NotFound, "secret %s not found", request.Name)
	case err == store.ErrNotExist:
		return nil, grpc.Errorf(codes.NotFound, "secret %s version %s not found", request.Name, request.Version)
	case err != nil:
		return nil, err
	default:
		return &api.RemoveSecretResponse{}, nil
	}
}

// ListSecrets returns a list of all secrets.
func (s *Server) ListSecrets(ctx context.Context, request *api.ListSecretsRequest) (*api.ListSecretsResponse, error) {
	var (
		secrets []*api.Secret
		err     error
		by      store.By
	)

	by = store.All
	// return all secrets that match either any of the names or any of the name prefixes (why would you give both?)
	if request.Filters != nil {
		var filters []store.By
		for _, name := range request.Filters.Names {
			filters = append(filters, store.ByName(name))
		}
		for _, prefix := range request.Filters.NamePrefixes {
			filters = append(filters, store.ByNamePrefix(prefix))
		}
		switch len(filters) {
		case 0:
			break
		case 1:
			by = filters[0]
		default:
			by = store.Or(filters...)
		}
	}

	s.memstore.View(func(tx store.ReadTx) {
		secrets, err = store.FindSecrets(tx, by)
	})

	if err != nil {
		return nil, err
	}

	for _, secret := range secrets {
		cleanSecret(secret)
	}

	return &api.ListSecretsResponse{
		Secrets: secrets,
	}, nil
}

func validateSecretSpec(spec *api.SecretSpec) error {
	if spec == nil {
		return grpc.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	}
	if err := validateAnnotations(spec.Annotations); err != nil {
		return err
	}
	if _, ok := api.SecretType_name[int32(spec.Type)]; !ok {
		return grpc.Errorf(codes.InvalidArgument, errInvalidArgument.Error())
	}

	if len(spec.Data) > MaxSecretSize {
		return grpc.Errorf(codes.InvalidArgument, "secret data too large: max %d bytes", MaxSecretSize)
	}
	return nil
}

func validateAnnotations(m api.Annotations) error {
	if m.Name == "" {
		return grpc.Errorf(codes.InvalidArgument, "meta: name must be provided")
	} else if !isValidName.MatchString(m.Name) {
		// if the name doesn't match the regex
		return grpc.Errorf(codes.InvalidArgument, "invalid name, only [a-zA-Z0-9][a-zA-Z0-9-_.]*[a-zA-Z0-9] are allowed")
	}
	return nil
}
