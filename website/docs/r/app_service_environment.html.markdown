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
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "ase" {
  name                 = "asesubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_subnet" "gateway" {
  name                 = "gatewaysubnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]
}

resource "azurerm_app_service_environment" "example" {
  name                         = "example-ase"
  subnet_id                    = azurerm_subnet.ase.id
  pricing_tier                 = "I2"
  front_end_scale_factor       = 10
  internal_load_balancing_mode = "Web, Publishing"
  allowed_user_ip_cidrs        = ["11.22.33.44/32", "55.66.77.0/24"]

  cluster_setting {
    name  = "DisableTls1.0"
    value = "1"
  }
}

```

## Argument Reference

* `name` - (Required) The name of the App Service Environment. Changing this forces a new resource to be created. 

* `subnet_id` - (Required) The ID of the Subnet which the App Service Environment should be connected to. Changing this forces a new resource to be created.

~> **NOTE:** a /24 or larger CIDR is required. Once associated with an ASE this size cannot be changed.

~> **NOTE:** When using `ASEV3` the Subnet used for the ASE must have a delegation for `"Microsoft.Web/hostingEnvironments"` configured.

* `cluster_setting` - (Optional) Zero or more `cluster_setting` blocks as defined below. 

* `internal_load_balancing_mode` - (Optional) Specifies which endpoints to serve internally in the Virtual Network for the App Service Environment. Possible values are `None`, `Web`, `Publishing` and combined value `"Web, Publishing"`. Defaults to `None`.

* `pricing_tier` - (Optional) Pricing tier for the front end instances. Possible values are `I1`, `I2` and `I3`. Defaults to `I1`.

* `front_end_scale_factor` - (Optional) Scale factor for front end instances. Possible values are between `5` and `15`. Defaults to `15`.

* `allowed_user_ip_cidrs` - (Optional) Allowed user added IP ranges on the ASE database. Use the addresses you want to set as the explicit egress address ranges.

~> **NOTE:** `allowed_user_ip_cidrs` The addresses that will be used for all outbound traffic from your App Service Environment to the internet to avoid asymmetric routing challenge. If you're routing the traffic on premises, these addresses are your NATs or gateway IPs. If you want to route the App Service Environment outbound traffic through an NVA, the egress address is the public IP of the NVA. Please visit [Create your ASE with the egress addresses](https://docs.microsoft.com/en-us/azure/app-service/environment/forced-tunnel-support#add-your-own-ips-to-the-ase-azure-sql-firewall)

* `resource_group_name` - (Optional) The name of the Resource Group where the App Service Environment exists. Defaults to the Resource Group of the Subnet (specified by `subnet_id`).

* `version` - (Optional) The Version of the Application Service Environment. Possible values are `ASEV2` and `ASEV3`. Defaults to `ASEV2`.

~> **NOTE:** When using `ASEV3` values for the properties `allowed_user_ip_cidrs`, and `pricing_tier` are not allowed.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created. 

---

A `cluster_setting` block supports the following:

* `name` - (Required) The name of the Cluster Setting. 

* `value` - (Required) The value for the Cluster Setting. 

## Attribute Reference

* `id` - The ID of the App Service Environment.

* `location` - The location where the App Service Environment exists.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 4 hours) Used when creating the App Service Environment.
* `update` - (Defaults to 4 hours) Used when updating the App Service Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Environment.
* `delete` - (Defaults to 4 hours) Used when deleting the App Service Environment.

## Import

The App Service Environment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_environment.myAppServiceEnv /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Web/hostingEnvironments/myAppServiceEnv
```
