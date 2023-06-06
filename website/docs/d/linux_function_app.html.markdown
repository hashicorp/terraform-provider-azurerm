---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_linux_function_app"
description: |-
  Gets information about an existing Linux Function App.
---

# Data Source: azurerm_linux_function_app

Use this data source to access information about an existing Linux Function App.

## Example Usage

```hcl
data "azurerm_linux_function_app" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_linux_function_app.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Linux Function App.

* `resource_group_name` - (Required) The name of the Resource Group where the Linux Function App should exist.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Linux Function App.

* `location` -  The Azure Region where the Linux Function App exists.

* `app_settings` - A map of key-value pairs for [App Settings](https://docs.microsoft.com/azure/azure-functions/functions-app-settings) and custom values.

* `auth_settings` - A `auth_settings` block as defined below.

* `auth_settings_v2` - A `auth_settings_v2` block as defined below.

* `availability` - The current availability state. Possible values are `Normal`, `Limited`, and `DisasterRecoveryMode`.

* `backup` - A `backup` block as defined below.

* `builtin_logging_enabled` - Is built in logging enabled?

* `client_certificate_enabled` - Are Client Certificates enabled?

* `client_certificate_mode` -  The mode of the Function App's client certificates requirement for incoming requests.

* `client_certificate_exclusion_paths` - Paths to exclude when using client certificates, separated by ;

* `connection_string` -  A `connection_string` blocks as defined below.

* `content_share_force_disabled` - Are the settings for linking the Function App to storage suppressed?

* `custom_domain_verification_id` - The identifier used by App Service to perform domain ownership verification via DNS TXT record.

* `daily_memory_time_quota` -  The amount of memory in gigabyte-seconds that your application is allowed to consume per day.

* `default_hostname` - The default hostname of the Linux Function App.

* `hosting_environment_id` - The ID of the App Service Environment used by Function App.

* `enabled` - Is the Function App enabled?

* `functions_extension_version` - The runtime version associated with the Function App.

* `https_only` - Can the Function App only be accessed via HTTPS?

* `identity` - A `identity` block as defined below.

* `kind` - The Kind value for this Linux Function App.

* `outbound_ip_address_list` - A list of outbound IP addresses. For example `["52.23.25.3", "52.143.43.12"]`

* `outbound_ip_addresses` - A comma separated list of outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12`.

* `possible_outbound_ip_address_list` - A list of possible outbound IP addresses, not all of which are necessarily in use. This is a superset of `outbound_ip_address_list`. For example `["52.23.25.3", "52.143.43.12"]`.

* `possible_outbound_ip_addresses` - A comma separated list of possible outbound IP addresses as a string. For example `52.23.25.3,52.143.43.12,52.143.43.17`. This is a superset of `outbound_ip_addresses`.

* `service_plan_id` - The ID of the App Service Plan within which this Function App has been created.

* `site_config` -  A `site_config` block as defined below.

* `site_credential` - A `site_credential` block as defined below.

* `sticky_settings` - A `sticky_settings` block as defined below.

* `storage_account_name` - The backend storage account name used by this Function App.

* `storage_account_access_key` -  The access key used to access the backend storage account for the Function App.

* `storage_key_vault_secret_id` - The Key Vault Secret ID, including version, that contains the Connection String to connect to the storage account for this Function App.

* `storage_uses_managed_identity` - Does the Function App use Managed Identity to access the storage account?

* `tags` - A mapping of tags which are assigned to the Linux Function App.

* `usage` - The current usage state. Possible values are `Normal` and `Exceeded`.

* `virtual_network_subnet_id` - The subnet id which the Linux Function App is vNet Integrated with.

---

An `active_directory` block exports the following:

* `client_id` - The ID of the Client used to authenticate with Azure Active Directory.

* `allowed_audiences` - A list of Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

* `client_secret` -  The Client Secret of the Client ID.

* `client_secret_setting_name` - The App Setting name that contains the client secret of the Client.

---

A `application_stack` block exports the following:

* `docker` -  One or more `docker` blocks as defined below.

* `dotnet_version` -  The version of .NET used.

* `java_version` - The Version of Java used.

