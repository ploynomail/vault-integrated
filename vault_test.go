package vaultintegrated

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

type vault_test_var struct {
	info *VaultInfo
}

var vault_test_var_vautl_factory_func *vault_test_var = &vault_test_var{
	info: &VaultInfo{
		Dsn:             "http://10.0.2.4:8200",
		Username:        "uu",
		Password:        "ff",
		RoleName:        "my-role",
		AuthType:        "up",
		RenewLeadSec:    10,
		RotationLeadSec: 60,
	},
}

var log_test log.Logger = log.DefaultLogger
var v *Vault
var cleanup func()
var err error

func TestVaultFactoryFunc(t *testing.T) {
	v, cleanup, err = NewVault(vault_test_var_vautl_factory_func.info, log_test, true)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
}
