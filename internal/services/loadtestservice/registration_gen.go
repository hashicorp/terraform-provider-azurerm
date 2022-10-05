package loadtestservice

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.TypedServiceRegistration = autoRegistration{}

type autoRegistration struct {
}

func (a autoRegistration) Name() string {
	return "Load Test"
}

func (a autoRegistration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (a autoRegistration) Resources() []sdk.Resource {
	return []sdk.Resource{
		LoadTestResource{},
	}
}

func (a autoRegistration) WebsiteCategories() []string {
	return []string{
		"Load Test",
	}
}
