// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func BackupPolicyName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][-\da-zA-Z]{2,149}$`), "must be between 3 and 150 characters in length and start with letters and contains only letters, numbers or hyphens")(v, k)
}
