// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func VideoAnalyzerName() func(interface{}, string) ([]string, []error) {
	return validation.StringMatch(
		regexp.MustCompile("^[-a-z0-9]{3,24}$"),
		"Video Analyzer name must be 3 - 24 characters long, contain only lowercase letters and numbers.",
	)
}
