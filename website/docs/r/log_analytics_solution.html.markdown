---
subcategory: "Log Analytics"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_solution"
description: |-
  Manages a Log Analytics (formally Operational Insights) Solution.
---

# azurerm_log_analytics_solution

Manages a Log Analytics (formally Operational Insights) Solution.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "k8s-log-analytics-test"
  location = "West Europe"
}

resource "random_id" "workspace" {
  keepers = {
    # Generate a new id each time we switch to a new resource group
    group_name = azurerm_resource_group.example.name
  }

  byte_length = 8
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "k8s-workspace-${random_id.workspace.hex}"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_solution" "example" {
  solution_name         = "ContainerInsights"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  workspace_resource_id = azurerm_log_analytics_workspace.example.id
  workspace_name        = azurerm_log_analytics_workspace.example.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/ContainerInsights"
  }
}
```

## Argument Reference

The following arguments are supported:

* `solution_name` - (Required) Specifies the name of the solution to be deployed. See [here for options](https://docs.microsoft.com/en-us/azure/log-analytics/log-analytics-add-solutions).Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Log Analytics solution is created. Changing this forces a new resource to be created. Note: The solution and its related workspace can only exist in the same resource group.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `workspace_resource_id` - (Required) The full resource ID of the Log Analytics workspace with which the solution will be linked. Changing this forces a new resource to be created.

* `workspace_name` - (Required) The full name of the Log Analytics workspace with which the solution will be linked. Changing this forces a new resource to be created.

* `plan` - (Required) A `plan` block as documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `plan` block includes:

* `publisher` - (Required) The publisher of the solution. For example `Microsoft`. Changing this forces a new resource to be created.

* `product` - (Required) The product name of the solution. For example `OMSGallery/Containers`. Changing this forces a new resource to be created.

* `promotion_code` - (Optional) A promotion code to be used with the solution.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Log Analytics Solution.
* `update` - (Defaults to 30 minutes) Used when updating the Log Analytics Solution.
* `read` - (Defaults to 5 minutes) Used when retrieving the Log Analytics Solution.
* `delete` - (Defaults to 30 minutes) Used when deleting the Log Analytics Solution.

## Import

Log Analytics Solutions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_solution.solution1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationsManagement/solutions/solution1
```
