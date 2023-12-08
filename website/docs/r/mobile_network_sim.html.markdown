---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_sim"
description: |-
  Manages a Mobile Network Sim.
---

# azurerm_mobile_network_sim

Manages a Mobile Network Sim.

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

resource "azurerm_mobile_network_sim_group" "example" {
  name              = "example-mnsg"
  location          = azurerm_resource_group.example.location
  mobile_network_id = azurerm_mobile_network.example.id
}

resource "azurerm_mobile_network_slice" "example" {
  name              = "example-slice"
  mobile_network_id = azurerm_mobile_network.example.id
  location          = azurerm_resource_group.example.location

  single_network_slice_selection_assistance_information {
    slice_service_type = 1
  }
}

resource "azurerm_mobile_network_attached_data_network" "example" {
  mobile_network_data_network_name            = azurerm_mobile_network_data_network.example.name
  mobile_network_packet_core_data_plane_id    = azurerm_mobile_network_packet_core_data_plane.example.id
  location                                    = azurerm_resource_group.example.location
  dns_addresses                               = ["1.1.1.1"]
  user_equipment_address_pool_prefixes        = ["2.4.0.0/24"]
  user_equipment_static_address_pool_prefixes = ["2.4.1.0/24"]
  user_plane_access_name                      = "test"
  user_plane_access_ipv4_address              = "10.204.141.4"
  user_plane_access_ipv4_gateway              = "10.204.141.1"
  user_plane_access_ipv4_subnet               = "10.204.141.0/24"
}

resource "azurerm_mobile_network_sim" "example" {
  name                                     = "example-sim"
  mobile_network_sim_group_id              = azurerm_mobile_network_sim_group.example.id
  authentication_key                       = "00000000000000000000000000000000"
  integrated_circuit_card_identifier       = "8900000000000000000"
  international_mobile_subscriber_identity = "000000000000000"
  operator_key_code                        = "00000000000000000000000000000000"

  static_ip_configuration {
    attached_data_network_id = data.azurerm_mobile_network_attached_data_network.test.id
    slice_id                 = azurerm_mobile_network_slice.test.id
    static_ipv4_address      = "2.4.0.1"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Mobile Network Sim. Changing this forces a new Mobile Network Sim to be created.

* `mobile_network_sim_group_id` - (Required) The ID of the Mobile Network which the Mobile Network Sim belongs to. Changing this forces a new Mobile Network Sim to be created.

* `authentication_key` - (Required) The Ki value for the SIM.

* `international_mobile_subscriber_identity` - (Required) The international mobile subscriber identity (IMSI) for the SIM. Changing this forces a new Mobile Network Sim to be created.

* `integrated_circuit_card_identifier` - (Required) The integrated circuit card ID (ICCID) for the SIM. Changing this forces a new Mobile Network Sim to be created.

* `operator_key_code` - (Required) The Opc value for the SIM.

* `device_type` - (Optional) An optional free-form text field that can be used to record the device type this SIM is associated with, for example `Video camera`. The Azure portal allows SIMs to be grouped and filtered based on this value.

* `sim_policy_id` - (Optional) The ID of SIM policy used by this SIM.

* `static_ip_configuration` - (Optional) A `static_ip_configuration` block as defined below.

---

A `static_ip_configuration` block supports the following:

* `attached_data_network_id` - (Required) The ID of attached data network on which the static IP address will be used. The combination of attached data network and slice defines the network scope of the IP address.

* `slice_id` - (Required) The ID of network slice on which the static IP address will be used. The combination of attached data network and slice defines the network scope of the IP address.

* `static_ipv4_address` - (Optional) The IPv4 address assigned to the SIM at this network scope. This address must be in the userEquipmentStaticAddressPoolPrefix defined in the attached data network.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Sim.

* `sim_state` - The state of the SIM resource.

* `vendor_key_fingerprint` - The public key fingerprint of the SIM vendor who provided this SIM, if any.

* `vendor_name` - The name of the SIM vendor who provided this SIM, if any.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Mobile Network Sim.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Sim.
* `update` - (Defaults to 90 minutes) Used when updating the Mobile Network Sim.
* `delete` - (Defaults to 90 minutes) Used when deleting the Mobile Network Sim.

## Import

Mobile Network Sim can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_sim.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/simGroups/simGroup1/sims/sim1
```
