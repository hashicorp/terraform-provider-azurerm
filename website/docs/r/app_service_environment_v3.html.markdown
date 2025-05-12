---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_environment_v3"
description: |-
  Manages a 3rd Generation (v3) App Service Environment.

---

# azurerm_app_service_environment_v3

Manages a 3rd Generation (v3) App Service Environment.

## Example Usage

This example provisions an App Service Environment V3. Additional examples of how to use the `azurerm_app_service_environment_v3` resource can be found [in the `./examples/app-service-environment-v3` directory](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/app-service-environment-v3) within the GitHub Repository.

```hcl
resource "azurerm_resource_group" "example" {
  name     = "exampleRG1"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example-subnet"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "Microsoft.Web.hostingEnvironments"
    service_delegation {
      name    = "Microsoft.Web/hostingEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_app_service_environment_v3" "example" {
  name                = "example-asev3"
  resource_group_name = azurerm_resource_group.example.name
  subnet_id           = azurerm_subnet.example.id

  internal_load_balancing_mode = "Web, Publishing"

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
    value = "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
  }

  tags = {
    env         = "production"
    terraformed = "true"
  }
}

resource "azurerm_service_plan" "example" {
  name                       = "example"
  resource_group_name        = azurerm_resource_group.example.name
  location                   = azurerm_resource_group.example.location
  os_type                    = "Linux"
  sku_name                   = "I1v2"
  app_service_environment_id = azurerm_app_service_environment_v3.example.id
}
```

## Argument Reference

* `name` - (Required) The name of the App Service Environment. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the App Service Environment exists. Defaults to the Resource Group of the Subnet (specified by `subnet_id`). Changing this forces a new resource to be created.

* `subnet_id` - (Required) The ID of the Subnet which the App Service Environment should be connected to. Changing this forces a new resource to be created.

~> **Note:** a /24 or larger CIDR is required. Once associated with an ASE, this size cannot be changed.

~> **Note:** This Subnet requires a delegation to `Microsoft.Web/hostingEnvironments` as detailed in the example above.

* `allow_new_private_endpoint_connections` - (Optional) Should new Private Endpoint Connections be allowed. Defaults to `true`.

* `cluster_setting` - (Optional) Zero or more `cluster_setting` blocks as defined below.

* `dedicated_host_count` - (Optional) This ASEv3 should use dedicated Hosts. Possible values are `2`. Changing this forces a new resource to be created.

* `remote_debugging_enabled` - (Optional) Whether to enable remote debug. Defaults to `false`.

* `zone_redundant` - (Optional) Set to `true` to deploy the ASEv3 with availability zones supported. Zonal ASEs can be deployed in some regions, you can refer to [Availability Zone support for App Service Environments](https://docs.microsoft.com/azure/app-service/environment/zone-redundancy). You can only set either `dedicated_host_count` or `zone_redundant` but not both. Changing this forces a new resource to be created.

~> **Note:** Setting this value will provision 2 Physical Hosts for your App Service Environment V3, this is done at additional cost, please be aware of the pricing commitment in the [General Availability Notes](https://techcommunity.microsoft.com/t5/apps-on-azure/announcing-app-service-environment-v3-ga/ba-p/2517990)

* `internal_load_balancing_mode` - (Optional) Specifies which endpoints to serve internally in the Virtual Network for the App Service Environment. Possible values are `None` (for an External VIP Type), and `"Web, Publishing"` (for an Internal VIP Type). Defaults to `None`. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

~> **Note:** The underlying API does not currently support changing Tags on this resource. Making changes in the portal for tags will cause Terraform to detect a change that will force a recreation of the ASEV3 unless `ignore_changes` lifecycle meta-argument is used.

---

A `cluster_setting` block supports the following:

~> **Note:** If this block is specified it must contain the `FrontEndSSLCipherSuiteOrder` setting, with the value `TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`.

* `name` - (Required) The name of the Cluster Setting.

* `value` - (Required) The value for the Cluster Setting.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Service Environment.

* `dns_suffix` - the DNS suffix for this App Service Environment V3.

* `external_inbound_ip_addresses` - The external inbound IP addresses of the App Service Environment V3.

* `inbound_network_dependencies` - An `inbound_network_dependencies` block as defined below.

* `internal_inbound_ip_addresses` - The internal inbound IP addresses of the App Service Environment V3.

* `ip_ssl_address_count` - The number of IP SSL addresses reserved for the App Service Environment V3.

* `linux_outbound_ip_addresses` - Outbound addresses of Linux based Apps in this App Service Environment V3

* `location` - The location where the App Service Environment exists.

* `pricing_tier` - Pricing tier for the front end instances.

* `windows_outbound_ip_addresses` - Outbound addresses of Windows based Apps in this App Service Environment V3.

---

An `inbound_network_dependencies` block exports the following:

* `description` - A short description of the purpose of the network traffic.

* `ip_addresses` - A list of IP addresses that network traffic will originate from in CIDR notation.

* `ports` - The ports that network traffic will arrive to the App Service Environment V3 on.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 6 hours) Used when creating the 3rd Generation (v3) App Service Environment.
* `read` - (Defaults to 5 minutes) Used when retrieving the 3rd Generation (v3) App Service Environment.
* `update` - (Defaults to 6 hours) Used when updating the 3rd Generation (v3) App Service Environment.
* `delete` - (Defaults to 6 hours) Used when deleting the 3rd Generation (v3) App Service Environment.

## Import

A 3rd Generation (v3) App Service Environment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_environment_v3.myAppServiceEnv /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myResourceGroup/providers/Microsoft.Web/hostingEnvironments/myAppServiceEnv
```
