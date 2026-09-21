// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func FilePath(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^(.)+.cap$`), "must end with extension name '.cap'")(v, k)
}
