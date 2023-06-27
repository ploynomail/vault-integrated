package vaultintegrated

import (
	"context"
	"testing"
	"time"
)

var vault_test_var_vautl_read_DatabaseCredentials *vault_test_var = &vault_test_var{
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

func TestVaultReadDatabaseCredentials(t *testing.T) {
	ctx := context.Background()
	v, cleanup, err = NewVault(vault_test_var_vautl_read_DatabaseCredentials.info, log_test, true)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	databaseCredentials := NewDatabaseCredentials(v, "database", "my-role")
	v.RegisterVaultSecrets(databaseCredentials)
	err := databaseCredentials.ReadCredentials(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(databaseCredentials.Username, databaseCredentials.Password)
	for sf := range databaseCredentials.Ch {
		t.Log(sf.Username, sf.Password)

	}
	time.Sleep(time.Second * 6000)
}
