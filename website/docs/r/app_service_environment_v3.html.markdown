---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_environment_v3"
description: |-
  Manages a 3rd Generation (v3) App Service Environment.

---

# azurerm_app_service_environment

Manages a 3rd Generation (v3) App Service Environment.

## Example Usage

This example provisions an App Service Environment V3. Additional examples of how to use the azurerm_app_service_environment_v3 resource can be found in the ./examples/app-service-environment-v3 directory within the Github Repository.

```hcl
# terraform.tfvars
ase_resource_group_name               = "rg-asev3-terraform-ilb-demo"
use_existing_vnet_and_subnet          = false
vnet_resource_group_name              = "rg-asev3-terraform-demo"
virtual_network_name                  = "vnet-spoke-asev3"
location                              = "West US 2"
vnet_address_prefixes                 = ["172.16.0.0/16"]
subnet_name                           = "snet-asev3"
subnet_address_prefixes               = ["172.16.0.0/24"]
ase_name                              = "asev3-ilb-20211001-2"
dedicated_host_count                  = 0
zone_redundant                        = false
create_private_dns                    = true
internal_load_balancing_mode          = "Web, Publishing"
network_security_group_name           = "nsg-asev3"
network_security_group_security_rules = []

```
```hcl
# variables.tf

variable "ase_resource_group_name" {
  type = string
}

variable "use_existing_vnet_and_subnet" {
  type    = bool
  default = false
}

variable "vnet_resource_group_name" {
  type = string
}

variable "virtual_network_name" {
  type = string
}

variable "location" {
  type = string
}

variable "vnet_address_prefixes" {
  type    = list(string)
  default = ["172.16.0.0/16"]
}

variable "subnet_name" {
  type = string
}

variable "subnet_address_prefixes" {
  type    = list(string)
  default = ["172.16.0.0/24"]
}

variable "ase_name" {
  type = string
}

variable "dedicated_host_count" {
  type    = number
  default = 0
}

variable "zone_redundant" {
  type    = bool
  default = false
}

variable "create_private_dns" {
  type    = bool
  default = true
}

variable "internal_load_balancing_mode" {
  type    = string
  default = "Web, Publishing"
}

variable "network_security_group_name" {
  type = string
}

variable "network_security_group_security_rules" {
  type    = any
  default = []
}

```

```hcl
# main.tf
resource "azurerm_resource_group" "rg" {
  name     = var.ase_resource_group_name
  location = var.location
}

data "azurerm_subnet" "snet-exist" {
  count                = var.use_existing_vnet_and_subnet ? 1 : 0
  name                 = var.subnet_name
  virtual_network_name = var.virtual_network_name
  resource_group_name  = var.vnet_resource_group_name
}

resource "azurerm_network_security_group" "nsg" {
  count               = var.use_existing_vnet_and_subnet ? 0 : 1
  name                = var.network_security_group_name
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_virtual_network" "vnet" {
  count               = var.use_existing_vnet_and_subnet ? 0 : 1
  name                = var.virtual_network_name
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  address_space       = var.vnet_address_prefixes
}

resource "azurerm_subnet" "snet" {
  count                = var.use_existing_vnet_and_subnet ? 0 : 1
  name                 = var.subnet_name
  resource_group_name  = azurerm_resource_group.rg.name
  virtual_network_name = azurerm_virtual_network.vnet[0].name
  address_prefixes     = var.subnet_address_prefixes

  delegation {
    name = "Microsoft.Web.hostingEnvironments"
    service_delegation {
      name    = "Microsoft.Web/hostingEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }  
}

resource "azurerm_subnet" "snet-delegation" {
  count                = var.use_existing_vnet_and_subnet ? 1 : 0
  name                 = var.subnet_name
  virtual_network_name = var.virtual_network_name
  resource_group_name  = var.vnet_resource_group_name
  address_prefixes     = data.azurerm_subnet.snet-exist[0].address_prefixes

  delegation {
    name = "Microsoft.Web.hostingEnvironments"
    service_delegation {
      name    = "Microsoft.Web/hostingEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "nsg-association" {
  count                     = var.use_existing_vnet_and_subnet ? 0 : 1
  subnet_id                 = azurerm_subnet.snet[0].id
  network_security_group_id = azurerm_network_security_group.nsg[0].id
}

resource "azurerm_app_service_environment_v3" "asev3" {
  name                = var.ase_name
  resource_group_name = azurerm_resource_group.rg.name
  subnet_id           = var.use_existing_vnet_and_subnet ? azurerm_subnet.snet-delegation[0].id : azurerm_subnet.snet[0].id

  dedicated_host_count         = var.dedicated_host_count >= 2 ? var.dedicated_host_count : null
  zone_redundant               = var.zone_redundant ? var.zone_redundant : null
  internal_load_balancing_mode = var.internal_load_balancing_mode
}

resource "azurerm_private_dns_zone" "zone" {
  count               = (var.create_private_dns && var.internal_load_balancing_mode == "Web, Publishing") ? 1 : 0
  name                = azurerm_app_service_environment_v3.asev3.dns_suffix
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurerm_private_dns_a_record" "a-wildcard" {
  count               = (var.create_private_dns && var.internal_load_balancing_mode == "Web, Publishing") ? 1 : 0
  name                = "*"
  zone_name           = azurerm_private_dns_zone.zone[0].name
  resource_group_name = azurerm_resource_group.rg.name
  ttl                 = 3600
  records             = azurerm_app_service_environment_v3.asev3.internal_inbound_ip_addresses
}

resource "azurerm_private_dns_a_record" "a-scm" {
  count               = (var.create_private_dns && var.internal_load_balancing_mode == "Web, Publishing") ? 1 : 0
  name                = "*.scm"
  zone_name           = azurerm_private_dns_zone.zone[0].name
  resource_group_name = azurerm_resource_group.rg.name
  ttl                 = 3600
  records             = azurerm_app_service_environment_v3.asev3.internal_inbound_ip_addresses
}

resource "azurerm_private_dns_a_record" "a-at" {
  count               = (var.create_private_dns && var.internal_load_balancing_mode == "Web, Publishing") ? 1 : 0
  name                = "@"
  zone_name           = azurerm_private_dns_zone.zone[0].name
  resource_group_name = azurerm_resource_group.rg.name
  ttl                 = 3600
  records             = azurerm_app_service_environment_v3.asev3.internal_inbound_ip_addresses
}

```

