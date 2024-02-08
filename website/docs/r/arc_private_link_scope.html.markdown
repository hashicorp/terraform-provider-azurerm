---
subcategory: "Hybrid Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_arc_private_link_scope"
description: |-
  Manages an Azure Arc Private Link Scope.
---

# azurerm_arc_private_link_scope

Manages an Azure Arc Private Link Scope.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "west europe"
}

resource "azurerm_arc_private_link_scope" "example" {
  name                = "plsexample"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Arc Private Link Scope should exist. Changing this forces a new Azure Arc Private Link Scope to be created.

* `name` - (Required) The name which should be used for the Azure Arc Private Link Scope. Changing this forces a new Azure Arc Private Link Scope to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Arc Private Link Scope should exist. Changing this forces a new Azure Arc Private Link Scope to be created.

---

* `public_network_access_enabled` - (Optional) Indicates whether machines associated with the private link scope can also use public Azure Arc service endpoints. Defaults to `false`. Possible values are `true` and `false`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Arc Private Link Scope.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Arc Private Link Scope.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Arc Private Link Scope.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Arc Private Link Scope.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Arc Private Link Scope.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Arc Private Link Scope.

## Import

Azure Arc Private Link Scope can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_arc_private_link_scope.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.HybridCompute/privateLinkScopes/privateLinkScope1
```
