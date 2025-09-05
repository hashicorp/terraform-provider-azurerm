---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_custom_hostname_binding"
description: |-
  Manages a Hostname Binding within an App Service (or Function App).

---

# azurerm_app_service_custom_hostname_binding

Manages a Hostname Binding within an App Service (or Function App).

## Example Usage

```hcl
resource "random_id" "server" {
  keepers = {
    azi_id = 1
  }

  byte_length = 8
}

resource "azurerm_resource_group" "example" {
  name     = "some-resource-group"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "example" {
  name                = "some-app-service-plan"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "example" {
  name                = random_id.server.hex
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id
}

resource "azurerm_app_service_custom_hostname_binding" "example" {
  hostname            = "www.mywebsite.com"
  app_service_name    = azurerm_app_service.example.name
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required) Specifies the Custom Hostname to use for the App Service, example `www.example.com`. Changing this forces a new resource to be created.

~> **Note:** A CNAME needs to be configured from this Hostname to the Azure Website - otherwise Azure will reject the Hostname Binding.

* `app_service_name` - (Required) The name of the App Service in which to add the Custom Hostname Binding. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the App Service exists. Changing this forces a new resource to be created.

* `ssl_state` - (Optional) The SSL type. Possible values are `IpBasedEnabled` and `SniEnabled`. Changing this forces a new resource to be created.

* `thumbprint` - (Optional) The SSL certificate thumbprint. Changing this forces a new resource to be created.

-> **Note:** `thumbprint` must be specified when `ssl_state` is set.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Service Custom Hostname Binding

* `virtual_ip` - The virtual IP address assigned to the hostname if IP based SSL is enabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Custom Hostname Binding.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Custom Hostname Binding.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Custom Hostname Binding.

## Import

App Service Custom Hostname Bindings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_custom_hostname_binding.mywebsite /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1/hostNameBindings/mywebsite.com
```
