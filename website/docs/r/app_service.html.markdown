---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service"
sidebar_current: "docs-azurerm-resource-app-service-x"
description: |-
  Manages an App Service (within an App Service Plan).

---

# azurerm_app_service

Manages an App Service (within an App Service Plan).

-> **Note:** When using Slots - the `app_settings`, `connection_string` and `site_config` blocks on the `azurerm_app_service` resource will be overwritten when promoting a Slot using the `azurerm_app_service_active_slot` resource.

## Example Usage

This example provisions a Windows App Service. Other examples of the `azurerm_app_service` resource can be found in [the `./examples/app-service` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/app-service)

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_app_service_plan" "test" {
  name                = "example-appserviceplan"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "example-app-service"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  app_service_plan_id = "${azurerm_app_service_plan.test.id}"

  site_config {
    dotnet_framework_version = "v4.0"
    scm_type                 = "LocalGit"
  }

  app_settings = {
    "SOME_KEY" = "some-value"
  }

  connection_string {
    name  = "Database"
    type  = "SQLServer"
    value = "Server=some-server.mydomain.com;Integrated Security=SSPI"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the App Service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the App Service.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `app_service_plan_id` - (Required) The ID of the App Service Plan within which to create this App Service.

* `app_settings` - (Optional) A key-value pair of App Settings.

* `connection_string` - (Optional) One or more `connection_string` blocks as defined below.

* `client_affinity_enabled` - (Optional) Should the App Service send session affinity cookies, which route client requests in the same session to the same instance?

* `client_cert_enabled` - (Optional) Does the App Service require client certificates for incoming requests? Defaults to `false`.

* `enabled` - (Optional) Is the App Service Enabled?

* `https_only` - (Optional) Can the App Service only be accessed via HTTPS? Defaults to `false`.

* `site_config` - (Optional) A `site_config` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `identity` - (Optional) A Managed Service Identity block as defined below.

---

A `connection_string` block supports the following:

* `name` - (Required) The name of the Connection String.

* `type` - (Required) The type of the Connection String. Possible values are `APIHub`, `Custom`, `DocDb`, `EventHub`, `MySQL`, `NotificationHub`, `PostgreSQL`, `RedisCache`, `ServiceBus`, `SQLAzure` and  `SQLServer`.

* `value` - (Required) The value for the Connection String.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the App Service. At this time the only allowed value is `SystemAssigned`.

~> The assigned `principal_id` and `tenant_id` can be retrieved after the App Service has been created. More details are available below.

---

A `site_config` block supports the following:

* `always_on` - (Optional) Should the app be loaded at all times? Defaults to `false`.

* `app_command_line` - (Optional) App command line to launch, e.g. `/sbin/myserver -b 0.0.0.0`.

* `cors` - (Optional) A `cors` block as defined below.

* `default_documents` - (Optional) The ordering of default documents to load, if an address isn't specified.

* `dotnet_framework_version` - (Optional) The version of the .net framework's CLR used in this App Service. Possible values are `v2.0` (which will use the latest version of the .net framework for the .net CLR v2 - currently `.net 3.5`) and `v4.0` (which corresponds to the latest version of the .net CLR v4 - which at the time of writing is `.net 4.7.1`). [For more information on which .net CLR version to use based on the .net framework you're targeting - please see this table](https://en.wikipedia.org/wiki/.NET_Framework_version_history#Overview). Defaults to `v4.0`.

* `ftps_state` - (Optional) State of FTP / FTPS service for this App Service. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`.

* `http2_enabled` - (Optional) Is HTTP2 Enabled on this App Service? Defaults to `false`.

* `ip_restriction` - (Optional) A [List of objects](/docs/configuration/attr-as-blocks.html) representing ip restrictions as defined below.

* `java_version` - (Optional) The version of Java to use. If specified `java_container` and `java_container_version` must also be specified. Possible values are `1.7`, `1.8` and `11`.

* `java_container` - (Optional) The Java Container to use. If specified `java_version` and `java_container_version` must also be specified. Possible values are `JETTY` and `TOMCAT`.

* `java_container_version` - (Optional) The version of the Java Container to use. If specified `java_version` and `java_container` must also be specified.

* `local_mysql_enabled` - (Optional) Is "MySQL In App" Enabled? This runs a local MySQL instance with your app and shares resources from the App Service plan.

~> **NOTE:** MySQL In App is not intended for production environments and will not scale beyond a single instance. Instead you may wish [to use Azure Database for MySQL](/docs/providers/azurerm/r/mysql_database.html).

* `linux_fx_version` - (Optional) Linux App Framework and version for the App Service. Possible options are a Docker container (`DOCKER|<user/image:tag>`), a base-64 encoded Docker Compose file (`COMPOSE|${base64encode(file("compose.yml"))}`) or a base-64 encoded Kubernetes Manifest (`KUBE|${base64encode(file("kubernetes.yml"))}`).

Additional examples of how to run Containers via the `azurerm_app_service` resource can be found in [the `./examples/app-service` directory within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/app-service).

* `managed_pipeline_mode` - (Optional) The Managed Pipeline Mode. Possible values are `Integrated` and `Classic`. Defaults to `Integrated`.

* `min_tls_version` - (Optional) The minimum supported TLS version for the app service. Possible values are `1.0`, `1.1`, and `1.2`. Defaults to `1.2` for new app services.

* `php_version` - (Optional) The version of PHP to use in this App Service. Possible values are `5.5`, `5.6`, `7.0`, `7.1` and `7.2`.

* `python_version` - (Optional) The version of Python to use in this App Service. Possible values are `2.7` and `3.4`.

* `remote_debugging_enabled` - (Optional) Is Remote Debugging Enabled? Defaults to `false`.

* `remote_debugging_version` - (Optional) Which version of Visual Studio should the Remote Debugger be compatible with? Possible values are `VS2012`, `VS2013`, `VS2015` and `VS2017`.

* `scm_type` - (Optional) The type of Source Control enabled for this App Service. Defaults to `None`. Possible values are: `BitbucketGit`, `BitbucketHg`, `CodePlexGit`, `CodePlexHg`, `Dropbox`, `ExternalGit`, `ExternalHg`, `GitHub`, `LocalGit`, `None`, `OneDrive`, `Tfs`, `VSO` and `VSTSRM`

* `use_32_bit_worker_process` - (Optional) Should the App Service run in 32 bit mode, rather than 64 bit mode?

~> **NOTE:** when using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `virtual_network_name` - (Optional) The name of the Virtual Network which this App Service should be attached to.

* `websockets_enabled` - (Optional) Should WebSockets be enabled?

---

A `cors` block supports the following:

* `allowed_origins` - (Optional) A list of origins which should be able to make cross-origin calls. `*` can be used to allow all calls.

* `support_credentials` - (Optional) Are credentials supported?

---

Elements of `ip_restriction` support:

* `ip_address` - (Required) The IP Address used for this IP Restriction.

* `subnet_mask` - (Optional) The Subnet mask used for this IP Restriction. Defaults to `255.255.255.255`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service.

* `default_site_hostname` - The Default Hostname associated with the App Service - such as `mysite.azurewebsites.net`

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12`

* `possible_outbound_ip_addresses` - A comma separated list of outbound IP addresses - such as `52.23.25.3,52.143.43.12,52.143.43.17` - not all of which are necessarily in use. Superset of `outbound_ip_addresses`.

* `source_control` - A `source_control` block as defined below, which contains the Source Control information when `scm_type` is set to `LocalGit`.

* `site_credential` - A `site_credential` block as defined below, which contains the site-level credentials used to publish to this App Service.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this App Service.

---

`identity` exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this App Service.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this App Service.

-> You can access the Principal ID via `${azurerm_app_service.test.identity.0.principal_id}` and the Tenant ID via `${azurerm_app_service.test.identity.0.tenant_id}`

---

`site_credential` exports the following:

* `username` - The username which can be used to publish to this App Service
* `password` - The password associated with the username, which can be used to publish to this App Service.

~> **NOTE:** both `username` and `password` for the `site_credential` block are only exported when `scm_type` is set to `LocalGit`

---

`source_control` exports the following:

* `repo_url` - URL of the Git repository for this App Service.
* `branch` - Branch name of the Git repository for this App Service.

## Import

App Services can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service.instance1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/instance1
```
