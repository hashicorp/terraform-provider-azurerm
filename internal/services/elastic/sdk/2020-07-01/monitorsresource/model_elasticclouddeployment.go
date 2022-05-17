package monitorsresource

type ElasticCloudDeployment struct {
	AzureSubscriptionId     *string `json:"azureSubscriptionId,omitempty"`
	DeploymentId            *string `json:"deploymentId,omitempty"`
	ElasticsearchRegion     *string `json:"elasticsearchRegion,omitempty"`
	ElasticsearchServiceUrl *string `json:"elasticsearchServiceUrl,omitempty"`
	KibanaServiceUrl        *string `json:"kibanaServiceUrl,omitempty"`
	KibanaSsoUrl            *string `json:"kibanaSsoUrl,omitempty"`
	Name                    *string `json:"name,omitempty"`
}
