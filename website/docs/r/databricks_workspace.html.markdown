---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_workspace"
sidebar_current: "docs-azurerm-resource-databricks-workspace"
description: |-
  Manages a Databricks Workspace
---

# azurerm_databricks_workspace

Manages a Databricks Workspace

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_databricks_workspace" "test" {
  name                = "databricks-test"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
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

* `sku` - (Required) The `sku` to use for the Databricks Workspace. Possible values are `standard` or `premium`. Changing this forces a new resource to be created.

* `managed_resource_group_name` - (Optional) The name of the resource group where Azure should place the managed Databricks resources. Changing this forces a new resource to be created.

~> **NOTE** Azure requires that this Resource Group does not exist in this Subscription (and that the Azure API creates it) - otherwise the deployment will fail.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Databricks Workspace.

* `managed_resource_group_id` - The ID of the Managed Resource Group created by the Databricks Workspace.

## Import

Databrick Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_databricks_workspace.workspace1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Databricks/workspaces/workspace1
```
