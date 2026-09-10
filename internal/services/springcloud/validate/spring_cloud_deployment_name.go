// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SpringCloudDeploymentName(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([a-z])([a-z\d-]{2,30})([a-z\d])$`), "must begin with a letter, end with a letter or number, contain only lowercase letters, numbers and hyphens. The value must be between 4 and 32 characters long")(i, k)
}
