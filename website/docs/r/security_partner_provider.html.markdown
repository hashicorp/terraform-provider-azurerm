---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_security_partner_provider"
description: |-
  Manages a Security Partner Provider.
---

# azurerm_security_partner_provider

Manages a Security Partner Provider.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "example" {
  name                = "example-vwan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_hub" "example" {
  name                = "example-vhub"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  virtual_wan_id      = azurerm_virtual_wan.example.id
  address_prefix      = "10.0.2.0/24"
}

resource "azurerm_vpn_gateway" "example" {
  name                = "example-vpngw"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  virtual_hub_id      = azurerm_virtual_hub.example.id
}

resource "azurerm_security_partner_provider" "example" {
  name                   = "example-spp"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  virtual_hub_id         = azurerm_virtual_hub.example.id
  security_provider_type = "IBoss"

  tags = {
    ENV = "Prod"
  }

  depends_on = [azurerm_vpn_gateway.example]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Security Partner Provider. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Security Partner Provider should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Security Partner Provider should exist. Changing this forces a new resource to be created.

* `security_provider_type` - (Required) The security provider name. Possible values are `ZScaler`, `IBoss` and `Checkpoint` is allowed. Changing this forces a new resource to be created.

* `virtual_hub_id` - (Optional) The ID of the Virtual Hub within which this Security Partner Provider should be created. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Security Partner Provider.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Security Partner Provider.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Security Partner Provider.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Partner Provider.
* `update` - (Defaults to 30 minutes) Used when updating the Security Partner Provider.
* `delete` - (Defaults to 30 minutes) Used when deleting the Security Partner Provider.

## Import

Security Partner Providers can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_security_partner_provider.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/securityPartnerProviders/securityPartnerProvider1
```
