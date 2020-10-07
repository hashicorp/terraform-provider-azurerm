---
subcategory: "Private DNS"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_dns_zone_virtual_network_link"
description: |-
  Manages a Private DNS Zone Virtual Network Link.
---

# azurerm_private_dns_zone_virtual_network_link

Enables you to manage Private DNS zone Virtual Network Links. These Links enable DNS resolution and registration inside Azure Virtual Networks using Azure Private DNS.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acceptanceTestResourceGroup1"
  location = "West US"
}

resource "azurerm_private_dns_zone" "example" {
  name                = "mydomain.com"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_private_dns_zone_virtual_network_link" "example" {
  name                  = "test"
  resource_group_name   = azurerm_resource_group.example.name
  private_dns_zone_name = azurerm_private_dns_zone.example.name
  virtual_network_id    = azurerm_virtual_network.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Private DNS Zone Virtual Network Link. Changing this forces a new resource to be created.

* `private_dns_zone_name` - (Required) The name of the Private DNS zone (without a terminating dot). Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the Private DNS Zone exists. Changing this forces a new resource to be created.

* `virtual_network_id` - (Required) The ID of the Virtual Network that should be linked to the DNS Zone. Changing this forces a new resource to be created.

* `registration_enabled` - (Optional) Is auto-registration of virtual machine records in the virtual network in the Private DNS zone enabled? Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Private DNS Zone Virtual Network Link.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Private DNS Zone Virtual Network Link.
* `update` - (Defaults to 30 minutes) Used when updating the Private DNS Zone Virtual Network Link.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private DNS Zone Virtual Network Link.
* `delete` - (Defaults to 30 minutes) Used when deleting the Private DNS Zone Virtual Network Link.

## Import

Private DNS Zone Virtual Network Links can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_private_dns_zone_virtual_network_link.link1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/privateDnsZones/zone1.com/virtualNetworkLinks/myVnetLink1
```
