// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_packet_capture

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

func resourceVirtualMachinePacketCaptureCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.PacketCaptures
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	watcherId, err := networkwatchers.ParseNetworkWatcherID(d.Get("network_watcher_id").(string))
	if err != nil {
		return err
	}

	id := packetcaptures.NewPacketCaptureID(subscriptionId, watcherId.ResourceGroupName, watcherId.NetworkWatcherName, d.Get("name").(string))

	targetResourceId := d.Get("virtual_machine_id").(string)
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
			return tf.ImportAsExistsError("azurerm_virtual_machine_packet_capture", id.ID())
		}
	}

	storageLocation := expandVirtualMachinePacketCaptureStorageLocation(d.Get("storage_location").([]any))
	payload := packetcaptures.PacketCapture{
		Properties: packetcaptures.PacketCaptureParameters{
			Target:                  targetResourceId,
			TargetType:              pointer.To(packetcaptures.PacketCaptureTargetTypeAzureVM),
			StorageLocation:         storageLocation,
			BytesToCapturePerPacket: pointer.To(int64(bytesToCapturePerPacket)),
			TimeLimitInSeconds:      pointer.To(int64(timeLimitInSeconds)),
			TotalBytesPerSession:    pointer.To(int64(totalBytesPerSession)),
			Filters:                 expandVirtualMachinePacketCaptureFilters(d.Get("filter").([]any)),
		},
	}

	if err := client.CreateCallbackThenPoll(ctx, id, payload, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualMachinePacketCaptureRead(d, meta)
}

func expandVirtualMachinePacketCaptureStorageLocation(input []any) packetcaptures.PacketCaptureStorageLocation {
	location := input[0].(map[string]any)

	storageLocation := packetcaptures.PacketCaptureStorageLocation{}

	if v := location["file_path"]; v != "" {
		storageLocation.FilePath = pointer.To(v.(string))
	}
	if v := location["storage_account_id"]; v != "" {
		storageLocation.StorageId = pointer.To(v.(string))
	}

	return storageLocation
}

func expandVirtualMachinePacketCaptureFilters(input []any) *[]packetcaptures.PacketCaptureFilter {
	if len(input) == 0 {
		return nil
	}

	filters := make([]packetcaptures.PacketCaptureFilter, 0)

	for _, v := range input {
		inputFilter := v.(map[string]any)

		localIPAddress := inputFilter["local_ip_address"].(string)
		localPort := inputFilter["local_port"].(string) // TODO: should this be an int?
		protocol := inputFilter["protocol"].(string)
		remoteIPAddress := inputFilter["remote_ip_address"].(string)
		remotePort := inputFilter["remote_port"].(string)

		filter := packetcaptures.PacketCaptureFilter{
			LocalIPAddress:  pointer.To(localIPAddress),
			LocalPort:       pointer.To(localPort),
			Protocol:        pointer.ToEnum[packetcaptures.PcProtocol](protocol),
			RemoteIPAddress: pointer.To(remoteIPAddress),
			RemotePort:      pointer.To(remotePort),
		}
		filters = append(filters, filter)
	}

	return &filters
}
