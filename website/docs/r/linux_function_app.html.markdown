---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_linux_function_app"
description: |-
  Manages a Linux Function App.
---

# azurerm_linux_function_app

Manages a Linux Function App.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "linuxfunctionappsa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "example" {
  name                = "example-app-service-plan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Linux"
  sku_name            = "Y1"
}

resource "azurerm_linux_function_app" "example" {
  name                = "example-linux-function-app"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  storage_account_name = azurerm_storage_account.example.name
  service_plan_id      = azurerm_service_plan.example.id

  site_config {}
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Linux Function App should exist. Changing this forces a new Linux Function App to be created.

* `name` - (Required) The name which should be used for this Linux Function App. Changing this forces a new Linux Function App to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Linux Function App should exist. Changing this forces a new Linux Function App to be created.

* `service_plan_id` - (Required) The ID of the App Service Plan within which to create this Function App.

* `site_config` - (Required) A `site_config` block as defined below.

---

* `app_settings` - (Optional) A map of key-value pairs for [App Settings](https://docs.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.

* `auth_settings` - (Optional) A `auth_settings` block as defined below.

* `backup` - (Optional) A `backup` block as defined below.

* `builtin_logging_enabled` - (Optional) Should built in logging be enabled. Configures `AzureWebJobsDashboard` app setting based on the configured storage setting.

* `client_certificate_enabled` - (Optional) Should the function app use Client Certificates.

* `client_certificate_mode` - (Optional) The mode of the Function App's client certificates requirement for incoming requests. Possible values are `Required`, `Optional`, and `OptionalInteractiveUser`.

* `connection_string` - (Optional) One or more `connection_string` blocks as defined below.

* `daily_memory_time_quota` - (Optional) The amount of memory in gigabyte-seconds that your application is allowed to consume per day. Setting this value only affects function apps under the consumption plan. Defaults to `0`.

* `enabled` - (Optional) Is the Function App enabled?

* `content_share_force_disabled` - (Optional) Should the settings for linking the Function App to storage be suppressed. 

* `functions_extension_version` - (Optional) The runtime version associated with the Function App. Defaults to `~4`.

* `https_only` - (Optional) Can the Function App only be accessed via HTTPS? Defaults to `false`.

* `identity` - (Optional) A `identity` block as defined below.

* `key_vault_reference_identity_id` - (Optional) The User Assigned Identity ID used for accessing KeyVault secrets. The identity must be assigned to the application in the `identity` block. [For more information see - Access vaults with a user-assigned identity](https://docs.microsoft.com/azure/app-service/app-service-key-vault-references#access-vaults-with-a-user-assigned-identity)

* `sticky_settings` - A `sticky_settings` block as defined below.

* `storage_account_access_key` - (Optional) The access key which will be used to access the backend storage account for the Function App. Conflicts with `storage_uses_managed_identity`.

* `storage_account_name` - (Optional) The backend storage account name which will be used by this Function App.

* `storage_uses_managed_identity` - (Optional) Should the Function App use Managed Identity to access the storage account. Conflicts with `storage_account_access_key`.

~> **NOTE:** One of `storage_account_access_key` or `storage_uses_managed_identity` must be specified when using `storage_account_name`.

* `storage_key_vault_secret_id` - (Optional) The Key Vault Secret ID, optionally including version, that contains the Connection String to connect to the storage account for this Function App.

~> **NOTE:** `storage_key_vault_secret_id` cannot be used with `storage_account_name`. 

~> **NOTE:** `storage_key_vault_secret_id` used without a version will use the latest version of the secret, however, the service can take up to 24h to pick up a rotation of the latest version. See the [official docs](https://docs.microsoft.com/azure/app-service/app-service-key-vault-references#rotation) for more information.

* `tags` - (Optional) A mapping of tags which should be assigned to the Linux Function App.

* `virtual_network_subnet_id` - (Optional) The subnet id which will be used by this Function App for [regional virtual network integration](https://docs.microsoft.com/en-us/azure/app-service/overview-vnet-integration#regional-virtual-network-integration).

~> **NOTE on regional virtual network integration:** The AzureRM Terraform provider provides regional virtual network integration via the standalone resource [app_service_virtual_network_swift_connection](app_service_virtual_network_swift_connection.html) and in-line within this resource using the `virtual_network_subnet_id` property. You cannot use both methods simutaneously.

~> **Note:** Assigning the `virtual_network_subnet_id` property requires [RBAC permissions on the subnet](https://docs.microsoft.com/en-us/azure/app-service/overview-vnet-integration#permissions)

---

An `active_directory` block supports the following:

* `client_id` - (Required) The ID of the Client to use to authenticate with Azure Active Directory.

* `allowed_audiences` - (Optional) Specifies a list of Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

~> **Note:** The `client_id` value is always considered an allowed audience.

* `client_secret` - (Optional) The Client Secret for the Client ID. Cannot be used with `client_secret_setting_name`.

* `client_secret_setting_name` - (Optional) The App Setting name that contains the client secret of the Client. Cannot be used with `client_secret`.

---

A `application_stack` block supports the following:

* `docker` - (Optional) One or more `docker` blocks as defined below.

* `dotnet_version` - (Optional) The version of .NET to use. Possible values include `3.1` and `6.0`.

* `use_dotnet_isolated_runtime` - (Optional) Should the DotNet process use an isolated runtime. Defaults to `false`.

* `java_version` - (Optional) The Version of Java to use. Supported versions include `8`, and `11`.

* `node_version` - (Optional) The version of Node to run. Possible values include `12`, `14`, and `16`.

* `python_version` - (Optional) The version of Python to run. Possible values include `3.6`, `3.7`, `3.8`, and `3.9`.

* `powershell_core_version` - (Optional) The version of PowerShell Core to run. Possible values are `7`, and `7.2`.

* `use_custom_runtime` - (Optional) Should the Linux Function App use a custom runtime?

---

An `app_service_logs` block supports the following:

* `disk_quota_mb` - (Required) The amount of disk space to use for logs. Valid values are between `25` and `100`. 

* `retention_period_days` - (Optional) The retention period for logs in days. Valid values are between `0` and `99999`. Defaults to `0` (never delete). 

~> **NOTE:** This block is not supported on Consumption plans.

---

An `auth_settings` block supports the following:

* `enabled` - (Required) Should the Authentication / Authorization feature be enabled for the Linux Web App?

* `active_directory` - (Optional) An `active_directory` block as defined above.

* `additional_login_parameters` - (Optional) Specifies a map of login Parameters to send to the OpenID Connect authorization endpoint when a user logs in.

* `allowed_external_redirect_urls` - (Optional) Specifies a list of External URLs that can be redirected to as part of logging in or logging out of the Linux Web App.

* `default_provider` - (Optional) The default authentication provider to use when multiple providers are configured. Possible values include: `AzureActiveDirectory`, `Facebook`, `Google`, `MicrosoftAccount`, `Twitter`, `Github`

~> **NOTE:** This setting is only needed if multiple providers are configured, and the `unauthenticated_client_action` is set to "RedirectToLoginPage".

* `facebook` - (Optional) A `facebook` block as defined below.

* `github` - (Optional) A `github` block as defined below.

* `google` - (Optional) A `google` block as defined below.

* `issuer` - (Optional) The OpenID Connect Issuer URI that represents the entity which issues access tokens for this Linux Web App.

~> **NOTE:** When using Azure Active Directory, this value is the URI of the directory tenant, e.g. https://sts.windows.net/{tenant-guid}/.

* `microsoft` - (Optional) A `microsoft` block as defined below.

* `runtime_version` - (Optional) The RuntimeVersion of the Authentication / Authorization feature in use for the Linux Web App.

* `token_refresh_extension_hours` - (Optional) The number of hours after session token expiration that a session token can be used to call the token refresh API. Defaults to `72` hours.

* `token_store_enabled` - (Optional) Should the Linux Web App durably store platform-specific security tokens that are obtained during login flows? Defaults to `false`.

* `twitter` - (Optional) A `twitter` block as defined below.

* `unauthenticated_client_action` - (Optional) The action to take when an unauthenticated client attempts to access the app. Possible values include: `RedirectToLoginPage`, `AllowAnonymous`.

---

A `backup` block supports the following:

* `name` - (Required) The name which should be used for this Backup.

* `schedule` - (Required) A `schedule` block as defined below.

* `storage_account_url` - (Required) The SAS URL to the container.

* `enabled` - (Optional) Should this backup job be enabled?

---

A `connection_string` block supports the following:

* `name` - (Required) The name which should be used for this Connection.

* `type` - (Required) Type of database. Possible values include: `MySQL`, `SQLServer`, `SQLAzure`, `Custom`, `NotificationHub`, `ServiceBus`, `EventHub`, `APIHub`, `DocDb`, `RedisCache`, and `PostgreSQL`.

* `value` - (Required) The connection string value.


---

A `cors` block supports the following:

* `allowed_origins` - (Required) Specifies a list of origins that should be allowed to make cross-origin calls.

* `support_credentials` - (Optional) Are credentials allowed in CORS requests? Defaults to `false`.

---

A `docker` block supports the following:

* `registry_url` - (Required) The URL of the docker registry.

* `image_name` - (Required) The name of the Docker image to use.

* `image_tag` - (Required) The image tag of the image to use.

* `registry_username` - (Optional) The username to use for connections to the registry.

~> **NOTE:** This value is required if `container_registry_use_managed_identity` is not set to `true`.

* `registry_password` - (Optional) The password for the account to use to connect to the registry.

~> **NOTE:** This value is required if `container_registry_use_managed_identity` is not set to `true`.

---

A `facebook` block supports the following:

* `app_id` - (Required) The App ID of the Facebook app used for login.

* `app_secret` - (Optional) The App Secret of the Facebook app used for Facebook login. Cannot be specified with `app_secret_setting_name`.

* `app_secret_setting_name` - (Optional) The app setting name that contains the `app_secret` value used for Facebook login. Cannot be specified with `app_secret`.

* `oauth_scopes` - (Optional) Specifies a list of OAuth 2.0 scopes to be requested as part of Facebook login authentication.

---

A `github` block supports the following:

* `client_id` - (Required) The ID of the GitHub app used for login.

* `client_secret` - (Optional) The Client Secret of the GitHub app used for GitHub login. Cannot be specified with `client_secret_setting_name`.

* `client_secret_setting_name` - (Optional) The app setting name that contains the `client_secret` value used for GitHub login. Cannot be specified with `client_secret`.

* `oauth_scopes` - (Optional) Specifies a list of OAuth 2.0 scopes that will be requested as part of GitHub login authentication.

---

A `google` block supports the following:

* `client_id` - (Required) The OpenID Connect Client ID for the Google web application.

* `client_secret` - (Optional) The client secret associated with the Google web application.  Cannot be specified with `client_secret_setting_name`.

* `client_secret_setting_name` - (Optional) The app setting name that contains the `client_secret` value used for Google login. Cannot be specified with `client_secret`.

* `oauth_scopes` - (Optional) Specifies a list of OAuth 2.0 scopes that will be requested as part of Google Sign-In authentication. If not specified, `openid`, `profile`, and `email` are used as default scopes.

---

A `headers` block supports the following:

~> **NOTE:** Please see the [official Azure Documentation](https://docs.microsoft.com/azure/app-service/app-service-ip-restrictions#filter-by-http-header) for details on using header filtering.

* `x_azure_fdid` - (Optional) Specifies a list of Azure Front Door IDs.

* `x_fd_health_probe` - (Optional) Specifies if a Front Door Health Probe should be expected.

* `x_forwarded_for` - (Optional) Specifies a list of addresses for which matching should be applied. Omitting this value means allow any.

* `x_forwarded_host` - (Optional) Specifies a list of Hosts for which matching should be applied.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Linux Function App. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Linux Function App.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

An `ip_restriction` block supports the following:

* `action` - (Optional) The action to take. Possible values are `Allow` or `Deny`.

* `headers` - (Optional) A `headers` block as defined above.

* `ip_address` - (Optional) The CIDR notation of the IP or IP Range to match. For example: `10.0.0.0/24` or `192.168.10.1/32`

* `name` - (Optional) The name which should be used for this `ip_restriction`.

* `priority` - (Optional) The priority value of this `ip_restriction`.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

~> **NOTE:** One and only one of `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified.

---

A `microsoft` block supports the following:

* `client_id` - (Required) The OAuth 2.0 client ID that was created for the app used for authentication.

* `client_secret` - (Optional) The OAuth 2.0 client secret that was created for the app used for authentication. Cannot be specified with `client_secret_setting_name`.

* `client_secret_setting_name` - (Optional) The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication. Cannot be specified with `client_secret`.

* `oauth_scopes` - (Optional) Specifies a list of OAuth 2.0 scopes that will be requested as part of Microsoft Account authentication. If not specified, `wl.basic` is used as the default scope.

---

A `schedule` block supports the following:

* `frequency_interval` - (Required) How often the backup should be executed (e.g. for weekly backup, this should be set to `7` and `frequency_unit` should be set to `Day`).

~> **NOTE:** Not all intervals are supported on all Linux Function App SKUs. Please refer to the official documentation for appropriate values.

* `frequency_unit` - (Required) The unit of time for how often the backup should take place. Possible values include: `Day` and `Hour`. 

* `keep_at_least_one_backup` - (Optional) Should the service keep at least one backup, regardless of age of backup. Defaults to `false`.

* `retention_period_days` - (Optional) After how many days backups should be deleted.

* `start_time` - (Optional) When the schedule should start working in RFC-3339 format.

---

A `scm_ip_restriction` block supports the following:

* `action` - (Optional) The action to take. Possible values are `Allow` or `Deny`.

* `headers` - (Optional) A `headers` block as defined above.

* `ip_address` - (Optional) The CIDR notation of the IP or IP Range to match. For example: `10.0.0.0/24` or `192.168.10.1/32`

* `name` - (Optional) The name which should be used for this `ip_restriction`.

* `priority` - (Optional) The priority value of this `ip_restriction`.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

~> **NOTE:** One and only one of `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified.

---

A `site_config` block supports the following:

* `always_on` - (Optional) If this Linux Web App is Always On enabled. Defaults to `false`.

* `api_definition_url` - (Optional) The URL of the API definition that describes this Linux Function App.

* `api_management_api_id` - (Optional) The ID of the API Management API for this Linux Function App.

* `app_command_line` - (Optional) The App command line to launch.

* `app_scale_limit` - (Optional) The number of workers this function app can scale out to. Only applicable to apps on the Consumption and Premium plan.

* `application_insights_connection_string` - (Optional) The Connection String for linking the Linux Function App to Application Insights.

* `application_insights_key` - (Optional) The Instrumentation Key for connecting the Linux Function App to Application Insights.

* `application_stack` - (Optional) An `application_stack` block as defined above.

~> **Note:** If this is set, there must not be an application setting `FUNCTIONS_WORKER_RUNTIME`.

* `app_service_logs` - (Optional) An `app_service_logs` block as defined above.

* `container_registry_managed_identity_client_id` - (Optional) The Client ID of the Managed Service Identity to use for connections to the Azure Container Registry.

* `container_registry_use_managed_identity` - (Optional) Should connections for Azure Container Registry use Managed Identity.

* `cors` - (Optional) A `cors` block as defined above.

* `default_documents` - (Optional) Specifies a list of Default Documents for the Linux Web App.

* `elastic_instance_minimum` - (Optional) The number of minimum instances for this Linux Function App. Only affects apps on Elastic Premium plans.

* `ftps_state` - (Optional) State of FTP / FTPS service for this function app. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`. Defaults to `Disabled`.

* `health_check_path` - (Optional) The path to be checked for this function app health.

* `health_check_eviction_time_in_min` - (Optional) The amount of time in minutes that a node can be unhealthy before being removed from the load balancer. Possible values are between `2` and `10`. Only valid in conjunction with `health_check_path`.

* `http2_enabled` - (Optional) Specifies if the HTTP2 protocol should be enabled. Defaults to `false`.

* `ip_restriction` - (Optional) One or more `ip_restriction` blocks as defined above.

* `load_balancing_mode` - (Optional) The Site load balancing mode. Possible values include: `WeightedRoundRobin`, `LeastRequests`, `LeastResponseTime`, `WeightedTotalTraffic`, `RequestHash`, `PerSiteRoundRobin`. Defaults to `LeastRequests` if omitted.

* `managed_pipeline_mode` - (Optional) Managed pipeline mode. Possible values include: `Integrated`, `Classic`. Defaults to `Integrated`.

* `minimum_tls_version` - (Optional) The configures the minimum version of TLS required for SSL requests. Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.

* `pre_warmed_instance_count` - (Optional) The number of pre-warmed instances for this function app. Only affects apps on an Elastic Premium plan.

* `remote_debugging_enabled` - (Optional) Should Remote Debugging be enabled. Defaults to `false`.

* `remote_debugging_version` - (Optional) The Remote Debugging Version. Possible values include `VS2017` and `VS2019`.

* `runtime_scale_monitoring_enabled` - (Optional) Should Scale Monitoring of the Functions Runtime be enabled? 

* `scm_ip_restriction` - (Optional) One or more `scm_ip_restriction` blocks as defined above.

* `scm_minimum_tls_version` - (Optional) Configures the minimum version of TLS required for SSL requests to the SCM site Possible values include: `1.0`, `1.1`, and  `1.2`. Defaults to `1.2`.

* `scm_use_main_ip_restriction` - (Optional) Should the Linux Function App `ip_restriction` configuration be used for the SCM also.

* `use_32_bit_worker` - (Optional) Should the Linux Web App use a 32-bit worker process. Defaults to `true`.

* `vnet_route_all_enabled` - (Optional) Should all outbound traffic to have NAT Gateways, Network Security Groups and User Defined Routes applied? Defaults to `false`.

* `websockets_enabled` - (Optional) Should Web Sockets be enabled. Defaults to `false`.

* `worker_count` - (Optional) The number of Workers for this Linux Function App.

---

A `sticky_settings` block supports the following:

* `app_setting_names` - (Optional) A list of `app_setting` names that the Linux Function App will not swap between Slots when a swap operation is triggered.

* `connection_string_names` - (Optional) A list of `connection_string` names that the Linux Function App will not swap between Slots when a swap operation is triggered.

---

A `twitter` block supports the following:

* `consumer_key` - (Required) The OAuth 1.0a consumer key of the Twitter application used for sign-in.

* `consumer_secret` - (Optional) The OAuth 1.0a consumer secret of the Twitter application used for sign-in. Cannot be specified with `consumer_secret_setting_name`.

* `consumer_secret_setting_name` - (Optional) The app setting name that contains the OAuth 1.0a consumer secret of the Twitter application used for sign-in. Cannot be specified with `consumer_secret`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Linux Function App.

* `custom_domain_verification_id` - The identifier used by App Service to perform domain ownership verification via DNS TXT record.

* `default_hostname` - The default hostname of the Linux Function App.

* `identity` - An `identity` block as defined below.

* `kind` - The Kind value for this Linux Function App.

* `outbound_ip_address_list` - A list of outbound IP addresses. For example `["52.23.25.3", "52.143.43.12"]`

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12`.

* `possible_outbound_ip_address_list` - A list of possible outbound IP addresses, not all of which are necessarily in use. This is a superset of `outbound_ip_address_list`. For example `["52.23.25.3", "52.143.43.12"]`.

* `possible_outbound_ip_addresses` - A comma separated list of possible outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12,52.143.43.17`. This is a superset of `outbound_ip_addresses`. For example `["52.23.25.3", "52.143.43.12","52.143.43.17"]`.

* `site_credential` - A `site_credential` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

---

A `site_credential` block exports the following:

* `name` - The Site Credentials Username used for publishing.

* `password` - The Site Credentials Password used for publishing.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Linux Function App.
* `read` - (Defaults to 5 minutes) Used when retrieving the Linux Function App.
* `update` - (Defaults to 30 minutes) Used when updating the Linux Function App.
* `delete` - (Defaults to 30 minutes) Used when deleting the Linux Function App.

## Import

Linux Function Apps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_linux_function_app.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
```
