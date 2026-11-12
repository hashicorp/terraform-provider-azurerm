// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datafactory

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_credential_service_principal"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_credential_user_managed_identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_custom_dataset"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_customer_managed_key"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_data_flow"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_azure_blob"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_azure_sql_table"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_binary"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_cosmosdb_sqlapi"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_delimited_text"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_http"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_json"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_mysql"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_parquet"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_postgresql"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_snowflake"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_dataset_sql_server_table"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_flowlet_data_flow"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_integration_runtime_azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_integration_runtime_azure_ssis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_integration_runtime_self_hosted"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_custom_service"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_azure_blob_storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_azure_databricks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_azure_file_storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_azure_function"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_azure_search"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_azure_sql_database"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_azure_table_storage"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_cosmosdb"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_cosmosdb_mongoapi"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_data_lake_storage_gen2"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_key_vault"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_kusto"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_mysql"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_odata"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_odbc"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_postgresql"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_sftp"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_snowflake"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_sql_managed_instance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_sql_server"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_synapse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_linked_service_web"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_managed_private_endpoint"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_pipeline"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_trigger_blob_event"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_trigger_custom_event"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_trigger_schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_trigger_schedules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/resources/data_factory_trigger_tumbling_window"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.FrameworkServiceRegistration               = Registration{}
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
)

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
		data_factory_trigger_schedule.DataSource{},
		data_factory_trigger_schedules.DataSource{},
	}
}

func (Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		data_factory_credential_service_principal.Resource{},
		data_factory_credential_user_managed_identity.Resource{},
		data_factory_customer_managed_key.Resource{},
		data_factory_dataset_azure_sql_table.Resource{},
		data_factory_linked_service_sql_managed_instance.Resource{},
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_factory": data_factory.RegisterDataSource(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_data_factory":                                       data_factory.RegisterResource(),
		"azurerm_data_factory_custom_dataset":                        data_factory_custom_dataset.RegisterResource(),
		"azurerm_data_factory_data_flow":                             data_factory_data_flow.RegisterResource(),
		"azurerm_data_factory_dataset_azure_blob":                    data_factory_dataset_azure_blob.RegisterResource(),
		"azurerm_data_factory_dataset_binary":                        data_factory_dataset_binary.RegisterResource(),
		"azurerm_data_factory_dataset_cosmosdb_sqlapi":               data_factory_dataset_cosmosdb_sqlapi.RegisterResource(),
		"azurerm_data_factory_dataset_delimited_text":                data_factory_dataset_delimited_text.RegisterResource(),
		"azurerm_data_factory_dataset_http":                          data_factory_dataset_http.RegisterResource(),
		"azurerm_data_factory_dataset_json":                          data_factory_dataset_json.RegisterResource(),
		"azurerm_data_factory_dataset_mysql":                         data_factory_dataset_mysql.RegisterResource(),
		"azurerm_data_factory_dataset_parquet":                       data_factory_dataset_parquet.RegisterResource(),
		"azurerm_data_factory_dataset_postgresql":                    data_factory_dataset_postgresql.RegisterResource(),
		"azurerm_data_factory_dataset_snowflake":                     data_factory_dataset_snowflake.RegisterResource(),
		"azurerm_data_factory_dataset_sql_server_table":              data_factory_dataset_sql_server_table.RegisterResource(),
		"azurerm_data_factory_flowlet_data_flow":                     data_factory_flowlet_data_flow.RegisterResource(),
		"azurerm_data_factory_integration_runtime_azure":             data_factory_integration_runtime_azure.RegisterResource(),
		"azurerm_data_factory_integration_runtime_azure_ssis":        data_factory_integration_runtime_azure_ssis.RegisterResource(),
		"azurerm_data_factory_integration_runtime_self_hosted":       data_factory_integration_runtime_self_hosted.RegisterResource(),
		"azurerm_data_factory_linked_custom_service":                 data_factory_linked_custom_service.RegisterResource(),
		"azurerm_data_factory_linked_service_azure_blob_storage":     data_factory_linked_service_azure_blob_storage.RegisterResource(),
		"azurerm_data_factory_linked_service_azure_databricks":       data_factory_linked_service_azure_databricks.RegisterResource(),
		"azurerm_data_factory_linked_service_azure_file_storage":     data_factory_linked_service_azure_file_storage.RegisterResource(),
		"azurerm_data_factory_linked_service_azure_function":         data_factory_linked_service_azure_function.RegisterResource(),
		"azurerm_data_factory_linked_service_azure_search":           data_factory_linked_service_azure_search.RegisterResource(),
		"azurerm_data_factory_linked_service_azure_sql_database":     data_factory_linked_service_azure_sql_database.RegisterResource(),
		"azurerm_data_factory_linked_service_azure_table_storage":    data_factory_linked_service_azure_table_storage.RegisterResource(),
		"azurerm_data_factory_linked_service_cosmosdb":               data_factory_linked_service_cosmosdb.RegisterResource(),
		"azurerm_data_factory_linked_service_cosmosdb_mongoapi":      data_factory_linked_service_cosmosdb_mongoapi.RegisterResource(),
		"azurerm_data_factory_linked_service_data_lake_storage_gen2": data_factory_linked_service_data_lake_storage_gen2.RegisterResource(),
		"azurerm_data_factory_linked_service_key_vault":              data_factory_linked_service_key_vault.RegisterResource(),
		"azurerm_data_factory_linked_service_kusto":                  data_factory_linked_service_kusto.RegisterResource(),
		"azurerm_data_factory_linked_service_mysql":                  data_factory_linked_service_mysql.RegisterResource(),
		"azurerm_data_factory_linked_service_odata":                  data_factory_linked_service_odata.RegisterResource(),
		"azurerm_data_factory_linked_service_odbc":                   data_factory_linked_service_odbc.RegisterResource(),
		"azurerm_data_factory_linked_service_postgresql":             data_factory_linked_service_postgresql.RegisterResource(),
		"azurerm_data_factory_linked_service_sftp":                   data_factory_linked_service_sftp.RegisterResource(),
		"azurerm_data_factory_linked_service_snowflake":              data_factory_linked_service_snowflake.RegisterResource(),
		"azurerm_data_factory_linked_service_sql_server":             data_factory_linked_service_sql_server.RegisterResource(),
		"azurerm_data_factory_linked_service_synapse":                data_factory_linked_service_synapse.RegisterResource(),
		"azurerm_data_factory_linked_service_web":                    data_factory_linked_service_web.RegisterResource(),
		"azurerm_data_factory_managed_private_endpoint":              data_factory_managed_private_endpoint.RegisterResource(),
		"azurerm_data_factory_pipeline":                              data_factory_pipeline.RegisterResource(),
		"azurerm_data_factory_trigger_blob_event":                    data_factory_trigger_blob_event.RegisterResource(),
		"azurerm_data_factory_trigger_custom_event":                  data_factory_trigger_custom_event.RegisterResource(),
		"azurerm_data_factory_trigger_schedule":                      data_factory_trigger_schedule.RegisterResource(),
		"azurerm_data_factory_trigger_tumbling_window":               data_factory_trigger_tumbling_window.RegisterResource(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{}
}
