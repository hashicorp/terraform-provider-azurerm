## 2.20.0 (Unreleased)

UPGRADE NOTES

* **Enhanced Validation for Locations** - the Azure Provider now validates that the value for the `location` argument is a supported Azure Region within the Azure Environment being used (from the Azure Metadata Service) - which allows us to catch configuration errors for this field at `terraform plan` time, rather than during a `terraform apply`. This functionality is now enabled by default, and can be opted-out of by setting the Environment Variable `ARM_PROVIDER_ENHANCED_VALIDATION` to `false`

FEATURES: 

* **New Resource:** `azurerm_kusto_cluster_customer_managed_key` [GH-7520]

ENHANCEMENTS:

* * dependencies: updating to v44.1.0 of `github.com/Azure/azure-sdk-for-go` [GH-7774]

BUG FIXES

* Data Source: `azurerm_private_dns_zone` - fix a crash when the zone does not exist [GH-7783]
* `azurerm_application_gateway` - Fix crash with `gateway_ip_configuration` [GH-7789]

## 2.19.0 (July 16, 2020)

UPGRADE NOTES:

* HDInsight 3.6 will be retired (in Azure Public) on 2020-12-30 - HDInsight 4.0 does not support ML Services, RServer or Storm Clusters - as such the `azurerm_hdinsight_ml_services_cluster`, `azurerm_hdinsight_rserver_cluster` and `azurerm_hdinsight_storm_cluster` resources are deprecated and will be removed in the next major version of the Azure Provider. ([#7706](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7706))
* provider: no longer auto register the Microsoft.StorageCache RP ([#7768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7768))

FEATURES:

* **New Data source:** `azurerm_route_filter` ([#6341](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6341))
* **New Resource:** `azurerm_route_filter` ([#6341](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6341))

ENHANCEMENTS:

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
* `azurerm_storage_container` - container creation will retry if a container of the same name has not completed it's delete operation ([#7179](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7179))
* `azurerm_storage_share` - share creation will retry if a share of the same name has not completed it's previous delete operation ([#7179](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7179))
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

ENHANCEMENTS:

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

ENHANCEMENTS:

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

ENHANCEMENTS:

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

ENHANCEMENTS:

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

---

For information on changes between the v2.10.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.44.0 changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
