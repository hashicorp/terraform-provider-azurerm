---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager"
description: |-
  Manages a Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager.
---

# azurerm_palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager

Manages a Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager.

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
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
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

resource "azurerm_palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager" "example" {
  name                             = "example-ngfwvh"
  resource_group_name              = azurerm_resource_group.example.name
  location                         = azurerm_resource_group.example.location
  strata_cloud_manager_tenant_name = "example-scm-tenant"

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

* `location` - (Required) The Azure Region where the Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager should exist. Changing this forces a new Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager to be created.

* `name` - (Required) The name which should be used for this Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager. Changing this forces a new Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager to be created.

* `network_profile` - (Required) A `network_profile` block as defined below.

* `resource_group_name` - (Required) The name of the Resource Group where the Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager should exist. Changing this forces a new Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager to be created.

* `strata_cloud_manager_tenant_name` - (Required) Strata Cloud Manager name which is intended to manage the policy for this firewall.

---

* `dns_settings` - (Optional) A `dns_settings` block as defined below.

* `marketplace_offer_id` - (Optional) The marketplace offer ID. Defaults to `pan_swfw_cloud_ngfw`. Changing this forces a new resource to be created.

* `plan_id` - (Optional) The billing plan ID as published by Liftr.PAN. Defaults to `panw-cngfw-payg`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager.

---

A `dns_settings` block supports the following:

* `dns_servers` - (Optional) Specifies a list of DNS servers to use. Conflicts with `dns_settings[0].use_azure_dns`.

* `use_azure_dns` - (Optional) Should the Firewall use Azure Supplied DNS servers. Conflicts with `dns_settings[0].dns_servers`. Defaults to `false`.

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

* `id` - The ID of the Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager.
* `update` - (Defaults to 3 hours) Used when updating the Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager.
* `delete` - (Defaults to 2 hours) Used when deleting the Palo Alto Next Generation Firewall Virtual Network Strata Cloud Manager.

## Import

Palo Alto Next Generation Firewall Virtual Network Strata Cloud Managers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/PaloAltoNetworks.Cloudngfw/firewalls/myVNetStrataCloudManagerFW
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `PaloAltoNetworks.Cloudngfw` - 2025-05-23
