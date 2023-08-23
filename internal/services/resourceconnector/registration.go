package resourceconnector

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/resource-connector"
}

func (r Registration) Name() string {
	return "Resource Connector"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ResourceConnectorApplianceResource{},
	}
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Resource Connector",
	}
}
