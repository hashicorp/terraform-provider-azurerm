package redisenterprise

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DatabaseProperties struct {
	ClientProtocol    *Protocol                         `json:"clientProtocol,omitempty"`
	ClusteringPolicy  *ClusteringPolicy                 `json:"clusteringPolicy,omitempty"`
	EvictionPolicy    *EvictionPolicy                   `json:"evictionPolicy,omitempty"`
	GeoReplication    *DatabasePropertiesGeoReplication `json:"geoReplication"`
	Modules           *[]Module                         `json:"modules,omitempty"`
	Persistence       *Persistence                      `json:"persistence"`
	Port              *int64                            `json:"port,omitempty"`
	ProvisioningState *ProvisioningState                `json:"provisioningState,omitempty"`
	ResourceState     *ResourceState                    `json:"resourceState,omitempty"`
}
