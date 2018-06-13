---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service"
sidebar_current: "docs-azurerm-datasource-app-service-x"
description: |-
  Get information about an App Service.
---

# Data Source: azurerm_app_service

Use this data source to obtain information about an App Service.

## Example Usage

```hcl
data "azurerm_app_service" "test" {
  name                = "search-app-service"
  resource_group_name = "search-service"
}

output "app_service_id" {
  value = "${data.azurerm_app_service.test.id}"
}
```

## Argument Reference

* `name` - (Required) The name of the App Service.

* `resource_group_name` - (Required) The Name of the Resource Group where the App Service exists.

## Attributes Reference

* `id` - The ID of the App Service.

* `location` - The Azure location where the App Service exists.

* `app_service_plan_id` - The ID of the App Service Plan within which the App Service exists.

* `app_settings` - A key-value pair of App Settings for the App Service.

* `connection_string` - An `connection_string` block as defined below.

* `client_affinity_enabled` - Does the App Service send session affinity cookies, which route client requests in the same session to the same instance?

* `enabled` - Is the App Service Enabled?

* `https_only` - Can the App Service only be accessed via HTTPS?

* `site_config` - A `site_config` block as defined below.

* `tags` - A mapping of tags to assign to the resource.

---

`connection_string` supports the following:

* `name` - The name of the Connection String.

* `type` - The type of the Connection String.

* `value` - The value for the Connection String.

---

`site_config` supports the following:

* `always_on` - Is the app be loaded at all times?

* `default_documents` - The ordering of default documents to load, if an address isn't specified.

* `dotnet_framework_version` - The version of the .net framework's CLR used in this App Service.

* `http2_enabled` - Is HTTP2 Enabled on this App Service?

* `ip_restriction` - One or more `ip_restriction` blocks as defined below.

* `java_version` - The version of Java in use.

* `java_container` - The Java Container in use.

* `java_container_version` - The version of the Java Container in use.

* `local_mysql_enabled` - Is "MySQL In App" Enabled? This runs a local MySQL instance with your app and shares resources from the App Service plan.

* `managed_pipeline_mode` - The Managed Pipeline Mode used in this App Service.

* `php_version` - The version of PHP used in this App Service.

* `python_version` - The version of Python used in this App Service.

* `remote_debugging_enabled` - Is Remote Debugging Enabled in this App Service?

* `remote_debugging_version` - Which version of Visual Studio is the Remote Debugger compatible with?

* `scm_type` - The type of Source Control enabled for this App Service.

* `use_32_bit_worker_process` - Does the App Service run in 32 bit mode, rather than 64 bit mode?

* `websockets_enabled` - Are WebSockets enabled for this App Service?

---

`ip_restriction` exports the following:

* `ip_address` - The IP Address used for this IP Restriction.

* `subnet_mask` - The Subnet mask used for this IP Restriction.
