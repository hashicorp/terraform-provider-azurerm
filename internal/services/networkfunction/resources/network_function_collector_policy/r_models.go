// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_function_collector_policy

type NetworkFunctionCollectorPolicyModel struct {
	Name                                   string               `tfschema:"name"`
	NetworkFunctionAzureTrafficCollectorId string               `tfschema:"traffic_collector_id"`
	IpfxEmission                           []IpfxEmissionModel  `tfschema:"ipfx_emission"`
	IpfxIngestion                          []IpfxIngestionModel `tfschema:"ipfx_ingestion"`
	Location                               string               `tfschema:"location"`
	Tags                                   map[string]any       `tfschema:"tags"`
}

type IpfxEmissionModel struct {
	DestinationTypes []string `tfschema:"destination_types"`
}

type IpfxIngestionModel struct {
	SourceResourceIds []string `tfschema:"source_resource_ids"`
}

func (r Resource) ModelObject() any {
	return &NetworkFunctionCollectorPolicyModel{}
}
