package vaultintegrated

type VaultInfo struct {
	Dsn                  string // vault dsn
	Username             string // vault username
	Password             string // vault password
	JwtToken             string // jwt token
	RoleName             string // vault role name
	RenewLeadSec         int64  // renew token lead time default 60
	RotationLeadSec      int64  // rotation lead time default 60
	AuthType             string // auth type
	AuthTTLIncrement     int64  // auth ttl increment default 3600
	CredentialsIncrement int64  // credentials ttl increment default 3600
}

// GetRenewLeadSec returns the renew lead time default 60
func (x *VaultInfo) GetRenewLeadSec() int64 {
	if x.RenewLeadSec == 0 {
		return 60
	}
	return x.RenewLeadSec
}

// GetRotationLeadSec returns the rotation lead time default 60
func (x *VaultInfo) GetRotationLeadSec() int64 {
	if x.RotationLeadSec == 0 {
		return 60
	}
	return x.RotationLeadSec
}

// GetAuthTTLIncrement returns the auth ttl increment default 3600
func (x *VaultInfo) GetAuthTTLIncrement() int64 {
	if x.AuthTTLIncrement == 0 {
		return 3600
	}
	return x.AuthTTLIncrement
}

// GetCredentialsIncrement returns the credentials ttl increment default 3600
func (x *VaultInfo) GetCredentialsIncrement() int64 {
	if x.CredentialsIncrement == 0 {
		return 3600
	}
	return x.CredentialsIncrement
}
