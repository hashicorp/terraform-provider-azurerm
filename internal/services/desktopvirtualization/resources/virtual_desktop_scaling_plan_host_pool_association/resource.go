// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_scaling_plan_host_pool_association

import (
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var scalingPlanResourceType = "azurerm_virtual_desktop_scaling_plan"

var hostPoolResourceType = "azurerm_virtual_desktop_host_pool"

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualDesktopScalingPlanHostPoolAssociationCreate,
		Read:   resourceVirtualDesktopScalingPlanHostPoolAssociationRead,
		Update: resourceVirtualDesktopScalingPlanHostPoolAssociationUpdate,
		Delete: resourceVirtualDesktopScalingPlanHostPoolAssociationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ScalingPlanHostPoolAssociationID(id)
			return err
		}),

		Schema: virtualDesktopScalingPlanHostPoolAssociationSchema(),
	}
}
