---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mysql_server_key"
description: |-
  Manages a MySQL Server Key.
---

# azurerm_mysql_server_key

Manages a Customer Managed Key for a MySQL Server.

~> **Note:** Azure Database for MySQL Single Server and its sub resources are scheduled for retirement by 2024-09-16 and will migrate to using Azure Database for MySQL Flexible Server: https://go.microsoft.com/fwlink/?linkid=2216041. The `azurerm_mysql_server_key` resource is deprecated and will be removed in v4.0 of the AzureRM Provider. Please use the `customer_managed_key` property of the `azurerm_mysql_flexible_server` resource instead.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_key_vault" "example" {
  name                     = "examplekv"
  location                 = azurerm_resource_group.example.location
  resource_group_name      = azurerm_resource_group.example.name
  tenant_id                = data.azurerm_client_config.current.tenant_id
  sku_name                 = "premium"
  purge_protection_enabled = true
}

resource "azurerm_key_vault_access_policy" "server" {
  key_vault_id       = azurerm_key_vault.example.id
  tenant_id          = data.azurerm_client_config.current.tenant_id
  object_id          = azurerm_mysql_server.example.identity[0].principal_id
  key_permissions    = ["Get", "UnwrapKey", "WrapKey"]
  secret_permissions = ["Get"]
}

resource "azurerm_key_vault_access_policy" "client" {
  key_vault_id       = azurerm_key_vault.example.id
  tenant_id          = data.azurerm_client_config.current.tenant_id
  object_id          = data.azurerm_client_config.current.object_id
  key_permissions    = ["Get", "Create", "Delete", "List", "Restore", "Recover", "UnwrapKey", "WrapKey", "Purge", "Encrypt", "Decrypt", "Sign", "Verify", "GetRotationPolicy"]
  secret_permissions = ["Get"]
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

resource "azurerm_mysql_server" "example" {
  name                             = "example-mysql-server"
  location                         = azurerm_resource_group.example.location
  resource_group_name              = azurerm_resource_group.example.name
  sku_name                         = "GP_Gen5_2"
  administrator_login              = "acctestun"
  administrator_login_password     = "H@Sh1CoR3!"
  ssl_enforcement_enabled          = true
  ssl_minimal_tls_version_enforced = "TLS1_1"
  storage_mb                       = 51200
  version                          = "5.7"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_mysql_server_key" "example" {
  server_id        = azurerm_mysql_server.example.id
  key_vault_key_id = azurerm_key_vault_key.example.id
}
```

## Argument Reference

The following arguments are supported:

* `server_id` - (Required) The ID of the MySQL Server. Changing this forces a new resource to be created.

* `key_vault_key_id` - (Required) The URL to a Key Vault Key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MySQL Server Key.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the MySQL Server Key.
* `update` - (Defaults to 60 minutes) Used when updating the MySQL Server Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the MySQL Server Key.
* `delete` - (Defaults to 60 minutes) Used when deleting the MySQL Server Key.

## Import

A MySQL Server Key can be imported using the `resource id` of the MySQL Server Key, e.g.

```shell
terraform import azurerm_mysql_server_key.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DBforMySQL/servers/server1/keys/keyvaultname_key-name_keyversion
```
