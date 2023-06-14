package vaultintegrated

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/vault-client-go"
)

const (
	UPAUTH int = iota
	AppRoleAUTH
	JWTAUTH
	TokenAUTH
	K8SAUTH
)

type keyType int

var (
	AuthKey = keyType(0)
	// databaseKey = keyType(1)
)

type Vault struct {
	client                *vault.Client
	log                   *log.Helper
	toBeUpdateCredentials map[keyType]*vault.Response[map[string]interface{}]
	ctx                   context.Context
	info                  *VaultInfo
}

func NewVault(info *VaultInfo, logger log.Logger) (*Vault, func(), error) {
	client, err := vault.New(vault.WithAddress(info.Dsn))
	if err != nil {
		return nil, nil, err
	}
	ctx := context.Background()

	vault := &Vault{
		client:                client,
		log:                   log.NewHelper(logger),
		toBeUpdateCredentials: make(map[keyType]*vault.Response[map[string]interface{}], 0),
		ctx:                   ctx,
		info:                  info,
	}

	resp, err := AuthTypeMap[info.AuthType](ctx, vault.log, client, info)
	if err != nil {
		return nil, nil, err
	}
	vault.log.Infof("Vault auth success of: %s", info.AuthType)
	if err := vault.client.SetToken(resp.Auth.ClientToken); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		vault.ReovkeAllCredentials(ctx)
		client.ClearToken()
		vault.log.Info("Vault client token cleared")
	}
	return vault, cleanup, nil
}
