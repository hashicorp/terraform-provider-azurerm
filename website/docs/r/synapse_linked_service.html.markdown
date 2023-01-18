---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_synapse_linked_service"
description: |-
  Manages a Linked Service (connection) between a resource and Azure Synapse. This is a generic resource that supports all different Linked Service Types.
---

# azurerm_synapse_linked_service

Manages a Synapse Linked Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_kind             = "BlobStorage"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "example"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_synapse_workspace" "example" {
  name                                 = "example"
  resource_group_name                  = azurerm_resource_group.example.name
  location                             = azurerm_resource_group.example.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
  managed_virtual_network_enabled      = true

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_synapse_firewall_rule" "example" {
  name                 = "allowAll"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  start_ip_address     = "0.0.0.0"
  end_ip_address       = "255.255.255.255"
}

resource "azurerm_synapse_integration_runtime_azure" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  location             = azurerm_resource_group.example.location
}

resource "azurerm_synapse_linked_service" "example" {
  name                 = "example"
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  type                 = "AzureBlobStorage"
  type_properties_json = <<JSON
{
  "connectionString": "${azurerm_storage_account.example.primary_connection_string}"
}
JSON
  integration_runtime {
    name = azurerm_synapse_integration_runtime_azure.example.name
  }

  depends_on = [
    azurerm_synapse_firewall_rule.example,
  ]
}


```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Synapse Linked Service. Changing this forces a new Synapse Linked Service to be created.

* `synapse_workspace_id` - (Required) The Synapse Workspace ID in which to associate the Linked Service with. Changing this forces a new Synapse Linked Service to be created.

* `type` - (Required) The type of data stores that will be connected to Synapse. Valid Values include `AmazonMWS`, `AmazonRdsForOracle`, `AmazonRdsForSqlServer`, `AmazonRedshift`, `AmazonS3`, `AzureBatch`. Changing this forces a new resource to be created.
`AzureBlobFS`, `AzureBlobStorage`, `AzureDataExplorer`, `AzureDataLakeAnalytics`, `AzureDataLakeStore`, `AzureDatabricks`, `AzureDatabricksDeltaLake`, `AzureFileStorage`, `AzureFunction`,
`AzureKeyVault`, `AzureML`, `AzureMLService`, `AzureMariaDB`, `AzureMySql`, `AzurePostgreSql`, `AzureSqlDW`, `AzureSqlDatabase`, `AzureSqlMI`, `AzureSearch`, `AzureStorage`,
`AzureTableStorage`, `Cassandra`, `CommonDataServiceForApps`, `Concur`, `CosmosDb`, `CosmosDbMongoDbApi`, `Couchbase`, `CustomDataSource`, `Db2`, `Drill`, 
`Dynamics`, `DynamicsAX`, `DynamicsCrm`, `Eloqua`, `FileServer`, `FtpServer`, `GoogleAdWords`, `GoogleBigQuery`, `GoogleCloudStorage`, `Greenplum`, `HBase`, `HDInsight`,
`HDInsightOnDemand`, `HttpServer`, `Hdfs`, `Hive`, `Hubspot`, `Impala`, `Informix`, `Jira`, `LinkedService`, `Magento`, `MariaDB`, `Marketo`, `MicrosoftAccess`, `MongoDb`,
`MongoDbAtlas`, `MongoDbV2`, `MySql`, `Netezza`, `OData`, `Odbc`, `Office365`, `Oracle`, `OracleServiceCloud`, `Paypal`, `Phoenix`, `PostgreSql`, `Presto`, `QuickBooks`, 
`Responsys`, `RestService`, `SqlServer`, `Salesforce`, `SalesforceMarketingCloud`, `SalesforceServiceCloud`, `SapBW`, `SapCloudForCustomer`, `SapEcc`, `SapHana`, `SapOpenHub`,
`SapTable`, `ServiceNow`, `Sftp`, `SharePointOnlineList`, `Shopify`, `Snowflake`, `Spark`, `Square`, `Sybase`, `Teradata`, `Vertica`, `Web`, `Xero`, `Zoho`.

* `type_properties_json` - (Required) A JSON object that contains the properties of the Synapse Linked Service.

---

* `additional_properties` - (Optional) A map of additional properties to associate with the Synapse Linked Service.

* `annotations` - (Optional) List of tags that can be used for describing the Synapse Linked Service.

* `description` - (Optional) The description for the Synapse Linked Service.

* `integration_runtime` - (Optional) A `integration_runtime` block as defined below.

* `parameters` - (Optional) A map of parameters to associate with the Synapse Linked Service.

---

A `integration_runtime` block supports the following:

* `name` - (Required) The integration runtime reference to associate with the Synapse Linked Service.

* `parameters` - (Optional) A map of parameters to associate with the integration runtime.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Synapse Linked Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Synapse Linked Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Linked Service.
* `update` - (Defaults to 30 minutes) Used when updating the Synapse Linked Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the Synapse Linked Service.

## Import

Synapse Linked Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_synapse_linked_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Synapse/workspaces/workspace1/linkedServices/linkedservice1
```
