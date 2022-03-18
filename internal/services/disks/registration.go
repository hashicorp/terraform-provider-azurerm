package disks

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}

type Registration struct{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/disks"
}

func (r Registration) Name() string {
	return "Disks"
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{}
}

func (r Registration) Resources() []sdk.Resource {
	resources := []sdk.Resource{
		DiskPoolResource{},
		DiskPoolManagedDiskAttachmentResource{},
		DisksPoolIscsiTargetResource{},
		DiskPoolIscsiTargetLunModel{},
	}

	if !features.ThreePointOhBeta() {
		resources = append(resources, StorageDisksPoolResource{})
	}

	return resources
}

func (r Registration) WebsiteCategories() []string {
	return []string{
		"Disks",
		"Storage",
	}
}
