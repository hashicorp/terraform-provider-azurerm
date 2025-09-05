subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_subnet"
description: |-
  Manages an Azure subnet. Subnets represent network segments within the IP space defined by the virtual network.

---

# azurerm_subnet

Manages an Azure subnet. Subnets represent network segments within the IP space defined by the virtual network.

~> **Note:** Terraform currently provides both a standalone [Subnet resource](subnet.html) and allows for Subnets to be defined in-line within the [Virtual Network resource](virtual_network.html).  
> **Tip:** Avoid using both in-line subnets and standalone subnet resources together to prevent configuration conflicts, which may overwrite existing subnets.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
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
    name = "delegation"

    service_delegation {
      name    = "Microsoft.ContainerInstance/containerGroups"
      actions = [
        "Microsoft.Network/virtualNetworks/subnets/join/action",
        "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action"
      ]
    }
  }
}
