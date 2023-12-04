package systemcentervirtualmachinemanager

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/systemcentervirtualmachinemanager"
}

// Name is the name of this Service
func (r Registration) Name() string {
	return "System Center Virtual Machine Manager"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"System Center Virtual Machine Manager",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		SystemCenterVirtualMachineManagerVirtualMachineInstanceGuestAgentResource{},
	}
}
