## 3.0.0 (Unreleased)

NOTES:

* **Major Version**: Version 3.0 of the Azure Provider is a major version - some behaviours have changed and some deprecated fields/resources have been removed - please refer to [the 3.0 upgrade guide for more information](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/3.0-upgrade-guide).
* When upgrading to v3.0 of the AzureRM Provider, we recommend upgrading to the latest version of Terraform Core ([which can be found here](https://www.terraform.io/downloads)) - the next major release of the AzureRM Provider (v4.0) will require Terraform 1.0 or later.

FEATURES:

* **New Data Source**: `azurerm_healthcare_workspace` [GH-15759]
* **New Data Source**: `azurerm_managed_api` [GH-15797]
* **New Resource**: `azurerm_api_connection` [GH-15797]
* **New Resource**: `azurerm_healthcare_workspace` [GH-15759]
* **New Resource**: `azurerm_stream_analytics_function_javascript_uda` [GH-15831]

ENHANCEMENTS:

* dependencies: upgrading to `v0.26.0` of `github.com/hashicorp/go-azure-helpers` [GH-15889]
* provider: MSAL (and Microsoft Graph) is now used for authentication instead of ADAL (and Azure Active Directory Graph) [GH-12443]
* provider: added a new feature flag within the `key_vault` block for `recover_soft_deleted_certificates`, for configuring whether a soft-deleted `azurerm_key_vault_certificate` should be recovered during creation [GH-10273]
* provider: added a new feature flag within the `key_vault` block for `purge_soft_deleted_certificates_on_destroy`, for configuring whether a deleted `azurerm_key_vault_certificate` should be purged during deletion [GH-10273]
* provider: added a new feature flag within the `key_vault` block for `recover_soft_deleted_keys`, for configuring whether a soft-deleted `azurerm_key_vault_key` should be recovered during creation [GH-10273]
* provider: added a new feature flag within the `key_vault` block for `purge_soft_deleted_keys_on_destroy`, for configuring whether a deleted `azurerm_key_vault_key` should be purged during deletion [GH-10273]
* provider: added a new feature flag within the `key_vault` block for `recover_soft_deleted_secrets`, for configuring whether a soft-deleted `azurerm_key_vault_secret` should be recovered during creation [GH-10273]
* provider: added a new feature flag within the `key_vault` block for `purge_soft_deleted_secrets_on_destroy`, for configuring whether a deleted `azurerm_key_vault_secret` should be purged during deletion [GH-10273]
* provider: added a new feature flag within the `resource_group` block for `prevent_deletion_if_contains_resources`, for configuring whether Terraform should prevent the deletion of a Resource Group which still contains items [GH-13777]
* Resources supporting Availability Zones: Zones are now treated consistently across the Provider and the field within Terraform has been renamed to either `zone` (for a single Zone) or `zones` (where multiple can be defined) - the complete list of resources can be found in the 3.0 Upgrade Guide [GH-14588]
* Resources supporting Managed Identity: Identity blocks are now treated consistently across the Provider - the complete list of resources can be found in the 3.0 Upgrade Guide [GH-15187]
* provider: removing the `network` and `relaxed_locking` feature flags, since this is now enabled by default [GH-15719]
* `azurerm_eventgrid_system_topic_event_subscription` - support for the `delivery_property` property [GH-15559]
* `azurerm_iothub` - add support for `authentication_type` and `identity_id` properties in the `file_upload` block [GH-15874]
* `azurerm_kubernetes_cluster` - the `kube_admin_config` block is now marked as sensitive in addition to all items within it [GH-4105]
* `azurerm_kubernetes_cluster` - add support for `key_vault_secrets_provider` and `open_service_mesh_enabled` in Azure China and Azure Government [GH-15878]
* `azurerm_linux_function_app` - add support for `storage_key_vault_secret_id` [GH-15793]
* `azurerm_linux_function_app` - updating the read timeout to be 5m [GH-15867]
* `azurerm_linux_function_app_slot` - add support for `storage_key_vault_secret_id` [GH-15793]
* `azurerm_linux_function_app_slot` - updating the read timeout to be 5m [GH-15867]
* Data Source: `azurerm_linux_function_app` - add support for `storage_key_vault_secret_id` [GH-15793]
* Data Source: `azurerm_windows_function_app` - add support for `storage_key_vault_secret_id` [GH-15793]
* `azurerm_management_group_policy_assignment` - support for User Assigned Identities [GH-15376]
* `azurerm_mssql_server` - `minimum_tls_version` now defaults to `1.2` [GH-10276]
* `azurerm_mysql_server` - `ssl_minimal_tls_version_enforced` now defaults to `1.2` [GH-10276]
* `azurerm_network_security_rule` - no longer locking on the network security group name [GH-15719]
`azurerm_postgresql_server` - `ssl_minimal_tls_version_enforced` now defaults to `1.2` [GH-10276]
* `azurerm_redis_cache` - `minimum_tls_version` now defaults to `1.2` [GH-10276]
* `azurerm_resource_group` - Terraform now checks during the deletion of a Resource Group if there's any items remaining and will raise an error if so by default (to avoid deleting items unintentionally). This behaviour can be controlled using the `prevent_deletion_if_contains_resources` feature-flag within the `resource_group` block within the `features` block. [GH-13777]
* `azurerm_resource_group_policy_assignment` - support for User Assigned Identities [GH-15376]
* `azurerm_resource_policy_assignment` - support for User Assigned Identities [GH-15376]
* `azurerm_static_site` - `identity` now supports `SystemAssigned, UserAssigned` [GH-15834]
* `azurerm_storage_account` - `min_tls_version` now defaults to `1.2` [GH-10276]
* `azurerm_subscription_policy_assignment` - support for User Assigned Identities [GH-15376]
* `azurerm_windows_function_app` - add support for `storage_key_vault_secret_id` [GH-15793]
* `azurerm_windows_function_app` - updating the read timeout to be 5m [GH-15867]
* `azurerm_windows_function_app_slot` - add support for `storage_key_vault_secret_id` [GH-15793]
* `azurerm_windows_function_app_slot` - updating the read timeout to be 5m [GH-15867]

BUG FIXES:

* `azurerm_application_gateway` - the `backend_address_pool` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the field `fqdns` within the `backend_address_pool` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the field `ip_addresses` within the `backend_address_pool` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the `backend_http_settings` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the `frontend_port` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the field `host_names` within the `frontend_port` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the `http_listener` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the `private_endpoint_connection` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the `private_link_configuration` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the `probe` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the `redirect_configuration` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the `request_routing_rule` block is now a Set rather than a List [GH-6896]
* `azurerm_application_gateway` - the `ssl_certificate` block is now a Set rather than a List [GH-6896]
* `azurerm_cosmosdb_mongo_collection` - the `default_ttl_seconds` can now be set to `-1` [GH-15736]

---

For information on changes between the v2.99.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v1.44.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
