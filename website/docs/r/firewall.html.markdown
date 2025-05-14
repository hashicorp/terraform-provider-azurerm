---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_firewall"
description: |-
  Manages an Azure Firewall.

---

# azurerm_firewall

Manages an Azure Firewall.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "testvnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "example" {
  name                = "testpip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "example" {
  name                = "testfirewall"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku_name            = "AZFW_VNet"
  sku_tier            = "Standard"

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.example.id
    public_ip_address_id = azurerm_public_ip.example.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Firewall. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the resource. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku_name` - (Required) SKU name of the Firewall. Possible values are `AZFW_Hub` and `AZFW_VNet`. Changing this forces a new resource to be created.

* `sku_tier` - (Required) SKU tier of the Firewall. Possible values are `Premium`, `Standard` and `Basic`.

* `firewall_policy_id` - (Optional) The ID of the Firewall Policy applied to this Firewall.

* `ip_configuration` - (Optional) An `ip_configuration` block as documented below.

* `dns_servers` - (Optional) A list of DNS servers that the Azure Firewall will direct DNS traffic to the for name resolution.

* `dns_proxy_enabled` - (Optional) Whether DNS proxy is enabled. It will forward DNS requests to the DNS servers when set to `true`. It will be set to `true` if `dns_servers` provided with a not empty list.

* `private_ip_ranges` - (Optional) A list of SNAT private CIDR IP ranges, or the special string `IANAPrivateRanges`, which indicates Azure Firewall does not SNAT when the destination IP address is a private range per IANA RFC 1918.

* `management_ip_configuration` - (Optional) A `management_ip_configuration` block as documented below, which allows force-tunnelling of traffic to be performed by the firewall. Adding or removing this block or changing the `subnet_id` in an existing block forces a new resource to be created. Changing this forces a new resource to be created.

* `threat_intel_mode` - (Optional) The operation mode for threat intelligence-based filtering. Possible values are: `Off`, `Alert` and `Deny`. Defaults to `Alert`.

* `virtual_hub` - (Optional) A `virtual_hub` block as documented below.

* `zones` - (Optional) Specifies a list of Availability Zones in which this Azure Firewall should be located. Changing this forces a new Azure Firewall to be created.

-> **Note:** Availability Zones are [only supported in several regions at this time](https://docs.microsoft.com/azure/availability-zones/az-overview).

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `ip_configuration` block supports the following:

* `name` - (Required) Specifies the name of the IP Configuration.

* `subnet_id` - (Optional) Reference to the subnet associated with the IP Configuration. Changing this forces a new resource to be created.

-> **Note:** The Subnet used for the Firewall must have the name `AzureFirewallSubnet` and the subnet mask must be at least a `/26`.

-> **Note:** At least one and only one `ip_configuration` block may contain a `subnet_id`.

* `public_ip_address_id` - (Optional) The ID of the Public IP Address associated with the firewall.

-> **Note:** A public ip address is required unless a `management_ip_configuration` block is specified.

-> **Note:** When multiple `ip_configuration` blocks with `public_ip_address_id` are configured, `terraform apply` will raise an error when one or some of these `ip_configuration` blocks are removed. because the `public_ip_address_id` is still used by the `firewall` resource until the `firewall` resource is updated. and the destruction of `azurerm_public_ip` happens before the update of firewall by default. to destroy of `azurerm_public_ip` will cause the error. The workaround is to set `create_before_destroy=true` to the `azurerm_public_ip` resource `lifecycle` block. See more detail: [destroying.md#create-before-destroy](https://github.com/hashicorp/terraform/blob/main/docs/destroying.md#create-before-destroy)

-> **Note:** The Public IP must have a `Static` allocation and `Standard` SKU.

---

A `management_ip_configuration` block supports the following:

* `name` - (Required) Specifies the name of the IP Configuration.

* `subnet_id` - (Required) Reference to the subnet associated with the IP Configuration. Changing this forces a new resource to be created.

-> **Note:** The Management Subnet used for the Firewall must have the name `AzureFirewallManagementSubnet` and the subnet mask must be at least a `/26`.

* `public_ip_address_id` - (Required) The ID of the Public IP Address associated with the firewall.

-> **Note:** The Public IP must have a `Static` allocation and `Standard` SKU.

---

A `virtual_hub` block supports the following:

* `virtual_hub_id` - (Required) Specifies the ID of the Virtual Hub where the Firewall resides in.

* `public_ip_count` - (Optional) Specifies the number of public IPs to assign to the Firewall. Defaults to `1`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Firewall.

* `ip_configuration` - A `ip_configuration` block as defined below.

* `virtual_hub` - A `virtual_hub` block as defined below.

---

A `ip_configuration` block exports the following:

* `private_ip_address` - The Private IP address of the Azure Firewall.

---

A `virtual_hub` block exports the following:

* `private_ip_address` - The private IP address associated with the Firewall.

* `public_ip_addresses` - The list of public IP addresses associated with the Firewall.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 minutes) Used when creating the Firewall.
* `read` - (Defaults to 5 minutes) Used when retrieving the Firewall.
* `update` - (Defaults to 90 minutes) Used when updating the Firewall.
* `delete` - (Defaults to 90 minutes) Used when deleting the Firewall.

## Import

Azure Firewalls can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_firewall.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/azureFirewalls/testfirewall
```