* `node_version` - The version of Node used.

* `python_version` - The version of Python used.

* `powershell_core_version` - The version of PowerShell Core used.

* `use_custom_runtime` - Does the Linux Function App use a custom runtime?

---

An `app_service_logs` block exports the following:

* `disk_quota_mb` -  The amount of disk space used for logs.

* `retention_period_days` - The retention period for logs in days.

---

An `auth_settings` block exports the following:

* `enabled` -  Is the Authentication / Authorization feature enabled for the Linux Web App?

* `active_directory` - An `active_directory` block as defined above.

* `additional_login_parameters` - A map of login parameters sent to the OpenID Connect authorization endpoint when a user logs in.

* `allowed_external_redirect_urls` - A list of External URLs that can be redirected to as part of logging in or logging out of the Linux Web App.

* `default_provider` - The default authentication provider used when multiple providers are configured.

* `facebook` - A `facebook` block as defined below.

* `github` - A `github` block as defined below.

* `google` - A `google` block as defined below.

* `issuer` - The OpenID Connect Issuer URI that represents the entity which issues access tokens for this Linux Web App.

* `microsoft` - A `microsoft` block as defined below.

* `runtime_version` - The RuntimeVersion of the Authentication / Authorization feature in use for the Linux Web App.

* `token_refresh_extension_hours` - The number of hours after session token expiration that a session token can be used to call the token refresh API.

* `token_store_enabled` - Does the Linux Web App durably store platform-specific security tokens that are obtained during login flows?

* `twitter` - A `twitter` block as defined below.

* `unauthenticated_client_action` - The action to taken when an unauthenticated client attempts to access the app.

---

An `auth_settings_v2` block exports the following:

* `auth_enabled` - Are the AuthV2 Settings enabled.

* `runtime_version` - The Runtime Version of the Authentication and Authorisation feature of this App.

* `config_file_path` - The path to the App Auth settings.

* `require_authentication` - Is the authentication flow used for all requests.

* `unauthenticated_action` - The action to take for requests made without authentication.

* `default_provider` -The Default Authentication Provider used when more than one Authentication Provider is configured and the `unauthenticated_action` is set to `RedirectToLoginPage`.

* `excluded_paths` - The paths which should be excluded from the `unauthenticated_action` when it is set to `RedirectToLoginPage`.

* `require_https` -Is HTTPS required on connections?

* `http_route_api_prefix` - The prefix that should precede all the authentication and authorisation paths.

* `forward_proxy_convention` - The convention used to determine the url of the request made.

* `forward_proxy_custom_host_header_name` -The name of the custom header containing the host of the request.

* `forward_proxy_custom_scheme_header_name` - The name of the custom header containing the scheme of the request.

* `apple_v2` - An `apple_v2` block as defined below.

* `active_directory_v2` - An `active_directory_v2` block as defined below.

* `azure_static_web_app_v2` - An `azure_static_web_app_v2` block as defined below.

* `custom_oidc_v2` - Zero or more `custom_oidc_v2` blocks as defined below.

* `facebook_v2` - A `facebook_v2` block as defined below.

* `github_v2` - A `github_v2` block as defined below.

* `google_v2` - A `google_v2` block as defined below.

* `microsoft_v2` - A `microsoft_v2` block as defined below.

* `twitter_v2` - A `twitter_v2` block as defined below.

* `login` - A `login` block as defined below.

---

An `apple_v2` block supports the following:

* `client_id` - The OpenID Connect Client ID for the Apple web application.

* `client_secret_setting_name` - The app setting name that contains the `client_secret` value used for Apple Login.

* `login_scopes` - A list of Login Scopes provided by this Authentication Provider.

---

An `active_directory_v2` block supports the following:

* `client_id` - The ID of the Client used to authenticate with Azure Active Directory.

* `tenant_auth_endpoint` - The Azure Tenant Endpoint for the Authenticating Tenant. e.g. `https://login.microsoftonline.com/v2.0/{tenant-guid}/`

* `client_secret_setting_name` - The App Setting name that contains the client secret of the Client.

* `client_secret_certificate_thumbprint` - The thumbprint of the certificate used for signing purposes.

