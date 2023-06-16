package vaultintegrated

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/vault-client-go"
)

var AuthTypeMap map[string]AuthTypeFunc

func init() {
	AuthTypeMap = make(map[string]AuthTypeFunc, 0)
}

type AuthTypeFunc func(context.Context, *log.Helper, *vault.Client, *VaultInfo) (*vault.Response[map[string]interface{}], error)

func RegisterAuthType(name string, authType AuthTypeFunc) {
	AuthTypeMap[name] = authType
}
