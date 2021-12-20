package disks

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistration = Registration{}

type Registration struct{}

func (r Registration) Name() string {
	return "Disks"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		StorageDisksPoolResource{},
	}
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		// TODO: introduce "Disks"
		"Storage",
	}
}
