---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_environment"
description: |-
  Manages an App Service Environment.

---

# azurerm_app_service_environment

Manages a App Service Environment

**WARNING** Deleting an App Service Environment resource will also delete App Service Plans and App Services associated with it. 

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

  subnet {
    name           = "asesubnet"
    address_prefix = "10.0.1.0/24"
  }

  subnet {
    name           = "gatewaysubnet"
    address_prefix = "10.0.2.0/24"
  }
}

data "azurerm_subnet" "example" {
  name                 = "asesubnet"
  virtual_network_name = "${azurerm_virtual_network.example.name}"
  resource_group_name  = "${azurerm_resource_group.example.name}"
}

resource "azurerm_app_service_environment" "example" {
  name                   = "example-ase"
  subnet_id              = "${data.azurerm_subnet.example.id}"
  pricing_tier           = "I2"
  front_end_scale_factor = 10
}

```

## Argument Reference

* `name` - (Required) name of the App Service Environment. 

~> **NOTE** Must meet DNS name specification.

* `subnet_id` - (Required) Resource ID for the ASE subnet.

~> **NOTE** a /24 or larger CIDR is required. Once associated with an ASE this size cannot be changed.

* `pricing_tier` - (Optional) Pricing tier for the front end instances. Possible values are `I1` (default), `I2` and `I3`. 

~> **NOTE** Azure currently utilises Dv2 instances for Isolated SKUs, being `Standard_D1_V2`, `Standard_D2_V2`, and `Standard_D3_V2`.

* `front_end_scale_factor` - (Optional) Scale factor for front end instances. Possible values are between `15` (default) and `5`.

~> **NOTE** Lowering/changing this value has cost implications, see https://docs.microsoft.com/en-us/azure/app-service/environment/using-an-ase#front-end-scaling for details.

## Attribute Reference

* `id` - The ID of the App Services Environment.

* `resource_group_name` - The name of the resource group.

* `location` - The location the App Service Environment is deployed into.

## Import

```shell
terraform import azurerm_app_service_environment.myAppServiceEnv /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Web/hostingEnvironments/myAppServiceEnv
```