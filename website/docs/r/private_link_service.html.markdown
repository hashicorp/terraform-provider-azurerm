---
subcategory: ""
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service"
sidebar_current: "docs-azurerm-resource-private-link-service"
description: |-
  Manages an Azure Private Link Service.
---

# azurerm_private_link_service

Manages an Azure Private Link Service.


## Private Link Service Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "exampleRG"
  location = "Eastus2"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-avn"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                                          = "example-snet"
  resource_group_name                           = azurerm_resource_group.example.name
  virtual_network_name                          = azurerm_virtual_network.example.name
  address_prefix                                = "10.5.1.0/24"
  enforce_private_link_service_network_policies = true
}

resource "azurerm_public_ip" "example" {
  name                = "example-api"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "example-lb"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  frontend_ip_configuration {
    name                 = azurerm_public_ip.example.name
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_private_link_service" "example" {
  name                = "myPrivateLinkService"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  auto_approval_subscription_ids              = ["00000000-0000-0000-0000-000000000000"]
  visibility_subscription_ids                 = ["00000000-0000-0000-0000-000000000000"]
  load_balancer_frontend_ip_configuration_ids = [azurerm_lb.example.frontend_ip_configuration.0.id]

  nat_ip_configuration {
    name                       = "primaryIpConfiguration"
    private_ip_address         = "10.5.1.17"
    private_ip_address_version = "IPv4"
    subnet_id                  = azurerm_subnet.example.id
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondaryIpConfiguration"
    private_ip_address         = "10.5.1.18"
    private_ip_address_version = "IPv4"
    subnet_id                  = azurerm_subnet.example.id
    primary                    = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private link service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the private link service resides. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `auto_approval_subscription_ids` - (Optional) A list of subscription globally unique identifiers(GUID) that will be automatically be able to use this service.

* `visibility_subscription_ids` - (Optional) A list of subscription globally unique identifiers(GUID) that will be able to see this service. If left undefined all Azure subscriptions will be able to see this service.

* `nat_ip_configuration` - (Required) A `nat_ip_configuration` block as defined below. There maybe upto 8 nat_ip_configuration blocks per private link service.

* `load_balancer_frontend_ip_configuration_ids` - (Required) A list of Standard Load Balancer(SLB) resource IDs. The Private Link service is tied to the frontend IP address of a SLB. All traffic destined for the private link service will reach the frontend of the SLB. You can configure SLB rules to direct this traffic to appropriate backend pools where your applications are running.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

---

The `nat_ip_configuration` block supports the following:

* `name` - (Required) The name of primary private link service NAT IP configuration. Changing this forces a new resource to be created.

* `private_ip_address` - (Optional) The private IP address of the NAT IP configuration.

* `private_ip_address_version` - (Optional) The ip address version of the `ip_configuration`, the supported value is `IPv4`. Defaults to `IPv4`.

-> **NOTE:** Private Link Service Supports `IPv4` traffic only.

* `subnet_id` - (Required) The resource ID of the subnet to be used by the service.

-> **NOTE:** Verify that the subnets `enforce_private_link_service_network_policies` attribute is set to `true`.

* `primary` - (Required) Specifies if the `nat_ip_configuration` block is the primary ip configuration for the service or not. Valid values are `true` or `false`. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `alias` - The alias is a globally unique name for your private link service which Azure generates for you. Your can use this alias to request a connection to your private link service.

* `network_interfaces` - A list of network interface resource ids that are being used by the service.


## Import

Private Link Service can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_private_link_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG/providers/Microsoft.Network/privateLinkServices/privatelinkservicename
```
