// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_scaling_plan_host_pool_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func virtualDesktopScalingPlanHostPoolAssociationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scaling_plan_id": {
			Description:  "The resource ID for the Virtual Desktop Scaling Plan.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: scalingplan.ValidateScalingPlanID,
		},

		"host_pool_id": {
			Description:  "The resource ID for the Virtual Desktop Host Pool.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: scalingplan.ValidateHostPoolID,
		},

		"enabled": {
			Description: "Should the Scaling Plan be enabled on this Host Pool?",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},
	}
}
