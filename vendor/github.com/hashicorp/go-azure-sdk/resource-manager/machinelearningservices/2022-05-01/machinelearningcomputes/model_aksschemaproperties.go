package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AKSSchemaProperties struct {
	AgentCount                 *int64                      `json:"agentCount,omitempty"`
	AgentVMSize                *string                     `json:"agentVmSize,omitempty"`
	AksNetworkingConfiguration *AksNetworkingConfiguration `json:"aksNetworkingConfiguration,omitempty"`
	ClusterFqdn                *string                     `json:"clusterFqdn,omitempty"`
	ClusterPurpose             *ClusterPurpose             `json:"clusterPurpose,omitempty"`
	LoadBalancerSubnet         *string                     `json:"loadBalancerSubnet,omitempty"`
	LoadBalancerType           *LoadBalancerType           `json:"loadBalancerType,omitempty"`
	SslConfiguration           *SslConfiguration           `json:"sslConfiguration,omitempty"`
	SystemServices             *[]SystemService            `json:"systemServices,omitempty"`
}
