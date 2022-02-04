## 2.95.0 (February 04, 2022)

FEATURES: 

* **New Data Source:** `azurerm_container_group` ([#14946](https://github.com/hashicorp/terraform-provider-azurerm/issues/14946))
* **New Data Source:** `azurerm_logic_app_standard` ([#15199](https://github.com/hashicorp/terraform-provider-azurerm/issues/15199))
* **New Beta Resource:** - `azurerm_disk_pool_iscsi_target` ([#14975](https://github.com/hashicorp/terraform-provider-azurerm/issues/14975))
* **New Beta Resource:** - `azurerm_linux_function_app_slot` ([#14940](https://github.com/hashicorp/terraform-provider-azurerm/issues/14940))
* **New Beta Resource:** - `azurerm_windows_function_app_slot` ([#14940](https://github.com/hashicorp/terraform-provider-azurerm/issues/14940))
* **New Beta Resource:** - `azurerm_windows_web_app_slot` ([#14613](https://github.com/hashicorp/terraform-provider-azurerm/issues/14613))
* **New Beat Resource:** - `azurerm_traffic_manager_azure_endpoint` ([#15178](https://github.com/hashicorp/terraform-provider-azurerm/issues/15178))
* **New Beat Resource:** - `azurerm_traffic_manager_external_endpoint` ([#15178](https://github.com/hashicorp/terraform-provider-azurerm/issues/15178))
* **New Beat Resource:** - `azurerm_traffic_manager_nested_endpoint` ([#15178](https://github.com/hashicorp/terraform-provider-azurerm/issues/15178))

ENHANCEMENTS:

* dependencies: upgrading to `v0.22.0` of `github.com/hashicorp/go-azure-helpers` ([#15207](https://github.com/hashicorp/terraform-provider-azurerm/issues/15207))
* dependencies: updating `backup` to API Version `2021-07-01` ([#14980](https://github.com/hashicorp/terraform-provider-azurerm/issues/14980))
* `azurerm_storage_account` - the `identity` block is no longer computed ([#15207](https://github.com/hashicorp/terraform-provider-azurerm/issues/15207))
* `azurerm_linux_virtual_machine` - support for the `dedicated_host_group_id` property ([#14936](https://github.com/hashicorp/terraform-provider-azurerm/issues/14936))
* `azurerm_recovery_services_vault` - support Zone Redundant storage ([#14980](https://github.com/hashicorp/terraform-provider-azurerm/issues/14980))
* `azurerm_web_pubsub_hub` - the `managed_identity_id` property within the `auth` block now accepts UUIDs ([#15183](https://github.com/hashicorp/terraform-provider-azurerm/issues/15183))
* `azurerm_windows_virtual_machine` - support for the `dedicated_host_group_id` property ([#14936](https://github.com/hashicorp/terraform-provider-azurerm/issues/14936))

BUG FIXES:

* `azurerm_container_group` - fixing parallel provisioning failures with the same `network_profile_id` ([#15098](https://github.com/hashicorp/terraform-provider-azurerm/issues/15098))
* `azurerm_frontdoor` - fixing the validation for `resource_group_name` ([#15174](https://github.com/hashicorp/terraform-provider-azurerm/issues/15174))
* `azurerm_kubernetes_cluster` - prevent panic when updating `sku_tier` ([#15229](https://github.com/hashicorp/terraform-provider-azurerm/issues/15229))
* `azurerm_hdinsight_interactive_query_cluster` - support for the `storage_resource_id` property to fix missing storage account errors ([#15039](https://github.com/hashicorp/terraform-provider-azurerm/issues/15039))
* `azurerm_hdinsight_hadoop_cluster` - support for the `storage_resource_id` property to fix missing storage account errors ([#15039](https://github.com/hashicorp/terraform-provider-azurerm/issues/15039))
* `azurerm_hdinsight_spark_cluster` - support for the `storage_resource_id` property to fix missing storage account errors ([#15039](https://github.com/hashicorp/terraform-provider-azurerm/issues/15039))
* `azurerm_hdinsight_hbase_cluster` - support for the `storage_resource_id` property to fix missing storage account errors ([#15039](https://github.com/hashicorp/terraform-provider-azurerm/issues/15039))
* `azurerm_log_analytics_datasource_windows_event` - adding a state migration to fix `ID was missing the dataSources element` ([#15194](https://github.com/hashicorp/terraform-provider-azurerm/issues/15194))
* `azurerm_policy_definition` - fix the deprecation of `management_group_name` in favour of `management_group_id` ([#15209](https://github.com/hashicorp/terraform-provider-azurerm/issues/15209))
* `azurerm_policy_set_definition` - fix the deprecation of `management_group_name` in favour of `management_group_id` ([#15209](https://github.com/hashicorp/terraform-provider-azurerm/issues/15209))
* `azurerm_static_site` - fixing the creation of a Free tier Static Site ([#15141](https://github.com/hashicorp/terraform-provider-azurerm/issues/15141))
* `azurerm_storage_share` - fixing the `ShareBeingDeleted` error when the Storage Share is recreated ([#15180](https://github.com/hashicorp/terraform-provider-azurerm/issues/15180))
* 
## 2.94.0 (January 28, 2022)

UPGRADE NOTES:

* provider: support for the Azure German cloud has been removed in this release as this environment is no longer operational ([#14403](https://github.com/hashicorp/terraform-provider-azurerm/issues/14403))
* `azurerm_api_management_policy` - resources that were created with v2.92.0 will be marked as tainted due to a [bug](https://github.com/hashicorp/terraform-provider-azurerm/issues/15042). This version addresses the underlying issue, but the actual resource needs to either be untainted (via `terraform untaint`) or allow Terraform to delete the resource and create it again.
* `azurerm_hdinsight_kafka_cluster` - the `security_group_name` property in the `rest_proxy` block is conditionally required when the `use_msal` provider property is enabled ([#14403](https://github.com/hashicorp/terraform-provider-azurerm/issues/14403))

FEATURES:

* **New Data Source:** `azurerm_linux_function_app` ([#15009](https://github.com/hashicorp/terraform-provider-azurerm/issues/15009))
* **New Data Source** `azurerm_web_pubsub` ([#14731](https://github.com/hashicorp/terraform-provider-azurerm/issues/14731))
* **New Data Source** `azurerm_web_pubsub_hub` ([#14731](https://github.com/hashicorp/terraform-provider-azurerm/issues/14731))
* **New Resource:** `azurerm_web_pubsub` ([#14731](https://github.com/hashicorp/terraform-provider-azurerm/issues/14731))
* **New Resource:** `azurerm_web_pubsub_hub` ([#14731](https://github.com/hashicorp/terraform-provider-azurerm/issues/14731))
* **New Resource:** `azurerm_virtual_desktop_host_pool_registration_info` ([#14134](https://github.com/hashicorp/terraform-provider-azurerm/issues/14134))

ENHANCEMENTS:

* dependencies: updating to `v61.3.0` of `github.com/Azure/azure-sdk-for-go` ([#15080](https://github.com/hashicorp/terraform-provider-azurerm/issues/15080))
* dependencies: updating to `v0.21.0` of `github.com/hashicorp/go-azure-helpers` ([#15043](https://github.com/hashicorp/terraform-provider-azurerm/issues/15043))
* dependencies: updating `kusto` to API Version `2021-08-27` ([#15040](https://github.com/hashicorp/terraform-provider-azurerm/issues/15040))
* provider: opt-in support for v2 authentication tokens via the `use_msal` provider property ([#14403](https://github.com/hashicorp/terraform-provider-azurerm/issues/14403))
* `azurerm_app_service_slot`- support for the `storage_account` block ([#15084](https://github.com/hashicorp/terraform-provider-azurerm/issues/15084))
* `azurerm_stream_analytics_stream_input_eventhub` - support for the `partition_key` property ([#15019](https://github.com/hashicorp/terraform-provider-azurerm/issues/15019))

BUG FIXES:

* `data.image_source` - fix a regression around `id` ([#15119](https://github.com/hashicorp/terraform-provider-azurerm/issues/15119))
* `azurerm_api_management_backend` fix a crash caused by `backend_credentials` ([#15123](https://github.com/hashicorp/terraform-provider-azurerm/issues/15123))
* `azurerm_api_management_policy` - fixing the Resource ID for the `api_management_policy` block when this was provisioned using version `2.92.0` of the Azure Provider ([#15060](https://github.com/hashicorp/terraform-provider-azurerm/issues/15060))
* `azurerm_bastion_host` - fix a crash by adding nil check for the `copy_paste_enabled` property ([#15074](https://github.com/hashicorp/terraform-provider-azurerm/issues/15074))
* `azurerm_dev_test_lab` - fix an unexpected diff on with the `key_vault_id` property ([#15054](https://github.com/hashicorp/terraform-provider-azurerm/issues/15054))
* `azurerm_subscription_cost_management_export` - now sents the `ETag` when updating a cost management export ([#15017](https://github.com/hashicorp/terraform-provider-azurerm/issues/15017))
* `azurerm_template_deployment` - fixes a potential bug occuring during the deletion of a template deployment ([#15085](https://github.com/hashicorp/terraform-provider-azurerm/issues/15085))
* `azurerm_eventhub` - the `partition_count` property can now be changed when using Premium `sku` ([#15088](https://github.com/hashicorp/terraform-provider-azurerm/issues/15088))

## 2.93.1 (January 24, 2022)

BUG FIXES:

* `azurerm_app_service` - fix name availability check request ([#15062](https://github.com/hashicorp/terraform-provider-azurerm/issues/15062))

## 2.93.0 (January 21, 2022)

FEATURES:

* **New Data Source**: `azurerm_mysql_flexible_server` ([#14976](https://github.com/hashicorp/terraform-provider-azurerm/issues/14976))
* **New Beta Data Source**: `azurerm_windows_function_app` ([#14964](https://github.com/hashicorp/terraform-provider-azurerm/issues/14964))

ENHANCEMENTS: 

* dependencies: upgrading to `v61.1.0` of `github.com/Azure/azure-sdk-for-go` ([#14828](https://github.com/hashicorp/terraform-provider-azurerm/issues/14828))
* dependencies: updating `containerregistry` to API version `2021-08-01-preview` ([#14961](https://github.com/hashicorp/terraform-provider-azurerm/issues/14961))
* Data Source `azurerm_logic_app_workflow` - exporting the `identity` block ([#14896](https://github.com/hashicorp/terraform-provider-azurerm/issues/14896))
* `azurerm_bastion_host` - support for the `copy_paste_enabled`, `file_copy_enabled`, `ip_connect_enabled`, `shareable_link_enabled`, and `tunneling_enabled` properties ([#14987](https://github.com/hashicorp/terraform-provider-azurerm/issues/14987))
* `azurerm_bastion_host` - support for the `scale_units` property ([#14968](https://github.com/hashicorp/terraform-provider-azurerm/issues/14968))
* `azurerm_security_center_automation ` - the `event_source` property can now be set to `AssessmentsSnapshot`,
`RegulatoryComplianceAssessment`, `RegulatoryComplianceAssessmentSnapshot`, `SecureScoreControlsSnapshot`, `SecureScoresSnapshot`, and `SubAssessmentsSnapshot` ([#14996](https://github.com/hashicorp/terraform-provider-azurerm/issues/14996))
* `azurerm_static_site` - support for the `identity` block ([#14911](https://github.com/hashicorp/terraform-provider-azurerm/issues/14911))
* `azurerm_iothub` - Support for Identity-Based Endpoints ([#14705](https://github.com/hashicorp/terraform-provider-azurerm/issues/14705))
* `azurerm_servicebus_namespace_network_rule_set` -  support for the `public_network_access_enabled` property ([#14967](https://github.com/hashicorp/terraform-provider-azurerm/issues/14967))

BUG FIXES:

* `azurerm_machine_learning_compute_instance` - add validation for `tenant_id` and `object_id` properties to prevent null values and subsequent panic ([#14982](https://github.com/hashicorp/terraform-provider-azurerm/issues/14982))
* `azurerm_linux_function_app` - (beta) fix potential panic in `application_stack` when that block is not in config ([#14844](https://github.com/hashicorp/terraform-provider-azurerm/issues/14844))
* `azurerm_storage_share_file` changing the `content_md5` property will now trigger recreation and the `content_length` property of share file will now be set when updating properties. ([#15007](https://github.com/hashicorp/terraform-provider-azurerm/issues/15007))

## 2.92.0 (January 14, 2022)

FEATURES:

* **New Resource:** `azurerm_api_management_api_tag` ([#14711](https://github.com/hashicorp/terraform-provider-azurerm/issues/14711))
* **New Resource:** `azurerm_disk_pool_managed_disk_attachment` ([#14268](https://github.com/hashicorp/terraform-provider-azurerm/issues/14268))

ENHANCEMENTS:

* dependencies: upgrading `eventgrid` to API version `2021-12-01` ([#14433](https://github.com/hashicorp/terraform-provider-azurerm/issues/14433))
* `azurerm_api_management_custom_domain` - the `proxy` property has been deprecated in favour of the `gateway` for the 3.0 release ([#14628](https://github.com/hashicorp/terraform-provider-azurerm/issues/14628))
* `azurerm_databricks_workspace_customer_managed_key` - allow creation of resource when `infrastructure_encryption_enabled` is set to `true` for the databricks workspace ([#14915](https://github.com/hashicorp/terraform-provider-azurerm/issues/14915))
* `azurerm_eventgrid_domain` - support for the `local_auth_enabled`, `auto_create_topic_with_first_subscription`, and `auto_delete_topic_with_last_subscription` properties ([#14433](https://github.com/hashicorp/terraform-provider-azurerm/issues/14433))
* `azurerm_monitor_action_group` - support for the `event_hub_receiver` block ([#14771](https://github.com/hashicorp/terraform-provider-azurerm/issues/14771))
* `azurerm_mssql_server_extended_auditing_policy` - support storing audit data in storage account that is behind a firewall and VNet ([#14656](https://github.com/hashicorp/terraform-provider-azurerm/issues/14656))
* `azurerm_purview_account` - export the `managed_resources` block ([#14865](https://github.com/hashicorp/terraform-provider-azurerm/issues/14865))
* `azurerm_recovery_services_vault`- support for customer-managed keys (CMK) with the `encryption` block ([#14718](https://github.com/hashicorp/terraform-provider-azurerm/issues/14718))
* `azurerm_storage_account` - support for the `infrastructure_encryption_enabled` property ([#14864](https://github.com/hashicorp/terraform-provider-azurerm/issues/14864))

BUG FIXES:

* `azurerm_aadb2c_directory` - fix importing existing resources ([#14879](https://github.com/hashicorp/terraform-provider-azurerm/issues/14879))
* `azurerm_consumption_budget_subscription` - fix issue in migration logic ([#14898](https://github.com/hashicorp/terraform-provider-azurerm/issues/14898))
* `azurerm_cosmosdb_account` - only force ForceMongo when kind is set to mongo ([#14924](https://github.com/hashicorp/terraform-provider-azurerm/issues/14924))
* `azurerm_cosmosdb_mongo_collection` - now validates that "_id" is included as an index key ([#14857](https://github.com/hashicorp/terraform-provider-azurerm/issues/14857))
* `azurem_hdinsight` - hdinsight resources using oozie metastore can now be created without error ([#14880](https://github.com/hashicorp/terraform-provider-azurerm/issues/14880))
* `azurerm_log_analytics_datasource_windows_performance_counter` - state migration for case conversion of ID element ([#14916](https://github.com/hashicorp/terraform-provider-azurerm/issues/14916))
* `azurerm_monitor_aad_diagnostic_setting` - use the correct parser function for event hub rule IDs ([#14944](https://github.com/hashicorp/terraform-provider-azurerm/issues/14944))
* `azurerm_mysql_server_key` - fix issue when checking for existing resource on create ([#14883](https://github.com/hashicorp/terraform-provider-azurerm/issues/14883))
* `azurerm_spring_cloud_service` - fix panic when removing git repos ([#14900](https://github.com/hashicorp/terraform-provider-azurerm/issues/14900))
* `azurerm_log_analytics_workspace` - the `reservation_capcity_in_gb_per_day` has been deprecated and renamed to `reservation_capacity_in_gb_per_day` ([#14910](https://github.com/hashicorp/terraform-provider-azurerm/issues/14910))
* `azurerm_iothub_dps` - fixed default value of `allocation_weight` to match azure default ([#14943](https://github.com/hashicorp/terraform-provider-azurerm/issues/14943))
* `azurerm_iothub` - now exports `event_hub_events_namespace` and has a fallback route by default ([#14942](https://github.com/hashicorp/terraform-provider-azurerm/issues/14942))

## 2.91.0 (January 07, 2022)

FEATURES:

* **New Data Source:** `azurerm_aadb2c_directory` ([#14671](https://github.com/hashicorp/terraform-provider-azurerm/issues/14671))
* **New Data Source:** `azurerm_sql_managed_instance` ([#14739](https://github.com/hashicorp/terraform-provider-azurerm/issues/14739))
* **New Resource:** `azurerm_aadb2c_directory` ([#14671](https://github.com/hashicorp/terraform-provider-azurerm/issues/14671))
* **New Resource:** `azurerm_app_service_slot_custom_hostname_binding` ([#13097](https://github.com/hashicorp/terraform-provider-azurerm/issues/13097))
* **New Resource:** `azurerm_data_factory_linked_service_odbc` ([#14787](https://github.com/hashicorp/terraform-provider-azurerm/issues/14787))
* **New Resource:** `azurerm_disk_pool` ([#14675](https://github.com/hashicorp/terraform-provider-azurerm/issues/14675))
* **New Resource:** `azurerm_load_test` ([#14724](https://github.com/hashicorp/terraform-provider-azurerm/issues/14724))
* **New Resource:** `azurerm_virtual_desktop_scaling_plan` ([#14188](https://github.com/hashicorp/terraform-provider-azurerm/issues/14188))

ENHANCEMENTS:

* dependencies: upgrading `appplatform` to API version `2021-09-01-preview` ([#14365](https://github.com/hashicorp/terraform-provider-azurerm/issues/14365))
* dependencies: upgrading `network` to API Version `2021-05-01` ([#14164](https://github.com/hashicorp/terraform-provider-azurerm/issues/14164))
* dependencies: upgrading to `v60.2.0` of `github.com/Azure/azure-sdk-for-go` ([#14688](https://github.com/hashicorp/terraform-provider-azurerm/issues/14688)] and [[#14667](https://github.com/hashicorp/terraform-provider-azurerm/issues/14667))
* dependencies: upgrading to `v2.10.1` of `github.com/hashicorp/terraform-plugin-sdk` ([#14666](https://github.com/hashicorp/terraform-provider-azurerm/issues/14666))
* `azurerm_application_gateway` - support for the `key_vault_secret_id` and `force_firewall_policy_association` properties ([#14413](https://github.com/hashicorp/terraform-provider-azurerm/issues/14413))
* `azurerm_application_gateway` - support the `fips_enagled` property ([#14797](https://github.com/hashicorp/terraform-provider-azurerm/issues/14797))
* `azurerm_cdn_endpoint_custom_domain` - support for HTTPS ([#13283](https://github.com/hashicorp/terraform-provider-azurerm/issues/13283))
* `azurerm_hdinsight_hbase_cluster` - support for the `network` property ([#14825](https://github.com/hashicorp/terraform-provider-azurerm/issues/14825))
* `azurerm_iothub` - support for the `identity` block ([#14354](https://github.com/hashicorp/terraform-provider-azurerm/issues/14354))
* `azurerm_iothub_endpoint_servicebus_queue_resource` - depracating the `iothub_name` propertyin favour of `iothub_id` property ([#14690](https://github.com/hashicorp/terraform-provider-azurerm/issues/14690))
* `azurerm_iothub_endpoint_storage_container_resource` - depracating the `iothub_name` property in favour of `iothub_id` property [[#14690](https://github.com/hashicorp/terraform-provider-azurerm/issues/14690)] 
* `azurerm_iot_fallback_route` - support for the `source` property ([#14836](https://github.com/hashicorp/terraform-provider-azurerm/issues/14836))
* `azurerm_kubernetes_cluster` - support for the `public_network_access_enabled`, `scale_down_mode`, and `workload_runtime` properties ([#14386](https://github.com/hashicorp/terraform-provider-azurerm/issues/14386))
* `azurerm_linux_function_app` - (Beta Resource) fix the filtering of `app_settings` for `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#14815](https://github.com/hashicorp/terraform-provider-azurerm/issues/14815))
* `azurerm_linux_virtual_machine` - support for the `user_data` property ([#13888](https://github.com/hashicorp/terraform-provider-azurerm/issues/13888))
* `azurerm_linux_virtual_machine_scale_set` - support for the `user_data` property ([#13888](https://github.com/hashicorp/terraform-provider-azurerm/issues/13888))
* `azurerm_managed_disk` - support for the `gallery_image_reference_id` property ([#14121](https://github.com/hashicorp/terraform-provider-azurerm/issues/14121))
* `azurerm_mysql_server` - support capacities up to `16TB` for the `storage_mb` property ([#14838](https://github.com/hashicorp/terraform-provider-azurerm/issues/14838))
* `azurerm_postgresql_flexible_server` - support for the `geo_redundant_backup_enabled` property ([#14661](https://github.com/hashicorp/terraform-provider-azurerm/issues/14661))
* `azurerm_recovery_services_vault` - support for the `storage_mode_type` property ([#14659](https://github.com/hashicorp/terraform-provider-azurerm/issues/14659))
* `azurerm_spring_cloud_certificate` - support for the `certificate_content` property ([#14689](https://github.com/hashicorp/terraform-provider-azurerm/issues/14689))
* `azurerm_servicebus_namespace_authorization_rule` - the `resource_group_name` and `namespace_name` properties have been deprecated in favour of the `namespace_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_namespace_network_rule_set` - the `resource_group_name` and `namespace_name` properties have been deprecated in favour of the `namespace_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_namespace_authorization_rule` - the `resource_group_name` and `namespace_name` properties have been deprecated in favour of the `namespace_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_queue` - the `resource_group_name` and `namespace_name` properties have been deprecated in favour of the `namespace_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_queue_authorization_rule` - the `resource_group_name`, `namespace_name`, and `queue_name` properties have been deprecated in favour of the `queue_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_subscription` - the `resource_group_name`, `namespace_name`, and `topic_name` properties have been deprecated in favour of the `topic_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_subscription_rule` - the `resource_group_name`, `namespace_name`, `topic_name`, and `subscription_name` properties have been deprecated in favour of the `subscription_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_topic` - the `resource_group_name` and `namespace_name` properties have been deprecated in favour of the `namespace_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_servicebus_topic_authorization_rule` - the `resource_group_name`, `namespace_name`, and `topic_name` properties have been deprecated in favour of the `topic_id` property ([#14784](https://github.com/hashicorp/terraform-provider-azurerm/issues/14784))
* `azurerm_shared_image_version` - images can now be sorted by semver ([#14708](https://github.com/hashicorp/terraform-provider-azurerm/issues/14708))
* `azurerm_virtual_network_gateway_connection` - support for the `connection_mode` property ([#14738](https://github.com/hashicorp/terraform-provider-azurerm/issues/14738))
* `azurerm_web_application_firewall_policy` - the `file_upload_limit_in_mb` property within the `policy_settings` block can now be set to `4000` ([#14715](https://github.com/hashicorp/terraform-provider-azurerm/issues/14715))
* `azurerm_windows_virtual_machine` - support for the `user_data` property ([#13888](https://github.com/hashicorp/terraform-provider-azurerm/issues/13888))
* `azurerm_windows_virtual_machine_scale_set` - support for the `user_data` property ([#13888](https://github.com/hashicorp/terraform-provider-azurerm/issues/13888))

BUG FIXES:

* `azurerm_app_service_environment_v3` - fix the default value of the `allow_new_private_endpoint_connections` property ([#14805](https://github.com/hashicorp/terraform-provider-azurerm/issues/14805))
* `azurerm_consumption_budget_subscription` - added an additional state migration to fix the bug introduced by the first one and to parse the `subscription_id` from the resource's ID ([#14803](https://github.com/hashicorp/terraform-provider-azurerm/issues/14803))
* `azurerm_network_interface_security_group_association` - checking the ID matches the expected format during import ([#14753](https://github.com/hashicorp/terraform-provider-azurerm/issues/14753))
* `azurerm_storage_management_policy` - handle the unexpected deletion of the storage account ([#14799](https://github.com/hashicorp/terraform-provider-azurerm/issues/14799))

## 2.90.0 (December 17, 2021)

FEATURES:

* **New Data Source:** `azurerm_app_configuration_key` ([#14484](https://github.com/hashicorp/terraform-provider-azurerm/issues/14484))
* **New Resource:** `azurerm_container_registry_task` ([#14533](https://github.com/hashicorp/terraform-provider-azurerm/issues/14533))
* **New Resource:** `azurerm_maps_creator` ([#14566](https://github.com/hashicorp/terraform-provider-azurerm/issues/14566))
* **New Resource:** `azurerm_netapp_snapshot_policy` ([#14230](https://github.com/hashicorp/terraform-provider-azurerm/issues/14230))
* **New Resource:** `azurerm_synapse_sql_pool_workload_classifier` ([#14412](https://github.com/hashicorp/terraform-provider-azurerm/issues/14412))
* **New Resource:** `azurerm_synapse_workspace_sql_aad_admin` ([#14341](https://github.com/hashicorp/terraform-provider-azurerm/issues/14341))
* **New Resource:** `azurerm_vpn_gateway_nat_rule` ([#14527](https://github.com/hashicorp/terraform-provider-azurerm/issues/14527))

ENHANCEMENTS:

* dependencies: updating `apimanagement` to API Version `2021-08-01` ([#14312](https://github.com/hashicorp/terraform-provider-azurerm/issues/14312))
* dependencies: updating `managementgroups` to API Version `2020-05-01` ([#14635](https://github.com/hashicorp/terraform-provider-azurerm/issues/14635))
* dependencies: updating `redisenterprise` to use an Embedded SDK ([#14502](https://github.com/hashicorp/terraform-provider-azurerm/issues/14502))
* dependencies: updating to `v0.19.1` of `github.com/hashicorp/go-azure-helpers` ([#14627](https://github.com/hashicorp/terraform-provider-azurerm/issues/14627))
* dependencies: updating to `v2.10.0` of `github.com/hashicorp/terraform-plugin-sdk` ([#14596](https://github.com/hashicorp/terraform-provider-azurerm/issues/14596))
* Data Source: `azurerm_function_app_host_keys` - support for `signalr_extension_key` and `durabletask_extension_key` ([#13648](https://github.com/hashicorp/terraform-provider-azurerm/issues/13648))
* `azurerm_application_gateway ` - support for private link configurations ([#14583](https://github.com/hashicorp/terraform-provider-azurerm/issues/14583))
* `azurerm_blueprint_assignment` - support for the `lock_exclude_actions` property ([#14648](https://github.com/hashicorp/terraform-provider-azurerm/issues/14648))
* `azurerm_container_group` - support for `ip_address_type = None` ([#14460](https://github.com/hashicorp/terraform-provider-azurerm/issues/14460))
* `azurerm_cosmosdb_account` - support for the `create_mode` property and `restore` block ([#14362](https://github.com/hashicorp/terraform-provider-azurerm/issues/14362))
* `azurerm_data_factory_dataset_*` - deprecate `data_factory_name` in favour of `data_factory_id` for consistency across all data factory dataset resources ([#14610](https://github.com/hashicorp/terraform-provider-azurerm/issues/14610))
* `azurerm_data_factory_integration_runtime_*`- deprecate `data_factory_name` in favour of `data_factory_id` for consistency across all data factory integration runtime resources ([#14610](https://github.com/hashicorp/terraform-provider-azurerm/issues/14610))
* `azurerm_data_factory_trigger_*`- deprecate `data_factory_name` in favour of `data_factory_id` for consistency across all data factory trigger resources ([#14610](https://github.com/hashicorp/terraform-provider-azurerm/issues/14610))
* `azurerm_data_factory_pipeline`- deprecate `data_factory_name` in favour of `data_factory_id` for consistency across all data factory resources ([#14610](https://github.com/hashicorp/terraform-provider-azurerm/issues/14610))
* `azurerm_iothub` - support for the `cloud_to_device` block ([#14546](https://github.com/hashicorp/terraform-provider-azurerm/issues/14546))
* `azurerm_iothub_endpoint_eventhub` - the `iothub_name` property has been deprecated in favour of the `iothub_id` property ([#14632](https://github.com/hashicorp/terraform-provider-azurerm/issues/14632))
* `azurerm_logic_app_workflow` - support for the `open_authentication_policy` block ([#14007](https://github.com/hashicorp/terraform-provider-azurerm/issues/14007))
* `azurerm_signalr` - support for the `live_trace_enabled` property ([#14646](https://github.com/hashicorp/terraform-provider-azurerm/issues/14646))
* `azurerm_xyz_policy_assignment` add support for `non_compliance_message` ([#14518](https://github.com/hashicorp/terraform-provider-azurerm/issues/14518))

BUG FIXES:

* `azurerm_cosmosdb_account` - will now set a default value for `default_identity_type` when the API return a nil value ([#14643](https://github.com/hashicorp/terraform-provider-azurerm/issues/14643))
* `azurerm_function_app` - address `app_settings` during creation rather than just updates ([#14638](https://github.com/hashicorp/terraform-provider-azurerm/issues/14638))
* `azurerm_marketplace_agreement` - fix crash when the import check triggers ([#14614](https://github.com/hashicorp/terraform-provider-azurerm/issues/14614))
* `azurerm_postgresql_configuration` - now locks during write operations to prevent conflicts ([#14619](https://github.com/hashicorp/terraform-provider-azurerm/issues/14619))
* `azurerm_postgresql_flexible_server_configuration` - now locks during write operations to prevent conflicts ([#14607](https://github.com/hashicorp/terraform-provider-azurerm/issues/14607))

---

For information on changes between the v2.89.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v2.00.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
