---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_slice"
description: |-
  Manages a Mobile Network Slice.
---

# azurerm_mobile_network_slice

Manages a Mobile Network Slice.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mobile_network" "example" {
  name                = "example-mn"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  mobile_country_code = "001"
  mobile_network_code = "01"
}


resource "azurerm_mobile_network_slice" "example" {
  name              = "example-mns"
  mobile_network_id = azurerm_mobile_network.example.id
  location          = azurerm_resource_group.example.location
  description       = "an example slice"

  single_network_slice_selection_assistance_information {
    slice_service_type = 1
  }

  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Slice. Changing this forces a new Mobile Network Slice to be created.

* `mobile_network_id` - (Required) The ID of Mobile Network which the Mobile Network Slice belongs to. Changing this forces a new Mobile Network Slice to be created.

* `location` - (Required) Specifies the Azure Region where the Mobile Network Slice should exist. Changing this forces a new Mobile Network Slice to be created.

* `single_network_slice_selection_assistance_information` - (Required) A `single_network_slice_selection_assistance_information` block as defined below. Single-network slice selection assistance information (S-NSSAI). Unique at the scope of a mobile network.

* `description` - (Optional) A description for this Mobile Network Slice.

* `tags` - (Optional) A mapping of tags which should be assigned to the Mobile Network Slice.

---

A `single_network_slice_selection_assistance_information` block supports the following:

* `slice_differentiator` - (Optional) Slice differentiator (SD). Must be a 6 digit hex string.

* `slice_service_type` - (Required) Slice/service type (SST). Must be between `0` and `255`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Slice.



## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Mobile Network Slice.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Slice.
* `update` - (Defaults to 3 hours) Used when updating the Mobile Network Slice.
* `delete` - (Defaults to 3 hours) Used when deleting the Mobile Network Slice.

## Import

Mobile Network Slice can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_slice.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/mobileNetworks/mobileNetwork1/slices/slice1
```
