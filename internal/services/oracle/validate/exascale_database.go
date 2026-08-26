// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

const exascaleDatabaseVirtualMachineClusterSSHPublicKeysMaxCombinedLength = 10000

func ExascaleDatabaseResourceName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z_](?:[a-zA-Z0-9_]*(?:-[a-zA-Z0-9_]+)*-?)?$`), "must begin with a letter or underscore (_), contain only letters, numbers, underscores (_) and cannot contain any consecutive hyphens (--)")(i, k)
}

func ExascaleDatabaseVirtualMachineClusterSSHPublicKeys(i interface{}, k string) (warnings []string, errors []error) {
	var keys []string
	switch v := i.(type) {
	case []interface{}:
		keys = make([]string, 0, len(v))
		for _, keyRaw := range v {
			key, ok := keyRaw.(string)
			if !ok {
				return nil, append(errors, fmt.Errorf("expected %s to contain only strings", k))
			}
			keys = append(keys, key)
		}
	case []string:
		keys = v
	default:
		return nil, append(errors, fmt.Errorf("expected type of %s to be list", k))
	}

	totalLength := 0
	for _, key := range keys {
		totalLength += len(key)
		if totalLength > exascaleDatabaseVirtualMachineClusterSSHPublicKeysMaxCombinedLength {
			errors = append(errors, fmt.Errorf("the combined length of all provided public SSH keys cannot exceed 10000 characters"))
			return
		}
	}

	return
}
