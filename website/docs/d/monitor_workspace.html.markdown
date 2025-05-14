---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_data_collection_endpoint"
description: |-
  Get information about the specified Workspace.
---

# Data Source: azurerm_monitor_workspace

Use this data source to access information about an existing Workspace.

## Example Usage

```hcl
data "azurerm_monitor_workspace" "example" {
  name                = "example-workspace"
  resource_group_name = azurerm_resource_group.example.name
}

output "query_endpoint" {
  value = data.azurerm_monitor_workspace.example.query_endpoint
}
```

## Argument Reference

- `name` - Specifies the name of the Workspace.

- `resource_group_name` - Specifies the name of the resource group the Workspace is located in.

## Attributes Reference

- `id` - The ID of the Resource.

- `kind` - The kind of the Workspace. Possible values are `Linux` and `Windows`.

- `location` - The Azure Region where the Workspace is located.

- `query_endpoint` - The query endpoint for the Azure Monitor Workspace.

- `public_network_access_enabled` - Whether network access from public internet to the Workspace are allowed.

- `default_data_collection_endpoint_id` - The ID of the managed default Data Collection Endpoint created with the Azure Monitor Workspace.

- `default_data_collection_rule_id` - The ID of the managed default Data Collection Rule created with the Azure Monitor Workspace.

- `tags` - A mapping of tags that are assigned to the Workspace.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Workspace.
