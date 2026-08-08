// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"
)

const exascaleDatabaseVirtualMachineClusterSSHPublicKeysMaxCombinedLength = 10000

func ExascaleDatabaseResourceName(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, append(errors, fmt.Errorf("expected type of %s to be string", k))
	}

	// The value must start with a letter or an underscore (_).
	// The value can only include letters, numbers, underscores (_), and hyphens (-).
	// The value cannot contain consecutive hyphens (--).
	if !regexp.MustCompile(`^[a-zA-Z_]([a-zA-Z0-9_]*(-[a-zA-Z0-9_]+)*-?)?$`).MatchString(v) {
		errors = append(errors, fmt.Errorf("%s must begin with a letter or underscore (_), contain only letters, numbers, underscores (_), and hyphens (-), and cannot contain consecutive hyphens (--)", k))
	}

	return
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
