// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func BatchAccountIpRange(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(/([0-9]|[1-2][0-9]|30))?$`), "must start with IPV4 address and/or slash, number of bits (0-30) as prefix. Example: 127.0.0.1/8")(v, k)
}
