## 2.57.0 (Unreleased)

FEATURES:

* **Data Source:** `azurerm_postgresql_flexible_server` [GH-11081]
* **Data Source:** `azurerm_key_vault_managed_hardware_security_module` [GH-10873]
* **New Resource:** `azurerm_data_factory_dataset_snowflake `  [GH-11116]
* **New Resource:** `azurerm_key_vault_managed_hardware_security_module `  [GH-10873]
* **New Resource:** `azurerm_mssql_job_agent` [GH-11248]
* **New Resource:** `azurerm_mssql_transparent_data_encryption` [GH-11148]
* **New Resource:** `azurerm_postgresql_flexible_server` [GH-11081]
* **New Resource:** `azurerm_spring_cloud_app_cosmosdb_association` [GH-11307]
* **New Resource:** `azurerm_sentinel_data_connector_microsoft_defender_advanced_threat_protection` [GH-10669]

ENHANCEMENTS:

* dependencies: updating to `v53.3.0` of `github.com/Azure/azure-sdk-for-go` [GH-11365]
* `azurerm_container_registry` - deprecating the `georeplication_locations` property in favour of the `georeplications` property GH-11200]
* `azurerm_database_migration` - switching to using an ID Formatter [GH-11378]
* `azurerm_database_migration_project` - switching to using an ID Formatter [GH-11378]
* `azurerm_databricks_workspace` - switching to using an ID Formatter [GH-11378]
* `azurerm_databricks_workspace` - fixes propagation of tags to connected resources [GH-11405]
* `azurerm_data_factory_linked_service_azure_file_storage` - support for the `key_vault_password` property [GH-11436]
* `azurerm_frontdoor` - sync `MaxItems` on various attributes to match azure docs [GH-11421]
* `azurerm_hdinsight_interactive_query_cluster` - add support for private link endpoint [GH-11300]
* `azurerm_hdinsight_hadoop_cluster` - add support for private link endpoint [GH-11300]
* `azurerm_hdinsight_spark_cluster` - add support for private link endpoint [GH-11300]
* `azurerm_kubernetes_cluster` support for the `empty_bulk_delete_max` in the `auto_scaler_profile` block #[GH-11060]
* `azurerm_lighthouse_definition` - support for the `delegated_role_definition_ids` property [GH-11269]
* `azurerm_redis_enterprise_cluster` - support for the `minimum_tls_version` and `hostname` properties [GH-11203]

BUG FIXES:

* `azurerm_api_management` - will no longer panic with an empty `hostname_configuration` [GH-11426]
* `azurerm_api_management_diagnostic` - fix a crash with the `frontend_request`, `frontend_response`, `backend_request`, `backend_response` blocks [GH-11402]
* `azurerm_eventgrid_system_topic` - remove strict validation on `topic_type` [GH-11352]
* `azurerm_iothub` - change `filter_rule` from TypeSet to TypeList to resolve an ordering issue [GH-10341]
* `azurerm_linux_virtual_machine_scale_set` - the default value for the `priority` property will no longer force a replacement of the resource [GH-11362]
* `azurerm_monitor_activity_log_alert` - fix a persistent diff for the `service_health` block [GH-11383]
* `azurerm_mssql_database ` - error when secondary database uses `max_size_gb` [GH-11401]
* `azurerm_virtual_network_gatewa` - updating the `custom_route` block no longer creates a new resource [GH- 11433]

## 2.56.0 (April 15, 2021)

FEATURES:

