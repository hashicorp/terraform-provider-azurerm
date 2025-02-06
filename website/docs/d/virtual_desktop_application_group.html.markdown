---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_application_group"
description: |-
  Gets information about an existing Application Group.
---

# Data Source: azurerm_virtual_desktop_application_group

Use this data source to access information about an existing Application Group.

## Example Usage

```hcl
data "azurerm_virtual_desktop_application_group" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_virtual_desktop_application_group.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Application Group.

* `resource_group_name` - (Required) The name of the Resource Group where the Application Group exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Application Group.

* `description` - The description of the Application Group.

* `friendly_name` - The friendly name of the Application Group.

* `host_pool_id` - The Virtual Desktop Host Pool ID the Application Group is associated to.

* `location` - The Azure Region where the Application Group exists.

* `tags` - A mapping of tags assigned to the Application Group.

* `type` - The type of Application Group (`RemoteApp` or `Desktop`).

* `workspace_id` - The Virtual Desktop Workspace ID the Application Group is associated to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Application Group.
