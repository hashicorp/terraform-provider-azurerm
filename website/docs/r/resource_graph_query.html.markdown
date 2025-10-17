---
subcategory: "Resource Graph"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_graph_query"
description: |-
  Manages a Resource Graph Query.
---

# azurerm_resource_graph_query

Manages a Resource Graph Query Example.

## Example Usage

```hcl
resource "azurerm_resource_graph_query" "example" {
  resource_group_name = "example"
  location            = "West Europe"
  query               = "\nresources\n| summarize count() by type, location"
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Resource Graph Query should exist. Changing this forces a new Resource Graph Query to be created.

* `query` - (Required) The query for the Resource Graph Query.

* `resource_group_name` - (Required) The name of the Resource Group where the Resource Graph Query should exist. Changing this forces a new Resource Graph Query to be created.

---

* `description` - (Optional) The description of the Resource Graph Query.

* `name` - (Optional) The name which should be used for this Resource Graph Query Example. Changing this forces a new Resource Graph Query to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Graph Query Example.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Graph Query Example.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Graph Query Example.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Graph Query Example.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Graph Query Example.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Graph Query Example.

## Import

Resource Graph Query Examples can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_graph_query.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ResourceGraph/queries/query1
```
