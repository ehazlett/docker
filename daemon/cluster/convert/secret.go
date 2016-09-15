package convert

import (
	"fmt"

	types "github.com/docker/docker/api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/protobuf/ptypes"
)

// SecretFromGRPC converts a grpc Service to a Service.
func SecretFromGRPC(s *swarmapi.Secret) types.Secret {
	secret := types.Secret{
		ID:            s.ID,
		Name:          s.Name,
		LatestVersion: s.LatestVersion,
		SecretData:    map[string]*types.SecretData{},
	}

	// Meta
	secret.Version.Index = s.Meta.Version.Index
	secret.CreatedAt, _ = ptypes.Timestamp(s.Meta.CreatedAt)
	secret.UpdatedAt, _ = ptypes.Timestamp(s.Meta.UpdatedAt)

	// Data
	for k, v := range s.SecretData {
		spec := types.SecretSpec{
			Data: v.Spec.Data,
		}
		switch v.Spec.Type {
		case swarmapi.SecretType_ContainerSecret:
			spec.Type = types.ContainerSecret
		case swarmapi.SecretType_NodeSecret:
			spec.Type = types.NodeSecret
		}
		// Annotations
		spec.Name = v.Spec.Annotations.Name
		spec.Labels = v.Spec.Annotations.Labels

		secret.SecretData[k] = &types.SecretData{
			ID:         v.ID,
			Spec:       spec,
			Digest:     v.Digest,
			SecretSize: v.SecretSize,
		}
	}

	return secret
}

// SecretFromGRPC converts a grpc Service to a Service.
func SecretSpecToGRPC(s types.SecretSpec) (swarmapi.SecretSpec, error) {
	spec := swarmapi.SecretSpec{
		Annotations: swarmapi.Annotations{
			Name:   s.Name,
			Labels: s.Labels,
		},
		Data: s.Data,
	}

	switch s.Type {
	case types.ContainerSecret:
		spec.Type = swarmapi.SecretType_ContainerSecret
	case types.NodeSecret:
		spec.Type = swarmapi.SecretType_NodeSecret
	default:
		return swarmapi.SecretSpec{}, fmt.Errorf("unknown secret type specified: %s", s.Type)
	}

	return spec, nil
}
