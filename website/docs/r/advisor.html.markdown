---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_advisor"
description: |-
  Configure an Advisor
---

# azurerm_advisor

Configure an Advisor

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_advisor" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  exclude             = false
}


```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Optional) The name of the Resource Group. The advisor of the resource group that you want to configure. If not assigned, the advisor of the subscription will be configured.

* `exclude` - (Optional) Should exclude the resource from Advisor evaluations?

* `low_cpu_threshold` - (Optional)  Minimum percentage threshold for Advisor low CPU utilization evaluation. Valid only for subscriptions. Possible values are 5, 10, 15 or 20.

