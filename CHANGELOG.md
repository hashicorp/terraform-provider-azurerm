## 2.65.0 (Unreleased)

ENHANCEMENTS:

* dependencies: updating to `v2.6.1` of `github.com/hashicorp/terraform-plugin-sdk` [GH-12209]
* dependencies: upgrading to `v55.3.0` of `github.com/Azure/azure-sdk-for-go` [GH-12263]
* dependencies: updating to `v0.11.19` of `github.com/Azure/go-autorest/autorest` [GH-12209]
* dependencies: updating to `v0.9.14` of `github.com/Azure/go-autorest/autorest/adal` [GH-12209]
* dependencies: updating the embedded SDK for Eventhub Namespaces to use API Version `2021-01-01-preview` [GH-12290]

BUG FIXES:

* `azurerm_data_factory` - fix a bug where the `name` property was stored with the wrong casing [GH-12128]

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
