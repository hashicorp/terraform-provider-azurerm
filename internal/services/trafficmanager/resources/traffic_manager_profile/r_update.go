// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package traffic_manager_profile

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/trafficmanager/2022-04-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	trafficmanagerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmTrafficManagerProfileUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := trafficmanagerClient.From(meta).ProfilesClient
	ctx, cancel := timeouts.ForUpdate(account.StopContext(meta), d)
	defer cancel()

	id, err := profiles.ParseTrafficManagerProfileID(d.Id())
	if err != nil {
		return err
	}

	update := profiles.Profile{
		Properties: &profiles.ProfileProperties{},
	}
	if d.HasChange("tags") {
		update.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("profile_status") {
		status := profiles.ProfileStatus(d.Get("profile_status").(string))
		update.Properties.ProfileStatus = &status
	}

	if d.HasChange("traffic_routing_method") {
		routingMethod := profiles.TrafficRoutingMethod(d.Get("traffic_routing_method").(string))
		update.Properties.TrafficRoutingMethod = &routingMethod
	}

	if d.HasChange("max_return") {
		if maxReturn, ok := d.GetOk("max_return"); ok {
			update.Properties.MaxReturn = pointer.To(int64(maxReturn.(int)))
		}
	}

	if d.HasChange("dns_config") {
		update.Properties.DnsConfig = expandArmTrafficManagerDNSConfig(d)
	}

	if d.HasChange("monitor_config") {
		update.Properties.MonitorConfig = expandArmTrafficManagerMonitorConfig(d)
	}

	if d.HasChange("traffic_view_enabled") {
		update.Properties.TrafficViewEnrollmentStatus = expandArmTrafficManagerTrafficView(d.Get("traffic_view_enabled").(bool))
	}

	if _, err := client.Update(ctx, *id, update); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceArmTrafficManagerProfileRead(d, meta)
}
