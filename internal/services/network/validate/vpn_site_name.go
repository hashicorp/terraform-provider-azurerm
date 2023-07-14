// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func VpnSiteName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[^'<>%&:?/+]+$`), "The value must not contain characters from '<>%&:?/+.")
}
