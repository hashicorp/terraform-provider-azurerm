package scclusterrecords

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SCClusterSpecEntity struct {
	ApiEndpoint            *string                            `json:"apiEndpoint,omitempty"`
	Availability           *string                            `json:"availability,omitempty"`
	Byok                   *SCClusterByokEntity               `json:"byok,omitempty"`
	Cloud                  *string                            `json:"cloud,omitempty"`
	Config                 *ClusterConfigEntity               `json:"config,omitempty"`
	Environment            *SCClusterNetworkEnvironmentEntity `json:"environment,omitempty"`
	HTTPEndpoint           *string                            `json:"httpEndpoint,omitempty"`
	KafkaBootstrapEndpoint *string                            `json:"kafkaBootstrapEndpoint,omitempty"`
	Name                   *string                            `json:"name,omitempty"`
	Network                *SCClusterNetworkEnvironmentEntity `json:"network,omitempty"`
	Package                *Package                           `json:"package,omitempty"`
	Region                 *string                            `json:"region,omitempty"`
	Zone                   *string                            `json:"zone,omitempty"`
}
