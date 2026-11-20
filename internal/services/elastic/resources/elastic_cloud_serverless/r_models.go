// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package elastic_cloud_serverless

type ElasticCloudServerlessResourceModel struct {
	Name                     string            `tfschema:"name"`
	ResourceGroupName        string            `tfschema:"resource_group_name"`
	Location                 string            `tfschema:"location"`
	ConfigurationType        string            `tfschema:"configuration_type"`
	ElasticCloudEmailAddress string            `tfschema:"elastic_cloud_email_address"`
	Kind                     string            `tfschema:"kind"`
	OfferID                  string            `tfschema:"offer_id"`
	ProjectType              string            `tfschema:"project_type"`
	Sku                      string            `tfschema:"sku"`
	TermID                   string            `tfschema:"term_id"`
	GenerateAPIKey           bool              `tfschema:"generate_api_key"`
	MonitoringEnabled        bool              `tfschema:"monitoring_enabled"`
	PlanID                   string            `tfschema:"plan_id"`
	PublisherID              string            `tfschema:"publisher_id"`
	Tags                     map[string]string `tfschema:"tags"`
	ElasticCloudDeploymentID string            `tfschema:"elastic_cloud_deployment_id"`
	ElasticsearchServiceURL  string            `tfschema:"elasticsearch_service_url"`
	KibanaServiceURL         string            `tfschema:"kibana_service_url"`
	KibanaSSOURI             string            `tfschema:"kibana_sso_uri"`
}

func (r Resource) ModelObject() any {
	return &ElasticCloudServerlessResourceModel{}
}
