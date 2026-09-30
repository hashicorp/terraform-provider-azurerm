---
subcategory: "Discovery"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_discovery_bookshelf_private_endpoint_connection"
description: |-
  Manages the connection between a Private Endpoint and a Microsoft Discovery Bookshelf.
---

# azurerm_discovery_bookshelf_private_endpoint_connection

Manages the connection between a Private Endpoint and a Microsoft Discovery Bookshelf.

-> **Note:** The Private Endpoint Connection is created by Azure when a Private Endpoint targets the Bookshelf - this resource manages the state of that connection, it does not create it. A `azurerm_private_endpoint` targeting the Bookshelf must therefore exist first.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example-identity"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "private_endpoint" {
  name                 = "privateEndpointSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_subnet" "search" {
  name                 = "searchSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/23"]

  delegation {
    name = "Microsoft.App/environments"

    service_delegation {
      name    = "Microsoft.App/environments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_subnet" "endpoint" {
  name                 = "endpointSubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.4.0/24"]
}

resource "azurerm_discovery_bookshelf" "example" {
  name                       = "example-bookshelf"
  resource_group_name        = azurerm_resource_group.example.name
  location                   = azurerm_resource_group.example.location
  private_endpoint_subnet_id = azurerm_subnet.private_endpoint.id
  search_subnet_id           = azurerm_subnet.search.id
  workload_identity_ids      = [azurerm_user_assigned_identity.example.id]
}

data "azurerm_discovery_bookshelf_private_link_resources" "example" {
  bookshelf_id = azurerm_discovery_bookshelf.example.id
}

resource "azurerm_private_endpoint" "example" {
  name                = "example-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = "example-connection"
    private_connection_resource_id = azurerm_discovery_bookshelf.example.id
    subresource_names              = [data.azurerm_discovery_bookshelf_private_link_resources.example.private_link_resources[0].group_id]
    is_manual_connection           = true
    request_message                = "Please approve this connection"
  }
}

resource "azurerm_discovery_bookshelf_private_endpoint_connection" "example" {
  name         = azurerm_private_endpoint.example.name
  bookshelf_id = azurerm_discovery_bookshelf.example.id
  status       = "Approved"
  description  = "Approved by Terraform"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Private Endpoint Connection. Changing this forces a new Microsoft Discovery Bookshelf Private Endpoint Connection to be created.

* `bookshelf_id` - (Required) The ID of the Microsoft Discovery Bookshelf which the Private Endpoint is connected to. Changing this forces a new Microsoft Discovery Bookshelf Private Endpoint Connection to be created.

* `status` - (Required) The status of the Private Endpoint Connection. Possible values are `Approved` and `Rejected`.

---

* `description` - (Optional) A message describing why the Private Endpoint Connection is in its current state.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Microsoft Discovery Bookshelf Private Endpoint Connection.

* `actions_required` - A message indicating whether changes on the Private Endpoint require any updates on the consumer.

* `group_ids` - A list of Private Link Resource group IDs which the Private Endpoint is connected to.

* `private_endpoint_id` - The ID of the Private Endpoint which this connection belongs to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Microsoft Discovery Bookshelf Private Endpoint Connection.
* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft Discovery Bookshelf Private Endpoint Connection.
* `update` - (Defaults to 30 minutes) Used when updating the Microsoft Discovery Bookshelf Private Endpoint Connection.
* `delete` - (Defaults to 30 minutes) Used when deleting the Microsoft Discovery Bookshelf Private Endpoint Connection.

## Import

Microsoft Discovery Bookshelf Private Endpoint Connections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_discovery_bookshelf_private_endpoint_connection.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/example-resources/providers/Microsoft.Discovery/bookshelves/example-bookshelf/privateEndpointConnections/example-connection
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Discovery` - 2026-06-01
