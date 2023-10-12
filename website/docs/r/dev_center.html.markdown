---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_dev_center"
description: |-
  Manages a Dev Center resource.
---

# azurerm_dev_center

Manages a Dev Center.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

resource "azurerm_dev_center" "example" {
  name                = "example-devcenter"
  resource_group_name = azurerm_resource_group.example.name
  location            = "westeurope"
  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Dev Center. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Dev Center exists. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Dev Center exists. Changing this forces a new resource to be created.

* `identity` - (Optional) A `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `identity` block supports the following:

* `type` - (Required) The type of identity used for the Dev Center. Possible values are `SystemAssigned` `UserAssigned` and `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) A list of User Assigned Identity IDs to be associated with the Dev Center.

* `principal_id` - (Optional) The Principal ID of the Dev Center.

* `tenant_id` - (Optional) The Tenant ID of the Dev Center.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dev Center.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center.
* `update` - (Defaults to 30 minutes) Used when updating the Dev Center.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dev Center.

## Import

Dev Center can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_dev_center.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevCenter/devCenters/devCenter1
```
