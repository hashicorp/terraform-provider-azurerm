// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ApplicationLoadBalancerSubnetAssociationName() pluginsdk.SchemaValidateFunc {
	return validation.All(
		validation.StringLenBetween(1, 64),
		validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9]`), "the name must begin with a letter or number."),
		validation.StringMatch(regexp.MustCompile(`[a-zA-Z0-9]$`), "the name must end with a letter or number."),
		validation.StringMatch(regexp.MustCompile(`[a-zA-Z0-9_.-]{0,64}`), "the name may contain only letters, numbers, underscores, periods, or hyphens."),
	)
}
