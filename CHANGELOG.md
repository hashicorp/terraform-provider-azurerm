## 3.86.0 (Unreleased)

FEATURES:

* New Resource: `azurerm_log_analytics_workspace_table` [GH-24229]

ENHANCEMENTS:

* dependencies: updating to `v0.20231215.1114251` of `hashicorp/go-azure-sdk` [GH-24251]
* `azurerm_logic_app_standard` - add support for `site_config.0.public_network_access_enabled` [GH-24257]

BUG FIXES:

* `azurerm_bot_service` - fixing a bug where `public_network_access_enabled` was always set to `true` in Read func [GH-24255]

## 3.85.0 (December 14, 2023)

FEATURES:

* New Data Source: `azurerm_locations` ([#23324](https://github.com/hashicorp/terraform-provider-azurerm/issues/23324))
* New Resource: `azurerm_iotcentral_organization` ([#23132](https://github.com/hashicorp/terraform-provider-azurerm/issues/23132))

ENHANCEMENTS:

* provider: support for authenticating using Azure Kubernetes Service Workload Identity ([#23965](https://github.com/hashicorp/terraform-provider-azurerm/issues/23965))
* dependencies: updating to `v0.65.0` of `github.com/hashicorp/go-azure-helpers` ([#24222](https://github.com/hashicorp/terraform-provider-azurerm/issues/24222))
* dependencies: updating to `v0.20231214.1220802` of `github.com/hashicorp/go-azure-sdk` ([#24246](https://github.com/hashicorp/terraform-provider-azurerm/issues/24246))
* dependencies: updating to version `v0.20231214.1160726` of `github.com/hashicorp/go-azure-sdk` ([#24241](https://github.com/hashicorp/terraform-provider-azurerm/issues/24241))
* dependencies: update `security/automation` to use `hashicorp/go-azure-sdk` ([#24156](https://github.com/hashicorp/terraform-provider-azurerm/issues/24156))
* `dataprotection`: updating to API Version `2023-05-01` ([#24143](https://github.com/hashicorp/terraform-provider-azurerm/issues/24143))
* `kusto`: removing the remnants of the old Resource ID Parsers now this uses `hashicorp/go-azure-sdk` ([#24238](https://github.com/hashicorp/terraform-provider-azurerm/issues/24238))
* Data Source: `azurerm_cognitive_account` - export the `identity` block ([#24214](https://github.com/hashicorp/terraform-provider-azurerm/issues/24214))
* Data Source: `azurerm_monitor_workspace` - add support for the `default_data_collection_endpoint_id` and `default_data_collection_rule_id` properties ([#24153](https://github.com/hashicorp/terraform-provider-azurerm/issues/24153))
* Data Source: `azurerm_shared_image_gallery` - add support for the `image_names` property ([#24176](https://github.com/hashicorp/terraform-provider-azurerm/issues/24176))
* `azurerm_dns_txt_record` - allow up to `4096` characters for the property `record.value` ([#24169](https://github.com/hashicorp/terraform-provider-azurerm/issues/24169))
* `azurerm_container_app` - support for the `workload_profile_name` property ([#24219](https://github.com/hashicorp/terraform-provider-azurerm/issues/24219))
* `azurerm_container_app` - suppot for the `init_container` block ([#23955](https://github.com/hashicorp/terraform-provider-azurerm/issues/23955))
* `azurerm_hpc_cache_blob_nfs_target` - support for the `verification_timer_in_seconds` and `write_back_timer_in_seconds` properties ([#24207](https://github.com/hashicorp/terraform-provider-azurerm/issues/24207))
* `azurerm_hpc_cache_nfs_target` - support for the `verification_timer_in_seconds` and `write_back_timer_in_seconds` properties ([#24208](https://github.com/hashicorp/terraform-provider-azurerm/issues/24208))
* `azurerm_linux_web_app` - make `client_secret_setting_name` optional and conflict with `client_secret_certificate_thumbprint` ([#21834](https://github.com/hashicorp/terraform-provider-azurerm/issues/21834))
* `azurerm_linux_web_app_slot` - make `client_secret_setting_name` optional and conflict with `client_secret_certificate_thumbprint` ([#21834](https://github.com/hashicorp/terraform-provider-azurerm/issues/21834))
* `azurerm_linux_web_app` - fix a bug in `app_settings` where settings could be lost ([#24221](https://github.com/hashicorp/terraform-provider-azurerm/issues/24221))
* `azurerm_linux_web_app_slot` - fix a bug in `app_settings` where settings could be lost ([#24221](https://github.com/hashicorp/terraform-provider-azurerm/issues/24221))
* `azurerm_log_analytics_workspace` - add support for the `immediate_data_purge_on_30_days_enabled` property ([#24015](https://github.com/hashicorp/terraform-provider-azurerm/issues/24015))
* `azurerm_mssql_server` - support for other identity types for the key vault key ([#24236](https://github.com/hashicorp/terraform-provider-azurerm/issues/24236))
* `azurerm_machine_learning_datastore_blobstorage` - resource now skips validation when being created ([#24078](https://github.com/hashicorp/terraform-provider-azurerm/issues/24078))
* `azurerm_machine_learning_datastore_datalake_gen2` - resource now skips validation when being created ([#24078](https://github.com/hashicorp/terraform-provider-azurerm/issues/24078))
* `azurerm_machine_learning_datastore_fileshare` - resource now skips validation when being created ([#24078](https://github.com/hashicorp/terraform-provider-azurerm/issues/24078))
* `azurerm_monitor_workspace` - support for the `default_data_collection_endpoint_id` and `default_data_collection_rule_id` properties ([#24153](https://github.com/hashicorp/terraform-provider-azurerm/issues/24153))
* `azurerm_redis_cache` - support for the `storage_account_subscription_id` property ([#24101](https://github.com/hashicorp/terraform-provider-azurerm/issues/24101))
* `azurerm_storage_blob` -  support for the `source_content` type `Page` ([#24177](https://github.com/hashicorp/terraform-provider-azurerm/issues/24177))
* `azurerm_web_application_firewall_policy` - support new values to the `rule_group_name` property ([#24194](https://github.com/hashicorp/terraform-provider-azurerm/issues/24194))
* `azurerm_windows_web_app` - make the `client_secret_setting_name` property optional and conflicts with the `client_secret_certificate_thumbprint` property ([#21834](https://github.com/hashicorp/terraform-provider-azurerm/issues/21834))
* `azurerm_windows_web_app_slot` - make the `client_secret_setting_name` property optional and conflicts with the `client_secret_certificate_thumbprint` property ([#21834](https://github.com/hashicorp/terraform-provider-azurerm/issues/21834))
* `azurerm_windows_web_app` - fix a bug in `app_settings` where settings could be lost ([#24221](https://github.com/hashicorp/terraform-provider-azurerm/issues/24221))
* `azurerm_windows_web_app_slot` - fix a bug in `app_settings` where settings could be lost ([#24221](https://github.com/hashicorp/terraform-provider-azurerm/issues/24221))
* `azurerm_cognitive_account` - add `ContentSafety` to the `kind` property validation ([#24205](https://github.com/hashicorp/terraform-provider-azurerm/issues/24205))

BUG FIXES:

* provider: fix an authentication issue with Azure Storage when running in Azure China cloud ([#24246](https://github.com/hashicorp/terraform-provider-azurerm/issues/24246))
* Data Source: `azurerm_role_definition` - fix bug where `role_definition_id` and `scope` were being incorrectly set ([#24211](https://github.com/hashicorp/terraform-provider-azurerm/issues/24211))
* `azurerm_batch_account` - fix bug where `UserAssigned, SystemAssigned` could be passed to the resource even though it isn't supported ([#24204](https://github.com/hashicorp/terraform-provider-azurerm/issues/24204))
* `azurerm_batch_pool` - fix bug where `settings_json` and `protected_settings` were not being unmarshaled ([#24075](https://github.com/hashicorp/terraform-provider-azurerm/issues/24075))
* `azurerm_bot_service_azure_bot` - fix bug where `public_network_access_enabled` was being set as the value for `LuisKey` ([#24164](https://github.com/hashicorp/terraform-provider-azurerm/issues/24164))
* `azurerm_cognitive_account_customer_managed_key` - `identity_client_id` is no longer passed to the api when it is empty ([#24231](https://github.com/hashicorp/terraform-provider-azurerm/issues/24231))
* `azurerm_linux_web_app_slot` - error when `service_plan_id` is identical to the parent `service_plan_id` ([#23403](https://github.com/hashicorp/terraform-provider-azurerm/issues/23403))
* `azurerm_management_group_template_deployment` - fixing a bug where `template_spec_version_id` couldn't be updated ([#24072](https://github.com/hashicorp/terraform-provider-azurerm/issues/24072))
* `azurerm_pim_active_role_assignment` - fix an importing issue by filtering available role assignments based on the provided `scope` ([#24077](https://github.com/hashicorp/terraform-provider-azurerm/issues/24077))
* `azurerm_pim_eligible_role_assignment` - fix an importing issue by filtering available role assignments based on the provided `scope` ([#24077](https://github.com/hashicorp/terraform-provider-azurerm/issues/24077))
* `azurerm_resource_group_template_deployment` - fixing a bug where `template_spec_version_id` couldn't be updated ([#24072](https://github.com/hashicorp/terraform-provider-azurerm/issues/24072))
* `azurerm_security_center_setting` - fix the casing for the `setting_name` `Sentinel` ([#24210](https://github.com/hashicorp/terraform-provider-azurerm/issues/24210))
* `azurerm_storage_account` - Fix crash when checking for `routingInputs.PublishInternetEndpoints` and `routingInputs.PublishMicrosoftEndpoints` ([#24228](https://github.com/hashicorp/terraform-provider-azurerm/issues/24228))
* `azurerm_storage_share_file` - prevent panic when the file specified by `source` is empty ([#24179](https://github.com/hashicorp/terraform-provider-azurerm/issues/24179))
* `azurerm_subscription_template_deployment` - fixing a bug where `template_spec_version_id` couldn't be updated ([#24072](https://github.com/hashicorp/terraform-provider-azurerm/issues/24072))
* `azurerm_tenant_template_deployment` - fixing a bug where `template_spec_version_id` couldn't be updated ([#24072](https://github.com/hashicorp/terraform-provider-azurerm/issues/24072))
* `azurerm_virtual_machine` - prevent a panic by nil checking the first element of `additional_capabilities` ([#24159](https://github.com/hashicorp/terraform-provider-azurerm/issues/24159))
* `azurerm_windows_web_app_slot` - error when `service_plan_id` is identical to the parent `service_plan_id` ([#23403](https://github.com/hashicorp/terraform-provider-azurerm/issues/23403))

## 3.84.0 (December 07, 2023)

FEATURES:

* **New Data Source:** `azurerm_storage_containers` ([#24061](https://github.com/hashicorp/terraform-provider-azurerm/issues/24061))
* **New Resource:** `azurerm_elastic_san` ([#23619](https://github.com/hashicorp/terraform-provider-azurerm/issues/23619))
* **New Resource:** `azurerm_key_vault_managed_hardware_security_module_role_assignment` ([#22332](https://github.com/hashicorp/terraform-provider-azurerm/issues/22332))
* **New Resource:** `azurerm_key_vault_managed_hardware_security_module_role_definition` ([#22332](https://github.com/hashicorp/terraform-provider-azurerm/issues/22332))

ENHANCEMENTS:

* dependencies: updating mssql elasticpools from `v5.0` to `2023-05-01-preview`
* dependencies: updating to `v0.20231207.1122031` of `github.com/hashicorp/go-azure-sdk` ([#24149](https://github.com/hashicorp/terraform-provider-azurerm/issues/24149))
* Data Source: `azurerm_storage_account` - export the primary and secondary internet and microsoft hostnames for blobs, dfs, files, queues, tables and web ([#23517](https://github.com/hashicorp/terraform-provider-azurerm/issues/23517))
* Data Source: `azurerm_cosmosdb_account` - export the `connection_strings`, `primary_sql_connection_string`, `secondary_sql_connection_string`, `primary_readonly_sql_connection_string`, `secondary_readonly_sql_connection_string`, `primary_mongodb_connection_string`, `secondary_mongodb_connection_string`, `primary_readonly_mongodb_connection_string`, and `secondary_readonly_mongodb_connection_string` attributes ([#24129](https://github.com/hashicorp/terraform-provider-azurerm/issues/24129))
* `azurerm_bot_service_azure_bot` - support for the `public_network_access_enabled` property ([#24125](https://github.com/hashicorp/terraform-provider-azurerm/issues/24125))
* `azurerm_container_app_environment` - support for the `workload_profile` property ([#23478](https://github.com/hashicorp/terraform-provider-azurerm/issues/23478))
* `azurerm_cosmosdb_cassandra_datacenter` - support for the `seed_node_ip_addresses` property ([#24076](https://github.com/hashicorp/terraform-provider-azurerm/issues/24076))
* `azurerm_firewall` - support for the `dns_proxy_enabled` property ([#20519](https://github.com/hashicorp/terraform-provider-azurerm/issues/20519))
* `azurerm_kubernetes_cluster` - support for the  `support_plan` property and the `sku_tier` `Premium` ([#23970](https://github.com/hashicorp/terraform-provider-azurerm/issues/23970))
* `azurerm_mssql_database` - support for `enclave_type` field ([#24054](https://github.com/hashicorp/terraform-provider-azurerm/issues/24054))
* `azurerm_mssql_elasticpool` - support for `enclave_type` field ([#24054](https://github.com/hashicorp/terraform-provider-azurerm/issues/24054))
* `azurerm_mssql_managed_instance` - support for more `vcores`: `6`, `10`, `12`, `20`, `48`, `56`, `96`, `128` ([#24085](https://github.com/hashicorp/terraform-provider-azurerm/issues/24085))
* `azurerm_redis_linked_server` - support for the property `geo_replicated_primary_host_name` ([#23984](https://github.com/hashicorp/terraform-provider-azurerm/issues/23984))
* `azurerm_storage_account` - expose the primary and secondary internet and microsoft hostnames for blobs, dfs, files, queues, tables and web ([#23517](https://github.com/hashicorp/terraform-provider-azurerm/issues/23517))
* `azurerm_synapse_role_assignment` - support for the `principal_type` property ([#24089](https://github.com/hashicorp/terraform-provider-azurerm/issues/24089))
* `azurerm_spring_cloud_build_deployment` - support for the `application_performance_monitoring_ids` property ([#23969](https://github.com/hashicorp/terraform-provider-azurerm/issues/23969))
* `azurerm_virtual_network_gateway` - support for the `bgp_route_translation_for_nat_enabled`, `dns_forwarding_enabled`, `ip_sec_replay_protection_enabled`, `remote_vnet_traffic_enabled`, `virtual_wan_traffic_enabled`, `radius_server`, `virtual_network_gateway_client_connection`, `policy_group`, and `ipsec_policy` property ([#23220](https://github.com/hashicorp/terraform-provider-azurerm/issues/23220))

BUG FIXES:

* `azurerm_application_insights_api_key` - prevent a bug where multiple keys couldn't be created for an Application Insights instance ([#23463](https://github.com/hashicorp/terraform-provider-azurerm/issues/23463))
* `azurerm_container_registry` - the `network_rule_set.virtual_network` property has been deprecated ([#24140](https://github.com/hashicorp/terraform-provider-azurerm/issues/24140))
* `azurerm_hdinsight_hadoop_cluster` - set `roles.edge_node.install_script_action.parameters` into state by retrieving the value provided in the user config since this property isn't returned by the API ([#23971](https://github.com/hashicorp/terraform-provider-azurerm/issues/23971))
* `azurerm_kubernetes_cluster` - prevent a bug where maintenance window start date was always recalculated and sent to the API ([#23985](https://github.com/hashicorp/terraform-provider-azurerm/issues/23985))
* `azurerm_mssql_database` - will no longer send all long retention values in payload unless set ([#24124](https://github.com/hashicorp/terraform-provider-azurerm/issues/24124))
* `azurerm_mssql_managed_database` - will no longer send all long retention values in payload unless set ([#24124](https://github.com/hashicorp/terraform-provider-azurerm/issues/24124))
* `azurerm_mssql_server_microsoft_support_auditing_policy` - only include storage endpoint in payload if set ([#24122](https://github.com/hashicorp/terraform-provider-azurerm/issues/24122))
* `azurerm_mobile_network_packet_core_control_plane` - prevent a panic if the HTTP Response is nil ([#24083](https://github.com/hashicorp/terraform-provider-azurerm/issues/24083))
* `azurerm_storage_account` - revert plan time name validation `(#23799)` ([#24142](https://github.com/hashicorp/terraform-provider-azurerm/issues/24142))
* `azurerm_web_application_firewall_policy` - split create and update function to fix lifecycle - ignore changes ([#23412](https://github.com/hashicorp/terraform-provider-azurerm/issues/23412))

## 3.83.0 (November 30, 2023)

UPGRADE NOTES

* Key Vaults are now loaded using [the `ListBySubscription` API within the Key Vault Resource Provider](https://learn.microsoft.com/en-us/rest/api/keyvault/keyvault/vaults/list-by-subscription?view=rest-keyvault-keyvault-2022-07-01&tabs=HTTP) rather than [the Resources API](https://learn.microsoft.com/en-us/rest/api/keyvault/keyvault/vaults/list?view=rest-keyvault-keyvault-2022-07-01&tabs=HTTP). This change means that the Provider now caches the list of Key Vaults available within a Subscription, rather than loading these piecemeal to workaround stale data returned from the Resources API ([#24019](https://github.com/hashicorp/terraform-provider-azurerm/issues/24019))

FEATURES:

* New Data Source: `azurerm_stack_hci_cluster` ([#24032](https://github.com/hashicorp/terraform-provider-azurerm/issues/24032))

ENHANCEMENTS:

* dependencies: updating to `v0.20231129.1103252` of `github.com/hashicorp/go-azure-sdk` ([#24063](https://github.com/hashicorp/terraform-provider-azurerm/issues/24063))
* `automation`: updating to API Version `2023-11-01` ([#24017](https://github.com/hashicorp/terraform-provider-azurerm/issues/24017))
* `keyvault`: the cache is now populated using the `ListBySubscription` endpoint on the KeyVault Resource Provider rather than via the `Resources` API ([#24019](https://github.com/hashicorp/terraform-provider-azurerm/issues/24019)).
* `keyvault`: updating the cache to populate all Key Vaults available within the Subscription to reduce the number of API calls ([#24019](https://github.com/hashicorp/terraform-provider-azurerm/issues/24019))
* Data Source `azurerm_private_dns_zone`: refactoring to use the `ListBySubscription` API rather than the Resources API when `resource_group_name` is omitted ([#24024](https://github.com/hashicorp/terraform-provider-azurerm/issues/24024))
* `azurerm_dashboard_grafana` - support for `grafana_major_version` ([#24014](https://github.com/hashicorp/terraform-provider-azurerm/issues/24014))
* `azurerm_linux_web_app` - add support for dotnet 8 ([#23893](https://github.com/hashicorp/terraform-provider-azurerm/issues/23893))
* `azurerm_linux_web_app_slot` - add support for dotnet 8 ([#23893](https://github.com/hashicorp/terraform-provider-azurerm/issues/23893))
* `azurerm_media_transform` -  deprecate `face_detector_preset` and `video_analyzer_preset` ([#24002](https://github.com/hashicorp/terraform-provider-azurerm/issues/24002))
* `azurerm_postgresql_database` - update the validation of `collation` to include `Norwegian_Norway.1252` ([#24070](https://github.com/hashicorp/terraform-provider-azurerm/issues/24070))
* `azurerm_postgresql_flexible_server` - updating to API Version `2023-06-01-preview` ([#24016](https://github.com/hashicorp/terraform-provider-azurerm/issues/24016))
* `azurerm_redis_cache` - support for the `active_directory_authentication_enabled` property ([#23976](https://github.com/hashicorp/terraform-provider-azurerm/issues/23976))
* `azurerm_windows_web_app` - add support for dotnet 8 ([#23893](https://github.com/hashicorp/terraform-provider-azurerm/issues/23893))
* `azurerm_windows_web_app_slot` - add support for dotnet 8 ([#23893](https://github.com/hashicorp/terraform-provider-azurerm/issues/23893))
* `azurerm_storage_account` -  add `name` validation in custom diff ([#23799](https://github.com/hashicorp/terraform-provider-azurerm/issues/23799))

BUG FIXES:

* authentication: fix a bug where auxiliary tenants were not correctly authorized ([#24063](https://github.com/hashicorp/terraform-provider-azurerm/issues/24063))
* `azurerm_app_configuration` - normalize location in `replica` block ([#24074](https://github.com/hashicorp/terraform-provider-azurerm/issues/24074))
* `azurerm_cosmosdb_account` - cosmosdb version and capabilities can now be updated at the same time ([#24029](https://github.com/hashicorp/terraform-provider-azurerm/issues/24029))
* `azurerm_data_factory_flowlet_data_flow` - `source` and `sink` properties are now optional ([#23987](https://github.com/hashicorp/terraform-provider-azurerm/issues/23987))
* `azurerm_datadog_monitor_tag_rule` - correctly handle default rule ([#22806](https://github.com/hashicorp/terraform-provider-azurerm/issues/22806))
* `azurerm_ip_group`: fixing a crash when `firewall_ids` and `firewall_policy_ids` weren't parsed correctly from the API Response ([#24031](https://github.com/hashicorp/terraform-provider-azurerm/issues/24031))
* `azurerm_nginx_deployment` - add default value of `20` for `capacity` ([#24033](https://github.com/hashicorp/terraform-provider-azurerm/issues/24033))

## 3.82.0 (November 23, 2023)

FEATURES:

* New Data Source: `azurerm_monitor_workspace` ([#23928](https://github.com/hashicorp/terraform-provider-azurerm/issues/23928))
* New Resource: `azurerm_application_load_balancer_subnet_association` ([#23628](https://github.com/hashicorp/terraform-provider-azurerm/issues/23628))

ENHANCEMENTS:

* dependencies: updating to `v0.20231117.1130141` of `github.com/hashicorp/go-azure-sdk` ([#23945](https://github.com/hashicorp/terraform-provider-azurerm/issues/23945))
* `azurestackhci`: updating to API Version `2023-08-01` ([#23939](https://github.com/hashicorp/terraform-provider-azurerm/issues/23939))
* `dashboard`: updating to API Version `2023-09-01` ([#23929](https://github.com/hashicorp/terraform-provider-azurerm/issues/23929))
* `hpccache`: updating to API version `2023-05-01` ([#24005](https://github.com/hashicorp/terraform-provider-azurerm/issues/24005))
* `mssql`: updating resources using `hashicorp/go-azure-sdk` to API Version `2023-02-01-preview` ([#23721](https://github.com/hashicorp/terraform-provider-azurerm/issues/23721))
* `templatespecversions`: updating to API Version `2022-02-01` ([#24007](https://github.com/hashicorp/terraform-provider-azurerm/issues/24007))
*  Data Source: `azurerm_template_spec_version` - refactoring to use `hashicorp/go-azure-sdk` ([#24007](https://github.com/hashicorp/terraform-provider-azurerm/issues/24007))
* `azurerm_cosmosdb_postgresql_cluster` - `coordinator_storage_quota_in_mb` and `coordinator_vcore_count` are no longer required for read replicas ([#23928](https://github.com/hashicorp/terraform-provider-azurerm/issues/23928))
* `azurerm_dashboard_grafana` - `sku` can now be set to `Essential` ([#23934](https://github.com/hashicorp/terraform-provider-azurerm/issues/23934))
* `azurerm_gallery_application_version` - add support for the `config_file`, `package_file` and `target_region.exclude_from_latest` properties ([#23816](https://github.com/hashicorp/terraform-provider-azurerm/issues/23816))
* `azurerm_hdinsight_hadoop_cluster` - `script_actions` is no longer Force New ([#23888](https://github.com/hashicorp/terraform-provider-azurerm/issues/23888))
* `azurerm_hdinsight_hbase_cluster` - `script_actions` is no longer Force New ([#23888](https://github.com/hashicorp/terraform-provider-azurerm/issues/23888))
* `azurerm_hdinsight_interactive_query_cluster` - `script_actions` is no longer Force New ([#23888](https://github.com/hashicorp/terraform-provider-azurerm/issues/23888))
* `azurerm_hdinsight_kafka_cluster` - `script_actions` is no longer Force New ([#23888](https://github.com/hashicorp/terraform-provider-azurerm/issues/23888))
* `azurerm_hdinsight_spark_cluster` - `script_actions` is no longer Force New ([#23888](https://github.com/hashicorp/terraform-provider-azurerm/issues/23888))
* `azurerm_kubernetes_cluster` - add support for the `gpu_instance` property ([#23887](https://github.com/hashicorp/terraform-provider-azurerm/issues/23887))
* `azurerm_kubernetes_cluster_node_pool` - add support for the `gpu_instance` property ([#23887](https://github.com/hashicorp/terraform-provider-azurerm/issues/23887))
* `azurerm_log_analytics_workspace` - add support for the `identity` property ([#23864](https://github.com/hashicorp/terraform-provider-azurerm/issues/23864))
* `azurerm_linux_function_app` - add support for dotnet 8 ([#23638](https://github.com/hashicorp/terraform-provider-azurerm/issues/23638))
* `azurerm_linux_function_app_slot` - add support for dotnet 8 ([#23638](https://github.com/hashicorp/terraform-provider-azurerm/issues/23638))
* `azurerm_managed_lustre_file_system` - export attribute `mgs_address` ([#23942](https://github.com/hashicorp/terraform-provider-azurerm/issues/23942))
* `azurerm_mssql_database` - support for Hyperscale SKUs ([#23974](https://github.com/hashicorp/terraform-provider-azurerm/issues/23974))
* `azurerm_mssql_database` - refactoring to use `hashicorp/go-azure-sdk` ([#23721](https://github.com/hashicorp/terraform-provider-azurerm/issues/23721))
* `azurerm_mssql_server` - refactoring to use `hashicorp/go-azure-sdk` ([#23721](https://github.com/hashicorp/terraform-provider-azurerm/issues/23721))
* `azurerm_shared_image` - add support for `trusted_launch_supported` ([#23781](https://github.com/hashicorp/terraform-provider-azurerm/issues/23781))
* `azurerm_spring_cloud_container_deployment` - add support for the `application_performance_monitoring_ids` property ([#23862](https://github.com/hashicorp/terraform-provider-azurerm/issues/23862))
* `azurerm_spring_cloud_customized_accelerator` - add support for the `accelerator_type` and `path` properties ([#23797](https://github.com/hashicorp/terraform-provider-azurerm/issues/23797))
* `azurerm_point_to_site_vpn_gateway` - allow multiple `connection_configurations` blocks ([#23936](https://github.com/hashicorp/terraform-provider-azurerm/issues/23936))
* `azurerm_private_dns_cname_record` - `ttl` can now be set to 0 ([#23918](https://github.com/hashicorp/terraform-provider-azurerm/issues/23918))
* `azurerm_windows_function_app` - add support for dotnet 8 ([#23638](https://github.com/hashicorp/terraform-provider-azurerm/issues/23638))
* `azurerm_windows_function_app_slot` - add support for dotnet 8 ([#23638](https://github.com/hashicorp/terraform-provider-azurerm/issues/23638))

BUG FIXES:
* `azurerm_api_management` - correct a bug with additional location zones within the `additional_location` block ([#23943](https://github.com/hashicorp/terraform-provider-azurerm/issues/23943))
* `azurerm_dev_test_linux_virtual_machine` - `storage_type` is now ForceNew to match the updated API behaviour ([#23973](https://github.com/hashicorp/terraform-provider-azurerm/issues/23973))
* `azurerm_dev_test_windows_virtual_machine` - `storage_type` is now ForceNew to match the updated API behaviour ([#23973](https://github.com/hashicorp/terraform-provider-azurerm/issues/23973))
* `azurerm_disk_encryption_set` - resource will recreate if `identity` changes from `SystemAssigned` to `UserAssigned` ([#23904](https://github.com/hashicorp/terraform-provider-azurerm/issues/23904))
* `azurerm_eventhub_cluster`: `sku_name` is no longer ForceNew ([#24009](https://github.com/hashicorp/terraform-provider-azurerm/issues/24009))
* `azurerm_firewall` - recasing the value for `firewall_policy_id` to workaround the API returning the incorrect casing ([#23993](https://github.com/hashicorp/terraform-provider-azurerm/issues/23993))
* `azurerm_security_center_subscription_pricing` - fix a bug preventing removal of `extensions` and downgrading `tier` to `Free` ([#23821](https://github.com/hashicorp/terraform-provider-azurerm/issues/23821))
* `azurerm_windows_web_app` - fix an issue of incorrect application stack settings during update ([#23372](https://github.com/hashicorp/terraform-provider-azurerm/issues/23372))

## 3.81.0 (November 16, 2023)

ENHANCEMENTS:

* dependencies: updating to `v0.20231116.1162710` of `github.com/hashicorp/go-azure-sdk` ([#23922](https://github.com/hashicorp/terraform-provider-azurerm/issues/23922))
* `managedservices`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#23890](https://github.com/hashicorp/terraform-provider-azurerm/issues/23890))
* `network`: updating to API Version `2023-06-01` ([#23875](https://github.com/hashicorp/terraform-provider-azurerm/issues/23875))
* `servicelinker`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#23890](https://github.com/hashicorp/terraform-provider-azurerm/issues/23890))
* `storage`: refactoring usages of `github.com/hashicorp/go-azure-sdk` to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#23890](https://github.com/hashicorp/terraform-provider-azurerm/issues/23890))
* Data Source: `azurerm_network_ddos_protection_plan`: refactoring to use `hashicorp/go-azure-sdk` ([#23849](https://github.com/hashicorp/terraform-provider-azurerm/issues/23849))
* `azurerm_linux_function_app` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled` ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_linux_function_app_slot` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_linux_web_app` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_linux_web_app_slot` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_logic_app_integration_account_certificate` - `name` now accepts underscores ([#23866](https://github.com/hashicorp/terraform-provider-azurerm/issues/23866))
* `azurerm_logic_app_integration_account_partner` - `business_identity.value` now accepts underscores ([#23866](https://github.com/hashicorp/terraform-provider-azurerm/issues/23866))
* `azurerm_monitor_data_collection_rule` - added support for `WorkspaceTransforms` as `kind` ([#23873](https://github.com/hashicorp/terraform-provider-azurerm/issues/23873))
* `azurerm_network_ddos_protection_plan`: refactoring to use `hashicorp/go-azure-sdk` ([#23849](https://github.com/hashicorp/terraform-provider-azurerm/issues/23849))
* `azurerm_windows_function_app` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_windows_function_app_slot` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_windows_web_app` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_windows_web_app_slot` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))

## 3.80.0 (November 09, 2023)

ENHANCEMENTS:

* `internal/sdk` - Added support for pointer Types in resource models ([#23810](https://github.com/hashicorp/terraform-provider-azurerm/issues/23810))
* dependencies: updating to `v0.63.0` of `github.com/hashicorp/go-azure-helpers` ([#23785](https://github.com/hashicorp/terraform-provider-azurerm/issues/23785))
* dependencies: updating to `v0.20231106.1151347` of `github.com/hashicorp/go-azure-sdk` ([#23787](https://github.com/hashicorp/terraform-provider-azurerm/issues/23787))
* `azurerm_cognitive_deployment` - support for the `version_upgrade_option` property ([#22520](https://github.com/hashicorp/terraform-provider-azurerm/issues/22520))
* `azurerm_firewall_policy_rule_collection_group` - add support for the property `http_headers` ([#23641](https://github.com/hashicorp/terraform-provider-azurerm/issues/23641))
* `azurerm_kubernetes_cluster` - `fips_enabled` can be updated in the `default_node_pool` without recreating the cluster ([#23612](https://github.com/hashicorp/terraform-provider-azurerm/issues/23612))
* `azurerm_kusto_cluster` - the cluster `name` can now include dashes ([#23790](https://github.com/hashicorp/terraform-provider-azurerm/issues/23790))
* `azurerm_postgresql_database` - update the validation of `collation` to include support for `French_France.1252` ([#23783](https://github.com/hashicorp/terraform-provider-azurerm/issues/23783))

BUG FIXES:

* Data Source: `azurerm_data_protection_backup_vault` - removing `import` support, since Data Sources don't support being imported ([#23820](https://github.com/hashicorp/terraform-provider-azurerm/issues/23820))
* Data Source: `azurerm_kusto_database` - removing `import` support, since Data Sources don't support being imported ([#23820](https://github.com/hashicorp/terraform-provider-azurerm/issues/23820))
* Data Source: `azurerm_virtual_hub_route_table` - removing `import` support, since Data Sources don't support being imported ([#23820](https://github.com/hashicorp/terraform-provider-azurerm/issues/23820))
* `azurerm_windows_web_app` - prevent a panic with the `auto_heal.actions` property ([#23836](https://github.com/hashicorp/terraform-provider-azurerm/issues/23836))
* `azurerm_windows_web_app` - prevent a panic with the `auto_heal.triggers` property ([#23812](https://github.com/hashicorp/terraform-provider-azurerm/issues/23812))

## 3.79.0 (November 02, 2023)

ENHANCEMENTS:

* provider: log instead of error when RPs are unavailable when validating RP registrations ([#23380](https://github.com/hashicorp/terraform-provider-azurerm/issues/23380))
* `azurerm_arc_kuberenetes_cluster_extension_resource` - the `version` and `release_train` properties can now be set simultaneously ([#23692](https://github.com/hashicorp/terraform-provider-azurerm/issues/23692))
* `azurerm_container_apps` - support for the `ingress.exposed_port` property ([#23752](https://github.com/hashicorp/terraform-provider-azurerm/issues/23752))
* `azurerm_cosmosdb_postgresql_cluster` - read replica clusters can be created without specifying `administrator_login_password` property ([#23750](https://github.com/hashicorp/terraform-provider-azurerm/issues/23750))
* `azurerm_managed_application` - arrays can be supplied in the `parameter_values` property ([#23754](https://github.com/hashicorp/terraform-provider-azurerm/issues/23754))
* `azurerm_storage_management_policy` - support for properties `rule.*.actions.*.base_blob.0.tier_to_cold_after_days_since_{modification|last_access_time|creation}_greater_than and rule.*.actions.*.{snapshot|version}.0.tier_to_cold_after_days_since_creation_greater_than` ([#23574](https://github.com/hashicorp/terraform-provider-azurerm/issues/23574))

BUG FIXES:

* `azurerm_api_management_diagnostic` - the `operation_name_format` attribute will only be sent if `identifier` is set to `applicationinsights` ([#23736](https://github.com/hashicorp/terraform-provider-azurerm/issues/23736))
* `azurerm_backup_policy_vm` - fix payload by using current datetime ([#23586](https://github.com/hashicorp/terraform-provider-azurerm/issues/23586))
* `azurerm_kubernetes_cluster` - the `custom_ca_trust_certificates_base64` property can not be removed, only updated ([#23737](https://github.com/hashicorp/terraform-provider-azurerm/issues/23737))

## 3.78.0 (October 26, 2023)

FEATURES:

* New Resource: `azurerm_resource_management_private_link_association` ([#23546](https://github.com/hashicorp/terraform-provider-azurerm/issues/23546))

ENHANCEMENTS:

* dependencies: updating to `v0.20231025.1113325` of `github.com/hashicorp/go-azure-sdk` ([#23684](https://github.com/hashicorp/terraform-provider-azurerm/issues/23684))
* dependencies: updating to `v1.58.3` of `google.golang.org/grpc` ([#23691](https://github.com/hashicorp/terraform-provider-azurerm/issues/23691))
* dependencies: updating search service from `2022-09-01` to `2023-11-01` ([#23698](https://github.com/hashicorp/terraform-provider-azurerm/issues/23698)) 
* Data Source: `azurerm_monitor_workspace` - export `query_endpoint` ([#23629](https://github.com/hashicorp/terraform-provider-azurerm/issues/23629))
* `azurerm_express_route_port` - support for `macsec_sci_enabled` ([#23625](https://github.com/hashicorp/terraform-provider-azurerm/issues/23625))
* `azurerm_eventhub_namespace_customer_managed_key` - support for the `user_assigned_identity_id` property ([#23635](https://github.com/hashicorp/terraform-provider-azurerm/issues/23635))
* `azurerm_postgresql_flexible_server` - `private_dns_zone_id` is no longer ForceNew and case is suppressed  ([#23660](https://github.com/hashicorp/terraform-provider-azurerm/issues/23660))
* `azurerm_synapse_workspace` - add support for `azuread_authentication_only` ([#23659](https://github.com/hashicorp/terraform-provider-azurerm/issues/23659))
* `azurerm_redis_enterprise_cluster` - support for new location `Japan East` ([#23696](https://github.com/hashicorp/terraform-provider-azurerm/issues/23696))
* `azurerm_search_service` - support for `semantic_search_sku` field ([#23698](https://github.com/hashicorp/terraform-provider-azurerm/issues/23698))

BUG FIXES:

* `azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack` - added lock for ruleStackID ([#23601](https://github.com/hashicorp/terraform-provider-azurerm/issues/23601))
* `azurerm_cognitive_deployment` - remove forceNew tag from `rai_policy_name` ([#23697](https://github.com/hashicorp/terraform-provider-azurerm/issues/23697))

## 3.77.0 (October 19, 2023)

FEATURES:

* New Resources: `azurerm_application_load_balancer_frontend` ([#23411](https://github.com/hashicorp/terraform-provider-azurerm/issues/23411))
* New Resources: `azurerm_dev_center` ([#23538](https://github.com/hashicorp/terraform-provider-azurerm/issues/23538))
* New Resources: `azurerm_dev_center_project` ([#23538](https://github.com/hashicorp/terraform-provider-azurerm/issues/23538))

ENHANCEMENTS:

* dependencies: updating to `v0.62.0` of `github.com/hashicorp/go-azure-helpers` ([#23581](https://github.com/hashicorp/terraform-provider-azurerm/issues/23581))
* dependencies: updating Kusto SDK from `2023-05-02` to `2023-08-15` ([#23598](https://github.com/hashicorp/terraform-provider-azurerm/issues/23598))
* dependencies: updating nginx from `2022-08-01` to `2023-04-01` ([#23583](https://github.com/hashicorp/terraform-provider-azurerm/issues/23583))
* `netapp`: updating to use API Version `2023-05-01` ([#23576](https://github.com/hashicorp/terraform-provider-azurerm/issues/23576))
* `springcloud`: updating to use API Version `2023-09-01-preview` ([#23544](https://github.com/hashicorp/terraform-provider-azurerm/issues/23544))
* `storage`: updating to use API Version `2023-01-01` ([#23543](https://github.com/hashicorp/terraform-provider-azurerm/issues/23543))
* `internal/sdk`: fixing an issue where struct fields containing `removedInNextMajorVersion` wouldn't be decoded correctly ([#23564](https://github.com/hashicorp/terraform-provider-azurerm/issues/23564))
* `internal/sdk`: struct tag parsing is now handled consistently during both encoding and decoding ([#23568](https://github.com/hashicorp/terraform-provider-azurerm/issues/23568))
* provider: the `roll_instances_when_required` provider feature in the `virtual_machine_scale_set` block is now optional ([#22976](https://github.com/hashicorp/terraform-provider-azurerm/issues/22976))
* Data Source: `azurerm_automation_account`: refactoring the remaining usage of `Azure/azure-sdk-for-go` to use `hashicorp/go-azure-sdk` ([#23555](https://github.com/hashicorp/terraform-provider-azurerm/issues/23555))
* `azurerm_automation_account`: refactoring the remaining usage of `Azure/azure-sdk-for-go` to use `hashicorp/go-azure-sdk` ([#23555](https://github.com/hashicorp/terraform-provider-azurerm/issues/23555))
* `azurerm_resource_deployment_script_azure_cli` - improve validation for the `version` property to support newer versions ([#23370](https://github.com/hashicorp/terraform-provider-azurerm/issues/23370))
* `azurerm_resource_deployment_script_azure_power_shell` - improve validation for the `version` property to support newer versions ([#23370](https://github.com/hashicorp/terraform-provider-azurerm/issues/23370))
* `azurerm_nginx_deployment` - support for the `capacity` and `email` properties ([#23596](https://github.com/hashicorp/terraform-provider-azurerm/issues/23596))

BUG FIXES:

* Data Source: `azurerm_virtual_hub_connection` - export the `inbound_route_map_id`, `outbound_route_map_id`, and `static_vnet_local_route_override_criteria` attributes in the `routing` block, and fix a bug where these attributes could not be set ([#23491](https://github.com/hashicorp/terraform-provider-azurerm/issues/23491))
* `azurerm_cdn_frontdoor_rule` - the `url_filename_condition` properties `match_values` is now optional if `operator` is set to `Any` ([#23541](https://github.com/hashicorp/terraform-provider-azurerm/issues/23541))
* `azurerm_shared_image_gallery` - added the `Private` and `Groups` options for the `sharing.permission` property ([#23570](https://github.com/hashicorp/terraform-provider-azurerm/issues/23570))
* `azurerm_redis_cache` - fixed incorrect ssl values for `redis_primary_connection_string` and `secondary_connection_string` ([#23575](https://github.com/hashicorp/terraform-provider-azurerm/issues/23575))
* `azurerm_monitor_activity_log_alert` - the `recommend_category` property now can be set to `HighAvailability` ([#23605](https://github.com/hashicorp/terraform-provider-azurerm/issues/23605))
* `azurerm_recovery_services_vault` - the `encryption` property can now be used with the `cross_region_restore_enabled` property ([#23618](https://github.com/hashicorp/terraform-provider-azurerm/issues/23618))
* `azurerm_storage_account_customer_managed_key` - prevent a panic when the keyvault id is empty ([#23599](https://github.com/hashicorp/terraform-provider-azurerm/issues/23599))

## 3.76.0 (October 12, 2023)

FEATURES:

* New Resource: `azurerm_security_center_storage_defender` ([#23242](https://github.com/hashicorp/terraform-provider-azurerm/issues/23242))
* New Resource: `azurerm_spring_cloud_application_insights_application_performance_monitoring` ([#23107](https://github.com/hashicorp/terraform-provider-azurerm/issues/23107))

ENHANCEMENTS:

* provider: updating to build using Go `1.21.3` ([#23514](https://github.com/hashicorp/terraform-provider-azurerm/issues/23514))
* dependencies: updating to `v0.20231012.1141427` of `github.com/hashicorp/go-azure-sdk` ([#23534](https://github.com/hashicorp/terraform-provider-azurerm/issues/23534))
* Data Source: `azurerm_application_gateway` - support for `backend_http_settings`, `global`, `gateway_ip_configuration` and additional attributes ([#23318](https://github.com/hashicorp/terraform-provider-azurerm/issues/23318))
* Data Source: `azurerm_network_service_tags` - export the `name` attribute ([#23382](https://github.com/hashicorp/terraform-provider-azurerm/issues/23382))
* `azurerm_cosmosdb_postgresql_cluster` - add support for `sql_version` of `16` and `citus_version` of `12.1` ([#23476](https://github.com/hashicorp/terraform-provider-azurerm/issues/23476))
* `azurerm_palo_alto_local_rulestack` - correctly normalize the `location` property ([#23483](https://github.com/hashicorp/terraform-provider-azurerm/issues/23483))
* `azurerm_static_site` - add support for `app_settings` ([#23421](https://github.com/hashicorp/terraform-provider-azurerm/issues/23421))

BUG FIXES:

* `azurerm_automation_schedule` - fix a bug when updating `start_time` ([#23494](https://github.com/hashicorp/terraform-provider-azurerm/issues/23494))
* `azurerm_eventhub` - remove ForceNew and check `partition_count` is not decreased ([#23499](https://github.com/hashicorp/terraform-provider-azurerm/issues/23499))
* `azurerm_managed_lustre_file_system` - update validation for `storage_capacity_in_tb` according to `sku_name` in use ([#23428](https://github.com/hashicorp/terraform-provider-azurerm/issues/23428))
* `azurerm_virtual_machine` - fix a crash when the API response for the `os_profile` block contains nil properties ([#23535](https://github.com/hashicorp/terraform-provider-azurerm/issues/23535))

## 3.75.0 (September 28, 2023)

FEATURES:

* New Resource: `azurerm_application_load_balancer` ([#22517](https://github.com/hashicorp/terraform-provider-azurerm/issues/22517))
* New Resource: `azurerm_resource_management_private_link` ([#23098](https://github.com/hashicorp/terraform-provider-azurerm/issues/23098))

ENHANCEMENTS:

* dependencies: `firewall` migrated to `hashicorp/go-azure-sdk` ([#22863](https://github.com/hashicorp/terraform-provider-azurerm/issues/22863))
* `azurerm_bot_service_azure_bot` - add support for the `icon_url` property ([#23114](https://github.com/hashicorp/terraform-provider-azurerm/issues/23114))
* `azurerm_cognitive_deployment` - `capacity` property is now updateable ([#23251](https://github.com/hashicorp/terraform-provider-azurerm/issues/23251))
* `azurerm_container_group` - added support for `key_vault_user_identity_id` ([#23332](https://github.com/hashicorp/terraform-provider-azurerm/issues/23332))
* `azurerm_data_factory` - added support for the `publish_enabled` property ([#2334](https://github.com/hashicorp/terraform-provider-azurerm/issues/2334))
* `azurerm_firewall_policy_rule_collection_group` - add support for the `description` property ([#23354](https://github.com/hashicorp/terraform-provider-azurerm/issues/23354))
* `azurerm_kubernetes_cluster` - `network_profile.network_policy` can be migrated to `cilium` ([#23342](https://github.com/hashicorp/terraform-provider-azurerm/issues/23342))
* `azurerm_log_analytics_workspace` - add support for the `data_collection_rule_id` property ([#23347](https://github.com/hashicorp/terraform-provider-azurerm/issues/23347))
* `azurerm_mysql_flexible_server` - add support for the `io_scaling_enabled` property ([#23329](https://github.com/hashicorp/terraform-provider-azurerm/issues/23329))

BUG FIXES:

* `azurerm_api_management_api` - fix importing `openapi` format content file issue ([#23348](https://github.com/hashicorp/terraform-provider-azurerm/issues/23348))
* `azurerm_cdn_frontdoor_rule` - allow a `cache_duration` of `00:00:00` ([#23384](https://github.com/hashicorp/terraform-provider-azurerm/issues/23384))
* `azurerm_cosmosdb_cassandra_datacenter` - `sku_name` is now updatable ([#23419](https://github.com/hashicorp/terraform-provider-azurerm/issues/23419))
* `azurerm_key_vault_certificate` - fix a bug that prevented soft-deleted certificates from being recovered ([#23204](https://github.com/hashicorp/terraform-provider-azurerm/issues/23204))
* `azurerm_log_analytics_solution` - fix create and update lifecycle of resource by splitting methods ([#23333](https://github.com/hashicorp/terraform-provider-azurerm/issues/23333))
* `azurerm_management_group_subscription_association` - mark resource as gone correctly if not found when retrieving ([#23335](https://github.com/hashicorp/terraform-provider-azurerm/issues/23335))
* `azurerm_management_lock` - add polling after create and delete to check for RP propagation ([#23345](https://github.com/hashicorp/terraform-provider-azurerm/issues/23345))
* `azurerm_monitor_diagnostic_setting` - added validation to ensure at least one of `category` or `category_group` is supplied ([#23308](https://github.com/hashicorp/terraform-provider-azurerm/issues/23308))
* `azurerm_palo_alto_local_rulestack_prefix_list` - fix rulestack not being committed on delete ([#23362](https://github.com/hashicorp/terraform-provider-azurerm/issues/23362))
* `azurerm_palo_alto_local_rulestack_fqdn_list` - fix rulestack not being committed on delete ([#23362](https://github.com/hashicorp/terraform-provider-azurerm/issues/23362))
* `security_center_subscription_pricing_resource` - disabled extensions logic now works as expected ([#22997](https://github.com/hashicorp/terraform-provider-azurerm/issues/22997))


## 3.74.0 (September 21, 2023)

NOTES:

* `azurerm_synapse_sql_pool` - users that have imported `azurerm_synapse_sql_pool` resources that were created outside of Terraform using an `LRS` storage account type will need to use `ignore_changes` to avoid the resource from being destroyed and recreated.

FEATURES:

* **New Resource**: `azurerm_arc_resource_bridge_appliance` ([#23108](https://github.com/hashicorp/terraform-provider-azurerm/issues/23108))
* **New Resource**: `azurerm_data_factory_dataset_azure_sql_table` ([#23264](https://github.com/hashicorp/terraform-provider-azurerm/issues/23264))
* **New Resource**: `azurerm_function_app_connection` ([#23127](https://github.com/hashicorp/terraform-provider-azurerm/issues/23127))

ENHANCEMENTS:

* dependencies: updating to `v0.20230918.1115907` of `github.com/hashicorp/go-azure-sdk` ([#23337](https://github.com/hashicorp/terraform-provider-azurerm/issues/23337))
* dependencies: downgrading to `v1.12.5` of `github.com/rickb777/date` ([#23296](https://github.com/hashicorp/terraform-provider-azurerm/issues/23296))
* `mysql`: updating to use API Version `2022-01-01` ([#23320](https://github.com/hashicorp/terraform-provider-azurerm/issues/23320))
* `azurerm_app_configuration` - support for the `replica` block ([#22452](https://github.com/hashicorp/terraform-provider-azurerm/issues/22452))
* `azurerm_bot_channel_directline` - support for `user_upload_enabled`, `endpoint_parameters_enabled`, and `storage_enabled` ([#23149](https://github.com/hashicorp/terraform-provider-azurerm/issues/23149))
* `azurerm_container_app` - support for scale rules ([#23294](https://github.com/hashicorp/terraform-provider-azurerm/issues/23294))
* `azurerm_container_app_environment` - support for zone redundancy ([#23313](https://github.com/hashicorp/terraform-provider-azurerm/issues/23313))
* `azurerm_container_group` - support for the `key_vault_user_identity_id` property for Customer Managed Keys ([#23332](https://github.com/hashicorp/terraform-provider-azurerm/issues/23332))
* `azurerm_cosmosdb_account` - support for MongoDB connection strings ([#23331](https://github.com/hashicorp/terraform-provider-azurerm/issues/23331))
* `azurerm_data_factory_dataset_delimited_text` - support for the `dynamic_file_system_enabled`, `dynamic_path_enabled`, and `dynamic_filename_enabled` properties ([#23261](https://github.com/hashicorp/terraform-provider-azurerm/issues/23261))
* `azurerm_data_factory_dataset_parquet` - support for the `azure_blob_fs_location` block ([#23261](https://github.com/hashicorp/terraform-provider-azurerm/issues/23261))
* `azurerm_monitor_diagnostic_setting` - validation to ensure either `category` or `category_group` are supplied in `enabled_log` and `log` blocks ([#23308](https://github.com/hashicorp/terraform-provider-azurerm/issues/23308))
* `azurerm_network_interface` - support for the `auxiliary_mode` and `auxiliary_sku` properties ([#22979](https://github.com/hashicorp/terraform-provider-azurerm/issues/22979))
* `azurerm_postgresql_flexible_server` - increased the maximum supported value for `storage_mb` ([#23277](https://github.com/hashicorp/terraform-provider-azurerm/issues/23277))
* `azurerm_shared_image_version` - support for the `replicated_region_deletion_enabled` and `target_region.exclude_from_latest_enabled` properties ([#23147](https://github.com/hashicorp/terraform-provider-azurerm/issues/23147))
* `azurerm_storage_account` - support for setting `domain_name` and `domain_guid` for `AADKERB` ([#22833](https://github.com/hashicorp/terraform-provider-azurerm/issues/22833))
* `azurerm_storage_account_customer_managed_key` - support for cross-tenant customer-managed keys with the `federated_identity_client_id`, and `key_vault_uri` properties ([#20356](https://github.com/hashicorp/terraform-provider-azurerm/issues/20356))
* `azurerm_web_application_firewall_policy` - support for the `rate_limit_duration`, `rate_limit_threshold`, `group_rate_limit_by`, and `request_body_inspect_limit_in_kb` properties ([#23239](https://github.com/hashicorp/terraform-provider-azurerm/issues/23239))

BUG FIXES:

* Data Source: `azurerm_container_app_environment`: fix `log_analytics_workspace_name` output to correct value ([#23298](https://github.com/hashicorp/terraform-provider-azurerm/issues/23298))
* `azurerm_api_management_api` - set the `service_url` property when importing the resource ([#23011](https://github.com/hashicorp/terraform-provider-azurerm/issues/23011))
* `azurerm_app_configuration` - prevent crash by nil checking the encryption configuration ([#23302](https://github.com/hashicorp/terraform-provider-azurerm/issues/23302))
* `azurerm_app_configuration_feature` - update `percentage_filter_value` to accept correct type of float ([#23263](https://github.com/hashicorp/terraform-provider-azurerm/issues/23263))
* `azurerm_container_app` - fix an issue with `commands` and `args` being overwritten when using multiple containers ([#23338](https://github.com/hashicorp/terraform-provider-azurerm/issues/23338))
* `azurerm_key_vault_certificate` - fix issue where certificates couldn't be recovered anymore ([#23204](https://github.com/hashicorp/terraform-provider-azurerm/issues/23204))
* `azurerm_key_vault_key` - the ForceNew when `expiration_date` is removed from the config file ([#23327](https://github.com/hashicorp/terraform-provider-azurerm/issues/23327))
* `azurerm_linux_function_app` - fix a bug in setting the storage settings when using Elastic Premium plans ([#21212](https://github.com/hashicorp/terraform-provider-azurerm/issues/21212))
* `azurerm_linux_web_app` - fix docker app stack update ([#23303](https://github.com/hashicorp/terraform-provider-azurerm/issues/23303))
* `azurerm_linux_web_app` - fix crash in auto heal expansion ([#21328](https://github.com/hashicorp/terraform-provider-azurerm/issues/21328))
* `azurerm_linux_web_app_slot` - fix docker app stack update ([#23303](https://github.com/hashicorp/terraform-provider-azurerm/issues/23303))
* `azurerm_linux_web_app_slot` - fix crash in auto heal expansion ([#21328](https://github.com/hashicorp/terraform-provider-azurerm/issues/21328))
* `azurerm_log_analytics_solution` - fix bug where the resource wasn't handling successful creation on subsequent applies ([#23312](https://github.com/hashicorp/terraform-provider-azurerm/issues/23312))
* `azurerm_management_group_subscription_association` - fix bug to correctly mark resource as gone if not found during read ([#23335](https://github.com/hashicorp/terraform-provider-azurerm/issues/23335))
* `azurerm_mssql_elasticpool` - remove check that prevents `license_type` from being set for certain skus ([#23262](https://github.com/hashicorp/terraform-provider-azurerm/issues/23262))
* `azurerm_servicebus_queue` - fixing an issue where `auto_delete_on_idle` couldn't be set to `P10675199DT2H48M5.4775807S` ([#23296](https://github.com/hashicorp/terraform-provider-azurerm/issues/23296))
* `azurerm_servicebus_topic` - fixing an issue where `auto_delete_on_idle` couldn't be set to `P10675199DT2H48M5.4775807S` ([#23296](https://github.com/hashicorp/terraform-provider-azurerm/issues/23296))
* `azurerm_storage_account` - prevent sending unsupported blob properties in payload for `Storage` account kind ([#23288](https://github.com/hashicorp/terraform-provider-azurerm/issues/23288))
* `azurerm_synapse_sql_pool` - expose `storage_account_type` ([#23217](https://github.com/hashicorp/terraform-provider-azurerm/issues/23217))
* `azurerm_windows_function_app` - fix a bug in setting the storage settings when using Elastic Premium plans ([#21212](https://github.com/hashicorp/terraform-provider-azurerm/issues/21212))
* `azurerm_windows_web_app` - fix docker app stack update ([#23303](https://github.com/hashicorp/terraform-provider-azurerm/issues/23303))
* `azurerm_windows_web_app_slot` - fix docker app stack update ([#23303](https://github.com/hashicorp/terraform-provider-azurerm/issues/23303))

DEPRECATIONS:

* `azurerm_application_gateway` - deprecate `Standard` and `WAF` skus ([#23310](https://github.com/hashicorp/terraform-provider-azurerm/issues/23310))
* `azurerm_bot_channel_web_chat` - deprecate `site_names` in favour of `site` block ([#23161](https://github.com/hashicorp/terraform-provider-azurerm/issues/23161))
* `azurerm_monitor_diagnostic_setting` - deprecate `retention_policy` in favour of `azurerm_storage_management_policy` ([#23260](https://github.com/hashicorp/terraform-provider-azurerm/issues/23260))

## 3.73.0 (September 14, 2023)

FEATURES:

* **New Resource**: `azurerm_iothub_endpoint_cosmosdb_account` ([#23065](https://github.com/hashicorp/terraform-provider-azurerm/issues/23065))
* **New Resource**: `azurerm_virtual_hub_routing_intent` ([#23138](https://github.com/hashicorp/terraform-provider-azurerm/issues/23138))

ENHANCEMENTS:

* dependencies: updating to `v0.1.1` of `github.com/btubbs/datetime` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v1.3.1` of `github.com/google/uuid` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v0.61.0` of `github.com/hashicorp/go-azure-helpers` ([#23249](https://github.com/hashicorp/terraform-provider-azurerm/issues/23249))
* dependencies: updating to `v0.20230907.1113401` of `github.com/hashicorp/go-azure-sdk` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v1.5.0` of `github.com/hashicorp/go-hclog` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v2.29.0` of `github.com/hashicorp/terraform-plugin-sdk/v2` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v1.5.1` of `github.com/hashicorp/terraform-plugin-testing` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v1.20.2` of `github.com/rickb777/date` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v0.13.0` of `golang.org/x/crypto` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v0.15.0` of `golang.org/x/net` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v0.13.0` of `golang.org/x/tools` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* `azurerm_bot_channel_ms_teams` - support for `deployment_environment` ([#23122](https://github.com/hashicorp/terraform-provider-azurerm/issues/23122))
* `azurerm_managed_disk` - updating to use API Version `2023-04-02` ([#23233](https://github.com/hashicorp/terraform-provider-azurerm/issues/23233))
* `azurerm_managed_disk` - support for `optimized_frequent_attach_enabled` ([#23241](https://github.com/hashicorp/terraform-provider-azurerm/issues/23241))
* `azurerm_managed_disk` - support for `performance_plus_enabled` ([#23241](https://github.com/hashicorp/terraform-provider-azurerm/issues/23241))
* `azurerm_maps_account` - support for `local_authentication_enabled` ([#23216](https://github.com/hashicorp/terraform-provider-azurerm/issues/23216))
* `azurerm_mssql_elasticpool` - support for configuring `license_type` when using the `Hyperscale` sku ([#23256](https://github.com/hashicorp/terraform-provider-azurerm/issues/23256))
* `azurerm_security_center_assessment_policy` - refactoring to use `hashicorp/go-azure-sdk` ([#23158](https://github.com/hashicorp/terraform-provider-azurerm/issues/23158))

BUG FIXES:

* `azurerm_api_management` - split create and update methods  ([#23259](https://github.com/hashicorp/terraform-provider-azurerm/issues/23259))
* `azurerm_api_management_backend` - fixing a panic when flattening the `credentials` block ([#23219](https://github.com/hashicorp/terraform-provider-azurerm/issues/23219))
* `azurerm_key_vault_certificate` - fixing a regression where certificates from a custom/unknown issuer would be polled indefinitely ([#23214](https://github.com/hashicorp/terraform-provider-azurerm/issues/23214))
* `azurerm_redis_cache` - prevent sending `redis_configuration.aof_backup_enabled` when the sku is not `Premium` to avoid API error ([#22774](https://github.com/hashicorp/terraform-provider-azurerm/issues/22774))
* `azurerm_web_application_firewall_policy` - capture and toggle state of `custom_rule` blocks with an `enabled` field ([#23163](https://github.com/hashicorp/terraform-provider-azurerm/issues/23163))

## 3.72.0 (September 07, 2023)

FEATURES:

* Provider Feature: subscription cancellation on `destroy` can now be disabled via the provider `features` block ([#19936](https://github.com/hashicorp/terraform-provider-azurerm/issues/19936))
* **New Data Source**: `netapp_volume_quota_rule` ([#23042](https://github.com/hashicorp/terraform-provider-azurerm/issues/23042))
* **New Resource**: `azurerm_automation_python3_package` ([#23087](https://github.com/hashicorp/terraform-provider-azurerm/issues/23087))
* **New Resource**: `netapp_volume_quota_rule` ([#23042](https://github.com/hashicorp/terraform-provider-azurerm/issues/23042))

ENHANCEMENTS:

* dependencies: updating to `v0.20230906.1160501` of `github.com/hashicorp/go-azure-sdk` ([#23191](https://github.com/hashicorp/terraform-provider-azurerm/issues/23191))
* `containerapps`: updating to API Version `2023-05-01` ([#22804](https://github.com/hashicorp/terraform-provider-azurerm/issues/22804))
* `keyvault`: upgrade remaining resources to `2023-02-01` ([#23089](https://github.com/hashicorp/terraform-provider-azurerm/issues/23089))
* `redisenterprise`: updating to API Version `2023-07-01` ([#23178](https://github.com/hashicorp/terraform-provider-azurerm/issues/23178))
* `vpngateway`: updating to use `hashicorp/go-azure-sdk` ([#22906](https://github.com/hashicorp/terraform-provider-azurerm/issues/22906))
* `internal/sdk`: typed resources using a custom importer now get a timed context ([#23160](https://github.com/hashicorp/terraform-provider-azurerm/issues/23160))
* `azurerm_batch_pool` - support for `accelerated_networking_enabled` ([#23021](https://github.com/hashicorp/terraform-provider-azurerm/issues/23021))
* `azurerm_batch_pool` - support for `automatic_upgrade_enabled` ([#23021](https://github.com/hashicorp/terraform-provider-azurerm/issues/23021))
* `azurerm_bot_channel_direct_line_speech` - support for the `cognitive_account_id` property ([#23106](https://github.com/hashicorp/terraform-provider-azurerm/issues/23106))
* `azurerm_bot_service_azure_bot` - support for the `local_authentication_enabled` property ([#23096](https://github.com/hashicorp/terraform-provider-azurerm/issues/23096))
* `azurerm_container_app_environment` - support for the `dapr_application_insights_connection_string` ([#23080](https://github.com/hashicorp/terraform-provider-azurerm/issues/23080))
* `azurerm_cosmosdb_cassandra_datacenter` - refactoring to use `hashicorp/go-azure-sdk` ([#23110](https://github.com/hashicorp/terraform-provider-azurerm/issues/23110))
* `azurerm_cosmosdb_cassandra_datacenter` - updating to API Version `2023-04-15` ([#23110](https://github.com/hashicorp/terraform-provider-azurerm/issues/23110))
* `azurerm_kubernetes_cluster` - Azure CNI can be updated to use `overlay` ([#22709](https://github.com/hashicorp/terraform-provider-azurerm/issues/22709))
* `azurerm_monitor_diagnostic_setting` - deprecating `retention_policy` within `enabled_log` ([#23029](https://github.com/hashicorp/terraform-provider-azurerm/issues/23029))
* `azurerm_mssql_database` - split create and update methods ([#23209](https://github.com/hashicorp/terraform-provider-azurerm/issues/23209))
* `azurerm_postgresql_database` - `collation` can now be set to `English_United Kingdom.1252` ([#23171](https://github.com/hashicorp/terraform-provider-azurerm/issues/23171))
* `azurerm_postgresql_flexible_database` - `collation` can now be set to `English_United Kingdom.1252` ([#23171](https://github.com/hashicorp/terraform-provider-azurerm/issues/23171))
* `azurerm_postgresql_flexible_server` - support for the `auto_grow_enabled` property ([#23069](https://github.com/hashicorp/terraform-provider-azurerm/issues/23069))
* `azurerm_redis_enterprise_cluster` - support for Flash clusters in Brazil South ([#23200](https://github.com/hashicorp/terraform-provider-azurerm/issues/23200))
* `azurerm_resource_provider_registration` - refactoring to use `hashicorp/go-azure-sdk` ([#23072](https://github.com/hashicorp/terraform-provider-azurerm/issues/23072))
* `azurerm_virtual_machine_extension` - support for `provision_after_extensions` ([#23124](https://github.com/hashicorp/terraform-provider-azurerm/issues/23124))
* `azurerm_virtual_network_gateway` - increasing the default timeout for create to `90m` ([#23003](https://github.com/hashicorp/terraform-provider-azurerm/issues/23003))
* `azurerm_virtual_hub_connection` - support for `inbound_route_map_id`, `outbound_route_map_id`, and `static_vnet_local_route_override_criteria` properties ([#23049](https://github.com/hashicorp/terraform-provider-azurerm/issues/23049))

BUG FIXES:

* `azurerm_api_management_api_policy` - added state migration to mutate id's ending in `policies/policy` ([#23128](https://github.com/hashicorp/terraform-provider-azurerm/issues/23128))
* `azurerm_api_management_api_operation_policy` - added state migration to mutate id's ending in `policies/policy` ([#23128](https://github.com/hashicorp/terraform-provider-azurerm/issues/23128))
* `azurerm_api_management_product_policy` - added state migration to mutate id's ending in `policies/policy` ([#23128](https://github.com/hashicorp/terraform-provider-azurerm/issues/23128))
* `azurerm_automation_account` - fixes logic for `local_authentication_enabled` ([#23082](https://github.com/hashicorp/terraform-provider-azurerm/issues/23082))
* `azurerm_key_vault_managed_storage_account` - check id can be parsed correctly before setting it in state ([#23022](https://github.com/hashicorp/terraform-provider-azurerm/issues/23022))
* `azurerm_monitor_diagnostic_setting` - fix `enabled_log` feature flagged schema ([#23093](https://github.com/hashicorp/terraform-provider-azurerm/issues/23093))
* `azurerm_pim_active_role_assignment`: polling for the duration of the timeout, rather than a fixed 5 minute value ([#22932](https://github.com/hashicorp/terraform-provider-azurerm/issues/22932))
* `azurerm_policy_set_definition` - only sending `parameters` when a value is configured ([#23155](https://github.com/hashicorp/terraform-provider-azurerm/issues/23155))
* `azurerm_synapse_workspace` - fixes index out-of-range panic when parsing `storage_data_lake_gen2_filesystem_id` ([#23019](https://github.com/hashicorp/terraform-provider-azurerm/issues/23019))
* `machine_learning_datastore_*` - fixes container ids ([#23140](https://github.com/hashicorp/terraform-provider-azurerm/issues/23140))
* `azurerm_key_vault_certificate` - id now points to new version when certificate is updated ([#23135](https://github.com/hashicorp/terraform-provider-azurerm/issues/23135))
* `azurerm_site_recovery_replicated_vm` - update `network_interface` diff so replicated items now can be updated ([#23199](https://github.com/hashicorp/terraform-provider-azurerm/issues/23199))

DEPRECATION:

* Data Source: `azure_monitor_log_profile` - Azure is retiring Azure Log Profiles on the 30th of September 2026 ([#23146](https://github.com/hashicorp/terraform-provider-azurerm/issues/23146))
* `azure_monitor_log_profile` - Azure is retiring Azure Log Profiles on the 30th of September 2026 ([#23146](https://github.com/hashicorp/terraform-provider-azurerm/issues/23146))

## 3.71.0 (August 24, 2023)

BREAKING CHANGES:

* **App Service `win32_status` property** - Due to a change made in the service to the underlying type of the Auto Heal property `win32_status` combined with a prior bug (in `v3.62.1` and earlier) causing the value of this property to be stored incorrectly in state as an empty string, the value of this property could not be updated or state migrated to accommodate the necessary type change in the state. This results in the resources named above returning an error of a number is needed when decoding the state for this value. Unfortunately, this is a breaking change and will require users of this field to change their Terraform Configuration. The field `win32_status` has been replaced by `win32_status_code` (this remains an int, as in 3.63.0 onwards) for `azurerm_linux_web_app`, `azurerm_linux_web_app_slot`, `azurerm_windows_web_app`, `azurerm_windows_web_app_slot resources`. ([#23075](https://github.com/hashicorp/terraform-provider-azurerm/issues/23075))

FEATURES:

* **New Resource**: `azurerm_databricks_workspace_root_dbfs_customer_managed_key` ([#22579](https://github.com/hashicorp/terraform-provider-azurerm/issues/22579))

ENHANCEMENTS:

* dependencies: updating to `v0.20230824.1130652` of `github.com/hashicorp/go-azure-sdk` ([#23076](https://github.com/hashicorp/terraform-provider-azurerm/issues/23076))
* `trafficmanager`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22579](https://github.com/hashicorp/terraform-provider-azurerm/issues/22579))
* `webpubsub`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22579](https://github.com/hashicorp/terraform-provider-azurerm/issues/22579))
* `automation`: upgrade remaining resources to `2022-08-08` ([#22989](https://github.com/hashicorp/terraform-provider-azurerm/issues/22989))
* `azurerm_storage_management_policy` - move to `hashicorp/go-azure-sdk` ([#23035](https://github.com/hashicorp/terraform-provider-azurerm/issues/23035))
* Data Source: `azurerm_disk_encryption_set` -  support for the `identity` block ([#23005](https://github.com/hashicorp/terraform-provider-azurerm/issues/23005))
* `azurerm_container_group` - support for the `sku` and `(init_)container.*.security` properties ([#23034](https://github.com/hashicorp/terraform-provider-azurerm/issues/23034))
* `azurerm_kubernetes_cluster` -  extend allowed ranges for various `sysctl_config` attribute ranges ([#23077](https://github.com/hashicorp/terraform-provider-azurerm/issues/23077))
* `azurerm_kubernetes_cluster_node_pool` -  extend allowed ranges for various `sysctl_config` attribute ranges ([#23077](https://github.com/hashicorp/terraform-provider-azurerm/issues/23077))
* `azurerm_kubernetes_cluster` - clusters can be updated to use the `cilium` dataplane by setting the value in `ebpf_data_plane` ([#22952](https://github.com/hashicorp/terraform-provider-azurerm/issues/22952))
* `azurerm_linux_virtual_machine_scale_set` - cancel rolling upgrades that are in progress before destroying the resource ([#22991](https://github.com/hashicorp/terraform-provider-azurerm/issues/22991))
* `azurerm_servicebus_namespace` - support for `network_rule_set` block ([#23057](https://github.com/hashicorp/terraform-provider-azurerm/issues/23057))
* `azurerm_windows_virtual_machine_scale_set` - cancel rolling upgrades that are in progress before destroying the resource ([#22991](https://github.com/hashicorp/terraform-provider-azurerm/issues/22991))
* `azurerm_synapse_spark_pool` - support addtional values for the `node_size_family` property ([#23040](https://github.com/hashicorp/terraform-provider-azurerm/issues/23040))

BUG FIXES:

* `azurerm_api_management_policy` - fixes an error caused by a migration ([#23018](https://github.com/hashicorp/terraform-provider-azurerm/issues/23018))
* `azurerm_kubernetes_cluster` - deprecate `public_network_access_enabled` and prevent sending it to the API since it isn't functional ([#22478](https://github.com/hashicorp/terraform-provider-azurerm/issues/22478))

## 3.70.0 (August 17, 2023)

FEATURES:

* **New Resource**: `azurerm_mssql_virtual_machine_availability_group_listener` ([#22808](https://github.com/hashicorp/terraform-provider-azurerm/issues/22808))
* **New Resource**: `azurerm_mssql_virtual_machine_group` ([#22808](https://github.com/hashicorp/terraform-provider-azurerm/issues/22808))

ENHANCEMENTS:

* dependencies: updating to `v0.20230815.1165905` of `github.com/hashicorp/go-azure-sdk` ([#22981](https://github.com/hashicorp/terraform-provider-azurerm/issues/22981))
* `apimanagement`: updating to use `hashicorp/go-azure-sdk` ([#22783](https://github.com/hashicorp/terraform-provider-azurerm/issues/22783))
* `cosmos`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22874](https://github.com/hashicorp/terraform-provider-azurerm/issues/22874))
* `devtestlabs`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22981](https://github.com/hashicorp/terraform-provider-azurerm/issues/22981))
* `policy`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22874](https://github.com/hashicorp/terraform-provider-azurerm/issues/22874))
* `postgresql`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22874](https://github.com/hashicorp/terraform-provider-azurerm/issues/22874))
* `recoveryservices`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22874](https://github.com/hashicorp/terraform-provider-azurerm/issues/22874))
* `resources`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22874](https://github.com/hashicorp/terraform-provider-azurerm/issues/22874))
* `storage`: updating Storage Account and Storage Blob Container to use Common IDs to enable migrating to `hashicorp/go-azure-sdk` in the future ([#22915](https://github.com/hashicorp/terraform-provider-azurerm/issues/22915))
* Data Source: `azurerm_kubernetes_cluster` - add support for the `current_kubernetes_version` property ([#22986](https://github.com/hashicorp/terraform-provider-azurerm/issues/22986))
* `azurerm_mssql_virtual_machine` - add support for the `sql_virtual_machine_group_id` and `wsfc_domain_credential` properties ([#22808](https://github.com/hashicorp/terraform-provider-azurerm/issues/22808))
* `azurerm_netapp_pool` - `size_in_tb` can be sized down to 2 TB ([#22943](https://github.com/hashicorp/terraform-provider-azurerm/issues/22943))
* `azurerm_stack_hci_cluster` - add support for the `automanage_configuration_id` property ([#22857](https://github.com/hashicorp/terraform-provider-azurerm/issues/22857))
* Data Source: `azurerm_disk_encryption_set` - now exports `key_vault_key_url` ([#22893](https://github.com/hashicorp/terraform-provider-azurerm/issues/22893))
* `azurerm_disk_encryption_set` - now exports `key_vault_key_url` ([#22893](https://github.com/hashicorp/terraform-provider-azurerm/issues/22893))

BUG FIXES:

* `azurerm_cognitive_deployment` - add lock on parent resource to prevent errors when deleting the resource ([#22940](https://github.com/hashicorp/terraform-provider-azurerm/issues/22940))
* `azurerm_cost_management_scheduled_action` - fix update for `email_address_sender` ([#22930](https://github.com/hashicorp/terraform-provider-azurerm/issues/22930))
* `azurerm_disk_encryption_set` - now correctly supports key rotation by specifying a versionless Key ID when setting `auto_key_rotation_enabled` to `true` ([#22893](https://github.com/hashicorp/terraform-provider-azurerm/issues/22893))
* `azurerm_iothub_dps` - updating the validation for `target` within the `ip_filter_rule` block to match the values defined in the Azure API Definitions ([#22891](https://github.com/hashicorp/terraform-provider-azurerm/issues/22891))
* `azurerm_postgresql_database` - reworking the validation for database collation ([#22928](https://github.com/hashicorp/terraform-provider-azurerm/issues/22928))
* `azurerm_postgresql_flexible_database` - reworking the validation for database collation ([#22928](https://github.com/hashicorp/terraform-provider-azurerm/issues/22928))
* `azurerm_storage_management_policy` - check for an existing resource to prevent overwriting property values ([#22966](https://github.com/hashicorp/terraform-provider-azurerm/issues/22966))
* `azurerm_virtual_network_gateway_connection` - `custom_bgp_addresses.secondary` is now `Optional` rather than `Required` ([#22912](https://github.com/hashicorp/terraform-provider-azurerm/issues/22912))
* `azurerm_web_application_firewall_policy` - fix handling not found in read ([#22982](https://github.com/hashicorp/terraform-provider-azurerm/issues/22982))

---

For information on changes between the v3.69.0 and v3.0.0 releases, please see [the previous v3.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v3.md).

For information on changes between the v2.99.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v1.44.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
