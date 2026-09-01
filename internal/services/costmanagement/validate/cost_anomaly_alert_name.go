// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CostAnomalyAlertName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([a-z\d-]*)$`), "must contain only lowercase letters, numbers and hyphens")(i, k)
}
