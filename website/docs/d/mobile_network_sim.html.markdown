---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_sim"
description: |-
  Get information about a Mobile Network Sim.
---

# azurerm_mobile_network_sim

Get information about a Mobile Network Sim.

## Example Usage

```hcl
data "azurerm_mobile_network_sim_group" "example" {
  name                = "example-mnsg"
  resource_group_name = "example-rg"
}


data "azurerm_mobile_network_sim" "example" {
  name                        = "example-sim"
  mobile_network_sim_group_id = data.azurerm_mobile_network_sim_group.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name which should be used for this Mobile Network Sim. 

* `mobile_network_sim_group_id` - The ID of the Mobile Network which the Mobile Network Sim belongs to. 

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Sim.

* `sim_state` - The state of the SIM resource.

* `vendor_key_fingerprint` - The public key fingerprint of the SIM vendor who provided this SIM.

* `vendor_name` - The name of the SIM vendor who provided this SIM, if any.

* `international_mobile_subscriber_identity` - The international mobile subscriber identity (IMSI) for the SIM.

* `integrated_circuit_card_identifier` - The integrated circuit card ID (ICCID) for the SIM.

* `device_type` -  The device type this SIM is associated with.

* `sim_policy_id` - The ID of SIM policy used by this SIM.

* `static_ip_configuration` - A `static_ip_configuration` block as defined below.

---

A `static_ip_configuration` block supports the following:

* `attached_data_network_id` - The ID of attached data network on which the static.

* `slice` - The ID of network slice on which the static IP address will be used. 

* `static_ipv4_address` - The IPv4 address assigned to the SIM at this network scope.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Sim.

## Import

Mobile Network Sim can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_sim.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/simGroups/simGroup1/sims/sim1
```
