// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_host_pool

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/hostpool"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualDesktopHostPoolCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DesktopVirtualization.HostPoolsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := hostpool.NewHostPoolID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_desktop_host_pool", id.ID())
		}
	}

	personalDesktopAssignmentType := hostpool.PersonalDesktopAssignmentType(d.Get("personal_desktop_assignment_type").(string))
	vmTemplate := d.Get("vm_template").(string)
	if vmTemplate != "" {
		// we have no use with the json object as azure accepts string only
		// merely here for validation
		if _, err := pluginsdk.ExpandJsonFromString(vmTemplate); err != nil {
			return fmt.Errorf("expanding JSON for `vm_template`: %+v", err)
		}
	}
	payload := hostpool.HostPool{
		Location: location.Normalize(d.Get("location").(string)),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
		Properties: hostpool.HostPoolProperties{
			HostPoolType:                  hostpool.HostPoolType(d.Get("type").(string)),
			FriendlyName:                  pointer.To(d.Get("friendly_name").(string)),
			Description:                   pointer.To(d.Get("description").(string)),
			ValidationEnvironment:         pointer.To(d.Get("validate_environment").(bool)),
			CustomRdpProperty:             pointer.To(d.Get("custom_rdp_properties").(string)),
			MaxSessionLimit:               pointer.To(int64(d.Get("maximum_sessions_allowed").(int))),
			StartVMOnConnect:              pointer.To(d.Get("start_vm_on_connect").(bool)),
			LoadBalancerType:              hostpool.LoadBalancerType(d.Get("load_balancer_type").(string)),
			PersonalDesktopAssignmentType: &personalDesktopAssignmentType,
			PreferredAppGroupType:         hostpool.PreferredAppGroupType(d.Get("preferred_app_group_type").(string)),
			PublicNetworkAccess:           pointer.ToEnum[hostpool.HostpoolPublicNetworkAccess](d.Get("public_network_access").(string)),
			AgentUpdate:                   expandAgentUpdateCreate(d.Get("scheduled_agent_updates").([]any)),
			VMTemplate:                    &vmTemplate,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceVirtualDesktopHostPoolRead(d, meta)
}

func expandAgentUpdateCreate(input []any) *hostpool.AgentUpdateProperties {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)

	props := hostpool.AgentUpdateProperties{}

	useSessionHostLocalTime := *pointer.To(raw["use_session_host_timezone"].(bool))
	updateScheduleTimeZone := pointer.To(raw["timezone"].(string))

	if raw["enabled"].(bool) {
		props.Type = pointer.To(hostpool.SessionHostComponentUpdateTypeScheduled)
		if !useSessionHostLocalTime { // based on the priority used in the Azure Portal, if Session Host time is selected, this overrides the explicit TimeZone setting
			props.MaintenanceWindowTimeZone = updateScheduleTimeZone
			props.UseSessionHostLocalTime = &useSessionHostLocalTime
			props.MaintenanceWindows = expandAgentUpdateSchedule(raw["schedule"].([]any))
		}
	} else {
		props.Type = pointer.To(hostpool.SessionHostComponentUpdateTypeDefault)
		props.MaintenanceWindows = &[]hostpool.MaintenanceWindowProperties{}
		props.UseSessionHostLocalTime = &useSessionHostLocalTime // required by REST API even when set to Default/Disabled
		props.MaintenanceWindowTimeZone = updateScheduleTimeZone // required by REST API even when set to Default/Disabled
	}

	return &props
}
