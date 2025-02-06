---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack"
description: |-
  Manages a Palo Alto Next Generation Firewall VHub Local Rulestack.
---

# azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack

Manages a Palo Alto Next Generation Firewall VHub Local Rulestack.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                = "acceptanceTestPublicIp1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  allocation_method   = "Static"

  tags = {
    environment = "Production"
  }
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-virtualwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-virtualhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.0.0/23"

  tags = {
    "hubSaaSPreview" = "true"
  }
}

resource "azurerm_palo_alto_virtual_network_appliance" "example" {
  name           = "example-appliance"
  virtual_hub_id = azurerm_virtual_hub.example.id
}
resource "azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack" "example" {
  name                = "example-ngfwvn"
  resource_group_name = azurerm_resource_group.example.name
  rulestack_id        = azurerm_palo_alto_local_rulestack.example.id

  network_profile {
    public_ip_address_ids        = [azurerm_public_ip.example.id]
    virtual_hub_id               = azurerm_virtual_hub.example.id
    network_virtual_appliance_id = azurerm_palo_alto_virtual_network_appliance.example.id
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Palo Alto Next Generation Firewall VHub Local Rulestack. Changing this forces a new Palo Alto Next Generation Firewall VHub Local Rulestack to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Palo Alto Next Generation Firewall VHub Local Rulestack should exist. Changing this forces a new Palo Alto Next Generation Firewall VHub Local Rulestack to be created.

* `rulestack_id` - (Required) The ID of the Local Rulestack to be used for this Next Generation Firewall.

* `network_profile` - (Required) A `network_profile` block as defined below.

---

* `destination_nat` - (Optional) One or more `destination_nat` blocks as defined below.

* `dns_settings` - (Optional) A `dns_settings` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Palo Alto Next Generation Firewall VHub Local Rulestack.

---

A `backend_config` block supports the following:

* `port` - (Required) The port number to send traffic to.

* `public_ip_address` - (Required) The Public IP Address to send the traffic to.

---

A `destination_nat` block supports the following:

* `name` - (Required) The name which should be used for this NAT.

* `protocol` - (Required) The protocol used for this Destination NAT. Possible values include `TCP` and `UDP`.

* `backend_config` - (Optional) A `backend_config` block as defined above.

* `frontend_config` - (Optional) A `frontend_config` block as defined below.

---

A `dns_settings` block supports the following:

* `dns_servers` - (Optional) Specifies a list of DNS servers to proxy. Conflicts with `dns_settings[0].use_azure_dns`.

* `use_azure_dns` - (Optional) Should Azure DNS servers be used? Conflicts with `dns_settings[0].dns_servers`. Defaults to `false`.

---

A `frontend_config` block supports the following:

* `port` - (Required) The port on which traffic will be receiveed.

* `public_ip_address_id` - (Required) The ID of the Public IP Address resource the traffic will be received on.

---

A `network_profile` block supports the following:

* `network_virtual_appliance_id` - (Required) The ID of the Palo Alto Network Virtual Appliance in the VHub. Changing this forces a new Palo Alto Next Generation Firewall VHub Local Rulestack to be created.

* `public_ip_address_ids` - (Required) Specifies a list of Public IP IDs to use for this Next Generation Firewall.

* `virtual_hub_id` - (Required) The ID of the Virtual Hub this Next generation Fireall will be deployed in. Changing this forces a new Palo Alto Next Generation Firewall VHub Local Rulestack to be created.

* `egress_nat_ip_address_ids` - (Optional) Specifies a list of Public IP IDs to use for Egress NAT.

* `trusted_address_ranges` - (Optional) Specifies a list of trusted ranges to use for the Network.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Next Generation Firewall VHub Local Rulestack.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Palo Alto Next Generation Firewall VHub Local Rulestack.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Next Generation Firewall VHub Local Rulestack.
* `update` - (Defaults to 3 hours) Used when updating the Palo Alto Next Generation Firewall VHub Local Rulestack.
* `delete` - (Defaults to 2 hours) Used when deleting the Palo Alto Next Generation Firewall VHub Local Rulestack.

## Import

Palo Alto Next Generation Firewall VHub Local Rulestacks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/PaloAltoNetworks.Cloudngfw/firewalls/myVhubRulestackFW
```
