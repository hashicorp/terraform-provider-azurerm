---
subcategory: "API Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_api_management_standalone_gateway"
description: |-
  Manages an Azure API Management Standalone Gateway.
---

# azurerm_api_management_standalone_gateway

Manages an Azure API Management Standalone Gateway.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-network"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "apim-delegation"
    service_delegation {
      name = "Microsoft.Web/serverFarms"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/action"
      ]
    }
  }
}

resource "azurerm_api_management_standalone_gateway" "example" {
  name                 = "example-gateway-flexible"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  virtual_network_type = "External"
  backend_subnet_id    = azurerm_subnet.example.id

  sku {
    capacity = 1
    name     = "WorkspaceGatewayPremium"
  }

  tags = {
    Hello = "World"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this API Management Standalone Gateway. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the API Management Standalone Gateway should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the API Management Standalone Gateway should exist. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `backend_subnet_id` - (Optional) Specifies the subnet id which the backend systems are hosted. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the API Management Standalone Gateway. Changing this forces a new resource to be created.

* `virtual_network_type` - (Optional) Specifies the type of VPN in which API Management gateway needs to be configured in. Possible values are `External` and `Internal`. Changing this forces a new resource to be created.

---

A `sku` block supports the following:

* `name` - (Required) The Name of the Sku. The only possible value is `WorkspaceGatewayPremium`.

* `capacity` - (Optional) The number of deployed units of the Sku. Defaults to `1`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the API Management Standalone Gateway.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the API Management Standalone Gateway.
* `read` - (Defaults to 5 minutes) Used when retrieving the API Management Standalone Gateway.
* `update` - (Defaults to 30 minutes) Used when updating the API Management Standalone Gateway.
* `delete` - (Defaults to 30 minutes) Used when deleting the API Management Standalone Gateway.

## Import

API Management Standalone Gateway can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_api_management_standalone_gateway.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ApiManagement/gateways/gateway1
```
