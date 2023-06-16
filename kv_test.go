package vaultintegrated

import (
	"context"
	"testing"
)

var vault_test_var_vautl_read_KVSecrets *vault_test_var = &vault_test_var{
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

func TestVaultReadKVSecrets(t *testing.T) {
	ctx := context.Background()
	v, cleanup, err = NewVault(vault_test_var_vautl_read_KVSecrets.info, log_test)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	kvReadTest := NewKVSecret(ctx, v, "app", "kv")
	err := kvReadTest.ReadCredentials(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(kvReadTest.Value)
}
