package vaultintegrated

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

var (
	ErrNeedUserOrPassword = errors.New("need username or password")
)

func init() {
	RegisterAuthType("up", UPLogin)
}

func UPLogin(ctx context.Context, log *log.Helper, c *vault.Client, info *VaultInfo) (*vault.Response[map[string]interface{}], error) {
	if info.Username == "" || info.Password == "" {
		return nil, ErrNeedUserOrPassword
	}

	resp, err := c.Auth.UserpassLogin(ctx, info.Username, schema.UserpassLoginRequest{
		Password: info.Password,
	})
	if err != nil {
		return nil, err
	}
	log.Debugf("Vault init UPLogin Success!!!")
	return resp, nil
}