* **New Resource:** `azurerm_data_factory_linked_service_azure_databricks` ([#10962](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10962))
* **New Resource:** `azurerm_data_lake_store_virtual_network_rule` ([#10430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10430))
* **New Resource:** `azurerm_media_live_event_output` ([#10917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10917))
* **New Resource:** `azurerm_spring_cloud_app_mysql_association` ([#11229](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11229))

ENHANCEMENTS:

* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v53.0.0` ([#11302](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11302))
* dependencies: updating `containerservice` to API version `2021-02-01` ([#10972](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10972))
* `azurerm_app_service` - fix broken `ip_restrictions` and `scm_ip_restrictions` ([#11170](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11170))
* `azurerm_application_gateway` - support for configuring `firewall_policy_id` within the `path_rule` block ([#11239](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11239))
* `azurerm_firewall_policy_rule_collection_group` - allow `*` for the `network_rule_collection.destination_ports` property ([#11326](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11326))
* `azurerm_function_app` - fix broken `ip_restrictions` and `scm_ip_restrictions` ([#11170](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11170))
* `azurerm_data_factory_linked_service_sql_database` - support managed identity and service principal auth and add the `keyvault_password` property ([#10735](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10735))
* `azurerm_hpc_cache` - support for `tags` ([#11268](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11268))
* `azurerm_linux_virtual_machine_scale_set` - Support health extension for rolling ugrade mode ([#9136](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9136))
* `azurerm_monitor_activity_log_alert` - support for `service_health` ([#10978](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10978))
* `azurerm_mssql_database` - support for the `geo_backup_enabled` property ([#11177](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11177))
* `azurerm_public_ip` - support for `ip_tags` ([#11270](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11270))
* `azurerm_windows_virtual_machine_scale_set` - Support health extension for rolling ugrade mode ([#9136](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9136))

BUG FIXES:

* `azurerm_app_service_slot` - fix crash bug when given empty `http_logs` ([#11267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11267))

## 2.55.0 (April 08, 2021)

FEATURES:

* **New Resource:** `azurerm_api_management_email_template` ([#10914](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10914))
* **New Resource:** `azurerm_communication_service` ([#11066](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11066))
* **New Resource:** `azurerm_express_route_port` ([#10074](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10074))
* **New Resource:** `azurerm_spring_cloud_app_redis_association` ([#11154](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11154))

ENHANCEMENTS:

* Data Source: `azurerm_user_assigned_identity` - exporting `tenant_id` ([#11253](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11253))
* Data Source: `azurerm_function_app` - exporting `client_cert_mode` ([#11161](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11161))
* `azurerm_eventgrid_data_connection` - support for the `table_name`, `mapping_rule_name`, and `data_format` properties ([#11157](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11157))
* `azurerm_hpc_cache` - support for configuring `dns` ([#11236](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11236))
* `azurerm_hpc_cache` - support for configuring `ntp_server` ([#11236](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11236))
* `azurerm_hpc_cache_nfs_target` - support for the `access_policy_name` property ([#11186](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11186))
* `azurerm_hpc_cache_nfs_target` - `usage_model` can now be set to `READ_HEAVY_CHECK_180`, `WRITE_WORKLOAD_CHECK_30`, `WRITE_WORKLOAD_CHECK_60` and `WRITE_WORKLOAD_CLOUDWS` ([#11247](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11247))
* `azurerm_function_app` - support for configuring `client_cert_mode` ([#11161](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11161))
* `azurerm_netapp_volume` - adding `root_access_enabled` to the `export_policy_rule` block ([#11105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11105))
* `azurerm_private_endpoint` - allows for an alias to specified ([#10779](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10779))
* `azurerm_user_assigned_identity` - exporting `tenant_id` ([#11253](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11253))
* `azurerm_web_application_firewall_policy` - `version` within the `managed_rule_set` block can now be set to (OWASP) `3.2` ([#11244](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11244))

BUG FIXES:

* Data Source: `azurerm_dns_zone` - fixing a bug where the Resource ID wouldn't contain the Resource Group name when looking this up ([#11221](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11221))
* `azurerm_media_service_account` - `storage_authentication_type` correctly accepts both `ManagedIdentity` and `System` ([#11222](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11222))
* `azurerm_web_application_firewall_policy` - `http_listener_ids` and `path_based_rule_ids` are now Computed only ([#11196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11196))

## 2.54.0 (April 02, 2021)

FEATURES:

* **New Resource:** `azurerm_hpc_cache_access_policy` ([#11083](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11083))
* **New Resource:** `azurerm_management_group_subscription_association` ([#11069](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11069))
* **New Resource:** `azurerm_media_live_event` ([#10724](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10724))

ENHANCEMENTS:

* dependencies: updating to `v52.6.0` of `github.com/Azure/azure-sdk-for-go` ([#11108](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11108))
* dependencies: updating `storage` to API version `2021-01-01` ([#11094](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11094))
* dependencies: updating `storagecache` (a.k.a `hpc`) to API version `2021-03-01` ([#11083](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11083))
* `azurerm_application_gateway` - support for rewriting urls with the `url` block ([#10950](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10950))
* `azurerm_cognitive_account` - Add support for `network_acls` ([#11164](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11164))
* `azurerm_container_registry` - support for the `quarantine_policy_enabled` property ([#11011](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11011))
* `azurerm_firewall` - support for the `private_ip_ranges` property [p[#10627](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10627)]
* `azurerm_log_analytics_workspace` - Fix issue where -1 couldn't be specified for `daily_quota_gb` ([#11182](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11182))
* `azurerm_spring_cloud_service` - supports for the `sample_rate` property ([#11106](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11106))
* `azurerm_storage_account` - support for the `container_delete_retention_policy` property ([#11131](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11131))
* `azurerm_virtual_desktop_host_pool` - support for the `custom_rdp_properties` property ([#11160](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11160))
* `azurerm_web_application_firewall_policy` - support for the `http_listener_ids` and `path_based_rule_ids` properties ([#10860](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10860))

BUG FIXES:

* `azurerm_api_management` - the `certificate_password` property is now optional ([#11139](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11139))
* `azurerm_data_factory_linked_service_azure_blob_storage` - correct managed identity implementation by implementing the `service_endpoint` property ([#10830](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10830))
* `azurerm_machine_learning_workspace` - deprecate the `Enterprise` sku as it has been deprecated by Azure ([#11063](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11063))
* `azurerm_machine_learning_workspace` - support container registries in other subscriptions ([#11065](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11065))
* `azurerm_site_recovery_fabric` - Fixes error in checking for existing resource ([#11130](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11130))
* `azurerm_spring_cloud_custom_domain` - `thumbprint` is required when specifying `certificate_name` ([#11145](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11145))
* `azurerm_subscription` - fixes broken timeout on destroy ([#11124](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11124))

## 2.53.0 (March 26, 2021)

FEATURES:

* **New Resource:** `azurerm_management_group_template_deployment` ([#10603](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10603))
* **New Resource:** `azurerm_tenant_template_deployment` ([#10603](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10603))
* **New Data Source:** `azurerm_template_spec_version` ([#10603](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10603))

ENHANCEMENTS:

* dependencies: updating to `v52.5.0` of `github.com/Azure/azure-sdk-for-go` ([#11015](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11015))
* Data Source: `azurerm_key_vault_secret` - support for the `versionless_id` attribute ([#11091](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11091))
* `azurerm_container_registry` - support for the `public_network_access_enabled` property ([#10969](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10969))
* `azurerm_kusto_eventhub_data_connection` - support for the `event_system_properties` block ([#11006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11006))
* `azurerm_logic_app_trigger_recurrence` - Add support for `schedule`  ([#11055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11055))
* `azurerm_resource_group_template_deployment` - add support for `template_spec_version_id` property ([#10603](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10603))
* `azurerm_role_definition` - the `permissions` block is now optional ([#9850](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9850))
* `azurerm_subscription_template_deployment` - add support for `template_spec_version_id` property ([#10603](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10603))


BUG FIXES:

* `azurerm_frontdoor_custom_https_configuration` - fixing a crash during update ([#11046](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11046))
* `azurerm_resource_group_template_deployment` - always sending `parameters_content` during an update ([#11001](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11001))
* `azurerm_role_definition` - fixing crash when permissions are empty ([#9850](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9850))
* `azurerm_subscription_template_deployment` - always sending `parameters_content` during an update ([#11001](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11001))
* `azurerm_spring_cloud_app` - supports for the `tls_enabled` property ([#11064](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11064))

## 2.52.0 (March 18, 2021)

FEATURES:

* **New Resource:** `azurerm_mssql_firewall_rule` ([#10954](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10954))
* **New Resource:** `azurerm_mssql_virtual_network_rule` ([#10954](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10954))

ENHANCEMENTS:

* dependencies: updating to `v52.4.0` of `github.com/Azure/azure-sdk-for-go` ([#10982](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10982))
* `azurerm_api_management_subscription` - making `user_id` property optional [[#10638](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10638)}

BUG FIXES:

* `azurerm_cosmosdb_account_resource` - marking `connection_string` as sensitive ([#10942](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10942))
*  `azurerm_eventhub_namespace_disaster_recovery_config` - deprecating the `alternate_name` property due to a service side API bug ([#11013](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11013))
* `azurerm_local_network_gateway` - making the `address_space` property optional ([#10983](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10983))
* `azurerm_management_group` - validation for `subscription_id` list property entries ([#10948](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10948))

## 2.51.0 (March 12, 2021)

FEATURES:

* **New Resource:** `azurerm_purview_account` ([#10395](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10395))
* **New Resource:** `azurerm_data_factory_dataset_parquet` ([#10852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10852))
* **New Resource:** `azurerm_security_center_server_vulnerability_assessment` ([#10030](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10030))
* **New Resource:** `azurerm_security_center_assessment` ([#10694](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10694))
* **New Resource:** `azurerm_security_center_assessment_policy` ([#10694](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10694))
* **New Resource:** `azurerm_sentinel_data_connector_azure_advanced_threat_protection` ([#10666](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10666))
* **New Resource:** `azurerm_sentinel_data_connector_azure_security_center` ([#10667](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10667))
* **New Resource:** `azurerm_sentinel_data_connector_microsoft_cloud_app_security` ([#10668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10668))

ENHANCEMENTS:

* dependencies: updating to v52.3.0 of `github.com/Azure/azure-sdk-for-go` ([#10829](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10829))
* `azurerm_role_assignment` - support enrollment ids in `scope` argument ([#10890](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10890))
* `azurerm_kubernetes_cluster` - support `None` for the `private_dns_zone_id` property ([#10774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10774))
* `azurerm_kubernetes_cluster` - support for `expander` in the `auto_scaler_profile` block ([#10777](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10777))
* `azurerm_linux_virtual_machine` - support for configuring `platform_fault_domain` ([#10803](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10803))
* `azurerm_linux_virtual_machine_scale_set` - will no longer recreate the resource when `rolling_upgrade_policy` or `health_probe_id` is updated ([#10856](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10856))
* `azurerm_netapp_volume` - support creating from a snapshot via the `create_from_snapshot_resource_id` property ([#10906](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10906))
* `azurerm_role_assignment` - support for the `description`, `condition`, and `condition_version` ([#10804](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10804))
* `azurerm_windows_virtual_machine` - support for configuring `platform_fault_domain` ([#10803](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10803))
* `azurerm_windows_virtual_machine_scale_set` - will no longer recreate the resource when `rolling_upgrade_policy` or `health_probe_id` is updated ([#10856](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10856))

BUG FIXES:

* Data Source: `azurerm_function_app_host_keys` - retrying reading the keys to work around a broken API ([#10894](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10894))
* Data Source: `azurerm_log_analytics_workspace` - ensure the `id` is returned with the correct casing ([#10892](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10892))
* Data Source: `azurerm_monitor_action_group` - add support for `aad_auth` attribute ([#10876](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10876))
* `azurerm_api_management_custom_domain` - prevent a perpetual diff ([#10636](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10636))
* `azurerm_eventhub_consumer_group` - detecting as removed when deleted in Azure ([#10900](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10900))
* `azurerm_key_vault_access_policy` - Fix destroy where permissions casing on service does not match config / state ([#10931](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10931))
* `azurerm_key_vault_secret` - setting the value of the secret after recovering it ([#10920](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10920))
* `azurerm_kusto_eventhub_data_connection` - make `table_name` and `data_format` optional ([#10913](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10913))
* `azurerm_mssql_virtual_machine` - workaround for inconsistent API value for `log_backup_frequency_in_minutes` in the `manual_schedule` block ([#10899](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10899))
* `azurerm_postgres_server` - support for replicaset scaling ([#10754](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10754))
* `azurerm_postgresql_aad_administrator` - prevent invalid usernames for the `login` property ([#10757](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10757))

## 2.50.0 (March 05, 2021)

FEATURES:

* **New Data Source:** `azurerm_vmware_private_cloud` ([#9284](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9284))
* **New Resource:** `azurerm_kusto_eventgrid_data_connection` ([#10712](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10712))
* **New Resource:** `azurerm_sentinel_data_connector_aws_cloud_trail` ([#10664](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10664))
* **New Resource:** `azurerm_sentinel_data_connector_azure_active_directory` ([#10665](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10665))
* **New Resource:** `azurerm_sentinel_data_connector_office_365` ([#10671](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10671))
* **New Resource:** `azurerm_sentinel_data_connector_threat_intelligence` ([#10670](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10670))
* **New Resource:** `azurerm_subscription` ([#10718](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10718))
* **New Resource:** `azurerm_vmware_private_cloud` ([#9284](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9284))

ENHANCEMENTS:
* dependencies: updating to `v52.0.0` of `github.com/Azure/azure-sdk-for-go` ([#10787](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10787))
* dependencies: updating `compute` to API version `2020-12-01` ([#10650](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10650))
* Data Source: `azurerm_dns_zone` - updating to use a consistent Terraform Resource ID to avoid API issues ([#10786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10786))
* `azurerm_dns_a_record` - updating to use a consistent Terraform Resource ID to avoid API issues ([#10786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10786))
* `azurerm_dns_aaaa_record` - updating to use a consistent Terraform Resource ID to avoid API issues ([#10786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10786))
* `azurerm_dns_caa_record` - updating to use a consistent Terraform Resource ID to avoid API issues ([#10786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10786))
* `azurerm_dns_cname_record` - updating to use a consistent Terraform Resource ID to avoid API issues ([#10786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10786))
* `azurerm_dns_mx_record` - updating to use a consistent Terraform Resource ID to avoid API issues ([#10786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10786))
* `azurerm_dns_ns_record` - updating to use a consistent Terraform Resource ID to avoid API issues ([#10786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10786))
* `azurerm_dns_ptr_record` - updating to use a consistent Terraform Resource ID to avoid API issues ([#10786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10786))
* `azurerm_dns_srv_record` - updating to use a consistent Terraform Resource ID to avoid API issues ([#10786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10786))
* `azurerm_dns_txt_record` - updating to use a consistent Terraform Resource ID to avoid API issues ([#10786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10786))
* `azurerm_dns_zone` - updating to use a consistent Terraform Resource ID to avoid API issues ([#10786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10786))
* `azurerm_function_app_host_keys` - support for `event_grid_extension_config_key` ([#10823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10823))
* `azurerm_keyvault_secret` - support for the `versionless_id` property ([#10738](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10738))
* `azurerm_kubernetes_cluster` - support `private_dns_zone_id` when using a `service_principal` ([#10737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10737))
* `azurerm_kusto_cluster` - supports for the `double_encryption_enabled` property ([#10264](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10264))
* `azurerm_linux_virtual_machine` - support for configuring `license_type` ([#10776](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10776))
* `azurerm_log_analytics_workspace_resource` - support permanent deletion of workspaces with the `permanently_delete_on_destroy` feature flag ([#10235](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10235))
* `azurerm_monitor_action_group` - support for secure webhooks via the `aad_auth` block ([#10509](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10509))
* `azurerm_mssql_database` - support for the `log_monitoring_enabled` property within the `extended_auditing_policy` block ([#10324](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10324))
* `azurerm_mssql_database_extended_auditing_policy ` - support for the `log_monitoring_enabled` property ([#10324](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10324))
* `azurerm_mssql_server` - support for the `log_monitoring_enabled` property within the `extended_auditing_policy` block ([#10324](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10324))
* `azurerm_mssql_server_extended_auditing_policy ` - support for the `log_monitoring_enabled` property [[#10324](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10324)] 
* `azurerm_signalr_service` - support for the `upstream_endpoint` block ([#10459](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10459))
* `azurerm_sql_server` - support for the `log_monitoring_enabled` property within the `extended_auditing_policy` block ([#10324](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10324))
* `azurerm_sql_database` - support for the `log_monitoring_enabled` property within the `extended_auditing_policy` block ([#10324](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10324))
* `azurerm_spring_cloud_java_deployment` - supporting delta updates ([#10729](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10729))
* `azurerm_virtual_network_gateway` - deprecate `peering_address` in favour of `peering_addresses` ([#10381](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10381))

BUG FIXES:

* Data Source: `azurerm_netapp_volume` - fixing a crash when setting `data_protection_replication` ([#10795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10795))
* `azurerm_api_management` - changing the `sku_name` property no longer forces a new resouce to be created ([#10747](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10747))
* `azurerm_api_management` - the field `tenant_access` can only be configured when not using a Consumption SKU ([#10766](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10766))
* `azurerum_frontdoor` - removed the MaxItems validation from the Backend Pools ([#10828](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10828))
* `azurerm_kubernetes_cluster_resource` - allow windows passwords as short as `8` charaters long ([#10816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10816))
* `azurerm_cosmosdb_mongo_collection` - ignore throughput if Cosmos DB provisioned in 'serverless' capacity mode ([#10389](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10389))
* `azurerm_linux_virtual_machine` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue ([#10722](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10722))
* `azurerm_linux_virtual_machine_scale_set` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue ([#10722](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10722))
* `azurerm_netapp_volume` - fixing a crash when setting `data_protection_replication` ([#10795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10795))
* `azurerm_virtual_machine` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue ([#10722](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10722))
* `azurerm_virtual_machine_scale_set` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue ([#10722](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10722))
* `azurerm_windows_virtual_machine` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue ([#10722](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10722))
* `azurerm_windows_virtual_machine_scale_set` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue ([#10722](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10722))

---

For information on changes between the v2.49.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
