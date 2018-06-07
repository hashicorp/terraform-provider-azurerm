---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_slot"
sidebar_current: "docs-azurerm-resource-app-service-slot"
description: |-
  Manages an App Service Slot (within an App Service).

---

# azurerm_app_service_slot

Manages an App Service Slot (within an App Service).

-> **Note:** When using Slots - the `app_settings`, `connection_string` and `site_config` blocks on the `azurerm_app_service` resource will be overwritten when promoting a Slot using the `azurerm_app_service_active_slot` resource.


## Example Usage (.net 4.x)

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

  site_config {
    dotnet_framework_version = "v4.0"
  }

  app_settings {
    "SOME_KEY" = "some-value"
  }

  connection_string {
    name  = "Database"
    type  = "SQLServer"
    value = "Server=some-server.mydomain.com;Integrated Security=SSPI"
  }
}

resource "azurerm_app_service_slot" "test" {
  name                = "${random_id.server.hex}"
  app_service_name    = "${azurerm_app_service.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    dotnet_framework_version = "v4.0"
  }

  app_settings {
    "SOME_KEY" = "some-value"
  }

  connection_string {
    name  = "Database"
    type  = "SQLServer"
    value = "Server=some-server.mydomain.com;Integrated Security=SSPI"
  }
}
```

## Example Usage (Java 1.8)

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

  site_config {
    java_version           = "1.8"
    java_container         = "JETTY"
    java_container_version = "9.3"
  }
}

resource "azurerm_app_service_slot" "test" {
  name                = "${random_id.server.hex}"
  app_service_name    = "${azurerm_app_service.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    java_version           = "1.8"
    java_container         = "JETTY"
    java_container_version = "9.3"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the App Service Slot component. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the App Service Slot component.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `app_service_plan_id` - (Required) The ID of the App Service Plan within which to create this App Service Slot. Changing this forces a new resource to be created.

* `app_service_name` - (Required) The name of the App Service within which to create the App Service Slot.  Changing this forces a new resource to be created.

* `app_settings` - (Optional) A key-value pair of App Settings.

* `connection_string` - (Optional) An `connection_string` block as defined below.

* `client_affinity_enabled` - (Optional) Should the App Service Slot send session affinity cookies, which route client requests in the same session to the same instance? Changing this forces a new resource to be created.

* `enabled` - (Optional) Is the App Service Slot Enabled? Changing this forces a new resource to be created.

* `https_only` - (Optional) Can the App Service Slot only be accessed via HTTPS? Defaults to `false`. Changing this forces a new resource to be created.

* `site_config` - (Optional) A `site_config` object as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource. Changing this forces a new resource to be created.

---

`connection_string` supports the following:

* `name` - (Required) The name of the Connection String.
* `type` - (Required) The type of the Connection String. Possible values are `APIHub`, `Custom`, `DocDb`, `EventHub`, `MySQL`, `NotificationHub`, `PostgreSQL`, `RedisCache`, `ServiceBus`, `SQLAzure` and  `SQLServer`.
* `value` - (Required) The value for the Connection String.

---

`site_config` supports the following:

* `always_on` - (Optional) Should the app be loaded at all times? Defaults to `false`.
* `default_documents` - (Optional) The ordering of default documents to load, if an address isn't specified.
* `dotnet_framework_version` - (Optional) The version of the .net framework's CLR used in this App Service Slot. Possible values are `v2.0` (which will use the latest version of the .net framework for the .net CLR v2 - currently `.net 3.5`) and `v4.0` (which corresponds to the latest version of the .net CLR v4 - which at the time of writing is `.net 4.7.1`). [For more information on which .net CLR version to use based on the .net framework you're targeting - please see this table](https://en.wikipedia.org/wiki/.NET_Framework_version_history#Overview). Defaults to `v4.0`.
* `http2_enabled` - (Optional) Is HTTP2 Enabled on this App Service? Defaults to `false`.
* `ip_restriction` - (Optional) One or more `ip_restriction` blocks as defined below.
* `java_version` - (Optional) The version of Java to use. If specified `java_container` and `java_container_version` must also be specified. Possible values are `1.7` and `1.8`.
* `java_container` - (Optional) The Java Container to use. If specified `java_version` and `java_container_version` must also be specified. Possible values are `JETTY` and `TOMCAT`.
* `java_container_version` - (Optional) The version of the Java Container to use. If specified `java_version` and `java_container` must also be specified.

* `local_mysql_enabled` - (Optional) Is "MySQL In App" Enabled? This runs a local MySQL instance with your app and shares resources from the App Service plan.

~> **NOTE:** MySQL In App is not intended for production environments and will not scale beyond a single instance. Instead you may wish [to use Azure Database for MySQL](/docs/providers/azurerm/r/mysql_database.html).

* `managed_pipeline_mode` - (Optional) The Managed Pipeline Mode. Possible values are `Integrated` and `Classic`. Defaults to `Integrated`.
* `php_version` - (Optional) The version of PHP to use in this App Service Slot. Possible values are `5.5`, `5.6`, `7.0` and `7.1`.
* `python_version` - (Optional) The version of Python to use in this App Service Slot. Possible values are `2.7` and `3.4`.
* `remote_debugging_enabled` - (Optional) Is Remote Debugging Enabled? Defaults to `false`.
* `remote_debugging_version` - (Optional) Which version of Visual Studio should the Remote Debugger be compatible with? Possible values are `VS2012`, `VS2013`, `VS2015` and `VS2017`.
* `use_32_bit_worker_process` - (Optional) Should the App Service Slot run in 32 bit mode, rather than 64 bit mode?

~> **Note:** Deployment Slots are not supported in the `Free`, `Shared`, or `Basic` App Service Plans.

* `websockets_enabled` - (Optional) Should WebSockets be enabled?

---

`ip_restriction` supports the following:

* `ip_address` - (Required) The IP Address used for this IP Restriction.

* `subnet_mask` - (Optional) The Subnet mask used for this IP Restriction. Defaults to `255.255.255.255`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service Slot.

* `default_site_hostname` - The Default Hostname associated with the App Service Slot - such as `mysite.azurewebsites.net`

## Import

App Service Slots can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_slot.instance1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/website1/slots/instance1
```
