## 2.69.0 (Unreleased)

FEATURES:

* **New Resource** `azurerm_batch_job` [GH-12573]
* **New Resource** `azurerm_data_factory_managed_private_endpoint` [GH-12618]
* **New Resource** `azurerm_data_protection_backup_policy_blob_storage` [GH-12362]

ENHANCEMENTS:

* dependencies: Updgrading to `v55.6.0` of `github.com/Azure/azure-sdk-for-go` [GH-12565]
* `azurerm_api_management_named_value` - the field `secret_id` can now be set to a versionless Key Vault Key [GH-12641]
* `azurerm_data_factory_integration_runtime_azure_ssis` - support for the `public_ips`, `express_custom_setup`, `package_store`, and `proxy` blocks [GH-12545]
* `azurerm_bot_channels_registration` - support for the `cmk_key_vault_url`, `description`, `icon_url`, and `isolated_network_enabled` [GH-12560]
* `azurerm_data_factory_integration_runtime_azure` - support for the `virtual_network_enabled` property [GH-12619]
* `azurerm_kubernetes_cluster` - support for downgrading `sku_tier` from `Paid` to `Free` without recreating the Cluster [GH-12651]
* `azurerm_postgresql_flexible_server` - support for the `high_availability` block [GH-12587]

BUG FIXES:

* `data.azurerm_redis_cache` - fix a bug that caused the data source to raise an error [GH-12666]
* `azurerm_application_gateway` - return an error when ssl policy is not properly configured  [GH-12647]
* `azurerm_data_factory_linked_custom_service` - fix a bug causing `additional_properties` to be read incorrectly into state [GH-12664]
* `azurerm_eventhub_authorization_rule` - fixing the error "empty non-retryable error received" [GH-12642]
* `azurerm_machine_learning_compute_cluster` - fix a crash when creating a cluster without specifying `subnet_resource_id` [GH-12658]

## 2.68.0 (July 16, 2021)

FEATURES:

