---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_vpn_server_configuration_policy_group"
description: |-
  Manages a VPN Server Configuration Policy Group.
---

# azurerm_vpn_server_configuration_policy_group

Manages a VPN Server Configuration Policy Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_vpn_server_configuration" "example" {
  name                     = "example-VPNSC"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  vpn_authentication_types = ["Radius"]

  radius {
    server {
      address = "10.105.1.1"
      secret  = "vindicators-the-return-of-worldender"
      score   = 15
    }
  }
}

resource "azurerm_vpn_server_configuration_policy_group" "example" {
  name                        = "example-VPNSCPG"
  vpn_server_configuration_id = azurerm_vpn_server_configuration.example.id

  policy {
    name  = "policy1"
    type  = "RadiusAzureGroupId"
    value = "6ad1bd08"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this VPN Server Configuration Policy Group. Changing this forces a new resource to be created.

* `vpn_server_configuration_id` - (Required) The ID of the VPN Server Configuration that the VPN Server Configuration Policy Group belongs to. Changing this forces a new resource to be created.

* `policy` - (Required) One or more `policy` blocks as documented below.

* `is_default` - (Optional) Is this a default VPN Server Configuration Policy Group? Defaults to `false`. Changing this forces a new resource to be created.

* `priority` - (Optional) The priority of this VPN Server Configuration Policy Group. Defaults to `0`.

---

A `policy` block supports the following:

* `name` - (Required) The name of the VPN Server Configuration Policy member.

* `type` - (Required) The attribute type of the VPN Server Configuration Policy member. Possible values are `AADGroupId`, `CertificateGroupId` and `RadiusAzureGroupId`.

* `value` - (Required) The value of the attribute that is used for the VPN Server Configuration Policy member.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The ID of the VPN Server Configuration Policy Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the VPN Server Configuration Policy Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the VPN Server Configuration Policy Group.
* `update` - (Defaults to 30 minutes) Used when updating the VPN Server Configuration Policy Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the VPN Server Configuration Policy Group.

## Import

VPN Server Configuration Policy Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_vpn_server_configuration_policy_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.Network/vpnServerConfigurations/serverConfiguration1/configurationPolicyGroups/configurationPolicyGroup1
```
