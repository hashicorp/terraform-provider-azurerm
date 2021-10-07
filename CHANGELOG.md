## 2.80 (Unreleased)

BUG FIXES:

* `azurerm_function_app` - fix regressions in function app storage introduced in v2.77 [GH-13580]
* `azurerm_managed_application` - fixed typecasting bug [GH-13641]

IMPROVEMENTS:

* Data Source `azurerm_public_ips` - Deprecate `attached` for `attachment_status` to improve filtering [GH-13500]
* Data Source `azurerm_public_ips` - Return public IPs associated with NAT gateways when `attached` set to `true` or `attachment_status` set to `Attached` [GH-13610]
* `azurerm_stream_analytics_output_eventhub` - support for the `partition_key` property [GH-13562]
* `azurerm_managed_disk` - support for the `logical_sector_size` property [GH-13637]
* `azurerm_kusto_eventhub_data_connection supports` - support for `identity_id` property [GH-13488]

## 2.79.1 (October 01, 2021)

BUG FIXES: 

* `azurerm_managed_disk` - the `max_shares` propety is now `Computed` to account for managed disks that are already managed by Terraform ([#13587](https://github.com/hashicorp/terraform-provider-azurerm/issues/13587))

## 2.79.0 (October 01, 2021)

FEATURES: 

* **New Resource:** `azurerm_app_configuration_feature` ([#13452](https://github.com/hashicorp/terraform-provider-azurerm/issues/13452))
* **New Resource:** `azurerm_logic_app_standard` ([#13196](https://github.com/hashicorp/terraform-provider-azurerm/issues/13196))

IMPROVEMENTS:

* Data Source: `azurerm_key_vault_certificate` - exporting the `expires` and `not_before` attributes ([#13527](https://github.com/hashicorp/terraform-provider-azurerm/issues/13527))
* Data Source: `azurerm_key_vault_certificate_data` - exporting the `not_before` attribute ([#13527](https://github.com/hashicorp/terraform-provider-azurerm/issues/13527))
* `azurerm_communication_service` - export the `primary_connection_string`, `secondary_connection_string`, `primary_key`, and `secondary_key` attributes ([#13549](https://github.com/hashicorp/terraform-provider-azurerm/issues/13549))
* `azurerm_consumption_budget_subscription`  support for the `Forecasted` threshold type ([#13567](https://github.com/hashicorp/terraform-provider-azurerm/issues/13567))
* `azurerm_consumption_budget_resource_group  support for the `Forecasted` threshold type ([#13567](https://github.com/hashicorp/terraform-provider-azurerm/issues/13567))
* `azurerm_managed_disk` - support for the `max_shares` property ([#13571](https://github.com/hashicorp/terraform-provider-azurerm/issues/13571))
* `azurerm_mssql_database` - will now update replicated databases SKUs first ([#13478](https://github.com/hashicorp/terraform-provider-azurerm/issues/13478))
* `azurerm_virtual_hub_connection` - optimized state change refresh function ([#13548](https://github.com/hashicorp/terraform-provider-azurerm/issues/13548))

BUG FIXES:

* `azurerm_cosmosdb_account` - the `mongo_server_version` can now be changed without creating a new resouce ([#13520](https://github.com/hashicorp/terraform-provider-azurerm/issues/13520))
* `azurerm_iothub` - correctly suppress diffs for the `connection_string` property ([#13517](https://github.com/hashicorp/terraform-provider-azurerm/issues/13517))
* `azurerm_kubernetes_cluster` - explicitly setting `upgrade_channel` to `None` when it's unset to workaround a breaking behavioural change in AKS ([#13493](https://github.com/hashicorp/terraform-provider-azurerm/issues/13493))
* `azurerm_linux_virtual_machine_scale_set` - will not correctly ignore the `protected_setting` block withing the `extension` block ([#13440](https://github.com/hashicorp/terraform-provider-azurerm/issues/13440))
* `azurerm_windows_virtual_machine_scale_set` - will not correctly ignore the `protected_setting` block withing the `extension` block ([#13440](https://github.com/hashicorp/terraform-provider-azurerm/issues/13440))
* `azurerm_app_configuration_key` - correctly set the `etag` property ([#13534](https://github.com/hashicorp/terraform-provider-azurerm/issues/13534))

## 2.78.0 (September 23, 2021)

UPGRADE NOTES

* The `azurerm_data_factory_dataset_snowflake` has been updated to set the correct `schema_column` api property with the correct schema - to retain the old behaviour please switch to the `structure_column` property ([#13344](https://github.com/hashicorp/terraform-provider-azurerm/issues/13344))

FEATURES: 

* **New Resource:** `azurerm_frontdoor_rules_engine` ([#13249](https://github.com/hashicorp/terraform-provider-azurerm/issues/13249))
* **New Resource:** `azurerm_key_vault_managed_storage_account` ([#13271](https://github.com/hashicorp/terraform-provider-azurerm/issues/13271))
* **New Resource:** `azurerm_key_vault_managed_storage_account_sas_token_definition` ([#13271](https://github.com/hashicorp/terraform-provider-azurerm/issues/13271))
* **New Resource:** `azurerm_mssql_failover_group` ([#13446](https://github.com/hashicorp/terraform-provider-azurerm/issues/13446))
* **New Resource:** `azurerm_synapse_sql_pool_extended_auditing_policy` ([#12952](https://github.com/hashicorp/terraform-provider-azurerm/issues/12952))
* **New Resource:** `azurerm_synapse_workspace_extended_auditing_policy` ([#12952](https://github.com/hashicorp/terraform-provider-azurerm/issues/12952))

ENHANCEMENTS:

* upgrading `iothub` to API Version `2021-03-31` ([#13324](https://github.com/hashicorp/terraform-provider-azurerm/issues/13324))
* `data.azurerm_private_endpoint_connection` - Export `network_interface` attributes from private endpoints ([#13421](https://github.com/hashicorp/terraform-provider-azurerm/issues/13421))
* `azurerm_app_service` - support for the `vnet_route_all_enabled` property ([#13310](https://github.com/hashicorp/terraform-provider-azurerm/issues/13310))
* `azurerm_bot_channel_slack` - support for the `signing_secret` property ([#13454](https://github.com/hashicorp/terraform-provider-azurerm/issues/13454))
* `azurerm_data_factory` - support for `identity` being `SystemAssiged` and `UserAssigned` ([#13473](https://github.com/hashicorp/terraform-provider-azurerm/issues/13473))
* `azurerm_function_app` - support for the `vnet_route_all_enabled` property ([#13310](https://github.com/hashicorp/terraform-provider-azurerm/issues/13310))
* `azurerm_machine_learning_workspace` - support for `public_network_access_enabled`, `public_network_access_enabled`, and `discovery_url` properties ([#13268](https://github.com/hashicorp/terraform-provider-azurerm/issues/13268))
* `azurerm_private_endpoint_connection` - export the `network_interface` attribute from private endpoints ([#13421](https://github.com/hashicorp/terraform-provider-azurerm/issues/13421))
* `azurerm_storage_account_network_rules ` - Deprecate `storage_account_name` and `resource_group_name` in favor of `storage_account_id` ([#13307](https://github.com/hashicorp/terraform-provider-azurerm/issues/13307))
* `azurerm_storage_share_file` - will now recreate and upload deleted/missing files ([#13269](https://github.com/hashicorp/terraform-provider-azurerm/issues/13269))
* `azurerm_synapse_workspace` - the `tenant_id` property is now computed ([#13464](https://github.com/hashicorp/terraform-provider-azurerm/issues/13464))

BUG FIXES:

* Data Source: `azurerm_app_service_certificate` - prevent panics if the API returns a nil `issue_date` or `expiration_date` ([#13401](https://github.com/hashicorp/terraform-provider-azurerm/issues/13401))
* `azurerm_app_service_certificate` - prevent panics if the API returns a nil `issue_date` or `expiration_date` ([#13401](https://github.com/hashicorp/terraform-provider-azurerm/issues/13401))
* `azurerm_app_service_certificate_binding` - reverted a change that introduced a bug in certificate selection for non-managed certificates ([#13455](https://github.com/hashicorp/terraform-provider-azurerm/issues/13455))
* `azurerm_container_group` - allow creation of shared volume between containers in multi container group ([#13374](https://github.com/hashicorp/terraform-provider-azurerm/issues/13374))
* `azurerm_kubernetes_cluster` - changing the `private_cluster_public_fqdn_enabled` no longer created a new resource ([#13413](https://github.com/hashicorp/terraform-provider-azurerm/issues/13413))
* `azurerm_app_configuration_key` - fix nil pointer for removed key ([#13483](https://github.com/hashicorp/terraform-provider-azurerm/issues/13483))

## 2.77.0 (September 17, 2021)

FEATURES:

* **New Data Source:** `azurerm_policy_virtual_machine_configuration_assignment` ([#13311](https://github.com/hashicorp/terraform-provider-azurerm/issues/13311))
* **New Resource:** `azurerm_synapse_integration_runtime_self_hosted` ([#13264](https://github.com/hashicorp/terraform-provider-azurerm/issues/13264))
* **New Resource:** `azurerm_synapse_integration_runtime_azure` ([#13341](https://github.com/hashicorp/terraform-provider-azurerm/issues/13341))
* **New Resource:** `azurerm_synapse_linked_service` ([#13204](https://github.com/hashicorp/terraform-provider-azurerm/issues/13204))
* **New Resource:** `azurerm_synapse_sql_pool_security_alert_policy` ([#13276](https://github.com/hashicorp/terraform-provider-azurerm/issues/13276))
* **New Resource:** `azurerm_synapse_sql_pool_vulnerability_assessment` ([#13276](https://github.com/hashicorp/terraform-provider-azurerm/issues/13276))
* **New Resource:** `azurerm_synapse_workspace_security_alert_policy` ([#13276](https://github.com/hashicorp/terraform-provider-azurerm/issues/13276))
* **New Resource:** `azurerm_synapse_workspace_vulnerability_assessment` ([#13276](https://github.com/hashicorp/terraform-provider-azurerm/issues/13276))

ENHANCEMENTS:

* Data Source: `azurerm_mssql_elasticpool` - export the `sku` block ([#13336](https://github.com/hashicorp/terraform-provider-azurerm/issues/13336))
* `azurerm_api_management` - now supports purging soft deleted instances via the `purge_soft_delete_on_destroy` provider level feature ([#12850](https://github.com/hashicorp/terraform-provider-azurerm/issues/12850))
* `azurerm_data_factory_trigger_schedule` - support for the `activated` property ([#13390](https://github.com/hashicorp/terraform-provider-azurerm/issues/13390))
* `azurerm_logic_app_workflow` - support for the `enabled` and `access_control` properties ([#13265](https://github.com/hashicorp/terraform-provider-azurerm/issues/13265))
* `azurerm_monitor_scheduled_query_rules_alert` - support `auto_mitigation_enabled` property ([#13213](https://github.com/hashicorp/terraform-provider-azurerm/issues/13213))
* `azurerm_machine_learning_inference_cluster` - support for the `identity` block ([#12833](https://github.com/hashicorp/terraform-provider-azurerm/issues/12833))
* `azurerm_machine_learning_compute_cluster` - support for the `ssh_public_access_enabled enhancement` property and the `identity` and `ssh` blocks ([#12833](https://github.com/hashicorp/terraform-provider-azurerm/issues/12833))
* `azurerm_spring_cloud_service` - support for the `connection_string` property ([#13262](https://github.com/hashicorp/terraform-provider-azurerm/issues/13262))

BUG FIXES:

* `azurerm_app_service_certificate_binding` - rework for removal of thumbprint from service ([#13379](https://github.com/hashicorp/terraform-provider-azurerm/issues/13379))
* `azurerm_app_service_managed_certificate`: Fix for empty `issue_date` ([#13357](https://github.com/hashicorp/terraform-provider-azurerm/issues/13357))
* `azurerm_cosmosdb_sql_container`: fix crash when deleting ([#13339](https://github.com/hashicorp/terraform-provider-azurerm/issues/13339))
* `azurerm_frontdoor` - Fix crash when cache is disabled ([#13338](https://github.com/hashicorp/terraform-provider-azurerm/issues/13338))
* `azurerm_function_app` - fix `app_settings` for `WEBSITE_CONTENTSHARE` ([#13349](https://github.com/hashicorp/terraform-provider-azurerm/issues/13349))
* `azurerm_function_app_slot` - fix `app_settings` for `WEBSITE_CONTENTSHARE` ([#13349](https://github.com/hashicorp/terraform-provider-azurerm/issues/13349))
* `azurerm_kubernetes_cluster_node_pool` - `os_sku` is now computed ([#13321](https://github.com/hashicorp/terraform-provider-azurerm/issues/13321))
* `azurerm_linux_virtual_machine_scale_set` - fixed crash when `automatic_os_policy` was nil ([#13335](https://github.com/hashicorp/terraform-provider-azurerm/issues/13335))
* `azurerm_lb` - support for adding or replacing a `frontend_ip_configuration` with an `availability_zone` ([#13305](https://github.com/hashicorp/terraform-provider-azurerm/issues/13305))
* `azurerm_virtual_hub_connection` - fixing race condition in the creation of virtual network resources ([#13294](https://github.com/hashicorp/terraform-provider-azurerm/issues/13294))

## 2.76.0 (September 10, 2021)

NOTES
* Opt-In Beta: Version 2.76 of the Azure Provider introduces an opt-in Beta for some of the new functionality coming in 3.0 - more information can be found [in the 3.0 Notes](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/website/docs/guides/3.0-beta.html.markdown) and [3.0 Upgrade Guide](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/website/docs/guides/3.0-upgrade-guide.html.markdown) ([#12132](https://github.com/hashicorp/terraform-provider-azurerm/issues/12132))

FEATURES:

* **New Data Source:** `azurerm_eventgrid_domain` ([#13033](https://github.com/hashicorp/terraform-provider-azurerm/issues/13033))
* **New Resource:** `azurerm_data_protection_backup_instance_blob_storage` ([#12683](https://github.com/hashicorp/terraform-provider-azurerm/issues/12683))
* **New Resource:** `azurerm_logic_app_integration_account_assembly` ([#13239](https://github.com/hashicorp/terraform-provider-azurerm/issues/13239))
* **New Resource:** `azurerm_logic_app_integration_account_batch_configuration` ([#13215](https://github.com/hashicorp/terraform-provider-azurerm/issues/13215))
* **New Resource:** `azurerm_logic_app_integration_account_agreement` ([#13287](https://github.com/hashicorp/terraform-provider-azurerm/issues/13287))
* **New Resource:** `azurerm_sql_managed_database` ([#12431](https://github.com/hashicorp/terraform-provider-azurerm/issues/12431))

ENHANCEMENTS:

* upgrading `cdn` to API Version `2021-09-01` ([#13282](https://github.com/hashicorp/terraform-provider-azurerm/issues/13282))
* upgrading `cosmos` to API Version `2021-06-15` ([#13188](https://github.com/hashicorp/terraform-provider-azurerm/issues/13188))
* `azurerm_app_service_certificate` - support argument `app_service_plan_id` for usage with ASE ([#13101](https://github.com/hashicorp/terraform-provider-azurerm/issues/13101))
* `azurerm_application_gateway` - mTLS support for Application Gateways ([#13273](https://github.com/hashicorp/terraform-provider-azurerm/issues/13273))
* `azurerm_cosmosdb_account` support for the `local_authentication_disabled` property ([#13237](https://github.com/hashicorp/terraform-provider-azurerm/issues/13237))
* `azurerm_data_factory_integration_runtime_azure` -  support for the `cleanup_enabled` and `subnet_id` properties ([#13222](https://github.com/hashicorp/terraform-provider-azurerm/issues/13222))
* `azurerm_data_factory_trigger_schedule` - support for the `schedule` and `description` properties ([#13243](https://github.com/hashicorp/terraform-provider-azurerm/issues/13243))
* `azurerm_firewall_policy_rule_collection_group` - support for the `description`, `destination_addresses`, `destination_urls`, `terminate_tls`, and `web_categories` properties ([#13190](https://github.com/hashicorp/terraform-provider-azurerm/issues/13190))
* `azurerm_eventgrid_event_subscription` - support for the `delivery_identity` and `dead_letter_identity` blocks ([#12945](https://github.com/hashicorp/terraform-provider-azurerm/issues/12945))
* `azurerm_eventgrid_system_topic_event_subscription` - support for the `delivery_identity` and `dead_letter_identity` blocks ([#12945](https://github.com/hashicorp/terraform-provider-azurerm/issues/12945))
* `azurerm_eventgrid_domain` support for the `identity` block ([#12951](https://github.com/hashicorp/terraform-provider-azurerm/issues/12951))
* `azurerm_eventgrid_topic` support for the `identity` block ([#12951](https://github.com/hashicorp/terraform-provider-azurerm/issues/12951))
* `azurerm_eventgrid_system_topic` support for the `identity` block ([#12951](https://github.com/hashicorp/terraform-provider-azurerm/issues/12951))
* `azurerm_kubernetes_cluster` - support for the `os_sku` property ([#13284](https://github.com/hashicorp/terraform-provider-azurerm/issues/13284))
* `azurerm_synapse_workspace` - support for the `tenant_id` property ([#13290](https://github.com/hashicorp/terraform-provider-azurerm/issues/13290))
* `azurerm_site_recovery_network_mapping`- refactoring to use an ID Formatter/Parser ([#13277](https://github.com/hashicorp/terraform-provider-azurerm/issues/13277))
* `azurerm_stream_analytics_output_blob` - support for the `Parquet` type and the `batch_max_wait_time` and `batch_min_rows` properties ([#13245](https://github.com/hashicorp/terraform-provider-azurerm/issues/13245))
* `azurerm_virtual_network_gateway_resource` - support for multiple vpn authentication types ([#13228](https://github.com/hashicorp/terraform-provider-azurerm/issues/13228))

BUG FIXES:

* Data Source: `azurerm_kubernetes_cluster` - correctly read resource when `local_account_disabled` is `true` ([#13260](https://github.com/hashicorp/terraform-provider-azurerm/issues/13260))
* `azurerm_api_management_subscription` - relax `subscription_id` validation ([#13203](https://github.com/hashicorp/terraform-provider-azurerm/issues/13203))
* `azurerm_app_configuration_key` - fix KV import with no label ([#13253](https://github.com/hashicorp/terraform-provider-azurerm/issues/13253))
* `azurerm_synapse_sql_pool` - properly support UTF-8 characters for the `name` property ([#13289](https://github.com/hashicorp/terraform-provider-azurerm/issues/13289))

## 2.75.0 (September 02, 2021)

FEATURES:

* **New Data Source:** `azurerm_cosmosdb_mongo_database` ([#13123](https://github.com/hashicorp/terraform-provider-azurerm/issues/13123))
* **New Resource:** `azurerm_cognitive_account_customer_managed_key` ([#12901](https://github.com/hashicorp/terraform-provider-azurerm/issues/12901))
* **New Resource:** `azurerm_logic_app_integration_account_partner` ([#13157](https://github.com/hashicorp/terraform-provider-azurerm/issues/13157))
* **New Resource:** `azurerm_logic_app_integration_account_map` ([#13187](https://github.com/hashicorp/terraform-provider-azurerm/issues/13187))
* **New Resource:** `azurerm_app_configuration_key` ([#13118](https://github.com/hashicorp/terraform-provider-azurerm/issues/13118))

ENHANCEMENTS:

* dependencies: upgrading to `v57.0.0` of `github.com/Azure/azure-sdk-for-go` ([#13160](https://github.com/hashicorp/terraform-provider-azurerm/issues/13160))
* upgrading `dataprotection` to API Version `2021-07-01` ([#13161](https://github.com/hashicorp/terraform-provider-azurerm/issues/13161))
* `azurerm_application_insights` - support the `local_authentication_disabled` property ([#13174](https://github.com/hashicorp/terraform-provider-azurerm/issues/13174))
* `azurerm_data_factory_linked_service_azure_blob_storage` - support for the `key_vault_sas_token` property ([#12880](https://github.com/hashicorp/terraform-provider-azurerm/issues/12880))
* `azurerm_data_factory_linked_service_azure_function` support for the `key_vault_key` block ([#13159](https://github.com/hashicorp/terraform-provider-azurerm/issues/13159))
* `azurerm_data_protection_backup_instance_postgresql` - support the `database_credential_key_vault_secret_id` property ([#13183](https://github.com/hashicorp/terraform-provider-azurerm/issues/13183))
* `azurerm_hdinsight_hadoop_cluster` - support for the `security_profile` block ([#12866](https://github.com/hashicorp/terraform-provider-azurerm/issues/12866))
* `azurerm_hdinsight_hbase_cluster` - support for the `security_profile` block ([#12866](https://github.com/hashicorp/terraform-provider-azurerm/issues/12866))
* `azurerm_hdinsight_interactive_query_cluster` - support for the `security_profile` block ([#12866](https://github.com/hashicorp/terraform-provider-azurerm/issues/12866))
* `azurerm_hdinsight_kafka_cluster` - support for the `security_profile` block ([#12866](https://github.com/hashicorp/terraform-provider-azurerm/issues/12866))
* `azurerm_hdinsight_spark_cluster` - support for the `security_profile` block ([#12866](https://github.com/hashicorp/terraform-provider-azurerm/issues/12866))
* `azurerm_mssql_server`- refactoring to use an ID Formatter/Parser ([#13151](https://github.com/hashicorp/terraform-provider-azurerm/issues/13151))
* `azurerm_policy_virtual_machine_configuration_assignment` - support for the `assignment_type`, `content_uri`, and `content_hash` properties ([#13176](https://github.com/hashicorp/terraform-provider-azurerm/issues/13176))
* `azurerm_storage_account` - handle nil values for AllowBlobPublicAccess ([#12689](https://github.com/hashicorp/terraform-provider-azurerm/issues/12689))
* `azurerm_synapse_spark_pool` - add support spark for `3.1` ([#13181](https://github.com/hashicorp/terraform-provider-azurerm/issues/13181))

## 2.74.0 (August 27, 2021)

FEATURES:

* **New Resource:** `azurerm_logic_app_integration_account_schema` ([#13100](https://github.com/hashicorp/terraform-provider-azurerm/issues/13100))
* **New Resource:** `azurerm_relay_namespace_authorization_rule` ([#13116](https://github.com/hashicorp/terraform-provider-azurerm/issues/13116))
* **New Resource:** `azurerm_relay_hybrid_connection_authorization_rule` ([#13116](https://github.com/hashicorp/terraform-provider-azurerm/issues/13116))

ENHANCEMENTS:

* dependencies: upgrading `monitor` to API Version `2021-07-01-preview` ([#13121](https://github.com/hashicorp/terraform-provider-azurerm/issues/13121))
* dependencies: upgrading `devtestlabs` to API Version `2018-09-15` ([#13074](https://github.com/hashicorp/terraform-provider-azurerm/issues/13074))
* Data Source: `azurerm_servicebus_namespace_authorization_rule` - support for the `primary_connection_string_alias` and `secondary_connection_string_alias` properties ([#12997](https://github.com/hashicorp/terraform-provider-azurerm/issues/12997))
* Data Source: `azurerm_servicebus_queue_authorization_rule` - support for the `primary_connection_string_alias` and `secondary_connection_string_alias` properties ([#12997](https://github.com/hashicorp/terraform-provider-azurerm/issues/12997))
* Data Source: `azurerm_network_service_tags` - new properties `ipv4_cidrs` and `ipv6_cidrs` ([#13058](https://github.com/hashicorp/terraform-provider-azurerm/issues/13058))
* `azurerm_api_management` - now exports certificate `expiry`, `thumbprint` and `subject` attributes ([#12262](https://github.com/hashicorp/terraform-provider-azurerm/issues/12262))
* `azurerm_app_configuration` - support for user assigned identities ([#13080](https://github.com/hashicorp/terraform-provider-azurerm/issues/13080))
* `azurerm_app_service` - add support for `vnet_route_all_enabled` property ([#13073](https://github.com/hashicorp/terraform-provider-azurerm/issues/13073))
* `azurerm_app_service_plan` - support for the `zone_redundant` property  ([#13145](https://github.com/hashicorp/terraform-provider-azurerm/issues/13145))
* `azurerm_data_factory_dataset_binary` -  support for `dynamic_path_enabled` and `dynamic_path_enabled`  properties ([#13117](https://github.com/hashicorp/terraform-provider-azurerm/issues/13117))
* `azurerm_data_factory_dataset_delimited_text` -  support for `dynamic_path_enabled` and `dynamic_path_enabled`  properties ([#13117](https://github.com/hashicorp/terraform-provider-azurerm/issues/13117))
* `azurerm_data_factory_dataset_json` -  support for `dynamic_path_enabled` and `dynamic_path_enabled`  properties ([#13117](https://github.com/hashicorp/terraform-provider-azurerm/issues/13117))
* `azurerm_data_factory_dataset_parquet` -  support for `dynamic_path_enabled` and `dynamic_path_enabled`  properties ([#13117](https://github.com/hashicorp/terraform-provider-azurerm/issues/13117))
* `azurerm_firewall_policy` - support for the `intrusion_detection`, `identity` and `tls_certificate` blocks ([#12769](https://github.com/hashicorp/terraform-provider-azurerm/issues/12769))
* `azurerm_kubernetes_cluster` - support for the `pod_subnet_id` property ([#12313](https://github.com/hashicorp/terraform-provider-azurerm/issues/12313))
* `azurerm_kubernetes_cluster_node_pool` - support for the `pod_subnet_id` property ([#12313](https://github.com/hashicorp/terraform-provider-azurerm/issues/12313))
* `azurerm_monitor_autoscale_setting` - support for the field `divide_by_instance_count` within the `metric_trigger` block ([#13121](https://github.com/hashicorp/terraform-provider-azurerm/issues/13121))
* `azurerm_redis_enterprise_cluster` - the `tags` property can now be updated ([#13084](https://github.com/hashicorp/terraform-provider-azurerm/issues/13084))
* `azurerm_storage_account` - add support for `shared_key_access_enabled` property ([#13014](https://github.com/hashicorp/terraform-provider-azurerm/issues/13014))
* `azurerm_servicebus_namespace_authorization_rule` - support for the `primary_connection_string_alias` and `secondary_connection_string_alias` properties ([#12997](https://github.com/hashicorp/terraform-provider-azurerm/issues/12997))
* `azurerm_servicebus_topic_authorization_rule` - support for the `primary_connection_string_alias` and `secondary_connection_string_alias` properties ([#12997](https://github.com/hashicorp/terraform-provider-azurerm/issues/12997))
* `azurerm_dev_test_global_vm_shutdown_schedule` - support for the `mail` property ([#13074](https://github.com/hashicorp/terraform-provider-azurerm/issues/13074))

BUG FIXES:

* `azurerm_data_factory_dataset_delimited_text` - support empty values for the `column_delimiter`, `row_delimiter`, `quote_character`, `escape_character`, and `encoding` propeties ([#13149](https://github.com/hashicorp/terraform-provider-azurerm/issues/13149))
* `azurerm_cosmosdb_cassandra_table` - correctly update `throughput` ([#13102](https://github.com/hashicorp/terraform-provider-azurerm/issues/13102))
* `azurerm_private_dns_a_record` - fix regression in `name` validation and add max recordset limit validation ([#13093](https://github.com/hashicorp/terraform-provider-azurerm/issues/13093))
* `azurerm_postgresql_flexible_server_database` the `charset` and `collation` properties are now optional ([#13110](https://github.com/hashicorp/terraform-provider-azurerm/issues/13110))
* `azurerm_spring_cloud_app` - Fix crash when identity is not present ([#13125](https://github.com/hashicorp/terraform-provider-azurerm/issues/13125))

## 2.73.0 (August 20, 2021)

FEATURES:

* **New Data Source:** `azurerm_vpn_gateway` ([#12844](https://github.com/hashicorp/terraform-provider-azurerm/issues/12844))
* **New Data Source:** `azurerm_data_protection_backup_vault` ([#13062](https://github.com/hashicorp/terraform-provider-azurerm/issues/13062))
* **New Resource:** `azurerm_api_management_notification_recipient_email` ([#12849](https://github.com/hashicorp/terraform-provider-azurerm/issues/12849))
* **New Resource:** `azurerm_logic_app_integration_account_session` ([#12982](https://github.com/hashicorp/terraform-provider-azurerm/issues/12982))
* **New Resource:** `azurerm_machine_learning_synapse_spark` ([#13022](https://github.com/hashicorp/terraform-provider-azurerm/issues/13022))
* **New Resource:** `azurerm_machine_learning_compute_instance` ([#12834](https://github.com/hashicorp/terraform-provider-azurerm/issues/12834))
* **New Resource:** `azurerm_vpn_gateway` ([#13003](https://github.com/hashicorp/terraform-provider-azurerm/issues/13003))

ENHANCEMENTS:

* Dependencies: upgrade `github.com/Azure/azure-sdk-for-go` to `v56.2.0` ([#12969](https://github.com/hashicorp/terraform-provider-azurerm/issues/12969))
* Dependencies: updating `frontdoor` to use API version `2020-05-01` ([#12831](https://github.com/hashicorp/terraform-provider-azurerm/issues/12831))
* Dependencies: updating `web` to use API version `2021-02-01` ([#12970](https://github.com/hashicorp/terraform-provider-azurerm/issues/12970))
* Dependencies: updating `kusto` to use API version `2021-01-01` ([#12967](https://github.com/hashicorp/terraform-provider-azurerm/issues/12967))
* Dependencies: updating `machinelearning` to use API version `2021-07-01` ([#12833](https://github.com/hashicorp/terraform-provider-azurerm/issues/12833))
* Dependencies: updating `network` to use API version `2021-02-01` ([#13002](https://github.com/hashicorp/terraform-provider-azurerm/issues/13002))
* appconfiguration: updating to use the latest embedded SDK ([#12950](https://github.com/hashicorp/terraform-provider-azurerm/issues/12950))
* eventhub: updating to use the latest embedded SDK ([#12946](https://github.com/hashicorp/terraform-provider-azurerm/issues/12946))
* Data Source: `azurerm_iothub` - support for the property `hostname` ([#13001](https://github.com/hashicorp/terraform-provider-azurerm/issues/13001))
* Data Source: `azurerm_application_security_group` - refactoring to use an ID Formatter/Parser ([#13028](https://github.com/hashicorp/terraform-provider-azurerm/issues/13028))
* `azurerm_active_directory_domain_service` - export the `resource_id` attribute ([#13011](https://github.com/hashicorp/terraform-provider-azurerm/issues/13011))
* `azurerm_app_service_environment_v3` - updated for GA changes, including support for `internal_load_balancing_mode`, `zone_redundant`, `dedicated_host_count`, and several new exported properties ([#12932](https://github.com/hashicorp/terraform-provider-azurerm/issues/12932))
* `azurerm_application_security_group` - refactoring to use an ID Formatter/Parser ([#13028](https://github.com/hashicorp/terraform-provider-azurerm/issues/13028))
* `azurerm_data_lake_store` - support for the `identity` block ([#13050](https://github.com/hashicorp/terraform-provider-azurerm/issues/13050))
* `azurerm_kubernetes_cluster` - support for the `ultra_ssd_enabled` and `private_cluster_public_fqdn_enabled` properties ([#12780](https://github.com/hashicorp/terraform-provider-azurerm/issues/12780))
* `azurerm_kubernetes_cluster_node_pool` - supportfor the `ultra_ssd_enabled` property ([#12780](https://github.com/hashicorp/terraform-provider-azurerm/issues/12780))
* `azurerm_logic_app_trigger_http_request` - support for the `callback_url` attribute ([#13057](https://github.com/hashicorp/terraform-provider-azurerm/issues/13057))
* `azurerm_netapp_volume` - support for the `snapshot_directory_visible` property ([#12961](https://github.com/hashicorp/terraform-provider-azurerm/issues/12961))
* `azurerm_sql_server` - support for configuring `threat_detection_policy` ([#13048](https://github.com/hashicorp/terraform-provider-azurerm/issues/13048))
* `azurerm_stream_analytics_output_eventhub` - support for the `property_columns` property ([#12947](https://github.com/hashicorp/terraform-provider-azurerm/issues/12947))

BUG FIXES:

* `azurerm_frontdoor` - expose support for `cache_duration` and `cache_query_parameters` fields ([#12831](https://github.com/hashicorp/terraform-provider-azurerm/issues/12831))
* `azurerm_network_watcher_flow_log` - correctly truncate name by ensuring it doesn't end in a `-` ([#12984](https://github.com/hashicorp/terraform-provider-azurerm/issues/12984))
* `azurerm_databricks_workspace` - corrent logic for the `public_network_access_enabled` property ([#13034](https://github.com/hashicorp/terraform-provider-azurerm/issues/13034))
* `azurerm_databricks_workspace` - fix potential crash in Read ([#13025](https://github.com/hashicorp/terraform-provider-azurerm/issues/13025))
* `azurerm_private_dns_zone_id` - correctly handle inconsistant case ([#13000](https://github.com/hashicorp/terraform-provider-azurerm/issues/13000))
* `azurerm_private_dns_a_record_resource` - currently validate the name property by allowing `@`s ([#13042](https://github.com/hashicorp/terraform-provider-azurerm/issues/13042))
* `azurerm_eventhub_namespace` - support upto `40` for the `maximum_throughput_units` property ([#13065](https://github.com/hashicorp/terraform-provider-azurerm/issues/13065))
* `azurerm_kubernetes_cluster` - fix crash in update when previously configured AAD Profile is now `nil` ([#13043](https://github.com/hashicorp/terraform-provider-azurerm/issues/13043))
* `azurerm_redis_enterprise_cluster` - changing the tags property no longer creates a new resource ([#12956](https://github.com/hashicorp/terraform-provider-azurerm/issues/12956))
* `azurerm_storage_account` - allow 0 for the `cors.max_age_in_seconds` property ([#13010](https://github.com/hashicorp/terraform-provider-azurerm/issues/13010))
* `azurerm_servicebus_topic` - correctyl validate the `name` property ([#13026](https://github.com/hashicorp/terraform-provider-azurerm/issues/13026))
* `azurerm_virtual_hub_connection` - will not correctly lock it's cirtual network during updates ([#12999](https://github.com/hashicorp/terraform-provider-azurerm/issues/12999))
* `azurerm_linux_virtual_machine_scale_set` - fix potential crash in updates to the `rolling_upgrade_policy` block ([#13029](https://github.com/hashicorp/terraform-provider-azurerm/issues/13029))



## 2.72.0 (August 12, 2021)

UPGRADE NOTES

* This version of the Azure Provider introduces the `prevent_deletion_if_contains_resources` feature flag (which is disabled by default) which (when enabled) means that Terraform will check for Resources nested within the Resource Group during the deletion of the Resource Group and require that these Resources are deleted first. This avoids the unintentional deletion of unmanaged Resources within a Resource Group - and is defaulted off in 2.x versions of the Azure Provider but **will be enabled by default in version 3.0 of the Azure Provider**, see [the `features` block documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#features) for more information. ([#12657](https://github.com/hashicorp/terraform-provider-azurerm/issues/12657))


FEATURES:

* **New Resource:** `azurerm_video_analyzer` ([#12665](https://github.com/hashicorp/terraform-provider-azurerm/issues/12665))
* **New Resource:** `azurerm_video_analyzer_edge_module` ([#12911](https://github.com/hashicorp/terraform-provider-azurerm/issues/12911))

ENHANCEMENTS:

* `azurerm_api_management_named_value` - support for system managed identities ([#12938](https://github.com/hashicorp/terraform-provider-azurerm/issues/12938))
* `azurerm_application_insights_smart_detection_rule` - support all currenly availible rules in the SDK ([#12857](https://github.com/hashicorp/terraform-provider-azurerm/issues/12857))
* `azurerm_function_app` - add support for `dotnet_framework_version` in ([#12883](https://github.com/hashicorp/terraform-provider-azurerm/issues/12883))
* `azurerm_resource_group` - conditionally (based on the `prevent_deletion_if_contains_resources` features flag - see the 'Upgrade Notes' section) checking for nested Resources during deletion of the Resource Group and raising an error if Resources are found ([#12657](https://github.com/hashicorp/terraform-provider-azurerm/issues/12657))

BUG FIXES:

* Data Source: `azurerm_key_vault_certificate_data` - updating the PEM Header when using a RSA Private Key so this validates with OpenSSL ([#12896](https://github.com/hashicorp/terraform-provider-azurerm/issues/12896))
* `azurerm_active_directory_domain_service` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_app_service_environment` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_cdn_profile` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_container_registry_scope_map` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_container_registry_token` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_container_registry_webhook` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_container_registry` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_data_factory_dataset_delimited_text` - correctly send optional optional values to the API ([#12921](https://github.com/hashicorp/terraform-provider-azurerm/issues/12921))
* `azurerm_data_lake_analytics_account` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_data_lake_store` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_data_protection_backup_instance_disk` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_database_migration_service` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_dns_zone` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_eventgrid_domain_topic` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_eventgrid_domain` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_eventgrid_event_subscription` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_eventgrid_system_topic_event_subscription` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_eventgrid_system_topic` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_eventgrid_topic` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_express_route_circuit_authorization` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_express_route_circuit_peering` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_express_route_gateway` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_express_route_port` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_frontdoor_firewall_policy` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_hpc_cache_blob_nfs_target` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_iothub` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_key_vault_managed_hardware_security_module` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_kubernetes_cluster` - prevent nil panic when rbac config is empty ([#12881](https://github.com/hashicorp/terraform-provider-azurerm/issues/12881))
* `azurerm_iot_dps` - fixing a crash during creation ([#12919](https://github.com/hashicorp/terraform-provider-azurerm/issues/12919))
* `azurerm_local_network_gateway` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_logic_app_trigger_recurrence` - update time zone strings to match API behaviour, and use the timezone even when `start_time` is not specified ([#12453](https://github.com/hashicorp/terraform-provider-azurerm/issues/12453))
* `azurerm_mariadb_database` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_mariadb_server` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_mariadb_virtual_network_rule` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_mssql_database` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_mssql_virtual_network_rule` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_mysql_server` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_nat_gateway` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_network_packet_capture` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_packet_capture` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_postgresql_configuration` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_postgresql_firewall_rule` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_postgresql_server` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_postgresql_virtual_network_rule` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_private_dns_zone_virtual_network_link` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_private_endpoint` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_private_link_service` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_shared_image_gallery` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_sql_virtual_network_rule` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_virtual_machine_scale_set_extension` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_virtual_wan` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_vpn_gateway_connection` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))
* `azurerm_web_application_firewall_policy` - removing an unnecessary check during deletion ([#12879](https://github.com/hashicorp/terraform-provider-azurerm/issues/12879))

## 2.71.0 (August 06, 2021)

FEATURES:

* **New Data Source:** `azurerm_databricks_workspace_private_endpoint_connection` ([#12543](https://github.com/hashicorp/terraform-provider-azurerm/issues/12543))
* **New Resource:** `azurerm_api_management_tag` ([#12535](https://github.com/hashicorp/terraform-provider-azurerm/issues/12535))
* **New Resource:** `azurerm_bot_channel_line` ([#12746](https://github.com/hashicorp/terraform-provider-azurerm/issues/12746))
* **New Resource:** `azurerm_cdn_endpoint_custom_domain` ([#12496](https://github.com/hashicorp/terraform-provider-azurerm/issues/12496))
* **New Resource:** `azurerm_data_factory_data_flow` ([#12588](https://github.com/hashicorp/terraform-provider-azurerm/issues/12588))
* **New Resource:** `azurerm_postgresql_flexible_server_database` ([#12550](https://github.com/hashicorp/terraform-provider-azurerm/issues/12550))

ENHANCEMENTS:

* dependencies: upgrading to `v56.0.0` of `github.com/Azure/azure-sdk-for-go` ([#12781](https://github.com/hashicorp/terraform-provider-azurerm/issues/12781))
* dependencies: updating `appinsights` to use API Version `2020-02-02` ([#12818](https://github.com/hashicorp/terraform-provider-azurerm/issues/12818))
* dependencies: updating `containerservice` to use API Version `2021-05-1` ([#12747](https://github.com/hashicorp/terraform-provider-azurerm/issues/12747))
* dependencies: updating `machinelearning` to use API Version `2021-04-01` ([#12804](https://github.com/hashicorp/terraform-provider-azurerm/issues/12804))
* dependencies: updating `databricks` to use API Version `2021-04-01-preview` ([#12543](https://github.com/hashicorp/terraform-provider-azurerm/issues/12543))
* PowerBI: refactoring to use an Embedded SDK ([#12787](https://github.com/hashicorp/terraform-provider-azurerm/issues/12787))
* SignalR: refactoring to use an Embedded SDK ([#12785](https://github.com/hashicorp/terraform-provider-azurerm/issues/12785))
* `azurerm_api_management_api_diagnostic` - support for the `operation_name_format` property ([#12782](https://github.com/hashicorp/terraform-provider-azurerm/issues/12782))
* `azurerm_app_service` - support for the acr_use_managed_identity_credentials and acr_user_managed_identity_client_id properties ([#12745](https://github.com/hashicorp/terraform-provider-azurerm/issues/12745))
* `azurerm_app_service` - support `v6.0` for the `dotnet_framework_version` property ([#12788](https://github.com/hashicorp/terraform-provider-azurerm/issues/12788))
* `azurerm_application_insights` - support for the `workspace_id` property ([#12818](https://github.com/hashicorp/terraform-provider-azurerm/issues/12818))
* `azurerm_databricks_workspace` - support for private link endpoint ([#12543](https://github.com/hashicorp/terraform-provider-azurerm/issues/12543))
* `azurerm_databricks_workspace` - add support for `Customer Managed Keys for Managed Services` ([#12799](https://github.com/hashicorp/terraform-provider-azurerm/issues/12799))
* `azurerm_data_factory_linked_service_data_lake_storage_gen2` - don't send a secure connection string when using a managed identity ([#12359](https://github.com/hashicorp/terraform-provider-azurerm/issues/12359))
* `azurerm_function_app` - support for the `elastic_instance_minimum`, `app_scale_limit`, and `runtime_scale_monitoring_enabled` properties ([#12741](https://github.com/hashicorp/terraform-provider-azurerm/issues/12741))
* `azurerm_kubernetes_cluster` - support for the `local_account_disabled` property ([#12386](https://github.com/hashicorp/terraform-provider-azurerm/issues/12386))
* `azurerm_kubernetes_cluster` - support for the `maintenance_window` block ([#12762](https://github.com/hashicorp/terraform-provider-azurerm/issues/12762))
* `azurerm_kubernetes_cluster` - the field `automatic_channel_upgrade` can now be set to `node-image` ([#12667](https://github.com/hashicorp/terraform-provider-azurerm/issues/12667))
* `azurerm_logic_app_workflow` - support for the `workflow_parameters` ([#12314](https://github.com/hashicorp/terraform-provider-azurerm/issues/12314))
* `azurerm_mssql_database` - support for the `Free` and `FSV2` SKU's ([#12835](https://github.com/hashicorp/terraform-provider-azurerm/issues/12835))
* `azurerm_network_security_group` - the `protocol` property now supports `Ah` and `Esp` values ([#12865](https://github.com/hashicorp/terraform-provider-azurerm/issues/12865))
* `azurerm_public_ip_resource` - support for sku_tier property ([#12775](https://github.com/hashicorp/terraform-provider-azurerm/issues/12775))
* `azurerm_redis_cache` - support for the `replicas_per_primary`, `redis_version`, and `tenant_settings` properties and blocks ([#12820](https://github.com/hashicorp/terraform-provider-azurerm/issues/12820))
* `azurerm_redis_enterprise_cluster` - this can now be provisioned in `Canada Central` ([#12842](https://github.com/hashicorp/terraform-provider-azurerm/issues/12842))
* `azurerm_static_site` - support `Standard` SKU ([#12510](https://github.com/hashicorp/terraform-provider-azurerm/issues/12510))

BUG FIXES:

* Data Source `azurerm_ssh_public_key` - normalising the SSH Public Key ([#12800](https://github.com/hashicorp/terraform-provider-azurerm/issues/12800))
* `azurerm_api_management_api_subscription` - fixing the default scope to be `/apis` rather than `all_apis` as required by the latest API ([#12829](https://github.com/hashicorp/terraform-provider-azurerm/issues/12829))
* `azurerm_app_service_active_slot` - fix 404 not found on read for slot ([#12792](https://github.com/hashicorp/terraform-provider-azurerm/issues/12792))
* `azurerm_linux_virtual_machine_scale_set` - fix crash in checking for latest image ([#12808](https://github.com/hashicorp/terraform-provider-azurerm/issues/12808))
* `azurerm_kubernetes_cluster` - corrently valudate the `net_ipv4_ip_local_port_range_max` property ([#12859](https://github.com/hashicorp/terraform-provider-azurerm/issues/12859))
* `azurerm_local_network_gateway` - fixing a crash where the `LocalNetworkAddressSpace` block was nil ([#12822](https://github.com/hashicorp/terraform-provider-azurerm/issues/12822))
* `azurerm_notification_hub_authorization_rule` - switching to use an ID Formatter ([#12845](https://github.com/hashicorp/terraform-provider-azurerm/issues/12845))
* `azurerm_notification_hub` - switching to use an ID Formatter ([#12845](https://github.com/hashicorp/terraform-provider-azurerm/issues/12845))
* `azurerm_notification_hub_namespace` - switching to use an ID Formatter ([#12845](https://github.com/hashicorp/terraform-provider-azurerm/issues/12845))
* `azurerm_postgresql_database` - fixing a crash in the Azure SDK ([#12823](https://github.com/hashicorp/terraform-provider-azurerm/issues/12823))
* `azurerm_private_dns_zone` - fixing a crash during deletion ([#12824](https://github.com/hashicorp/terraform-provider-azurerm/issues/12824))
* `azurerm_resource_group_template_deployment` - fixing deletion of nested items when using non-top level items ([#12421](https://github.com/hashicorp/terraform-provider-azurerm/issues/12421))
* `azurerm_subscription_template_deployment` - fixing deletion of nested items when using non-top level items ([#12421](https://github.com/hashicorp/terraform-provider-azurerm/issues/12421))
* `azurerm_virtual_machine_extension` - changing the `publisher` property now creates a new resource ([#12790](https://github.com/hashicorp/terraform-provider-azurerm/issues/12790))

## 2.70.0 (July 30, 2021)

FEATURES:

* **New Data Source** `azurerm_storage_share` ([#12693](https://github.com/hashicorp/terraform-provider-azurerm/issues/12693))
* **New Resource** `azurerm_bot_channel_alexa` ([#12682](https://github.com/hashicorp/terraform-provider-azurerm/issues/12682))
* **New Resource** `azurerm_bot_channel_direct_line_speech` ([#12735](https://github.com/hashicorp/terraform-provider-azurerm/issues/12735))
* **New Resource** `azurerm_bot_channel_facebook` ([#12709](https://github.com/hashicorp/terraform-provider-azurerm/issues/12709))
* **New Resource** `azurerm_bot_channel_sms` ([#12713](https://github.com/hashicorp/terraform-provider-azurerm/issues/12713))
* **New Resource** `azurerm_data_factory_trigger_custom_event` ([#12448](https://github.com/hashicorp/terraform-provider-azurerm/issues/12448))
* **New Resource** `azurerm_data_factory_trigger_tumbling_window` ([#12437](https://github.com/hashicorp/terraform-provider-azurerm/issues/12437))
* **New Resource** `azurerm_data_protection_backup_instance_disk` ([#12617](https://github.com/hashicorp/terraform-provider-azurerm/issues/12617))

ENHANCEMENTS:

* dependencies: Upgrade `web` (App Service) API to `2021-01-15` ([#12635](https://github.com/hashicorp/terraform-provider-azurerm/issues/12635))
* analysisservices: refactoring to use an Embedded SDK ([#12771](https://github.com/hashicorp/terraform-provider-azurerm/issues/12771))
* maps: refactoring to use an Embedded SDK ([#12716](https://github.com/hashicorp/terraform-provider-azurerm/issues/12716))
* msi: refactoring to use an Embedded SDK ([#12715](https://github.com/hashicorp/terraform-provider-azurerm/issues/12715))
* relay: refactoring to use an Embedded SDK ([#12772](https://github.com/hashicorp/terraform-provider-azurerm/issues/12772))
* vmware: refactoring to use an Embedded SDK ([#12751](https://github.com/hashicorp/terraform-provider-azurerm/issues/12751))
* Data Source: `azurerm_storage_account_sas` - support for the property `ip_addresses` ([#12705](https://github.com/hashicorp/terraform-provider-azurerm/issues/12705))
* `azurerm_api_management_diagnostic` - support for the property `operation_name_format` ([#12736](https://github.com/hashicorp/terraform-provider-azurerm/issues/12736))
* `azurerm_automation_certificate` - the `exportable` property can now be set ([#12738](https://github.com/hashicorp/terraform-provider-azurerm/issues/12738))
* `azurerm_data_factory_dataset_binary` - the blob `path` and `filename` propeties are now optional ([#12676](https://github.com/hashicorp/terraform-provider-azurerm/issues/12676))
* `azurerm_data_factory_trigger_blob_event` - support for the `activation` property ([#12644](https://github.com/hashicorp/terraform-provider-azurerm/issues/12644))
* `azurerm_data_factory_pipeline` - support for the `concurrency` and `moniter_metrics_after_duration` properties ([#12685](https://github.com/hashicorp/terraform-provider-azurerm/issues/12685))
* `azurerm_hdinsight_interactive_query_cluster` - support for the `encryption_in_transit_enabled` property ([#12767](https://github.com/hashicorp/terraform-provider-azurerm/issues/12767))
* `azurerm_hdinsight_spark_cluster` - support for the `encryption_in_transit_enabled` property ([#12767](https://github.com/hashicorp/terraform-provider-azurerm/issues/12767))
* `azurerm_firewall_polcy` - support for property `private_ip_ranges` ([#12696](https://github.com/hashicorp/terraform-provider-azurerm/issues/12696))

BUG FIXES:

* `azurerm_cdn_endpoint` - fixing a crash when the future is nil ([#12743](https://github.com/hashicorp/terraform-provider-azurerm/issues/12743))
* `azurerm_private_endpoint` - working around a casing issue in `private_connection_resource_id` for MariaDB, MySQL and PostgreSQL resources ([#12761](https://github.com/hashicorp/terraform-provider-azurerm/issues/12761))

---

For information on changes between the v2.69.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).
