package monitorsresource

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ElasticCloudDeployment struct {
	AzureSubscriptionId     *string `json:"azureSubscriptionId,omitempty"`
	DeploymentId            *string `json:"deploymentId,omitempty"`
	ElasticsearchRegion     *string `json:"elasticsearchRegion,omitempty"`
	ElasticsearchServiceUrl *string `json:"elasticsearchServiceUrl,omitempty"`
	KibanaServiceUrl        *string `json:"kibanaServiceUrl,omitempty"`
	KibanaSsoUrl            *string `json:"kibanaSsoUrl,omitempty"`
	Name                    *string `json:"name,omitempty"`
}
