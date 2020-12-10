---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_postgresql_server_key"
description: |-
  Manages a PostgreSQL Server Key.
---

# azurerm_postgresql_server_key

Manages a Customer Managed Key for a PostgreSQL Server.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                = "examplekv"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  soft_delete_enabled      = true
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "server" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_postgresql_server.example.identity.0.principal_id

  key_permissions    = ["get", "unwrapkey", "wrapkey"]
  secret_permissions = ["get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions    = ["get", "create", "delete", "list", "restore", "recover", "unwrapkey", "wrapkey", "purge", "encrypt", "decrypt", "sign", "verify"]
  secret_permissions = ["get"]
}

resource "azurerm_key_vault_key" "example" {
  name         = "tfex-key"
  key_vault_id = azurerm_key_vault.example.id
  key_type     = "RSA"
  key_size     = 2048
  key_opts     = ["decrypt", "encrypt", "sign", "unwrapKey", "verify", "wrapKey"]

  depends_on = [
    azurerm_key_vault_access_policy.client,
    azurerm_key_vault_access_policy.server,
  ]
}

resource "azurerm_postgresql_server" "example" {
  name                = "example-postgre-server"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  administrator_login          = "psqladminun"
  administrator_login_password = "H@Sh1CoR3!"

  sku_name   = "GP_Gen5_2"
  version    = "11"
  storage_mb = 51200

  ssl_enforcement_enabled = true

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_postgresql_server_key" "example" {
  server_id        = azurerm_postgresql_server.example.id
  key_vault_key_id = azurerm_key_vault_key.example.id
}
```

## Argument Reference

The following arguments are supported:

* `server_id` - (Required) The ID of the PostgreSQL Server. Changing this forces a new resource to be created.

* `key_vault_key_id` - (Required) The URL to a Key Vault Key.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the PostgreSQL Server Key.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the PostgreSQL Server Key.
* `update` - (Defaults to 30 minutes) Used when updating the PostgreSQL Server Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the PostgreSQL Server Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the PostgreSQL Server Key.

## Import

A PostgreSQL Server Key can be imported using the `resource id` of the PostgreSQL Server Key, e.g.

```shell
terraform import azurerm_postgresql_server_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforPostgreSQL/servers/server1/keys/keyvaultname_key-name_keyversion
```
