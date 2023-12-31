package vaultintegrated

import (
	"context"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

type KVSecret struct {
	Key       string
	Value     interface{}
	v         *Vault
	MountPath string
	ctx       context.Context
}

func NewKVSecret(ctx context.Context, v *Vault, key string, mountPath string) *KVSecret {
	return &KVSecret{
		ctx:       ctx,
		v:         v,
		Key:       key,
		MountPath: mountPath,
	}
}

// ReadCredentials read credentials from vault
func (k *KVSecret) ReadCredentials(ctx context.Context, opts ...vault.RequestOption) error {
	opts = append(opts, vault.WithMountPath(k.MountPath))
	resp, err := k.v.client.Secrets.KvV2Read(ctx, k.Key, opts...)
	if err != nil {
		return err
	}
	k.Value = resp.Data.Data
	return nil
}

// RotateCredential KV Secret don't need to rotate, becuase it's not dynamic, so just return
func (k *KVSecret) RotateCredentials(context.Context, *Vault) {
}

func (k *KVSecret) WriteCredentials(ctx context.Context, keyPair map[string]interface{}, opts ...vault.RequestOption) error {
	opts = append(opts, vault.WithMountPath(k.MountPath))
	_, err := k.v.client.Secrets.KvV2Write(ctx, k.Key, schema.KvV2WriteRequest{
		Data: keyPair,
	}, opts...)
	return err
}
