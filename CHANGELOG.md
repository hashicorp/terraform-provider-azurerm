## 4.7.0 (October 24, 2024)

FEATURES:

* **New Data Source**: `azurerm_oracle_adbs_character_sets` ([#27698](https://github.com/hashicorp/terraform-provider-azurerm/issues/27698))
* **New Data Source**: `azurerm_oracle_adbs_national_character_sets` ([#27698](https://github.com/hashicorp/terraform-provider-azurerm/issues/27698))
* **New Data Source**: `azurerm_oracle_autonomous_database` ([#27696](https://github.com/hashicorp/terraform-provider-azurerm/issues/27696))
* **New Data Source**: `azurerm_oracle_db_nodes` ([#27698](https://github.com/hashicorp/terraform-provider-azurerm/issues/27698))
* **New Data Source**: `azurerm_oracle_db_system_shapes` ([#27698](https://github.com/hashicorp/terraform-provider-azurerm/issues/27698))
* **New Data Source**: `azurerm_oracle_gi_versions` ([#27698](https://github.com/hashicorp/terraform-provider-azurerm/issues/27698))
* **New Resource**: `azurerm_dev_center_project_pool` ([#27706](https://github.com/hashicorp/terraform-provider-azurerm/issues/27706))
* **New Resource**: `azurerm_oracle_autonomous_database` ([#27696](https://github.com/hashicorp/terraform-provider-azurerm/issues/27696))
* **New Resource**: `azurerm_video_indexer_account` ([#27632](https://github.com/hashicorp/terraform-provider-azurerm/issues/27632))

ENHANCEMENTS:

* Dependencies - update `go-azure-sdk` to `v0.20241021.1074254` ([#27713](https://github.com/hashicorp/terraform-provider-azurerm/issues/27713))
* `newrelic` - upgrade api version to `2024-03-01`  ([#27135](https://github.com/hashicorp/terraform-provider-azurerm/issues/27135))
* `cosmosdb` - upgrade api version to `2024-08-15` ([#27659](https://github.com/hashicorp/terraform-provider-azurerm/issues/27659))
* `azurerm_application_gateway` - support for the new `Basic` SKU value ([#27440](https://github.com/hashicorp/terraform-provider-azurerm/issues/27440))
* `azurerm_consumption_budget_management_group` - the property `notification.threshold_type` can now be updated ([#27511](https://github.com/hashicorp/terraform-provider-azurerm/issues/27511))
* `azurerm_consumption_budget_resource_group` - the property `notification.threshold_type` can now be updated ([#27511](https://github.com/hashicorp/terraform-provider-azurerm/issues/27511))
* `azurerm_container_app` - add support for the `template.container.readiness_probe.initial_delay` and `template.container.startup_probe.initial_delay` properties ([#27551](https://github.com/hashicorp/terraform-provider-azurerm/issues/27551))
* `azurerm_mssql_managed_instance` - the `storage_account_type` property can now be updated ([#27737](https://github.com/hashicorp/terraform-provider-azurerm/issues/27737))

BUG FIXES:

* `azurerm_automation_software_update_configuration` - correct validation to not allow `5` and allow `-1` ([#25574](https://github.com/hashicorp/terraform-provider-azurerm/issues/25574))
* `azurerm_cosmosdb_sql_container` - fix recreation logic for `partition_key_version` ([#27692](https://github.com/hashicorp/terraform-provider-azurerm/issues/27692))
* `azurerm_mssql_database` - updating short term retention policy now works as expected ([#27714](https://github.com/hashicorp/terraform-provider-azurerm/issues/27714))
* `azurerm_network_watcher_flow_log` - fix issue where `tags` were not being updated ([#27389](https://github.com/hashicorp/terraform-provider-azurerm/issues/27389))
* `azurerm_postgresql_flexible_server_virtual_endpoint` - retrieve and parse `replica_server_id` for cross-region scenarios as well as remove custom poller for the delete operation ([#27509](https://github.com/hashicorp/terraform-provider-azurerm/issues/27509))

## 4.6.0 (October 18, 2024)

FEATURES:

* **New Resource**: `azurerm_dev_center_attached_network` ([#27638](https://github.com/hashicorp/terraform-provider-azurerm/issues/27638))
* **New Resource**: `azurerm_oracle_cloud_vm_cluster` ([#27678](https://github.com/hashicorp/terraform-provider-azurerm/issues/27678))
* **New Resource**: `azurerm_oracle_exadata_infrastructure` ([#27678](https://github.com/hashicorp/terraform-provider-azurerm/issues/27678))
* **New Data Source**: `azurerm_oracle_cloud_vm_cluster` ([#27678](https://github.com/hashicorp/terraform-provider-azurerm/issues/27678))
* **New Data Source**: `azurerm_oracle_db_servers` ([#27678](https://github.com/hashicorp/terraform-provider-azurerm/issues/27678))
* **New Data Source**: `azurerm_oracle_exadata_infrastructure` ([#27678](https://github.com/hashicorp/terraform-provider-azurerm/issues/27678))

ENHANCEMENTS:

* `redisenterprise` - upgrade api version to `2024-06-01-preview`  ([#27597](https://github.com/hashicorp/terraform-provider-azurerm/issues/27597))
* `azurerm_app_configuration` - support for premium sku ([#27674](https://github.com/hashicorp/terraform-provider-azurerm/issues/27674))
* `azurerm_container_app` - support for the `max_inactive_revisions` property ([#27598](https://github.com/hashicorp/terraform-provider-azurerm/issues/27598))
* `azurerm_kubernetes_cluster` - remove lock on subnets ([#27583](https://github.com/hashicorp/terraform-provider-azurerm/issues/27583))
* `azurerm_nginx_deployment` - allow updates for `sku` ([#27604](https://github.com/hashicorp/terraform-provider-azurerm/issues/27604))
* `azurerm_fluid_relay_server` - support for the `customer_managed_key` property ([#27581](https://github.com/hashicorp/terraform-provider-azurerm/issues/27581))
* `azurerm_linux_virtual_machine` - support the `UBUNTU_PRO` value for the `license_type` property ([#27534](https://github.com/hashicorp/terraform-provider-azurerm/issues/27534))


BUGS:

* `azurerm_api_management_api_diagnostic` - do not set `OperationNameFormat` when the `identifier` property is `azuremonitor` ([#27456](https://github.com/hashicorp/terraform-provider-azurerm/issues/27456))
* `azurerm_api_management` - prevent a panic ([#27649](https://github.com/hashicorp/terraform-provider-azurerm/issues/27649))
* `azurerm_mssql_database` - make `short_term_retention_policy.backup_interval_in_hours` computed ([#27656](https://github.com/hashicorp/terraform-provider-azurerm/issues/27656))

## 4.5.0 (October 10, 2024)

FEATURES:

* **New Resource**: `azurerm_stack_hci_virtual_hard_disk` ([#27474](https://github.com/hashicorp/terraform-provider-azurerm/issues/27474))

ENHANCEMENTS:

* `azurerm_bastion_host` - support for the `Premium` SKU and `session_recording_enabled` property ([#27278](https://github.com/hashicorp/terraform-provider-azurerm/issues/27278))
* `azurerm_log_analytics_cluster` - the `size_gb` property now supports all of 100, 200, 300, 400, 500, 1000, 2000, 5000, 10000, 25000, and 50000 ([#27616](https://github.com/hashicorp/terraform-provider-azurerm/issues/27616))
* `azurerm_mssql_elasticpool` - allow `PRMS` for the `family` property ([#27615](https://github.com/hashicorp/terraform-provider-azurerm/issues/27615))


BUG FIXES:

* `azurerm_mssql_database` - now creates successfully when elastic pool is hyperscale ([#27505](https://github.com/hashicorp/terraform-provider-azurerm/issues/27505))
* `azurerm_postgresql_flexible_server_configuration` - now locks to prevent conflicts when deploying multiple ([#27355](https://github.com/hashicorp/terraform-provider-azurerm/issues/27355))


## 4.4.0 (October 04, 2024)

ENHANCEMENTS: 

* dependencies - update `github.com/hashicorp/go-azure-sdk` to `v0.20240923.1151247` ([#27491](https://github.com/hashicorp/terraform-provider-azurerm/issues/27491))
* `azurerm_site_recovery_replicated_vm` - support for the `target_virtual_machine_size` property ([#27480](https://github.com/hashicorp/terraform-provider-azurerm/issues/27480))

BUG FIXES:

* `azurerm_app_service_certificate` - `key_vault_secret_id` can now be versionless ([#27537](https://github.com/hashicorp/terraform-provider-azurerm/issues/27537))
* `azurerm_linux_virtual_machine_scale_set` - prevent crash when `auto_upgrade_minor_version_enabled` is nil ([#27353](https://github.com/hashicorp/terraform-provider-azurerm/issues/27353))
* `azurerm_role_assignment` - correctly parse ID when it's a root or provider scope ([#27237](https://github.com/hashicorp/terraform-provider-azurerm/issues/27237))
* `azurerm_storage_blob` - `source_content` is now ForceNew ([#27508](https://github.com/hashicorp/terraform-provider-azurerm/issues/27508))
* `azurerm_virtual_network_gateway_connection` - revert `shared_key` to Optional and Computed ([#27560](https://github.com/hashicorp/terraform-provider-azurerm/issues/27560))

## 4.3.0 (September 19, 2024)

FEATURES:

* **New Resource**: `azurerm_advisor_suppression` ([#26177](https://github.com/hashicorp/terraform-provider-azurerm/issues/26177))
* **New Resource**: `azurerm_data_protection_backup_policy_mysql_flexible_server` ([#26955](https://github.com/hashicorp/terraform-provider-azurerm/issues/26955))
* **New Resource**: `azurerm_key_vault_managed_hardware_security_module_key_rotation_policy` ([#27306](https://github.com/hashicorp/terraform-provider-azurerm/issues/27306))
* **New Resource**: `azurerm_stack_hci_deployment_setting` ([#25646](https://github.com/hashicorp/terraform-provider-azurerm/issues/25646))
* **New Resource**: `azurerm_stack_hci_storage_path` ([#26509](https://github.com/hashicorp/terraform-provider-azurerm/issues/26509))
* **New Data Source**: `azurerm_vpn_server_configuration` ([#27054](https://github.com/hashicorp/terraform-provider-azurerm/issues/27054))

ENHANCEMENTS: 

* `managementgroups` - migrate to `hashicorp/go-azure-sdk` ([#26430](https://github.com/hashicorp/terraform-provider-azurerm/issues/26430))
* `nginx` - upgrade api version to `2024-06-01-preview`  ([#27345](https://github.com/hashicorp/terraform-provider-azurerm/issues/27345))
* `azurerm_linux[windows]_web[function]_app[app_slot]` - upgrade api version from `2023-01-01` to `2023-12-01` ([#27196](https://github.com/hashicorp/terraform-provider-azurerm/issues/27196))
* `azurerm_cosmosdb_account` - support for the capability `EnableNoSQLVectorSearch` ([#27357](https://github.com/hashicorp/terraform-provider-azurerm/issues/27357))azurerm_container_app_custom_domain - fix parsing the certificate ID error #25972
* `azurerm_container_app_custom_domain` - support other certificate types ([#25972](https://github.com/hashicorp/terraform-provider-azurerm/issues/25972))
* `azurerm_linux_virtual_machine_scale_set` - the `zones` property can now be updated without creating a new resource ([#27288](https://github.com/hashicorp/terraform-provider-azurerm/issues/27288))
* `azurerm_orchestrated_virtual_machine_scale_set` - the `zones` property can now be updated without creating a new resource ([#27288](https://github.com/hashicorp/terraform-provider-azurerm/issues/27288))
* `azurerm_role_management_policy` - support for resource scope ([#27205](https://github.com/hashicorp/terraform-provider-azurerm/issues/27205))
* `azurerm_spring_cloud_gateway` - changing the `environment_variables` and `sensitive_environment_variables` properties no longer creates a new resource ([#27404](https://github.com/hashicorp/terraform-provider-azurerm/issues/27404))
* `azurerm_static_web_app` - support for the `public_network_access_enabled` property ([#26345](https://github.com/hashicorp/terraform-provider-azurerm/issues/26345))
* `azurerm_shared_image` - support for the `disk_controller_type_nvme_enabled` property ([#26370](https://github.com/hashicorp/terraform-provider-azurerm/issues/26370))
* `azurerm_storage_blob` - changing the `source` property no longer creates a new resource ([#27394](https://github.com/hashicorp/terraform-provider-azurerm/issues/27394))
* `azurerm_storage_object_replication` - changing the `rules.x. source_container_name` and `rules.x. destination_container_name` properties no longer creates a new resource ([#27394](https://github.com/hashicorp/terraform-provider-azurerm/issues/27394))
* `azurerm_windows_virtual_machine_scale_set` - the `zones` property can now be updated without creating a new resource ([#27288](https://github.com/hashicorp/terraform-provider-azurerm/issues/27288)) 

BUG FIXES:

* `azurerm_application_insights` - fix crash when read for `DataVolumeCap` is `nil` ([#27352](https://github.com/hashicorp/terraform-provider-azurerm/issues/27352))
* `azurerm_container_app` - relax validation on the ingress traffic property ([#27396](https://github.com/hashicorp/terraform-provider-azurerm/issues/27396))
* `azurerm_log_analytics_workspace_table` - will now correctly set `total_retention_in_days` when `sku` is `Basic` ([#27420](https://github.com/hashicorp/terraform-provider-azurerm/issues/27420))

## 4.2.0 (September 12, 2024)

FEATURES:

* **New Resource**: `azurerm_arc_machine` ([#26647](https://github.com/hashicorp/terraform-provider-azurerm/issues/26647))
* **New Resource**: `azurerm_arc_machine_automanage_configuration_assignment` ([#26657](https://github.com/hashicorp/terraform-provider-azurerm/issues/26657)) 

ENHANCEMENTS:

* `network/bastionhosts` - upgrade api version from `2023-11-01` to `2024-01-01` ([#27277](https://github.com/hashicorp/terraform-provider-azurerm/issues/27277))
* `recoveryservices` - upgrade `recoveryservicessiterecovery` from `2022-10-0`1 to `2024-04-01` ([#27281](https://github.com/hashicorp/terraform-provider-azurerm/issues/27281))
* `azurerm_data_protection_backup_vault` - support for the `property cross_region_restore_enabled` property ([#27197](https://github.com/hashicorp/terraform-provider-azurerm/issues/27197))
* `azurem_mssql_managed_instance` - support for the `service_principal_type` property ([#27240](https://github.com/hashicorp/terraform-provider-azurerm/issues/27240))

BUG FIXES:

* `azurerm_cosmosdb_account` - fix crash during state migration ([#27302](https://github.com/hashicorp/terraform-provider-azurerm/issues/27302))
* `azurerm_servicebus_queue` - fix defaults of the `default_message_ttl` and `auto_delete_on_idle` properties ([#27305](https://github.com/hashicorp/terraform-provider-azurerm/issues/27305))

## 4.1.0 (September 05, 2024)

ENHANCEMENTS:

* dependencies - bump `hashicorp/go-azure-sdk` to `v0.20240903.1111904` ([#27268](https://github.com/hashicorp/terraform-provider-azurerm/issues/27268))
* Virtual Machine Scale Sets - upgrade api version from `2024-03-01` to `2024-07-01` ([#27230](https://github.com/hashicorp/terraform-provider-azurerm/issues/27230))
* `hdinsights` - update the HDInsights Node definition validation of VM sizes to include new V5 types ([#27270](https://github.com/hashicorp/terraform-provider-azurerm/issues/27270))
* `azurerm_api_management_logger` - support for the `application_insights.connection_string` property ([#27137](https://github.com/hashicorp/terraform-provider-azurerm/issues/27137))
* `azurerm_bot_service_azure_bot` - will now send the value for the `developer_app_insights_api_key` property ([#27280](https://github.com/hashicorp/terraform-provider-azurerm/issues/27280))
* `azurerm_netapp_volume` - support for the `smb3_protocol_encryption_enabled` property ([#27228](https://github.com/hashicorp/terraform-provider-azurerm/issues/27228))
* `azurerm_subnet` - support `Microsoft.DevOpsInfrastructure` as delegation service ([#27259](https://github.com/hashicorp/terraform-provider-azurerm/issues/27259))

BUG FIXES:

* `azurerm_mysql_flexible_server` - correctly set `source_server_id` in the state file ([#27295](https://github.com/hashicorp/terraform-provider-azurerm/issues/27295))
* `azurerm_cosmosdb_account` - the `ip_range_filter` property now supports IPV4 addresses ([#27208](https://github.com/hashicorp/terraform-provider-azurerm/issues/27208))
* `azurerm_cosmosdb_account` - added state migration for `ip_range_filter` underlying type change from `string` to `set` ([#27276](https://github.com/hashicorp/terraform-provider-azurerm/issues/27276))
* `azurerm_linux_virtual_machine` - the `admin_ssh_key.public_key` property now supports ed25519 ssh keys ([#27202](https://github.com/hashicorp/terraform-provider-azurerm/issues/27202))
* `azurerm_sentinel_automation_rule` - no longer panics when using `condition_json`  ([#27269](https://github.com/hashicorp/terraform-provider-azurerm/issues/27269))
* `azurerm_kubernetes_cluster` -  the `host_encryption_enabled` and `node_public_ip_enabled` properties are now set correctly ([#27218](https://github.com/hashicorp/terraform-provider-azurerm/issues/27218))

## 4.0.1 (August 23, 2024)

BUG FIXES:

* provider: fix a validation bug that prevents `terraform validate` from working when `subscription_id` is not specified ([#27178](https://github.com/hashicorp/terraform-provider-azurerm/issues/27178))
* `azurerm_cognitive_deployment` - fixed replacement of `scale` block with `sku` ([#27173](https://github.com/hashicorp/terraform-provider-azurerm/issues/27173))
* `azurerm_kubernetes_cluster` - prevent a panic ([#27183](https://github.com/hashicorp/terraform-provider-azurerm/issues/27183))
* `azurerm_kubernetes_cluster_node_pool` - prevent a panic caused by renamed `enable_*` properties ([#27164](https://github.com/hashicorp/terraform-provider-azurerm/issues/27164))
* `azurerm_sentinel_data_connector_microsoft_threat_intelligence` - prevent error by removing deprecated property `bing_safety_phishing_url_lookback_date` ([#27171](https://github.com/hashicorp/terraform-provider-azurerm/issues/27171))

## 4.0.0 (August 22, 2024)

NOTES:

* **Major Version**: Version 4.0 of the Azure Provider is a major version - some behaviours have changed and some deprecated fields/resources have been removed - please refer to [the 4.0 upgrade guide for more information](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/4.0-upgrade-guide).
* When upgrading to v4.0 of the AzureRM Provider, we recommend upgrading to the latest version of Terraform Core ([which can be found here](https://www.terraform.io/downloads)).

ENHANCEMENTS:

* Data Source: `azurerm_shared_image` - add support for the `trusted_launch_supported`, `trusted_launch_enabled`, `confidential_vm_supported`, `confidential_vm_enabled`, `accelerated_network_support_enabled` and `hibernation_enabled` properties ([#26975](https://github.com/hashicorp/terraform-provider-azurerm/issues/26975))
* dependencies: updating `hashicorp/go-azure-sdk` to `v0.20240819.1075239` ([#27107](https://github.com/hashicorp/terraform-provider-azurerm/issues/27107))
* `applicationgateways` - updating to use `2023-11-01` ([#26776](https://github.com/hashicorp/terraform-provider-azurerm/issues/26776))
* `containerregistry` - updating to use `2023-06-01-preview` ([#23393](https://github.com/hashicorp/terraform-provider-azurerm/issues/23393))
* `containerservice` - updating to `2024-05-01` ([#27105](https://github.com/hashicorp/terraform-provider-azurerm/issues/27105))
* `mssql` - updating to use `hashicorp/go-azure-sdk` and `023-08-01-preview` ([#27073](https://github.com/hashicorp/terraform-provider-azurerm/issues/27073))
* `mssqlmanagedinstance` - updating to use `hashicorp/go-azure-sdk` and `2023-08-01-preview` ([#26872](https://github.com/hashicorp/terraform-provider-azurerm/issues/26872))
* `azurerm_image` - add support for the `disk_encryption_set_id` property to the `data_disk` block ([#27015](https://github.com/hashicorp/terraform-provider-azurerm/issues/27015))
* `azurerm_log_analytics_workspace_table` - add support for more `total_retention_in_days` and `retention_in_days` values ([#27053](https://github.com/hashicorp/terraform-provider-azurerm/issues/27053))
* `azurerm_mssql_elasticpool` - add support for the `HS_MOPRMS` and `MOPRMS` skus ([#27085](https://github.com/hashicorp/terraform-provider-azurerm/issues/27085))
* `azurerm_netapp_pool` - allow `1` as a valid value for `size_in_tb` ([#27095](https://github.com/hashicorp/terraform-provider-azurerm/issues/27095))
* `azurerm_notification_hub` - add support for the `browser_credential` property ([#27058](https://github.com/hashicorp/terraform-provider-azurerm/issues/27058))
* `azurerm_redis_cache` - add support for the `access_keys_authentication_enabled` property ([#27039](https://github.com/hashicorp/terraform-provider-azurerm/issues/27039))
* `azurerm_role_assignment` - add support for the `/`, `/providers/Microsoft.Capacity` and `/providers/Microsoft.BillingBenefits` scopes ([#26663](https://github.com/hashicorp/terraform-provider-azurerm/issues/26663))
* `azurerm_shared_image` - add support for the `hibernation_enabled` property ([#26975](https://github.com/hashicorp/terraform-provider-azurerm/issues/26975))
* `azurerm_storage_account` - support `queue_encryption_key_type` and `table_encryption_key_type` for more storage account kinds ([#27112](https://github.com/hashicorp/terraform-provider-azurerm/issues/27112))
* `azurerm_web_application_firewall_policy` - add support for the `request_body_enforcement` property ([#27094](https://github.com/hashicorp/terraform-provider-azurerm/issues/27094))

BUG FIXES:

* `azurerm_ip_group_cidr` - fixed the position of the CIDR check to correctly refresh the resource when it's no longer present ([#27103](https://github.com/hashicorp/terraform-provider-azurerm/issues/27103))
* `azurerm_monitor_diagnostic_setting` - add further polling to work around an eventual consistency issue when creating the resource ([#27088](https://github.com/hashicorp/terraform-provider-azurerm/issues/27088))
* `azurerm_storage_account` - prevent API error by populating `infrastructure_encryption_enabled` when updating `customer_managed_key` ([#26971](https://github.com/hashicorp/terraform-provider-azurerm/issues/26971))
* `azurerm_storage_blob_inventory_policy` - the `filter` property can now be set when `scope` is `container` ([#27113](https://github.com/hashicorp/terraform-provider-azurerm/issues/27113))
* `azurerm_virtual_network_dns_servers` - moved locks to prevent the creation of subnets with stale data ([#27036](https://github.com/hashicorp/terraform-provider-azurerm/issues/27036))
* `azurerm_virtual_network_gateway_connection` - allow `0` as a valid value for `ipsec_policy.sa_datasize` ([#27056](https://github.com/hashicorp/terraform-provider-azurerm/issues/27056))

---

For information on changes between the v3.116.0 and v3.0.0 releases, please see [the previous v3.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v3.md).

For information on changes between the v2.99.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v1.44.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
