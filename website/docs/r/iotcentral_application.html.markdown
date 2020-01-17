---
subcategory: "IoT Central"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotcentral_application"
description: |-
  Manages an IotCentral Application
---

# azurerm_iotcentral_application

Manages an IotCentral Application

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_iotcentral_application" "example" {
  name                = "iocentral-app1"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  sub_domain          = "iocentral-app1-subdomian"
  display_name        = "iocentral-app1-display-name"
  sku                 = "S1"
  template            = "iotc-default@1.0.0" 		

  tags = {
    purpose = "testing"
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

* `sku` - (Optional) A `sku` name. Possible values is `S1`, Default value is `S1`

* `template` - (Optional) A `template` name. IoT Central application template name. Default is a custom application.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IoTCentralApplication.

## Import

IoTCentralApplication can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotcentral_application.app1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.IoTCentral/IoTApps/app1
```