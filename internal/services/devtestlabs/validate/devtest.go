// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func DevTestLabName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[A-Za-z0-9_-]+$"),
		"Lab Name can only include alphanumeric characters, underscores, hyphens.",
	)
}

func DevTestVirtualMachineName(maxLength int) pluginsdk.SchemaValidateFunc {
	return validation.All(
		validation.StringLenBetween(1, maxLength),
		validation.StringMatch(regexp.MustCompile("^([a-zA-Z0-9]{1})([a-zA-Z0-9-]+)([a-zA-Z0-9]{1})$"), "may contain letters, numbers, or '-', must begin and end with a letter or number"),
		validation.StringMatch(regexp.MustCompile("([a-zA-Z-]+)"), "cannot be all numbers"),
	)
}

func DevTestVirtualNetworkUsagePermissionType() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice(virtualnetworks.PossibleValuesForUsagePermissionType(), false)
}
