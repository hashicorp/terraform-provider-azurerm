package elasticsan

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct {
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/elasticsan"
}
func (Registration) Name() string {
	return "ElasticSan"
}

func (Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ElasticSANResource{},
	}
}

func (Registration) WebsiteCategories() []string {
	return []string{
		"Elastic SAN",
	}
}
