---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_project_connection"
description: |-
  Manages a Cognitive Account Project Connection.
---

# azurerm_cognitive_account_project_connection

Manages a Cognitive Account Project Connection.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                       = "example-account"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "example-account-subdomain"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_project" "example" {
  name                 = "example-project"
  cognitive_account_id = azurerm_cognitive_account.example.id
  location             = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account" "openai" {
  name                = "example-openai"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "OpenAI"
  sku_name            = "S0"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_project_connection" "example" {
  name                 = "example-connection"
  cognitive_project_id = azurerm_cognitive_account_project.example.id
  auth_type            = "ApiKey"
  category             = "AzureOpenAI"
  target               = azurerm_cognitive_account.openai.endpoint
  api_key              = azurerm_cognitive_account.openai.primary_access_key

  metadata = {
    apiType    = "Azure"
    resourceId = azurerm_cognitive_account.openai.id
    location   = azurerm_cognitive_account.openai.location
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Account Project Connection. Changing this forces a new resource to be created.

* `cognitive_project_id` - (Required) The ID of the Cognitive Account Project where the Connection should exist. Changing this forces a new resource to be created.

* `auth_type` - (Required) The authentication type for the connection. Possible values are `AAD`, `ApiKey`, `CustomKeys`, and `OAuth2`.

* `category` - (Required) The category of the connection. Possible values are `ADLSGen2`, `AIServices`, `AmazonMws`, `AmazonRdsForOracle`, `AmazonRdsForSqlServer`, `AmazonRedshift`, `AmazonS3Compatible`, `ApiKey`, `AzureBlob`, `AzureDataExplorer`, `AzureDatabricksDeltaLake`, `AzureMariaDb`, `AzureMySqlDb`, `AzureOneLake`, `AzureOpenAI`, `AzurePostgresDb`, `AzureSqlDb`, `AzureSqlMi`, `AzureSynapseAnalytics`, `AzureTableStorage`, `BingLLMSearch`, `Cassandra`, `CognitiveSearch`, `CognitiveService`, `Concur`, `ContainerRegistry`, `CosmosDb`, `CosmosDbMongoDbApi`, `Couchbase`, `CustomKeys`, `Db2`, `Drill`, `Dynamics`, `DynamicsAx`, `DynamicsCrm`, `Elasticsearch`, `Eloqua`, `FileServer`, `FtpServer`, `GenericContainerRegistry`, `GenericHttp`, `GenericRest`, `Git`, `GoogleAdWords`, `GoogleBigQuery`, `GoogleCloudStorage`, `Greenplum`, `Hbase`, `Hdfs`, `Hive`, `Hubspot`, `Impala`, `Informix`, `Jira`, `Magento`, `ManagedOnlineEndpoint`, `MariaDb`, `Marketo`, `MicrosoftAccess`, `MongoDbAtlas`, `MongoDbV2`, `MySql`, `Netezza`, `ODataRest`, `Odbc`, `Office365`, `OpenAI`, `Oracle`, `OracleCloudStorage`, `OracleServiceCloud`, `PayPal`, `Phoenix`, `Pinecone`, `PostgreSql`, `Presto`, `PythonFeed`, `QuickBooks`, `Redis`, `Responsys`, `S3`, `Salesforce`, `SalesforceMarketingCloud`, `SalesforceServiceCloud`, `SapBw`, `SapCloudForCustomer`, `SapEcc`, `SapHana`, `SapOpenHub`, `SapTable`, `Serp`, `Serverless`, `ServiceNow`, `Sftp`, `SharePointOnlineList`, `Shopify`, `Snowflake`, `Spark`, `SqlServer`, `Square`, `Sybase`, `Teradata`, `Vertica`, `WebTable`, `Xero`, and `Zoho`. Changing this forces a new resource to be created.

* `metadata` - (Required) A mapping of metadata key-value pairs for the connection. The required keys depend on the `category` specified.

* `target` - (Required) The target endpoint URL for the connection.

* `api_key` - (Optional) The API key for authentication. This field is sensitive.

~> **Note:** `api_key` is required when `auth_type` is set to `ApiKey`. `api_key` cannot be set together with `oauth2` or `custom_keys`.

* `custom_keys` - (Optional) A mapping of custom keys for authentication. All values in this map are sensitive.

~> **Note:** `custom_keys` is required when `auth_type` is set to `CustomKeys`. `custom_keys` cannot be set together with `api_key` or `oauth2`.

* `oauth2` - (Optional) An `oauth2` block as defined below.

~> **Note:** `oauth2` is required when `auth_type` is set to `OAuth2`. `oauth2` cannot be set together with `api_key` or `custom_keys`.

---

An `oauth2` block supports the following:

* `auth_url` - (Required) The OAuth2 authorization URL.

* `client_id` - (Optional) The OAuth2 client ID.

* `client_secret` - (Optional) The OAuth2 client secret. This field is sensitive.

* `developer_token` - (Optional) The OAuth2 developer token. This field is sensitive.

* `password` - (Optional) The OAuth2 password. This field is sensitive.

* `refresh_token` - (Optional) The OAuth2 refresh token. This field is sensitive.

* `tenant_id` - (Optional) The OAuth2 tenant ID.

* `username` - (Optional) The OAuth2 username.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Account Project Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Account Project Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Account Project Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Account Project Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Account Project Connection.

## Import

Cognitive Account Project Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_project_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1/projects/project1/connections/connection1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CognitiveServices` - 2025-06-01
