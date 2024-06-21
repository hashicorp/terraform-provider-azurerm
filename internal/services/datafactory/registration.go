// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var _ sdk.TypedServiceRegistrationWithAGitHubLabel = Registration{}
var _ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/data-factory"
}

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

func (Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		TriggerScheduleDataSource{},
		TriggerSchedulesDataSource{},
	}
}

func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		DataFactoryDatasetAzureSQLTableResource{},
		DataFactoryCredentialServicePrincipalResource{},
		DataFactoryCredentialUserAssignedManagedIdentityResource{},
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_factory": dataSourceDataFactory(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	resources := map[string]*pluginsdk.Resource{
		"azurerm_data_factory":                                       resourceDataFactory(),
		"azurerm_data_factory_data_flow":                             resourceDataFactoryDataFlow(),
		"azurerm_data_factory_flowlet_data_flow":                     resourceDataFactoryFlowletDataFlow(),
		"azurerm_data_factory_dataset_azure_blob":                    resourceDataFactoryDatasetAzureBlob(),
		"azurerm_data_factory_dataset_binary":                        resourceDataFactoryDatasetBinary(),
		"azurerm_data_factory_dataset_cosmosdb_sqlapi":               resourceDataFactoryDatasetCosmosDbSQLAPI(),
		"azurerm_data_factory_dataset_delimited_text":                resourceDataFactoryDatasetDelimitedText(),
		"azurerm_data_factory_dataset_http":                          resourceDataFactoryDatasetHTTP(),
		"azurerm_data_factory_dataset_json":                          resourceDataFactoryDatasetJSON(),
		"azurerm_data_factory_dataset_mysql":                         resourceDataFactoryDatasetMySQL(),
		"azurerm_data_factory_dataset_parquet":                       resourceDataFactoryDatasetParquet(),
		"azurerm_data_factory_dataset_postgresql":                    resourceDataFactoryDatasetPostgreSQL(),
		"azurerm_data_factory_dataset_snowflake":                     resourceDataFactoryDatasetSnowflake(),
		"azurerm_data_factory_dataset_sql_server_table":              resourceDataFactoryDatasetSQLServerTable(),
		"azurerm_data_factory_custom_dataset":                        resourceDataFactoryCustomDataset(),
		"azurerm_data_factory_integration_runtime_azure":             resourceDataFactoryIntegrationRuntimeAzure(),
		"azurerm_data_factory_integration_runtime_azure_ssis":        resourceDataFactoryIntegrationRuntimeAzureSsis(),
		"azurerm_data_factory_integration_runtime_self_hosted":       resourceDataFactoryIntegrationRuntimeSelfHosted(),
		"azurerm_data_factory_linked_custom_service":                 resourceDataFactoryLinkedCustomService(),
		"azurerm_data_factory_linked_service_azure_blob_storage":     resourceDataFactoryLinkedServiceAzureBlobStorage(),
		"azurerm_data_factory_linked_service_azure_databricks":       resourceDataFactoryLinkedServiceAzureDatabricks(),
		"azurerm_data_factory_linked_service_azure_file_storage":     resourceDataFactoryLinkedServiceAzureFileStorage(),
		"azurerm_data_factory_linked_service_azure_function":         resourceDataFactoryLinkedServiceAzureFunction(),
		"azurerm_data_factory_linked_service_azure_search":           resourceDataFactoryLinkedServiceAzureSearch(),
		"azurerm_data_factory_linked_service_azure_sql_database":     resourceDataFactoryLinkedServiceAzureSQLDatabase(),
		"azurerm_data_factory_linked_service_azure_table_storage":    resourceDataFactoryLinkedServiceAzureTableStorage(),
		"azurerm_data_factory_linked_service_cosmosdb":               resourceDataFactoryLinkedServiceCosmosDb(),
		"azurerm_data_factory_linked_service_cosmosdb_mongoapi":      resourceDataFactoryLinkedServiceCosmosDbMongoAPI(),
		"azurerm_data_factory_linked_service_data_lake_storage_gen2": resourceDataFactoryLinkedServiceDataLakeStorageGen2(),
		"azurerm_data_factory_linked_service_key_vault":              resourceDataFactoryLinkedServiceKeyVault(),
		"azurerm_data_factory_linked_service_kusto":                  resourceDataFactoryLinkedServiceKusto(),
		"azurerm_data_factory_linked_service_mysql":                  resourceDataFactoryLinkedServiceMySQL(),
		"azurerm_data_factory_linked_service_odata":                  resourceArmDataFactoryLinkedServiceOData(),
		"azurerm_data_factory_linked_service_odbc":                   resourceDataFactoryLinkedServiceOdbc(),
		"azurerm_data_factory_linked_service_postgresql":             resourceDataFactoryLinkedServicePostgreSQL(),
		"azurerm_data_factory_linked_service_sftp":                   resourceDataFactoryLinkedServiceSFTP(),
		"azurerm_data_factory_linked_service_snowflake":              resourceDataFactoryLinkedServiceSnowflake(),
		"azurerm_data_factory_linked_service_sql_server":             resourceDataFactoryLinkedServiceSQLServer(),
		"azurerm_data_factory_linked_service_synapse":                resourceDataFactoryLinkedServiceSynapse(),
		"azurerm_data_factory_linked_service_web":                    resourceDataFactoryLinkedServiceWeb(),
		"azurerm_data_factory_managed_private_endpoint":              resourceDataFactoryManagedPrivateEndpoint(),
		"azurerm_data_factory_pipeline":                              resourceDataFactoryPipeline(),
		"azurerm_data_factory_trigger_blob_event":                    resourceDataFactoryTriggerBlobEvent(),
		"azurerm_data_factory_trigger_custom_event":                  resourceDataFactoryTriggerCustomEvent(),
		"azurerm_data_factory_trigger_schedule":                      resourceDataFactoryTriggerSchedule(),
		"azurerm_data_factory_trigger_tumbling_window":               resourceDataFactoryTriggerTumblingWindow(),
	}

	if !features.FourPointOhBeta() {
		resources["azurerm_data_factory_integration_runtime_managed"] = resourceDataFactoryIntegrationRuntimeManaged()
	}

	return resources
}
