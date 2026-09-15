// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DatadogMonitorsEmailAddress(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9._%+-]+@(?:[A-Za-z0-9-]+\.)+[A-Za-z]{2,}$`), "expected value not match regular expression")(i, k)
}
