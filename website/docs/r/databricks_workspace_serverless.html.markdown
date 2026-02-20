---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_workspace_serverless"
description: |-
  Manages a Databricks Serverless Workspace.
---

# azurerm_databricks_workspace_serverless

Manages a Databricks Workspace with `serverless` `compute_mode`. It is only supported with premium SKU. To create Databricks Workspace with `hybrid` `compute_mode`, please use `azurerm_databricks_workspace` instead.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_databricks_workspace_serverless" "example" {
  name                = "databricks-test"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  tags = {
    Environment = "Production"
  }
}
```

-> **Note:** You can use [the Databricks Terraform Provider](https://registry.terraform.io/providers/databrickslabs/databricks/latest/docs) to manage resources within the Databricks Serverless Workspace.

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Databricks Serverless Workspace. Changing this forces a new Databricks Serverless Workspace to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Databricks Serverless Workspace should exist. Changing this forces a new Databricks Serverless Workspace to be created.

* `location` - (Required) The Azure Region where the Databricks Serverless Workspace should exist. Changing this forces a new Databricks Serverless Workspace to be created.

---

* `enhanced_security_compliance` - (Optional) An `enhanced_security_compliance` block as defined below.

* `managed_services_cmk_key_vault_id` - (Optional) Resource ID of the Key Vault which contains the `managed_services_cmk_key_vault_key_id` key.

-> **Note:** The `managed_services_cmk_key_vault_id` field is only required if the Key Vault exists in a different subscription than the Databricks Serverless Workspace. If the `managed_services_cmk_key_vault_id` field is not specified it is assumed that the `managed_services_cmk_key_vault_key_id` is hosted in the same subscriptioin as the Databricks Serverless Workspace.

-> **Note:** If you are using multiple service principals to execute Terraform across subscriptions you will need to add an additional `azurerm_key_vault_access_policy` resource granting the service principal access to the key vault in that subscription.

* `managed_services_cmk_key_vault_key_id` - (Optional) Customer managed encryption properties for the Databricks Serverless Workspace managed resources(e.g. Notebooks and Artifacts).

* `public_network_access_enabled` - (Optional) Allow public access for accessing the Databricks Serverless Workspace. Set value to `false` to access the Databricks Serverless Workspace only via private link endpoint. Possible values include `true` or `false`. Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Databricks Serverless Workspace.

---

A `enhanced_security_compliance` block supports the following:

* `automatic_cluster_update_enabled` - (Optional) Enables automatic cluster updates for the Databricks Serverless Workspace. Defaults to `false`.

* `compliance_security_profile_enabled` - (Optional) Enables compliance security profile for the Databricks Serverless Workspace. Defaults to `false`.

~> **Note:** Changing the value of `compliance_security_profile_enabled` from `true` to `false` forces a replacement of the Databricks Serverless Workspace.

~> **Note:** The attributes `automatic_cluster_update_enabled` and `enhanced_security_monitoring_enabled` must be set to `true` in order to set `compliance_security_profile_enabled` to `true`.

* `compliance_security_profile_standards` - (Optional) A list of standards to enforce on the Databricks Serverless Workspace. Possible value include `HIPAA`.

~> **Note:** `compliance_security_profile_enabled` must be set to `true` in order to use `compliance_security_profile_standards`.

~> **Note:** Removing a standard from the `compliance_security_profile_standards` list forces a replacement of the Databricks Serverless Workspace.

* `enhanced_security_monitoring_enabled` - (Optional) Enables enhanced security monitoring for the Databricks Serverless Workspace. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Databricks Serverless Workspace.

* `workspace_id` - The unique identifier of the Databricks Serverless Workspace in Databricks control plane.

* `workspace_url` - The Databricks Serverless Workspace URL which is of the format 'adb-{workspaceId}.{random}.azuredatabricks.net'

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Databricks Serverless Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Serverless Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Databricks Serverless Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Databricks Serverless Workspace.

## Import

Databricks Serverless Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_databricks_workspace_serverless.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Databricks/workspaces/workspace1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Databricks` - 2026-01-01
