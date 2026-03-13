---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_discovery_workspace"
description: |-
  Manages a Storage Discovery workspace.
---

# azurerm_storage_discovery_workspace

Manages a Storage Discovery workspace. A workspace defines which storage resources to scan across your Microsoft Entra tenant and how to segment reporting for them.

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

* `name` - (Required) Specifies the name of the Discovery workspace. Must be 4-64 characters long, start with a letter, and contain only letters, numbers, and hyphens (no consecutive hyphens). Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Storage Discovery Workspace should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Discovery resource is created. Changing this forces a new resource to be created.

* `workspace_root` - (Required) A set of top-level Azure resource identifiers (Subscription IDs or Resource Group IDs) where Storage Discovery initiates its scan for storage accounts. You cannot specify both a subscription and its child resource group. Maximum of 100 items.

-> **Note:** The user or service principal must have at least Reader access (`Microsoft.Storage/storageAccounts/read`) on each specified root. The default limit of 100 can be increased by contacting Azure Support.

* `scopes` - (Required) One or more `scopes` blocks as defined below.

* `description` - (Optional) A description for the Discovery workspace resource.

* `sku` - (Optional) Specifies the Storage Discovery pricing plan. Possible values are `Free` and `Standard`. Defaults to `Standard`. See [Understand Storage Discovery Pricing](https://learn.microsoft.com/azure/storage-discovery/understand-pricing) for details.

* `tags` - (Optional) A mapping of tags which should be assigned to the resource.

---

A `scopes` block supports the following:

* `display_name` - (Required) The display name for this scope.

* `resource_types` - (Required) A list of Azure resource type strings to include in this scope. For example, `"Microsoft.Storage/storageAccounts"`.

* `tag_keys_only` - (Optional) A list of tag keys that will be used to filter resources. Resources with any of these tag keys will be included.

* `tags` - (Optional) A map of tag key-value pairs that resources must have to be included in this scope.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Discovery Workspace.

-> **Note:** It can take up to 24 hours after workspace creation for metrics to begin appearing in reports.

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
