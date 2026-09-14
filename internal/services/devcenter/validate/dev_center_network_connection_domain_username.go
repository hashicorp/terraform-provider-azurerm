// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DevCenterNetworkConnectionDomainUsername(i interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`), "is not a valid email")(i, k)
}
