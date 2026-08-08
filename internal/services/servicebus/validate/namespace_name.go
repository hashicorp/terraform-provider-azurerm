// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func NamespaceName(i interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringMatch(regexp.MustCompile("^[a-zA-Z][-a-zA-Z0-9]{4,48}[a-zA-Z0-9]$"), "can contain only letters, numbers, and hyphens. The namespace must start with a letter, and it must end with a letter or number and be between 6 and 50 characters long"),
		// The name cannot end with "-", "-sb" or "-mgmt".
		// See more details from link https://docs.microsoft.com/en-us/rest/api/servicebus/create-namespace.
		validation.StringDoesNotMatch(regexp.MustCompile(`(-|-sb|-mgmt)$`), "cannot end with a hyphen, -sb, or -mgmt"),
	)(i, k)
}
