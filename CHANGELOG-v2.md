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

## 2.49.0 (February 26, 2021)

FEATURES:

* **New Data Source:** `azurerm_spring_cloud_app` ([#10678](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10678))
* **New Resource:** `azurerm_databox_edge_device` ([#10730](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10730))
* **New Resource:** `azurerm_databox_edge_order` ([#10730](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10730))
* **New Resource:** `azurerm_kusto_iothub_data_connection` ([#8626](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8626))
* **New Resource:** `azurerm_redis_enterprise_cluster` ([#10706](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10706))
* **New Resource:** `azurerm_redis_enterprise_database` ([#10706](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10706))
* **New Resource:** `azurerm_security_center_assessment_metadata` ([#10124](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10124))
* **New Resource:** `azurerm_spring_cloud_custom_domain` ([#10404](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10404))

ENHANCEMENTS:

* dependencies: updating `github.com/hashicorp/terraform-plugin-sdk` to the latest `1.x` branch ([#10692](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10692))
* dependencies: updating `github.com/hashicorp/go-azure-helpers` to `v0.14.0` ([#10740](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10740))
* dependencies: updating `github.com/Azure/go-autorest/autorest` to `v0.11.18` ([#10740](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10740))
* testing: updating the tests to use the Terraform release binaries when running acceptance tests ([#10523](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10523))
* `azurerm_api_management` - support for the  `tenant_access` block ([#10475](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10475))
* `azurerm_api_management_logger` - support for configuring a `resource_id` ([#10652](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10652))
* `azurerm_data_factory_linked_service_azure_blob_storage` - now supports the `sas_uri` property ([#10551](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10551))
* `azurerm_data_factory_linked_service_azure_blob_storage` - now supports Managed Identity and Service Principal authentication ([#10551](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10551))
* `azurerm_monitor_smart_detector_alert_rule` - supports for the `tags` property ([#10646](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10646))
* `azurerm_netapp_volume` - support for the `data_protection_replication` block ([#10610](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10610))
* `azurerm_sentinel_alert_rule_ms_security_incident` - support `Microsoft Defender Advanced Threat Protection` and `Office 365 Advanced Threat Protection` values for the `product_filter` property ([#10725](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10725))
* `azurerm_service_fabric_cluster` - Add support for the `upgrade policy` block ([#10713](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10713))

BUG FIXES:

* provider: fixing support for Azure Cloud Shell ([#10740](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10740))
* provider: MSI authentication is explicitly unavailable in Azure App Service and Function Apps as these are intentionally not supported ([#10740](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10740))
* provider: only showing the deprecation message if `skip_credentials_registration` is explicitly configured ([#10699](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10699))
* `azurerm_batch_certificate` - allow empty `password` when format is pfx ([#10642](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10642))
* `azurerm_data_factory_integration_runtime_azure_ssis` - the `administrator_login` and `administrator_password` properties are now optional ([#10474](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10474))
* `azurerm_data_factory_integration_runtime_managed` - the `administrator_login` and `administrator_password` properties are now optional ([#10640](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10640))
* `azurerm_eventhub_namespace` - the `capacity` property can now be greater than `50` ([#10734](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10734))
* `azurerm_key_vault_certificate` - waiting for deletion to complete before purging ([#10577](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10577))
* `azurerm_key_vault_key` - now waits for deletion to complete before purging ([#10577](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10577))
* `azurerm_key_vault_secret` - now waits for deletion to complete before purging ([#10577](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10577))
* `azurerm_kusto_cluster` - changing the `virtual_network_configuration` property forces a new resource to be created ([#10640](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10640))
* `azurerm_lb_outbound_rule` - fixing a crash when `frontendIPConfigurations` is omitted in the API response ([#10696](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10696))
* `azurerm_media_content_key_policy` - fix an encoding bug which prevented configuring `ask` in the `fairplay_configuration` block ([#10684](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10684))

## 2.48.0 (February 18, 2021)

FEATURES:

* **New Data Source:** `azurerm_application_gateway` ([#10268](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10268))

ENHANCEMENTS:

* dependencies: updating to build using Go 1.16 which adds support for `darwin/arm64` (Apple Silicon) ([#10615](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10615))
* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v51.2.0` ([#10561](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10561))
* Data Source: `azurerm_bastion_host` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* Data Source: `azurerm_point_to_site_vpn_gateway` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* Data Source: `azurerm_kubernetes_cluster` - exposing the `upgrade_settings` block ([#10376](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10376))
* Data Source: `azurerm_kubernetes_cluster_node_pool` - exposing the `upgrade_settings` block ([#10376](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10376))
* Data Source: `azurerm_route` - pdating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* Data Source: `azurerm_subnet ` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* Data Source: `azurerm_subscriptions` - adding the field `id` to the `subscriptions` block ([#10598](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10598))
* Data Source: `azurerm_virtual_network` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_bastion_host` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_bastion_host` - support for enhanced import validation ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_kubernetes_cluster` - support for configuring the `upgrade_settings` block ([#10376](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10376))
* `azurerm_kubernetes_cluster` - support for `automatic_channel_upgrade` ([#10530](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10530))
* `azurerm_kubernetes_cluster` - support for `skip_nodes_with_local_storage` within the `auto_scaler_profile` block ([#10531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10531))
* `azurerm_kubernetes_cluster` - support for `skip_nodes_with_system_pods` within the `auto_scaler_profile` block ([#10531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10531))
* `azurerm_kubernetes_cluster_node_pool` - support for configuring the `upgrade_settings` block ([#10376](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10376))
* `azurerm_lighthouse_definition` - add support for `principal_id_display_name` property ([#10613](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10613))
* `azurerm_log_analytics_workspace` - Support for `capacity_reservation_level` property and `CapacityReservation` SKU ([#10612](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10612))
* `azurerm_point_to_site_vpn_gateway` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_point_to_site_vpn_gateway` - support for enhanced import validation ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_route` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_route` - support for enhanced import validation ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_subnet` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_subnet` - support for enhanced import validation ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `synapse_workspace_resource` - support for the `azure_devops_repo` and `github_repo` blocks ([#10157](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10157))
* `azurerm_virtual_network` - updating to use a Resource ID Formatter ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))
* `azurerm_virtual_network` - support for enhanced import validation ([#10570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10570))

BUG FIXES:

* `azurerm_eventgrid_event_subscription` - change the number of possible `advanced_filter` items from `5` to `25` ([#10625](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10625))
* `azurerm_key_vault` - normalizing the casing on the `certificate_permissions`, `key_permissions`, `secret_permissions` and `storage_permissions` fields within the `access_policy` block ([#10593](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10593))
* `azurerm_key_vault_access_policy ` - normalizing the casing on the `certificate_permissions`, `key_permissions`, `secret_permissions` and `storage_permissions` fields ([#10593](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10593))
* `azurerm_mariadb_firewall_rule` - correctly validate the `name` property ([#10579](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10579))
* `azurerm_postgresql_server` - correctly change `ssl_minimal_tls_version_enforced` on update ([#10606](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10606))
* `azurerm_private_endpoint` - only updating the associated Private DNS Zone Group when there's changes ([#10559](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10559))
* `azurerm_resource_group_template_deployment` - fixing an issue where the API version for nested items couldn't be found during deletion ([#10565](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10565))

## 2.47.0 (February 11, 2021)

UPGRADE NOTES

* `azurerm_frontdoor` & `azurerm_frontdoor_custom_https_configuration` - the new fields `backend_pool_health_probes`, `backend_pool_load_balancing_settings`, `backend_pools`, `frontend_endpoints`, `routing_rules` have been added to the `azurerm_frontdoor` resource, which are a map of name-ID references. An upcoming version of the Azure Provider will change the blocks `backend_pool`, `backend_pool_health_probe`, `backend_pool_load_balancing`, `frontend_endpoint` and `routing_rule` from a List to a Set to work around an ordering issue within the Azure API - as such you should update your Terraform Configuration to reference these new Maps, rather than the Lists directly, due to the upcoming breaking change. For example, changing `azurerm_frontdoor.example.frontend_endpoint[1].id` to `azurerm_frontdoor.example.frontend_endpoints["exampleFrontendEndpoint2"]` ([#9357](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9357))
* `azurerm_lb_backend_address_pool` - the field `backend_addresses` has been deprecated and is no longer functional - instead the `azurerm_lb_backend_address_pool_address` resource offers the same functionality. ([#10488](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10488))
* `azurerm_linux_virtual_machine_scale_set` & `azurerm_windows_virtual_machine_scale_set` - the in-line `extension` block is now GA - the environment variable `ARM_PROVIDER_VMSS_EXTENSIONS_BETA` no longer has any effect and can be removed ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_data_factory_integration_runtime_managed` - this resource has been renamed/deprecated in favour of `azurerm_data_factory_integration_runtime_azure_ssis` ([#10236](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10236))
* The provider-block field `skip_credentials_validation` is now deprecated since this was non-functional and will be removed in 3.0 of the Azure Provider ([#10464](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10464))

FEATURES:

* **New Data Source:** `azurerm_key_vault_certificate_data` ([#8184](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8184))
* **New Resource:** `azurerm_application_insights_smart_detection_rule` ([#10539](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10539))
* **New Resource:** `azurerm_data_factory_integration_runtime_azure` ([#10236](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10236))
* **New Resource:** `azurerm_data_factory_integration_runtime_azure_ssis` ([#10236](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10236))
* **New Resource:** `azurerm_lb_backend_address_pool_address` ([#10488](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10488))

ENHANCEMENTS:

* dependencies: updating `github.com/hashicorp/terraform-plugin-sdk` to `v1.16.0` ([#10521](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10521))
* `azurerm_frontdoor` - added the new fields `backend_pool_health_probes`, `backend_pool_load_balancing_settings`, `backend_pools`, `frontend_endpoints`, `routing_rules` which are a map of name-ID references ([#9357](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9357))
* `azurerm_kubernetes_cluster` - updating the validation for the `log_analytics_workspace_id` field within the `oms_agent` block within the `addon_profile` block ([#10520](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10520))
* `azurerm_kubernetes_cluster` - support for configuring `only_critical_addons_enabled` ([#10307](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10307))
* `azurerm_kubernetes_cluster` - support for configuring `private_dns_zone_id` ([#10201](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10201))
* `azurerm_linux_virtual_machine_scale_set` - the `extension` block is now GA and available without enabling the beta ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_media_streaming_endpoint` - exporting the field `host_name` ([#10527](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10527))
* `azurerm_mssql_virtual_machine` - support for `auto_backup` ([#10460](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10460))
* `azurerm_windows_virtual_machine_scale_set` - the `extension` block is now GA and available without enabling the beta ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_site_recovery_replicated_vm` - support for the `recovery_public_ip_address_id` property and changing `target_static_ip` or `target_static_ip` force a new resource to be created ([#10446](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10446))

BUG FIXES:

* provider: the provider-block field `skip_credentials_validation` is now deprecated since this was non-functional. This will be removed in 3.0 of the Azure Provider ([#10464](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10464))
* Data Source: `azurerm_shared_image_versions` - retrieving all versions of the image prior to filtering ([#10519](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10519))
* `azurerm_app_service` - the `ip_restriction.x.ip_address` propertynow accepts anything other than an empty string ([#10440](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10440))
* `azurerm_cosmosdb_account` - validate the `key_vault_key_id` property is versionless ([#10420](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10420))
* `azurerm_cosmosdb_account` - will no longer panic if the response is nil ([#10525](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10525))
* `azurerm_eventhub_namespace` - correctly downgrade to the `Basic` sku ([#10536](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10536))
* `azurerm_key_vault_key` - export the `versionless_id` attribute ([#10420](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10420))
* `azurerm_lb_backend_address_pool` - the `backend_addresses` block is now deprecated and non-functional - use the `azurerm_lb_backend_address_pool_address` resource instead ([#10488](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10488))
* `azurerm_linux_virtual_machine_scale_set` - fixing a bug when `protected_settings` within the `extension` block was an empty string ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_linux_virtual_machine_scale_set` - fixing a bug when `settings` within the `extension` block was an empty string ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_monitor_diagnostic_setting` - changing the `log_analytics_workspace_id` property no longer creates a new resource ([#10512](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10512))
* `azurerm_storage_data_lake_gen2_filesystem` - do not set/retrieve ACLs when HNS is not enabled  ([#10470](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10470))
* `azurerm_windows_virtual_machine_scale_set ` - fixing a bug when `protected_settings` within the `extension` block was an empty string ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))
* `azurerm_windows_virtual_machine_scale_set ` - fixing a bug when `settings` within the `extension` block was an empty string ([#10528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10528))

## 2.46.1 (February 05, 2021)

BUG FIXES:

* `azurerm_lb_backend_address_pool` - mark `backend_address` as computed ([#10481](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10481))

## 2.46.0 (February 04, 2021)

FEATURES:

* **New Resource:** `azurerm_api_management_identity_provider_aadb2c` ([#10240](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10240))
* **New Resource:** `azurerm_cosmosdb_cassandra_table` ([#10328](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10328))

ENHANCEMENTS:

* dependencies: updating `recoveryservices` to API version `2018-07-10` ([#10373](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10373))
* `azurerm_api_management_diagnostic` - support for the `always_log_errors`, `http_correlation_protocol`, `log_client_ip`, `sampling_percentage` and `verbosity` properties ([#10325](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10325))
* `azurerm_api_management_diagnostic` - support for the `frontend_request`, `frontend_response`, `backend_request` and `backend_response` blocks ([#10325](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10325))
* `azurerm_kubernetes_cluster` - support for configuring the field `enable_host_encryption` within the `default_node_pool` block ([#10398](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10398))
* `azurerm_kubernetes_cluster` - added length validation to the `admin_password` field within the `windows_profile` block ([#10452](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10452))
* `azurerm_kubernetes_cluster_node_pool` - support for `enable_host_encryption` ([#10398](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10398))
* `azurerm_lb_backend_address_pool` - support for the `backend_address` block ([#10291](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10291))
* `azurerm_redis_cache` - support for the `public_network_access_enabled` property ([#10410](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10410))
* `azurerm_role_assignment` - adding validation for that the `scope` is either a Management Group, Subscription, Resource Group or Resource ID ([#10438](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10438))
* `azurerm_service_fabric_cluster` - support for the `reverse_proxy_certificate_common_names` block ([#10367](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10367))
* `azurerm_monitor_metric_alert` - support for the `skip_metric_validation` property ([#10422](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10422))

BUG FIXES:

* Data Source: `azurerm_api_management` fix an exception with User Assigned Managed Identities ([#10429](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10429))
* `azurerm_api_management_api_diagnostic` - fix a bug where specifying `log_client_ip = false` would not disable the setting ([#10325](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10325))
* `azurerm_key_vault` - fixing a race condition when setting the cache ([#10447](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10447))
* `azurerm_key_vault_certificate` - fixing a race condition when setting the cache ([#10447](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10447))
* `azurerm_key_vault_key` - fixing a race condition when setting the cache ([#10447](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10447))
* `azurerm_key_vault_secret` - fixing a race condition when setting the cache ([#10447](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10447))
* `azurerm_mssql_virtual_machine` - fixing a crash where the KeyVault was nil in the API response ([#10469](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10469))
* `azurerm_storage_account_datasource` - prevent panics from passing in an empty `name` ([#10370](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10370))
* `azurerm_storage_data_lake_gen2_filesystem` - change the `ace` property to a TypeSet to ensure consistent ordering ([#10372](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10372))
* `azurerm_storage_data_lake_gen2_path` - change the `ace` property to a TypeSet to ensure consistent ordering ([#10372](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10372))

## 2.45.1 (January 28, 2021)

BUG FIXES:

* `azurerm_app_service_environment` - prevent a panic when the API returns a nil cluster settings ([#10365](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10365))

## 2.45.0 (January 28, 2021)

FEATURES:

* **New Data Source** `azurerm_search_service` ([#10181](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10181))
* **New Resource:** `azurerm_data_factory_linked_service_snowflake` ([#10239](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10239))
* **New Resource:** `azurerm_data_factory_linked_service_azure_table_storage` ([#10305](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10305))
* **New Resource:** `azurerm_iothub_enrichment` ([#9239](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9239))
* **New Resource:** `azurerm_iot_security_solution` ([#10034](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10034))
* **New Resource:** `azurerm_media_streaming_policy` ([#10133](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10133))
* **New Resource:** `azurerm_spring_cloud_active_deployment` ([#9959](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9959))
* **New Resource:** `azurerm_spring_cloud_java_deployment` ([#9959](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9959))

IMPROVEMENTS:

* dependencies: updating to `v0.11.17` of `github.com/Azure/go-autorest/autorest` ([#10259](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10259))
* dependencies: updating the `firewall` resources to use the Networking API `2020-07-01` ([#10252](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10252))
* dependencies: updating the `load balancer` resources to use the Networking API version `2020-05-01` ([#10263](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10263))
* Data Source: `azurerm_app_service_environment` - export the `cluster_setting` block ([#10303](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10303))
* Data Source: `azurerm_key_vault_certificate` - support for the `certificate_data_base64` attribute ([#10275](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10275))
* `azurerm_app_service` - support for the propety `number_of_workers` ([#10143](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10143))
* `azurerm_app_service_environment` - support for the `cluster_setting` block ([#10303](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10303))
* `azurerm_data_factory_dataset_delimited_text` - support for the `compression_codec` property ([#10182](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10182))
* `azurerm_firewall_policy` - support for the `sku` property ([#10186](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10186))
* `azurerm_iothub` - support for the `enrichment` property ([#9239](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9239))
* `azurerm_key_vault` - optimised loading of and added caching when retrieving the Key Vault ([#10330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10330))
* `azurerm_key_vault` - support both ipv4 and cidr formats for the `network_acls.ip_rules` property ([#10266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10266))
* `azurerm_key_vault_certificate` - optimised loading of and added caching when retrieving the Key Vault ([#10330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10330))
* `azurerm_key_vault_key` - optimised loading of and added caching when retrieving the Key Vault ([#10330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10330))
* `azurerm_key_vault_secret` - optimised loading of and added caching when retrieving the Key Vault ([#10330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10330))
* `azurerm_key_vault_certificate` - support for the `certificate_data_base64` attribute ([#10275](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10275))
* `azurerm_linux_virtual_machine` - skipping shutdown for a machine in a failed state ([#10189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10189))
* `azurerm_media_services_account` - support for setting the `storage_authentication_type` field to `System` ([#10133](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10133))
* `azurerm_redis_cache` - support multiple availability zones ([#10283](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10283))
* `azurerm_storage_data_lake_gen2_filesystem` - support for the `ace` block ([#9917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9917))
* `azurerm_servicebus_namespace` - will now allow a capacity of `16` for the `Premium` SKU ([#10337](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10337))
* `azurerm_windows_virtual_machine` - skipping shutdown for a machine in a failed state ([#10189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10189))
* `azurerm_linux_virtual_machine_scale_set` - support for the `extensions_time_budget` property ([#10298](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10298))
* `azurerm_windows_virtual_machine_scale_set` - support for the `extensions_time_budget` property ([#10298](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10298))

BUG FIXES:

* `azurerm_iot_time_series_insights_reference_data_set` - the field `data_string_comparison_behavior` is now `ForceNew` ([#10343](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10343))
* `azurerm_iot_time_series_insights_reference_data_set` - the `key_property` block is now `ForceNew` ([#10343](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10343))
* `azurerm_linux_virtual_machine_scale_set` - fixing an issue where `protected_settings` field within the `extension` block couldn't be empty ([#10351](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10351))
* `azurerm_linux_virtual_machine_scale_set` - fixing an issue where `settings` field within the `extension` block couldn't be empty ([#10351](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10351))
* `azurerm_media_streaming_endpoint` - stopping the streaming endpoint prior to deletion if the endpoint is in a running state ([#10216](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10216))
* `azurerm_role_definition` - don't add `scope` to `assignable_scopes` unless none are specified ([#8624](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8624))
* `azurerm_windows_virtual_machine_scale_set` - fixing an issue where `protected_settings` field within the `extension` block couldn't be empty ([#10351](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10351))
* `azurerm_windows_virtual_machine_scale_set` - fixing an issue where `settings` field within the `extension` block couldn't be empty ([#10351](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10351))

## 2.44.0 (January 21, 2021)

FEATURES:

* **New Data Source:** `azurerm_iothub` ([#10228](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10228))
* **New Resource:** `azurerm_media_content_key_policy` ([#9971](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9971))

IMPROVEMENTS:

* dependencies: updating `github.com/Azure/go-autorest` to `v0.11.16` ([#10164](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10164))
* dependencies: updating `appconfiguration` to API version `2020-06-01` ([#10176](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10176))
* dependencies: updating `appplatform` to API version `2020-07-01` ([#10175](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10175))
* dependencies: updating `containerservice` to API version `2020-12-01` ([#10171](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10171))
* dependencies: updating `msi` to API version `2018-11-30` ([#10174](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10174))
* Data Source: `azurerm_kubernetes_cluster` - support for the field `user_assigned_identity_id` within the `identity` block ([#8737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8737))
* `azurerm_api_management` - support additional TLS ciphers within the `security` block ([#9276](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9276))
* `azurerm_api_management_api_diagnostic` - support the `sampling_percentage` property ([#9321](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9321))
* `azurerm_container_group` - support for updating `tags` ([#10210](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10210))
* `azurerm_kubernetes_cluster` - the field `type` within the `identity` block can now be set to `UserAssigned` ([#8737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8737))
* `azurerm_kubernetes_cluster` - support for the field `new_pod_scale_up_delay` within the `auto_scaler_profile` block ([#9291](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9291))
* `azurerm_kubernetes_cluster` - support for the field `user_assigned_identity_id` within the `identity` block ([#8737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8737))
* `azurerm_monitor_autoscale_setting` - now supports the `dimensions` property ([#9795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9795))
* `azurerm_sentinel_alert_rule_scheduled` - now supports the `event_grouping_setting` property ([#10078](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10078))

BUG FIXES:

* `azurerm_backup_protected_file_share` - updating to account for a breaking API change ([#9015](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9015))
* `azurerm_key_vault_certificate` - fixing a crash when `subject` within the `certificate_policy` block was nil ([#10200](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10200))
* `azurerm_user_assigned_identity` - adding a state migration to update the ID format ([#10196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10196))

## 2.43.0 (January 14, 2021)

FEATURES:

* **New Data Source:** `azurerm_sentinel_alert_rule_template` ([#7020](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7020))

IMPROVEMENTS:

* Data Source: `azurerm_api_management` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* Data Source: `azurerm_kubernetes_cluster` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* Data Source: `azurerm_virtual_machine` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* Data Source: `azurerm_virtual_machine_scale_set` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_api_management` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_app_service` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_app_service_slot` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_container_group` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_cosmosdb_account` - support for `analytical_storage_enabled property` ([#10055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10055))
* `azurerm_cosmosdb_gremlin_graph` - support the `default_ttl` property ([#10159](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10159))
* `azurerm_data_factory` - support for `public_network_enabled` ([#9605](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9605))
* `azurerm_data_factory_dataset_delimited_text` - support for the `compression_type` property ([#10070](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10070))
* `azurerm_data_factory_linked_service_sql_server`: support for the `key_vault_password` block ([#10032](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10032))
* `azurerm_eventgrid_domain` - support for the `public_network_access_enabled` and `inbound_ip_rule` properties  ([#9922](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9922))
* `azurerm_eventgrid_topic` - support for the `public_network_access_enabled` and `inbound_ip_rule` properties  ([#9922](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9922))
* `azurerm_eventhub_namespace` - support the `trusted_service_access_enabled` property ([#10169](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10169))
* `azurerm_function_app` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_function_app_slot` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_kusto_cluster` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_linux_virtual_machine` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_linux_virtual_machine_scale_set` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_security_center_automation` - the field `event_source` within the `source` block now supports `SecureScoreControls ` and `SecureScores` ([#10126](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10126))
* `azurerm_synapse_workspace` - support for the `sql_identity_control_enabled` property ([#10033](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10033))
* `azurerm_virtual_machine` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_virtual_machine_scale_set` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_windows_virtual_machine` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_windows_virtual_machine_scale_set` - adding validation on the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))

BUG FIXES:

* Data Source: `azurerm_log_analytics_workspace` - returning the Resource ID in the correct casing ([#10162](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10162))
* `azurerm_advanced_threat_protection` - fix a regression in the Resouce ID format ([#10190](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10190))
* `azurerm_api_management` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_app_service` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_app_service_slot` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_application_gateway` - ensuring the casing on `identity_ids` within the `identity` block ([#10031](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10031))
* `azurerm_blueprint_assignment` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_container_group` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_databricks_workspace` - changing the sku no longer always forces a new resource to be created ([#9541](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9541))
* `azurerm_function_app` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_function_app_slot` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_kubernetes_cluster` - ensuring the casing of the `user_assigned_identity_id` field within the `kubelet_identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_kusto_cluster` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_linux_virtual_machine` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_linux_virtual_machine_scale_set` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_monitor_diagnostic_setting` - handling mixed casing of the EventHub Namespace Authorization Rule ID ([#10104](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10104))
* `azurerm_mssql_virtual_machine` - address persistent diff and use relative expiry for service principal password ([#10125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10125))
* `azurerm_role_assignment` - fix race condition in read after create ([#10134](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10134))
* `azurerm_role_definition` - address eventual consistency issues in update and delete ([#10170](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10170))
* `azurerm_virtual_machine` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_virtual_machine_scale_set` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_windows_virtual_machine` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))
* `azurerm_windows_virtual_machine_scale_set` - ensuring the casing of the `identity_ids` field within the `identity` block ([#10105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10105))

## 2.42.0 (January 08, 2021)

BREAKING CHANGES

* `azurerm_key_vault` - the field `soft_delete_enabled` is now defaulted to `true` to match the breaking change in the Azure API where Key Vaults now have Soft Delete enabled by default, which cannot be disabled. This property is now non-functional, defaults to `true` and will be removed in version 3.0 of the Azure Provider. ([#10088](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10088))
* `azurerm_key_vault` - the field `soft_delete_retention_days` is now defaulted to `90` days to match the Azure API behaviour, as the Azure API does not return a value for this field when not explicitly configured, so defaulting this removes a diff with `0`. ([#10088](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10088))

FEATURES:

* **New Data Source:** `azurerm_eventgrid_domain_topic` ([#10050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10050))
* **New Data Source:** `azurerm_ssh_public_key` ([#9842](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9842))
* **New Resource:** `azurerm_data_factory_linked_service_synapse` ([#9928](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9928))
* **New Resource:** `azurerm_disk_access` ([#9889](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9889))
* **New Resource:** `azurerm_media_streaming_locator` ([#9992](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9992))
* **New Resource:** `azurerm_sentinel_alert_rule_fusion` ([#9829](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9829))
* **New Resource:** `azurerm_ssh_public_key` ([#9842](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9842))

IMPROVEMENTS:

* batch: updating to API version `2020-03-01` ([#10036](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10036))
* dependencies: upgrading to `v49.2.0` of `github.com/Azure/azure-sdk-for-go` ([#10042](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10042))
* dependencies: upgrading to `v0.15.1` of `github.com/tombuildsstuff/giovanni` ([#10035](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10035))
* Data Source: `azurerm_hdinsight_cluster` - support for the `kafka_rest_proxy_endpoint` property ([#8064](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8064))
* Data Source: `azurerm_databricks_workspace` - support for the `tags` property ([#9933](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9933))
* Data Source: `azurerm_subscription` - support for the `tags` property ([#8064](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8064))
* `azurerm_app_service` - now supports  `detailed_error_mesage_enabled` and `failed_request_tracing_enabled ` logs settings ([#9162](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9162))
* `azurerm_app_service` - now supports  `service_tag` in `ip_restriction` blocks ([#9609](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9609))
* `azurerm_app_service_slot` - now supports  `detailed_error_mesage_enabled` and `failed_request_tracing_enabled ` logs settings ([#9162](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9162))
* `azurerm_batch_pool` support for the `public_address_provisioning_type` property ([#10036](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10036))
* `azurerm_api_management` - support `Consumption_0` for the `sku_name` property ([#6868](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6868))
* `azurerm_cdn_endpoint` - only send `content_types_to_compress` and `geo_filter` to the API when actually set ([#9902](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9902))
* `azurerm_cosmosdb_mongo_collection` - correctly read back the `_id` index when mongo 3.6 ([#8690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8690))
* `azurerm_container_group` - support for the `volume.empty_dir` property ([#9836](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9836))
* `azurerm_data_factory_linked_service_azure_file_storage` - support for the `file_share` property ([#9934](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9934))
* `azurerm_dedicated_host` - support for addtional `sku_name` values ([#9951](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9951))
* `azurerm_devspace_controller` - deprecating since new DevSpace Controllers can no longer be provisioned, this will be removed in version 3.0 of the Azure Provider ([#10049](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10049))
* `azurerm_function_app` - make `pre_warmed_instance_count` computed to use azure's default ([#9069](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9069))
* `azurerm_function_app` - now supports  `service_tag` in `ip_restriction` blocks ([#9609](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9609))
* `azurerm_hdinsight_hadoop_cluster` - allow the value `Standard_D4a_V4` for the `vm_type` property ([#10000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10000))
* `azurerm_hdinsight_kafka_cluster` - support for the `rest_proxy` and `kafka_management_node` blocks ([#8064](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8064))
* `azurerm_key_vault` - the field `soft_delete_enabled` is now defaulted to `true` to match the Azure API behaviour where Soft Delete is force-enabled and can no longer be disabled. This field is deprecated, can be safely removed from your Terraform Configuration, and will be removed in version 3.0 of the Azure Provider. ([#10088](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10088))
* `azurerm_kubernetes_cluster` - add support for network_mode ([#8828](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8828))
* `azurerm_log_analytics_linked_service` - add validation for resource ID type ([#9932](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9932))
* `azurerm_log_analytics_linked_service` - update validation to use generated validate functions ([#9950](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9950))
* `azurerm_monitor_diagnostic_setting` - validation that `eventhub_authorization_rule_id` is a EventHub Namespace Authorization Rule ID ([#9914](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9914))
* `azurerm_monitor_diagnostic_setting` - validation that `log_analytics_workspace_id` is a Log Analytics Workspace ID ([#9914](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9914))
* `azurerm_monitor_diagnostic_setting` - validation that `storage_account_id` is a Storage Account ID ([#9914](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9914))
* `azurerm_network_security_rule` - increase allowed the number of `application_security_group` blocks allowed ([#9884](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9884))
* `azurerm_sentinel_alert_rule_ms_security_incident` - support the `alert_rule_template_guid` and `display_name_exclude_filter` properties ([#9797](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9797))
* `azurerm_sentinel_alert_rule_scheduled` - support for the `alert_rule_template_guid` property ([#9712](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9712))
* `azurerm_sentinel_alert_rule_scheduled` - support for creating incidents ([#8564](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8564))
* `azurerm_spring_cloud_app` - support the properties `https_only`, `is_public`, and `persistent_disk` ([#9957](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9957))
* `azurerm_subscription` - support for the `tags` property ([#9047](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9047))
* `azurerm_synapse_workspace` - support for the `managed_resource_group_name` property ([#10017](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10017))
* `azurerm_traffic_manager_profile` - support for the `traffic_view_enabled` property ([#10005](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10005))

BUG FIXES:

provider: will not correctly register the `Microsoft.Blueprint` and `Microsoft.HealthcareApis` RPs ([#10062](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10062))
* `azurerm_application_gateway` - allow `750` for `file_upload_limit_mb` when the sku is `WAF_v2` ([#8753](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8753))
* `azurerm_firewall_policy_rule_collection_group` - correctly validate the `network_rule_collection.destination_ports` property ([#9490](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9490))
* `azurerm_cdn_endpoint` - changing many `delivery_rule` condition `match_values` to optional ([#8850](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8850))
* `azurerm_cosmosdb_account` - always include `key_vault_id` in update requests for azure policy enginer compatibility ([#9966](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9966))
* `azurerm_cosmosdb_table` - do not call the throughput api when serverless ([#9749](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9749))
* `azurerm_key_vault` - the field `soft_delete_retention_days` is now defaulted to `90` days to match the Azure API behaviour. ([#10088](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10088))
* `azurerm_kubernetes_cluster` - parse oms `log_analytics_workspace_id` to ensure correct casing ([#9976](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9976))
* `azurerm_role_assignment` fix crash in retry logic ([#10051](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10051))
* `azurerm_storage_account` - allow hns when `account_tier` is `Premium` ([#9548](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9548))
* `azurerm_storage_share_file` - allowing files smaller than 4KB to be uploaded ([#10035](https://github.com/terraform-providers/terraform-provider-azurerm/issues/10035))

## 2.41.0 (December 17, 2020)

UPGRADE NOTES:

* `azurerm_key_vault` - Azure will be introducing a breaking change on December 31st, 2020 by force-enabling Soft Delete on all new and existing Key Vaults. To workaround this, this release of the Azure Provider still allows you to configure Soft Delete on before this date (but once this is enabled this cannot be disabled). Since new Key Vaults will automatically be provisioned using Soft Delete in the future, and existing Key Vaults will be upgraded - a future release will deprecate the `soft_delete_enabled` field and default this to true early in 2021. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_certificate` - Terraform will now attempt to `purge` Certificates during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_key` - Terraform will now attempt to `purge` Keys during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_secret` - Terraform will now attempt to `purge` Secrets during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))

FEATURES:

* **New Resource:** `azurerm_eventgrid_system_topic_event_subscription` ([#9852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9852))
* **New Resource:** `azurerm_media_job` ([#9859](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9859))
* **New Resource:** `azurerm_media_streaming_endpoint` ([#9537](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9537))
* **New Resource:** `azurerm_subnet_service_endpoint_storage_policy` ([#8966](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8966))
* **New Resource:** `azurerm_synapse_managed_private_endpoint` ([#9260](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9260))

IMPROVEMENTS:

* `azurerm_app_service` - Add support for `outbound_ip_address_list` and `possible_outbound_ip_address_list` ([#9871](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9871))
* `azurerm_disk_encryption_set` - support for updating `key_vault_key_id` ([#7913](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7913))
* `azurerm_iot_time_series_insights_gen2_environment` - exposing `data_access_fqdn` ([#9848](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9848))
* `azurerm_key_vault_certificate` - performing a "purge" of the Certificate during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_key` - performing a "purge" of the Key during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_secret` - performing a "purge" of the Secret during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_log_analytics_linked_service` - Add new fields `workspace_id`, `read_access_id`, and `write_access_id` ([#9410](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9410))
* `azurerm_linux_virtual_machine` - Normalise SSH keys to cover VM import cases ([#9897](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9897))
* `azurerm_subnet` - support for the `service_endpoint_policy` block ([#8966](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8966))
* `azurerm_traffic_manager_profile` - support for new field `max_return` and support for `traffic_routing_method` to be `MultiValue` ([#9487](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9487))

BUG FIXES:

* `azurerm_key_vault_certificate` - reading `dns_names` and `emails` within the `subject_alternative_names` block from the Certificate if not returned from the API ([#8631](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8631))
* `azurerm_key_vault_certificate` - polling until the Certificate is fully deleted during deletion ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_key` - polling until the Key is fully deleted during deletion ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_secret` -  polling until the Secret is fully deleted during deletion ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_log_analytics_workspace` - adding a state migration to correctly update the Resource ID ([#9853](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9853))

## 2.40.0 (December 10, 2020)

FEATURES:

* **New Resource:** `azurerm_app_service_certificate_binding` ([#9415](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9415))
* **New Resource:** `azurerm_digital_twins_endpoint_eventhub` ([#9673](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9673))
* **New Resource:** `azurerm_digital_twins_endpoint_servicebus`  ([#9702](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9702))
* **New Resource:** `azurerm_media_asset` ([#9387](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9387))
* **New Resource:** `azurerm_media_transform` ([#9663](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9663))
* **New Resource:** `azurerm_resource_provider` ([#7951](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7951))
* **New Resource:** `azurerm_stack_hci_cluster` ([#9134](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9134))
* **New Resource:** `azurerm_storage_share_file` ([#9406](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9406))
* **New Resource:** `azurerm_storage_sync_cloud_endpoint` ([#8540](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8540))

IMPROVEMENTS:

* dependencies: upgrading `github.com/Azure/go-autorest/validation` to `v0.3.1` ([#9783](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9783))
* dependencies: updating Log Analytics to API version `2020-08-01` ([#9764](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9764))
* internal: disabling the Azure SDK's validation since it's superfluous ([#9783](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9783))
* `azurerm_app_service` - support for PHP version `7.4` ([#9727](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9727))
* `azurerm_bot_channel_directline` - support for enhanced import validation ([#9690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9690))
* `azurerm_bot_channel_email` - support for enhanced import validation ([#9690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9690))
* `azurerm_bot_channel_ms_teams` - support for enhanced import validation ([#9690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9690))
* `azurerm_bot_channel_slack` - support for enhanced import validation ([#9690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9690))
* `azurerm_bot_channels_registration` - support for enhanced import validation ([#9690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9690))
* `azurerm_bot_connection` - support for enhanced import validation ([#9690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9690))
* `azurerm_bot_web_app` - support for enhanced import validation ([#9690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9690))
* `azurerm_cosmosdb_sql_container` - support for the `partition_key_version` property ([#9496](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9496))
* `azurerm_kusto_cluster` - support for the `engine` property ([#9696](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9696))
* `azurerm_kusto_eventhub_data_connection` - support for `compression` ([#9692](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9692))
* `azurerm_iothub` - support for the `min_tls_version` property ([#9670](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9670))
* `azurerm_recovery_services_vault` - support for the `identity` block ([#9689](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9689))
* `azurerm_redis_cache` - adding enhanced import validation ([#9771](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9771))
* `azurerm_redis_cache` - adding validation that `subnet_id` is a valid Subnet ID ([#9771](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9771))
* `azurerm_redis_firewall_rule` - adding enhanced import validation ([#9771](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9771))
* `azurerm_redis_linked_server` - adding enhanced import validation ([#9771](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9771))
* `azurerm_redis_linked_server` - adding validation that `linked_redis_cache_id` is a valid Redis Cache ID ([#9771](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9771))
* `azurerm_security_center_automation` - support for the `description` and `tags` properties ([#9676](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9676))
* `azurerm_stream_analytics_reference_input_blob` - support for enhanced import validation ([#9735](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9735))
* `azurerm_stream_analytics_stream_input_blob` - support for enhanced import validation ([#9735](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9735))
* `azurerm_stream_analytics_stream_input_iothub` - support for enhanced import validation ([#9735](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9735))
* `azurerm_stream_analytics_stream_input_eventhub` - support for enhanced import validation ([#9735](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9735))
* `azurerm_storage_account` - enable the `allow_blob_public_access` and `azurerm_storage_account` properties in US Government Cloud ([#9540](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9540))

BUG FIXES:

* `azurerm_app_service_managed_certificate` - create certificate in service plan resource group to prevent diff loop ([#9701](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9701))
* `azurerm_bot_channel_directline` - the field `bot_name` is now ForceNew to match the documentation/API behaviour ([#9690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9690))
* `azurerm_bot_channel_ms_teams` - the field `bot_name` is now ForceNew to match the documentation/API behaviour ([#9690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9690))
* `azurerm_bot_channel_slack` - the field `bot_name` is now ForceNew to match the documentation/API behaviour ([#9690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9690))
* `azurerm_bot_connection` - the field `bot_name` is now ForceNew to match the documentation/API behaviour ([#9690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9690))
* `azurerm_frontdoor` - working around an upstream API issue by rewriting the returned ID's within Terraform ([#9750](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9750))
* `azurerm_frontdoor_custom_https_configuration` - working around an upstream API issue by rewriting the returned ID's within Terraform ([#9750](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9750))
* `azurerm_frontdoor_firewall_policy` - working around an upstream API issue by rewriting the returned ID's within Terraform ([#9750](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9750))
* `azurerm_media_services_account` - fixing a bug where `storage_authentication_type` wasn't set ([#9663](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9663))
* `azurerm_media_service_account` - checking for the presence of an existing account during creation ([#9802](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9802))
* `azurerm_postgresql_server` - changing the `geo_redundant_backup_enabled` property now forces a new resource ([#9694](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9694))
* `azurerm_postgresql_server` - Fix issue when specifying empty threat detection list attributes ([#9739](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9739))
* `azurerm_signar_service` - having an empty `allowed_origins` in the `cors` block will no longer cause a panic ([#9671](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9671))

## 2.39.0 (December 04, 2020)

FEATURES:

* **New Resource:** `azurerm_api_management_policy` ([#9215](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9215))
* **New Resource:** `azurerm_digital_twins_endpoint_eventgrid` ([#9489](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9489))
* **New Resource:** `azurerm_iot_time_series_insights_gen2_environment` ([#9616](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9616))

IMPROVEMENTS: 

* `azurerm_dashboard` - adding validation at import time to ensure the ID is for a Dashboard ([#9530](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9530))
* `azurerm_keyvault_certificate` - add `3072` to allowed values for `key_size` ([#9524](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9524))
* `azurerm_media_services_account` - support for the `identity`, `tags`, and `storage_authentication` properties ([#9457](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9457))
* `azurerm_notification_hub_authorization_rule` - adding validation at import time to ensure the ID is for a Notification Hub Authorization Rule ([#9529](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9529))
* `azurerm_notification_hub_namespace` - adding validation at import time to ensure the ID is for a Notification Hub Namespace ([#9529](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9529))
* `azurerm_postgresql_active_directory_administrator` - validating during import that the ID is for a PostgreSQL Active Directory Administrator ([#9532](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9532))
* `azurerm_postgresql_configuration` - validating during import that the ID is for a PostgreSQL Configuration ([#9532](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9532))
* `azurerm_postgresql_database` - validating during import that the ID is for a PostgreSQL Database ([#9532](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9532))
* `azurerm_postgresql_firewall_rule` - validating during import that the ID is for a PostgreSQL Firewall Rule ([#9532](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9532))
* `azurerm_postgresql_virtual_network_rule` - validating during import that the ID is for a PostgreSQL Virtual Network Rule ([#9532](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9532))
* `azurerm_traffic_manager_profile` - allow up to `2147483647` for the `ttl` property ([#9522](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9522))

BUG FIXES:

* `azurerm_security_center_workspace` - fixing the casing on the `workspace_id` ([#9651](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9651))
* `azurerm_eventhub_dedicated_cluster` - the `sku_name` capacity can be greater then `1` ([#9649](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9649))

## 2.38.0 (November 27, 2020)

FEATURES:

* **New Resource** `azurerm_app_service_managed_certificate` ([#9378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9378))
* **New Data Source:** `azurerm_digital_twins_instance` ([#9430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9430))
* **New Data Source:** `azurerm_virtual_wan` ([#9382](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9382))
* **New Resource:** `azurerm_digital_twins_instance` ([#9430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9430))

IMPROVEMENTS: 

* dependencies: updating App Service to API version `2020-06-01` ([#9409](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9409))
* Data Source `azurerm_app_service` now exports the `custom_domain_verification_id` attribute ([#9378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9378))
* Data Source`azurerm_function_app` now exports the `custom_domain_verification_id` attribute ([#9378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9378))
* Data Source: `azurerm_spring_cloud_service` - now exports the `outbound_public_ip_addresses` attribute ([#9261](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9261))
* `azurerm_app_service` now exports `custom_domain_verification_id` ([#9378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9378))
* `azurerm_application_insights` - validating the resource ID is correct during import ([#9446](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9446))
* `azurerm_application_insights_web_test` - validating the resource ID is correct during import ([#9446](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9446))
* `azurerm_express_route_circuit_peering` - support for the `ipv6` block  ([#9235](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9235))
* `azurerm_function_app` now exports the `custom_domain_verification_id` attribute ([#9378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9378))
* `azurerm_vpn_server_configuration` - deprecate the `radius_server` block in favour of the `radius` block which supports multiple servers ([#9308](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9308))
* `azurerm_spring_cloud_service` - now exports the `outbound_public_ip_addresses` attribute ([#9261](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9261))
* `azurerm_virtual_network_gateway` - support for the `dpd_timeout_seconds` and `local_azure_ip_address_enabled` properties ([#9330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9330))
* `azurerm_virtual_network_gateway_connection` - support for the `private_ip_address_enabled` propeties and the `custom_route` block ([#9330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9330))

BUG FIXES:

* `azurerm_api_management` - fixing an issue where developer portal certificates are updated on every apply ([#7299](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7299))
* `azurerm_cosmosdb_account` - corrently updates the `zone_redundant` property during updates ([#9485](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9485))
* `azurerm_search_service` - `allowed_ips` now supports specifying a CIDR Block in addition to an IPv4 address ([#9493](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9493))
* `azurerm_virtual_desktop_application_group` - adding a state migration to avoid a breaking change when upgrading from `v2.35.0` or later ([#9495](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9495))
* `azurerm_virtual_desktop_host_pool` - adding a state migration to avoid a breaking change when upgrading from `v2.35.0` or later ([#9495](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9495))
* `azurerm_virtual_desktop_workspace` - adding a state migration to avoid a breaking change when upgrading from `v2.35.0` or later ([#9495](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9495))
* `azurerm_virtual_desktop_workspace_application_group_association` - adding a state migration to avoid a breaking change when upgrading from `v2.35.0` or later ([#9495](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9495))
* `azurerm_windows_virtual_machine` - no longer sets `patch_mode` on creation if it is the default value ([#9495](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9432))

## 2.37.0 (November 20, 2020)

FEATURES:

* **New Data Source:** `azurerm_servicebus_subscription` ([#9272](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9272))
* **New Data Source:** `azurerm_storage_encryption_scope` ([#8894](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8894))
* **New Resource:** `azurerm_log_analytics_cluster` ([#8946](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8946))
* **New Resource:** `azurerm_log_analytics_cluster_customer_managed_key` ([#8946](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8946))
* **New Resource:** `azurerm_security_center_automation` ([#8781](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8781))
* **New Resource:** `azurerm_storage_data_lake_gen2_path` ([#7521](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7521))
* **New Resource:** `azurerm_storage_encryption_scope` ([#8894](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8894))
* **New Resource:** `azurerm_vpn_gateway_connection` ([#9160](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9160))

IMPROVEMENTS:

* storage: foundational improvements to support toggling between the Data Plane and Resource Manager Storage API's in the future ([#9314](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9314))
* Data Source: `azurerm_firewall`-  exposing `dns_servers`, `firewall_policy_id`, `sku_name`, `sku_tier`, `threat_intel_mode`, `virtual_hub` and `zones` ([#8879](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8879))
* Data Source: `azurerm_firewall`-  exposing `public_ip_address_id` and `private_ip_address_id` within the `ip_configuration` block ([#8879](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8879))
* Data Source: `azurerm_firewall`-  exposing `name` within the `management_ip_configuration` block ([#8879](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8879))
* Data Source: `azurerm_kubernetes_node_pool` - exposing `os_disk_type` ([#9166](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9166))
* `azurerm_api_management_api_diagnostic` - support for the `always_log_errors`, `http_correlation_protocol`, `log_client_ip` and `verbosity` attributes ([#9172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9172))
* `azurerm_api_management_api_diagnostic` - support the `frontend_request`, `frontend_response`, `backend_request` and `backend_response` blocks ([#9172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9172))
* `azurerm_container_group` - support for secret container volumes with the `container.#.volume.#.secret` attribute ([#9117](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9117))
* `azurerm_cosmosdb_account` - support for the `public_network_access_enabled` property ([#9236](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9236))
* `azurerm_cosmosdb_cassandra_keyspace` - `throughput` can now be set to higher than `1000000` if enabled by Azure Support ([#9050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9050))
* `azurerm_cosmosdb_gremlin_database` - `throughput` can now be set to higher than `1000000` if enabled by Azure Support ([#9050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9050))
* `azurerm_cosmosdb_mongo_database` - `throughput` can now be set to higher than `1000000` if enabled by Azure Support ([#9050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9050))
* `azurerm_cosmosdb_sql_container` - `max_throughput` within the `autoscale_settings` block can now be set to higher than `1000000` if enabled by Azure Support ([#9050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9050))
* `azurerm_cosmosdb_sql_database` - `throughput` can now be set to higher than `1000000` if enabled by Azure Support ([#9050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9050))
* `azurerm_cosmosdb_table` - `throughput` can now be set to higher than `1000000` if enabled by Azure Support ([#9050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9050))
* `azurerm_dns_zone` - support for the `soa_record` block ([#9319](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9319))
* `azurerm_firewall` - support for `firewall_policy_id`, `sku_name`, `sku_tier` and `virtual_hub` ([#8879](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8879))
* `azurerm_kubernetes_cluster` - support for configuring `os_disk_type` within the `default_node_pool` block ([#9166](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9166))
* `azurerm_kubernetes_cluster` - `max_count` within the `default_node_pool` block can now be set to a maximum value of `1000` ([#9227](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9227))
* `azurerm_kubernetes_cluster` - `min_count` within the `default_node_pool` block can now be set to a maximum value of `1000` ([#9227](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9227))
* `azurerm_kubernetes_cluster` - `node_count` within the `default_node_pool` block can now be set to a maximum value of `1000` ([#9227](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9227))
* `azurerm_kubernetes_cluster` - the block `http_application_routing` within the `addon_profile` block can now be updated/removed ([#9358](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9358))
* `azurerm_kubernetes_node_pool` - support for configuring `os_disk_type` ([#9166](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9166))
* `azurerm_kubernetes_node_pool` - `max_count` can now be set to a maximum value of `1000` ([#9227](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9227))
* `azurerm_kubernetes_node_pool` - `min_count` can now be set to a maximum value of `1000` ([#9227](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9227))
* `azurerm_kubernetes_node_pool` - `node_count` can now be set to a maximum value of `1000` ([#9227](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9227))
* `azurerm_linux_virtual_machine` - support for the `extensions_time_budget` property ([#9257](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9257))
* `azurerm_linux_virtual_machine` - updating the `dedicated_host_id` no longer forces a new resource ([#9264](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9264))
* `azurerm_linux_virtual_machine` - support for graceful shutdowns (via the features block) ([#8470](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8470))
* `azurerm_linux_virtual_machine_scale_set` - support for the `platform_fault_domain_count`, `disk_iops_read_write`, and `disk_mbps_read_write` properties ([#9262](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9262))
* `azurerm_mssql_database` - `sku_name` supports more `DWxxxc` options ([#9370](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9370))
* `azurerm_policy_set_definition` - support for the `policy_definition_group` block ([#9259](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9259))
* `azurerm_postgresql_server` - increase max storage to 16TiB ([#9373](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9373))
* `azurerm_private_dns_zone` - support for the `soa_record` block ([#9319](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9319))
* `azurerm_storage_blob` - support for `content_md5` ([#7786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7786))
* `azurerm_windows_virtual_machine` - support for the `extensions_time_budget` property ([#9257](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9257))
* `azurerm_windows_virtual_machine` - updating the `dedicated_host_id` nolonger forces a new resource ([#9264](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9264))
* `azurerm_windows_virtual_machine` - support for graceful shutdowns (via the features block) ([#8470](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8470))
* `azurerm_windows_virtual_machine` - support for the `patch_mode` property ([#9258](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9258))
* `azurerm_windows_virtual_machine_scale_set` - support for the `platform_fault_domain_count`, `disk_iops_read_write`, and `disk_mbps_read_write` properties ([#9262](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9262))

BUG FIXES:

* Data Source: `azurerm_key_vault_certificate` - fixing a crash when serializing the certificate policy block ([#9355](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9355))
* `azurerm_api_management` - the field `xml_content` within the `policy` block now supports C#/.net interpolations ([#9296](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9296))
* `azurerm_cosmosdb_sql_container` - no longer attempts to get throughput settings when cosmos account is serverless ([#9311](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9311))
* `azurerm_firewall_policy` - deprecate the `dns.network_rule_fqdn_enabled` property as the API no longer allows it to be set ([#9332](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9332))
* `azurerm_key_vault_certificate` - fixing a crash when serializing the certificate policy block ([#9355](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9355))
* `azurerm_mssql_virtual_machine` - fixing a crash when serializing `auto_patching` ([#9388](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9388))
* `azurerm_resource_group_template_deployment` - fixing an issue during deletion where the API version of nested resources couldn't be determined ([#9364](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9364))

## 2.36.0 (November 12, 2020)

UPGRADE NOTES:

* `azurerm_network_connection_monitor` - has been updated to work with v2 of the resource as the service team is deprecating v1 - all v1 properties have been deprecated and will be removed in version `3.0` of the provider and v2 propeties added. ([#8640](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8640))

FEATURES:

* **New Data Source:** `azurerm_data_share_dataset_kusto_database` ([#8544](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8544))
* **New Data Source:** `azurerm_traffic_manager_profile` ([#9229](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9229))
* **New Resource:** `azurerm_api_management_custom_domain` ([#8228](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8228))
* **New Resource:** `azurerm_data_share_dataset_kusto_database` ([#8544](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8544))
* **New Resource:** `azurerm_log_analytics_storage_insights` ([#9014](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9014))
* **New Resource:** `azurerm_monitor_smart_detector_alert_rule` ([#9032](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9032))
* **New Resource:** `azurerm_virtual_hub_security_partner_provider` ([#8978](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8978))
* **New Resource:** `azurerm_virtual_hub_bgp_connection` ([#8959](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8959))

IMPROVEMENTS:

* dependencies: upgrading to `v0.4.2` of `github.com/Azure/go-autorest/autorest/azure/cli` ([#9168](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9168))
* dependencies: upgrading to `v48.1.0` of `github.com/Azure/azure-sdk-for-go` ([#9213](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9213))
* dependencies: upgrading to `v0.13.0` of `github.com/hashicorp/go-azure-helpers` ([#9191](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9191))
* dependencies: upgrading to `v0.14.0` of `github.com/tombuildsstuff/giovanni` ([#9189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9189))
* storage: upgrading the Data Plane API's to API Version `2019-12-12` ([#9192](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9192))
* Data Source `azurerm_kubernetes_node_pool` - exporting `proximity_placement_group_id` ([#9195](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9195))
* `azurerm_app_service` support `v5.0` for the `dotnet_framework_version` ([#9251](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9251))
* `azurerm_availability_set` - adding validation to the `name` field ([#9279](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9279))
* `azurerm_cosmosdb_account` - support for the `key_vault_key_id` property allowing use of Customer Managed Keys ([#8919](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8919))
* `azurerm_eventgrid_domain` - adding validation to the `name` field ([#9281](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9281))
* `azurerm_eventgrid_domain_topic` - adding validation to the `name` field ([#9281](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9281))
* `azurerm_eventgrid_domain_topic` - adding validation to the `domain_name` field ([#9281](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9281))
* `azurerm_eventgrid_event_subscription` - adding validation to the `name` field ([#9281](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9281))
* `azurerm_eventgrid_topic` - adding validation to the `name` field ([#9281](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9281))
* `azurerm_eventgrid_system_topic` - adding validation to the `name` field ([#9281](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9281))
* `azurerm_function_app` - support for the `health_check_path` property under site_config ([#9233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9233))
* `azurerm_linux_virtual_machine` - support for managed boot diagnostics by leaving the `storage_account_uri` property empty ([#8917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8917))
* `azurerm_linux_virtual_machine_scale_set` - support for managed boot diagnostics by leaving the `storage_account_uri` property empty ([#8917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8917))
* `azurerm_log_analytics_workspace` - support for the `internet_ingestion_enabled` and `internet_query_enabled` properties ([#9033](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9033))
* `azurerm_logic_app_workflow` added logicapp name validation ([#9282](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9282))
* `azurerm_kubernetes_cluster` - support for `proximity_placement_group_id` within the `default_node_pool` block ([#9195](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9195))
* `azurerm_kubernetes_node_pool` - support for `proximity_placement_group_id` ([#9195](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9195))
* `azurerm_policy_remediation` - support for the `resource_discovery_mode` property ([#9210](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9210))
* `azurerm_point_to_site_vpn_gateway` - support for the `route` block ([#9158](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9158))
* `azurerm_virtual_network` - support for the `bgp_community` and `vnet_protection_enabled` ([#8979](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8979))
* `azurerm_vpn_gateway` - support for the `instance_0_bgp_peering_addresses` and `instance_1_bgp_peering_addresses` blocks ([#9035](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9035))
* `azurerm_windows_virtual_machine` - support for managed boot diagnostics by leaving the `storage_account_uri` property empty ([#8917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8917))
* `azurerm_windows_virtual_machine_scale_set` - support for managed boot diagnostics by leaving the `storage_account_uri` property empty ([#8917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8917))

BUG FIXES:

* `azurerm_cosmosdb_sql_database`  no longer attempts to get throughput settings when cosmos account is serverless ([#9187](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9187))
* `azurerm_kubernetes_cluster` - changing the field `availability_zones` within the `default_node_pool` block now requires recreating the resource to match the behaviour of the Azure API ([#8814](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8814))
* `azurerm_kubernetes_cluster_node_pool` - changing the field `availability_zones` now requires recreating the resource to match the behaviour of the Azure API ([#8814](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8814))
* `azurerm_log_analytics_workspace` - fix the `Free` tier from setting the `daily_quota_gb` property ([#9228](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9228))
* `azurerm_linux_virtual_machine` - the field `disk_size_gb` within the `os_disk` block can now be configured up to `4095` ([#9202](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9202))
* `azurerm_linux_virtual_machine_scale_set` - the field `disk_size_gb` within the `os_disk` block can now be configured up to `4095` ([#9202](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9202))
* `azurerm_linux_virtual_machine_scale_set` - the field `computer_name_prefix` can now end with a dash ([#9182](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9182))
* `azurerm_windows_virtual_machine` - the field `disk_size_gb` within the `os_disk` block can now be configured up to `4095` ([#9202](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9202))
* `azurerm_windows_virtual_machine_scale_set` - the field `disk_size_gb` within the `os_disk` block can now be configured up to `4095` ([#9202](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9202))

## 2.35.0 (November 05, 2020)

UPGRADE NOTES:

* `azurerm_kubernetes_cluster` - the field `enable_pod_security_policy` and `node_taints` (within the `default_node_pool` block) can no longer be configured - see below for more details ([#8982](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8982))

FEATURES:

* **New Data Source:** `azurerm_images` ([#8629](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8629))
* **New Resource:** `azurerm_firewall_policy_rule_collection_group` ([#8603](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8603))
* **New Resource:** `azurerm_virtual_hub_ip_configuration` ([#8912](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8912))
* **New Resource:** `azurerm_virtual_hub_route_table` ([#8939](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8939))

IMPROVEMENTS:

* dependencies: updating `containerservice` to API version `2020-09-01` ([#8982](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8982))
* dependencies: updating `iottimeseriesinsights` to API Version `2020-05-15` ([#9129](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9129))
* `azurerm_data_factory_linked_service_data_lake_storage_gen2` - Supports managed identity auth through `use_managed_identity ` ([#8938](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8938))
* `azurerm_firewall` - support the `dns_servers` property ([#8878](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8878))
* `azurerm_firewall_network_rule_collection` - support the `destination_fqdns` property in the `rule` block ([#8878](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8878))
* `azurerm_virtual_hub_connection` - support for the `routing` block ([#8950](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8950))

BUG FIXES:

* Fixed regression that prevented Synapse client registering in all Azure environments ([#9100](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9100))
* `azurerm_cosmosdb_mongo_database` no longer attempts to get throughput settings when cosmos account is serverless ([#8673](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8673))
* `azurerm_key_vault_access_policy` - check access policy consistency before committing to state ([#9125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9125))
* `azurerm_kubernetes_cluster` - the field `enable_pod_security_policy` can no longer be set, due to this functionality being removed from AKS as of `2020-10-15` ([#8982](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8982))
* `azurerm_kubernetes_cluster` - the field `node_taints` can no longer be set on the `default_node_pool` block, to match the behaviour of AKS ([#8982](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8982))
* `azurerm_virtual_desktop_application_group` - adding validation to the `host_pool_id` field ([#9057](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9057))
* `azurerm_virtual_desktop_workspace_application_group_association` - adding validation to the `application_group_id` field ([#9057](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9057))
* `azurerm_virtual_desktop_workspace_application_group_association` - adding validation to the `workspace_id` field ([#9057](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9057))
* `azurerm_virtual_desktop_workspace_application_group_association` - validating the ID during import is a Workspace Application Group Association ID ([#9057](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9057))
* `azurerm_postgresql_firewall_rule` - add validation for `start_ip_address` and `end_ip_address` properties ([#8963](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8963))


## 2.34.0 (October 29, 2020)

UPGRADE NOTES

* `azurerm_api_management_api` - fixing a regression introduced in v2.16 where this value for `subscription_required` was defaulted to `false` instead of `true` ([#7963](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7963))

FEATURES: 

* **New Data Source:** `azurerm_cognitive_account` ([#8773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8773))
* **New Resource:** `azurerm_log_analytics_data_export_rule` ([#8995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8995))
* **New Resource:** `azurerm_log_analytics_linked_storage_account` ([#9002](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9002))
* **New Resource:** `azurerm_security_center_auto_provisioning` ([#8595](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8595))
* **New Resource:** `azurerm_synapse_role_assignment` ([#8863](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8863))
* **New Resource:** `azurerm_vpn_site` ([#8896](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8896))

IMPROVEMENTS:

* Data Source: `azurerm_policy_definition` - can now look up built-in policy by name ([#9078](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9078))
* `azurerm_backup_policy_vm` - support for the property `instant_restore_retention_days` ([#8822](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8822))
* `azurerm_container_group` - support for the property `git_repo` within the `volume` block ([#7924](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7924))
* `azurerm_iothub` - support for the `resource_group` property within the `endpoint` block ([#8032](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8032))
* `azurerm_key_vault` - support for the `contact` block ([#8937](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8937))
* `azurerm_log_analytics_saved_search` - support for `tags` ([#9034](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9034))
* `azurerm_log_analytics_solution` - support for `tags` ([#9048](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9048))
* `azurerm_logic_app_trigger_recurrence` - support for `time_zone` [[#8829](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8829)] 
* `azurerm_policy_definition` - can now look up builtin policy by name ([#9078](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9078))

BUG FIXES: 

* `azurerm_automation_module` - raising the full error from the Azure API during creation ([#8498](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8498))
* `azurerm_api_management_api` - fixing a regression introduced in v2.16 where the value for `subscription_required` was defaulted to `false` instead of `true` ([#7963](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7963))
* `azurerm_app_service` - fixing a crash when provisioning an app service inside an App Service Environment which doesn't exist ([#8993](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8993))
* `azurerm_cdn_endpoint` - disable persisting default value for `is_compression_enabled` to state file ([#8610](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8610))
* `azurerm_databricks_workspace` correctly validate the `name` property ([#8997](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8997))
* `azurerm_dev_test_policy` - now correctly deletes ([#9077](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9077))
* `azurerm_log_analytics_workspace` - support for the `daily_quota_gb` property ([#8861](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8861))
* `azurerm_local_network_gateway` - support for the `gateway_fqdn` property ([#8998](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8998))
* `azurerm_key_vault` - prevent unwanted diff due to inconsistent casing for the `sku_name` property ([#8983](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8983))
* `azurerm_kubernetes_cluster` - fix issue where `min_count` and `max_count` couldn't be equal ([#8957](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8957))
* `azurerm_kubernetes_cluster` - `min_count` can be updated when `enable_auto_scaling` is set to true ([#8619](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8619))
* `azurerm_private_dns_zone_virtual_network_link` - fixes case issue in `name` ([#8617](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8617))
* `azurerm_private_endpoint` - fix crash when deleting private endpoint ([#9068](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9068))
* `azurerm_signalr_service` - switching the`features` block to a set so order is irrelevant ([#8815](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8815))
* `azurerm_virtual_desktop_application_group` - correctly validate the `name`property ([#9030](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9030))

## 2.33.0 (October 22, 2020)

UPGRADE NOTES

* This release includes a workaround for [a breaking change in Azure’s API related to the Extended Auditing Policy](https://github.com/Azure/azure-rest-api-specs/issues/11271) of the SQL and MSSQL resources. The Service Team have confirmed that this Regression will first roll out to all regions before the bug fix is deployed - as such this workaround will be removed in a future release once the fix for the Azure API has been rolled out to all regions.

FEATURES: 

* **New Resource:** `azurerm_service_fabric_mesh_secret` ([#8933](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8933))
* **New Resource:** `azurerm_service_fabric_mesh_secret_value` ([#8933](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8933))

IMPROVEMENTS:

* Data Source: `azurerm_shared_image_version` - exposing `os_disk_image_size_gb` ([#8904](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8904))
* `azurerm_app_configuration` - support for the `identity` block ([#8875](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8875))
* `azurerm_cosmosdb_sql_container` - support for composite indexes ([#8792](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8792))
* `azurerm_mssql_database` - do not set longterm and shortterm retention policies when using the `DW` SKUs ([#8899](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8899))
* `azurerm_mysql_firewall_rule` - validating the `start_ip_address` and `end_ip_address` fields are IP Addresses ([#8948](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8948))
* `azurerm_redis_firewall_rule` - validating the `start_ip` and `end_ip` fields are IP Addresses ([#8948](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8948))
* `azurerm_search_service` - support for the `identity` block ([#8907](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8907))
* `azurerm_sql_firewall_rule` - adding validation for the `start_ip_address` and `end_ip_address` fields ([#8935](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8935))

BUG FIXES:

* `azurerm_application_gateway` - now supports `ignore_changes` for `ssl_certificate` when using pre-existing certificates ([#8761](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8761))
* `azurerm_mssql_database` - working around a breaking change/regression in the Azure API ([#8975](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8975))
* `azurerm_mssql_database_extended_auditing_policy` - working around a breaking change/regression in the Azure API ([#8975](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8975))
* `azurerm_mssql_server` - working around a breaking change/regression in the Azure API ([#8975](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8975))
* `azurerm_mssql_server_extended_auditing_policy` - working around a breaking change/regression in the Azure API ([#8975](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8975))
* `azurerm_sql_database` - working around a breaking change/regression in the Azure API ([#8975](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8975))
* `azurerm_sql_server` - working around a breaking change/regression in the Azure API ([#8975](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8975))
* `azurerm_policy_set_definition` - Fix updates for `parameters` and `parameter_values` in `policy_definition_reference` blocks ([#8882](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8882))

## 2.32.0 (October 15, 2020)

FEATURES:

* **New data source:** `azurerm_mysql_server` ([#8787](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8787))
* **New resource:** `azurerm_security_center_setting` ([#8783](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8783))
* **New Resource:** `azurerm_service_fabric_mesh_local_network` ([#8838](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8838))
* **New resource:** `azurerm_eventgrid_system_topic` ([#8735](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8735))

IMPROVEMENTS:

* `azurerm_container_registry` - support for the `trust_policy` and `retention_policy` blocks ([#8698](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8698))
* `azurerm_security_center_contact` - override SDK creat function to handle `201` response code ([#8774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8774))

## 2.31.1 (October 08, 2020)

IMPROVEMENTS:

* `azurerm_cognitive_account` - `kind` now supports `Personalizer` ([#8860](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8860))
* `azurerm_search_service` - `sku` now supports `storage_optimized_l1` and `storage_optimized_l2` ([#8859](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8859))
* `azurerm_storage_share` - set `metadata` to `Computed` and set `acl` `start` and `expiry` to `Optional` ([#8811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8811))

BUG FIXES:

* `azurerm_dedicated_hardware_security_module` - `stamp_id` now optional to allow use in Locations which use `zones` ([#8826](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8826))
* `azurerm_storage_account`-`large_file_share_enabled` marked as computed to prevent existing storage shares from attempting to disable the default ([#8807](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8807))

## 2.31.0 (October 08, 2020)

UPGRADE NOTES

* This release updates the `azurerm_security_center_subscription_pricing` resource to use the latest version of the Security API which now allows configuring multiple Resource Types - as such a new field `resource_type` is now available. Configurations default the `resource_type` to `VirtualMachines` which matches the behaviour of the previous release - but your Terraform Configuration may need updating.

FEATURES:

* **New Resource:** `azurerm_service_fabric_mesh_application` ([#6761](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6761))
* **New Resource:** `azurerm_virtual_desktop_application_group` ([#8605](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8605))
* **New Resource:** `azurerm_virtual_desktop_workspace_application_group_association` ([#8605](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8605))
* **New Resource:** `azurerm_virtual_desktop_host_pool` ([#8605](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8605))
* **New Resource:** `azurerm_virtual_desktop_workspace` ([#8605](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8605))

IMPROVEMENTS:

* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v46.4.0` ([#8642](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8642))
* `data.azurerm_application_insights` - support for the `connection_string` property ([#8699](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8699))
* `azurerm_app_service` - support for IPV6 addresses in the `ip_restriction` property ([#8599](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8599))
* `azurerm_application_insights` - support for the `connection_string` property ([#8699](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8699))
* `azurerm_backup_policy_vm` - validate daily backups is > `7` ([#7898](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7898))
* `azurerm_dedicated_host` - add support for the `DSv4-Type1` and `sku_name` properties ([#8718](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8718))
* `azurerm_iothub` - Support for the `public_network_access_enabled` property ([#8586](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8586))
* `azurerm_key_vault_certificate_issuer` - the `org_id` property is now optional ([#8687](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8687))
* `azurerm_kubernetes_cluster_node_pool` - the `max_node`, `min_node`, and `node_count` properties can now be set to `0` ([#8300](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8300))
* `azurerm_mssql_database` - the `min_capacity` property can now be set to `0` ([#8308](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8308))
* `azurerm_mssql_database` - support for `long_term_retention_policy` and `short_term_retention_policy` blocks [[#8765](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8765)] 
* `azurerm_mssql_server` - support the `minimum_tls_version` property ([#8361](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8361))
* `azurerm_mssql_virtual_machine` - support for `storage_configuration_settings` ([#8623](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8623))
* `azurerm_security_center_subscription_pricing` - now supports per `resource_type` pricing ([#8549](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8549))
* `azurerm_storage_account` - support for the `large_file_share_enabled` property ([#8789](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8789))
* `azurerm_storage_share` - support for large quotas (up to `102400` GB) ([#8666](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8666))

BUG FIXES:

* `azurerm_function_app` - mark the `app_settings` block as computed ([#8682](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8682))
* `azurerm_function_app_slot` - mark the `app_settings` block as computed ([#8682](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8682))
* `azurerm_policy_set_definition` - corrects issue with empty `parameter_values` attribute ([#8668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8668))
* `azurerm_policy_definition` - `mode` property now enforces correct case ([#8795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8795))

## 2.30.0 (October 01, 2020)

UPGRADE NOTES

* This release renames certain fields within the `azurerm_cosmosdb_account` (data source & resource) and `azurerm_function_app_host_keys` data source to follow HashiCorp's [inclusive language guidelines](https://discuss.hashicorp.com/t/inclusive-language-changes) - where fields have been renamed, existing fields will continue to remain available until the next major version of the Azure Provider (`v3.0`)

FEATURES: 

* **New Data Source:** `azurerm_cosmosdb_sql_storedprocedure` ([#6189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6189))
* **New Data Source:** `azurerm_ip_groups` ([#8556](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8556))
* **New Resource:** `azurerm_ip_groups` ([#8556](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8556))
* **New Resource:** `azurerm_resource_group_template_deployment` ([#8672](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8672))
* **New Resource:** `azurerm_subscription_template_deployment` ([#8672](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8672))

IMPROVEMENTS:

* dependencies: updating `iothub` to `2020-03-01` ([#8688](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8688))
* dependencies: updating `storagecache` to `2020-03-01` ([#8078](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8078))
* dependencies: updating `resources` to API Version `2020-06-01` ([#8672](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8672))
* `azurerm_analysis_services_server` - support for the `S8v2` and `S9v2` SKU's ([#8707](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8707))
* `azurerm_cognitive_account` - support for the `S` `sku` ([#8639](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8639))
* `azurerm_container_group` - support for the `dns_config` block ([#7912](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7912))
* `azurerm_cosmosdb_account` - support the `zone_reduntant` property ([#8295](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8295))
* `azurerm_cosmosdb_mongo_collection` - will now respect the order of the `keys` property in the `index` block ([#8602](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8602))
* `azurerm_hpc_cache` -  support the `mtu` and `root_squash_enabled` properties ([#8078](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8078))
* `azurerm_key_vault` - add support for `enable_rbac_authorization` ([#8670](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8670))
* `azurerm_lighthouse_assignment` - limit the `scope` property to subsriptions ([#8601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8601))
* `azurerm_logic_app_workflow` - support for the `integration_service_environment_id` property ([#8504](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8504))
* `azurerm_servicebus_topic` - validate the `max_size_in_megabytes` property ([#8648](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8648))
* `azurerm_servicebus_queue` - validate the `max_size_in_megabytes` property ([#8648](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8648))
* `azurerm_servicebus_subscription_rule` - support the `correlation_filter.properties` property ([#8646](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8646))
* `azurerm_storage_management_policy` - support the `appendBlob` value for `blob_types` ([#8659](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8659))


BUG FIXES:

* `azurerm_monitor_metric_alert` - property wait when creating/updating multiple monitor metric alerts ([#8667](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8667))
* `azurerm_linux_virtual_machine_scale_set` - fix empty JSON error in `settings` and `protected_settings` when these values are not used ([#8627](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8627))

## 2.29.0 (September 24, 2020)

UPGRADE NOTES:

* `azurerm_api_management` - the value `None` has been removed from the `identity` block to match other resources, to specify an API Management Service with no Managed Identity remove the `identity` block ([#8411](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8411))
* `azurerm_container_registry` -  the `storage_account_id` property now forces a new resource as required by the updated API version ([#8477](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8477))
* `azurerm_virtual_hub_connection` - deprecating the field `vitual_network_to_hub_gateways_traffic_allowed` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_virtual_hub_connection` - deprecating the field `hub_to_vitual_network_traffic_allowed` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_virtual_wan` - deprecating the field `allow_vnet_to_vnet_traffic` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))

FEATURES: 

* **New Data Source:** `azurerm_data_share_dataset_kusto_cluster` ([#8464](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8464))
* **New Data Source:** `azurerm_databricks_workspace` ([#8502](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8502))
* **New Data Source:** `azurerm_firewall_policy` ([#7390](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7390))
* **New Data Source:** `azurerm_storage_sync_group` ([#8462](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8462))
* **New Data Source:** `azurerm_mssql_server` ([#7917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7917))
* **New Resource:** `azurerm_data_share_dataset_kusto_cluster` ([#8464](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8464))
* **New Resource:** `azurerm_firewall_policy` ([#7390](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7390))
* **New Resource:** `azurerm_mysql_server_key` ([#8125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8125))
* **New Resource:** `azurerm_postgresql_server_key` ([#8126](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8126))

IMPROVEMENTS:

* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v46.3.0` ([#8592](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8592))
* dependencies: updating `containerregistry` to `2019-05-01` ([#8477](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8477))
* Data Source: `azurerm_api_management` - export the `private_ip_addresses` attribute for primary and additional locations ([#8290](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8290))
* `azurerm_api_management` - support the `virtual_network_configuration` block for additional locations ([#8290](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8290))
* `azurerm_api_management` - export the `private_ip_addresses` attribute for additional locations ([#8290](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8290))
* `azurerm_cosmosdb_account` - support the `Serverless` value for the `capabilities` property ([#8533](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8533))
* `azurerm_cosmosdb_sql_container` - support for the `indexing_policy` property ([#8461](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8461))
* `azurerm_mssql_server` - support for the `recover_database_id` and `restore_dropped_database_id` properties ([#7917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7917))
* `azurerm_policy_set_definition` - support for typed parameter values other then string in `the policy_definition_reference` block deprecating `parameters` in favour of `parameter_vcaluess` ([#8270](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8270))
* `azurerm_search_service` - Add support for `allowed_ips` ([#8557](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8557))
* `azurerm_service_fabric_cluster` - Remove two block limit for `client_certificate_thumbprint` ([#8521](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8521))
* `azurerm_signalr_service` - support for delta updates ([#8541](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8541))
* `azurerm_spring_cloud_service` - support for configuring the `network` block ([#8568](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8568))
* `azurerm_virtual_hub_connection` - deprecating the field `vitual_network_to_hub_gateways_traffic_allowed` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_virtual_hub_connection` - deprecating the field `hub_to_vitual_network_traffic_allowed` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_virtual_hub_connection` - switching to use the now separate API for provisioning these resources ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_virtual_wan` - deprecating the field `allow_vnet_to_vnet_traffic` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_windows_virtual_machine` - support for updating the `license_type` field ([#8542](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8542))

BUG FIXES:

* `azurerm_api_management` - the value `None` for the field `type` within the `identity` block has been removed - to remove a managed identity remove the `identity` block ([#8411](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8411))
* `azurerm_app_service` - don't try to manage source_control when scm_type is `VSTSRM` ([#8531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8531))
* `azurerm_function_app` - don't try to manage source_control when scm_type is `VSTSRM` ([#8531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8531))
* `azurerm_kubernetes_cluster` - picking the first system node pool if the original `default_node_pool` has been removed ([#8503](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8503))

## 2.28.0 (September 17, 2020)

UPGRADE NOTES

* The `id` field for the `azurerm_role_definition` changed in release 2.27.0 to work around a bug in the Azure API when using management groups, where the Scope isn't returned - the existing `id` field is available as `role_definition_resource_id` from this version of the Azure Provider.

FEATURES:

* **New Data Source:** `azurerm_data_share_dataset_data_lake_gen2` [[#7907](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7907)] 
* **New Data Source:** `azurerm_servicebus_queue_authorization_rule` ([#8438](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8438))
* **New Data Source:** `azurerm_storage_sync` [[#7843](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7843)] 
* **New Resource:** `azurerm_data_share_dataset_data_lake_gen2` ([#7907](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7907))
* **New Resource:** `azurerm_lighthouse_definition` ([#6560](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6560))
* **New Resource:** `azurerm_lighthouse_assignment` ([#6560](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6560))
* **New Resource:** `azurerm_mssql_server_extended_auditing_policy`  ([#8447](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8447))
* **New Resource:** `azurerm_storage_sync` ([#7843](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7843))
* **New Resource:** `azurerm_synapse_sql_pool` ([#8095](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8095))

IMPROVEMENTS:

* Data Source: `azurerm_app_service_environment` - Expose vip information of an app service environment ([#8487](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8487))
* Data Source: `azurerm_function_app` - export the `identity` block ([#8389](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8389))
* `azurerm_app_service_hybrid_connection` - support relays in different namespaces ([#8370](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8370))
* `azurerm_cosmosdb_cassandra_keyspace` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_gremlin_database` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_gremlin_graph` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_mongo_collection` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_mongo_database` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_sql_container` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_sql_database` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_table` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_firewall` - support the `management_ip_configuration` block ([#8235](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8235))
* `azurerm_storage_account_customer_managed_key` - support for key rotation ([#7836](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7836))

BUG FIXES:

* Data Source: `azurerm_function_app_host_keys` - Fix a crash when null ID sometimes returned by API ([#8430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8430))
* `azurerm_cognitive_account` - correctly wait on update logic ([#8386](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8386))
* `azurerm_eventhub_consumer_group` - allow the `name` property to be set to `$Default` ([#8388](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8388))
* `azurerm_kubernetes_cluster` - ensure the OMS Agent Log Analytics Workspace case is preserved after disabling/enabling ([#8374](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8374))
* `azurerm_management_group_id` - loosen case restritions during parsing of management group ID ([#8024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8024))
* `azurerm_packet_capture` - fix to ID path to match change in API ([#8167](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8167))
* `azurerm_role_definition` - expose `role_definition_resource_id` ([#8492](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8492))

## 2.27.0 (September 10, 2020)

UPGRADE NOTES

* The `id` field for the `azurerm_role_definition` has changed in this release to work around a bug in the Azure API when using management groups, where the Scope isn't returned - the existing `id` field is available as `role_definition_resource_id` on the new resource from version 2.28.0 of the Azure Provider.

FEATURES:

* **New Data Source:** `azurerm_attestation_provider` ([#7885](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7885))
* **New Data Source:** `azurerm_function_app_host_keys` ([#7902](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7902))
* **New Data Source:** `azurerm_lb_rule` ([#8365](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8365))
* **New Resource:** `azurerm_mssql_database_extended_auditing_policy` ([#7793](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7793))
* **New Resource:** `azurerm_attestation_provider` ([#7885](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7885))
* **New Resource:** `azurerm_api_management_api_diagnostic` ([#7873](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7873))
* **New Resource:** `azurerm_data_factory_linked_service_azure_sql_database` ([#8349](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8349))

IMPROVEMENTS:

* Data Source: `azurerm_virtual_network_gateway` - exposing `aad_audience`, `aad_issuer` and `aad_tenant` within the `vpn_client_configuration` block ([#8294](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8294))
* `azurerm_cosmosdb_account` - supporting the value `AllowSelfServeUpgradeToMongo36` for the `name` field within the `capabilities` block ([#8335](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8335))
* `azurerm_linux_virtual_machine` - Add support for `encryption_at_host_enabled` ([#8322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8322))
* `azurerm_linux_virtual_machine_scale_set` - Add support for `encryption_at_host_enabled` ([#8322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8322))
* `azurerm_servicebus_subscription` - add support for `dead_lettering_on_filter_evaluation_error` ([#8412](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8412))
* `azurerm_spring_cloud_app` - support for the `identity` block ([#8336](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8336))
* `azurerm_storage_share_directory` - Update name validation ([#8366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8366))
* `azurerm_virtual_network_gateway` - support for `aad_audience`, `aad_issuer` and `aad_tenant` within the `vpn_client_configuration` block ([#8294](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8294))
* `azurerm_windows_virtual_machine` - Add support for `encryption_at_host_enabled` ([#8322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8322))
* `azurerm_windows_virtual_machine_scale_set` - Add support for `encryption_at_host_enabled` ([#8322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8322))

BUG FIXES:

* `azurerm_api_management_x.y.api_name` - validation fix ([#8409](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8409))
* `azurerm_application_insights_webtests` - Fix an issue where the `kind` property is sometimes set to `null` ([#8372](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8372))
* `azurerm_cognitive_account` - Fixes a crash when provisioning a QnAMaker and supports AnomalyDetector ([#8357](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8357))
* `azurerm_linux_virtual_machine` - Add WaitForState on VM delete ([#8383](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8383))
* `azurerm_network_security_group` - fixed issue where updates would fail for resource ([#8384](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8384))
* `azurerm_role_definition` - fixed delete operation when role is scoped to Management Group ([#6107](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6107))
* `azurerm_windows_virtual_machine` - Add WaitForState on VM delete ([#8383](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8383))

## 2.26.0 (September 04, 2020)

UPGRADE NOTES:

* **Opt-In Beta:** This release introduces an opt-in beta for in-line Virtual Machine Scale Set Extensions. This functionality enables the resource to be used with Azure Service Fabric and other extensions that may require creation time inclusion on Scale Set members. Please see the documentation for `azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` for information.

FEATURES:

* **New Resource:** `azurerm_log_analytics_saved_search` ([#8253](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8253))

IMPROVEMENTS:

* dependencies: updating `loganalytics` to `2020-03-01-preview` ([#8234](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8234))
* `azurerm_api_management_subscription` - Support `allow_tracing property` ([#7969](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7969))
* `azurerm_application_gateway ` - Add support for `probe.properties.port` ([#8278](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8278))
* `azurerm_linux_virtual_machine_scale_set` - Beta support for `extension` blocks ([#8222](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8222))
* `azurerm_log_analytics_workspace`- the `sku` value is now optional and defaults to `PerGB2018` ([#8272](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8272))
* `azurerm_windows_virtual_machine_scale_set` - Beta support for `extension` blocks ([#8222](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8222))

BUG FIXES:

* `azurerm_cdn_endpoint` - fixing the casing of the Resource ID to be consistent ([#8237](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8237))
* `azurerm_cdn_profile` - fixing the casing of the Resource ID to be consistent ([#8237](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8237))
* `azurerm_key_vault_key` - updating the latest version of the key when updating metadata ([#8304](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8304))
* `azurerm_key_vault_secret` - updating the latest version of the secret when updating metadata ([#8304](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8304))
* `azurerm_linux_virtual_machine` - allow updating `allow_extension_operations` regardless of the value of `provision_vm_agent` (for when the VM Agent has been installed manually) ([#8001](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8001))
* `azurerm_linux_virtual_machine_scale_set` - working around a bug in the Azure API by always sending the existing Storage Image Reference during updates ([#7983](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7983))
* `azurerm_network_interface_application_gateway_association` - handling the Network Interface being deleted during a refresh ([#8267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8267))
* `azurerm_network_interface_application_security_group_association` - handling the Network Interface being deleted during a refresh ([#8267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8267))
* `azurerm_network_interface_backend_address_pool_association` - handling the Network Interface being deleted during a refresh ([#8267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8267))
* `azurerm_network_interface_nat_rule_association_resource` - handling the Network Interface being deleted during a refresh ([#8267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8267))
* `azurerm_network_interface_network_security_group_association` - handling the Network Interface being deleted during a refresh ([#8267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8267))
* `azurerm_windows_virtual_machine` - allow updating `allow_extension_operations` regardless of the value of `provision_vm_agent` (for when the VM Agent has been installed manually) ([#8001](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8001))
* `azurerm_windows_virtual_machine_scale_set` - working around a bug in the Azure API by always sending the existing Storage Image Reference during updates ([#7983](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7983))

## 2.25.0 (August 27, 2020)

UPGRADE NOTES:

* `azurerm_container_group` - The `secure_environment_variables` field within the `container` now maps keys with empty values, which differs from previous versions of this provider which ignored empty values ([#8151](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8151))

FEATURES:

* **New Resource** `azurerm_spring_cloud_certificate` ([#8067](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8067))

IMPROVEMENTS:

* dependencies: updating `keyvault` to `2019-09-01` ([#7822](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7822))
* `azurerm_app_service_slot_virtual_network_swift_connection` - adding validation that the `app_service_id` is an App Service / Function App ID ([#8111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8111))
* `azurerm_app_service_slot_virtual_network_swift_connection` - adding validation that the `subnet` is a Subnet ID ([#8111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8111))
* `azurerm_batch_pool` - Remove `network_configuration` from update payload ([#8189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8189))
* `azurerm_frontdoor_firewall_policy` - `match_variable` within the `match_condition` block can now be set to `SocketAddr` ([#8244](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8244))
* `azurerm_linux_virtual_machine_scale_set` - `upgrade_mode="Automatic"` no longer requires health probe ([#6667](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6667))
* `azurerm_key_vault` - support for `soft_delete_retention_days` ([#7822](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7822))
* `azurerm_shared_image` - Support for `purchase_plan` ([#8124](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8124))
* `azurerm_shared_image_gallery` - validating at import time that the ID is for a Shared Image Gallery ([#8240](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8240))
* `azurerm_windows_virtual_machine_scale_set` - `upgrade_mode="Automatic"` no longer requires health probe ([#6667](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6667))

BUG FIXES:

* Data Source: `azurerm_app_service` - ensuring the `site_config` block is correctly set into the state ([#8212](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8212))
* Enhanced Validation: supporting "centralindia", "southindia" and "westindia" as valid regions in Azure Public (working around invalid data from the Azure API) ([#8217](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8217))
* `azurerm_application_gateway` - allow setting `ip_addresses` within the `backend_address_pool` block to an empty list ([#8210](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8210))
* `azurerm_application_gateway` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_container_group` - the `secure_environment_variables` field within the `container` now maps keys with empty values ([#8151](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8151))
* `azurerm_dedicated_host` - waiting for the resource to be gone 20 times rather than 10 to work around an API issue ([#8221](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8221))
* `azurerm_dedicated_host_group` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_firewall` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_hardware_security_module` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_lb` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_linux_virtual_machine` - support for updating `ultra_ssd_enabled` within the `additional_capabilities` block without recreating the virtual machine ([#8015](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8015))
* `azurerm_linux_virtual_machine_scale_set` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_managed_disk` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_nat_gateway` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_orchestrated_virtual_machine_scale_set` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_public_ip_prefix` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_public_ip` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_redis_cache` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_virtual_machine` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_virtual_machine_scale_set` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_windows_virtual_machine` - support for updating `ultra_ssd_enabled` within the `additional_capabilities` block without recreating the virtual machine ([#8015](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8015))
* `azurerm_windows_virtual_machine_scale_set` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))

## 2.24.0 (August 20, 2020)

FEATURES:

* **New Resource:** `azurerm_synapse_spark_pool` ([#7886](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7886))

IMPROVEMENTS:

* dependencies: update `containerinstance` to API version `2019-12-01` ([#8110](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8110))
* `azurerm_api_management_api` - now supports `oauth2_authorization` and `openid_authentication` ([#7617](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7617))
* `azurerm_policy_definition` - `mode` can now be updated without recreating the resource ([#7976](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7976))

BUG FIXES:

* `azurerm_frontdoor` - ensuring all fields are set into the state ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor` - rewriting case-inconsistent Resource ID's to ensure they're reliable ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor_firewall_policy` - ensuring all fields are set into the state ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor_firewall_policy` - rewriting case-inconsistent Resource ID's to ensure they're reliable ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor_custom_https_configuration` - ensuring all fields are set into the state ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor_custom_https_configuration` - ensuring the `resource_group_name` field is set into the state ([#8173](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8173))
* `azurerm_frontdoor_custom_https_configuration` - rewriting case-inconsistent Resource ID's to ensure they're reliable ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor_custom_https_configuration` - updating the ID to use the frontendEndpoint's Resource ID rather than a custom Resource ID ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_lb` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_backend_address_pool` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_backend_address_pool` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_backend_address_pool` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_nat_pool` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_nat_pool` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_nat_pool` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_nat_rule` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_nat_rule` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_nat_rule` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_outbound_rule` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_outbound_rule` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_outbound_rule` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_probe` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_probe` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_probe` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_rule` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_rule` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_rule` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_storage_account` - only sending `allow_blob_public_access` and `min_tls_version` in Azure Public since these are currently not supported in other regions ([#8148](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8148))

## 2.23.0 (August 13, 2020)

FEATURES:

* **New Resource:** `azurerm_integration_service_environment` ([#7763](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7763))
* **New Resource:** `azurerm_redis_linked_server` ([#8026](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8026))
* **New Resource:** `azurerm_synapse_firewall_rule` ([#7904](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7904))

IMPROVEMENTS:

* dependencies: updating `containerservice` to `2020-04-01` ([#7894](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7894))
* dependencies: updating `mysql` to `2020-01-01` ([#8062](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8062))
* dependencies: updating `postgresql` to `2020-01-01` ([#8045](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8045))
* Data Source: `azurerm_app_service` now exports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* Data Source: `azurerm_function_app` now exports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* Data Source: `azurerm_function_app` now exports `site_config` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_app_service` now supports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports full `ip_restriction` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports full `scm_ip_restriction` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_eventhub_namespace` - support for the `identity` block ([#8065](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8065))
* `azurerm_postgresql_server` - support for the `identity` block ([#8044](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8044))
* `azurerm_site_recovery_replicated_vm` - support setting `target_network_id` and `network_interface` on failover ([#5688](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5688))
* `azurerm_storage_account` - support `static_website` for `BlockBlobStorage` account type ([#7890](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7890))
* `azurerm_storage_account` - filter `allow_blob_public_access` and `min_tls_version` from Azure US Government ([#8092](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8092))

BUG FIXES:

* All resources using a `location` field - allowing the value `global` when using enhanced validation ([#8042](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8042))
* Data Source: `azurerm_api_management_user` - `user_id` now accepts single characters ([#7975](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7975))
* `azurerm_application_gateway` - enforce case for the `rule_type` property ([#8061](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8061))
* `azurerm_iothub_consumer_group` - lock during creation and deletion to workaround an API issue ([#8041](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8041))
* `azurerm_iothub` - the `endpoint` and `route` lists can now be cleared by setting them to `[]` ([#8028](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8028))
* `azurerm_linux_virtual_machine` - handling machines which are already stopped/deallocated ([#8000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8000))
* `azurerm_mariadb_virtual_network_rule` will now work across subscriptions ([#8100](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8100))
* `azurerm_monitor_metric_alert_resource` - continue using `SingleResourceMultiMetricCriteria` for existing alerts ([#7995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7995))
* `azurerm_mysql_server` - prevent a non empty plan when using `threat_detection_policy` ([#7981](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7981))
* `azurerm_orchestrated_virtual_machine_scale_set` - allow `single_placement_group` to be `true` ([#7821](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7821))
* `azurerm_mysql_server` - support for the `identity` block ([#8059](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8059))
* `azurerm_storage_account` - set default for `min_tls_version` to `TLS_10` ([#8152](https://github.com/terraform-providers/terraform-provider-azurerm/pull/8152))
* `azurerm_traffic_manager_profile` - updating no longer clears all endpoints ([#7846](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7846))
* `azurerm_windows_virtual_machine` - handling machines which are already stopped/deallocated [[#8000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8000)]'
* `azurerm_data_factory_dataset_delimited_text` - fix issue with property `azure_blob_storage_account` ([#7953](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7953))

## 2.22.0 (August 07, 2020)

DEPENDENCIES:

* updating `github.com/Azure/azure-sdk-for-go` to `v44.2.0` ([#7933](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7933))

IMPROVEMENTS:

* `azurerm_cosmosdb_account` - support `DisableRateLimitingResponses` with the `capabilities` property ([#8016](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8016))
* `azurerm_storage_account` - support for the `min_tls_version` property ([#7879](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7879))
* `azurerm_storage_account_sas` - support for the `signed_version attribute` property ([#8020](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8020))
* `azurerm_servicebus_queue` - support for the `enable_batched_operations`, `status`, `forward_to`, and `forward_dead_lettered_messages_to` ([#7990](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7990))

BUG FIXES:

* Data Source: `azurerm_key_vault_certificate` - fixing a crash when using acmebot certificates ([#8029](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8029))
* `azurerm_iothub_shared_access_policy` - prevent `primary_connection_string` & `secondary_connection_string` from regenerating during every apply ([#8017](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8017))

## 2.21.0 (July 31, 2020)

DEPENDENCIES:

* updating `search` to `2020-03-13` ([#7867](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7867))
* updating `go-azure-helpers` to `v0.11.2` ([#7911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7911))

FEATURES:

* **New Data Source:** `azurerm_data_share_dataset_data_lake_gen1` ([#7840](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7840))
* **New Resource:** `azurerm_dedicated_hardware_security_module` ([#7727](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7727))

IMPROVEMENTS:
* `azurerm_api_management_identity_provider_aad` - Support for `signin_tenant` ([#7901](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7901))
* `azurerm_app_service_plan` - update the relation between `kind` and `reserved` ([#7943](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7943))
* `azurerm_automation_runbook` - recreate `azurerm_automation_job_schedule` after an update ([#7555](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7555))
* `azurerm_app_service_slot` - support for the `application_logs.file_system` ([#7311](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7311))
* `azurerm_firewall` - no longer requires a `zone` ([#7817](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7817))
* `azurerm_function_app_slot` - support for the `site_config.auto_swap_slot_name` property ([#7859](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7859))
* `azurerm_kubernetes_cluster` - support for in-place upgrade from `Free` to `Paid` for `sku_tier` ([#7927](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7927))
* `azurerm_monitor_scheduled_query_rules_alert` - `action.0.custom_webhook_payload` is now sent as empty to allow for Azure's default to take effect([#7838](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7838))
* `azurerm_search_service` - support for the `public_network_access_enabled` property ([#7867](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7867))
* `azurerm_servicebus_subscription` - support for the `status` property ([#7852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7852))

BUG FIXES:

* `azurerm_automation_runbook` - allow `publish_content_link` resource to not be set ([#7824](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7824))
* `azurerm_api_management_named_value` - the `value` has been marked as sensitive to hide secret values ([#7819](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7819))
* `azurerm_cognitive_account` - allow `qname_runtime_endpoint` to not be set ([#7916](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7916))
* `azurerm_iothub_dps` - the only valid value for the `sku` property for the API is now `S1` ([#7847](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7847))
* `azurerm_eventgrid_event_subscription` - deprecate the `topic_name` as it is now readonly in the API ([#7871](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7871))
* `azurerm_kubernetes_cluster` - updates will no longer fail when using managed AAD integration ([#7874](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7874))

## 2.20.0 (July 23, 2020)

UPGRADE NOTES

* **Enhanced Validation for Locations** - the Azure Provider now validates that the value for the `location` argument is a supported Azure Region within the Azure Environment being used (from the Azure Metadata Service) - which allows us to catch configuration errors for this field at `terraform plan` time, rather than during a `terraform apply`. This functionality is now enabled by default, and can be opted-out of by setting the Environment Variable `ARM_PROVIDER_ENHANCED_VALIDATION` to `false`
* `azurerm_storage_account` - will now default `allow_blob_public_access` to false to align with the portal and be secure by default ([#7784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7784))

DEPENDENCIES:

* updating `github.com/Azure/azure-sdk-for-go` to `v44.1.0` ([#7774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7774))
* updating `cosmos` to `2020-04-01` ([#7597](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7597))

FEATURES: 

* **New Data Source:** `azurerm_synapse_workspace` ([#7517](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7517))
* **New Resource:** `azurerm_data_share_dataset_data_lake_gen1` - add `dataset_data_lake_gen1` suppport for `azurerm_data_share` ([#7511](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7511))
* **New Resource:** `azurerm_frontdoor_custom_https_configuration` - move the front door `custom_https_configuration` to its own resource to allow for parallel creation/update of custom https certificates. ([#7498](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7498))
* **New Resource:** `azurerm_kusto_cluster_customer_managed_key` ([#7520](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7520))
* **New Resource:** `azurerm_synapse_workspace` ([#7517](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7517))

IMPROVEMENTS:

* `azurerm_cosmos_db_account` - add support for the `enable_free_tier` property ([#7814](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7814))

BUG FIXES:

* Data Source: `azurerm_private_dns_zone` - fix a crash when the zone does not exist ([#7783](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7783))
* `azurerm_application_gateway` - fix crash with `gateway_ip_configuration` ([#7789](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7789))
* `azurerm_cosmos_account` - the `geo_location.prefix` property has been deprecated as service no longer accepts it as an input since Apr 25, 2019 ([#7597](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7597))
* `azurerm_monitor_autoscale_setting` - fix crash in `notification` ([#7835](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7835))
* `azurerm_storage_account` - will now default `allow_blob_public_access` to false to align with the portal and be secure by default ([#7784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7784))

## 2.19.0 (July 16, 2020)

UPGRADE NOTES:

* HDInsight 3.6 will be retired (in Azure Public) on 2020-12-30 - HDInsight 4.0 does not support ML Services, RServer or Storm Clusters - as such the `azurerm_hdinsight_ml_services_cluster`, `azurerm_hdinsight_rserver_cluster` and `azurerm_hdinsight_storm_cluster` resources are deprecated and will be removed in the next major version of the Azure Provider. ([#7706](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7706))
* provider: no longer auto register the Microsoft.StorageCache RP ([#7768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7768))

FEATURES:

* **New Data source:** `azurerm_route_filter` ([#6341](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6341))
* **New Resource:** `azurerm_route_filter` ([#6341](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6341))

IMPROVEMENTS:

* dependencies: updating to v44.0.0 of `github.com/Azure/azure-sdk-for-go` ([#7616](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7616))
* dependencies: updating the `machinelearning` API to version `2020-04-01` ([#7703](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7703))
* Data Source: `azurerm_storage_account` - exposing `allow_blob_public_access` ([#7739](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7739))
* Data Source: `azurerm_dns_zone` - now provides feedback if a `resource_group_name` is needed to resolve ambiguous zone ([#7680](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7680))
* `azurerm_automation_schedule` - Updated validation for timezone strings ([#7754](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7754))
* `azurerm_express_route_circuit_peering` - support for the `route_filter_id` property ([#6341](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6341))
* `azurerm_kubernetes_cluster` - no longer sending the `kubernetes_dashboard` addon in Azure China since this is not supported in this region ([#7714](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7714))
* `azurerm_local_network_gateway`- `address_space` order can now be changed ([#7745](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7745))
* `azurerm_machine_learning_workspace` - adding the field `high_business_impact` ([#7703](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7703))
* `azurerm_monitor_metric_alert` - support for multiple scopes and associated criteria ([#7159](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7159))
* `azurerm_mssql_database` `elastic_pool_id` remove forcenew ([#7628](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7628))
* `azurerm_policy_assignment` - support for `metadata` property ([#7725](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7725))
* `azurerm_policy_set_definition` - support for the `policy_definition_reference_id` property ([#7018](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7018))
* `azurerm_storage_account` - support for configuring `allow_blob_public_access` ([#7739](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7739))
* `azurerm_storage_container` - container creation will retry if a container of the same name has not completed its delete operation ([#7179](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7179))
* `azurerm_storage_share` - share creation will retry if a share of the same name has not completed its previous delete operation ([#7179](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7179))
* `azurerm_virtual_network_gateway_connection` - support for the `traffic_selector_policy` block ([#6586](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6586))
* `azurerm_orchestrated_virtual_machine_scale_set` - support for the `proximity_placement_group_id` property ([#7510](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7510))


BUG FIXES:

* provider: deprecating `metadata_url` to `metadata_host` since this is a hostname ([#7740](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7740))
* `azurerm_*_virtual_machine` - `allow_extensions_operations` can now be updated ([#7749](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7749))
* `azurerm_eventhub_namespace` - changing to `zone_redundant` now force a new resource ([#7612](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7612))
* `azurerm_express_route_circuit` - fix eventual consistency issue in create ([#7753](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7753))
* `azurerm_express_route_circuit` - fix potential crash ([#7776](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7776))
* `azurerm_managed_disk` - allow up to `65536` GB for the `disk_size_gb` property ([#7689](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7689))
* `azurerm_machine_learning_workspace` - waiting until the Machine Learning Workspace has been fully deleted ([#7635](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7635))
* `azurerm_mysql_server` - `ssl_minimal_tls_version_enforced` now correctly set in updates ([#7307](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7307))
* `azurerm_notification_hub` - validating that the ID is in the correct format when importing the resource ([#7690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7690))
* `azurerm_redis_cache` - fixing a bug when provisioning with authentication disabled ([#7734](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7734))
* `azurerm_virtual_hub` - the field `address_prefix` is now `ForceNew` to match the behaviour of the Azure API ([#7713](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7713))
* `azurerm_virtual_hub_connection` - using the delete timeout if specified ([#7731](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7731))

## 2.18.0 (July 10, 2020)

FEATURES:

* `metadata_url` can be set at the provider level to use an environment provided by a specific url ([#7664](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7664))
* **New Data Source:** `azurerm_key_vault_certificate_issuer` ([#7074](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7074))
* **New Data Source:** `azurerm_web_application_firewall_policy` ([#7469](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7469))
* **New Resource:** `azurerm_automation_connection` ([#6847](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6847))
* **New Resource:** `azurerm_automation_connection_certificate` ([#6847](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6847))
* **New Resource:** `azurerm_automation_connection_classic_certificate` ([#6847](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6847))
* **New Resource:** `azurerm_automation_connection_service_pricipal` ([#6847](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6847))
* **New Resource:** `azurerm_app_service_slot_virtual_network_swift_connection` ([#5916](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5916))
* **New Resource:** `azurerm_data_factory_dataset_azure_blob` ([#6366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6366))
* **New Resource:** `azurerm_data_factory_dataset_cosmosdb_sqlapi` ([#6366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6366))
* **New Resource:** `azurerm_data_factory_dataset_delimited_text` ([#6366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6366))
* **New Resource:** `azurerm_data_factory_dataset_http` ([#6366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6366))
* **New Resource:** `azurerm_data_factory_dataset_json` ([#6366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6366))
* **New Resource:** `azurerm_data_factory_linked_service_azure_blob_storage` ([#6366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6366))
* **New Resource:** `azurerm_data_factory_linked_service_azure_file_storage` ([#6366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6366))
* **New Resource:** `azurerm_data_factory_linked_service_azure_file_storage` ([#6366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6366))
* **New Resource:** `azurerm_data_factory_linked_service_cosmosdb` ([#6366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6366))
* **New Resource:** `azurerm_data_factory_linked_service_sftp` ([#6366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6366))
* **New Resource:** `azurerm_data_factory_linked_service_sftp` ([#6366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6366))
* **New Resource:** `azurerm_key_vault_certificate_issuer` ([#7074](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7074))
* **New Resource:** `azurerm_kusto_attached_database_configuration` ([#7377](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7377))
* **New Resource:** `azurerm_kusto_database_principal_assignment` ([#7484](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7484))
* **New Resource:** `azurerm_mysql_active_directory_administrator` ([#7621](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7621))

IMPROVEMENTS:

* dependencies: updating `github.com/tombuildsstuff/giovanni` to `v0.11.0` ([#7608](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7608))
* dependencies: updating `network` to `2020-05-01` ([#7585](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7585))
* Data Source: `azurerm_eventhub_namespace` - exposing the `dedicated_cluster_id` field ([#7548](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7548))
* `azurerm_cosmosdb_account` - support for the `ignore_missing_vnet_service_endpoint` property ([#7348](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7348))
* `azurerm_application_gateway` - support for the `firewall_policy_id` attribute within the `http_listener` block ([#7580](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7580))
* `azurerm_eventhub_namespace` - support for configuring the `dedicated_cluster_id` field ([#7548](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7548))
* `azurerm_eventhub_namespace` - support for setting `partition_count` to `1024` when using a Dedicated Cluster ([#7548](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7548))
* `azurerm_eventhub_namespace` - support for setting `retention_count` to `90` when using a Dedicated Cluster ([#7548](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7548))
* `azurerm_hdinsight_hadoop_cluster` - now supports Azure Monitor ([#7045](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7045))
* `azurerm_hdinsight_hbase_cluster` - now supports external metastores ([#6969](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6969))
* `azurerm_hdinsight_hbase_cluster` - now supports Azure Monitor ([#7045](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7045))
* `azurerm_hdinsight_interactive_query_cluster` - now supports external metastores ([#6969](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6969))
* `azurerm_hdinsight_interactive_query_cluster` - now supports Azure Monitor ([#7045](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7045))
* `azurerm_hdinsight_kafka_cluster` - now supports external metastores ([#6969](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6969))
* `azurerm_hdinsight_kafka_cluster` - now supports external Azure Monitor ([#7045](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7045))
* `azurerm_hdinsight_spark_cluster` - now supports external metastores ([#6969](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6969))
* `azurerm_hdinsight_spark_cluster` - now supports external Azure Monitor ([#7045](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7045))
* `azurerm_hdinsight_storm_cluster` - now supports external metastores ([#6969](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6969))
* `azurerm_hdinsight_storm_cluster` - now supports external Azure Monitor ([#7045](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7045))
* `azurerm_policy_set_definition` - the `management_group_id` property has been deprecated in favour of `management_group_name` to align with the behaviour in `azurerm_policy_definition` ([#6943](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6943))
* `azurerm_kusto_cluster` - support for the `language_extensions` property ([#7421](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7421))
* `azurerm_kusto_cluster` - Support for the `optimized_auto_scale` property ([#7371](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7371))
* `azurerm_mysql_server` - support for the `threat_detection_policy` property ([#7156](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7156))
* `azurerm_mssql_database` - the `sku_name` property now only forces a new resource for the `HS` (HyperScale) family ([#7559](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7559))
* `azurerm_web_application_firewall_policy` - allow setting `version` to `0.1` (for when `type` is set to `Microsoft_BotManagerRuleSet`) ([#7579](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7579))
* `azurerm_web_application_firewall_policy` - support the `transforms` property in the `custom_rules.match_conditions` block ([#7545](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7545))
* `azurerm_web_application_firewall_policy` - support the `request_body_check`, `file_upload_limit_in_mb`, and `max_request_body_size_in_kb` properties in the `policy_settings` block ([#7363](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7363))

BUG FIXES: 

* `azurerm_api_management_api_operation_policy` - correctly parse XLM ([#7345](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7345))
* `azurerm_application_insights_api_key` - now correctly checks if the resource exists upon creation ([#7650](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7650))
* `azurerm_api_management_identity_provider_aad` - fix perpetual diff on the `client_secret` property ([#7529](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7529))
* `azurerm_eventhub_namespace_authorization_rule` - correctly update old resource IDs ([#7622](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7622))
* `azurerm_policy_remediation` - removing the validation for the `policy_definition_reference_id` field since this isn't a Resource ID ([#7600](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7600))
* `azurerm_storage_data_lake_gen2_filesystem` - prevent a crash during plan if storage account was deleted ([#7378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7378))

## 2.17.0 (July 03, 2020)

UPGRADE NOTES:

* `azurerm_hdinsight_hadoop_cluster` - the `enabled` property within the `gateway` block now defaults to `true` and cannot be disabled, due to a behavioural change in the Azure API ([#7111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7111))
* `azurerm_hdinsight_hbase_cluster` - the `enabled` property within the `gateway` block now defaults to `true` and cannot be disabled, due to a behavioural change in the Azure API ([#7111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7111))
* `azurerm_hdinsight_interactive_query_cluster` - the `enabled` property within the `gateway` block now defaults to `true` and cannot be disabled, due to a behavioural change in the Azure API ([#7111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7111))
* `azurerm_hdinsight_kafka_cluster` - the `enabled` property within the `gateway` block now defaults to `true` and cannot be disabled, due to a behavioural change in the Azure API ([#7111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7111))
* `azurerm_hdinsight_ml_services_cluster` - the `enabled` property within the `gateway` block now defaults to `true` and cannot be disabled, due to a behavioural change in the Azure API ([#7111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7111))
* `azurerm_hdinsight_rserver_cluster` - the `enabled` property within the `gateway` block now defaults to `true` and cannot be disabled, due to a behavioural change in the Azure API ([#7111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7111))
* `azurerm_hdinsight_spark_cluster` - the `enabled` property within the `gateway` block now defaults to `true` and cannot be disabled, due to a behavioural change in the Azure API ([#7111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7111))
* `azurerm_hdinsight_storm_cluster` - the `enabled` property within the `gateway` block now defaults to `true` and cannot be disabled, due to a behavioural change in the Azure API ([#7111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7111))

FEATURES: 

* **New Resource:** `azurerm_kusto_cluster_principal_assignment` ([#7533](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7533))

IMPROVEMENTS:

* dependencies: updating to v43.2.0 of `github.com/Azure/azure-sdk-for-go` ([#7546](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7546))
* Data Source: `azurerm_eventhub_namespace` - exposing the `zone_redundant` property ([#7534](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7534))
* Data Source: `azurerm_postgresql_server` - exposing `sku_name` ([#7523](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7523))
* `azurerm_app_service_environment` - the property `user_whitelisted_ip_ranges` has been deprecated and renamed to `allowed_user_ip_cidrs` to clarify the function and expected format ([#7499](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7499))
* `azurerm_eventhub_namespace` - support for the `zone_redundant` property ([#7534](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7534))
* `azurerm_key_vault_certificate` - exposing the `certificate_attribute` block ([#7387](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7387))
* `azurerm_kusto_cluster` - Support `trusted_external_tenants` ([#7374](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7374))
* `azurerm_sentinel_alert_rule_ms_security_incident` - the property `text_whitelist` has been deprecated and renamed to `display_name_filter` to better match the api ([#7499](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7499))
* `azurerm_shared_image` - support for specialized images via the `specialized` property ([#7277](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7277))
* `azurerm_shared_image_version` - support for specialized images via the `specialized` property ([#7277](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7277))
* `azurerm_spring_cloud_service` - support for `sku_name` ([#7531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7531))
* `azurerm_spring_cloud_service` - support for the `trace` block ([#7531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7531))

BUG FIXES: 

* `azurerm_api_management_named_value` - polling until the property is fully created ([#7547](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7547))
* `azurerm_api_management_property` - polling until the property is fully created ([#7547](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7547))
* `azurerm_linux_virtual_machine_scale_set` - using the provider feature `roll_instances_when_required` when `upgrade_mode` is set to `Manual` ([#7513](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7513))
* `azurerm_marketplace_agreement` - fix issue around import ([#7515](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7515))
* `azurerm_windows_virtual_machine_scale_set` - using the provider feature `roll_instances_when_required` when `upgrade_mode` is set to `Manual` ([#7513](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7513))

## 2.16.0 (June 25, 2020)

DEPENDENCIES:

* updating `github.com/Azure/go-autorest/azure/cli` to `v0.3.1` ([#7433](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7433))

FEATURES:

* **New Resource:** `azurerm_postgresql_active_directory_administrator` ([#7411](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7411))

IMPROVEMENTS:

* authentication: Azure CLI - support for access tokens in custom directories ([#7433](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7433))
* `azurerm_api_management_api` - support for the `subscription_required` property ([#4885](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4885))
* `azurerm_app_service_environment` - support a value of `Web, Publishing` for the `internal_load_balancing_mode` property ([#7346](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7346))
* `azurerm_kusto_cluster` - support for the `identity` block ([#7367](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7367))
* `azurerm_kusto_cluster` - support for `virtual_network_configuration` block ([#7369](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7369))
* `azurerm_kusto_cluster` - supoport for the `zone` property ([#7373](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7373))
* `azurerm_firewall` - support for configuring `threat_intel_mode` ([#7437](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7437))
* `azurerm_management_group` - waiting until the Management Group has been fully replicated after creating ([#7473](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7473))
* `azurerm_monitor_activity_log_alert` - support for the fields `recommendation_category`, `recommendation_impact` and `recommendation_type` in the `criteria` block ([#7458](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7458))
* `azurerm_mssql_database` - support up to `5` for the `min_capacity` property ([#7457](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7457))
* `azurerm_mssql_database` - support `GP_S_Gen5` SKUs up to `GP_S_Gen5_40` ([#7453](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7453))

BUG FIXES: 

* `azurerm_api_management_api` - allowing dots as a prefix of the `name` field ([#7478](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7478))
* `azurerm_function_app` - state fixes for `app_settings` ([#7440](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7440))
* `azurerm_hdinsight_hadoop_cluster` - fixes for node and instance count validation ([#7430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7430))
* `azurerm_hdinsight_hbase_cluster` - fixes for node and instance count validation ([#7430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7430))
* `azurerm_hdinsight_interactive_query_cluster` - fixes for node and instance count validation ([#7430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7430))
* `azurerm_hdinsight_kafka_cluster` - fixes for node and instance count validation ([#7430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7430))
* `azurerm_hdinsight_ml_services_cluster` - fixes for node and instance count validation ([#7430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7430))
* `azurerm_hdinsight_rserver_cluster` - fixes for node and instance count validation ([#7430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7430))
* `azurerm_hdinsight_spark_cluster` - fixes for node and instance count validation ([#7430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7430))
* `azurerm_hdinsight_storm_cluster` - fixes for node and instance count validation ([#7430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7430))
* `azurerm_monitor_autoscale_settings` - support for setting `time_aggregation` to `Last` as per the documentation ([#7480](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7480))
* `azurerm_postgresql_server` - can now update the tier of `sku_name` by recreating the resource ([#7456](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7456))
* `azurerm_network_interface_security_group_association` - is now considered delete whtn the  network interfact is notfound ([#7459](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7459))
* `azurerm_role_definition` - terraform import now sets scope to prevent a force recreate ([#7424](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7424))
* `azurerm_storage_account_network_rules` - corretly clear `ip_rules`, `virtual_network_subnet_ids` when set to `[]` ([#7385](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7385))

## 2.15.0 (June 19, 2020)

UPGRADE NOTES:

* `azurerm_orchestrated_virtual_machine_scale_set` - the `single_placement_group` property is now required to be `false` by the service team in the `2019-12-01` compute API ([#7188](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7188))

DEPENDENCIES

* updating to `v43.1.0` of `github.com/Azure/azure-sdk-for-go` ([#7188](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7188))
* upgrading `kusto` to`2019-12-01` ([#7101](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7101))
* upgrading `kusto` to`2020-02-15` ([#6838](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6838))

FEATURES

* **New Data Source:** `azurerm_data_share_dataset_blob_storage` ([#7107](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7107))
* **New Resource:** `azurerm_data_factory_integration_runtime_self_hosted` ([#6535](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6535))
* **New Resource:** `azurerm_data_share_dataset_blob_storage` ([#7107](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7107))
* **New Resource:** `azurerm_eventhub_cluster` ([#7306](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7306))
* **New Resource:** `azurerm_maintenance_assignment_dedicated_host` ([#6713](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6713))
* **New Resource:** `azurerm_maintenance_assignment_virtual_machine` ([#6713](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6713))

IMPROVEMENTS:

* Data Source: `azurerm_management_group` - support lookup via `display_name` ([#6845](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6845))
* `azurerm_api_management` - support for the `developer_portal_url` property ([#7263](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7263))
* `azurerm_app_service` - support for `scm_ip_restriction` ([#6955](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6955))
* `azurerm_app_service_certificate `- support for the `hosting_environment_profile_id` propety ([#7087](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7087))
* `azurerm_app_service_environment` - support for the `user_whitelisted_ip_ranges` property ([#7324](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7324))
* `azurerm_kusto_cluster` - Support for `enable_purge` ([#7375](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7375))
* `azurerm_kusto_cluster` - Support for extended Kusto Cluster SKUs ([#7372](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7372))
* `azurerm_policy_assignment` - added support for `enforcement_mode`  ([#7331](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7331))
* `azurerm_private_endpoint` - support for the `private_dns_zone_group`, `private_dns_zone_configs`, and `custom_dns_configs` blocks ([#7246](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7246))
* `azurerm_storage_share_directory ` - `name` can now contain one nested directory ([#7382](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7382))

BUG FIXES:

* `azurerm_api_management_api` - correctly wait for future on create/update ([#7273](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7273))
* `azurerm_bot_connection` - adding a runtime check for the available service providers in the Azure Region being used ([#7279](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7279))
* `azurerm_healthcare_service` - the `access_policy_object_ids` property is now optional ([#7296](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7296))
* `azurerm_hdinsight_cluster` - deprecating the `min_instance_count` property ([#7272](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7272))
* `azurerm_network_watcher_flow_log` - propertly disable the flowlog on destroy ([#7154](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7154))

## 2.14.0 (June 11, 2020)

UPGRADE NOTES:

* `azurerm_kubernetes_cluster` - the Azure Policy add-on now only supports `v2` (as per the Azure API) ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))

DEPENDENCIES: 

* `containerservice` - updating to `2020-03-01` ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `policy` - updating to `2019-09-01` ([#7211](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7211))

FEATURES:

* **New Data Source:** `azurerm_blueprint_definition` ([#6930](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6930))
* **New Data Source:** `azurerm_blueprint_published_version` ([#6930](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6930))
* **New Data Source:** `azurerm_key_vault_certificate` ([#7285](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7285))
* **New Data Source:** `azurerm_kubernetes_cluster_node_pool` ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* **New Resource:** `azurerm_blueprint_assignment` ([#6930](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6930))
* **New Resource:** `azurerm_data_factory_linked_service_key_vault` ([#6971](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6971))
* **New Resource:** `azurerm_iot_time_series_insights_access_policy` ([#7202](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7202))
* **New Resource:** `azurerm_iot_time_series_insights_reference_data_set` ([#7112](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7112))
* **New Resource:** `azurerm_app_service_hybrid_connection` ([#7224](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7224))

ENHANCEMENTS:

* Data Source: `azurerm_kubernetes_cluster` - exposing the `version` of the Azure Policy add-on ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* Data Source: `azurerm_kubernetes_cluster` - exposing the `orchestrator_version` being used for each Node Pool ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* Data Source: `azurerm_kubernetes_cluster` - exposing the `disk_encryption_set_id` field ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_api_management_api` - ensuring `wsdl_selector` is populated when `content_format` is `wsdl` ([#7076](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7076))
* `azurerm_cosmosdb_account` modifying `geo_location` no longer triggers a recreation of the resource ([#7217](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7217))
* `azurerm_eventgrid_event_subscription` - support for `azure_function_endpoint` ([#7182](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7182))
* `azurerm_eventgrid_event_subscription` - exposing `base_url`, `max_events_per_batch`, `preferred_batch_size_in_kilobytes`, `active_directory_tenant_id` and `active_directory_app_id_or_uri` in the `webhook_endpoint` block ([#7207](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7207))
* `azurerm_kubernetes_cluster` - support for configuring/updating the version of Kubernetes used in the Default Node Pool ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster` - support for Azure Active Directory (Managed) Integration v2 ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster` - support for using a Disk Encryption Set ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster` - support for configuring the Auto-Scale Profile ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster` - support for configuring `outbound_ports_allocated` and `idle_timeout_in_minutes` within the `load_balancer_profile` block ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster` - support for the Uptime SLA / Paid SKU ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster` - exposing the `private_fqdn` of the cluster ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster_node_pool` - support for configuring/updating the version of Kubernetes ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster_node_pool` - support for Spot Node Pools ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster_node_pool` - support for System & User Node Pools ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_web_application_firewall_policy` - Add support for `GeoMatch` operator in request filter ([#7181](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7181))

BUG FIXES:

* Data Source: `azurerm_kubernetes_cluster` - fixing an issue where some read-only fields were unintentionally marked as user-configurable ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_application_gateway` - support for specifying the ID of a Key Vault Secret without a version ([#7095](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7095))
* `azurerm_bot_channel_ms_teams` - only sending `calling_web_hook` when it's got a value ([#7294](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7294))
* `azurerm_eventhub_namespace_authorization_rule` - handling the Resource ID changing on Azure's side from `authorizationRules` to `AuthorizationRules` ([#7248](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7248))
* `azurerm_eventgrid_event_subscription` - fixing a crash when `subject_filter` was omitted ([#7222](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7222))
* `azurerm_function_app` - fix app_settings when using linux consumption plan ([#7230](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7230))
* `azurerm_linux_virtual_machine_scale_set` - adding validation for the `max_bid_price` field ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster` - the Azure Policy add-on is not supported in Azure China and no longer sent ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster` - the Azure Policy add-on is not supported in Azure US Government and no longer sent ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster` - the Kubernetes Dashboard add-on is not supported in Azure US Government and no longer sent ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster` - searching for a system node pool when importing the `default_node_pool` ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_kubernetes_cluster_node_pool` - changes to the `node_taints` field now force a new resource, matching the updated API behaviour ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))
* `azurerm_management_group` - using the Subscription ID rather than Subscription Resource ID when detaching Subscriptions from Management Groups during deletion ([#7216](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7216))
* `azurerm_windows_virtual_machine_scale_set` - adding validation for the `max_bid_price` field ([#7233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7233))

## 2.13.0 (June 04, 2020)

FEATURES:

* **New Data Source**: `azurerm_logic_app_integration_account` ([#7099](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7099))
* **New Data Source:** `azurerm_virtual_machine_scale_set` ([#7141](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7141))
* **New Resource**: `azurerm_logic_app_integration_account` ([#7099](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7099))
* **New Resource**: `azurerm_monitor_action_rule_action_group` ([#6563](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6563))
* **New Resource**: `azurerm_monitor_action_rule_suppression` ([#6563](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6563))

IMPROVEMENTS:

* `azurerm_data_factory_pipeline` - Support for `activities` ([#6224](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6224))
* `azurerm_eventgrid_event_subscription` - support for advanced filtering ([#6861](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6861))
* `azurerm_signalr_service` - support for `EnableMessagingLogs` feature ([#7094](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7094))

BUG FIXES:

* `azurerm_app_service` - default priority now set on ip restricitons when not explicitly specified ([#7059](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7059))
* `azurerm_app_service` - App Services check correct scope for name availability in ASE ([#7157](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7157))
* `azurerm_cdn_endpoint` - `origin_host_header` can now be set to empty ([#7164](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7164))
* `azurerm_cosmosdb_account` - workaround for CheckNameExists 500 response code bug ([#7189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7189))
* `azurerm_eventhub_authorization_rule` - Fix intermittent 404 errors ([#7122](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7122))
* `azurerm_eventgrid_event_subscription` - fixing an error when setting the `hybrid_connection_endpoint` block ([#7203](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7203))
* `azurerm_function_app` - correctly set `Kind` when `os_type` is `linux` ([#7140](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7140))
* `azurerm_key_vault_certificate` - always setting the `certificate_data` and `thumbprint` fields ([#7204](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7204))
* `azurerm_role_assignment` - support for Preview role assignments ([#7205](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7205))
* `azurerm_virtual_network_gateway` - `vpn_client_protocols` is now also computed to prevent permanent diffs ([#7168](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7168))

## 2.12.0 (May 28, 2020)

FEATURES:

* **New Data Source:** `azurerm_advisor_recommendations` ([#6867](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6867))
* **New Resource:** `azurerm_dev_test_global_shutdown_schedule` ([#5536](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5536))
* **New Resource:** `azurerm_nat_gateway_public_ip_association` ([#6450](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6450))

IMPROVEMENTS:

* Data Source: `azurerm_kubernetes_cluster` - exposing the `oms_agent_identity` block within the `addon_profile` block ([#7056](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7056))
* Data Source: `azurerm_kubernetes_cluster` - exposing the `identity` and `kubelet_identity` properties ([#6527](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6527))
* `azurerm_batch_pool` - support the `container_image_names` property ([#6689](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6689))
* `azurerm_eventgrid_event_subscription` - support for the `expiration_time_utc`, `service_bus_topic_endpoint`, and `service_bus_queue_endpoint`, property ([#6860](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6860))
* `azurerm_eventgrid_event_subscription` - the `eventhub_endpoint` was deprecated in favour of the `eventhub_endpoint_id` property ([#6860](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6860))
* `azurerm_eventgrid_event_subscription` - the `hybrid_connection_endpoint` was deprecated in favour of the `hybrid_connection_endpoint_id` property ([#6860](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6860))
* `azurerm_eventgrid_topic` - support for `input_schema`, `input_mapping_fields`, and `input_mapping_default_values` ([#6858](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6858))
* `azurerm_kubernetes_cluster` - exposing the `oms_agent_identity` block within the `addon_profile` block ([#7056](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7056))
* `azurerm_logic_app_action_http` - support for the `run_after` property ([#7079](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7079))
* `azurerm_storage_account` - support `RAGZRS` and `GZRS` for the `account_replication_type` property ([#7080](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7080))

BUG FIXES:

* `azurerm_api_management_api_version_set` - handling changes to the Azure Resource ID ([#7071](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7071))
* `azurerm_key_vault_certificate` - fixing a bug when using externally-signed certificates (using the `Unknown` issuer) where polling would continue indefinitely ([#6979](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6979))
* `azurerm_linux_virtual_machine` - correctly validating the rsa ssh `public_key` properties length ([#7061](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7061))
* `azurerm_linux_virtual_machine` - allow setting `virtual_machine_scale_set_id` in non-zonal deployment ([#7057](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7057))
* `azurerm_servicebus_topic` - support for numbers in the `name` field ([#7027](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7027))
* `azurerm_shared_image_version` - `target_region.x.storage_account_type` is now defaulted and multiple `target_region`s can be added/removed ([#6940](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6940))
* `azurerm_sql_virtual_network_rule` - updating the validation for the `name` field ([#6968](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6968))
* `azurerm_windows_virtual_machine` - allow setting `virtual_machine_scale_set_id` in non-zonal deployment ([#7057](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7057))
* `azurerm_windows_virtual_machine` - correctly validating the rsa ssh `public_key` properties length ([#7061](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7061))

## 2.11.0 (May 21, 2020)

DEPENDENCIES:

* updating `github.com/Azure/azure-sdk-for-go` to `v42.1.0` ([#6725](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6725))
* updating `network` to `2020-03-01` ([#6727](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6727))

FEATURES:

* **Opt-In/Experimental Enhanced Validation for Locations:** This allows validating that the `location` field being specified is a valid Azure Region within the Azure Environment being used - which can be caught via `terraform plan` rather than `terraform apply`. This can be enabled by setting the Environment Variable `ARM_PROVIDER_ENHANCED_VALIDATION` to `true` and will be enabled by default in a future release of the AzureRM Provider ([#6927](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6927))
* **Data Source:** `azurerm_data_share` ([#6789](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6789))
* **New Resource:** `azurerm_data_share` ([#6789](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6789))
* **New Resource:** `azurerm_iot_time_series_insights_standard_environment` ([#7012](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7012))
* **New Resource:** `azurerm_orchestrated_virtual_machine_scale_set` ([#6626](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6626))

IMPROVEMENTS:

* Data Source: `azurerm_platform_image` - support for `version` filter ([#6948](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6948))
* `azurerm_api_management_api_version_set` - updating the validation for the `name` field ([#6947](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6947))
* `azurerm_app_service` - the `ip_restriction` block now supports the `action` property ([#6967](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6967))
* `azurerm_databricks_workspace` - exposing `workspace_id` and `workspace_url` ([#6973](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6973))
* `azurerm_netapp_volume` - support the `mount_ip_addresses` property ([#5526](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5526))
* `azurerm_redis_cache` - support new maxmemory policies `allkeys-lfu` & `volatile-lfu` ([#7031](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7031))
* `azurerm_storage_account` - allowing the value `PATCH` for `allowed_methods` within the `cors_rule` block within the `blob_properties` block ([#6964](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6964))

BUG FIXES:

* Data Source: `azurerm_api_management_group` - raising an error when the Group cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_image` - raising an error when the Image cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_data_lake_store` - raising an error when Data Lake Store cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_data_share_account` - raising an error when Data Share Account cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_hdinsight_cluster` - raising an error when the HDInsight Cluster cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_healthcare_service` - raising an error when the HealthCare Service cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_healthcare_service` - ensuring all blocks are set in the response ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_firewall` - raising an error when the Firewall cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_maintenance_configuration` - raising an error when the Maintenance Configuration cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_private_endpoint_connection` - raising an error when the Private Endpoint Connection cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_resources` - does not return all matched resources sometimes ([#7036](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7036))
* Data Source: `azurerm_shared_image_version` - raising an error when the Image Version cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_shared_image_versions` - raising an error when Image Versions cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* Data Source: `azurerm_user_assigned_identity` - raising an error when the User Assigned Identity cannot be found ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* `azurerm_api_management_subscription` - fix the export of `primary_key` and `secondary_key` ([#6938](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6938))
* `azurerm_eventgrid_event_subscription` - correctly parsing the ID ([#6958](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6958))
* `azurerm_healthcare_service` - ensuring all blocks are set in the response ([#7024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7024))
* `azurerm_linux_virtual_machine` - allowing name to end with a capital letter ([#7023](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7023))
* `azurerm_linux_virtual_machine_scale_set` - allowing name to end with a capital ([#7023](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7023))
* `azurerm_management_group` - workaround for 403 bug in service response ([#6668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6668))
* `azurerm_postgresql_server` - do not attempt to get the threat protection when the `sku` is `basic` ([#7015](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7015))
* `azurerm_windows_virtual_machine` - allowing name to end with a capital ([#7023](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7023))
* `azurerm_windows_virtual_machine_scale_set` - allowing name to end with a capital ([#7023](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7023))

## 2.10.0 (May 14, 2020)

DEPENDENCIES: 

* updating `eventgrid` to `2020-04-01-preview` ([#6837](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6837))
* updating `iothub` to `2019-03-22-preview` ([#6875](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6875))

FEATURES:

* **New Data Source:** `azurerm_eventhub` ([#6841](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6841))
* **New Resource:** `azurerm_eventgrid_domain_topic` ([#6859](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6859))

IMPROVEMENTS:

* All Data Sources: adding validation for the `resource_group_name` field to not be empty where it's Required ([#6864](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6864))
* Data Source: `azurerm_virtual_machine` - export `identity` attribute ([#6826](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6826))
* `azurerm_api_management` - support for configuring the Developer Portal ([#6724](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6724))
* `azurerm_api_management` - support for user assigned managed identities ([#6783](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6783))
* `azurerm_api_management` - support `key_vault_id` that do not have a version ([#6723](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6723))
* `azurerm_api_management_diagnostic` - support required property `api_management_logger_id` ([#6682](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6682))
* `azurerm_application_gateway` - support for WAF policies ([#6105](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6105))
* `azurerm_app_service_environment` - support specifying explicit resource group ([#6821](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6821))
* `azurerm_express_route_circuit` - de-provision and re-provision circuit when changing the bandwidth reduction ([#6601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6601))
* `azurerm_frontdoor` - expose the `header_frontdoor_id` attribute ([#6916](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6916))
* `azurerm_log_analytics_workspace` - add support for `rentention_in_days` for Free Tier ([#6844](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6844))
* `azurerm_mariadb_server` - support for the `create_mode` property allowing the creation of replicas, point in time restores, and geo restors ([#6865](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6865))
* `azurerm_mariadb_server` - support for the `public_network_access_enabled` property ([#6865](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6865))
* `azurerm_mariadb_server` - all properties in the `storage_profile` block have been moved to the top level ([#6865](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6865))
* `azurerm_mariadb_server` - the following properties were renamed and changed to a boolean type: `ssl_enforcement` to `ssl_enforcement_enabled`, `geo_redundant_backup` to `geo_redundant_backup_enabled`, and `auto_grow` 
* `azurerm_mysql_server` - support for the `create_mode` property allowing the creation of replicas, point in time restores, and geo restors ([#6833](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6833))
* `azurerm_mysql_server` - support for the `public_network_access_enabled` property ([#6833](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6833))
* `azurerm_mysql_server` - all properties in the `storage_profile` block have been moved to the top level ([#6833](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6833))
* `azurerm_mysql_server` - the following properties were renamed and changed to a boolean type: `ssl_enforcement` to `ssl_enforcement_enabled`, `geo_redundant_backup` to `geo_redundant_backup_enabled`, and `auto_grow` to `auto_grow_enabled` ([#6833](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6833))
* `azurerm_mssql_server`  - add support for the `azuread_administrator` property ([#6822](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6822))
* `azurerm_postgres_server` - support for the `threat_detection_policy` property ([#6721](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6721))
* `azurerm_storage_account` - enable migration of `account_kind` from `Storage` to `StorageV2` ([#6580](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6580))
* `azurerm_windows_virtual_machine` - the `os_disk.disk_encryption_set_id` can now be updated ([#6846](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6846))

BUG FIXES:

* Data Source: `azurerm_automation_account` - using the ID of the Automation Account, rather than the ID of the Automation Account's Registration Info ([#6848](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6848))
* Data Source: `azurerm_security_group` - fixing crash where id is nil ([#6910](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6910))
* Data Source: `azurerm_mysql_server` - remove `administrator_login_password` property as it is not returned from the api ([#6865](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6865))
* `azurerm_api_management` - fixing a crash when `policy` is nil ([#6862](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6862))
* `azurerm_api_management` - only sending the `hostname_configuration` properties if they are not empty ([#6850](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6850))
* `azurerm_api_management_diagnostic` - can now be provision again by supporting `api_management_logger_id` ([#6682](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6682))
* `azurerm_api_management_named_value` - fix the non empty plan when `secret` is true ([#6834](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6834))
* `azurerm_application_insights` - `retention_in_days` defaults to 90 ([#6851](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6851))
* `azurerm_data_factory_trigger_schedule` - setting the `type` required for Pipeline References ([#6871](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6871))
* `azurerm_kubernetes_cluster` - fixes the `InvalidLoadbalancerProfile` error ([#6534](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6534))
* `azurerm_linux_virtual_machine_scale_set` - support for updating the `do_not_run_extensions_on_overprovisioned_machines` property ([#6917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6917))
* `azurerm_monitor_diagnostic_setting` - fix possible crash with `retention_policy` ([#6911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6911))
* `azurerm_mariadb_server` - the `storage_mb` property is now optional when `auto_grow` is enabled ([#6865](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6865))
* `azurerm_mysql_server` - the `storage_mb` property is now optional when `auto_grow` is enabled ([#6833](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6833))
* `azurerm_role_assignment` - added evential consistency check to assignment creation ([#6925](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6925))
* `azurerm_windows_virtual_machine_scale_set` - support for updating the `do_not_run_extensions_on_overprovisioned_machines` property ([#6917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6917))

## 2.9.0 (May 07, 2020)

FEATURES:

* **New Data Source:** `azurerm_data_share_account` ([#6575](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6575))
* **New Resource:** `azurerm_data_share_account` ([#6575](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6575))
* **New Resource:** `azurerm_function_app_slot` ([#6435](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6435))
* **New Resource:** `azurerm_sentinel_alert_rule_scheduled` ([#6650](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6650))

IMPROVEMENTS:

* Data Source: `azurerm_eventhub_authorization_rule` - support for the `primary_connection_string_alias` an `secondary_connection_string_alias` propeties  ([#6708](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6708))
* Data Source: `azurerm_eventhub_namespace_authorization_rule` - support for the `primary_connection_string_alias` an `secondary_connection_string_alias` propeties  ([#6708](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6708))
* Data Source: `azurerm_eventhub_namespace` - support for the `default_primary_connection_string_alias` an `_defaultsecondary_connection_string_alias` propeties  ([#6708](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6708))
* `azurerm_analysis_services_server` - support updating when the Server is paused ([#6786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6786))
* `azurerm_app_service` - support for health_check_path preview feature added ([#6661](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6661))
* `azurerm_app_service` - support for `name` and `priority` on `ip_restrictions` ([#6705](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6705))
* `azurerm_application_gateway` - support for SSL Certificates without passwords ([#6742](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6742))
* `azurerm_eventhub_authorization_rule` - support for the `primary_connection_string_alias` an `secondary_connection_string_alias` propeties  ([#6708](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6708))
* `azurerm_eventhub_namespace_authorization_rule` - support for the `primary_connection_string_alias` an `secondary_connection_string_alias` propeties  ([#6708](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6708))
* `azurerm_eventhub_namespace` - support for the `default_primary_connection_string_alias` an `_defaultsecondary_connection_string_alias` propeties  ([#6708](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6708))
* `azurerm_hdinsight_hadoop_cluster` - support for metastores on cluster creation ([#6145](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6145))
* `azurerm_key_vault_certificate` - support for recovering a soft-deleted certificate if the `features` flag `recover_soft_deleted_key_vaults` is set to `true` ([#6716](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6716))
* `azurerm_key_vault_key` - support for recovering a soft-deleted key if the `features` flag `recover_soft_deleted_key_vaults` is set to `true` ([#6716](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6716))
* `azurerm_key_vault_secret` - support for recovering a soft-deleted secret if the `features` flag `recover_soft_deleted_key_vaults` is set to `true` ([#6716](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6716))
* `azurerm_linux_virtual_machine_scale_set` - support for configuring `create_mode` for data disks ([#6744](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6744))
* `azurerm_monitor_diagnostic_setting` - `log_analytics_destination_type` supports `AzureDiagnostics` ([#6769](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6769))
* `azurerm_windows_virtual_machine_scale_set` - support for configuring `create_mode` for data disks ([#6744](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6744))

BUG FIXES:

* provider: raising an error when the environment is set to `AZURESTACKCLOUD` ([#6817](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6817))
* `azurerm_analysis_services_server` - ip restriction name field no longer case sensitive ([#6774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6774))
* `azurerm_automation_runbook` - the `publish_content_link` property is now optional ([#6813](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6813))
* `azurerm_eventhub_namespace_authorization_rule` - lock to prevent multiple resources won't clash ([#6701](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6701))
* `azurerm_network_interface` - changes to dns servers no longer use incremental update ([#6624](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6624))
* `azurerm_policy_assignment` - allow polices with scopes without `subscription/<id>` (built-in policies) ([#6792](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6792))
* `azurerm_policy_definition` - changes to the dynamic fields (`createdBy`, `createdOn`, `updatedBy`, `updatedOn`) keys in the `metadata` field are excluded from diff's ([#6734](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6734))
* `azurerm_redis_cache` - ensure `rdb_storage_connection_string` is set when `rdb_backup_enabled` is enabled ([#6819](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6819))
* `azurerm_site_recovery_network_mapping` - handling an API Error when checking for the presence of an existing Network Mapping ([#6747](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6747))

## 2.8.0 (April 30, 2020)

FEATURES:

* **New Data Source:** `azurerm_sentinel_alert_rule_ms_security_incident` ([#6606](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6606))
* **New Data Source:** `azurerm_shared_image_versions` ([#6700](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6700))
* **New Resource:** `azurerm_managed_application` ([#6386](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6386))
* **New Resource:** `azurerm_mssql_server` ([#6677](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6677))
* **New Resource:** `azurerm_sentinel_alert_rule_ms_security_incident` ([#6606](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6606))

IMPROVEMENTS:

* `azurerm_api_management` - `sku_name` supports the `Consumption` value for `sku` ([#6602](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6602))
* `azurerm_api_management_api` - support for openapi v3 content formats ([#6618](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6618))
* `azurerm_application_gateway` - support `host_names` property ([#6630](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6630))
* `azurerm_express_route_circuit_peering` - support for the `customer_asn` and `routing_registry_name` propeties ([#6596](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6596))
* `azurerm_frontdoor` - Add support for `backend_pools_send_receive_timeout_seconds` ([#6604](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6604))
* `azurerm_mssql_server` -support the `public_network_access_enabled` property ([#6678](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6678))
* `azurerm_mssql_database` - support for the `extended_auditing_policy` block ([#6402](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6402))
* `azurerm_mssql_elasticpool` - support `license_type` ([#6631](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6631))
* `azurerm_subnet`: Support for multiple prefixes with `address_prefixes` ([#6493](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6493))
* `data.azurerm_shared_image_version` - `name` supports `latest` and `recent`  ([#6707](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6707))

BUG FIXES:

* `azurerm_key_vault` - can now be created without subscription level permissions ([#6260](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6260))
* `azurerm_linux_virtual_machine` - fix validation for `name` to allow full length resource names ([#6639](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6639))
* `azurerm_linux_virtual_machine_scale_set` - fix validation for `name` to allow full length resource names ([#6639](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6639))
* `azurerm_monitor_diagnostic_setting` - make `retention_policy` and `retention_policy` optional ([#6603](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6603))
* `azurerm_redis_cache` - correctly build connection strings when SSL is disabled ([#6635](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6635))
* `azurerm_sql_database` - prevent extended auditing policy for secondary databases ([#6402](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6402))
* `azurerm_web_application_firewall_policy` - support for the `managed_rules` property which is required by the new API version ([#6126](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6126))
* `azurerm_windows_virtual_machine` - fix validation for `name` to allow full length resource names ([#6639](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6639))
* `azurerm_windows_virtual_machine_scale_set` - fix validation for `name` to allow full length resource names ([#6639](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6639))
* `azurerm_virtual_network_gateway_connection` - `shared_key` is now optional when `type` is `IPSec` ([#6565](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6565))

## 2.7.0 (April 23, 2020)

FEATURES:

* **New Data Source:** `azurerm_private_dns_zone` ([#6512](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6512))
* **New Resource:** `azurerm_maintenance_configuration` ([#6038](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6038))
* **New Resource:** `azurerm_servicebus_namespace_network_rule_set` ([#6379](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6379))
* **New Resource:** `azurerm_spring_cloud_app` ([#6384](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6384))

DEPENDENCIES:

* updating `apimanagement` to `2019-12-01` ([#6479](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6479))
* updating the fork of `github.com/Azure/go-autorest` ([#6509](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6509))

IMPROVEMENTS:

* Data Source: `app_service_environment` - export the `location` property ([#6538](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6538))
* Data Source: `azurerm_notification_hub_namespace` - export `tags` ([#6578](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6578))
* `azurerm_api_management` - support for virtual network integrations ([#5769](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5769))
* `azurerm_cosmosdb_mongo_collection` - support for the `index` and `system_index` properties ([#6426](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6426))
* `azurerm_function_app` - added `storage_account_id` and `storage_account_access_key` ([#6304](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6304))
* `azurerm_kubernetes_cluster` - deprecating `private_link_enabled` in favour of `private_cluster_enabled ` ([#6431](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6431))
* `azurerm_mysql_server` - support for the `public_network_access_enabled` property ([#6590](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6590))
* `azurerm_notification_hub` - support for `tags` ([#6578](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6578))
* `azurerm_notification_hub_namespace` - support for `tags` ([#6578](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6578))
* `azurerm_postgres_server` - support for the `create_mode` property allowing replicas, point in time restores, and geo restores to be created ([#6459](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6459))
* `azurerm_postgres_server` - support for the `infrastructure_encryption_enabled`, `public_network_access_enabled`, and `ssl_minimal_tls_version_enforced` properties ([#6459](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6459))
* `azurerm_postgres_server` - all properties in the `storage_profile` block have been moved to the top level ([#6459](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6459))
* `azurerm_postgres_server` - the following properties were renamed and changed to a boolean type: `ssl_enforcement` to `ssl_enforcement_enabled`, `geo_redundant_backup` to `geo_redundant_backup_enabled`, and `auto_grow` to `auto_grow_enabled` ([#6459](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6459))
* `azurerm_private_endpoint` - Add support for `tags` ([#6574](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6574))
* `azurerm_shared_image` - support `hyper_v_generation` property ([#6511](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6511))
* `azurerm_linux_virtual_machine_scale_set` - support for the `automatic_instance_repair` property ([#6346](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6346))
* `azurerm_windows_virtual_machine_scale_set` - support for the `automatic_instance_repair` property ([#6346](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6346))

BUG FIXES:

* Data Source: `azurerm_private_link_service` - fixing a crash when parsing the response ([#6504](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6504))
* `azurerm_application_gateway` - prevent panic by disallowing empty values for `backend_address_pool.#.fqdns` ([#6549](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6549))
* `azurerm_application_gateway` - block reordering without changes no longer causes update ([#6476](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6476))
* `azurerm_cdn_endpoint` - `origin_host_header` is now required ([#6550](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6550))
* `azurerm_cdn_endpoint` - setting the `request_header_condition` block ([#6541](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6541))
* `azurerm_iothub_dps` - fix crash when path isn't cased correctly ([#6570](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6570))
* `azurerm_linux_virtual_machine_scale_set` - fixes crash with `boot_diagnositics` ([#6569](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6569))
* `azurerm_policy_assignment` - allow scopes that don't start with `subscription/<id>` ([#6576](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6576))
* `azurerm_postgres_server` - the `storage_mb` property is now optional when `auto_grow` is enabled ([#6459](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6459))
* `azurerm_public_ip_prefix` - update `prefix_length` validation to accept all valid IPv4 address ranges ([#6589](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6589))
* `azurerm_route` - add validation to the `name` and `route_table_name`propeties ([#6055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6055))
* `azurerm_virtual_network_gateway` - per api requirements, `public_ip_address_id` is required ([#6548](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6548))

## 2.6.0 (April 16, 2020)

FEATURES:

* **New Data Source:** `azurerm_policy_set_definition` ([#6305](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6305))

DEPENDENCIES:

* updating `github.com/Azure/azure-sdk-for-go` to `v41.2.0` ([#6419](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6419))

IMPROVEMENTS:

* Data Source: `azurerm_policy_definition` - can now lookup with `name` ([#6275](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6275))
* Data Source: `azurerm_policy_definition` - the field `management_group_id` has been deprecated and renamed to `management_group_name` ([#6275](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6275))
* `azurerm_application_insights` - support for the `disable_ip_masking` property ([#6354](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6354))
* `azurerm_cdn_endpoint` - support for configuring `delivery_rule` ([#6163](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6163))
* `azurerm_cdn_endpoint` - support for configuring `global_delivery_rule` ([#6163](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6163))
* `azurerm_function_app` - support for the `pre_warmed_instance_count` property ([#6333](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6333))
* `azurerm_hdinsight_hadoop_cluster` - support for the `tls_min_version` property ([#6440](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6440))
* `azurerm_hdinsight_hbase_cluster` - support for the `tls_min_version` property ([#6440](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6440))
* `azurerm_hdinsight_interactive_query_cluster` - support for the `tls_min_version` property ([#6440](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6440))
* `azurerm_hdinsight_kafka_cluster` - support for the `tls_min_version` property ([#6440](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6440))
* `azurerm_hdinsight_ml_services_cluster` - support for the `tls_min_version` property ([#6440](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6440))
* `azurerm_hdinsight_rserver_cluster` - support for the `tls_min_version` property ([#6440](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6440))
* `azurerm_hdinsight_spark_cluster` - support for the `tls_min_version` property ([#6440](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6440))
* `azurerm_hdinsight_storm_cluster` - support the `threat_detection_policy` property ([#6437](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6437))
* `azurerm_kubernetes_cluster` - exporting the `kubelet_identity` ([#6393](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6393))
* `azurerm_kubernetes_cluster` - support for updating the `managed_outbound_ip_count`, `outbound_ip_prefix_ids` and `outbound_ip_address_ids` fields within the `load_balancer_profile` block ([#5847](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5847))
* `azurerm_network_interface` - export the `internal_domain_name_suffix` property ([#6455](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6455))
* `azurerm_policy_definition` - the `management_group_id` has been deprecated and renamed to `management_group_name` ([#6275](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6275))
* `azurerm_sql_server` - support for the `connection_policy` property ([#6438](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6438))
* `azurerm_virtual_network` - export the `guid` attribute ([#6445](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6445))

BUG FIXES:

* Data Source: `azurerm_data_factory`- fixing a bug where the ID wasn't set ([#6492](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6492))
* Data Source: `azurerm_eventhub_namespace_authorization_rule` - ensuring the `id` field is set ([#6496](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6496))
* Data Source: `azurerm_mariadb_server` - ensuring the `id` field is set ([#6496](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6496))
* Data Source: `azurerm_network_ddos_protection_plan` - ensuring the `id` field is set ([#6496](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6496))
* `azurerm_function_app` - prevent a panic from the API returning an empty IP Security Restriction ([#6442](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6442))
* `azurerm_machine_learning_workspace` - the `Enterprise` sku will now properly work ([#6397](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6397))
* `azurerm_managed_disk`-  fixing a bug where the machine would be stopped regardless of whether it was currently shut down or not ([#4690](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4690))

## 2.5.0 (April 09, 2020)

BREAKING CHANGES:

* Azure Kubernetes Service
	* Due to a breaking change in the AKS API, the `azurerm_kubernetes_cluster` resource features a significant behavioural change where creating Mixed-Mode Authentication clusters (e.g. using a Service Principal with a Managed Identity) is no longer supported.
	* The AKS Team have confirmed that existing clusters will be updated by the Azure API to use only MSI when a change is made to the Cluster (but not the Node Pool). Whilst Terraform could perform this automatically some environments have restrictions on which tags can be added/removed - as such this operation will need to be performed out-of-band. Instead, upon detecting a Mixed-Mode Cluster which has not yet been updated - or upon detecting a former Mixed-Mode Cluster where the Terraform Configuration still contains a `service_principal` block - Terraform will output instructions on how to proceed.
	* `azurerm_kubernetes_cluster_node_pool` - clusters with auto-scale disabled must ensure that `min_count` and `max_count` are set to `null` (or omitted) rather than `0` (since 0 isn't a valid value for these fields).

NOTES:

* There's currently a bug in the Azure Kubernetes Service (AKS) API where the Tags on Node Pools are returned in the incorrect case - [this bug is being tracked in this issue](https://github.com/Azure/azure-rest-api-specs/issues/8952). This affects the `tags` field within the `default_node_pool` block for `azurerm_kubernetes_clusters` and the `tags` field for the `azurerm_kubernetes_cluster_node_pool` resource.

IMPROVEMENTS:

* dependencies: updating to use version `2020-02-01` of the Containers API ([#6095](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6095))
* **New Resource:** `azurerm_private_dns_txt_record` ([#6309](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6309))
* `azurerm_kubernetes_cluster` - making the `service_principal` block optional - so it's now possible to create MSI-only clusters ([#6095](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6095))
* `azurerm_kubernetes_cluster` - making the `windows_profile` block computed as Windows credentials are now generated by Azure if unspecified ([#6095](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6095))
* `azurerm_kubernetes_cluster` - support for `outbound_type` within the `network_profile` block ([#6120](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6120))
* `azurerm_linux_virtual_machine` - OS disk encryption settings can no be updated ([#6230](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6230))
* `azurerm_windows_virtual_machine` - OS disk encryption settings can no be updated ([#6230](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6230))

BUG FIXES:

* `azurerm_kubernetes_cluster` - requiring that `min_count` and `max_count` within the `default_node_pool` block are set to `null` rather than `0` when auto-scaling is disabled ([#6095](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6095))
* `azurerm_kubernetes_cluster` - ensuring that a value for `node_count` within the `default_node_pool` block is always passed to the API to match a requirement in the API ([#6095](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6095))
* `azurerm_kubernetes_cluster` - ensuring that `tags` are set into the state for the `default_node_pool` ([#6095](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6095))
* `azurerm_kubernetes_cluster` - conditionally sending the `aci_connector_linux` block for Azure China ([#6370](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6370))
* `azurerm_kubernetes_cluster` - conditionally sending the `http_application_routing` block for Azure China & Azure US Government ([#6370](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6370))
* `azurerm_kubernetes_cluster_node_pool` - requiring that `min_count` and `max_count` are set to `null` rather than `0` when auto-scaling is disabled ([#6095](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6095))
* `azurerm_linux_virtual_machine` - if the `priority` property on read is empty assume it to be `Regular` ([#6301](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6301))
* `azurerm_windows_virtual_machine` - if the `priority` property on read is empty assume it to be `Regular` ([#6301](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6301))

## 2.4.0 (April 02, 2020)

FEATURES:

* **New Data Source:** `azurerm_managed_application_definition` ([#6211](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6211))
* **New Resource:** `azurerm_hpc_cache_nfs_target` ([#6191](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6191))
* **New Resource:** `azurerm_log_analytics_datasource_windows_event ` ([#6321](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6321))
* **New Resource:** `azurerm_log_analytics_datasource_windows_performance_counter` ([#6274](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6274))
* **New Resource:** `azurerm_managed_application_definition` ([#6211](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6211))
* **New Resource:** `azurerm_spring_cloud_service` ([#4928](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4928))

IMPROVEMENTS:

* `azurerm_network_interface` - always send `enable_accelerated_networking` to the api ([#6289](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6289))
* `azurerm_management_group` - deprecated and rename the `group_id` property to `name` to better match what it represents ([#6276](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6276))

BUGS:

* `azurerm_application_gateway` - can now set `include_path` with `target_url` ([#6175](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6175))
* `azurerm_policy_set_definition` - mark `metadata` as computed ([#6266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6266))

## 2.3.0 (March 27, 2020)

FEATURES:

* **New Data Source:** `azurerm_mssql_database` ([#6083](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6083))
* **New Data source:** `azurerm_network_service_tags` ([#6229](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6229))
* **New Resource:** `azurerm_custom_resource_provider` ([#6234](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6234))
* **New Resource:** `azurerm_hpc_cache_blob_target` ([#6035](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6035))
* **New Resource:** `azurerm_machine_learning_workspace` ([#5696](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5696))
* **New Resource:** `azurerm_mssql_database` ([#6083](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6083))
* **New Resource:** `azurerm_mssql_virtual_machine` ([#5263](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5263))
* **New resource:** `azurerm_policy_remediation` ([#5746](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5746))

IMPROVEMENTS:

* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v40.3.0` ([#6134](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6134))
* dependencies: updating `github.com/terraform-providers/terraform-provider-azuread` to `v0.8.0` ([#6134](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6134))
* dependencies: updating `github.com/tombuildsstuff/giovanni` to `v0.10.0` ([#6169](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6169))
* all resources using the `location` field - adding validation to ensure this is not an empty string where this field is Required ([#6242](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6242))
* Data Source `azurerm_storage_container` - exposing the `resource_manager_id` field ([#6170](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6170))
* `azurerm_automation_schedule` - adding validation for the timezone field ([#5759](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5759))
* `azurerm_cognitive_account` - support for the `qna_runtime_endpoint` property ([#5778](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5778))
* `azurerm_hpc_cache` - exposing the `mount_addresses` field ([#6214](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6214))
* `azurerm_lb` - allow ipv6 addresses for the `private_ip_address` property ([#6125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6125))
* `azurerm_managed_disk` - the `disk_encryption_set_id` field is no longer ForceNew ([#6207](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6207))
* `azurerm_public_ip` - support for Dynamic IPv6 Addresses ([#6140](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6140))
* `azurerm_service_fabric_cluster` - support for the `client_certificate_common_name` property ([#6097](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6097))
* `azurerm_storage_container` - exposing the `resource_manager_id` field ([#6170](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6170))
* `azurerm_storage_share` - exposing the `resource_manager_id` field ([#6170](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6170))
* `azurerm_traffic_manager_profile` - support for the `custom_header` property ([#5923](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5923))

BUG FIXES:

* `azurerm_analysis_server` - switching the `ipv4_firewall_rule` block to a Set rather than a List to handle this being unordered ([#6179](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6179))
* `azurerm_linux_virtual_machine` - making the `custom_data` field sensitive ([#6225](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6225))
* `azurerm_linux_virtual_machine_scale_set` - making the `custom_data` field sensitive ([#6225](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6225))
* `azurerm_managed_disk`- only rebooting the attached Virtual Machine when changing the Disk Size, Disk Encryption Set ID or Storage Account Type ([#6162](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6162))
* `azurerm_netapp_volume` - allow up to `102400` MB for the `storage_quota_in_gb` property ([#6228](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6228))
* `azurerm_policy_definition` - fixing a bug when parsing the Management Group ID ([#5981](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5981))
* `azurerm_postgresql_server` - updating the validation for the `name` field ([#6064](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6064))
* `azurerm_sql_database` - use the correct base URI for the Extended Auditing Policies Client ([#6233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6233))
* `azurerm_storage_management_policy` - conditionally setting values within the `base_blob` block ([#6250](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6250))
* `azurerm_virtual_machine_data_disk_attachment` - detecting the disk attachment as gone when the VM is no longer available ([#6237](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6237))
* `azurerm_windows_virtual_machine` - making the `custom_data` field sensitive ([#6225](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6225))
* `azurerm_windows_virtual_machine_scale_set` - making the `custom_data` field sensitive ([#6225](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6225))

## 2.2.0 (March 18, 2020)

FEATURES:

* **New Data Source:** `azurerm_app_configuration` ([#6133](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6133))
* **New Data Source:** `azurerm_powerbi_embedded` ([#5152](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5152))
* **New Resource:** `azurerm_cost_management_export_resource_group` ([#6131](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6131))
* **New Resource:** `azurerm_powerbi_embedded` ([#5152](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5152))
* **New Resource:** `azurerm_virtual_hub_connection` ([#5951](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5951))

IMPROVEMENTS:

* Data Source: * `azurerm_logic_app_workflow`  - expose computed field: `endpoint_configuration` ([#5862](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5862))
* `azurerm_application_gateway` - support for key vault SSL certificate via the `key_value_secret_id` property ([#4366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4366))
* `azurerm_function_app` - support for configuring `daily_memory_time_quota` ([#6100](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6100))
* `azurerm_logic_app_workflow`  - expose computed field: `endpoint_configuration` ([#5862](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5862))
* `azurerm_linux_virtual_machine_scale_set` - support for `scale_in_policy` and `terminate_notification` ([#5391](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5391))
* `azurerm_sql_database` - support for the `extended_auditing_policy` property ([#5049](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5049))
* `azurerm_windows_virtual_machine_scale_set` - support for `scale_in_policy` and `terminate_notification` ([#5391](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5391))

BUG FIXES:

* Data Source: `azurerm_iothub_dps_shared_access_policy` - building the `primary_connection_string` and `secondary_connection_string` from the Service endpoint rather than the Devices endpoint ([#6108](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6108))
* `azurerm_function_app` - Add `WEBSITE_CONTENT` & `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` for premium plans ([#5761](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5761))
* `azurerm_iothub_dps_shared_access_policy` - building the `primary_connection_string` and `secondary_connection_string` from the Service endpoint rather than the Devices endpoint ([#6108](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6108))
* `azurerm_linux_virtual_machine` - updating the validation for `name` to allow periods ([#5966](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5966))
* `azurerm_linux_virtual_machine_scale_set` - updating the validation for `name` to allow periods ([#5966](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5966))
* `azurerm_storage_management_policy` - Fixed the use of single blob rule actions ([#5803](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5803))

## 2.1.0 (March 11, 2020)

NOTES:

The `azurerm_frontdoor` resource has introduced a breaking change due to the underlying service API which enforces `location` attributes must be set to 'Global' on all newly deployed Front Door services.

FEATURES:

* **New Data Source:** `azurerm_database_migration_project` ([#5993](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5993))
* **New Data Source:** `azurerm_database_migration_service` ([#5258](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5258))
* **New Data Source:** `azurerm_kusto_cluster` ([#5942](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5942))
* **New Data Source:** `azurerm_servicebus_topic_authorization_rule` ([#6017](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6017))
* **New Resource:** `azurerm_bot_channel_directline` ([#5445](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5445))
* **New Resource:** `azurerm_database_migration_project` ([#5993](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5993))
* **New Resource:** `azurerm_database_migration_service` ([#5258](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5258))
* **New Resource:** `azurerm_hpc_cache` ([#5528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5528))
* **New Resource:** `azurerm_iotcentral_application` ([#5446](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5446))
* **New Resource:** `azurerm_monitor_scheduled_query_rules_alert` ([#5053](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5053))
* **New Resource:** `azurerm_monitor_scheduled_query_rules_log` ([#5053](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5053))
* **New Resource:** `azurerm_spatial_anchors_account` ([#6011](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6011))

IMPROVEMENTS:

* batch: upgrading to API version `2019-08-01` ([#5967](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5967))
* containerservice: upgrading to API version `2019-11-01` ([#5531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5531))
* netapp: upgrading to API version `2019-10-01` ([#5531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5531))
* dependencies: temporarily switching to using a fork of `github.com/Azure/go-autorest` to workaround an issue in the storage authorizer ([#6050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6050))
* dependencies: updating `github.com/tombuildsstuff/giovanni` to `v0.9.0` ([#6050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6050))
* `azurerm_application_gateway` - support up to `125` for the `capacity` property with V2 SKU's ([#5906](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5906))
* `azurerm_automation_dsc_configuration` - support for the `tags` property ([#5827](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5827))
* `azurerm_batch_pool` - support for the `public_ips` property ([#5967](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5967))
* `azurerm_frontdoor` - exposed new attributes in `backend_pool_health_probe` block `enabled` and `probe_method` ([#5924](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5924))
* `azurerm_function_app` - Added `os_type` field to facilitate support of `linux` function apps ([#5839](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5839))
* `azurerm_kubernetes_cluster`: Support for the `node_labels` property ([#5531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5531))
* `azurerm_kubernetes_cluster`: Support for the `tags` property ([#5931](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5931))
* `azurerm_kubernetes_cluster_node_pool`: Support for the `node_labels` property ([#5531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5531))
* `azurerm_kubernetes_cluster_node_pool`: Support for the `tags` property ([#5931](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5931))
* `azurerm_kusto_cluster` - support for `enable_disk_encryption` and `enable_streaming_ingest` properties ([#5855](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5855))
* `azurerm_lb` - support for the `private_ip_address_version` property ([#5590](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5590))
* `azurerm_mariadb_server` - changing the `geo_redundant_backup` property now forces a new resource ([#5961](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5961))
* `azurerm_netapp_account` - support for the `tags` property ([#5995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5995))
* `azurerm_netapp_pool` - support for the `tags` property ([#5995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5995))
* `azurerm_netapp_snapshot` - support for the `tags` property ([#5995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5995))
* `azurerm_netapp_volume` - support for the `tags` property ([#5995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5995))
* `azurerm_netapp_volume` - support for the `protocol_types` property ([#5485](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5485))
* `azurerm_netapp_volume` - deprecated the `cifs_enabled`, `nfsv3_enabled`, and `nfsv4_enabled` properties in favour of `protocols_enabled` ([#5485](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5485))
* `azurerm_network_watcher_flow_log` - support for the traffic analysis `interval_in_minutes` property ([#5851](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5851))
* `azurerm_private_dns_a_record` - export the `fqdn` property ([#5949](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5949))
* `azurerm_private_dns_aaaa_record` - export the `fqdn` property ([#5949](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5949))
* `azurerm_private_dns_cname_record` - export the `fqdn` property ([#5949](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5949))
* `azurerm_private_dns_mx_record` - export the `fqdn` property ([#5949](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5949))
* `azurerm_private_dns_ptr_record` - export the `fqdn` property ([#5949](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5949))
* `azurerm_private_dns_srv_record` - export the `fqdn` property ([#5949](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5949))
* `azurerm_private_endpoint` - exposed `private_ip_address` as a computed attribute ([#5838](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5838))
* `azurerm_redis_cache` - support for the `primary_connection_string` and `secondary_connection_string` properties ([#5958](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5958))
* `azurerm_sql_server` - support for the `extended_auditing_policy` property ([#5036](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5036))
* `azurerm_storage_account` - support up to 50 tags ([#5934](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5934))
* `azurerm_virtual_wan` - support for the `type` property ([#5877](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5877))

BUG FIXES:

* `azurerm_app_service_plan` - no longer sends an empty `app_service_environment_id` property on update ([#5915](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5915))
* `azurerm_automation_schedule` - fix time validation ([#5876](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5876))
* `azurerm_batch_pool` - `frontend_port_range ` is now set correctly. ([#5941](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5941))
* `azurerm_dns_txt_record` - support records up to `1024` characters in length ([#5837](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5837))
* `azurerm_frontdoor` - fix the way `backend_pool_load_balancing`/`backend_pool_health_probe` ([#5924](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5924))
* `azurerm_frontdoor` - all new front door resources to be created in the `Global` location ([#6015](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6015))
* `azurerm_frontdoor_firewall_policy` - add validation for Frontdoor WAF Name Restrictions ([#5943](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5943))
* `azurerm_linux_virtual_machine_scale_set` - correct `source_image_id` validation ([#5901](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5901))
* `azurerm_netapp_volume` - support volmes uoto `100TB` in size ([#5485](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5485))
* `azurerm_search_service` - changing the properties `replica_count` & `partition_count` properties no longer force a new resource ([#5935](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5935))
* `azurerm_storage_account` - fixing a crash when an empty `static_website` block was specified ([#6050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6050))
* `azurerm_storage_account` - using SharedKey Authorization for reading/updating the Static Website when not using AzureAD authentication ([#6050](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6050))

## 2.0.0 (February 24, 2020)

NOTES:

* **Major Version:** Version 2.0 of the Azure Provider is a major version - some deprecated fields/resources have been removed - please [refer to the 2.0 upgrade guide for more information](https://www.terraform.io/docs/providers/azurerm/guides/2.0-upgrade-guide.html).
* **Provider Block:** The Azure Provider now requires that a `features` block is specified within the Provider block, which can be used to alter the behaviour of certain resources - [more information on the `features` block can be found in the documentation](https://www.terraform.io/docs/providers/azurerm/index.html#features).
* **Terraform 0.10/0.11:** Version 2.0 of the Azure Provider no longer supports Terraform 0.10 or 0.11 - you must upgrade to Terraform 0.12 to use version 2.0 of the Azure Provider.

FEATURES:

* **Custom Timeouts:** - all resources within the Azure Provider now allow configuring custom timeouts - please [see Terraform's Timeout documentation](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) and the documentation in each data source resource for more information.
* **Requires Import:** The Azure Provider now checks for the presence of an existing resource prior to creating it - which means that if you try and create a resource which already exists (without importing it) you'll be prompted to import this into the state.
* **New Data Source:** `azurerm_app_service_environment` ([#5508](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5508))
* **New Data Source:** `azurerm_eventhub_authorization_rule` ([#5805](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5805))
* **New Resource:** `azurerm_app_service_environment` ([#5508](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5508))
* **New Resource:** `azurerm_express_route_gateway` ([#5523](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5523))
* **New Resource:** `azurerm_linux_virtual_machine` ([#5705](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5705))
* **New Resource:** `azurerm_linux_virtual_machine_scale_set` ([#5705](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5705))
* **New Resource:** `azurerm_network_interface_security_group_association` ([#5784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5784))
* **New Resource:** `azurerm_storage_account_customer_managed_key` ([#5668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5668))
* **New Resource:** `azurerm_virtual_machine_scale_set_extension` ([#5705](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5705))
* **New Resource:** `azurerm_windows_virtual_machine` ([#5705](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5705))
* **New Resource:** `azurerm_windows_virtual_machine_scale_set` ([#5705](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5705))

BREAKING CHANGES:

* The Environment Variable `DISABLE_CORRELATION_REQUEST_ID` has been renamed to `ARM_DISABLE_CORRELATION_REQUEST_ID` to match the other Environment Variables
* The field `tags` is no longer `computed`
* Data Source: `azurerm_api_management` - removing the deprecated `sku` block ([#5725](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5725))
* Data Source: `azurerm_app_service` - removing the deprecated field `subnet_mask` from the `site_config` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* Data Source: `azurerm_app_service_plan` - the deprecated `properties` block has been removed since these properties have been moved to the top level ([#5717](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5717))
* Data Source: `azurerm_azuread_application` - This data source has been removed since it was deprecated ([#5748](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5748))
* Data Source: `azurerm_azuread_service_principal` - This data source has been removed since it was deprecated ([#5748](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5748))
* Data Source: `azurerm_builtin_role_definition` - the deprecated data source has been removed ([#5844](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5844))
* Data Source: `azurerm_dns_zone` - removing the deprecated `zone_type` field ([#5794](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5794))
* Data Source: `azurerm_dns_zone` - removing the deprecated `registration_virtual_network_ids` field ([#5794](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5794))
* Data Source: `azurerm_dns_zone` - removing the deprecated `resolution_virtual_network_ids` field ([#5794](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5794))
* Data Source: `azurerm_key_vault` - removing the `sku` block since this has been deprecated in favour of the `sku_name` field ([#5774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5774))
* Data Source: `azurerm_key_vault_key` - removing the deprecated `vault_uri` field ([#5774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5774))
* Data Source: `azurerm_key_vault_secret` - removing the deprecated `vault_uri` field ([#5774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5774))
* Data Source: `azurerm_kubernetes_cluster` - removing the field `dns_prefix` from the `agent_pool_profile` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* Data Source: `azurerm_network_interface` - removing the deprecated field `internal_fqdn` ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* Data Source: `azurerm_private_link_service` - removing the deprecated field `network_interface_ids` ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* Data Source: `azurerm_private_link_endpoint_connection` - the deprecated data source has been removed ([#5844](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5844))
* Data Source: `azurerm_recovery_services_protection_policy_vm` has been renamed to `azurerm_backup_policy_vm` ([#5816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5816))
* Data Source: `azurerm_role_definition` - removing the alias `VirtualMachineContributor` which has been deprecated in favour of the full name `Virtual Machine Contributor` ([#5733](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5733))
* Data Source: `azurerm_storage_account` - removing the `account_encryption_source` field since this is no longer configurable by Azure ([#5668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5668))
* Data Source: `azurerm_storage_account` - removing the `enable_blob_encryption` field since this is no longer configurable by Azure ([#5668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5668))
* Data Source: `azurerm_storage_account` - removing the `enable_file_encryption` field since this is no longer configurable by Azure ([#5668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5668))
* Data Source: `azurerm_scheduler_job_collection` - This data source has been removed since it was deprecated ([#5712](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5712))
* Data Source: `azurerm_subnet` - removing the deprecated `ip_configuration` field ([#5801](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5801))
* Data Source: `azurerm_virtual_network` - removing the deprecated `address_spaces` field ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_api_management` - removing the deprecated `sku` block ([#5725](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5725))
* `azurerm_api_management` - removing the deprecated fields in the `security` block ([#5725](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5725))
* `azurerm_application_gateway` - the field `fqdns` within the `backend_address_pool` block is no longer computed ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_application_gateway` - the field `ip_addresses` within the `backend_address_pool` block is no longer computed ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_application_gateway` - the deprecated field `fqdn_list` within the `backend_address_pool` block has been removed ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_application_gateway` - the deprecated field `ip_address_list` within the `backend_address_pool` block has been removed ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_application_gateway` - the deprecated field `disabled_ssl_protocols` has been removed ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_application_gateway` - the field `disabled_protocols` within the `ssl_policy` block is no longer computed ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_app_service` - removing the field `subnet_mask` from the `site_config` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_app_service` - the field `ip_address` within the `site_config` block now refers to a CIDR block, rather than an IP Address to match the Azure API ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_app_service` - removing the field `virtual_network_name` from the `site_config` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_app_service_plan` - the deprecated `properties` block has been removed since these properties have been moved to the top level ([#5717](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5717))
* `azurerm_app_service_slot` - removing the field `subnet_mask` from the `site_config` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_app_service_slot` - the field `ip_address` within the `site_config` block now refers to a CIDR block, rather than an IP Address to match the Azure API ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_app_service_slot` - removing the field `virtual_network_name` from the `site_config` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_application_gateway` - updating the default value for the `body` field within the `match` block from `*` to an empty string ([#5752](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5752))
* `azurerm_automation_account` - removing the `sku` block which has been deprecated in favour of the `sku_name` field ([#5781](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5781))
* `azurerm_automation_credential` - removing the deprecated `account_name` field ([#5781](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5781))
* `azurerm_automation_runbook` - removing the deprecated `account_name` field ([#5781](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5781))
* `azurerm_automation_schedule` - removing the deprecated `account_name` field ([#5781](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5781))
* `azurerm_autoscale_setting` - the deprecated resource has been removed ([#5844](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5844))
* `azurerm_availability_set` - updating the default value for `managed` from `false` to `true` ([#5724](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5724))
* `azurerm_azuread_application` - This resource has been removed since it was deprecated ([#5748](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5748))
* `azurerm_azuread_service_principal_password` - This resource has been removed since it was deprecated ([#5748](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5748))
* `azurerm_azuread_service_principal` - This resource has been removed since it was deprecated ([#5748](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5748))
* `azurerm_client_config` - removing the deprecated field `service_principal_application_id` ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_client_config` - removing the deprecated field `service_principal_object_id` ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_cognitive_account` - removing the deprecated `sku_name` block ([#5797](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5797))
* `azurerm_connection_monitor` - the deprecated resource has been removed ([#5844](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5844))
* `azurerm_container_group` - removing the `port` field from the `container` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_container_group` - removing the `protocol` field from the `container` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_container_group` - the `ports` field is no longer Computed ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_container_group` - the `protocol` field within the `ports` block is no longer Computed and now defaults to `TCP` ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_container_group` - removing the deprecated field `command` ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_container_registry` - removing the deprecated `storage_account` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_container_service` - This resource has been removed since it was deprecated ([#5709](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5709))
* `azurerm_cosmosdb_mongo_collection` - removing the deprecated `indexes` block ([#5853](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5853))
* `azurerm_ddos_protection_plan` - the deprecated resource has been removed ([#5844](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5844))
* `azurerm_devspace_controller` - removing the deprecated `sku` block ([#5795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5795))
* `azurerm_dns_cname_record` - removing the deprecated `records` field ([#5794](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5794))
* `azurerm_dns_ns_record` - removing the deprecated `records` field ([#5794](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5794))
* `azurerm_dns_zone` - removing the deprecated `zone_type` field ([#5794](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5794))
* `azurerm_dns_zone` - removing the deprecated `registration_virtual_network_ids` field ([#5794](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5794))
* `azurerm_dns_zone` - removing the deprecated `resolution_virtual_network_ids` field ([#5794](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5794))
* `azurerm_eventhub` - removing the deprecated `location` field ([#5793](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5793))
* `azurerm_eventhub_authorization_rule` - removing the deprecated `location` field ([#5793](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5793))
* `azurerm_eventhub_consumer_group` - removing the deprecated `location` field ([#5793](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5793))
* `azurerm_eventhub_namespace` - removing the deprecated `kafka_enabled` field since this is now managed by Azure ([#5793](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5793))
* `azurerm_eventhub_namespace_authorization_rule` - removing the deprecated `location` field ([#5793](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5793))
* `azurerm_firewall` - removing the deprecated field `internal_public_ip_address_id` from the `ip_configuration` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_firewall` - the field `public_ip_address_id` within the `ip_configuration` block is now required ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_frontdoor` -  field `cache_enabled` within the `forwarding_configuration` block now defaults to `false` rather than `true` ([#5852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5852))
* `azurerm_frontdoor` - the field `cache_query_parameter_strip_directive` within the `forwarding_configuration` block now defaults to `StripAll` rather than `StripNone`. ([#5852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5852))
* `azurerm_frontdoor` - the field `forwarding_protocol` within the `forwarding_configuration` block now defaults to `HttpsOnly` rather than `MatchRequest` ([#5852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5852))
* `azurerm_function_app` - removing the field `virtual_network_name` from the `site_config` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_function_app` - updating the field `ip_address` within the `ip_restriction` block to accept a CIDR rather than an IP Address to match the updated API behaviour ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_iot_dps` - This resource has been removed since it was deprecated ([#5753](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5753))
* `azurerm_iot_dps_certificate` - This resource has been removed since it was deprecated ([#5753](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5753))
* `azurerm_iothub`- The deprecated `sku.tier` property will be removed. ([#5790](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5790))
* `azurerm_iothub_dps` - The deprecated `sku.tier` property will be removed. ([#5790](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5790))
* `azurerm_key_vault` - removing the `sku` block since this has been deprecated in favour of the `sku_name` field ([#5774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5774))
* `azurerm_key_vault_access_policy` - removing the deprecated field `vault_name` which has been superseded by the `key_vault_id` field ([#5774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5774))
* `azurerm_key_vault_access_policy` - removing the deprecated field `resource_group_name ` which has been superseded by the `key_vault_id` field ([#5774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5774))
* `azurerm_key_vault_certificate` - removing the deprecated `vault_uri` field ([#5774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5774))
* `azurerm_key_vault_key` - removing the deprecated `vault_uri` field ([#5774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5774))
* `azurerm_key_vault_secret` - removing the deprecated `vault_uri` field ([#5774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5774))
* `azurerm_kubernetes_cluster` - updating the default value for `load_balancer_sku` to `Standard` from `Basic` ([#5747](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5747))
* `azurerm_kubernetes_cluster` - the block `default_node_pool` is now required ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_kubernetes_cluster` - removing the deprecated `agent_pool_profile` block ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_kubernetes_cluster` - the field `enable_pod_security_policy` is no longer computed ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_lb_backend_address_pool` - removing the deprecated `location` field ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_lb_nat_pool` - removing the deprecated `location` field ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_lb_nat_rule` - removing the deprecated `location` field ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_lb_probe` - removing the deprecated `location` field ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_lb_rule` - removing the deprecated `location` field ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_log_analytics_workspace_linked_service` - This resource has been removed since it was deprecated ([#5754](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5754))
* `azurerm_log_analytics_linked_service` - The `resource_id` field has been moved from the `linked_service_properties` block to the top-level and the deprecated field `linked_service_properties` will be removed. This has been replaced by the `resource_id` resource ([#5775](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5775))
* `azurerm_maps_account` - the `sku_name` field is now case-sensitive ([#5776](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5776))
* `azurerm_mariadb_server` - removing the `sku` block since it's been deprecated in favour of the `sku_name` field ([#5777](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5777))
* `azurerm_metric_alertrule` - the deprecated resource has been removed ([#5844](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5844))
* `azurerm_monitor_metric_alert` - updating the default value for `auto_mitigate` from `false` to `true` ([#5773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5773))
* `azurerm_monitor_metric_alertrule` - the deprecated resource has been removed ([#5844](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5844))
* `azurerm_mssql_elasticpool` - removing the deprecated `elastic_pool_properties` block ([#5744](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5744))
* `azurerm_mysql_server` - removing the deprecated `sku` block ([#5743](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5743))
* `azurerm_network_interface` - removing the deprecated `application_gateway_backend_address_pools_ids` field from the `ip_configurations` block ([#5784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5784))
* `azurerm_network_interface` - removing the deprecated `application_security_group_ids ` field from the `ip_configurations` block ([#5784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5784))
* `azurerm_network_interface` - removing the deprecated `load_balancer_backend_address_pools_ids ` field from the `ip_configurations` block ([#5784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5784))
* `azurerm_network_interface` - removing the deprecated `load_balancer_inbound_nat_rules_ids ` field from the `ip_configurations` block ([#5784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5784))
* `azurerm_network_interface` - removing the deprecated `internal_fqdn` field ([#5784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5784))
* `azurerm_network_interface` - removing the `network_security_group_id` field in favour of a new split-out resource `azurerm_network_interface_security_group_association` ([#5784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5784))
* `azurerm_network_interface_application_security_group_association` - removing the `ip_configuration_name` field associations between Network Interfaces and Application Security Groups now need to be made to all IP Configurations ([#5815](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5815))
* `azurerm_network_interface` - the `virtual_machine_id` field is now computed-only since it's not setable ([#5784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5784))
* `azurerm_notification_hub_namesapce` - removing the `sku` block in favour of the `sku_name` argument ([#5722](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5722))
* `azurerm_postgresql_server` - removing the `sku` block which has been deprecated in favour of the `sku_name` field ([#5721](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5721))
* `azurerm_private_link_endpoint` - the deprecated resource has been removed ([#5844](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5844))
* `azurerm_private_link_service` - removing the deprecated field `network_interface_ids` ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_public_ip` - making the `allocation_method` field required ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_public_ip` - removing the deprecated field `public_ip_address_allocation` ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* `azurerm_recovery_network_mapping` - the deprecated resource has been removed ([#5816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5816))
* `azurerm_recovery_replicated_vm` - the deprecated resource has been removed ([#5816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5816))
* `azurerm_recovery_services_fabric` - the deprecated resource has been removed ([#5816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5816))
* `azurerm_recovery_services_protected_vm` - the deprecated resource has been removed ([#5816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5816))
* `azurerm_recovery_services_protection_container` - the deprecated resource has been removed ([#5816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5816))
* `azurerm_recovery_services_protection_container_mapping` - the deprecated resource has been removed ([#5816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5816))
* `azurerm_recovery_services_protection_policy_vm` - the deprecated resource has been removed ([#5816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5816))
* `azurerm_recovery_services_replication_policy` - the deprecated resource has been removed ([#5816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5816))
* `azurerm_relay_namespace` - removing the `sku` block in favour of the `sku_name` field ([#5719](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5719))
* `azurerm_scheduler_job` - This resource has been removed since it was deprecated ([#5712](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5712))
* `azurerm_scheduler_job_collection` - This resource has been removed since it was deprecated ([#5712](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5712))
* `azurerm_storage_account` - updating the default value for `account_kind` from `Storage` to `StorageV2` ([#5850](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5850))
* `azurerm_storage_account` - removing the deprecated `account_type` field ([#5710](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5710))
* `azurerm_storage_account` - removing the deprecated `enable_advanced_threat_protection` field ([#5710](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5710))
* `azurerm_storage_account` - updating the default value for `enable_https_traffic_only` from `false` to `true` ([#5808](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5808))
* `azurerm_storage_account` - removing the `account_encryption_source` field since this is no longer configurable by Azure ([#5668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5668))
* `azurerm_storage_account` - removing the `enable_blob_encryption` field since this is no longer configurable by Azure ([#5668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5668))
* `azurerm_storage_account` - removing the `enable_file_encryption` field since this is no longer configurable by Azure ([#5668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5668))
* `azurerm_storage_blob` - making the `type` field case-sensitive ([#5710](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5710))
* `azurerm_storage_blob` - removing the deprecated `attempts` field ([#5710](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5710))
* `azurerm_storage_blob` - removing the deprecated `resource_group_name` field ([#5710](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5710))
* `azurerm_storage_container` - removing the deprecated `resource_group_name` field ([#5710](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5710))
* `azurerm_storage_container` - removing the deprecated `properties` block ([#5710](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5710))
* `azurerm_storage_queue` - removing the deprecated `resource_group_name` field ([#5710](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5710))
* `azurerm_storage_share` - removing the deprecated `resource_group_name` field ([#5710](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5710))
* `azurerm_storage_table` - removing the deprecated `resource_group_name` field ([#5710](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5710))
* `azurerm_subnet` - removing the deprecated `ip_configuration` field ([#5801](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5801))
* `azurerm_subnet` - removing the deprecated `network_security_group_id` field ([#5801](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5801))
* `azurerm_subnet` - removing the deprecated `route_table_id` field ([#5801](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5801))
* `azurerm_subnet` - making the `actions` list within the `service_delegation` block within the `service_endpoints` block non-computed ([#5801](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5801))
* `azurerm_virtual_network_peering` - `allow_virtual_network_access` now defaults to true, matching the API and Portal behaviours. ([#5832](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5832))
* `azurerm_virtual_wan` - removing the deprecated field `security_provider_name` ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))

IMPROVEMENTS:

* web: updating to API version `2019-08-01` ([#5823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5823))
* Data Source: `azurerm_kubernetes_service_version` - support for filtering of preview releases ([#5662](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5662))
* `azurerm_dedicated_host` - support for setting `sku_name` to `DSv3-Type2` and `ESv3-Type2` ([#5768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5768))
* `azurerm_key_vault` - support for configuring `purge_protection_enabled` ([#5344](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5344))
* `azurerm_key_vault` - support for configuring `soft_delete_enabled` ([#5344](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5344))
* `azurerm_sql_database` - support for configuring `zone_redundant` ([#5772](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5772))
* `azurerm_storage_account` - support for configuring the `static_website` block ([#5649](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5649))
* `azurerm_storage_account` - support for configuring `cors_rules` within the `blob_properties` block ([#5425](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5425))
* `azurerm_subnet` - support for delta updates ([#5801](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5801))
* `azurerm_windows_virtual_machine` - fixing a bug when provisioning from a Shared Gallery image ([#5661](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5661))

BUG FIXES:

* `azurerm_application_insights` - the `application_type` field is now case sensitive as documented ([#5817](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5817))
* `azurerm_api_management_api` - allows blank `path` field ([#5833](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5833))
* `azurerm_eventhub_namespace` - the field `ip_rule` within the `network_rulesets` block now supports a maximum of 128 items ([#5831](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5831))
* `azurerm_eventhub_namespace` - the field `virtual_network_rule` within the `network_rulesets` block now supports a maximum of 128 items ([#5831](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5831))
* `azurerm_linux_virtual_machine` - using the delete custom timeout during deletion ([#5764](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5764))
* `azurerm_netapp_account` - allowing the `-` character to be used in the `name` field ([#5842](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5842))
* `azurerm_network_interface` - the `dns_servers` field now respects ordering ([#5784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5784))
* `azurerm_public_ip_prefix` - fixing the validation for the `prefix_length` to match the Azure API ([#5693](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5693))
* `azurerm_recovery_services_vault` - using the requested cloud rather than the default ([#5825](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5825))
* `azurerm_role_assignment` - validating that the `name` is a UUID ([#5624](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5624))
* `azurerm_signalr_service` - ensuring the SignalR segment is parsed in the correct case ([#5737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5737))
* `azurerm_storage_account` - locking on the storage account resource when updating the storage account ([#5668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5668))
* `azurerm_subnet` - supporting updating of the `enforce_private_link_endpoint_network_policies` field ([#5801](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5801))
* `azurerm_subnet` - supporting updating of the `enforce_private_link_service_network_policies` field ([#5801](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5801))
* `azurerm_windows_virtual_machine` - using the delete custom timeout during deletion ([#5764](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5764))

---

For information on v1.44.0 and prior releases, please see [the v1.44.0 changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
