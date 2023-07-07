// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SubscriptionName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[_a-zA-Z0-9][-._a-zA-Z0-9]{0,48}([_a-zA-Z0-9])?$"),
		"The name can contain only letters, numbers, periods, hyphens and underscores. The name must start and end with a letter, number or underscore and be a maximum of 50 characters long.",
	)
}
