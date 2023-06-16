package vaultintegrated

import (
	"context"

	"github.com/hashicorp/vault-client-go/schema"
)

// ReovkeAllCredentials revokes all credentials
func (v *Vault) ReovkeAllCredentials(ctx context.Context) {
	v.log.Debugf("Reovke Credentials cycle: begin")
	defer v.log.Debugf("Reovke Credentials cycle: end")
	for _, x := range v.toBeUpdateCredentials {
		for _, f := range x {
			_, err := v.client.System.LeasesRevokeLease(ctx, schema.LeasesRevokeLeaseRequest{
				LeaseId: f.LeaseID,
			})
			if err != nil {
				v.log.Errorf("renew cycle: error reading lease: %v", err)
			}
		}
	}
}
