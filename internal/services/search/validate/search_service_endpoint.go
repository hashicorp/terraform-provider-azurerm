// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SearchServiceEndpoint(v any, k string) ([]string, []error) {
	return validation.All(
		validation.IsURLWithHTTPS,
		validation.StringDoesNotMatch(regexp.MustCompile(`/$`), fmt.Sprintf("`%s` must not end with a trailing `/`", k)),
	)(v, k)
}
