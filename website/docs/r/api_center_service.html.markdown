---
subcategory: "ApiCenter"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_center_service"
description: |-
  Manages an API Center Service.
---

# azurerm_api_center_service

Manages an API Center Service.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_center_service" "example" {
  name                = "apicenter-example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this API Center Service. Changing this forces a new API Center Service to be created.

* `resource_group_name` - (Required) The name of the Resource Group where this API Center Service should exist. Changing this forces a new API Center Service to be created.

* `location` - (Required) The Azure Region where this API Center Service should exist. Changing this forces a new API Center Service to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the API Center Service.

---

An `identity` block supports the following:

* `type` - (Required) The type of managed identity to assign. Possible values are `SystemAssigned`, `UserAssigned`, and `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) - A list of one or more Resource IDs for User Assigned Managed identities to assign. Required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Center Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Center Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Center Service.
* `update` - (Defaults to 30 minutes) Used when updating the API Center Service.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Center Service.

## Import

API Center Service can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_center_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.ApiCenter/services/example/center
```
