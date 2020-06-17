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
  sku          = "S1"
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

* `sku` - (Optional) A `sku` name. Possible values is `ST1`, `ST2`, Default value is `ST1`

* `template` - (Optional) A `template` name. IoT Central application template name. Default is a custom application.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoT Central Application.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Central Application.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Central Application.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Central Application.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Central Application.

## Import

The IoT Central Application can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotcentral_application.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.IoTCentral/IoTApps/app1
```
