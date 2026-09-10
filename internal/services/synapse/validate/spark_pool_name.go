// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SparkPoolName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z\d]{0,14}$`), "can contain only letters or numbers, must start with a letter, and be between 1 and 15 characters long")(i, k)
}
