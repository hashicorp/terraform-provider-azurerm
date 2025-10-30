---
subcategory: "Dashboard"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dashboard"
description: |-
  Manages a Dashboard.
---

# azurerm_dashboard

Manages a Dashboard.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dashboard" "example" {
  name                = "example-dashboard"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  tags = {
    environment = "Production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Dashboard. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Dashboard should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Dashboard should exist. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dashboard.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dashboard.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dashboard.
* `update` - (Defaults to 30 minutes) Used when updating the Dashboard.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dashboard.

## Import

Dashboards can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dashboard.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Dashboard/dashboards/dashboard1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Dashboard` - 2025-08-01
