## 2.37.0 (Unreleased)

FEATURES:

* **New Resource:** `azurerm_log_analytics_cluster` [GH-8946]
* **New Resource:** `azurerm_log_analytics_cluster_customer_managed_key` [GH-8946]
* **New Resource:** `azurerm_security_center_automation` [GH-8781]

IMPROVEMENTS:

* storage: foundational improvements to support toggling between the Data Plane and Resource Manager Storage API's in the future [GH-9314]
* Data Source: `azurerm_kubernetes_node_pool` - exposing `os_disk_type` [GH-9166]
* `azurerm_api_management_api_diagnostic` - support for the `always_log_errors`, `http_correlation_protocol`, `log_client_ip` and `verbosity` attributes [GH-9172]
* `azurerm_api_management_api_diagnostic` - support the `frontend_request`, `frontend_response`, `backend_request` and `backend_response` blocks [GH-9172]
* `azurerm_cosmosdb` - removing the cosmosdb autoscale upper cap [GH-9050]
* `azurerm_firewall` - supports for firewall manager policies [GH-8879]
* `azurerm_kubernetes_cluster` - support for configuring `os_disk_type` within the `default_node_pool` block [GH-9166]
* `azurerm_kubernetes_node_pool` - support for configuring `os_disk_type` [GH-9166]
* `azurerm_linux_virtual_machine` - Support `extensions_time_budget` property [GH-9257]
* `azurerm_linux_virtual_machine` - updating the `dedicated_host_id` no longer forces a new resource [GH-9264]
* `azurerm_kubernetes_cluster` - the block `http_application_routing` within the `addon_profile` block can now be updated/removed [GH-9358]
* `azurerm_mssql_database` - `sku_name` supports more `DWxxxc` options [GH-9370]
* `azurerm_postgresql_server` - increase max storage to 16TiB [GH-9373]
* `azurerm_windows_virtual_machine` - Support `extensions_time_budget` property [GH-9257]
* `azurerm_windows_virtual_machine` - updating the `dedicated_host_id` nolonger forces a new resource [GH-9264]
* `azurerm_windows_virtual_machine` - support for the `patch_mode` property [GH-9258]

BUG FIXES:

* Data Source: `azurerm_key_vault_certificate` - fixing a crash when serializing the certificate policy block [GH-9355]
* `azurerm_cosmosdb_sql_container` - no longer attempts to get throughput settings when cosmos account is serverless [GH-9311]
* `azurerm_key_vault_certificate` - fixing a crash when serializing the certificate policy block [GH-9355]
* `azurerm_resource_group_template_deployment` - fixing an issue during deletion where the API version of nested resources couldn't be determined [GH-9364]

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

---

For information on changes between the v2.30.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.44.0 changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
