// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_packet_capture

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkwatchers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/packetcaptures"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineScaleSetPacketCaptureCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PacketCaptures
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	watcherId, err := networkwatchers.ParseNetworkWatcherID(d.Get("network_watcher_id").(string))
	if err != nil {
		return err
	}

	id := packetcaptures.NewPacketCaptureID(subscriptionId, watcherId.ResourceGroupName, watcherId.NetworkWatcherName, d.Get("name").(string))

	targetResourceId := d.Get("virtual_machine_scale_set_id").(string)
	bytesToCapturePerPacket := d.Get("maximum_bytes_per_packet").(int)
	totalBytesPerSession := d.Get("maximum_bytes_per_session").(int)
	timeLimitInSeconds := d.Get("maximum_capture_duration_in_seconds").(int)

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_machine_scale_set_packet_capture", id.ID())
		}
	}

	storageLocation := expandVirtualMachineScaleSetPacketCaptureStorageLocation(d.Get("storage_location").([]any))
	payload := packetcaptures.PacketCapture{
		Properties: packetcaptures.PacketCaptureParameters{
			Target:                  targetResourceId,
			TargetType:              pointer.To(packetcaptures.PacketCaptureTargetTypeAzureVMSS),
			StorageLocation:         storageLocation,
			BytesToCapturePerPacket: pointer.To(int64(bytesToCapturePerPacket)),
			TimeLimitInSeconds:      pointer.To(int64(timeLimitInSeconds)),
			TotalBytesPerSession:    pointer.To(int64(totalBytesPerSession)),
			Filters:                 expandVirtualMachineScaleSetPacketCaptureFilters(d.Get("filter").([]any)),
		},
	}

	if v, ok := d.GetOk("machine_scope"); ok {
		payload.Properties.Scope = expandVirtualMachineScaleSetPacketCaptureMachineScope(v.([]any))
	}

	if err := client.CreateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualMachineScaleSetPacketCaptureRead(d, meta)
}
