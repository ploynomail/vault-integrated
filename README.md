# vault golang 集成封装

本库在vault-client-go[https://github.com/hashicorp/vault-client-go]基础上做了一些封装，方便应用集成。
注意： 该集成库只支持vault 1.13.3版本进行测试开始，且不是一个充分完备的集成库，只是一个简单的封装，如果有需要，可以自行修改。prod使用前请进行充分测试。

## 使用方法

### 1. 安装

```shell
go install github.com/hashicorp/vault-client-go
```

### 2. 基本使用方法
```go
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

    ctx := context.Background()
	v, cleanup, err = NewVault(vault_test_var_vautl_read_alicloudAK.info, log_test)
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
```

### 3. 内置认证方式
userpass: 用户名密码认证 AuthType: up
K8S: k8s集群认证 AuthType: k8s

### 4. 自定义认证方式
实现AuthTypeFunc函数类型，注册到Vault中即可

```go
type AuthTypeFunc func(context.Context, *log.Helper, *vault.Client, *VaultInfo) (*vault.Response[map[string]interface{}], error)
RegisterAuthType("auth name", AuthTypeFunc)
```

例如：
```go
var (
	ErrNeedUserOrPassword = errors.New("need username or password")
)

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

func init() {
    vaultintegrated.RegisterAuthType("up", UPLogin)
}
```
### 5. 内置secret
alicloudAK: 阿里云AK认证，使用方法参照：alicloudAK_test.go
kv: kv存储，使用方法参照：kv_test.go
database: 数据库认证，使用方法参照：database_test.go

### 6. 自定义secret
实现下列接口，实例化后注册到Vault中即可
测试PR
```go
type VaultSecrets interface {
	ReadCredentials(ctx context.Context, opts ...vault.RequestOption) error
	RotateCredentials(context.Context, *Vault)
}
```
