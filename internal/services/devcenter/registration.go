package devcenter

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/dev-center"
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Dev Center",
	}
}

func (r Registration) Name() string {
	return "Dev Center"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		DevCenterResource{},
	}
}

var _ sdk.TypedServiceRegistration = Registration{}
