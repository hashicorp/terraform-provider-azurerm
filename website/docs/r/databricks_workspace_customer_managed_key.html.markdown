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
  key_vault_key_id = azurerm_key_vault_key.example.id
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
  depends_on = [azurerm_key_vault_access_policy.terraform]

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

resource "azurerm_key_vault_access_policy" "terraform" {
  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_key_vault.example.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign"
  ]
}

resource "azurerm_key_vault_access_policy" "databricks" {
  depends_on = [azurerm_databricks_workspace.example]

  key_vault_id = azurerm_key_vault.example.id
  tenant_id    = azurerm_databricks_workspace.example.storage_account_identity.0.tenant_id
  object_id    = azurerm_databricks_workspace.example.storage_account_identity.0.principal_id

  key_permissions = [
    "Create",
    "Delete",
    "Get",
    "Purge",
    "Recover",
    "Update",
    "List",
    "Decrypt",
    "Sign"
  ]
}
```
## Example HCL Configurations

* [Databricks Workspace with Databricks File System Customer Managed Keys](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/databricks/customer-managed-key/dbfs)
* [Databricks Workspace with Customer Managed Keys for Managed Services](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/databricks/customer-managed-key/managed-services)
* [Databricks Workspace with Private Endpoint, Customer Managed Keys for Managed Services and Databricks File System Customer Managed Keys](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/private-endpoint/databricks/managed-services)


## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required) The ID of the Databricks Workspace..

* `key_vault_key_id` - (Required) The ID of the Key Vault.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Databricks Workspace.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Customer Managed Key for this Databricks Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Customer Managed Key for this Databricks Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Customer Managed Key for this Databricks Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Customer Managed Key for this Databricks Workspace.

## Import

Databricks Workspace Customer Managed Key can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_databricks_workspace_customer_managed_key.workspace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Databricks/workspaces/workspace1
```
