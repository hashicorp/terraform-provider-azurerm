---
subcategory: "Load Test"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_playwright_workspace"
description: |-
  Manages a Playwright Workspace.
---

# azurerm_playwright_workspace

Manages a Playwright Workspace.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "West Europe"
}
resource "azurerm_playwright_workspace" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Playwright Workspace. Changing this forces a new Playwright Workspace to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Playwright Workspace should exist. Changing this forces a new Playwright Workspace to be created.

* `location` - (Required) The Azure Region where the Playwright Workspace should exist. Changing this forces a new Playwright Workspace to be created.

---

* `local_auth_enabled` - (Optional) Whether to enable local authentication through service access tokens for operations. Defaults to `false`.

* `regional_affinity_enabled` - (Optional) Whether the regional affinity is enabled. When enabled, workers connect to browsers in the closest Azure region for lower latency. When disabled, workers connect to browsers in the Azure region where the Playwright Workspace was created. Defaults to `true`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Playwright Workspace.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Playwright Workspace.

* `dataplane_uri` - The data plane service API URI of the Playwright Workspace.

* `workspace_id` - The ID in GUID format of the Playwright Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Playwright Workspace.
* `read` - (Defaults to 5 minutes) Used when retrieving the Playwright Workspace.
* `update` - (Defaults to 30 minutes) Used when updating the Playwright Workspace.
* `delete` - (Defaults to 30 minutes) Used when deleting the Playwright Workspace.

## Import

Playwright Workspaces can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_playwright_workspace.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.LoadTestService/playwrightWorkspaces/workspace1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.LoadTestService` - 2025-09-01
