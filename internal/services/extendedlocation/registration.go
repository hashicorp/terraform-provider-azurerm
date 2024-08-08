package extendedlocation

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/extendedlocation"
}

func (Registration) Name() string {
	return "ExtendedLocation"
}

func (Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		CustomLocationDataSource{},
	}
}

func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{}
}

func (Registration) WebsiteCategories() []string {
	return []string{
		"Extended Location",
	}
}