* `jwt_allowed_groups` - The list of Allowed Groups in the JWT Claim.

* `jwt_allowed_client_applications` - The list of Allowed Client Applications in the JWT Claim.

* `www_authentication_disabled` - Is the www-authenticate provider omitted from the request?

* `allowed_groups` -The list of allowed Group Names for the Default Authorisation Policy.

* `allowed_identities` - The list of allowed Identities for the Default Authorisation Policy.

* `allowed_applications` - The list of allowed Applications for the Default Authorisation Policy.

* `login_parameters` - A map of key-value pairs sent to the Authorisation Endpoint when a user logs in.

* `allowed_audiences` - Specifies a list of Allowed audience values to consider when validating JWTs issued by Azure Active Directory.

---

An `azure_static_web_app_v2` block supports the following:

* `client_id` - The ID of the Client to use to authenticate with Azure Static Web App Authentication.

---

A `custom_oidc_v2` block supports the following:

* `name` - The name of the Custom OIDC Authentication Provider.

* `client_id` - The ID of the Client to use to authenticate with the Custom OIDC.

* `openid_configuration_endpoint` - The app setting name that contains the `client_secret` value used for the Custom OIDC Login.

* `name_claim_type` - The name of the claim that contains the users name.

* `scopes` - The list of the scopes that are requested while authenticating.

* `client_credential_method` - The Client Credential Method used.

* `client_secret_setting_name` - The App Setting name that contains the secret for this Custom OIDC Client. This is generated from `name` above and suffixed with `_PROVIDER_AUTHENTICATION_SECRET`.

* `authorisation_endpoint` - The endpoint to make the Authorisation Request as supplied by `openid_configuration_endpoint` response.

* `token_endpoint` - The endpoint used to request a Token as supplied by `openid_configuration_endpoint` response.

* `issuer_endpoint` - The endpoint that issued the Token as supplied by `openid_configuration_endpoint` response.

* `certification_uri` - The endpoint that provides the keys necessary to validate the token as supplied by `openid_configuration_endpoint` response.

---

A `facebook_v2` block supports the following:

* `app_id` - The App ID of the Facebook app used for login.

* `app_secret_setting_name` - The app setting name that contains the `app_secret` value used for Facebook Login.

* `graph_api_version` - The version of the Facebook API to be used while logging in.

* `login_scopes` - The list of scopes that are requested as part of Facebook Login authentication.

---

A `github_v2` block supports the following:

* `client_id` - The ID of the GitHub app used for login..

* `client_secret_setting_name` - The app setting name that contains the `client_secret` value used for GitHub Login.

* `login_scopes` - The list of OAuth 2.0 scopes that are requested as part of GitHub Login authentication.

---

A `google_v2` block supports the following:

* `client_id` - The OpenID Connect Client ID for the Google web application.

* `client_secret_setting_name` - The app setting name that contains the `client_secret` value used for Google Login.

* `allowed_audiences` - The list of Allowed Audiences that are requested as part of Google Sign-In authentication.

* `login_scopes` - (Optional) The list of OAuth 2.0 scopes that should be requested as part of Google Sign-In authentication.

---

A `microsoft_v2` block supports the following:

* `client_id` - The OAuth 2.0 client ID that was created for the app used for authentication.

* `client_secret_setting_name` - The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication.

* `allowed_audiences` - The list of Allowed Audiences that are be requested as part of Microsoft Sign-In authentication.

* `login_scopes` - The list of Login scopes that are requested as part of Microsoft Account authentication.

---

A `twitter_v2` block supports the following:

* `consumer_key` - The OAuth 1.0a consumer key of the Twitter application used for sign-in.

* `consumer_secret_setting_name` - The app setting name that contains the OAuth 1.0a consumer secret of the Twitter application used for sign-in.

---

A `login` block supports the following:

* `logout_endpoint` - The endpoint to which logout requests are made.

* `token_store_enabled` - Is the Token Store configuration Enabled.

* `token_refresh_extension_time` - The number of hours after session token expiration that a session token can be used to call the token refresh API.

* `token_store_path` - The directory path in the App Filesystem in which the tokens are stored.

