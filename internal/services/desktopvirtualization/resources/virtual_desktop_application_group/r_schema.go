// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_application_group

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/applicationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualDesktopApplicationGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Virtual Desktop Application Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApplicationGroupName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"type": {
			Description:  "Type of Virtual Desktop Application Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(applicationgroup.PossibleValuesForApplicationGroupType(), false),
		},

		"host_pool_id": {
			Description:  "Resource ID for a Virtual Desktop Host Pool to associate with the Virtual Desktop Application Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: hostpool.ValidateHostPoolID,
		},

		"friendly_name": {
			Description:  "Option to set a friendly name for the Virtual Desktop Application Group.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},

		"default_desktop_display_name": {
			Description:  "Option to set the display name for the default sessionDesktop desktop when `type` is set to `Desktop`. A value here is mandatory for connections to the desktop using the Windows 365 portal. Without it the connection will hang at 'Loading Client'.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Description:  "Option to set a description for the Virtual Desktop Application Group.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"tags": commonschema.Tags(),
	}
}
