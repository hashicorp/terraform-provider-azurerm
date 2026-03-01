// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func RaiBlocklistName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_-]{0,62}[a-zA-Z0-9]$"),
		"The RAI Blocklist Name must be between 2 and 64 characters long, start and end with an alphanumeric character, and contain only alphanumeric characters, hyphens(-) or underscores(_).",
	)
}
