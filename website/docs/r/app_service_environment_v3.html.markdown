---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_environment_v3"
description: |-
  Manages a 3rd Generation (v3) App Service Environment.

---

# azurerm_app_service_environment

Manages a 3rd Generation (v3) App Service Environment.

~> **NOTE:** App Service Environment V3 is currently in Preview.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "exampleRG1"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "inbound" {
  name                 = "inbound"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_subnet" "outbound" {
  name                 = "outbound"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  service_delegation {
    name    = "Microsoft.Web/hostingEnvironments"
    actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
  }
}

resource "azurerm_app_service_environment_v3" "example" {
  name                = "example-asev3"
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.outbound.id

  cluster_setting {
    name  = "DisableTls1.0"
    value = "1"
  }

  cluster_setting {
    name  = "InternalEncryption"
    value = "true"
  }

  cluster_setting {
    name  = "FrontEndSSLCipherSuiteOrder"
    value = "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384_P256,TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256_P256,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA384_P256,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256_P256,TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA_P256,TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA_P256"
  }

  tags = {
    env         = "production"
    terraformed = "true"
  }

}

```

## Argument Reference

* `name` - (Required) The name of the App Service Environment. Changing this forces a new resource to be created. 

* `resource_group_name` - (Required) The name of the Resource Group where the App Service Environment exists. Defaults to the Resource Group of the Subnet (specified by `subnet_id`).

* `subnet_id` - (Required) The ID of the Subnet which the App Service Environment should be connected to. Changing this forces a new resource to be created.

~> **NOTE** a /24 or larger CIDR is required. Once associated with an ASE, this size cannot be changed.

~> **NOTE:** This is the "outbound" Subnet which is required to have a delegation to `Microsoft.Web/hostingEnvironments` as detailed in the example above. Additionally, an "inbound" subnet is required in the Virtual Network which must not have `enforce_private_link_endpoint_network_policies` enabled.   

* `cluster_setting` - (Optional) Zero or more `cluster_setting` blocks as defined below. 

* `tags` - (Optional) A mapping of tags to assign to the resource. 

---

A `cluster_setting` block supports the following:

* `name` - (Required) The name of the Cluster Setting. 

* `value` - (Required) The value for the Cluster Setting. 

## Attribute Reference

* `id` - The ID of the App Service Environment.

* `location` - The location where the App Service Environment exists.

* `pricing_tier` - Pricing tier for the front end instances.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 6 hours) Used when creating the 3rd Generation (v3) App Service Environment.
* `update` - (Defaults to 6 hours) Used when updating the 3rd Generation (v3) App Service Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the 3rd Generation (v3) App Service Environment.
* `delete` - (Defaults to 6 hours) Used when deleting the 3rd Generation (v3) App Service Environment.

## Import

A 3rd Generation (v3) App Service Environment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_environment.myAppServiceEnv /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Web/hostingEnvironments/myAppServiceEnv
```
