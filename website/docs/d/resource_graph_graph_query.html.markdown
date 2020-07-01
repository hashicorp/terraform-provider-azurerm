---
subcategory: "Resource Graph"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_resource_graph_graph_query"
description: |-
  Gets information about an existing Resource Graph.
---

# Data Source: azurerm_resource_graph_graph_query

Use this data source to access information about an existing Resource Graph.

## Example Usage

```hcl
data "azurerm_resource_graph_graph_query" "example" {
  resource_group_name = "existing"
  resource_name       = "existing"
}

output "id" {
  value = data.azurerm_resource_graph_graph_query.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the Resource Graph exists. Changing this forces a new Resource Graph to be created.

* `resource_name` - (Required) The name of the Graph Query resource. Changing this forces a new Resource Graph to be created.

---

* `description` - (Optional) The description of a graph query.

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Graph.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Graph.

* `name` - The Name of this Resource Graph.

* `result_kind` - Enum indicating a type of graph query.

* `time_modified` - Date and time in UTC of the last modification that was made to this graph query definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Graph.
