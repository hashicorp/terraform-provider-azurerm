---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_private_endpoint"
description: |-
  Manages a Private Endpoint.
---

# azurerm_private_endpoint

Manages a Private Endpoint.

-> **NOTE** Private Endpoint is currently in Public Preview.

Azure Private Endpoint is a network interface that connects you privately and securely to a service powered by Azure Private Link. Private Endpoint uses a private IP address from your VNet, effectively bringing the service into your VNet. The service could be an Azure service such as Azure Storage, SQL, etc. or your own Private Link Service.

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
  address_prefixes     = ["10.0.1.0/24"]

  enforce_private_link_service_network_policies = true
}

resource "azurerm_subnet" "endpoint" {
  name                 = "endpoint"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  enforce_private_link_endpoint_network_policies = true
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

resource "azurerm_private_endpoint" "example" {
  name                = "example-endpoint"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.endpoint.id

  private_service_connection {
    name                           = "example-privateserviceconnection"
    private_connection_resource_id = azurerm_private_link_service.example.id
    is_manual_connection           = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the Name of the Private Endpoint. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the Name of the Resource Group within which the Private Endpoint should exist. Changing this forces a new resource to be created.

* `location` - (Required) The supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet from which Private IP Addresses will be allocated for this Private Endpoint. Changing this forces a new resource to be created.

* `private_dns_zone_group` - (Optional) A `private_dns_zone_group` block as defined below.

* `private_service_connection` - (Required) A `private_service_connection` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `private_dns_zone_group` supports the following:

* `name` - (Required) Specifies the Name of the Private DNS Zone Group. Changing this forces a new `private_dns_zone_group` resource to be created.

* `private_dns_zone_ids` - (Required) Specifies the list of Private DNS Zones to include within the `private_dns_zone_group`.

---

A `private_service_connection` supports the following:

* `name` - (Required) Specifies the Name of the Private Service Connection. Changing this forces a new resource to be created.

* `is_manual_connection` - (Required) Does the Private Endpoint require Manual Approval from the remote resource owner? Changing this forces a new resource to be created.

-> **NOTE:** If you are trying to connect the Private Endpoint to a remote resource without having the correct RBAC permissions on the remote resource set this value to `true`.

* `private_connection_resource_id` - (Required) The ID of the Private Link Enabled Remote Resource which this Private Endpoint should be connected to. Changing this forces a new resource to be created.

* `subresource_names` - (Required) A list of subresource names which the Private Endpoint is able to connect to. `subresource_names` corresponds to `group_id`. Changing this forces a new resource to be created.

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

* `request_message` - (Optional) A message passed to the owner of the remote resource when the private endpoint attempts to establish the connection to the remote resource. The request message can be a maximum of `140` characters in length. Only valid if `is_manual_connection` is set to `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Private Endpoint.

---

A `private_dns_zone_group` block exports:

* `id` - The ID of the Private DNS Zone Group.

---

A `custom_dns_configs` block exports:

* `fqdn` - The fully qualified domain name to the `private_endpoint`.

* `ip_addresses` - A list of all IP Addresses that map to the `private_endpoint` fqdn.

-> **NOTE:** If a Private DNS Zone Group has been defined and is currently connected correctly this block will be empty.

---

A `private_dns_zone_configs` block exports:

* `name` - The name of the Private DNS Zone that the config belongs to.

* `id` - The ID of the Private DNS Zone Config.

* `private_dns_zone_id` - A list of IP Addresses

* `record_sets` - A `record_sets` block as defined below.

---

A `private_service_connection` block exports:

* `private_ip_address` - (Computed) The private IP address associated with the private endpoint, note that you will have a private IP address assigned to the private endpoint even if the connection request was `Rejected`.

---

A `record_sets` block exports:

* `name` - The name of the Private DNS Zone that the config belongs to.

* `type` - The type of DNS record.

* `fqdn` - The fully qualified domain name to the `private_dns_zone`.

* `ttl` - The time to live for each connection to the `private_dns_zone`.

* `ip_addresses` - A list of all IP Addresses that map to the `private_dns_zone` fqdn.

-> **NOTE:** If a Private DNS Zone Group has not been configured correctly the `record_sets` attributes will be empty.

---

## Example HCL Configurations

* How to connect a `Private Endpoint` to a [Cosmos MongoDB](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/private-endpoint/cosmos-db)
* How to connect a `Private Endpoint` to a [PostgreSQL Server](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/private-endpoint/postgresql)
* How to connect a `Private Endpoint` to a [Private Link Service](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/private-endpoint/private-link-service)
* How to connect a `Private Endpoint` to a [Private DNS Group](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/private-endpoint/private-dns-group)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Private Endpoint.
* `update` - (Defaults to 60 minutes) Used when updating the Private Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Private Endpoint.
* `delete` - (Defaults to 60 minutes) Used when deleting the Private Endpoint.

## Import

Private Endpoints can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_private_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/privateEndpoints/endpoint1
```
