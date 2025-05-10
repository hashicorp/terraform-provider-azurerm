---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network"
description: |-
  Manages an Azure Mobile Network.
---

# azurerm_mobile_network

Manages a Mobile Network.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "east us"
}

resource "azurerm_mobile_network" "example" {
  name                = "example-mn"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  mobile_country_code = "001"
  mobile_network_code = "01"

  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network. Changing this forces a new Mobile Network to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Mobile Network should exist. Changing this forces a new Mobile Network to be created.

* `location` - (Required) Specifies the Azure Region where the Mobile Network should exist. Changing this forces a new Mobile Network to be created. The possible values are `eastus` and `northeurope`.

* `mobile_country_code` - (Required) Mobile country code (MCC), defined in https://www.itu.int/rec/T-REC-E.212 . Changing this forces a new resource to be created.

* `mobile_network_code` - (Required) Mobile network code (MNC), defined in https://www.itu.int/rec/T-REC-E.212 . Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Mobile Network.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network.

* `service_key` - The mobile network resource identifier.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Mobile Network.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network.
* `update` - (Defaults to 3 hours) Used when updating the Mobile Network.
* `delete` - (Defaults to 3 hours) Used when deleting the Mobile Network.

## Import

Mobile Network can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/mobileNetworks/mobileNetwork1
```
