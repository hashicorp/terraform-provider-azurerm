---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_virtual_network_swift_connection"
description: |-
  Manages an App Service Virtual Network Association.

---

# azurerm_app_service_virtual_network_swift_connection

Manages an App Service Virtual Network Association for [Regional VNet Integration](https://docs.microsoft.com/azure/app-service/web-sites-integrate-with-vnet#regional-vnet-integration).

This resource can be used for both App Services and Function Apps.

~> **Note:** The following resources support associating the vNet for Regional vNet Integration directly on the resource and via the `azurerm_app_service_virtual_network_swift_connection` resource. You can't use both simultaneously.

- [azurerm_linux_function_app](linux_function_app.html)
- [azurerm_linux_function_app_slot](linux_function_app_slot.html)
- [azurerm_linux_web_app](linux_web_app.html)
- [azurerm_linux_web_app_slot](linux_web_app_slot.html)
- [azurerm_logic_app_standard](logic_app_standard.html)
- [azurerm_windows_function_app](windows_function_app.html)
- [azurerm_windows_function_app_slot](windows_function_app_slot.html)
- [azurerm_windows_web_app](windows_web_app.html)
- [azurerm_windows_web_app_slot](windows_web_app_slot.html)

This resource requires the `Microsoft.Network/virtualNetworks/subnets/write` permission scope on the subnet.  

The resource specific vNet integration requires the `Microsoft.Network/virtualNetworks/subnets/join/action` permission scope.

There is a hard limit of [one VNet integration per App Service Plan](https://docs.microsoft.com/azure/app-service/web-sites-integrate-with-vnet#regional-vnet-integration).
Multiple apps in the same App Service plan can use the same VNet.

## Example Usage (with App Service)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-virtual-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "example-delegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_app_service_plan" "example" {
  name                = "example-app-service-plan"
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

resource "azurerm_app_service_virtual_network_swift_connection" "example" {
  app_service_id = azurerm_app_service.example.id
  subnet_id      = azurerm_subnet.example.id
}
```

## Example Usage (with Function App)

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-virtual-network"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "example-delegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_app_service_plan" "example" {
  name                = "example-app-service-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "functionsappexamplesa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_function_app" "example" {
  name                       = "example-function-app"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  app_service_plan_id        = azurerm_app_service_plan.example.id
  storage_account_name       = azurerm_storage_account.example.name
  storage_account_access_key = azurerm_storage_account.example.primary_access_key
}

resource "azurerm_app_service_virtual_network_swift_connection" "example" {
  app_service_id = azurerm_function_app.example.id
  subnet_id      = azurerm_subnet.example.id
}
```

## Argument Reference

The following arguments are supported:

* `app_service_id` - (Required) The ID of the App Service or Function App to associate to the VNet. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet the app service will be associated to (the subnet must have a `service_delegation` configured for `Microsoft.Web/serverFarms`).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Service Virtual Network Association

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Virtual Network Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Virtual Network Association.
* `update` - (Defaults to 30 minutes) Used when updating the App Service Virtual Network Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Virtual Network Association.

## Import

App Service Virtual Network Associations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_virtual_network_swift_connection.myassociation /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1/config/virtualNetwork
```
