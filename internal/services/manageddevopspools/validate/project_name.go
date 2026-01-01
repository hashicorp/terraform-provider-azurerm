// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ProjectName(i interface{}, k string) (warnings []string, errors []error) {
	return validation.All(
		validation.StringLenBetween(1, 64),
		validation.StringMatch(
			regexp.MustCompile(`^[^_.]`),
			"project name must not start with underscore or period",
		),
		validation.StringMatch(
			regexp.MustCompile(`[^.]$`),
			"project name must not end with period",
		),
		validation.StringMatch(
			regexp.MustCompile(`^[^\\/:*?"'<>;#$*{}+=\[\]|,\x00-\x1F\x7F]*$`),
			"project name must not contain special characters: \\ / : * ? \" ' < > ; # $ * { } + = [ ] | , or control characters",
		),
		validation.StringNotInSlice([]string{
			"App_Browsers", "App_Code", "App_Data", "App_GlobalResources",
			"App_LocalResources", "App_Themes", "App_WebResources",
			"bin", "web.config",
		}, true),
	)(i, k)
}
