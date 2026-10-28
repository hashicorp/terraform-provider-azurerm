// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_application

import (
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/application"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/applicationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualDesktopApplicationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				// NOTE: name string length is incorrect in the API specs, see: https://github.com/Azure/azure-rest-api-specs/issues/46813
				validation.StringLenBetween(1, 260),
				validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9@._ -]*$`),
					"Virtual desktop application name must be 1 - 260 characters long and may only contain letters, numbers, spaces, periods, underscores, hyphens, and @.",
				),
			),
		},

		"application_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: applicationgroup.ValidateApplicationGroupID,
		},

		"friendly_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
			// NOTE: O+C The API will use the value in `name` as the default
			Computed: true,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"path": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"command_line_argument_policy": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(application.PossibleValuesForCommandLineSetting(), false),
		},

		"command_line_arguments": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"show_in_portal": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"icon_path": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C The API will use the value in `path` as the default
			Computed: true,
		},

		"icon_index": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},
	}
}
