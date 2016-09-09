package convert

import (
	types "github.com/docker/engine-api/types/swarm"
	swarmapi "github.com/docker/swarmkit/api"
	"github.com/docker/swarmkit/protobuf/ptypes"
)

// SecretFromGRPC converts a grpc Service to a Service.
func SecretFromGRPC(s *swarmapi.Secret) types.Secret {
	secret := types.Secret{
		ID:            s.ID,
		Name:          s.Name,
		LatestVersion: s.LatestVersion,
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
