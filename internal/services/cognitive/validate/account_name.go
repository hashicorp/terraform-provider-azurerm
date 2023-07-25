// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AccountName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^([a-zA-Z0-9]{1}[a-zA-Z0-9_.-]{1,})$"),
		"The Cognitive Services Account Name can only start with an alphanumeric character, and must only contain alphanumeric characters, periods, dashes or underscores.",
	)
}
