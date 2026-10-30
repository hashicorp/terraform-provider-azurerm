// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_scaling_plan

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopScalingPlanCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.ScalingPlansClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := scalingplan.NewScalingPlanID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s): %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_scaling_plan", id.ID())
		}
	}

	location := location.Normalize(d.Get("location").(string))
	t := d.Get("tags").(map[string]any)

	payload := scalingplan.ScalingPlan{
		Name:     pointer.To(d.Get("name").(string)),
		Location: location,
		Tags:     tags.Expand(t),
		Properties: scalingplan.ScalingPlanProperties{
			Description:        pointer.To(d.Get("description").(string)),
			FriendlyName:       pointer.To(d.Get("friendly_name").(string)),
			TimeZone:           d.Get("time_zone").(string),
			HostPoolType:       pointer.To(scalingplan.ScalingHostPoolTypePooled),
			ExclusionTag:       pointer.To(d.Get("exclusion_tag").(string)),
			Schedules:          expandScalingPlanSchedule(d.Get("schedule").([]any)),
			HostPoolReferences: expandScalingPlanHostpoolReference(d.Get("host_pool").([]any)),
		},
	}

	if _, err := client.Create(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualDesktopScalingPlanRead(d, meta)
}
