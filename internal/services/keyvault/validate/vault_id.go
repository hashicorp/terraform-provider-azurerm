package validate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// TODO: this file will be removed in time, but exists to allow for a gradual migration
// since this Resource ID is used all over the place

func VaultID(input interface{}, key string) (warnings []string, errors []error) {
	return commonids.ValidateKeyVaultID(input, key)
}
