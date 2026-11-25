// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopHostPoolRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hostpool.ParseHostPoolID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.HostPoolName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}

		props := model.Properties
		maxSessionLimit := 0
		if props.MaxSessionLimit != nil {
			maxSessionLimit = int(*props.MaxSessionLimit)
		}

		d.Set("custom_rdp_properties", props.CustomRdpProperty)
		d.Set("description", props.Description)
		d.Set("friendly_name", props.FriendlyName)
		d.Set("maximum_sessions_allowed", maxSessionLimit)
		d.Set("load_balancer_type", string(props.LoadBalancerType))
		personalDesktopAssignmentType := ""
		if props.PersonalDesktopAssignmentType != nil {
			personalDesktopAssignmentType = string(*props.PersonalDesktopAssignmentType)
		}
		d.Set("personal_desktop_assignment_type", personalDesktopAssignmentType)
		d.Set("preferred_app_group_type", string(props.PreferredAppGroupType))
		d.Set("public_network_access", pointer.FromEnum(props.PublicNetworkAccess))
		d.Set("start_vm_on_connect", props.StartVMOnConnect)
		d.Set("type", string(props.HostPoolType))
		d.Set("validate_environment", props.ValidationEnvironment)
		d.Set("scheduled_agent_updates", flattenAgentUpdate(props.AgentUpdate))
		d.Set("vm_template", props.VMTemplate)
	}

	return nil
}

func flattenAgentUpdateSchedule(input *[]hostpool.MaintenanceWindowProperties) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		dayOfWeek := ""
		if item.DayOfWeek != nil {
			dayOfWeek = string(*item.DayOfWeek)
		}
		hourOfDay := pointer.To(int64(0))
		if item.Hour != nil {
			hourOfDay = item.Hour
		}
		results = append(results, map[string]any{
			"day_of_week": dayOfWeek,
			"hour_of_day": hourOfDay,
		})
	}
	return results
}

var hostPoolResourceType = "azurerm_virtual_desktop_host_pool"
