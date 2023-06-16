package vaultintegrated

import (
	"context"
	"time"

	"github.com/hashicorp/vault-client-go"
)

type VaultSecrets interface {
	ReadCredentials(ctx context.Context, opts ...vault.RequestOption) error
	RotateCredentials(context.Context, *Vault)
}

type RotateCredentialsInotify[T VaultSecrets] chan T

func (v *Vault) RegisterVaultSecrets(secrets VaultSecrets) {
	v.toBeRotateCredentials = append(v.toBeRotateCredentials, secrets)
}

func (v *Vault) RotateALLCredentials(ctx context.Context) {
	for _, secret := range v.toBeRotateCredentials {
		secret.RotateCredentials(ctx, v)
	}
}

// RotateAuthToken is used to rotate the auth token
func (v *Vault) RotateAuthToken(ctx context.Context) {
	resp, err := v.client.Auth.TokenLookUpSelf(ctx)
	if err != nil {
		v.log.Errorf("rotatae cycle: error looking up token: %v", err)
	}
	t, err := time.Parse(time.RFC3339Nano, resp.Data["expire_time"].(string))
	if err != nil {
		v.log.Errorf("rotatae cycle: error parsing expire_time: %v; Force login again;", err)
		v.AuthToken, err = AuthTypeMap[v.info.AuthType](ctx, v.log, v.client, v.info)
		if err != nil {
			v.log.Errorf("rotatae cycle: error login again: %v", err)
		}
		v.log.Infof("rotatae cycle: vault auth success of: %s", v.info.AuthType)
		if err := v.client.SetToken(v.AuthToken.Auth.ClientToken); err != nil {
			v.log.Errorf("rotatae cycle: error setting auth token: %v", err)
		}
	}
	if !v.isExpiredExpireTime(t) {
		v.log.Debugf("rotatae cycle: auth token is not expired")
		return
	}
	v.log.Debugf("rotatae cycle: auth token is expired")
	v.log.Debugf("rotatae cycle: rotating auth token")

	v.AuthToken, err = AuthTypeMap[v.info.AuthType](ctx, v.log, v.client, v.info)
	if err != nil {
		v.log.Errorf("rotatae cycle: error login again: %v", err)
	}
	v.log.Infof("rotatae cycle: Vault auth success of: %s", v.info.AuthType)
	if err := v.client.SetToken(v.AuthToken.Auth.ClientToken); err != nil {
		v.log.Errorf("rotatae cycle: error setting auth token: %v", err)
	}
}
