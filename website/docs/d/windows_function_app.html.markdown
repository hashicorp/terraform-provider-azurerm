---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_windows_function_app"
description: |-
  Gets information about an existing Windows Function App.
---

# Data Source: azurerm_windows_function_app

Use this data source to access information about an existing Windows Function App.

## Example Usage

```hcl
data "azurerm_windows_function_app" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_windows_function_app.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Windows Function App.

* `resource_group_name` - (Required) The name of the Resource Group where the Windows Function App exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Windows Function App.

* `app_settings` - A `map of key-value pairs for App Settings and custom values.

* `auth_settings` - A `auth_settings` block as defined below.

* `backup` - A `backup` block as defined below.

* `builtin_logging_enabled` - Is the built-in logging enabled?

* `client_certificate_enabled` - Is the use of Client Certificates enabled?

* `client_certificate_mode` - The mode of the Function App's client certificates requirement for incoming requests.

* `connection_string` - One or more `connection_string` blocks as defined below.

* `content_share_force_disabled` - Are Content Share Settings disabled?

* `custom_domain_verification_id` - The identifier used by App Service to perform domain ownership verification via DNS TXT record.

* `daily_memory_time_quota` - The amount of memory in gigabyte-seconds that your application is allowed to consume per day.

* `default_hostname` - The default hostname of the Windows Function App.

* `enabled` - Is the Function App enabled?

* `functions_extension_version` - The runtime version associated with the Function App.

* `https_only` - Is the Function App only accessible via HTTPS?

* `identity` - A `identity` block as defined below.

* `kind` - The Kind value for this Windows Function App.

* `location` - The Azure Region where the Windows Function App exists.

