---
subcategory: "Maintenance"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maintenance_assignment_arc_machine"
description: |-
  Lists Maintenance Assignments for an Arc Machine.
---

# List resource: azurerm_maintenance_assignment_arc_machine

Lists Maintenance Assignments for an Arc Machine.

## Example Usage

```hcl
list "azurerm_maintenance_assignment_arc_machine" "example" {
  provider = azurerm
  config {
    arc_machine_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.HybridCompute/machines/example-machine"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `arc_machine_id` - (Required) The ID of the Arc Machine whose Maintenance Assignments should be listed.
