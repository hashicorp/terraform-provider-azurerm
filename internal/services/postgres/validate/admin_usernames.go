// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AdminUsernames(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringNotInSlice([]string{"azure_superuser", "azure_pg_admin", "admin", "administrator", "root", "guest", "public"}, true),
		validation.StringDoesNotMatch(regexp.MustCompile(`(?i)^pg_`), "cannot start with 'pg_'"),
	)(i, k)
}
