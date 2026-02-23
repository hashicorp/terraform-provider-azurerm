---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_discovery_workspace"
description: |-
  Manages a Storage Discovery Workspace.
---

# azurerm_storage_discovery_workspace

Manages a Storage Discovery Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_subscription" "current" {}

resource "azurerm_storage_discovery_workspace" "example" {
  name                = "example-sdw"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  description         = "Example Storage Discovery Workspace"
  sku                 = "Standard"

  workspace_root = [data.azurerm_subscription.current.id]

  scopes {
    display_name   = "Production Storage"
    resource_types = ["Microsoft.Storage/storageAccounts"]
    tag_keys_only  = ["environment", "department"]
    tags = {
      criticality = "high"
    }
  }

  tags = {
    environment = "production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Storage Discovery Workspace. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Storage Discovery Workspace should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Storage Discovery Workspace should exist. Changing this forces a new resource to be created.

* `workspace_root` - (Required) A set of Azure Resource IDs that define the root scope for storage discovery. Each ID can be either a Subscription ID or a Resource Group ID. You cannot specify both a subscription and its child resource group. Maximum of 100 items.

* `scopes` - (Required) One or more `scopes` blocks as defined below.

---

* `description` - (Optional) A description for the Storage Discovery Workspace.

* `sku` - (Optional) The SKU for the Storage Discovery Workspace. Possible values are `Free` and `Standard`. Defaults to `Standard`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `scopes` block supports the following:

* `display_name` - (Required) The display name for this scope.

* `resource_types` - (Required) A list of Azure resource type strings to include in this scope. For example, `"Microsoft.Storage/storageAccounts"`.

* `tag_keys_only` - (Optional) A list of tag keys that will be used to filter resources. Resources with any of these tag keys will be included.

* `tags` - (Optional) A map of tag key-value pairs that resources must have to be included in this scope.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Discovery Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Discovery Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Discovery Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Discovery Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Discovery Workspace.

## Import

Storage Discovery Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_discovery_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StorageDiscovery/storageDiscoveryWorkspaces/workspace1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.StorageDiscovery` - 2025-09-01
