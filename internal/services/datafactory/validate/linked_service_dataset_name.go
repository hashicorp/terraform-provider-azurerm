// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func LinkedServiceDatasetName(v interface{}, k string) ([]string, []error) {
	return validation.StringDoesNotMatch(regexp.MustCompile(`^[-.+?/<>*%&:\\]+$`), "any of '-' '.', '+', '?', '/', '<', '>', '*', '%', '&', ':', '\\', are not allowed")(v, k)
}
