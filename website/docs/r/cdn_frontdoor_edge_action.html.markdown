---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_edge_action"
description: |-
  Manages a Front Door (standard/premium) Edge Action.
---

# azurerm_cdn_frontdoor_edge_action

Manages a Front Door (standard/premium) Edge Action.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_edge_action" "example" {
  name                = "exampleedgeaction"
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    environment = "Production"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Edge Action. Possible values must be between 1 and 50 characters in length and may only contain letters and numbers. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Edge Action should exist. Changing this forces a new resource to be created.

---

* `sku` - (Optional) The SKU of the Edge Action. The only possible value is `Standard`. Defaults to `Standard`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Edge Action.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Edge Action.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Edge Action.
* `read` - (Defaults to 5 minutes) Used when retrieving the Edge Action.
* `update` - (Defaults to 30 minutes) Used when updating the Edge Action.
* `delete` - (Defaults to 30 minutes) Used when deleting the Edge Action.

## Import

Edge Actions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_edge_action.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/edgeActions/edgeAction1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Cdn` - 2026-10-01
