---
subcategory: "IoT Central"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotcentral_ad_group_user"
description: |-
  Manages an IotCentral AD Group User
---

# azurerm_iotcentral_ad_group_user

Manages an IoT Central AD Group User

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource"
  location = "West Europe"
}

resource "azurerm_iotcentral_application" "example" {
  name                = "example-iotcentral-app"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sub_domain          = "example-iotcentral-app-subdomain"

  display_name = "example-iotcentral-app-display-name"
  sku          = "ST1"
  template     = "iotc-default@1.0.0"

  tags = {
    Foo = "Bar"
  }
}

data "azurerm_client_config" "current" {}

resource "azuread_group" "example" {
  display_name     = "example"
  security_enabled = true
}

data "azurerm_iotcentral_role" "app_admin" {
  sub_domain   = azurerm_iotcentral_application.example.sub_domain
  display_name = "App Administrator"
}

resource "azurerm_iotcentral_ad_group_user" "example" {
  sub_domain = azurerm_iotcentral_application.example.sub_domain
  object_id  = azuread_group.example.object_id
  tenant_id  = data.azurerm_client_config.current.tenant_id

  roles {
    role = data.azurerm_iotcentral_role.app_admin.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `sub_domain` - (Required) The application `sub_domain`. Changing this forces a new resource to be created.

* `object_id` - (Required) The AAD `object_id` of the AD Group. Changing this forces a new resource to be created.

* `tenant_id` - (Required) The AAD `tenant_id` of the AD Group. Changing this forces a new resource to be created.

* `roles` - (Required) A `roles` block as defined below that specify the permissions to access the application.

---

The `roles` block supports the following:

* `role` - (Required) The `id` of the role for this role assignment.

* `organization` - (Optional) The `organization_id` of the organization for this role assignment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID reference of the user, formated as `{subdomain}.{baseDomain}/api/users/{userId}`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Central Application.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Central Application.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Central Application.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Central Application.

## Import

The IoT Central AD Group User can be imported using the `id`, e.g.

```shell
terraform import azurerm_iotcentral_ad_group_user.example {subdomain}.{baseDomain}/api/users/{userId}
```
