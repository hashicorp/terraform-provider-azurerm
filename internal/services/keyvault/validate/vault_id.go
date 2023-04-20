package validate

import (
	"fmt"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// TODO: this file will be removed in time, but exists to allow for a gradual migration
// since this Resource ID is used all over the place

func VaultID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := commonids.ParseKeyVaultID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
