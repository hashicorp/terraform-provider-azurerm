// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// Domain labels can be 63 characters long per the Network API, the compute team adds a dash and a UUID when deploying to multiple
// Zones which causes a validation error in the RP. Updating the validation code to be artificially constrictive to account for the
// RPs behavior...
func OrchestratedDomainNameLabel(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-z][a-z0-9-]{1,24}[a-z0-9]$`), "must be between 1 - 26 characters long, start with a lower case letter, end with a lower case letter or number and contains only a-z, 0-9 and hyphens")(v, k)
}
