// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CapabilityHostName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_-]{0,254}$"),
		"`name` must be between 1 and 255 characters long, start with an alphanumeric character, and contain only alphanumeric characters, dashes(-), or underscores(_)",
	)
}
