package vaultintegrated

import (
	"context"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type DatabaseStaticCredentials struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	v         *Vault
	MountPath string
	RoleName  string
	Ch        RotateCredentialsInotify[*DatabaseStaticCredentials]
}

func NewDatabaseStaticCredentials(v *Vault, mountPath, roleName string) *DatabaseStaticCredentials {
	return &DatabaseStaticCredentials{
		v:         v,
		MountPath: mountPath,
		RoleName:  roleName,
		Ch:        make(RotateCredentialsInotify[*DatabaseStaticCredentials]),
	}
}

func (d *DatabaseStaticCredentials) ReadCredentials(ctx context.Context, opts ...vault.RequestOption) error {
	resp, err := d.v.client.Secrets.DatabaseReadStaticRoleCredentials(ctx, d.RoleName, opts...)
	if err != nil {
		d.v.log.Errorf("Get Database Crednetials Err:%v", err)
	}
	d.v.toBeUpdateCredentials[DatabaseKey] = append(d.v.toBeUpdateCredentials[DatabaseKey], resp)
	d.Username = resp.Data["username"].(string)
	d.Password = resp.Data["password"].(string)
	return nil
}

func (d *DatabaseStaticCredentials) RotateCredentials(ctx context.Context, v *Vault) {
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
