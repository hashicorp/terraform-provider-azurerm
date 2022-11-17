## 3.32.0 (November 17, 2022)

DEPRECATIONS

* The `azurerm_integration_service_environment` resource is now deprecated as the underlying Azure Service is being retired on `2024-08-31` and new instances cannot be provisioned (by default) after `2022-11-01` ([#19265](https://github.com/hashicorp/terraform-provider-azurerm/issues/19265))

ENHANCEMENTS:

* dependencies: updating to `v0.20221116.1175352` of `github.com/hashicorp/go-azure-sdk` ([#19319](https://github.com/hashicorp/terraform-provider-azurerm/issues/19319))
* `azurerm_security_center_subscription_pricing` - support for the `subplan` property ([#19273](https://github.com/hashicorp/terraform-provider-azurerm/issues/19273))
* `azurerm_storage_account` - support for the `sas_policy` block ([#19222](https://github.com/hashicorp/terraform-provider-azurerm/issues/19222))
* `azurerm_windows_web_app`, `azurerm_windows_web_app_slot` - aupport for `17` value for `java_version` property ([#19249](https://github.com/hashicorp/terraform-provider-azurerm/issues/19249))
* `azurerm_storage_blob_inventory_policy` - support for `include_deleted` property ([#19286](https://github.com/hashicorp/terraform-provider-azurerm/issues/19286))

BUG FIXES:

* `azurerm_app_service_public_certificate` - add custom poller to prevent `Root resource was present, but now absent.` result ([#19348](https://github.com/hashicorp/terraform-provider-azurerm/issues/19348))
* `azurerm_eventhub_namespace` - correct `zone_redundant` property ([#19164](https://github.com/hashicorp/terraform-provider-azurerm/issues/19164))
* `azurerm_orchestrated_virtual_machine_scale_set` - allow no image to be specified ([#19263](https://github.com/hashicorp/terraform-provider-azurerm/issues/19263))
* `azurerm_synapse_firewall_rule` - wait for the firewall to be ready ([#19227](https://github.com/hashicorp/terraform-provider-azurerm/issues/19227))
* `azurerm_service_fabric_managed_cluster` - correctly define `active_directory` as a List ([#19163](https://github.com/hashicorp/terraform-provider-azurerm/issues/19163))
* `azurerm_orchestrated_virtual_machine_scale_set` -  instance parameter is now set on update ([#19337](https://github.com/hashicorp/terraform-provider-azurerm/issues/19337))


## 3.31.0 (November 10, 2022)

FEATURES:

* **New Resource:** `azurerm_federated_identity_credential` ([#19199](https://github.com/hashicorp/terraform-provider-azurerm/issues/19199))
* **New Resource:** `azurerm_stream_analytics_stream_input_eventhub_v2` ([#19150](https://github.com/hashicorp/terraform-provider-azurerm/issues/19150))

ENHANCEMENTS

* dependencies: updating to `v0.20221108.1145701` of `github.com/hashicorp/go-azure-sdk` ([#19193](https://github.com/hashicorp/terraform-provider-azurerm/issues/19193))
* dependencies: updating `network` to API Version `2022-05-01` ([#19124](https://github.com/hashicorp/terraform-provider-azurerm/issues/19124))
* dependencies: updating `sentinel` to API Version `2022-10-01-preview` ([#19161](https://github.com/hashicorp/terraform-provider-azurerm/issues/19161))
* `azurerm_disk_encryption_set` - support for the `federated_client_id` property ([#19184](https://github.com/hashicorp/terraform-provider-azurerm/issues/19184))
* `azurerm_linux_web_app` - support for .NET 7 ([#19232](https://github.com/hashicorp/terraform-provider-azurerm/issues/19232))
* `azurerm_linux_function_app` - support for .NET 7 ([#19232](https://github.com/hashicorp/terraform-provider-azurerm/issues/19232))
* `azurerm_managed_disk` - support for expanding data disks without downtime ([#17245](https://github.com/hashicorp/terraform-provider-azurerm/issues/17245))
* `azurerm_mssql_virtual_machine` - support for the `sql_instance` block ([#19123](https://github.com/hashicorp/terraform-provider-azurerm/issues/19123))
* `azurerm_public_ip` - support for the `ddos_protection_mode` and `ddos_protection_plan_id` properties ([#19206](https://github.com/hashicorp/terraform-provider-azurerm/issues/19206))
* `azurerm_sentinel_alert_rule_nrt` - support for the `techniques` property ([#19142](https://github.com/hashicorp/terraform-provider-azurerm/issues/19142))
* `azurerm_sentinel_alert_rule_fusion` - support for the source block ([#19093](https://github.com/hashicorp/terraform-provider-azurerm/issues/19093))
* `azurerm_windows_web_app` - support for .NET 7 ([#19232](https://github.com/hashicorp/terraform-provider-azurerm/issues/19232))
* `azurerm_windows_function_app` - support for .NET 7 ([#19232](https://github.com/hashicorp/terraform-provider-azurerm/issues/19232))

BUG FIXES:

* `azurerm_cdn_frontdoor_route` - update read function to parse `cdn_frontdoor_origin_group_id` insensitively ([#19178](https://github.com/hashicorp/terraform-provider-azurerm/issues/19178))
* `azurerm_cdn_frontdoor_rule` - update `url_redirect_action` to allow `query_string` field to pass multiple query string parameters ([#19180](https://github.com/hashicorp/terraform-provider-azurerm/issues/19180))
* `azurerm_cdn_frontdoor_profile` - the `response_timeout_seconds` property is no longer force new ([#19175](https://github.com/hashicorp/terraform-provider-azurerm/issues/19175))

## 3.30.0 (November 03, 2022)

FEATURES:

* **New Resource:** `azurerm_kubernetes_fleet_manager` ([#19111](https://github.com/hashicorp/terraform-provider-azurerm/issues/19111))
* **New Resource:** `azurerm_mssql_server_microsoft_support_auditing_policy` ([#18609](https://github.com/hashicorp/terraform-provider-azurerm/issues/18609))
* **New Resource:** `azurerm_private_dns_resolver_virtual_network_link` ([#19029](https://github.com/hashicorp/terraform-provider-azurerm/issues/19029))
* **New Resource:** `azurerm_private_dns_resolver_forwarding_rule` ([#19028](https://github.com/hashicorp/terraform-provider-azurerm/issues/19028))

ENHANCEMENTS

* dependencies: `iothub` updating to `2022-04-30-preview` ([#19070](https://github.com/hashicorp/terraform-provider-azurerm/issues/19070))
* dependencies: updating to `v0.47.0` of `github.com/hashicorp/go-azure-helpers` ([#19107](https://github.com/hashicorp/terraform-provider-azurerm/issues/19107))
* dependencies: updating to `v0.20221102.1171058` of `github.com/hashicorp/go-azure-sdk` ([#19108](https://github.com/hashicorp/terraform-provider-azurerm/issues/19108))
* webpubsub: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#18892](https://github.com/hashicorp/terraform-provider-azurerm/issues/18892))
* Data Source: `azurerm_application_gateway` - export the `backend_address_pool` block ([#19026](https://github.com/hashicorp/terraform-provider-azurerm/issues/19026))
* Data Source: `azurerm_function_app_host_keys` - export `webpubsub_extension_key` property ([#19073](https://github.com/hashicorp/terraform-provider-azurerm/issues/19073))
* `azurerm_iothub` - support for `DigitalTwinChangeEvents` as `source` and `fallback_route.source` ([#19070](https://github.com/hashicorp/terraform-provider-azurerm/issues/19070))
* `azurerm_iothub_fallback_route` - support for `DigitalTwinChangeEvents` as `source` ([#19070](https://github.com/hashicorp/terraform-provider-azurerm/issues/19070))
* `azurerm_iothub_route` - support for `DigitalTwinChangeEvents` as `source` ([#19070](https://github.com/hashicorp/terraform-provider-azurerm/issues/19070))
* `azurerm_kubernetes_cluster` - support for the `web_app_routing` block ([#18667](https://github.com/hashicorp/terraform-provider-azurerm/issues/18667))
* `azurerm_linux_virtual_machine_scale_set` - support for the `protected_settings_from_key_vault` blovk ([#19098](https://github.com/hashicorp/terraform-provider-azurerm/issues/19098))
* `azurerm_linux_virtual_machine_scale_set` - support for `StandardSSD_ZRS`, `PremiumV2_LRS`, and `Premium_ZRS` storage account types ([#19091](https://github.com/hashicorp/terraform-provider-azurerm/issues/19091))
* `azurerm_mssql_virtual_machine` - support for the `system_db_on_data_disk_enabled` property ([#19115](https://github.com/hashicorp/terraform-provider-azurerm/issues/19115))
* `azurerm_monitor_diagnostic_setting` - support for the `partner_solution_id` property ([#19114](https://github.com/hashicorp/terraform-provider-azurerm/issues/19114))
* `azurerm_policy_definition` - reverse the order of policies lookup to favour builtin ([#18338](https://github.com/hashicorp/terraform-provider-azurerm/issues/18338))
* `azurerm_policy_set_definition` - reverse the order of policies lookup to favour builtin ([#18338](https://github.com/hashicorp/terraform-provider-azurerm/issues/18338))
* `azurerm_security_center_contact` - support for the `name` property ([#18999](https://github.com/hashicorp/terraform-provider-azurerm/issues/18999))
* `azurerm_stream_analytics_job` - support for the `job_storage_account` block ([#19120](https://github.com/hashicorp/terraform-provider-azurerm/issues/19120))
* `azurerm_virtual_machine_extension` - support for the `protected_settings_from_key_vault` blovk ([#19098](https://github.com/hashicorp/terraform-provider-azurerm/issues/19098))
* `azurerm_virtual_machine_scale_set_extension` - support for the `protected_settings_from_key_vault` blovk ([#19098](https://github.com/hashicorp/terraform-provider-azurerm/issues/19098))
* `azurerm_windows_virtual_machine_scale_set` - support for the `protected_settings_from_key_vault` blovk ([#19098](https://github.com/hashicorp/terraform-provider-azurerm/issues/19098))
* `azurerm_windows_virtual_machine_scale_set` - support for `StandardSSD_ZRS`, `PremiumV2_LRS`, and `Premium_ZRS` storage account types ([#19091](https://github.com/hashicorp/terraform-provider-azurerm/issues/19091))

BUG FIXES:

* Data Source: `azurerm_app_configuration_keys` - fix a crash when `label` is not set ([#19032](https://github.com/hashicorp/terraform-provider-azurerm/issues/19032))
* `azurerm_api_management` - correct the api return `subnet_id` with the wrong case ([#18988](https://github.com/hashicorp/terraform-provider-azurerm/issues/18988))
* `azurerm_cdn_frontdoor_firewall_policy` - expose `AnomalyScoring` in override rule action for DRS `2.0` ([#19095](https://github.com/hashicorp/terraform-provider-azurerm/issues/19095))
* `azurerm_eventhub_namespace_disaster_recovery_config` - will now correctly break the pairing ([#19030](https://github.com/hashicorp/terraform-provider-azurerm/issues/19030))
* `azurerm_kubernetes_cluster` - set a valid default value for `auto_scaler_profile.expander` ([#19057](https://github.com/hashicorp/terraform-provider-azurerm/issues/19057))
* `azurerm_linux_virtual_machine_scale_set` - can now set `automatic_os_upgrade_policy` with rolling upgrades enables ([#18605](https://github.com/hashicorp/terraform-provider-azurerm/issues/18605))
* `azurerm_mssql_database` - handle the `license_type` property no longer being returned by API ([#19084](https://github.com/hashicorp/terraform-provider-azurerm/issues/19084))
* `azurerm_postgresql_flexible_server_database` - is now correctly removed from state on deletion ([#19081](https://github.com/hashicorp/terraform-provider-azurerm/issues/19081))
* `azurerm_virtual_network_gateway_connection` - correctly set `authorization_key` from state as the API returnes `*`s ([#19071](https://github.com/hashicorp/terraform-provider-azurerm/issues/19071))
* `azurerm_windows_virtual_machine_scale_set` - can now set `automatic_os_upgrade_policy` with rolling upgrades enables ([#18605](https://github.com/hashicorp/terraform-provider-azurerm/issues/18605))

## 3.29.1 (October 28, 2022)

BUG FIXES:

* `azurerm_kubernetes_cluster` - prevent panic when setting `public_network_access_enabled` ([#19048](https://github.com/hashicorp/terraform-provider-azurerm/issues/19048))

## 3.29.0 (October 27, 2022)

FEATURES:

* **New Data Source:** `azurerm_api_management_gateway_host_name_configuration` ([#17166](https://github.com/hashicorp/terraform-provider-azurerm/issues/17166))
* **New Data Source:** `azurerm_cdn_frontdoor_firewall_policy` ([#18903](https://github.com/hashicorp/terraform-provider-azurerm/issues/18903))
* **New Resource:** `azurerm_datadog_monitor_tag_rule` ([#17825](https://github.com/hashicorp/terraform-provider-azurerm/issues/17825))
* **New Resource:** `azurerm_datadog_monitor_sso_configuration` ([#17825](https://github.com/hashicorp/terraform-provider-azurerm/issues/17825))
* **New Resource:** `azurerm_iothub_device_update_account` ([#18789](https://github.com/hashicorp/terraform-provider-azurerm/issues/18789))
* **New Resource:** `azurerm_iothub_device_update_instance` ([#18789](https://github.com/hashicorp/terraform-provider-azurerm/issues/18789))
* **New Resource:** `azurerm_nginx_configuration` ([#18761](https://github.com/hashicorp/terraform-provider-azurerm/issues/18761))
* **New Resource:** `azurerm_nginx_certificate` ([#18762](https://github.com/hashicorp/terraform-provider-azurerm/issues/18762))
* **New Resource:** `azurerm_private_dns_resolver` ([#18473](https://github.com/hashicorp/terraform-provider-azurerm/issues/18473))
* **New Resource:** `azurerm_private_dns_resolver_dns_forwarding_ruleset` ([#19012](https://github.com/hashicorp/terraform-provider-azurerm/issues/19012))
* **New Resource:** `azurerm_private_dns_resolver_inbound_endpoint` ([#18983](https://github.com/hashicorp/terraform-provider-azurerm/issues/18983))
* **New Resource:** `azurerm_private_dns_resolver_outbound_endpoint` ([#18986](https://github.com/hashicorp/terraform-provider-azurerm/issues/18986))

ENHANCEMENTS:

* Dependencies `compute` - updating to `2022-08-01` ([#18994](https://github.com/hashicorp/terraform-provider-azurerm/issues/18994))
* Dependencies `containerinstance` - updating to `2021-10-01` ([#17785](https://github.com/hashicorp/terraform-provider-azurerm/issues/17785))
* Dependencies: update `go-azure-helpers` to `v0.45.0` ([#18968](https://github.com/hashicorp/terraform-provider-azurerm/issues/18968))
* containerservice: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#18705](https://github.com/hashicorp/terraform-provider-azurerm/issues/18705))
* customproviders - refactoring to use `github.com/hashicorp/go-azure-sdk` ([#18978](https://github.com/hashicorp/terraform-provider-azurerm/issues/18978))
* snapshot - refactoring to use `github.com/hashicorp/go-azure-sdk` ([#18957](https://github.com/hashicorp/terraform-provider-azurerm/issues/18957))
* disks: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#18928](https://github.com/hashicorp/terraform-provider-azurerm/issues/18928))
* Data Source: `azurerm_storage_management_policy` - add support for `tier_to_archive_after_days_since_last_tier_change_greater_than` ([#18898](https://github.com/hashicorp/terraform-provider-azurerm/issues/18898))
* `azurerm_container_group` - the `network_profile_id` property hasbeen deprecated in favour of `subnet_ids` as the newer versions of the API no longer support it ([#17785](https://github.com/hashicorp/terraform-provider-azurerm/issues/17785))
* `azurerm_cdn_frontdoor_rule` - allow the `cdn_frontdoor_origin_group_id` field to be optional in the `route_configuration_override_action` ([#18906](https://github.com/hashicorp/terraform-provider-azurerm/issues/18906))
* `azurerm_cdn_frontdoor_rule` - expose `Disabled` as a possible value of `cache_behavior` in the `route_configuration_override_action` ([#18906](https://github.com/hashicorp/terraform-provider-azurerm/issues/18906))
* `azurerm_disk_encryption_set` - support for identities `UserAssigned` and `SystemAssigned,UserAssgined` ([#18525](https://github.com/hashicorp/terraform-provider-azurerm/issues/18525))
* `azurerm_hdinsight_kafka_cluster` - support for the `compute_isolation` block ([#17449](https://github.com/hashicorp/terraform-provider-azurerm/issues/17449))
* `azurerm_hdinsight_spark_cluster` - support for the `compute_isolation` block ([#17449](https://github.com/hashicorp/terraform-provider-azurerm/issues/17449))
* `azurerm_hdinsight_interactive_query_cluster` - support for the `compute_isolation` block ([#17449](https://github.com/hashicorp/terraform-provider-azurerm/issues/17449))
* `azurerm_hdinsight_hbase_cluster` - support for the `compute_isolation` block ([#17449](https://github.com/hashicorp/terraform-provider-azurerm/issues/17449))
* `azurerm_hdinsight_hadoop_cluster` - support for the `compute_isolation` block ([#17449](https://github.com/hashicorp/terraform-provider-azurerm/issues/17449))
* `azurerm_container_group` - support for the `dns_name_label_reuse_policy` block ([#17785](https://github.com/hashicorp/terraform-provider-azurerm/issues/17785))
* `azurerm_kubernetes_cluster` - support for the `workload_autoscaler_profile` block ([#18967](https://github.com/hashicorp/terraform-provider-azurerm/issues/18967))
* `azurerm_linux_function_app` - support for using `storage_account` external Azure Storage Account configurations ([#18760](https://github.com/hashicorp/terraform-provider-azurerm/issues/18760))
* `azurerm_linux_function_app` - support for Java 17 ([#18689](https://github.com/hashicorp/terraform-provider-azurerm/issues/18689))
* `azurerm_linux_function_app_slot` - support for using `storage_account` external Azure Storage Account configurations ([#18760](https://github.com/hashicorp/terraform-provider-azurerm/issues/18760))
* `azurerm_logic_app_action_http` - support for the `queries` property ([#18934](https://github.com/hashicorp/terraform-provider-azurerm/issues/18934))
* `azurerm_monitor_scheduled_query_rules_alert_v2` - add `evaluation_frequency`, `window_duration`, `mute_actions_after_alert_duration`, and `query_time_range_override`validation ([#18960](https://github.com/hashicorp/terraform-provider-azurerm/issues/18960))
* `azurerm_mssql_virtual_machine` - =support for the `assessment` block ([#18923](https://github.com/hashicorp/terraform-provider-azurerm/issues/18923))
* `azurerm_mssql_server_transparent_data_encryption` - support for autorotation of keyvault keys ([#18523](https://github.com/hashicorp/terraform-provider-azurerm/issues/18523))
* `azurerm_logic_app_standard` - support for the `scm_ip_restriction` block and the `scm_use_main_ip_restriction`, `scm_min_tls_version`, `scm_type` properties ([#18853](https://github.com/hashicorp/terraform-provider-azurerm/issues/18853))
* `azurerm_postgresql_server` - can now set `public_network_access_enabled` during creation when point in time restore is used ([#18962](https://github.com/hashicorp/terraform-provider-azurerm/issues/18962))
* `azurerm_servicebus_namespace_disaster_recovery_config` - support the `alias_authorization_rule_id` property ([#18729](https://github.com/hashicorp/terraform-provider-azurerm/issues/18729))
* `azurerm_synapse_workspace` - `sql_administrator_login` and `sql_administrator_login_password` are now no longer required for the creation of this resource ([#18850](https://github.com/hashicorp/terraform-provider-azurerm/issues/18850))
* `azurerm_synapse_workspace` - enable user assigned managed identity ([#19007](https://github.com/hashicorp/terraform-provider-azurerm/issues/19007))
* `azurerm_windows_function_app` - support for using `storage_account` external Azure Storage Account configurations ([#18760](https://github.com/hashicorp/terraform-provider-azurerm/issues/18760))
* `azurerm_windows_function_app` - support for Java 17 ([#18689](https://github.com/hashicorp/terraform-provider-azurerm/issues/18689))
* `azurerm_windows_function_app_slot` - support for using `storage_account` external Azure Storage Account configurations ([#18760](https://github.com/hashicorp/terraform-provider-azurerm/issues/18760))

BUG FIXES:

*  provider: will no loner automatically register the `Microsoft.StoragePool` provider as Azure has halted the preview of Azure Disk Pools, and it will not be made generally available ([#18905](https://github.com/hashicorp/terraform-provider-azurerm/issues/18905))
*  `azurerm_app_configuration_keys` - will now correctly retrieve result if more than 100 entries are returned ([#19020](https://github.com/hashicorp/terraform-provider-azurerm/issues/19020))
* `azurerm_data_factory_dataset_parquet` - `azure_blob_storage_location.path` and `http_server_location.path` are now Optional ([#19009](https://github.com/hashicorp/terraform-provider-azurerm/issues/19009))
* `azurerm_disk_pool` - has been deprecated as Azure has halted the preview of Azure Disk Pools, and it will not be made generally available ([#18905](https://github.com/hashicorp/terraform-provider-azurerm/issues/18905))
* `azurerm_disk_pool_iscsi_target` - has been deprecated as Azure has halted the preview of Azure Disk Pools, and it will not be made generally available ([#18905](https://github.com/hashicorp/terraform-provider-azurerm/issues/18905))
* `azurerm_disk_pool_iscsi_target_lun` - has been deprecated as Azure has halted the preview of Azure Disk Pools, and it will not be made generally available ([#18905](https://github.com/hashicorp/terraform-provider-azurerm/issues/18905))
* `azurerm_disk_pool_managed_disk_attachment` - has been deprecated as Azure has halted the preview of Azure Disk Pools, and it will not be made generally available ([#18905](https://github.com/hashicorp/terraform-provider-azurerm/issues/18905))
* `azurerm_linux_virtual_machine_scale_set` - the `gallery_applications` block has been renamted to `gallery_application` ([#19014](https://github.com/hashicorp/terraform-provider-azurerm/issues/19014))
* `azurerm_managed_disk` - `logical_sector_size`, `disk_iops_read_write`, `disk_mbps_read_write`, `disk_iops_read_only` and `disk_mbps_read_only` can be set when `storage_account_type` is `PremiumV2_LRS` ([#18991](https://github.com/hashicorp/terraform-provider-azurerm/issues/18991))
* `azurerm_monitor_data_collection_rule` - correctly support streams ([#18966](https://github.com/hashicorp/terraform-provider-azurerm/issues/18966))
* `azurerm_netapp_volume` - correctly set snapshot ID when `create_from_snapshot_resource_id` is specified ([#18996](https://github.com/hashicorp/terraform-provider-azurerm/issues/18996))
* `azurerm_key_vault_certificate` - new versions of key vault certs can now be imported ([#18848](https://github.com/hashicorp/terraform-provider-azurerm/issues/18848))
* `azurerm_postgresql_server` - correctly create replicas when CMK is enabled ([#18805](https://github.com/hashicorp/terraform-provider-azurerm/issues/18805))
* `azurerm_stream_analytics_stream_input_eventhub` - `shared_access_policy_key` and `shared_access_policy_name` are no longer required ([#18959](https://github.com/hashicorp/terraform-provider-azurerm/issues/18959))
* `azurerm_windows_virtual_machine_scale_set` - the `gallery_applications` block has been renamted to `gallery_application` ([#19014](https://github.com/hashicorp/terraform-provider-azurerm/issues/19014))

## 3.28.0 (October 20, 2022)

FEATURES:

* **New Data Source:** `azurerm_cdn_frontdoor_secret` ([#18817](https://github.com/hashicorp/terraform-provider-azurerm/issues/18817))
* **New Resource:** `azurerm_databricks_access_connector` ([#18709](https://github.com/hashicorp/terraform-provider-azurerm/issues/18709))
* **New Resource:** `azurerm_sentinel_data_connector_dynamics_365` ([#18859](https://github.com/hashicorp/terraform-provider-azurerm/issues/18859))
* **New Resource:** `azurerm_sentinel_data_connector_iot` ([#18862](https://github.com/hashicorp/terraform-provider-azurerm/issues/18862))
* **New Resource:** `azurerm_sentinel_data_connector_office_365_project` ([#18858](https://github.com/hashicorp/terraform-provider-azurerm/issues/18858))
* **New Resource:** `azurerm_sentinel_data_connector_office_irm` ([#18856](https://github.com/hashicorp/terraform-provider-azurerm/issues/18856))
* **New Resource:** `azurerm_sentinel_data_connector_office_power_bi` ([#18857](https://github.com/hashicorp/terraform-provider-azurerm/issues/18857))

ENHANCEMENTS:

* dependencies: updating to `v0.20221018.1075906` of `github.com/hashicorp/go-azure-sdk` ([#18833](https://github.com/hashicorp/terraform-provider-azurerm/issues/18833))
* `azurestackhci`: updating to API Version `2022-09-01` ([#18759](https://github.com/hashicorp/terraform-provider-azurerm/issues/18759))
* Data Source: `azurerm_linux_function_app` - add support for `client_certificate_exclusion_paths ` ([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))
* Data Source: `azurerm_linux_web_app` - add support for `client_certificate_exclusion_paths ` ([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))
* Data Source: `azurerm_windows_function_app` - add support for `client_certificate_exclusion_paths ` ([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))
* Data Source: `azurerm_windows_web_app` - add support for `client_certificate_exclusion_paths ` ([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))
* `azurerm_cdn_frontdoor_firewall_policy` - managed rules can now exclude matches on `RequestBodyJsonArgNames` ([#18874](https://github.com/hashicorp/terraform-provider-azurerm/issues/18874))
* `azurerm_cosmosdb_account` - support for the `primary_sql_connection_string`, `secondary_sql_connection_string`, `primary_readonly_sql_connection_string`, and `secondary_readonly_sql_connection_string` attributes ([#17810](https://github.com/hashicorp/terraform-provider-azurerm/issues/17810))
* `azurerm_fluid_relay_server` - support for the `service_endpoint` property ([#18763](https://github.com/hashicorp/terraform-provider-azurerm/issues/18763))
* `azurerm_fluid_relay_server` - support for the `primary_key` and `secondary_key` properties ([#18765](https://github.com/hashicorp/terraform-provider-azurerm/issues/18765))
* `azurerm_linux_function_app` - correctly set `use_32_bit_worker` during Create ([#18680](https://github.com/hashicorp/terraform-provider-azurerm/issues/18680))
* `azurerm_linux_function_app` - add support for the `client_certificate_exclusion_paths` property ([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))
* `azurerm_linux_function_app` - add `VS2022` to `remote_debugging_version` valid values ([#18684](https://github.com/hashicorp/terraform-provider-azurerm/issues/18684))
* `azurerm_linux_function_app_slot` - add support for the `client_certificate_exclusion_paths` property ([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))
* `azurerm_linux_web_app` - add support for the `client_certificate_exclusion_paths` property([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))
* `azurerm_linux_web_app_slot` - add support for the `client_certificate_exclusion_paths` property ([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))
* `azurerm_storage_account` - support for the `immutability_policy` block ([#18774](https://github.com/hashicorp/terraform-provider-azurerm/issues/18774))
* `azurerm_storage_account` - customer managed keys can be now enabled when `account_tier` is set to `Premium` ([#18872](https://github.com/hashicorp/terraform-provider-azurerm/issues/18872))
* `azurerm_storage_management_policy` - support for the `tier_to_archive_after_days_since_last_tier_change_greater_than` property ([#18792](https://github.com/hashicorp/terraform-provider-azurerm/issues/18792))
* `azurerm_subnet` - add support for `Microsoft.LabServices/labplans` ([#18822](https://github.com/hashicorp/terraform-provider-azurerm/issues/18822))
* `azurerm_windows_virtual_machine_scale_set` - allow disabling secure boot when creating a virtual machine scale set with disk encryption type `VMGuestStateOnly` ([#18749](https://github.com/hashicorp/terraform-provider-azurerm/issues/18749))
* `azurerm_windows_function_app` - correctly  set `use_32_bit_worker` during Create ([#18680](https://github.com/hashicorp/terraform-provider-azurerm/issues/18680))
* `azurerm_windows_function_app` - add support for the `client_certificate_exclusion_paths` property ([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))
* `azurerm_windows_function_app` - add `VS2022` to `remote_debugging_version` valid values ([#18684](https://github.com/hashicorp/terraform-provider-azurerm/issues/18684))
* `azurerm_windows_function_app_slot` - add support for the `client_certificate_exclusion_paths` correctly ([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))
* `azurerm_windows_web_app` - add support for the `client_certificate_exclusion_paths` correctly ([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))
* `azurerm_windows_web_app_slot` - add support for the `client_certificate_exclusion_paths` correctly  ([#16603](https://github.com/hashicorp/terraform-provider-azurerm/issues/16603))

BUG FIXES:

* `azurerm_automation_software_update_configuration` - parse subscription IDs correctly when set in `scope` ([#18860](https://github.com/hashicorp/terraform-provider-azurerm/issues/18860))
* `azurerum_cdn_frontdoor_route` - fix a panic on import ([#18824](https://github.com/hashicorp/terraform-provider-azurerm/issues/18824))
* `azurerm_eventhub_namespace` - ignore case for `network_rulesets.x.virtual_network_rule.x.subnet_id` ([#18818](https://github.com/hashicorp/terraform-provider-azurerm/issues/18818))
* `azurerm_firewall_policy_rule_collection_group` - limit the number of destination ports in a NAT rule to one ([#18766](https://github.com/hashicorp/terraform-provider-azurerm/issues/18766))
* Data Source: `azurerm_linux_function_app` - fix missing error on data source not found ([#18876](https://github.com/hashicorp/terraform-provider-azurerm/issues/18876))
* `azurerm_linux_function_app` - fix an issue where `app_settings` would show a diff when setting `vnet_route_all_enabled` to true ([#18836](https://github.com/hashicorp/terraform-provider-azurerm/issues/18836))
* `azurerm_linux_function_app_slot` - fix an issue where `app_settings` would show a diff when setting `vnet_route_all_enabled` to true ([#18836](https://github.com/hashicorp/terraform-provider-azurerm/issues/18836))
* `azurerm_linux_virtual_machine` - allow disabling secure boot when creating a virtual machine with disk encryption type `VMGuestStateOnly` ([#18749](https://github.com/hashicorp/terraform-provider-azurerm/issues/18749))
* `azurerm_linux_virtual_machine_scale_set` - allow disabling secure boot when creating a virtual machine scale set with disk encryption type `VMGuestStateOnly` ([#18749](https://github.com/hashicorp/terraform-provider-azurerm/issues/18749))
* `azurerm_network_security_group` - correct the casing of the `protocol` property ([#18799](https://github.com/hashicorp/terraform-provider-azurerm/issues/18799))
* `azurerm_network_security_rule` - correct the casing of the `protocol` property ([#18799](https://github.com/hashicorp/terraform-provider-azurerm/issues/18799))
* `azurerm_recovery_services_vault` - fix issue where `soft_delete_enabled` is reset to the default value when the `identity` block is updated ([#18871](https://github.com/hashicorp/terraform-provider-azurerm/issues/18871))
* `azurerm_windows_virtual_machine` - allow disabling secure boot when creating a virtual machine with disk encryption type `VMGuestStateOnly` ([#18749](https://github.com/hashicorp/terraform-provider-azurerm/issues/18749))
* `azurerm_windows_function_app` - fix an issue where `app_settings` would show a diff when setting `vnet_route_all_enabled` to true ([#18836](https://github.com/hashicorp/terraform-provider-azurerm/issues/18836))
* `azurerm_windows_function_app_slot` - fix an issue where `app_settings` would show a diff when setting `vnet_route_all_enabled` to true ([#18836](https://github.com/hashicorp/terraform-provider-azurerm/issues/18836))
* `azurerm_windows_web_app` - fix parsing of `docker_container_name` and `docker_container_registry` on read ([#18251](https://github.com/hashicorp/terraform-provider-azurerm/issues/18251))

## 3.27.0 (October 13, 2022)

BREAKING CHANGES:

* `azurerm_cdn_frontdoor_custom_domain` - removed the `associate_with_cdn_frontdoor_route_id` field to allow for a custom domain to be associated with multiple routes. ([#18600](https://github.com/hashicorp/terraform-provider-azurerm/issues/18600))

FEATURES:

* **New DataSource:** `data.azurerm_cosmosdb_sql_role_definition` ([#18728](https://github.com/hashicorp/terraform-provider-azurerm/issues/18728))
* **New DataSource:** `data.azurerm_cosmosdb_sql_database` ([#18728](https://github.com/hashicorp/terraform-provider-azurerm/issues/18728))
* **New Resource:** `azurerm_cdn_frontdoor_custom_domain_association` ([#18600](https://github.com/hashicorp/terraform-provider-azurerm/issues/18600))
* **New Resource:** `azurerm_nginx_deployment` ([#18510](https://github.com/hashicorp/terraform-provider-azurerm/issues/18510))
* **New Resource:** `azurerm_orbital_contact_profile` ([#18317](https://github.com/hashicorp/terraform-provider-azurerm/issues/18317))
* **New Resource:** `azurerm_sentinel_data_connector_office_atp` ([#18708](https://github.com/hashicorp/terraform-provider-azurerm/issues/18708))

ENHANCEMENTS:

* dependencies: updating to version `v0.44.` of `github.com/hashicorp/go-azure-helpers` ([#18716](https://github.com/hashicorp/terraform-provider-azurerm/issues/18716))
* dependencies: updating to version `v0.50.0` of `github.com/manicminer/hamilton` ([#18716](https://github.com/hashicorp/terraform-provider-azurerm/issues/18716))
* `azurerm_automation_runbook` - support for the `draft` block and `log_activity_trace` propertry ([#17961](https://github.com/hashicorp/terraform-provider-azurerm/issues/17961))
* `azurerm_app_configuration` - support for the `encrption`, `local_auth_enabled`, `public_network_access_enabled`, `purge_protection_enabled`, and `soft_delete_retention_days` properties ([#17714](https://github.com/hashicorp/terraform-provider-azurerm/issues/17714))
* `azurerm_api_management_api` - support for the `contact` and `license` blocks ([#18472](https://github.com/hashicorp/terraform-provider-azurerm/issues/18472))
* `azurerm_cdn_frontdoor_route` - exposed `cdn_frontdoor_custom_domain_ids` and `link_to_default_domain` ([#18600](https://github.com/hashicorp/terraform-provider-azurerm/issues/18600))
* `azurerm_data_factory_integration_runtime_azure_ssis` - support for `elastic_pool_namr` property ([#18696](https://github.com/hashicorp/terraform-provider-azurerm/issues/18696))
* `azurerm_dedicated_hardware_security_module` - support the `management_network_profile` block ([#18702](https://github.com/hashicorp/terraform-provider-azurerm/issues/18702))
* `azurerm_hdinsight_hadoop_cluster`, - support for the `script_actions` block ([#18670](https://github.com/hashicorp/terraform-provider-azurerm/issues/18670))
* `azurerm_hdinsight_hbase_cluster`,  - support for the `script_actions` block ([#18670](https://github.com/hashicorp/terraform-provider-azurerm/issues/18670))
* `azurerm_hdinsight_interactive_query_cluster`, - support for the `script_actions` block ([#18670](https://github.com/hashicorp/terraform-provider-azurerm/issues/18670))
* `azurerm_spark_cluster` - support for the `script_actions` block ([#18670](https://github.com/hashicorp/terraform-provider-azurerm/issues/18670))
* `azurerm_kubernetes_cluster` - support the `workload_identity_enabled` property ([#18742](https://github.com/hashicorp/terraform-provider-azurerm/issues/18742))
* `azurerm_firewall_policy_rule_collection_group`- add `Mssql` as an option for `type` validation ([#18746](https://github.com/hashicorp/terraform-provider-azurerm/issues/18746))
* `azurerm_log_analytics_cluster` - ensuring that the `identity` block is always set ([#18700](https://github.com/hashicorp/terraform-provider-azurerm/issues/18700))
* `azurerm_linux_web_app` - support for python `3.10` ([#18744](https://github.com/hashicorp/terraform-provider-azurerm/issues/18744))
* `azurerm_linux_web_app_slot` - support for python `3.10` ([#18744](https://github.com/hashicorp/terraform-provider-azurerm/issues/18744))
* `azurerm_mssql_database` - support for the `import` block ([#18588](https://github.com/hashicorp/terraform-provider-azurerm/issues/18588))
* `azurerm_stream_analytics_output_servicebus_queue` - support for the `authentication_mode` property ([#18491](https://github.com/hashicorp/terraform-provider-azurerm/issues/18491))

BUG FIXES: 

* `azurerm_kubernetes_cluster` - `orchestrator_version` is set properly for clusters created with an older API version ([#18130](https://github.com/hashicorp/terraform-provider-azurerm/issues/18130))
* `azurerm_kubernetes_cluster_node_pool` - `orchestrator_version` is set properly for node pools created with an older API version ([#18130](https://github.com/hashicorp/terraform-provider-azurerm/issues/18130))
* `azurerm_log_analytics_cluster` - fixing an issue when checking for the presence of an existing Log Analytics Cluster ([#18700](https://github.com/hashicorp/terraform-provider-azurerm/issues/18700))
* `azurerm_logic_app_workflow` - can now be updated when associated with `integration_service_environment_id` ([#18660](https://github.com/hashicorp/terraform-provider-azurerm/issues/18660))
* `azurerm_spring_cloud_connection` - correctly parse storage blob resource id ([#18699](https://github.com/hashicorp/terraform-provider-azurerm/issues/18699))
* `azurerm_app_service_connection` - correctly parse storage blob resource id ([#18699](https://github.com/hashicorp/terraform-provider-azurerm/issues/18699))


## 3.26.0 (October 06, 2022)

BREAKING CHANGES:

* `azurerm_load_test` - the computed attribute `dataplane_uri` has been renamed to `data_plane_uri` for consistency ([#18654](https://github.com/hashicorp/terraform-provider-azurerm/issues/18654))

FEATURES:

* **New Resource:** `azurerm_iotcentral_application_network_rule_set` ([#18589](https://github.com/hashicorp/terraform-provider-azurerm/issues/18589))

ENHANCEMENTS:

* dependencies: updating to `v0.43.0` of `github.com/hashicorp/go-azure-helpers` ([#18630](https://github.com/hashicorp/terraform-provider-azurerm/issues/18630))
* dependencies: updating to `v0.20221004.1155444` of `github.com/hashicorp/go-azure-sdk` ([#18628](https://github.com/hashicorp/terraform-provider-azurerm/issues/18628))
* provider: support for auto-registering SDK Clients and Services ([#18629](https://github.com/hashicorp/terraform-provider-azurerm/issues/18629))
* `azurerm_batch_pool` - support for the `node_deallocation_method`, `dynamic_vnet_assignment_scope`, and `source_port_ranges` properties ([#18436](https://github.com/hashicorp/terraform-provider-azurerm/issues/18436))
* `azurerm_kubernetes_cluster` - support for `pod_cidrs` and `service_cidrs` properties ([#16657](https://github.com/hashicorp/terraform-provider-azurerm/issues/16657))
* `azurerm_kubernetes_cluster` - support for `message_of_the_day`, `managed_outbound_ipv6_count`, `scale_down_mode` and `workload_runtime` properties ([#16741](https://github.com/hashicorp/terraform-provider-azurerm/issues/16741))
* `azurerm_kubernetes_cluster_node_pool` - support for `message_of_the_day`, `scale_down_mode` and `workload_runtime` properties ([#16741](https://github.com/hashicorp/terraform-provider-azurerm/issues/16741))
* `azurerm_load_test` - switching to an auto-generated resource ([#18654](https://github.com/hashicorp/terraform-provider-azurerm/issues/18654))
* `azurerm_load_test` - the computed attribute `dataplane_uri` has been renamed to `data_plane_uri` for consistency ([#18654](https://github.com/hashicorp/terraform-provider-azurerm/issues/18654))
* `azurerm_load_test` - support for the `description` field ([#18654](https://github.com/hashicorp/terraform-provider-azurerm/issues/18654))
* `azurerm_user_assigned_identity` - switching to an auto-generated resource ([#18654](https://github.com/hashicorp/terraform-provider-azurerm/issues/18654))

BUG FIXES:

* `azurerm_linux_function_app_slot` - read app settings from the correct endpoint ([#18396](https://github.com/hashicorp/terraform-provider-azurerm/issues/18396))
* `azurerm_load_test` - changing the `name` field now forces a new resource to be created ([#18654](https://github.com/hashicorp/terraform-provider-azurerm/issues/18654))
* `azurerm_windows_function_app_slot` - read app settings from the correct endpoint ([#18396](https://github.com/hashicorp/terraform-provider-azurerm/issues/18396))

## 3.25.0 (September 29, 2022)

FEATURES:

* **New Resource:** `azurerm_cdn_frontdoor_route` ([#18231](https://github.com/hashicorp/terraform-provider-azurerm/issues/18231))
* **New Resource:** `azurerm_cdn_frontdoor_custom_domain` ([#18231](https://github.com/hashicorp/terraform-provider-azurerm/issues/18231))
* **New Resource:** `azurerm_cdn_route_disable_link_to_default_domain` ([#18231](https://github.com/hashicorp/terraform-provider-azurerm/issues/18231))

ENHANCEMENTS:

* dependencies: `machinelearning` - updating to use `2022-05-01` ([#17671](https://github.com/hashicorp/terraform-provider-azurerm/issues/17671))
* dependencies: updating to version `v0.20220921.1082044` of `github.com/hashicorp/go-azure-sdk` ([#18557](https://github.com/hashicorp/terraform-provider-azurerm/issues/18557))
* provider: support for the `oidc_token_file_path` property and `ARM_OIDC_TOKEN_FILE_PATH` environment variable ([#18335](https://github.com/hashicorp/terraform-provider-azurerm/issues/18335))
* Data Source: `azurerm_databricks_workspace` - exports the `location` propertuy ([#18521](https://github.com/hashicorp/terraform-provider-azurerm/issues/18521))
* `azurerm_api_management` - support for the `additional_location.gateway_disabled`, `certificate_source`, and `certificate_status` properties ([#18508](https://github.com/hashicorp/terraform-provider-azurerm/issues/18508))
* `azurerm_automation_software_update_configuration` - the `classification` property has been deprecated in favour of the `classifications` property that supports multiple values ([#18539](https://github.com/hashicorp/terraform-provider-azurerm/issues/18539))
* `azurerm_healthcare_fhir_service` - support for the `oci_artifact` block ([#18571](https://github.com/hashicorp/terraform-provider-azurerm/issues/18571))
* `azurerm_healthcare_fhir` - support for the `public_network_access_enabled` property ([#18566](https://github.com/hashicorp/terraform-provider-azurerm/issues/18566))
* `azurerm_iotcentral_application` - support for the `identity` and `public_network_access_enabled` properties ([#18564](https://github.com/hashicorp/terraform-provider-azurerm/issues/18564))
* `azurerm_linux_virtual_machine` - support for the `gallery_application` property ([#18406](https://github.com/hashicorp/terraform-provider-azurerm/issues/18406))
* `azurerm_machine_learning_workspace` - support for the `public_network_access_enabled` and `v1_legacy_mode` properties ([#18469](https://github.com/hashicorp/terraform-provider-azurerm/issues/18469))
* `azurerm_storage_account` - support for the `multichannel_enabled` property ([#17999](https://github.com/hashicorp/terraform-provider-azurerm/issues/17999))
* `azurerm_virtual_hub_bgp_connection` - support for the `virtual_network_connection_id` property ([#18469](https://github.com/hashicorp/terraform-provider-azurerm/issues/18469))
* `azurerm_windows_virtual_machine` - support for the `gallery_application` property ([#18406](https://github.com/hashicorp/terraform-provider-azurerm/issues/18406))

BUG FIXES:

* Data Source: `azurerm_key_vault_certificate_data` - correctly create PEM private key block header for EC keys ([#18419](https://github.com/hashicorp/terraform-provider-azurerm/issues/18419))
* `azurerm_log_analytics_linked_storage_account` - correctly `data_source_type` case handling ([#18116](https://github.com/hashicorp/terraform-provider-azurerm/issues/18116))

## 3.24.0 (September 22, 2022)

FEATURES:

* **New Resource**: `azurerm_automation_software_update_configuration` ([#17902](https://github.com/hashicorp/terraform-provider-azurerm/issues/17902))
* **New Resource**: `azurerm_monitor_alert_processing_rule_action_group` ([#17006](https://github.com/hashicorp/terraform-provider-azurerm/issues/17006))
* **New Resource**: `azurerm_monitor_alert_processing_rule_suppression` ([#17006](https://github.com/hashicorp/terraform-provider-azurerm/issues/17006))

ENHANCEMENTS:

* dependencies: updating to version `v0.20220916.1125744` of `github.com/hashicorp/go-azure-sdk` ([#18446](https://github.com/hashicorp/terraform-provider-azurerm/issues/18446))
* dependencies: `disks` - updating to use `2022-03-02` ([#17671](https://github.com/hashicorp/terraform-provider-azurerm/issues/17671))
* Data Source: `azurerm_automation_account` - exports the `identity` attribute ([#18478](https://github.com/hashicorp/terraform-provider-azurerm/issues/18478))
* Data Source: `azurerm_storage_account` - export the `azure_files_identity_based_auth` property ([#18405](https://github.com/hashicorp/terraform-provider-azurerm/issues/18405))
* `azurerm_api_management_api_operation` - support the `example`, `schema_id`, and `type_name` properties ([#18409](https://github.com/hashicorp/terraform-provider-azurerm/issues/18409))
* `azurerm_cognitive_account` - support for the `customer_managed_key` property ([#18516](https://github.com/hashicorp/terraform-provider-azurerm/issues/18516))
* `azurerm_data_factory_flowlet_data_flow` - support for the `rejected_linked_service` property ([#18056](https://github.com/hashicorp/terraform-provider-azurerm/issues/18056))
* `azurerm_data_factory_data_flow` - support for the `rejected_linked_service` property ([#18056](https://github.com/hashicorp/terraform-provider-azurerm/issues/18056))
* `azurerm_sentinel_alert_rule_scheduled` - support for the `techniques` property ([#18430](https://github.com/hashicorp/terraform-provider-azurerm/issues/18430))
* `azurerm_linux_virtual_machine` - support for the `patch_assessment_mode` property ([#18437](https://github.com/hashicorp/terraform-provider-azurerm/issues/18437))
* `azurerm_managed_disk` - support for the `PremiumV2_LRS` type ([#17671](https://github.com/hashicorp/terraform-provider-azurerm/issues/17671))
* `azurerm_orchestrated_virtual_machine_scale_set` - support for the `user_data_base64` property ([#18486](https://github.com/hashicorp/terraform-provider-azurerm/issues/18486))
* `azurerm_private_endpoint` - support for the `custom_network_interface_name` property ([#18025](https://github.com/hashicorp/terraform-provider-azurerm/issues/18025))
* `azurerm_virtual_machine_extension` - support for the `failure_suppression_enabled` property ([#18441](https://github.com/hashicorp/terraform-provider-azurerm/issues/18441))
* `azurerm_virtual_machine_scale_set_extension` - support for the `failure_suppression_enabled` property ([#18441](https://github.com/hashicorp/terraform-provider-azurerm/issues/18441))
* `azurerm_windows_virtual_machine` - support for the `patch_assessment_mode` property ([#18437](https://github.com/hashicorp/terraform-provider-azurerm/issues/18437))

BUG FIXES:

* `azurerm_monitor_metric_alert` - pass multi criteria to the API in the correct order ([#18438](https://github.com/hashicorp/terraform-provider-azurerm/issues/18438))
* `azurerm_monitor_diagnostic_settings` - correctly parsing the case for the `workspace_id` property ([#18467](https://github.com/hashicorp/terraform-provider-azurerm/issues/18467))
* `azurerm_security_center_workspace` - correctly parsing the case for the `workspace_id` property ([#18467](https://github.com/hashicorp/terraform-provider-azurerm/issues/18467))

## 3.23.0 (September 15, 2022)

FEATURES:

* **New Data Source**: `azurerm_private_dns_zone_virtual_network_link` ([#18045](https://github.com/hashicorp/terraform-provider-azurerm/issues/18045))
* **New Data Source**: `azurerm_monitor_data_collection_rule` ([#18318](https://github.com/hashicorp/terraform-provider-azurerm/issues/18318))

ENHANCEMENTS:

* `azurerm_api_management_api_schema` - support for the `components` and `definitions` properties ([#18394](https://github.com/hashicorp/terraform-provider-azurerm/issues/18394))
* `azurerm_automation_account` - support for the `hybrid_service_url` property ([#18320](https://github.com/hashicorp/terraform-provider-azurerm/issues/18320))
* `azurerm_batch_pool` - support for the `user_assigned_identity_id` property ([#17104](https://github.com/hashicorp/terraform-provider-azurerm/issues/17104))
* `azurerm_batch_pool` - support for the `data_disks`, `disk_encryption`, `extensions`, `node_placement`, `task_scheduling_policy`, `user_accounts`, and `windows` blocks ([#18226](https://github.com/hashicorp/terraform-provider-azurerm/issues/18226))
* `azurerm_cosmosdb_account` - support for  User Assigned identities ([#18378](https://github.com/hashicorp/terraform-provider-azurerm/issues/18378))
* `azurerm_eventhub_namespace` - support for the `public_network_access_enabled` property ([#18314](https://github.com/hashicorp/terraform-provider-azurerm/issues/18314))
* `azurerm_logic_app_standard` - support for the `virtual_network_subnet_id` property for vNet integration ([#17731](https://github.com/hashicorp/terraform-provider-azurerm/issues/17731))
* `azurerm_management_group_policy_remediation` - the `policy_definition_id` property has been deprecated in favour of the more accuractly named `policy_definition_reference_id` property ([#18037](https://github.com/hashicorp/terraform-provider-azurerm/issues/18037))
* `azurerm_resource_policy_remediation` - the `policy_definition_id` property has been deprecated in favour of the more accuractly named `policy_definition_reference_id` property ([#18037](https://github.com/hashicorp/terraform-provider-azurerm/issues/18037))
* `azurerm_resource_group_policy_remediation` - the `policy_definition_id` property has been deprecated in favour of the more accuractly named `policy_definition_reference_id` property ([#18037](https://github.com/hashicorp/terraform-provider-azurerm/issues/18037))
* `azurerm_subscription_policy_remediation` - the `policy_definition_id` property has been deprecated in favour of the more accuractly named `policy_definition_reference_id` property ([#18037](https://github.com/hashicorp/terraform-provider-azurerm/issues/18037))

BUG FIXES:

* `azurerm_netapp_volume`: add extra validation when `data_protection_snapshot_policy.0. snapshot_policy_id` is empty ([#18348](https://github.com/hashicorp/terraform-provider-azurerm/issues/18348))
 
## 3.22.0 (September 08, 2022)

FEATURES:

* **New Resource**: `azurerm_api_management_api_tag_description` ([#17876](https://github.com/hashicorp/terraform-provider-azurerm/issues/17876))
* **New Resource**: `azurerm_api_management_schema` ([#18158](https://github.com/hashicorp/terraform-provider-azurerm/issues/18158))
* **New Resource**: `azurerm_automation_watcher` ([#17927](https://github.com/hashicorp/terraform-provider-azurerm/issues/17927))
* **New Resource**: `azurerm_automation_source_control` ([#18175](https://github.com/hashicorp/terraform-provider-azurerm/issues/18175))
* **New Resource**: `azurerm_container_registry_token_password` ([#15939](https://github.com/hashicorp/terraform-provider-azurerm/issues/15939))
* **New Resource**: `azurerm_monitor_data_collection_rule_association` ([#17948](https://github.com/hashicorp/terraform-provider-azurerm/issues/17948))
* **New Resource**: `azurerm_orbital_spacecraft` ([#17860](https://github.com/hashicorp/terraform-provider-azurerm/issues/17860))

ENHANCEMENTS:

* dependencies: updating to version `v0.20220907.1111434` of `github.com/hashicorp/go-azure-sdk` ([#18282](https://github.com/hashicorp/terraform-provider-azurerm/issues/18282))
* dependencies: `desktopvirtualization` - updating to use `2022-02-10` ([#17489](https://github.com/hashicorp/terraform-provider-azurerm/issues/17489))
* dependencies: `iothub.dps` - update to use `hashicorp/go-azure-sdk` ([#18299](https://github.com/hashicorp/terraform-provider-azurerm/issues/18299))
* `azurerm_api_management_api` - the `soap_pass_through` property has been deprecated in favour of the `api_type` property ([#17812](https://github.com/hashicorp/terraform-provider-azurerm/issues/17812))
* `azurerm_kubernetes_cluster` - support for the `edge_zone` property ([#18115](https://github.com/hashicorp/terraform-provider-azurerm/issues/18115))
* `azurerm_kubernetes_cluster` - support for the `windows_profile.gmsa` block ([#16437](https://github.com/hashicorp/terraform-provider-azurerm/issues/16437))
* `azurerm_mssql_database` - support for the `maintenance_configuration_name` property ([#18247](https://github.com/hashicorp/terraform-provider-azurerm/issues/18247))
* `azurerm_virtual_desktop_host_pool` - support for the `scheduled_agent_updates` block ([#17489](https://github.com/hashicorp/terraform-provider-azurerm/issues/17489))
* `azurerm_hdinsight_kafka_cluster` - support for the `extension` property ([#17846](https://github.com/hashicorp/terraform-provider-azurerm/issues/17846))
* `azurerm_hdinsight_spark_cluster` - support for the `extension` property ([#17846](https://github.com/hashicorp/terraform-provider-azurerm/issues/17846))
* `azurerm_hdinsight_interactive_query_cluster` - support for the `extension` property ([#17846](https://github.com/hashicorp/terraform-provider-azurerm/issues/17846))
* `azurerm_hdinsight_hbase_cluster` - support for the `extension` property ([#17846](https://github.com/hashicorp/terraform-provider-azurerm/issues/17846))
* `azurerm_hdinsight_hadoop_cluster` - support for the `extension` property ([#17846](https://github.com/hashicorp/terraform-provider-azurerm/issues/17846))

BUG FIXES:

* `azurerm_mssql_database` - the `license_type` property is now also Computed ([#18230](https://github.com/hashicorp/terraform-provider-azurerm/issues/18230))
* `azurerm_log_analytics_solution` - a state migration to work around the previously incorrect id casing ([#18291](https://github.com/hashicorp/terraform-provider-azurerm/issues/18291))

## 3.21.1 (September 02, 2022)

BREAKING CHANGES:

* `azurerm_container_registry` - the field `azuread_authentication_as_arm_policy_enabled` has been removed to fix a regression - support for this will be reintroduced in a future release.
* `azurerm_container_registry` - the field `soft_delete_policy` has been removed to fix a regression - support for this will be reintroduced in a future release.

NOTES:

* the `containerregistry` api version has been reverted to `2021-08-01-preview` to restore the `virtual_network` block meaning the `azuread_authentication_as_arm_policy_enabled` and `soft_delete_policy` properties had to be removed as they were not supported by the API version that supported virtual network rules. ([#18239](https://github.com/hashicorp/terraform-provider-azurerm/issues/18239))

BUG FIXES:

* `azurerm_container_registry` - the `virtual_network` block has been restored ([#18239](https://github.com/hashicorp/terraform-provider-azurerm/issues/18239))
* `azurerm_log_analytics_data_export_rule` - a state migration to work around the previously incorrect id casing ([#18240](https://github.com/hashicorp/terraform-provider-azurerm/issues/18240))

## 3.21.0 (September 01, 2022)

FEATURES:

* **New Data Source**: `azurerm_monitor_data_collection_endpoint` ([#17992](https://github.com/hashicorp/terraform-provider-azurerm/issues/17992))
* **New Resource**: `azurerm_app_service_connection` ([#16907](https://github.com/hashicorp/terraform-provider-azurerm/issues/16907))
* **New Resource**: `azurerm_automation_hybrid_runbook_worker` ([#17893](https://github.com/hashicorp/terraform-provider-azurerm/issues/17893))
* **New Resource**: `azurerm_api_management_gateway_certificate_authority` ([#17879](https://github.com/hashicorp/terraform-provider-azurerm/issues/17879))
* **New Resource**: `azurerm_api_management_gateway_host_name_configuration` ([#17962](https://github.com/hashicorp/terraform-provider-azurerm/issues/17962))
* **New Resource**: `azurerm_api_management_product_tag` ([#17798](https://github.com/hashicorp/terraform-provider-azurerm/issues/17798))
* **New Resource**: `azurerm_automation_connection_type` ([#17538](https://github.com/hashicorp/terraform-provider-azurerm/issues/17538))
* **New Resource**: `azurerm_automation_hybrid_runbook_worker_group` ([#17881](https://github.com/hashicorp/terraform-provider-azurerm/issues/17881))
* **New Resource:** `azurerm_cdn_frontdoor_rule` ([#18010](https://github.com/hashicorp/terraform-provider-azurerm/issues/18010))
* **New Resource:** `azurerm_cdn_frontdoor_secret` ([#18010](https://github.com/hashicorp/terraform-provider-azurerm/issues/18010))
* **New Resource**: `azurerm_container_registry_task_schedule_run_now` ([#15120](https://github.com/hashicorp/terraform-provider-azurerm/issues/15120))
* **New Resource**: `azurerm_cosmosdb_sql_dedicated_gateway` ([#18133](https://github.com/hashicorp/terraform-provider-azurerm/issues/18133))
* **New Resource**: `azurerm_dashboard_grafana` ([#17840](https://github.com/hashicorp/terraform-provider-azurerm/issues/17840))
* **New Resource**: `azurerm_healthcare_medtech_service` ([#15967](https://github.com/hashicorp/terraform-provider-azurerm/issues/15967))
* **New Resource**: `azurerm_log_analytics_query_pack_query` ([#17929](https://github.com/hashicorp/terraform-provider-azurerm/issues/17929))
* **New Resource**: `azurerm_spring_cloud_connection` ([#16907](https://github.com/hashicorp/terraform-provider-azurerm/issues/16907))
* **New Resource**: `azurerm_search_shared_private_link_service` ([#17744](https://github.com/hashicorp/terraform-provider-azurerm/issues/17744))
* **New Resource**: `azurerm_sentinel_alert_rule_nrt` ([#15999](https://github.com/hashicorp/terraform-provider-azurerm/issues/15999))

ENHANCEMENTS:

* dependencies: updating to version `v0.20220830.1105041` of `github.com/hashicorp/go-azure-sdk` ([#18183](https://github.com/hashicorp/terraform-provider-azurerm/issues/18183))
* dependencies: `log_analytics` - update to use `hashicorp/go-azure-sdk` ([#18098](https://github.com/hashicorp/terraform-provider-azurerm/issues/18098))
* `azurerm_batch_pool` - support for the `mount` property ([#18042](https://github.com/hashicorp/terraform-provider-azurerm/issues/18042))
* `azurerm_container_registry` - support for the `azuread_authentication_as_arm_policy_enabled` and `soft_delete_policy` properties ([#17926](https://github.com/hashicorp/terraform-provider-azurerm/issues/17926))
* `azurerm_cosmosdb_cassandra_cluster` - support for the `HoursBetweenBackups` property ([#18154](https://github.com/hashicorp/terraform-provider-azurerm/issues/18154))
* `azurerm_hdinsight_kafka_cluster` - add support for the `disk_encryption` property ([#17351](https://github.com/hashicorp/terraform-provider-azurerm/issues/17351))
* `azurerm_hdinsight_spark_cluster` - add support for the `disk_encryption` property ([#17351](https://github.com/hashicorp/terraform-provider-azurerm/issues/17351))
* `azurerm_hdinsight_interactive_query_cluster` - add support for the `disk_encryption` property ([#17351](https://github.com/hashicorp/terraform-provider-azurerm/issues/17351))
* `azurerm_hdinsight_hbase_cluster` - add support for the `disk_encryption` property ([#17351](https://github.com/hashicorp/terraform-provider-azurerm/issues/17351))
* `azurerm_hdinsight_hadoop_cluster` - add support for the `disk_encryption` property ([#17351](https://github.com/hashicorp/terraform-provider-azurerm/issues/17351))
* `azurerm_iothub_dps` - support for the `resource_count`, `parallel_deployments`, and `failure_percentage` properties ([#18151](https://github.com/hashicorp/terraform-provider-azurerm/issues/18151))
* `azurerm_kubernetes_node_pool` - spot node pools can now be upgraded ([#18124](https://github.com/hashicorp/terraform-provider-azurerm/issues/18124))
* `azurerm_linux_virtual_machine` - the `source_image_id` property now supports both `Community Gallery Images`, and `Shared Gallery Images` resource IDs ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - support for the following properties `host_group_id`, and `extension_operations_enabled` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - support for the `data_disk` block property `name` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - support for the `scale_in` block properties `rule`, and `force_deletion_enabled` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - support for the `rolling_upgrade_policy` block properties `cross_zone_upgrade_enabled`, and `prioritize_unhealthy_instances_enabled`  ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - added support for the `spot_restore` block ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - support for the `spot_restore` block properties `enabled`, and `timeout` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - support for the `public_ip_address` block property `version` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - the `source_image_id` property now supports both `Community Gallery Images`, and `Shared Gallery Images` resource IDs ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - support for the `gallery_applications` code block ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - support for the `gallery_applications` block properties `configuration_reference_blob_uri`, `order`, `package_referenceId`, and `tag` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - deprecated the `scale_in_policy` property in favour of the `scale_in` block due to additional fields being added ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_linux_virtual_machine_scale_set` - support for the `scale_in` block property `rule` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_management_group_policy_remediation` - support for the `resource_count`, `parallel_deployments`, and `failure_percentage` properties ([#17313](https://github.com/hashicorp/terraform-provider-azurerm/issues/17313))
* `azurerm_monitor_diagnostic_setting` - support for the `category_group` property ([#16367](https://github.com/hashicorp/terraform-provider-azurerm/issues/16367))
* `azurerm_orchestrated_virtual_machine_scale_set` - support for the following properties `capacity_reservation_group_id`, `single_placement_group`, and `extension_operations_enabled` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_orchestrated_virtual_machine_scale_set` - support for the `extension` block property `suppress_failures_enabled` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_orchestrated_virtual_machine_scale_set` - support for the `additional_capabilities` block property `ultra_ssd_enabled` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_orchestrated_virtual_machine_scale_set` - support for the `public_ip_address` block properties `version`, and `sku_name` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_orchestrated_virtual_machine_scale_set` - support for `linux_configuration`, and `windows_configuration` code blocks property `patch_assessment_mode` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_orchestrated_virtual_machine_scale_set` - the `source_image_id` property now supports both `Community Gallery Images`, and `Shared Gallery Images` resource IDs ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_policy_definition - export the `role_definition_ids` attribute ([#18043](https://github.com/hashicorp/terraform-provider-azurerm/issues/18043))
* `azurerm_resource_group_policy_remediation` - support for the `resource_count`, `parallel_deployments`, and `failure_percentage` properties ([#17313](https://github.com/hashicorp/terraform-provider-azurerm/issues/17313))
* `azurerm_resource_policy_remediation` - support for the `resource_count`, `parallel_deployments`, and `failure_percentage` properties ([#17313](https://github.com/hashicorp/terraform-provider-azurerm/issues/17313))
* `azurerm_role_assignment` - support for `scope` to start with `/providers/Subscription` ([#17456](https://github.com/hashicorp/terraform-provider-azurerm/issues/17456))
* `azurerm_servicebus_namespace` - support for the `public_network_access_enabled` and `minimum_tls_version` properties ([#17805](https://github.com/hashicorp/terraform-provider-azurerm/issues/17805))
* `azurerm_storage_account` - support for the `public_network_access_enabled` property ([#18005](https://github.com/hashicorp/terraform-provider-azurerm/issues/18005))
* `azurerm_stream_analytics_output_eventhub` - support for the `authentication_mode` property ([#18096](https://github.com/hashicorp/terraform-provider-azurerm/issues/18096))
* `azurerm_stream_analytics_output_mssql` - support for the `authentication_mode` property ([#18096](https://github.com/hashicorp/terraform-provider-azurerm/issues/18096))
* `azurerm_stream_analytics_output_servicebus_topic` - support for the `authentication_mode` property ([#18096](https://github.com/hashicorp/terraform-provider-azurerm/issues/18096))
* `azurerm_stream_analytics_output_powerbi` - support for the `token_user_principal_name` and `token_user_display_name` properties ([#18117](https://github.com/hashicorp/terraform-provider-azurerm/issues/18117))
* `azurerm_stream_analytics_output_cosmosdb` - support for the `partition_key` property ([#18120](https://github.com/hashicorp/terraform-provider-azurerm/issues/18120))
* `azurerm_stream_analytics_reference_input_blob` - support for the `authentication_mode` property ([#18137](https://github.com/hashicorp/terraform-provider-azurerm/issues/18137))
* `azurerm_stream_analytics_reference_input_mssql` - support for the `table` property ([#18211](https://github.com/hashicorp/terraform-provider-azurerm/issues/18211))
* `azurerm_subscription_policy_remediation` - support for the `resource_count`, `parallel_deployments`, and `failure_percentage` properties ([#17313](https://github.com/hashicorp/terraform-provider-azurerm/issues/17313))
* `azurerm_windows_virtual_machine` - the `source_image_id` property now supports both `Community Gallery Images`, and `Shared Gallery Images` resource IDs ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - support for the following properties `host_group_id`, and `extension_operations_enabled` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - support for the `data_disk` block property `name` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - support for the `scale_in` block properties `rule`, and `force_deletion_enabled` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - support for the `rolling_upgrade_policy` block properties `cross_zone_upgrade_enabled`, and `prioritize_unhealthy_instances_enabled`  ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - added support for the `spot_restore` block ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - support for the `spot_restore` block properties `enabled`, and `timeout` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - support for the `public_ip_address` block property `version` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - the `source_image_id` property now supports both `Community Gallery Images`, and `Shared Gallery Images` resource IDs ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - support for the `gallery_applications` code block ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - support for the `gallery_applications` block properties `configuration_reference_blob_uri`, `order`, `package_referenceId`, and `tag` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - deprecated the `scale_in_policy` property in favour of the `scale_in` block due to additional fields being added ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))
* `azurerm_windows_virtual_machine_scale_set` - support for the `scale_in` block property `rule` ([#17571](https://github.com/hashicorp/terraform-provider-azurerm/issues/17571))

BUG FIXES:

* `azurerm_kubernetes_cluster` - `kube_config` is now set when AAD is enabled for a `v1.24` cluster ([#18142](https://github.com/hashicorp/terraform-provider-azurerm/issues/18142))
* `azurerm_redis_cache` - will now recreate the cache when downgrading the SKU ([#17767](https://github.com/hashicorp/terraform-provider-azurerm/issues/17767))
* `azurerm_spring_cloud_service` - ignore the default zero value for `read_timeout_seconds` ([#18161](https://github.com/hashicorp/terraform-provider-azurerm/issues/18161))

## 3.20.0 (August 25, 2022)

FEATURES:

* **Provider:** support for generic OIDC authentication providers ([#18118](https://github.com/hashicorp/terraform-provider-azurerm/issues/18118))
* **New Resource**: `azurerm_backup_policy_vm_workload` ([#17765](https://github.com/hashicorp/terraform-provider-azurerm/issues/17765))
* **New Resource**: `azurerm_monitor_scheduled_query_rules_alert_v2` ([#17772](https://github.com/hashicorp/terraform-provider-azurerm/issues/17772))

ENHANCEMENTS:

* Dependencies: update `go-azure-sdk` to `v0.20220824.1090858` ([#18100](https://github.com/hashicorp/terraform-provider-azurerm/issues/18100))
* Dependencies: `consumption` - updating to use `hashicorp/go-azure-sdk` ([#18101](https://github.com/hashicorp/terraform-provider-azurerm/issues/18101))
* `azurerm_data_factory_dataset_json` - `filename` and `path` in `azure_blob_storage_location` block can now be empty ([#18061](https://github.com/hashicorp/terraform-provider-azurerm/issues/18061))

BUG FIXES:

* `data.azurerm_kubernetes_cluster` - `kube_config` is now set when AAD is enabled for a v1.24 cluster ([#18131](https://github.com/hashicorp/terraform-provider-azurerm/issues/18131))
* `azurerm_cosmosdb_sql_database` - prevent panic in autoacale settings ([#18070](https://github.com/hashicorp/terraform-provider-azurerm/issues/18070))
* `azurerm_kubernetes_cluster_node_pool` - fix a crash in expanding upgrade settings ([#18074](https://github.com/hashicorp/terraform-provider-azurerm/issues/18074))
* `azurerm_mssql_elastic_pool` - list of values for `maintenance_configuration_name` is now correct ([#18041](https://github.com/hashicorp/terraform-provider-azurerm/issues/18041))
* `azurerm_postgresql_flexible_server` - `point_in_time_restore_time_in_utc` correctly converts to RFC3339 ([#18106](https://github.com/hashicorp/terraform-provider-azurerm/issues/18106))

## 3.19.1 (August 19, 2022)

BUG FIXES:

* `azurerm_dns_a_record` - parse resource IDs insensitively in the read functions due to casing on the dnsZones segment ([#18048](https://github.com/hashicorp/terraform-provider-azurerm/issues/18048))
* `azurerm_dns_aaaa_record` - parse resource IDs insensitively in the read functions due to casing on the dnsZones segment ([#18048](https://github.com/hashicorp/terraform-provider-azurerm/issues/18048))
* `azurerm_dns_caa_record` - parse resource IDs insensitively in the read functions due to casing on the dnsZones segment ([#18048](https://github.com/hashicorp/terraform-provider-azurerm/issues/18048))
* `azurerm_dns_cname_record` - parse resource IDs insensitively in the read functions due to casing on the dnsZones segment ([#18048](https://github.com/hashicorp/terraform-provider-azurerm/issues/18048))
* `azurerm_dns_mx_record` - parse resource IDs insensitively in the read functions due to casing on the dnsZones segment ([#18048](https://github.com/hashicorp/terraform-provider-azurerm/issues/18048))
* `azurerm_dns_ns_record` - parse resource IDs insensitively in the read functions due to casing on the dnsZones segment ([#18048](https://github.com/hashicorp/terraform-provider-azurerm/issues/18048))
* `azurerm_dns_ptr_record` - parse resource IDs insensitively in the read functions due to casing on the dnsZones segment ([#18048](https://github.com/hashicorp/terraform-provider-azurerm/issues/18048))
* `azurerm_dns_srv_record` - parse resource IDs insensitively in the read functions due to casing on the dnsZones segment ([#18048](https://github.com/hashicorp/terraform-provider-azurerm/issues/18048))
* `azurerm_dns_txt_record` - parse resource IDs insensitively in the read functions due to casing on the dnsZones segment ([#18048](https://github.com/hashicorp/terraform-provider-azurerm/issues/18048))
* `azurerm_dns_zone` - parse resource IDs insensitively in the read functions due to casing on the dnsZones segment ([#18048](https://github.com/hashicorp/terraform-provider-azurerm/issues/18048))

## 3.19.0 (August 18, 2022)

FEATURES:

* **New Data Source**: `azurerm_dns_a_record` ([#17477](https://github.com/hashicorp/terraform-provider-azurerm/issues/17477))
* **New Data Source**: `azurerm_dns_aaaa_record` ([#17477](https://github.com/hashicorp/terraform-provider-azurerm/issues/17477))
* **New Data Source**: `azurerm_dns_caa_record` ([#17477](https://github.com/hashicorp/terraform-provider-azurerm/issues/17477))
* **New Data Source**: `azurerm_dns_cname_record` ([#17477](https://github.com/hashicorp/terraform-provider-azurerm/issues/17477))
* **New Data Source**: `azurerm_dns_mx_record` ([#17477](https://github.com/hashicorp/terraform-provider-azurerm/issues/17477))
* **New Data Source**: `azurerm_dns_ns_record` ([#17477](https://github.com/hashicorp/terraform-provider-azurerm/issues/17477))
* **New Data Source**: `azurerm_dns_ptr_record` ([#17477](https://github.com/hashicorp/terraform-provider-azurerm/issues/17477))
* **New Data Source**: `azurerm_dns_soa_record` ([#17477](https://github.com/hashicorp/terraform-provider-azurerm/issues/17477))
* **New Data Source**: `azurerm_dns_srv_record` ([#17477](https://github.com/hashicorp/terraform-provider-azurerm/issues/17477))
* **New Data Source**: `azurerm_dns_txt_record` ([#17477](https://github.com/hashicorp/terraform-provider-azurerm/issues/17477))
* **New Data Source**: `azurerm_private_dns_a_record` ([#18036](https://github.com/hashicorp/terraform-provider-azurerm/issues/18036))
* **New Data Source**: `azurerm_private_dns_aaaa_record` ([#18036](https://github.com/hashicorp/terraform-provider-azurerm/issues/18036))
* **New Data Source**: `azurerm_private_dns_cname_record` ([#18036](https://github.com/hashicorp/terraform-provider-azurerm/issues/18036))
* **New Data Source**: `azurerm_private_dns_mx_record` ([#18036](https://github.com/hashicorp/terraform-provider-azurerm/issues/18036))
* **New Data Source**: `azurerm_private_dns_ptr_record` ([#18036](https://github.com/hashicorp/terraform-provider-azurerm/issues/18036))
* **New Data Source**: `azurerm_private_dns_soa_record` ([#18036](https://github.com/hashicorp/terraform-provider-azurerm/issues/18036))
* **New Data Source**: `azurerm_private_dns_srv_record` ([#18036](https://github.com/hashicorp/terraform-provider-azurerm/issues/18036))
* **New Data Source**: `azurerm_private_dns_txt_record` ([#18036](https://github.com/hashicorp/terraform-provider-azurerm/issues/18036))
* **New Resource**: `azurerm_eventhub_namespace_schema_group` ([#17635](https://github.com/hashicorp/terraform-provider-azurerm/issues/17635))
* **New Resource**: `azurerm_cdn_frontdoor_firewall_policy` ([#17715](https://github.com/hashicorp/terraform-provider-azurerm/issues/17715))
* **New Resource**: `azurerm_cdn_frontdoor_security_policy` ([#17715](https://github.com/hashicorp/terraform-provider-azurerm/issues/17715))
* **New Resource**: `azurerm_data_factory_flowlet_data_flow` ([#16987](https://github.com/hashicorp/terraform-provider-azurerm/issues/16987))

ENHANCEMENTS:

* Dependencies: update `go-azure-helpers` to `v0.39.1` ([#18015](https://github.com/hashicorp/terraform-provider-azurerm/issues/18015))
* Dependencies: update `go-azure-sdk` to `v0.20220815.1092453` ([#17998](https://github.com/hashicorp/terraform-provider-azurerm/issues/17998))
* Dependencies: `dedicated_host_*` to use `hashicorp/go-azure-sdk` ([#17616](https://github.com/hashicorp/terraform-provider-azurerm/issues/17616))
* Dependencies: `dataprotection`: updating to use `hashicorp/go-azure-sdk` ([#17700](https://github.com/hashicorp/terraform-provider-azurerm/issues/17700))
* Dependencies: `dns` - updating to use `hashicorp/go-azure-sdk` ([#17986](https://github.com/hashicorp/terraform-provider-azurerm/issues/17986))
* Dependencies: `maintenance` - updating to use `hashicorp/go-azure-sdk` ([#17954](https://github.com/hashicorp/terraform-provider-azurerm/issues/17954))
* Data Source: `azurerm_images` - now uses a logical id ([#17766](https://github.com/hashicorp/terraform-provider-azurerm/issues/17766))
* Data Source: `azurerm_management_group` - now exports the `management_group_ids`, `all_management_group_ids`, and `all_subscription_ids` attributes ([#16208](https://github.com/hashicorp/terraform-provider-azurerm/issues/16208))
* `azurerm_active_directory_domain_service` - support for the `kerberos_armoring_enabled` and `kerberos_rc4_encryption_enabled` properties ([#17853](https://github.com/hashicorp/terraform-provider-azurerm/issues/17853))
* `azurerm_application_gateway` - support for the `global` block ([#17651](https://github.com/hashicorp/terraform-provider-azurerm/issues/17651))
* `azurerm_application_gateway` - support for `components` in `rewrite_rule_set.rewrite_rule.url` ([#13899](https://github.com/hashicorp/terraform-provider-azurerm/issues/13899))
* `azurerm_automation_account` - support for the `private_endpoint_connection` property ([#17934](https://github.com/hashicorp/terraform-provider-azurerm/issues/17934))
* `azurerm_automation_account` - support for the `encryption` block and `local_authentication_enabled` property ([#17454](https://github.com/hashicorp/terraform-provider-azurerm/issues/17454))
* `azurerm_batch_account` - support for the `storage_account_authentication_mode`, `storage_account_node_identit`, and `allowed_authentication_modes` properties ([#16758](https://github.com/hashicorp/terraform-provider-azurerm/issues/16758))
* `azurerm_batch_pool` - support for identity referencees in container registries ([#17416](https://github.com/hashicorp/terraform-provider-azurerm/issues/17416))
* `azurerm_data_factory_data_flow` - support for the `flowlet` block ([#16987](https://github.com/hashicorp/terraform-provider-azurerm/issues/16987))
* `azurerm_data_factory_integration_runtime_azure_ssis` - support for the `express_vnet_injection` property ([#17756](https://github.com/hashicorp/terraform-provider-azurerm/issues/17756))
* `azurerm_firewall_policy_resource` - support for the `private_ranges` and `allow_sql_redirect` properties ([#17842](https://github.com/hashicorp/terraform-provider-azurerm/issues/17842))
* `azurerm_key_vault` - support for the `public_network_access_enabled` property ([#17552](https://github.com/hashicorp/terraform-provider-azurerm/issues/17552))
* `azurerm_linux_virtual_machine` - now supports delete Eviction policies ([#17226](https://github.com/hashicorp/terraform-provider-azurerm/issues/17226))
* `azurerm_linux_virtual_machine_scale_set` - now supports delete Eviction policies ([#17226](https://github.com/hashicorp/terraform-provider-azurerm/issues/17226))
* `azurerm_mssql_elastic_pool` - support for the `maintenance_configuration_name` property ([#17790](https://github.com/hashicorp/terraform-provider-azurerm/issues/17790))
* `azurerm_mssql_server` - support `Disabled` for the `minimum_tls_version` property ([#16595](https://github.com/hashicorp/terraform-provider-azurerm/issues/16595))
* `azurerm_spring_cloud_app` - support the `public_endpoint_enabled` property ([#17630](https://github.com/hashicorp/terraform-provider-azurerm/issues/17630))
* `azurerm_spring_cloud_gateway_route_config` - support for the `open_api;azurerm_spring_cloud_service`  and `log_stream_public_endpoint_enabledread_timeout_seconds` properties ([#17630](https://github.com/hashicorp/terraform-provider-azurerm/issues/17630))
* `azurerm_shared_image` - support for the `architecture` property ([#17250](https://github.com/hashicorp/terraform-provider-azurerm/issues/17250))
* `azurerm_storage_account` - support for the `default_to_oauth_authentication` property ([#17116](https://github.com/hashicorp/terraform-provider-azurerm/issues/17116))
* `azurerm_storage_table_entity` - support for specifying data types on entity properties ([#15782](https://github.com/hashicorp/terraform-provider-azurerm/issues/15782))
* `azurerm_shared_image_version` - support for `blob_uri` and `storage_account_id` ([#17768](https://github.com/hashicorp/terraform-provider-azurerm/issues/17768))
* `azurerm_windows_virtual_machine` - now supports delete Eviction policies ([#17226](https://github.com/hashicorp/terraform-provider-azurerm/issues/17226))
* `azurerm_windows_virtual_machine_scale_set` - now supports delete Eviction policies ([#17226](https://github.com/hashicorp/terraform-provider-azurerm/issues/17226))
* `azurerm_web_application_firewall_policy` - support for the `excluded_rule_set` property ([#17757](https://github.com/hashicorp/terraform-provider-azurerm/issues/17757))
* `azurerm_log_analytics_workspace` - support for the `cmk_for_query_forced` property ([#17365](https://github.com/hashicorp/terraform-provider-azurerm/issues/17365))
* `azurerm_lb_backend_address_pool_address` - support for the `backend_address_ip_configuration_id` property ([#17770](https://github.com/hashicorp/terraform-provider-azurerm/issues/17770))

BUG FIXES:

* Data Source: `azurerm_windows_web_app` - add missing schema definition for 'virtual_network_subnet_id' ([#18028](https://github.com/hashicorp/terraform-provider-azurerm/issues/18028))
* `azurerm_cdn_endpoint_custom_domain` - deprecating the `key_vault_certificate_id` property in favour of the `key_vault_secret_id` property withing the `user_managed https_allows` block ([#17114](https://github.com/hashicorp/terraform-provider-azurerm/issues/17114))
* `azurerm_data_protection_backup_policy_postgresql_resource` - prevent a crash when given an empty criteria block ([#17904](https://github.com/hashicorp/terraform-provider-azurerm/issues/17904))
* `azurerm_disk_encryption_set` - prevent an issue during creation when the disk encryption set and key vault are in different subscriptions ([#17964](https://github.com/hashicorp/terraform-provider-azurerm/issues/17964))
* `azurerm_windows_function_app` fix a bug with setting values for `WindowsFxString` ([#18014](https://github.com/hashicorp/terraform-provider-azurerm/issues/18014))
* `azurerm_windows_function_app_slot` - fix a bug with setting values for `WindowsFxString` ([#18014](https://github.com/hashicorp/terraform-provider-azurerm/issues/18014))
* `azurerm_linux_function_app` - correctly send `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_linux_function_app` - fix content settings when `storage_uses_managed_identity` is set to `true` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_linux_function_app_slot` - correctly send `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_linux_function_app_slot` - fix content settings when `storage_uses_managed_identity` is set to `true` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_windows_function_app` - correctly send `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_windows_function_app` - fix content settings when `storage_uses_managed_identity` is set to `true` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_windows_function_app_slot` - correctly send `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_windows_function_app_slot` - fix content settings when `storage_uses_managed_identity` is set to `true` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))

## 3.18.0 (August 11, 2022)

FEATURES: 

* **New Resource**: `azurerm_monitor_data_collection_endpoint` ([#17684](https://github.com/hashicorp/terraform-provider-azurerm/issues/17684))

ENHANCEMENTS:

* dependencies: updating `github.com/hashicorp/go-azure-sdk` to `v0.20220809.1122626` ([#17905](https://github.com/hashicorp/terraform-provider-azurerm/issues/17905))
* storage: updating to use API Version `2021-09-01` ([#17523](https://github.com/hashicorp/terraform-provider-azurerm/issues/17523))
* `azurerm_express_route_circuit_peering` - support for the `ipv4_enabled` and `gateway_manager_etag` properties ([#17338](https://github.com/hashicorp/terraform-provider-azurerm/issues/17338))
* `azurerm_site_recovery_replicated_vm` - support for the `target_disk_encryption` block ([#15783](https://github.com/hashicorp/terraform-provider-azurerm/issues/15783))
* `azurerm_subnet` - deprecate `enforce_private_link_endpoint_network_policies` property in favour of `private_endpoint_network_policies_enabled` ([#17464](https://github.com/hashicorp/terraform-provider-azurerm/issues/17464))
* `azurerm_subnet` - deprecate `enforce_private_link_service_network_policies` property in favour of `private_link_service_network_policies_enabled` ([#17464](https://github.com/hashicorp/terraform-provider-azurerm/issues/17464))
* `azurerm_servicebus_subscription` - support for the `client_scoped_subscription_enabled` property and the `client_scoped_subscription` block ([#17101](https://github.com/hashicorp/terraform-provider-azurerm/issues/17101))

BUG FIXES:

* `azurerm_backup_policy_vm` - now prevents crash when `frequency` is set to Hourly and, `hour_interval` and `hour_duration`are not set ([#17880](https://github.com/hashicorp/terraform-provider-azurerm/issues/17880))
* Data Source: `azurerm_blueprint_definition` - Fix `version` property output ([#16299](https://github.com/hashicorp/terraform-provider-azurerm/issues/16299))

## 3.17.0 (August 04, 2022)

ENHANCEMENTS:

* domainservice: updating to use API Version `2021-05-01` ([#17737](https://github.com/hashicorp/terraform-provider-azurerm/issues/17737))
* Data Source: `azurerm_proximity_placement_group` - refactoring to use `hashicorp/go-azure-sdk` ([#17776](https://github.com/hashicorp/terraform-provider-azurerm/issues/17776))
* `azurerm_api_management` - update the `sku_name` property validation to accept newer Premium SKUs ([#17887](https://github.com/hashicorp/terraform-provider-azurerm/issues/17887))
* `azurerm_firewall` - the property `sku_tier` is now updateable ([#17577](https://github.com/hashicorp/terraform-provider-azurerm/issues/17577))
* `azurerm_linux_virtual_machine_scale_set` - the property `instances` is now Optional and defaults to `0` ([#17836](https://github.com/hashicorp/terraform-provider-azurerm/issues/17836))
* `azurerm_log_analytics_cluster` - updated validation for the `size_gb` property ([#17780](https://github.com/hashicorp/terraform-provider-azurerm/issues/17780))
* `azurerm_proximity_placement_group` - refactoring to use `hashicorp/go-azure-sdk` ([#17776](https://github.com/hashicorp/terraform-provider-azurerm/issues/17776))
* `azurerm_shared_image` - improved validation for the `publisher`, `offer` and `sku` properties in the `identifier` block ([#17547](https://github.com/hashicorp/terraform-provider-azurerm/issues/17547))
* `azurerm_subnet` - support for the service delegation `Microsoft.Orbital/orbitalGateway` ([#17854](https://github.com/hashicorp/terraform-provider-azurerm/issues/17854))
* `azurerm_eventhub_namespace` - support for the `local_authentication_enabled`, `public_network_access_enabled`, and `minimum_tls_version` properties ([#17194](https://github.com/hashicorp/terraform-provider-azurerm/issues/17194))

BUG FIXES:

* Data Source: `azurerm_private_dns_zone` - returning the correct Resource ID when not specifying the `resource_group_name` ([#17729](https://github.com/hashicorp/terraform-provider-azurerm/issues/17729))

## 3.16.0 (July 28, 2022)

FEATURES: 

* **New Resource**: `azurerm_datadog_monitor` ([#16131](https://github.com/hashicorp/terraform-provider-azurerm/issues/16131))
* **New Resource**: `azurerm_kusto_cluster_managed_private_endpoint` ([#17667](https://github.com/hashicorp/terraform-provider-azurerm/issues/17667))
* **New Resource**: `azurerm_log_analytics_query_pack` ([#17685](https://github.com/hashicorp/terraform-provider-azurerm/issues/17685))
* **New Resource**: `azurerm_logz_sub_account_tag_rule` ([#17557](https://github.com/hashicorp/terraform-provider-azurerm/issues/17557))
* **New Resource**: `azurerm_signalr_shared_private_link_resource` ([#16187](https://github.com/hashicorp/terraform-provider-azurerm/issues/16187))

ENHANCEMENTS:

* dependencies: updating to version `v0.20220725.1163004` of `github.com/hashicorp/go-azure-sdk` ([#17753](https://github.com/hashicorp/terraform-provider-azurerm/issues/17753))
* automationaccount: updating to use `hashicorp/go-azure-sdk` ([#17347](https://github.com/hashicorp/terraform-provider-azurerm/issues/17347))
* Data Source: `azurerm_linux_function_app` - support the `virtual_network_subnet_id` property for for vNet integration ([#17494](https://github.com/hashicorp/terraform-provider-azurerm/issues/17494))
* Data Source: `azurerm_windows_function_app` - support the `virtual_network_subnet_id` property for for vNet integration ([#17572](https://github.com/hashicorp/terraform-provider-azurerm/issues/17572))
* Data Source: `azurerm_windows_web_app` - support the `virtual_network_subnet_id` property for for vNet integration ([#17576](https://github.com/hashicorp/terraform-provider-azurerm/issues/17576))
* `eventhub`: updating all data sources/resources onto single API Version `2021-11-01` ([#17719](https://github.com/hashicorp/terraform-provider-azurerm/issues/17719))
* `azurerm_bot_service_azure_bot` - support for the `streaming_endpoint_enabled` property ([#17423](https://github.com/hashicorp/terraform-provider-azurerm/issues/17423))
* `azurerm_cognitive_account` - support for the `custom_question_answering_search_service_key` property ([#17683](https://github.com/hashicorp/terraform-provider-azurerm/issues/17683))
* `asurerm_iothub_dps_certificate` - support for the `is_verified` property ([#17106](https://github.com/hashicorp/terraform-provider-azurerm/issues/17106))
* `azurerm_linux_web_app`  - the `virtual_network_subnet_id` property is no longer `ForceNew` ([#17584](https://github.com/hashicorp/terraform-provider-azurerm/issues/17584))
* `azurerm_linux_web_app_slot` - the `virtual_network_subnet_id` property is no longer `ForceNew` ([#17584](https://github.com/hashicorp/terraform-provider-azurerm/issues/17584))
* `azurerm_linux_function_app` support the `virtual_network_subnet_id` property for for vNet integration ([#17494](https://github.com/hashicorp/terraform-provider-azurerm/issues/17494))
* `azurerm_linux_function_app_slot` support the `virtual_network_subnet_id` property for for vNet integration ([#17494](https://github.com/hashicorp/terraform-provider-azurerm/issues/17494))
* `azurerm_stream_analytics_stream_input_eventhub` - support for the `authentication_mode` property ([#17739](https://github.com/hashicorp/terraform-provider-azurerm/issues/17739))
* `azurerm_windows_function_app` support the `virtual_network_subnet_id` property for for vNet integration ([#17572](https://github.com/hashicorp/terraform-provider-azurerm/issues/17572))
* `azurerm_windows_function_app_slot` support the `virtual_network_subnet_id` property for for vNet integration ([#17572](https://github.com/hashicorp/terraform-provider-azurerm/issues/17572))
* `azurerm_windows_web_app` support the `virtual_network_subnet_id` property for for vNet integration ([#17576](https://github.com/hashicorp/terraform-provider-azurerm/issues/17576))
* `azurerm_windows_web_app_slot` support the `virtual_network_subnet_id` property for for vNet integration ([#17576](https://github.com/hashicorp/terraform-provider-azurerm/issues/17576))

BUG FIXES:

* `azurerm_linux_function_app` - fix casing bug with the `linux_fx_string` property for Node apps ([#17789](https://github.com/hashicorp/terraform-provider-azurerm/issues/17789))
* `azurerm_linux_function_app_slot` - fix casing bug with the `linux_fx_string` property for Node apps ([#17789](https://github.com/hashicorp/terraform-provider-azurerm/issues/17789))
* `azurerm_resource_group_template_deployment` - fixing a bug where the same Resource Provider defined in different casings would cause the API Version to not be identified ([#17707](https://github.com/hashicorp/terraform-provider-azurerm/issues/17707))

## 3.15.1 (July 25, 2022)

BUG FIXES: 

* `data.azurerm_servicebus_queue` - fix a regression around `namespace_id` ([#17755](https://github.com/hashicorp/terraform-provider-azurerm/issues/17755))
* `azurerm_postgresql_aad_administrator` - fix the state migration ([#17732](https://github.com/hashicorp/terraform-provider-azurerm/issues/17732))
* `azurerm_postgresql_server` - fix a regression around `id` ([#17755](https://github.com/hashicorp/terraform-provider-azurerm/issues/17755))

## 3.15.0 (July 21, 2022)

FEATURES: 

* **New Data Source**: `azurerm_cdn_frontdoor_origin_group` ([#17089](https://github.com/hashicorp/terraform-provider-azurerm/issues/17089))
* **New Data Source**: `azurerm_cdn_frontdoor_origin` ([#17089](https://github.com/hashicorp/terraform-provider-azurerm/issues/17089))
* **New Resource**: `azurerm_cdn_frontdoor_origin_group` ([#17089](https://github.com/hashicorp/terraform-provider-azurerm/issues/17089))
* **New Resource**: `azurerm_cdn_frontdoor_origin` ([#17089](https://github.com/hashicorp/terraform-provider-azurerm/issues/17089))
* **New Resource**: `azurerm_application_insights_workbook` ([#17368](https://github.com/hashicorp/terraform-provider-azurerm/issues/17368))
* **New Resource**: `azurerm_monitor_data_collection_rule` ([#17342](https://github.com/hashicorp/terraform-provider-azurerm/issues/17342))
* **New Resource**: `azurerm_route_server` ([#16578](https://github.com/hashicorp/terraform-provider-azurerm/issues/16578))
* **New Resource**: `azurerm_route_server_bgp_connection` ([#16578](https://github.com/hashicorp/terraform-provider-azurerm/issues/16578))
* **New Resource**: `azurerm_web_pubsub_private_link_resource` ([#15550](https://github.com/hashicorp/terraform-provider-azurerm/issues/15550))

ENHANCEMENTS:

* dependencies: updating to `v0.20220715.1071215` of `github.com/hashicorp/go-azure-sdk` ([#17645](https://github.com/hashicorp/terraform-provider-azurerm/issues/17645))
* domainservice: to use `hashicorp/go-azure-sdk` ([#17595](https://github.com/hashicorp/terraform-provider-azurerm/issues/17595))
* servicebus: refactoring to use `hashicorp/go-azure-sdk` ([#17628](https://github.com/hashicorp/terraform-provider-azurerm/issues/17628))
* postgres: refactoring to use `hashicorp/go-azure-sdk` ([#17625](https://github.com/hashicorp/terraform-provider-azurerm/issues/17625))
* `azurerm_kusto_cluster_resource` - support for the `allowed_fqdns`, `allowed_ip_ranges`, and `outbound_network_access_restricted` properties ([#17581](https://github.com/hashicorp/terraform-provider-azurerm/issues/17581))
* `azurerm_storage_account` - supports for the `change_feed_retention_in_days` property ([#17130](https://github.com/hashicorp/terraform-provider-azurerm/issues/17130))

## 3.14.0 (July 14, 2022)

FEATURES:

* **New Resource**: `azurerm_application_insights_workbook_template` ([#17433](https://github.com/hashicorp/terraform-provider-azurerm/issues/17433))
* **New Resource**: `azurerm_gallery_application` ([#17394](https://github.com/hashicorp/terraform-provider-azurerm/issues/17394))
* **New Resource**: `azurerm_gallery_application_version` ([#17394](https://github.com/hashicorp/terraform-provider-azurerm/issues/17394))
 
ENHANCEMENTS:

* dependencies: updating to `v0.20220712.1111122` of `github.com/hashicorp/go-azure-sdk` ([#17606](https://github.com/hashicorp/terraform-provider-azurerm/issues/17606))
* dependencies: updating to `v0.37.0` of `github.com/hashicorp/go-azure-helpers` ([#17588](https://github.com/hashicorp/terraform-provider-azurerm/issues/17588))
* dependencies: updating to `v2.18.0` of `github.com/hashicorp/terraform-plugin-sdk` ([#17141](https://github.com/hashicorp/terraform-provider-azurerm/issues/17141))
* appconfiguration: updating to use API Version `2022-05-01` ([#17467](https://github.com/hashicorp/terraform-provider-azurerm/issues/17467))
* spring: updating to use API Version `2022-05-01-preview` ([#17467](https://github.com/hashicorp/terraform-provider-azurerm/issues/17467))
* databricks: refactoring to use `hashicorp/go-azure-sdk` ([#17475](https://github.com/hashicorp/terraform-provider-azurerm/issues/17475))
* lighthouse: refactoring to use `hashicorp/go-azure-sdk` ([#17590](https://github.com/hashicorp/terraform-provider-azurerm/issues/17590))
* policyremediation: updated to use version `2021-10-01` ([#17298](https://github.com/hashicorp/terraform-provider-azurerm/issues/17298))
* signalr: refactoring to use `hashicorp/go-azure-sdk` ([#17463](https://github.com/hashicorp/terraform-provider-azurerm/issues/17463))
* storage: refactoring `objectreplicationpolicy` to use `hashicorp/go-azure-sdk` ([#17471](https://github.com/hashicorp/terraform-provider-azurerm/issues/17471))
* Data Source: `azurerm_availability_set` - updating to use `hashicorp/go-azure-sdk` ([#17608](https://github.com/hashicorp/terraform-provider-azurerm/issues/17608))
* Data Source: `azurerm_ssh_public_key` - refactoring to use `hashicorp/go-azure-sdk` ([#17609](https://github.com/hashicorp/terraform-provider-azurerm/issues/17609))
* `azurerm_availability_set` - updating to use `hashicorp/go-azure-sdk` ([#17608](https://github.com/hashicorp/terraform-provider-azurerm/issues/17608))
* `azurerm_container_group` - support for the `http_headers` property ([#17519](https://github.com/hashicorp/terraform-provider-azurerm/issues/17519))
* `azurerm_dashboard` - refactoring to use `hashicorp/go-azure-sdk` ([#17598](https://github.com/hashicorp/terraform-provider-azurerm/issues/17598))
* `azurerm_kusto_cluster` - support for the `public_ip_address` property ([#17520](https://github.com/hashicorp/terraform-provider-azurerm/issues/17520))
* `azurerm_kusto_script` - support for the `script_content` property ([#17522](https://github.com/hashicorp/terraform-provider-azurerm/issues/17522))
* `azurerm_kusto_iothub_data_connection` - support for the `database_routing_type` property ([#17526](https://github.com/hashicorp/terraform-provider-azurerm/issues/17526))
* `azurerm_kusto_eventhub_data_connection` - support for the `database_routing_type` property ([#17525](https://github.com/hashicorp/terraform-provider-azurerm/issues/17525))
* `azurerm_kusto_eventgrid_data_connection` - support for the `database_routing_type`, `eventgrid_resource_id`, and `managed_identity_resource_id` properties ([#17524](https://github.com/hashicorp/terraform-provider-azurerm/issues/17524))
* `azurerm_kubernetes_cluster` - support for the `host_group_id` property ([#17496](https://github.com/hashicorp/terraform-provider-azurerm/issues/17496))
* `azurerm_kubernetes_cluster_node_pool` - support for the `host_group_id` property ([#17496](https://github.com/hashicorp/terraform-provider-azurerm/issues/17496))
* `azurerm_linux_virtual_machine_scale_set` - support for `capacity_reservation_group_id` property ([#17530](https://github.com/hashicorp/terraform-provider-azurerm/issues/17530))
* `azurerm_linux_virtual_machine_scale_set` - support for the `placement` property for os disks ([#17013](https://github.com/hashicorp/terraform-provider-azurerm/issues/17013))
* `azurerm_orchestrated_virtual_machine_scale_set` - support for the `placement` property for os disks ([#17013](https://github.com/hashicorp/terraform-provider-azurerm/issues/17013))
* `azurerm_shared_image` - support for the `end_of_life_date` `disk_types_not_allowed`, `max_recommended_vcpu_count`, `max_recommended_vcpu_count`, `max_recommended_memory_in_gb`, `min_recommended_memory_in_gb` ([#17300](https://github.com/hashicorp/terraform-provider-azurerm/issues/17300))
* `azurerm_signalr_service` - Add support for `live_trace` ([#17629](https://github.com/hashicorp/terraform-provider-azurerm/issues/17629))
* `azurerm_ssh_public_key` - refactoring to use `hashicorp/go-azure-sdk` ([#17609](https://github.com/hashicorp/terraform-provider-azurerm/issues/17609))
* `azurerm_stream_analytics_output_blob` - support for the `authentication_mode` property ([#16652](https://github.com/hashicorp/terraform-provider-azurerm/issues/16652))
* `azurerm_windows_virtual_machine_scale_set` - support for `capacity_reservation_group_id` property ([#17530](https://github.com/hashicorp/terraform-provider-azurerm/issues/17530))
* `azurerm_windows_virtual_machine_scale_set` - support for the `placement` property for os disks ([#17013](https://github.com/hashicorp/terraform-provider-azurerm/issues/17013))
 
BUG FIXES:

* `azurerm_api_management` - correct set the API Management Cipher `TLS_RSA_WITH_3DES_EDE_CBC_SHA` ([#17554](https://github.com/hashicorp/terraform-provider-azurerm/issues/17554))
* `azurerm_dev_test_lab_schedule` - deleting the schedule during deletion ([#17614](https://github.com/hashicorp/terraform-provider-azurerm/issues/17614))
* `azurerm_linux_function_app` - set the `default_hostname` properly on read ([#17498](https://github.com/hashicorp/terraform-provider-azurerm/issues/17498))
* `azurerm_linux_function_app_slot` - set the `default_hostname` properly on read ([#17498](https://github.com/hashicorp/terraform-provider-azurerm/issues/17498))
* `azurerm_windows_function_app` - set the `default_hostname` properly on read ([#17498](https://github.com/hashicorp/terraform-provider-azurerm/issues/17498))
* `azurerm_windows_function_app` - correctly create function apps when custom handlers are used ([#17498](https://github.com/hashicorp/terraform-provider-azurerm/issues/17498))
* `azurerm_windows_function_app_slot` - set the `default_hostname` properly on read ([#17498](https://github.com/hashicorp/terraform-provider-azurerm/issues/17498))
* `azurerm_windows_function_app_slot` - correctly create function apps when custom handlers are used ([#17498](https://github.com/hashicorp/terraform-provider-azurerm/issues/17498))

## 3.13.0 (July 08, 2022)

FEATURES:

* **New Data Source**: `azurerm_public_maintenance_configurations` ([#16810](https://github.com/hashicorp/terraform-provider-azurerm/issues/16810))
* **New Resource**: `azurerm_fluid_relay_server` ([#17238](https://github.com/hashicorp/terraform-provider-azurerm/issues/17238))
* **New Resource**: `azurerm_logz_sub_account` ([#16581](https://github.com/hashicorp/terraform-provider-azurerm/issues/16581))

ENHANCEMENTS:

* azurestackhci: refactoring to use `hashicorp/go-azure-sdk` ([#17469](https://github.com/hashicorp/terraform-provider-azurerm/issues/17469))
* containerinstance: refactoring to use `hashicorp/go-azure-sdk` ([#17499](https://github.com/hashicorp/terraform-provider-azurerm/issues/17499))
* eventhub: refactoring to use `hashicorp/go-azure-sdk` ([#17445](https://github.com/hashicorp/terraform-provider-azurerm/issues/17445))
* hardwaresecuritymodules: refactoring to use `hashicorp/go-azure-sdk` ([#17470](https://github.com/hashicorp/terraform-provider-azurerm/issues/17470))
* netapp: refactoring to use `hashicorp/go-azure-sdk` ([#17465](https://github.com/hashicorp/terraform-provider-azurerm/issues/17465))
* privatedns: refactoring to use `hashicorp/go-azure-sdk` ([#17436](https://github.com/hashicorp/terraform-provider-azurerm/issues/17436))
* Data Source: `azurerm_container_registry` - add support for the `data_endpoint_enabled` property ([#17466](https://github.com/hashicorp/terraform-provider-azurerm/issues/17466))
* `azurerm_hdinsight_kafka_cluster` -support for the `network` block ([#17259](https://github.com/hashicorp/terraform-provider-azurerm/issues/17259))
* `azurerm_key_vault_certificate` - will now correctly recover certificates on import ([#17415](https://github.com/hashicorp/terraform-provider-azurerm/issues/17415))
* `azurerm_kubernetes_clusterl`- support for the `capacity_reservation_group_id` property ([#17395](https://github.com/hashicorp/terraform-provider-azurerm/issues/17395))
* `azurerm_kubernetes_node_pool`- support for the `capacity_reservation_group_id` property ([#17395](https://github.com/hashicorp/terraform-provider-azurerm/issues/17395))
* `azurerm_linux_virtual_machine` - support for the `capacity_reservation_group_id` property ([#17236](https://github.com/hashicorp/terraform-provider-azurerm/issues/17236))
* `azurerm_spring_cloud_deployment` - support for the `addon_json` property ([#16984](https://github.com/hashicorp/terraform-provider-azurerm/issues/16984))
* `azurerm_synapse_integration_runtime_azure` - the `location` property now supports `Auto Resolve` ([#17111](https://github.com/hashicorp/terraform-provider-azurerm/issues/17111))
* `azurerm_windows_virtual_machine` - support for the `capacity_reservation_group_id` property ([#17236](https://github.com/hashicorp/terraform-provider-azurerm/issues/17236))

BUG FIXES:

* `azurerm_application_gateway` -  the `request_routing_rule.x.priority` property is now optional ([#17380](https://github.com/hashicorp/terraform-provider-azurerm/issues/17380))

## 3.12.0 (June 30, 2022)

FEATURES:

* **New Resource**: `azurerm_active_directory_domain_service_trust` ([#17045](https://github.com/hashicorp/terraform-provider-azurerm/issues/17045))
* **New Resource**: `azurerm_data_protection_resource_guard` ([#17325](https://github.com/hashicorp/terraform-provider-azurerm/issues/17325))
* **New Resource**: `azurerm_spring_cloud_api_portal_custom_domain` ([#16966](https://github.com/hashicorp/terraform-provider-azurerm/issues/16966))

ENHANCEMENTS:

* dependencies: updating to `v0.20220628.1190740` of `github.com/hashicorp/go-azure-sdk` ([#17399](https://github.com/hashicorp/terraform-provider-azurerm/issues/17399))
* appservice: replacing usages of `ioutil` with `io` ([#17392](https://github.com/hashicorp/terraform-provider-azurerm/issues/17392))
* containerservice: updated to use version `2022-03-02-preview` ([#17084](https://github.com/hashicorp/terraform-provider-azurerm/issues/17084))
* elastic: refactoring to use `hashicorp/go-azure-sdk` ([#17431](https://github.com/hashicorp/terraform-provider-azurerm/issues/17431))
* loadtest: refactoring to use `hashicorp/go-azure-sdk` ([#17432](https://github.com/hashicorp/terraform-provider-azurerm/issues/17432))
* maps: refactoring to use `hashicorp/go-azure-sdk` ([#17434](https://github.com/hashicorp/terraform-provider-azurerm/issues/17434))
* mixedreality: switching to use `hashicorp/go-azure-sdk` ([#17417](https://github.com/hashicorp/terraform-provider-azurerm/issues/17417))
* msi: refactoring to use `hashicorp/go-azure-sdk` ([#17430](https://github.com/hashicorp/terraform-provider-azurerm/issues/17430))
* powerbi: refactoring to use `hashicorp/go-azure-sdk` ([#17435](https://github.com/hashicorp/terraform-provider-azurerm/issues/17435))
* purview: refactoring to use `hashicorp/go-azure-sdk` ([#17419](https://github.com/hashicorp/terraform-provider-azurerm/issues/17419))
* redisenterprise: refactoring to use `hashicorp/go-azure-sdk` ([#17387](https://github.com/hashicorp/terraform-provider-azurerm/issues/17387))
* relay: refactoring to use `hashicorp/go-azure-sdk` ([#17385](https://github.com/hashicorp/terraform-provider-azurerm/issues/17385))
* search: refactoring to use `hashicorp/go-azure-sdk` ([#17386](https://github.com/hashicorp/terraform-provider-azurerm/issues/17386))
* servicefabricmanaged: refactoring to use `hashicorp/go-azure-sdk` ([#17384](https://github.com/hashicorp/terraform-provider-azurerm/issues/17384))
* trafficmanager: refactoring to use `hashicorp/go-azure-sdk` ([#17383](https://github.com/hashicorp/terraform-provider-azurerm/issues/17383))
* videoanalyzer: refactoring to use `hashicorp/go-azure-sdk` ([#17382](https://github.com/hashicorp/terraform-provider-azurerm/issues/17382))
* vmware: refactoring to use `hashicorp/go-azure-sdk` ([#17381](https://github.com/hashicorp/terraform-provider-azurerm/issues/17381))
* Data Source: `azurerm_key_vault_key` - exporting the `resource_id` and `resource_versionless_id` attributes ([#17424](https://github.com/hashicorp/terraform-provider-azurerm/issues/17424))
* Data Source: `azurerm_key_vault_secret` - exporting the `resource_id` and `resource_versionless_id` attributes ([#17424](https://github.com/hashicorp/terraform-provider-azurerm/issues/17424))
* Data Source: `azurerm_spatial_anchors_account` - exposing the `tags` attribute ([#17417](https://github.com/hashicorp/terraform-provider-azurerm/issues/17417))
* `azurerm_bot_service_azure_bot` - support new bot type with the `microsoft_app_msi_id`, `microsoft_app_tenant_id`,  and `microsoft_app_type` properties ([#17077](https://github.com/hashicorp/terraform-provider-azurerm/issues/17077))
* `azurerm_bot_channels_registration` - support for the `streaming_endpoint_enabled` property ([#17369](https://github.com/hashicorp/terraform-provider-azurerm/issues/17369))
* `azurerm_data_factory` - support for the `purview_id` property ([#17001](https://github.com/hashicorp/terraform-provider-azurerm/issues/17001))
* `azurerm_digital_twins_instance` - support for the `identity` block ([#17076](https://github.com/hashicorp/terraform-provider-azurerm/issues/17076))
* `azurerm_key_vault_key` - exporting the `resource_id` and `resource_versionless_id` attributes ([#17424](https://github.com/hashicorp/terraform-provider-azurerm/issues/17424))
* `azurerm_key_vault_secret` - exporting the `resource_id` and `resource_versionless_id` attributes ([#17424](https://github.com/hashicorp/terraform-provider-azurerm/issues/17424))
* `azurerm_kubernetes_cluster` - support for version aliases ([#17084](https://github.com/hashicorp/terraform-provider-azurerm/issues/17084))
* `azurerm_linux_web_app` - support for the `virtual_network_subnet_id` property ([#17354](https://github.com/hashicorp/terraform-provider-azurerm/issues/17354))
* `azurerm_linux_web_app_slot` - support for the `virtual_network_subnet_id` property ([#17354](https://github.com/hashicorp/terraform-provider-azurerm/issues/17354))
* `azurerm_private_link_service` - support for the `fqdns` property ([#17366](https://github.com/hashicorp/terraform-provider-azurerm/issues/17366))
* `azurerm_shared_image_version` - support `Premium_LRS` for the `storage_account_type` property ([#17390](https://github.com/hashicorp/terraform-provider-azurerm/issues/17390))
* `azurerm_shared_image_version` - support for the `disk_encryption_set_id`, `end_of_life_date`, and `replication_mode` properties ([#17295](https://github.com/hashicorp/terraform-provider-azurerm/issues/17295))
* `azurerm_static_site_custom_domain` - the `validation_type` propety is now optional ([#15849](https://github.com/hashicorp/terraform-provider-azurerm/issues/15849))
* `azurerm_vpn_site` - support for the `o365_policy` block ([#16820](https://github.com/hashicorp/terraform-provider-azurerm/issues/16820))

BUG FIXES:

* Data Source: `azurerm_key_vault` - caching the Key Vault URI when the Key Vault has been retrieved ([#17407](https://github.com/hashicorp/terraform-provider-azurerm/issues/17407))
* `azurerm_application_gateway` - prevent a crash when the `waf_configuration` block is removed ([#17241](https://github.com/hashicorp/terraform-provider-azurerm/issues/17241))
* `azurerm_data_factory_dataset_snowflake` - ensuring `schema` is sent to the API to fix a UI bug in the Azure Data Factory Portal ([#17346](https://github.com/hashicorp/terraform-provider-azurerm/issues/17346))
* `azurerm_data_factory_linked_service_azure_file_storage` - corredctly assign `user_id`([#17398](https://github.com/hashicorp/terraform-provider-azurerm/issues/17398))
* `azurerm_key_vault` - ensuring that `soft_delete_enabled` is explicitly set when `purge_protection_enabled` is set ([#16368](https://github.com/hashicorp/terraform-provider-azurerm/issues/16368))
* `azurerm_linux_function_app` - correctly validate the `app_setting_names` and `connection_string_names` properties within the `sticky_settings` block ([#17209](https://github.com/hashicorp/terraform-provider-azurerm/issues/17209))
* `azurerm_linux_web_app` - correctly configure `auto_heal` and `slow_request` ([#17296](https://github.com/hashicorp/terraform-provider-azurerm/issues/17296))
* `azurerm_linux_web_app` - correctly validate the `app_setting_names` and `connection_string_names` properties within the `sticky_settings` block ([#17209](https://github.com/hashicorp/terraform-provider-azurerm/issues/17209))
* `azurerm_management_group_policy_assignment` - the `name` property can no longer contain `/` ([#16484](https://github.com/hashicorp/terraform-provider-azurerm/issues/16484))
* `azurerm_policy_assignment` - the `name` property can no longer contain `/` ([#16484](https://github.com/hashicorp/terraform-provider-azurerm/issues/16484))
* `azurerm_resource_group_policy_assignment` - the `name` property can no longer contain `/` ([#16484](https://github.com/hashicorp/terraform-provider-azurerm/issues/16484))
* `azurerm_subscription_policy_assignment` - the `name` property can no longer contain `/` ([#16484](https://github.com/hashicorp/terraform-provider-azurerm/issues/16484))
* `azurerm_windows_function_app` - correctly validate the `app_setting_names` and `connection_string_names` properties within the `sticky_settings` block ([#17209](https://github.com/hashicorp/terraform-provider-azurerm/issues/17209))
* `azurerm_windows_web_app` - correctly configure `auto_heal` and `slow_request` ([#17296](https://github.com/hashicorp/terraform-provider-azurerm/issues/17296))
* `azurerm_windows_web_app` - correctly validate the `app_setting_names` and `connection_string_names` properties within the `sticky_settings` block ([#17209](https://github.com/hashicorp/terraform-provider-azurerm/issues/17209))

## 3.11.0 (June 23, 2022)

FEATURES:

* **New Data Source**: `azurerm_management_group_template_deployment` ([#14524](https://github.com/hashicorp/terraform-provider-azurerm/issues/14524))
* **New Data Source**: `azurerm_policy_assignment` ([#16527](https://github.com/hashicorp/terraform-provider-azurerm/issues/16527))
* **New Data Source**: `azurerm_resource_group_template_deployment` ([#14524](https://github.com/hashicorp/terraform-provider-azurerm/issues/14524))
* **New Data Source**: `azurerm_subscription_template_deployment` ([#14524](https://github.com/hashicorp/terraform-provider-azurerm/issues/14524))
* **New Data Source**: `azurerm_tenant_template_deployment` ([#14524](https://github.com/hashicorp/terraform-provider-azurerm/issues/14524))

ENHANCEMENTS:

* dependencies: updating to `v0.20220623.1064317` of `github.com/hashicorp/go-azure-sdk` ([#17348](https://github.com/hashicorp/terraform-provider-azurerm/issues/17348))
* batch: updating to use API Version `2022-01-01` ([#17219](https://github.com/hashicorp/terraform-provider-azurerm/issues/17219))
* confidentialledger: updating to use API Version `2022-05-13` ([#17146](https://github.com/hashicorp/terraform-provider-azurerm/issues/17146))
* desktopvirtualization: refactoring to use `hashicorp/go-azure-sdk` ([#17340](https://github.com/hashicorp/terraform-provider-azurerm/issues/17340))
* Data Source: `azurerm_managed_disk` - exporting the `disk_access_id` attribute ([#17270](https://github.com/hashicorp/terraform-provider-azurerm/issues/17270))
* Data Source: `azurerm_managed_disk` - exporting the `network_access_policy` attribute ([#17270](https://github.com/hashicorp/terraform-provider-azurerm/issues/17270))
* Data Source: `azurerm_storage_account` - add support for the `identity` property ([#17215](https://github.com/hashicorp/terraform-provider-azurerm/issues/17215))

BUG FIXES:

* Data Source: `azurerm_mysql_flexible_server` - generate the correct terraform resource ID ([#17301](https://github.com/hashicorp/terraform-provider-azurerm/issues/17301))
* `azurerm_shared_image` - the `privacy_statement_uri`, `publisher`, `offer`, and `sku` fields are now ForceNew ([#17289](https://github.com/hashicorp/terraform-provider-azurerm/issues/17289))
* `azurerm_shared_image_*` - correctly validate the `gallery_name` property ([#17201](https://github.com/hashicorp/terraform-provider-azurerm/issues/17201))
* `azurerm_time_series_insights_gen2_environment` - correctly order `id_properties` ([#17234](https://github.com/hashicorp/terraform-provider-azurerm/issues/17234))

## 3.10.0 (June 09, 2022)

FEATURES:

* **New Data Source**: `azurerm_cdn_frontdoor_rule_set` ([#17094](https://github.com/hashicorp/terraform-provider-azurerm/issues/17094))
* **New Resource**: `azurerm_capacity_reservation_group` ([#16464](https://github.com/hashicorp/terraform-provider-azurerm/issues/16464))
* **New Resource**: `azurerm_capacity_reservation` ([#16464](https://github.com/hashicorp/terraform-provider-azurerm/issues/16464))
* **New Resource**: `azurerm_cdn_frontdoor_rule_set` ([#17094](https://github.com/hashicorp/terraform-provider-azurerm/issues/17094))

ENHANCEMENTS:

* `azurerm_cosmosdb_cassandra_cluster` - support for the `authentication_method`, `client_certificate`, `external_gossip_certificate`, `external_seed_node`, `identity`, `repair_enabled` and `version` properties ([#16799](https://github.com/hashicorp/terraform-provider-azurerm/issues/16799))
* `azurerm_key_vault_managed_hardware_security_module` - support for purging when soft deleted ([#17148](https://github.com/hashicorp/terraform-provider-azurerm/issues/17148))
* `azurerm_hpc_cache` - support for `identity` block and the `key_vault_key_id` and `automatically_rotate_key_to_latest_enabled` properties ([#16972](https://github.com/hashicorp/terraform-provider-azurerm/issues/16972))

BUG FIXES:

* `azurerm_api_management` - default hostname proxy configuration is no longer ignored ([#16524](https://github.com/hashicorp/terraform-provider-azurerm/issues/16524))
* `azurerm_application_gateway` - add default value for `backend_http_settings.0.request_timeout` ([#17162](https://github.com/hashicorp/terraform-provider-azurerm/issues/17162))
* `azurerm_applicaton_gateway` -`priority` is now required ([#16849](https://github.com/hashicorp/terraform-provider-azurerm/issues/16849))
* `azurerm_container_group` - Double the delete check timeout for nic ([#17115](https://github.com/hashicorp/terraform-provider-azurerm/issues/17115))
* `azurerm_windows_function_app_x` - `custom_domain_verification_id` is now written to state file (([#17183](https://github.com/hashicorp/terraform-provider-azurerm/issues/17183))

## 3.9.0 (June 02, 2022)

FEATURES:

* **New Data Source**: `azurerm_app_configuration_keys` ([#17053](https://github.com/hashicorp/terraform-provider-azurerm/issues/17053))
* **New Data Source**: `azurerm_cdn_frontdoor_endpoint` ([#17078](https://github.com/hashicorp/terraform-provider-azurerm/issues/17078))
* **New Data Source**: `azurerm_cdn_frontdoor_profile` ([#17061](https://github.com/hashicorp/terraform-provider-azurerm/issues/17061))
* **New Resource**: `azurerm_cdn_frontdoor_endpoint` ([#17078](https://github.com/hashicorp/terraform-provider-azurerm/issues/17078))
* **New Resource**: `azurerm_cdn_frontdoor_profile` ([#17061](https://github.com/hashicorp/terraform-provider-azurerm/issues/17061))
* **New Resource**: `azurerm_sentinel_data_connector_office_atp` ([#16825](https://github.com/hashicorp/terraform-provider-azurerm/issues/16825))
* **New Resource**: `azurerm_vpn_server_configuration_policy_group` ([#16911](https://github.com/hashicorp/terraform-provider-azurerm/issues/16911))

ENHANCEMENTS:

* dependencies: upgrading to `v0.33.0` of `github.com/hashicorp/go-azure-hepers` ([#17074](https://github.com/hashicorp/terraform-provider-azurerm/issues/17074))
* dependencies: upgrading to `v1.6.1` of `github.com/hashicorp/go-getter` ([#17074](https://github.com/hashicorp/terraform-provider-azurerm/issues/17074))
* dependencies: upgrade `netapp` to `2021-10-01` ([#17043](https://github.com/hashicorp/terraform-provider-azurerm/issues/17043))
* `azurerm_batch_job` - refactor to split `create` and `update` ([#17138](https://github.com/hashicorp/terraform-provider-azurerm/issues/17138))
* `azurerm_data_factory_trigger_schedule` - support for the `pipeline` block ([#16922](https://github.com/hashicorp/terraform-provider-azurerm/issues/16922))
* `azurerm_backup_policy_vm` - support for `V2` policies viu the `policy_type` property, supporting Enhanced Policies of the hourly type ([#16940](https://github.com/hashicorp/terraform-provider-azurerm/issues/16940))
* `azurerm_log_analytics_workspace` - allow property updates when a workspace is linked to a cluster ([#17069](https://github.com/hashicorp/terraform-provider-azurerm/issues/17069))
* `azurerm_netapp_volume` - support for the `network_features` property ([#17043](https://github.com/hashicorp/terraform-provider-azurerm/issues/17043))
* `azurerm_provider_registration` - refactor to split `create` and `update` ([#17138](https://github.com/hashicorp/terraform-provider-azurerm/issues/17138))
* `azurerm_web_pubsub_hub` - the `event_handler` block is now optional ([#17037](https://github.com/hashicorp/terraform-provider-azurerm/issues/17037))
* `azurerm_redis_cache` - support the `identity` block ([#16990](https://github.com/hashicorp/terraform-provider-azurerm/issues/16990))
* `azurerm_service_fabric_managed_cluster` - refactor to split `create` and `update` ([#17138](https://github.com/hashicorp/terraform-provider-azurerm/issues/17138))
* `azurerm_synapse_role_assignment` - the `role_name` property now supports `Synapse Monitoring Operator` ([#17024](https://github.com/hashicorp/terraform-provider-azurerm/issues/17024))
* `azurerm_vpn_gateway_nat_rule` - support for the `port_range` property ([#16724](https://github.com/hashicorp/terraform-provider-azurerm/issues/16724))

BUG FIXES:

* `azurerm_container_registry_task` - sending `authentication` within the `source_trigger` block when updating ([#17002](https://github.com/hashicorp/terraform-provider-azurerm/issues/17002))
* `azurerm_eventhub_authorization_rule` - extend regex char limit for `name` ([#17057](https://github.com/hashicorp/terraform-provider-azurerm/issues/17057))
* `azurerm_kubernetes_cluster` - prevent a potential crash during import of a cluster that doesn't have an API Server Access Profile ([#17005](https://github.com/hashicorp/terraform-provider-azurerm/issues/17005))

## 3.8.0 (May 26, 2022)

FEATURES:

* **New Resource**: `azurerm_mssql_server_dns_alias` ([#16861](https://github.com/hashicorp/terraform-provider-azurerm/issues/16861))
* **New Resource**: `azurerm_spring_cloud_gateway_route_config` ([#16721](https://github.com/hashicorp/terraform-provider-azurerm/issues/16721))
* **New Resource**: `azurerm_spring_cloud_api_portal` ([#16719](https://github.com/hashicorp/terraform-provider-azurerm/issues/16719))
* **New Resource**: `azurerm_spring_cloud_build_deployment` ([#16730](https://github.com/hashicorp/terraform-provider-azurerm/issues/16730))

ENHANCEMENTS:

* dependencies: upgrade `botservice` to `2021-05-01-preview` ([#16665](https://github.com/hashicorp/terraform-provider-azurerm/issues/16665))
* dependencies: upgrade `keyvault` to `2021-10-01` ([#16955](https://github.com/hashicorp/terraform-provider-azurerm/issues/16955))
* `azurerm_active_directory_domain_service` - supports for the `domain_configuration_type` property ([#16920](https://github.com/hashicorp/terraform-provider-azurerm/issues/16920))
* `azurerm_backup_protected_vm` - allow the attached vm to be disassociated from the backup ([#16939](https://github.com/hashicorp/terraform-provider-azurerm/issues/16939))
* `azurerm_backup_protected_vm` - the backup is now removed from state when it is soft deleted ([#16939](https://github.com/hashicorp/terraform-provider-azurerm/issues/16939))
* `azurerm_portal_dashboard` - now supports the `display_name` argument ([#16406](https://github.com/hashicorp/terraform-provider-azurerm/issues/16406))
* `azurerm_data_factory_trigger_schedule` - support for the `time_zone` property ([#16918](https://github.com/hashicorp/terraform-provider-azurerm/issues/16918))
* `azurerm_linux_virtual_machine` - add support for Confidential VMs ([#16905](https://github.com/hashicorp/terraform-provider-azurerm/issues/16905))
* `azurerm_linux_virtual_machine_scale_set` - add support for Confidential VMs ([#16916](https://github.com/hashicorp/terraform-provider-azurerm/issues/16916))
* `azurerm_linux_web_app` - add support for `zip_deploy_file` property ([#16779](https://github.com/hashicorp/terraform-provider-azurerm/issues/16779))
* `azurerm_linux_web_app_slot` - add support for `zip_deploy_file` property ([#16779](https://github.com/hashicorp/terraform-provider-azurerm/issues/16779))
* `azurerm_managed_disk` - add support for Confidential VM ([#16908](https://github.com/hashicorp/terraform-provider-azurerm/issues/16908))
* `azurerm_spring_cloud_service` - suppport the `build_agent_pool_size` property ([#16841](https://github.com/hashicorp/terraform-provider-azurerm/issues/16841))
* `azurerm_spring_cloud_service`- support the `zone_redundant` property ([#16872](https://github.com/hashicorp/terraform-provider-azurerm/issues/16872))
* `azurerm_synapse_spark_pool` - the `spark_version` property now supports `3.2` ([#16906](https://github.com/hashicorp/terraform-provider-azurerm/issues/16906))
* `azurerm_virtual_network_gateway_connection` - support for the `egress_nat_rule_ids` and `ingress_nat_rule_ids` properties ([#16862](https://github.com/hashicorp/terraform-provider-azurerm/issues/16862))
* `azurerm_vpn_gateway` - support for the `bgp_route_translation_for_nat_enabled` property ([#16817](https://github.com/hashicorp/terraform-provider-azurerm/issues/16817))
* `azurerm_vpn_gateway_connection` - support for the `custom_bgp_address` block ([#16960](https://github.com/hashicorp/terraform-provider-azurerm/issues/16960))
* `azurerm_windows_virtual_machine` - add support for Confidential VMs ([#16905](https://github.com/hashicorp/terraform-provider-azurerm/issues/16905))
* `azurerm_windows_virtual_machine_scale_set` - add support for Confidential VM ([#16916](https://github.com/hashicorp/terraform-provider-azurerm/issues/16916))
* `azurerm_windows_web_app` - add support for `zip_deploy_file` property ([#16779](https://github.com/hashicorp/terraform-provider-azurerm/issues/16779))
* `azurerm_windows_web_app_slot` - add support for `zip_deploy_file` property ([#16779](https://github.com/hashicorp/terraform-provider-azurerm/issues/16779))

BUG FIXES:

* `azurerm_mysql_server` -  fix an error updating `public_network_access_enabled` with replicas ([#16506](https://github.com/hashicorp/terraform-provider-azurerm/issues/16506))
* `azurerm_linux_function_app_slot` - correctly check for name availability during creation ([#16410](https://github.com/hashicorp/terraform-provider-azurerm/issues/16410))
* `azurerm_windows_function_app_slot` - correctly check for name availability during creation ([#16410](https://github.com/hashicorp/terraform-provider-azurerm/issues/16410))
* `azurerm_windows_virtual_machine` - changing the `timezone` property now creates a new resources ([#16866](https://github.com/hashicorp/terraform-provider-azurerm/issues/16866))

## 3.7.0 (May 19, 2022)

FEATURES:

* **New Authentication Method:** OIDC ([#16555](https://github.com/hashicorp/terraform-provider-azurerm/issues/16555))
* **New Data Source**: `azurerm_elastic_cloud_elasticsearch` ([#14821](https://github.com/hashicorp/terraform-provider-azurerm/issues/14821))
* **New Resource**: `azurerm_elastic_cloud_elasticsearch` ([#14821](https://github.com/hashicorp/terraform-provider-azurerm/issues/14821))
* **New Resource**: `azurerm_healthcare_fhir_service` ([#15913](https://github.com/hashicorp/terraform-provider-azurerm/issues/15913))
* **New Resource**: `azurerm_virtual_network_gateway_nat_rule` ([#15720](https://github.com/hashicorp/terraform-provider-azurerm/issues/15720))

ENHANCEMENTS:

* dependencies: upgrade `redis` to `2020-12-01` ([#16532](https://github.com/hashicorp/terraform-provider-azurerm/issues/16532))
* `azurerm_container_registry` - support changing replications ([#16678](https://github.com/hashicorp/terraform-provider-azurerm/issues/16678))
* `azurerm_disk_encryption_set` - the `encryption_type` property now supports `ConfidentialVmEncryptedWithCustomerKey` ([#16870](https://github.com/hashicorp/terraform-provider-azurerm/issues/16870))
* `azurerm_linux_function_app` - add support for PowerShell `7.2`  ([#16718](https://github.com/hashicorp/terraform-provider-azurerm/issues/16718))
* `azurerm_signalr_service` - support the `Premium_P1` SKU ([#16875](https://github.com/hashicorp/terraform-provider-azurerm/issues/16875))
* `azurerm_spring_cloud_app` - support for the `identity` block ([#16280](https://github.com/hashicorp/terraform-provider-azurerm/issues/16280))
* `azurerm_spring_cloud_app` - support for the `addon_json` property ([#16722](https://github.com/hashicorp/terraform-provider-azurerm/issues/16722))
* `azurerm_windows_function_app` - support for PowerShell `7.2`  ([#16718](https://github.com/hashicorp/terraform-provider-azurerm/issues/16718))
* `azurerm_mssql_managed_instance` - support for the `maintenance_configuration_name` property ([#16832](https://github.com/hashicorp/terraform-provider-azurerm/issues/16832))

BUG FIXES:

* Data Source: `azurerm_databricks_workspace` - prevent a panic when the SKU field is missing ([#16819](https://github.com/hashicorp/terraform-provider-azurerm/issues/16819))
* `azurerm_application_insights_web_test` - working around a breaking change in the API where creation would fail ([#16845](https://github.com/hashicorp/terraform-provider-azurerm/issues/16845))
* `azurerm_express_route_gateway` - handle gateway connections not found error ([#16804](https://github.com/hashicorp/terraform-provider-azurerm/issues/16804))
* `azurerm_shared_image` - changing the `eula` property now creates a new resource ([#16868](https://github.com/hashicorp/terraform-provider-azurerm/issues/16868))

DEPRECATIONS:

* `azurerm_video_analyzer` - Video Analyzer (Preview) is now Deprecated and will be Retired on 2022-11-30 - as such this resource is deprecated and will be removed in v4.0 of the AzureRM Provider ([#16847](https://github.com/hashicorp/terraform-provider-azurerm/issues/16847))
* `azurerm_video_analyzer_edge_module` - Video Analyzer (Preview) is now Deprecated and will be Retired on 2022-11-30 - as such this resource is deprecated and will be removed in v4.0 of the AzureRM Provider ([#16847](https://github.com/hashicorp/terraform-provider-azurerm/issues/16847))

## 3.6.0 (May 12, 2022)

FEATURES:

* **New Resource**: `azurerm_confidential_ledger` ([#15420](https://github.com/hashicorp/terraform-provider-azurerm/issues/15420))
* **New Resource**: `azurerm_managed_disk_sas_token` ([#15558](https://github.com/hashicorp/terraform-provider-azurerm/issues/15558))
* **New Resource**: `azurerm_spring_cloud_gateway` ([#16175](https://github.com/hashicorp/terraform-provider-azurerm/issues/16175))
* **New Resource**: `azurerm_spring_cloud_build_pack_binding` ([#16673](https://github.com/hashicorp/terraform-provider-azurerm/issues/16673))
* **New Resource**: `azurerm_spring_cloud_gateway_custom_domain` ([#16720](https://github.com/hashicorp/terraform-provider-azurerm/issues/16720))
* **New Resource**: `azurerm_stream_analytics_output_powerbi` ([#16439](https://github.com/hashicorp/terraform-provider-azurerm/issues/16439))

ENHANCEMENTS:

* dependencies: updating to `v64.0.0` of `github.com/Azure/azure-sdk-for-go` ([#16631](https://github.com/hashicorp/terraform-provider-azurerm/issues/16631))
* dependencies: upgrade `network` to `2021-08-01` ([#16631](https://github.com/hashicorp/terraform-provider-azurerm/issues/16631))
* `azurerm_container_group` - support for the `key_vault_key_id` property (Customer Managed Key encryption) ([#16709](https://github.com/hashicorp/terraform-provider-azurerm/issues/16709))
* `azurerm_cosmosdb_account` - support mongo version `4.2` ([#16738](https://github.com/hashicorp/terraform-provider-azurerm/issues/16738))
* `azurerm_cosmosdb_cassandra_cluster` - support for the `tags` property ([#16743](https://github.com/hashicorp/terraform-provider-azurerm/issues/16743))
* `azurerm_kubernetes_cluster_node_pool` - the property `node_labels` can now be updated ([#16360](https://github.com/hashicorp/terraform-provider-azurerm/issues/16360))
* `azurerm_kubernetes_cluster` - the property `default_node_pool.node_labels` can now be updated ([#16360](https://github.com/hashicorp/terraform-provider-azurerm/issues/16360))
* `azurerm_kubernetes_cluster` - allow value `none` for `network_profile.network_plugin` ([#16250](https://github.com/hashicorp/terraform-provider-azurerm/issues/16250))
* `azurerm_kusto_script` - lock kusto cluster so multiple scripts can be applied ([#16690](https://github.com/hashicorp/terraform-provider-azurerm/issues/16690))
* `azurerm_storage_share` - support the `access_tier` attribute ([#16462](https://github.com/hashicorp/terraform-provider-azurerm/issues/16462))
* `azurerm_snapshot` - support for the `trusted_launch_enabled` propertyu ([#16679](https://github.com/hashicorp/terraform-provider-azurerm/issues/16679))
* `azurerm_stream_analytics_function_javascript_uda` - support for the `input.configuration_parameter` property ([#16575](https://github.com/hashicorp/terraform-provider-azurerm/issues/16575))
* `azurerm_stream_analytics_function_javascript_udf` - support for the `input.configuration_parameter` property ([#16579](https://github.com/hashicorp/terraform-provider-azurerm/issues/16579))
* `azurerm_linux_virtual_machine` - correctly support for the update the `diff_disk_settings.placement` property ([#14847](https://github.com/hashicorp/terraform-provider-azurerm/issues/14847))
* `azurerm_virtual_network_gateway_connection` - support for the `custom_bgp_addresses` property ([#16631](https://github.com/hashicorp/terraform-provider-azurerm/issues/16631))
* `azurerm_windows_virtual_machine` - correctly support for the update the `diff_disk_settings.placement` property ([#14847](https://github.com/hashicorp/terraform-provider-azurerm/issues/14847))

BUG FIXES:

* `azurerm_app_configuration_feature` - allow successful creation of resource without specifying any optional filters ([#16459](https://github.com/hashicorp/terraform-provider-azurerm/issues/16459))
* `azurerm_mssql_managed_instance_failover_group` - correctly import resource and sent primary isntance id ([#16705](https://github.com/hashicorp/terraform-provider-azurerm/issues/16705))

## 3.5.0 (May 05, 2022)

FEATURES:

* **New Data Source**: `azurerm_healthcare_dicom_service` ([#15887](https://github.com/hashicorp/terraform-provider-azurerm/issues/15887))
* **New Resource**: `azurerm_healthcare_dicom_service` ([#15887](https://github.com/hashicorp/terraform-provider-azurerm/issues/15887))
* **New Resource**: `azurerm_mssql_managed_instance_vulnerability_assessment` ([#16639](https://github.com/hashicorp/terraform-provider-azurerm/issues/16639))
* **New resource**: `azurerm_sentinel_data_connector_aws_s3` ([#16440](https://github.com/hashicorp/terraform-provider-azurerm/issues/16440))
* **New Resource**: `azurerm_spring_cloud_builder` ([#16036](https://github.com/hashicorp/terraform-provider-azurerm/issues/16036))
* **New Resource**: `azurerm_spring_cloud_configuration_service` ([#16087](https://github.com/hashicorp/terraform-provider-azurerm/issues/16087))

ENHANCEMENTS:

* dependencies: updating to `v63.4.0` of `github.com/Azure/azure-sdk-for-go` ([#16533](https://github.com/hashicorp/terraform-provider-azurerm/issues/16533))
* dependencies: updating to `v1.5.11` of `github.com/hashicorp/go-getter` ([#16659](https://github.com/hashicorp/terraform-provider-azurerm/issues/16659))
* dependencies: upgrade `recoveryservices` to `2021-12-01` ([#16001](https://github.com/hashicorp/terraform-provider-azurerm/issues/16001))
* `azurerm_linux_virtual_machine_scale_set` - improve validation on the `termination_notification.timeout` property ([#16594](https://github.com/hashicorp/terraform-provider-azurerm/issues/16594))
* `azurerm_orchestrated_virtual_machine_scale_set` - improve validation on the `termination_notification.timeout` property ([#16594](https://github.com/hashicorp/terraform-provider-azurerm/issues/16594))
* `azurerm_servicebus_namespace` - the `sku` property can now be updated to `Basic` or `Standard` without recreating the resource ([#16523](https://github.com/hashicorp/terraform-provider-azurerm/issues/16523))
* `azurerm_storage_account` - support for the `cross_tenant_replication_enabled` property ([#16351](https://github.com/hashicorp/terraform-provider-azurerm/issues/16351))
* `azurerm_windows_virtual_machine_scale_set` - improve validation on the `termination_notification.timeout` property ([#16594](https://github.com/hashicorp/terraform-provider-azurerm/issues/16594))
* `azurerm_virtual_network_gateway_connection` - the `traffic_selector_policy` property can now be specified ([#15938](https://github.com/hashicorp/terraform-provider-azurerm/issues/15938))
* `azurerm_stream_analytics_output_servicebus_queue` - support for the `property_columns` and `system_property_columns` properties ([#16572](https://github.com/hashicorp/terraform-provider-azurerm/issues/16572))

BUG FIXES:

* Data Source: `azurerm_servicebus_queue_authorization_rule` - prevent a possible crash by setting `queue_name` correctly ([#16561](https://github.com/hashicorp/terraform-provider-azurerm/issues/16561))
* Data Source: `azurerm_service_plan:` - correctly populate the `kind` and `os_type` attributes ([#16431](https://github.com/hashicorp/terraform-provider-azurerm/issues/16431))
* `azurerm_data_factory_dataset_delimited_text` - set defaults properly for `column_delimiter`, `quote_character`, `escape_character`, `first_row_as_header` and `null_value` ([#16543](https://github.com/hashicorp/terraform-provider-azurerm/issues/16543))
* `azurerm_linux_function_app` - correctly deduplicate user `app_settings` ([#15740](https://github.com/hashicorp/terraform-provider-azurerm/issues/15740))
* `azurerm_linux_function_app` -  fix `app_settings.WEBSITE_RUN_FROM_PACKAGE` handling from external sources ([#16641](https://github.com/hashicorp/terraform-provider-azurerm/issues/16641))
* `azurerm_linux_function_app_slot` - correctly deduplicate user `app_settings` ([#15740](https://github.com/hashicorp/terraform-provider-azurerm/issues/15740))
* `azurerm_linux_function_app_slot` - fix `app_settings.WEBSITE_RUN_FROM_PACKAGE` handling from external sources ([#16641](https://github.com/hashicorp/terraform-provider-azurerm/issues/16641))
* `azurerm_machine_learning_compute_cluster` - resource will now be deleted instead of just detached ([#16640](https://github.com/hashicorp/terraform-provider-azurerm/issues/16640))
* `azurerm_windows_function_app` - correctly deduplicate user `app_settings` ([#15740](https://github.com/hashicorp/terraform-provider-azurerm/issues/15740))
* `azurerm_windows_function_app_slot` - correctly deduplicate user `app_settings` ([#15740](https://github.com/hashicorp/terraform-provider-azurerm/issues/15740))

## 3.4.0 (April 28, 2022)

FEATURES:

* **New Resource**: `azurerm_stream_analytics_output_cosmosdb` ([#16441](https://github.com/hashicorp/terraform-provider-azurerm/issues/16441))

ENHANCEMENTS:

* dependencies: updating to `v63.1.0` of `github.com/Azure/azure-sdk-for-go` ([#16283](https://github.com/hashicorp/terraform-provider-azurerm/issues/16283))
* dependencies: updating to `v0.11.26` of `github.com/Azure/go-autorest` ([#16458](https://github.com/hashicorp/terraform-provider-azurerm/issues/16458))
* dependencies: upgrading to `v0.30.0` of `github.com/hashicorp/go-azure-helpers` ([#16504](https://github.com/hashicorp/terraform-provider-azurerm/issues/16504))
* dependencies: upgrade `sqlvirtualmachine` to `2021-11-01-preview` ([#15835](https://github.com/hashicorp/terraform-provider-azurerm/issues/15835))
* Data Source: `azurerm_linux_function_app` - add support for `sticky_settings` ([#16546](https://github.com/hashicorp/terraform-provider-azurerm/issues/16546))
* Data Source: `azurerm_linux_web_app` - add support for `sticky_settings` ([#16546](https://github.com/hashicorp/terraform-provider-azurerm/issues/16546))
* Data Source: `azurerm_windows_function_app` - add support for `sticky_settings` ([#16546](https://github.com/hashicorp/terraform-provider-azurerm/issues/16546))
* Data Source: `azurerm_windows_web_app` - add support for `sticky_settings` ([#16546](https://github.com/hashicorp/terraform-provider-azurerm/issues/16546))
* `azurerm_kubernetes_cluster` - support for the `run_command_enabled` property ([#15029](https://github.com/hashicorp/terraform-provider-azurerm/issues/15029))
* `azurerm_linux_function_app` - add support for `sticky_settings` ([#16546](https://github.com/hashicorp/terraform-provider-azurerm/issues/16546))
* `azurerm_linux_web_app` - add support for `sticky_settings` ([#16546](https://github.com/hashicorp/terraform-provider-azurerm/issues/16546))
* `azurerm_monitor_aad_diagnostic_setting` - remove validation on `log.category` to allow for new log categories that are available in Azure ([#16534](https://github.com/hashicorp/terraform-provider-azurerm/issues/16534))
* `azurerm_mssql_database` - Support for `short_term_retention_policy.0.backup_interval_in_hours` ([#16528](https://github.com/hashicorp/terraform-provider-azurerm/issues/16528))
* `azurerm_postgresql_server` - add validation for `public_network_access_enabled` ([#16516](https://github.com/hashicorp/terraform-provider-azurerm/issues/16516))
* `azurerm_stream_analytics_job` - support for the `type` property ([#16548](https://github.com/hashicorp/terraform-provider-azurerm/issues/16548))
* `azurerm_windows_function_app` - add support for `sticky_settings` ([#16546](https://github.com/hashicorp/terraform-provider-azurerm/issues/16546))
* `azurerm_windows_web_app` - add support for `sticky_settings` ([#16546](https://github.com/hashicorp/terraform-provider-azurerm/issues/16546))
* `azurerm_linux_virtual_machine_scale_set` -  the `terminate_notification` property has been renamed to `termination_notification` ([#15570](https://github.com/hashicorp/terraform-provider-azurerm/issues/15570))
* `azurerm_windows_virtual_machine_scale_set` -  the `terminate_notification` property has been renamed to `termination_notification` ([#15570](https://github.com/hashicorp/terraform-provider-azurerm/issues/15570))

BUG FIXES:

* `azurerm_datafactory_dataset_x` - Fix crash around `azure_blob_storage_location.0.dynamic_container_enabled` ([#16514](https://github.com/hashicorp/terraform-provider-azurerm/issues/16514))
* `azurerm_kubernetes_cluster` - allow updates to a cluster running a deprecated version of kubernetes ([#16551](https://github.com/hashicorp/terraform-provider-azurerm/issues/16551))
* `azurerm_resource_policy_remediation` - will no longer try to cancel a completed remediation task during deletion ([#16478](https://github.com/hashicorp/terraform-provider-azurerm/issues/16478))

## 3.3.0 (April 21, 2022)

FEATURES:

* **New Resource**: `azurerm_spring_cloud_container_deployment` ([#16181](https://github.com/hashicorp/terraform-provider-azurerm/issues/16181))

ENHANCEMENTS:

* dependencies: updating to `v0.19.0` of `github.com/tombuildsstuff/giovanni` ([#16460](https://github.com/hashicorp/terraform-provider-azurerm/issues/16460))
* Data Source: `azurerm_kubernetes_cluster` - exporting the `microsoft_defender` block ([#16218](https://github.com/hashicorp/terraform-provider-azurerm/issues/16218))
* Data Source: `azurerm_storage_account` - exporting the `nfsv3_enabled` attribute ([#16404](https://github.com/hashicorp/terraform-provider-azurerm/issues/16404))
* `azurerm_data_factory_linked_service_azure_blob_storage` - support for the `storage_kind` property ([#16403](https://github.com/hashicorp/terraform-provider-azurerm/issues/16403))
* `azurerm_data_factory_linked_service_azure_blob_storage` - support for the `service_principal_linked_key_vault_key` property ([#16414](https://github.com/hashicorp/terraform-provider-azurerm/issues/16414))
* `data_factory_linked_service_sql_server_resource` - support for the `user_name` property ([#16118](https://github.com/hashicorp/terraform-provider-azurerm/issues/16118))
* `azurerm_kubernetes_cluster` - support for the `microsoft_defender` block ([#16218](https://github.com/hashicorp/terraform-provider-azurerm/issues/16218))
* `azurerm_redis_enterprise_cluster` - support for the `linked_database_id` and `linked_database_group_nickname` properties ([#16045](https://github.com/hashicorp/terraform-provider-azurerm/issues/16045))
* `azurerm_spring_cloud_service` - support for the `service_registry_enabled` property ([#16277](https://github.com/hashicorp/terraform-provider-azurerm/issues/16277))
* `azurerm_stream_analytics_output_mssql` - support for the `system_property_columns` property ([#16425](https://github.com/hashicorp/terraform-provider-azurerm/issues/16425))
* `azurerm_stream_analytics_output_servicebus_topic` - support for the `max_batch_count` and `max_writer_count` properties ([#16409](https://github.com/hashicorp/terraform-provider-azurerm/issues/16409))
* `azurerm_stream_analytics_output_table` - support for the `columns_to_remove` property ([#16389](https://github.com/hashicorp/terraform-provider-azurerm/issues/16389))
* `azurerm_virtual_hub_connection` - the `internet_security_enabled` property can now be updated ([#16430](https://github.com/hashicorp/terraform-provider-azurerm/issues/16430))

BUG FIXES:

* `azurerm_cdn_endpoint` - the `origin.http` and `origin.https_ports` properties now have thed efault values of `80` and `443` respectivly ([#16143](https://github.com/hashicorp/terraform-provider-azurerm/issues/16143))
* `azurerm_key_vault_certificate` - now authenticates and manages resources correctly within the US Gov Cloud ([#16455](https://github.com/hashicorp/terraform-provider-azurerm/issues/16455))
* `azurerm_key_vault_key` - now authenticates and manages resources correctly within the US Gov Cloud ([#16455](https://github.com/hashicorp/terraform-provider-azurerm/issues/16455))
* `azurerm_key_vault_managed_storage_account` - now authenticates and manages resources correctly within the US Gov Cloud ([#16455](https://github.com/hashicorp/terraform-provider-azurerm/issues/16455))
* `azurerm_key_vault_secret` - now authenticates and manages resources correctly within the US Gov Cloud ([#16455](https://github.com/hashicorp/terraform-provider-azurerm/issues/16455))
* `azurerm_kubernetes_cluster` - the `role_based_access_control_enabled` property can now be disabled ([#16488](https://github.com/hashicorp/terraform-provider-azurerm/issues/16488))
* `azurerm_linux_function_app` - the `ip_address` property is now correctly set into state when the `service_tag` property is specified ([#16426](https://github.com/hashicorp/terraform-provider-azurerm/issues/16426))
* `azurerm_linux_function_app` - fix a bug in updates to `app_settings` where settings could be lost ([#16442](https://github.com/hashicorp/terraform-provider-azurerm/issues/16442))
* `azurerm_linux_function_app_slot` -  this `ip_address` property is now correctly set into state when the `service_tag` property is specified ([#16426](https://github.com/hashicorp/terraform-provider-azurerm/issues/16426))
* `azurerm_linux_web_app` -  the `ip_address` property is correctly set into state when the `service_tag` property is specified ([#16426](https://github.com/hashicorp/terraform-provider-azurerm/issues/16426))
* `azurerm_linux_web_app` - fix a potential crash when an empty `app_stack` block is used ([#16446](https://github.com/hashicorp/terraform-provider-azurerm/issues/16446))
* `azurerm_linux_web_app_slot` -  the `ip_address` property is now correctly set into state when the `service_tag` property is specified ([#16426](https://github.com/hashicorp/terraform-provider-azurerm/issues/16426))
* `azurerm_linux_web_app_slot` - fix a potential crash when an empty `app_stack` block is used ([#16446](https://github.com/hashicorp/terraform-provider-azurerm/issues/16446))
* `azurerm_sentinel_alert_rule_fusion` - will no longer send the `etag` property during updates as it is longer required ([#16428](https://github.com/hashicorp/terraform-provider-azurerm/issues/16428))
* `azurerm_sentinel_alert_rule_machine_learning_behavior_analytics` - will no longer send the `etag` property during updates as it is longer required ([#16428](https://github.com/hashicorp/terraform-provider-azurerm/issues/16428))
* `azurerm_sentinel_alert_rule_ms_security_incident` - will no longer send the `etag` property during updates as it is longer required ([#16428](https://github.com/hashicorp/terraform-provider-azurerm/issues/16428))
* `azurerm_sentinel_alert_rule_scheduled` - will no longer send the `etag` property during updates as it is longer required ([#16428](https://github.com/hashicorp/terraform-provider-azurerm/issues/16428))
* `azurerm_sentinel_data_connector_aws_cloud_trail` - will no longer send the `etag` property during updates as it is longer required ([#16428](https://github.com/hashicorp/terraform-provider-azurerm/issues/16428))
* `azurerm_sentinel_data_connector_microsoft_cloud_app_security` - will no longer send the `etag` property during updates as it is longer required ([#16428](https://github.com/hashicorp/terraform-provider-azurerm/issues/16428))
* `azurerm_sentinel_data_connector_office_365` - will no longer send the `etag` property during updates as it is longer required ([#16428](https://github.com/hashicorp/terraform-provider-azurerm/issues/16428))
* `azurerm_storage_account` - will now update `identity` before `customer_managed_key` enabling adding a new identity with access to the CMK ([#16419](https://github.com/hashicorp/terraform-provider-azurerm/issues/16419))
* `azurerm_subnet` - the `address_prefixes` property is now (explicitly) required ([#16402](https://github.com/hashicorp/terraform-provider-azurerm/issues/16402))
* `azurerm_windows_function_app` - the `ip_address` property is now correctly set into state when the `service_tag` property is specified ([#16426](https://github.com/hashicorp/terraform-provider-azurerm/issues/16426))
* `azurerm_windows_function_app` - fix a bug in updates to `app_settings` where settings could be lost ([#16442](https://github.com/hashicorp/terraform-provider-azurerm/issues/16442))
* `azurerm_windows_function_app_slot` - the `ip_address` property is now correctly set into state when the `service_tag` property is specified ([#16426](https://github.com/hashicorp/terraform-provider-azurerm/issues/16426))
* `azurerm_windows_web_app` - the `ip_address` property is now correctly set into state when the `service_tag` property is specified ([#16426](https://github.com/hashicorp/terraform-provider-azurerm/issues/16426))
* `azurerm_windows_web_app` - prevent a potential crash when an empty `app_stack` block is used ([#16446](https://github.com/hashicorp/terraform-provider-azurerm/issues/16446))
* `azurerm_windows_web_app_slot` - the `ip_address` property is now correctly set into state when the `service_tag` property is specified ([#16426](https://github.com/hashicorp/terraform-provider-azurerm/issues/16426))
* `azurerm_windows_web_app_slot` - prevent a potential crash when an empty `app_stack` block is used ([#16446](https://github.com/hashicorp/terraform-provider-azurerm/issues/16446))

## 3.2.0 (April 14, 2022)

FEATURES:

* **New Datasource**: `azurerm_kusto_database` ([#16180](https://github.com/hashicorp/terraform-provider-azurerm/issues/16180))
* **New Resource**: `azurerm_container_connected_registry` ([#15731](https://github.com/hashicorp/terraform-provider-azurerm/issues/15731))
* **New Resource**: `azurerm_managment_group_policy_exemption` ([#16293](https://github.com/hashicorp/terraform-provider-azurerm/issues/16293))
* **New Resource**: `azurerm_resource_group_policy_exemption` ([#16293](https://github.com/hashicorp/terraform-provider-azurerm/issues/16293))
* **New Resource**: `azurerm_resource_policy_exemption` ([#16293](https://github.com/hashicorp/terraform-provider-azurerm/issues/16293))
* **New Resource**: `azurerm_stream_analytics_job_schedule` ([#16349](https://github.com/hashicorp/terraform-provider-azurerm/issues/16349))
* **New Resource**: `azurerm_subscription_policy_exemption` ([#16293](https://github.com/hashicorp/terraform-provider-azurerm/issues/16293))

ENHANCEMENTS:

* Data Source: `azurerm_stream_analytics_job` - support for the `last_output_time`, `start_mode`, and `start_time` properties ([#16349](https://github.com/hashicorp/terraform-provider-azurerm/issues/16349))
* `azurerm_container_group` - support for the `init_container` block ([#16204](https://github.com/hashicorp/terraform-provider-azurerm/issues/16204))
* `azurerm_machine_learning_workspace` - renamed the `public_network_access_enabled` property to `public_access_behind_virtual_network_enabled` to better reflect what this property does ([#16288](https://github.com/hashicorp/terraform-provider-azurerm/issues/16288))
* `azurerm_media_streaming_endpoint` support Standard Streaming Endpoints ([#16304](https://github.com/hashicorp/terraform-provider-azurerm/issues/16304))
* `azurerm_cdn_endpoint` - the `url_path_condition` property now allows the `RegEx` and `Wildcard` values ([#16385](https://github.com/hashicorp/terraform-provider-azurerm/issues/16385))

BUG FIXES:

* Data Source: `azurerm_log_analytics_linked_storage_account` - correctly set the `data_source_type` property ([#16313](https://github.com/hashicorp/terraform-provider-azurerm/issues/16313))
* `azurerm_lb_outbound_rule` - allow `0` for the `allocated_outbound_ports` property ([#16369](https://github.com/hashicorp/terraform-provider-azurerm/issues/16369))
* `azurerm_mysql_flexible_server` - `backup_retention_days` can now be set any value from `1`-`35` ([#16312](https://github.com/hashicorp/terraform-provider-azurerm/issues/16312))
* `azurerm_sentinel_watchlist` - support for the required property `item_search_key` ([#15861](https://github.com/hashicorp/terraform-provider-azurerm/issues/15861))
* `azurerm_vpn_server_configuration` - the `server_root_certificate` property is now optional ([#16366](https://github.com/hashicorp/terraform-provider-azurerm/issues/16366))
* `azurerm_storage_data_lake_gen2_path` - support `$superuser` as an option for `owner` and `group` ([#16370](https://github.com/hashicorp/terraform-provider-azurerm/issues/16370))
* `azurerm_eventhub_namespace` - can now be updated when customer managed keys are being used ([#16371](https://github.com/hashicorp/terraform-provider-azurerm/issues/16371))
* `azurerm_postgresql_flexible_server` - `high_availability` blocks can now be added and removed ([#16328](https://github.com/hashicorp/terraform-provider-azurerm/issues/16328))

## 3.1.0 (April 07, 2022)

FEATURES:

* **New Resource**: `azurerm_container_registry_agent_pool` ([#16258](https://github.com/hashicorp/terraform-provider-azurerm/issues/16258))

ENHANCEMENTS:

* dependencies: updating to `v63.0.0` of `github.com/Azure/azure-sdk-for-go` ([#16147](https://github.com/hashicorp/terraform-provider-azurerm/issues/16147))
* dependencies: updating `digitaltwins` to use API Version `2020-12-01` ([#16044](https://github.com/hashicorp/terraform-provider-azurerm/issues/16044))
* dependencies: updating `streamanalytics` to use API Version `2020-03-01` ([#16270](https://github.com/hashicorp/terraform-provider-azurerm/issues/16270))
* provider: upgrading to Go `1.18` ([#16247](https://github.com/hashicorp/terraform-provider-azurerm/issues/16247))
* Data Source: `azurerm_kubernetes_cluster` - support for the `oidc_issuer_enabled` and `oidc_issuer_url` properties [[#16130](https://github.com/hashicorp/terraform-provider-azurerm/issues/16130)] 
* Data Source: `azurerm_service_plan` - add support for `zone_balancing_enabled` ([#16156](https://github.com/hashicorp/terraform-provider-azurerm/issues/16156))
* `azurerm_application_gateway` - add `KNOWN-CVES` to accepted values for the `rule_group_name` property ([#16080](https://github.com/hashicorp/terraform-provider-azurerm/issues/16080))
* `azurerm_automation_account` - the `dsc_primary_access_key` and `dsc_secondary_access_key` properties are now marked as sensitive ([#16161](https://github.com/hashicorp/terraform-provider-azurerm/issues/16161))
* `azurerm_cognitive_account` - support for the `custom_question_answering_search_service_id` property ([#15804](https://github.com/hashicorp/terraform-provider-azurerm/issues/15804))
* `azurerm_consumption_budget_management_group` - support for `SubscriptionID` and `SubscriptionName` options in the `dimension` block ([#16074](https://github.com/hashicorp/terraform-provider-azurerm/issues/16074))
* `azurerm_cosmosdb_gremlin_graph` - the property `indexing_mode` is now case-sensitive ([#16152](https://github.com/hashicorp/terraform-provider-azurerm/issues/16152))
* `azurerm_cosmosdb_sql_container` - the property `indexing_mode` is now case-sensitive ([#16152](https://github.com/hashicorp/terraform-provider-azurerm/issues/16152))
* `azurerm_dedicated_host` - support for the the `DSv3-Type4` and `ESv3-Type4` SKUs ([#16253](https://github.com/hashicorp/terraform-provider-azurerm/issues/16253))
* `azurerm_kubernetes_cluster` - support for the `oidc_issuer_enabled` and `oidc_issuer_url` properties [[#16130](https://github.com/hashicorp/terraform-provider-azurerm/issues/16130)] 
* `azurerm_kubernetes_cluster` - the `network_profile` block now supports the `ip_versions` property ([#16088](https://github.com/hashicorp/terraform-provider-azurerm/issues/16088))
* `azurerm_mssql_database` - support for the `ledger_enabled` property ([#16214](https://github.com/hashicorp/terraform-provider-azurerm/issues/16214))
* `azurerm_service_plan` - support for the `zone_balancing_enabled` property ([#16156](https://github.com/hashicorp/terraform-provider-azurerm/issues/16156))
* `azurerm_servicebus_namespace` - support for the `customer_managed_key` block ([#15601](https://github.com/hashicorp/terraform-provider-azurerm/issues/15601))
* `azurerm_web_application_firewall_policy` - add `KNOWN-CVES` to accepted values for `rule_group_name` ([#16080](https://github.com/hashicorp/terraform-provider-azurerm/issues/16080))
* `azurerm_servicebus_namespace` - add support for the `local_auth_enabled` property ([#16268](https://github.com/hashicorp/terraform-provider-azurerm/issues/16268))

BUG FIXES:

* `azurerm_api_management_api_operation_tag` - now retrieves tags from the correct API ([#16006](https://github.com/hashicorp/terraform-provider-azurerm/issues/16006))
* `azurerm_api_management_api_operation` - prevent a potential panic when parsing `representation` ([#14848](https://github.com/hashicorp/terraform-provider-azurerm/issues/14848))
* `azurerm_application_gateway` - a `frontend_ip_configuration` blocks can now be updated ([#16132](https://github.com/hashicorp/terraform-provider-azurerm/issues/16132))
* `azurerm_application_insights` - remove the disable logic for the created Action Groups ([#16170](https://github.com/hashicorp/terraform-provider-azurerm/issues/16170))
* `azurerm_cosmosdb_sql_container` - disabling the `analytical_storage_ttl` property now forces a new resoruce to be created ([#16229](https://github.com/hashicorp/terraform-provider-azurerm/issues/16229))
* `azurerm_linux_function_app` - only one of `application_insights_key` or `application_insights_connection_string` needs to be optionally specified ([#16134](https://github.com/hashicorp/terraform-provider-azurerm/issues/16134))
* `azurerm_linux_function_app_slot` - only one of `application_insights_key` or `application_insights_connection_string` needs to be optionally specified ([#16134](https://github.com/hashicorp/terraform-provider-azurerm/issues/16134))
* `azurerm_windows_function_app` - fix the import check for Service Plan OS type ([#16164](https://github.com/hashicorp/terraform-provider-azurerm/issues/16164))
* `azurerm_linux_web_app_slot ` - fix `container_registry_managed_identity_client_id` property validation ([#16149](https://github.com/hashicorp/terraform-provider-azurerm/issues/16149))
* `azurerm_windows_web_app` - add support for `dotnetcore` in site metadata property `current_stack` ([#16129](https://github.com/hashicorp/terraform-provider-azurerm/issues/16129))
* `azurerm_windows_web_app` - fix docker `windowsFXVersion` when `docker_container_registry` is specified ([#16192](https://github.com/hashicorp/terraform-provider-azurerm/issues/16192))
* `azurerm_windows_web_app_slot` - add support for `dotnetcore` in site metadata property `current_stack` ([#16129](https://github.com/hashicorp/terraform-provider-azurerm/issues/16129))
* `azurerm_windows_web_app_slot` - fix docker `windowsFXVersion` when `docker_container_registry` is specified ([#16192](https://github.com/hashicorp/terraform-provider-azurerm/issues/16192))
* `azurerm_storage_data_lake_gen2_filesystem` - add support for `$superuser` in `group` and `owner` properties ([#16215](https://github.com/hashicorp/terraform-provider-azurerm/issues/16215))

## 3.0.2 (March 26, 2022)

BUG FIXES:

* `azurerm_cosmosdb_account` - prevent a panic when the API returns an empty list of read or write locations ([#16031](https://github.com/hashicorp/terraform-provider-azurerm/issues/16031))
* `azurerm_cdn_endpoint` - prevent a panic when there is an empty `country_codes` property ([#16066](https://github.com/hashicorp/terraform-provider-azurerm/issues/16066))
* `azurerm_key_vault` - fix the `authorizer was not an auth.CachedAuthorizer ` error ([#16078](https://github.com/hashicorp/terraform-provider-azurerm/issues/16078))
* `azurerm_linux_function_app` - correctly update storage settings when using MSI ([#16046](https://github.com/hashicorp/terraform-provider-azurerm/issues/16046))
* `azurerm_managed_disk` - changing the `zone` property now correctly creates a new resource ([#16070](https://github.com/hashicorp/terraform-provider-azurerm/issues/16070))
* `azurerm_resource_group` - will now during deletion if there are still resources found in the group it will wait a little bit and check again to handle eventually consistancy bugs ([#16073](https://github.com/hashicorp/terraform-provider-azurerm/issues/16073))
* `azurerm_windows_function_app` - correctly update the storage settings when using MSI authentication ([#16046](https://github.com/hashicorp/terraform-provider-azurerm/issues/16046))

## 3.0.1 (March 24, 2022)

BUG FIXES:

* provider: the `prevent_deletion_if_contains_resources` feature flag within the `resource_group` block now defaults to `true` ([#16021](https://github.com/hashicorp/terraform-provider-azurerm/issues/16021))

## 3.0.0 (March 24, 2022)

NOTES:

* **Major Version**: Version 3.0 of the Azure Provider is a major version - some behaviours have changed and some deprecated fields/resources have been removed - please refer to [the 3.0 upgrade guide for more information](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/3.0-upgrade-guide).
* When upgrading to v3.0 of the AzureRM Provider, we recommend upgrading to the latest version of Terraform Core ([which can be found here](https://www.terraform.io/downloads)) - the next major release of the AzureRM Provider (v4.0) will require Terraform 1.0 or later.

FEATURES:

* **New Data Source**: `azurerm_healthcare_workspace` ([#15759](https://github.com/hashicorp/terraform-provider-azurerm/issues/15759))
* **New Data Source**: `azurerm_key_vault_encrypted_value` ([#15873](https://github.com/hashicorp/terraform-provider-azurerm/issues/15873))
* **New Data Source**: `azurerm_managed_api` ([#15797](https://github.com/hashicorp/terraform-provider-azurerm/issues/15797))
* **New Resource**: `azurerm_api_connection` ([#15797](https://github.com/hashicorp/terraform-provider-azurerm/issues/15797))
* **New Resource**: `azurerm_healthcare_workspace` ([#15759](https://github.com/hashicorp/terraform-provider-azurerm/issues/15759))
* **New Resource**: `azurerm_stream_analytics_function_javascript_uda` ([#15831](https://github.com/hashicorp/terraform-provider-azurerm/issues/15831))
* **New Resource**: `azurerm_security_center_server_vulnerability_assessment_virtual_machine` ([#15747](https://github.com/hashicorp/terraform-provider-azurerm/issues/15747))

ENHANCEMENTS:

* dependencies: updating to `v62.3.0` of `github.com/Azure/azure-sdk-for-go` ([#15927](https://github.com/hashicorp/terraform-provider-azurerm/issues/15927))
* dependencies: updating to `v0.26.0` of `github.com/hashicorp/go-azure-helpers` ([#15889](https://github.com/hashicorp/terraform-provider-azurerm/issues/15889))
* dependencies: updating `appplatform` to API Version `2022-01-01-preview` ([#15597](https://github.com/hashicorp/terraform-provider-azurerm/issues/15597))
* provider: MSAL (and Microsoft Graph) is now used for authentication instead of ADAL (and Azure Active Directory Graph) ([#12443](https://github.com/hashicorp/terraform-provider-azurerm/issues/12443))
* provider: all (non-deprecated) resources now validate the Resource ID during import ([#15989](https://github.com/hashicorp/terraform-provider-azurerm/issues/15989))
* provider: added a new feature flag within the `api_management` block for `recover_soft_deleted`, for configuring whether a soft-deleted `azurerm_api_management` should be recovered during creation ([#15871](https://github.com/hashicorp/terraform-provider-azurerm/issues/15871))
* provider: added a new feature flag within the `key_vault` block for `recover_soft_deleted_certificates`, for configuring whether a soft-deleted `azurerm_key_vault_certificate` should be recovered during creation ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `key_vault` block for `purge_soft_deleted_certificates_on_destroy`, for configuring whether a deleted `azurerm_key_vault_certificate` should be purged during deletion ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `key_vault` block for `recover_soft_deleted_keys`, for configuring whether a soft-deleted `azurerm_key_vault_key` should be recovered during creation ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `key_vault` block for `purge_soft_deleted_keys_on_destroy`, for configuring whether a deleted `azurerm_key_vault_key` should be purged during deletion ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `key_vault` block for `recover_soft_deleted_secrets`, for configuring whether a soft-deleted `azurerm_key_vault_secret` should be recovered during creation ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `key_vault` block for `purge_soft_deleted_secrets_on_destroy`, for configuring whether a deleted `azurerm_key_vault_secret` should be purged during deletion ([#10273](https://github.com/hashicorp/terraform-provider-azurerm/issues/10273))
* provider: added a new feature flag within the `resource_group` block for `prevent_deletion_if_contains_resources`, for configuring whether Terraform should prevent the deletion of a Resource Group which still contains items ([#13777](https://github.com/hashicorp/terraform-provider-azurerm/issues/13777))
* provider: the feature flag `permanently_delete_on_destroy` within the `log_analytics_workspace` block now defaults to `true` ([#15948](https://github.com/hashicorp/terraform-provider-azurerm/issues/15948))
* Resources supporting Availability Zones: Zones are now treated consistently across the Provider and the field within Terraform has been renamed to either `zone` (for a single Zone) or `zones` (where multiple can be defined) - the complete list of resources can be found in the 3.0 Upgrade Guide ([#14588](https://github.com/hashicorp/terraform-provider-azurerm/issues/14588))
* Resources supporting Managed Identity: Identity blocks are now treated consistently across the Provider - the complete list of resources can be found in the 3.0 Upgrade Guide ([#15187](https://github.com/hashicorp/terraform-provider-azurerm/issues/15187))
* provider: removing the `network` and `relaxed_locking` feature flags, since this is now enabled by default ([#15719](https://github.com/hashicorp/terraform-provider-azurerm/issues/15719))
* Data Source: `azurerm_linux_function_app` - support for the `storage_key_vault_secret_id` property ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* Data Source: `azurerm_storage_account_sas` - now exports the `tag` and `filter` attributes ([#15863](https://github.com/hashicorp/terraform-provider-azurerm/issues/15863))
* Data Source: `azurerm_windows_function_app` - support for `storage_key_vault_secret_id` property ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* `azurerm_application_insights` - can now disable Rule and Action Groups that are automatically created ([#15892](https://github.com/hashicorp/terraform-provider-azurerm/issues/15892))
* `azurerm_cdn_endpoint` - the `host_name` property has been renamed to `fqdn` ([#15992](https://github.com/hashicorp/terraform-provider-azurerm/issues/15992))
* `azurerm_eventgrid_system_topic_event_subscription` - support for the `delivery_property` property ([#15559](https://github.com/hashicorp/terraform-provider-azurerm/issues/15559))
* `azurerm_iothub` - add support for the `authentication_type` and `identity_id` properties in the `file_upload` block ([#15874](https://github.com/hashicorp/terraform-provider-azurerm/issues/15874))
* `azurerm_kubernetes_cluster` - the `kube_admin_config` block is now marked as sensitive in addition to all items within it ([#4105](https://github.com/hashicorp/terraform-provider-azurerm/issues/4105))
* `azurerm_kubernetes_cluster` - add support for the `key_vault_secrets_provider` and `open_service_mesh_enabled` property in Azure China and Azure Government ([#15878](https://github.com/hashicorp/terraform-provider-azurerm/issues/15878))
* `azurerm_linux_function_app` - add support for the `storage_key_vault_secret_id` property ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* `azurerm_linux_function_app` - updating the read timeout to be `5m` ([#15867](https://github.com/hashicorp/terraform-provider-azurerm/issues/15867))
* `azurerm_linux_function_app` - support for node version `16` preview ([#15884](https://github.com/hashicorp/terraform-provider-azurerm/issues/15884))
* `azurerm_linux_function_app` - add support for `use_dotnet_isolated_runtime` ([#15969](https://github.com/hashicorp/terraform-provider-azurerm/issues/15969))
* `azurerm_linux_function_app_slot` - add support for `use_dotnet_isolated_runtime` ([#15969](https://github.com/hashicorp/terraform-provider-azurerm/issues/15969))
* `azurerm_linux_function_app_slot` - add support for `storage_key_vault_secret_id` ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* `azurerm_linux_function_app_slot` - updating the read timeout to be 5m ([#15867](https://github.com/hashicorp/terraform-provider-azurerm/issues/15867))
* `azurerm_linux_virtual_machine` - support for the `termination_notification` property ([#14933](https://github.com/hashicorp/terraform-provider-azurerm/issues/14933))
* `azurerm_linux_virtual_machine ` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_linux_virtual_machine_scale_set` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_linux_web_app` - support for PHP version 8.0 ([#15933](https://github.com/hashicorp/terraform-provider-azurerm/issues/15933))
* `azurerm_loadbalancer` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_managed_disk` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_management_group_policy_assignment` - support for User Assigned Identities ([#15376](https://github.com/hashicorp/terraform-provider-azurerm/issues/15376))
* `azurerm_mssql_server` - the `minimum_tls_version` property now defaults to `1.2` ([#10276](https://github.com/hashicorp/terraform-provider-azurerm/issues/10276))
* `azurerm_mysql_server` - the `ssl_minimal_tls_version_enforced` property now defaults to `1.2` ([#10276](https://github.com/hashicorp/terraform-provider-azurerm/issues/10276))
* `azurerm_network_interface` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_network_security_rule` - no longer locks on the network security group name ([#15719](https://github.com/hashicorp/terraform-provider-azurerm/issues/15719))
* `azurerm_postgresql_server` - the `ssl_minimal_tls_version_enforced` property now defaults to `1.2` ([#10276](https://github.com/hashicorp/terraform-provider-azurerm/issues/10276))
* `azurerm_public_ip` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_redis_cache` - the `minimum_tls_version` property  now defaults to `1.2` ([#10276](https://github.com/hashicorp/terraform-provider-azurerm/issues/10276))
* `azurerm_resource_group` - Terraform now checks during the deletion of a Resource Group if there's any items remaining and will raise an error if so by default (to avoid deleting items unintentionally). This behaviour can be controlled using the `prevent_deletion_if_contains_resources` feature-flag within the `resource_group` block within the `features` block. ([#13777](https://github.com/hashicorp/terraform-provider-azurerm/issues/13777))
* `azurerm_resource_group_policy_assignment` - support for User Assigned Identities ([#15376](https://github.com/hashicorp/terraform-provider-azurerm/issues/15376))
* `azurerm_resource_policy_assignment` - support for User Assigned Identities ([#15376](https://github.com/hashicorp/terraform-provider-azurerm/issues/15376))
* `azurerm_sentinel_alert_rule_scheduled` - support for `alert_details_override` ([#15901](https://github.com/hashicorp/terraform-provider-azurerm/issues/15901))
* `azurerm_sentinel_alert_rule_scheduled` - support for `entity_mapping` [[#15901](https://github.com/hashicorp/terraform-provider-azurerm/issues/15901)] 
* `azurerm_sentinel_alert_rule_scheduled` - support for `custom_details` ([#15901](https://github.com/hashicorp/terraform-provider-azurerm/issues/15901))
* `azurerm_sentinel_alert_rule_scheduled` - support for `group_by_alert_details` ([#15901](https://github.com/hashicorp/terraform-provider-azurerm/issues/15901))
* `azurerm_sentinel_alert_rule_scheduled` - support for `group_by_custom_details` ([#15901](https://github.com/hashicorp/terraform-provider-azurerm/issues/15901))
* `azurerm_site_recovery_replicated_vm` - support for the `target_availability_zone` property ([#15617](https://github.com/hashicorp/terraform-provider-azurerm/issues/15617))
* `azurerm_shared_image` - support for the `support_accelerated_network` property ([#15562](https://github.com/hashicorp/terraform-provider-azurerm/issues/15562))
* `azurerm_static_site` - the `identity` property now supports `SystemAssigned` and `UserAssigned` ([#15834](https://github.com/hashicorp/terraform-provider-azurerm/issues/15834))
* `azurerm_storage_account` - the `allow_blob_public_access` property has been renamed to `allow_nested_items_to_be_public` to better represent what is being enabled ([#12689](https://github.com/hashicorp/terraform-provider-azurerm/issues/12689))
* `azurerm_storage_account` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_storage_account` - `ZRS` is no longer supported when using `StorageV1` ([#16004](https://github.com/hashicorp/terraform-provider-azurerm/issues/16004))
* `azurerm_storage_account` - the `min_tls_version` property now defaults to `1.2` ([#10276](https://github.com/hashicorp/terraform-provider-azurerm/issues/10276))
* `azurerm_storage_share` - `quota` is now required ([#15982](https://github.com/hashicorp/terraform-provider-azurerm/issues/15982))
* `azurerm_subscription_policy_assignment` - support for User Assigned Identities ([#15376](https://github.com/hashicorp/terraform-provider-azurerm/issues/15376))
* `azurerm_virtual_network` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_virtual_network_gateway` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_virtual_hub` - support for the `virtual_router_asn` and `virtual_router_ips` properties ([#15741](https://github.com/hashicorp/terraform-provider-azurerm/issues/15741))
* `azurerm_windows_function_app` - add support for `storage_key_vault_secret_id` property ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* `azurerm_windows_function_app` - updating the read timeout to be `5m` ([#15867](https://github.com/hashicorp/terraform-provider-azurerm/issues/15867))
* `azurerm_windows_function_app` node version validation string can not be prefixed with `~` ([#15884](https://github.com/hashicorp/terraform-provider-azurerm/issues/15884))
* `azurerm_windows_function_app` support for node version `16` preview support ([#15884](https://github.com/hashicorp/terraform-provider-azurerm/issues/15884))
* `azurerm_windows_function_app` - add support for `use_dotnet_isolated_runtime` ([#15969](https://github.com/hashicorp/terraform-provider-azurerm/issues/15969))
* `azurerm_windows_function_app_slot` - add support for `use_dotnet_isolated_runtime` ([#15969](https://github.com/hashicorp/terraform-provider-azurerm/issues/15969))
* `azurerm_windows_function_app_slot` - add support for the `storage_key_vault_secret_id` property ([#15793](https://github.com/hashicorp/terraform-provider-azurerm/issues/15793))
* `azurerm_windows_function_app_slot` - updating the read timeout to be 5m ([#15867](https://github.com/hashicorp/terraform-provider-azurerm/issues/15867))
* `azurerm_windows_virtual_machine` - support for the `termination_notification` property ([#14933](https://github.com/hashicorp/terraform-provider-azurerm/issues/14933))
* `azurerm_windows_virtual_machine` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))
* `azurerm_windows_virtual_machine_scale_set` - support for the `edge_zone` property ([#15890](https://github.com/hashicorp/terraform-provider-azurerm/issues/15890))

BUG FIXES:

* provider: the `recover_soft_deleted_key_vaults` feature flag within the `key_vault` block now defaults to `true` ([#15984](https://github.com/hashicorp/terraform-provider-azurerm/issues/15984))
* provider: the `purge_soft_delete_on_destroy ` feature flag within the `key_vault` block now defaults to `true` [[#15984](https://github.com/hashicorp/terraform-provider-azurerm/issues/15984)] 
* `azurerm_app_configuration_feature` - detecting that the key is gone when the App Configuration has been deleted ([#15973](https://github.com/hashicorp/terraform-provider-azurerm/issues/15973))
* `azurerm_app_configuration_key` - detecting that the key is gone when the App Configuration has been deleted ([#15973](https://github.com/hashicorp/terraform-provider-azurerm/issues/15973))
* `azurerm_application_gateway` - the `backend_address_pool` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the field `fqdns` within the `backend_address_pool` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the field `ip_addresses` within the `backend_address_pool` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `backend_http_settings` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `frontend_port` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the field `host_names` within the `frontend_port` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `http_listener` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `private_endpoint_connection` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `private_link_configuration` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `probe` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `redirect_configuration` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `request_routing_rule` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_application_gateway` - the `ssl_certificate` block is now a Set rather than a List ([#6896](https://github.com/hashicorp/terraform-provider-azurerm/issues/6896))
* `azurerm_container_registry` - validate the `georepliactions` property does not include the location of the Container Registry ([#15847](https://github.com/hashicorp/terraform-provider-azurerm/issues/15847))
* `azurerm_cosmosdb_mongo_collection` - the `default_ttl_seconds` property can now be set to `-1` ([#15736](https://github.com/hashicorp/terraform-provider-azurerm/issues/15736))
* `azurerm_eventhub` - prevent panic when the `capture_description` block is removed ([#15930](https://github.com/hashicorp/terraform-provider-azurerm/issues/15930))
* `azurerm_key_vault_access_policy` - validating the Resource ID during import ([#15989](https://github.com/hashicorp/terraform-provider-azurerm/issues/15989))
* `azurerm_linux_function_app` - fixed update handling of `app_settings` for `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#15907](https://github.com/hashicorp/terraform-provider-azurerm/issues/15907))
* `azurerm_linux_function_app_slot` - fixed update handling of `app_settings` for `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#15907](https://github.com/hashicorp/terraform-provider-azurerm/issues/15907))
* `azurerm_local_network_gateway` - fix for `address_space` cannot be updated ([#15159](https://github.com/hashicorp/terraform-provider-azurerm/issues/15159))
* `azurerm_log_analytics_cluster_customer_managed_key` - detecting when the Customer Managed Key has been removed ([#15973](https://github.com/hashicorp/terraform-provider-azurerm/issues/15973))
* `azurerm_mssql_database_vulnerability_assessment_rule_baseline` - prevent the resource from being replaced every apply ([#14759](https://github.com/hashicorp/terraform-provider-azurerm/issues/14759))
* `azurerm_security_center_auto_provisioning ` - validating the Resource ID during import [[#15989](https://github.com/hashicorp/terraform-provider-azurerm/issues/15989)] 
* `azurerm_security_center_setting` - changing the `setting_name` property now forces a new resource ([#15983](https://github.com/hashicorp/terraform-provider-azurerm/issues/15983))
* `azurerm_synapse_workspace` - fixing a bug where workspaces created from a Dedicated SQL Pool / SQL Data Warehouse couldn't be retrieved ([#15829](https://github.com/hashicorp/terraform-provider-azurerm/issues/15829))
* `azurerm_synapse_workspace_key` - keys can now be correctly rotated ([#15897](https://github.com/hashicorp/terraform-provider-azurerm/issues/15897))
* `azurerm_windows_function_app` - fixed update handling of `app_settings` for `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#15907](https://github.com/hashicorp/terraform-provider-azurerm/issues/15907))
* `azurerm_windows_function_app_slot` - fixed update handling of `app_settings` for `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#15907](https://github.com/hashicorp/terraform-provider-azurerm/issues/15907))

---

For information on changes between the v2.99.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v1.44.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
