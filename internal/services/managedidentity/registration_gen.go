package managedidentity

import "github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

var _ sdk.TypedServiceRegistration = autoRegistration{}

type autoRegistration struct {
}

func (a autoRegistration) Name() string {
	return "Managed Service Identities"
}

func (a autoRegistration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (a autoRegistration) Resources() []sdk.Resource {
	return []sdk.Resource{
		UserAssignedIdentityResource{},
	}
}

func (a autoRegistration) WebsiteCategories() []string {
	return []string{
		"Authorization",
	}
}
