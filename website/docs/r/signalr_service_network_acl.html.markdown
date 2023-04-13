---
subcategory: "Messaging"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_signalr_service_network_acl"
description: |-
  Manages the Network ACL for a SignalR service.
---

# azurerm_signalr_service_network_acl

Manages the Network ACL for a SignalR service.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_signalr_service" "example" {
  name                = "example-signalr"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    name     = "Standard_S1"
    capacity = 1
  }
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.5.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.5.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
}

resource "azurerm_private_endpoint" "example" {
  name                = "example-privateendpoint"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  subnet_id           = azurerm_subnet.example.id

  private_service_connection {
    name                           = "psc-sig-test"
    is_manual_connection           = false
    private_connection_resource_id = azurerm_signalr_service.example.id
    subresource_names              = ["signalr"]
  }
}

resource "azurerm_signalr_service_network_acl" "example" {
  signalr_service_id = azurerm_signalr_service.example.id
  default_action     = "Deny"

  public_network {
    allowed_request_types = ["ClientConnection"]
  }

  private_endpoint {
    id                    = azurerm_private_endpoint.example.id
    allowed_request_types = ["ServerConnection"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `signalr_service_id` - (Required) The ID of the SignalR service. Changing this forces a new resource to be created.

* `default_action` - (Required) The default action to control the network access when no other rule matches. Possible values are `Allow` and `Deny`.

* `public_network` - (Required) A `public_network` block as defined below.

* `private_endpoint` - (Optional) A `private_endpoint` block as defined below.

---

A `public_network` block supports the following:

* `allowed_request_types` - (Optional) The allowed request types for the public network. Possible values are `ClientConnection`, `ServerConnection`, `RESTAPI` and `Trace`.

~> **Note:** When `default_action` is `Allow`, `allowed_request_types`cannot be set.

* `denied_request_types` - (Optional) The denied request types for the public network. Possible values are `ClientConnection`, `ServerConnection`, `RESTAPI` and `Trace`.

~> **Note:** When `default_action` is `Deny`, `denied_request_types`cannot be set.

~> **Note:** `allowed_request_types` - (Optional) and `denied_request_types` cannot be set together.

---

A `private_endpoint` block supports the following:

* `id` - (Required) The ID of the Private Endpoint which is based on the SignalR service.

* `allowed_request_types` - (Optional) The allowed request types for the Private Endpoint Connection. Possible values are `ClientConnection`, `ServerConnection`, `RESTAPI` and `Trace`.

~> **Note:** When `default_action` is `Allow`, `allowed_request_types`cannot be set.

* `denied_request_types` - (Optional) The denied request types for the Private Endpoint Connection. Possible values are `ClientConnection`, `ServerConnection`, `RESTAPI` and `Trace`.

~> **Note:** When `default_action` is `Deny`, `denied_request_types`cannot be set.

~> **Note:** `allowed_request_types` - (Optional) and `denied_request_types` cannot be set together.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SignalR service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network ACL of the SignalR service
* `read` - (Defaults to 5 minutes) Used when retrieving the Network ACL of the SignalR service
* `update` - (Defaults to 30 minutes) Used when updating the Network ACL of the SignalR service
* `delete` - (Defaults to 30 minutes) Used when deleting the Network ACL of the SignalR service

## Import

Network ACLs for a SignalR service can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_signalr_service_network_acl.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.SignalRService/signalR/signalr1
```
