## 2.87.0 (November 26, 2021)

FEATURES:

* **New Resource:** `azurerm_api_management_notification_recipient_user` ([#14239](https://github.com/hashicorp/terraform-provider-azurerm/issues/14239))
* **New Resource:** `azurerm_app_service_public_certificate` ([#14337](https://github.com/hashicorp/terraform-provider-azurerm/issues/14337))
* **New Resource:** `azurerm_service_fabric_managed_cluster` ([#14131](https://github.com/hashicorp/terraform-provider-azurerm/issues/14131))
* **New Resource:** `azurerm_sentinel_watchlist` ([#14258](https://github.com/hashicorp/terraform-provider-azurerm/issues/14258))
* **New Resource:** `azurerm_static_site_custom_domain` ([#12764](https://github.com/hashicorp/terraform-provider-azurerm/issues/12764))
* **New Resource:** `azurerm_stream_analytics_cluster` ([#14082](https://github.com/hashicorp/terraform-provider-azurerm/issues/14082))
* **New Resource:** `azurerm_stream_analytics_managed_private_endpoint` ([#14082](https://github.com/hashicorp/terraform-provider-azurerm/issues/14082))

ENHANCEMENTS:

* dependencies: upgrading to `v0.18.0` of `github.com/hashicorp/go-azure-helpers` ([#14261](https://github.com/hashicorp/terraform-provider-azurerm/issues/14261))
* `azurerm_automation_rule` - support for the `expiration` property ([#14262](https://github.com/hashicorp/terraform-provider-azurerm/issues/14262))
* `azurerm_cosmosdb_account` - support for the `analytical_storage` and `capacity` blocks, `default_identity_type` and `storage_redundancy` properties ([#14346](https://github.com/hashicorp/terraform-provider-azurerm/issues/14346))
* `azurerm_eventgrid_event_subscription` - support the `queue_message_time_to_live_in_seconds` and `user_assigned_identity` properties ([#14318](https://github.com/hashicorp/terraform-provider-azurerm/issues/14318))
* `azurerm_firewall_policy` - allow cidr ranges for the `threat_intelligence_allowlist` property ([#14340](https://github.com/hashicorp/terraform-provider-azurerm/issues/14340))
* `azurerm_managed_disk` - support for the `public_network_access_enabled` property ([#14199](https://github.com/hashicorp/terraform-provider-azurerm/issues/14199))
* `azurerm_mssql_elasticpool` - support for the `DC` family ([#14270](https://github.com/hashicorp/terraform-provider-azurerm/issues/14270))
* `azurerm_mssql_server` - groundwork for the (currently disabled) 3.0 feature to set the default TLS version to 1.2 ([#14229](https://github.com/hashicorp/terraform-provider-azurerm/issues/14229))
* `azurerm_mysql_server` - groundwork for the (currently disabled) 3.0 feature to set the default TLS version to 1.2 ([#14229](https://github.com/hashicorp/terraform-provider-azurerm/issues/14229))
* `azurerm_orchestrated_virtual_machine_scale_set` - add extension support ([#14236](https://github.com/hashicorp/terraform-provider-azurerm/issues/14236))
* `azurerm_postgresql_server` - groundwork for the (currently disabled) 3.0 feature to set the default TLS version to 1.2 ([#14229](https://github.com/hashicorp/terraform-provider-azurerm/issues/14229))
* `azurerm_redis_cache` - groundwork for the (currently disabled) 3.0 feature to set the default TLS version to 1.2 ([#14229](https://github.com/hashicorp/terraform-provider-azurerm/issues/14229))
* `azurerm_service_plan` (beta) - add Logic App SKUs to validation. ([#14288](https://github.com/hashicorp/terraform-provider-azurerm/issues/14288))
* `azurerm_site_recovery_replication_policy` - now supports disabling of snapshots and their retention ([#14329](https://github.com/hashicorp/terraform-provider-azurerm/issues/14329))
* `azurerm_storage_account` - groundwork for the (currently disabled) 3.0 feature to set the default TLS version to 1.2 ([#14229](https://github.com/hashicorp/terraform-provider-azurerm/issues/14229))
* `azurerm_stream_analytics_job` - `compatibility_level` now accepts 1.2 ([#14294](https://github.com/hashicorp/terraform-provider-azurerm/issues/14294))

BUG FIXES: 

* `azurerm_function_app_slot` - fix a bug in `app_settings` for `WEBSITE_CONTENTSHARE` incorrectly updating ([#14211](https://github.com/hashicorp/terraform-provider-azurerm/issues/14211))
* `azurerm_monitor_diagnostic_setting` - Swap Relay parser and validator with EventHub ([#14277](https://github.com/hashicorp/terraform-provider-azurerm/issues/14277))
* `azurerm_stream_analytics_stream_input_eventhub` - correctly support creation with the default `eventhub_consumer_group_name` ([#14264](https://github.com/hashicorp/terraform-provider-azurerm/issues/14264))
* `azurerm_synapse_workspace` - fix a crash during update when `sql_aad_admin` was configured ([#14275](https://github.com/hashicorp/terraform-provider-azurerm/issues/14275))
* `azurerm_linux_virtual_machine` - `patch_mode` is now properly supported [GH0-14042]

## 2.86.0 (November 19, 2021)

FEATURES:

* **New Beta Resource:** `azurerm_linux_function_app` ([#13806](https://github.com/hashicorp/terraform-provider-azurerm/issues/13806))
* **New Resource:** `azurerm_automation_webhook` ([#13893](https://github.com/hashicorp/terraform-provider-azurerm/issues/13893))
* **New Resource:** `azurerm_resource_group_cost_management_export` ([#14140](https://github.com/hashicorp/terraform-provider-azurerm/issues/14140))
* **New Resource:** `azurerm_subscription_cost_management_export` ([#14140](https://github.com/hashicorp/terraform-provider-azurerm/issues/14140))
* **New Resource:** `azurerm_logz_tag_rule` ([#14020](https://github.com/hashicorp/terraform-provider-azurerm/issues/14020))
* **New Resource:** `azurerm_monitor_private_link_scoped_service` ([#14119](https://github.com/hashicorp/terraform-provider-azurerm/issues/14119))
* **New Resource:** `azurerm_storage_disks_pool` ([#14145](https://github.com/hashicorp/terraform-provider-azurerm/issues/14145))

ENHANCEMENTS:

* compute: updating to use API Version `2021-07-01` ([#14174](https://github.com/hashicorp/terraform-provider-azurerm/issues/14174))
* databricks: updating the embedded SDK to use the new Resource ID Parsers ([#14157](https://github.com/hashicorp/terraform-provider-azurerm/issues/14157))
* datalake: updating the embedded SDK to use the new Resource ID Parsers ([#14158](https://github.com/hashicorp/terraform-provider-azurerm/issues/14158))
* maps: updating the embedded SDK to use the new Resource ID Parsers ([#14155](https://github.com/hashicorp/terraform-provider-azurerm/issues/14155))
* powerbi: updating the embedded SDK to use the new Resource ID Parsers ([#14154](https://github.com/hashicorp/terraform-provider-azurerm/issues/14154))
* relay: updating the embedded SDK to use the new Resource ID Parsers ([#14153](https://github.com/hashicorp/terraform-provider-azurerm/issues/14153))
* signalr: updating the embedded SDK to use the new Resource ID Parsers ([#14150](https://github.com/hashicorp/terraform-provider-azurerm/issues/14150))
* storage: updating to use API Version `2021-04-01` ([#14083](https://github.com/hashicorp/terraform-provider-azurerm/issues/14083))
* videoanalyzer: updating the embedded SDK to use the new Resource ID Parsers ([#14135](https://github.com/hashicorp/terraform-provider-azurerm/issues/14135))
* Data Source: `azurerm_storage_account` - support for the `table_encryption_key_type` and `queue_encryption_key_type` attributes ([#14080](https://github.com/hashicorp/terraform-provider-azurerm/issues/14080))
* `azurerm_container_registry` - support for the `anonymous_pull_enabled`, `data_endpoint_enabled`, and `network_rule_bypass_option` properties ([#14096](https://github.com/hashicorp/terraform-provider-azurerm/issues/14096))
* `azurerm_cosmosdb_cassandra_datacenter` - support the `availabilit_zones_enabled` property ([#14235](https://github.com/hashicorp/terraform-provider-azurerm/issues/14235))
* `azurerm_cost_management_export_resource_group` - has been deprecated in favour of the `azurerm_resource_group_cost_management_export` resource ([#14140](https://github.com/hashicorp/terraform-provider-azurerm/issues/14140))
* `azurerm_disk_encryption_set` - add support for the `encryption_type` property ([#14218](https://github.com/hashicorp/terraform-provider-azurerm/issues/14218))
* `azurerm_elastic_pool` - support for the `Fsv2` family SKUs ([#14250](https://github.com/hashicorp/terraform-provider-azurerm/issues/14250))
* `azurerm_key_vault_certificate` - groundwork for the (currently disabled) 3.0 feature to support more granular configuration of soft-delete and purge protection ([#13682](https://github.com/hashicorp/terraform-provider-azurerm/issues/13682))
* `azurerm_key_vault_key` - groundwork for the (currently disabled) 3.0 feature to support more granular configuration of soft-delete and purge protection ([#13682](https://github.com/hashicorp/terraform-provider-azurerm/issues/13682))
* `azurerm_key_vault_secret` - groundwork for the (currently disabled) 3.0 feature to support more granular configuration of soft-delete and purge protection ([#13682](https://github.com/hashicorp/terraform-provider-azurerm/issues/13682))
* `azurerm_key_vault_certificate` - the `certificate_policy` property is now optional for imported certificates ([#14225](https://github.com/hashicorp/terraform-provider-azurerm/issues/14225))
* `azurerm_kubernetes_cluster` - support for `outbound_type` = `*NATGateway` and the `nat_gateway_profile` block ([#14142](https://github.com/hashicorp/terraform-provider-azurerm/issues/14142))
* `azurerm_linux_web_app` - (Beta) add support for `health_check_eviction_time_in_mins` and `vnet_route_all_enabled` ([#14202](https://github.com/hashicorp/terraform-provider-azurerm/issues/14202))
* `azurerm_managed_disk` - support for the `on_demand_bursting_enabled` property ([#14137](https://github.com/hashicorp/terraform-provider-azurerm/issues/14137))
* `azurerm_mssql_server` - support for the `azuread_authentication_only` property on creation ([#14169](https://github.com/hashicorp/terraform-provider-azurerm/issues/14169))
* `azurerm_machine_learning_workspace` - support for the `encryption` block ([#14120](https://github.com/hashicorp/terraform-provider-azurerm/issues/14120))
* `azurerm_orchestrated_virtual_machine_scale_set` - added support for VMSS Flex public preview ([#14003](https://github.com/hashicorp/terraform-provider-azurerm/issues/14003))
* `azurerm_postgresql_flexible_server` - the `zone` and `standby_availability_zone` properties are no longer computed ([#13843](https://github.com/hashicorp/terraform-provider-azurerm/issues/13843))
* `azurerm_public_ip_prefix` - support for the `ip_version` property ([#14228](https://github.com/hashicorp/terraform-provider-azurerm/issues/14228))
* `azurerm_purview_account` - support for the `managed_resource_group_name` property ([#14217](https://github.com/hashicorp/terraform-provider-azurerm/issues/14217))
* `azurerm_resource_provider_registration` - support for managing `features` ([#12385](https://github.com/hashicorp/terraform-provider-azurerm/issues/12385))
* `azurerm_windows_virtual_machine` - support for the `vtpm_enabled` and `secure_boot_enabled` properties ([#13713](https://github.com/hashicorp/terraform-provider-azurerm/issues/13713))
* `azurerm_windows_virtual_machine_scale_set` - support for the `vtpm_enabled` and `secure_boot_enabled` properties ([#13713](https://github.com/hashicorp/terraform-provider-azurerm/issues/13713))
* `azurerm_windows_web_app` - (Beta) add support for the `health_check_eviction_time_in_mins` and `vnet_route_all_enabled` properties ([#14202](https://github.com/hashicorp/terraform-provider-azurerm/issues/14202))
* `azurerm_stream_analytics_output_servicebus_topic` - support for the `property_columns` property ([#14252](https://github.com/hashicorp/terraform-provider-azurerm/issues/14252))
* `azurerm_storage_account` - support for `table_encryption_key_type` and `queue_encryption_key_type` properties ([#14080](https://github.com/hashicorp/terraform-provider-azurerm/issues/14080))
* `azurerm_storage_account` - (Beta) add a state migration for the renaming of `allow_blob_public_access` to `allow_nested_items_to_be_public` ([#13607](https://github.com/hashicorp/terraform-provider-azurerm/issues/13607))
* `azurerm_sql_active_directory_administrator` - support for the `azuread_authentication_only` property ([#14172](https://github.com/hashicorp/terraform-provider-azurerm/issues/14172))
* `azurerm_virtual_network` - support for the `flow_timeout_in_minutes` property ([#14200](https://github.com/hashicorp/terraform-provider-azurerm/issues/14200))
* `azurerm_virtual_desktop_application_group` - support for the `default_desktop_display_name` property ([#14227](https://github.com/hashicorp/terraform-provider-azurerm/issues/14227))

BUG FIXES: 

* `azurerm_backup_protected_file_share` - correctly list file shares that are added to an existing storage account not returned by the Backup Protectable Items API ([#14238](https://github.com/hashicorp/terraform-provider-azurerm/issues/14238))
* `azurerm_frontdoor` - validation for `probe_method` allows the default value ([#14204](https://github.com/hashicorp/terraform-provider-azurerm/issues/14204))
* `azurerm_key_vault_managed_hardware_security_module` - extend context timeouts for creation and deletion ([#14253](https://github.com/hashicorp/terraform-provider-azurerm/issues/14253))
* `azurerm_key_vault_certificate` - changing the `tags` property no longer forces a new resource to be created ([#14079](https://github.com/hashicorp/terraform-provider-azurerm/issues/14079))
* `azurerm_linux_virtual_machine_scale_set` - changing the `source_image_reference.offer` and `source_image_reference.publisher` now creates a new resource ([#14165](https://github.com/hashicorp/terraform-provider-azurerm/issues/14165))
* `azurerm_mssql_database` - corrert an error when using `OnlineSecondary` with auditing on the primary database ([#14192](https://github.com/hashicorp/terraform-provider-azurerm/issues/14192))
* `azurerm_network_watcher_flow_log` - now locks on the network security group to prevent `AnotherOperationInProgress` errors ([#14160](https://github.com/hashicorp/terraform-provider-azurerm/issues/14160))
* `azurerm_windows_virtual_machine_scale_set` - `source_image_reference.offer` and `source_image_reference.publisher` are now ForceNew ([#14165](https://github.com/hashicorp/terraform-provider-azurerm/issues/14165))

## 2.85.0 (November 12, 2021)

FEATURES:

* **New Data Source:** `azurerm_batch_application` ([#14043](https://github.com/hashicorp/terraform-provider-azurerm/issues/14043))
* **New Resource:** `azurerm_monitor_private_link_scope` ([#14098](https://github.com/hashicorp/terraform-provider-azurerm/issues/14098))
* **New Resource:** `azurerm_mysql_flexible_server_firewall_rule` ([#14136](https://github.com/hashicorp/terraform-provider-azurerm/issues/14136))
* **New Resource:** `azurerm_synapse_workspace_aad_admin` ([#13600](https://github.com/hashicorp/terraform-provider-azurerm/issues/13600))

IMPROVEMENTS:

* dependencies: upgrading to `v0.17.1` of `github.com/hashicorp/go-azure-helpers` ([#14141](https://github.com/hashicorp/terraform-provider-azurerm/issues/14141))
* dependencies: upgrading to `v2.8.0` of `github.com/hashicorp/terraform-plugin-sdk` ([#14060](https://github.com/hashicorp/terraform-provider-azurerm/issues/14060))
* `azurerm_application_insights` - support for the `internet_ingestion_enabled` and `internet_query_enabled` properties ([#14035](https://github.com/hashicorp/terraform-provider-azurerm/issues/14035))
* `azurerm_backup_protected_vm` - support for the `exclude_disk_luns` and `include_disk_luns` properties ([#14097](https://github.com/hashicorp/terraform-provider-azurerm/issues/14097))
* `azurerm_managed_disk_resource` - support for the `disk_iops_read_only` and `disk_mbps_read_only` properties ([#14025](https://github.com/hashicorp/terraform-provider-azurerm/issues/14025))
* `azurerm_security_center_subscription_pricing` - `resource_type` can now be set to `OpenSourceRelationalDatabases` ([#14103](https://github.com/hashicorp/terraform-provider-azurerm/issues/14103))
* `azurerm_storage_encryption_scope` - allow versionless `key_vault_key_id` ([#14085](https://github.com/hashicorp/terraform-provider-azurerm/issues/14085))
* `azurerm_sql_managed_instance` - support for the `identity` block ([#14052](https://github.com/hashicorp/terraform-provider-azurerm/issues/14052))
* `azurerm_virtual_network_gateway` - enable configuration of an active-active zone redundant gateway with P2S ([#14124](https://github.com/hashicorp/terraform-provider-azurerm/issues/14124))

BUG FIXES:

* Data Source: `azurerm_redis_cache` - parsing the `subnet_id` response value case-insensitively ([#14108](https://github.com/hashicorp/terraform-provider-azurerm/issues/14108))
* Data Source: `azurerm_redis_cache` - ensuring that `shard_count` always has a value set ([#14108](https://github.com/hashicorp/terraform-provider-azurerm/issues/14108))
* Data Source: `azurerm_consumption_budget_resource_group` - add missing `threshold_type` property in the schema ([#14125](https://github.com/hashicorp/terraform-provider-azurerm/issues/14125))
* Data Source: `azurerm_consumption_budget_subscription` - add missing `threshold_type` property in the schema ([#14125](https://github.com/hashicorp/terraform-provider-azurerm/issues/14125))
* `azurerm_api_management_certificate` - set `subject` property from correct field ([#14026](https://github.com/hashicorp/terraform-provider-azurerm/issues/14026))
* `azurerm_app_service_virtual_network_swift_connection` - fixing a panic when checking for an existing resource during creation ([#14070](https://github.com/hashicorp/terraform-provider-azurerm/issues/14070))
* `azurerm_frontdoor_resource` - route engines are no longer removed on update ([#14093](https://github.com/hashicorp/terraform-provider-azurerm/issues/14093))
* `azurerm_redis_cache` - parsing the `subnet_id` response value case-insensitively ([#14108](https://github.com/hashicorp/terraform-provider-azurerm/issues/14108))
* `azurerm_redis_cache` - ensuring that `shard_count` always has a value set ([#14108](https://github.com/hashicorp/terraform-provider-azurerm/issues/14108))
* `azurerm_storage_blob` - ensuring that `cache_control` is sent during updates ([#14100](https://github.com/hashicorp/terraform-provider-azurerm/issues/14100))

## 2.84.0 (November 05, 2021)

FEATURES:

* **New Resource:** `azurerm_cosmosdb_cassandra_cluster` ([#14019](https://github.com/hashicorp/terraform-provider-azurerm/issues/14019))
* **New Resource:** `azurerm_cosmosdb_cassandra_datacenter` ([#14019](https://github.com/hashicorp/terraform-provider-azurerm/issues/14019))
* **New Resource:** `logz_monitor` ([#13874](https://github.com/hashicorp/terraform-provider-azurerm/issues/13874))
* **New Resource:** `azurerm_stream_analytics_output_synapse` ([#14013](https://github.com/hashicorp/terraform-provider-azurerm/issues/14013))

IMPROVEMENTS:

* upgrading `cosmos` to API Version `2021-10-15` ([#13785](https://github.com/hashicorp/terraform-provider-azurerm/issues/13785))
* upgrading `aks` to API Version `2021-08-01` ([#13465](https://github.com/hashicorp/terraform-provider-azurerm/issues/13465))
* upgrading `purview` to API Version `2021-07-01` ([#13785](https://github.com/hashicorp/terraform-provider-azurerm/issues/13785))
* Data Source: `azurerm_key_vault_key` - export the `cureve`, `x`, `y`, `public_key_pem`, and `public_key_openssh` attributes ([#13934](https://github.com/hashicorp/terraform-provider-azurerm/issues/13934))
* `azurerm_app_service_slot` - support for the `key_vault_reference_identity_id` property  ([#13988](https://github.com/hashicorp/terraform-provider-azurerm/issues/13988))
* `azurerm_cosmosdb_account` - the backup backup type can now be changed from `Periodic` to `Continuous` without creating a new resource ([#13967](https://github.com/hashicorp/terraform-provider-azurerm/issues/13967))
* `azurerm_firewall_policy_rule_collection_group` - support for the `translated_fqdn` property ([#13976](https://github.com/hashicorp/terraform-provider-azurerm/issues/13976))
* `azurerm_firewall_policy` - support for the `insights` block ([#14004](https://github.com/hashicorp/terraform-provider-azurerm/issues/14004))
* `azurerm_logic_app_integration_account` - support the `integration_service_environment_id` property ([#14015](https://github.com/hashicorp/terraform-provider-azurerm/issues/14015))
* `azurerm_function_app` - support for the `key_vault_reference_identity_id` property ([#13962](https://github.com/hashicorp/terraform-provider-azurerm/issues/13962))
* `azurerm_key_vault_key` - support for the `public_key_pem` and `public_key_openssh` attributes ([#13934](https://github.com/hashicorp/terraform-provider-azurerm/issues/13934))
* `azurerm_linux_virtual_machine` - support for the `patch_mode` property ([#13866](https://github.com/hashicorp/terraform-provider-azurerm/issues/13866))
* `azurerm_machine_learning_compute_cluster` - support for the `local_auth_enabled` property ([#13820](https://github.com/hashicorp/terraform-provider-azurerm/issues/13820))
* `azurerm_machine_learning_compute_cluster` - support for the `local_auth_enabled` property ([#13820](https://github.com/hashicorp/terraform-provider-azurerm/issues/13820))
* `azurerm_machine_learning_synapse_spark` - support for the `local_auth_enabled` property ([#13820](https://github.com/hashicorp/terraform-provider-azurerm/issues/13820))
* `azurerm_monitor_smart_detector_alert_rule` - support additional detector types ([#13998](https://github.com/hashicorp/terraform-provider-azurerm/issues/13998))
* `azurerm_mssql_elasticpool` - support `GP_FSv2` for the `sku` property ([#13973](https://github.com/hashicorp/terraform-provider-azurerm/issues/13973))
* `azurerm_synapse_workspace` - supports for the `sql_aad_admin` block ([#13659](https://github.com/hashicorp/terraform-provider-azurerm/issues/13659))
* `azurerm_sql_managed_instance` - support for the `dns_zone_partner_id` property ([#13951](https://github.com/hashicorp/terraform-provider-azurerm/issues/13951))
* `azurerm_storage_blob` - support for the `cache_control` property ([#13946](https://github.com/hashicorp/terraform-provider-azurerm/issues/13946))
* `azurerm_storage_share` - support for the `enabled_protocol` property ([#13938](https://github.com/hashicorp/terraform-provider-azurerm/issues/13938))

BUG FIXES:

* `azurerm_application_insights` - correct vlaidation for the `daily_data_cap_in_gb` property ([#13971](https://github.com/hashicorp/terraform-provider-azurerm/issues/13971))
* `azurerm_logic_app_standard` - will no longer error when working on private networks ([#13964](https://github.com/hashicorp/terraform-provider-azurerm/issues/13964))
* `azurerm_managed_disk_resource` - the validation for the `disk_iops_read_write` and `disk_mbps_read_write` properties ensures values greater then 0 ([#14028](https://github.com/hashicorp/terraform-provider-azurerm/issues/14028))
* `azurerm_purview_account` - deprecate the `sku_name` property ([#13897](https://github.com/hashicorp/terraform-provider-azurerm/issues/13897))
* `azurerm_synapse_workspace_key` - deprecated the `cusomter_managed_key_name` property in favour of the correctly spelled `customer_managed_key_name` one ([#13881](https://github.com/hashicorp/terraform-provider-azurerm/issues/13881))

## 2.83.0 (October 29, 2021)

FEATURES:

* **New Data Source:** `azurerm_eventgrid_system_topic` ([#13851](https://github.com/hashicorp/terraform-provider-azurerm/issues/13851))
* **New Data Source:** `azurerm_billing_mpa_account_scope` ([#13723](https://github.com/hashicorp/terraform-provider-azurerm/issues/13723))
* **New Resource:** `azurerm_kusto_script` ([#13692](https://github.com/hashicorp/terraform-provider-azurerm/issues/13692))
* **New Resource:** `azurerm_iot_time_series_insights_event_source_eventhub` ([#13917](https://github.com/hashicorp/terraform-provider-azurerm/issues/13917))
* **New Resource:** `azurerm_stream_analytics_reference_input_mssql` ([#13822](https://github.com/hashicorp/terraform-provider-azurerm/issues/13822))
* **New Resource:** `azurerm_sentinel_automation_rule` ([#11502](https://github.com/hashicorp/terraform-provider-azurerm/issues/11502))
* **New Resource:** `azurerm_stream_analytics_output_table` ([#13854](https://github.com/hashicorp/terraform-provider-azurerm/issues/13854))

IMPROVEMENTS:

* upgrading `mysql` to API Version `2021-05-01` ([#13818](https://github.com/hashicorp/terraform-provider-azurerm/issues/13818))
* `azurerm_application_gateway` - support for the `priority` property ([#13498](https://github.com/hashicorp/terraform-provider-azurerm/issues/13498))
* `azurerm_firewall_application_rule_collection` - the `port` property is now required instead of optional ([#13869](https://github.com/hashicorp/terraform-provider-azurerm/issues/13869))
* `azurerm_kubernetes_cluster` - expose the `portal_fqdn` attribute ([#13887](https://github.com/hashicorp/terraform-provider-azurerm/issues/13887))
* `azurerm_linux_virtual_machine_scale_set` -  support for `automatic_upgrade_enabled` in extensions ([#13394](https://github.com/hashicorp/terraform-provider-azurerm/issues/13394))
* `azurerm_linux_virtual_machine_scale_set` - added feature for `scale_to_zero_before_deletion`([#13635](https://github.com/hashicorp/terraform-provider-azurerm/issues/13635))
* `azurerm_managed_disk` - support for the `trusted_launch_enabled` property ([#13849](https://github.com/hashicorp/terraform-provider-azurerm/issues/13849))
* `azurerm_postgres_flexible_server` - enhanced validation for the `administrator_login` property ([#13942](https://github.com/hashicorp/terraform-provider-azurerm/issues/13942))
* `azurerm_servicebus_queue` - support for the `max_message_size_in_kilobytes` property ([#13762](https://github.com/hashicorp/terraform-provider-azurerm/issues/13762))
* `azurerm_servicebus_topic` - support for the `max_message_size_in_kilobytes` property ([#13762](https://github.com/hashicorp/terraform-provider-azurerm/issues/13762))
* `azurerm_servicebus_namespace_network_rule_set` - support for the `trusted_services_allowed` property ([#13853](https://github.com/hashicorp/terraform-provider-azurerm/issues/13853))
* `azurerm_windows_virtual_machine_scale_set` - added feature for `scale_to_zero_before_deletion`([#13635](https://github.com/hashicorp/terraform-provider-azurerm/issues/13635))
* `azurerm_synapse_workspace` - support for the `linking_allowed_for_aad_tenant_ids`, `compute_subnet_id`, `public_network_access_enabled`, `purview_id`, and `last_commit_id` properties  ([#13817](https://github.com/hashicorp/terraform-provider-azurerm/issues/13817))
* `azurerm_spring_cloud_java_deployment` – the `cpu` and `memory_in_gb` properties have been deprecated in favour of the `quota` block ([#12924](https://github.com/hashicorp/terraform-provider-azurerm/issues/12924))
* `azurerm_vpn_gateway` - support for the `routing_preference` property ([#13882](https://github.com/hashicorp/terraform-provider-azurerm/issues/13882))
* `azurerm_virtual_hub` - support for the `default_route_table_id` property ([#13840](https://github.com/hashicorp/terraform-provider-azurerm/issues/13840))
* `azurerm_virtual_machine_scale_set_extension ` - support for `automatic_upgrade_enabled` ([#13394](https://github.com/hashicorp/terraform-provider-azurerm/issues/13394))
* `azurerm_windows_virtual_machine_scale_set` -  support for `automatic_upgrade_enabled` in extensions ([#13394](https://github.com/hashicorp/terraform-provider-azurerm/issues/13394))

BUG FIXES:

* `azurerm_automation_schedule_resource` - allow `Etc/UTC` for the `timezone` property ([#13906](https://github.com/hashicorp/terraform-provider-azurerm/issues/13906))
* `azurerm_app_configuration_key` - now supports forward slashes in the `key` ([#13859](https://github.com/hashicorp/terraform-provider-azurerm/issues/13859))
* `azurerm_application_gateway` - prevent multiple `ssl_policy` blocks ([#13929](https://github.com/hashicorp/terraform-provider-azurerm/issues/13929))
* `azurerm_cosmosdb_account` - the `capabilities` property is now computed ([#13936](https://github.com/hashicorp/terraform-provider-azurerm/issues/13936))
* `azurerm_cognitive_account` - will now handle the unexpected state `Accepted` when waiting for creats ([#13925](https://github.com/hashicorp/terraform-provider-azurerm/issues/13925))
* `azurerm_data_factory` - can now read global parameter values ([#13519](https://github.com/hashicorp/terraform-provider-azurerm/issues/13519))
* `azurerm_firewall_policy` - will now correctly import ([#13862](https://github.com/hashicorp/terraform-provider-azurerm/issues/13862))
* `azurerm_firewall_policy` - changing the identity will no longer create a new resource ([#13904](https://github.com/hashicorp/terraform-provider-azurerm/issues/13904))

## 2.82.0 (October 21, 2021)

FEATURES: 

* **New Resource:** `azurerm_mysql_flexible_server_configuration` ([#13831](https://github.com/hashicorp/terraform-provider-azurerm/issues/13831))
* **New Resource:** `azurerm_synapse_sql_pool_vulnerability_assessment_baseline` ([#13744](https://github.com/hashicorp/terraform-provider-azurerm/issues/13744))
* **New Resource:** `azurerm_virtual_hub_route_table_route` ([#13743](https://github.com/hashicorp/terraform-provider-azurerm/issues/13743))

IMPROVEMENTS:

* dependencies: upgrading to `v58.0.0` of `github.com/Azure/azure-sdk-for-go` ([#13613](https://github.com/hashicorp/terraform-provider-azurerm/issues/13613))
* upgrading `netapp` to API Version `2021-06-01` ([#13812](https://github.com/hashicorp/terraform-provider-azurerm/issues/13812))
* upgrading `servicebus` to API Version `2021-06-01-preview` ([#13701](https://github.com/hashicorp/terraform-provider-azurerm/issues/13701))
* Data Source: `azurerm_disk_encryption_set` - support for the `auto_key_rotation_enabled` property ([#13747](https://github.com/hashicorp/terraform-provider-azurerm/issues/13747))
* Data Source: `azurerm_virtual_machine` - expose IP addresses as data source outputs ([#13773](https://github.com/hashicorp/terraform-provider-azurerm/issues/13773))
* `azurerm_batch_account` - support for the `identity` block ([#13742](https://github.com/hashicorp/terraform-provider-azurerm/issues/13742))
* `azurerm_batch_pool` - support for the `identity` block ([#13779](https://github.com/hashicorp/terraform-provider-azurerm/issues/13779))
* `azurerm_container_registry` - supports for the `regiononal_endpoint_enabled` property ([#13767](https://github.com/hashicorp/terraform-provider-azurerm/issues/13767))
* `azurerm_data_factory_integration_runtime_azure` - support `AutoResolve` for the `location` property ([#13731](https://github.com/hashicorp/terraform-provider-azurerm/issues/13731))
* `azurerm_disk_encryption_set` - support for the `auto_key_rotation_enabled` property ([#13747](https://github.com/hashicorp/terraform-provider-azurerm/issues/13747))
* `azurerm_iot_security_solution` - support for the `additional_workspace` and `disabled_data_sources` properties ([#13783](https://github.com/hashicorp/terraform-provider-azurerm/issues/13783))
* `azurerm_kubernetes_cluster` - support for the `open_service_mesh` block ([#13462](https://github.com/hashicorp/terraform-provider-azurerm/issues/13462))
* `azurerm_lb` - support for the `gateway_load_balancer_frontend_ip_configuration_id` property ([#13559](https://github.com/hashicorp/terraform-provider-azurerm/issues/13559))
* `azurerm_lb_backend_address_pool` - support for the `tunnel_interface` block ([#13559](https://github.com/hashicorp/terraform-provider-azurerm/issues/13559))
* `azurerm_lb_rule` - the `backend_address_pool_ids` property has been deprecated in favour of the `backend_address_pool_ids` property ([#13559](https://github.com/hashicorp/terraform-provider-azurerm/issues/13559))
* `azurerm_lb_nat_pool` - support for the `floating_ip_enabled`, `tcp_reset_enabled`, and `idle_timeout_in_minutes` properties ([#13674](https://github.com/hashicorp/terraform-provider-azurerm/issues/13674))
* `azurerm_mssql_server` - support for the `azuread_authentication_only` property ([#13754](https://github.com/hashicorp/terraform-provider-azurerm/issues/13754))
* `azurerm_network_interface` - support for the `gateway_load_balancer_frontend_ip_configuration_id` property ([#13559](https://github.com/hashicorp/terraform-provider-azurerm/issues/13559))
* `azurerm_synapse_spark_pool` - support for the `cache_size`, `compute_isolation_enabled`, `dynamic_executor_allocation_enabled`, `session_level_packages_enabled` and `spark_config` properties ([#13690](https://github.com/hashicorp/terraform-provider-azurerm/issues/13690))

BUG FIXES:

* `azurerm_app_configuration_feature` - fix default value handling for percentage appconfig feature filters. ([#13771](https://github.com/hashicorp/terraform-provider-azurerm/issues/13771))
* `azurerm_cosmosdb_account` - force `MongoEnabled` feature when enabling `MongoDBv3.4`. ([#13757](https://github.com/hashicorp/terraform-provider-azurerm/issues/13757))
* `azurerm_mssql_server` - will now configure the `azuread_administrator` during resource creation ([#13753](https://github.com/hashicorp/terraform-provider-azurerm/issues/13753))
* `azurerm_mssql_database` - fix failure by preventing `extended_auditing_policy` from being configured for secondaries ([#13799](https://github.com/hashicorp/terraform-provider-azurerm/issues/13799))
* `azurerm_postgresql_flexible_server` - changing the `standby_availability_zone` no longer forces a new resource ([#13507](https://github.com/hashicorp/terraform-provider-azurerm/issues/13507))
* `azurerm_servicebus_subscription` - the `name` field can now start & end with an underscore ([#13797](https://github.com/hashicorp/terraform-provider-azurerm/issues/13797))

## 2.81.0 (October 14, 2021)

FEATURES: 

* **New Data Source:** `azurerm_consumption_budget_resource_group` ([#12538](https://github.com/hashicorp/terraform-provider-azurerm/issues/12538))
* **New Data Source:** `azurerm_consumption_budget_subscription` ([#12540](https://github.com/hashicorp/terraform-provider-azurerm/issues/12540))
* **New Resource:** `azurerm_data_factory_linked_service_cosmosdb_mongoapi` ([#13636](https://github.com/hashicorp/terraform-provider-azurerm/issues/13636))
* **New Resource:** `azurerm_mysql_flexible_server` ([#13678](https://github.com/hashicorp/terraform-provider-azurerm/issues/13678))

IMPROVEMENTS:

* upgrading `batch` to API Version `2021-06-01`([#13718](https://github.com/hashicorp/terraform-provider-azurerm/issues/13718))
* upgrading `mssql` to API Version `v5.0`([#13622](https://github.com/hashicorp/terraform-provider-azurerm/issues/13622))
* Data Source: `azurerm_key_vault` - exports the `enable_rbac_authorization` attribute ([#13717](https://github.com/hashicorp/terraform-provider-azurerm/issues/13717))
* `azurerm_app_service` - support for the `key_vault_reference_identity_id` property ([#13720](https://github.com/hashicorp/terraform-provider-azurerm/issues/13720))
* `azurerm_lb` - support for the `sku_tier` property ([#13680](https://github.com/hashicorp/terraform-provider-azurerm/issues/13680))
* `azurerm_eventgrid_event_subscription` - support the `delivery_property` block ([#13595](https://github.com/hashicorp/terraform-provider-azurerm/issues/13595))
* `azurerm_mssql_server` - support for the `user_assigned_identity_ids` and `primary_user_assigned_identity_id` properties ([#13683](https://github.com/hashicorp/terraform-provider-azurerm/issues/13683))
* `azurerm_network_connection_monitor` - add support for the `destination_port_behavior` property ([#13518](https://github.com/hashicorp/terraform-provider-azurerm/issues/13518))
* `azurerm_security_center_workspace` - now supports the `Free` pricing tier ([#13710](https://github.com/hashicorp/terraform-provider-azurerm/issues/13710))
* `azurerm_kusto_attached_database_configuration` - support for the `sharing` property ([#13487](https://github.com/hashicorp/terraform-provider-azurerm/issues/13487))

BUG FIXES:

* Data Source: `azurerm_cosmosdb_account`- prevent a panic from an index out of range error ([#13560](https://github.com/hashicorp/terraform-provider-azurerm/issues/13560))
* `azurerm_function_app_slot` - the `client_affinity` property has been deprecated as it is no longer configurable in the service's API ([#13711](https://github.com/hashicorp/terraform-provider-azurerm/issues/13711))
* `azurerm_kubernetes_cluster` - the `kube_config` and `kube_admin_config` blocks can now be marked entirely as `Sensitive` via an environment variable ([#13732](https://github.com/hashicorp/terraform-provider-azurerm/issues/13732))
* `azurerm_logic_app_workflow` - will not check for `nil` and empty access control properties ([#13689](https://github.com/hashicorp/terraform-provider-azurerm/issues/13689))
* `azurerm_management_group` - will not nil check child management groups when deassociating a subscription from a management group ([#13540](https://github.com/hashicorp/terraform-provider-azurerm/issues/13540))
* `azurerm_subnet_resource` - will now lock the virtual network and subnet on updates ([#13726](https://github.com/hashicorp/terraform-provider-azurerm/issues/13726))
* `azurerm_app_configuration_key` - can now mix labeled and unlabeled keys ([#13736](https://github.com/hashicorp/terraform-provider-azurerm/issues/13736))
 
## 2.80.0 (October 08, 2021)

FEATURES: 

* **New Data Source:** `backup_policy_file_share` ([#13444](https://github.com/hashicorp/terraform-provider-azurerm/issues/13444))

IMPROVEMENTS:

* Data Source `azurerm_public_ips` - deprecate the `attached` property infavour of the `attachment_status` property to improve filtering ([#13500](https://github.com/hashicorp/terraform-provider-azurerm/issues/13500))
* Data Source `azurerm_public_ips` - return public IPs associated with NAT gateways when `attached` set to `true` or `attachment_status` set to `Attached` ([#13610](https://github.com/hashicorp/terraform-provider-azurerm/issues/13610))
* `azurerm_kusto_eventhub_data_connection supports` - support for the `identity_id` property ([#13488](https://github.com/hashicorp/terraform-provider-azurerm/issues/13488))
* `azurerm_managed_disk` - support for the `logical_sector_size` property ([#13637](https://github.com/hashicorp/terraform-provider-azurerm/issues/13637))
* `azurerm_service_fabric_cluster` - support for the `service_fabric_zonal_upgrade_mode` and `service_fabric_zonal_upgrade_mode` properties ([#13399](https://github.com/hashicorp/terraform-provider-azurerm/issues/13399))
* `azurerm_stream_analytics_output_eventhub` - support for the `partition_key` property ([#13562](https://github.com/hashicorp/terraform-provider-azurerm/issues/13562))
* `azurerm_linux_virtual_machine_scale_set` - correctly update the `overprovision` property ([#13653](https://github.com/hashicorp/terraform-provider-azurerm/issues/13653))

BUG FIXES:

* `azurerm_function_app` - fix regressions in function app storage introduced in v2.77 ([#13580](https://github.com/hashicorp/terraform-provider-azurerm/issues/13580))
* `azurerm_managed_application` - fixed typecasting bug ([#13641](https://github.com/hashicorp/terraform-provider-azurerm/issues/13641))

---

For information on changes between the v2.69.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v2.00.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
