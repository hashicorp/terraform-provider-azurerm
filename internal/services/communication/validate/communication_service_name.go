// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CommunicationServiceName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^(([a-zA-Z])|([a-zA-Z][0-9a-zA-Z-]{0,62}[0-9a-zA-Z]))$`), "must be between 1 and 64 characters in length and start with letters and contain only letters, numbers and hyphens. And it cannot end with hyphen")(i, k)
}
