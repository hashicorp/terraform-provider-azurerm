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

type NetworkSecurityGroupRuleModel struct {
	DestinationPortRange []PortRangeModel `tfschema:"destination_port_range"`
	SourceIpRange        string           `tfschema:"source_ip_range"`
}

type PortRangeModel struct {
	Maximum int64 `tfschema:"maximum"`
	Minimum int64 `tfschema:"minimum"`
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

func FlattenNetworkSecurityGroupCidr(input *[]exadbvmclusters.NsgCidr) []NetworkSecurityGroupRuleModel {
	output := make([]NetworkSecurityGroupRuleModel, 0)

	if input != nil {
		for _, nsgCidr := range *input {
			var portRangeModel []PortRangeModel
			if nsgCidr.DestinationPortRange != nil {
				portRangeModel = append(portRangeModel, PortRangeModel{
					Maximum: nsgCidr.DestinationPortRange.Max,
					Minimum: nsgCidr.DestinationPortRange.Min,
				})
			}
			output = append(output, NetworkSecurityGroupRuleModel{
				DestinationPortRange: portRangeModel,
				SourceIpRange:        nsgCidr.Source,
			})
		}
	}

	return output
}
