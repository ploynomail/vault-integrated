package vaultintegrated

import (
	"time"
)

// isExpiredAuthToken is a helper function that returns true if the given expiration
func (v *Vault) isExpiredTTL(ttl int64) bool {
	return ttl < v.info.GetRenewLeadSec()
}

// isExpired is a helper function that returns true if the given expiration
func (v *Vault) isExpiredExpireTime(expiration time.Time) bool {
	now := time.Now().UTC()
	expiration = expiration.Add(time.Duration(-v.info.GetRotationLeadSec()) * time.Second)
	return now.After(expiration)
}
