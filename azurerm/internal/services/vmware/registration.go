package vmware

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Vmware"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"VMware (AVS)",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_vmware_private_cloud": dataSourceVmwarePrivateCloud(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_vmware_private_cloud": resourceVmwarePrivateCloud(),
	}
}
