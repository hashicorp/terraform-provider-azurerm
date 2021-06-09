---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_group"
description: |-
  Manages an API Management Group.
---

# azurerm_api_management_group

Manages an API Management Group.


## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_api_management" "example" {
  name                = "example-apim"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_group" "example" {
  name                = "example-apimg"
  resource_group_name = azurerm_resource_group.example.name
  api_management_name = azurerm_api_management.example.name
  display_name        = "Example Group"
  description         = "This is an example API management group."
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the API Management Group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Group should exist. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the [API Management Service](api_management.html) in which the API Management Group should exist. Changing this forces a new resource to be created.

* `display_name` - (Required) The display name of this API Management Group.

* `description` - (Optional) The description of this API Management Group.

* `external_id` - (Optional) The identifier of the external Group. For example, an Azure Active Directory group `aad://<tenant>.onmicrosoft.com/groups/<group object id>`.

* `type` - (Optional) The type of this API Management Group. Possible values are `custom` and `external`. Default is `custom`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Group.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Group.

## Import

API Management Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-resources/providers/Microsoft.ApiManagement/service/example-apim/groups/example-apimg
```
