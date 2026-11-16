// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopHostPoolUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := hostpool.ParseHostPoolID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.HostPoolName, hostPoolResourceType)
	defer locks.UnlockByName(id.HostPoolName, hostPoolResourceType)

	payload := hostpool.HostPoolPatch{}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	// lintignore:R019 // deliberate subset: only the fields feeding HostPoolPatchProperties; tags are applied separately above
	if d.HasChanges(
		"custom_rdp_properties",
		"description",
		"friendly_name",
		"load_balancer_type",
		"maximum_sessions_allowed",
		"preferred_app_group_type",
		"public_network_access",
		"start_vm_on_connect",
		"validate_environment",
		"scheduled_agent_updates",
	) {
		payload.Properties = &hostpool.HostPoolPatchProperties{}

		if d.HasChange("custom_rdp_properties") {
			payload.Properties.CustomRdpProperty = pointer.To(d.Get("custom_rdp_properties").(string))
		}

		if d.HasChange("description") {
			payload.Properties.Description = pointer.To(d.Get("description").(string))
		}

		if d.HasChange("friendly_name") {
			payload.Properties.FriendlyName = pointer.To(d.Get("friendly_name").(string))
		}

		if d.HasChange("load_balancer_type") {
			loadBalancerType := hostpool.LoadBalancerType(d.Get("load_balancer_type").(string))
			payload.Properties.LoadBalancerType = &loadBalancerType
		}

		if d.HasChange("maximum_sessions_allowed") {
			payload.Properties.MaxSessionLimit = pointer.To(int64(d.Get("maximum_sessions_allowed").(int)))
		}

		if d.HasChange("preferred_app_group_type") {
			preferredAppGroupType := hostpool.PreferredAppGroupType(d.Get("preferred_app_group_type").(string))
			payload.Properties.PreferredAppGroupType = &preferredAppGroupType
		}

		if d.HasChange("public_network_access") {
			payload.Properties.PublicNetworkAccess = pointer.ToEnum[hostpool.HostpoolPublicNetworkAccess](d.Get("public_network_access").(string))
		}

		if d.HasChange("start_vm_on_connect") {
			payload.Properties.StartVMOnConnect = pointer.To(d.Get("start_vm_on_connect").(bool))
		}

		if d.HasChange("validate_environment") {
			payload.Properties.ValidationEnvironment = pointer.To(d.Get("validate_environment").(bool))
		}

		if d.HasChanges("scheduled_agent_updates") {
			payload.Properties.AgentUpdate = expandAgentUpdatePatch(d.Get("scheduled_agent_updates").([]any))
		}

		if d.HasChanges("vm_template") {
			payload.Properties.VMTemplate = pointer.To(d.Get("vm_template").(string))
		}
	}

	if _, err := client.Update(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceVirtualDesktopHostPoolRead(d, meta)
}

func flattenAgentUpdate(input *hostpool.AgentUpdateProperties) []any {
	if input == nil {
		return []any{}
	}
	enabled := false
	if input.Type != nil {
		if *input.Type == hostpool.SessionHostComponentUpdateTypeScheduled {
			enabled = true
		}
	}

	return []any{
		map[string]any{
			"enabled":                   enabled,
			"timezone":                  input.MaintenanceWindowTimeZone,
			"use_session_host_timezone": input.UseSessionHostLocalTime,
			"schedule":                  flattenAgentUpdateSchedule(input.MaintenanceWindows),
		},
	}
}
