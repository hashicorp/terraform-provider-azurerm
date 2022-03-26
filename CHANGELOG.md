
## 3.0.2 (March 26, 2022)

BUG FIXES:

* `azurerm_cosmosdb_account` - prevent a panic when the API returns an nil list of read or write locations ([#16031](https://github.com/hashicorp/terraform-provider-azurerm/issues/16031))
* `azurerm_cdn_endpoint` - prevent a panic when there is an empty `country_codes` property ([#16066](https://github.com/hashicorp/terraform-provider-azurerm/issues/16066))
* `azurerm_key_vault` - fix the `authorizer was not an auth.CachedAuthorizer ` error ([#16078](https://github.com/hashicorp/terraform-provider-azurerm/issues/16078))
* `azurerm_linux_function_app` - correctly update storage settings when using MSI ([#16046](https://github.com/hashicorp/terraform-provider-azurerm/issues/16046))
* `azurerm_managed_disk` - changing the `zone` property now correctly create a new resource ([#16070](https://github.com/hashicorp/terraform-provider-azurerm/issues/16070))
* `azurerm_resource_group` - wait for eventual consistency when deleting ([#16073](https://github.com/hashicorp/terraform-provider-azurerm/issues/16073))
* `azurerm_windows_function_app` - correctly update storage settings when using MSI ([#16046](https://github.com/hashicorp/terraform-provider-azurerm/issues/16046))

## 3.0.1 (March 24, 2022)

BUG FIXES:

* provider: the `prevent_deletion_if_contains_resources` feature flag within the `resource_group` block now defaults to `true` ([#16021](https://github.com/hashicorp/terraform-provider-azurerm/issues/16021))

## 3.0.0 (March 24, 2022)

NOTES:

* **Major Version**: Version 3.0 of the Azure Provider is a major version - some behaviours have changed and some deprecated fields/resources have been removed - please refer to [the 3.0 upgrade guide for more information](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/3.0-upgrade-guide).
* When upgrading to v3.0 of the AzureRM Provider, we recommend upgrading to the latest version of Terraform Core ([which can be found here](https://www.terraform.io/downloads)) - the next major release of the AzureRM Provider (v4.0) will require Terraform 1.0 or later.

FEATURES:

* **New Data Source**: `azurerm_healthcare_workspace` ([#15759](https://github.com/hashicorp/terraform-provider-azurerm/issues/15759))
* **New Data Source**: `azurerm_key_vault_encrypted_value` ([#15873](https://github.com/hashicorp/terraform-provider-azurerm/issues/15873))
* **New Data Source**: `azurerm_managed_api` ([#15797](https://github.com/hashicorp/terraform-provider-azurerm/issues/15797))
* **New Resource**: `azurerm_api_connection` ([#15797](https://github.com/hashicorp/terraform-provider-azurerm/issues/15797))
* **New Resource**: `azurerm_healthcare_workspace` ([#15759](https://github.com/hashicorp/terraform-provider-azurerm/issues/15759))
* **New Resource**: `azurerm_stream_analytics_function_javascript_uda` ([#15831](https://github.com/hashicorp/terraform-provider-azurerm/issues/15831))
* **New Resource**: `azurerm_security_center_server_vulnerability_assessment_virtual_machine` ([#15747](https://github.com/hashicorp/terraform-provider-azurerm/issues/15747))

ENHANCEMENTS:

* dependencies: updating to `v62.3.0` of `github.com/Azure/azure-sdk-for-go` ([#15927](https://github.com/hashicorp/terraform-provider-azurerm/issues/15927))
* dependencies: updating to `v0.26.0` of `github.com/hashicorp/go-azure-helpers` ([#15889](https://github.com/hashicorp/terraform-provider-azurerm/issues/15889))
* dependencies: updating `appplatform` to API Version `2022-01-01-preview` ([#15597](https://github.com/hashicorp/terraform-provider-azurerm/issues/15597))
* provider: MSAL (and Microsoft Graph) is now used for authentication instead of ADAL (and Azure Active Directory Graph) ([#12443](https://github.com/hashicorp/terraform-provider-azurerm/issues/12443))
* provider: all (non-deprecated) resources now validate the Resource ID during import ([#15989](https://github.com/hashicorp/terraform-provider-azurerm/issues/15989))
* provider: added a new feature flag within the `api_management` block for `recover_soft_deleted`, for configuring whether a soft-deleted `azurerm_api_management` should be recovered during creation ([#15871](https://github.com/hashicorp/terraform-provider-azurerm/issues/15871))
* provider: added a new feature flag within the `key_vault` block for `recover_soft_deleted_certificates`, for configuring whether a soft-deleted `azurerm_key_vault_certificate` should be recovered during creation ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `key_vault` block for `purge_soft_deleted_certificates_on_destroy`, for configuring whether a deleted `azurerm_key_vault_certificate` should be purged during deletion ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `key_vault` block for `recover_soft_deleted_keys`, for configuring whether a soft-deleted `azurerm_key_vault_key` should be recovered during creation ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `key_vault` block for `purge_soft_deleted_keys_on_destroy`, for configuring whether a deleted `azurerm_key_vault_key` should be purged during deletion ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `key_vault` block for `recover_soft_deleted_secrets`, for configuring whether a soft-deleted `azurerm_key_vault_secret` should be recovered during creation ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `key_vault` block for `purge_soft_deleted_secrets_on_destroy`, for configuring whether a deleted `azurerm_key_vault_secret` should be purged during deletion ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `resource_group` block for `prevent_deletion_if_contains_resources`, for configuring whether Terraform should prevent the deletion of a Resource Group which still contains items ([#13777](https://github.com/hashicorp/terraform-provider-azurerm/issues/13777))
* provider: the feature flag `permanently_delete_on_destroy` within the `log_analytics_workspace` block now defaults to `true` ([#15948](https://github.com/hashicorp/terraform-provider-azurerm/issues/15948))
* Resources supporting Availability Zones: Zones are now treated consistently across the Provider and the field within Terraform has been renamed to either `zone` (for a single Zone) or `zones` (where multiple can be defined) - the complete list of resources can be found in the 3.0 Upgrade Guide ([#14588](https://github.com/hashicorp/terraform-provider-azurerm/issues/14588))
* Resources supporting Managed Identity: Identity blocks are now treated consistently across the Provider - the complete list of resources can be found in the 3.0 Upgrade Guide ([#15187](https://github.com/hashicorp/terraform-provider-azurerm/issues/15187))
* provider: removing the `network` and `relaxed_locking` feature flags, since this is now enabled by default ([#15719](https://github.com/hashicorp/terraform-provider-azurerm/issues/15719))
* Data Source: `azurerm_linux_function_app` - support for the `storage_key_vault_secret_id` property ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* Data Source: `azurerm_storage_account_sas` - now exports the `tag` and `filter` attributes ([#15863](https://github.com/hashicorp/terraform-provider-azurerm/issues/15863))
* Data Source: `azurerm_windows_function_app` - support for `storage_key_vault_secret_id` property ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* `azurerm_application_insights` - can now disable Rule and Action Groups that are automatically created ([#15892](https://github.com/hashicorp/terraform-provider-azurerm/issues/15892))
* `azurerm_cdn_endpoint` - the `host_name` property has been renamed to `fqdn` ([#15992](https://github.com/hashicorp/terraform-provider-azurerm/issues/15992))
* `azurerm_eventgrid_system_topic_event_subscription` - support for the `delivery_property` property ([#15559](https://github.com/hashicorp/terraform-provider-azurerm/issues/15559))
* `azurerm_iothub` - add support for the `authentication_type` and `identity_id` properties in the `file_upload` block ([#15874](https://github.com/hashicorp/terraform-provider-azurerm/issues/15874))
* `azurerm_kubernetes_cluster` - the `kube_admin_config` block is now marked as sensitive in addition to all items within it ([#4105](https://github.com/hashicorp/terraform-provider-azurerm/issues/4105))
* `azurerm_kubernetes_cluster` - add support for the `key_vault_secrets_provider` and `open_service_mesh_enabled` property in Azure China and Azure Government ([#15878](https://github.com/hashicorp/terraform-provider-azurerm/issues/15878))
* `azurerm_linux_function_app` - add support for the `storage_key_vault_secret_id` property ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* `azurerm_linux_function_app` - updating the read timeout to be `5m` ([#15867](https://github.com/hashicorp/terraform-provider-azurerm/issues/15867))
* `azurerm_linux_function_app` - support for node version `16` preview ([#15884](https://github.com/hashicorp/terraform-provider-azurerm/issues/15884))
* `azurerm_linux_function_app` - add support for `use_dotnet_isolated_runtime` ([#15969](https://github.com/hashicorp/terraform-provider-azurerm/issues/15969))
* `azurerm_linux_function_app_slot` - add support for `use_dotnet_isolated_runtime` ([#15969](https://github.com/hashicorp/terraform-provider-azurerm/issues/15969))
* `azurerm_linux_function_app_slot` - add support for `storage_key_vault_secret_id` ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* `azurerm_linux_function_app_slot` - updating the read timeout to be 5m ([#15867](https://github.com/hashicorp/terraform-provider-azurerm/issues/15867))
* `azurerm_linux_virtual_machine` - support for the `termination_notification` property ([#14933](https://github.com/hashicorp/terraform-provider-azurerm/issues/14933))
* `azurerm_linux_virtual_machine ` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_linux_virtual_machine_scale_set` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_linux_web_app` - support for PHP version 8.0 ([#15933](https://github.com/hashicorp/terraform-provider-azurerm/issues/15933))
* `azurerm_loadbalancer` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_managed_disk` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_management_group_policy_assignment` - support for User Assigned Identities ([#15376](https://github.com/hashicorp/terraform-provider-azurerm/issues/15376))
* `azurerm_mssql_server` - the `minimum_tls_version` property now defaults to `1.2` ([#10276](https://github.com/hashicorp/terraform-provider-azurerm/issues/10276))
* `azurerm_mysql_server` - the `ssl_minimal_tls_version_enforced` property now defaults to `1.2` ([#10276](https://github.com/hashicorp/terraform-provider-azurerm/issues/10276))
* `azurerm_network_interface` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_network_security_rule` - no longer locks on the network security group name ([#15719](https://github.com/hashicorp/terraform-provider-azurerm/issues/15719))
* `azurerm_postgresql_server` - the `ssl_minimal_tls_version_enforced` property now defaults to `1.2` ([#10276](https://github.com/hashicorp/terraform-provider-azurerm/issues/10276))
* `azurerm_public_ip` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_redis_cache` - the `minimum_tls_version` property  now defaults to `1.2` ([#10276](https://github.com/hashicorp/terraform-provider-azurerm/issues/10276))
* `azurerm_resource_group` - Terraform now checks during the deletion of a Resource Group if there's any items remaining and will raise an error if so by default (to avoid deleting items unintentionally). This behaviour can be controlled using the `prevent_deletion_if_contains_resources` feature-flag within the `resource_group` block within the `features` block. ([#13777](https://github.com/hashicorp/terraform-provider-azurerm/issues/13777))
* `azurerm_resource_group_policy_assignment` - support for User Assigned Identities ([#15376](https://github.com/hashicorp/terraform-provider-azurerm/issues/15376))
* `azurerm_resource_policy_assignment` - support for User Assigned Identities ([#15376](https://github.com/hashicorp/terraform-provider-azurerm/issues/15376))
* `azurerm_sentinel_alert_rule_scheduled` - support for `alert_details_override` ([#15901](https://github.com/hashicorp/terraform-provider-azurerm/issues/15901))
* `azurerm_sentinel_alert_rule_scheduled` - support for `entity_mapping` [[#15901](https://github.com/hashicorp/terraform-provider-azurerm/issues/15901)] 
* `azurerm_sentinel_alert_rule_scheduled` - support for `custom_details` ([#15901](https://github.com/hashicorp/terraform-provider-azurerm/issues/15901))
* `azurerm_sentinel_alert_rule_scheduled` - support for `group_by_alert_details` ([#15901](https://github.com/hashicorp/terraform-provider-azurerm/issues/15901))
* `azurerm_sentinel_alert_rule_scheduled` - support for `group_by_custom_details` ([#15901](https://github.com/hashicorp/terraform-provider-azurerm/issues/15901))
* `azurerm_site_recovery_replicated_vm` - support for the `target_availability_zone` property ([#15617](https://github.com/hashicorp/terraform-provider-azurerm/issues/15617))
* `azurerm_shared_image` - support for the `support_accelerated_network` property ([#15562](https://github.com/hashicorp/terraform-provider-azurerm/issues/15562))
* `azurerm_static_site` - the `identity` property now supports `SystemAssigned` and `UserAssigned` ([#15834](https://github.com/hashicorp/terraform-provider-azurerm/issues/15834))
* `azurerm_storage_account` - the `allow_blob_public_access` property has been renamed to `allow_nested_items_to_be_public` to better represent what is being enabled ([#12689](https://github.com/hashicorp/terraform-provider-azurerm/issues/12689))
* `azurerm_storage_account` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_storage_account` - `ZRS` is no longer supported when using `StorageV1` ([#16004](https://github.com/hashicorp/terraform-provider-azurerm/issues/16004))
* `azurerm_storage_account` - the `min_tls_version` property now defaults to `1.2` ([#10276](https://github.com/hashicorp/terraform-provider-azurerm/issues/10276))
* `azurerm_storage_share` - `quota` is now required ([#15982](https://github.com/hashicorp/terraform-provider-azurerm/issues/15982))
* `azurerm_subscription_policy_assignment` - support for User Assigned Identities ([#15376](https://github.com/hashicorp/terraform-provider-azurerm/issues/15376))
* `azurerm_virtual_network` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_virtual_network_gateway` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_virtual_hub` - support for the `virtual_router_asn` and `virtual_router_ips` properties ([#15741](https://github.com/hashicorp/terraform-provider-azurerm/issues/15741))
* `azurerm_windows_function_app` - add support for `storage_key_vault_secret_id` property ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* `azurerm_windows_function_app` - updating the read timeout to be `5m` ([#15867](https://github.com/hashicorp/terraform-provider-azurerm/issues/15867))
* `azurerm_windows_function_app` node version validation string can not be prefixed with `~` ([#15884](https://github.com/hashicorp/terraform-provider-azurerm/issues/15884))
* `azurerm_windows_function_app` support for node version `16` preview support ([#15884](https://github.com/hashicorp/terraform-provider-azurerm/issues/15884))
* `azurerm_windows_function_app` - add support for `use_dotnet_isolated_runtime` ([#15969](https://github.com/hashicorp/terraform-provider-azurerm/issues/15969))
* `azurerm_windows_function_app_slot` - add support for `use_dotnet_isolated_runtime` ([#15969](https://github.com/hashicorp/terraform-provider-azurerm/issues/15969))
* `azurerm_windows_function_app_slot` - add support for the `storage_key_vault_secret_id` property ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* `azurerm_windows_function_app_slot` - updating the read timeout to be 5m ([#15867](https://github.com/hashicorp/terraform-provider-azurerm/issues/15867))
* `azurerm_windows_virtual_machine` - support for the `termination_notification` property ([#14933](https://github.com/hashicorp/terraform-provider-azurerm/issues/14933))
* `azurerm_windows_virtual_machine` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_windows_virtual_machine_scale_set` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))

BUG FIXES:

* provider: the `recover_soft_deleted_key_vaults` feature flag within the `key_vault` block now defaults to `true` ([#15984](https://github.com/hashicorp/terraform-provider-azurerm/issues/15984))
* provider: the `purge_soft_delete_on_destroy ` feature flag within the `key_vault` block now defaults to `true` [[#15984](https://github.com/hashicorp/terraform-provider-azurerm/issues/15984)] 
* `azurerm_app_configuration_feature` - detecting that the key is gone when the App Configuration has been deleted ([#15973](https://github.com/hashicorp/terraform-provider-azurerm/issues/15973))
* `azurerm_app_configuration_key` - detecting that the key is gone when the App Configuration has been deleted ([#15973](https://github.com/hashicorp/terraform-provider-azurerm/issues/15973))
* `azurerm_application_gateway` - the `backend_address_pool` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the field `fqdns` within the `backend_address_pool` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the field `ip_addresses` within the `backend_address_pool` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `backend_http_settings` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `frontend_port` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the field `host_names` within the `frontend_port` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `http_listener` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `private_endpoint_connection` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `private_link_configuration` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `probe` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `redirect_configuration` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `request_routing_rule` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `ssl_certificate` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_container_registry` - validate the `georepliactions` property does not include the location of the Container Registry ([#15847](https://github.com/hashicorp/terraform-provider-azurerm/issues/15847))
* `azurerm_cosmosdb_mongo_collection` - the `default_ttl_seconds` property can now be set to `-1` ([#15736](https://github.com/hashicorp/terraform-provider-azurerm/issues/15736))
* `azurerm_eventhub` - prevent panic when the `capture_description` block is removed ([#15930](https://github.com/hashicorp/terraform-provider-azurerm/issues/15930))
* `azurerm_key_vault_access_policy` - validating the Resource ID during import ([#15989](https://github.com/hashicorp/terraform-provider-azurerm/issues/15989))
* `azurerm_linux_function_app` - fixed update handling of `app_settings` for `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#15907](https://github.com/hashicorp/terraform-provider-azurerm/issues/15907))
* `azurerm_linux_function_app_slot` - fixed update handling of `app_settings` for `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#15907](https://github.com/hashicorp/terraform-provider-azurerm/issues/15907))
* `azurerm_local_network_gateway` - fix for `address_space` cannot be updated ([#15159](https://github.com/hashicorp/terraform-provider-azurerm/issues/15159))
* `azurerm_log_analytics_cluster_customer_managed_key` - detecting when the Customer Managed Key has been removed ([#15973](https://github.com/hashicorp/terraform-provider-azurerm/issues/15973))
* `azurerm_mssql_database_vulnerability_assessment_rule_baseline` - prevent the resource from being replaced every apply ([#14759](https://github.com/hashicorp/terraform-provider-azurerm/issues/14759))
* `azurerm_security_center_auto_provisioning ` - validating the Resource ID during import [[#15989](https://github.com/hashicorp/terraform-provider-azurerm/issues/15989)] 
* `azurerm_security_center_setting` - changing the `setting_name` property now forces a new resource ([#15983](https://github.com/hashicorp/terraform-provider-azurerm/issues/15983))
* `azurerm_synapse_workspace` - fixing a bug where workspaces created from a Dedicated SQL Pool / SQL Data Warehouse couldn't be retrieved ([#15829](https://github.com/hashicorp/terraform-provider-azurerm/issues/15829))
* `azurerm_synapse_workspace_key` - keys can now be correctly rotated ([#15897](https://github.com/hashicorp/terraform-provider-azurerm/issues/15897))
* `azurerm_windows_function_app` - fixed update handling of `app_settings` for `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#15907](https://github.com/hashicorp/terraform-provider-azurerm/issues/15907))
* `azurerm_windows_function_app_slot` - fixed update handling of `app_settings` for `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#15907](https://github.com/hashicorp/terraform-provider-azurerm/issues/15907))

---

For information on changes between the v2.99.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v1.44.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
