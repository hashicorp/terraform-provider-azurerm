// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func GoogleClientID(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9-]+\.apps\.googleusercontent\.com$`), "must start with an identifier containing alphanumeric characters and hyphens and end with '.apps.googleusercontent.com'")(v, k)
}
