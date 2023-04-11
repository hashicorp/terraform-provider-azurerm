---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_custom_ip_prefix"
description: Manages a Custom IP Prefix
---

# azurerm_custom_ip_prefix

Manages a Custom IP prefix.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_custom_ip_prefix" "example" {
  name                = "example-CustomIPPrefix"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  cidr                          = "1.2.3.4/22"
  commissioning_enabled         = true
  roa_validity_end_date         = "2099-12-12"
  wan_validation_signed_message = "signed message for WAN validation"
  zones                         = ["1", "2", "3"]
  tags = {
    env = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Custom IP Prefix. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Custom IP Prefix should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which to create the Custom IP Prefix. Changing this forces a new resource to be created.

* `cidr` - (Required) The `cidr` of the Custom IP Prefix, either IPv4 or IPv6. Changing this forces a new resource to be created.

-> **IPv6 Note** This resource has only been tested with IPv4 and does not yet support IPv6 prefixes.

* `roa_validity_end_date` - (Required) The expiration date of the Route Origin Authorization (ROA) document which has been filed with the Routing Internet Registry (RIR) for this prefix. The expected format is `YYYY-MM-DD`. Changing this forces a new resource to be created.

* `wan_validation_signed_message` - (Required) The signed base64-encoded authorization message, which will be sent to Microsoft for WAN verification. Refer to [Azure documentation](https://learn.microsoft.com/en-us/azure/virtual-network/ip-services/create-custom-ip-address-prefix-cli#certificate-readiness) for more details about the process for your RIR. Changing this forces a new resource to be created.

* `commissioning_enabled` - (Optional) Specifies that the custom IP prefix should be commissioned after provisioning in Azure. Defaults to `false`.

!> **Warning** Changing the value of `commissioning_enabled` from `true` to `false` causes the IP prefix to stop being advertised by Azure and is functionally equivalent to deleting it when used in a production setting.

* `internet_advertising_disabled` - (Optional) Specifies that the custom IP prefix should not be publicly advertised on the Internet when commissioned (regional commissioning feature). Defaults to `false`.

!> **Warning** Changing the value of `internet_advertising_disabled` from `true` to `false` causes the IP prefix to stop being advertised by Azure and is functionally equivalent to deleting it when used in a production setting.

* `tags` - (Optional) A mapping of tags to assign to the Custom IP Prefix.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Custom IP Prefix should be located. Changing this forces a new resource to be created.

-> **Note:** In regions with [availability zones](https://docs.microsoft.com/en-us/azure/availability-zones/az-overview), the Custom IP Prefix must be specified as either `Zone-redundant` or assigned to a specific zone. It can't be created with no zone specified in these regions. All IPs from the prefix must have the same zonal properties.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Custom IP Prefix.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 hours) Used when creating the Custom IP Prefix.
* `update` - (Defaults to 17 hours) Used when updating the Custom IP Prefix.
* `read` - (Defaults to 5 minutes) Used when retrieving the Custom IP Prefix.
* `delete` - (Defaults to 17 hours) Used when deleting the Custom IP Prefix.

## Import

A Custom IP Prefix can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_custom_ip_prefix.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/customIPPrefixes/customIPPrefix1
```
