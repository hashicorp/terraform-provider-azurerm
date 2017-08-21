---
layout: "azurerm"
page_title: "Azure Resource Manager: azure_application_gateway"
sidebar_current: "docs-azurerm-resource-application-gateway"
description: |-
  Creates a new application gateway based on a previously created virtual network with configured subnets. 
---

# azurerm\_application\_gateway

Creates a new application gateway ibased on a previously created virtual network with configured subnets.

~> **NOTE on Application Gateways:** Terraform currently
provides both a standalone [Subnet resource](subnet.html), and allows for Subnets to be defined in-line within the [Virtual Network resource](virtual_network.html).
At this time you cannot use a Virtual Network with in-line Subnets in conjunction with any Subnet resources. Doing so will cause a conflict of Subnet configurations and will overwrite Subnet's.

## Example Usage

```hcl
# Create a resource group
resource "azurerm_resource_group" "rg" {
  name     = "my-rg-application-gateway-12345"
  location = "West US"
}
 
# Create a virtual network in the web_servers resource group
resource "azurerm_virtual_network" "vnet" {
  name                = "my-vnet-12345"
  resource_group_name = "${azurerm_resource_group.rg.name}"
  address_space       = ["10.254.0.0/16"]
  location            = "${azurerm_resource_group.rg.location}"
}

resource "azurerm_subnet" "sub1" {
  name                 = "my-subnet-1"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet.name}"
  address_prefix       = "10.254.0.0/24"
}

resource "azurerm_subnet" "sub2" {
  name                 = "my-subnet-2"
  resource_group_name  = "${azurerm_resource_group.rg.name}"
  virtual_network_name = "${azurerm_virtual_network.vnet.name}"
  address_prefix       = "10.254.2.0/24"
}

resource "azurerm_public_ip" "pip" {
  name                         = "my-pip-12345"
  location                     = "${azurerm_resource_group.rg.location}"
  resource_group_name          = "${azurerm_resource_group.rg.name}"
  public_ip_address_allocation = "dynamic"
}

# Create an application gateway
resource "azurerm_application_gateway" "network" {
  name                = "my-application-gateway-12345"
  location            = "West US"
  resource_group_name = "${azurerm_resource_group.rg.name}"
 
  sku {
    name           = "Standard_Small"
    tier           = "Standard"
    capacity       = 2
  }
 
  gateway_ip_configuration {
      name         = "my-gateway-ip-configuration"
      subnet_id    = "${azurerm_virtual_network.vnet.id}/subnets/${azurerm_subnet.sub1.name}"
  }
 
  frontend_port {
      name         = "my-frontend-port"
      port         = 80
  }
 
  frontend_ip_configuration {
      name         = "my-frontend-ip-configuration"  
      public_ip_address_id = "${azurerm_public_ip.pip.id}"
  }

  backend_address_pool {
      name = "my-backend-address-pool"
  }
 
  backend_http_settings {
      name                  = "${azurerm_virtual_network.vnet.name}-be-htst"
      cookie_based_affinity = "Disabled"
      port                  = 80
      protocol              = "Http"
     request_timeout        = 1
  }
 
  http_listener {
        name                                  = "${azurerm_virtual_network.vnet.name}-httplstn"
        frontend_ip_configuration_name        = "${azurerm_virtual_network.vnet.name}-feip"
        frontend_port_name                    = "${azurerm_virtual_network.vnet.name}-feport"
        protocol                              = "Http"
  }
 
request_routing_rule {
        name                       = "${azurerm_virtual_network.vnet.name}-rqrt"
        rule_type                  = "Basic"
        http_listener_name         = "${azurerm_virtual_network.vnet.name}-httplstn"
        backend_address_pool_name  = "${azurerm_virtual_network.vnet.name}-beap"
        backend_http_settings_name = "${azurerm_virtual_network.vnet.name}-be-htst"
}
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the virtual network. Changing this forces a
    new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to
    create the virtual network.

* `address_space` - (Required) The address space that is used the virtual
    network. You can supply more than one address space. Changing this forces
    a new resource to be created.

* `location` - (Required) The location/region where the virtual network is
    created. Changing this forces a new resource to be created.

* `dns_servers` - (Optional) List of IP addresses of DNS servers

* `subnet` - (Optional) Can be specified multiple times to define multiple
    subnets. Each `subnet` block supports fields documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

The `subnet` block supports:

* `name` - (Required) The name of the subnet.

* `address_prefix` - (Required) The address prefix to use for the subnet.

* `security_group` - (Optional) The Network Security Group to associate with
    the subnet. (Referenced by `id`, ie. `azurerm_network_security_group.test.id`)

## Attributes Reference

The following attributes are exported:

* `id` - The virtual NetworkConfiguration ID.

* `name` - The name of the virtual network.

* `resource_group_name` - The name of the resource group in which to create the virtual network.

* `location` - The location/region where the virtual network is created

* `address_space` - The address space that is used the virtual network.


## Import

Virtual Networks can be imported using the `resource id`, e.g.

```
terraform import azurerm_virtual_network.testNetwork /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1
```
