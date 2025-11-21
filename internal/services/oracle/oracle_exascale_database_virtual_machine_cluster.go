// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package oracle

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/exadbvmclusters"
)

type ExascaleDatabaseDataCollectionModel struct {
	DiagnosticsEventsEnabled bool `tfschema:"diagnostics_events_enabled"`
	HealthMonitoringEnabled  bool `tfschema:"health_monitoring_enabled"`
	IncidentLogsEnabled      bool `tfschema:"incident_logs_enabled"`
}

type ExascaleDatabaseVirtualMachineClusterStorageModel struct {
	TotalSizeInGb int64 `tfschema:"total_size_in_gb"`
}

type NetworkSecurityGroupCidrModel struct {
	DestinationPortRange []PortRangeModel `tfschema:"destination_port_range"`
	Source               string           `tfschema:"source"`
}

type PortRangeModel struct {
	Max int64 `tfschema:"max"`
	Min int64 `tfschema:"min"`
}

func FlattenExadbDataCollectionOption(dataCollectionOptions *exadbvmclusters.DataCollectionOptions) []ExascaleDatabaseDataCollectionModel {
	output := make([]ExascaleDatabaseDataCollectionModel, 0)
	if dataCollectionOptions != nil {
		return append(output, ExascaleDatabaseDataCollectionModel{
			DiagnosticsEventsEnabled: pointer.From(dataCollectionOptions.IsDiagnosticsEventsEnabled),
			HealthMonitoringEnabled:  pointer.From(dataCollectionOptions.IsHealthMonitoringEnabled),
			IncidentLogsEnabled:      pointer.From(dataCollectionOptions.IsIncidentLogsEnabled),
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
