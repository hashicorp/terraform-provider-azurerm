---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_user"
description: |-
  Manages an API Management User.
---

# azurerm_api_management_user

Manages an API Management User.

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
  publisher_name      = "My Company"
  publisher_email     = "company@terraform.io"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_user" "example" {
  user_id             = "5931a75ae4bbd512288c680b"
  api_management_name = azurerm_api_management.example.name
  resource_group_name = azurerm_resource_group.example.name
  first_name          = "Example"
  last_name           = "User"
  email               = "tom+tfdev@hashicorp.com"
  state               = "active"
}
```

## Argument Reference

The following arguments are supported:

* `api_management_name` - (Required) The name of the API Management Service in which the User should be created. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

* `email` - (Required) The email address associated with this user.

* `first_name` - (Required) The first name for this user.

* `last_name` - (Required) The last name for this user.

* `user_id` - (Required) The Identifier for this User, which must be unique within the API Management Service. Changing this forces a new resource to be created.

---

* `confirmation` - (Optional) The kind of confirmation email which will be sent to this user. Possible values are `invite` and `signup`. Changing this forces a new resource to be created.

* `note` - (Optional) A note about this user.

* `password` - (Optional) The password associated with this user.

* `state` - (Optional) The state of this user. Possible values are `active`, `blocked` and `pending`.

-> **Note:** the State can be changed from Pending -> Active/Blocked but not from Active/Blocked -> Pending.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management User.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 45 minutes) Used when creating the API Management User.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management User.
* `update` - (Defaults to 45 minutes) Used when updating the API Management User.
* `delete` - (Defaults to 45 minutes) Used when deleting the API Management User.

## Import

API Management Users can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_user.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.ApiManagement/service/instance1/users/abc123
```