* `token_store_sas_setting_name` - The name of the app setting which contains the SAS URL of the blob storage containing the tokens.

* `preserve_url_fragments_for_logins` - Are the fragments from the request preserved after the login request is made.

* `allowed_external_redirect_urls` - External URLs that can be redirected to as part of logging in or logging out of the app.

* `cookie_expiration_convention` - The method by which cookies expire.

* `cookie_expiration_time` - The time after the request is made when the session cookie should expire.

* `validate_nonce` - Is the nonce validated while completing the login flow.

* `nonce_expiration_time` - The time after the request is made when the nonce should expire.

---

A `backup` block exports the following:

* `name` - The name of this Backup.

* `schedule` - A `schedule` block as defined below.

* `storage_account_url` - The SAS URL to the container.

* `enabled` - Is this backup job enabled?

---

A `connection_string` block exports the following:

* `name` - The name of this Connection.

* `type` -  Type of database.

* `value` - The connection string value.

---

A `cors` block exports the following:

* `allowed_origins` - A list of origins that are allowed to make cross-origin calls.

* `support_credentials` - Are credentials allowed in CORS requests?

---

A `docker` block exports the following:

* `registry_url` - The URL of the docker registry.

* `image_name` -  The name of the Docker image used.

* `image_tag` - The image tag of the image used.

* `registry_username` - The username used for connections to the registry.

* `registry_password` - The password for the account to use to connect to the registry.

---

A `facebook` block exports the following:

* `app_id` - The App ID of the Facebook app used for login.

* `app_secret` - The App Secret of the Facebook app used for Facebook login.

* `app_secret_setting_name` - The app setting name that contains the `app_secret` value used for Facebook login.

* `oauth_scopes` - Specifies a list of OAuth 2.0 scopes requested as part of Facebook login authentication.

---

A `github` block exports the following:

* `client_id` - The ID of the GitHub app used for login.

* `client_secret` - The Client Secret of the GitHub app used for GitHub login.

* `client_secret_setting_name` - The app setting name that contains the `client_secret` value used for GitHub login.

* `oauth_scopes` - Specifies a list of OAuth 2.0 scopes that are requested as part of GitHub login authentication.

---

A `google` block exports the following:

* `client_id` - The OpenID Connect Client ID for the Google web application.

* `client_secret` - The client secret associated with the Google web application.

* `client_secret_setting_name` - The app setting name that contains the `client_secret` value used for Google login.

* `oauth_scopes` - A list of OAuth 2.0 scopes that are requested as part of Google Sign-In authentication.

---

A `headers` block exports the following:

* `x_azure_fdid` - A list of Azure Front Door IDs.

* `x_fd_health_probe` - Should a Front Door Health Probe be expected?

* `x_forwarded_for` - A list of addresses for which matching is applied.

* `x_forwarded_host` - A list of Hosts for which matching is applied.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Linux Function App.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Linux Function App.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Linux Function App.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Linux Function App.

---

An `ip_restriction` block exports the following:

* `action` - The action to take.

* `headers` - A `headers` block as defined above.

* `ip_address` -  The CIDR notation of the IP or IP Range that is matched.

* `name` - The name which is used for this `ip_restriction`.

* `priority` - The priority value of this `ip_restriction`.

* `service_tag` - The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - The Virtual Network Subnet ID used for this IP Restriction.

---

A `microsoft` block exports the following:

* `client_id` -  The OAuth 2.0 client ID that was created for the app used for authentication.

* `client_secret` -  The OAuth 2.0 client secret that was created for the app used for authentication.

* `client_secret_setting_name` - The app setting name containing the OAuth 2.0 client secret that was created for the app used for authentication.

* `oauth_scopes` - A list of OAuth 2.0 scopes that will be requested as part of Microsoft Account authentication.

---

A `schedule` block exports the following:

* `frequency_interval` -  How often the backup is executed.

* `frequency_unit` - The unit of time for how often the backup takes place.

* `keep_at_least_one_backup` - Does the service keep at least one backup, regardless of age of backup?

* `retention_period_days` - After how many days backups are deleted.

* `start_time` -  When the schedule starts working in RFC-3339 format.

