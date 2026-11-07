// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_scaling_plan

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopScalingPlanUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scalingplan.ParseScalingPlanID(d.Id())
	if err != nil {
		return err
	}

	t := d.Get("tags").(map[string]any)

	payload := scalingplan.ScalingPlanPatch{
		Tags: tags.Expand(t),
		Properties: &scalingplan.ScalingPlanPatchProperties{
			Description:        pointer.To(d.Get("description").(string)),
			FriendlyName:       pointer.To(d.Get("friendly_name").(string)),
			TimeZone:           pointer.To(d.Get("time_zone").(string)),
			ExclusionTag:       pointer.To(d.Get("exclusion_tag").(string)),
			Schedules:          expandScalingPlanSchedule(d.Get("schedule").([]any)),
			HostPoolReferences: expandScalingPlanHostpoolReference(d.Get("host_pool").([]any)),
		},
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceVirtualDesktopScalingPlanRead(d, meta)
}
