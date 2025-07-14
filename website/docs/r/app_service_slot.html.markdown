---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_app_service_slot"
description: |-
  Manages an App Service Slot (within an App Service).

---

# azurerm_app_service_slot

Manages an App Service Slot (within an App Service).

!> **Note:** This resource has been deprecated in version 3.0 of the AzureRM provider and will be removed in version 4.0. Please use [`azurerm_linux_web_app_slot`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/linux_web_app_slot) and [`azurerm_windows_web_app_slot`](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/windows_web_app_slot) resources instead.

-> **Note:** When using Slots - the `app_settings`, `connection_string` and `site_config` blocks on the `azurerm_app_service` resource will be overwritten when promoting a Slot using the `azurerm_app_service_active_slot` resource.

## Example Usage (.NET 4.x)

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

* `resource_group_name` - (Required) The name of the resource group in which to create the App Service Slot component. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `app_service_plan_id` - (Required) The ID of the App Service Plan within which to create this App Service Slot. Changing this forces a new resource to be created.

* `app_service_name` - (Required) The name of the App Service within which to create the App Service Slot. Changing this forces a new resource to be created.

* `app_settings` - (Optional) A key-value pair of App Settings.

* `auth_settings` - (Optional) A `auth_settings` block as defined below.

* `connection_string` - (Optional) An `connection_string` block as defined below.

* `client_affinity_enabled` - (Optional) Should the App Service Slot send session affinity cookies, which route client requests in the same session to the same instance?

* `enabled` - (Optional) Is the App Service Slot Enabled? Defaults to `true`.

* `https_only` - (Optional) Can the App Service Slot only be accessed via HTTPS? Defaults to `false`.

* `site_config` - (Optional) A `site_config` object as defined below.

* `storage_account` - (Optional) One or more `storage_account` blocks as defined below.

