package fluidrelay

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

type Registration struct{}

func (r Registration) Name() string {
	return "Fluid Relay"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		Server{},
	}
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Fluid Relay",
	}
}

var (
	_ sdk.TypedServiceRegistration = (*Registration)(nil)
)
