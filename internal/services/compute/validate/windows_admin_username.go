// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// WindowsAdminUsername validates that admin_username meets the Azure API requirements for Windows Virtual Machines.
func WindowsAdminUsername(i interface{}, k string) (warnings []string, errors []error) {
	// adminUsername must not be empty, can be at most 20 characters, cannot end with a dot and cannot match a disallowed name.
	return validation.All(
		validation.StringIsNotWhiteSpace,
		validation.StringLenBetween(1, 20),
		validation.StringDoesNotMatch(regexp.MustCompile(`\.$`), "cannot end with a dot"),
		validation.StringNotInSlice([]string{"administrator", "admin", "user", "user1", "test", "user2", "test1", "user3", "admin1", "1", "123", "a", "actuser", "adm", "admin2", "aspnet", "backup", "console", "david", "guest", "john", "owner", "root", "server", "sql", "support", "support_388945a0", "sys", "test2", "test3", "user4", "user5"}, false),
	)(i, k)
}