* **New Data Source** `azurerm_local_network_gateway` ([#12579](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12579))
* **New Resource** `azurerm_api_management_api_release` ([#12562](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12562))
* **New Resource** `azurerm_data_protection_backup_policy_disk` ([#12361](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12361))
* **New Resource** `azurerm_data_factory_custom_dataset` ([#12484](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12484))
* **New Resource** `azurerm_data_factory_dataset_binary` ([#12369](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12369))
* **New Resource** `azurerm_maintenance_assignment_virtual_machine_scale_set` ([#12273](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12273))
* **New Resource** `azurerm_postgresql_flexible_server_configuration` ([#12294](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12294))
* **New Resource** `azurerm_synapse_private_link_hub` ([#12495](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12495))

ENHANCEMENTS:

* dependencies: upgrading to `v55.5.0` of `github.com/Azure/azure-sdk-for-go` ([#12435](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12435))
* dependencies: updating `bot` to use API Version `2021-03-01` ([#12449](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12449))
* dependencies: updating `maintenance` to use API Version `2021-05-01` ([#12273](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12273))
* `azurerm_api_management_named_value` - support for the `value_from_key_vault` block ([#12309](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12309))
* `azurerm_api_management_api_diagnostic` - support for the `data_masking`1 property ([#12419](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12419))
* `azurerm_cognitive_account` - support for the `identity`, `storage`, `disable_local_auth`, `fqdns`, `public_network_access_enabled`, and `restrict_outbound_network_access` properties ([#12469](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12469))
* `azurerm_cognitive_account` - the `virtual_network_subnet_ids` property has been deprecated in favour of `virtual_network_rules` block to supoport the `ignore_missing_vnet_service_endpoint` property ([#12600](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12600))
* `azurerm_container_registry` - now exports the `principal_id` and `tenant_id` attributes in the `identity` block ([#12378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12378))
* `azurerm_data_factory` - support for the `managed_virtual_network_enabled` property ([#12343](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12343))
* `azurerm_linux_virtual_machine_scale_set` - Fix un-necessary VMSS instance rolling request ([#12590](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12590))
* `azurerm_maintenance_configuration` - support for the `window`, `visibility`, and `properties` blocks ([#12273](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12273))
* `azurerm_powerbi_embedded` - support for the `mode` property ([#12394](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12394))
* `azurerm_redis_cache` - support for the `maintenance_window` property in the `patch_schedule` block ([#12472](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12472))
* `azurerm_storage_account_customer_managed_key` - support for the `user_assigned_identity_id` property ([#12516](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12516))

BUG FIXES:

* `azurerm_api_management` - no longer forces a new resource when changing the `subnet_id` property ([#12611](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12611))
* `azurerm_function_app` - set a default value for `os_type` and allow a blank string to be specified as per documentation ([#12482](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12482))
* `azurerm_key_vault_access_policy` - prevent a possible panic on delete ([#12616](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12616))
* `azurerm_postgresql_flexible_server` - add new computed property `private_dns_zone_id` to work around a upcomming breaking change in the API ([#12288](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12288))
* `machine_learning_compute_cluster` - make the `subnet_resource_id` property actually optional ([#12558](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12558))
* `azurerm_mssql_database` - don't allow license_type to be set for serverless SQL databases ([#12555](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12555))
* `azurerm_subnet_network_security_group_association` - prevent potential deadlocks when using multiple association resources ([#12267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12267))

## 2.67.0 (July 09, 2021)

FEATURES:

* **New Data Source** `azurerm_api_management_gateway` ([#12297](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12297))
* **New Resource** `azurerm_api_management_gateway` ([#12297](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12297))
* **New Resource** `azurerm_databricks_workspace_customer_managed_key`([#12331](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12331))

ENHANCEMENTS:

* dependencies: updating `postgresqlflexibleservers` to use API Version `2021-06-01` ([#12405](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12405))
* `azurerm_databricks_workspace` - add support for `machine_learning_workspace_id`, `customer_managed_key_enabled`, `infrastructure_encryption_enabled` and `storage_account_identity` ([#12331](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12331))
* `azurerm_security_center_assessment_policy` - support for the `categories` propety ([#12383](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12383))

BUG FIXES:

* `azurerm_api_management` - fix an issue where changing the location of an `additional_location` would force a new resource ([#12468](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12468))
* `azurerm_app_service` - fix crash when resource group or ASE is missing. ([#12518](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12518))
* `azurerm_automation_variable_int` - fixed value parsing order causing `1` to be considered a bool ([#12511](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12511))
* `azurerm_automation_variable_bool` - fixed value parsing order causing `1` to be considered a bool ([#12511](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12511))
* `azurerm_data_factory_dataset_parquet` - the `azure_blob_storage_location.filename` property cis now optional ([#12414](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12414))
* `azurerm_kusto_eventhub_data_connection` - `APACHEAVRO` can now be used as a `data_format` option ([#12480](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12480))
* `azurerm_site_recovery_replicated_vm ` - Fix potential crash in reading `managed_disk` properties ([#12509](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12509))
* `azurerm_storage_account` - `account_replication_type` can now be updated ([#12479](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12479))
* `azurerm_storage_management_policy` - fix crash in read of properties ([#12487](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12487))
* `azurerm_storage_share_directory` now allows underscore in property `name` [[#12454](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12454)] 
* `azurerm_security_center_subscription_pricing` - removed Owner permission note from documentation ([#12481](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12481))

DEPRECATIONS:

* `azurerm_postgresql_flexible_server` - the `cmk_enabled` property has been deprecated as it has been removed from the API ([#12405](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12405))
* `azurerm_virtual_machine_configuration_policy_assignment` - has been deprecated and renamed to `azurerm_policy_virtual_machine_configuration_assignment` ([#12497](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12497))

## 2.66.0 (July 02, 2021)

FEATURES:

* **New Resource** `azurerm_api_management_api_operation_tag` ([#12384](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12384))
* **New Resource** `azurerm_data_factory_linked_custom_service` ([#12224](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12224))
* **New Resource** `azurerm_data_factory_trigger_blob_event` ([#12330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12330))
* **New Resource** `azurerm_express_route_connection` ([#11320](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11320))
* **New Resource** `azurerm_express_route_circuit_connection` ([#11303](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11303))
* **New Resource** `azurerm_management_group_policy_assignment` ([#12349](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12349))
* **New Resource** `azurerm_resource_group_policy_assignment` ([#12349](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12349))
* **New Resource** `azurerm_resource_policy_assignment` ([#12349](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12349))
* **New Resource** `azurerm_subscription_policy_assignment` ([#12349](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12349))
* **New resource** `azurerm_tenant_configuration` ([#11697](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11697))
* Cognitive Service now supports purging soft delete accounts ([#12281](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12281))

ENHANCEMENTS:

* dependencies: updating `cognitive` to use API Version `2021-03-01` ([#12281](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12281))
* dependencies: updating `trafficmanager` to use API Version `2018-08-01` ([#12400](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12400))
* `azurerm_api_management_backend` - support for the `client_certificate_id` property  ([#12402](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12402))
* `azurerm_api_management_api` - support for the `revision_description`, `version_description`, and `source_api_id` properties ([#12266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12266))
* `azurerm_batch_account` - support for the `public_network_access_enabled` property ([#12401](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12401))
* `azurerm_eventgrid_event_subscription` - support for additional advanced filters `string_not_begins_with`, `string_not_ends_with`, `string_not_contains`, `is_not_null`, `is_null_or_undefined`, `number_in_range` and `number_not_in_range` ([#12167](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12167))
* `azurerm_eventgrid_system_topic_event_subscription` - support for additional advanced filters `string_not_begins_with`, `string_not_ends_with`, `string_not_contains`, `is_not_null`, `is_null_or_undefined`, `number_in_range` and `number_not_in_range` ([#12167](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12167))
* `azurerm_kubernetes_cluster` - support for the `fips_enabled`, `kubelet_disk_type`, and `license` properties ([#11835](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11835))
* `azurerm_kubernetes_cluster_node_pool` - support for the `fips_enabled`, and `kubelet_disk_type` properties ([#11835](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11835))
* `azurerm_lighthouse_definition` - support for the `plan` block ([#12360](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12360))
* `azurerm_site_recovery_replicated_vm` - Add support for `target_disk_encryption_set_id` in `managed_disk` ([#12374](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12374))
* `azurerm_traffic_manager_endpoint` - supports for the `minimum_required_child_endpoints_ipv4` and `minimum_required_child_endpoints_ipv6` ([#12400](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12400))

BUG FIXES:

* `azurerm_app_service` - fix app_setting and SCM setting ordering ([#12280](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12280))
* `azurerm_hdinsight_kafka_cluster` - will no longer panic from an empty `component_version` property ([#12261](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12261))
* `azurerm_spatial_anchors_account` - the `tags` property can now be updated without creating a new resource ([#11985](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11985))
* **Data Source** `azurerm_app_service_environment_v3` - fix id processing for Read ([#12436](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12436))


## 2.65.0 (June 25, 2021)

FEATURES:

* **New Resource** `azurerm_data_protection_backup_instance_postgresql` ([#12220](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12220))
* **New Resource** `azurerm_hpc_cache_blob_nfs_target` ([#11671](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11671))
* **New Resource** `azurerm_nat_gateway_public_ip_prefix_association` ([#12353](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12353))

ENHANCEMENTS:

* dependencies: updating to `v2.6.1` of `github.com/hashicorp/terraform-plugin-sdk` ([#12209](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12209))
* dependencies: upgrading to `v55.3.0` of `github.com/Azure/azure-sdk-for-go` ([#12263](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12263))
* dependencies: updating to `v0.11.19` of `github.com/Azure/go-autorest/autorest` ([#12209](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12209))
* dependencies: updating to `v0.9.14` of `github.com/Azure/go-autorest/autorest/adal` ([#12209](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12209))
* dependencies: updating the embedded SDK for Eventhub Namespaces to use API Version `2021-01-01-preview` ([#12290](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12290))
* `azurerm_express_route_circuit_peering` - support for the `bandwidth_in_gbps` and `express_route_port_id` properties ([#12289](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12289))
* `azurerm_kusto_iothub_data_connection` - support for the `data_format`, `mapping_rule_name` and `table_name` properties ([#12293](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12293))
* `azurerm_linux_virtual_machine` - updating `proximity_placement_group_id` will no longer create a new resoruce ([#11790](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11790))
* `azurerm_security_center_assessment_metadata` - support for the `categories` property ([#12278](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12278))
* `azurerm_windows_virtual_machine` - updating `proximity_placement_group_id` will no longer create a new resoruce ([#11790](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11790))

BUG FIXES:

* `azurerm_data_factory` - fix a bug where the `name` property was stored with the wrong casing ([#12128](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12128))

## 2.64.0 (June 18, 2021)

FEATURES:

* **New Data Source** `azurerm_key_vault_secrets` ([#12147](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12147))
* **New Resource** `azurerm_api_management_redis_cache` ([#12174](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12174))
* **New Resource** `azurerm_data_factory_linked_service_odata` ([#11556](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11556))
* **New Resource** `azurerm_data_protection_backup_policy_postgresql` ([#12072](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12072))
* **New Resource** `azurerm_machine_learning_compute_cluster` ([#11675](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11675))
* **New Resource** `azurerm_eventhub_namespace_customer_managed_key` ([#12159](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12159))
* **New Resource** `azurerm_virtual_desktop_application` ([#12077](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12077))

ENHANCEMENTS:

* dependencies: updating to `v55.2.0` of `github.com/Azure/azure-sdk-for-go` ([#12153](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12153))
* dependencies: updating `synapse` to use API Version `2021-03-01` ([#12183](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12183))
* `azurerm_api_management` - support for the `client_certificate_enabled`, `gateway_disabled`, `min_api_version`, and `zones` propeties ([#12125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12125))
* `azurerm_api_management_api_schema` - prevent plan not empty after apply for json definitions  ([#12039](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12039))
* `azurerm_application_gateway` - correctly poopulat the `identity` block ([#12226](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12226))
* `azurerm_container_registry` - support for the `zone_redundancy_enabled` field ([#11706](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11706))
* `azurerm_cosmosdb_sql_container` - support for the `spatial_index` block ([#11625](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11625))
* `azurerm_cosmos_gremlin_graph` - support for the `spatial_index` property ([#12176](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12176))
* `azurerm_data_factory` - support for `global_parameter` ([#12178](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12178))
* `azurerm_kubernetes_cluster` - support for the `kubelet_config` and `linux_os_config` blocks ([#11119](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11119))
* `azurerm_monitor_metric_alert` - support the `StartsWith` dimension operator ([#12181](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12181))
* `azurerm_private_link_service`  - changing `load_balancer_frontend_ip_configuration_ids` list no longer creates a new resource ([#12250](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12250))
* `azurerm_stream_analytics_job` - supports for the `identity` block ([#12171](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12171))
* `azurerm_storage_account` - support for the `share_properties` block ([#12103](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12103))
* `azurerm_synapse_workspace` - support for the `data_exfiltration_protection_enabled` property ([#12183](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12183))
* `azurerm_synapse_role_assignment` - support for scopes and new role types ([#11690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11690))

BUG FIXES:

* `azurerm_synapse_role_assignment` - support new roles and scopes ([#11690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11690))
* `azurerm_lb` - fix zone behaviour bug introduced in recent API upgrade ([#12208](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12208))

## 2.63.0 (June 11, 2021)

FEATURES:

* **New Resource** `azurerm_data_factory_linked_service_azure_search` ([#12122](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12122))
* **New Resource** `azurerm_data_factory_linked_service_kusto` ([#12152](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12152))

ENHANCEMENTS:

* dependencies: updating `streamanalytics` to use API Version `2020-03-01-preview` ([#12133](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12133))
* dependencies: updating `virtualdesktop` to use API Version `2020-11-02-preview` ([#12160](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12160))
* `data.azurerm_synapse_workspace` - support for the `identity` attribute ([#12098](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12098))
* `azurerm_cosmosdb_gremlin_graph` - support for the `composite_index` and `partition_key_version` properties ([#11693](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11693))
* `azurerm_data_factory_dataset_azure_blob` - support for the `dynamic_filename_enabled` and `dynamic_path_enabled` properties ([#12034](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12034))
* `azurerm_data_factory_dataset_delimited_text` - supports the `azure_blob_fs_location` property ([#12041](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12041))
* `azurerm_data_factory_linked_service_azure_sql_database` - support for the `key_vault_connection_string` property ([#12139](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12139))
* `azurerm_data_factory_linked_service_sql_server` - add `key_vault_connection_string` argument ([#12117](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12117))
* `azurerm_data_factory_linked_service_data_lake_storage_gen2` - supports for the `storage_account_key` property ([#12136](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12136))
* `azurerm_eventhub` - support for the `status` property ([#12043](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12043))
* `azurerm_kubernetes_cluster` - support migration of `service_principal` to `identity` ([#12049](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12049))
* `azurerm_kubernetes_cluster` -support for BYO `kubelet_identity` ([#12037](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12037))
* `azurerm_kusto_cluster_customer_managed_key` - supports for the `user_identity` property ([#12135](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12135))
* `azurerm_network_watcher_flow_log` - support for the `location` and `tags` properties ([#11670](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11670))
* `azurerm_storage_account` - support for user assigned identities ([#11752](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11752))
* `azurerm_storage_account_customer_managed_key` - support the use of keys from key vaults in remote subscription ([#12142](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12142))
* `azurerm_virtual_desktop_host_pool` - support for the `start_vm_on_connect` property ([#12160](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12160))
* `azurerm_vpn_server_configuration` - now supports multiple `auth` blocks ([#12085](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12085))

BUG FIXES:

* Service: App Configuration - Fixed a bug in tags on resources all being set to the same value ([#12062](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12062))
* Service: Event Hubs - Fixed a bug in tags on resources all being set to the same value ([#12062](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12062))
* `azurerm_subscription` - fix ability to specify `DevTest` as `workload` ([#12066](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12066))
* `azurerm_sentinel_alert_rule_scheduled` - the query frequency duration can noe be up to 14 days ([#12164](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12164))

## 2.62.1 (June 08, 2021)

BUG FIXES:

* `azurerm_role_assignment` - use the correct ID when assigning roles to resources ([#12076](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12076))


## 2.62.0 (June 04, 2021)

FEATURES:

* **New Resource** `azurerm_data_protection_backup_vault` ([#11955](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11955))
* **New Resource** `azurerm_postgresql_flexible_server_firewall_rule` ([#11834](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11834))
* **New Resource** `azurerm_vmware_express_route_authorization` ([#11812](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11812))
* **New Resource** `azurerm_storage_object_replication_policy` ([#11744](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11744))

ENHANCEMENTS:

* dependencies: updating `network` to use API Version `2020-11-01` ([#11627](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11627))
* `azurerm_app_service_environment` - support for the `internal_ip_address`, `service_ip_address`, and `outbound_ip_addresses`properties ([#12026](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12026))
* `azurerm_api_management_api_subscription` - support for the `api_id` property ([#12025](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12025))
* `azurerm_container_registry` - support for  versionless encryption keys for ACR ([#11856](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11856))
* `azurerm_kubernetes_cluster` -  support for `gateway_name` for Application Gateway add-on ([#11984](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11984))
* `azurerm_kubernetes_cluster` - support update of `azure_rbac_enabled` ([#12029](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12029))
* `azurerm_kubernetes_cluster` - support for `node_public_ip_prefix_id` ([#11635](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11635))
* `azurerm_kubernetes_cluster_node_pool` - support for `node_public_ip_prefix_id` ([#11635](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11635))
* `azurerm_machine_learning_inference_cluster` - support for the `ssl.leaf_domain_label` and `ssl.overwrite_existing_domain` properties ([#11830](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11830))
* `azurerm_role_assignment` - support the `delegated_managed_identity_resource_id` property ([#11848](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11848))

BUG FIXES:

* `azuerrm_postgres_server` - do no update `password` unless its changed ([#12008](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12008))
* `azuerrm_storage_acount` - prevent `containerDeleteRetentionPolicy` and `lastAccessTimeTrackingPolicy` not supported in `AzureUSGovernment` errors ([#11960](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11960))

## 2.61.0 (May 27, 2021)

FEATURES:

* **New Data Source:** `azurerm_spatial_anchors_account` ([#11824](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11824))

ENHANCEMENTS:

* dependencies: updating to `v54.3.0` of `github.com/Azure/azure-sdk-for-go` ([#11813](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11813))
* dependencies: updating `mixedreality` to use API Version `2021-01-01` ([#11824](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11824))
* refactor: switching to use an embedded SDK for `appconfiguration` ([#11959](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11959))
* refactor: switching to use an embedded SDK for `eventhub` ([#11973](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11973))
* provider: support for the Virtual Machine `skip_shutdown_and_force_delete` feature ([#11216](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11216))
* provider: support for the Virtual Machine Scale Set `force_delete` feature ([#11216](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11216))
* provider: no longer auto register the Microsoft.DevSpaces RP ([#11822](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11822))
* Data Source: `azurerm_key_vault_certificate_data` - support certificate bundles and add support for ECDSA keys ([#11974](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11974))
* `azurerm_data_factory_linked_service_sftp` - support for hostkey related properties ([#11825](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11825))
* `azurerm_spatial_anchors_account` - support for `account_domain` and `account_id` ([#11824](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11824))
* `azurerm_static_site` - Add support for `tags` attribute ([#11849](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11849))
* `azurerm_storage_account` - `private_link_access` supports more values ([#11957](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11957))
* `azurerm_storage_account_network_rules`: `private_link_access` supports more values ([#11957](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11957))
* `azurerm_synapse_spark_pool` - `spark_version` now supports `3.0` ([#11972](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11972))

BUG FIXES:

* `azurerm_cdn_endpoint` - do not send an empty `origin_host_header` to the api ([#11852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11852))
* `azurerm_linux_virtual_machine_scale_set`: changing the `disable_automatic_rollback` and `enable_automatic_os_upgrade` properties no longer created a new resource ([#11723](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11723))
* `azurerm_storage_share`: Fix ID for `resource_manager_id` ([#11828](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11828))
* `azurerm_windows_virtual_machine_scale_set`: changing the `disable_automatic_rollback` and `enable_automatic_os_upgrade` properties no longer created a new resource ([#11723](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11723))

## 2.60.0 (May 20, 2021)

FEATURES:

* **New Data Source:** `azurerm_eventhub_cluster` ([#11763](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11763))
* **New Data Source:** `azurerm_redis_enterprise_database` ([#11734](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11734))
* **New Resource:** `azurerm_static_site` ([#7150](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7150))
* **New Resource:** `azurerm_machine_learning_inference_cluster` ([#11550](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11550))

ENHANCEMENTS:

* dependencies: updating `aks` to use API Version `2021-03-01` ([#11708](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11708))
* dependencies: updating `eventgrid` to use API Version `2020-10-15-preview` ([#11746](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11746))
* `azurerm_cosmosdb_mongo_collection` - support for the `analytical_storage_ttl` property ([#11735](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11735))
* `azurerm_cosmosdb_cassandra_table` - support for the `analytical_storage_ttl` property ([#11755](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11755))
* `azurerm_healthcare_service` - support for the `public_network_access_enabled` property ([#11736](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11736))
* `azurerm_hdinsight_kafka_cluster` - support for the `encryption_in_transit_enabled` property ([#11737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11737))
* `azurerm_media_services_account` - support for the `key_delivery_access_control` block ([#11726](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11726))
* `azurerm_monitor_activity_log_alert` - support for `Security` event type for Azure Service Health alerts ([#11802](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11802))
* `azurerm_netapp_volume` - support for the `security_style` property - ([#11684](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11684))
* `azurerm_redis_cache` - suppot for the `replicas_per_master` peoperty ([#11714](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11714))
* `azurerm_spring_cloud_service` - support for the `required_network_traffic_rules` block ([#11633](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11633))
* `azurerm_storage_account_management_policy` - the `name` property can now contain `-` ([#11792](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11792))

BUG FIXES:

* `azurerm_frontdoor` - added a check for `nil` to avoid panic on destroy ([#11720](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11720))
* `azurerm_linux_virtual_machine_scale_set` - the `extension` blocks are now a set ([#11425](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11425))
* `azurerm_virtual_network_gateway_connection` - fix a bug where `shared_key` was not being updated ([#11742](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11742))
* `azurerm_windows_virtual_machine_scale_set` - the `extension` blocks are now a set ([#11425](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11425))
* `azurerm_windows_virtual_machine_scale_set` - changing the `license_type` will no longer create a new resource ([#11731](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11731))

---

For information on changes between the v2.59.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
