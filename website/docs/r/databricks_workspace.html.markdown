---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_workspace"
description: |-
  Manages a Databricks Workspace
---

# azurerm_databricks_workspace

Manages a Databricks Workspace

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_databricks_workspace" "example" {
  name                = "databricks-test"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "standard"

  tags = {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Databricks Workspace resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Databricks Workspace should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource has to be created. Changing this forces a new resource to be created.

* `sku` - (Required) The `sku` to use for the Databricks Workspace. Possible values are `standard`, `premium`, or `trial`. Changing this can force a new resource to be created in some circumstances.

~> **NOTE** Downgrading the `sku` to the `trial` version of the Databricks Workspace will force a new resource to be created.

* `managed_resource_group_name` - (Optional) The name of the resource group where Azure should place the managed Databricks resources. Changing this forces a new resource to be created.

~> **NOTE** Azure requires that this Resource Group does not exist in this Subscription (and that the Azure API creates it) - otherwise the deployment will fail.

* `custom_parameters` - (Optional) A `custom_parameters` block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `custom_parameters` block supports the following:

* `aml_workspace_id` - (Optional) The ID of a Azure Machine Learning workspace to link with Databricks workspace. Changing this forces a new resource to be created.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as documented below.

* `customer_managed_key_enabled` - (Optional) Is the workspace enabled for CMK encryption? If `true` this enables the Managed Identity for the managed storage account. Possible values are `true` or `false`. Defaults to `false`.

~> **NOTE** Once `customer_managed_key_enabled` has been set to `true` it cannot be set back to `false`. If you wish to remove your customer managed key encryption from your workspace you will need to update the `customer_managed_key` blocks `source` field to be `Default` and remove the `customer_managed_key_enabled` attribute from the `custom_parameters` block. `customer_managed_key_enabled` is olny available with the `premium` Databricks Workspace `sku`.

* `infrastructure_encryption_enabled`- (Optional) Is the DBFS root file system enabled with a secondary layer of encryption with platform managed keys for data at rest? Possible values are `true` or `false`. Defaults to `false`.

~> **NOTE** Once `infrastructure_encryption_enabled` has been set to `true` it cannot be set back to `false`. `infrastructure_encryption_enabled` is olny available with the `premium` Databricks Workspace `sku`.

* `no_public_ip` - (Optional) Are public IP Addresses not allowed? Possible values are `true` or `false`. Defaults to `false`. Changing this forces a new resource to be created.

* `public_subnet_name` - (Optional) The name of the Public Subnet within the Virtual Network. Required if `virtual_network_id` is set. Changing this forces a new resource to be created.

* `private_subnet_name` - (Optional) The name of the Private Subnet within the Virtual Network. Required if `virtual_network_id` is set. Changing this forces a new resource to be created.

* `virtual_network_id` - (Optional) The ID of a Virtual Network where this Databricks Cluster should be created. Changing this forces a new resource to be created.

~> **NOTE** Databricks requires that a network security group is associated with public and private subnets when `virtual_network_id` is set. Also, both public and private subnets must be delegated to `Microsoft.Databricks/workspaces`.

---

A `customer_managed_key` block supports the following:

* `source` - (Optional) The encryption key source. Possible values include: `Default` or `Microsoft.Keyvault`. Defaults to `Default`.

* `name` - (Optional) The name of Key Vault key.

* `version` - (Optional) The version of Key Vault key.

* `vault_uri` - (Optional) The Uri of Key Vault.

~> **NOTE** To successfully provision `customer_managed_key` you must first set the `customer_managed_key_enabled` field to `true`. Once the `customer_managed_key_enabled` has been set to `true` you will then need to add the `customer_managed_key` block into your configuration file and `apply` the changes.

---

A `provider_authorization` block supports the following:

* `principal_id` - (Required) The provider's principal UUID. This is the identity that the provider will use to call ARM to manage the workspace resources. Changing this forces a new resource to be created.

* `role_definition_id` - (Required) The provider's role definition UUID. This role will define all the permissions that the provider must have on the workspace's container resource group. This role definition cannot have permission to delete the resource group. Changing this forces a new resource to be created.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Databricks Workspace in the Azure management plane.

* `managed_resource_group_id` - The ID of the Managed Resource Group created by the Databricks Workspace.

* `workspace_url` - The workspace URL which is of the format 'adb-{workspaceId}.{random}.azuredatabricks.net'

* `workspace_id` - The unique identifier of the databricks workspace in Databricks control plane.

* `storage_account_identity` - A `storage_account_identity` block as documented below.

---

A `storage_account_identity` block exports the following:

* `principal_id` - The principal UUID for the internal databricks storage account needed to provide access to the workspace for enabling Customer Managed Keys.

* `tenant_id` - The UUID of the tenant where the internal databricks storage account was created.

* `type` - The type of the internal databricks storage account.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Databricks Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Databricks Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Databricks Workspace.

## Import

Databrick Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_databricks_workspace.workspace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Databricks/workspaces/workspace1
```
