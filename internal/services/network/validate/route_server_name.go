// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func RouteServerName() pluginsdk.SchemaValidateFunc {
	return validation.All(
		validation.StringIsNotEmpty,
		validation.StringLenBetween(1, 80),
		validation.StringMatch(
			regexp.MustCompile(`^[A-Za-z\d][A-Za-z\d.\-_]*[A-Za-z\d_]$`),
			"The name must begin with a letter or number, end with a letter, number or underscore, and may contain only letters, numbers, underscores, periods, or hyphens.",
		),
	)
}