* `outbound_ip_address_list` - A list of outbound IP addresses.

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12`.

* `possible_outbound_ip_address_list` - AA list of possible outbound IP addresses, not all of which are necessarily in use.

* `possible_outbound_ip_addresses` - A list of possible outbound IP addresses, not all of which are necessarily in use. This is a superset of `outbound_ip_address_list`. For example `["52.23.25.3", "52.143.43.12"]`.

* `service_plan_id` - The ID of the App Service Plan.

* `site_config` - A `site_config` block as defined below.

* `site_credential` - A `site_credential` block as defined below.

* `sticky_settings` - A `sticky_settings` block as defined below.

* `storage_account_access_key` - The access key which is used to access the backend storage account for the Function App.

* `storage_account_name` - The backend storage account name which is used by this Function App.

* `storage_key_vault_secret_id` - The Key Vault Secret ID, including version, that contains the Connection String used to connect to the storage account for this Function App.

* `storage_uses_managed_identity` - Is the Function App using a Managed Identity to access the storage account?

* `tags` - A mapping of tags assigned to the Windows Function App.

---

A `active_directory` block exports the following:

* `allowed_audiences` - A list of Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

* `client_id` - The ID of the Client to use to authenticate with Azure Active Directory.

* `client_secret` - The Client Secret for the Client ID.

* `client_secret_setting_name` - The App Setting name that contains the client secret of the Client.

---

A `app_service_logs` block exports the following:

* `disk_quota_mb` - The amount of disk space to use for logs.

* `retention_period_days` - The retention period for logs in days.

---

A `application_stack` block exports the following:

* `dotnet_version` - The version of .Net to use.

* `java_version` - The version of Java to use.

* `node_version` - The version of Node to use.

* `powershell_core_version` - The version of PowerShell Core to use.

* `use_custom_runtime` - Is the Windows Function App using a custom runtime?.

---

A `auth_settings` block exports the following:

* `active_directory` - A `active_directory` block as defined above.

* `additional_login_parameters` - A map of Login Parameters to send to the OpenID Connect authorization endpoint when a user logs in.

* `allowed_external_redirect_urls` - A list of External URLs that can be redirected to as part of logging in or logging out of the Windows Function App.

* `default_provider` - The default authentication provider to use when multiple providers are configured.

* `enabled` - Is the Authentication / Authorization feature for the Windows Function enabled?

* `facebook` - A `facebook` block as defined below.

* `github` - A `github` block as defined below.

* `google` - A `google` block as defined below.

* `issuer` - The OpenID Connect Issuer URI that represents the entity which issues access tokens for this Windows Function App.

* `microsoft` - A `microsoft` block as defined below.

* `runtime_version` - The Runtime Version of the Authentication / Authorization feature in use for the Windows Function App.

* `token_refresh_extension_hours` - The number of hours after session token expiration that a session token can be used to call the token refresh API.

* `token_store_enabled` - Is the durable storing of platform-specific security token that are obtained during login flows enabled?

* `twitter` - A `twitter` block as defined below.

* `unauthenticated_client_action` - The action to take when an unauthenticated client attempts to access the app.

---

A `backup` block exports the following:

* `enabled` - Is the Backup Job enabled?

* `name` - The name of this Backup.

* `schedule` - A `schedule` block as defined below.

* `storage_account_url` - The SAS URL to the container.

---

A `connection_string` block exports the following:

* `name` - The name of this Connection.

* `type` - Type of database.

* `value` - The connection string value.

---

A `cors` block exports the following:

* `allowed_origins` - A list of origins that should be allowed to make cross-origin calls.

* `support_credentials` - Are credentials allows in CORS requests?.

---

A `facebook` block exports the following:

* `app_id` - The App ID of the Facebook app used for login.

* `app_secret` - The App Secret of the Facebook app used for Facebook Login.

* `app_secret_setting_name` - The app setting name that contains the `app_secret` value used for Facebook Login.

* `oauth_scopes` - A list of OAuth 2.0 scopes to be requested as part of Facebook Login authentication.

---

A `github` block exports the following:

* `client_id` - The ID of the GitHub app used for login.

* `client_secret` - The Client Secret of the GitHub app used for GitHub Login.

* `client_secret_setting_name` - The app setting name that contains the `client_secret` value used for GitHub Login.

* `oauth_scopes` - A list of OAuth 2.0 scopes that will be requested as part of GitHub Login authentication.

---

A `google` block exports the following:

* `client_id` - The OpenID Connect Client ID for the Google web application.

* `client_secret` - The client secret associated with the Google web application.

* `client_secret_setting_name` - The app setting name that contains the `client_secret` value used for Google Login.

* `oauth_scopes` - A list of OAuth 2.0 scopes that will be requested as part of Google Sign-In authentication.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Windows Function App.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Windows Function App.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Windows Function App.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Windows Function App.

---

A `microsoft` block exports the following:

* `client_id` - The OAuth 2.0 client ID that was created for the app used for authentication.

* `client_secret` - The OAuth 2.0 client secret that was created for the app used for authentication.

* `client_secret_setting_name` - The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication.

* `oauth_scopes` - A list of OAuth 2.0 scopes that will be requested as part of Microsoft Account authentication.

---

A `schedule` block exports the following:

* `frequency_interval` - How often the backup is executed.

* `frequency_unit` - The unit of time the backup should take place.

* `keep_at_least_one_backup` - Should the service keep at least one backup.

* `retention_period_days` - After how many days backups is deleted.

* `start_time` - When the schedule should start working in RFC-3339 format.

---

A `site_config` block exports the following:

* `always_on` - Is this Windows Function App Always On?.

* `api_definition_url` - The URL of the API definition that describes this Windows Function App.

* `api_management_api_id` - The ID of the API Management API for this Windows Function App.

* `app_command_line` - The App command line to launch.

* `app_scale_limit` - The number of workers this function app can scale out to.

* `app_service_logs` - A `app_service_logs` block as defined above.

* `application_insights_connection_string` - The Connection String for linking the Windows Function App to Application Insights.

* `application_insights_key` - The Instrumentation Key for connecting the Windows Function App to Application Insights.

* `application_stack` - A `application_stack` block as defined above.

* `cors` - A `cors` block as defined above.

* `default_documents` - A list of Default Documents for the Windows Web App.

* `detailed_error_logging_enabled` - Is detailed error logging enabled?

* `elastic_instance_minimum` - The number of minimum instances for this Windows Function App.

* `ftps_state` - State of FTP / FTPS service for this Windows Function App.

* `health_check_eviction_time_in_min` - The amount of time in minutes that a node can be unhealthy before being removed from the load balancer.

* `health_check_path` - The path to be checked for this Windows Function App health.

* `http2_enabled` - Is the HTTP2 protocol enabled?

* `ip_restriction` - One or more `ip_restriction` blocks as defined above.

* `load_balancing_mode` - The Site load balancing mode.

* `managed_pipeline_mode` - The Managed pipeline mode.

* `minimum_tls_version` - The minimum version of TLS required for SSL requests.

* `pre_warmed_instance_count` - The number of pre-warmed instances for this Windows Function App.

* `remote_debugging_enabled` - Is Remote Debugging enabled?

* `remote_debugging_version` - The Remote Debugging Version.

* `runtime_scale_monitoring_enabled` - Is Scale Monitoring of the Functions Runtime enabled?

* `scm_ip_restriction` - One or more `scm_ip_restriction` blocks as defined above.

* `scm_minimum_tls_version` - The minimum version of TLS required for SSL requests to the SCM site.

* `scm_type` - The SCM type.

* `scm_use_main_ip_restriction` - Is the `ip_restriction` configuration used for the SCM?.

* `use_32_bit_worker` - Is the Windows Function App using a 32-bit worker process?

* `vnet_route_all_enabled` - Are all outbound traffic to NAT Gateways, Network Security Groups and User Defined Routes applied?

* `websockets_enabled` - Are Web Sockets enabled?

* `windows_fx_version` - The Windows FX version.

* `worker_count` - The number of Workers for this Windows Function App.

---

A `site_credential` block exports the following:

* `name` - The Site Credentials Username used for publishing.

* `password` - The Site Credentials Password used for publishing.

---

A `sticky_settings` block exports the following:

* `app_setting_names` - A list of `app_setting` names that the Windows Function App will not swap between Slots when a swap operation is triggered.

* `connection_string_names` - A list of `connection_string` names that the Windows Function App will not swap between Slots when a swap operation is triggered.

---

A `twitter` block exports the following:

* `consumer_key` - The OAuth 1.0a consumer key of the Twitter application used for sign-in.

* `consumer_secret` - The OAuth 1.0a consumer secret of the Twitter application used for sign-in.

* `consumer_secret_setting_name` - The app setting name that contains the OAuth 1.0a consumer secret of the Twitter application used for sign-in.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 10 minutes) Used when retrieving the Windows Function App.
