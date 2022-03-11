## 2.99.0 (March 11, 2022)

NOTES

* **Preparation for 3.0**: We intend for v2.99.0 to be the last release in the 2.x line - we’ll be turning our focus to 3.0 with the next release. We recommend [consulting the list of changes coming in 3.0](https://registry.terraform.io/providers/hashicorp/azurerm/2.99.0/docs/guides/3.0-upgrade-guide) to be aware and [trialling the Beta available in the latest 2.x releases](https://registry.terraform.io/providers/hashicorp/azurerm/2.99.0/docs/guides/3.0-beta) if you’re interested.

FEATURES:

* New Beta Resource: `azurerm_function_app_function` ([#15605](https://github.com/hashicorp/terraform-provider-azurerm/issues/15605))
* New Beta Resource: `azurerm_function_app_hybrid_connection` ([#15702](https://github.com/hashicorp/terraform-provider-azurerm/issues/15702))
* New Beta Resource: `azurerm_web_app_hybrid_connection` ([#15702](https://github.com/hashicorp/terraform-provider-azurerm/issues/15702))
* New Resource: `azurerm_cosmosdb_sql_role_assignment` ([#15038](https://github.com/hashicorp/terraform-provider-azurerm/issues/15038))
* New Resource: `azurerm_cosmosdb_sql_role_definition` ([#15035](https://github.com/hashicorp/terraform-provider-azurerm/issues/15035))

ENHANCEMENTS:

* dependencies: updating to `v62.1.0` of `github.com/Azure/azure-sdk-for-go` ([#15716](https://github.com/hashicorp/terraform-provider-azurerm/issues/15716))
* dependencies: updating `compute` to `2021-11-01` ([#15099](https://github.com/hashicorp/terraform-provider-azurerm/issues/15099))
* dependencies: updating `kubernetescluster` to `2022-01-02-preview` ([#15648](https://github.com/hashicorp/terraform-provider-azurerm/issues/15648))
* dependencies: updating `sentinel` to `2021-09-01-preview` ([#14983](https://github.com/hashicorp/terraform-provider-azurerm/issues/14983))
* Data Source: `azurerm_kubernetes_cluster` - deprecated the `addon_profile` block in favour of `aci_connector_linux`, `azure_policy_enabled`, `http_application_routing_enabled`, `ingress_application_gateway`, `key_vault_secrets_provider`, `oms_agent` and `open_service_mesh_enabled` properties ([#15584](https://github.com/hashicorp/terraform-provider-azurerm/issues/15584))
* Data Source: `azurerm_kubernetes_cluster` - deprecated the `role_based_access_control` block in favour of `azure_active_directory_role_based_access_control` and `role_based_access_control_enabled` properties ([#15584](https://github.com/hashicorp/terraform-provider-azurerm/issues/15584))
* Data Source: `azurerm_servicebus_namespace_authorization_rule` - support for the `namespace_id` property ([#15671](https://github.com/hashicorp/terraform-provider-azurerm/issues/15671))
* Data Source: `azurerm_servicebus_namespace_disaster_recovery_config` - support for the `namespace_id` property ([#15671](https://github.com/hashicorp/terraform-provider-azurerm/issues/15671))
* Data Source: `azurerm_servicebus_queue` - support for the `namespace_id` property ([#15671](https://github.com/hashicorp/terraform-provider-azurerm/issues/15671))
* Data Source: `azurerm_servicebus_queue_authorization_rule` - support for the `queue_id` property ([#15671](https://github.com/hashicorp/terraform-provider-azurerm/issues/15671))
* Data Source: `azurerm_servicebus_subscription` - support for the `topic_id` property ([#15671](https://github.com/hashicorp/terraform-provider-azurerm/issues/15671))
* Data Source: `azurerm_servicebus_topic` - support for the `namespace_id` property ([#15671](https://github.com/hashicorp/terraform-provider-azurerm/issues/15671))
* Data Source: `azurerm_servicebus_topic_authorization_rule` - support for the `topic_id` property ([#15671](https://github.com/hashicorp/terraform-provider-azurerm/issues/15671))
* Data Source: `azurerm_virtual_network` - support for the `tags` property ([#14882](https://github.com/hashicorp/terraform-provider-azurerm/issues/14882))
* `azurerm_batch_account` - support for customer managed keys ([#14749](https://github.com/hashicorp/terraform-provider-azurerm/issues/14749))
* `azurerm_container_registry` support for the `export_policy_enabled` property ([#15036](https://github.com/hashicorp/terraform-provider-azurerm/issues/15036))
* `azurerm_kubernetes_cluster` - deprecate the `role_based_access_control` block in favour of `role_based_access_control_enabled` and `azure_active_directory_role_based_access_control` ([#15546](https://github.com/hashicorp/terraform-provider-azurerm/issues/15546))
* `azurerm_iothub` - deprecate the `ip_filter_rule` property in favour of the `network_rule_set` property ([#15590](https://github.com/hashicorp/terraform-provider-azurerm/issues/15590))
* `azurerm_lb_nat_rule` - the `frontend_port` and `backend_port` properties now support `0` ([#15694](https://github.com/hashicorp/terraform-provider-azurerm/issues/15694))
* `azurerm_machine_learning_compute_instance` - updating the validation on the `name` property ([#14839](https://github.com/hashicorp/terraform-provider-azurerm/issues/14839))
* `azurerm_mssql_database_extended_auditing_policy` - support for the `enabled` property ([#15624](https://github.com/hashicorp/terraform-provider-azurerm/issues/15624))
* `azurerm_mssql_server_extended_auditing_policy` - support for the `enabled` property ([#15624](https://github.com/hashicorp/terraform-provider-azurerm/issues/15624))
* `azurerm_management_group_policy_assignment` - the `parameters` property can now be updated ([#15623](https://github.com/hashicorp/terraform-provider-azurerm/issues/15623))
* `azurerm_mssql_server` - the `administrator_login` and `administrator_login_password` properties are now optional when Azure AD authentication is enforced ([#15771](https://github.com/hashicorp/terraform-provider-azurerm/issues/15771))
* `azurerm_resource_policy_assignment`  - the `parameters` property can now be updated ([#15623](https://github.com/hashicorp/terraform-provider-azurerm/issues/15623))
* `azurerm_resource_group_policy_assignment` - the `parameters` property can now be updated ([#15623](https://github.com/hashicorp/terraform-provider-azurerm/issues/15623))
* `azurerm_recovery_service_vault` - support for the `cross_region_restore_enabled` property ([#15757](https://github.com/hashicorp/terraform-provider-azurerm/issues/15757))
* `azurerm_subscription_policy_assignment` - the `parameters` property can now be updated ([#15623](https://github.com/hashicorp/terraform-provider-azurerm/issues/15623))
* `azurerm_storage_object_replication` - support for replicating containers across subscriptions ([#15603](https://github.com/hashicorp/terraform-provider-azurerm/issues/15603))

BUG FIXES:

* `azurerm_backup_protected_vm` - the `source_vm_id` property is now case insensitive ([#15656](https://github.com/hashicorp/terraform-provider-azurerm/issues/15656))
* `azurerm_batch_job` - will not longer fail during creation if multiple `common_environment_properties` are set ([#15686](https://github.com/hashicorp/terraform-provider-azurerm/issues/15686))
* `azurerm_container_group` - correctly parse empty or omitted `dns_config.options` and `dns_config.search_domains` properties ([#15618](https://github.com/hashicorp/terraform-provider-azurerm/issues/15618))
* `azurerm_key_vault_key` - correctly set the vault id on import ([#15670](https://github.com/hashicorp/terraform-provider-azurerm/issues/15670))
* `azurerm_monitor_diagnostic_setting` - will now correctly parse the `eventhub_authorization_rule_id` property ([#15582](https://github.com/hashicorp/terraform-provider-azurerm/issues/15582))
* `azurerm_mssql_managed_instance_active_directory_administrator` - prevent a perpetual diff with the instance ID ([#15725](https://github.com/hashicorp/terraform-provider-azurerm/issues/15725))
* `azurerm_orchestrated_virtual_machine_scale_set` - prevent a crash when the 3.0 beta was enabled ([#15637](https://github.com/hashicorp/terraform-provider-azurerm/issues/15637))
* `azurerm_storage_data_lake_gen2_filesystem` - support configuring the `group` and `owner` properties ([#15598](https://github.com/hashicorp/terraform-provider-azurerm/issues/15598))
* `azurerm_virtual_network_gateway` - prevent a panic with `bgp_settings.0.peering_address` ([#15689](https://github.com/hashicorp/terraform-provider-azurerm/issues/15689))

## 2.98.0 (February 25, 2022)

FEATURES:

* New Beta Resource: `azurerm_function_app_active_slot` ([#15246](https://github.com/hashicorp/terraform-provider-azurerm/issues/15246))
* New Beta Resource: `azurerm_web_app_active_slot` ([#15246](https://github.com/hashicorp/terraform-provider-azurerm/issues/15246))

ENHANCEMENTS:

* dependencies: upgrading to `v0.18.0` of `github.com/tombuildsstuff/giovanni` ([#15507](https://github.com/hashicorp/terraform-provider-azurerm/issues/15507))
* `azurerm_linux_function_app` - adds `key_vault_reference_identity_id` support ([#15553](https://github.com/hashicorp/terraform-provider-azurerm/issues/15553))
* `azurerm_linux_function_app_slot` - adds `key_vault_reference_identity_id` support ([#15553](https://github.com/hashicorp/terraform-provider-azurerm/issues/15553))
* `azurerm_windows_function_app` - adds `key_vault_reference_identity_id` support ([#15553](https://github.com/hashicorp/terraform-provider-azurerm/issues/15553))
* `azurerm_windows_function_app_slot` - adds `key_vault_reference_identity_id` support ([#15553](https://github.com/hashicorp/terraform-provider-azurerm/issues/15553))

BUG FIXES:

* `azurerm_cosmosdb_mongo_collection` - can now set the `autoscale_settings` property without setting a `shard_key` when creating a cosmos DB mongo collection ([#15529](https://github.com/hashicorp/terraform-provider-azurerm/issues/15529))
* `azurerm_firewall_policy` - will not wait for resource to finish provisioning after creation ([#15561](https://github.com/hashicorp/terraform-provider-azurerm/issues/15561))

## 2.97.0 (February 18, 2022)

UPGRADE NOTES:

* **3.0 Beta:** This release includes a new feature-flag to opt-into the 3.0 Beta - which (when enabled) introduces a number of new data sources/resources, behavioural changes, field renames and removes some older deprecated resources. The 3.0 Beta is still a work-in-progress at this time and as such the changes listed [in the 3.0 Upgrade Guide](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/3.0-upgrade-guide) may change, however we're interested to hear your feedback and [instructions on how to opt-into the 3.0 Beta can be found here](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/3.0-beta).

FEATURES:

* **New Data Source:** `azurerm_extended_locations` ([#15181](https://github.com/hashicorp/terraform-provider-azurerm/issues/15181))
* **New Data Source:** `azurerm_mssql_managed_instance` ([#15203](https://github.com/hashicorp/terraform-provider-azurerm/issues/15203))
* **New Resource:** `azurerm_iothub_certificate` ([#15461](https://github.com/hashicorp/terraform-provider-azurerm/issues/15461))
* **New Resource:** `azurerm_mssql_outbound_firewall_rule` ([#14795](https://github.com/hashicorp/terraform-provider-azurerm/issues/14795))
* **New Resource:** `azurerm_mssql_managed_database` ([#15203](https://github.com/hashicorp/terraform-provider-azurerm/issues/15203))
* **New Resource:** `azurerm_mssql_managed_instance` ([#15203](https://github.com/hashicorp/terraform-provider-azurerm/issues/15203))
* **New Resource:** `azurerm_mssql_managed_instance_active_directory_administrator` ([#15203](https://github.com/hashicorp/terraform-provider-azurerm/issues/15203))
* **New Resource:** `azurerm_mssql_managed_instance_failover_group` ([#15203](https://github.com/hashicorp/terraform-provider-azurerm/issues/15203))
* **New Resource:** `azurerm_spring_cloud_storage` ([#15375](https://github.com/hashicorp/terraform-provider-azurerm/issues/15375))

ENHANCEMENTS:

* dependencies: upgrading to `v0.24.1` of `github.com/hashicorp/go-azure-helpers` ([#15430](https://github.com/hashicorp/terraform-provider-azurerm/issues/15430))
* `azurerm_automation_account` - add support for the `public_network_access_enabled` property ([#15429](https://github.com/hashicorp/terraform-provider-azurerm/issues/15429))
* `azurerm_kubernetes_cluster` - deprecate the `addon_profile` block, moving all properties to the top level as well as removing the `enabled` field for all add-ons ([#15108](https://github.com/hashicorp/terraform-provider-azurerm/issues/15108))
* `azurerm_kusto_cluster` - supports for the `public_network_access_enabled` property ([#15428](https://github.com/hashicorp/terraform-provider-azurerm/issues/15428))
* `azurerm_machine_learning_workspace` - support for both `SystemAssigned, UserAssigned` and `UserAssigned` Identities ([#14181](https://github.com/hashicorp/terraform-provider-azurerm/issues/14181))
* `azurerm_machine_learning_workspace` - support for encryption using a User Assigned Identity ([#14181](https://github.com/hashicorp/terraform-provider-azurerm/issues/14181))
* `azurerm_monitor_activity_log_alert` support for the `resource_health` block ([#14917](https://github.com/hashicorp/terraform-provider-azurerm/issues/14917))
* `azurerm_iothub_dps` - support for the `ip_filter_rule` block and the `public_network_access_enabled` property ([#15343](https://github.com/hashicorp/terraform-provider-azurerm/issues/15343))
* `azurerm_spring_cloud_app` - support for the `custom_persistent_disk` block ([#15400](https://github.com/hashicorp/terraform-provider-azurerm/issues/15400))
* `azurerm_servicebus_namespace` - support for the `identity` block ([#15371](https://github.com/hashicorp/terraform-provider-azurerm/issues/15371))
* `azurerm_storage_account` - add support for creating a customer managed key upon creation of a storage account ([#15082](https://github.com/hashicorp/terraform-provider-azurerm/issues/15082))
* `azurerm_storage_management_policy` - add support for `tier_to_cool_after_days_since_last_access_time_greater_than`, `tier_to_archive_after_days_since_last_access_time_greater_than,` and `delete_after_days_since_last_access_time_greater_than` ([#15423](https://github.com/hashicorp/terraform-provider-azurerm/issues/15423))
* `azurerm_web_pubsub` - support for the `identity` block ([#15288](https://github.com/hashicorp/terraform-provider-azurerm/issues/15288))

BUG FIXES:

* `azurerm_application_gateway` - fixing a regression where the `identity` block wasn't set into the state ([#15412](https://github.com/hashicorp/terraform-provider-azurerm/issues/15412))
* `azurerm_automation_account` - fixing a crash where the `keys` weren't returned from the API ([#15482](https://github.com/hashicorp/terraform-provider-azurerm/issues/15482))
* `azurerm_kusto_cluster` - ranaming the properties `enable_auto_stop` to `auto_stop_enabled`, `enable_disk_encryption` to `disk_encryption_enabled`, `enable_streaming_ingest` to `streaming_ingestion_enabled`, and `enable_purge` to `purge_enabled` with the orginal properties being deprecated ([#15368](https://github.com/hashicorp/terraform-provider-azurerm/issues/15368))
* `azurerm_log_analytics_linked_storage_account` - correct casing for `data_source_type` when using `ingestion` ([#15451](https://github.com/hashicorp/terraform-provider-azurerm/issues/15451))
* `azurerm_logic_app_integration_account_map` - set `content_type` to `text/plain` when `map_type` is `Liquid` ([#15370](https://github.com/hashicorp/terraform-provider-azurerm/issues/15370))
* `azurerm_stream_analytics_cluster` - fix an issue where the `tags` were not being set in the state ([#15380](https://github.com/hashicorp/terraform-provider-azurerm/issues/15380))
* `azurerm_virtual_desktop_host_pool` - the `registration_info` info block is deprecated in favour of the `azurerm_virtual_desktop_host_pool_registration_info` resource due to changes in the API ([#14953](https://github.com/hashicorp/terraform-provider-azurerm/issues/14953))
* `azurerm_virtual_machine_data_disk_attachment` - fixing a panic when an incorrect `disk_id` is provided ([#15470](https://github.com/hashicorp/terraform-provider-azurerm/issues/15470))
* `azurerm_web_application_firewall_policy` - `disabled_rules` is now optional ([#15386](https://github.com/hashicorp/terraform-provider-azurerm/issues/15386))

## 2.96.0 (February 11, 2022)

FEATURES: 

* **New Data Source:** `azurerm_portal_dashboard` ([#15326](https://github.com/hashicorp/terraform-provider-azurerm/issues/15326))
* **New Data Source:** `azurerm_site_recovery_fabric` ([#15349](https://github.com/hashicorp/terraform-provider-azurerm/issues/15349))
* **New Data Source:** `azurerm_site_recovery_protection_container` ([#15349](https://github.com/hashicorp/terraform-provider-azurerm/issues/15349))
* **New Data Source:** `azurerm_site_recovery_replication_policy` ([#15349](https://github.com/hashicorp/terraform-provider-azurerm/issues/15349))
* **New Resource:** `azurerm_disk_pool_iscsi_target_lun` ([#15329](https://github.com/hashicorp/terraform-provider-azurerm/issues/15329))
* **New Resource:** `azurerm_sentinel_watchlist_item` ([#14366](https://github.com/hashicorp/terraform-provider-azurerm/issues/14366))
* **New Resource:** `azurerm_stream_analytics_output_function` ([#15162](https://github.com/hashicorp/terraform-provider-azurerm/issues/15162))
* **New Resource:** `azurerm_web_pubsub_network_acl` ([#14827](https://github.com/hashicorp/terraform-provider-azurerm/issues/14827))
* **New Beta Resource:** `azurerm_app_service_source_control_slot` ([#15301](https://github.com/hashicorp/terraform-provider-azurerm/issues/15301))

ENHANCEMENTS: 

* dependencies: updating to `v0.23.1` of `github.com/hashicorp/go-azure-helpers` ([#15314](https://github.com/hashicorp/terraform-provider-azurerm/issues/15314))
* `azurerm_application_gateway` - the `type` property within the `identity` block is now required when an `identity` block is specified ([#15337](https://github.com/hashicorp/terraform-provider-azurerm/issues/15337))
* `azurerm_application_insights` - support for the `force_customer_storage_for_profiler` property ([#15254](https://github.com/hashicorp/terraform-provider-azurerm/issues/15254))
* `azurerm_automation_account` - support for managed identities ([#15072](https://github.com/hashicorp/terraform-provider-azurerm/issues/15072))
* `azurerm_data_factory` - refactoring the `identity` block to be consistant across resources ([#15344](https://github.com/hashicorp/terraform-provider-azurerm/issues/15344))
* `azurerm_kusto_cluster` - support for the `enable_auto_stop` ([#15332](https://github.com/hashicorp/terraform-provider-azurerm/issues/15332))
* `azurerm_linux_virtual_machine` - support the `StandardSSD_ZRS` and `Premium_ZRS` values for the `storage_account_type` property ([#15360](https://github.com/hashicorp/terraform-provider-azurerm/issues/15360))
* `azurerm_linux_virtual_machine` - full support for Automatic VM Guest Patching ([#14906](https://github.com/hashicorp/terraform-provider-azurerm/issues/14906))
* `azurerm_network_watcher_flow_log` - the `name` property can now be set for new resources ([#15016](https://github.com/hashicorp/terraform-provider-azurerm/issues/15016))
* `azurerm_orchestrated_virtual_machine_scale_set` - full support for Automatic VM Guest Patching and Hotpatching ([#14935](https://github.com/hashicorp/terraform-provider-azurerm/issues/14935))
* `azurerm_windows_virtual_machine` - support the `StandardSSD_ZRS` and `Premium_ZRS` values for the `storage_account_type` property ([#15360](https://github.com/hashicorp/terraform-provider-azurerm/issues/15360))
* `azurerm_windows_virtual_machine` - full support for Automatic VM Guest Patching and Hotpaching ([#14796](https://github.com/hashicorp/terraform-provider-azurerm/issues/14796))

BUG FIXES:

* `azurerm_application_insights_api_key` - prevent panic by checking for the id of an existing API Key ([#15297](https://github.com/hashicorp/terraform-provider-azurerm/issues/15297))
* `azurerm_app_service_active_slot` - fix regression in ID set in creation of new resource ([#15291](https://github.com/hashicorp/terraform-provider-azurerm/issues/15291))
* `azurerm_firewall` - working around an Azure API issue when deleting the Firewall ([#15330](https://github.com/hashicorp/terraform-provider-azurerm/issues/15330))
* `azurerm_kubernetes_cluster` - unsetting `outbound_ip_prefix_ids` or `outbound_ip_address_ids` with an empty slice will default the `load_balancer_profile` to a managed outbound IP ([#15338](https://github.com/hashicorp/terraform-provider-azurerm/issues/15338))
* `azurerm_orchestrated_virtual_machine_scale_set` - fixing a crash when the `computer_name_prefix` wasn't specified ([#15312](https://github.com/hashicorp/terraform-provider-azurerm/issues/15312))
* `azurerm_recovery_services_vault` - fixing an issue where the subscription couldn't be found when running in Azure Government ([#15316](https://github.com/hashicorp/terraform-provider-azurerm/issues/15316))

## 2.95.0 (February 04, 2022)

FEATURES: 

* **New Data Source:** `azurerm_container_group` ([#14946](https://github.com/hashicorp/terraform-provider-azurerm/issues/14946))
* **New Data Source:** `azurerm_logic_app_standard` ([#15199](https://github.com/hashicorp/terraform-provider-azurerm/issues/15199))
* **New Resource:** `azurerm_disk_pool_iscsi_target` ([#14975](https://github.com/hashicorp/terraform-provider-azurerm/issues/14975))
* **New Beta Resource:** `azurerm_linux_function_app_slot` ([#14940](https://github.com/hashicorp/terraform-provider-azurerm/issues/14940))
* **New Resource:** `azurerm_traffic_manager_azure_endpoint` ([#15178](https://github.com/hashicorp/terraform-provider-azurerm/issues/15178))
* **New Resource:** `azurerm_traffic_manager_external_endpoint` ([#15178](https://github.com/hashicorp/terraform-provider-azurerm/issues/15178))
* **New Resource:** `azurerm_traffic_manager_nested_endpoint` ([#15178](https://github.com/hashicorp/terraform-provider-azurerm/issues/15178))
* **New Beta Resource:** `azurerm_windows_function_app_slot` ([#14940](https://github.com/hashicorp/terraform-provider-azurerm/issues/14940))
* **New Beta Resource:** `azurerm_windows_web_app_slot` ([#14613](https://github.com/hashicorp/terraform-provider-azurerm/issues/14613))

ENHANCEMENTS:

* dependencies: upgrading to `v0.22.0` of `github.com/hashicorp/go-azure-helpers` ([#15207](https://github.com/hashicorp/terraform-provider-azurerm/issues/15207))
* dependencies: updating `backup` to API Version `2021-07-01` ([#14980](https://github.com/hashicorp/terraform-provider-azurerm/issues/14980))
* `azurerm_storage_account` - the `identity` block is no longer computed ([#15207](https://github.com/hashicorp/terraform-provider-azurerm/issues/15207))
* `azurerm_linux_virtual_machine` - support for the `dedicated_host_group_id` property ([#14936](https://github.com/hashicorp/terraform-provider-azurerm/issues/14936))
* `azurerm_recovery_services_vault` - support Zone Redundant storage ([#14980](https://github.com/hashicorp/terraform-provider-azurerm/issues/14980))
* `azurerm_web_pubsub_hub` - the `managed_identity_id` property within the `auth` block now accepts UUIDs ([#15183](https://github.com/hashicorp/terraform-provider-azurerm/issues/15183))
* `azurerm_windows_virtual_machine` - support for the `dedicated_host_group_id` property ([#14936](https://github.com/hashicorp/terraform-provider-azurerm/issues/14936))

BUG FIXES:

* `azurerm_container_group` - fixing parallel provisioning failures with the same `network_profile_id` ([#15098](https://github.com/hashicorp/terraform-provider-azurerm/issues/15098))
* `azurerm_frontdoor` - fixing the validation for `resource_group_name` ([#15174](https://github.com/hashicorp/terraform-provider-azurerm/issues/15174))
* `azurerm_kubernetes_cluster` - prevent panic when updating `sku_tier` ([#15229](https://github.com/hashicorp/terraform-provider-azurerm/issues/15229))
* `azurerm_hdinsight_interactive_query_cluster` - support for the `storage_resource_id` property to fix missing storage account errors ([#15039](https://github.com/hashicorp/terraform-provider-azurerm/issues/15039))
* `azurerm_hdinsight_hadoop_cluster` - support for the `storage_resource_id` property to fix missing storage account errors ([#15039](https://github.com/hashicorp/terraform-provider-azurerm/issues/15039))
* `azurerm_hdinsight_spark_cluster` - support for the `storage_resource_id` property to fix missing storage account errors ([#15039](https://github.com/hashicorp/terraform-provider-azurerm/issues/15039))
* `azurerm_hdinsight_hbase_cluster` - support for the `storage_resource_id` property to fix missing storage account errors ([#15039](https://github.com/hashicorp/terraform-provider-azurerm/issues/15039))
* `azurerm_log_analytics_datasource_windows_event` - adding a state migration to fix `ID was missing the dataSources element` ([#15194](https://github.com/hashicorp/terraform-provider-azurerm/issues/15194))
* `azurerm_policy_definition` - fix the deprecation of `management_group_name` in favour of `management_group_id` ([#15209](https://github.com/hashicorp/terraform-provider-azurerm/issues/15209))
* `azurerm_policy_set_definition` - fix the deprecation of `management_group_name` in favour of `management_group_id` ([#15209](https://github.com/hashicorp/terraform-provider-azurerm/issues/15209))
* `azurerm_static_site` - fixing the creation of a Free tier Static Site ([#15141](https://github.com/hashicorp/terraform-provider-azurerm/issues/15141))
* `azurerm_storage_share` - fixing the `ShareBeingDeleted` error when the Storage Share is recreated ([#15180](https://github.com/hashicorp/terraform-provider-azurerm/issues/15180))
* 
## 2.94.0 (January 28, 2022)

UPGRADE NOTES:

* provider: support for the Azure German cloud has been removed in this release as this environment is no longer operational ([#14403](https://github.com/hashicorp/terraform-provider-azurerm/issues/14403))
* `azurerm_api_management_policy` - resources that were created with v2.92.0 will be marked as tainted due to a [bug](https://github.com/hashicorp/terraform-provider-azurerm/issues/15042). This version addresses the underlying issue, but the actual resource needs to either be untainted (via `terraform untaint`) or allow Terraform to delete the resource and create it again.
* `azurerm_hdinsight_kafka_cluster` - the `security_group_name` property in the `rest_proxy` block is conditionally required when the `use_msal` provider property is enabled ([#14403](https://github.com/hashicorp/terraform-provider-azurerm/issues/14403))

FEATURES:

* **New Data Source:** `azurerm_linux_function_app` ([#15009](https://github.com/hashicorp/terraform-provider-azurerm/issues/15009))
* **New Data Source** `azurerm_web_pubsub` ([#14731](https://github.com/hashicorp/terraform-provider-azurerm/issues/14731))
* **New Data Source** `azurerm_web_pubsub_hub` ([#14731](https://github.com/hashicorp/terraform-provider-azurerm/issues/14731))
* **New Resource:** `azurerm_web_pubsub` ([#14731](https://github.com/hashicorp/terraform-provider-azurerm/issues/14731))
* **New Resource:** `azurerm_web_pubsub_hub` ([#14731](https://github.com/hashicorp/terraform-provider-azurerm/issues/14731))
* **New Resource:** `azurerm_virtual_desktop_host_pool_registration_info` ([#14134](https://github.com/hashicorp/terraform-provider-azurerm/issues/14134))

ENHANCEMENTS:

* dependencies: updating to `v61.3.0` of `github.com/Azure/azure-sdk-for-go` ([#15080](https://github.com/hashicorp/terraform-provider-azurerm/issues/15080))
* dependencies: updating to `v0.21.0` of `github.com/hashicorp/go-azure-helpers` ([#15043](https://github.com/hashicorp/terraform-provider-azurerm/issues/15043))
* dependencies: updating `kusto` to API Version `2021-08-27` ([#15040](https://github.com/hashicorp/terraform-provider-azurerm/issues/15040))
* provider: opt-in support for v2 authentication tokens via the `use_msal` provider property ([#14403](https://github.com/hashicorp/terraform-provider-azurerm/issues/14403))
* `azurerm_app_service_slot`- support for the `storage_account` block ([#15084](https://github.com/hashicorp/terraform-provider-azurerm/issues/15084))
* `azurerm_stream_analytics_stream_input_eventhub` - support for the `partition_key` property ([#15019](https://github.com/hashicorp/terraform-provider-azurerm/issues/15019))

BUG FIXES:

* `data.image_source` - fix a regression around `id` ([#15119](https://github.com/hashicorp/terraform-provider-azurerm/issues/15119))
* `azurerm_api_management_backend` fix a crash caused by `backend_credentials` ([#15123](https://github.com/hashicorp/terraform-provider-azurerm/issues/15123))
* `azurerm_api_management_policy` - fixing the Resource ID for the `api_management_policy` block when this was provisioned using version `2.92.0` of the Azure Provider ([#15060](https://github.com/hashicorp/terraform-provider-azurerm/issues/15060))
* `azurerm_bastion_host` - fix a crash by adding nil check for the `copy_paste_enabled` property ([#15074](https://github.com/hashicorp/terraform-provider-azurerm/issues/15074))
* `azurerm_dev_test_lab` - fix an unexpected diff on with the `key_vault_id` property ([#15054](https://github.com/hashicorp/terraform-provider-azurerm/issues/15054))
* `azurerm_subscription_cost_management_export` - now sents the `ETag` when updating a cost management export ([#15017](https://github.com/hashicorp/terraform-provider-azurerm/issues/15017))
* `azurerm_template_deployment` - fixes a potential bug occuring during the deletion of a template deployment ([#15085](https://github.com/hashicorp/terraform-provider-azurerm/issues/15085))
* `azurerm_eventhub` - the `partition_count` property can now be changed when using Premium `sku` ([#15088](https://github.com/hashicorp/terraform-provider-azurerm/issues/15088))

## 2.93.1 (January 24, 2022)

BUG FIXES:

* `azurerm_app_service` - fix name availability check request ([#15062](https://github.com/hashicorp/terraform-provider-azurerm/issues/15062))

## 2.93.0 (January 21, 2022)

FEATURES:

* **New Data Source**: `azurerm_mysql_flexible_server` ([#14976](https://github.com/hashicorp/terraform-provider-azurerm/issues/14976))
* **New Beta Data Source**: `azurerm_windows_function_app` ([#14964](https://github.com/hashicorp/terraform-provider-azurerm/issues/14964))

ENHANCEMENTS: 

* dependencies: upgrading to `v61.1.0` of `github.com/Azure/azure-sdk-for-go` ([#14828](https://github.com/hashicorp/terraform-provider-azurerm/issues/14828))
* dependencies: updating `containerregistry` to API version `2021-08-01-preview` ([#14961](https://github.com/hashicorp/terraform-provider-azurerm/issues/14961))
* Data Source `azurerm_logic_app_workflow` - exporting the `identity` block ([#14896](https://github.com/hashicorp/terraform-provider-azurerm/issues/14896))
* `azurerm_bastion_host` - support for the `copy_paste_enabled`, `file_copy_enabled`, `ip_connect_enabled`, `shareable_link_enabled`, and `tunneling_enabled` properties ([#14987](https://github.com/hashicorp/terraform-provider-azurerm/issues/14987))
* `azurerm_bastion_host` - support for the `scale_units` property ([#14968](https://github.com/hashicorp/terraform-provider-azurerm/issues/14968))
* `azurerm_security_center_automation ` - the `event_source` property can now be set to `AssessmentsSnapshot`,
`RegulatoryComplianceAssessment`, `RegulatoryComplianceAssessmentSnapshot`, `SecureScoreControlsSnapshot`, `SecureScoresSnapshot`, and `SubAssessmentsSnapshot` ([#14996](https://github.com/hashicorp/terraform-provider-azurerm/issues/14996))
* `azurerm_static_site` - support for the `identity` block ([#14911](https://github.com/hashicorp/terraform-provider-azurerm/issues/14911))
* `azurerm_iothub` - Support for Identity-Based Endpoints ([#14705](https://github.com/hashicorp/terraform-provider-azurerm/issues/14705))
* `azurerm_servicebus_namespace_network_rule_set` -  support for the `public_network_access_enabled` property ([#14967](https://github.com/hashicorp/terraform-provider-azurerm/issues/14967))

BUG FIXES:

* `azurerm_machine_learning_compute_instance` - add validation for `tenant_id` and `object_id` properties to prevent null values and subsequent panic ([#14982](https://github.com/hashicorp/terraform-provider-azurerm/issues/14982))
* `azurerm_linux_function_app` - (beta) fix potential panic in `application_stack` when that block is not in config ([#14844](https://github.com/hashicorp/terraform-provider-azurerm/issues/14844))
* `azurerm_storage_share_file` changing the `content_md5` property will now trigger recreation and the `content_length` property of share file will now be set when updating properties. ([#15007](https://github.com/hashicorp/terraform-provider-azurerm/issues/15007))

## 2.92.0 (January 14, 2022)

FEATURES:

* **New Resource:** `azurerm_api_management_api_tag` ([#14711](https://github.com/hashicorp/terraform-provider-azurerm/issues/14711))
* **New Resource:** `azurerm_disk_pool_managed_disk_attachment` ([#14268](https://github.com/hashicorp/terraform-provider-azurerm/issues/14268))

ENHANCEMENTS:

* dependencies: upgrading `eventgrid` to API version `2021-12-01` ([#14433](https://github.com/hashicorp/terraform-provider-azurerm/issues/14433))
* `azurerm_api_management_custom_domain` - the `proxy` property has been deprecated in favour of the `gateway` for the 3.0 release ([#14628](https://github.com/hashicorp/terraform-provider-azurerm/issues/14628))
* `azurerm_databricks_workspace_customer_managed_key` - allow creation of resource when `infrastructure_encryption_enabled` is set to `true` for the databricks workspace ([#14915](https://github.com/hashicorp/terraform-provider-azurerm/issues/14915))
* `azurerm_eventgrid_domain` - support for the `local_auth_enabled`, `auto_create_topic_with_first_subscription`, and `auto_delete_topic_with_last_subscription` properties ([#14433](https://github.com/hashicorp/terraform-provider-azurerm/issues/14433))
* `azurerm_monitor_action_group` - support for the `event_hub_receiver` block ([#14771](https://github.com/hashicorp/terraform-provider-azurerm/issues/14771))
* `azurerm_mssql_server_extended_auditing_policy` - support storing audit data in storage account that is behind a firewall and VNet ([#14656](https://github.com/hashicorp/terraform-provider-azurerm/issues/14656))
* `azurerm_purview_account` - export the `managed_resources` block ([#14865](https://github.com/hashicorp/terraform-provider-azurerm/issues/14865))
* `azurerm_recovery_services_vault`- support for customer-managed keys (CMK) with the `encryption` block ([#14718](https://github.com/hashicorp/terraform-provider-azurerm/issues/14718))
* `azurerm_storage_account` - support for the `infrastructure_encryption_enabled` property ([#14864](https://github.com/hashicorp/terraform-provider-azurerm/issues/14864))

BUG FIXES:

* `azurerm_aadb2c_directory` - fix importing existing resources ([#14879](https://github.com/hashicorp/terraform-provider-azurerm/issues/14879))
* `azurerm_consumption_budget_subscription` - fix issue in migration logic ([#14898](https://github.com/hashicorp/terraform-provider-azurerm/issues/14898))
* `azurerm_cosmosdb_account` - only force ForceMongo when kind is set to mongo ([#14924](https://github.com/hashicorp/terraform-provider-azurerm/issues/14924))
* `azurerm_cosmosdb_mongo_collection` - now validates that "_id" is included as an index key ([#14857](https://github.com/hashicorp/terraform-provider-azurerm/issues/14857))
* `azurem_hdinsight` - hdinsight resources using oozie metastore can now be created without error ([#14880](https://github.com/hashicorp/terraform-provider-azurerm/issues/14880))
* `azurerm_log_analytics_datasource_windows_performance_counter` - state migration for case conversion of ID element ([#14916](https://github.com/hashicorp/terraform-provider-azurerm/issues/14916))
* `azurerm_monitor_aad_diagnostic_setting` - use the correct parser function for event hub rule IDs ([#14944](https://github.com/hashicorp/terraform-provider-azurerm/issues/14944))
* `azurerm_mysql_server_key` - fix issue when checking for existing resource on create ([#14883](https://github.com/hashicorp/terraform-provider-azurerm/issues/14883))
* `azurerm_spring_cloud_service` - fix panic when removing git repos ([#14900](https://github.com/hashicorp/terraform-provider-azurerm/issues/14900))
* `azurerm_log_analytics_workspace` - the `reservation_capcity_in_gb_per_day` has been deprecated and renamed to `reservation_capacity_in_gb_per_day` ([#14910](https://github.com/hashicorp/terraform-provider-azurerm/issues/14910))
* `azurerm_iothub_dps` - fixed default value of `allocation_weight` to match azure default ([#14943](https://github.com/hashicorp/terraform-provider-azurerm/issues/14943))
* `azurerm_iothub` - now exports `event_hub_events_namespace` and has a fallback route by default ([#14942](https://github.com/hashicorp/terraform-provider-azurerm/issues/14942))

## 2.91.0 (January 07, 2022)

FEATURES:

* **New Data Source:** `azurerm_aadb2c_directory` ([#14671](https://github.com/hashicorp/terraform-provider-azurerm/issues/14671))
* **New Data Source:** `azurerm_sql_managed_instance` ([#14739](https://github.com/hashicorp/terraform-provider-azurerm/issues/14739))
* **New Resource:** `azurerm_aadb2c_directory` ([#14671](https://github.com/hashicorp/terraform-provider-azurerm/issues/14671))
* **New Resource:** `azurerm_app_service_slot_custom_hostname_binding` ([#13097](https://github.com/hashicorp/terraform-provider-azurerm/issues/13097))
* **New Resource:** `azurerm_data_factory_linked_service_odbc` ([#14787](https://github.com/hashicorp/terraform-provider-azurerm/issues/14787))
* **New Resource:** `azurerm_disk_pool` ([#14675](https://github.com/hashicorp/terraform-provider-azurerm/issues/14675))
* **New Resource:** `azurerm_load_test` ([#14724](https://github.com/hashicorp/terraform-provider-azurerm/issues/14724))
* **New Resource:** `azurerm_virtual_desktop_scaling_plan` ([#14188](https://github.com/hashicorp/terraform-provider-azurerm/issues/14188))

ENHANCEMENTS:

* dependencies: upgrading `appplatform` to API version `2021-09-01-preview` ([#14365](https://github.com/hashicorp/terraform-provider-azurerm/issues/14365))
* dependencies: upgrading `network` to API Version `2021-05-01` ([#14164](https://github.com/hashicorp/terraform-provider-azurerm/issues/14164))
* dependencies: upgrading to `v60.2.0` of `github.com/Azure/azure-sdk-for-go` ([#14688](https://github.com/hashicorp/terraform-provider-azurerm/issues/14688)] and [[#14667](https://github.com/hashicorp/terraform-provider-azurerm/issues/14667))
* dependencies: upgrading to `v2.10.1` of `github.com/hashicorp/terraform-plugin-sdk` ([#14666](https://github.com/hashicorp/terraform-provider-azurerm/issues/14666))
* `azurerm_application_gateway` - support for the `key_vault_secret_id` and `force_firewall_policy_association` properties ([#14413](https://github.com/hashicorp/terraform-provider-azurerm/issues/14413))
* `azurerm_application_gateway` - support the `fips_enagled` property ([#14797](https://github.com/hashicorp/terraform-provider-azurerm/issues/14797))
* `azurerm_cdn_endpoint_custom_domain` - support for HTTPS ([#13283](https://github.com/hashicorp/terraform-provider-azurerm/issues/13283))
* `azurerm_hdinsight_hbase_cluster` - support for the `network` property ([#14825](https://github.com/hashicorp/terraform-provider-azurerm/issues/14825))
* `azurerm_iothub` - support for the `identity` block ([#14354](https://github.com/hashicorp/terraform-provider-azurerm/issues/14354))
* `azurerm_iothub_endpoint_servicebus_queue_resource` - depracating the `iothub_name` propertyin favour of `iothub_id` property ([#14690](https://github.com/hashicorp/terraform-provider-azurerm/issues/14690))
* `azurerm_iothub_endpoint_storage_container_resource` - depracating the `iothub_name` property in favour of `iothub_id` property [[#14690](https://github.com/hashicorp/terraform-provider-azurerm/issues/14690)] 
* `azurerm_iot_fallback_route` - support for the `source` property ([#14836](https://github.com/hashicorp/terraform-provider-azurerm/issues/14836))
* `azurerm_kubernetes_cluster` - support for the `public_network_access_enabled`, `scale_down_mode`, and `workload_runtime` properties ([#14386](https://github.com/hashicorp/terraform-provider-azurerm/issues/14386))
* `azurerm_linux_function_app` - (Beta Resource) fix the filtering of `app_settings` for `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#14815](https://github.com/hashicorp/terraform-provider-azurerm/issues/14815))
* `azurerm_linux_virtual_machine` - support for the `user_data` property ([#13888](https://github.com/hashicorp/terraform-provider-azurerm/issues/13888))
* `azurerm_linux_virtual_machine_scale_set` - support for the `user_data` property ([#13888](https://github.com/hashicorp/terraform-provider-azurerm/issues/13888))
* `azurerm_managed_disk` - support for the `gallery_image_reference_id` property ([#14121](https://github.com/hashicorp/terraform-provider-azurerm/issues/14121))
* `azurerm_mysql_server` - support capacities up to `16TB` for the `storage_mb` property ([#14838](https://github.com/hashicorp/terraform-provider-azurerm/issues/14838))
* `azurerm_postgresql_flexible_server` - support for the `geo_redundant_backup_enabled` property ([#14661](https://github.com/hashicorp/terraform-provider-azurerm/issues/14661))
* `azurerm_recovery_services_vault` - support for the `storage_mode_type` property ([#14659](https://github.com/hashicorp/terraform-provider-azurerm/issues/14659))
* `azurerm_spring_cloud_certificate` - support for the `certificate_content` property ([#14689](https://github.com/hashicorp/terraform-provider-azurerm/issues/14689))
* `azurerm_servicebus_namespace_authorization_rule` - the `resource_group_name` and `namespace_name` properties have been deprecated in favour of the `namespace_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_namespace_network_rule_set` - the `resource_group_name` and `namespace_name` properties have been deprecated in favour of the `namespace_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_namespace_authorization_rule` - the `resource_group_name` and `namespace_name` properties have been deprecated in favour of the `namespace_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_queue` - the `resource_group_name` and `namespace_name` properties have been deprecated in favour of the `namespace_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_queue_authorization_rule` - the `resource_group_name`, `namespace_name`, and `queue_name` properties have been deprecated in favour of the `queue_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_subscription` - the `resource_group_name`, `namespace_name`, and `topic_name` properties have been deprecated in favour of the `topic_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_subscription_rule` - the `resource_group_name`, `namespace_name`, `topic_name`, and `subscription_name` properties have been deprecated in favour of the `subscription_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_topic` - the `resource_group_name` and `namespace_name` properties have been deprecated in favour of the `namespace_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_topic_authorization_rule` - the `resource_group_name`, `namespace_name`, and `topic_name` properties have been deprecated in favour of the `topic_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_shared_image_version` - images can now be sorted by semver ([#14708](https://github.com/hashicorp/terraform-provider-azurerm/issues/14708))
* `azurerm_virtual_network_gateway_connection` - support for the `connection_mode` property ([#14738](https://github.com/hashicorp/terraform-provider-azurerm/issues/14738))
* `azurerm_web_application_firewall_policy` - the `file_upload_limit_in_mb` property within the `policy_settings` block can now be set to `4000` ([#14715](https://github.com/hashicorp/terraform-provider-azurerm/issues/14715))
* `azurerm_windows_virtual_machine` - support for the `user_data` property ([#13888](https://github.com/hashicorp/terraform-provider-azurerm/issues/13888))
* `azurerm_windows_virtual_machine_scale_set` - support for the `user_data` property ([#13888](https://github.com/hashicorp/terraform-provider-azurerm/issues/13888))

BUG FIXES:

* `azurerm_app_service_environment_v3` - fix the default value of the `allow_new_private_endpoint_connections` property ([#14805](https://github.com/hashicorp/terraform-provider-azurerm/issues/14805))
* `azurerm_consumption_budget_subscription` - added an additional state migration to fix the bug introduced by the first one and to parse the `subscription_id` from the resource's ID ([#14803](https://github.com/hashicorp/terraform-provider-azurerm/issues/14803))
* `azurerm_network_interface_security_group_association` - checking the ID matches the expected format during import ([#14753](https://github.com/hashicorp/terraform-provider-azurerm/issues/14753))
* `azurerm_storage_management_policy` - handle the unexpected deletion of the storage account ([#14799](https://github.com/hashicorp/terraform-provider-azurerm/issues/14799))

## 2.90.0 (December 17, 2021)

FEATURES:

* **New Data Source:** `azurerm_app_configuration_key` ([#14484](https://github.com/hashicorp/terraform-provider-azurerm/issues/14484))
* **New Resource:** `azurerm_container_registry_task` ([#14533](https://github.com/hashicorp/terraform-provider-azurerm/issues/14533))
* **New Resource:** `azurerm_maps_creator` ([#14566](https://github.com/hashicorp/terraform-provider-azurerm/issues/14566))
* **New Resource:** `azurerm_netapp_snapshot_policy` ([#14230](https://github.com/hashicorp/terraform-provider-azurerm/issues/14230))
* **New Resource:** `azurerm_synapse_sql_pool_workload_classifier` ([#14412](https://github.com/hashicorp/terraform-provider-azurerm/issues/14412))
* **New Resource:** `azurerm_synapse_workspace_sql_aad_admin` ([#14341](https://github.com/hashicorp/terraform-provider-azurerm/issues/14341))
* **New Resource:** `azurerm_vpn_gateway_nat_rule` ([#14527](https://github.com/hashicorp/terraform-provider-azurerm/issues/14527))

ENHANCEMENTS:

* dependencies: updating `apimanagement` to API Version `2021-08-01` ([#14312](https://github.com/hashicorp/terraform-provider-azurerm/issues/14312))
* dependencies: updating `managementgroups` to API Version `2020-05-01` ([#14635](https://github.com/hashicorp/terraform-provider-azurerm/issues/14635))
* dependencies: updating `redisenterprise` to use an Embedded SDK ([#14502](https://github.com/hashicorp/terraform-provider-azurerm/issues/14502))
* dependencies: updating to `v0.19.1` of `github.com/hashicorp/go-azure-helpers` ([#14627](https://github.com/hashicorp/terraform-provider-azurerm/issues/14627))
* dependencies: updating to `v2.10.0` of `github.com/hashicorp/terraform-plugin-sdk` ([#14596](https://github.com/hashicorp/terraform-provider-azurerm/issues/14596))
* Data Source: `azurerm_function_app_host_keys` - support for `signalr_extension_key` and `durabletask_extension_key` ([#13648](https://github.com/hashicorp/terraform-provider-azurerm/issues/13648))
* `azurerm_application_gateway ` - support for private link configurations ([#14583](https://github.com/hashicorp/terraform-provider-azurerm/issues/14583))
* `azurerm_blueprint_assignment` - support for the `lock_exclude_actions` property ([#14648](https://github.com/hashicorp/terraform-provider-azurerm/issues/14648))
* `azurerm_container_group` - support for `ip_address_type = None` ([#14460](https://github.com/hashicorp/terraform-provider-azurerm/issues/14460))
* `azurerm_cosmosdb_account` - support for the `create_mode` property and `restore` block ([#14362](https://github.com/hashicorp/terraform-provider-azurerm/issues/14362))
* `azurerm_data_factory_dataset_*` - deprecate `data_factory_name` in favour of `data_factory_id` for consistency across all data factory dataset resources ([#14610](https://github.com/hashicorp/terraform-provider-azurerm/issues/14610))
* `azurerm_data_factory_integration_runtime_*`- deprecate `data_factory_name` in favour of `data_factory_id` for consistency across all data factory integration runtime resources ([#14610](https://github.com/hashicorp/terraform-provider-azurerm/issues/14610))
* `azurerm_data_factory_trigger_*`- deprecate `data_factory_name` in favour of `data_factory_id` for consistency across all data factory trigger resources ([#14610](https://github.com/hashicorp/terraform-provider-azurerm/issues/14610))
* `azurerm_data_factory_pipeline`- deprecate `data_factory_name` in favour of `data_factory_id` for consistency across all data factory resources ([#14610](https://github.com/hashicorp/terraform-provider-azurerm/issues/14610))
* `azurerm_iothub` - support for the `cloud_to_device` block ([#14546](https://github.com/hashicorp/terraform-provider-azurerm/issues/14546))
* `azurerm_iothub_endpoint_eventhub` - the `iothub_name` property has been deprecated in favour of the `iothub_id` property ([#14632](https://github.com/hashicorp/terraform-provider-azurerm/issues/14632))
* `azurerm_logic_app_workflow` - support for the `open_authentication_policy` block ([#14007](https://github.com/hashicorp/terraform-provider-azurerm/issues/14007))
* `azurerm_signalr` - support for the `live_trace_enabled` property ([#14646](https://github.com/hashicorp/terraform-provider-azurerm/issues/14646))
* `azurerm_xyz_policy_assignment` add support for `non_compliance_message` ([#14518](https://github.com/hashicorp/terraform-provider-azurerm/issues/14518))

BUG FIXES:

* `azurerm_cosmosdb_account` - will now set a default value for `default_identity_type` when the API return a nil value ([#14643](https://github.com/hashicorp/terraform-provider-azurerm/issues/14643))
* `azurerm_function_app` - address `app_settings` during creation rather than just updates ([#14638](https://github.com/hashicorp/terraform-provider-azurerm/issues/14638))
* `azurerm_marketplace_agreement` - fix crash when the import check triggers ([#14614](https://github.com/hashicorp/terraform-provider-azurerm/issues/14614))
* `azurerm_postgresql_configuration` - now locks during write operations to prevent conflicts ([#14619](https://github.com/hashicorp/terraform-provider-azurerm/issues/14619))
* `azurerm_postgresql_flexible_server_configuration` - now locks during write operations to prevent conflicts ([#14607](https://github.com/hashicorp/terraform-provider-azurerm/issues/14607))

---

For information on changes between the v2.89.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v2.00.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
