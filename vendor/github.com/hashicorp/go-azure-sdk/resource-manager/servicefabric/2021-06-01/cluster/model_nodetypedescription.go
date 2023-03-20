package cluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodeTypeDescription struct {
	ApplicationPorts             *EndpointRangeDescription `json:"applicationPorts,omitempty"`
	Capacities                   *map[string]string        `json:"capacities,omitempty"`
	ClientConnectionEndpointPort int64                     `json:"clientConnectionEndpointPort"`
	DurabilityLevel              *DurabilityLevel          `json:"durabilityLevel,omitempty"`
	EphemeralPorts               *EndpointRangeDescription `json:"ephemeralPorts,omitempty"`
	HTTPGatewayEndpointPort      int64                     `json:"httpGatewayEndpointPort"`
	IsPrimary                    bool                      `json:"isPrimary"`
	IsStateless                  *bool                     `json:"isStateless,omitempty"`
	MultipleAvailabilityZones    *bool                     `json:"multipleAvailabilityZones,omitempty"`
	Name                         string                    `json:"name"`
	PlacementProperties          *map[string]string        `json:"placementProperties,omitempty"`
	ReverseProxyEndpointPort     *int64                    `json:"reverseProxyEndpointPort,omitempty"`
	VMInstanceCount              int64                     `json:"vmInstanceCount"`
}
