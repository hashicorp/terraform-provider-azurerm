---
subcategory: "Advisor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_advisor_configurations"
sidebar_current: "docs-azurerm-resource-advisor-configurations"
description: |-
  Configure a subscription's advisor or a resource group's advisor.
---

# azurerm_advisor_configurations

Configure a subscription's advisor or a resource group's advisor.

Advisor analyzes the recent usage patterns of your virtual machines and uses rules to identify low usage virtual machines. You can customize these rules to better match your business needs.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_advisor_configurations" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  exclude             = false
  low_cpu_threshold   = "5"
}


```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Optional) The name of the Resource Group. The advisor of which you want to configure. If not assigned, the advisor of the subscription will be configured.

* `exclude` - (Optional) Exclude the resource from Advisor evaluations.Valid values: false (default) or true.

* `low_cpu_threshold` - (Optional)  Minimum percentage threshold for Advisor low CPU utilization evaluation. Valid only for subscriptions. Valid values: 5 (default), 10, 15 or 20.

