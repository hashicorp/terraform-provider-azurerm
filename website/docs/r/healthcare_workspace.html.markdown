---
subcategory: "Healthcare"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_healthcare_workspace"
description: |-
  Manages a Healthcare Workspace.
---

# azurerm_healthcare_workspace

Manages a Healthcare workspace

## Example Usage

```hcl
resource "azurerm_healthcare_workspace" "test" {
  name                = "tfexworkspace"
  resource_group_name = "tfex-resource_group"
  location            = "east us"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Healthcare Workspace. Changing this forces a new Healthcare Workspace to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Healthcare Workspace should exist. Changing this forces a new Healthcare Workspace to be created.

* `location` - (Required) Specifies the Azure Region where the Healthcare Workspace should be created. Changing this forces a new Healthcare Workspace to be created.

* `tags` - (Optional) A mapping of tags to assign to the Healthcare Workspace.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Healthcare Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Healthcare Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Healthcare Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Healthcare Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Healthcare Workspace.

## Import

Healthcare Workspaces can be imported using the resource`id`, e.g.

```shell
terraform import azurerm_healthcare_workspace.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.HealthcareApis/workspaces/workspace1
```
