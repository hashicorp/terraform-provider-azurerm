package batch

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Batch"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Batch",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_batch_account":     dataSourceBatchAccount(),
		"azurerm_batch_certificate": dataSourceBatchCertificate(),
		"azurerm_batch_pool":        dataSourceBatchPool(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_batch_account":     resourceBatchAccount(),
		"azurerm_batch_application": resourceBatchApplication(),
		"azurerm_batch_certificate": resourceBatchCertificate(),
		"azurerm_batch_pool":        resourceBatchPool(),
	}
}
