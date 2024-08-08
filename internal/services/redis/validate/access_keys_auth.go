package validate

import (
	"fmt"
	"log"
)

// Entra (AD) auth has to be set to disable access keys auth
// https://learn.microsoft.com/en-us/azure/azure-cache-for-redis/cache-azure-active-directory-for-authentication
func ValidateAccessKeysAuth(accessKeysAuthenticationDisabled bool, activeDirectoryAuthenticationEnabled bool) error {
	log.Printf("[DEBUG] ValidateAccessKeysAuth: accessKeysAuthenticationDisabled: %v, activeDirectoryAuthenticationEnabled: %v", accessKeysAuthenticationDisabled, activeDirectoryAuthenticationEnabled)

	if accessKeysAuthenticationDisabled && !activeDirectoryAuthenticationEnabled {
		return fmt.Errorf("microsoft entra authorization (active_directory_authentication_enabled) must be enabled in order to disable access key authentication (access_keys_authentication_disabled): https://learn.microsoft.com/en-us/azure/azure-cache-for-redis/cache-azure-active-directory-for-authentication#disable-access-key-authentication-on-your-cache")
	}

	return nil
}
