---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_slot_virtual_network_swift_connection"
description: |-
  Manages an App Service's Slot Virtual Network Association.

---

# azurerm_app_service_slot_virtual_network_swift_connection

Manages an App Service Slot's Virtual Network Association (this is for the [Regional VNet Integration](https://docs.microsoft.com/en-us/azure/app-service/web-sites-integrate-with-vnet#regional-vnet-integration) which is still in preview).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "uksouth"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-virtual-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name               = "example-subnet"
  virtual_network_id = azurerm_virtual_network.example.id
  address_prefixes   = ["10.0.1.0/24"]

  delegation {
    name = "example-delegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_app_service_plan" "example" {
  name                = "example-service-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = "example-app-service"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id
}

resource "azurerm_app_service_slot" "example-staging" {
  name                = "staging"
  app_service_name    = azurerm_app_service.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id
}

resource "azurerm_app_service_slot_virtual_network_swift_connection" "example" {
  slot_name      = azurerm_app_service_slot.example-staging.name
  app_service_id = azurerm_app_service.example.id
  subnet_id      = azurerm_subnet.example.id
}
```

## Argument Reference

The following arguments are supported:

* `app_service_id` - (Required) The ID of the App Service or Function App to associate to the VNet. Changing this forces a new resource to be created.

* `slot_name` - (Required) The name of the App Service Slot or Function App Slot. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet the app service will be associated to (the subnet must have a `service_delegation` configured for `Microsoft.Web/serverFarms`).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service Slot Virtual Network Association

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Virtual Network Association.
* `update` - (Defaults to 30 minutes) Used when updating the App Service Virtual Network Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Virtual Network Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Virtual Network Association.

## Import

App Service Slot Virtual Network Associations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_slot_virtual_network_swift_connection.myassociation /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1/slots/stageing/networkconfig/virtualNetwork
```
