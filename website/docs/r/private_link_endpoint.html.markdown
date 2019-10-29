---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_link_endpoint"
sidebar_current: "docs-azurerm-resource-private-endpoint"
description: |-
  Manages a Azure Private Link Endpoint instance.
---

# azurerm_private_link_endpoint

Manages a Azure `Private Link Endpoint` instance.

Azure `Private Link Endpoint` is a network interface that connects you privately and securely to a service powered by `Azure Private Link`. `Private Link Endpoint` uses a private IP address from your VNet, effectively bringing the service into your VNet. The service could be an Azure service such as Azure Storage, SQL, etc. or your own `Private Link Service`.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "exampleRG"
  location = "East US"
}

resource "azurerm_virtual_network" "example" {
  name                = "examplevnet"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "exampleSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.0.1.0/24"

  disable_private_link_service_network_policies = true
  disable_private_link_endpoint_network_policies     = true
}

resource "azurerm_public_ip" "example" {
  name                = "examplePip"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "exampleLb"
  sku                 = "Standard"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  frontend_ip_configuration {
    name                 = azurerm_public_ip.example.name
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_private_link_service" "example" {
  name                = "examplepls"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  nat_ip_configuration {
    name      = azurerm_public_ip.example.name
    subnet_id = azurerm_subnet.example.id
  }

  load_balancer_frontend_ip_configuration_ids = [azurerm_lb.test.frontend_ip_configuration.0.id]
}

resource "azurerm_private_link_endpoint" "example" {
  name                = "examplepe"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.example.id

  tags = {
    env = "example"
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the Name of the Private Link Endpoint. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Name of the Resource Group within which the Private Link Endpoint exists. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `subnet_id` - (Required) Specifies the resource ID of the subnet from which the private IP addresses will be allocated for the private link endpoint.

* `private_service_connection` - (Optional) A `private_service_connection` block as defined below. Once defined it becomes a required argument.

* `tags` - (Optional) A mapping of tags assigned to the resource. Changing this forces a new resource to be created.

---

A `private_service_connection` contains:

* `name` - (Required) Specifies the Name of the Private Service Connection. Changing this forces a new resource to be created.

* `is_manual_connection` - (Required) Specifies if the private link endpoint requires manaul approval via the remote resource owner or not. If you are trying to connect the private link endpoint to a remote resource without having the correct RBAC permissions on the remote resource set this value to `true`. Changing this forces a new resource to be created.

* `private_connection_resource_id` - (Required) The Azure resource ID of the private link enabled remote resource to connect the private link endpoint to. Changing this forces a new resource to be created.

* `subresource_names` - (Required) The subresource name(s) that the Private Link Endpoint is allowed to connect to, see `Subresource Names` below. Changing this forces a new resource to be created.

* `request_message` - (Optional) A message passed to the owner of the remote resource when the private link endpoint attempts to establish the connection to the remote resource. The request message can be a maximum of `140` characters in length.

## Subresource Names

Resource type | Subresource name | Subresource name
-- | -- | --
Sql DB/DW | sqlServer | 
Storage Account  | blob | blob_secondary
Storage Account  | table | table_secondary
Storage Account  | queue | queue_secondary
Storage Account  | file | file_secondary
Storage Account  | web | web_secondary
Data Lake File System Gen2 | dfs | dfs_secondary

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Prviate Link Endpoint.

* `network_interface_ids` - Displays an list of network interface resource IDs that have been created for this Private Link Endpoint.

A `private_service_connection` - block exports the following:

* `request_response` - Possible values are as follows:
   * `Auto-Approved` - The remote resource owner has added you to the `Auto-Approved` RBAC permission list for the remote resouce, all private link endpoint connection requests will be automatically `Approved`.

   * `Deleted state` - The resource owner has `Rejected` the private link endpoint connection request and has removed your private link endpoint request from the remote resource.

   * `request/response message` - If you submitted a manual private link endpoint connection request, while in the `Pending` status the `request_response` will display the same text from your `request_message` in the `private_service_connection` block above. If the private link endpoint connection request was `Rejected` by the owner of the remote resource, the text for the rejection will be displayed as the `request_response` text, if If the private link endpoint connection request was `Approved` by the owner of the remote resource, the text for the approval will be displayed as the `request_response` text

   * `status` - The current status of the privat link endpoint request, possible values will be `Pending`, `Approved`, `Rejected`, or `Disconnected`.

* `private_ip_address` - The private IP address associated with the private link endpoint, note that you will have a private IP address assigned to the private link endpoint even if the connection request was `Rejected`.

## Import

The `Private Link Endpoint` can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_private_link_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/privateEndpoints/example-private-link-endpoint
```
