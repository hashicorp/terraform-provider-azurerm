---
subcategory: "Orbital"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_orbital_contact"
description: |-
  Manages an orbital contact resource.
---

# azurerm_orbital_contact

Manages an orbital contact.

~> **Note:** The `azurerm_orbital_contact` resource has been deprecated and will be removed in v5.0 of the AzureRM Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

resource "azurerm_orbital_spacecraft" "example" {
  name                = "example-spacecraft"
  resource_group_name = azurerm_resource_group.example.name
  location            = "westeurope"
  norad_id            = "12345"

  links {
    bandwidth_mhz        = 100
    center_frequency_mhz = 101
    direction            = "Uplink"
    polarization         = "LHCP"
    name                 = "examplename"
  }

  two_line_elements = ["1 23455U 94089A   97320.90946019  .00000140  00000-0  10191-3 0  2621", "2 23455  99.0090 272.6745 0008546 223.1686 136.8816 14.11711747148495"]
  title_line        = "AQUA"

  tags = {
    aks-managed-cluster-name = "9a57225d-a405-4d40-aa46-f13d2342abef"
  }
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "orbitalgateway"

    service_delegation {
      name = "Microsoft.Orbital/orbitalGateways"
      actions = [
        "Microsoft.Network/publicIPAddresses/join/action",
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/read",
        "Microsoft.Network/publicIPAddresses/read",
      ]
    }
  }
}

resource "azurerm_orbital_contact_profile" "example" {
  name                              = "example-contactprofile"
  resource_group_name               = azurerm_resource_group.example.name
  location                          = azurerm_resource_group.example.location
  minimum_variable_contact_duration = "PT1M"
  auto_tracking                     = "disabled"
  links {
    channels {
      name                 = "channelname"
      bandwidth_mhz        = 100
      center_frequency_mhz = 101
      end_point {
        end_point_name = "AQUA_command"
        ip_address     = "10.0.1.0"
        port           = "49153"
        protocol       = "TCP"
      }
    }
    direction    = "Uplink"
    name         = "RHCP_UL"
    polarization = "RHCP"
  }
  network_configuration_subnet_id = azurerm_subnet.example.id
}

resource "azurerm_orbital_contact" "example" {
  name                   = "example-contact"
  spacecraft_id          = azurerm_orbital_spacecraft.example.id
  reservation_start_time = "2020-07-16T20:35:00.00Z"
  reservation_end_time   = "2020-07-16T20:55:00.00Z"
  ground_station_name    = "WESTUS2_0"
  contact_profile_id     = azurerm_orbital_contact_profile.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Contact. Changing this forces a new resource to be created.

* `spacecraft_id` - (Required) The ID of the spacecraft which the contact will be made to. Changing this forces a new resource to be created.

* `reservation_start_time` - (Required) Reservation start time of the Contact. Changing this forces a new resource to be created.

* `reservation_end_time` - (Required) Reservation end time of the Contact. Changing this forces a new resource to be created.

* `ground_station_name` - (Required) Name of the Azure ground station. Changing this forces a new resource to be created.

* `contact_profile_id` - (Required) ID of the orbital contact profile. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Contact.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Contact.
* `read` - (Defaults to 5 minutes) Used when retrieving the Contact.
* `delete` - (Defaults to 30 minutes) Used when deleting the Contact.

## Import

Spacecraft can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_orbital_contact.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Orbital/spacecrafts/spacecraft1/contacts/contact1
```
