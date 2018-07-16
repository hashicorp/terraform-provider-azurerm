---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_log_analytics_solution"
sidebar_current: "docs-azurerm-resource-oms-log-analytics-solution"
description: |-
  Manages a Log Analytics (formally Operational Insights) Solution.
---

# azurerm_log_analytics_solution

Manages a Log Analytics (formally Operational Insights) Solution.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "k8s-log-analytics-test"
  location = "westeurope"
}

resource "random_id" "workspace" {
  keepers = {
    # Generate a new id each time we switch to a new resource group
    group_name = "${azurerm_resource_group.test.name}"
  }

  byte_length = 8
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "k8s-workspace-${random_id.workspace.hex}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Free"
}

resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "Containers"
  location              = "${azurerm_resource_group.test.location}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  workspace_resource_id = "${azurerm_log_analytics_workspace.test.id}"
  workspace_name        = "${azurerm_log_analytics_workspace.test.name}"

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/Containers"
  }
}
```

## Argument Reference

The following arguments are supported:

* `solution_name` - (Required) Specifies the name of the solution to be deployed. See [here for options](https://docs.microsoft.com/en-us/azure/log-analytics/log-analytics-add-solutions).Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Log Analytics solution is created. Changing this forces a new resource to be created. Note: The solution and it's related workspace can only exist in the same resource group.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `workspace_resource_id` - (Required) The full resource ID of the Log Analytics workspace with which the solution will be linked. Changing this forces a new resource to be created.

* `workspace_resource_name` - (Required) The full name of the Log Analytics workspace with which the solution will be linked. Changing this forces a new resource to be created.

* `plan` - A `plan` block as documented below.

---

A `plan` block includes:

* `publisher` - (Required) The publisher of the solution. For example `Microsoft`. Changing this forces a new resource to be created.

* `product` - (Required) The product name of the solution. For example `OMSGallery/Containers`. Changing this forces a new resource to be created.

* `promotion_code` - (Optional) A promotion code to be used with the solution.


## Import

Log Analytics Solutions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_log_analytics_solution.solution1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.OperationsManagement/solutions/solution1
```
