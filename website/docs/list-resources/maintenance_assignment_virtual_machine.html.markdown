---
subcategory: "Maintenance"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maintenance_assignment_virtual_machine"
description: |-
  Lists Maintenance Configuration Assignments on a Virtual Machine.
---

# List resource: azurerm_maintenance_assignment_virtual_machine

Lists Maintenance Configuration Assignments on a Virtual Machine.

## Example Usage

### List all Maintenance Configuration Assignments on a Virtual Machine

```hcl
list "azurerm_maintenance_assignment_virtual_machine" "example" {
  provider = azurerm
  config {
    virtual_machine_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `virtual_machine_id` - (Required) The ID of the Virtual Machine to query Maintenance Configuration Assignments for.
