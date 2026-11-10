// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_scaling_plan_host_pool_association

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopScalingPlanHostPoolAssociationCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	scalingPlanId, err := scalingplan.ParseScalingPlanID(d.Get("scaling_plan_id").(string))
	if err != nil {
		return err
	}
	hostPoolId, err := scalingplan.ParseHostPoolID(d.Get("host_pool_id").(string))
	if err != nil {
		return err
	}
	associationId := parse.NewScalingPlanHostPoolAssociationId(*scalingPlanId, *hostPoolId).ID()

	locks.ByName(scalingPlanId.ScalingPlanName, scalingPlanResourceType)
	defer locks.UnlockByName(scalingPlanId.ScalingPlanName, scalingPlanResourceType)

	locks.ByName(hostPoolId.HostPoolName, hostPoolResourceType)
	defer locks.UnlockByName(hostPoolId.HostPoolName, hostPoolResourceType)

	existing, err := client.Get(ctx, *scalingPlanId)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("%s was not found", *scalingPlanId)
		}

		return fmt.Errorf("retrieving %s: %+v", *scalingPlanId, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: model was nil", *scalingPlanId)
	}
	model := *existing.Model

	hostPoolAssociations := []scalingplan.ScalingHostPoolReference{}
	if v := model.Properties.HostPoolReferences; v != nil {
		hostPoolAssociations = *v
	}

	hostPoolStr := hostPoolId.ID()
	if scalingPlanHostPoolAssociationExists(model.Properties, hostPoolStr) {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_scaling_plan_host_pool_association", associationId)
		}
	}
	hostPoolAssociations = append(hostPoolAssociations, scalingplan.ScalingHostPoolReference{
		HostPoolArmPath:    &hostPoolStr,
		ScalingPlanEnabled: pointer.To(d.Get("enabled").(bool)),
	})

	payload := scalingplan.ScalingPlanPatch{
		Properties: &scalingplan.ScalingPlanPatchProperties{
			HostPoolReferences: &hostPoolAssociations,
			Schedules:          model.Properties.Schedules,
		},
		Tags: model.Tags,
	}
	if _, err = client.Update(ctx, *scalingPlanId, payload); err != nil {
		return fmt.Errorf("creating association between %s and %s: %+v", *scalingPlanId, *hostPoolId, err)
	}

	d.SetId(associationId)
	return resourceVirtualDesktopScalingPlanHostPoolAssociationRead(d, meta)
}
