---
subcategory: "IoT Central"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotcentral_application"
description: |-
  Manages an IotCentral Application
---

# azurerm_iotcentral_application

Manages an IoT Central Application

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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the IotHub resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the IotHub resource has to be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource has to be create. Changing this forces a new resource to be created.

* `sub_domain` - (Required) A `sub_domain` name. Subdomain for the IoT Central URL. Each application must have a unique subdomain.

* `display_name` - (Optional) A `display_name` name. Custom display name for the IoT Central application. Default is resource name.

~> **Note:** Due to a bug in the provider, the default value of `display_name` of a newly created IoT Central App will be the Resource Group Name, it will be fixed and use resource name in 4.0. For an existing IoT Central App, this could be fixed by specifying the `display_name` explicitly.

* `identity` - (Optional) An `identity` block as defined below.

* `public_network_access_enabled` - (Optional) Whether public network access is allowed for the IoT Central Application. Defaults to `true`.

* `sku` - (Optional) A `sku` name. Possible values is `ST0`, `ST1`, `ST2`, Default value is `ST1`

* `template` - (Optional) A `template` name. IoT Central application template name. Defaults to `iotc-pnp-preview@1.0.0`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this IoT Central Application. The only possible value is `SystemAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoT Central Application.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Central Application.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Central Application.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Central Application.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Central Application.

## Import

The IoT Central Application can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotcentral_application.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.IoTCentral/iotApps/app1
```
