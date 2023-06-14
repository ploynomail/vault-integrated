package vaultintegrated

import (
	"context"
	"time"

	"github.com/hashicorp/vault-client-go/schema"
)

func (v *Vault) RotateTokenOrCredentials(ctx context.Context) {
	v.log.Debugf("rotate token or credentials cycle: begin")
	defer v.log.Debugf("rotate token or credentials cycle: end")
	for i, x := range v.toBeUpdateCredentials {
		// rotate auth token
		if x.Auth != nil {
			resp, err := v.client.Auth.TokenLookUpSelf(ctx)
			if err != nil {
				v.log.Errorf("renew cycle: error looking up token: %v", err)
			}
			t, err := time.Parse(time.RFC3339Nano, resp.Data["expire_time"].(string))
			if err != nil {
				v.log.Errorf("renew cycle: error parsing expire_time: %v", err)
				continue
			}
			if !v.isExpiredExpireTime(t) {
				v.log.Debugf("rotate token or credentials cycle: auth token %s is not expired", x.Auth.ClientToken)
				continue
			}
			v.log.Debugf("rotate token or credentials cycle: auth token %s is expired", x.Auth.ClientToken)
			v.log.Debugf("rotate token or credentials cycle: rotating auth token %s", x.Auth.ClientToken)
			// rotate the auth token
			resp, err = v.UPLogin(ctx, v.c.Username, v.c.Password)
			if err != nil {
				v.log.Errorf("rotate cycle: error rotating auth token: %v", err)
			}
			if err := v.client.SetToken(resp.Auth.ClientToken); err != nil {
				v.log.Errorf("rotate cycle: error setting auth token: %v", err)
			}
			continue
		}
		// rotate credentials
		rs, err := v.client.System.LeasesReadLease(ctx, schema.LeasesReadLeaseRequest{
			LeaseId: x.LeaseID,
		})
		if err != nil {
			v.log.Errorf("rotate cycle: error reading lease: %v", err)
			v.toBeUpdateCredentials = append(v.toBeUpdateCredentials[:i], v.toBeUpdateCredentials[i+1:]...)
		}
		if !v.isExpiredExpireTime(rs.Data.ExpireTime) {
			v.log.Debugf("rotate token or credentials cycle: lease %s is not expired", x.LeaseID)
			continue
		}
		v.log.Debugf("rotate token or credentials cycle: lease %s is expired", x.LeaseID)
		v.log.Debugf("rotate token or credentials cycle: rotating lease %s", x.LeaseID)
		// rotate the lease(Read A New Database Credentials)
		DatabaseCredentials, resp, err := v.GetDatabaseCredentials(ctx)
		if err != nil {
			v.log.Errorf("rotate cycle: error rotating lease: %v", err)
		}
		v.toBeUpdateCredentials[i] = resp
		v.dataBaseCredentialsRoatateChan <- DatabaseCredentials
	}
}