## Argument Reference

* `ase_resource_group_name` - (Required) The name of the Resource Group where the App Service Environment exists. 

* `location` - (Required) The location of the Resource Group where the App Service Environment exists. 

* `use_existing_vnet_and_subnet` - (Required) Set to `true` if the ASEv3 will be resided in an existing VNet and subnet.

* `vnet_resource_group_name` - (Optional) The name of the Resource Group where the existing VNet exists. Only required if `use_existing_vnet_and_subnet` set to `true`. 

* `virtual_network_name` - (Required) The name of the VNet where the App Service Environment exists.

* `vnet_address_prefixes` - (Optional) The address spaces of VNet where the App Service Environment exists. Only required if `use_existing_vnet_and_subnet` set to `false`.

* `subnet_name` - (Required) The name of the subnet where the App Service Environment exists.

* `subnet_address_prefixes` - (Optional)  The address spaces of subnet where the App Service Environment exists. Only required if `use_existing_vnet_and_subnet` set to `false`.

* `ase_name` - (Required) The name of the App Service Environment. This is a global Azure resource name should be unique.

* `dedicated_host_count` - (Optional) This ASEv3 should use dedicated Hosts. Possible vales are `2`. Changing this forces a new resource to be created. You can only set either `dedicated_host_count` or `zone_redundant` but not both.

* `zone_redundant` - (Optional) Set to `true` to deplyed ASEv3 with availability zones supported. Zonal ASEs can be deployed in some regions, you can refer to [Availability Zone support for App Service Environments](https://docs.microsoft.com/en-us/azure/app-service/environment/zone-redundancy). You can only set either `dedicated_host_count` or `zone_redundant` but not both.

~> **NOTE:** Setting this value will provision 2 Physical Hosts for your App Service Environment V3, this is done at additional cost, please be aware of the pricing commitment in the [General Availability Notes](https://techcommunity.microsoft.com/t5/apps-on-azure/announcing-app-service-environment-v3-ga/ba-p/2517990)

* `create_private_dns` - (Required) Set to `true` to create private dns zone after ASEv3 had been created. Only available for internal load balancing mode ASE.

* `internal_load_balancing_mode` - (Required) Specifies which endpoints to serve internally in the Virtual Network for the App Service Environment. Possible values are `None` (for an External VIP Type), and `"Web, Publishing"` (for an Internal VIP Type). Defaults to `None`.

* `network_security_group_name` - (Optional) The network security group name of the subnet where App Service Environment exists.

* `network_security_group_security_rules` - (Optional) The network security group rules of the subnet network security group where App Service Environment exists.

~> **NOTE** a /24 or larger CIDR is required. Once associated with an ASE, this size cannot be changed.

~> **NOTE:** This Subnet requires a delegation to `Microsoft.Web/hostingEnvironments` as detailed in the example above.    

---

## Attribute Reference

* `allow_new_private_endpoint_connections` - New Private Endpoint Connections be allowed. Defaults to `true`.

* `dns_suffix` - the DNS suffix for this App Service Environment V3. 

* `external_inbound_ip_addresses` - The external outbound IP addresses of the App Service Environment V3.

* `id` - The ID of the App Service Environment.

* `inbound_network_dependencies` - An inbound Network Dependencies block as defined below.

* `internal_inbound_ip_addresses` - The internal outbound IP addresses of the App Service Environment V3.

* `internal_load_balancing_mode` - Endpoints to serve internally in the Virtual Network for the App Service Environment. Possible values are `None` (for an External VIP Type), and `"Web, Publishing"` (for an Internal VIP Type). Defaults to `None`.

* `ip_ssl_address_count` - The number of IP SSL addresses reserved for the App Service Environment V3.

* `linux_outbound_ip_addresses` - Outbound addresses of Linux based Apps in this App Service Environment V3

* `location` - The location where the App Service Environment exists.

* `name` - The name of the App Service Environment V3.

* `pricing_tier` - Pricing tier for the front end instances.

* `resource_group_name` - The name of the Resource Group where the App Service Environment exists.

* `subnet_id` - The ID of the Subnet which the App Service Environment connected.

* `windows_outbound_ip_addresses` - Outbound addresses of Windows based Apps in this App Service Environment V3. 

* `zone_redundant` - ASEv3 deployed in availability zones or not.

--- 
A `cluster_setting` block supports the following:

~> **NOTE:** If this block is specified it must contain the `FrontEndSSLCipherSuiteOrder` setting, with the value `TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`.

* `name` - (Required) The name of the Cluster Setting. 

* `value` - (Required) The value for the Cluster Setting. 


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
