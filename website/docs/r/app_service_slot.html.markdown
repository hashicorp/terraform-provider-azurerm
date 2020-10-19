---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_slot"
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

  site_config {
    dotnet_framework_version = "v4.0"
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

resource "azurerm_app_service_slot" "example" {
  name                = random_id.server.hex
  app_service_name    = azurerm_app_service.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id

  site_config {
    dotnet_framework_version = "v4.0"
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

## Example Usage (Java 1.8)

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

  site_config {
    java_version           = "1.8"
    java_container         = "JETTY"
    java_container_version = "9.3"
  }
}

resource "azurerm_app_service_slot" "example" {
  name                = random_id.server.hex
  app_service_name    = azurerm_app_service.example.name
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  app_service_plan_id = azurerm_app_service_plan.example.id

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

* `auth_settings` - (Optional) A `auth_settings` block as defined below.

* `connection_string` - (Optional) An `connection_string` block as defined below.

* `client_affinity_enabled` - (Optional) Should the App Service Slot send session affinity cookies, which route client requests in the same session to the same instance?

* `enabled` - (Optional) Is the App Service Slot Enabled?

* `https_only` - (Optional) Can the App Service Slot only be accessed via HTTPS? Defaults to `false`.

* `site_config` - (Optional) A `site_config` object as defined below.

* `identity` - (Optional) A Managed Service Identity block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`connection_string` supports the following:

* `name` - (Required) The name of the Connection String.
* `type` - (Required) The type of the Connection String. Possible values are `APIHub`, `Custom`, `DocDb`, `EventHub`, `MySQL`, `NotificationHub`, `PostgreSQL`, `RedisCache`, `ServiceBus`, `SQLAzure`, and  `SQLServer`.
* `value` - (Required) The value for the Connection String.

---

`site_config` supports the following:

* `app_command_line` - (Optional) App command line to launch, e.g. `/sbin/myserver -b 0.0.0.0`.

* `always_on` - (Optional) Should the app be loaded at all times? Defaults to `false`.

* `cors` - (Optional) A `cors` block as defined below.

* `default_documents` - (Optional) The ordering of default documents to load, if an address isn't specified.

* `dotnet_framework_version` - (Optional) The version of the .net framework's CLR used in this App Service Slot. Possible values are `v2.0` (which will use the latest version of the .net framework for the .net CLR v2 - currently `.net 3.5`) and `v4.0` (which corresponds to the latest version of the .net CLR v4 - which at the time of writing is `.net 4.7.1`). [For more information on which .net CLR version to use based on the .net framework you're targeting - please see this table](https://en.wikipedia.org/wiki/.NET_Framework_version_history#Overview). Defaults to `v4.0`.

* `http2_enabled` - (Optional) Is HTTP2 Enabled on this App Service? Defaults to `false`.

* `ip_restriction` - (Optional) A [List of objects](/docs/configuration/attr-as-blocks.html) representing ip restrictions as defined below.

-> **NOTE** User has to explicitly set `ip_restriction` to empty slice (`[]`) to remove it.

* `java_container` - (Optional) The Java Container to use. If specified `java_version` and `java_container_version` must also be specified. Possible values are `JETTY` and `TOMCAT`.

* `java_container_version` - (Optional) The version of the Java Container to use. If specified `java_version` and `java_container` must also be specified.

* `java_version` - (Optional) The version of Java to use. If specified `java_container` and `java_container_version` must also be specified. Possible values are `1.7`, `1.8`, and `11` and their specific versions - except for Java 11 (e.g. `1.7.0_80`, `1.8.0_181`, `11`)

* `local_mysql_enabled` - (Optional) Is "MySQL In App" Enabled? This runs a local MySQL instance with your app and shares resources from the App Service plan.

~> **NOTE:** MySQL In App is not intended for production environments and will not scale beyond a single instance. Instead you may wish [to use Azure Database for MySQL](/docs/providers/azurerm/r/mysql_database.html).

* `managed_pipeline_mode` - (Optional) The Managed Pipeline Mode. Possible values are `Integrated` and `Classic`. Defaults to `Integrated`.

* `min_tls_version` - (Optional) The minimum supported TLS version for the app service. Possible values are `1.0`, `1.1`, and `1.2`. Defaults to `1.2` for new app services.

* `php_version` - (Optional) The version of PHP to use in this App Service Slot. Possible values are `5.5`, `5.6`, `7.0`, `7.1`, `7.2`, and `7.3`.

* `python_version` - (Optional) The version of Python to use in this App Service Slot. Possible values are `2.7` and `3.4`.

* `remote_debugging_enabled` - (Optional) Is Remote Debugging Enabled? Defaults to `false`.

* `remote_debugging_version` - (Optional) Which version of Visual Studio should the Remote Debugger be compatible with? Possible values are `VS2012`, `VS2013`, `VS2015`, and `VS2017`.

* `scm_type` - (Optional) The type of Source Control enabled for this App Service Slot. Defaults to `None`. Possible values are: `BitbucketGit`, `BitbucketHg`, `CodePlexGit`, `CodePlexHg`, `Dropbox`, `ExternalGit`, `ExternalHg`, `GitHub`, `LocalGit`, `None`, `OneDrive`, `Tfs`, `VSO`, and `VSTSRM`

* `use_32_bit_worker_process` - (Optional) Should the App Service Slot run in 32 bit mode, rather than 64 bit mode?

~> **Note:** Deployment Slots are not supported in the `Free`, `Shared`, or `Basic` App Service Plans.

* `virtual_network_name` - (Optional) The name of the Virtual Network which this App Service Slot should be attached to.

* `websockets_enabled` - (Optional) Should WebSockets be enabled?

* `auto_swap_slot_name` - (Optional) The name of the slot to automatically swap to during deployment

---

A `cors` block supports the following:

* `allowed_origins` - (Optional) A list of origins which should be able to make cross-origin calls. `*` can be used to allow all calls.

* `support_credentials` - (Optional) Are credentials supported?

---

A `auth_settings` block supports the following:

* `enabled` - (Required) Is Authentication enabled?

* `active_directory` - (Optional) A `active_directory` block as defined below.

* `additional_login_params` - (Optional) Login parameters to send to the OpenID Connect authorization endpoint when a user logs in. Each parameter must be in the form "key=value".

* `allowed_external_redirect_urls` - (Optional) External URLs that can be redirected to as part of logging in or logging out of the app.

* `default_provider` - (Optional) The default provider to use when multiple providers have been set up. Possible values are `AzureActiveDirectory`, `Facebook`, `Google`, `MicrosoftAccount` and `Twitter`.

~> **NOTE:** When using multiple providers, the default provider must be set for settings like `unauthenticated_client_action` to work.

* `facebook` - (Optional) A `facebook` block as defined below.

* `google` - (Optional) A `google` block as defined below.

* `issuer` - (Optional) Issuer URI. When using Azure Active Directory, this value is the URI of the directory tenant, e.g. https://sts.windows.net/{tenant-guid}/.

* `microsoft` - (Optional) A `microsoft` block as defined below.

* `runtime_version` - (Optional) The runtime version of the Authentication/Authorization module.

* `token_refresh_extension_hours` - (Optional) The number of hours after session token expiration that a session token can be used to call the token refresh API. Defaults to 72.

* `token_store_enabled` - (Optional) If enabled the module will durably store platform-specific security tokens that are obtained during login flows. Defaults to false.

* `twitter` - (Optional) A `twitter` block as defined below.

* `unauthenticated_client_action` - (Optional) The action to take when an unauthenticated client attempts to access the app. Possible values are `AllowAnonymous` and `RedirectToLoginPage`.

---

A `active_directory` block supports the following:

* `client_id` - (Required) The Client ID of this relying party application. Enables OpenIDConnection authentication with Azure Active Directory.

* `client_secret` - (Optional) The Client Secret of this relying party application. If no secret is provided, implicit flow will be used.

* `allowed_audiences` (Optional) Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

---

A `facebook` block supports the following:

* `app_id` - (Required) The App ID of the Facebook app used for login

* `app_secret` - (Required) The App Secret of the Facebook app used for Facebook Login.

* `oauth_scopes` (Optional) The OAuth 2.0 scopes that will be requested as part of Facebook Login authentication. https://developers.facebook.com/docs/facebook-login

---

A `google` block supports the following:

* `client_id` - (Required) The OpenID Connect Client ID for the Google web application.

* `client_secret` - (Required) The client secret associated with the Google web application.

* `oauth_scopes` (Optional) The OAuth 2.0 scopes that will be requested as part of Google Sign-In authentication. https://developers.google.com/identity/sign-in/web/

---

A `ip_restriction` block supports the following:

* `ip_address` - (Optional) The IP Address used for this IP Restriction in CIDR notation.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

-> **NOTE:** One of either `ip_address` or `virtual_network_subnet_id` must be specified

* `name` - (Optional) The name for this IP Restriction.

* `priority` - (Optional) The priority for this IP Restriction. Restrictions are enforced in priority order. By default, priority is set to 65000 if not specified.

* `action` - (Optional) Does this restriction `Allow` or `Deny` access for this IP range. Defaults to `Allow`.  

---

A `microsoft` block supports the following:

* `client_id` - (Required) The OAuth 2.0 client ID that was created for the app used for authentication.

* `client_secret` - (Required) The OAuth 2.0 client secret that was created for the app used for authentication.

* `oauth_scopes` (Optional) The OAuth 2.0 scopes that will be requested as part of Microsoft Account authentication. https://msdn.microsoft.com/en-us/library/dn631845.aspx

---

A `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the App Service. Possible values are `SystemAssigned` (where Azure will generate a Service Principal for you), `UserAssigned` where you can specify the Service Principal IDs in the `identity_ids` field, and `SystemAssigned, UserAssigned` which assigns both a system managed identity as well as the specified user assigned identities.

~> **NOTE:** When `type` is set to `SystemAssigned`, The assigned `principal_id` and `tenant_id` can be retrieved after the App Service has been created. More details are available below.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned. Required if `type` is `UserAssigned`.

---

A `logs` block supports the following:

* `application_logs` - (Optional) An `application_logs` block as defined below.

* `http_logs` - (Optional) An `http_logs` block as defined below.

---

An `application_logs` block supports the following:

* `file_system_level` - (Optional) The file system log level. Possible values are `Off`, `Error`, `Warning`, `Information`, and `Verbose`.

* `azure_blob_storage` - (Optional) An `azure_blob_storage` block as defined below.

---

An `http_logs` block supports *one* of the following:

* `file_system` - (Optional) A `file_system` block as defined below.

* `azure_blob_storage` - (Optional) An `azure_blob_storage` block as defined below.

---

An `azure_blob_storage` block supports the following:

* `level` - (Required) The level at which to log. Possible values include `Error`, `Warning`, `Information`, `Verbose` and `Off`. **NOTE:** this field is not available for `http_logs`

* `sas_url` - (Required) The URL to the storage container, with a Service SAS token appended. **NOTE:** there is currently no means of generating Service SAS tokens with the `azurerm` provider.

* `retention_in_days` - (Required) The number of days to retain logs for.

---

A `file_system` block supports the following:

* `retention_in_days` - (Required) The number of days to retain logs for.

* `retention_in_mb` - (Required) The maximum size in megabytes that http log files can use before being removed.

---

Elements of `ip_restriction` [block](/docs/configuration/attr-as-blocks.html) support:

* `ip_address` - (Required) The IP Address used for this IP Restriction.

* `subnet_mask` - (Optional) The Subnet mask used for this IP Restriction. Defaults to `255.255.255.255`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the App Service Slot.

* `default_site_hostname` - The Default Hostname associated with the App Service Slot - such as `mysite.azurewebsites.net`

* `site_credential` - A `site_credential` block as defined below, which contains the site-level credentials used to publish to this App Service.

---

The `site_credential` block exports the following:

* `username` - The username which can be used to publish to this App Service
* `password` - The password associated with the username, which can be used to publish to this App Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Slot.
* `update` - (Defaults to 30 minutes) Used when updating the App Service Slot.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Slot.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Slot.

## Import

App Service Slots can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_slot.instance1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/website1/slots/instance1
```
