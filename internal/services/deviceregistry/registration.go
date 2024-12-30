package deviceregistry

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistration = Registration{}
)

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		AssetDataSource{},
		AssetEndpointProfileDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		AssetResource{},
		AssetEndpointProfileResource{},
	}
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Device Registry"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Device Registry",
	}
}
