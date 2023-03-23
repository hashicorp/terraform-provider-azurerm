---
subcategory: "IoT Central"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_iotcentral_application_network_rule_set"
description: |-
  Manages an IoT Central Application Network Rule Set.
---

# azurerm_iotcentral_application_network_rule_set

Manages an IoT Central Application Network Rule Set.

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

  tags = {
    Foo = "Bar"
  }
}

resource "azurerm_iotcentral_application_network_rule_set" "example" {
  iotcentral_application_id = azurerm_iotcentral_application.example.id

  ip_rule {
    name    = "rule1"
    ip_mask = "10.0.1.0/24"
  }

  ip_rule {
    name    = "rule2"
    ip_mask = "10.1.1.0/24"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `iotcentral_application_id` - (Required) The ID of the IoT Central Application. Changing this forces a new resource to be created.

* `apply_to_device` - (Optional) Whether these IP Rules apply for device connectivity to IoT Hub and Device Provisioning Service associated with this IoT Central Application. Possible values are `true`, `false`. Defaults to `true`

* `default_action` - (Optional) Specifies the default action for the IoT Central Application Network Rule Set. Possible values are `Allow` and `Deny`. Defaults to `Deny`.

* `ip_rule` - (Optional) One or more `ip_rule` blocks as defined below.

---

A `ip_rule` block supports the following:

* `name` - (Required) The name of the IP Rule

* `ip_mask` - (Required) The IP address range in CIDR notation for the IP Rule.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the IoT Central Application Network Rule Set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the IoT Central Application Network Rule Set.
* `read` - (Defaults to 5 minutes) Used when retrieving the IoT Central Application Network Rule Set.
* `update` - (Defaults to 30 minutes) Used when updating the IoT Central Application Network Rule Set.
* `delete` - (Defaults to 30 minutes) Used when deleting the IoT Central Application Network Rule Set.

## Import

IoT Central Application Network Rule Sets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_iotcentral_application_network_rule_set.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.IoTCentral/iotApps/app1
```
