---
subcategory: "IoT Central"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotcentral_group_user"
description: |-
  Manages a IoT Central Group User.
---

# azurerm_iotcentral_group_user

Manages a IoT Central Group User.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "example-resource"
  location = "West Europe"
}

resource "azurerm_iotcentral_application" "example" {
  name                = "example-iotcentral-app"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sub_domain          = "example-iotcentral-app-subdomain"
  display_name        = "example-iotcentral-app-display-name"
  sku                 = "ST1"
  template            = "iotc-default@1.0.0"
  tags = {
    Foo = "Bar"
  }
}

resource "azuread_group" "example" {
  display_name     = "Example AD Group"
  security_enabled = true
}

resource "azurerm_iotcentral_group_user" "example" {
  iotcentral_application_id = azurerm_iotcentral_application.example.id
  user_id                   = "example-user-id"
  object_id                 = azuread_group.example.object_id
  tenant_id                 = data.azurerm_client_config.current.tenant_id

  role {
    role_id = "344138e9-8de4-4497-8c54-5237e96d6aaf"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `iotcentral_application_id` - (Required) The application `id`. Changing this forces a new IoT Central Group User to be created.

* `user_id` - (Required) The ID of the user. Changing this forces a new IoT Central Group User to be created.

---

* `object_id` - (Required) The object_id of the Group. Changing this forces a new IoT Central Group User to be created.

* `role` - (Required) One or more `role` blocks as defined below.

* `tenant_id` - (Required) The tenant_id of the Group. Changing this forces a new IoT Central Group User to be created.

---

A `role` block supports the following:

* `role_id` - (Required) The ID of the role for this role assignment.

* `organization_id` - (Optional) The ID of the organization for this role assignment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID reference of the user, formated as `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroup}/providers/Microsoft.IoTCentral/iotApps/{application}/users/{userId}`

* `type` - The type of user.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Central Group User.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Central Group User.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Central Group User.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Central Group User.

## Import

IoT Central Group Users can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotcentral_group_user.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.IoTCentral/iotApps/application1/users/user1
```