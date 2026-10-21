// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_connection_monitor

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectionmonitors"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkwatchers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNetworkConnectionMonitorRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.ConnectionMonitors
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := connectionmonitors.ParseConnectionMonitorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.ConnectionMonitorName)

	if model := resp.Model; model != nil {
		networkWatcherId := networkwatchers.NewNetworkWatcherID(id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName)
		d.Set("network_watcher_id", networkWatcherId.ID())

		if loc := model.Location; loc != nil {
			d.Set("location", location.Normalize(*loc))
		}

		if props := model.Properties; props != nil {
			if props.ConnectionMonitorType != nil && *props.ConnectionMonitorType == connectionmonitors.ConnectionMonitorTypeSingleSourceDestination {
				return errors.New("the resource created via API version 2019-06-01 or before (a.k.a v1) isn't compatible to this version of provider. Please migrate to v2 pluginsdk")
			}
			d.Set("notes", props.Notes)

			if err := d.Set("endpoint", flattenNetworkConnectionMonitorEndpoint(props.Endpoints)); err != nil {
				return fmt.Errorf("setting `endpoint`: %+v", err)
			}

			if err := d.Set("output_workspace_resource_ids", flattenNetworkConnectionMonitorOutput(props.Outputs)); err != nil {
				return fmt.Errorf("setting `output`: %+v", err)
			}

			if err := d.Set("test_configuration", flattenNetworkConnectionMonitorTestConfiguration(props.TestConfigurations)); err != nil {
				return fmt.Errorf("setting `test_configuration`: %+v", err)
			}

			if err := d.Set("test_group", flattenNetworkConnectionMonitorTestGroup(props.TestGroups)); err != nil {
				return fmt.Errorf("setting `test_group`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}
