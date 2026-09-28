// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FlexibleServerAdministratorLogin(i interface{}, k string) (warnings []string, errors []error) {
	return validation.All(
		validation.StringLenBetween(1, 32),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_]*$`), "must only contain characters, numbers or '_'"),
		validation.StringNotInSlice([]string{"azure_superuser", "admin", "administrator", "root", "guest", "public"}, false),
	)(i, k)
}
