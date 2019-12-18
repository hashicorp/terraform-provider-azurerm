---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_endpoint"
sidebar_current: "docs-azurerm-resource-private-endpoint"
description: |-
  Manages an Endpoint within a Private Link Service.
---
-> **NOTE:** The 'azurerm_private_link_endpoint' resource is being deprecated in favour of the renamed version 'azurerm_private_endpoint'.
Information on migrating to the renamed resource can be found here: https://terraform.io/docs/providers/azurerm/guides/migrating-between-renamed-resources.html
As such the existing 'azurerm_private_link_endpoint' resource is deprecated and will be removed in the next major version of the AzureRM Provider (2.0).

# azurerm_private_link_endpoint

Manages an Endpoint within a Private Link Service.

-> **NOTE** Private Link is currently in Public Preview.

Azure Private Link Endpoint is a network interface that connects you privately and securely to a service powered by Azure Private Link. Private Link Endpoint uses a private IP address from your VNet, effectively bringing the service into your VNet. The service could be an Azure service such as Azure Storage, SQL, etc. or your own Private Link Service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "service" {
  name                 = "service"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.0.1.0/24"

  disable_private_link_service_network_policy_enforcement  = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "endpoint"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.0.2.0/24"

  disable_private_link_endpoint_network_policy_enforcement = true
}

resource "azurerm_public_ip" "example" {
  name                = "example-pip"
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
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  nat_ip_configuration {
    name      = azurerm_public_ip.example.name
    primary   = true
    subnet_id = azurerm_subnet.service.id
  }

  load_balancer_frontend_ip_configuration_ids = [
    azurerm_lb.example.frontend_ip_configuration.0.id,
  ]
}

resource "azurerm_private_link_endpoint" "example" {
  name                = "example-endpoint"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.endpoint.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the Name of the Private Link Endpoint. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Name of the Resource Group within which the Private Link Endpoint should exist. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet from which Private IP Addresses will be allocated for this Private Link Endpoint. Changing this forces a new resource to be created.

* `private_service_connection` - (Required) A `private_service_connection` block as defined below.

---

A `private_service_connection` supports the following:

* `name` - (Required) Specifies the Name of the Private Service Connection. Changing this forces a new resource to be created.

* `is_manual_connection` - (Required) Does the Private Link Endpoint require Manual Approval from the remote resource owner? Changing this forces a new resource to be created.

-> **NOTE:** If you are trying to connect the Private Link Endpoint to a remote resource without having the correct RBAC permissions on the remote resource set this value to `true`.

* `private_connection_resource_id` - (Required) The ID of the Private Link Enabled Remote Resource which this Private Link Endpoint should be connected to. Changing this forces a new resource to be created.

* `subresource_names` - (Optional) A list of subresource names which the Private Link Endpoint is able to connect to. Changing this forces a new resource to be created.

-> Several possible values for this field are shown below, however this is not extensive:

| Resource Type                 | SubResource Name | Secondary SubResource Name |
| ----------------------------- | ---------------- | -------------------------- |
| Data Lake File System Gen2    | dfs              | dfs_secondary              |
| Sql Database / Data Warehouse | sqlServer        |                            |
| Storage Account               | blob             | blob_secondary             |
| Storage Account               | file             | file_secondary             |
| Storage Account               | queue            | queue_secondary            |
| Storage Account               | table            | table_secondary            |
| Storage Account               | web              | web_secondary              |

See the product [documentation](https://docs.microsoft.com/en-us/azure/private-link/private-endpoint-overview#dns-configuration) for more information.

* `request_message` - (Optional) A message passed to the owner of the remote resource when the private link endpoint attempts to establish the connection to the remote resource. The request message can be a maximum of `140` characters in length. Only valid if `is_manual_connection` is set to `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Private Link Endpoint.

## Import

Private Link Endpoints can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_private_link_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/privateEndpoints/endpoint1
```