---

A `scm_ip_restriction` block exports the following:

* `action` - The action taken.

* `headers` - A `headers` block as defined above.

* `ip_address` - The CIDR notation of the IP or IP Range matched.

* `name` - The name used for this `ip_restriction`.

* `priority` - The priority value of this `ip_restriction`.

* `service_tag` - The Service Tag used for this IP Restriction.

* `virtual_network_subnet_id` - The Virtual Network Subnet ID used for this IP Restriction.

---

A `sticky_settings` block exports the following:

* `app_setting_names` - A list of `app_setting` names that the Linux Function App will not swap between Slots when a swap operation is triggered.

* `connection_string_names` - A list of `connection_string` names that the Linux Function App will not swap between Slots when a swap operation is triggered.

---

A `site_config` block exports the following:

* `always_on` - If this Linux Web App is Always On enabled.

* `api_definition_url` -  The URL of the API definition that describes this Linux Function App.

* `api_management_api_id` - The ID of the API Management API for this Linux Function App.

* `app_command_line` -  The App command line that is launched.

* `app_scale_limit` - The number of workers this function app can scale out to.

* `application_insights_connection_string` - The Connection String that links the Linux Function App to Application Insights.

* `application_insights_key` -  The Instrumentation Key that connects the Linux Function App to Application Insights.

* `application_stack` -  An `application_stack` block as defined above.

* `app_service_logs` - An `app_service_logs` block as defined above.

* `auto_swap_slot_name` -  The Linux Function App Slot Name that is automatically swapped to when deployment to that slot is successfully completed.

* `container_registry_managed_identity_client_id` - The Client ID of the Managed Service Identity that is used for connections to the Azure Container Registry.

* `container_registry_use_managed_identity` - Do connections for Azure Container Registry use Managed Identity?

* `cors` -  A `cors` block as defined above.

* `default_documents` -  A list of Default Documents for the Linux Web App.

* `elastic_instance_minimum` -  The number of minimum instances for this Linux Function App.

* `ftps_state` - State of FTP / FTPS service for this function app.

* `health_check_path` - The path that is checked for this function app health.

* `health_check_eviction_time_in_min` - The amount of time in minutes that a node can be unhealthy before being removed from the load balancer.

* `http2_enabled` - Is the HTTP2 protocol enabled?

* `ip_restriction` - One or more `ip_restriction` blocks as defined above.

* `load_balancing_mode` -  The Site load balancing mode.

* `managed_pipeline_mode` - Managed pipeline mode.

* `minimum_tls_version` -  The minimum version of TLS required for SSL requests.

* `pre_warmed_instance_count` - The number of pre-warmed instances for this function app.

* `remote_debugging_enabled` -  Is Remote Debugging enabled?

* `remote_debugging_version` - The Remote Debugging Version.

* `runtime_scale_monitoring_enabled` - Is Scale Monitoring of the Functions Runtime enabled?

* `scm_ip_restriction` - One or more `scm_ip_restriction` blocks as defined above.

* `scm_minimum_tls_version` - The minimum version of TLS for SSL requests to the SCM site.

* `scm_use_main_ip_restriction` -  Is the Linux Function App `ip_restriction` configuration used for the SCM also?

* `use_32_bit_worker` - Does the Linux Web App use a 32-bit worker process?

* `vnet_route_all_enabled` - Are all outbound traffic to NAT Gateways, Network Security Groups and User Defined Routes applied?

* `websockets_enabled` - Are Web Sockets enabled?

* `worker_count` - The number of Workers for this Linux Function App.

---

A `twitter` block exports the following:

* `consumer_key` - The OAuth 1.0a consumer key of the Twitter application used for sign-in.

* `consumer_secret` - The OAuth 1.0a consumer secret of the Twitter application used for sign-in.

* `consumer_secret_setting_name` - The app setting name that contains the OAuth 1.0a consumer secret of the Twitter application used for sign-in.

---

A `site_credential` block exports the following:

* `name` - The Site Credentials Username used for publishing.

* `password` - The Site Credentials Password used for publishing.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 25 minutes) Used when retrieving the Linux Function App.
