package vaultintegrated

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

var (
	ErrNeedJWT = errors.New("need JWT Token")
)

func K8SLogin(ctx context.Context, log *log.Helper, c *vault.Client, info *VaultInfo) (*vault.Response[map[string]interface{}], error) {
	if info.JwtToken == "" {
		return nil, ErrNeedJWT
	}
	resp, err := c.Auth.KubernetesLogin(ctx, schema.KubernetesLoginRequest{
		Jwt:  info.JwtToken,
		Role: info.RoleName,
	})
	if err != nil {
		return nil, err
	}
	log.Debugf("Vault init K8SLogin Success!!!")
	return resp, nil
}
