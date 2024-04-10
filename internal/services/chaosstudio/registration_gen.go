package chaosstudio

// NOTE: this file is generated - manual changes will be overwritten.

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.TypedServiceRegistration = autoRegistration{}

type autoRegistration struct {
}

func (autoRegistration) Name() string {
	return "ChaosStudio"
}

func (autoRegistration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (autoRegistration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ChaosStudioTargetResource{},
	}
}

func (autoRegistration) WebsiteCategories() []string {
	return []string{
		"Chaos Studio",
	}
}
