---
subcategory: "App Service (Web Apps)"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_windows_web_app"
description: |-
  Gets information about an existing Windows Web App.
---

# Data Source: azurerm_windows_web_app

Use this data source to access information about an existing Windows Web App.

## Example Usage

```hcl
data "azurerm_windows_web_app" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_windows_web_app.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Windows Web App.

* `resource_group_name` - (Required) The name of the Resource Group where the Windows Web App exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Windows Web App.

* `app_settings` - A map of key-value pairs of App Settings.

* `auth_settings` - A `auth_settings` block as defined below.

* `auth_settings_v2` - An `auth_settings_v2` block as defined below.

* `backup` - A `backup` block as defined below.

* `client_affinity_enabled` - Is Client Affinity enabled?

* `client_certificate_enabled` - Are Client Certificates enabled?

* `client_certificate_mode` - The Client Certificate mode.

* `client_certificate_exclusion_paths` - Paths to exclude when using client certificates, separated by ;

* `connection_string` - A `connection_string` block as defined below.

* `custom_domain_verification_id` - The identifier used by App Service to perform domain ownership verification via DNS TXT record.

* `hosting_environment_id` - The ID of the App Service Environment used by App Service.

* `default_hostname` - The Default Hostname of the Windows Web App.

* `enabled` - Is the Windows Web App enabled?

* `https_only` - Does the Windows Web App require HTTPS connections.

* `identity` - A `identity` block as defined below.

* `kind` - The string representation of the Windows Web App Kind.

* `location` - The Azure Region where the Windows Web App exists.

* `logs` - A `logs` block as defined below.

* `outbound_ip_address_list` - The list of Outbound IP Addresses for this Windows Web App.

* `outbound_ip_addresses` - A string representation of the list of Outbound IP Addresses for this Windows Web App.

* `possible_outbound_ip_address_list` - The list of Possible Outbound IP Addresses that could be used by this Windows Web App.

* `possible_outbound_ip_addresses` - The string representation of the list of Possible Outbound IP Addresses that could be used by this Windows Web App.

* `service_plan_id` - The ID of the Service Plan in which this Windows Web App resides.

* `site_config` - A `site_config` block as defined below.

* `site_credential` - A `site_credential` block as defined below.

* `sticky_settings` - A `sticky_settings` block as defined below.

* `storage_account` - A `storage_account` block as defined below.

* `tags` - A mapping of tags assigned to the Windows Web App.

* `virtual_network_subnet_id` - The subnet id which the Windows Web App is vNet Integrated with.

---

A `action` block exports the following:

* `action_type` - The predefined action to be taken to an Auto Heal trigger.

* `custom_action` - A `custom_action` block as defined below.

* `minimum_process_execution_time` - The minimum amount of time in `hh:mm:ss` the Windows Web App must have been running before the defined action will be run in the event of a trigger.

---

An `active_directory` block exports the following:

* `allowed_audiences` - An `allowed_audiences` block as defined below.

* `client_id` -  The ID of the Client used to authenticate with Azure Active Directory.

* `client_secret` - The Client Secret for the Client ID.

* `client_secret_setting_name` - The App Setting name that contains the client secret of the Client.

---

An `application_logs` block exports the following:

* `azure_blob_storage` - An `azure_blob_storage` block as defined below.

* `file_system_level` - The logging level.

---

An `application_stack` block exports the following:

* `current_stack` - The Current Stack value of the Windows Web App.

* `docker_image_name` - The docker image, including tag, used by this Windows Web App.

* `docker_registry_url` - The URL of the container registry where the `docker_image_name` is located.

* `docker_registry_username` - The User Name to use for authentication against the registry to pull the image.

* `docker_registry_password` - The User Name to use for authentication against the registry to pull the image.

* `dotnet_version` - The version of .NET in use.

* `java_container` - The Java Container in use.

* `java_container_version` - The Version of the Java Container in use.

* `java_version` - The Version of Java in use.

* `node_version` - The Version of Node in use.

* `php_version` - The Version of the PHP in use.

* `python_version` - The Version of Python in use.

---

A `auth_settings` block exports the following:

* `active_directory` - A `active_directory` block as defined above.

* `additional_login_parameters` - A `additional_login_parameters` block as defined above.

* `allowed_external_redirect_urls` - A `allowed_external_redirect_urls` block as defined above.

* `default_provider` - The default authentication provider in use when multiple providers are configured.

* `enabled` - Is the Authentication / Authorization feature is enabled for the Windows Web App?

* `facebook` - A `facebook` block as defined below.

* `github` - A `github` block as defined below.

* `google` - A `google` block as defined below.

* `issuer` - The OpenID Connect Issuer URI that represents the entity which issues access tokens for this Windows Web App.

* `microsoft` - A `microsoft` block as defined below.

* `runtime_version` - The RuntimeVersion of the Authentication / Authorization feature in use for the Windows Web App.

* `token_refresh_extension_hours` - The number of hours after session token expiration that a session token can be used to call the token refresh API.

* `token_store_enabled` - Does Windows Web App durably store platform-specific security tokens that are obtained during login flows enabled?

* `twitter` - A `twitter` block as defined below.

* `unauthenticated_client_action` - The action to take when an unauthenticated client attempts to access the app.

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

A `auto_heal_setting` block exports the following:

* `action` - A `action` block as defined above.

* `trigger` - A `trigger` block as defined below.

---

A `azure_blob_storage` block exports the following:

* `level` - The level at which to log. Possible values include `Error`, `Warning`, `Information`, `Verbose` and `Off`. **NOTE:** this field is not available for `http_logs`

* `retention_in_days` - The time in days after which blobs will be removed.

* `sas_url` - The SAS url to the Azure Blob container.

---

A `backup` block exports the following:

* `enabled` - Is the Backup enabled?

* `name` - The name of this Backup.

* `schedule` - A `schedule` block as defined below.

* `storage_account_url` - The SAS URL to the container.

---

A `connection_string` block exports the following:

* `name` - The name of this Connection String.

* `type` - The type of Database.

* `value` - The Connection String value.

---

A `cors` block exports the following:

* `allowed_origins` - A `allowed_origins` block as defined above.

* `support_credentials` - Whether CORS requests with credentials are allowed.

---

A `custom_action` block exports the following:

* `executable` - The command run when this `auto_heal` action is triggered.

* `parameters` - The parameters passed to the `executable`.

---

A `facebook` block exports the following:

* `app_id` - The App ID of the Facebook app used for login.

* `app_secret` - The App Secret of the Facebook app used for Facebook login.

* `app_secret_setting_name` - The app setting name that contains the `app_secret` value used for Facebook login.

* `oauth_scopes` - A list of OAuth 2.0 scopes that are part of Facebook login authentication.

---

A `file_system` block exports the following:

* `retention_in_days` - The retention period in days.

* `retention_in_mb` - The maximum size in megabytes that log files can use.

---

A `github` block exports the following:

* `client_id` - The ID of the GitHub app used for login.

* `client_secret` - The Client Secret of the GitHub app used for GitHub login.

* `client_secret_setting_name` - The app setting name that contains the `client_secret` value used for GitHub login.

* `oauth_scopes` - A list of OAuth 2.0 scopes in the GitHub login authentication.

---

A `google` block exports the following:

* `client_id` - The OpenID Connect Client ID for the Google web application.

* `client_secret` - The client secret associated with the Google web application.

* `client_secret_setting_name` - The app setting name that contains the `client_secret` value used for Google login.

* `oauth_scopes` - A list of OAuth 2.0 scopes that are part of Google Sign-In authentication.

---

A `http_logs` block exports the following:

* `azure_blob_storage` - A `azure_blob_storage` block as defined above.

* `file_system` - A `file_system` block as defined above.

---

A `identity` block exports the following:

* `identity_ids` - A `identity_ids` block as defined below.

* `principal_id` - The Principal ID Managed Service Identity.

* `tenant_id` - The Tenant ID of the Managed Service Identity.

* `type` - The type of Managed Service Identity.

---

A `logs` block exports the following:

* `application_logs` - A `application_logs` block as defined above.

* `detailed_error_messages` - Is Detailed Error Messaging enabled.

* `failed_request_tracing` - Is Failed Request Tracing enabled.

* `http_logs` - An `http_logs` block as defined above.

---

A `microsoft` block exports the following:

* `client_id` - The OAuth 2.0 client ID used by the app for authentication.

* `client_secret` - The OAuth 2.0 client secret used by the app for authentication.

* `client_secret_setting_name` - The app setting name containing the OAuth 2.0 client secret used by the app for authentication.

* `oauth_scopes` - A list of OAuth 2.0 scopes requested as part of Microsoft Account authentication.

---

A `requests` block exports the following:

* `count` - The number of requests in the specified `interval` to trigger this rule.

* `interval` - The interval in `hh:mm:ss`.

---

A `schedule` block exports the following:

* `frequency_interval` - How often the backup will be executed.

* `frequency_unit` - The unit of time for how often the backup should take place.

* `keep_at_least_one_backup` - Will the service keep at least one backup, regardless of age of backup.

* `last_execution_time` - The time of the last backup attempt.

* `retention_period_days` - After how many days backups should be deleted.

* `start_time` - When the schedule should start in RFC-3339 format.

---

A `site_config` block exports the following:

* `always_on` - Is this Windows Web App is Always On enabled.

* `api_definition_url` - The ID of the APIM configuration for this Windows Web App.

* `api_management_api_id` - The ID of the API Management setting linked to the Windows Web App.

* `app_command_line` - The command line used to launch this app.

* `application_stack` - A `application_stack` block as defined above.

* `auto_heal_enabled` - Are Auto heal rules to be enabled.

* `auto_heal_setting` - A `auto_heal_setting` block as defined above.

* `auto_swap_slot_name` - The Windows Web App Slot Name to automatically swap to when deployment to that slot is successfully completed.

* `container_registry_managed_identity_client_id` - The Client ID of the Managed Service Identity used for connections to the Azure Container Registry.

* `container_registry_use_managed_identity` - Do connections for Azure Container Registry use Managed Identity.

* `cors` - A `cors` block as defined above.

* `default_documents` - The list of Default Documents for the Windows Web App.

* `detailed_error_logging_enabled` - Is Detailed Error Logging enabled.

* `ftps_state` - The State of FTP / FTPS service.

* `health_check_path` - The path to the Health Check endpoint.

* `health_check_eviction_time_in_min` - (Optional) The amount of time in minutes that a node can be unhealthy before being removed from the load balancer. Possible values are between `2` and `10`. Only valid in conjunction with `health_check_path`.

* `http2_enabled` - Is HTTP2.0 enabled.

* `ip_restriction` - A `ip_restriction` block as defined above.

* `load_balancing_mode` - The site Load Balancing Mode.

* `local_mysql_enabled` - Is the Local MySQL enabled.

* `managed_pipeline_mode` - The Managed Pipeline Mode.

* `minimum_tls_version` - The Minimum version of TLS for requests.

* `remote_debugging` - Is Remote Debugging enabled.

* `remote_debugging_version` - The Remote Debugging Version.

* `scm_ip_restriction` - A `scm_ip_restriction` block as defined above.

* `scm_minimum_tls_version` - The Minimum version of TLS for requests to SCM.

* `scm_type` - The Source Control Management Type in use.

* `scm_use_main_ip_restriction` - Is the Windows Web App `ip_restriction` configuration used for the SCM also.

* `use_32_bit_worker` - Does the Windows Web App use a 32-bit worker.

* `virtual_application` - A `virtual_application` block as defined below.

* `vnet_route_all_enabled` - Are all outbound traffic to NAT Gateways, Network Security Groups and User Defined Routes applied?

* `websockets_enabled` - Are Web Sockets enabled?

* `windows_fx_version` - The string representation of the Windows FX Version.

* `worker_count` - The number of Workers for this Windows App Service.

---

A `site_credential` block exports the following:

* `name` - The Site Credentials Username used for publishing.

* `password` - The Site Credentials Password used for publishing.

---

A `slow_request` block exports the following:

* `count` - The number of requests within the interval at which to trigger.

* `interval` - The time interval.

* `path` - The App Path for which this rule applies.

* `time_taken` - The amount of time that qualifies as slow for this rule.

---

A `status_code` block exports the following:

* `count` - The number of occurrences of the defined `status_code` in the specified `interval` on which to trigger this rule.

* `interval` - The time interval in the form `hh:mm:ss`.

* `path` - The path to which this rule status code applies.

* `status_code_range` - The status code or range for this rule.

* `sub_status` - The Request Sub Status of the Status Code.

* `win32_status_code` - The Win32 Status Code of the Request.

---

A `sticky_settings` block exports the following:

* `app_setting_names` - A list of `app_setting` names that the Windows Web App will not swap between Slots when a swap operation is triggered.

* `connection_string_names` - A list of `connection_string` names that the Windows Web App will not swap between Slots when a swap operation is triggered.

---

A `storage_account` block exports the following:

* `access_key` - The Access key for the storage account.

* `account_name` - The Name of the Storage Account.

* `mount_path` - The path at which to mount the Storage Share.

* `name` - The name of this Storage Account.

* `share_name` - The Name of the File Share.

* `type` - The Azure Storage Type.

---

A `trigger` block exports the following:

* `private_memory_kb` - The amount of Private Memory used.

* `requests` - A `requests` block as defined above.

* `slow_request` - A `slow_request` block as defined above.

* `status_code` - A `status_code` block as defined above.

---

A `twitter` block exports the following:

* `consumer_key` - The OAuth 1.0a consumer key of the Twitter application used for sign-in.

* `consumer_secret` - The OAuth 1.0a consumer secret of the Twitter application used for sign-in.

* `consumer_secret_setting_name` - The app setting name that contains the OAuth 1.0a consumer secret of the Twitter application used for sign-in.

---

A `virtual_application` block exports the following:

* `physical_path` - The path on disk to the Virtual Application.

* `preload` - Is this Application Pre-loaded at startup.

* `virtual_directory` - A `virtual_directory` block as defined below.

* `virtual_path` - The Virtual Path of the Virtual Application on the service.

---

A `virtual_directory` block exports the following:

* `physical_path` - The path on disk to the Virtual Directory

* `virtual_path` - The Virtual Path of the Virtual Directory.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 10 minutes) Used when retrieving the Windows Web App.
