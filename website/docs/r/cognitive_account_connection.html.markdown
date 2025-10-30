---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_connection"
description: |-
  Manages a Cognitive Services Account Connection.
---

# azurerm_cognitive_account_connection

Manages a Cognitive Services Account Connection.

## Example Usage

### AAD Authentication

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                       = "example-aiservices"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "exampleaiservices"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "examplecsc"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_cognitive_account_connection" "example" {
  name                 = "example-connection"
  cognitive_account_id = azurerm_cognitive_account.example.id
  auth_type            = "AAD"
  category             = "AzureBlob"
  target               = azurerm_storage_account.example.primary_blob_endpoint

  metadata = {
    accountName   = azurerm_storage_account.example.name
    containerName = azurerm_storage_container.example.name
  }
}
```

### API Key Authentication

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                       = "example-aiservices"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "exampleaiservices"

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

resource "azurerm_cognitive_account_connection" "example" {
  name                 = "example-connection"
  cognitive_account_id = azurerm_cognitive_account.example.id
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

### OAuth2 Authentication

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                       = "example-aiservices"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "exampleaiservices"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "examplecsc"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_cognitive_account_connection" "example" {
  name                 = "example-connection"
  cognitive_account_id = azurerm_cognitive_account.example.id
  auth_type            = "OAuth2"
  category             = "AzureBlob"
  target               = azurerm_storage_account.example.primary_blob_endpoint

  metadata = {
    containerName = azurerm_storage_container.example.name
    accountName   = azurerm_storage_account.example.name
  }

  oauth2 {
    auth_url        = "https://login.microsoftonline.com/00000000-0000-0000-0000-000000000000/oauth2/v2.0/token"
    client_id       = "00000000-0000-0000-0000-000000000000"
    client_secret   = "placeHolderClientSecret"
    tenant_id       = "00000000-0000-0000-0000-000000000000"
    developer_token = "placeHolderDevToken"
    refresh_token   = "placeRefreshToken"
    username        = "placeHolderUsername"
    password        = "placeHolderPassword"
  }
}
```

### Managed Identity Authentication

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-identity"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_cognitive_account" "example" {
  name                       = "example-aiservices"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "exampleaiservices"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "examplecsc"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_cognitive_account_connection" "example" {
  name                 = "example-connection"
  cognitive_account_id = azurerm_cognitive_account.example.id
  auth_type            = "ManagedIdentity"
  category             = "AzureBlob"
  target               = azurerm_storage_account.example.primary_blob_endpoint

  metadata = {
    containerName = azurerm_storage_container.example.name
    accountName   = azurerm_storage_account.example.name
  }

  managed_identity {
    client_id   = azurerm_user_assigned_identity.example.client_id
    resource_id = azurerm_user_assigned_identity.example.id
  }
}
```

### Custom Keys Authentication

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                       = "example-aiservices"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "exampleaiservices"

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

resource "azurerm_cognitive_account_connection" "example" {
  name                 = "example-connection"
  cognitive_account_id = azurerm_cognitive_account.example.id
  auth_type            = "CustomKeys"
  category             = "CustomKeys"
  target               = azurerm_cognitive_account.openai.endpoint

  metadata = {
    apiType    = "Azure"
    resourceId = azurerm_cognitive_account.openai.id
    location   = azurerm_cognitive_account.openai.location
  }

  custom_keys = {
    primaryKey   = azurerm_cognitive_account.openai.primary_access_key
    secondaryKey = azurerm_cognitive_account.openai.secondary_access_key
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Services Account Connection. Changing this forces a new resource to be created.

* `cognitive_account_id` - (Required) The ID of the Cognitive Services Account. Changing this forces a new resource to be created.

* `auth_type` - (Required) The authentication type for the connection. Possible values are `AAD`, `ApiKey`, `CustomKeys`, `ManagedIdentity`, and `OAuth2`.

* `category` - (Required) The category of the connection. Possible values are `ADLSGen2`, `AIServices`, `AmazonMws`, `AmazonRdsForOracle`, `AmazonRdsForSqlServer`, `AmazonRedshift`, `AmazonS3Compatible`, `ApiKey`, `AzureBlob`, `AzureDataExplorer`, `AzureDatabricksDeltaLake`, `AzureMariaDb`, `AzureMySqlDb`, `AzureOneLake`, `AzureOpenAI`, `AzurePostgresDb`, `AzureSqlDb`, `AzureSqlMi`, `AzureSynapseAnalytics`, `AzureTableStorage`, `BingLLMSearch`, `Cassandra`, `CognitiveSearch`, `CognitiveService`, `Concur`, `ContainerRegistry`, `CosmosDb`, `CosmosDbMongoDbApi`, `Couchbase`, `CustomKeys`, `Db2`, `Drill`, `Dynamics`, `DynamicsAx`, `DynamicsCrm`, `Elasticsearch`, `Eloqua`, `FileServer`, `FtpServer`, `GenericContainerRegistry`, `GenericHttp`, `GenericRest`, `Git`, `GoogleAdWords`, `GoogleBigQuery`, `GoogleCloudStorage`, `Greenplum`, `Hbase`, `Hdfs`, `Hive`, `Hubspot`, `Impala`, `Informix`, `Jira`, `Magento`, `ManagedOnlineEndpoint`, `MariaDb`, `Marketo`, `MicrosoftAccess`, `MongoDbAtlas`, `MongoDbV2`, `MySql`, `Netezza`, `ODataRest`, `Odbc`, `Office365`, `OpenAI`, `Oracle`, `OracleCloudStorage`, `OracleServiceCloud`, `PayPal`, `Phoenix`, `Pinecone`, `PostgreSql`, `Presto`, `PythonFeed`, `QuickBooks`, `Redis`, `Responsys`, `S3`, `Salesforce`, `SalesforceMarketingCloud`, `SalesforceServiceCloud`, `SapBw`, `SapCloudForCustomer`, `SapEcc`, `SapHana`, `SapOpenHub`, `SapTable`, `Serp`, `Serverless`, `ServiceNow`, `Sftp`, `SharePointOnlineList`, `Shopify`, `Snowflake`, `Spark`, `SqlServer`, `Square`, `Sybase`, `Teradata`, `Vertica`, `WebTable`, `Xero`, and `Zoho`. Changing this forces a new resource to be created.

* `metadata` - (Required) A mapping of metadata key-value pairs for the connection.

* `target` - (Required) The target endpoint or resource for the connection.

* `api_key` - (Optional) The API key for authentication. This field is sensitive. Must be set with `auth_type` of `ApiKey`. Cannot be used with `managed_identity`, `oauth2`, or `custom_keys`.

* `managed_identity` - (Optional) A `managed_identity` block as defined below. Must be set with `auth_type` of `ManagedIdentity`. Cannot be used with `api_key`, `oauth2`, or `custom_keys`.

* `oauth2` - (Optional) An `oauth2` block as defined below. Must be set with `auth_type` of `OAuth2`. Cannot be used with `api_key`, `managed_identity`, or `custom_keys`.

* `custom_keys` - (Optional) A mapping of custom keys for authentication. Must be set with `auth_type` of `CustomKeys`. All values in this map are sensitive. Cannot be used with `api_key`, `managed_identity`, or `oauth2`.

---

A `managed_identity` block supports the following:

* `client_id` - (Required) The client ID of the managed identity.

* `resource_id` - (Required) The resource ID of the managed identity.

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

* `id` - The ID of the Cognitive Services Account Connection.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Services Account Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Services Account Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Services Account Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Services Account Connection.

## Import

Cognitive Services Account Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_connection.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1/connections/connection1
```
