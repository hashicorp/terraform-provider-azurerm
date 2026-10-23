// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_packet_capture

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2025-04-01/virtualmachinescalesetvms"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/packetcaptures"
)

func flattenVirtualMachineScaleSetPacketCaptureStorageLocation(input packetcaptures.PacketCaptureStorageLocation) []any {
	return []any{
		map[string]any{
			"file_path":          pointer.From(input.FilePath),
			"storage_account_id": pointer.From(input.StorageId),
			"storage_path":       pointer.From(input.StoragePath),
		},
	}
}

func flattenVirtualMachineScaleSetPacketCaptureFilters(input *[]packetcaptures.PacketCaptureFilter) []any {
	filters := make([]any, 0)

	if inFilter := input; inFilter != nil {
		for _, v := range *inFilter {
			protocol := ""
			if v.Protocol != nil {
				protocol = string(*v.Protocol)
			}

			filters = append(filters, map[string]any{
				"local_ip_address":  pointer.From(v.LocalIPAddress),
				"local_port":        pointer.From(v.LocalPort),
				"protocol":          protocol,
				"remote_ip_address": pointer.From(v.RemoteIPAddress),
				"remote_port":       pointer.From(v.RemotePort),
			})
		}
	}

	return filters
}

func flattenVirtualMachineScaleSetPacketCaptureMachineScope(input *packetcaptures.PacketCaptureMachineScope) ([]any, error) {
	outputs := make([]any, 0)
	if input == nil || (input.Exclude == nil && input.Include == nil) || (len(*input.Exclude) == 0 && len(*input.Include) == 0) {
		return outputs, nil
	}

	excludedInstanceIds, err := flattenVirtualMachineScaleSetPacketCaptureScopeInstanceIds(input.Exclude)
	if err != nil {
		return nil, err
	}

	includedInstanceIds, err := flattenVirtualMachineScaleSetPacketCaptureScopeInstanceIds(input.Include)
	if err != nil {
		return nil, err
	}

	outputs = append(outputs, map[string]any{
		"exclude_instance_ids": excludedInstanceIds,
		"include_instance_ids": includedInstanceIds,
	})

	return outputs, nil
}

func flattenVirtualMachineScaleSetPacketCaptureScopeInstanceIds(input *[]string) ([]string, error) {
	instances := make([]string, 0)
	if input == nil {
		return instances, nil
	}

	for _, instance := range *input {
		vmssInstanceId, err := virtualmachinescalesetvms.ParseVirtualMachineScaleSetVirtualMachineIDInsensitively(instance)
		if err != nil {
			return nil, err
		}

		instances = append(instances, vmssInstanceId.InstanceId)
	}

	return instances, nil
}
