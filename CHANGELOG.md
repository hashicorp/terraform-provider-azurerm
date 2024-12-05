## 4.13.0 (Unreleased)

ENHANCEMENTS:

* `azurerm_cognitive_deployment` - add support for the `dynamic_throttling_enabled` property [GH-28100]
* `azurerm_key_vault_managed_hardware_security_module_key` - `key_type` now supports `oct-HSM` [GH-28171]
* `azurerm_machine_learning_datastore_datalake_gen2` - can now be used with storage account in a different subscription [GH-28123]
* `azurerm_network_watcher_flow_log` - `target_resource_id` supports subnets and network interfaces [GH-28177]

BUG:

* Data Source: `azurerm_logic_app_standard` - update `identity` to support User Assigned Identities [GH-28158]
* `azurerm_cdn_frontdoor_origin_group` - update validation of `interval_in_seconds` to match API behaviour [GH-28143]
* `azurerm_container_group` - retrieve log analytics workspace key from config when updating resource [GH-28025]
* `azurerm_search_service` -  the `partition_count` property can now be up to `3` when using basic sku [GH-28105]

## 4.12.0 (November 28, 2024)

FEATURES:

* **New Data Source**: `azurerm_mssql_managed_database` ([#27026](https://github.com/hashicorp/terraform-provider-azurerm/issues/27026))

BUG FIXES:

* `azurerm_application_insights_api_key` - fix condition that nil checks the list of available API keys to prevent an indefinate loop when keys created outside of Terraform are present ([#28037](https://github.com/hashicorp/terraform-provider-azurerm/issues/28037))
* `azurerm_data_factory_linked_service_azure_sql_database` - send `tenant_id` only if it has been specified ([#28120](https://github.com/hashicorp/terraform-provider-azurerm/issues/28120))
* `azurerm_eventgrid_event_subscription` - fix crash when flattening `advanced_filter` ([#28110](https://github.com/hashicorp/terraform-provider-azurerm/issues/28110))
* `azurerm_virtual_network_gateway` - fix crash issue when specifying `root_certificate ` or `revoked_certificate` ([#28099](https://github.com/hashicorp/terraform-provider-azurerm/issues/28099))

ENHANCEMENTS:

* dependencies - update `go-azure-sdk` to `v0.20241128.1112539` ([#28137](https://github.com/hashicorp/terraform-provider-azurerm/issues/28137))
* `containerapps` - update api version to `2024-03-01` ([#28074](https://github.com/hashicorp/terraform-provider-azurerm/issues/28074))
* `Search` - update api version to `2024-06-01-preview` ([#27803](https://github.com/hashicorp/terraform-provider-azurerm/issues/27803))
* Data Source: `azurerm_logic_app_standard` - add support for the `public_network_access` property ([#27913](https://github.com/hashicorp/terraform-provider-azurerm/issues/27913))
* Data Source: `azurerm_search_service` - add support for the `customer_managed_key_encryption_compliance_status` property ([#27478](https://github.com/hashicorp/terraform-provider-azurerm/issues/27478))
* `azurerm_container_registry_task` - add validation on `cpu` as well as on `agent_pool_name`and `agent_setting` ([#28098](https://github.com/hashicorp/terraform-provider-azurerm/issues/28098))
* `azurerm_databricks_workspace` - add support for the `enhanced_security_compliance` block ([#26606](https://github.com/hashicorp/terraform-provider-azurerm/issues/26606))
* `azurerm_eventhub` - deprecate `namespace_name` and `resource_group_name` in favour of `namespace_id` ([#28055](https://github.com/hashicorp/terraform-provider-azurerm/issues/28055))
* `azurerm_logic_app_standard` - add support for the `public_network_access` property ([#27913](https://github.com/hashicorp/terraform-provider-azurerm/issues/27913))
* `azurerm_search_service` - add support for the `customer_managed_key_encryption_compliance_status` property ([#27478](https://github.com/hashicorp/terraform-provider-azurerm/issues/27478))
* `azurerm_cosmosdb_account` - add support for value `EnableNoSQLFullTextSearch` in the  `capabilities.name` property ([#28114](https://github.com/hashicorp/terraform-provider-azurerm/issues/28114))

## 4.11.0 (November 22, 2024)

NOTES:
* New [ephemeral resources](https://developer.hashicorp.com/terraform/language/v1.10.x/resources/ephemeral) `azurerm_key_vault_certificate` and `azurerm_key_vault_secret` now support [ephemeral values](https://developer.hashicorp.com/terraform/language/v1.10.x/values/variables#exclude-values-from-state)

FEATURES:

* **New Ephemeral Resource**: `azurerm_key_vault_certificate` ([#28083](https://github.com/hashicorp/terraform-provider-azurerm/issues/28083))
* **New Ephemeral Resource**: `azurerm_key_vault_secret` ([#28083](https://github.com/hashicorp/terraform-provider-azurerm/issues/28083))
* **New Resource**: `azurerm_eventgrid_namespace` ([#27682](https://github.com/hashicorp/terraform-provider-azurerm/issues/27682))

ENHANCEMENTS:

* dependencies: update `hashicorp/go-azure-sdk` to `v0.20241118.1115603` ([#28075](https://github.com/hashicorp/terraform-provider-azurerm/issues/28075))
* `batch` - upgrade api version to `2024-07-01` ([#27982](https://github.com/hashicorp/terraform-provider-azurerm/issues/27982))
* `containerregistry` - upgrade api version to `2023-11-01-preview` ([#27983](https://github.com/hashicorp/terraform-provider-azurerm/issues/27983))
* `azurerm_application_gateway` - `1.1` is now accepted as a valid `rule_set_version` in the `waf_configuration` block ([#28039](https://github.com/hashicorp/terraform-provider-azurerm/issues/28039))
* `azurerm_arc_machine` - add support for the `identity` and `tags` properties ([#27987](https://github.com/hashicorp/terraform-provider-azurerm/issues/27987))
* `azurerm_container_app` - `secret.name` now accepts up to 253 characters and `.` ([#27935](https://github.com/hashicorp/terraform-provider-azurerm/issues/27935))
* `azurerm_network_manager` - `scope_accesses` now accepts `Routing` ([#28033](https://github.com/hashicorp/terraform-provider-azurerm/issues/28033))
* `azurerm_network_watcher_flow_log` - add support for the `target_resource_id` property ([#26015](https://github.com/hashicorp/terraform-provider-azurerm/issues/26015))
* `azurerm_role_assignment` - `condition_version` will be defaulted to `2.0` when `condition` has been set ([#27189](https://github.com/hashicorp/terraform-provider-azurerm/issues/27189))
* `azurerm_subnet` - `Informatica.DataManagement/organizations` is a valid `service_delegation` ([#27993](https://github.com/hashicorp/terraform-provider-azurerm/issues/27993))
* `azurerm_virtual_network` - `Informatica.DataManagement/organizations` is a valid `service_delegation` ([#27993](https://github.com/hashicorp/terraform-provider-azurerm/issues/27993))
* `azurerm_web_application_firewall_policy` - `1.1` is now accepted as a valid `version` for `Microsoft_BotManagerRuleSet` rule types ([#28039](https://github.com/hashicorp/terraform-provider-azurerm/issues/28039))

BUG FIXES:

* `azurerm_api_management` - `public_ip_address_id` is no longer required when `zone` has been set ([#27976](https://github.com/hashicorp/terraform-provider-azurerm/issues/27976))
* `azurerm_api_management_diagnostic` - raise and error when `operation_name_format` is used with and `identity` that is not `applicationinsights` ([#27630](https://github.com/hashicorp/terraform-provider-azurerm/issues/27630))
* `azurerm_api_management_api_diagnostic` - raise and error when `operation_name_format` is used with and `identity` that is not `applicationinsights` ([#27630](https://github.com/hashicorp/terraform-provider-azurerm/issues/27630))
* `azurerm_application_gateway` - `rewrite_rule_set` can be supplied when using `Basic` sku ([#28011](https://github.com/hashicorp/terraform-provider-azurerm/issues/28011))
* `azurerm_container_registry_token_password` - correctly mark as gone if container registry token doesn't exist ([#27232](https://github.com/hashicorp/terraform-provider-azurerm/issues/27232))
* `azurerm_kusto_cluster` - `allowed_fqdn` and `allowed_ip_ranges` can now be set to empty lists ([#27529](https://github.com/hashicorp/terraform-provider-azurerm/issues/27529))
* `azurerm_linux_function_app_slot` - create content settings when using a consumpton plan ([#25412](https://github.com/hashicorp/terraform-provider-azurerm/issues/25412))
* `azurerm_virtual_network_gatway` - updating `ip_configuration` now recreates the resource ([#27828](https://github.com/hashicorp/terraform-provider-azurerm/issues/27828))


## 4.10.0 (November 14, 2024)

BREAKING CHANGES:

* dependencies - update `cognitive` to `2024-10-01`, due to a behavioural change in this version of the API, the `primary_access_key` and `secondary_access_key` can not be retrieved if `local_authentication_enabled` has been set to `false`. These properties that may have had values previously will now be empty. This has affected the `azurerm_ai_services` and `azurerm_cognitive_account` resources as well as the `azurerm_cognitive_account` data source ([#27851](https://github.com/hashicorp/terraform-provider-azurerm/issues/27851))

FEATURES:

* **New Data Source**: `azurerm_key_vault_managed_hardware_security_module_key` ([#27827](https://github.com/hashicorp/terraform-provider-azurerm/issues/27827))
* **New Resource**: `azurerm_netapp_backup_vault` ([#27188](https://github.com/hashicorp/terraform-provider-azurerm/issues/27188))
* **New Resource**: `azurerm_netapp_backup_policy` ([#27188](https://github.com/hashicorp/terraform-provider-azurerm/issues/27188))

ENHANCEMENTS:

* dependencies: update `terraform-plugin-framework` to version `v1.13.0` ([#27936](https://github.com/hashicorp/terraform-provider-azurerm/issues/27936))
* dependencies: update `terraform-plugin-framework-validators` to version `v0.14.0` ([#27936](https://github.com/hashicorp/terraform-provider-azurerm/issues/27936))
* dependencies: update `terraform-plugin-go` to version `v0.25.0` ([#27936](https://github.com/hashicorp/terraform-provider-azurerm/issues/27936))
* dependencies: update `terraform-plugin-mux` to version `v0.17.0` ([#27936](https://github.com/hashicorp/terraform-provider-azurerm/issues/27936))
* dependencies: update `terraform-plugin-sdk/v2` to version `v2.35.0` ([#27936](https://github.com/hashicorp/terraform-provider-azurerm/issues/27936))
* Data Source: `azurerm_bastion_host` - add support for the `zones` property ([#27909](https://github.com/hashicorp/terraform-provider-azurerm/issues/27909))
* `azurerm_application_gateway` - support more values for the `status_code` property ([#27535](https://github.com/hashicorp/terraform-provider-azurerm/issues/27535))
* `azurerm_bastion_host` - support for the `zones` property ([#27909](https://github.com/hashicorp/terraform-provider-azurerm/issues/27909))
* `azurerm_communication_service` - support for `usgov` region ([#27919](https://github.com/hashicorp/terraform-provider-azurerm/issues/27919))
* `azurerm_email_communication_service` - support for `usgov` region added ([#27919](https://github.com/hashicorp/terraform-provider-azurerm/issues/27919))
* `azurerm_linux_function_app` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_linux_function_app_slot` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_linux_web_app` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_linux_web_app_slot` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_windows_web_app` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_windows_web_app_slot` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_windows_function_app` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))
* `azurerm_windows_function_app_slot` - support for .NET 9  ([#27879](https://github.com/hashicorp/terraform-provider-azurerm/issues/27879))

BUG FIXES:

* `azurerm_log_analytics_workspace_table` - use the subscription from workspace ID ([#27590](https://github.com/hashicorp/terraform-provider-azurerm/issues/27590))
* `azurerm_traffic_manager_external_endpoint` - the value for `priority` will be dynamically assigned by the API ([#27966](https://github.com/hashicorp/terraform-provider-azurerm/issues/27966))
* `azurerm_traffic_manager_azure_endpoint` - the value for `priority` will be dynamically assigned by the API ([#27966](https://github.com/hashicorp/terraform-provider-azurerm/issues/27966))


## 4.9.0 (November 08, 2024)

FEATURES:

* **New Resource**: `azurerm_dynatrace_monitor` ([#27432](https://github.com/hashicorp/terraform-provider-azurerm/issues/27432))
* **New Resource**: `azurerm_dashboard_grafana_managed_private_endpoint` ([#27781](https://github.com/hashicorp/terraform-provider-azurerm/issues/27781))
* **New Resource**: `azurerm_data_protection_backup_instance_mysql_flexible_server` ([#27464](https://github.com/hashicorp/terraform-provider-azurerm/issues/27464))
* **New Resource**: `azurerm_mongo_cluster` ([#27636](https://github.com/hashicorp/terraform-provider-azurerm/issues/27636))
* **New Resource**: `azurerm_stack_hci_network_interface` ([#26888](https://github.com/hashicorp/terraform-provider-azurerm/issues/26888))

ENHANCEMENTS:

* dependencies - update `go-azure-sdk` to `v0.20241104.1140654` ([#27896](https://github.com/hashicorp/terraform-provider-azurerm/issues/27896))
* dependencies - update `go-azure-helpers` to `v0.71.0` ([#27897](https://github.com/hashicorp/terraform-provider-azurerm/issues/27897))
* dependencies - update `golang-jwt` to `v4.5.1` ([#27938](https://github.com/hashicorp/terraform-provider-azurerm/issues/27938))
* `storage` - allow `azurerm_storage_account` to be used in Data Plane restrictive environments ([#27818](https://github.com/hashicorp/terraform-provider-azurerm/issues/27818))
* `azurerm_cognitive_deployment` - `sku.0.name` now supports `DataZoneStandard` ([#27926](https://github.com/hashicorp/terraform-provider-azurerm/issues/27926))
* `azurerm_mssql_managed_database` - support for the `tags` property ([#27857](https://github.com/hashicorp/terraform-provider-azurerm/issues/27857))
* `azurerm_oracle_cloud_vm_cluster` - support for the `domain`, `scan_listener_port_tcp`, `scan_listener_port_tcp_ssl` and `zone_id` properties ([#27808](https://github.com/hashicorp/terraform-provider-azurerm/issues/27808))
* `azurerm_public_ip_prefix` - support for the `sku_tier` property ([#27882](https://github.com/hashicorp/terraform-provider-azurerm/issues/27882))
* `azurerm_public_ip` - support for the `domain_name_label_scope` property ([#27748](https://github.com/hashicorp/terraform-provider-azurerm/issues/27748))
* `azurerm_subnet` - `default_outbound_access_enabled` can now be updated ([#27858](https://github.com/hashicorp/terraform-provider-azurerm/issues/27858))
* `azurerm_storage_container` - support for the `storage_account_id` property ([#27733](https://github.com/hashicorp/terraform-provider-azurerm/issues/27733))
* `azurerm_storage_share` - support for the `storage_account_id` property ([#27733](https://github.com/hashicorp/terraform-provider-azurerm/issues/27733))

## 4.8.0 (October 31, 2024)

FEATURES:

* **New Data Source**: `azurerm_virtual_network_peering` ([#27530](https://github.com/hashicorp/terraform-provider-azurerm/issues/27530))
* **New Resource**: `azurerm_machine_learning_workspace_network_outbound_rule_fqdn` ([#27384](https://github.com/hashicorp/terraform-provider-azurerm/issues/27384))
* **New Resource**: `azurerm_stack_hci_extension` ([#26929](https://github.com/hashicorp/terraform-provider-azurerm/issues/26929))
* **New Resource**: `azurerm_stack_hci_marketplace_gallery_image` ([#27532](https://github.com/hashicorp/terraform-provider-azurerm/issues/27532))
* **New Resource**: `azurerm_trusted_signing_account` ([#27720](https://github.com/hashicorp/terraform-provider-azurerm/issues/27720))

ENHANCEMENTS:

* `mysql` - upgrade api version to `2023-12-30` ([#27767](https://github.com/hashicorp/terraform-provider-azurerm/issues/27767))
* `network` - upgrade api version to `2024-03-01 ` ([#27746](https://github.com/hashicorp/terraform-provider-azurerm/issues/27746))
* `azurerm_cosmosdb_account`: support for CMK through `managed_hsm_key_id` property ([#26521](https://github.com/hashicorp/terraform-provider-azurerm/issues/26521))
* `azurerm_cosmosdb_account` - support further versions for `mongo_server_version` ([#27763](https://github.com/hashicorp/terraform-provider-azurerm/issues/27763))
* `azurerm_container_app_environment` - changing the `log_analytics_workspace_id` property no longer creates a new resource ([#27794](https://github.com/hashicorp/terraform-provider-azurerm/issues/27794))
* `azurerm_data_factory_linked_service_azure_sql_database` - add support for the `credential_name` property ([#27629](https://github.com/hashicorp/terraform-provider-azurerm/issues/27629))
* `azurerm_key_vault_key` - `expiration_date` only recreates the resource when it is removed from the config file ([#27813](https://github.com/hashicorp/terraform-provider-azurerm/issues/27813))
* `azurerm_kubernetes_cluster` - fix issue where`maintenance_window_auto_upgrade`/`maintenance_window_auto_upgrade`/`maintenance_window_node_os ` might not be read into state ([#26915](https://github.com/hashicorp/terraform-provider-azurerm/issues/26915))
* `azurerm_kubernetes_cluster` - support for the `backend_pool_type` property ([#27596](https://github.com/hashicorp/terraform-provider-azurerm/issues/27596))
* `azurerm_kubernetes_cluster` - support for the `daemonset_eviction_for_empty_nodes_enabled`, `daemonset_eviction_for_occupied_nodes_enabled`, and `ignore_daemonsets_utilization_enabled` properties ([#27588](https://github.com/hashicorp/terraform-provider-azurerm/issues/27588))
* `azurerm_load_test` - `description` can now be updated ([#27800](https://github.com/hashicorp/terraform-provider-azurerm/issues/27800))
* `azurerm_oracle_cloud_vm_cluster` - export the `ocid` property ([#27785](https://github.com/hashicorp/terraform-provider-azurerm/issues/27785))
* `azurerm_orchestrated_virtual_machine_scale_set` - add support for `sku_profile` block ([#27599](https://github.com/hashicorp/terraform-provider-azurerm/issues/27599))
* `azurerm_web_application_firewall_policy` - add support for `policy_settings.0.file_upload_enforcement` ([#27774](https://github.com/hashicorp/terraform-provider-azurerm/issues/27774))

BUG FIXES:

* `azurerm_automation_hybrid_runbook_worker_group` - correctly mark resource as gone if it's absent when reading it ([#27797](https://github.com/hashicorp/terraform-provider-azurerm/issues/27797))
* `azurerm_automation_hybrid_runbook_worker` - correctly mark resource as gone if it's absent when reading it ([#27797](https://github.com/hashicorp/terraform-provider-azurerm/issues/27797))
* `azurerm_automation_python3_package` - correctly mark resource as gone if it's absent when reading it ([#27797](https://github.com/hashicorp/terraform-provider-azurerm/issues/27797))
* `azurerm_data_protection_backup_vault` - prevent panic when checking value of `cross_region_restore_enabled` ([#27762](https://github.com/hashicorp/terraform-provider-azurerm/issues/27762))
* `azurerm_role_management_policy` - fix panic when unmarshalling the policy into a specific type ([#27731](https://github.com/hashicorp/terraform-provider-azurerm/issues/27731))
* `azurerm_security_center_subscription_pricing` - correctly type assert the `additional_extension_properties` property when building the payload ([#27721](https://github.com/hashicorp/terraform-provider-azurerm/issues/27721))
* `azurerm_synapse_workspace_aad_admin` - will no correctly delete when using `azurerm_synapse_workspace_aad_admin` with `azurerm_synapse_workspace` ([#27606](https://github.com/hashicorp/terraform-provider-azurerm/issues/27606))
* `azurerm_windows_function_app_slot` - fixed panic in state migration ([#27700](https://github.com/hashicorp/terraform-provider-azurerm/issues/27700))

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

* dependencies - update `go-azure-sdk` to `v0.20241021.1074254` ([#27713](https://github.com/hashicorp/terraform-provider-azurerm/issues/27713))
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
