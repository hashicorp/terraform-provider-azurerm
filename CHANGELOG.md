## 2.55.0 (Unreleased)

FEATURES:

* **New Resource:** `azurerm_api_management_email_template` [GH-10914]
* **New Resource:** `azurerm_communication_service` [GH-11066]
* **New Resource:** `azurerm_express_route_port` [GH-10074]
* **New Resource:** `azurerm_spring_cloud_app_redis_association` [GH-11154]

ENHANCEMENTS:

* `azurerm_hpc_cache_nfs_target` - support for the `access_policy_name` property [GH-11186]
* `azurerm_private_endpoint` - allows for an alias to specified [GH-10779]


BUG FIXES:

* Data Source: `azurerm_dns_zone` - fixing a bug where the Resource ID wouldn't contain the Resource Group name when looking this up [GH-11221]
* `azurerm_media_service_account` - `storage_authentication_type` correctly accepts both `ManagedIdentity` and `System` [GH-11222]
* `azurerm_web_application_firewall_policy` - `http_listener_ids` and `path_based_rule_ids` are now Computed only [GH-11196]

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
