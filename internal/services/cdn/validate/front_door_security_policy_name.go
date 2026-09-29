// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FrontDoorSecurityPolicyName(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[\da-zA-Z](?:[-\da-zA-Z]*[\da-zA-Z])?$`), "must begin and end with an alphanumeric character, and may contain only alphanumeric characters and hyphens")(v, k)
}
