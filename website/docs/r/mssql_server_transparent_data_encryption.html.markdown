---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_server_transparent_data_encryption"
description: |-
  Manages the transparent data encryption configuration for a MSSQL Server
---

# azurerm_mssql_server_transparent_data_encryption

Manages the transparent data encryption configuration for a MSSQL Server

~> **NOTE:** Once transparent data encryption is enabled on a MS SQL instance, it is not possible to remove TDE. You will be able to switch between 'ServiceManaged' and 'CustomerManaged' keys, but will not be able to remove encryption. For safety when this resource is deleted, the TDE mode will automatically be set to 'ServiceManaged'. See `key_vault_uri` for more information on how to specify the key types. As SQL Server only supports a single configuration for encryption settings, this resource will replace the current encryption settings on the server. 

~> **Note:** See [documentation](https://docs.microsoft.com/en-us/azure/azure-sql/database/transparent-data-encryption-byok-overview) for important information on how handle lifecycle management of the keys to prevent data lockout. 

## Example Usage with Service Managed Key

```hcl

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "EastUs"
}

resource "azurerm_mssql_server" "example" {
  name                         = "mssqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
  minimum_tls_version          = "1.2"

  azuread_administrator {
    login_username = "AzureAD Admin"
    object_id      = "00000000-0000-0000-0000-000000000000"
  }

  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.example.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.example.primary_access_key
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }

  tags = {
    environment = "production"
  }
}

resource "azurerm_mssql_server_transparent_data_encryption" "example" {
  server_id = azurerm_mssql_server.example.id
}
```

## Example Usage with Customer Managed Key

```hcl

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "EastUs"
}

resource "azurerm_mssql_server" "example" {
  name                         = "mssqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "thisIsKat11"
  minimum_tls_version          = "1.2"

  azuread_administrator {
    login_username = "AzureAD Admin"
    object_id      = "00000000-0000-0000-0000-000000000000"
  }

  extended_auditing_policy {
    storage_endpoint                        = azurerm_storage_account.example.primary_blob_endpoint
    storage_account_access_key              = azurerm_storage_account.example.primary_access_key
    storage_account_access_key_is_secondary = true
    retention_in_days                       = 6
  }

  tags = {
    environment = "production"
  }

  identity {
    type = "SystemAssigned"
  }
}

# Create a key vault with policies for the deployer to create a key & SQL Server to wrap/unwrap/get key
resource "azurerm_key_vault" "example" {
  name                        = "example"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"

  access_policy {
    tenant_id = data.azurerm_client_config.current.tenant_id
    object_id = data.azurerm_client_config.current.object_id

    key_permissions = [
      "Get", "List", "Create", "Delete", "Update", "Recover", "Purge",
    ]
  }
  access_policy {
    tenant_id = azurerm_mssql_server.example.identity[0].tenant_id
    object_id = azurerm_mssql_server.example.identity[0].principal_id

    key_permissions = [
      "Get", "WrapKey", "UnwrapKey"
    ]
  }
}

resource "azurerm_key_vault_key" "example" {
  name         = "byok"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "unwrapKey",
    "wrapKey",
  ]

  depends_on = [
    azurerm_key_vault.example
  ]
}

resource "azurerm_mssql_server_transparent_data_encryption" "example" {
  server_id        = azurerm_mssql_server.example.id
  key_vault_key_id = azurerm_key_vault_key.example.id
}


```

## Arguments Reference

The following arguments are supported:

* `server_id` - (Required) Specifies the name of the MS SQL Server.

---

* `key_vault_key_id` - (Optional) To use customer managed keys from Azure Key Vault, provide the AKV Key ID. To use service managed keys, omit this field.

~> **NOTE:** In order to use customer managed keys, the identity of the MSSQL server must have the following permissions on the key vault: 'get', 'wrapKey' and 'unwrapKey' 

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the MSSQL encryption protector

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MSSQL.
* `read` - (Defaults to 5 minutes) Used when retrieving the MSSQL.
* `update` - (Defaults to 30 minutes) Used when updating the MSSQL.
* `delete` - (Defaults to 30 minutes) Used when deleting the MSSQL.

## Import

~> **NOTE:** This resource does not need to be imported to manage it, however the import will work. 

SQL Server Transparent Data Encryption can be imported using the resource id, e.g.

```shell
terraform import azurerm_mssql_server_transparent_data_encryption.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Sql/servers/server1/encryptionProtector/current
```

