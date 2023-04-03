---
subcategory: "Voice Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_voice_services_communications_gateway_test_line"
description: |-
  Manages a Voice Services Communications Gateway Test Line.
---

# azurerm_voice_services_communications_gateway_test_line

Manages a Voice Services Communications Gateway Test Line.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Central US"
}

resource "azurerm_voice_services_communications_gateway" "example" {
  name                = "example-vcg"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_voice_services_communications_gateway_test_line" "example" {
  name                                     = "example-vtl"
  location                                 = "West Central US"
  voice_services_communications_gateway_id = azurerm_voice_services_communications_gateway.example.id
  phone_number                             = "123456789"
  purpose                                  = "Automated"
  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Voice Services Communications Gateway Test Line. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Voice Services Communications Gateway Test Line should exist. Changing this forces a new resource to be created.

* `voice_services_communications_gateway_id` - (Required) Specifies the ID of the Voice Services Communications Gateway. Changing this forces a new resource to be created.

* `phone_number` - (Required) Specifies the phone number.

* `purpose` - (Required) The purpose of the Voice Services Communications Gateway Test Line. Possible values are `Automated` or `Manual`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Voice Services Communications Gateway Test Line.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Voice Services Communications Gateway Test Line.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Voice Services Communications Gateway Test Line.
* `read` - (Defaults to 5 minutes) Used when retrieving the Voice Services Communications Gateway Test Line.
* `update` - (Defaults to 30 minutes) Used when updating the Voice Services Communications Gateway Test Line.
* `delete` - (Defaults to 30 minutes) Used when deleting the Voice Services Communications Gateway Test Line.

## Import

Voice Services Communications Gateway Test Line can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_voice_services_communications_gateway_test_line.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.VoiceServices/communicationsGateways/communicationsGateway1/testLines/testLine1
```
