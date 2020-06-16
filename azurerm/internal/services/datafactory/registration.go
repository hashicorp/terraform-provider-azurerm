package datafactory

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Data Factory"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Data Factory",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_data_factory": dataSourceArmDataFactory(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_data_factory":                                       resourceArmDataFactory(),
		"azurerm_data_factory_dataset_mysql":                         resourceArmDataFactoryDatasetMySQL(),
		"azurerm_data_factory_dataset_postgresql":                    resourceArmDataFactoryDatasetPostgreSQL(),
		"azurerm_data_factory_dataset_sql_server_table":              resourceArmDataFactoryDatasetSQLServerTable(),
		"azurerm_data_factory_integration_runtime_managed":           resourceArmDataFactoryIntegrationRuntimeManaged(),
		"azurerm_data_factory_integration_runtime_self_hosted":       resourceArmDataFactoryIntegrationRuntimeSelfHosted(),
		"azurerm_data_factory_linked_service_data_lake_storage_gen2": resourceArmDataFactoryLinkedServiceDataLakeStorageGen2(),
		"azurerm_data_factory_linked_service_key_vault":              resourceArmDataFactoryLinkedServiceKeyVault(),
		"azurerm_data_factory_linked_service_mysql":                  resourceArmDataFactoryLinkedServiceMySQL(),
		"azurerm_data_factory_linked_service_postgresql":             resourceArmDataFactoryLinkedServicePostgreSQL(),
		"azurerm_data_factory_linked_service_sql_server":             resourceArmDataFactoryLinkedServiceSQLServer(),
		"azurerm_data_factory_pipeline":                              resourceArmDataFactoryPipeline(),
		"azurerm_data_factory_trigger_schedule":                      resourceArmDataFactoryTriggerSchedule(),
	}
}
