---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_function_app_flex_consumption"
description: |-
  Manages a Function App Running on a Flex Consumption Plan.
---

# azurerm_function_app_flex_consumption

Manages a Function App Running on a Flex Consumption Plan.

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
  name                     = "examplelinuxfunctionappsa"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "example-flexcontainer"
  storage_account_id    = azurerm_storage_account.example.id
  container_access_type = "private"
}

resource "azurerm_service_plan" "example" {
  name                = "example-app-service-plan"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku_name            = "FC1"
  os_type             = "Linux"
}

resource "azurerm_function_app_flex_consumption" "example" {
  name                = "example-linux-function-app"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  service_plan_id     = azurerm_service_plan.example.id

  storage_container_type      = "blobContainer"
  storage_container_endpoint  = "${azurerm_storage_account.example.primary_blob_endpoint}${azurerm_storage_container.example.name}"
  storage_authentication_type = "StorageAccountConnectionString"
  storage_access_key          = azurerm_storage_account.example.primary_access_key
  runtime_name                = "node"
  runtime_version             = "20"
  maximum_instance_count      = 50
  instance_memory_in_mb       = 2048

  site_config {}
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region where the Function App should exist. Changing this forces a new Function App to be created.

* `name` - (Required) The name which should be used for this Function App. Changing this forces a new Function App to be created. Limit the function name to 32 characters to avoid naming collisions. For more information about [Function App naming rule](https://docs.microsoft.com/en-us/azure/azure-resource-manager/management/resource-name-rules#microsoftweb) and [Host ID Collisions](https://github.com/Azure/azure-functions-host/wiki/Host-IDs#host-id-collisions)

* `resource_group_name` - (Required) The name of the Resource Group where the Function App should exist. Changing this forces a new Linux Function App to be created.

* `service_plan_id` - (Required) The ID of the App Service Plan within which to create this Function App. Changing this forces a new Linux Function App to be created.

* `site_config` - (Required) A `site_config` block as defined below.

* `storage_container_type` - (Required) The storage container type used for the Function App. The current supported type is `blobContainer`.

* `storage_container_endpoint` - (Required) The backend storage container endpoint which will be used by this Function App.

* `storage_authentication_type` - (Required) The authentication type which will be used to access the backend storage account for the Function App. Possible values are `StorageAccountConnectionString`, `SystemAssignedIdentity`, and `UserAssignedIdentity`.

* `runtime_name` - (Required) The Runtime of the Linux Function App. Possible values are `node`, `dotnet-isolated`, `powershell`, `python`, `java` and `custom`.

* `runtime_version` - (Required) The Runtime version of the Linux Function App. The values are diff from different runtime version. The supported values are `8.0`, `9.0` for `dotnet-isolated`, `20` for `node`, `3.10`, `3.11` for `python`, `11`, `17` for `java`, `7.4` for `powershell`.

---

* `app_settings` - (Optional) A map of key-value pairs for [App Settings](https://docs.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.

~> **Note:** For storage related settings, please use related properties that are available such as `storage_access_key`, terraform will assign the value to keys such as `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING`, `AzureWebJobsStorage` in app_setting.

~> **Note:** For application insight related settings, please use `application_insights_connection_string` and `application_insights_key`, terraform will assign the value to the key `APPINSIGHTS_INSTRUMENTATIONKEY` and `APPLICATIONINSIGHTS_CONNECTION_STRING` in app setting.

~> **Note:** For health check related settings, please use `health_check_eviction_time_in_min`, terraform will assign the value to the key `WEBSITE_HEALTHCHECK_MAXPINGFAILURES` in app setting.

~> **Note:** For those app settings that are deprecated or replaced by another properties for flex consumption function app, please check https://learn.microsoft.com/en-us/azure/azure-functions/functions-app-settings.

* `auth_settings` - (Optional) A `auth_settings` block as defined below.

* `auth_settings_v2` - (Optional) An `auth_settings_v2` block as defined below.

* `client_certificate_enabled` - (Optional) Should the function app use Client Certificates.

* `client_certificate_mode` - (Optional) The mode of the Function App's client certificates requirement for incoming requests. Possible values are `Required`, `Optional`, and `OptionalInteractiveUser`. Defaults to `Optional`.

* `client_certificate_exclusion_paths` - (Optional) Paths to exclude when using client certificates, separated by ;

* `connection_string` - (Optional) One or more `connection_string` blocks as defined below.

* `enabled` - (Optional) Is the Function App enabled? Defaults to `true`.

* `public_network_access_enabled` - (Optional) Should public network access be enabled for the Function App. Defaults to `true`.

* `identity` - (Optional) A `identity` block as defined below.

* `sticky_settings` - (Optional) A `sticky_settings` block as defined below.

* `storage_access_key` - (Optional) The access key which will be used to access the backend storage account for the Function App.

~> **Note:** The `storage_access_key` must be specified when `storage_authentication_type` is set to `StorageAccountConnectionString`.

* `storage_user_assigned_identity_id` - (Optional) The user assigned Managed Identity to access the storage account. Conflicts with `storage_access_key`.

~> **Note:** The `storage_user_assigned_identity_id` must be specified when `storage_authentication_type` is set to `UserAssignedIdentity`.

* `maximum_instance_count` - (Optional) The number of workers this function app can scale out to.

* `instance_memory_in_mb` - (Optional) The memory size of the instances on which your app runs. The [currently supported values](https://learn.microsoft.com/en-us/azure/azure-functions/flex-consumption-plan#instance-memory) are `2048` or `4096`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Linux Function App.

* `virtual_network_subnet_id` - (Optional) The subnet id which will be used by this Function App for [regional virtual network integration](https://docs.microsoft.com/en-us/azure/app-service/overview-vnet-integration#regional-virtual-network-integration).

~> **Note:** The AzureRM Terraform provider provides regional virtual network integration via the standalone resource [azurerm_app_service_virtual_network_swift_connection](app_service_virtual_network_swift_connection.html) and in-line within this resource using the `virtual_network_subnet_id` property. You cannot use both methods simultaneously. If the virtual network is set via the resource `app_service_virtual_network_swift_connection` then `ignore_changes` should be used in the function app configuration.

~> **Note:** Assigning the `virtual_network_subnet_id` property requires [RBAC permissions on the subnet](https://docs.microsoft.com/en-us/azure/app-service/overview-vnet-integration#permissions)

* `webdeploy_publish_basic_authentication_enabled` - (Optional) Should the default WebDeploy Basic Authentication publishing credentials enabled. Defaults to `true`.

~> **Note:** Setting this value to true will disable the ability to use `zip_deploy_file` which currently relies on the default publishing profile.

* `zip_deploy_file` - (Optional) The local path and filename of the Zip packaged application to deploy to this Linux Function App.

~> **Note:** Using this value requires either `WEBSITE_RUN_FROM_PACKAGE=1` or `SCM_DO_BUILD_DURING_DEPLOYMENT=true` to be set on the App in `app_settings`. Refer to the [Azure docs](https://learn.microsoft.com/en-us/azure/azure-functions/functions-deployment-technologies) for further details.

---

An `active_directory` block supports the following:

* `client_id` - (Required) The ID of the Client to use to authenticate with Azure Active Directory.

* `allowed_audiences` - (Optional) Specifies a list of Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

~> **Note:** The `client_id` value is always considered an allowed audience.

* `client_secret` - (Optional) The Client Secret for the Client ID. Cannot be used with `client_secret_setting_name`.

* `client_secret_setting_name` - (Optional) The App Setting name that contains the client secret of the Client. Cannot be used with `client_secret`.

---

An `app_service_logs` block supports the following:

* `disk_quota_mb` - (Optional) The amount of disk space to use for logs. Valid values are between `25` and `100`. Defaults to `35`.

* `retention_period_days` - (Optional) The retention period for logs in days. Valid values are between `0` and `99999`.(never delete).

~> **Note:** This block is not supported on Consumption plans.

---

An `auth_settings` block supports the following:

* `enabled` - (Required) Should the Authentication / Authorization feature be enabled for the Linux Web App?

* `active_directory` - (Optional) An `active_directory` block as defined above.

* `additional_login_parameters` - (Optional) Specifies a map of login Parameters to send to the OpenID Connect authorization endpoint when a user logs in.

* `allowed_external_redirect_urls` - (Optional) Specifies a list of External URLs that can be redirected to as part of logging in or logging out of the Linux Web App.

* `default_provider` - (Optional) The default authentication provider to use when multiple providers are configured. Possible values include: `AzureActiveDirectory`, `Facebook`, `Google`, `MicrosoftAccount`, `Twitter`, `Github`

~> **Note:** This setting is only needed if multiple providers are configured, and the `unauthenticated_client_action` is set to "RedirectToLoginPage".

* `facebook` - (Optional) A `facebook` block as defined below.

* `github` - (Optional) A `github` block as defined below.

* `google` - (Optional) A `google` block as defined below.

* `issuer` - (Optional) The OpenID Connect Issuer URI that represents the entity which issues access tokens for this Linux Web App.

~> **Note:** When using Azure Active Directory, this value is the URI of the directory tenant, e.g. <https://sts.windows.net/{tenant-guid}/>.

* `microsoft` - (Optional) A `microsoft` block as defined below.

* `runtime_version` - (Optional) The RuntimeVersion of the Authentication / Authorization feature in use for the Linux Web App.

* `token_refresh_extension_hours` - (Optional) The number of hours after session token expiration that a session token can be used to call the token refresh API. Defaults to `72` hours.

* `token_store_enabled` - (Optional) Should the Linux Web App durably store platform-specific security tokens that are obtained during login flows? Defaults to `false`.

* `twitter` - (Optional) A `twitter` block as defined below.

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

* `client_id` - (Required) The ID of the GitHub app used for login..

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

A `connection_string` block supports the following:

* `name` - (Required) The name which should be used for this Connection.

* `type` - (Required) Type of database. Possible values include: `MySQL`, `SQLServer`, `SQLAzure`, `Custom`, `NotificationHub`, `ServiceBus`, `EventHub`, `APIHub`, `DocDb`, `RedisCache`, and `PostgreSQL`.

* `value` - (Required) The connection string value.

---

A `cors` block supports the following:

* `allowed_origins` - (Optional) Specifies a list of origins that should be allowed to make cross-origin calls.

* `support_credentials` - (Optional) Are credentials allowed in CORS requests? Defaults to `false`.

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

* `client_secret` - (Optional) The client secret associated with the Google web application. Cannot be specified with `client_secret_setting_name`.

* `client_secret_setting_name` - (Optional) The app setting name that contains the `client_secret` value used for Google login. Cannot be specified with `client_secret`.

* `oauth_scopes` - (Optional) Specifies a list of OAuth 2.0 scopes that will be requested as part of Google Sign-In authentication. If not specified, `openid`, `profile`, and `email` are used as default scopes.

---

A `headers` block supports the following:

~> **Note:** Please see the [official Azure Documentation](https://docs.microsoft.com/azure/app-service/app-service-ip-restrictions#filter-by-http-header) for details on using header filtering.

* `x_azure_fdid` - (Optional) Specifies a list of Azure Front Door IDs.

* `x_fd_health_probe` - (Optional) Specifies if a Front Door Health Probe should be expected. The only possible value is `1`.

* `x_forwarded_for` - (Optional) Specifies a list of addresses for which matching should be applied. Omitting this value means allow any.

* `x_forwarded_host` - (Optional) Specifies a list of Hosts for which matching should be applied.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Linux Function App. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of User Assigned Managed Identity IDs to be assigned to this Linux Function App.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

An `ip_restriction` block supports the following:

* `action` - (Optional) The action to take. Possible values are `Allow` or `Deny`. Defaults to `Allow`.

* `headers` - (Optional) A `headers` block as defined above.

* `ip_address` - (Optional) The CIDR notation of the IP or IP Range to match. For example: `10.0.0.0/24` or `192.168.10.1/32`

* `name` - (Optional) The name which should be used for this `ip_restriction`.

* `priority` - (Optional) The priority value of this `ip_restriction`. Defaults to `65000`.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

~> **Note:** One and only one of `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified.

* `description` - (Optional) The Description of this IP Restriction.

---

A `microsoft` block supports the following:

* `client_id` - (Required) The OAuth 2.0 client ID that was created for the app used for authentication.

* `client_secret` - (Optional) The OAuth 2.0 client secret that was created for the app used for authentication. Cannot be specified with `client_secret_setting_name`.

* `client_secret_setting_name` - (Optional) The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication. Cannot be specified with `client_secret`.

* `oauth_scopes` - (Optional) Specifies a list of OAuth 2.0 scopes that will be requested as part of Microsoft Account authentication. If not specified, `wl.basic` is used as the default scope.

---

A `schedule` block supports the following:

* `frequency_interval` - (Required) How often the backup should be executed (e.g. for weekly backup, this should be set to `7` and `frequency_unit` should be set to `Day`).

~> **Note:** Not all intervals are supported on all Linux Function App SKUs. Please refer to the official documentation for appropriate values.

* `frequency_unit` - (Required) The unit of time for how often the backup should take place. Possible values include: `Day` and `Hour`.

* `keep_at_least_one_backup` - (Optional) Should the service keep at least one backup, regardless of age of backup. Defaults to `false`.

* `retention_period_days` - (Optional) After how many days backups should be deleted. Defaults to `30`.

* `start_time` - (Optional) When the schedule should start working in RFC-3339 format.

---

A `scm_ip_restriction` block supports the following:

* `action` - (Optional) The action to take. Possible values are `Allow` or `Deny`. Defaults to `Allow`.

* `headers` - (Optional) A `headers` block as defined above.

* `ip_address` - (Optional) The CIDR notation of the IP or IP Range to match. For example: `10.0.0.0/24` or `192.168.10.1/32`

* `name` - (Optional) The name which should be used for this `ip_restriction`.

* `priority` - (Optional) The priority value of this `ip_restriction`. Defaults to `65000`.

* `service_tag` - (Optional) The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - (Optional) The Virtual Network Subnet ID used for this IP Restriction.

~> **Note:** One and only one of `ip_address`, `service_tag` or `virtual_network_subnet_id` must be specified.

* `description` - (Optional) The Description of this IP Restriction.

---

A `site_config` block supports the following:

* `api_definition_url` - (Optional) The URL of the API definition that describes this Linux Function App.

* `api_management_api_id` - (Optional) The ID of the API Management API for this Linux Function App.

* `app_command_line` - (Optional) The App command line to launch.

* `application_insights_connection_string` - (Optional) The Connection String for linking the Linux Function App to Application Insights.

* `application_insights_key` - (Optional) The Instrumentation Key for connecting the Linux Function App to Application Insights.

* `app_service_logs` - (Optional) An `app_service_logs` block as defined above.

* `container_registry_managed_identity_client_id` - (Optional) The Client ID of the Managed Service Identity to use for connections to the Azure Container Registry.

* `container_registry_use_managed_identity` - (Optional) Should connections for Azure Container Registry use Managed Identity.

* `cors` - (Optional) A `cors` block as defined above.

* `default_documents` - (Optional) Specifies a list of Default Documents for the Linux Web App.

* `health_check_path` - (Optional) The path to be checked for this function app health.

* `health_check_eviction_time_in_min` - (Optional) The amount of time in minutes that a node can be unhealthy before being removed from the load balancer. Possible values are between `2` and `10`. Only valid in conjunction with `health_check_path`.

* `http2_enabled` - (Optional) Specifies if the HTTP2 protocol should be enabled. Defaults to `false`.

* `ip_restriction` - (Optional) One or more `ip_restriction` blocks as defined above.

* `ip_restriction_default_action` - (Optional) The Default action for traffic that does not match any `ip_restriction` rule. possible values include `Allow` and `Deny`. Defaults to `Allow`.

* `load_balancing_mode` - (Optional) The Site load balancing mode. Possible values include: `WeightedRoundRobin`, `LeastRequests`, `LeastResponseTime`, `WeightedTotalTraffic`, `RequestHash`, `PerSiteRoundRobin`. Defaults to `LeastRequests` if omitted.

* `managed_pipeline_mode` - (Optional) Managed pipeline mode. Possible values include: `Integrated`, `Classic`. Defaults to `Integrated`.

* `minimum_tls_version` - (Optional) The configures the minimum version of TLS required for SSL requests. Possible values include: `1.0`, `1.1`, `1.2` and `1.3`. Defaults to `1.2`.

* `remote_debugging_enabled` - (Optional) Should Remote Debugging be enabled. Defaults to `false`.

* `remote_debugging_version` - (Optional) The Remote Debugging Version. Possible values include `VS2017`, `VS2019`, and `VS2022`.

* `runtime_scale_monitoring_enabled` - (Optional) Should Scale Monitoring of the Functions Runtime be enabled?

~> **Note:** Functions runtime scale monitoring can only be enabled for Elastic Premium Function Apps or Workflow Standard Logic Apps and requires a minimum prewarmed instance count of 1.

* `scm_ip_restriction` - (Optional) One or more `scm_ip_restriction` blocks as defined above.

* `scm_ip_restriction_default_action` - (Optional) The Default action for traffic that does not match any `scm_ip_restriction` rule. possible values include `Allow` and `Deny`. Defaults to `Allow`.

* `scm_minimum_tls_version` - (Optional) The minimum version of TLS required for SSL requests to the SCM site. Possible values include `1.0`, `1.1`, `1.2` and `1.3`. Defaults to `1.2`.

* `scm_use_main_ip_restriction` - (Optional) Should the Linux Function App `ip_restriction` configuration be used for the SCM also.

* `use_32_bit_worker` - (Optional) Should the Linux Web App use a 32-bit worker. Defaults to `false`.

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

* `hosting_environment_id` - The ID of the App Service Environment used by Function App.

* `identity` - An `identity` block as defined below.

* `kind` - The Kind value for this Linux Function App.

* `outbound_ip_address_list` - A list of outbound IP addresses. For example `["52.23.25.3", "52.143.43.12"]`

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12`.

* `possible_outbound_ip_address_list` - A list of possible outbound IP addresses, not all of which are necessarily in use. This is a superset of `outbound_ip_address_list`. For example `["52.23.25.3", "52.143.43.12"]`.

* `possible_outbound_ip_addresses` - A comma separated list of possible outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12,52.143.43.17`. This is a superset of `outbound_ip_addresses`.

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

* `create` - (Defaults to 30 minutes) Used when creating the Function Flex Consumption App.
* `read` - (Defaults to 5 minutes) Used when retrieving the Function Flex Consumption App.
* `update` - (Defaults to 30 minutes) Used when updating the Function Flex Consumption App.
* `delete` - (Defaults to 30 minutes) Used when deleting the Function Flex Consumption App.

## Import

The Function Apps can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_function_app_flex_consumption.example /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1
```
