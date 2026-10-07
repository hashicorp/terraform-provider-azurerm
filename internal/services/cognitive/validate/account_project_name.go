// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AccountProjectName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_.-]{1,63}$"),
		"The Cognitive Services Account Project Name must be between 2 and 64 characters long, start with an alphanumeric character, and contain only alphanumeric characters, dashes(-), periods(.) or underscores(_).",
	)
}
