---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_next_generation_firewall_virtual_network_panorama"
description: |-
  Manages a Palo Alto Next Generation Firewall Virtual Network Panorama.
---

# azurerm_palo_alto_next_generation_firewall_virtual_network_panorama

Manages a Palo Alto Next Generation Firewall Virtual Network Panorama.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "westeurope"
}

resource "azurerm_public_ip" "example" {
  name                = "example-public-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_network_security_group" "example" {
  name                = "example-nsg"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    environment = "Production"
  }
}

resource "azurerm_subnet" "trust" {
  name                 = "example-trust-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "trusted"

    service_delegation {
      name = "PaloAltoNetworks.Cloudngfw/firewalls"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "trust" {
  subnet_id                 = azurerm_subnet.trust.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_subnet" "untrust" {
  name                 = "example-untrust-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "untrusted"

    service_delegation {
      name = "PaloAltoNetworks.Cloudngfw/firewalls"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
      ]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "untrust" {
  subnet_id                 = azurerm_subnet.untrust.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_panorama" "example" {
  name                   = "example-ngfwvh"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  panorama_base64_config = "e2RnbmFtZTogY25nZnctYXotZXhhbXBsZSwgdHBsbmFtZTogY25nZnctZXhhbXBsZS10ZW1wbGF0ZS1zdGFjaywgZXhhbXBsZS1wYW5vcmFtYS1zZXJ2ZXI6IDE5Mi4xNjguMC4xLCB2bS1hdXRoLWtleTogMDAwMDAwMDAwMDAwMDAwLCBleHBpcnk6IDIwMjQvMDcvMzF9Cg=="

  network_profile {
    public_ip_address_ids = [azurerm_public_ip.example.id]

    vnet_configuration {
      virtual_network_id  = azurerm_virtual_network.example.id
      trusted_subnet_id   = azurerm_subnet.trust.id
      untrusted_subnet_id = azurerm_subnet.untrust.id
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Palo Alto Next Generation Firewall Virtual Network Panorama should exist. Changing this forces a new Palo Alto Next Generation Firewall Virtual Network Panorama to be created.

* `name` - (Required) The name which should be used for this Palo Alto Next Generation Firewall Virtual Network Panorama. Changing this forces a new Palo Alto Next Generation Firewall Virtual Network Panorama to be created.

* `network_profile` - (Required) A `network_profile` block as defined below.

* `panorama_base64_config` - (Required) The base64 encoded configuration registration string as defined by your Panorama Server for your Cloud Device Group.

* `resource_group_name` - (Required) The name of the Resource Group where the Palo Alto Next Generation Firewall Virtual Network Panorama should exist. Changing this forces a new Palo Alto Next Generation Firewall Virtual Network Panorama to be created.

* `marketplace_offer_id` - (Optional) The marketplace offer ID. Defaults to `pan_swfw_cloud_ngfw`. Changing this forces a new resource to be created.

* `plan_id` - (Optional) The billing plan ID as published by Liftr.PAN. Defaults to `panw-cloud-ngfw-payg`.

~> **Note:** The former `plan_id` `panw-cloud-ngfw-payg` is defined as stop sell, but has been set as the default to not break any existing resources that were originally provisioned with it. Users need to explicitly set `plan_id` to `panw-cngfw-payg` when creating new resources.

---

* `destination_nat` - (Optional) One or more `destination_nat` blocks as defined below.

* `dns_settings` - (Optional) A `dns_settings` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Palo Alto Next Generation Firewall Virtual Network Panorama.

---

A `backend_config` block supports the following:

* `port` - (Required) The port number to send traffic to.

* `public_ip_address` - (Required) The IP Address to send the traffic to.

---

A `destination_nat` block supports the following:

* `name` - (Required) The name which should be used for this Destination NAT.

* `protocol` - (Required) The Protocol for this Destination NAT configuration. Possible values include `TCP` and `UDP`.

* `backend_config` - (Optional) A `backend_config` block as defined above.

* `frontend_config` - (Optional) A `frontend_config` block as defined below.

---

A `dns_settings` block supports the following:

* `dns_servers` - (Optional) Specifies a list of DNS servers to use. Conflicts with `dns_settings[0].use_azure_dns`.

* `use_azure_dns` - (Optional) Should the Firewall use Azure Supplied DNS servers. Conflicts with `dns_settings[0].dns_servers`. Defaults to `false`.

---

A `frontend_config` block supports the following:

* `port` - (Required) The port on which to receive traffic.

* `public_ip_address_id` - (Required) The ID of the Public IP Address on which to receive traffic.

~> **Note:** This must be an Azure Public IP address ID also specified in the `public_ip_address_ids` list.

---

A `network_profile` block supports the following:

* `public_ip_address_ids` - (Required) Specifies a list of Azure Public IP Address IDs.

* `vnet_configuration` - (Required) A `vnet_configuration` block as defined below.

* `egress_nat_ip_address_ids` - (Optional) Specifies a list of Azure Public IP Address IDs that can be used for Egress (Source) Network Address Translation.

* `trusted_address_ranges` - (Optional) Specifies a list of trusted ranges to use for the Network.

---

A `vnet_configuration` block supports the following:

* `virtual_network_id` - (Required) The ID of the Virtual Network.

* `trusted_subnet_id` - (Optional) The ID of the Trust subnet.

* `untrusted_subnet_id` - (Optional) The ID of the UnTrust subnet.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Next Generation Firewall Virtual Network Panorama.

* `panorama` - A `panorama` block as defined below.

---

A `panorama` block exports the following:

* `device_group_name` - The Device Group Name to which this Firewall Resource is registered.

* `host_name` - The Host Name of this Firewall Resource.

* `name` - The name of this Firewall Resource.

* `panorama_server_1` - The name of the First Panorana server.

* `panorama_server_2` - The name of the Second Panorana server.

* `template_name` - The name of the Panorama Template applied to this Firewall Resource.

* `virtual_machine_ssh_key` - The SSH Key to connect to the Firewall Resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Palo Alto Next Generation Firewall Virtual Network Panorama.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Next Generation Firewall Virtual Network Panorama.
* `update` - (Defaults to 3 hours) Used when updating the Palo Alto Next Generation Firewall Virtual Network Panorama.
* `delete` - (Defaults to 3 hours) Used when deleting the Palo Alto Next Generation Firewall Virtual Network Panorama.

## Import

Palo Alto Next Generation Firewall Virtual Network Panoramas can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_next_generation_firewall_virtual_network_panorama.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/PaloAltoNetworks.Cloudngfw/firewalls/myVNetPanoramaFW
```
