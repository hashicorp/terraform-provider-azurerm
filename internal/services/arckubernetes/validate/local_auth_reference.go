// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func LocalAuthReference(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(
		regexp.MustCompile(`^[a-z\d]([-a-z\d]{0,61}[a-z\d])?$`),
		"`local_auth_reference` must be between 1 and 63 characters. It can contain only lowercase letters, numbers, and hyphens (-). It must start and end with a lowercase letter or number.",
	)(i, k)
}
