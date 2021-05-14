## 2.60.0 (Unreleased)

FEATURES:

ENHANCEMENTS:

BUG FIXES:

* `azurerm_linux_virtual_machine_scale_set` - the `extension` blocks are now a set [GH-11425]
* `azurerm_windows_virtual_machine_scale_set` - the `extension` blocks are now a set [GH-11425]

## 2.59.0 (May 14, 2021)

FEATURES:

* **New Resource:** `azurerm_consumption_budget_resource_group` ([#9201](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9201))
* **New Resource:** `azurerm_consumption_budget_subscription` ([#9201](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9201))
* **New Resource:** `azurerm_monitor_aad_diagnostic_setting` ([#11660](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11660))
* **New Resource:** `azurerm_sentinel_alert_rule_machine_learning_behavior_analytics` ([#11552](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11552))
* **New Resource:** `azurerm_servicebus_namespace_disaster_recovery_config` ([#11638](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11638))

ENHANCEMENTS:

* dependencies: updating to `v54.4.0` of `github.com/Azure/azure-sdk-for-go` ([#11593](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11593))
* dependencies: updating `databox` to API version `2020-12-01` ([#11626](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11626))
* dependencies: updating `maps` to API version `2021-02-01` ([#11676](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11676))
* Data Source: `azurerm_kubernetes_cluster` - Add `ingress_application_gateway_identity` export for add-on `ingress_application_gateway` ([#11622](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11622))
* `azurerm_cosmosdb_account` - support for the `identity` and `cors_rule` blocks ([#11653](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11653))
* `azurerm_cosmosdb_account` - support for the `backup` property ([#11597](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11597))
* `azurerm_cosmosdb_sql_container` - support for the `analytical_storage_ttl` property ([#11655](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11655))
* `azurerm_container_registry` - support for the `identity` and `encryption` blocks ([#11661](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11661))
* `azurerm_frontdoor_custom_https_configuration` - Add support for resource import. ([#11642](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11642))
* `azurerm_kubernetes_cluster` - export the `ingress_application_gateway_identity` attribute for the `ingress_application_gateway` add-on ([#11622](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11622))
* `azurerm_managed_disk` - support for the `tier` property ([#11634](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11634))
* `azurerm_storage_account` - support for the `azure_files_identity_based_authentication` and `routing_preference` blocks ([#11485](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11485))
* `azurerm_storage_account` - support for the `private_link_access` property ([#11629](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11629))
* `azurerm_storage_account` - support for the `change_feed_enabled` property ([#11695](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11695))

BUG FIXES

* Data Source: `azurerm_container_registry_token` - updating the validation for the `name` field ([#11607](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11607))
* `azurerm_bastion_host` - updating the `ip_configuration` block properties now forces a new resource ([#11700](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11700))
* `azurerm_container_registry_token` - updating the validation for the `name` field ([#11607](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11607))
* `azurerm_mssql_database` - wil now correctly import the `creation_source_database_id` property for Secondary databases ([#11703](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11703))
* `azurerm_storage_account` - allow empty/blank values for the `allowed_headers` and `exposed_headers` properties ([#11692](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11692))

## 2.58.0 (May 07, 2021)

UPGRADE NOTES

* `azurerm_frontdoor` - The `custom_https_provisioning_enabled` field and the `custom_https_configuration` block have been deprecated and has been removed as they are no longer supported. ([#11456](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11456))
* `azurerm_frontdoor_custom_https_configuration` - The `resource_group_name` has been deprecated and has been removed as it is no longer supported. ([#11456](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11456))

FEATURES:

* **New Data Source:** `azurerm_storage_table_entity` ([#11562](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11562))
* **New Resource:** `azurerm_app_service_environment_v3` ([#11174](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11174))
* **New Resource:** `azurerm_cosmosdb_notebook_workspace` ([#11536](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11536))
* **New Resource:** `azurerm_cosmosdb_sql_trigger` ([#11535](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11535))
* **New Resource:** `azurerm_cosmosdb_sql_user_defined_function` ([#11537](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11537))
* **New Resource:** `azurerm_iot_time_series_insights_event_source_iothub` ([#11484](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11484))
* **New Resource:** `azurerm_storage_blob_inventory_policy` ([#11533](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11533))

ENHANCEMENTS:

* dependencies: updating `network-db` to API version `2020-07-01` ([#10767](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10767))
* `azurerm_cosmosdb_account` - support for the `access_key_metadata_writes_enabled`, `mongo_server_version`, and `network_acl_bypass` properties ([#11486](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11486))
* `azurerm_data_factory` - support for the `customer_managed_key_id` property ([#10502](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10502))
* `azurerm_data_factory_pipeline` - support for the `folder` property ([#11575](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11575))
* `azurerm_frontdoor` - Fix for Frontdoor resource elements being returned out of order. ([#11456](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11456))
* `azurerm_hdinsight_*_cluster` - support for autoscale  #8104 ([#11547](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11547))
* `azurerm_network_security_rule` - support for the protocols `Ah` and `Esp` ([#11581](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11581))
* `azurerm_network_connection_monitor` - support for the `coverage_level`, `excluded_ip_addresses`, `included_ip_addresses`, `target_resource_id`, and `resource_type` propeties ([#11540](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11540))

## 2.57.0 (April 30, 2021)

UPGRADE NOTES

* `azurerm_api_management_authorization_server` - due to a bug in the `2020-12-01` version of the API Management API, changes to `resource_owner_username` and `resource_owner_password` in Azure will not be noticed by Terraform ([#11146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11146))
* `azurerm_cosmosdb_account` - the `2021-02-01` version of the cosmos API defaults new MongoDB accounts to `v3.6` rather then `v3.2` ([#10926](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10926))
* `azurerm_cosmosdb_mongo_collection` - the `_id` index is now required by the new API/MongoDB version ([#10926](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10926))
* `azurerm_cosmosdb_gremlin_graph` and `azurerm_cosmosdb_sql_container` - the `patition_key_path` property is now required ([#10926](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10926))
 
FEATURES:

* **Data Source:** `azurerm_container_registry_scope_map`  ([#11350](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11350))
* **Data Source:** `azurerm_container_registry_token`  ([#11350](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11350))
* **Data Source:** `azurerm_postgresql_flexible_server` ([#11081](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11081))
* **Data Source:** `azurerm_key_vault_managed_hardware_security_module` ([#10873](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10873))
* **New Resource:** `azurerm_container_registry_scope_map`  ([#11350](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11350))
* **New Resource:** `azurerm_container_registry_token`  ([#11350](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11350))
* **New Resource:** `azurerm_data_factory_dataset_snowflake `  ([#11116](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11116))
* **New Resource:** `azurerm_healthbot` ([#11002](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11002))
* **New Resource:** `azurerm_key_vault_managed_hardware_security_module `  ([#10873](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10873))
* **New Resource:** `azurerm_media_asset_filter` ([#11110](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11110))
* **New Resource:** `azurerm_mssql_job_agent` ([#11248](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11248))
* **New Resource:** `azurerm_mssql_job_credential` ([#11363](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11363))
* **New Resource:** `azurerm_mssql_transparent_data_encryption` ([#11148](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11148))
* **New Resource:** `azurerm_postgresql_flexible_server` ([#11081](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11081))
* **New Resource:** `azurerm_spring_cloud_app_cosmosdb_association` ([#11307](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11307))
* **New Resource:** `azurerm_sentinel_data_connector_microsoft_defender_advanced_threat_protection` ([#10669](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10669))
* **New Resource:** `azurerm_virtual_machine_configuration_policy_assignment` ([#11334](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11334))
* **New Resource:** `azurerm_vmware_cluster` ([#10848](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10848))

ENHANCEMENTS:

* dependencies: updating to `v53.4.0` of `github.com/Azure/azure-sdk-for-go` ([#11439](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11439))
* dependencies: updating to `v1.17.2` of `github.com/hashicorp/terraform-plugin-sdk` ([#11431](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11431))
* dependencies: updating `cosmos-db` to API version `2021-02-01` ([#10926](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10926))
* dependencies: updating `keyvault` to API version `v7.1` ([#10926](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10926))
* Data Source: `azurerm_healthcare_service` - export the `cosmosdb_key_vault_key_versionless_id` attribute ([#11481](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11481))
* Data Source: `azurerm_key_vault_certificate` - export the `curve` attribute in the `key_properties` block ([#10867](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10867))
* Data Source: `azurerm_virtual_machine_scale_set` - now exports the `network_interfaces` ([#10585](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10585))
* `azurerm_app_service` - support for the `site_config.ip_restrictions.headers` and  `site_config.scm_ip_restrictions.headers` properties ([#11209](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11209))
* `azurerm_app_service_slot` - support for the `site_config.ip_restrictions.headers` and  `site_config.scm_ip_restrictions.headers` properties ([#11209](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11209))
* `azurerm_backup_policy_file_share` - support for the `retention_weekly`, `retention_monthly`, and `retention_yearly` blocks ([#10733](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10733))
* `azurerm_cosmosdb_sql_container` - support for the `conflict_resolution_policy` block ([#11517](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11517))
* `azurerm_container_group` - support for the `exposed_port` block ([#10491](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10491))
* `azurerm_container_registry` - deprecating the `georeplication_locations` property in favour of the `georeplications` property [#11200](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11200)]
* `azurerm_database_migration` - switching to using an ID Formatter ([#11378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11378))
* `azurerm_database_migration_project` - switching to using an ID Formatter ([#11378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11378))
* `azurerm_databricks_workspace` - switching to using an ID Formatter ([#11378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11378))
* `azurerm_databricks_workspace` - fixes propagation of tags to connected resources ([#11405](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11405))
* `azurerm_data_factory_linked_service_azure_file_storage` - support for the `key_vault_password` property ([#11436](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11436))
* `azurerm_dedicated_host_group` - support for the `automatic_placement_enabled` property ([#11428](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11428))
* `azurerm_frontdoor` - sync `MaxItems` on various attributes to match azure docs ([#11421](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11421))
* `azurerm_frontdoor_custom_https_configuration` - removing secret version validation when using azure key vault as the certificate source ([#11310](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11310))
* `azurerm_function_app` - support for the `site_config.ip_restrictions.headers` and  `site_config.scm_ip_restrictions.headers` properties ([#11209](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11209))
* `azurerm_function_app` - support the `java_version` property ([#10495](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10495))
* `azurerm_hdinsight_interactive_query_cluster` - add support for private link endpoint ([#11300](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11300))
* `azurerm_hdinsight_hadoop_cluster` - add support for private link endpoint ([#11300](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11300))
* `azurerm_hdinsight_spark_cluster` - add support for private link endpoint ([#11300](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11300))
* `azurerm_healthcare_service` - support for the `cosmosdb_key_vault_key_versionless_id` property ([#11481](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11481))
* `azurerm_kubernetes_cluster` - support for the `ingress_application_gateway` addon ([#11376](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11376))
* `azurerm_kubernetes_cluster` - support for the `azure_rbac_enabled` property ([#10441](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10441))
* `azurerm_hpc_cache` - support for the `directory_active_directory`, `directory_flat_file`, and `directory_ldap` blocks ([#11332](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11332))
* `azurerm_key_vault_certificate` - support additional values for the `key_size` property in the `key_properties` block ([#10867](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10867))
* `azurerm_key_vault_certificate` - support the `curve` property in the `key_properties` block ([#10867](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10867))
* `azurerm_key_vault_certificate` - the `key_size` property in the `key_properties` block is now optional ([#10867](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10867))
* `azurerm_kubernetes_cluster` - support for the `dns_prefix_private_cluster` property ([#11321](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11321))
* `azurerm_kubernetes_cluster` - support for the `max_node_provisioning_time`, `max_unready_percentage`, and `max_unready_nodes` properties ([#11406](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11406))
* `azurerm_storage_encryption_scope` - support for the `infrastructure_encryption_required` property ([#11462](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11462))
* `azurerm_kubernetes_cluster` support for the `empty_bulk_delete_max` in the `auto_scaler_profile` block #([#11060](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11060))
* `azurerm_lighthouse_definition` - support for the `delegated_role_definition_ids` property ([#11269](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11269))
* `azurerm_managed_application` - support for the `parameter_values` property ([#8632](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8632))
* `azurerm_managed_disk` - support for the `network_access_policy` and `disk_access_id` properties ([#9862](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9862))
* `azurerm_postgresql_server` - wait for replica restarts when needed ([#11458](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11458))
* `azurerm_redis_enterprise_cluster` - support for the `minimum_tls_version` and `hostname` properties ([#11203](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11203))
* `azurerm_storage_account` -  support for the `versioning_enabled`, `default_service_version`, and `last_access_time_enabled` properties within the `blob_properties` block ([#11301](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11301))
* `azurerm_storage_account` - support for the `nfsv3_enabled` property ([#11387](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11387))
* `azurerm_storage_management_policy` - support for the `version` block ([#11163](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11163))
* `azurerm_synapse_workspace` - support for the `customer_managed_key_versionless_id` property ([#11328](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11328))

BUG FIXES:

* `azurerm_api_management` - will no longer panic with an empty `hostname_configuration` ([#11426](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11426))
* `azurerm_api_management_diagnostic` - fix a crash with the `frontend_request`, `frontend_response`, `backend_request`, `backend_response` blocks ([#11402](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11402))
* `azurerm_eventgrid_system_topic` - remove strict validation on `topic_type` ([#11352](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11352))
* `azurerm_iothub` - change `filter_rule` from TypeSet to TypeList to resolve an ordering issue ([#10341](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10341))
* `azurerm_linux_virtual_machine_scale_set` - the default value for the `priority` property will no longer force a replacement of the resource ([#11362](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11362))
* `azurerm_monitor_activity_log_alert` - fix a persistent diff for the `service_health` block ([#11383](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11383))
* `azurerm_mssql_database ` - return an error when secondary database uses `max_size_gb` ([#11401](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11401))
* `azurerm_mssql_database` - correctly import the `create_mode` property ([#11026](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11026))
* `azurerm_netap_volume` - correctly set the `replication_frequency` attribute in the `data_protection_replication` block ([#11530](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11530))
* `azurerm_postgresql_server` - ensure `public_network_access_enabled` is correctly set for replicas ([#11465](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11465))
* `azurerm_postgresql_server` - can now correctly disable replication if required when `create_mode` is changed ([#11467](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11467))
* `azurerm_virtual_network_gatewa` - updating the `custom_route` block no longer forces a new resource to be created [GH- 11433]

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
