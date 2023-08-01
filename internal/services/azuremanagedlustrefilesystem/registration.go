package azuremanagedlustrefilesystem

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/azuremanagedlustrefilesystem"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Azure Managed Lustre File System"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Azure Managed Lustre File System",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		ManagedLustreFileSystemResource{},
	}
}
