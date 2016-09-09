package sorting

import (
	"sort"
	"time"

	"github.com/docker/swarmkit/api"
)

// We want to secret data in decreasing creation time order
type secretDataSorter []*api.SecretData

func (k secretDataSorter) Len() int      { return len(k) }
func (k secretDataSorter) Swap(i, j int) { k[i], k[j] = k[j], k[i] }
func (k secretDataSorter) Less(i, j int) bool {
	iTime := time.Unix(k[i].Meta.CreatedAt.Seconds, int64(k[i].Meta.CreatedAt.Nanos))
	jTime := time.Unix(k[j].Meta.CreatedAt.Seconds, int64(k[j].Meta.CreatedAt.Nanos))
	return jTime.Before(iTime)
}

// GetSortedSecretVersions takes a secret, and returns all the SecretData objects in order of most recently created to least
func GetSortedSecretVersions(secret *api.Secret) []*api.SecretData {
	sds := make(secretDataSorter, len(secret.SecretData))
	i := 0
	for _, secret := range secret.SecretData {
		sds[i] = secret
		i++
	}
	sort.Sort(sds)
	return sds
}
