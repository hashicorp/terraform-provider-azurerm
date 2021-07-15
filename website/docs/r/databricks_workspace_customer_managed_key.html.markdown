---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_workspace_customer_managed_key"
description: |-
  Manages a Customer Managed Key for a Databricks Workspace
---

# azurerm_databricks_workspace_customer_managed_key

Manages a Customer Managed Key for a Databricks Workspace

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_databricks_workspace" "example" {
  name                = "databricks-test"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "premium"

  customer_managed_key_enabled = true

  tags = {
    Environment = "Production"
  }
}

resource "azurerm_databricks_workspace_customer_managed_key" "example" {
  depends_on = [azurerm_key_vault_access_policy.databricks]

  workspace_id     = azurerm_databricks_workspace.example.id
  key_vault_key_id = azurerm_key_vault.example.id
}

resource "azurerm_key_vault" "example" {
  name                = "examplekeyvault"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "premium"

  purge_protection_enabled = true
}

resource "azurerm_key_vault_key" "example" {
  depends_on = [azurerm_key_vault_access_policy.example]

  name         = "example-certificate"
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

resource "azurerm_key_vault_access_policy" "example" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_key_vault.example.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "get",
    "list",
    "create",
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
    "delete",
    "restore",
    "recover",
    "update",
  ]
}

resource "azurerm_key_vault_access_policy" "databricks" {
  depends_on = [azurerm_databricks_workspace.example]

  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_databricks_workspace.example.storage_account_identity.0.tenant_id
  object_id    = azurerm_databricks_workspace.example.storage_account_identity.0.principal_id

  key_permissions = [
    "get",
    "unwrapKey",
    "wrapKey",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) The ID of the Databricks workspace.

* `key_vault_key_id` - (Required) The ID of the Key Vault.


#Attributes Reference

The following attributes are exported:

* `id` - The ID of the Databricks Customer Managed Key.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Databricks Customer Managed Key.
* `update` - (Defaults to 30 minutes) Used when updating the Databricks Customer Managed Key.
* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Customer Managed Key.
* `delete` - (Defaults to 30 minutes) Used when deleting the Databricks Customer Managed Key.

## Import

Databricks Workspace Customer Managed Key can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_databricks_workspace_customer_managed_key.workspace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Databricks/customerManagedKey/workspace1
```
