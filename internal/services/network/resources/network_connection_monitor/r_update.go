// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_connection_monitor

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectionmonitors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkConnectionMonitorUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.ConnectionMonitors
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connectionmonitors.ParseConnectionMonitorID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := connectionmonitors.ConnectionMonitor{
		Location: existing.Model.Location,
		Properties: connectionmonitors.ConnectionMonitorParameters{
			Endpoints:          existing.Model.Properties.Endpoints,
			TestConfigurations: existing.Model.Properties.TestConfigurations,
			TestGroups:         existing.Model.Properties.TestGroups,
			Notes:              existing.Model.Properties.Notes,
			Outputs:            existing.Model.Properties.Outputs,
		},
		Tags: existing.Model.Tags,
	}

	if d.HasChange("endpoint") {
		payload.Properties.Endpoints = expandNetworkConnectionMonitorEndpoint(d.Get("endpoint").(*pluginsdk.Set).List())
	}

	if d.HasChange("test_configuration") {
		payload.Properties.TestConfigurations = expandNetworkConnectionMonitorTestConfiguration(d.Get("test_configuration").(*pluginsdk.Set).List())
	}

	if d.HasChange("test_group") {
		payload.Properties.TestGroups = expandNetworkConnectionMonitorTestGroup(d.Get("test_group").(*pluginsdk.Set).List())
	}

	if d.HasChange("notes") {
		payload.Properties.Notes = pointer.To(d.Get("notes").(string))
	}

	if d.HasChange("output_workspace_resource_ids") {
		payload.Properties.Outputs = expandNetworkConnectionMonitorOutput(d.Get("output_workspace_resource_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if err = client.CreateOrUpdateThenPoll(ctx, *id, payload, connectionmonitors.DefaultCreateOrUpdateOperationOptions()); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkConnectionMonitorRead(d, meta)
}
