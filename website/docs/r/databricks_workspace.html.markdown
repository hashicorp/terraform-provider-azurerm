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
  location = "West US"
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

* `sku` - (Required) The `sku` to use for the Databricks Workspace. Possible values are `standard`, `premium`, or `trial`.

~> **NOTE** While downgrading to `trial`, the Databricks Workspace resource would be recreated.

* `managed_resource_group_name` - (Optional) The name of the resource group where Azure should place the managed Databricks resources. Changing this forces a new resource to be created.

~> **NOTE** Azure requires that this Resource Group does not exist in this Subscription (and that the Azure API creates it) - otherwise the deployment will fail.

* `custom_parameters` - (Optional) A `custom_parameters` block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`custom_parameters` supports the following:

* `no_public_ip` - (Optional) Are public IP Addresses not allowed?

* `public_subnet_name` - (Optional) The name of the Public Subnet within the Virtual Network. Required if `virtual_network_id` is set.

* `private_subnet_name` - (Optional) The name of the Private Subnet within the Virtual Network. Required if `virtual_network_id` is set.

* `virtual_network_id` - (Optional) The ID of a Virtual Network where this Databricks Cluster should be created.

~> **NOTE** Databricks requires that a network security group is associated with public and private subnets when `virtual_network_id` is set.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Databricks Workspace in the Azure management plane.

* `managed_resource_group_id` - The ID of the Managed Resource Group created by the Databricks Workspace.

* `workspace_url` - The workspace URL which is of the format 'adb-{workspaceId}.{random}.azuredatabricks.net'

* `workspace_id` - The unique identifier of the databricks workspace in Databricks control plane.

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
