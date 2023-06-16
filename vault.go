package vaultintegrated

import (
	"context"
	"time"

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
	AuthKey        = keyType(0)
	DatabaseKey    = keyType(1)
	AliCloudAKSKey = keyType(2)
)

type Vault struct {
	client                *vault.Client
	log                   *log.Helper
	AuthToken             *vault.Response[map[string]interface{}]
	toBeUpdateCredentials map[keyType][]*vault.Response[map[string]interface{}]
	toBeRotateCredentials []VaultSecrets
	ctx                   context.Context
	info                  *VaultInfo
	renewTicker           *time.Ticker
	rotateTicker          *time.Ticker
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
		AuthToken:             nil,
		toBeUpdateCredentials: make(map[keyType][]*vault.Response[map[string]interface{}], 0),
		toBeRotateCredentials: make([]VaultSecrets, 0),
		ctx:                   ctx,
		info:                  info,
		renewTicker:           time.NewTicker(time.Second * time.Duration(info.GetRenewLeadSec())),
		rotateTicker:          time.NewTicker(time.Second * time.Duration(info.GetRotationLeadSec()-5)),
	}
	// Register auth type
	RegisterAuthType("up", UPLogin)
	vault.AuthToken, err = AuthTypeMap[info.AuthType](ctx, vault.log, client, info)
	if err != nil {
		return nil, nil, err
	}
	vault.log.Infof("Vault auth success of: %s", info.AuthType)
	if err := vault.client.SetToken(vault.AuthToken.Auth.ClientToken); err != nil {
		return nil, nil, err
	}

	// Register vault secrets
	// vault.RegisterVaultSecrets(&DatabaseCredentials{})
	// vault.RegisterVaultSecrets(&KVSecret{})
	cleanup := func() {
		vault.ReovkeAllCredentials(ctx)
		client.ClearToken()
		vault.log.Info("Vault client token cleared")
		vault.renewTicker.Stop()
		vault.rotateTicker.Stop()
	}
	go func() {
		for {
			select {
			case <-vault.renewTicker.C:
				vault.RenewAuthToken(ctx)
				vault.RenewLeases(ctx)
			case <-vault.rotateTicker.C:
				vault.RotateAuthToken(ctx)
				vault.RotateALLCredentials(ctx)
			}
		}
	}()
	return vault, cleanup, nil
}
