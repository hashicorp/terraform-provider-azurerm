// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func BackupPolicyMySQLFlexibleServerName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile("^[-a-zA-Z0-9]{3,150}$"), "must be 3 - 150 characters long, contain only letters, numbers and hyphens")(i, k)
}
