package devcenter

// NOTE: this file is temporary to enable the project to compile prior to auto-generation

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.TypedServiceRegistration = autoRegistration{}

type autoRegistration struct {
}

func (autoRegistration) Name() string {
	return "Dev Center"
}

func (autoRegistration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (autoRegistration) Resources() []sdk.Resource {
	return []sdk.Resource{}
}

func (autoRegistration) WebsiteCategories() []string {
	return []string{
		"Dev Center",
	}
}
