// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_packet_capture

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkwatchers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/packetcaptures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineScaleSetPacketCaptureRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PacketCaptures
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := packetcaptures.ParsePacketCaptureID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.PacketCaptureName)
	d.Set("network_watcher_id", networkwatchers.NewNetworkWatcherID(id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("virtual_machine_scale_set_id", props.Target)
			d.Set("maximum_bytes_per_packet", int(*props.BytesToCapturePerPacket))
			d.Set("maximum_bytes_per_session", int(*props.TotalBytesPerSession))
			d.Set("maximum_capture_duration_in_seconds", int(*props.TimeLimitInSeconds))

			if err := d.Set("storage_location", flattenVirtualMachineScaleSetPacketCaptureStorageLocation(props.StorageLocation)); err != nil {
				return fmt.Errorf("setting `storage_location`: %+v", err)
			}

			if err := d.Set("filter", flattenVirtualMachineScaleSetPacketCaptureFilters(props.Filters)); err != nil {
				return fmt.Errorf("setting `filter`: %+v", err)
			}

			scope, err := flattenVirtualMachineScaleSetPacketCaptureMachineScope(props.Scope)
			if err != nil {
				return err
			}
			if err := d.Set("machine_scope", scope); err != nil {
				return fmt.Errorf(`setting "machine_scope": %+v`, err)
			}
		}
	}

	return nil
}
