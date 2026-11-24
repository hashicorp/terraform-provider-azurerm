// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_scaling_plan_host_pool_association

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopScalingPlanHostPoolAssociationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScalingPlanHostPoolAssociationID(d.Id())
	if err != nil {
		return err
	}

	scalingPlan, err := client.Get(ctx, id.ScalingPlan)
	if err != nil {
		if response.WasNotFound(scalingPlan.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id.ScalingPlan)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id.ScalingPlan, err)
	}
	if model := scalingPlan.Model; model != nil {
		hostPoolId := id.HostPool.ID()
		exists := scalingPlanHostPoolAssociationExists(model.Properties, hostPoolId)
		if !exists {
			log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.ScalingPlan, id.HostPool)
			d.SetId("")
			return nil
		}
		if v := model.Properties.HostPoolReferences; v != nil {
			for _, referenceId := range *v {
				if referenceId.HostPoolArmPath != nil {
					if strings.EqualFold(*referenceId.HostPoolArmPath, hostPoolId) {
						d.Set("enabled", referenceId.ScalingPlanEnabled)
					}
				}
			}
		}

		d.Set("scaling_plan_id", id.ScalingPlan.ID())
		d.Set("host_pool_id", hostPoolId)
	}

	return nil
}

func scalingPlanHostPoolAssociationExists(props scalingplan.ScalingPlanProperties, applicationGroupId string) bool {
	if props.HostPoolReferences == nil {
		return false
	}

	for _, id := range *props.HostPoolReferences {
		if id.HostPoolArmPath != nil {
			if strings.EqualFold(*id.HostPoolArmPath, applicationGroupId) {
				return true
			}
		}
	}

	return false
}
