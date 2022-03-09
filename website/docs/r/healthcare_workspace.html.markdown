---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcareapis_workspace"
description: |-
  Manages a Healthcare API Workspace.
---

# azurerm_healthcareapis_workspace

Manages a Healthcare workspace

## Example Usage

```hcl
resource "azurerm_healthcareapis_workspace" "test" {
  name                = "tfexworkspace"
  resource_group_name = "tfex-resource_group"
  location            = "east us"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the workspace instance. 
* `resource_group_name` - (Required) The name of the Resource Group in which to create the Workspace.
* `location` - (Required) Specifies the supported Azure Region where the Workspace should be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Healthcare Workspace.

## Timeouts
The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Healthcare Workspace.
  * `update` - (Defaults to 30 minutes) Used when updating the Healthcare Workspace.
  * `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Workspace.
  * `delete` - (Defaults to 30 minutes) Used when deleting the Healthcare Workspace.

## Import

Healthcare Service can be imported using the resource`id`, e.g.

```shell
terraform import azurerm_healthcareapis_workspace.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1
```
