package vaultintegrated

import (
	"context"
	"testing"
)

var vault_test_var_vautl_read_alicloudAK *vault_test_var = &vault_test_var{
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

func TestVaultReadAlicloudAK(t *testing.T) {
	ctx := context.Background()
	v, cleanup, err = NewVault(vault_test_var_vautl_read_alicloudAK.info, log_test, true)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	alicloudAK := NewAliCloudAKS(v, "alicloud", "policy-based")
	v.RegisterVaultSecrets(alicloudAK)
	err := alicloudAK.ReadCredentials(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("alicloudAK: %s:%s", alicloudAK.AK, alicloudAK.SK)
	for sf := range alicloudAK.Ch {
		t.Log(sf.AK, sf.SK)
	}

}
