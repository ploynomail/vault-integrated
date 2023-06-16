package vaultintegrated

import (
	"context"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type AliCloudAKS struct {
	AK        string
	SK        string
	v         *Vault
	MountPath string
	RoleName  string
	Ch        RotateCredentialsInotify[*AliCloudAKS]
}

func NewAliCloudAKS(v *Vault, mountPath, roleName string) *AliCloudAKS {
	return &AliCloudAKS{
		v:         v,
		MountPath: mountPath,
		RoleName:  roleName,
		Ch:        make(RotateCredentialsInotify[*AliCloudAKS]),
	}
}

func (d *AliCloudAKS) ReadCredentials(ctx context.Context, opts ...vault.RequestOption) error {
	resp, err := d.v.client.Secrets.AliCloudGenerateCredentials(ctx, d.RoleName, opts...)
	if err != nil {
		d.v.log.Errorf("Get AliCloud AKS Crednetials Err:%v", err)
	}
	d.v.toBeUpdateCredentials[AliCloudAKSKey] = append(d.v.toBeUpdateCredentials[AliCloudAKSKey], resp)
	d.AK = resp.Data["access_key"].(string)
	d.SK = resp.Data["secret_key"].(string)
	return nil
}

func (d *AliCloudAKS) RotateCredentials(ctx context.Context, v *Vault) {
	for i, resp := range v.toBeUpdateCredentials[AliCloudAKSKey] {
		rs, err := v.client.System.LeasesReadLease(ctx, schema.LeasesReadLeaseRequest{
			LeaseId: resp.LeaseID,
		})
		if err != nil {
			v.log.Errorf("rotate cycle: error reading lease: %v", err)
			// remove the lease from the list
			v.toBeUpdateCredentials[AliCloudAKSKey] = append(v.toBeUpdateCredentials[AliCloudAKSKey][:i], v.toBeUpdateCredentials[AliCloudAKSKey][i+1:]...)
			continue
		}
		// check if the lease is expired for a lead seconds
		if !v.isExpiredExpireTime(rs.Data.ExpireTime) {
			v.log.Debugf("rotate cycle: lease %s is not expired", rs.LeaseID)
			continue
		}
		v.log.Debugf("rotate cycle: lease %s is expired", resp.LeaseID)
		v.log.Debugf("rotate cycle: rotating lease %s", resp.LeaseID)
		// rotate the lease(Read A New AliCloud AKS Credentials)
		err = d.ReadCredentials(ctx, vault.WithMountPath(d.MountPath))
		if err != nil {
			v.log.Errorf("rotate cycle: error reading lease: %v", err)
			// remove the lease from the list
			v.toBeUpdateCredentials[AliCloudAKSKey] = append(v.toBeUpdateCredentials[AliCloudAKSKey][:i], v.toBeUpdateCredentials[AliCloudAKSKey][i+1:]...)
			continue
		}
		// notify the new credentials
		d.Ch <- d
		// remove the lease from the list
		v.toBeUpdateCredentials[AliCloudAKSKey] = append(v.toBeUpdateCredentials[AliCloudAKSKey][:i], v.toBeUpdateCredentials[AliCloudAKSKey][i+1:]...)
	}
}
