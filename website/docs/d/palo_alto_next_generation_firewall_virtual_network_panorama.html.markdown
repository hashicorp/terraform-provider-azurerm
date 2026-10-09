---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_palo_alto_next_generation_firewall_virtual_network_panorama"
description: |-
  Gets information about an existing Palo Alto Next Generation Firewall Virtual Network Panorama.
---

# Data Source: azurerm_palo_alto_next_generation_firewall_virtual_network_panorama

Use this data source to access information about an existing Palo Alto Next Generation Firewall Virtual Network Panorama.

## Example Usage

```hcl
data "azurerm_palo_alto_next_generation_firewall_virtual_network_panorama" "example" {
  name                = "existing-firewall"
  resource_group_name = "existing-resource-group"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Palo Alto Next Generation Firewall Virtual Network Panorama.

* `resource_group_name` - (Required) The name of the Resource Group where the Palo Alto Next Generation Firewall Virtual Network Panorama exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Palo Alto Next Generation Firewall Virtual Network Panorama.

* `destination_nat` - One or more `destination_nat` blocks as defined below.

* `dns_settings` - A `dns_settings` block as defined below.

* `location` - The Azure Region where the Firewall exists.

* `marketplace_offer_id` - The marketplace offer ID.

* `network_profile` - A `network_profile` block as defined below.

* `panorama` - A `panorama` block as defined below.

* `panorama_base64_config` - The Base64 encoded Panorama configuration registration string.

* `plan_id` - The billing plan ID.

* `tags` - A mapping of tags assigned to the Firewall.

---

A `backend_config` block exports the following:

* `port` - The port number to which traffic is sent.

* `public_ip_address` - The IP Address to which traffic is sent.

---

A `destination_nat` block exports the following:

* `backend_config` - A `backend_config` block as defined above.

* `frontend_config` - A `frontend_config` block as defined below.

* `name` - The name of this Destination NAT configuration.

* `protocol` - The protocol used by this Destination NAT configuration.

---

A `dns_settings` block exports the following:

* `azure_dns_servers` - The Azure DNS Servers used by the Firewall.

* `dns_servers` - The custom DNS Servers used by the Firewall.

* `use_azure_dns` - Whether the Firewall uses Azure DNS Servers.

---

A `frontend_config` block exports the following:

* `port` - The port on which traffic is received.

* `public_ip_address_id` - The ID of the Public IP Address on which traffic is received.

---

A `network_profile` block exports the following:

* `egress_nat_ip_address_ids` - The IDs of the Public IP Addresses used for Egress NAT.

* `egress_nat_ip_addresses` - The Public IP Addresses used for Egress NAT.

* `public_ip_address_ids` - The IDs of the Public IP Addresses used by the Firewall.

* `public_ip_addresses` - The Public IP Addresses used by the Firewall.

* `trusted_address_ranges` - The trusted address ranges used by the Firewall.

* `vnet_configuration` - A `vnet_configuration` block as defined below.

---

A `panorama` block exports the following:

* `device_group_name` - The Device Group Name to which the Firewall is registered.

* `host_name` - The Host Name of the Firewall.

* `name` - The name of the Firewall.

* `panorama_server_1` - The name of the first Panorama Server.

* `panorama_server_2` - The name of the second Panorama Server.

* `template_name` - The name of the Panorama Template applied to the Firewall.

* `virtual_machine_ssh_key` - The SSH Key used to connect to the Firewall.

---

A `vnet_configuration` block exports the following:

* `ip_of_trust_for_user_defined_routes` - The trusted IP Address used for user-defined routes.

* `trusted_subnet_id` - The ID of the trusted subnet.

* `untrusted_subnet_id` - The ID of the untrusted subnet.

* `virtual_network_id` - The ID of the Virtual Network in which the Firewall is deployed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Next Generation Firewall Virtual Network Panorama.
