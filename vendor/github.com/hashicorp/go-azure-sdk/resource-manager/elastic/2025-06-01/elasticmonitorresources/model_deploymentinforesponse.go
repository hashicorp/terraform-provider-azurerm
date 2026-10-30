package elasticmonitorresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentInfoResponse struct {
	ConfigurationType     *string                  `json:"configurationType,omitempty"`
	DeploymentURL         *string                  `json:"deploymentUrl,omitempty"`
	DiskCapacity          *string                  `json:"diskCapacity,omitempty"`
	ElasticsearchEndPoint *string                  `json:"elasticsearchEndPoint,omitempty"`
	MarketplaceSaasInfo   *MarketplaceSaaSInfo     `json:"marketplaceSaasInfo,omitempty"`
	MemoryCapacity        *string                  `json:"memoryCapacity,omitempty"`
	ProjectType           *string                  `json:"projectType,omitempty"`
	Status                *ElasticDeploymentStatus `json:"status,omitempty"`
	Version               *string                  `json:"version,omitempty"`
}
