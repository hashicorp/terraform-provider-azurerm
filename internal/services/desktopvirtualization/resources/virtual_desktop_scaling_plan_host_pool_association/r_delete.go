// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_scaling_plan_host_pool_association

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopScalingPlanHostPoolAssociationDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ScalingPlanHostPoolAssociationID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.ScalingPlan.ScalingPlanName, scalingPlanResourceType)
	defer locks.UnlockByName(id.ScalingPlan.ScalingPlanName, scalingPlanResourceType)

	locks.ByName(id.HostPool.HostPoolName, hostPoolResourceType)
	defer locks.UnlockByName(id.HostPool.HostPoolName, hostPoolResourceType)

	existing, err := client.Get(ctx, id.ScalingPlan)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("%s was not found", id.ScalingPlan)
		}

		return fmt.Errorf("retrieving %s: %+v", id.ScalingPlan, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", id.ScalingPlan)
	}
	model := *existing.Model

	hostPoolReferences := []scalingplan.ScalingHostPoolReference{}
	hostPoolId := id.HostPool.ID()
	if v := model.Properties.HostPoolReferences; v != nil {
		for _, referenceId := range *v {
			if referenceId.HostPoolArmPath != nil {
				if strings.EqualFold(*referenceId.HostPoolArmPath, hostPoolId) {
					continue
				}
			}

			hostPoolReferences = append(hostPoolReferences, referenceId)
		}
	}

	payload := scalingplan.ScalingPlanPatch{
		Properties: &scalingplan.ScalingPlanPatchProperties{
			HostPoolReferences: &hostPoolReferences,
			Schedules:          model.Properties.Schedules,
		},
		Tags: model.Tags,
	}
	if _, err = client.Update(ctx, id.ScalingPlan, payload); err != nil {
		return fmt.Errorf("removing association between %s and %s: %+v", id.ScalingPlan, id.HostPool, err)
	}

	return nil
}
