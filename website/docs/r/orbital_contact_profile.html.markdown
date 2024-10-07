---
subcategory: "Orbital"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_orbital_contact_profile"
description: |-
  Manages a orbital contact profile resource.
---

# azurerm_orbital_contact_profile

Manages a Contact profile.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "testvnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "testsubnet"
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
  name                              = "example-contact-profile"
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
        port           = "49513"
        protocol       = "TCP"
      }
    }
    direction    = "Uplink"
    name         = "RHCP_UL"
    polarization = "RHCP"
  }

  network_configuration_subnet_id = azurerm_subnet.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the contact profile. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the contact profile exists. Changing this forces a new resource to be created.

* `location` - (Required) The location where the contact profile exists. Changing this forces a new resource to be created.

* `minimum_variable_contact_duration` - (Required) Minimum viable contact duration in ISO 8601 format. Used for listing the available contacts with a spacecraft at a given ground station.

* `auto_tracking` - (Required) Auto-tracking configurations for a spacecraft. Possible values are `disabled`, `xBand` and `sBand`.

* `network_configuration_subnet_id` - (Required) ARM resource identifier of the subnet delegated to the Microsoft.Orbital/orbitalGateways. Needs to be at least a class C subnet, and should not have any IP created in it. Changing this forces a new resource to be created.

* `links` - (Required) A list of spacecraft links. A `links` block as defined below. Changing this forces a new resource to be created.

* `event_hub_uri` - (Optional) ARM resource identifier of the Event Hub used for telemetry. Requires granting Orbital Resource Provider the rights to send telemetry into the hub.

* `minimum_elevation_degrees` - (Optional) Maximum elevation of the antenna during the contact in decimal degrees.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `links` block supports the following:

* `channels` - (Required) A list of contact profile link channels. A `channels` block as defined below.

* `direction` - (Required) Direction of the link. Possible values are `Uplink` and `Downlink`.

* `name` - (Required) Name of the link.

* `polarization` - (Required) Polarization of the link. Possible values are `LHCP`, `RHCP`, `linearVertical` and `linearHorizontal`.

---

A `channels` block supports the following:

* `name` - (Required) Name of the channel.

* `center_frequency_mhz` - (Required) Center frequency in MHz.

* `bandwidth_mhz` - (Required) Bandwidth in MHz.

* `end_point` - (Required) Customer End point to store/retrieve data during a contact. An `end_point` block as defined below.

* `modulation_configuration` - (Optional) Copy of the modem configuration file such as Kratos QRadio. Only valid for uplink directions. If provided, the modem connects to the customer endpoint and accepts commands from the customer instead of a VITA.49 stream.

* `demodulation_configuration` - (Optional) Copy of the modem configuration file such as Kratos QRadio or Kratos QuantumRx. Only valid for downlink directions. If provided, the modem connects to the customer endpoint and sends demodulated data instead of a VITA.49 stream.

---

An `end_point` block supports the following:

* `end_point_name` - (Required) Name of an end point.

* `port` - (Required) TCP port to listen on to receive data.

* `protocol` - (Required) Protocol of an end point. Possible values are `TCP` and `UDP`.

* `ip_address` - (Optional) IP address of an end point.

---

## Attribute Reference

In addition to the Arguments listed above - the following attributes are exported:

* `id` - The ID of the contact profile.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Contact profile.
* `read` - (Defaults to 5 minutes) Used when retrieving the Contact profile.
* `update` - (Defaults to 30 minutes) Used when updating the Contact profile.
* `delete` - (Defaults to 30 minutes) Used when deleting the Contact profile.

## Import

Contact profile can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_orbital_contact_profile.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Orbital/contactProfiles/contactProfile1
```
