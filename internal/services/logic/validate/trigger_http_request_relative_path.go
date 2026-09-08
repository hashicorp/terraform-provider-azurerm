// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func TriggerHttpRequestRelativePath(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile("^[A-Za-z0-9_/}{]+$"), "can only contain alphanumeric characters, underscores, forward slashes and curly braces")(v, k)
}
