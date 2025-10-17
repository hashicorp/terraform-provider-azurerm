package organizationresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ClusterSpecEntity struct {
	ApiEndpoint            *string                   `json:"api_endpoint,omitempty"`
	Availability           *string                   `json:"availability,omitempty"`
	Byok                   *ClusterByokEntity        `json:"byok,omitempty"`
	Cloud                  *string                   `json:"cloud,omitempty"`
	Config                 *ClusterConfigEntity      `json:"config,omitempty"`
	DisplayName            *string                   `json:"display_name,omitempty"`
	Environment            *ClusterEnvironmentEntity `json:"environment,omitempty"`
	HTTPEndpoint           *string                   `json:"http_endpoint,omitempty"`
	KafkaBootstrapEndpoint *string                   `json:"kafka_bootstrap_endpoint,omitempty"`
	Network                *ClusterNetworkEntity     `json:"network,omitempty"`
	Region                 *string                   `json:"region,omitempty"`
	Zone                   *string                   `json:"zone,omitempty"`
}
