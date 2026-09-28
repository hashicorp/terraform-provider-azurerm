// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func EmbeddedAdministratorName(v interface{}, k string) (warnings []string, errors []error) {
	// a UUID is valid in addition to an email address
	return validation.Any(
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`), "isn't a valid email address"),
		validation.IsUUID,
	)(v, k)
}
