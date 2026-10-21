// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_connection_monitor

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectionmonitors"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkConnectionMonitorCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.ConnectionMonitors
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	watcherId, err := connectionmonitors.ParseNetworkWatcherID(d.Get("network_watcher_id").(string))
	if err != nil {
		return err
	}

	connectionMonitorId := connectionmonitors.NewConnectionMonitorID(subscriptionId, watcherId.ResourceGroupName, watcherId.NetworkWatcherName, d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, connectionMonitorId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", connectionMonitorId, err)
			}
		}

		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_network_connection_monitor", connectionMonitorId.ID())
		}
	}

	properties := connectionmonitors.ConnectionMonitor{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
		Properties: connectionmonitors.ConnectionMonitorParameters{
			Outputs:            expandNetworkConnectionMonitorOutput(d.Get("output_workspace_resource_ids").(*pluginsdk.Set).List()),
			TestConfigurations: expandNetworkConnectionMonitorTestConfiguration(d.Get("test_configuration").(*pluginsdk.Set).List()),
			TestGroups:         expandNetworkConnectionMonitorTestGroup(d.Get("test_group").(*pluginsdk.Set).List()),
		},
	}

	properties.Properties.Endpoints = expandNetworkConnectionMonitorEndpoint(d.Get("endpoint").(*pluginsdk.Set).List())

	if notes, ok := d.GetOk("notes"); ok {
		properties.Properties.Notes = pointer.To(notes.(string))
	}

	if err = client.CreateOrUpdateCallbackThenPoll(ctx, connectionMonitorId, properties, connectionmonitors.DefaultCreateOrUpdateOperationOptions(), sdk.SetIDCallback(meta, &connectionMonitorId, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", connectionMonitorId, err)
	}

	d.SetId(connectionMonitorId.ID())

	return resourceNetworkConnectionMonitorRead(d, meta)
}
