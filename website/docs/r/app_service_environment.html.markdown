---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_environment"
description: |-
  Manages an App Service Environment.

---

# azurerm_app_service_environment

Manages an App Service Environment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "exampleRG1"
  location = "westeurope"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet1"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "ase" {
  name                 = "asesubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_subnet" "gateway" {
  name                 = "gatewaysubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_app_service_environment" "example" {
  name                   = "example-ase"
  subnet_id              = azurerm_subnet.ase.id
  pricing_tier           = "I2"
  front_end_scale_factor = 10
}

```

## Argument Reference

* `name` - (Required) The name of the App Service Environment. Changing this forces a new resource to be created. 

* `subnet_id` - (Required) The ID of the Subnet which the App Service Environment should be connected to. Changing this forces a new resource to be created.

~> **NOTE** a /24 or larger CIDR is required. Once associated with an ASE this size cannot be changed.

* `pricing_tier` - (Optional) Pricing tier for the front end instances. Possible values are `I1`, `I2` and `I3`. Defaults to `I1`.

* `front_end_scale_factor` - (Optional) Scale factor for front end instances. Possible values are between `5` and `15`. Defaults to `15`.

## Attribute Reference

* `id` - The ID of the App Service Environment.

* `resource_group_name` - The name of the Resource Group where the App Service Environment exists.

* `location` - The location where the App Service Environment exists.

## Import

```shell
terraform import azurerm_app_service_environment.myAppServiceEnv /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Web/hostingEnvironments/myAppServiceEnv
```