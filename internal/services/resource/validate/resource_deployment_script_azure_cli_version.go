// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ResourceDeploymentScriptAzureCliVersion(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^\d+\.\d+\.\d+$`), "should be in the format `X.Y.Z` (e.g. `2.30.0`)")(i, k)
}
