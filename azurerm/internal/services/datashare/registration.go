package datashare

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

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
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_share_account":                dataSourceDataShareAccount(),
		"azurerm_data_share":                        dataSourceDataShare(),
		"azurerm_data_share_dataset_blob_storage":   dataSourceDataShareDatasetBlobStorage(),
		"azurerm_data_share_dataset_data_lake_gen1": dataSourceDataShareDatasetDataLakeGen1(),
		"azurerm_data_share_dataset_data_lake_gen2": dataSourceDataShareDatasetDataLakeGen2(),
		"azurerm_data_share_dataset_kusto_cluster":  dataSourceDataShareDatasetKustoCluster(),
		"azurerm_data_share_dataset_kusto_database": dataSourceDataShareDatasetKustoDatabase(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_share_account":                resourceDataShareAccount(),
		"azurerm_data_share":                        resourceDataShare(),
		"azurerm_data_share_dataset_blob_storage":   resourceDataShareDataSetBlobStorage(),
		"azurerm_data_share_dataset_data_lake_gen1": resourceDataShareDataSetDataLakeGen1(),
		"azurerm_data_share_dataset_data_lake_gen2": resourceDataShareDataSetDataLakeGen2(),
		"azurerm_data_share_dataset_kusto_cluster":  resourceDataShareDataSetKustoCluster(),
		"azurerm_data_share_dataset_kusto_database": resourceDataShareDataSetKustoDatabase(),
	}
}
