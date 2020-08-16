---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_server_key"
description: |-
  Manages an encryption key for a PostgreSQL server.
---

# azurerm_postgresql_server_key

Manages an encryption key for a PostgreSQL server.
Note: This feature requires the keyvault to have soft deletion and purge protection enabled.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "postgres-rg"
  location = "westeurope"
}

resource "azurerm_key_vault" "example" {
  name                        = "examplekeyvault"
  # ...  
  soft_delete_enabled         = true
  purge_protection_enabled    = true
}

resource "azurerm_key_vault_access_policy" "example" {
  key_vault_id = azurerm_key_vault.example.id
  object_id    = azurerm_postgresql_server.example.identity[0].principal_id
  tenant_id    = data.azurerm_client_config.current.tenant_id

  key_permissions = [
    "get", 
    "unwrapKey", 
    "wrapKey"
  ]
}

resource "azurerm_key_vault_key" "generated" {
  name         = "accexample-generated-key-%d"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]

}

resource "azurerm_postgresql_server" "example" {
  # ...

  identity {
    type = "SystemAssigned"
  }

}

resource "azurerm_postgresql_server_key" "example" {
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_postgresql_server.example.name
  key_type            = "AzureKeyVault"
  key_url             = azurerm_key_vault_key.generated.id
}
```

## Arguments Reference

The following arguments are supported:

* `key_type` - (Required) The key type. Currently only `AzureKeyVault` is supported. Changing this forces a new Database to be created.

* `key_url` - (Required) A key vault key url. Changing this forces a new Database to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Database should exist. Changing this forces a new Database to be created.

* `server_name` - (Required) Specifies the name of the PostgreSQL Server. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the server key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Database.
* `read` - (Defaults to 5 minutes) Used when retrieving the Database.
* `update` - (Defaults to 30 minutes) Used when updating the Database.
* `delete` - (Defaults to 30 minutes) Used when deleting the Database.

## Import

Server key can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_postgresql_server_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.DBforPostgreSQL/servers/server1/keys/key1
```
