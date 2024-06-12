---
subcategory: "Maintenance"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_maintenance_assignment_dynamic_scope"
description: |-
  Manages a Dynamic Maintenance Assignment
---

# azurerm_maintenance_assignment_dynamic_scope

Manages a Dynamic Maintenance Assignment.

~> **Note:** Only valid for `InGuestPatch` Maintenance Configuration Scopes.

## Example Usage

```hcl

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_maintenance_configuration" "example" {
  name                     = "example"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  scope                    = "InGuestPatch"
  in_guest_user_patch_mode = "User"

  window {
    start_date_time = formatdate("YYYY-MM-DD hh:mm", timestamp())
    time_zone       = "Greenwich Standard Time"
    recur_every     = "1Day"
  }

  install_patches {
    reboot = "Always"

    windows {
      classifications_to_include = ["Critical"]
      kb_numbers_to_exclude      = []
      kb_numbers_to_include      = []
    }
  }
}

resource "azurerm_maintenance_assignment_dynamic_scope" "example" {
  name                         = "example"
  maintenance_configuration_id = azurerm_maintenance_configuration.example.id

  filter {
    locations       = ["West Europe"]
    os_types        = ["Windows"]
    resource_groups = [azurerm_resource_group.example.name]
    resource_types  = ["Microsoft.Compute/virtualMachines"]
    tag_filter      = "Any"
    tags {
      tag    = "foo"
      values = ["barbar"]
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `maintenance_configuration_id` - (Required) The ID of the Maintenance Configuration Resource. Changing this forces a new Dynamic Maintenance Assignment to be created.

* `name` - (Required) The name which should be used for this Dynamic Maintenance Assignment. Changing this forces a new Dynamic Maintenance Assignment to be created.

~> **Note:** The `name` must be unique per subscription.

* `filter` - (Required) A `filter` block as defined below.

---

A `filter` block supports the following:

* `locations` - (Optional) Specifies a list of locations to scope the query to.

* `os_types` - (Optional) Specifies a list of allowed operating systems.

* `resource_groups` - (Optional) Specifies a list of allowed resource groups.

* `resource_types` - (Optional) Specifies a list of allowed resources.

* `tag_filter` - (Optional) Filter VMs by `Any` or `All` specified tags. Defaults to `Any`.

* `tags` - (Optional) A mapping of tags for the VM

---

A `tags` block supports the following:

* `tag` - (Required) Specifies the tag to filter by.

* `values` - (Required) Specifies a list of values the defined tag can have.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Dynamic Maintenance Assignment

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dynamic Maintenance Assignment
* `read` - (Defaults to 5 minutes) Used when retrieving the Dynamic Maintenance Assignment
* `update` - (Defaults to 30 minutes) Used when updating the Dynamic Maintenance Assignment
* `delete` - (Defaults to 10 minutes) Used when deleting the Dynamic Maintenance Assignment

## Import

Dynamic Maintenance Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_maintenance_assignment_dynamic_scope.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Maintenance/configurationAssignments/assignmentName
```
