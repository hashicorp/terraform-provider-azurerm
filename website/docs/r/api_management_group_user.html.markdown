---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_group_user"
sidebar_current: "docs-azurerm-resource-api-management-group-user"
description: |-
  Manages an API Management User Assignment to a Group.
---

# azurerm_api_management_group_user

Manages an API Management User Assignment to a Group.


## Example Usage

```hcl
data "azurerm_api_management_user" "example" {
  user_id             = "my-user"
  api_management_name = "example-apim"
  resource_group_name = "search-service"
}

resource "azurerm_api_management_group_user" "example" {
  user_id             = "${data.azurerm_api_management_user.example.id}"
  group_name          = "example-group"
  resource_group_name = "${data.azurerm_api_management_user.example.resource_group_name}"
  api_management_name = "${data.azurerm_api_management_user.example.api_management_name}"
}
```


## Argument Reference

The following arguments are supported:

* `user_id` - (Required) The ID of the API Management User which should be assigned to this API Management Group. Changing this forces a new resource to be created.

* `group_name` - (Required) The Name of the API Management Group within the API Management Service. Changing this forces a new resource to be created.

* `api_management_name` - (Required) The name of the API Management Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the API Management Service exists. Changing this forces a new resource to be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the API Management Group User.

## Import

API Management Group Users can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_group_user.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/service/service1/groups/groupId/users/user123
```
