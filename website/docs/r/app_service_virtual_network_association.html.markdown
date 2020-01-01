---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_virtual_network_swift_connection"
sidebar_current: "docs-azurerm-resource-app-service-virtual-network-association"
description: |-
  Manages an App Service Virtual Network Association.

---

# azurerm_app_service_virtual_network_swift_connection

Manages an App Service Virtual Network Association (this is for the [Regional VNet Integration](https://docs.microsoft.com/en-us/azure/app-service/web-sites-integrate-with-vnet#regional-vnet-integration) which is still in preview).

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "uksouth"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test1" {
  name                 = "acctestsubnet1"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"

  delegation {
    name = "acctestdelegation"

    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_app_service_plan" "test" {
  name                = "acctestasp"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctestas"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_virtual_network_swift_connection" "test" {
  app_service_id       = "${azurerm_app_service.test.id}"
  subnet_id            = "${azurerm_subnet.test1.id}"
}
```

## Argument Reference

The following arguments are supported:

* `app_service_id` - (Required) The ID of the App Service to associate to the VNet. Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the subnet the app service will be associated to (the subnet must have a `service_delegation` configured for `Microsoft.Web/serverFarms`).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service Virtual Network Association

## Import

App Service Virtual Network Associations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_virtual_network_swift_connection.myassociation /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1/networkconfig/virtualNetwork
```
