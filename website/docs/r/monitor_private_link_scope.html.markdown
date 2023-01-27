---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_private_link_scope"
description: |-
  Manages an Azure Monitor Private Link Scope
---

# azurerm_monitor_private_link_scope

Manages an Azure Monitor Private Link Scope.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_private_link_scope" "example" {
  name                = "example-ampls"
  resource_group_name = azurerm_resource_group.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Azure Monitor Private Link Scope. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Monitor Private Link Scope should exist. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Monitor Private Link Scope.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Monitor Private Link Scope.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Monitor Private Link Scope.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Monitor Private Link Scope.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Monitor Private Link Scope.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Monitor Private Link Scope.

## Import

Azure Monitor Private Link Scopes can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_private_link_scope.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Insights/privateLinkScopes/pls1
```
