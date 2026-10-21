// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_packet_capture

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/packetcaptures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandVirtualMachineScaleSetPacketCaptureStorageLocation(input []any) packetcaptures.PacketCaptureStorageLocation {
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

func expandVirtualMachineScaleSetPacketCaptureFilters(input []any) *[]packetcaptures.PacketCaptureFilter {
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

		filters = append(filters, packetcaptures.PacketCaptureFilter{
			LocalIPAddress:  pointer.To(localIPAddress),
			LocalPort:       pointer.To(localPort),
			Protocol:        pointer.ToEnum[packetcaptures.PcProtocol](protocol),
			RemoteIPAddress: pointer.To(remoteIPAddress),
			RemotePort:      pointer.To(remotePort),
		})
	}

	return &filters
}

func expandVirtualMachineScaleSetPacketCaptureMachineScope(input []any) *packetcaptures.PacketCaptureMachineScope {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]any)
	output := &packetcaptures.PacketCaptureMachineScope{}

	if exclude := raw["exclude_instance_ids"].([]any); len(exclude) > 0 {
		output.Exclude = pluginsdk.ExpandStringSlice(exclude)
	}

	if include := raw["include_instance_ids"].([]any); len(include) > 0 {
		output.Include = pluginsdk.ExpandStringSlice(include)
	}

	return output
}
