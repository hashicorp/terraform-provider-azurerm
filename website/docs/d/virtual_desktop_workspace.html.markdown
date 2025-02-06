---
subcategory: "Desktop Virtualization"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_virtual_desktop_workspace"
description: |-
  Gets information about an existing Virtual Desktop Workspace.
---

# Data Source: azurerm_virtual_desktop_workspace

Use this data source to access information about an existing Virtual Desktop Workspace.

## Example Usage

```hcl
data "azurerm_virtual_desktop_workspace" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_virtual_desktop_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Virtual Desktop Workspace to retrieve.

* `resource_group_name` - (Required) The name of the Resource Group where the Virtual Desktop Workspace exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Virtual Desktop Workspace.

* `description` - The description for the Virtual Desktop Workspace.

* `friendly_name` - The friendly name for the Virtual Desktop Workspace.

* `location` - The Azure Region where the Virtual Desktop Workspace exists.

* `public_network_access_enabled` - Is public network access enabled?

* `tags` - A mapping of tags assigned to the Virtual Desktop Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the azurerm_virtual_desktop_workspace.
