package batch

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_batch_account":     dataSourceArmBatchAccount(),
		"azurerm_batch_certificate": dataSourceArmBatchCertificate(),
		"azurerm_batch_pool":        dataSourceArmBatchPool(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_batch_account":     resourceArmBatchAccount(),
		"azurerm_batch_application": resourceArmBatchApplication(),
		"azurerm_batch_certificate": resourceArmBatchCertificate(),
		"azurerm_batch_pool":        resourceArmBatchPool(),
	}
}
