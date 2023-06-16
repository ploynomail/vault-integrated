package vaultintegrated

import (
	"context"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

// DatabaseCredentials is a set of dynamic credentials retrieved from Vault
type DatabaseCredentials struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	v         *Vault
	MountPath string
	RoleName  string
	Ch        RotateCredentialsInotify[*DatabaseCredentials]
}

func NewDatabaseCredentials(v *Vault, mountPath, roleName string) *DatabaseCredentials {
	return &DatabaseCredentials{
		v:         v,
		MountPath: mountPath,
		RoleName:  roleName,
		Ch:        make(RotateCredentialsInotify[*DatabaseCredentials]),
	}
}

func (d *DatabaseCredentials) ReadCredentials(ctx context.Context, opts ...vault.RequestOption) error {
	resp, err := d.v.client.Secrets.DatabaseGenerateCredentials(ctx, d.RoleName, opts...)
	if err != nil {
		d.v.log.Errorf("Get Database Crednetials Err:%v", err)
	}
	d.v.toBeUpdateCredentials[DatabaseKey] = append(d.v.toBeUpdateCredentials[DatabaseKey], resp)
	d.Username = resp.Data["username"].(string)
	d.Password = resp.Data["password"].(string)
	return nil
}

func (d *DatabaseCredentials) RotateCredentials(ctx context.Context, v *Vault) {
	for i, resp := range v.toBeUpdateCredentials[DatabaseKey] {
		rs, err := v.client.System.LeasesReadLease(ctx, schema.LeasesReadLeaseRequest{
			LeaseId: resp.LeaseID,
		})
		if err != nil {
			v.log.Errorf("rotate cycle: error reading lease: %v", err)
			// remove the lease from the list
			v.toBeUpdateCredentials[DatabaseKey] = append(v.toBeUpdateCredentials[DatabaseKey][:i], v.toBeUpdateCredentials[DatabaseKey][i+1:]...)
			continue
		}
		// check if the lease is expired for a lead seconds
		if !v.isExpiredExpireTime(rs.Data.ExpireTime) {
			v.log.Debugf("rotate cycle: lease %s is not expired", rs.LeaseID)
			continue
		}
		v.log.Debugf("rotate cycle: lease %s is expired", resp.LeaseID)
		v.log.Debugf("rotate cycle: rotating lease %s", resp.LeaseID)
		// rotate the lease(Read A New Database Credentials)
		err = d.ReadCredentials(ctx, vault.WithMountPath(d.MountPath))
		if err != nil {
			v.log.Errorf("rotate cycle: Rotate Database Crednetials Err:%v", err)
		}
		d.Ch <- d
	}
}
