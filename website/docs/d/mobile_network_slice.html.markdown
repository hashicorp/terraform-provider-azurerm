---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_slice"
description: |-
  Get information about a Mobile Network Slice.
---

# azurerm_mobile_network_slice

Get information about a Mobile Network Slice.

## Example Usage

```hcl
data "azurerm_mobile_network" "example" {
  name                = "example-mn"
  resource_group_name = "example-rg"
}

data "azurerm_mobile_network_slice" "example" {
  name              = "example-mns"
  mobile_network_id = data.azurerm_mobile_network.test.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Slice. 

* `mobile_network_id` - (Required) The ID of Mobile Network which the Mobile Network Slice belongs to.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Slice.

* `location` - The Azure Region where the Mobile Network Slice exists. 

* `single_network_slice_selection_assistance_information` - A `single_network_slice_selection_assistance_information` block as defined below. Single-network slice selection assistance information (S-NSSAI). 

* `description` - A description of this Mobile Network Slice.

* `tags` - A mapping of tags which are assigned to the Mobile Network Slice.

---

A `single_network_slice_selection_assistance_information` block supports the following:

* `slice_differentiator` - Slice differentiator (SD).

* `slice_service_type` - Slice/service type (SST).


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Slice.

