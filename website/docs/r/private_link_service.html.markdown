---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_service"
sidebar_current: "docs-azurerm-resource-private-link-service"
description: |-
  Manages an Azure PrivateLinkService instance.
---

# azurerm_private_link_service

Managea an Azure PrivateLinkService instance.


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
  name                 = "example-snet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.5.1.0/24"
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
  name                = "example-pls"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  fqdns               = ["testFqdns"]

  nat_ip_configuration {
    name                         = "primaryIpConfiguration"
    subnet_id                    = azurerm_subnet.example.id
    private_ip_address           = "10.5.1.17"
    private_ip_address_version   = "IPv4"
    private_ip_allocation_method = "Static"
  }

  load_balancer_frontend_ip_configuration_ids = [
    id = azurerm_lb.example.frontend_ip_configuration.0.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private link service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the private link service resides. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `auto_approval_subscription_ids` - (Optional) A list of subscription globally unique identifiers(GUID) that will be automatically be able to use this service.

* `nat_ip_configuration` - (Required) One or more `nat_ip_configuration` blocks as defined below.

* `load_balancer_frontend_ip_configuration_ids` - (Required) A list of `Standard` Load Balancer resource ids to direct the service network traffic toward.

* `private_endpoint_connection` - (Optional) One or more `private_endpoint_connection` blocks as defined below.

* `visibility_subscription_ids` - (Optional) A list of subscription globally unique identifiers(GUID) that will be able to see this service. If left undefined all Azure subscriptions will be able to see this service.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

---

The `nat_ip_configuration` block supports the following:

* `name` - (Required) The name of private link service ip configuration.

* `primary` - (Optional) If the `ip_configuration` is the primary ip configuration or not. Defaults to `true`.

* `private_ip_address` - (Optional) The private IP address of the IP configuration.

* `private_ip_allocation_method` - (Optional) The private IP address allocation method, supported values are `Static` and `Dynamic`. Defaults to `Dynamic`.

* `subnet_id` - (Required) The resource ID of the subnet to be used by the service.

* `private_ip_address_version` - (Optional) The ip address version of the `ip_configuration`, supported values are `IPv4` or `IPv6`. Defaults to `IPv4`.

-> **NOTE:** Private Link Service Supports `IPv4` traffic only.


---

The `private_endpoint_connection` block contains the following:

* `id` - (Computed) The resource ID of the `private_endpoint_connection`.

* `name` - (Required) The name of the resource that is unique within a resource group. This name can be used to access the resource.

* `private_endpoint_location` - (Required) The resource location of the `private_endpoint`.

* `private_endpoint_id` - (Required) The Private Endpoint ID.

* `state_action_required` - (Computed) A message indicating if changes on the service provider require any updates by the consumer.

* `state_description` - (Computed) The reason for approval/rejection of the connection.

* `state_status` - (Computed) Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.


## Attributes Reference

The following attributes are exported:

* `alias` - The alias of the private link service.

* `network_interfaces` - A list of network interface resource ids that are being used by the service.

---

* A `private_endpoint_connection` block contains:

* `id` - The resource ID of the `private_endpoint_connection`.

* `state_action_required` - A message indicating if changes on the service provider require any updates by the consumer.

* `state_description` - The reason for approval/rejection of the connection.

* `state_status` - Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.


## Import

Private Link Service can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_private_link_service.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG/providers/Microsoft.Network/privateLinkServices/privatelinkservicename
```
