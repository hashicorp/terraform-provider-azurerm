---
subcategory: "Resource Graph"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_graph_graph_query"
description: |-
  Manages a Resource Graph.
---

# azurerm_resource_graph_graph_query

Manages a Resource Graph.

## Example Usage

```hcl
resource "azurerm_resource_graph_graph_query" "example" {
  resource_group_name = "example"
  resource_name       = "example"
  query               = "where isnotnull(tags['Prod']) and properties.extensions[0].Name == 'docker'"
}
```

## Arguments Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the Resource Group where the Resource Graph should exist. Changing this forces a new Resource Graph to be created.

* `resource_name` - (Required) The name of the Graph Query resource. Changing this forces a new Resource Graph to be created.

* `query` - (Required) The query content of the Graph Query resource.

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

* `create` - (Defaults to 30 minutes) Used when creating the Resource Graph.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Graph.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Graph.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Graph.

## Import

Resource Graphs can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_graph_graph_query.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ResourceGraph/queries/resource1
```
