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