---
subcategory: "Orbital"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_orbital_spacecraft"
description: |-
  Manages a Spacecraft resource.
---

# azurerm_orbital_spacecraft

Manages a Spacecraft.

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
    bandwidth_mhz        = 30
    center_frequency_mhz = 2050
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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Spacecraft. Changing this forces a new resource to be created. 

* `resource_group_name` - (Required) The name of the Resource Group where the Spacecraft exists. Changing this forces a new resource to be created.

* `location` - (Required) The location where the Spacecraft exists. Changing this forces a new resource to be created.

* `norad_id` - (Required) NORAD ID of the Spacecraft.

* `links` - (Required) A `links` block as defined below. Changing this forces a new resource to be created.

* `two_line_elements` - (Required) A list of the two line elements (TLE), the first string being the first of the TLE, the second string being the second line of the TLE. Changing this forces a new resource to be created.

* `title_line` - (Required) Title of the two line elements (TLE).

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `links` block supports the following:

* `bandwidth_mhz` - (Required) Bandwidth in Mhz.

* `center_frequency_mhz` - (Required) Center frequency in Mhz.

~> **Note:** The value of `center_frequency_mhz +/- bandwidth_mhz / 2` should fall in one of these ranges: `Uplink/LHCP`: [2025, 2120]; `Uplink/Linear`: [399, 403],[435, 438],[449, 451]; `Uplink/RHCP`: [399, 403],[435, 438],[449, 451],[2025, 2120]; `Downlink/LHCP`: [2200, 2300], [7500, 8400]; `Downlink/Linear`: [399, 403], [435, 438], [449, 451]; Downlink/Linear`: [399, 403], [435, 438], [449, 451], [2200, 2300], [7500, 8400]

* `direction` - (Required) Direction if the communication. Possible values are `Uplink` and `Downlink`.

* `polarization` - (Required) Polarization. Possible values are `RHCP`, `LHCP`, `linearVertical` and `linearHorizontal`.

* `name` - (Required) Name of the link.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Spacecraft.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Spacecraft.
* `read` - (Defaults to 5 minutes) Used when retrieving the Spacecraft.
* `update` - (Defaults to 30 minutes) Used when updating the Spacecraft.
* `delete` - (Defaults to 30 minutes) Used when deleting the Spacecraft.

## Import

Spacecraft can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_orbital_spacecraft.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Orbital/spacecrafts/spacecraft1
```
