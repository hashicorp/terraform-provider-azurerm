// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ExascaleDatabaseStorageVaultName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z_](?:[a-zA-Z0-9_]*(?:-[a-zA-Z0-9_]+)*-?)?$`), "must begin with a letter or underscore (_), contain only letters, numbers, underscores (_) and cannot contain any consecutive hyphens (--)")(i, k)
}
