---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_workspace"
description: |-
  Gets information on an existing Databricks Workspace
---

# Data Source: azurerm_databricks_workspace

Use this data source to access information about an existing Databricks workspace.

## Example Usage

```hcl
data "azurerm_databricks_workspace" "example" {
  name                = "example-workspace"
  resource_group_name = "example-rg"
}

output "databricks_workspace_id" {
  value = data.azurerm_databricks_workspace.example.workspace_id
}
```

## Argument Reference

* `name` - The name of the Databricks Workspace.
* `resource_group_name` - The Name of the Resource Group where the Databricks Workspace exists.

## Attributes Reference

* `id` - The ID of the Databricks Workspace.

* `location` - The Azure location where the Databricks Workspace exists.

* `sku` - SKU of this Databricks Workspace.

* `workspace_id` - Unique ID of this Databricks Workspace in Databricks management plane.

* `workspace_url` - URL this Databricks Workspace is accessible on.

* `managed_disk_identity` - A `managed_disk_identity` block as documented below.

* `storage_account_identity` - A `storage_account_identity` block as documented below.

* `tags` - A mapping of tags to assign to the Databricks Workspace.

---

A `managed_disk_identity` block exports the following:

* `principal_id` - The principal UUID for the internal databricks disks identity needed to provide access to the workspace for enabling Customer Managed Keys.

* `tenant_id` - The UUID of the tenant where the internal databricks disks identity was created.

* `type` - The type of the internal databricks disk identity.

---

A `storage_account_identity` block exports the following:

* `principal_id` - The principal UUID for the internal databricks storage account needed to provide access to the workspace for enabling Customer Managed Keys.

* `tenant_id` - The UUID of the tenant where the internal databricks storage account was created.

* `type` - The type of the internal databricks storage account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Workspace.