* `logs` - (Optional) A `logs` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `key_vault_reference_identity_id` - (Optional) The User Assigned Identity Id used for looking up KeyVault secrets. The identity must be assigned to the application. See [Access vaults with a user-assigned identity](https://docs.microsoft.com/azure/app-service/app-service-key-vault-references#access-vaults-with-a-user-assigned-identity) for more information.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `storage_account` block supports the following:

* `name` - (Required) The name of the storage account identifier.

* `type` - (Required) The type of storage. Possible values are `AzureBlob` and `AzureFiles`.

* `account_name` - (Required) The name of the storage account.

* `share_name` - (Required) The name of the file share (container name, for Blob storage).

* `access_key` - (Required) The access key for the storage account.

* `mount_path` - (Optional) The path to mount the storage within the site's runtime environment.

---

The `connection_string` block supports the following:

* `name` - (Required) The name of the Connection String.
* `type` - (Required) The type of the Connection String. Possible values are `APIHub`, `Custom`, `DocDb`, `EventHub`, `MySQL`, `NotificationHub`, `PostgreSQL`, `RedisCache`, `ServiceBus`, `SQLAzure`, and `SQLServer`.
* `value` - (Required) The value for the Connection String.

---

A `site_config` block supports the following:

* `acr_use_managed_identity_credentials` - (Optional) Are Managed Identity Credentials used for Azure Container Registry pull

* `acr_user_managed_identity_client_id` - (Optional) If using User Managed Identity, the User Managed Identity Client Id

~> **Note:** When using User Managed Identity with Azure Container Registry the Identity will need to have the [ACRPull role assigned](https://docs.microsoft.com/azure/container-registry/container-registry-authentication-managed-identity#example-1-access-with-a-user-assigned-identity)

* `always_on` - (Optional) Should the slot be loaded at all times? Defaults to `false`.

~> **Note:** when using an App Service Plan in the `Free` or `Shared` Tiers `always_on` must be set to `false`.

* `app_command_line` - (Optional) App command line to launch, e.g. `/sbin/myserver -b 0.0.0.0`.

* `auto_swap_slot_name` - (Optional) The name of the slot to automatically swap to during deployment

* `cors` - (Optional) A `cors` block as defined below.

* `default_documents` - (Optional) The ordering of default documents to load, if an address isn't specified.

* `dotnet_framework_version` - (Optional) The version of the .NET framework's CLR used in this App Service Slot. Possible values are `v2.0` (which will use the latest version of the .NET framework for the .NET CLR v2 - currently `.net 3.5`), `v4.0` (which corresponds to the latest version of the .NET CLR v4 - which at the time of writing is `.net 4.7.1`), `v5.0` and `v6.0`. [For more information on which .NET CLR version to use based on the .NET framework you're targeting - please see this table](https://en.wikipedia.org/wiki/.NET_Framework_version_history#Overview). Defaults to `v4.0`.

* `ftps_state` - (Optional) State of FTP / FTPS service for this App Service Slot. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`.

* `health_check_path` - (Optional) The health check path to be pinged by App Service Slot. [For more information - please see App Service health check announcement](https://azure.github.io/AppService/2020/08/24/healthcheck-on-app-service.html).

* `number_of_workers` - (Optional) The scaled number of workers (for per site scaling) of this App Service Slot. Requires that `per_site_scaling` is enabled on the `azurerm_app_service_plan`. [For more information - please see Microsoft documentation on high-density hosting](https://docs.microsoft.com/azure/app-service/manage-scale-per-app).

* `http2_enabled` - (Optional) Is HTTP2 Enabled on this App Service? Defaults to `false`.

* `ip_restriction` - (Optional) A list of `ip_restriction` objects representing IP restrictions as defined below.

-> **Note:** User has to explicitly set `ip_restriction` to empty slice (`[]`) to remove it.

* `scm_use_main_ip_restriction` - (Optional) IP security restrictions for scm to use main. Defaults to `false`. 

-> **Note:** Any `scm_ip_restriction` blocks configured are ignored by the service when `scm_use_main_ip_restriction` is set to `true`. Any scm restrictions will become active if this is subsequently set to `false` or removed.

* `scm_ip_restriction` - (Optional) A list of `scm_ip_restriction` objects representing IP restrictions as defined below.

-> **Note:** User has to explicitly set `scm_ip_restriction` to empty slice (`[]`) to remove it.

* `java_version` - (Optional) The version of Java to use. If specified `java_container` and `java_container_version` must also be specified. Possible values are `1.7`, `1.8`, and `11` and their specific versions - except for Java 11 (e.g. `1.7.0_80`, `1.8.0_181`, `11`)

* `java_container` - (Optional) The Java Container to use. If specified `java_version` and `java_container_version` must also be specified. Possible values are `JAVA`, `JETTY`, and `TOMCAT`.

* `java_container_version` - (Optional) The version of the Java Container to use. If specified `java_version` and `java_container` must also be specified.

* `local_mysql_enabled` - (Optional) Is "MySQL In App" Enabled? This runs a local MySQL instance with your app and shares resources from the App Service plan.

~> **Note:** MySQL In App is not intended for production environments and will not scale beyond a single instance. Instead you may wish [to use Azure Database for MySQL](/docs/providers/azurerm/r/mysql_database.html).

* `linux_fx_version` - (Optional) Linux App Framework and version for the App Service Slot. Possible options are a Docker container (`DOCKER|<user/image:tag>`), a base-64 encoded Docker Compose file (`COMPOSE|${filebase64("compose.yml")}`) or a base-64 encoded Kubernetes Manifest (`KUBE|${filebase64("kubernetes.yml")}`).

~> **Note:** To set this property the App Service Plan to which the App belongs must be configured with `kind = "Linux"`, and `reserved = true` or the API will reject any value supplied.

* `windows_fx_version` - (Optional) The Windows Docker container image (`DOCKER|<user/image:tag>`)

Additional examples of how to run Containers via the `azurerm_app_service_slot` resource can be found in [the `./examples/app-service` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/app-service).

* `managed_pipeline_mode` - (Optional) The Managed Pipeline Mode. Possible values are `Integrated` and `Classic`. Defaults to `Integrated`.

* `min_tls_version` - (Optional) The minimum supported TLS version for the app service. Possible values are `1.0`, `1.1`, and `1.2`. Defaults to `1.2` for new app services.

* `php_version` - (Optional) The version of PHP to use in this App Service Slot. Possible values are `5.5`, `5.6`, `7.0`, `7.1`, `7.2`, `7.3`, and `7.4`.

* `python_version` - (Optional) The version of Python to use in this App Service Slot. Possible values are `2.7` and `3.4`.

* `remote_debugging_enabled` - (Optional) Is Remote Debugging Enabled? Defaults to `false`.

* `remote_debugging_version` - (Optional) Which version of Visual Studio should the Remote Debugger be compatible with? Currently only `VS2022` is supported.

* `scm_type` - (Optional) The type of Source Control enabled for this App Service Slot. Defaults to `None`. Possible values are: `BitbucketGit`, `BitbucketHg`, `CodePlexGit`, `CodePlexHg`, `Dropbox`, `ExternalGit`, `ExternalHg`, `GitHub`, `LocalGit`, `None`, `OneDrive`, `Tfs`, `VSO`, and `VSTSRM`

* `use_32_bit_worker_process` - (Optional) Should the App Service Slot run in 32 bit mode, rather than 64 bit mode?

~> **Note:** when using an App Service Plan in the `Free` or `Shared` Tiers `use_32_bit_worker_process` must be set to `true`.

* `vnet_route_all_enabled` - (Optional) Should all outbound traffic to have Virtual Network Security Groups and User Defined Routes applied? Defaults to `false`.

~> **Note:** This setting supersedes the previous mechanism of setting the `app_settings` value of `WEBSITE_VNET_ROUTE_ALL`. However, to prevent older configurations breaking Terraform will update this value if it not explicitly set to the value in `app_settings.WEBSITE_VNET_ROUTE_ALL`.

* `websockets_enabled` - (Optional) Should WebSockets be enabled?

---

A `cors` block supports the following:

* `allowed_origins` - (Required) A list of origins which should be able to make cross-origin calls. `*` can be used to allow all calls.

* `support_credentials` - (Optional) Are credentials supported?

---

A `auth_settings` block supports the following:

* `enabled` - (Required) Is Authentication enabled?

* `active_directory` - (Optional) A `active_directory` block as defined below.

* `additional_login_params` - (Optional) Login parameters to send to the OpenID Connect authorization endpoint when a user logs in. Each parameter must be in the form "key=value".

* `allowed_external_redirect_urls` - (Optional) External URLs that can be redirected to as part of logging in or logging out of the app.

* `default_provider` - (Optional) The default provider to use when multiple providers have been set up. Possible values are `AzureActiveDirectory`, `Facebook`, `Google`, `MicrosoftAccount` and `Twitter`.

~> **Note:** When using multiple providers, the default provider must be set for settings like `unauthenticated_client_action` to work.

* `facebook` - (Optional) A `facebook` block as defined below.

* `google` - (Optional) A `google` block as defined below.

* `issuer` - (Optional) Issuer URI. When using Azure Active Directory, this value is the URI of the directory tenant, e.g. <https://sts.windows.net/{tenant-guid}/>.

* `microsoft` - (Optional) A `microsoft` block as defined below.

* `runtime_version` - (Optional) The runtime version of the Authentication/Authorization module.

* `token_refresh_extension_hours` - (Optional) The number of hours after session token expiration that a session token can be used to call the token refresh API. Defaults to `72`.

* `token_store_enabled` - (Optional) If enabled the module will durably store platform-specific security tokens that are obtained during login flows. Defaults to `false`.

* `twitter` - (Optional) A `twitter` block as defined below.

* `unauthenticated_client_action` - (Optional) The action to take when an unauthenticated client attempts to access the app. Possible values are `AllowAnonymous` and `RedirectToLoginPage`.

---

A `active_directory` block supports the following:

* `client_id` - (Required) The Client ID of this relying party application. Enables OpenIDConnection authentication with Azure Active Directory.

* `client_secret` - (Optional) The Client Secret of this relying party application. If no secret is provided, implicit flow will be used.

* `allowed_audiences` - (Optional) Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

---

A `facebook` block supports the following:

* `app_id` - (Required) The App ID of the Facebook app used for login

* `app_secret` - (Required) The App Secret of the Facebook app used for Facebook login.

* `oauth_scopes` - (Optional) The OAuth 2.0 scopes that will be requested as part of Facebook login authentication. <https://developers.facebook.com/docs/facebook-login>

---

A `twitter` block supports the following:

* `consumer_key` - (Required) The consumer key of the Twitter app used for login

* `consumer_secret` - (Required) The consumer secret of the Twitter app used for login.

---

A `google` block supports the following:

* `client_id` - (Required) The OpenID Connect Client ID for the Google web application.

* `client_secret` - (Required) The client secret associated with the Google web application.

* `oauth_scopes` - (Optional) The OAuth 2.0 scopes that will be requested as part of Google Sign-In authentication. <https://developers.google.com/identity/sign-in/web/>

---

A `ip_restriction` block supports the following:

* `ip_address` - (Optional) The IP Address used for this IP Restriction in CIDR notation.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

-> **Note:** One of either `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified

* `name` - (Optional) The name for this IP Restriction.

* `priority` - (Optional) The priority for this IP Restriction. Restrictions are enforced in priority order. By default, priority is set to 65000 if not specified.

* `action` - (Optional) Does this restriction `Allow` or `Deny` access for this IP range. Defaults to `Allow`. 

* `headers` - (Optional) The `headers` block for this specific `ip_restriction` as defined below. The HTTP header filters are evaluated after the rule itself and both conditions must be true for the rule to apply.

---

A `headers` block supports the following:

* `x_azure_fdid` - (Optional) A list of allowed Azure FrontDoor IDs in UUID notation with a maximum of 8.

* `x_fd_health_probe` - (Optional) A list to allow the Azure FrontDoor health probe header. Only allowed value is "1".

* `x_forwarded_for` - (Optional) A list of allowed 'X-Forwarded-For' IPs in CIDR notation with a maximum of 8

* `x_forwarded_host` - (Optional) A list of allowed 'X-Forwarded-Host' domains with a maximum of 8.

---

A `scm_ip_restriction` block supports the following:

* `ip_address` - (Optional) The IP Address used for this IP Restriction in CIDR notation.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

-> **Note:** One of either `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified

* `name` - (Optional) The name for this IP Restriction.

* `priority` - (Optional) The priority for this IP Restriction. Restrictions are enforced in priority order. By default, priority is set to 65000 if not specified.

* `action` - (Optional) Allow or Deny access for this IP range. Defaults to `Allow`.

* `headers` - (Optional) The `headers` block for this specific `scm_ip_restriction` as defined below.

---

A `microsoft` block supports the following:

* `client_id` - (Required) The OAuth 2.0 client ID that was created for the app used for authentication.

* `client_secret` - (Required) The OAuth 2.0 client secret that was created for the app used for authentication.

* `oauth_scopes` - (Optional) The OAuth 2.0 scopes that will be requested as part of Microsoft Account authentication. <https://msdn.microsoft.com/en-us/library/dn631845.aspx>

---

A `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the App Service. Possible values are `SystemAssigned` (where Azure will generate a Service Principal for you), `UserAssigned` where you can specify the Service Principal IDs in the `identity_ids` field, and `SystemAssigned, UserAssigned` which assigns both a system managed identity as well as the specified user assigned identities.

~> **Note:** When `type` is set to `SystemAssigned`, The assigned `principal_id` and `tenant_id` can be retrieved after the App Service has been created. More details are available below.

* `identity_ids` - (Optional) Specifies a list of user managed identity ids to be assigned. Required if `type` is `UserAssigned`.

---

A `logs` block supports the following:

* `application_logs` - (Optional) An `application_logs` block as defined below.

* `http_logs` - (Optional) An `http_logs` block as defined below.

* `detailed_error_messages_enabled` - (Optional) Should `Detailed error messages` be enabled on this App Service slot? Defaults to `false`.

* `failed_request_tracing_enabled` - (Optional) Should `Failed request tracing` be enabled on this App Service slot? Defaults to `false`.

---

An `application_logs` block supports the following:

* `file_system_level` - (Optional) The file system log level. Possible values are `Off`, `Error`, `Warning`, `Information`, and `Verbose`. Defaults to `Off`.

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

* `retention_in_mb` - (Required) The maximum size in megabytes that HTTP log files can use before being removed.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the App Service Slot.

* `default_site_hostname` - The Default Hostname associated with the App Service Slot - such as `mysite.azurewebsites.net`

* `site_credential` - A `site_credential` block as defined below, which contains the site-level credentials used to publish to this App Service slot.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this App Service slot.

---

A `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this App Service slot.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this App Service slot.

-> **Note:** You can access the Principal ID via `azurerm_app_service_slot.example.identity[0].principal_id` and the Tenant ID via `azurerm_app_service_slot.example.identity[0].tenant_id`

---

The `site_credential` block exports the following:

* `username` - The username which can be used to publish to this App Service
* `password` - The password associated with the username, which can be used to publish to this App Service.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the App Service Slot.
* `read` - (Defaults to 5 minutes) Used when retrieving the App Service Slot.
* `update` - (Defaults to 30 minutes) Used when updating the App Service Slot.
* `delete` - (Defaults to 30 minutes) Used when deleting the App Service Slot.

## Import

App Service Slots can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_app_service_slot.instance1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Web/sites/website1/slots/instance1
```
