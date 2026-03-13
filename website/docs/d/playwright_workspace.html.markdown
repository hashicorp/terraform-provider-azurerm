---
subcategory: "Load Test"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_playwright_workspace"
description: |-
  Gets information about an existing Playwright Workspace.
---

# Data Source: azurerm_playwright_workspace

Use this data source to access information about an existing Playwright Workspace.

## Example Usage

```hcl
data "azurerm_playwright_workspace" "example" {
  name                = "existing"
  resource_group_name = "existing"
}
output "id" {
  value = data.azurerm_playwright_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Playwright Workspace. Changing this forces a new Playwright Workspace to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Playwright Workspace exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Playwright Workspace.

* `dataplane_uri` - The data plane service API URI of the Playwright Workspace.

* `local_auth_enabled` - Whether the local authentication is enabled.

* `location` - The Azure Region where the Playwright Workspace exists.

* `regional_affinity_enabled` - Whether the regional affinity is enabled.

* `tags` - A mapping of tags assigned to the Playwright Workspace.

* `workspace_id` - The ID in GUID format of the Playwright Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Playwright Workspace.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.LoadTestService` - 2025-09-01
