---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_windows_function_app_slot"
description: |-
  Manages a Windows Function App Slot.
---
# azurerm_windows_function_app_slot

Manages a Windows Function App Slot.

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
  name                     = "windowsfunctionappsa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "example" {
  name                = "example-app-service-plan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  os_type             = "Windows"
  sku_name            = "Y1"
}

resource "azurerm_windows_function_app" "example" {
  name                 = "example-windows-function-app"
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  storage_account_name = azurerm_storage_account.example.name
  service_plan_id      = azurerm_service_plan.example.id

  site_config {}
}

resource "azurerm_windows_function_app_slot" "example" {
  name                 = "example-slot"
  function_app_id      = azurerm_windows_function_app.example.id
  storage_account_name = azurerm_storage_account.example.name

  site_config {}
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Windows Function App Slot. Changing this forces a new resource to be created.

* `function_app_id` - (Required) The name of the Windows Function App this Slot is a member of. Changing this forces a new resource to be created.

* `site_config` - (Required) a `site_config` block as detailed below.

---

* `app_settings` - (Optional) A map of key-value pairs for [App Settings](https://docs.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.

* `auth_settings` - (Optional) an `auth_settings` block as detailed below.

* `auth_settings_v2` - (Optional) an `auth_settings_v2` block as detailed below.

* `backup` - (Optional) a `backup` block as detailed below.

* `builtin_logging_enabled` - (Optional) Should built-in logging be enabled. Configures `AzureWebJobsDashboard` app setting based on the configured storage setting. Defaults to `true`.

* `client_certificate_enabled` - (Optional) Should the Function App Slot use Client Certificates.

* `client_certificate_mode` - (Optional) The mode of the Function App Slot's client certificates requirement for incoming requests. Possible values are `Required`, `Optional`, and `OptionalInteractiveUser`. Defaults to `Optional`.

* `client_certificate_exclusion_paths` - (Optional) Paths to exclude when using client certificates, separated by ;

* `connection_string` - (Optional) a `connection_string` block as detailed below.

* `content_share_force_disabled` - (Optional) Force disable the content share settings.

* `daily_memory_time_quota` - (Optional) The amount of memory in gigabyte-seconds that your application is allowed to consume per day. Setting this value only affects function apps in Consumption Plans. Defaults to `0`.

* `enabled` - (Optional) Is the Windows Function App Slot enabled. Defaults to `true`.

* `ftp_publish_basic_authentication_enabled` - (Optional) Should the default FTP Basic Authentication publishing profile be enabled. Defaults to `true`.

* `functions_extension_version` - (Optional) The runtime version associated with the Function App Slot. Defaults to `~4`.

* `https_only` - (Optional) Can the Function App Slot only be accessed via HTTPS?. Defaults to `false`.

* `public_network_access_enabled` - (Optional) Should public network access be enabled for the Function App. Defaults to `true`.

* `identity` - (Optional) an `identity` block as detailed below.

* `key_vault_reference_identity_id` - (Optional) The User Assigned Identity ID used for accessing KeyVault secrets. The identity must be assigned to the application in the `identity` block. [For more information see - Access vaults with a user-assigned identity](https://docs.microsoft.com/azure/app-service/app-service-key-vault-references#access-vaults-with-a-user-assigned-identity)

* `service_plan_id` - (Optional) The ID of the Service Plan in which to run this slot. If not specified the same Service Plan as the Windows Function App will be used.

* `storage_account_access_key` - (Optional) The access key which will be used to access the storage account for the Function App Slot.

* `storage_account_name` - (Optional) The backend storage account name which will be used by this Function App Slot.

* `storage_account` - (Optional) One or more `storage_account` blocks as defined below.

* `storage_uses_managed_identity` - (Optional) Should the Function App Slot use its Managed Identity to access storage.

~> **Note:** One of `storage_account_access_key` or `storage_uses_managed_identity` must be specified when using `storage_account_name`.

* `storage_key_vault_secret_id` - (Optional) The Key Vault Secret ID, optionally including version, that contains the Connection String to connect to the storage account for this Function App Slot.

~> **Note:** `storage_key_vault_secret_id` cannot be used with `storage_account_name`.

~> **Note:** `storage_key_vault_secret_id` used without a version will use the latest version of the secret, however, the service can take up to 24h to pick up a rotation of the latest version. See the [official docs](https://docs.microsoft.com/azure/app-service/app-service-key-vault-references#rotation) for more information.

* `tags` - (Optional) A mapping of tags which should be assigned to the Windows Function App Slot.

* `virtual_network_backup_restore_enabled` - (Optional) Whether backup and restore operations over the linked virtual network are enabled. Defaults to `false`.

* `virtual_network_subnet_id` - (Optional) The subnet id which will be used by this Function App Slot for [regional virtual network integration](https://docs.microsoft.com/en-us/azure/app-service/overview-vnet-integration#regional-virtual-network-integration).

~> **Note:** The AzureRM Terraform provider provides regional virtual network integration via the standalone resource [app_service_virtual_network_swift_connection](app_service_virtual_network_swift_connection.html) and in-line within this resource using the `virtual_network_subnet_id` property. You cannot use both methods simultaneously. If the virtual network is set via the resource `app_service_virtual_network_swift_connection` then `ignore_changes` should be used in the function app slot configuration.

~> **Note:** Assigning the `virtual_network_subnet_id` property requires [RBAC permissions on the subnet](https://docs.microsoft.com/en-us/azure/app-service/overview-vnet-integration#permissions)

* `vnet_image_pull_enabled` - (Optional) Specifies whether traffic for the image pull should be routed over virtual network. Defaults to `false`.

~> **Note:** The feature can also be enabled via the app setting `WEBSITE_PULL_IMAGE_OVER_VNET`. The Setting is enabled by default for app running in the App Service Environment.

* `webdeploy_publish_basic_authentication_enabled` - (Optional) Should the default WebDeploy Basic Authentication publishing credentials enabled. Defaults to `true`.

---

An `auth_settings` block supports the following:

* `enabled` - (Required) Should the Authentication / Authorization feature be enabled?

* `active_directory` - (Optional) an `active_directory` block as detailed below.

* `additional_login_parameters` - (Optional) Specifies a map of login Parameters to send to the OpenID Connect authorization endpoint when a user logs in.

* `allowed_external_redirect_urls` - (Optional) Specifies a list of External URLs that can be redirected to as part of logging in or logging out of the Windows Web App.

* `default_provider` - (Optional) The default authentication provider to use when multiple providers are configured. Possible values include: `AzureActiveDirectory`, `Facebook`, `Google`, `MicrosoftAccount`, `Twitter`, `Github`.

~> **Note:** This setting is only needed if multiple providers are configured, and the `unauthenticated_client_action` is set to "RedirectToLoginPage".

* `facebook` - (Optional) a `facebook` block as detailed below.

* `github` - (Optional) a `github` block as detailed below.

* `google` - (Optional) a `google` block as detailed below.

* `issuer` - (Optional) The OpenID Connect Issuer URI that represents the entity which issues access tokens.

~> **Note:** When using Azure Active Directory, this value is the URI of the directory tenant, e.g. <https://sts.windows.net/{tenant-guid}/>.

* `microsoft` - (Optional) a `microsoft` block as detailed below.

* `runtime_version` - (Optional) The RuntimeVersion of the Authentication / Authorization feature in use.

* `token_refresh_extension_hours` - (Optional) The number of hours after session token expiration that a session token can be used to call the token refresh API. Defaults to `72` hours.

* `token_store_enabled` - (Optional) Should the Windows Web App durably store platform-specific security tokens that are obtained during login flows? Defaults to `false`.

* `twitter` - (Optional) a `twitter` block as detailed below.

* `unauthenticated_client_action` - (Optional) The action to take when an unauthenticated client attempts to access the app. Possible values include: `RedirectToLoginPage`, `AllowAnonymous`.

---

An `auth_settings_v2` block supports the following:

* `auth_enabled` - (Optional) Should the AuthV2 Settings be enabled. Defaults to `false`.

* `runtime_version` - (Optional) The Runtime Version of the Authentication and Authorisation feature of this App. Defaults to `~1`.

* `config_file_path` - (Optional) The path to the App Auth settings.

~> **Note:** Relative Paths are evaluated from the Site Root directory.

* `require_authentication` - (Optional) Should the authentication flow be used for all requests.

* `unauthenticated_action` - (Optional) The action to take for requests made without authentication. Possible values include `RedirectToLoginPage`, `AllowAnonymous`, `Return401`, and `Return403`. Defaults to `RedirectToLoginPage`.

* `default_provider` - (Optional) The Default Authentication Provider to use when the `unauthenticated_action` is set to `RedirectToLoginPage`. Possible values include: `apple`, `azureactivedirectory`, `facebook`, `github`, `google`, `twitter` and the `name` of your `custom_oidc_v2` provider.

~> **Note:** Whilst any value will be accepted by the API for `default_provider`, it can leave the app in an unusable state if this value does not correspond to the name of a known provider (either built-in value, or custom_oidc name) as it is used to build the auth endpoint URI.

* `excluded_paths` - (Optional) The paths which should be excluded from the `unauthenticated_action` when it is set to `RedirectToLoginPage`.

~> **Note:** This list should be used instead of setting `WEBSITE_WARMUP_PATH` in `app_settings` as it takes priority.

* `require_https` - (Optional) Should HTTPS be required on connections? Defaults to `true`.

* `http_route_api_prefix` - (Optional) The prefix that should precede all the authentication and authorisation paths. Defaults to `/.auth`.

* `forward_proxy_convention` - (Optional) The convention used to determine the url of the request made. Possible values include `NoProxy`, `Standard`, `Custom`. Defaults to `NoProxy`.

* `forward_proxy_custom_host_header_name` - (Optional) The name of the custom header containing the host of the request.

* `forward_proxy_custom_scheme_header_name` - (Optional) The name of the custom header containing the scheme of the request.

* `apple_v2` - (Optional) An `apple_v2` block as defined below.

* `active_directory_v2` - (Optional) An `active_directory_v2` block as defined below.

* `azure_static_web_app_v2` - (Optional) An `azure_static_web_app_v2` block as defined below.

* `custom_oidc_v2` - (Optional) Zero or more `custom_oidc_v2` blocks as defined below.

* `facebook_v2` - (Optional) A `facebook_v2` block as defined below.

* `github_v2` - (Optional) A `github_v2` block as defined below.

* `google_v2` - (Optional) A `google_v2` block as defined below.

* `microsoft_v2` - (Optional) A `microsoft_v2` block as defined below.

* `twitter_v2` - (Optional) A `twitter_v2` block as defined below.

* `login` - (Required) A `login` block as defined below.

---

An `apple_v2` block supports the following:

* `client_id` - (Required) The OpenID Connect Client ID for the Apple web application.

* `client_secret_setting_name` - (Required) The app setting name that contains the `client_secret` value used for Apple Login.

!> **Note:** A setting with this name must exist in `app_settings` to function correctly.

* `login_scopes` - A list of Login Scopes provided by this Authentication Provider.

~> **Note:** This is configured on the Authentication Provider side and is Read Only here.

---

An `active_directory_v2` block supports the following:

* `client_id` - (Required) The ID of the Client to use to authenticate with Azure Active Directory.

* `tenant_auth_endpoint` - (Required) The Azure Tenant Endpoint for the Authenticating Tenant. e.g. `https://login.microsoftonline.com/{tenant-guid}/v2.0/`

~> **Note:** [Here](https://learn.microsoft.com/en-us/entra/identity-platform/authentication-national-cloud#microsoft-entra-authentication-endpoints) is a list of possible authentication endpoints based on the cloud environment. [Here](https://learn.microsoft.com/en-us/azure/app-service/configure-authentication-provider-aad?tabs=workforce-tenant) is more information to better understand how to configure authentication for Azure App Service or Azure Functions.

* `client_secret_setting_name` - (Optional) The App Setting name that contains the client secret of the Client.

!> **Note:** A setting with this name must exist in `app_settings` to function correctly.

* `client_secret_certificate_thumbprint` - (Optional) The thumbprint of the certificate used for signing purposes.

!> **Note:** If one `client_secret_setting_name` or `client_secret_certificate_thumbprint` is specified, terraform won't write the client secret or secret certificate thumbprint back to `app_setting`, so make sure they are existed in `app_settings` to function correctly.

* `jwt_allowed_groups` - (Optional) A list of Allowed Groups in the JWT Claim.

* `jwt_allowed_client_applications` - (Optional) A list of Allowed Client Applications in the JWT Claim.

* `www_authentication_disabled` - (Optional) Should the www-authenticate provider should be omitted from the request? Defaults to `false`.

* `allowed_groups` - (Optional) The list of allowed Group Names for the Default Authorisation Policy.

* `allowed_identities` - (Optional) The list of allowed Identities for the Default Authorisation Policy.

* `allowed_applications` - (Optional) The list of allowed Applications for the Default Authorisation Policy.

* `login_parameters` - (Optional) A map of key-value pairs to send to the Authorisation Endpoint when a user logs in.

* `allowed_audiences` - (Optional) Specifies a list of Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

~> **Note:** This is configured on the Authentication Provider side and is Read Only here.

---

An `azure_static_web_app_v2` block supports the following:

* `client_id` - (Required) The ID of the Client to use to authenticate with Azure Static Web App Authentication.

---

A `custom_oidc_v2` block supports the following:

* `name` - (Required) The name of the Custom OIDC Authentication Provider.

~> **Note:** An `app_setting` matching this value in upper case with the suffix of `_PROVIDER_AUTHENTICATION_SECRET` is required. e.g. `MYOIDC_PROVIDER_AUTHENTICATION_SECRET` for a value of `myoidc`.

* `client_id` - (Required) The ID of the Client to use to authenticate with the Custom OIDC.

* `openid_configuration_endpoint` - (Required) The app setting name that contains the `client_secret` value used for the Custom OIDC Login.

* `name_claim_type` - (Optional) The name of the claim that contains the users name.

* `scopes` - (Optional) The list of the scopes that should be requested while authenticating.

* `client_credential_method` - The Client Credential Method used.

* `client_secret_setting_name` - The App Setting name that contains the secret for this Custom OIDC Client. This is generated from `name` above and suffixed with `_PROVIDER_AUTHENTICATION_SECRET`.

* `authorisation_endpoint` - The endpoint to make the Authorisation Request as supplied by `openid_configuration_endpoint` response.

* `token_endpoint` - The endpoint used to request a Token as supplied by `openid_configuration_endpoint` response.

* `issuer_endpoint` - The endpoint that issued the Token as supplied by `openid_configuration_endpoint` response.

* `certification_uri` - The endpoint that provides the keys necessary to validate the token as supplied by `openid_configuration_endpoint` response.

---

A `facebook_v2` block supports the following:

* `app_id` - (Required) The App ID of the Facebook app used for login.

* `app_secret_setting_name` - (Required) The app setting name that contains the `app_secret` value used for Facebook Login.

!> **Note:** A setting with this name must exist in `app_settings` to function correctly.

* `graph_api_version` - (Optional) The version of the Facebook API to be used while logging in.

* `login_scopes` - (Optional) The list of scopes that should be requested as part of Facebook Login authentication.

---

A `github_v2` block supports the following:

* `client_id` - (Required) The ID of the GitHub app used for login.

* `client_secret_setting_name` - (Required) The app setting name that contains the `client_secret` value used for GitHub Login.

!> **Note:** A setting with this name must exist in `app_settings` to function correctly.

* `login_scopes` - (Optional) The list of OAuth 2.0 scopes that should be requested as part of GitHub Login authentication.

---

A `google_v2` block supports the following:

* `client_id` - (Required) The OpenID Connect Client ID for the Google web application.

* `client_secret_setting_name` - (Required) The app setting name that contains the `client_secret` value used for Google Login.

!> **Note:** A setting with this name must exist in `app_settings` to function correctly.

* `allowed_audiences` - (Optional) Specifies a list of Allowed Audiences that should be requested as part of Google Sign-In authentication.

* `login_scopes` - (Optional) The list of OAuth 2.0 scopes that should be requested as part of Google Sign-In authentication.

---

A `microsoft_v2` block supports the following:

* `client_id` - (Required) The OAuth 2.0 client ID that was created for the app used for authentication.

* `client_secret_setting_name` - (Required) The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication.

!> **Note:** A setting with this name must exist in `app_settings` to function correctly.

* `allowed_audiences` - (Optional) Specifies a list of Allowed Audiences that will be requested as part of Microsoft Sign-In authentication.

* `login_scopes` - (Optional) The list of Login scopes that should be requested as part of Microsoft Account authentication.

---

A `twitter_v2` block supports the following:

* `consumer_key` - (Required) The OAuth 1.0a consumer key of the Twitter application used for sign-in.

* `consumer_secret_setting_name` - (Required) The app setting name that contains the OAuth 1.0a consumer secret of the Twitter application used for sign-in.

!> **Note:** A setting with this name must exist in `app_settings` to function correctly.

---

A `login` block supports the following:

* `logout_endpoint` - (Optional) The endpoint to which logout requests should be made.

* `token_store_enabled` - (Optional) Should the Token Store configuration Enabled. Defaults to `false`

* `token_refresh_extension_time` - (Optional) The number of hours after session token expiration that a session token can be used to call the token refresh API. Defaults to `72` hours.

* `token_store_path` - (Optional) The directory path in the App Filesystem in which the tokens will be stored.

* `token_store_sas_setting_name` - (Optional) The name of the app setting which contains the SAS URL of the blob storage containing the tokens.

* `preserve_url_fragments_for_logins` - (Optional) Should the fragments from the request be preserved after the login request is made. Defaults to `false`.

* `allowed_external_redirect_urls` - (Optional) External URLs that can be redirected to as part of logging in or logging out of the app. This is an advanced setting typically only needed by Windows Store application backends.

~> **Note:** URLs within the current domain are always implicitly allowed.

* `cookie_expiration_convention` - (Optional) The method by which cookies expire. Possible values include: `FixedTime`, and `IdentityProviderDerived`. Defaults to `FixedTime`.

* `cookie_expiration_time` - (Optional) The time after the request is made when the session cookie should expire. Defaults to `08:00:00`.

* `validate_nonce` - (Optional) Should the nonce be validated while completing the login flow. Defaults to `true`.

* `nonce_expiration_time` - (Optional) The time after the request is made when the nonce should expire. Defaults to `00:05:00`.

---

A `backup` block supports the following:

* `name` - (Required) The name which should be used for this Backup.

* `schedule` - (Required) a `schedule` block as detailed below.

* `storage_account_url` - (Required) The SAS URL to the container.

* `enabled` - (Optional) Should this backup job be enabled? Defaults to `true`.

---

A `connection_string` block supports the following:

* `name` - (Required) The name which should be used for this Connection.

* `type` - (Required) Type of database. Possible values include: `APIHub`, `Custom`, `DocDb`, `EventHub`, `MySQL`, `NotificationHub`, `PostgreSQL`, `RedisCache`, `ServiceBus`, `SQLAzure`, and `SQLServer`.

* `value` - (Required) The connection string value.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Windows Function App Slot. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Windows Function App Slot.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `site_config` block supports the following:

* `always_on` - (Optional) If this Windows Web App is Always On enabled. Defaults to `false`.

* `api_definition_url` - (Optional) The URL of the API definition that describes this Windows Function App.

* `api_management_api_id` - (Optional) The ID of the API Management API for this Windows Function App.

* `app_command_line` - (Optional) The program and any arguments used to launch this app via the command line. (Example `node myapp.js`).

* `app_scale_limit` - (Optional) The number of workers this function app can scale out to. Only applicable to apps on the Consumption and Premium plan.

* `app_service_logs` - (Optional) an `app_service_logs` block as detailed below.

* `application_insights_connection_string` - (Optional) The Connection String for linking the Windows Function App to Application Insights.

* `application_insights_key` - (Optional) The Instrumentation Key for connecting the Windows Function App to Application Insights.

* `application_stack` - (Optional) an `application_stack` block as detailed below.

* `auto_swap_slot_name` - (Optional) The name of the slot to automatically swap with when this slot is successfully deployed.

* `cors` - (Optional) a `cors` block as detailed below.

* `default_documents` - (Optional) Specifies a list of Default Documents for the Windows Web App.

* `detailed_error_logging_enabled` - Is detailed error logging enabled

* `elastic_instance_minimum` - (Optional) The number of minimum instances for this Windows Function App. Only affects apps on Elastic Premium plans.

* `ftps_state` - (Optional) State of FTP / FTPS service for this function app. Possible values include: `AllAllowed`, `FtpsOnly` and `Disabled`. Defaults to `Disabled`.

* `health_check_eviction_time_in_min` - (Optional) The amount of time in minutes that a node is unhealthy before being removed from the load balancer. Possible values are between `2` and `10`. Defaults to `0`. Only valid in conjunction with `health_check_path`.

* `health_check_path` - (Optional) The path to be checked for this function app health.

* `http2_enabled` - (Optional) Specifies if the HTTP2 protocol should be enabled. Defaults to `false`.

* `ip_restriction` - (Optional) an `ip_restriction` block as detailed below.

* `ip_restriction_default_action` - (Optional) The Default action for traffic that does not match any `ip_restriction` rule. possible values include `Allow` and `Deny`. Defaults to `Allow`.

* `load_balancing_mode` - (Optional) The Site load balancing mode. Possible values include: `WeightedRoundRobin`, `LeastRequests`, `LeastResponseTime`, `WeightedTotalTraffic`, `RequestHash`, `PerSiteRoundRobin`. Defaults to `LeastRequests` if omitted.

* `managed_pipeline_mode` - (Optional) The Managed Pipeline mode. Possible values include: `Integrated`, `Classic`. Defaults to `Integrated`.

* `minimum_tls_version` - (Optional) The configures the minimum version of TLS required for SSL requests. Possible values include: `1.0`, `1.1`, `1.2` and `1.3`. Defaults to `1.2`.

* `pre_warmed_instance_count` - (Optional) The number of pre-warmed instances for this function app. Only affects apps on an Elastic Premium plan.

* `remote_debugging_enabled` - (Optional) Should Remote Debugging be enabled. Defaults to `false`.

* `remote_debugging_version` - (Optional) The Remote Debugging Version. Currently only `VS2022` is supported.

* `runtime_scale_monitoring_enabled` - (Optional) Should Scale Monitoring of the Functions Runtime be enabled?

~> **Note:** Functions runtime scale monitoring can only be enabled for Elastic Premium Function Apps or Workflow Standard Logic Apps and requires a minimum prewarmed instance count of 1.

* `scm_ip_restriction` - (Optional) a `scm_ip_restriction` block as detailed below.

* `scm_ip_restriction_default_action` - (Optional) The Default action for traffic that does not match any `scm_ip_restriction` rule. possible values include `Allow` and `Deny`. Defaults to `Allow`.

* `scm_minimum_tls_version` - (Optional) Configures the minimum version of TLS required for SSL requests to the SCM site Possible values include: `1.0`, `1.1`, `1.2` and `1.3`. Defaults to `1.2`.

* `scm_type` - The SCM Type in use by the Windows Function App.

* `scm_use_main_ip_restriction` - (Optional) Should the Windows Function App `ip_restriction` configuration be used for the SCM also.

* `use_32_bit_worker` - (Optional) Should the Windows Web App use a 32-bit worker. Defaults to `true`.

* `vnet_route_all_enabled` - (Optional) Should all outbound traffic to have NAT Gateways, Network Security Groups and User Defined Routes applied? Defaults to `false`.

* `websockets_enabled` - (Optional) Should Web Sockets be enabled. Defaults to `false`.

* `windows_fx_version` - The Windows FX Version string.

* `worker_count` - (Optional) The number of Workers for this Windows Function App.

---

A `site_credential` block supports the following:

* `name` - The Site Credentials Username used for publishing.

* `password` - The Site Credentials Password used for publishing.

---

An `active_directory` block supports the following:

* `client_id` - (Required) The ID of the Client to use to authenticate with Azure Active Directory.

* `allowed_audiences` - (Optional) Specifies a list of Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

~> **Note:** The `client_id` value is always considered an allowed audience.

* `client_secret` - (Optional) The Client Secret for the Client ID. Cannot be used with `client_secret_setting_name`.

* `client_secret_setting_name` - (Optional) The App Setting name that contains the client secret of the Client. Cannot be used with `client_secret`.

---

A `facebook` block supports the following:

* `app_id` - (Required) The App ID of the Facebook app used for login.

* `app_secret` - (Optional) The App Secret of the Facebook app used for Facebook login. Cannot be specified with `app_secret_setting_name`.

* `app_secret_setting_name` - (Optional) The app setting name that contains the `app_secret` value used for Facebook login. Cannot be specified with `app_secret`.

* `oauth_scopes` - (Optional) Specifies a list of OAuth 2.0 scopes to be requested as part of Facebook Login authentication.

---

A `github` block supports the following:

* `client_id` - (Required) The ID of the GitHub app used for login.

* `client_secret` - (Optional) The Client Secret of the GitHub app used for GitHub login. Cannot be specified with `client_secret_setting_name`.

* `client_secret_setting_name` - (Optional) The app setting name that contains the `client_secret` value used for GitHub login. Cannot be specified with `client_secret`.

* `oauth_scopes` - (Optional) an `oauth_scopes`.

---

A `google` block supports the following:

* `client_id` - (Required) The OpenID Connect Client ID for the Google web application.

* `client_secret` - (Optional) The client secret associated with the Google web application. Cannot be specified with `client_secret_setting_name`.

* `client_secret_setting_name` - (Optional) The app setting name that contains the `client_secret` value used for Google login. Cannot be specified with `client_secret`.

* `oauth_scopes` - (Optional) Specifies a list of OAuth 2.0 scopes that will be requested as part of Google Sign-In authentication. If not specified, "openid", "profile", and "email" are used as default scopes.

---

A `microsoft` block supports the following:

* `client_id` - (Required) The OAuth 2.0 client ID that was created for the app used for authentication.

* `client_secret` - (Optional) The OAuth 2.0 client secret that was created for the app used for authentication. Cannot be specified with `client_secret_setting_name`.

* `client_secret_setting_name` - (Optional) The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication. Cannot be specified with `client_secret`.

* `oauth_scopes` - (Optional) Specifies a list of OAuth 2.0 scopes that will be requested as part of Microsoft Account authentication. If not specified, `wl.basic` is used as the default scope.

---

A `twitter` block supports the following:

* `consumer_key` - (Required) The OAuth 1.0a consumer key of the Twitter application used for sign-in.

* `consumer_secret` - (Optional) The OAuth 1.0a consumer secret of the Twitter application used for sign-in. Cannot be specified with `consumer_secret_setting_name`.

* `consumer_secret_setting_name` - (Optional) The app setting name that contains the OAuth 1.0a consumer secret of the Twitter application used for sign-in. Cannot be specified with `consumer_secret`.

---

A `schedule` block supports the following:

* `frequency_interval` - (Required) How often the backup should be executed (e.g. for weekly backup, this should be set to `7` and `frequency_unit` should be set to `Day`).

~> **Note:** Not all intervals are supported on all SKUs. Please refer to the official documentation for appropriate values.

* `frequency_unit` - (Required) The unit of time for how often the backup should take place. Possible values include: `Day` and `Hour`.

* `keep_at_least_one_backup` - (Optional) Should the service keep at least one backup, regardless of age of backup. Defaults to `false`.

* `retention_period_days` - (Optional) After how many days backups should be deleted. Defaults to `30`.

* `start_time` - (Optional) When the schedule should start working in RFC-3339 format.

* `last_execution_time` - The time the backup was last attempted.

---

An `app_service_logs` block supports the following:

* `disk_quota_mb` - (Optional) The amount of disk space to use for logs. Valid values are between `25` and `100`. Defaults to `35`.

* `retention_period_days` - (Optional) The retention period for logs in days. Valid values are between `0` and `99999`.(never delete).

~> **Note:** This block is not supported on Consumption plans.

---

An `application_stack` block supports the following:

* `dotnet_version` - (Optional) The version of .Net. Possible values are `v3.0`, `v4.0`, `v6.0`, `v7.0`, `v8.0` and `v9.0`. Defaults to `v4.0`.

* `use_dotnet_isolated_runtime` - (Optional) Should the DotNet process use an isolated runtime. Defaults to `false`.

* `java_version` - (Optional) The version of Java to use. Possible values are `1.8`, `11` and `17` (In-Preview).

* `node_version` - (Optional) The version of Node to use. Possible values are `~12`, `~14`, `~16`, `~18`, `~20`, and `~22`.

* `powershell_core_version` - (Optional) The PowerShell Core version to use. Possible values are `7`, `7.2`, and `7.4`.

* `use_custom_runtime` - (Optional) Does the Function App use a custom Application Stack?

---

A `cors` block supports the following:

* `allowed_origins` - (Optional) Specifies a list of origins that should be allowed to make cross-origin calls.

* `support_credentials` - (Optional) Are credentials allowed in CORS requests? Defaults to `false`.

---

An `ip_restriction` block supports the following:

* `action` - (Optional) The action to take. Possible values are `Allow` or `Deny`. Defaults to `Allow`.

* `headers` - (Optional) a `headers` block as detailed below.

* `ip_address` - (Optional) The CIDR notation of the IP or IP Range to match. For example: `10.0.0.0/24` or `192.168.10.1/32`

* `name` - (Optional) The name which should be used for this `ip_restriction`.

* `priority` - (Optional) The priority value of this `ip_restriction`. Defaults to `65000`.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

~> **Note:** One and only one of `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified.

* `description` - (Optional) The Description of this IP Restriction.

---

A `scm_ip_restriction` block supports the following:

* `action` - (Optional) The action to take. Possible values are `Allow` or `Deny`. Defaults to `Allow`.

* `headers` - (Optional) a `headers` block as detailed below.

* `ip_address` - (Optional) The CIDR notation of the IP or IP Range to match. For example: `10.0.0.0/24` or `192.168.10.1/32`

* `name` - (Optional) The name which should be used for this `ip_restriction`.

* `priority` - (Optional) The priority value of this `ip_restriction`. Defaults to `65000`.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

~> **Note:** Exactly one of `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified.

* `description` - (Optional) The Description of this IP Restriction.

---

A `headers` block supports the following:

~> **Note:** Please see the [official Azure Documentation](https://docs.microsoft.com/azure/app-service/app-service-ip-restrictions#filter-by-http-header) for details on using header filtering.

* `x_azure_fdid` - (Optional) Specifies a list of Azure Front Door IDs.

* `x_fd_health_probe` - (Optional) Specifies if a Front Door Health Probe should be expected. The only possible value is `1`.

* `x_forwarded_for` - (Optional) Specifies a list of addresses for which matching should be applied. Omitting this value means allow any.

* `x_forwarded_host` - (Optional) Specifies a list of Hosts for which matching should be applied.

---

A `storage_account` block supports the following:

* `access_key` - (Required) The Access key for the storage account.

* `account_name` - (Required) The Name of the Storage Account.

* `name` - (Required) The name which should be used for this Storage Account.

* `share_name` - (Required) The Name of the File Share or Container Name for Blob storage.

* `type` - (Required) The Azure Storage Type. Possible values include `AzureFiles`.

* `mount_path` - (Optional) The path at which to mount the storage share.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Windows Function App Slot

* `custom_domain_verification_id` - The identifier used by App Service to perform domain ownership verification via DNS TXT record.

* `default_hostname` - The default hostname of the Windows Function App Slot.

* `hosting_environment_id` - The ID of the App Service Environment used by Function App Slot.

* `identity` - An `identity` block as defined below.

* `kind` - The Kind value for this Windows Function App Slot.

* `outbound_ip_address_list` - A list of outbound IP addresses. For example `["52.23.25.3", "52.143.43.12"]`.

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

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Windows Function App Slot.
* `read` - (Defaults to 5 minutes) Used when retrieving the Windows Function App Slot.
* `update` - (Defaults to 30 minutes) Used when updating the Windows Function App Slot.
* `delete` - (Defaults to 30 minutes) Used when deleting the Windows Function App Slot.

## Import

A Windows Function App Slot can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_windows_function_app_slot.example "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1"
```
