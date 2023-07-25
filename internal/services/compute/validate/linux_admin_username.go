// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"
)

// LinuxAdminUsername validates that admin_username meets the Azure API requirements for Linux Virtual Machines.
func LinuxAdminUsername(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string but it wasn't", k))
		return
	}

	// adminUsername must not be empty.
	if strings.TrimSpace(v) == "" {
		errors = append(errors, fmt.Errorf("%q must not be empty", k))
		return
	}

	// adminUsername Max-length 64 characters.
	if len(v) > 64 {
		errors = append(errors, fmt.Errorf("%q most be between 1 and %d characters, got %d", k, 64, len(v)))
	}

	// adminUsername cannot match the following disallowed names.
	disallowedNames := []string{"administrator", "admin", "user", "user1", "test", "user2", "test1", "user3", "admin1", "1", "123", "a", "actuser", "adm", "admin2", "aspnet", "backup", "console", "david", "guest", "john", "owner", "root", "server", "sql", "support", "support_388945a0", "sys", "test2", "test3", "user4", "user5"}
	for _, value := range disallowedNames {
		if value == v {
			errors = append(errors, fmt.Errorf("%q specified is not allowed, got %q, cannot match: %q", k, v, strings.Join(disallowedNames, ", ")))
		}
	}

	return
}
