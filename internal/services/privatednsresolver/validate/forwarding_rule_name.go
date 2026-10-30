// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ForwardingRuleName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]{2,79}$`),
		"must be between 3 and 80 characters in length, begin with a letter, and may contain only letters, numbers, underscores or hyphens",
	)
}
