// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exadbvmclusters"
)

func FlattenExadbDataCollectionOption(dataCollectionOptions *exadbvmclusters.DataCollectionOptions) []ExascaleDatabaseDataCollectionOptionModel {
	output := make([]ExascaleDatabaseDataCollectionOptionModel, 0)
	if dataCollectionOptions != nil {
		return append(output, ExascaleDatabaseDataCollectionOptionModel{
			IsDiagnosticsEventsEnabled: pointer.From(dataCollectionOptions.IsDiagnosticsEventsEnabled),
			IsHealthMonitoringEnabled:  pointer.From(dataCollectionOptions.IsHealthMonitoringEnabled),
			IsIncidentLogsEnabled:      pointer.From(dataCollectionOptions.IsIncidentLogsEnabled),
		})
	}
	return output
}

func FlattenVMFileSystemStorage(input exadbvmclusters.ExadbVMClusterStorageDetails) []ExascaleDatabaseVirtualMachineClusterStorageModel {
	output := make([]ExascaleDatabaseVirtualMachineClusterStorageModel, 0)
	return append(output, ExascaleDatabaseVirtualMachineClusterStorageModel{
		TotalSizeInGb: input.TotalSizeInGbs,
	})
}

func FlattenNetworkSecurityGroupCidr(input *[]exadbvmclusters.NsgCidr) []NetworkSecurityGroupCidrModel {
	output := make([]NetworkSecurityGroupCidrModel, 0)

	if input != nil {
		for _, nsgCidr := range *input {
			var portRangeModel []PortRangeModel
			if nsgCidr.DestinationPortRange != nil {
				portRangeModel = append(portRangeModel, PortRangeModel{
					Max: nsgCidr.DestinationPortRange.Max,
					Min: nsgCidr.DestinationPortRange.Min,
				})
			}
			output = append(output, NetworkSecurityGroupCidrModel{
				DestinationPortRange: portRangeModel,
				Source:               nsgCidr.Source,
			})
		}
	}

	return output
}
