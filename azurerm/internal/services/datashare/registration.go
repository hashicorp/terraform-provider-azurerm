package datashare

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Data Share"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Data Share",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_data_share_account":              dataSourceDataShareAccount(),
		"azurerm_data_share":                      dataSourceDataShare(),
		"azurerm_data_share_dataset_blob_storage": dataSourceDataShareDatasetBlobStorage(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_data_share_account":              resourceArmDataShareAccount(),
		"azurerm_data_share":                      resourceArmDataShare(),
		"azurerm_data_share_dataset_blob_storage": resourceArmDataShareDataSetBlobStorage(),
	}
}
