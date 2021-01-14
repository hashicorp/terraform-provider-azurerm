---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service"
description: |-
  Manages a Private Link Service.
---

# azurerm_private_link_service

Manages a Private Link Service.

-> **NOTE** Private Link is now in [GA](https://docs.microsoft.com/en-gb/azure/private-link/).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                                          = "example-subnet"
  resource_group_name                           = azurerm_resource_group.example.name
  virtual_network_name                          = azurerm_virtual_network.example.name
  address_prefixes                              = ["10.5.1.0/24"]
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
  name                = "example-privatelink"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  auto_approval_subscription_ids              = ["00000000-0000-0000-0000-000000000000"]
  visibility_subscription_ids                 = ["00000000-0000-0000-0000-000000000000"]
  load_balancer_frontend_ip_configuration_ids = [azurerm_lb.example.frontend_ip_configuration.0.id]

  nat_ip_configuration {
    name                       = "primary"
    private_ip_address         = "10.5.1.17"
    private_ip_address_version = "IPv4"
    subnet_id                  = azurerm_subnet.example.id
    primary                    = true
  }

  nat_ip_configuration {
    name                       = "secondary"
    private_ip_address         = "10.5.1.18"
    private_ip_address_version = "IPv4"
    subnet_id                  = azurerm_subnet.example.id
    primary                    = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Private Link Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Private Link Service should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `nat_ip_configuration` - (Required) One or more (up to 8) `nat_ip_configuration` block as defined below.

* `load_balancer_frontend_ip_configuration_ids` - (Required) A list of Frontend IP Configuration ID's from a Standard Load Balancer, where traffic from the Private Link Service should be routed. You can use Load Balancer Rules to direct this traffic to appropriate backend pools where your applications are running.

---

* `auto_approval_subscription_ids` - (Optional) A list of Subscription UUID/GUID's that will be automatically be able to use this Private Link Service.

* `enable_proxy_protocol` - (Optional) Should the Private Link Service support the Proxy Protocol? Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

* `visibility_subscription_ids` - (Optional) A list of Subscription UUID/GUID's that will be able to see this Private Link Service.

-> **NOTE:** If no Subscription ID's are specified then Azure allows every Subscription to see this Private Link Service.

---

The `nat_ip_configuration` block supports the following:

* `name` - (Required) Specifies the name which should be used for the NAT IP Configuration. Changing this forces a new resource to be created.

* `subnet_id` - (Required) Specifies the ID of the Subnet which should be used for the Private Link Service.

-> **NOTE:** Verify that the Subnet's `enforce_private_link_service_network_policies` attribute is set to `true`.

* `primary` - (Required) Is this is the Primary IP Configuration? Changing this forces a new resource to be created.

* `private_ip_address` - (Optional) Specifies a Private Static IP Address for this IP Configuration.

* `private_ip_address_version` - (Optional) The version of the IP Protocol which should be used. At this time the only supported value is `IPv4`. Defaults to `IPv4`.

## Attributes Reference

The following attributes are exported:

* `alias` - A globally unique DNS Name for your Private Link Service. You can use this alias to request a connection to your Private Link Service.

* `network_interfaces` - A list of network interface resource ids that are being used by the service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Private Link Service.
* `update` - (Defaults to 60 minutes) Used when updating the Private Link Service.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private Link Service.
* `delete` - (Defaults to 60 minutes) Used when deleting the Private Link Service.

## Import

Private Link Services can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_private_link_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/privateLinkServices/service1
```
