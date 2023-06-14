package vaultintegrated

import (
	"context"

	"github.com/hashicorp/vault-client-go/schema"
)

// ReovkeAllCredentials revokes all credentials
func (v *Vault) ReovkeAllCredentials(ctx context.Context) {
	v.log.Debugf("Reovke Credentials cycle: begin")
	defer v.log.Debugf("Reovke Credentials cycle: end")
	for keyType, x := range v.toBeUpdateCredentials {
		// if auth is nil, it means that the credentials have been revoked
		// auht token will be revoke of client.ClearToken(),so we don't need to revoke it again
		if keyType == AuthKey && x.Auth == nil {
			continue
		}
		_, err := v.client.System.LeasesRevokeLease(ctx, schema.LeasesRevokeLeaseRequest{
			LeaseId: x.LeaseID,
		})
		if err != nil {
			v.log.Errorf("renew cycle: error reading lease: %v", err)
		}
	}
}
