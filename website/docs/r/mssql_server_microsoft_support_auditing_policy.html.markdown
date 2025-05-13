---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_server_microsoft_support_auditing_policy"
description: |-
  Manages a MS SQL Server Microsoft Support Auditing Policy.
---

# azurerm_mssql_server_microsoft_support_auditing_policy

Manages a MS SQL Server Microsoft Support Auditing Policy.


## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mssql_server" "example" {
  name                         = "example-sqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "AdminPassword123!"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_mssql_server_microsoft_support_auditing_policy" "example" {
  server_id                  = azurerm_mssql_server.example.id
  blob_storage_endpoint      = azurerm_storage_account.example.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.example.primary_access_key
}
```

## Example Usage with storage account behind VNet and firewall
```hcl
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_client_config" "example" {
}

resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "virtnetname-1"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                                           = "subnetname-1"
  resource_group_name                            = azurerm_resource_group.example.name
  virtual_network_name                           = azurerm_virtual_network.example.name
  address_prefixes                               = ["10.0.2.0/24"]
  service_endpoints                              = ["Microsoft.Sql", "Microsoft.Storage"]
  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_role_assignment" "example" {
  scope                = data.azurerm_subscription.primary.id
  role_definition_name = "Storage Blob Data Contributor"
  principal_id         = azurerm_mssql_server.example.identity[0].principal_id
}

resource "azurerm_mssql_server" "example" {
  name                         = "example-sqlserver"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  version                      = "12.0"
  administrator_login          = "missadministrator"
  administrator_login_password = "AdminPassword123!"
  minimum_tls_version          = "1.2"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_sql_virtual_network_rule" "sqlvnetrule" {
  name                = "sql-vnet-rule"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mssql_server.example.name
  subnet_id           = azurerm_subnet.example.id

}

resource "azurerm_sql_firewall_rule" "example" {
  name                = "FirewallRule1"
  resource_group_name = azurerm_resource_group.example.name
  server_name         = azurerm_mssql_server.example.name
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}

resource "azurerm_storage_account" "example" {
  name                = "examplesa"
  resource_group_name = azurerm_resource_group.example.name

  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"

  allow_nested_items_to_be_public = false

  network_rules {
    default_action             = "Deny"
    ip_rules                   = ["127.0.0.1"]
    virtual_network_subnet_ids = [azurerm_subnet.example.id]
    bypass                     = ["AzureServices"]
  }

  identity {
    type = "SystemAssigned"

  }
}

resource "azurerm_mssql_server_microsoft_support_auditing_policy" "example" {
  blob_storage_endpoint  = azurerm_storage_account.example.primary_blob_endpoint
  server_id              = azurerm_mssql_server.example.id
  log_monitoring_enabled = false

  storage_account_subscription_id = azurerm_subscription.primary.subscription_id

  depends_on = [
    azurerm_role_assignment.example,
    azurerm_storage_account.example,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `server_id` - (Required) The ID of the SQL Server to set the extended auditing policy. Changing this forces a new resource to be created.

* `enabled` - (Optional) Whether to enable the extended auditing policy. Possible values are `true` and `false`. Defaults to `true`.

-> **Note:** If `enabled` is `true`, `blob_storage_endpoint` or `log_monitoring_enabled` are required.

* `blob_storage_endpoint` - (Optional) The blob storage endpoint (e.g. https://example.blob.core.windows.net). This blob storage will hold all Microsoft support auditing logs.

* `storage_account_access_key` - (Optional) The access key to use for the auditing storage account.

* `log_monitoring_enabled` - (Optional) Enable audit events to Azure Monitor? To enable server audit events to Azure Monitor, please enable its main database audit events to Azure Monitor. Defaults to `true`.

* `storage_account_subscription_id` - (Optional) The ID of the Subscription containing the Storage Account.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the MS SQL Server Microsoft Support Auditing Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MS SQL Server Microsoft Support Auditing Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the MS SQL Server Microsoft Support Auditing Policy.
* `update` - (Defaults to 30 minutes) Used when updating the MS SQL Server Microsoft Support Auditing Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the MS SQL Server Microsoft Support Auditing Policy.

## Import

MS SQL Server Microsoft Support Auditing Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_server_microsoft_support_auditing_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Sql/servers/sqlServer1/devOpsAuditingSettings/default
```
