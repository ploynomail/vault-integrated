package vaultintegrated

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/vault-client-go/schema"
)

func (v *Vault) RenewLeases(ctx context.Context) {
	v.log.Debugf("renew cycle: begin")
	defer v.log.Debugf("renew cycle: end")
	for i, x := range v.toBeUpdateCredentials {
		// If the lease is an auth token, renew it
		// Auth, if non-nil, means that there was authentication information attached to this response
		if i == AuthKey && x.Auth != nil {
			resp, err := v.client.Auth.TokenLookUpSelf(ctx)
			if err != nil {
				v.log.Errorf("renew cycle: error looking up token: %v", err)
			}
			ttl, err := resp.Data["ttl"].(json.Number).Int64()
			if err != nil {
				v.log.Errorf("renew cycle: error parsing ttl: %v", err)
				continue
			}
			if !v.isExpiredTTL(ttl) {
				v.log.Debugf("renew cycle: auth token %s is not expired", x.Auth.ClientToken)
				continue
			}
			v.log.Debugf("renew cycle: auth token %s is expired", x.Auth.ClientToken)
			v.log.Debugf("renew cycle: renewing auth token %s", x.Auth.ClientToken)
			// renew the auth token 增加80秒ttl，但是不能超过最大ttl(增加总和也不能超过最大ttl)
			_, err = v.client.Auth.TokenRenewSelf(ctx, schema.TokenRenewSelfRequest{
				Increment: int32(v.info.GetRenewLeadSec()),
				Token:     x.Auth.ClientToken,
			})
			if err != nil {
				v.log.Errorf("renew cycle: error renewing auth token: %v", err)
			}
			v.log.Debugf("renew cycle: skipping auth token %s", x.Auth.ClientToken)
			continue
		} else {
			// If the lease is not an auth token, renew it
			rs, err := v.client.System.LeasesReadLease(ctx, schema.LeasesReadLeaseRequest{
				LeaseId: x.LeaseID,
			})
			if err != nil {
				v.log.Errorf("renew cycle: error reading lease: %v", err)
			}
			ttl := rs.Data.Ttl
			if ttl == 0 {
				v.log.Debugf("renew cycle: lease %s is not renewable", x.LeaseID)
				continue
			}
			if !v.isExpiredTTL(int64(ttl)) {
				v.log.Debugf("renew cycle: lease %s is not expired, expired time is: %v", x.LeaseID, rs.Data.ExpireTime)
				continue
			}
			v.log.Debugf("renew cycle: lease %s is expired", x.LeaseID)
			v.log.Debugf("renew cycle: renewing lease %s", x.LeaseID)
			// renew the lease
			_, err = v.client.System.LeasesRenewLeaseWithId(ctx, x.LeaseID, schema.LeasesRenewLeaseWithIdRequest{
				LeaseId:   x.LeaseID,
				Increment: 300,
			})
			if err != nil {
				v.log.Errorf("renew cycle: error renewing lease: %v", err)
			}
		}
	}
}
