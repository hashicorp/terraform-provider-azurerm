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

## Arguments Reference

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

- `private_endpoint_connections` - A list of `private_endpoint_connections` blocks as described below.

---

A `private_endpoint_connections` block exports the following:

- `name` - The name of the private endpoint connection.

- `id` - The ID of the private endpoint connection.

- `group_ids` - A list of group IDs (sometimes called subresource names) that this private endpoint connection allows access to.

-> **NOTE:** The Azure API does not provide a way to uniquely identify a private endpoint connection based on an ID alone. The `id` exported by the `private_endpoint_connection` block does not correspond to the `id` of the private endpoint created by the connecting resource (e.g., an `azurerm_dashboard_grafana_managed_private_endpoint`), and there is no such property available. In order to identify the private endpoint connection, both `id` and `group_ids` will be required.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Workspace.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Monitor` - 2023-04-03
