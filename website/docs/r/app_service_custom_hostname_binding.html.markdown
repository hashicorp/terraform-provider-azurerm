---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_custom_hostname_binding"
sidebar_current: "docs-azurerm-resource-app-service-custom-hostname-binding"
description: |-
  Manages a Hostname Binding within an App Service.

---

# azurerm_app_service_custom_hostname_binding

Manages a Hostname Binding within an App Service.

## Example Usage

```hcl
resource "random_id" "server" {
  keepers = {
    azi_id = 1
  }

  byte_length = 8
}

resource "azurerm_resource_group" "test" {
  name     = "some-resource-group"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "test" {
  name                = "some-app-service-plan"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "${random_id.server.hex}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"
}

resource "azurerm_app_service_custom_hostname_binding" "test" {
  hostname            = "www.mywebsite.com"
  app_service_name    = "${azurerm_app_service.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
```

## Argument Reference

The following arguments are supported:

* `hostname` - (Required) Specifies the Custom Hostname to use for the App Service, example `www.example.com`. Changing this forces a new resource to be created.

~> **NOTE:** A CNAME needs to be configured from this Hostname to the Azure Website - otherwise Azure will reject the Hostname Binding.

* `app_service_name` - (Required) The name of the App Service in which to add the Custom Hostname Binding. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the App Service exists. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service Custom Hostname Binding

## Import

App Service Custom Hostname Bindings can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_custom_hostname_binding.mywebsite /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1/hostNameBindings/mywebsite.com
```
