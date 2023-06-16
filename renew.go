package vaultintegrated

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/vault-client-go/schema"
)

// RenewAuthToken renews the auth token
func (v *Vault) RenewAuthToken(ctx context.Context) {
	resp, err := v.client.Auth.TokenLookUpSelf(ctx)
	if err != nil {
		v.log.Errorf("renew cycle: error looking up token: %v", err)
	}
	ttl, err := resp.Data["ttl"].(json.Number).Int64()
	if err != nil {
		v.log.Errorf("renew cycle: error parsing ttl: %v", err)

	}
	if !v.isExpiredTTL(ttl) {
		v.log.Debugf("renew cycle: auth token is not expired")
		return
	}
	v.log.Debugf("renew cycle: auth token is expired")
	// renew the auth token
	_, err = v.client.Auth.TokenRenewSelf(ctx, schema.TokenRenewSelfRequest{
		Increment: int32(v.info.GetRenewLeadSec()),
		Token:     v.AuthToken.Auth.ClientToken,
	})
	if err != nil {
		v.log.Errorf("renew cycle: error renewing auth token: %v", err)
	}
	v.log.Debugf("renew cycle: skipping auth token renewal")
}

// RenewLeases renews the leases
func (v *Vault) RenewLeases(ctx context.Context) {
	v.log.Debugf("renew cycle: begin, %s", time.Now().UTC())

	defer v.log.Debugf("renew cycle: end")
	for _, x := range v.toBeUpdateCredentials {
		for _, y := range x {
			if !y.Renewable {
				v.log.Debugf("renew cycle: lease %s is not renewable", y.LeaseID)
				continue
			}
			// If the lease is not an auth token, renew it
			rs, err := v.client.System.LeasesReadLease(ctx, schema.LeasesReadLeaseRequest{
				LeaseId: y.LeaseID,
			})
			if err != nil {
				v.log.Errorf("renew cycle: error reading lease: %v", err)
				continue
			}
			fmt.Println(rs.Data.Ttl)
			ttl := rs.Data.Ttl
			if ttl == 0 {
				v.log.Debugf("renew cycle: lease %s is not renewable", y.LeaseID)
				continue
			}
			if !v.isExpiredTTL(int64(ttl)) {
				v.log.Debugf("renew cycle: lease %s is not expired, expired time is: %v", y.LeaseID, rs.Data.ExpireTime)
				continue
			}
			v.log.Debugf("renew cycle: lease %s is expired", y.LeaseID)
			v.log.Debugf("renew cycle: renewing lease %s", y.LeaseID)
			// renew the lease
			_, err = v.client.System.LeasesRenewLeaseWithId(ctx, y.LeaseID, schema.LeasesRenewLeaseWithIdRequest{
				LeaseId:   y.LeaseID,
				Increment: int32(v.info.AuthTTLIncrement),
			})
			if err != nil {
				v.log.Errorf("renew cycle: error renewing lease: %v", err)
			}
		}
	}
}
