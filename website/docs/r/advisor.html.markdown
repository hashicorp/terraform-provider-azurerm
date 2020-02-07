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
resource "azurerm_resource_group" "test" {
  name     = "testaccRG-%d"
  location = "%s"
}

resource "azurerm_advisor" "test" {
  low_cpu_threshold       = "5"
  exclude_resource_groups = [azurerm_resource_group.test.name]
}
```

## Argument Reference

The following arguments are supported:

* `low_cpu_threshold` - (Optional)  Minimum percentage threshold for Advisor low CPU utilization evaluation. Possible values are 5, 10, 15 or 20.

* `exclude_resource_groups` - (Optional) A set of resource group name which should be excluded from Advisor?

