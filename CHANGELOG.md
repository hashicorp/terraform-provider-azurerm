## 3.104.2 (May 20, 2024)

NOTE: This is a re-release of `v3.104.1` to include missing changes, please refer to the changelog entries for `v3.104.1`.

## 3.104.1 (May 20, 2024)

BUG FIXES:

* `azurerm_pim_active_role_assignment` - fix a regression where roles assignments could not be created with no expiration ([#26029](https://github.com/hashicorp/terraform-provider-azurerm/issues/26029))
* `azurerm_pim_eligible_role_assignment` - fix a regression where roles assignments could not be created with no expiration ([#26029](https://github.com/hashicorp/terraform-provider-azurerm/issues/26029))

## 3.104.0 (May 16, 2024)

FEATURES:

* New Data Source: `azurerm_elastic_san` ([#25719](https://github.com/hashicorp/terraform-provider-azurerm/issues/25719))

ENHANCEMENTS:

* New Resource - `azurerm_key_vault_managed_hardware_security_module_key` ([#25935](https://github.com/hashicorp/terraform-provider-azurerm/issues/25935))
* Data Source - `azurerm_kubernetes_service_version` - support for the `default_version` property ([#25953](https://github.com/hashicorp/terraform-provider-azurerm/issues/25953))
* `network/applicationgateways` - update to use `hashicorp/go-azure-sdk` ([#25844](https://github.com/hashicorp/terraform-provider-azurerm/issues/25844))
* `dataprotection` - update API version to `2024-04-01` ([#25882](https://github.com/hashicorp/terraform-provider-azurerm/issues/25882))
* `databasemigration` - update API version to `2021-06-30` ([#25997](https://github.com/hashicorp/terraform-provider-azurerm/issues/25997))
* `network/ips` - update to use `hashicorp/go-azure-sdk` ([#25905](https://github.com/hashicorp/terraform-provider-azurerm/issues/25905))
* `network/localnetworkgateway` - update to use `hashicorp/go-azure-sdk` ([#25905](https://github.com/hashicorp/terraform-provider-azurerm/issues/25905))
* `network/natgateway` - update to use `hashicorp/go-azure-sdk` ([#25905](https://github.com/hashicorp/terraform-provider-azurerm/issues/25905))
* `network/networksecuritygroup` - update to use `hashicorp/go-azure-sdk` ([#25971](https://github.com/hashicorp/terraform-provider-azurerm/issues/25971))
* `network/publicips` - update to use `hashicorp/go-azure-sdk` ([#25971](https://github.com/hashicorp/terraform-provider-azurerm/issues/25971))
* `network/virtualwan` - update to use `hashicorp/go-azure-sdk` ([#25971](https://github.com/hashicorp/terraform-provider-azurerm/issues/25971))
* `network/vpn` - update to use `hashicorp/go-azure-sdk` ([#25971](https://github.com/hashicorp/terraform-provider-azurerm/issues/25971))
* `azurerm_databricks_workspace` - support for the `default_storage_firewall_enabled` property ([#25919](https://github.com/hashicorp/terraform-provider-azurerm/issues/25919))
* `azurerm_key_vault` - allow previously existing key vaults to continue to manage the `contact` field prior to the `v3.93.0` conditional polling change ([#25777](https://github.com/hashicorp/terraform-provider-azurerm/issues/25777))
* `azurerm_linux_function_app` - support for the PowerShell `7.4` ([#25980](https://github.com/hashicorp/terraform-provider-azurerm/issues/25980))
* `azurerm_log_analytics_cluster` - support for the value `UserAssigned` in the `identity.type` property ([#25940](https://github.com/hashicorp/terraform-provider-azurerm/issues/25940))
* `azurerm_pim_active_role_assignment` - remove hard dependency on the `roleAssignmentScheduleRequests` API, so that role assignments will not become unmanageable over time ([#25956](https://github.com/hashicorp/terraform-provider-azurerm/issues/25956))
* `azurerm_pim_eligible_role_assignment` - remove hard dependency on the `roleEligibilityScheduleRequests` API, so that role assignments will not become unmanageable over time ([#25956](https://github.com/hashicorp/terraform-provider-azurerm/issues/25956))
* `azurerm_windows_function_app` - support for the PowerShell `7.4` ([#25980](https://github.com/hashicorp/terraform-provider-azurerm/issues/25980))

BUG FIXES:

* `azurerm_container_app_job` - Allow `event_trigger_config.scale.min_executions` to be `0` ([#25931](https://github.com/hashicorp/terraform-provider-azurerm/issues/25931))
* `azurerm_container_app_job` - update validation to allow the `replica_retry_limit` property to be set to `0` ([#25984](https://github.com/hashicorp/terraform-provider-azurerm/issues/25984))
* `azurerm_data_factory_trigger_custom_event` - one of `subject_begins_with` and `subject_ends_with` no longer need to be set ([#25932](https://github.com/hashicorp/terraform-provider-azurerm/issues/25932))
* `azurerm_kubernetes_cluster_node_pool` - prevent race condition by checking the virtual network status when creating a node pool with a subnet ID ([#25888](https://github.com/hashicorp/terraform-provider-azurerm/issues/25888))
* `azurerm_postgresql_flexible_server` - fix for default `storage_tier` value when `storage_mb` field has been changed ([#25947](https://github.com/hashicorp/terraform-provider-azurerm/issues/25947))
* `azurerm_pim_active_role_assignment` - resolve a number of potential crashes ([#25956](https://github.com/hashicorp/terraform-provider-azurerm/issues/25956))
* `azurerm_pim_eligible_role_assignment` - resolve a number of potential crashes ([#25956](https://github.com/hashicorp/terraform-provider-azurerm/issues/25956))
* `azurerm_redis_enterprise_cluster_location_zone_support` - add `Central India` zones support ([#26000](https://github.com/hashicorp/terraform-provider-azurerm/issues/26000))
* `azurerm_sentinel_alert_rule_scheduled` - the `alert_rule_template_version` property is no longer `ForceNew` ([#25688](https://github.com/hashicorp/terraform-provider-azurerm/issues/25688))
* `azurerm_storage_sync_server_endpoint` - preventing a crashed due to `initial_upload_policy`  ([#25968](https://github.com/hashicorp/terraform-provider-azurerm/issues/25968))

## 3.103.1 (May 10, 2024)

BUG FIXES

* `loadtest` - fixing an issue where the SDK Clients weren't registered ([#25920](https://github.com/hashicorp/terraform-provider-azurerm/issues/25920))

## 3.103.0 (May 09, 2024)

FEATURES:

* New Resource: `azurerm_container_app_job` ([#23871](https://github.com/hashicorp/terraform-provider-azurerm/issues/23871))
* New Resource: `azurerm_container_app_environment_custom_domain` ([#24346](https://github.com/hashicorp/terraform-provider-azurerm/issues/24346))
* New Resource: `azurerm_data_factory_credential_service_principal` ([#25805](https://github.com/hashicorp/terraform-provider-azurerm/issues/25805))
* New Resource: `azurerm_network_manager_connectivity_configuration` ([#25746](https://github.com/hashicorp/terraform-provider-azurerm/issues/25746))
* New Resource: `azurerm_maintenance_assignment_dynamic_scope` ([#25467](https://github.com/hashicorp/terraform-provider-azurerm/issues/25467))
* New Resource: `azurerm_virtual_machine_gallery_application_assignment` ([#22945](https://github.com/hashicorp/terraform-provider-azurerm/issues/22945))
* New Resource: `azurerm_virtual_machine_automanage_configuration_assignment` ([#25480](https://github.com/hashicorp/terraform-provider-azurerm/issues/25480))

ENHANCEMENTS:

* provider - support for the `recover_soft_deleted_backup_protected_vm` feature ([#24157](https://github.com/hashicorp/terraform-provider-azurerm/issues/24157))
* dependencies: updating `github.com/hashicorp/go-azure-helpers` to `v0.69.0` ([#25903](https://github.com/hashicorp/terraform-provider-azurerm/issues/25903))
* `loganalytics` - update cluster resource to api version `2022-01-01` ([#25686](https://github.com/hashicorp/terraform-provider-azurerm/issues/25686))
* `azurerm_bastion_host` - support for the `kerberos_enabled` property ([#25823](https://github.com/hashicorp/terraform-provider-azurerm/issues/25823))
* `azurerm_container_app` - secrets can now be removed ([#25743](https://github.com/hashicorp/terraform-provider-azurerm/issues/25743))
* `azurerm_container_app_environment` - support for the `custom_domain_verification_id` property ([#24346](https://github.com/hashicorp/terraform-provider-azurerm/issues/24346))
* `azurerm_linux_virtual_machine` - support for the additional capability `hibernation_enabled` ([#25807](https://github.com/hashicorp/terraform-provider-azurerm/issues/25807))
* `azurerm_linux_virtual_machine` - support for additional values for the `license_type` property ([#25909](https://github.com/hashicorp/terraform-provider-azurerm/issues/25909))
* `azurerm_linux_virtual_machine_scale_set` - support for the `maximum_surge_instances` property for vmss rolling upgrades ([#24914](https://github.com/hashicorp/terraform-provider-azurerm/issues/24914))
* `azurerm_windows_virtual_machine` - support for the additional capability `hibernation_enabled` ([#25807](https://github.com/hashicorp/terraform-provider-azurerm/issues/25807))
* `azurerm_windows_virtual_machine_scale_set` - support for the `maximum_surge_instances_enabled` property for vmss rolling upgrades ([#24914](https://github.com/hashicorp/terraform-provider-azurerm/issues/24914))
* `azurerm_storage_account` - support for the `permanent_delete_enabled` property within retention policies ([#25778](https://github.com/hashicorp/terraform-provider-azurerm/issues/25778))

BUG FIXES:

* `azurerm_kubernetes_cluster` - erase `load_balancer_profile` when changing `network_profile.outbound_type` from `loadBalancer` to another outbound type ([#25530](https://github.com/hashicorp/terraform-provider-azurerm/issues/25530))
* `azurerm_log_analytics_saved_search` - the `function_parameters` property now repsects the order of elements ([#25869](https://github.com/hashicorp/terraform-provider-azurerm/issues/25869))
* `azurerm_linux_web_app` - fix `slow_request` with `path` issue in `auto_heal` by adding support for `slow_request_with_path` block ([#20049](https://github.com/hashicorp/terraform-provider-azurerm/issues/20049))
* `azurerm_linux_web_app_slot` - fix `slow_request` with `path` issue in `auto_heal` by adding support for `slow_request_with_path` block ([#20049](https://github.com/hashicorp/terraform-provider-azurerm/issues/20049))
* `azurerm_monitor_private_link_scoped_service` - normalize case of the `linked_resource_id` property during reads  ([#25787](https://github.com/hashicorp/terraform-provider-azurerm/issues/25787))
* `azurerm_role_assignment` - add addtional retry logic to assist with cross-tenant use ([#25853](https://github.com/hashicorp/terraform-provider-azurerm/issues/25853))
* `azurerm_web_pubsub_network_acl` - fixing a crash when `networkACL.PublicNetwork.Deny` was nil ([#25886](https://github.com/hashicorp/terraform-provider-azurerm/issues/25886))
* `azurerm_windows_web_app` - fix `slow_request` with `path` issue in `auto_heal` by adding support for `slow_request_with_path` block ([#20049](https://github.com/hashicorp/terraform-provider-azurerm/issues/20049))
* `azurerm_windows_web_app_slot` - fix `slow_request` with `path` issue in `auto_heal` by adding support for `slow_request_with_path` block ([#20049](https://github.com/hashicorp/terraform-provider-azurerm/issues/20049))

DEPRECATIONS:
* `azurerm_subnet` - the `private_endpoint_network_policies_enabled` property has been deprecated in favour of the `private_endpoint_network_policies` property ([#25779](https://github.com/hashicorp/terraform-provider-azurerm/issues/25779))

## 3.102.0 (May 02, 2024)

FEATURES:

* New Resource: `azurerm_storage_sync_server_endpoint` ([#25831](https://github.com/hashicorp/terraform-provider-azurerm/issues/25831))
* New Resource: `azurerm_storage_container_immutability_policy` ([#25804](https://github.com/hashicorp/terraform-provider-azurerm/issues/25804))

ENHANCEMENTS:

* `azurerm_load_test` - add support for `encryption` ([#25759](https://github.com/hashicorp/terraform-provider-azurerm/issues/25759))
* `azurerm_network_connection_monitor` - update validation for `target_resource_type` and `target_resource_id` ([#25745](https://github.com/hashicorp/terraform-provider-azurerm/issues/25745))
* `azurerm_mssql_managed_database` - support for a Restorable Database ID to be used as the `source_database_id` for point in time restore ([#25568](https://github.com/hashicorp/terraform-provider-azurerm/issues/25568))
* `azurerm_storage_account` - support for the `managed_hsm_key_id` property ([#25088](https://github.com/hashicorp/terraform-provider-azurerm/issues/25088))
* `azurerm_storage_account_customer_managed_key` - support for the `managed_hsm_key_id` property ([#25088](https://github.com/hashicorp/terraform-provider-azurerm/issues/25088))

BUG FIXES:

* `azurerm_linux_function_app` - now sets docker registry url in `linux_fx_version` by default ([#23911](https://github.com/hashicorp/terraform-provider-azurerm/issues/23911))
* `azurerm_resource_group` - work around sporadic eventual consistency errors ([#25758](https://github.com/hashicorp/terraform-provider-azurerm/issues/25758))

DEPRECATIONS:

* `azurerm_key_vault_managed_hardware_security_module_role_assignment` - the `vault_base_url` property has been deprecated in favour of the `managed_hsm_id` property ([#25601](https://github.com/hashicorp/terraform-provider-azurerm/issues/25601))

## 3.101.0 (April 25, 2024)

ENHANCEMENTS:

* dependencies: updating to `v0.20240424.1114424` of `github.com/hashicorp/go-azure-sdk` ([#25749](https://github.com/hashicorp/terraform-provider-azurerm/issues/25749))
* dependencies: updating to `v0.27.0` of `github.com/tombuildsstuff/giovanni` ([#25702](https://github.com/hashicorp/terraform-provider-azurerm/issues/25702))
* dependencies: updating `golang.org/x/net` to `0.23.0`
* `azurerm_cognitive_account` - the `kind` property now supports `ConversationalLanguageUnderstanding` ([#25735](https://github.com/hashicorp/terraform-provider-azurerm/issues/25735))
* `azurerm_container_app_custom_domain` - support the ability to use Azure Managed Certificates ([#25356](https://github.com/hashicorp/terraform-provider-azurerm/issues/25356))

BUG FIXES:

* Data Source: `azurerm_application_insights` - set correct AppID in data source ([#25687](https://github.com/hashicorp/terraform-provider-azurerm/issues/25687))
* `azurerm_virtual_network` - suppress diff in ordering for `address_space` due to inconsistent API response ([#23793](https://github.com/hashicorp/terraform-provider-azurerm/issues/23793))
* `azurerm_storage_data_lake_gen2_filesystem` - add context deadline for import ([#25712](https://github.com/hashicorp/terraform-provider-azurerm/issues/25712))
* `azurerm_virtual_network_gateway` - preserve existing `nat_rules` on updates ([#25690](https://github.com/hashicorp/terraform-provider-azurerm/issues/25690))

## 3.100.0 (April 18, 2024)

ENHANCEMENTS:

* dependencies: updating `hashicorp/go-azure-sdk` to `v0.20240417.1084633` ([#25659](https://github.com/hashicorp/terraform-provider-azurerm/issues/25659))
* `compute` - update Virtual Machine and Virtual Machine Scale Set resources and data sources to use `hashicorp/go-azure-sdk` ([#25533](https://github.com/hashicorp/terraform-provider-azurerm/issues/25533))
* `machine_learning` - Add new `machine_learning` block that supports `purge_soft_deleted_workspace_on_destroy` ([#25624](https://github.com/hashicorp/terraform-provider-azurerm/issues/25624))
* `loganalytics` - update cluster resource to use `hashicorp/go-azure-sdk` ([#23373](https://github.com/hashicorp/terraform-provider-azurerm/issues/23373))
* Data Source: `azurerm_management_group` - now exports the `tenant_scoped_id` attribute ([#25555](https://github.com/hashicorp/terraform-provider-azurerm/issues/25555))
* `azurerm_container_app` - the `ingress.ip_security_restriction.ip_address_range` property will now accept an IP address as valid input ([#25609](https://github.com/hashicorp/terraform-provider-azurerm/issues/25609))
* `azurerm_container_group` - the `identity` block can now be updated ([#25543](https://github.com/hashicorp/terraform-provider-azurerm/issues/25543))
* `azurerm_express_route_connection` - support for the `private_link_fast_path_enabled` property ([#25596](https://github.com/hashicorp/terraform-provider-azurerm/issues/25596))
* `azurerm_hdinsight_hadoop_cluster` - support for the `private_link_configuration` block ([#25629](https://github.com/hashicorp/terraform-provider-azurerm/issues/25629))
* `azurerm_hdinsight_hbase_cluster` - support for the `private_link_configuration` block ([#25629](https://github.com/hashicorp/terraform-provider-azurerm/issues/25629))
* `azurerm_hdinsight_interactive_query_cluster` - support for the `private_link_configuration` block ([#25629](https://github.com/hashicorp/terraform-provider-azurerm/issues/25629))
* `azurerm_hdinsight_kafka_cluster` - support for the `private_link_configuration` block ([#25629](https://github.com/hashicorp/terraform-provider-azurerm/issues/25629))
* `azurerm_hdinsight_spark_cluster` - support for the `private_link_configuration` block ([#25629](https://github.com/hashicorp/terraform-provider-azurerm/issues/25629))
* `azurerm_management_group` - now exports the `tenant_scoped_id` attribute ([#25555](https://github.com/hashicorp/terraform-provider-azurerm/issues/25555))
* `azurerm_monitor_activity_log_alert` - support for the `location` property ([#25389](https://github.com/hashicorp/terraform-provider-azurerm/issues/25389))
* `azurerm_mysql_flexible_server` - update validating regex for `sku_name` ([#25642](https://github.com/hashicorp/terraform-provider-azurerm/issues/25642))
* `azurerm_postgresql_flexible_server` - support for the `GeoRestore` `create_mode` ([#25664](https://github.com/hashicorp/terraform-provider-azurerm/issues/25664))
* `azurerm_virtual_network_gateway_connection` - support for the `private_link_fast_path_enabled` property ([#25650](https://github.com/hashicorp/terraform-provider-azurerm/issues/25650))
* `azurerm_windows_web_app` - support for the `handler_mapping` block ([#25631](https://github.com/hashicorp/terraform-provider-azurerm/issues/25631))
* `azurerm_windows_web_app_slot` - support for the `handler_mapping` block ([#25631](https://github.com/hashicorp/terraform-provider-azurerm/issues/25631))

BUG FIXES:

* storage: prevent a bug causing the second storage account key to be used for authentication instead of the first ([#25652](https://github.com/hashicorp/terraform-provider-azurerm/issues/25652))
* `azurerm_active_directory_domain_service` - prevent an issue where `filtered_sync_enabled` was not being updated ([#25594](https://github.com/hashicorp/terraform-provider-azurerm/issues/25594))
* `azurerm_application_insights` - add a state migration to fix the resource ID casing of Application Insights resources ([#25628](https://github.com/hashicorp/terraform-provider-azurerm/issues/25628))
* `azurerm_function_app_hybrid_connection` - can now use relay resources created in a different resource group ([#25541](https://github.com/hashicorp/terraform-provider-azurerm/issues/25541))
* `azurerm_kubernetes_cluster_node_pool` - prevent plan diff when the `windows_profile.outbound_nat_enabled` property is unset ([#25644](https://github.com/hashicorp/terraform-provider-azurerm/issues/25644))
* `azurerm_machine_learning_compute_cluster` - fix location to point to parent resource for computes ([#25643](https://github.com/hashicorp/terraform-provider-azurerm/issues/25643))
* `azurerm_machine_learning_compute_instance` - fix location to point to parent resource for computes ([#25643](https://github.com/hashicorp/terraform-provider-azurerm/issues/25643))
* `azurerm_storage_account` - check replication type when evaluating support level for shares and queues for V1 storage accounts ([#25581](https://github.com/hashicorp/terraform-provider-azurerm/issues/25581))
* `azurerm_storage_account` - added a sanity check for `dns_endpoint_type` and `blob_properties.restore_policy` ([#25450](https://github.com/hashicorp/terraform-provider-azurerm/issues/25450))
* `azurerm_web_app_hybrid_connection` - can now use relay resources created in a different resource group ([#25541](https://github.com/hashicorp/terraform-provider-azurerm/issues/25541))
* `azurerm_windows_web_app` - prevent removal of `site_config.application_stack.node_version` when `app_settings` are updated ([#25488](https://github.com/hashicorp/terraform-provider-azurerm/issues/25488))
* `azurerm_windows_web_app_slot` - prevent removal of `site_config.application_stack.node_version` when `app_settings` are updated ([#25489](https://github.com/hashicorp/terraform-provider-azurerm/issues/25489))

DEPRECATIONS:

* `logz` - the Logz resources are deprecated and will be removed in v4.0 of the AzureRM Provider since the API no longer allows new instances to be created ([#25405](https://github.com/hashicorp/terraform-provider-azurerm/issues/25405))
* `azurerm_machine_learning_compute_instance` - marked the `location` field as deprecated in v4.0 of the provider ([#25643](https://github.com/hashicorp/terraform-provider-azurerm/issues/25643))
* `azurerm_kubernetes_cluster` - the following properties have been deprecated since the API no longer supports cluster creation with legacy Azure Entra integration: `client_app_id`, `server_app_id`, `server_app_secret` and `managed` ([#25200](https://github.com/hashicorp/terraform-provider-azurerm/issues/25200))

## 3.99.0 (April 11, 2024)

BREAKING CHANGE: 

* `azurerm_linux_web_app` - `site_config.0.application_stack.0.java_version` must be specified with `java_server` and `java_server_version` ([#25553](https://github.com/hashicorp/terraform-provider-azurerm/issues/25553))

ENHANCEMENTS:

* dependencies: updating to `v0.20240411.1104331` of `github.com/hashicorp/go-azure-sdk/resourcemanager` and `github.com/hashicorp/go-azure-sdk/sdk` ([#25546](https://github.com/hashicorp/terraform-provider-azurerm/issues/25546))
* dependencies: updating to `v0.26.1` of `github.com/tombuildsstuff/giovanni` ([#25551](https://github.com/hashicorp/terraform-provider-azurerm/issues/25551))
* `azurerm_key_vault` - deprecate the `contact` property from v3.x provider and update properties to Computed & Optional ([#25552](https://github.com/hashicorp/terraform-provider-azurerm/issues/25552))
* `azurerm_key_vault_certificate_contacts` - in v4.0 make the `contact` property optional to allow for deletion of contacts from the key vault ([#25552](https://github.com/hashicorp/terraform-provider-azurerm/issues/25552))
* `azurerm_signalr_service` - support for setting the `sku` property to `Premium_P2` ([#25578](https://github.com/hashicorp/terraform-provider-azurerm/issues/25578))
* `azurerm_snapshot` - support for the `network_access_policy` and `public_network_access_enabled` properties ([#25421](https://github.com/hashicorp/terraform-provider-azurerm/issues/25421))
* `azurerm_storage_account` - extend the support level of `(blob|queue|share)_properties` for Storage kind ([#25427](https://github.com/hashicorp/terraform-provider-azurerm/issues/25427))
* `azurerm_storage_blob` - support for the `encryption_scope` property ([#25551](https://github.com/hashicorp/terraform-provider-azurerm/issues/25551))
* `azurerm_storage_container` - support for the `default_encryption_scope` and `encryption_scope_override_enabled` properties ([#25551](https://github.com/hashicorp/terraform-provider-azurerm/issues/25551))
* `azurerm_storage_data_lake_gen2_filesystem` - support for the `default_encryption_scope` property ([#25551](https://github.com/hashicorp/terraform-provider-azurerm/issues/25551))
* `azurerm_subnet` - the `delegation.x.service_delegation.x.name` property now supports `Oracle.Database/networkAttachments` ([#25571](https://github.com/hashicorp/terraform-provider-azurerm/issues/25571))
* `azurerm_web_pubsub` - support setting the `sku` property to `Premium_P2` ([#25578](https://github.com/hashicorp/terraform-provider-azurerm/issues/25578))

BUG FIXES:

* provider: fix an issue where the provider was not correctly configured when using a custom metadata host ([#25546](https://github.com/hashicorp/terraform-provider-azurerm/issues/25546))
* storage: fix a number of potential crashes during plan/apply with resources using the Storage data plane API ([#25525](https://github.com/hashicorp/terraform-provider-azurerm/issues/25525))
* `azurerm_application_insights` - fix issue where the wrong Application ID was set into the property `app_id` ([#25520](https://github.com/hashicorp/terraform-provider-azurerm/issues/25520))
* `azurerm_application_insights_api_key` - add a state migration to re-case static segments of the resource ID  ([#25567](https://github.com/hashicorp/terraform-provider-azurerm/issues/25567))
* `azurerm_container_app_environment_certificate` - the `subject_name` attribute is now correctly populated ([#25516](https://github.com/hashicorp/terraform-provider-azurerm/issues/25516))
* `azurerm_function_app_slot` - will now taint the resource when partially created ([#24520](https://github.com/hashicorp/terraform-provider-azurerm/issues/24520))
* `azurerm_linux_function_app` - will now taint the resource when partially created ([#24520](https://github.com/hashicorp/terraform-provider-azurerm/issues/24520))
* `azurerm_managed_disk` - filtering the Resource SKUs response to reduce the memory overhead, when determining whether a Managed Disk can be online resized or not ([#25549](https://github.com/hashicorp/terraform-provider-azurerm/issues/25549))
* `azurerm_monitor_alert_prometheus_rule_group` - the `severity` property is now set correctly when `0` ([#25408](https://github.com/hashicorp/terraform-provider-azurerm/issues/25408))
* `azurerm_monitor_smart_detector_alert_rule` - normalising the value for `id` within the `action_group` block ([#25559](https://github.com/hashicorp/terraform-provider-azurerm/issues/25559))
* `azurerm_redis_cache_access_policy_assignment` - the `object_id_alias` property now allows usernames ([#25523](https://github.com/hashicorp/terraform-provider-azurerm/issues/25523))
* `azurerm_windows_function_app` - will not taint the resource when partially created ([#24520](https://github.com/hashicorp/terraform-provider-azurerm/issues/24520))
* `azurerm_windows_function_app` - will not taint the resource when partially created ([#24520](https://github.com/hashicorp/terraform-provider-azurerm/issues/24520))

DEPRECATIONS:

* `azurerm_cosmosdb_account` - the `connection_strings` property has been superseded by the primary and secondary connection strings for sql, mongodb and readonly ([#25510](https://github.com/hashicorp/terraform-provider-azurerm/issues/25510))
* `azurerm_cosmosdb_account` - the `enable_free_tier` property has been superseded by `free_tier_enabled` ([#25510](https://github.com/hashicorp/terraform-provider-azurerm/issues/25510))
* `azurerm_cosmosdb_account` - the `enable_multiple_write_locations` property  has been superseded by `multiple_write_locations_enabled` ([#25510](https://github.com/hashicorp/terraform-provider-azurerm/issues/25510))
* `azurerm_cosmosdb_account` - the `enable_automatic_failover` property has been superseded by `automatic_failover_enabled` ([#25510](https://github.com/hashicorp/terraform-provider-azurerm/issues/25510))

## 3.98.0 (April 04, 2024)

FEATURES:

* New Resource: `azurerm_static_web_app_function_app_registration` ([#25331](https://github.com/hashicorp/terraform-provider-azurerm/issues/25331))
* New Resource: `azurerm_system_center_virtual_machine_manager_inventory_items` ([#25110](https://github.com/hashicorp/terraform-provider-azurerm/issues/25110))
* New Resource: `azurerm_workloads_sap_discovery_virtual_instance` ([#24342](https://github.com/hashicorp/terraform-provider-azurerm/issues/24342))
* New Resource: `azurerm_redis_cache_policy` ([#25477](https://github.com/hashicorp/terraform-provider-azurerm/issues/25477))
* New Resource: `azurerm_redis_cache_policy_assignment` ([#25477](https://github.com/hashicorp/terraform-provider-azurerm/issues/25477))

ENHANCEMENTS:

* dependencies: updating to `v0.20240402.1085733` of `github.com/hashicorp/go-azure-sdk` ([#25482](https://github.com/hashicorp/terraform-provider-azurerm/issues/25482))
* dependencies: updating to `v0.67.0` of `github.com/hashicorp/go-azure-helpers` ([#25446](https://github.com/hashicorp/terraform-provider-azurerm/issues/25446))
* dependencies: updating to `v0.25.4` of `github.com/tombuildsstuff/giovanni` ([#25404](https://github.com/hashicorp/terraform-provider-azurerm/issues/25404))
* `alertsmanagement` - updating remaining resources to use `hashicorp/go-azure-sdk` ([#25486](https://github.com/hashicorp/terraform-provider-azurerm/issues/25486))
* `applicationinsights` - updating remaining resources to use `hashicorp/go-azure-sdk` ([#25376](https://github.com/hashicorp/terraform-provider-azurerm/issues/25376))
* `compute` - update to API version `2024-03-01` ([#25436](https://github.com/hashicorp/terraform-provider-azurerm/issues/25436))
* `compute` - update shared image resources and data sources to use `hashicorp/go-azure-sdk` ([#25503](https://github.com/hashicorp/terraform-provider-azurerm/issues/25503))
* `containerinstance` - update to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#25416](https://github.com/hashicorp/terraform-provider-azurerm/issues/25416))
* `maintenance` -  updating to API Version `2023-04-01` ([#25388](https://github.com/hashicorp/terraform-provider-azurerm/issues/25388))
* `recovery_services` - Add `recovery_service` block to the provider that supports `vm_backup_stop_protection_and_retain_data_on_destroy` and `purge_protected_items_from_vault_on_destroy`([#25515](https://github.com/hashicorp/terraform-provider-azurerm/issues/25515))
* `storage` -  the Storage Account cache is now populated using `hashicorp/go-azure-sdk` ([#25437](https://github.com/hashicorp/terraform-provider-azurerm/issues/25437))
* `azurerm_bot_service_azure_bot` - support for the `cmk_key_vault_key_url` property ([#23640](https://github.com/hashicorp/terraform-provider-azurerm/issues/23640))
* `azurerm_capacity_reservation` - update validation for `capacity` ([#25471](https://github.com/hashicorp/terraform-provider-azurerm/issues/25471))
* `azurerm_container_app` - add support for `key_vault_id` and `identity` properties in the `secret` block ([#24773](https://github.com/hashicorp/terraform-provider-azurerm/issues/24773))
* `azurerm_databricks_workspace` - expose `managed_services_cmk_key_vault_id` and `managed_disk_cmk_key_vault_id and key_vault_id` to support cross subscription CMK's. ([#25091](https://github.com/hashicorp/terraform-provider-azurerm/issues/25091))
* `azurerm_databricks_workspace_root_dbfs_customer_managed_key` - expose `key_vault_id` to support cross subscription CMK's. ([#25091](https://github.com/hashicorp/terraform-provider-azurerm/issues/25091))
* `azurerm_managed_hsm_role_*_ids` - use specific resource id to replace generic nested item id ([#25323](https://github.com/hashicorp/terraform-provider-azurerm/issues/25323))
* `azurerm_mssql_database` - add support for `secondary_type` ([#25360](https://github.com/hashicorp/terraform-provider-azurerm/issues/25360))
* `azurerm_monitor_scheduled_query_rules_alert_v2` - support for the `identity` block ([#25365](https://github.com/hashicorp/terraform-provider-azurerm/issues/25365))
* `azurerm_mssql_server_extended_auditing_policy` - support for `audit_actions_and_groups` and `predicate_expression` ([#25425](https://github.com/hashicorp/terraform-provider-azurerm/issues/25425))
* `azurerm_netapp_account` - can now be imported ([#25384](https://github.com/hashicorp/terraform-provider-azurerm/issues/25384))
* `azurerm_netapp_volume` - support for the `kerberos_enabled`, `smb_continuous_availability_enabled`, `kerberos_5_read_only_enabled`, `kerberos_5_read_write_enabled`, `kerberos_5i_read_only_enabled`, `kerberos_5i_read_write_enabled`, `kerberos_5p_read_only_enabled`, and `kerberos_5p_read_write_enabled` properties ([#25385](https://github.com/hashicorp/terraform-provider-azurerm/issues/25385))
* `azurerm_recovery_services_vault` - upgrading to version `2024-01-01` ([#25325](https://github.com/hashicorp/terraform-provider-azurerm/issues/25325))
* `azurerm_stack_hci_cluster` - the `client_id` property is now optional ([#25407](https://github.com/hashicorp/terraform-provider-azurerm/issues/25407))
* `azurerm_storage_encryption_scope` - refactoring to use `hashicorp/go-azure-sdk` rather than `Azure/azure-sdk-for-go` ([#25437](https://github.com/hashicorp/terraform-provider-azurerm/issues/25437))
* `azurerm_mssql_elasticpool` - the `maintenance_configuration_name` property now supports values `SQL_SouthAfricaNorth_DB_1`, `SQL_SouthAfricaNorth_DB_2`, `SQL_WestUS3_DB_1` and `SQL_WestUS3_DB_2` ([#25500](https://github.com/hashicorp/terraform-provider-azurerm/issues/25500))
* `azurerm_lighthouse_assignment` - updating API Version from `2019-06-01` to `2022-10-01`  ([#25473](https://github.com/hashicorp/terraform-provider-azurerm/issues/25473))

BUG FIXES:

* `network` -  updating the `GatewaySubnet` validation to show the Subnet Name when the validation fails ([#25484](https://github.com/hashicorp/terraform-provider-azurerm/issues/25484))
* `azurerm_function_app_hybrid_connection` - fix an issue during creation when `send_key_name` is specified ([#25379](https://github.com/hashicorp/terraform-provider-azurerm/issues/25379))
* `azurerm_linux_web_app_slot` - fix a crash when upgrading the provider to v3.88.0 or later ([#25406](https://github.com/hashicorp/terraform-provider-azurerm/issues/25406))
* `azurerm_mssql_database` - update the behavior of the `enclave_type` field. ([#25508](https://github.com/hashicorp/terraform-provider-azurerm/issues/25508))
* `azurerm_mssql_elasticpool` - update the behavior of the `enclave_type` field. ([#25508](https://github.com/hashicorp/terraform-provider-azurerm/issues/25508))
* `azurerm_network_manager_deployment` - add locking ([#25368](https://github.com/hashicorp/terraform-provider-azurerm/issues/25368))
* `azurerm_resource_group_template_deployment` - changes to `parameters_content` and `template_content` now force `output_content` to be updated in the plan ([#25403](https://github.com/hashicorp/terraform-provider-azurerm/issues/25403))
* `azurerm_storage_blob` - fix a potential crash when the endpoint is unreachable ([#25404](https://github.com/hashicorp/terraform-provider-azurerm/issues/25404))
* `azurerm_storage_container` - fix a potential crash when the endpoint is unreachable ([#25404](https://github.com/hashicorp/terraform-provider-azurerm/issues/25404))
* `azurerm_storage_data_lake_gen2_filesystem` - fix a potential crash when the endpoint is unreachable ([#25404](https://github.com/hashicorp/terraform-provider-azurerm/issues/25404))
* `azurerm_storage_data_lake_gen2_filesystem_path` - fix a potential crash when the endpoint is unreachable ([#25404](https://github.com/hashicorp/terraform-provider-azurerm/issues/25404))
* `azurerm_storage_queue` - fix a potential crash when the endpoint is unreachable ([#25404](https://github.com/hashicorp/terraform-provider-azurerm/issues/25404))
* `azurerm_storage_share` - fix a potential crash when the endpoint is unreachable ([#25404](https://github.com/hashicorp/terraform-provider-azurerm/issues/25404))
* `azurerm_storage_share_directory` - fix a potential crash when the endpoint is unreachable ([#25404](https://github.com/hashicorp/terraform-provider-azurerm/issues/25404))
* `azurerm_storage_share_directory` - resolve an issue where directories might fail to destroy ([#25404](https://github.com/hashicorp/terraform-provider-azurerm/issues/25404))
* `azurerm_storage_share_file` - fix a potential crash when the endpoint is unreachable ([#25404](https://github.com/hashicorp/terraform-provider-azurerm/issues/25404))
* `azurerm_storage_share_file` - fix several bugs with path handling when creating files in subdirectories ([#25404](https://github.com/hashicorp/terraform-provider-azurerm/issues/25404))
* `azurerm_web_app_hybrid_connection` - fix an issue during creation when `send_key_name` is specified ([#25379](https://github.com/hashicorp/terraform-provider-azurerm/issues/25379))
* `azurerm_windows_web_app` - prevent a panic during resource upgrade ([#25509](https://github.com/hashicorp/terraform-provider-azurerm/issues/25509))

## 3.97.1 (March 22, 2024)

ENHANCEMENTS:

* `azurerm_nginx_deployment` - support for the `configuration` block ([#24276](https://github.com/hashicorp/terraform-provider-azurerm/issues/24276))

BUG FIXES:

* `azurerm_data_factory_integration_runtime_self_hosted` - ensure that autorizationh keys are exported ([#25246](https://github.com/hashicorp/terraform-provider-azurerm/issues/25246))
* `azurerm_storage_account` - defaulting the value for `dns_endpoint_type` to `Standard` when it's not returned from the Azure API ([#25367](https://github.com/hashicorp/terraform-provider-azurerm/issues/25367))


## 3.97.0 (March 21, 2024)

BREAKING CHANGES:

* `azurerm_linux_function_app` - `app_settings["WEBSITE_RUN_FROM_PACKAGE"]` must be added to `ignore_changes` for deployments where an external tool modifies the `WEBSITE_RUN_FROM_PACKAGE` property in the `app_settings` block. ([#24848](https://github.com/hashicorp/terraform-provider-azurerm/issues/24848))
* `azurerm_linux_function_app_slot` - `app_settings["WEBSITE_RUN_FROM_PACKAGE"]` must be added to `ignore_changes` for deployments where an external tool modifies the `WEBSITE_RUN_FROM_PACKAGE` property in the `app_settings` block. ([#24848](https://github.com/hashicorp/terraform-provider-azurerm/issues/24848))

FEATURES:

* New Resource: `azurerm_elastic_san_volume` ([#24802](https://github.com/hashicorp/terraform-provider-azurerm/issues/24802))

ENHANCEMENTS:

* dependencies: updating to `v0.25.3` of `github.com/tombuildsstuff/giovanni` ([#25362](https://github.com/hashicorp/terraform-provider-azurerm/issues/25362))
* dependencies: updating to `v0.20240321.1145953` of `github.com/hashicorp/go-azure-sdk` ([#25332](https://github.com/hashicorp/terraform-provider-azurerm/issues/25332))
* dependencies: updating to `v0.25.2` of `github.com/tombuildsstuff/giovanni` ([#25305](https://github.com/hashicorp/terraform-provider-azurerm/issues/25305))
* `azurestackhci`: updating to API Version `2024-01-01` ([#25279](https://github.com/hashicorp/terraform-provider-azurerm/issues/25279))
* `monitor/scheduledqueryrules`: updating to API version `2023-03-15-preview` ([#25350](https://github.com/hashicorp/terraform-provider-azurerm/issues/25350))
* `cosmosdb`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#25166](https://github.com/hashicorp/terraform-provider-azurerm/issues/25166))
* Data Source `azurerm_stack_hci_cluster`: refactoring the association to use `hashicorp/go-azure-sdk` ([#25293](https://github.com/hashicorp/terraform-provider-azurerm/issues/25293))
* `azurerm_app_configuration` - support for Environments other than Azure Public ([#25271](https://github.com/hashicorp/terraform-provider-azurerm/issues/25271))
* `azurerm_automanage_configuration` - refactoring to use `hashicorp/go-azure-sdk` ([#25293](https://github.com/hashicorp/terraform-provider-azurerm/issues/25293))
* `azurerm_container_app_environment` - add support for `Consumption` workload profile ([#25285](https://github.com/hashicorp/terraform-provider-azurerm/issues/25285))
* `azurerm_cosmosdb_postgresql_cluster` - expose list of server names and FQDN in the `servers` block ([#25240](https://github.com/hashicorp/terraform-provider-azurerm/issues/25240))
* `azurerm_data_share` - hyphens are now allowed in the resource's name ([#25242](https://github.com/hashicorp/terraform-provider-azurerm/issues/25242))
* `azurerm_data_factory_integration_runtime_azure_ssis` - support for the `copy_compute_scale` and `pipeline_external_compute_scale` blocks ([#25281](https://github.com/hashicorp/terraform-provider-azurerm/issues/25281))
* `azurerm_healthcare_service` - support for the `identity` and `configuration_export_storage_account_name` properties ([#25193](https://github.com/hashicorp/terraform-provider-azurerm/issues/25193))
* `azurerm_nginx_deployment` - support the `auto_scale_profile` block ([#24950](https://github.com/hashicorp/terraform-provider-azurerm/issues/24950))
* `azurerm_netapp_account_resource` - support for the `kerberos_ad_name`, `kerberos_kdc_ip property`, `enable_aes_encryption`, `local_nfs_users_with_ldap_allowed`, `server_root_ca_certificate`, `ldap_over_tls_enabled`, and `ldap_signing_enabled` properties ([#25340](https://github.com/hashicorp/terraform-provider-azurerm/issues/25340))
* `azurerm_netapp_account_resource` - support for [Support for Azure Netapp Files - AD Site Name #12462] via the `site_name` property ([#25340](https://github.com/hashicorp/terraform-provider-azurerm/issues/25340))
* `azurerm_stack_hci_cluster`: refactoring the association to use `hashicorp/go-azure-sdk` ([#25293](https://github.com/hashicorp/terraform-provider-azurerm/issues/25293))
* `azurerm_storage_account` - support for the `dns_endpoint_type` property ([#22583](https://github.com/hashicorp/terraform-provider-azurerm/issues/22583))
* `azurerm_storage_blob_inventory_policy` -  refactoring to use `hashicorp/go-azure-sdk` ([#25268](https://github.com/hashicorp/terraform-provider-azurerm/issues/25268))
* `azurerm_synapse_spark_pool` - added support for `3.4` ([#25319](https://github.com/hashicorp/terraform-provider-azurerm/issues/25319))

BUG FIXES:

* Data Source: `azurerm_storage_blob` - fix a bug that incorrectly parsed the endpoint in the resource ID ([#25283](https://github.com/hashicorp/terraform-provider-azurerm/issues/25283))
* Data Source: `azurerm_storage_table_entity` - fixing a regression when parsing the table endpoint ([#25307](https://github.com/hashicorp/terraform-provider-azurerm/issues/25307))
* `netapp_account_resource` - correct the `smb_server_name` property validation ([#25340](https://github.com/hashicorp/terraform-provider-azurerm/issues/25340))
* `azurerm_backup_policy_file_share` - prevent a bug when the `include_last_days` property does not work when `days` is empty ([#25280](https://github.com/hashicorp/terraform-provider-azurerm/issues/25280))
* `azurerm_backup_policy_vm` - prevent a bug when the `include_last_days` property does not work when `days` is empty ([#25280](https://github.com/hashicorp/terraform-provider-azurerm/issues/25280))
* `azurerm_container_app_custom_domain` - prevent an issue where the secret was not being passed through (#25196) ([#25251](https://github.com/hashicorp/terraform-provider-azurerm/issues/25251))
* `azurerm_data_protection_backup_instance_kubernetes_cluster` - prevent the protection errosr `ScenarioPluginInvalidWorkflowDataRequest` and `UserErrorKubernetesBackupExtensionUnhealthy` [azurerm_data_protection_backup_instance_kubernetes_cluster is created with message "Fix protection error for the backup instance" and code ScenarioPluginInvalidWorkflowDataRequest #25294] ([#25345](https://github.com/hashicorp/terraform-provider-azurerm/issues/25345))
* `azurerm_purview_account` - will now allow for PurView accounts with missing or disabled eventhubs without keys ([#25301](https://github.com/hashicorp/terraform-provider-azurerm/issues/25301))
* `azurerm_storage_account` - fix a crash when the storage account becomes unavailable whilst reading ([#25332](https://github.com/hashicorp/terraform-provider-azurerm/issues/25332))
* `azurerm_storage_blob` - fixing a regression where blobs within a nested directory wouldn't be parsed correctly ([#25305](https://github.com/hashicorp/terraform-provider-azurerm/issues/25305))
* `azurerm_storage_data_lake_gen2_path` - fixing a bug where there was no timeout available during import ([#25282](https://github.com/hashicorp/terraform-provider-azurerm/issues/25282))
* `azurerm_storage_queue` - fixing a bug where the Table URI was obtained rather than the Queue URI ([#25262](https://github.com/hashicorp/terraform-provider-azurerm/issues/25262))
* `azurerm_subscription` - fixing an issue when creating a subscription alias ([#25181](https://github.com/hashicorp/terraform-provider-azurerm/issues/25181))


## 3.96.0 (March 14, 2024)

ENHANCEMENTS:

* dependencies: updating to `v0.20240314.1083835` of `github.com/hashicorp/go-azure-sdk` ([#25255](https://github.com/hashicorp/terraform-provider-azurerm/issues/25255))
* dependencies: updating to `v0.25.1` of `github.com/tombuildsstuff/giovanni` ([#24798](https://github.com/hashicorp/terraform-provider-azurerm/issues/24798))
* dependencies: updating to `v1.33.0` of `google.golang.org/protobuf` ([#25243](https://github.com/hashicorp/terraform-provider-azurerm/issues/25243))
* `storage`: updating the data plane resources to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24798](https://github.com/hashicorp/terraform-provider-azurerm/issues/24798))
* Data Source: `azurerm_storage_table_entities` - support for AAD authentication ([#24798](https://github.com/hashicorp/terraform-provider-azurerm/issues/24798))
* Data Source: `azurerm_storage_table_entity` - support for AAD authentication ([#24798](https://github.com/hashicorp/terraform-provider-azurerm/issues/24798))
* `azurerm_kusto_cluster` - support `None` pattern for the `virtual_network_configuration` block ([#24733](https://github.com/hashicorp/terraform-provider-azurerm/issues/24733))
* `azurerm_linux_function_app` - support for the Node `20` runtime ([#24073](https://github.com/hashicorp/terraform-provider-azurerm/issues/24073))
* `azurerm_linux_function_app_slot` - support for the Node `20` runtime ([#24073](https://github.com/hashicorp/terraform-provider-azurerm/issues/24073))
* `azurerm_stack_hci_cluster` - support the `identity`, `cloud_id`, `service_endpoint` and `resource_provider_object_id` properties ([#25031](https://github.com/hashicorp/terraform-provider-azurerm/issues/25031))
* `azurerm_storage_share_file` - support for AAD authentication ([#24798](https://github.com/hashicorp/terraform-provider-azurerm/issues/24798))
* `azurerm_storage_share_directory` - support for AAD authentication, deprecate `share_name` and `storage_account_name` in favor of `storage_share_id` ([#24798](https://github.com/hashicorp/terraform-provider-azurerm/issues/24798))
* `azurerm_storage_table_entity` - support for AAD authentication, deprecate `share_name` and `storage_account_name` in favor of `storage_table_id` ([#24798](https://github.com/hashicorp/terraform-provider-azurerm/issues/24798))
* `azurerm_storage_table_entity` - support for AAD authentication ([#24798](https://github.com/hashicorp/terraform-provider-azurerm/issues/24798))
* `azurerm_windows_function_app` - support for the Node `20` runtime ([#24073](https://github.com/hashicorp/terraform-provider-azurerm/issues/24073))
* `azurerm_windows_function_app_slot` - support for the Node `20` runtime ([#24073](https://github.com/hashicorp/terraform-provider-azurerm/issues/24073))
* `azurerm_windows_web_app` - support for the Node `20` runtime ([#24073](https://github.com/hashicorp/terraform-provider-azurerm/issues/24073))
* `azurerm_windows_web_app_slot` - support for the Node `20` runtime ([#24073](https://github.com/hashicorp/terraform-provider-azurerm/issues/24073))

BUG FIXES:

* `azurerm_container_app_custom_domain` - fix resource ID parsing bug preventing import ([#25192](https://github.com/hashicorp/terraform-provider-azurerm/issues/25192))
* `azurerm_windows_web_app` - fix incorrect warning message when checking name availability ([#25214](https://github.com/hashicorp/terraform-provider-azurerm/issues/25214))
* `azurerm_virtual_machine_run_command` - prevent a bug during updates ([#25186](https://github.com/hashicorp/terraform-provider-azurerm/issues/25186))
* Data Source: `azurerm_storage_table_entities`  - Fix `items.x.properties` truncating to one entry ([#25211](https://github.com/hashicorp/terraform-provider-azurerm/issues/25211))

## 3.95.0 (March 08, 2024)

FEATURES:

* New Resource: `azurerm_container_app_custom_domain` ([#24421](https://github.com/hashicorp/terraform-provider-azurerm/issues/24421))
* New Resource: `azurerm_data_protection_backup_instance_kubernetes_cluster` ([#24940](https://github.com/hashicorp/terraform-provider-azurerm/issues/24940))
* New Resource: `azurerm_static_web_app` ([#25117](https://github.com/hashicorp/terraform-provider-azurerm/issues/25117))
* New resource: `azurerm_static_web_app_custom_domain` ([#25117](https://github.com/hashicorp/terraform-provider-azurerm/issues/25117))
* New resource: `azurerm_system_center_virtual_machine_manager_availability_set` ([#24975](https://github.com/hashicorp/terraform-provider-azurerm/issues/24975))
* New Resource: `azurerm_workloads_sap_three_tier_virtual_instance` ([#24384](https://github.com/hashicorp/terraform-provider-azurerm/issues/24384))
* New Resource: `azurerm_workloads_sap_single_node_virtual_instance` ([#24331](https://github.com/hashicorp/terraform-provider-azurerm/issues/24331))

ENHANCEMENTS:

* `dependencies`: updating to v0.20240229.1102109 of `github.com/hashicorp/go-azure-sdk` ([#25102](https://github.com/hashicorp/terraform-provider-azurerm/issues/25102))
* `monitor`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` [GH-#25102]
* `network`: updating to API Version `2023-09-01` ([#25095](https://github.com/hashicorp/terraform-provider-azurerm/issues/25095))
* `azurerm_data_factory_integration_runtime_managed` - support for the `credential_name` property ([#25033](https://github.com/hashicorp/terraform-provider-azurerm/issues/25033))
* `azurerm_linux_function_app` - support for the `description` property in the `ip_restriction` block ([#24527](https://github.com/hashicorp/terraform-provider-azurerm/issues/24527))
* `azurerm_linux_function_app` - support for the `ip_restriction_default_action` and `scm_ip_restriction_default_action` properties ([#25131](https://github.com/hashicorp/terraform-provider-azurerm/issues/25131))
* `azurerm_linux_function_app_slot` - support for the `description` property in the `ip_restriction` block ([#24527](https://github.com/hashicorp/terraform-provider-azurerm/issues/24527))
* `azurerm_linux_function_app_slot` - support for the `ip_restriction_default_action` and `scm_ip_restriction_default_action` properties ([#25131](https://github.com/hashicorp/terraform-provider-azurerm/issues/25131))
* `azurerm_linux_web_app` - support for the `description` property in the `ip_restriction` block ([#24527](https://github.com/hashicorp/terraform-provider-azurerm/issues/24527))
* `azurerm_linux_web_app` - support for the `ip_restriction_default_action` and `scm_ip_restriction_default_action` properties ([#25131](https://github.com/hashicorp/terraform-provider-azurerm/issues/25131))
* `azurerm_linux_web_app_slot` - support for the `description` property in the `ip_restriction` block ([#24527](https://github.com/hashicorp/terraform-provider-azurerm/issues/24527))
* `azurerm_linux_web_app_slot` - support for the `ip_restriction_default_action` and `scm_ip_restriction_default_action` properties ([#25131](https://github.com/hashicorp/terraform-provider-azurerm/issues/25131))
* `azurerm_mysql_flexible_server` - setting the `storage.size_gb` property to a smaller value now forces a new resource to be created ([#25074](https://github.com/hashicorp/terraform-provider-azurerm/issues/25074))
* `azurerm_orbital_contact_profile` - changing the `channels` property no longer creates a new resource ([#25129](https://github.com/hashicorp/terraform-provider-azurerm/issues/25129))
* `azurerm_private_dns_resolver_inbound_endpoint` - the `private_ip_address` property is no longer required when `private_ip_allocation_method` is `Dynamic` ([#25035](https://github.com/hashicorp/terraform-provider-azurerm/issues/25035))
* `stream_analytics_output_blob` - support for the `blob_write_mode` property ([#25127](https://github.com/hashicorp/terraform-provider-azurerm/issues/25127))
* `azurerm_windows_function_app` - support for the `description` property in the `ip_restriction` block ([#24527](https://github.com/hashicorp/terraform-provider-azurerm/issues/24527))
* `azurerm_windows_function_app` - support for the `ip_restriction_default_action` and `scm_ip_restriction_default_action` properties ([#25131](https://github.com/hashicorp/terraform-provider-azurerm/issues/25131))
* `azurerm_windows_function_app_slot` - support for the `description` property in the `ip_restriction` block ([#24527](https://github.com/hashicorp/terraform-provider-azurerm/issues/24527))
* `azurerm_windows_function_app_slot` - support for the `ip_restriction_default_action` and `scm_ip_restriction_default_action` properties ([#25131](https://github.com/hashicorp/terraform-provider-azurerm/issues/25131))
* `azurerm_windows_web_app` - support for the `description` property in the `ip_restriction` block ([#24527](https://github.com/hashicorp/terraform-provider-azurerm/issues/24527))
* `azurerm_windows_web_app` - support for the `ip_restriction_default_action` and `scm_ip_restriction_default_action` properties ([#25131](https://github.com/hashicorp/terraform-provider-azurerm/issues/25131))
* `azurerm_windows_web_app_slot` - support for the `description` property in the `ip_restriction` block ([#24527](https://github.com/hashicorp/terraform-provider-azurerm/issues/24527))
* `azurerm_windows_web_app_slot` - support for the `ip_restriction_default_action` and `scm_ip_restriction_default_action` properties ([#25131](https://github.com/hashicorp/terraform-provider-azurerm/issues/25131))

BUG FIXES:

* Data Source: `azurerm_function_app_host_keys` - correctly set `event_grid_extension_key` by searching for the renamed property in the API response ([#25108](https://github.com/hashicorp/terraform-provider-azurerm/issues/25108))
* `azurerm_app_service_public_certificate` - fix issue where certificate information was not being set correctly in the read ([#24943](https://github.com/hashicorp/terraform-provider-azurerm/issues/24943))
* `azurerm_container_registry` - prevent recreation of the resource when the `georeplication.tags` are updated ([#24994](https://github.com/hashicorp/terraform-provider-azurerm/issues/24994))
* `azurerm_firewall_policy_rule_collection_group` - fix issue where the client subscription ID was used to construct the `firewall_policy_id` ([#25145](https://github.com/hashicorp/terraform-provider-azurerm/issues/25145))
* `azurerm_function_app_hybrid_connection` - fix issue where `SendKeyValue` was not populated in the API payload ([#23761](https://github.com/hashicorp/terraform-provider-azurerm/issues/23761))
* `azurerm_orbital_contact_profile` - fix creation of the resource when `event_hub_uri` is not specified ([#25128](https://github.com/hashicorp/terraform-provider-azurerm/issues/25128))
* `azurerm_recovery_services_vault` - prevent a panic when `immutability` is updated ([#25132](https://github.com/hashicorp/terraform-provider-azurerm/issues/25132))
* `azurerm_storage_account` - fix issue where the queue encryption key type was set as the table encryption key type ([#25046](https://github.com/hashicorp/terraform-provider-azurerm/issues/25046))
* `azurerm_web_app_hybrid_connection` - fix issue where `SendKeyValue` was not populated in the API payload ([#23761](https://github.com/hashicorp/terraform-provider-azurerm/issues/23761))
* `azurerm_mssql_database` - fix incorrect error due to typo when using `restore_long_term_retention_backup_id` ([#25180](https://github.com/hashicorp/terraform-provider-azurerm/issues/25180))

DEPRECATIONS:

* Deprecated Resource: `azurerm_static_site` ([#25117](https://github.com/hashicorp/terraform-provider-azurerm/issues/25117))
* Deprecated Resource: `azurerm_static_site_custom_domain` ([#25117](https://github.com/hashicorp/terraform-provider-azurerm/issues/25117))
* `azurerm_kubernetes_fleet_manager` - the `hub_profile` property has been deprecated ([#25010](https://github.com/hashicorp/terraform-provider-azurerm/issues/25010))

## 3.94.0 (February 29, 2024)

FEATURES:

* **New Resource**: `azurerm_kubernetes_fleet_update_run` ([#24813](https://github.com/hashicorp/terraform-provider-azurerm/issues/24813))

ENHANCEMENTS:

* dependencies: updating to `v0.20240228.1142829` of `github.com/hashicorp/go-azure-sdk` ([#25081](https://github.com/hashicorp/terraform-provider-azurerm/issues/25081))
* `servicefabric`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#25002](https://github.com/hashicorp/terraform-provider-azurerm/issues/25002))
* `springcloud`: updating to API Version `2024-01-01-preview` ([#24937](https://github.com/hashicorp/terraform-provider-azurerm/issues/24937))
* `securitycenter`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#25081](https://github.com/hashicorp/terraform-provider-azurerm/issues/25081))
* Data Source: `azurerm_storage_table_entities` - support for `select` ([#24987](https://github.com/hashicorp/terraform-provider-azurerm/issues/24987))
* Data Source: `azurerm_netapp_volume` - support for the `smb_access_based_enumeration` and `smb_non_browsable` properties ([#24514](https://github.com/hashicorp/terraform-provider-azurerm/issues/24514))
* `azurerm_cosmosdb_account` - add support for the `minimal_tls_version` property ([#24966](https://github.com/hashicorp/terraform-provider-azurerm/issues/24966))
* `azurerm_federated_identity_credential` - the federated credentials can now be changed without creating a new resource ([#25003](https://github.com/hashicorp/terraform-provider-azurerm/issues/25003))
* `azurerm_kubernetes_cluster` - support for the `current_kubernetes_version` property ([#25079](https://github.com/hashicorp/terraform-provider-azurerm/issues/25079))
* `azurerm_kubernetes_cluster` - private DNS is now allowed for the `web_app_routing` property ([#25038](https://github.com/hashicorp/terraform-provider-azurerm/issues/25038))
* `azurerm_kubernetes_cluster` - migration between different `outbound_type`s is now allowed ([#25021](https://github.com/hashicorp/terraform-provider-azurerm/issues/25021))
* `azurerm_mssql_database` - support for the `recovery_point_id` and `restore_long_term_retention_backup_id` properties ([#24904](https://github.com/hashicorp/terraform-provider-azurerm/issues/24904))
* `azurerm_linux_virtual_machine` - support for the `automatic_upgrade_enabled`, `disk_controller_type`, `os_image_notification`, `treat_failure_as_deployment_failure_enabled`, and `vm_agent_platform_updates_enabled`properties ([#23394](https://github.com/hashicorp/terraform-provider-azurerm/issues/23394))
* `azurerm_nginx_deployment` - support for the `automatic_upgrade_channel` property ([#24867](https://github.com/hashicorp/terraform-provider-azurerm/issues/24867))
* `azurerm_netapp_volume` - support for the `smb_access_based_enumeration` and `smb_non_browsable` properties ([#24514](https://github.com/hashicorp/terraform-provider-azurerm/issues/24514))
* `azurerm_netapp_pool` - support for the `encryption_type` property ([#24993](https://github.com/hashicorp/terraform-provider-azurerm/issues/24993))
* `azurerm_role_definition` - upgrade to the API version `2022-05-01-preview` ([#25008](https://github.com/hashicorp/terraform-provider-azurerm/issues/25008))
* `azurerm_redis_cache` - allow AAD auth for all SKUs ([#25006](https://github.com/hashicorp/terraform-provider-azurerm/issues/25006))
* `azurerm_sql_managed_instance` - support for the `zone_redundant_enabled` property ([#25089](https://github.com/hashicorp/terraform-provider-azurerm/issues/25089))
* `azurerm_spring_cloud_gateway` - support for the `application_performance_monitoring_ids` property ([#24919](https://github.com/hashicorp/terraform-provider-azurerm/issues/24919))
* `azurerm_spring_cloud_configuration_service` - support for the `refresh_interval_in_seconds` property ([#25009](https://github.com/hashicorp/terraform-provider-azurerm/issues/25009))
* `azurerm_synapse_workspace` - support for using the `user_assigned_identity_id` property within the `customer_managed_key` block ([#25027](https://github.com/hashicorp/terraform-provider-azurerm/issues/25027))
* `azurerm_windows_virtual_machine` - support for the `automatic_upgrade_enabled`, `disk_controller_type`, `os_image_notification`, `treat_failure_as_deployment_failure_enabled`, and `vm_agent_platform_updates_enabled`properties ([#23394](https://github.com/hashicorp/terraform-provider-azurerm/issues/23394))

BUG FIXES:

* `azurerm_api_management_notification_recipient_email` - fixing an issue where response pages weren't iterated over correctly ([#25055](https://github.com/hashicorp/terraform-provider-azurerm/issues/25055))
* `azurerm_api_management_notification_recipient_user` - fixing an issue where response pages weren't iterated over correctly ([#25055](https://github.com/hashicorp/terraform-provider-azurerm/issues/25055))
* `azurerm_batch_pool` - fix setting the `extension.settings_json` property ([#24976](https://github.com/hashicorp/terraform-provider-azurerm/issues/24976))
* `azurerm_key_vault_key` - `expiration_date` can be updated if newer date is ahead ([#25000](https://github.com/hashicorp/terraform-provider-azurerm/issues/25000))
* `azurerm_pim_active_role_assignment` - fix an isue where the resource would disappear or fail to import after 45 days ([#24524](https://github.com/hashicorp/terraform-provider-azurerm/issues/24524))
* `azurerm_pim_eligible_role_assignment` - fix an isue where the resource would disappear or fail to import after 45 days ([#24524](https://github.com/hashicorp/terraform-provider-azurerm/issues/24524))
* `azurerm_recovery_services_vault` - validate that `use_system_assigned_identity` and `user_assigned_identity_id` cannot be set at the same time ([#24091](https://github.com/hashicorp/terraform-provider-azurerm/issues/24091))
* `azurerm_recovery_vaults` will now create properly with `SystemAssigned,UserAssigned` identity ([#24978](https://github.com/hashicorp/terraform-provider-azurerm/issues/24978))
* `azurerm_subscription` - fixing an issue where response pages weren't iterated over correctly ([#25055](https://github.com/hashicorp/terraform-provider-azurerm/issues/25055))

## 3.93.0 (February 22, 2024)

FEATURES:

* **New Data Source**: `azurerm_express_route_circuit_peering` ([#24971](https://github.com/hashicorp/terraform-provider-azurerm/issues/24971))
* **New Data Source**: `azurerm_storage_table_entities` ([#24973](https://github.com/hashicorp/terraform-provider-azurerm/issues/24973))
* **New Resource**: `azurerm_dev_center_catalog` ([#24833](https://github.com/hashicorp/terraform-provider-azurerm/issues/24833))
* **New Resource**: `azurerm_system_center_virtual_machine_manager_server` ([#24278](https://github.com/hashicorp/terraform-provider-azurerm/issues/24278))

BUG FIXES:

* `azurerm_key_vault` - conditionally polling the Data Plane endpoint when `public_network_access_enabled` is set to false ([#23823](https://github.com/hashicorp/terraform-provider-azurerm/issues/23823))
* `azurerm_storage_account` - allow the `identity.type` property to be `SystemAssigned, UserAssigned`  when using a Customer Managed Key ([#24923](https://github.com/hashicorp/terraform-provider-azurerm/issues/24923))
* `azurerm_automation_account` - prevent the `identity.identity_ids` User Assigned identity being set when not specified in config ([#24977](https://github.com/hashicorp/terraform-provider-azurerm/issues/24977))

ENHANCEMENTS:

* dependencies: updating to `v0.20240221.1170458` of `hashicorp/go-azure-sdk` ([#24967](https://github.com/hashicorp/terraform-provider-azurerm/issues/24967))
* dependencies: refactor `azurerm_spring_cloud_configuration_service` to use `go-azure-sdk` ([#24918](https://github.com/hashicorp/terraform-provider-azurerm/issues/24918))
* provider: support or the feature flag `virtual_machine_scale_set.reimage_on_manual_upgrade` ([#22975](https://github.com/hashicorp/terraform-provider-azurerm/issues/22975))
* `sentinel`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24962](https://github.com/hashicorp/terraform-provider-azurerm/issues/24962))
* `sqlvirtualmachines`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24912](https://github.com/hashicorp/terraform-provider-azurerm/issues/24912))
* `nginx` : updating to use `2024-01-01-preview` ([#24868](https://github.com/hashicorp/terraform-provider-azurerm/issues/24868))
* `azurerm_cosmosdb_account` - support for the `backup.tier` property ([#24595](https://github.com/hashicorp/terraform-provider-azurerm/issues/24595))
* `azurerm_linux_virtual_machine` - the `virtual_machine_scale_set_id` proeprty can now be changed without creating a new resource ([#24768](https://github.com/hashicorp/terraform-provider-azurerm/issues/24768))
* `azurerm_machine_learning_workspace` - support for the `managed_network.isolation_mode` property ([#24951](https://github.com/hashicorp/terraform-provider-azurerm/issues/24951))
* `azurerm_private_dns_resolver_inbound_endpoint` - support the `static` value for the `private_ip_allocation_method` property ([#24952](https://github.com/hashicorp/terraform-provider-azurerm/issues/24952))
* `azurerm_postgresql_flexible_server` - expose the `storage_tier` field ([#24892](https://github.com/hashicorp/terraform-provider-azurerm/issues/24892))
* `azurerm_redis_cache` - support for the `preferred_data_persistence_auth_method` property ([#24370](https://github.com/hashicorp/terraform-provider-azurerm/issues/24370))
* `azurerm_servicebus_namespace` - support for the `premium_messaging_partitions` property ([#24676](https://github.com/hashicorp/terraform-provider-azurerm/issues/24676))
* `azurerm_windows_virtual_machine` - the `virtual_machine_scale_set_id` proeprty can now be changed without creating a new resource ([#24768](https://github.com/hashicorp/terraform-provider-azurerm/issues/24768))

BUG FIXES:

* `azurerm_cognitive_deployment` - the `version_upgrade_option` property can not be updated without creating a new resource ([#24922](https://github.com/hashicorp/terraform-provider-azurerm/issues/24922))
* `azurerm_data_protection_backup_vault` - support or the `soft_delete` and `retention_duration_in_days` properties ([#24775](https://github.com/hashicorp/terraform-provider-azurerm/issues/24775))
* `azurerm_data_factory_pipeline` - correctly handle incorrect header values ([#24921](https://github.com/hashicorp/terraform-provider-azurerm/issues/24921))
* `azurerm_kusto_cluster` - `optimized_auto_scale` is now updated after `sku` has been updated ([#24906](https://github.com/hashicorp/terraform-provider-azurerm/issues/24906))
* `azurerm_key_vault_certificate` - will now only update the `lifetime_action` of the certificate block unless otherwise required ([#24755](https://github.com/hashicorp/terraform-provider-azurerm/issues/24755))
* `azurerm_linux_virtual_machine_scale_set` - correctly include `public_ip_prefix_id` during updates ([#24939](https://github.com/hashicorp/terraform-provider-azurerm/issues/24939))
* `azurerm_postgresql_flexible_server` - the `customer_managed_key.key_vault_key_id` property is now required ([#24981](https://github.com/hashicorp/terraform-provider-azurerm/issues/24981))
* `azurerm_nginx_deployment` - changing the `sku` property now creates a new resource ([#24905](https://github.com/hashicorp/terraform-provider-azurerm/issues/24905))
* `azurerm_orchestrated_virtual_machine_scale_set` - the `disk_size_gb` and `lun` parameters of `data_disks` are optional now ([#24944](https://github.com/hashicorp/terraform-provider-azurerm/issues/24944))
* `azurerm_storage_account` - change order of API calls to be GET-then-PUT ratehr then PATCHES ([#23935](https://github.com/hashicorp/terraform-provider-azurerm/issues/23935))
* `azurerm_storage_account` - improve the validation around the `immutability_policy` being used with `blob_properties` ([#24938](https://github.com/hashicorp/terraform-provider-azurerm/issues/24938))
* `azurerm_security_center_setting` - prevent a bug when name is `SENTINEL` ([#24497](https://github.com/hashicorp/terraform-provider-azurerm/issues/24497))
* `azurerm_windows_virtual_machine_scale_set` - correctly include `public_ip_prefix_id` during updates ([#24939](https://github.com/hashicorp/terraform-provider-azurerm/issues/24939))



## 3.92.0 (February 15, 2024)

FEATURES:

* **New Data Source**: `azurerm_virtual_desktop_application_group` ([#24771](https://github.com/hashicorp/terraform-provider-azurerm/issues/24771))

ENHANCEMENTS:

* provider: support for the feature flag `postgresql_flexible_server.restart_server_on_configuration_value_change property` ([#23811](https://github.com/hashicorp/terraform-provider-azurerm/issues/23811))
* dependencies: updating to v0.20240214.1142753 of `github.com/hashicorp/go-azure-sdk` ([#24889](https://github.com/hashicorp/terraform-provider-azurerm/issues/24889))
* `automation`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24858](https://github.com/hashicorp/terraform-provider-azurerm/issues/24858))
* `maintenance`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24819](https://github.com/hashicorp/terraform-provider-azurerm/issues/24819))
* `containerapps`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24862](https://github.com/hashicorp/terraform-provider-azurerm/issues/24862))
* `containerservices`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24872](https://github.com/hashicorp/terraform-provider-azurerm/issues/24872))
* `timeseriesinsights`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24889](https://github.com/hashicorp/terraform-provider-azurerm/issues/24889))
* `azurerm_container_app_environment`: support for the `infrastructure_resource_group_name` property ([#24361](https://github.com/hashicorp/terraform-provider-azurerm/issues/24361))
* `azurerm_cost_anomaly_alert` - support for the `subscription_id` property ([#24258](https://github.com/hashicorp/terraform-provider-azurerm/issues/24258))
* `azurerm_cosmosdb_account` - add default values for the `consistency_policy` code block ([#24830](https://github.com/hashicorp/terraform-provider-azurerm/issues/24830))
* `azurerm_dashboard_grafana` - support for the `smtp` block ([#24717](https://github.com/hashicorp/terraform-provider-azurerm/issues/24717))
* `azurerm_key_vault_certificates` - support for the `tags` property ([#24857](https://github.com/hashicorp/terraform-provider-azurerm/issues/24857))
* `azurerm_key_vault_secrets` - support for the `tags` property ([#24857](https://github.com/hashicorp/terraform-provider-azurerm/issues/24857))
* `azurerm_orchestrated_virtual_machine_scale_set` - support for the `additional_unattend_content` block ([#24292](https://github.com/hashicorp/terraform-provider-azurerm/issues/24292))
* `azurerm_virtual_desktop_host_pool` - support for the `vm_template` property ([#24369](https://github.com/hashicorp/terraform-provider-azurerm/issues/24369))

BUG FIXES:

* `azurerm_container_app_environment`: avoid unwanted changes when updating and using `log_analytics_workspace_id` ([#24303](https://github.com/hashicorp/terraform-provider-azurerm/issues/24303))
* `azurerm_cosmosdb_account` - fixed regression in the `backup` code block ([#24830](https://github.com/hashicorp/terraform-provider-azurerm/issues/24830))
* `azurerm_data_factory` - allow the `git_url` property to be blank/empty ([#24879](https://github.com/hashicorp/terraform-provider-azurerm/issues/24879))
* `azurerm_linux_web_app_slot` - the `worker_count` property now works correctly in the `site_config` block ([#24515](https://github.com/hashicorp/terraform-provider-azurerm/issues/24515))
* `azurerm_linux_web_app` - support `off` for the `file_system_level` property ([#24877](https://github.com/hashicorp/terraform-provider-azurerm/issues/24877))
* `azurerm_linux_web_app_slot` - support `off` for the `file_system_level` property ([#24877](https://github.com/hashicorp/terraform-provider-azurerm/issues/24877))
* `azurerm_private_endpoint` - fixing an issue where updating the Private Endpoint would remove any Application Security Group Association ([#24846](https://github.com/hashicorp/terraform-provider-azurerm/issues/24846))
* `azurerm_search_service` - fixed the update function to adjust for changed API behaviour ([#24837](https://github.com/hashicorp/terraform-provider-azurerm/issues/24837))
* `azurerm_search_service` - fixed the update function to adjust for changed API behaviour ([#24903](https://github.com/hashicorp/terraform-provider-azurerm/issues/24903))
* `azurerm_windows_web_app` - support `off` for the `file_system_level` property ([#24877](https://github.com/hashicorp/terraform-provider-azurerm/issues/24877))
* `azurerm_windows_web_app_slot` - support `off` for the `file_system_level` property ([#24877](https://github.com/hashicorp/terraform-provider-azurerm/issues/24877))

## 3.91.0 (February 08, 2024)

FEATURES:

* **New Data Source**: `azurerm_databricks_access_connector` ([#24769](https://github.com/hashicorp/terraform-provider-azurerm/issues/24769))
* **New Resource**: `azurerm_data_protection_backup_policy_kubernetes_cluster` ([#24718](https://github.com/hashicorp/terraform-provider-azurerm/issues/24718))
* **New Resource**: `azurerm_chaos_studio_experiment` ([#24779](https://github.com/hashicorp/terraform-provider-azurerm/issues/24779))
* **New Resource**: `azurerm_chaos_studio_capability` ([#24779](https://github.com/hashicorp/terraform-provider-azurerm/issues/24779))
* **New Resource**: `azurerm_dev_center_gallery` ([#23760](https://github.com/hashicorp/terraform-provider-azurerm/issues/23760))
* **New Resource:** `azurerm_kubernetes_fleet_member` ([#24792](https://github.com/hashicorp/terraform-provider-azurerm/issues/24792))
* **New Resource:** `azurerm_iotcentral_organization` ([#23132](https://github.com/hashicorp/terraform-provider-azurerm/issues/23132))
* **New Resource:** `azurerm_spring_cloud_app_dynamics_application_performance_monitoring` ([#24750](https://github.com/hashicorp/terraform-provider-azurerm/issues/24750))

ENHANCEMENTS:

* dependencies: updating to `v0.20240208.1095436` of `github.com/hashicorp/go-azure-sdk/resource-manager` ([#24819](https://github.com/hashicorp/terraform-provider-azurerm/issues/24819))
* dependencies: updating to `v0.20240208.1095436` of `github.com/hashicorp/go-azure-sdk/sdk` ([#24819](https://github.com/hashicorp/terraform-provider-azurerm/issues/24819))
* dependencies: refactor `azurerm_app_service_environment_v3` to use `go-azure-sdk` ([#24760](https://github.com/hashicorp/terraform-provider-azurerm/issues/24760))
* dependencies: refactor `azurerm_role_definition` to use `go-azure-sdk` ([#24266](https://github.com/hashicorp/terraform-provider-azurerm/issues/24266))
* `managedhsm`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24761](https://github.com/hashicorp/terraform-provider-azurerm/issues/24761))
* `hdinsight`: updating to API Version `2023-07-01` ([#24761](https://github.com/hashicorp/terraform-provider-azurerm/issues/24761))
* `streamanalytics`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24819](https://github.com/hashicorp/terraform-provider-azurerm/issues/24819))
* `azurerm_app_service_environment_v3` - support for the `remote_debugging_enabled` property ([#24760](https://github.com/hashicorp/terraform-provider-azurerm/issues/24760))
* `azurerm_storage_account` - support for the `local_user_enabled` property ([#24800](https://github.com/hashicorp/terraform-provider-azurerm/issues/24800))
* `azurerm_log_analytics_workspace_table` - support for the `total_retention_in_days` property ([#24513](https://github.com/hashicorp/terraform-provider-azurerm/issues/24513))
* `azurerm_maching_learning_workspace` - support for the `feature_store` and `kind` properties ([#24716](https://github.com/hashicorp/terraform-provider-azurerm/issues/24716))
* `azurerm_traffic_manager_azure_endpoint` - support for the `always_serve_enabled` property ([#24573](https://github.com/hashicorp/terraform-provider-azurerm/issues/24573))
* `azurerm_traffic_manager_external_endpoint` - support for the `always_serve_enabled` property ([#24573](https://github.com/hashicorp/terraform-provider-azurerm/issues/24573))

BUG FIXES:

* `azurerm_api_management` -  the `virtual_network_configuration` property now updates correctly outside of `virtual_network_type` ([#24569](https://github.com/hashicorp/terraform-provider-azurerm/issues/24569))

## 3.90.0 (February 01, 2024)

UPGRADE NOTES:

* provider - The provider will now automatically register the `AppConfiguration`, `DataFactory`, and `SignalRService` Resource Providers. When running Terraform with limited permissions, note that you [must disable automatic Resource Provider Registration and ensure that any Resource Providers Terraform requires are registered]([XXX](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#skip_provider_registration)). ([#24645](https://github.com/hashicorp/terraform-provider-azurerm/issues/24645))
  
FEATURES:

* **New Data Source**: `azurerm_nginx_configuration` ([#24642](https://github.com/hashicorp/terraform-provider-azurerm/issues/24642))
* **New Data Source**: `azurerm_virtual_desktop_workspace` ([#24732](https://github.com/hashicorp/terraform-provider-azurerm/issues/24732))
* **New Resource**: `azurerm_kubernetes_fleet_update_strategy` ([#24328](https://github.com/hashicorp/terraform-provider-azurerm/issues/24328))
* **New Resource**: `azurerm_site_recovery_vmware_replicated_vm` ([#22477](https://github.com/hashicorp/terraform-provider-azurerm/issues/22477))
* **New Resource**: `azurerm_spring_cloud_new_relic_application_performance_monitoring` ([#24699](https://github.com/hashicorp/terraform-provider-azurerm/issues/24699))

ENHANCEMENTS:

* provider: registering the Resource Provider `Microsoft.AppConfiguration` ([#24645](https://github.com/hashicorp/terraform-provider-azurerm/issues/24645))
* provider: registering the Resource Provider `Microsoft.DataFactory` ([#24645](https://github.com/hashicorp/terraform-provider-azurerm/issues/24645))
* provider: registering the Resource Provider `Microsoft.SignalRService` ([#24645](https://github.com/hashicorp/terraform-provider-azurerm/issues/24645))
* provider: the Provider is now built using Go 1.21.6 ([#24653](https://github.com/hashicorp/terraform-provider-azurerm/issues/24653))
* dependencies: the dependency `github.com/hashicorp/go-azure-sdk` has been split into multiple Go Modules - and as such will be referred to by those paths going forwards ([#24636](https://github.com/hashicorp/terraform-provider-azurerm/issues/24636))
* dependencies: updating to ``v0.20240201.1064937` of `github.com/hashicorp/go-azure-sdk/resource-manager` ([#24738](https://github.com/hashicorp/terraform-provider-azurerm/issues/24738))
* dependencies: updating to `v0.20240201.1064937` of `github.com/hashicorp/go-azure-sdk/sdk` ([#24738](https://github.com/hashicorp/terraform-provider-azurerm/issues/24738))
* `appservice`:  update to `go-azure-sdk` and API version `2023-01-01` ([#24688](https://github.com/hashicorp/terraform-provider-azurerm/issues/24688))
* `datafactory`: updating to use `tombuildsstuff/kermit` ([#24675](https://github.com/hashicorp/terraform-provider-azurerm/issues/24675))
* `hdinsight`: refactoring to use `github.com/hashicorp/go-azure-sdk/resource-manager` ([#24011](https://github.com/hashicorp/terraform-provider-azurerm/issues/24011))
* `hdinsight`: updating to API Version `2021-06-01` ([#24011](https://github.com/hashicorp/terraform-provider-azurerm/issues/24011))
* `loadbalancer`: updating to use `hashicorp/go-azure-sdk` ([#24291](https://github.com/hashicorp/terraform-provider-azurerm/issues/24291))
* `nginx`: updating to API Version `2023-09-01` ([#24640](https://github.com/hashicorp/terraform-provider-azurerm/issues/24640))
* `servicefabricmanagedcluster`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24654](https://github.com/hashicorp/terraform-provider-azurerm/issues/24654))
* `springcloud`: updating to use API Version `2023-11-01-preview` ([#24690](https://github.com/hashicorp/terraform-provider-azurerm/issues/24690))
* `subscriptions`: refactoring to use `hashicorp/go-azure-sdk` ([#24663](https://github.com/hashicorp/terraform-provider-azurerm/issues/24663))
* Data Source: `azurerm_stream_analytics_job` - support for User Assigned Identities ([#24738](https://github.com/hashicorp/terraform-provider-azurerm/issues/24738))
* `azurerm_cosmosdb_account` - support for the `gremlin_database` and `tables_to_restore` properties ([#24627](https://github.com/hashicorp/terraform-provider-azurerm/issues/24627))
* `azurerm_bot_channel_email` - support for the `magic_code` property ([#23129](https://github.com/hashicorp/terraform-provider-azurerm/issues/23129))
* `azurerm_cosmosdb_account` - support for the `partition_merge_enabled` property ([#24615](https://github.com/hashicorp/terraform-provider-azurerm/issues/24615))
* `azurerm_mssql_managed_database` - support for the `immutable_backups_enabled` property ([#24745](https://github.com/hashicorp/terraform-provider-azurerm/issues/24745))
* `azurerm_mssql_database` - support for the `immutable_backups_enabled` property ([#24745](https://github.com/hashicorp/terraform-provider-azurerm/issues/24745))
* `azurerm_palo_alto_next_generation_firewall_virtual_hub_panorama` - support for the `trusted_address_ranges` property ([#24459](https://github.com/hashicorp/terraform-provider-azurerm/issues/24459))
* `azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack` - support for the `trusted_address_ranges` property ([#24459](https://github.com/hashicorp/terraform-provider-azurerm/issues/24459))
* `azurerm_palo_alto_next_generation_firewall_virtual_network_panorama` - support for the `trusted_address_ranges` property ([#24459](https://github.com/hashicorp/terraform-provider-azurerm/issues/24459))
* `azurerm_servicebus_namespace` - updating to use API Version `2022-10-01-preview` ([#24650](https://github.com/hashicorp/terraform-provider-azurerm/issues/24650))
* `azurerm_spring_cloud_api_portal` - support for the `api_try_out_enabled` property ([#24696](https://github.com/hashicorp/terraform-provider-azurerm/issues/24696))
* `azurerm_spring_cloud_gateway` - support for the `local_response_cache_per_route` and `local_response_cache_per_instance` properties ([#24697](https://github.com/hashicorp/terraform-provider-azurerm/issues/24697))
* `azurerm_stream_analytics_job` - support for User Assigned Identities ([#24738](https://github.com/hashicorp/terraform-provider-azurerm/issues/24738))
* `azurerm_subscription` - refactoring to use `hashicorp/go-azure-sdk` to set tags on the subscription ([#24734](https://github.com/hashicorp/terraform-provider-azurerm/issues/24734))
* `azurerm_virtual_desktop_workspace` - correctly validate the `name` property ([#24668](https://github.com/hashicorp/terraform-provider-azurerm/issues/24668))


BUG FIXES:

* provider: skip registration for resource providers that are unavailable ([#24571](https://github.com/hashicorp/terraform-provider-azurerm/issues/24571))
* `azurerm_app_configuration` - no longer require `lifecycle_ignore_changes` for the `value` property when using a key vault reference ([#24702](https://github.com/hashicorp/terraform-provider-azurerm/issues/24702))
* `azurerm_app_service_managed_certificate` - fix casing issue in `app_service_plan_id` by parsing insensitively ([#24664](https://github.com/hashicorp/terraform-provider-azurerm/issues/24664))
* `azurerm_cognitive_deployment` - updates now include the `version` property ([#24700](https://github.com/hashicorp/terraform-provider-azurerm/issues/24700))
* `azurerm_dns_cname_record` - prevent casing issue in `target_resource_id` by parsing the ID insensitively ([#24181](https://github.com/hashicorp/terraform-provider-azurerm/issues/24181))
* `azurerm_mssql_managed_instance_failover_group` - prevent an issue when trying to create a failover group with a managed instance from a different subscription ([#24646](https://github.com/hashicorp/terraform-provider-azurerm/issues/24646))
* `azurerm_storage_account` - conditionally update properties only when needed ([#24669](https://github.com/hashicorp/terraform-provider-azurerm/issues/24669))
* `azurerm_storage_account` - change update order for `access_tier`to prevent errors when uploading blobs to the archive tier ([#22250](https://github.com/hashicorp/terraform-provider-azurerm/issues/22250))

## 3.89.0 (January 25, 2024)

FEATURES:

* New Data Source: `azurerm_data_factory_trigger_schedule` ([#24572](https://github.com/hashicorp/terraform-provider-azurerm/issues/24572))
* New Data Source: `azurerm_data_factory_trigger_schedules` ([#24572](https://github.com/hashicorp/terraform-provider-azurerm/issues/24572))
* New Data Source: `azurerm_ip_groups` ([#24540](https://github.com/hashicorp/terraform-provider-azurerm/issues/24540))
* New Data Source: `azurerm_nginx_certificate` ([#24577](https://github.com/hashicorp/terraform-provider-azurerm/issues/24577))
* New Resource: `azurerm_chaos_studio_target` ([#24580](https://github.com/hashicorp/terraform-provider-azurerm/issues/24580))
* New Resource: `azurerm_elastic_san_volume_group` ([#24166](https://github.com/hashicorp/terraform-provider-azurerm/issues/24166))
* New Resource: `azurerm_netapp_account_encryption` ([#23733](https://github.com/hashicorp/terraform-provider-azurerm/issues/23733))
* New Resource: `azurerm_redhat_openshift_cluster` ([#24375](https://github.com/hashicorp/terraform-provider-azurerm/issues/24375))

ENHANCEMENTS:

* dependencies: updating to `v0.66.1` of `github.com/hashicorp/go-azure-helpers` ([#24561](https://github.com/hashicorp/terraform-provider-azurerm/issues/24561))
* dependencies: updating to `v0.20240124.1115501` of `github.com/hashicorp/go-azure-sdk` ([#24619](https://github.com/hashicorp/terraform-provider-azurerm/issues/24619))
* `bot`: updating to API Version `2021-05-01-preview` ([#24555](https://github.com/hashicorp/terraform-provider-azurerm/issues/24555))
* `containerservice`: the SDK Clients now support logging ([#24564](https://github.com/hashicorp/terraform-provider-azurerm/issues/24564))
* `cosmosdb`: updating to API Version `2023-04-15` ([#24541](https://github.com/hashicorp/terraform-provider-azurerm/issues/24541))
* `loadtestservice`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` (and support logging) ([#24578](https://github.com/hashicorp/terraform-provider-azurerm/issues/24578))
* `managedidentity`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` (and support logging) ([#24578](https://github.com/hashicorp/terraform-provider-azurerm/issues/24578))
* `azurerm_api_management_api` - change the `id` format so specific `revision`s can be managed by Terraform ([#23031](https://github.com/hashicorp/terraform-provider-azurerm/issues/23031))
* `azurerm_data_protection_backup_vault` - the `redundancy` propety can now be set to `ZoneRedundant` ([#24556](https://github.com/hashicorp/terraform-provider-azurerm/issues/24556))
* `azurerm_data_factory_integration_runtime_azure_ssis` - support for the `credential_name` property ([#24458](https://github.com/hashicorp/terraform-provider-azurerm/issues/24458))
* `azurerm_orchestrated_virtual_machine_scale_set` - support `2022-datacenter-azure-edition-hotpatch` and `2022-datacenter-azure-edition-hotpatch-smalldisk` hotpatching images ([#23500](https://github.com/hashicorp/terraform-provider-azurerm/issues/23500))
* `azurerm_stream_analytics_job` - support for the `sku_name` property ([#24554](https://github.com/hashicorp/terraform-provider-azurerm/issues/24554))

BUG FIXES:

* Data Source: `azurerm_app_service` - parsing the API Response for `app_service_plan_id` case-insensitively ([#24626](https://github.com/hashicorp/terraform-provider-azurerm/issues/24626))
* Data Source: `azurerm_function_app` - parsing the API Response for `app_service_plan_id` case-insensitively ([#24626](https://github.com/hashicorp/terraform-provider-azurerm/issues/24626))
* `azurerm_app_configuration_key` - the value for the `value` property can now be removed/emptied ([#24582](https://github.com/hashicorp/terraform-provider-azurerm/issues/24582))

* `azurerm_app_service` - parsing the API Response for `app_service_plan_id` case-insensitively ([#24626](https://github.com/hashicorp/terraform-provider-azurerm/issues/24626))
* `azurerm_app_service_plan` - fix casing in `serverFarms` due to ID update ([#24562](https://github.com/hashicorp/terraform-provider-azurerm/issues/24562))
* `azurerm_app_service_slot` - parsing the API Response for `app_service_plan_id` case-insensitively ([#24626](https://github.com/hashicorp/terraform-provider-azurerm/issues/24626))
* `azurerm_automation_schedule` - only one `monthly_occurence` block can now be specified ([#24614](https://github.com/hashicorp/terraform-provider-azurerm/issues/24614))
* `azurerm_cognitive_deployment` - the `model.version` property is no longer required ([#24264](https://github.com/hashicorp/terraform-provider-azurerm/issues/24264))
* `azurerm_container_app` - multiple `custom_scale_rule` can not be updated ([#24509](https://github.com/hashicorp/terraform-provider-azurerm/issues/24509))
* `azurerm_container_registry_task_schedule_run_now` - prevent issue where the incorrect scheduled run in tracked if there have been multiple ([#24592](https://github.com/hashicorp/terraform-provider-azurerm/issues/24592))
* `azurerm_function_app` - parsing the API Response for `app_service_plan_id` case-insensitively ([#24626](https://github.com/hashicorp/terraform-provider-azurerm/issues/24626))
* `azurerm_function_app_slot` - parsing the API Response for `app_service_plan_id` case-insensitively ([#24626](https://github.com/hashicorp/terraform-provider-azurerm/issues/24626))
* `azurerm_logic_app_standard` - now will parse the app service ID insensitively ([#24562](https://github.com/hashicorp/terraform-provider-azurerm/issues/24562))
* `azurerm_logic_app_workflow` - the `workflow_parameters` will now correctly handle information specified by `$connections` ([#24141](https://github.com/hashicorp/terraform-provider-azurerm/issues/24141))
* `azurerm_mssql_managed_instance_security_alert_policy` - can not update empty storage attributes ([#24553](https://github.com/hashicorp/terraform-provider-azurerm/issues/24553))
* `azurerm_network_interface` - the `ip_configuration` properties are no longer added to a Load Balancer Backend if one of those `ip_configurations` is associated with a backend ([#24470](https://github.com/hashicorp/terraform-provider-azurerm/issues/24470))

## 3.88.0 (January 18, 2024)

FEATURES:

* New Data Source: `azurerm_nginx_deployment` ([#24492](https://github.com/hashicorp/terraform-provider-azurerm/issues/24492))
* New Resource: `azurerm_spring_cloud_dynatrace_application_performance_monitoring` ([#23889](https://github.com/hashicorp/terraform-provider-azurerm/issues/23889))
* New Resource: `azurerm_virtual_machine_run_command` ([#23377](https://github.com/hashicorp/terraform-provider-azurerm/issues/23377))

ENHANCEMENTS:

* dependencies: updating to `v0.20240117.1163544` of `github.com/hashicorp/go-azure-sdk` ([#24481](https://github.com/hashicorp/terraform-provider-azurerm/issues/24481))
* dependencies: updating to `v0.65.1` of `github.com/hashicorp/go-azure-helpers` ([#24479](https://github.com/hashicorp/terraform-provider-azurerm/issues/24479))
* `datashare`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24481](https://github.com/hashicorp/terraform-provider-azurerm/issues/24481))
* `kusto`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#24477](https://github.com/hashicorp/terraform-provider-azurerm/issues/24477))
* Data Source: `azurerm_application_gateway` - support for the `trusted_client_certificate.data` property ([#24474](https://github.com/hashicorp/terraform-provider-azurerm/issues/24474))
* `azurerm_service_plan`: refactoring to use `hashicorp/go-azure-sdk` ([#24483](https://github.com/hashicorp/terraform-provider-azurerm/issues/24483))
* `azurerm_container_group` - support for the `priority` property ([#24374](https://github.com/hashicorp/terraform-provider-azurerm/issues/24374))
* `azurerm_mssql_managed_database` - support for the `point_in_time_restore` property ([#24535](https://github.com/hashicorp/terraform-provider-azurerm/issues/24535))
* `azurerm_mssql_managed_instance` - now exports the `dns_zone` attribute ([#24435](https://github.com/hashicorp/terraform-provider-azurerm/issues/24435))
* `azurerm_linux_web_app_slot` - support for setting `python_version` to `3.12` ([#24363](https://github.com/hashicorp/terraform-provider-azurerm/issues/24363))
* `azurerm_linux_web_app` - support for setting `python_version` to `3.12` ([#24363](https://github.com/hashicorp/terraform-provider-azurerm/issues/24363))
* `azurerm_linux_function_app_slot` - support for setting `python_version` to `3.12` ([#24363](https://github.com/hashicorp/terraform-provider-azurerm/issues/24363))
*  `azurerm_linux_function_app` - support for setting `python_version` to `3.12` ([#24363](https://github.com/hashicorp/terraform-provider-azurerm/issues/24363))

BUG FIXES:

* `azurerm_application_gateway` - the `components` property within the `url` block is no longer computed ([#24480](https://github.com/hashicorp/terraform-provider-azurerm/issues/24480))
* `azurerm_cdn_frontdoor_route` - prevent an issue where `cdn_frontdoor_origin_path` gets removed on update if unchanged. ([#24488](https://github.com/hashicorp/terraform-provider-azurerm/issues/24488))
* `azurerm_cognitive_account` - fixing support for the `DC0` SKU ([#24526](https://github.com/hashicorp/terraform-provider-azurerm/issues/24526))

## 3.87.0 (January 11, 2024)

FEATURES:

* New Data Source: `azurerm_network_manager` ([#24398](https://github.com/hashicorp/terraform-provider-azurerm/issues/24398))
* New Resource: `azurerm_security_center_server_vulnerability_assessments_setting` ([#24299](https://github.com/hashicorp/terraform-provider-azurerm/issues/24299))

ENHANCEMENTS:

* dependencies: updating to `v0.20240111.1094251` of `github.com/hashicorp/go-azure-sdk` ([#24463](https://github.com/hashicorp/terraform-provider-azurerm/issues/24463))
* Data Source: `azurerm_mssql_database` - support for `identity`, `transparent_data_encryption_enabled`, `transparent_data_encryption_key_vault_key_id` and `transparent_data_encryption_key_automatic_rotation_enabled` ([#24412](https://github.com/hashicorp/terraform-provider-azurerm/issues/24412))
* Data Source: `azurerm_mssql_server` - support for `transparent_data_encryption_key_vault_key_id` ([#24412](https://github.com/hashicorp/terraform-provider-azurerm/issues/24412))
* `machinelearning`: updating to API Version `2023-10-01` ([#24416](https://github.com/hashicorp/terraform-provider-azurerm/issues/24416))
* `paloaltonetworks`: updating to API Version `2023-09-01` ([#24290](https://github.com/hashicorp/terraform-provider-azurerm/issues/24290))
* `azurerm_container_app` - update create time validations for `ingress.0.traffic_weight` ([#24042](https://github.com/hashicorp/terraform-provider-azurerm/issues/24042))
* `azurerm_container_app`- support for the `ip_security_restriction` block ([#23870](https://github.com/hashicorp/terraform-provider-azurerm/issues/23870))
* `azurerm_kubernetes_cluster` - properties in `default_node_pool.linux_os_config.sysctl_config` are now updateable via node pool cycling ([#24397](https://github.com/hashicorp/terraform-provider-azurerm/issues/24397))
* `azurerm_linux_web_app` - support the `VS2022` value for the `remote_debugging_version` property ([#24407](https://github.com/hashicorp/terraform-provider-azurerm/issues/24407))
* `azurerm_mssql_database` - support for `identity`, `transparent_data_encryption_key_vault_key_id` and `transparent_data_encryption_key_automatic_rotation_enabled` ([#24412](https://github.com/hashicorp/terraform-provider-azurerm/issues/24412))
* `azurerm_postgres_flexible_server` - the `sku_name` property now supports being set to `MO_Standard_E96ds_v5` ([#24367](https://github.com/hashicorp/terraform-provider-azurerm/issues/24367))
* `azurerm_role_assignment` - support for the `principal_type` property ([#24271](https://github.com/hashicorp/terraform-provider-azurerm/issues/24271))
* `azurerm_windows_web_app` - support the `VS2022` value for the `remote_debugging_version` property ([#24407](https://github.com/hashicorp/terraform-provider-azurerm/issues/24407))
* `azurerm_cdn_frontdoor_firewall_policy` - support for `request_body_check_enabled` property ([#24406](https://github.com/hashicorp/terraform-provider-azurerm/issues/24406))

BUG FIXES:

* Data Source: `azurerm_role_definition` - fix `role_definition_id` ([#24418](https://github.com/hashicorp/terraform-provider-azurerm/issues/24418))
* `azurerm_api_management` - the `sku_name` property can now be updated ([#24431](https://github.com/hashicorp/terraform-provider-azurerm/issues/24431))
* `azurerm_arc_kubernetes_flux_configuration` - prevent a bug where certain sensitive properties for `bucket` and `git_repository` were being overwritten after an update to the resource is made ([#24066](https://github.com/hashicorp/terraform-provider-azurerm/issues/24066))
* `azurerm_kubernetes_flux_configuration` - prevent a bug where certain sensitive properties for `bucket` and `git_repository` were being overwritten after an update to the resource is made ([#24066](https://github.com/hashicorp/terraform-provider-azurerm/issues/24066))
* `azure_linux_web_app` - prevent a bug in App Service processing of `application_stack` in updates to `site_config` ([#24424](https://github.com/hashicorp/terraform-provider-azurerm/issues/24424))
* `azure_linux_web_app_slot` - Fix bug in App Service processing of `application_stack` in updates to `site_config` ([#24424](https://github.com/hashicorp/terraform-provider-azurerm/issues/24424))
* `azurerm_network_manager_deployment` - update creation wait logic to better tolerate the api returning not found ([#24330](https://github.com/hashicorp/terraform-provider-azurerm/issues/24330))
* `azurerm_virtual_machine_data_disk_attachment` - do not update applications profile with disks ([#24145](https://github.com/hashicorp/terraform-provider-azurerm/issues/24145))
* `azure_windows_web_app` - prevent a bug in App Service processing of `application_stack` in updates to `site_config` ([#24424](https://github.com/hashicorp/terraform-provider-azurerm/issues/24424))
* `azure_windows_web_app_slot` - prevent a bug in App Service processing of `application_stack` in updates to `site_config` ([#24424](https://github.com/hashicorp/terraform-provider-azurerm/issues/24424))
* `azurerm_maintenance_configuration` - set the `reboot` property in flatten from `AlwaysReboot` to `Always` ([#24376](https://github.com/hashicorp/terraform-provider-azurerm/issues/24376))
* `azurerm_container_app_environment` - the `workload_profile` property can now be updated ([#24409](https://github.com/hashicorp/terraform-provider-azurerm/issues/24409))

## 3.86.0 (January 04, 2024)

FEATURES:

* New Data Source: `azurerm_dashboard_grafana` ([#24243](https://github.com/hashicorp/terraform-provider-azurerm/issues/24243))
* New Resource: `azurerm_log_analytics_workspace_table` ([#24229](https://github.com/hashicorp/terraform-provider-azurerm/issues/24229))
* New Resource: `azurerm_automation_powershell72_module` ([#23980](https://github.com/hashicorp/terraform-provider-azurerm/issues/23980))
* New Resource: `azurerm_data_factory_credential_user_managed_identity` ([#24307](https://github.com/hashicorp/terraform-provider-azurerm/issues/24307))

ENHANCEMENTS:

* dependencies: updating to `v0.20231215.1114251` of `hashicorp/go-azure-sdk` ([#24251](https://github.com/hashicorp/terraform-provider-azurerm/issues/24251))
* dependencies: `azurerm_spring_cloud_api_portal` - update to use `hashicorp/go-azure-sdk` ([#24321](https://github.com/hashicorp/terraform-provider-azurerm/issues/24321))
* Data Source: `azurerm_kusto_cluster` - now exports the `identity` block ([#24314](https://github.com/hashicorp/terraform-provider-azurerm/issues/24314))
* `azurerm_data_protection_backup_policy_postgresql` - support for the `time_zone` property ([#24312](https://github.com/hashicorp/terraform-provider-azurerm/issues/24312))
* `azurerm_data_protection_backup_policy_disk` - support for the `time_zone` property ([#24312](https://github.com/hashicorp/terraform-provider-azurerm/issues/24312))
* `azurerm_key_vault_managed_hardware_security_module` -the `tags` property can now be updated ([#24333](https://github.com/hashicorp/terraform-provider-azurerm/issues/24333))
* `azurerm_logic_app_standard` - support for the `site_config.0.public_network_access_enabled` property ([#24257](https://github.com/hashicorp/terraform-provider-azurerm/issues/24257))
* `azurerm_log_analytics_workspace_table` - support for the `plan` property ([#24341](https://github.com/hashicorp/terraform-provider-azurerm/issues/24341))
* `azurerm_linux_web_app` - support the value `20-lts` for the `node_version` property  ([#24289](https://github.com/hashicorp/terraform-provider-azurerm/issues/24289))
* `azurerm_recovery_services_vault` - support creation with immutability set to locked ([#23806](https://github.com/hashicorp/terraform-provider-azurerm/issues/23806))
* `azurerm_spring_cloud_service` - support for the `sku_tier` property ([#24103](https://github.com/hashicorp/terraform-provider-azurerm/issues/24103))

BUG FIXES:

* Data Source: `azurerm_role_definition` - correctly export the `role_definition_id` attribute ([#24320](https://github.com/hashicorp/terraform-provider-azurerm/issues/24320))
* `azurerm_bot_service` - fixing a bug where `public_network_access_enabled` was always set to `true` ([#24255](https://github.com/hashicorp/terraform-provider-azurerm/issues/24255))
* `azurerm_bot_service_azure_bot` - `tags` can now be updated ([#24332](https://github.com/hashicorp/terraform-provider-azurerm/issues/24332))
* `azurerm_cosmosdb_account` - fix validation for the `ip_range_filter` property ([#24306](https://github.com/hashicorp/terraform-provider-azurerm/issues/24306))
* `azurerm_linux_virtual_machine` - the `additional_capabilities.0.ultra_ssd_enabled` can now be changed during the update ([#24274](https://github.com/hashicorp/terraform-provider-azurerm/issues/24274))
* `azurerm_logic_app_standard` - update the default value of `version` from `~3` which is no longer supported to `~4` ([#24134](https://github.com/hashicorp/terraform-provider-azurerm/issues/24134))
* `azurerm_logic_app_standard` - fix a crash when setting the default `version` 4.0 flag ([#24322](https://github.com/hashicorp/terraform-provider-azurerm/issues/24322))
* `azurerm_iothub_device_update_account` - changing the `sku` property now creates a new resource ([#24324](https://github.com/hashicorp/terraform-provider-azurerm/issues/24324))
* `azurerm_iothub` - prevent an inconsistant value after an apply ([#24326](https://github.com/hashicorp/terraform-provider-azurerm/issues/24326))
* `azurerm_orchestrated_virtual_machine_scale_set` - correctly update the resource when hotpatch is enabled ([#24335](https://github.com/hashicorp/terraform-provider-azurerm/issues/24335))
* `azurerm_windows_virtual_machine` - the `additional_capabilities.0.ultra_ssd_enabled` can now be changed during the update ([#24274](https://github.com/hashicorp/terraform-provider-azurerm/issues/24274))
* `azurerm_scheduled_query_rules_alert` - changing the `data_source_id` now creates a new resource ([#24327](https://github.com/hashicorp/terraform-provider-azurerm/issues/24327))
* `azurerm_scheduled_query_rules_log` - changing the `data_source_id` now creates a new resource ([#24327](https://github.com/hashicorp/terraform-provider-azurerm/issues/24327))

## 3.85.0 (December 14, 2023)

FEATURES:

* New Data Source: `azurerm_locations` ([#23324](https://github.com/hashicorp/terraform-provider-azurerm/issues/23324))

ENHANCEMENTS:

* provider: support for authenticating using Azure Kubernetes Service Workload Identity ([#23965](https://github.com/hashicorp/terraform-provider-azurerm/issues/23965))
* dependencies: updating to `v0.65.0` of `github.com/hashicorp/go-azure-helpers` ([#24222](https://github.com/hashicorp/terraform-provider-azurerm/issues/24222))
* dependencies: updating to `v0.20231214.1220802` of `github.com/hashicorp/go-azure-sdk` ([#24246](https://github.com/hashicorp/terraform-provider-azurerm/issues/24246))
* dependencies: updating to version `v0.20231214.1160726` of `github.com/hashicorp/go-azure-sdk` ([#24241](https://github.com/hashicorp/terraform-provider-azurerm/issues/24241))
* dependencies: update `security/automation` to use `hashicorp/go-azure-sdk` ([#24156](https://github.com/hashicorp/terraform-provider-azurerm/issues/24156))
* dependencies `dataprotection`: updating to API Version `2023-05-01` ([#24143](https://github.com/hashicorp/terraform-provider-azurerm/issues/24143))
* `kusto`: removing the remnants of the old Resource ID Parsers now this uses `hashicorp/go-azure-sdk` ([#24238](https://github.com/hashicorp/terraform-provider-azurerm/issues/24238))
* Data Source: `azurerm_cognitive_account` - export the `identity` block ([#24214](https://github.com/hashicorp/terraform-provider-azurerm/issues/24214))
* Data Source: `azurerm_monitor_workspace` - add support for the `default_data_collection_endpoint_id` and `default_data_collection_rule_id` properties ([#24153](https://github.com/hashicorp/terraform-provider-azurerm/issues/24153))
* Data Source: `azurerm_shared_image_gallery` - add support for the `image_names` property ([#24176](https://github.com/hashicorp/terraform-provider-azurerm/issues/24176))
* `azurerm_dns_txt_record` - allow up to `4096` characters for the property `record.value` ([#24169](https://github.com/hashicorp/terraform-provider-azurerm/issues/24169))
* `azurerm_container_app` - support for the `workload_profile_name` property ([#24219](https://github.com/hashicorp/terraform-provider-azurerm/issues/24219))
* `azurerm_container_app` - suppot for the `init_container` block ([#23955](https://github.com/hashicorp/terraform-provider-azurerm/issues/23955))
* `azurerm_hpc_cache_blob_nfs_target` - support for the `verification_timer_in_seconds` and `write_back_timer_in_seconds` properties ([#24207](https://github.com/hashicorp/terraform-provider-azurerm/issues/24207))
* `azurerm_hpc_cache_nfs_target` - support for the `verification_timer_in_seconds` and `write_back_timer_in_seconds` properties ([#24208](https://github.com/hashicorp/terraform-provider-azurerm/issues/24208))
* `azurerm_linux_web_app` - make `client_secret_setting_name` optional and conflict with `client_secret_certificate_thumbprint` ([#21834](https://github.com/hashicorp/terraform-provider-azurerm/issues/21834))
* `azurerm_linux_web_app_slot` - make `client_secret_setting_name` optional and conflict with `client_secret_certificate_thumbprint` ([#21834](https://github.com/hashicorp/terraform-provider-azurerm/issues/21834))
* `azurerm_linux_web_app` - fix a bug in `app_settings` where settings could be lost ([#24221](https://github.com/hashicorp/terraform-provider-azurerm/issues/24221))
* `azurerm_linux_web_app_slot` - fix a bug in `app_settings` where settings could be lost ([#24221](https://github.com/hashicorp/terraform-provider-azurerm/issues/24221))
* `azurerm_log_analytics_workspace` - add support for the `immediate_data_purge_on_30_days_enabled` property ([#24015](https://github.com/hashicorp/terraform-provider-azurerm/issues/24015))
* `azurerm_mssql_server` - support for other identity types for the key vault key ([#24236](https://github.com/hashicorp/terraform-provider-azurerm/issues/24236))
* `azurerm_machine_learning_datastore_blobstorage` - resource now skips validation when being created ([#24078](https://github.com/hashicorp/terraform-provider-azurerm/issues/24078))
* `azurerm_machine_learning_datastore_datalake_gen2` - resource now skips validation when being created ([#24078](https://github.com/hashicorp/terraform-provider-azurerm/issues/24078))
* `azurerm_machine_learning_datastore_fileshare` - resource now skips validation when being created ([#24078](https://github.com/hashicorp/terraform-provider-azurerm/issues/24078))
* `azurerm_monitor_workspace` - support for the `default_data_collection_endpoint_id` and `default_data_collection_rule_id` properties ([#24153](https://github.com/hashicorp/terraform-provider-azurerm/issues/24153))
* `azurerm_redis_cache` - support for the `storage_account_subscription_id` property ([#24101](https://github.com/hashicorp/terraform-provider-azurerm/issues/24101))
* `azurerm_storage_blob` -  support for the `source_content` type `Page` ([#24177](https://github.com/hashicorp/terraform-provider-azurerm/issues/24177))
* `azurerm_web_application_firewall_policy` - support new values to the `rule_group_name` property ([#24194](https://github.com/hashicorp/terraform-provider-azurerm/issues/24194))
* `azurerm_windows_web_app` - make the `client_secret_setting_name` property optional and conflicts with the `client_secret_certificate_thumbprint` property ([#21834](https://github.com/hashicorp/terraform-provider-azurerm/issues/21834))
* `azurerm_windows_web_app_slot` - make the `client_secret_setting_name` property optional and conflicts with the `client_secret_certificate_thumbprint` property ([#21834](https://github.com/hashicorp/terraform-provider-azurerm/issues/21834))
* `azurerm_windows_web_app` - fix a bug in `app_settings` where settings could be lost ([#24221](https://github.com/hashicorp/terraform-provider-azurerm/issues/24221))
* `azurerm_windows_web_app_slot` - fix a bug in `app_settings` where settings could be lost ([#24221](https://github.com/hashicorp/terraform-provider-azurerm/issues/24221))
* `azurerm_cognitive_account` - add `ContentSafety` to the `kind` property validation ([#24205](https://github.com/hashicorp/terraform-provider-azurerm/issues/24205))

BUG FIXES:

* provider: fix an authentication issue with Azure Storage when running in Azure China cloud ([#24246](https://github.com/hashicorp/terraform-provider-azurerm/issues/24246))
* Data Source: `azurerm_role_definition` - fix bug where `role_definition_id` and `scope` were being incorrectly set ([#24211](https://github.com/hashicorp/terraform-provider-azurerm/issues/24211))
* `azurerm_batch_account` - fix bug where `UserAssigned, SystemAssigned` could be passed to the resource even though it isn't supported ([#24204](https://github.com/hashicorp/terraform-provider-azurerm/issues/24204))
* `azurerm_batch_pool` - fix bug where `settings_json` and `protected_settings` were not being unmarshaled ([#24075](https://github.com/hashicorp/terraform-provider-azurerm/issues/24075))
* `azurerm_bot_service_azure_bot` - fix bug where `public_network_access_enabled` was being set as the value for `LuisKey` ([#24164](https://github.com/hashicorp/terraform-provider-azurerm/issues/24164))
* `azurerm_cognitive_account_customer_managed_key` - `identity_client_id` is no longer passed to the api when it is empty ([#24231](https://github.com/hashicorp/terraform-provider-azurerm/issues/24231))
* `azurerm_linux_web_app_slot` - error when `service_plan_id` is identical to the parent `service_plan_id` ([#23403](https://github.com/hashicorp/terraform-provider-azurerm/issues/23403))
* `azurerm_management_group_template_deployment` - fixing a bug where `template_spec_version_id` couldn't be updated ([#24072](https://github.com/hashicorp/terraform-provider-azurerm/issues/24072))
* `azurerm_pim_active_role_assignment` - fix an importing issue by filtering available role assignments based on the provided `scope` ([#24077](https://github.com/hashicorp/terraform-provider-azurerm/issues/24077))
* `azurerm_pim_eligible_role_assignment` - fix an importing issue by filtering available role assignments based on the provided `scope` ([#24077](https://github.com/hashicorp/terraform-provider-azurerm/issues/24077))
* `azurerm_resource_group_template_deployment` - fixing a bug where `template_spec_version_id` couldn't be updated ([#24072](https://github.com/hashicorp/terraform-provider-azurerm/issues/24072))
* `azurerm_security_center_setting` - fix the casing for the `setting_name` `Sentinel` ([#24210](https://github.com/hashicorp/terraform-provider-azurerm/issues/24210))
* `azurerm_storage_account` - Fix crash when checking for `routingInputs.PublishInternetEndpoints` and `routingInputs.PublishMicrosoftEndpoints` ([#24228](https://github.com/hashicorp/terraform-provider-azurerm/issues/24228))
* `azurerm_storage_share_file` - prevent panic when the file specified by `source` is empty ([#24179](https://github.com/hashicorp/terraform-provider-azurerm/issues/24179))
* `azurerm_subscription_template_deployment` - fixing a bug where `template_spec_version_id` couldn't be updated ([#24072](https://github.com/hashicorp/terraform-provider-azurerm/issues/24072))
* `azurerm_tenant_template_deployment` - fixing a bug where `template_spec_version_id` couldn't be updated ([#24072](https://github.com/hashicorp/terraform-provider-azurerm/issues/24072))
* `azurerm_virtual_machine` - prevent a panic by nil checking the first element of `additional_capabilities` ([#24159](https://github.com/hashicorp/terraform-provider-azurerm/issues/24159))
* `azurerm_windows_web_app_slot` - error when `service_plan_id` is identical to the parent `service_plan_id` ([#23403](https://github.com/hashicorp/terraform-provider-azurerm/issues/23403))

## 3.84.0 (December 07, 2023)

FEATURES:

* **New Data Source:** `azurerm_storage_containers` ([#24061](https://github.com/hashicorp/terraform-provider-azurerm/issues/24061))
* **New Resource:** `azurerm_elastic_san` ([#23619](https://github.com/hashicorp/terraform-provider-azurerm/issues/23619))
* **New Resource:** `azurerm_key_vault_managed_hardware_security_module_role_assignment` ([#22332](https://github.com/hashicorp/terraform-provider-azurerm/issues/22332))
* **New Resource:** `azurerm_key_vault_managed_hardware_security_module_role_definition` ([#22332](https://github.com/hashicorp/terraform-provider-azurerm/issues/22332))

ENHANCEMENTS:

* dependencies: updating mssql elasticpools from `v5.0` to `2023-05-01-preview`
* dependencies: updating to `v0.20231207.1122031` of `github.com/hashicorp/go-azure-sdk` ([#24149](https://github.com/hashicorp/terraform-provider-azurerm/issues/24149))
* Data Source: `azurerm_storage_account` - export the primary and secondary internet and microsoft hostnames for blobs, dfs, files, queues, tables and web ([#23517](https://github.com/hashicorp/terraform-provider-azurerm/issues/23517))
* Data Source: `azurerm_cosmosdb_account` - export the `connection_strings`, `primary_sql_connection_string`, `secondary_sql_connection_string`, `primary_readonly_sql_connection_string`, `secondary_readonly_sql_connection_string`, `primary_mongodb_connection_string`, `secondary_mongodb_connection_string`, `primary_readonly_mongodb_connection_string`, and `secondary_readonly_mongodb_connection_string` attributes ([#24129](https://github.com/hashicorp/terraform-provider-azurerm/issues/24129))
* `azurerm_bot_service_azure_bot` - support for the `public_network_access_enabled` property ([#24125](https://github.com/hashicorp/terraform-provider-azurerm/issues/24125))
* `azurerm_container_app_environment` - support for the `workload_profile` property ([#23478](https://github.com/hashicorp/terraform-provider-azurerm/issues/23478))
* `azurerm_cosmosdb_cassandra_datacenter` - support for the `seed_node_ip_addresses` property ([#24076](https://github.com/hashicorp/terraform-provider-azurerm/issues/24076))
* `azurerm_firewall` - support for the `dns_proxy_enabled` property ([#20519](https://github.com/hashicorp/terraform-provider-azurerm/issues/20519))
* `azurerm_kubernetes_cluster` - support for the  `support_plan` property and the `sku_tier` `Premium` ([#23970](https://github.com/hashicorp/terraform-provider-azurerm/issues/23970))
* `azurerm_mssql_database` - support for `enclave_type` field ([#24054](https://github.com/hashicorp/terraform-provider-azurerm/issues/24054))
* `azurerm_mssql_elasticpool` - support for `enclave_type` field ([#24054](https://github.com/hashicorp/terraform-provider-azurerm/issues/24054))
* `azurerm_mssql_managed_instance` - support for more `vcores`: `6`, `10`, `12`, `20`, `48`, `56`, `96`, `128` ([#24085](https://github.com/hashicorp/terraform-provider-azurerm/issues/24085))
* `azurerm_redis_linked_server` - support for the property `geo_replicated_primary_host_name` ([#23984](https://github.com/hashicorp/terraform-provider-azurerm/issues/23984))
* `azurerm_storage_account` - expose the primary and secondary internet and microsoft hostnames for blobs, dfs, files, queues, tables and web ([#23517](https://github.com/hashicorp/terraform-provider-azurerm/issues/23517))
* `azurerm_synapse_role_assignment` - support for the `principal_type` property ([#24089](https://github.com/hashicorp/terraform-provider-azurerm/issues/24089))
* `azurerm_spring_cloud_build_deployment` - support for the `application_performance_monitoring_ids` property ([#23969](https://github.com/hashicorp/terraform-provider-azurerm/issues/23969))
* `azurerm_virtual_network_gateway` - support for the `bgp_route_translation_for_nat_enabled`, `dns_forwarding_enabled`, `ip_sec_replay_protection_enabled`, `remote_vnet_traffic_enabled`, `virtual_wan_traffic_enabled`, `radius_server`, `virtual_network_gateway_client_connection`, `policy_group`, and `ipsec_policy` property ([#23220](https://github.com/hashicorp/terraform-provider-azurerm/issues/23220))

BUG FIXES:

* `azurerm_application_insights_api_key` - prevent a bug where multiple keys couldn't be created for an Application Insights instance ([#23463](https://github.com/hashicorp/terraform-provider-azurerm/issues/23463))
* `azurerm_container_registry` - the `network_rule_set.virtual_network` property has been deprecated ([#24140](https://github.com/hashicorp/terraform-provider-azurerm/issues/24140))
* `azurerm_hdinsight_hadoop_cluster` - set `roles.edge_node.install_script_action.parameters` into state by retrieving the value provided in the user config since this property isn't returned by the API ([#23971](https://github.com/hashicorp/terraform-provider-azurerm/issues/23971))
* `azurerm_kubernetes_cluster` - prevent a bug where maintenance window start date was always recalculated and sent to the API ([#23985](https://github.com/hashicorp/terraform-provider-azurerm/issues/23985))
* `azurerm_mssql_database` - will no longer send all long retention values in payload unless set ([#24124](https://github.com/hashicorp/terraform-provider-azurerm/issues/24124))
* `azurerm_mssql_managed_database` - will no longer send all long retention values in payload unless set ([#24124](https://github.com/hashicorp/terraform-provider-azurerm/issues/24124))
* `azurerm_mssql_server_microsoft_support_auditing_policy` - only include storage endpoint in payload if set ([#24122](https://github.com/hashicorp/terraform-provider-azurerm/issues/24122))
* `azurerm_mobile_network_packet_core_control_plane` - prevent a panic if the HTTP Response is nil ([#24083](https://github.com/hashicorp/terraform-provider-azurerm/issues/24083))
* `azurerm_storage_account` - revert plan time name validation `(#23799)` ([#24142](https://github.com/hashicorp/terraform-provider-azurerm/issues/24142))
* `azurerm_web_application_firewall_policy` - split create and update function to fix lifecycle - ignore changes ([#23412](https://github.com/hashicorp/terraform-provider-azurerm/issues/23412))

## 3.83.0 (November 30, 2023)

UPGRADE NOTES:

* Key Vaults are now loaded using [the `ListBySubscription` API within the Key Vault Resource Provider](https://learn.microsoft.com/en-us/rest/api/keyvault/keyvault/vaults/list-by-subscription?view=rest-keyvault-keyvault-2022-07-01&tabs=HTTP) rather than [the Resources API](https://learn.microsoft.com/en-us/rest/api/keyvault/keyvault/vaults/list?view=rest-keyvault-keyvault-2022-07-01&tabs=HTTP). This change means that the Provider now caches the list of Key Vaults available within a Subscription, rather than loading these piecemeal to workaround stale data returned from the Resources API ([#24019](https://github.com/hashicorp/terraform-provider-azurerm/issues/24019))

FEATURES:

* New Data Source: `azurerm_stack_hci_cluster` ([#24032](https://github.com/hashicorp/terraform-provider-azurerm/issues/24032))

ENHANCEMENTS:

* dependencies: updating to `v0.20231129.1103252` of `github.com/hashicorp/go-azure-sdk` ([#24063](https://github.com/hashicorp/terraform-provider-azurerm/issues/24063))
* `automation`: updating to API Version `2023-11-01` ([#24017](https://github.com/hashicorp/terraform-provider-azurerm/issues/24017))
* `keyvault`: the cache is now populated using the `ListBySubscription` endpoint on the KeyVault Resource Provider rather than via the `Resources` API ([#24019](https://github.com/hashicorp/terraform-provider-azurerm/issues/24019)).
* `keyvault`: updating the cache to populate all Key Vaults available within the Subscription to reduce the number of API calls ([#24019](https://github.com/hashicorp/terraform-provider-azurerm/issues/24019))
* Data Source `azurerm_private_dns_zone`: refactoring to use the `ListBySubscription` API rather than the Resources API when `resource_group_name` is omitted ([#24024](https://github.com/hashicorp/terraform-provider-azurerm/issues/24024))
* `azurerm_dashboard_grafana` - support for `grafana_major_version` ([#24014](https://github.com/hashicorp/terraform-provider-azurerm/issues/24014))
* `azurerm_linux_web_app` - add support for dotnet 8 ([#23893](https://github.com/hashicorp/terraform-provider-azurerm/issues/23893))
* `azurerm_linux_web_app_slot` - add support for dotnet 8 ([#23893](https://github.com/hashicorp/terraform-provider-azurerm/issues/23893))
* `azurerm_media_transform` -  deprecate `face_detector_preset` and `video_analyzer_preset` ([#24002](https://github.com/hashicorp/terraform-provider-azurerm/issues/24002))
* `azurerm_postgresql_database` - update the validation of `collation` to include `Norwegian_Norway.1252` ([#24070](https://github.com/hashicorp/terraform-provider-azurerm/issues/24070))
* `azurerm_postgresql_flexible_server` - updating to API Version `2023-06-01-preview` ([#24016](https://github.com/hashicorp/terraform-provider-azurerm/issues/24016))
* `azurerm_redis_cache` - support for the `active_directory_authentication_enabled` property ([#23976](https://github.com/hashicorp/terraform-provider-azurerm/issues/23976))
* `azurerm_windows_web_app` - add support for dotnet 8 ([#23893](https://github.com/hashicorp/terraform-provider-azurerm/issues/23893))
* `azurerm_windows_web_app_slot` - add support for dotnet 8 ([#23893](https://github.com/hashicorp/terraform-provider-azurerm/issues/23893))
* `azurerm_storage_account` -  add `name` validation in custom diff ([#23799](https://github.com/hashicorp/terraform-provider-azurerm/issues/23799))

BUG FIXES:

* authentication: fix a bug where auxiliary tenants were not correctly authorized ([#24063](https://github.com/hashicorp/terraform-provider-azurerm/issues/24063))
* `azurerm_app_configuration` - normalize location in `replica` block ([#24074](https://github.com/hashicorp/terraform-provider-azurerm/issues/24074))
* `azurerm_cosmosdb_account` - cosmosdb version and capabilities can now be updated at the same time ([#24029](https://github.com/hashicorp/terraform-provider-azurerm/issues/24029))
* `azurerm_data_factory_flowlet_data_flow` - `source` and `sink` properties are now optional ([#23987](https://github.com/hashicorp/terraform-provider-azurerm/issues/23987))
* `azurerm_datadog_monitor_tag_rule` - correctly handle default rule ([#22806](https://github.com/hashicorp/terraform-provider-azurerm/issues/22806))
* `azurerm_ip_group`: fixing a crash when `firewall_ids` and `firewall_policy_ids` weren't parsed correctly from the API Response ([#24031](https://github.com/hashicorp/terraform-provider-azurerm/issues/24031))
* `azurerm_nginx_deployment` - add default value of `20` for `capacity` ([#24033](https://github.com/hashicorp/terraform-provider-azurerm/issues/24033))

## 3.82.0 (November 23, 2023)

FEATURES:

* New Data Source: `azurerm_monitor_workspace` ([#23928](https://github.com/hashicorp/terraform-provider-azurerm/issues/23928))
* New Resource: `azurerm_application_load_balancer_subnet_association` ([#23628](https://github.com/hashicorp/terraform-provider-azurerm/issues/23628))

ENHANCEMENTS:

* dependencies: updating to `v0.20231117.1130141` of `github.com/hashicorp/go-azure-sdk` ([#23945](https://github.com/hashicorp/terraform-provider-azurerm/issues/23945))
* `azurestackhci`: updating to API Version `2023-08-01` ([#23939](https://github.com/hashicorp/terraform-provider-azurerm/issues/23939))
* `dashboard`: updating to API Version `2023-09-01` ([#23929](https://github.com/hashicorp/terraform-provider-azurerm/issues/23929))
* `hpccache`: updating to API version `2023-05-01` ([#24005](https://github.com/hashicorp/terraform-provider-azurerm/issues/24005))
* `mssql`: updating resources using `hashicorp/go-azure-sdk` to API Version `2023-02-01-preview` ([#23721](https://github.com/hashicorp/terraform-provider-azurerm/issues/23721))
* `templatespecversions`: updating to API Version `2022-02-01` ([#24007](https://github.com/hashicorp/terraform-provider-azurerm/issues/24007))
*  Data Source: `azurerm_template_spec_version` - refactoring to use `hashicorp/go-azure-sdk` ([#24007](https://github.com/hashicorp/terraform-provider-azurerm/issues/24007))
* `azurerm_cosmosdb_postgresql_cluster` - `coordinator_storage_quota_in_mb` and `coordinator_vcore_count` are no longer required for read replicas ([#23928](https://github.com/hashicorp/terraform-provider-azurerm/issues/23928))
* `azurerm_dashboard_grafana` - `sku` can now be set to `Essential` ([#23934](https://github.com/hashicorp/terraform-provider-azurerm/issues/23934))
* `azurerm_gallery_application_version` - add support for the `config_file`, `package_file` and `target_region.exclude_from_latest` properties ([#23816](https://github.com/hashicorp/terraform-provider-azurerm/issues/23816))
* `azurerm_hdinsight_hadoop_cluster` - `script_actions` is no longer Force New ([#23888](https://github.com/hashicorp/terraform-provider-azurerm/issues/23888))
* `azurerm_hdinsight_hbase_cluster` - `script_actions` is no longer Force New ([#23888](https://github.com/hashicorp/terraform-provider-azurerm/issues/23888))
* `azurerm_hdinsight_interactive_query_cluster` - `script_actions` is no longer Force New ([#23888](https://github.com/hashicorp/terraform-provider-azurerm/issues/23888))
* `azurerm_hdinsight_kafka_cluster` - `script_actions` is no longer Force New ([#23888](https://github.com/hashicorp/terraform-provider-azurerm/issues/23888))
* `azurerm_hdinsight_spark_cluster` - `script_actions` is no longer Force New ([#23888](https://github.com/hashicorp/terraform-provider-azurerm/issues/23888))
* `azurerm_kubernetes_cluster` - add support for the `gpu_instance` property ([#23887](https://github.com/hashicorp/terraform-provider-azurerm/issues/23887))
* `azurerm_kubernetes_cluster_node_pool` - add support for the `gpu_instance` property ([#23887](https://github.com/hashicorp/terraform-provider-azurerm/issues/23887))
* `azurerm_log_analytics_workspace` - add support for the `identity` property ([#23864](https://github.com/hashicorp/terraform-provider-azurerm/issues/23864))
* `azurerm_linux_function_app` - add support for dotnet 8 ([#23638](https://github.com/hashicorp/terraform-provider-azurerm/issues/23638))
* `azurerm_linux_function_app_slot` - add support for dotnet 8 ([#23638](https://github.com/hashicorp/terraform-provider-azurerm/issues/23638))
* `azurerm_managed_lustre_file_system` - export attribute `mgs_address` ([#23942](https://github.com/hashicorp/terraform-provider-azurerm/issues/23942))
* `azurerm_mssql_database` - support for Hyperscale SKUs ([#23974](https://github.com/hashicorp/terraform-provider-azurerm/issues/23974))
* `azurerm_mssql_database` - refactoring to use `hashicorp/go-azure-sdk` ([#23721](https://github.com/hashicorp/terraform-provider-azurerm/issues/23721))
* `azurerm_mssql_server` - refactoring to use `hashicorp/go-azure-sdk` ([#23721](https://github.com/hashicorp/terraform-provider-azurerm/issues/23721))
* `azurerm_shared_image` - add support for `trusted_launch_supported` ([#23781](https://github.com/hashicorp/terraform-provider-azurerm/issues/23781))
* `azurerm_spring_cloud_container_deployment` - add support for the `application_performance_monitoring_ids` property ([#23862](https://github.com/hashicorp/terraform-provider-azurerm/issues/23862))
* `azurerm_spring_cloud_customized_accelerator` - add support for the `accelerator_type` and `path` properties ([#23797](https://github.com/hashicorp/terraform-provider-azurerm/issues/23797))
* `azurerm_point_to_site_vpn_gateway` - allow multiple `connection_configurations` blocks ([#23936](https://github.com/hashicorp/terraform-provider-azurerm/issues/23936))
* `azurerm_private_dns_cname_record` - `ttl` can now be set to 0 ([#23918](https://github.com/hashicorp/terraform-provider-azurerm/issues/23918))
* `azurerm_windows_function_app` - add support for dotnet 8 ([#23638](https://github.com/hashicorp/terraform-provider-azurerm/issues/23638))
* `azurerm_windows_function_app_slot` - add support for dotnet 8 ([#23638](https://github.com/hashicorp/terraform-provider-azurerm/issues/23638))

BUG FIXES:
* `azurerm_api_management` - correct a bug with additional location zones within the `additional_location` block ([#23943](https://github.com/hashicorp/terraform-provider-azurerm/issues/23943))
* `azurerm_dev_test_linux_virtual_machine` - `storage_type` is now ForceNew to match the updated API behaviour ([#23973](https://github.com/hashicorp/terraform-provider-azurerm/issues/23973))
* `azurerm_dev_test_windows_virtual_machine` - `storage_type` is now ForceNew to match the updated API behaviour ([#23973](https://github.com/hashicorp/terraform-provider-azurerm/issues/23973))
* `azurerm_disk_encryption_set` - resource will recreate if `identity` changes from `SystemAssigned` to `UserAssigned` ([#23904](https://github.com/hashicorp/terraform-provider-azurerm/issues/23904))
* `azurerm_eventhub_cluster`: `sku_name` is no longer ForceNew ([#24009](https://github.com/hashicorp/terraform-provider-azurerm/issues/24009))
* `azurerm_firewall` - recasing the value for `firewall_policy_id` to workaround the API returning the incorrect casing ([#23993](https://github.com/hashicorp/terraform-provider-azurerm/issues/23993))
* `azurerm_security_center_subscription_pricing` - fix a bug preventing removal of `extensions` and downgrading `tier` to `Free` ([#23821](https://github.com/hashicorp/terraform-provider-azurerm/issues/23821))
* `azurerm_windows_web_app` - fix an issue of incorrect application stack settings during update ([#23372](https://github.com/hashicorp/terraform-provider-azurerm/issues/23372))

## 3.81.0 (November 16, 2023)

ENHANCEMENTS:

* dependencies: updating to `v0.20231116.1162710` of `github.com/hashicorp/go-azure-sdk` ([#23922](https://github.com/hashicorp/terraform-provider-azurerm/issues/23922))
* `managedservices`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#23890](https://github.com/hashicorp/terraform-provider-azurerm/issues/23890))
* `network`: updating to API Version `2023-06-01` ([#23875](https://github.com/hashicorp/terraform-provider-azurerm/issues/23875))
* `servicelinker`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#23890](https://github.com/hashicorp/terraform-provider-azurerm/issues/23890))
* `storage`: refactoring usages of `github.com/hashicorp/go-azure-sdk` to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#23890](https://github.com/hashicorp/terraform-provider-azurerm/issues/23890))
* Data Source: `azurerm_network_ddos_protection_plan`: refactoring to use `hashicorp/go-azure-sdk` ([#23849](https://github.com/hashicorp/terraform-provider-azurerm/issues/23849))
* `azurerm_linux_function_app` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled` ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_linux_function_app_slot` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_linux_web_app` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_linux_web_app_slot` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_logic_app_integration_account_certificate` - `name` now accepts underscores ([#23866](https://github.com/hashicorp/terraform-provider-azurerm/issues/23866))
* `azurerm_logic_app_integration_account_partner` - `business_identity.value` now accepts underscores ([#23866](https://github.com/hashicorp/terraform-provider-azurerm/issues/23866))
* `azurerm_monitor_data_collection_rule` - added support for `WorkspaceTransforms` as `kind` ([#23873](https://github.com/hashicorp/terraform-provider-azurerm/issues/23873))
* `azurerm_network_ddos_protection_plan`: refactoring to use `hashicorp/go-azure-sdk` ([#23849](https://github.com/hashicorp/terraform-provider-azurerm/issues/23849))
* `azurerm_windows_function_app` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_windows_function_app_slot` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_windows_web_app` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))
* `azurerm_windows_web_app_slot` - add support for disabling Basic Auth for default Publishing Profile via new properties `ftp_publish_basic_authentication_enabled` and `webdeploy_publish_basic_authentication_enabled`  ([#23900](https://github.com/hashicorp/terraform-provider-azurerm/issues/23900))

## 3.80.0 (November 09, 2023)

ENHANCEMENTS:

* `internal/sdk` - Added support for pointer Types in resource models ([#23810](https://github.com/hashicorp/terraform-provider-azurerm/issues/23810))
* dependencies: updating to `v0.63.0` of `github.com/hashicorp/go-azure-helpers` ([#23785](https://github.com/hashicorp/terraform-provider-azurerm/issues/23785))
* dependencies: updating to `v0.20231106.1151347` of `github.com/hashicorp/go-azure-sdk` ([#23787](https://github.com/hashicorp/terraform-provider-azurerm/issues/23787))
* `azurerm_cognitive_deployment` - support for the `version_upgrade_option` property ([#22520](https://github.com/hashicorp/terraform-provider-azurerm/issues/22520))
* `azurerm_firewall_policy_rule_collection_group` - add support for the property `http_headers` ([#23641](https://github.com/hashicorp/terraform-provider-azurerm/issues/23641))
* `azurerm_kubernetes_cluster` - `fips_enabled` can be updated in the `default_node_pool` without recreating the cluster ([#23612](https://github.com/hashicorp/terraform-provider-azurerm/issues/23612))
* `azurerm_kusto_cluster` - the cluster `name` can now include dashes ([#23790](https://github.com/hashicorp/terraform-provider-azurerm/issues/23790))
* `azurerm_postgresql_database` - update the validation of `collation` to include support for `French_France.1252` ([#23783](https://github.com/hashicorp/terraform-provider-azurerm/issues/23783))

BUG FIXES:

* Data Source: `azurerm_data_protection_backup_vault` - removing `import` support, since Data Sources don't support being imported ([#23820](https://github.com/hashicorp/terraform-provider-azurerm/issues/23820))
* Data Source: `azurerm_kusto_database` - removing `import` support, since Data Sources don't support being imported ([#23820](https://github.com/hashicorp/terraform-provider-azurerm/issues/23820))
* Data Source: `azurerm_virtual_hub_route_table` - removing `import` support, since Data Sources don't support being imported ([#23820](https://github.com/hashicorp/terraform-provider-azurerm/issues/23820))
* `azurerm_windows_web_app` - prevent a panic with the `auto_heal.actions` property ([#23836](https://github.com/hashicorp/terraform-provider-azurerm/issues/23836))
* `azurerm_windows_web_app` - prevent a panic with the `auto_heal.triggers` property ([#23812](https://github.com/hashicorp/terraform-provider-azurerm/issues/23812))

## 3.79.0 (November 02, 2023)

ENHANCEMENTS:

* provider: log instead of error when RPs are unavailable when validating RP registrations ([#23380](https://github.com/hashicorp/terraform-provider-azurerm/issues/23380))
* `azurerm_arc_kuberenetes_cluster_extension_resource` - the `version` and `release_train` properties can now be set simultaneously ([#23692](https://github.com/hashicorp/terraform-provider-azurerm/issues/23692))
* `azurerm_container_apps` - support for the `ingress.exposed_port` property ([#23752](https://github.com/hashicorp/terraform-provider-azurerm/issues/23752))
* `azurerm_cosmosdb_postgresql_cluster` - read replica clusters can be created without specifying `administrator_login_password` property ([#23750](https://github.com/hashicorp/terraform-provider-azurerm/issues/23750))
* `azurerm_managed_application` - arrays can be supplied in the `parameter_values` property ([#23754](https://github.com/hashicorp/terraform-provider-azurerm/issues/23754))
* `azurerm_storage_management_policy` - support for properties `rule.*.actions.*.base_blob.0.tier_to_cold_after_days_since_{modification|last_access_time|creation}_greater_than and rule.*.actions.*.{snapshot|version}.0.tier_to_cold_after_days_since_creation_greater_than` ([#23574](https://github.com/hashicorp/terraform-provider-azurerm/issues/23574))

BUG FIXES:

* `azurerm_api_management_diagnostic` - the `operation_name_format` attribute will only be sent if `identifier` is set to `applicationinsights` ([#23736](https://github.com/hashicorp/terraform-provider-azurerm/issues/23736))
* `azurerm_backup_policy_vm` - fix payload by using current datetime ([#23586](https://github.com/hashicorp/terraform-provider-azurerm/issues/23586))
* `azurerm_kubernetes_cluster` - the `custom_ca_trust_certificates_base64` property can not be removed, only updated ([#23737](https://github.com/hashicorp/terraform-provider-azurerm/issues/23737))

## 3.78.0 (October 26, 2023)

FEATURES:

* New Resource: `azurerm_resource_management_private_link_association` ([#23546](https://github.com/hashicorp/terraform-provider-azurerm/issues/23546))

ENHANCEMENTS:

* dependencies: updating to `v0.20231025.1113325` of `github.com/hashicorp/go-azure-sdk` ([#23684](https://github.com/hashicorp/terraform-provider-azurerm/issues/23684))
* dependencies: updating to `v1.58.3` of `google.golang.org/grpc` ([#23691](https://github.com/hashicorp/terraform-provider-azurerm/issues/23691))
* dependencies: updating search service from `2022-09-01` to `2023-11-01` ([#23698](https://github.com/hashicorp/terraform-provider-azurerm/issues/23698)) 
* Data Source: `azurerm_monitor_workspace` - export `query_endpoint` ([#23629](https://github.com/hashicorp/terraform-provider-azurerm/issues/23629))
* `azurerm_express_route_port` - support for `macsec_sci_enabled` ([#23625](https://github.com/hashicorp/terraform-provider-azurerm/issues/23625))
* `azurerm_eventhub_namespace_customer_managed_key` - support for the `user_assigned_identity_id` property ([#23635](https://github.com/hashicorp/terraform-provider-azurerm/issues/23635))
* `azurerm_postgresql_flexible_server` - `private_dns_zone_id` is no longer ForceNew and case is suppressed  ([#23660](https://github.com/hashicorp/terraform-provider-azurerm/issues/23660))
* `azurerm_synapse_workspace` - add support for `azuread_authentication_only` ([#23659](https://github.com/hashicorp/terraform-provider-azurerm/issues/23659))
* `azurerm_redis_enterprise_cluster` - support for new location `Japan East` ([#23696](https://github.com/hashicorp/terraform-provider-azurerm/issues/23696))
* `azurerm_search_service` - support for `semantic_search_sku` field ([#23698](https://github.com/hashicorp/terraform-provider-azurerm/issues/23698))

BUG FIXES:

* `azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack` - added lock for ruleStackID ([#23601](https://github.com/hashicorp/terraform-provider-azurerm/issues/23601))
* `azurerm_cognitive_deployment` - remove forceNew tag from `rai_policy_name` ([#23697](https://github.com/hashicorp/terraform-provider-azurerm/issues/23697))

## 3.77.0 (October 19, 2023)

FEATURES:

* New Resources: `azurerm_application_load_balancer_frontend` ([#23411](https://github.com/hashicorp/terraform-provider-azurerm/issues/23411))
* New Resources: `azurerm_dev_center` ([#23538](https://github.com/hashicorp/terraform-provider-azurerm/issues/23538))
* New Resources: `azurerm_dev_center_project` ([#23538](https://github.com/hashicorp/terraform-provider-azurerm/issues/23538))

ENHANCEMENTS:

* dependencies: updating to `v0.62.0` of `github.com/hashicorp/go-azure-helpers` ([#23581](https://github.com/hashicorp/terraform-provider-azurerm/issues/23581))
* dependencies: updating Kusto SDK from `2023-05-02` to `2023-08-15` ([#23598](https://github.com/hashicorp/terraform-provider-azurerm/issues/23598))
* dependencies: updating nginx from `2022-08-01` to `2023-04-01` ([#23583](https://github.com/hashicorp/terraform-provider-azurerm/issues/23583))
* `netapp`: updating to use API Version `2023-05-01` ([#23576](https://github.com/hashicorp/terraform-provider-azurerm/issues/23576))
* `springcloud`: updating to use API Version `2023-09-01-preview` ([#23544](https://github.com/hashicorp/terraform-provider-azurerm/issues/23544))
* `storage`: updating to use API Version `2023-01-01` ([#23543](https://github.com/hashicorp/terraform-provider-azurerm/issues/23543))
* `internal/sdk`: fixing an issue where struct fields containing `removedInNextMajorVersion` wouldn't be decoded correctly ([#23564](https://github.com/hashicorp/terraform-provider-azurerm/issues/23564))
* `internal/sdk`: struct tag parsing is now handled consistently during both encoding and decoding ([#23568](https://github.com/hashicorp/terraform-provider-azurerm/issues/23568))
* provider: the `roll_instances_when_required` provider feature in the `virtual_machine_scale_set` block is now optional ([#22976](https://github.com/hashicorp/terraform-provider-azurerm/issues/22976))
* Data Source: `azurerm_automation_account`: refactoring the remaining usage of `Azure/azure-sdk-for-go` to use `hashicorp/go-azure-sdk` ([#23555](https://github.com/hashicorp/terraform-provider-azurerm/issues/23555))
* `azurerm_automation_account`: refactoring the remaining usage of `Azure/azure-sdk-for-go` to use `hashicorp/go-azure-sdk` ([#23555](https://github.com/hashicorp/terraform-provider-azurerm/issues/23555))
* `azurerm_resource_deployment_script_azure_cli` - improve validation for the `version` property to support newer versions ([#23370](https://github.com/hashicorp/terraform-provider-azurerm/issues/23370))
* `azurerm_resource_deployment_script_azure_power_shell` - improve validation for the `version` property to support newer versions ([#23370](https://github.com/hashicorp/terraform-provider-azurerm/issues/23370))
* `azurerm_nginx_deployment` - support for the `capacity` and `email` properties ([#23596](https://github.com/hashicorp/terraform-provider-azurerm/issues/23596))

BUG FIXES:

* Data Source: `azurerm_virtual_hub_connection` - export the `inbound_route_map_id`, `outbound_route_map_id`, and `static_vnet_local_route_override_criteria` attributes in the `routing` block, and fix a bug where these attributes could not be set ([#23491](https://github.com/hashicorp/terraform-provider-azurerm/issues/23491))
* `azurerm_cdn_frontdoor_rule` - the `url_filename_condition` properties `match_values` is now optional if `operator` is set to `Any` ([#23541](https://github.com/hashicorp/terraform-provider-azurerm/issues/23541))
* `azurerm_shared_image_gallery` - added the `Private` and `Groups` options for the `sharing.permission` property ([#23570](https://github.com/hashicorp/terraform-provider-azurerm/issues/23570))
* `azurerm_redis_cache` - fixed incorrect ssl values for `redis_primary_connection_string` and `secondary_connection_string` ([#23575](https://github.com/hashicorp/terraform-provider-azurerm/issues/23575))
* `azurerm_monitor_activity_log_alert` - the `recommend_category` property now can be set to `HighAvailability` ([#23605](https://github.com/hashicorp/terraform-provider-azurerm/issues/23605))
* `azurerm_recovery_services_vault` - the `encryption` property can now be used with the `cross_region_restore_enabled` property ([#23618](https://github.com/hashicorp/terraform-provider-azurerm/issues/23618))
* `azurerm_storage_account_customer_managed_key` - prevent a panic when the keyvault id is empty ([#23599](https://github.com/hashicorp/terraform-provider-azurerm/issues/23599))

## 3.76.0 (October 12, 2023)

FEATURES:

* New Resource: `azurerm_security_center_storage_defender` ([#23242](https://github.com/hashicorp/terraform-provider-azurerm/issues/23242))
* New Resource: `azurerm_spring_cloud_application_insights_application_performance_monitoring` ([#23107](https://github.com/hashicorp/terraform-provider-azurerm/issues/23107))

ENHANCEMENTS:

* provider: updating to build using Go `1.21.3` ([#23514](https://github.com/hashicorp/terraform-provider-azurerm/issues/23514))
* dependencies: updating to `v0.20231012.1141427` of `github.com/hashicorp/go-azure-sdk` ([#23534](https://github.com/hashicorp/terraform-provider-azurerm/issues/23534))
* Data Source: `azurerm_application_gateway` - support for `backend_http_settings`, `global`, `gateway_ip_configuration` and additional attributes ([#23318](https://github.com/hashicorp/terraform-provider-azurerm/issues/23318))
* Data Source: `azurerm_network_service_tags` - export the `name` attribute ([#23382](https://github.com/hashicorp/terraform-provider-azurerm/issues/23382))
* `azurerm_cosmosdb_postgresql_cluster` - add support for `sql_version` of `16` and `citus_version` of `12.1` ([#23476](https://github.com/hashicorp/terraform-provider-azurerm/issues/23476))
* `azurerm_palo_alto_local_rulestack` - correctly normalize the `location` property ([#23483](https://github.com/hashicorp/terraform-provider-azurerm/issues/23483))
* `azurerm_static_site` - add support for `app_settings` ([#23421](https://github.com/hashicorp/terraform-provider-azurerm/issues/23421))

BUG FIXES:

* `azurerm_automation_schedule` - fix a bug when updating `start_time` ([#23494](https://github.com/hashicorp/terraform-provider-azurerm/issues/23494))
* `azurerm_eventhub` - remove ForceNew and check `partition_count` is not decreased ([#23499](https://github.com/hashicorp/terraform-provider-azurerm/issues/23499))
* `azurerm_managed_lustre_file_system` - update validation for `storage_capacity_in_tb` according to `sku_name` in use ([#23428](https://github.com/hashicorp/terraform-provider-azurerm/issues/23428))
* `azurerm_virtual_machine` - fix a crash when the API response for the `os_profile` block contains nil properties ([#23535](https://github.com/hashicorp/terraform-provider-azurerm/issues/23535))

## 3.75.0 (September 28, 2023)

FEATURES:

* New Resource: `azurerm_application_load_balancer` ([#22517](https://github.com/hashicorp/terraform-provider-azurerm/issues/22517))
* New Resource: `azurerm_resource_management_private_link` ([#23098](https://github.com/hashicorp/terraform-provider-azurerm/issues/23098))

ENHANCEMENTS:

* dependencies: `firewall` migrated to `hashicorp/go-azure-sdk` ([#22863](https://github.com/hashicorp/terraform-provider-azurerm/issues/22863))
* `azurerm_bot_service_azure_bot` - add support for the `icon_url` property ([#23114](https://github.com/hashicorp/terraform-provider-azurerm/issues/23114))
* `azurerm_cognitive_deployment` - `capacity` property is now updateable ([#23251](https://github.com/hashicorp/terraform-provider-azurerm/issues/23251))
* `azurerm_container_group` - added support for `key_vault_user_identity_id` ([#23332](https://github.com/hashicorp/terraform-provider-azurerm/issues/23332))
* `azurerm_data_factory` - added support for the `publish_enabled` property ([#2334](https://github.com/hashicorp/terraform-provider-azurerm/issues/2334))
* `azurerm_firewall_policy_rule_collection_group` - add support for the `description` property ([#23354](https://github.com/hashicorp/terraform-provider-azurerm/issues/23354))
* `azurerm_kubernetes_cluster` - `network_profile.network_policy` can be migrated to `cilium` ([#23342](https://github.com/hashicorp/terraform-provider-azurerm/issues/23342))
* `azurerm_log_analytics_workspace` - add support for the `data_collection_rule_id` property ([#23347](https://github.com/hashicorp/terraform-provider-azurerm/issues/23347))
* `azurerm_mysql_flexible_server` - add support for the `io_scaling_enabled` property ([#23329](https://github.com/hashicorp/terraform-provider-azurerm/issues/23329))

BUG FIXES:

* `azurerm_api_management_api` - fix importing `openapi` format content file issue ([#23348](https://github.com/hashicorp/terraform-provider-azurerm/issues/23348))
* `azurerm_cdn_frontdoor_rule` - allow a `cache_duration` of `00:00:00` ([#23384](https://github.com/hashicorp/terraform-provider-azurerm/issues/23384))
* `azurerm_cosmosdb_cassandra_datacenter` - `sku_name` is now updatable ([#23419](https://github.com/hashicorp/terraform-provider-azurerm/issues/23419))
* `azurerm_key_vault_certificate` - fix a bug that prevented soft-deleted certificates from being recovered ([#23204](https://github.com/hashicorp/terraform-provider-azurerm/issues/23204))
* `azurerm_log_analytics_solution` - fix create and update lifecycle of resource by splitting methods ([#23333](https://github.com/hashicorp/terraform-provider-azurerm/issues/23333))
* `azurerm_management_group_subscription_association` - mark resource as gone correctly if not found when retrieving ([#23335](https://github.com/hashicorp/terraform-provider-azurerm/issues/23335))
* `azurerm_management_lock` - add polling after create and delete to check for RP propagation ([#23345](https://github.com/hashicorp/terraform-provider-azurerm/issues/23345))
* `azurerm_monitor_diagnostic_setting` - added validation to ensure at least one of `category` or `category_group` is supplied ([#23308](https://github.com/hashicorp/terraform-provider-azurerm/issues/23308))
* `azurerm_palo_alto_local_rulestack_prefix_list` - fix rulestack not being committed on delete ([#23362](https://github.com/hashicorp/terraform-provider-azurerm/issues/23362))
* `azurerm_palo_alto_local_rulestack_fqdn_list` - fix rulestack not being committed on delete ([#23362](https://github.com/hashicorp/terraform-provider-azurerm/issues/23362))
* `security_center_subscription_pricing_resource` - disabled extensions logic now works as expected ([#22997](https://github.com/hashicorp/terraform-provider-azurerm/issues/22997))


## 3.74.0 (September 21, 2023)

NOTES:

* `azurerm_synapse_sql_pool` - users that have imported `azurerm_synapse_sql_pool` resources that were created outside of Terraform using an `LRS` storage account type will need to use `ignore_changes` to avoid the resource from being destroyed and recreated.

FEATURES:

* **New Resource**: `azurerm_arc_resource_bridge_appliance` ([#23108](https://github.com/hashicorp/terraform-provider-azurerm/issues/23108))
* **New Resource**: `azurerm_data_factory_dataset_azure_sql_table` ([#23264](https://github.com/hashicorp/terraform-provider-azurerm/issues/23264))
* **New Resource**: `azurerm_function_app_connection` ([#23127](https://github.com/hashicorp/terraform-provider-azurerm/issues/23127))

ENHANCEMENTS:

* dependencies: updating to `v0.20230918.1115907` of `github.com/hashicorp/go-azure-sdk` ([#23337](https://github.com/hashicorp/terraform-provider-azurerm/issues/23337))
* dependencies: downgrading to `v1.12.5` of `github.com/rickb777/date` ([#23296](https://github.com/hashicorp/terraform-provider-azurerm/issues/23296))
* `mysql`: updating to use API Version `2022-01-01` ([#23320](https://github.com/hashicorp/terraform-provider-azurerm/issues/23320))
* `azurerm_app_configuration` - support for the `replica` block ([#22452](https://github.com/hashicorp/terraform-provider-azurerm/issues/22452))
* `azurerm_bot_channel_directline` - support for `user_upload_enabled`, `endpoint_parameters_enabled`, and `storage_enabled` ([#23149](https://github.com/hashicorp/terraform-provider-azurerm/issues/23149))
* `azurerm_container_app` - support for scale rules ([#23294](https://github.com/hashicorp/terraform-provider-azurerm/issues/23294))
* `azurerm_container_app_environment` - support for zone redundancy ([#23313](https://github.com/hashicorp/terraform-provider-azurerm/issues/23313))
* `azurerm_container_group` - support for the `key_vault_user_identity_id` property for Customer Managed Keys ([#23332](https://github.com/hashicorp/terraform-provider-azurerm/issues/23332))
* `azurerm_cosmosdb_account` - support for MongoDB connection strings ([#23331](https://github.com/hashicorp/terraform-provider-azurerm/issues/23331))
* `azurerm_data_factory_dataset_delimited_text` - support for the `dynamic_file_system_enabled`, `dynamic_path_enabled`, and `dynamic_filename_enabled` properties ([#23261](https://github.com/hashicorp/terraform-provider-azurerm/issues/23261))
* `azurerm_data_factory_dataset_parquet` - support for the `azure_blob_fs_location` block ([#23261](https://github.com/hashicorp/terraform-provider-azurerm/issues/23261))
* `azurerm_monitor_diagnostic_setting` - validation to ensure either `category` or `category_group` are supplied in `enabled_log` and `log` blocks ([#23308](https://github.com/hashicorp/terraform-provider-azurerm/issues/23308))
* `azurerm_network_interface` - support for the `auxiliary_mode` and `auxiliary_sku` properties ([#22979](https://github.com/hashicorp/terraform-provider-azurerm/issues/22979))
* `azurerm_postgresql_flexible_server` - increased the maximum supported value for `storage_mb` ([#23277](https://github.com/hashicorp/terraform-provider-azurerm/issues/23277))
* `azurerm_shared_image_version` - support for the `replicated_region_deletion_enabled` and `target_region.exclude_from_latest_enabled` properties ([#23147](https://github.com/hashicorp/terraform-provider-azurerm/issues/23147))
* `azurerm_storage_account` - support for setting `domain_name` and `domain_guid` for `AADKERB` ([#22833](https://github.com/hashicorp/terraform-provider-azurerm/issues/22833))
* `azurerm_storage_account_customer_managed_key` - support for cross-tenant customer-managed keys with the `federated_identity_client_id`, and `key_vault_uri` properties ([#20356](https://github.com/hashicorp/terraform-provider-azurerm/issues/20356))
* `azurerm_web_application_firewall_policy` - support for the `rate_limit_duration`, `rate_limit_threshold`, `group_rate_limit_by`, and `request_body_inspect_limit_in_kb` properties ([#23239](https://github.com/hashicorp/terraform-provider-azurerm/issues/23239))

BUG FIXES:

* Data Source: `azurerm_container_app_environment`: fix `log_analytics_workspace_name` output to correct value ([#23298](https://github.com/hashicorp/terraform-provider-azurerm/issues/23298))
* `azurerm_api_management_api` - set the `service_url` property when importing the resource ([#23011](https://github.com/hashicorp/terraform-provider-azurerm/issues/23011))
* `azurerm_app_configuration` - prevent crash by nil checking the encryption configuration ([#23302](https://github.com/hashicorp/terraform-provider-azurerm/issues/23302))
* `azurerm_app_configuration_feature` - update `percentage_filter_value` to accept correct type of float ([#23263](https://github.com/hashicorp/terraform-provider-azurerm/issues/23263))
* `azurerm_container_app` - fix an issue with `commands` and `args` being overwritten when using multiple containers ([#23338](https://github.com/hashicorp/terraform-provider-azurerm/issues/23338))
* `azurerm_key_vault_certificate` - fix issue where certificates couldn't be recovered anymore ([#23204](https://github.com/hashicorp/terraform-provider-azurerm/issues/23204))
* `azurerm_key_vault_key` - the ForceNew when `expiration_date` is removed from the config file ([#23327](https://github.com/hashicorp/terraform-provider-azurerm/issues/23327))
* `azurerm_linux_function_app` - fix a bug in setting the storage settings when using Elastic Premium plans ([#21212](https://github.com/hashicorp/terraform-provider-azurerm/issues/21212))
* `azurerm_linux_web_app` - fix docker app stack update ([#23303](https://github.com/hashicorp/terraform-provider-azurerm/issues/23303))
* `azurerm_linux_web_app` - fix crash in auto heal expansion ([#21328](https://github.com/hashicorp/terraform-provider-azurerm/issues/21328))
* `azurerm_linux_web_app_slot` - fix docker app stack update ([#23303](https://github.com/hashicorp/terraform-provider-azurerm/issues/23303))
* `azurerm_linux_web_app_slot` - fix crash in auto heal expansion ([#21328](https://github.com/hashicorp/terraform-provider-azurerm/issues/21328))
* `azurerm_log_analytics_solution` - fix bug where the resource wasn't handling successful creation on subsequent applies ([#23312](https://github.com/hashicorp/terraform-provider-azurerm/issues/23312))
* `azurerm_management_group_subscription_association` - fix bug to correctly mark resource as gone if not found during read ([#23335](https://github.com/hashicorp/terraform-provider-azurerm/issues/23335))
* `azurerm_mssql_elasticpool` - remove check that prevents `license_type` from being set for certain skus ([#23262](https://github.com/hashicorp/terraform-provider-azurerm/issues/23262))
* `azurerm_servicebus_queue` - fixing an issue where `auto_delete_on_idle` couldn't be set to `P10675199DT2H48M5.4775807S` ([#23296](https://github.com/hashicorp/terraform-provider-azurerm/issues/23296))
* `azurerm_servicebus_topic` - fixing an issue where `auto_delete_on_idle` couldn't be set to `P10675199DT2H48M5.4775807S` ([#23296](https://github.com/hashicorp/terraform-provider-azurerm/issues/23296))
* `azurerm_storage_account` - prevent sending unsupported blob properties in payload for `Storage` account kind ([#23288](https://github.com/hashicorp/terraform-provider-azurerm/issues/23288))
* `azurerm_synapse_sql_pool` - expose `storage_account_type` ([#23217](https://github.com/hashicorp/terraform-provider-azurerm/issues/23217))
* `azurerm_windows_function_app` - fix a bug in setting the storage settings when using Elastic Premium plans ([#21212](https://github.com/hashicorp/terraform-provider-azurerm/issues/21212))
* `azurerm_windows_web_app` - fix docker app stack update ([#23303](https://github.com/hashicorp/terraform-provider-azurerm/issues/23303))
* `azurerm_windows_web_app_slot` - fix docker app stack update ([#23303](https://github.com/hashicorp/terraform-provider-azurerm/issues/23303))

DEPRECATIONS:

* `azurerm_application_gateway` - deprecate `Standard` and `WAF` skus ([#23310](https://github.com/hashicorp/terraform-provider-azurerm/issues/23310))
* `azurerm_bot_channel_web_chat` - deprecate `site_names` in favour of `site` block ([#23161](https://github.com/hashicorp/terraform-provider-azurerm/issues/23161))
* `azurerm_monitor_diagnostic_setting` - deprecate `retention_policy` in favour of `azurerm_storage_management_policy` ([#23260](https://github.com/hashicorp/terraform-provider-azurerm/issues/23260))

## 3.73.0 (September 14, 2023)

FEATURES:

* **New Resource**: `azurerm_iothub_endpoint_cosmosdb_account` ([#23065](https://github.com/hashicorp/terraform-provider-azurerm/issues/23065))
* **New Resource**: `azurerm_virtual_hub_routing_intent` ([#23138](https://github.com/hashicorp/terraform-provider-azurerm/issues/23138))

ENHANCEMENTS:

* dependencies: updating to `v0.1.1` of `github.com/btubbs/datetime` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v1.3.1` of `github.com/google/uuid` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v0.61.0` of `github.com/hashicorp/go-azure-helpers` ([#23249](https://github.com/hashicorp/terraform-provider-azurerm/issues/23249))
* dependencies: updating to `v0.20230907.1113401` of `github.com/hashicorp/go-azure-sdk` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v1.5.0` of `github.com/hashicorp/go-hclog` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v2.29.0` of `github.com/hashicorp/terraform-plugin-sdk/v2` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v1.5.1` of `github.com/hashicorp/terraform-plugin-testing` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v1.20.2` of `github.com/rickb777/date` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v0.13.0` of `golang.org/x/crypto` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v0.15.0` of `golang.org/x/net` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* dependencies: updating to `v0.13.0` of `golang.org/x/tools` ([#23221](https://github.com/hashicorp/terraform-provider-azurerm/issues/23221))
* `azurerm_bot_channel_ms_teams` - support for `deployment_environment` ([#23122](https://github.com/hashicorp/terraform-provider-azurerm/issues/23122))
* `azurerm_managed_disk` - updating to use API Version `2023-04-02` ([#23233](https://github.com/hashicorp/terraform-provider-azurerm/issues/23233))
* `azurerm_managed_disk` - support for `optimized_frequent_attach_enabled` ([#23241](https://github.com/hashicorp/terraform-provider-azurerm/issues/23241))
* `azurerm_managed_disk` - support for `performance_plus_enabled` ([#23241](https://github.com/hashicorp/terraform-provider-azurerm/issues/23241))
* `azurerm_maps_account` - support for `local_authentication_enabled` ([#23216](https://github.com/hashicorp/terraform-provider-azurerm/issues/23216))
* `azurerm_mssql_elasticpool` - support for configuring `license_type` when using the `Hyperscale` sku ([#23256](https://github.com/hashicorp/terraform-provider-azurerm/issues/23256))
* `azurerm_security_center_assessment_policy` - refactoring to use `hashicorp/go-azure-sdk` ([#23158](https://github.com/hashicorp/terraform-provider-azurerm/issues/23158))

BUG FIXES:

* `azurerm_api_management` - split create and update methods  ([#23259](https://github.com/hashicorp/terraform-provider-azurerm/issues/23259))
* `azurerm_api_management_backend` - fixing a panic when flattening the `credentials` block ([#23219](https://github.com/hashicorp/terraform-provider-azurerm/issues/23219))
* `azurerm_key_vault_certificate` - fixing a regression where certificates from a custom/unknown issuer would be polled indefinitely ([#23214](https://github.com/hashicorp/terraform-provider-azurerm/issues/23214))
* `azurerm_redis_cache` - prevent sending `redis_configuration.aof_backup_enabled` when the sku is not `Premium` to avoid API error ([#22774](https://github.com/hashicorp/terraform-provider-azurerm/issues/22774))
* `azurerm_web_application_firewall_policy` - capture and toggle state of `custom_rule` blocks with an `enabled` field ([#23163](https://github.com/hashicorp/terraform-provider-azurerm/issues/23163))

## 3.72.0 (September 07, 2023)

FEATURES:

* Provider Feature: subscription cancellation on `destroy` can now be disabled via the provider `features` block ([#19936](https://github.com/hashicorp/terraform-provider-azurerm/issues/19936))
* **New Data Source**: `netapp_volume_quota_rule` ([#23042](https://github.com/hashicorp/terraform-provider-azurerm/issues/23042))
* **New Resource**: `azurerm_automation_python3_package` ([#23087](https://github.com/hashicorp/terraform-provider-azurerm/issues/23087))
* **New Resource**: `netapp_volume_quota_rule` ([#23042](https://github.com/hashicorp/terraform-provider-azurerm/issues/23042))

ENHANCEMENTS:

* dependencies: updating to `v0.20230906.1160501` of `github.com/hashicorp/go-azure-sdk` ([#23191](https://github.com/hashicorp/terraform-provider-azurerm/issues/23191))
* `containerapps`: updating to API Version `2023-05-01` ([#22804](https://github.com/hashicorp/terraform-provider-azurerm/issues/22804))
* `keyvault`: upgrade remaining resources to `2023-02-01` ([#23089](https://github.com/hashicorp/terraform-provider-azurerm/issues/23089))
* `redisenterprise`: updating to API Version `2023-07-01` ([#23178](https://github.com/hashicorp/terraform-provider-azurerm/issues/23178))
* `vpngateway`: updating to use `hashicorp/go-azure-sdk` ([#22906](https://github.com/hashicorp/terraform-provider-azurerm/issues/22906))
* `internal/sdk`: typed resources using a custom importer now get a timed context ([#23160](https://github.com/hashicorp/terraform-provider-azurerm/issues/23160))
* `azurerm_batch_pool` - support for `accelerated_networking_enabled` ([#23021](https://github.com/hashicorp/terraform-provider-azurerm/issues/23021))
* `azurerm_batch_pool` - support for `automatic_upgrade_enabled` ([#23021](https://github.com/hashicorp/terraform-provider-azurerm/issues/23021))
* `azurerm_bot_channel_direct_line_speech` - support for the `cognitive_account_id` property ([#23106](https://github.com/hashicorp/terraform-provider-azurerm/issues/23106))
* `azurerm_bot_service_azure_bot` - support for the `local_authentication_enabled` property ([#23096](https://github.com/hashicorp/terraform-provider-azurerm/issues/23096))
* `azurerm_container_app_environment` - support for the `dapr_application_insights_connection_string` ([#23080](https://github.com/hashicorp/terraform-provider-azurerm/issues/23080))
* `azurerm_cosmosdb_cassandra_datacenter` - refactoring to use `hashicorp/go-azure-sdk` ([#23110](https://github.com/hashicorp/terraform-provider-azurerm/issues/23110))
* `azurerm_cosmosdb_cassandra_datacenter` - updating to API Version `2023-04-15` ([#23110](https://github.com/hashicorp/terraform-provider-azurerm/issues/23110))
* `azurerm_kubernetes_cluster` - Azure CNI can be updated to use `overlay` ([#22709](https://github.com/hashicorp/terraform-provider-azurerm/issues/22709))
* `azurerm_monitor_diagnostic_setting` - deprecating `retention_policy` within `enabled_log` ([#23029](https://github.com/hashicorp/terraform-provider-azurerm/issues/23029))
* `azurerm_mssql_database` - split create and update methods ([#23209](https://github.com/hashicorp/terraform-provider-azurerm/issues/23209))
* `azurerm_postgresql_database` - `collation` can now be set to `English_United Kingdom.1252` ([#23171](https://github.com/hashicorp/terraform-provider-azurerm/issues/23171))
* `azurerm_postgresql_flexible_database` - `collation` can now be set to `English_United Kingdom.1252` ([#23171](https://github.com/hashicorp/terraform-provider-azurerm/issues/23171))
* `azurerm_postgresql_flexible_server` - support for the `auto_grow_enabled` property ([#23069](https://github.com/hashicorp/terraform-provider-azurerm/issues/23069))
* `azurerm_redis_enterprise_cluster` - support for Flash clusters in Brazil South ([#23200](https://github.com/hashicorp/terraform-provider-azurerm/issues/23200))
* `azurerm_resource_provider_registration` - refactoring to use `hashicorp/go-azure-sdk` ([#23072](https://github.com/hashicorp/terraform-provider-azurerm/issues/23072))
* `azurerm_virtual_machine_extension` - support for `provision_after_extensions` ([#23124](https://github.com/hashicorp/terraform-provider-azurerm/issues/23124))
* `azurerm_virtual_network_gateway` - increasing the default timeout for create to `90m` ([#23003](https://github.com/hashicorp/terraform-provider-azurerm/issues/23003))
* `azurerm_virtual_hub_connection` - support for `inbound_route_map_id`, `outbound_route_map_id`, and `static_vnet_local_route_override_criteria` properties ([#23049](https://github.com/hashicorp/terraform-provider-azurerm/issues/23049))

BUG FIXES:

* `azurerm_api_management_api_policy` - added state migration to mutate id's ending in `policies/policy` ([#23128](https://github.com/hashicorp/terraform-provider-azurerm/issues/23128))
* `azurerm_api_management_api_operation_policy` - added state migration to mutate id's ending in `policies/policy` ([#23128](https://github.com/hashicorp/terraform-provider-azurerm/issues/23128))
* `azurerm_api_management_product_policy` - added state migration to mutate id's ending in `policies/policy` ([#23128](https://github.com/hashicorp/terraform-provider-azurerm/issues/23128))
* `azurerm_automation_account` - fixes logic for `local_authentication_enabled` ([#23082](https://github.com/hashicorp/terraform-provider-azurerm/issues/23082))
* `azurerm_key_vault_managed_storage_account` - check id can be parsed correctly before setting it in state ([#23022](https://github.com/hashicorp/terraform-provider-azurerm/issues/23022))
* `azurerm_monitor_diagnostic_setting` - fix `enabled_log` feature flagged schema ([#23093](https://github.com/hashicorp/terraform-provider-azurerm/issues/23093))
* `azurerm_pim_active_role_assignment`: polling for the duration of the timeout, rather than a fixed 5 minute value ([#22932](https://github.com/hashicorp/terraform-provider-azurerm/issues/22932))
* `azurerm_policy_set_definition` - only sending `parameters` when a value is configured ([#23155](https://github.com/hashicorp/terraform-provider-azurerm/issues/23155))
* `azurerm_synapse_workspace` - fixes index out-of-range panic when parsing `storage_data_lake_gen2_filesystem_id` ([#23019](https://github.com/hashicorp/terraform-provider-azurerm/issues/23019))
* `machine_learning_datastore_*` - fixes container ids ([#23140](https://github.com/hashicorp/terraform-provider-azurerm/issues/23140))
* `azurerm_key_vault_certificate` - id now points to new version when certificate is updated ([#23135](https://github.com/hashicorp/terraform-provider-azurerm/issues/23135))
* `azurerm_site_recovery_replicated_vm` - update `network_interface` diff so replicated items now can be updated ([#23199](https://github.com/hashicorp/terraform-provider-azurerm/issues/23199))

DEPRECATION:

* Data Source: `azure_monitor_log_profile` - Azure is retiring Azure Log Profiles on the 30th of September 2026 ([#23146](https://github.com/hashicorp/terraform-provider-azurerm/issues/23146))
* `azure_monitor_log_profile` - Azure is retiring Azure Log Profiles on the 30th of September 2026 ([#23146](https://github.com/hashicorp/terraform-provider-azurerm/issues/23146))

## 3.71.0 (August 24, 2023)

BREAKING CHANGES:

* **App Service `win32_status` property** - Due to a change made in the service to the underlying type of the Auto Heal property `win32_status` combined with a prior bug (in `v3.62.1` and earlier) causing the value of this property to be stored incorrectly in state as an empty string, the value of this property could not be updated or state migrated to accommodate the necessary type change in the state. This results in the resources named above returning an error of a number is needed when decoding the state for this value. Unfortunately, this is a breaking change and will require users of this field to change their Terraform Configuration. The field `win32_status` has been replaced by `win32_status_code` (this remains an int, as in 3.63.0 onwards) for `azurerm_linux_web_app`, `azurerm_linux_web_app_slot`, `azurerm_windows_web_app`, `azurerm_windows_web_app_slot resources`. ([#23075](https://github.com/hashicorp/terraform-provider-azurerm/issues/23075))

FEATURES:

* **New Resource**: `azurerm_databricks_workspace_root_dbfs_customer_managed_key` ([#22579](https://github.com/hashicorp/terraform-provider-azurerm/issues/22579))

ENHANCEMENTS:

* dependencies: updating to `v0.20230824.1130652` of `github.com/hashicorp/go-azure-sdk` ([#23076](https://github.com/hashicorp/terraform-provider-azurerm/issues/23076))
* `trafficmanager`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22579](https://github.com/hashicorp/terraform-provider-azurerm/issues/22579))
* `webpubsub`: updating to use the transport layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22579](https://github.com/hashicorp/terraform-provider-azurerm/issues/22579))
* `automation`: upgrade remaining resources to `2022-08-08` ([#22989](https://github.com/hashicorp/terraform-provider-azurerm/issues/22989))
* `azurerm_storage_management_policy` - move to `hashicorp/go-azure-sdk` ([#23035](https://github.com/hashicorp/terraform-provider-azurerm/issues/23035))
* Data Source: `azurerm_disk_encryption_set` -  support for the `identity` block ([#23005](https://github.com/hashicorp/terraform-provider-azurerm/issues/23005))
* `azurerm_container_group` - support for the `sku` and `(init_)container.*.security` properties ([#23034](https://github.com/hashicorp/terraform-provider-azurerm/issues/23034))
* `azurerm_kubernetes_cluster` -  extend allowed ranges for various `sysctl_config` attribute ranges ([#23077](https://github.com/hashicorp/terraform-provider-azurerm/issues/23077))
* `azurerm_kubernetes_cluster_node_pool` -  extend allowed ranges for various `sysctl_config` attribute ranges ([#23077](https://github.com/hashicorp/terraform-provider-azurerm/issues/23077))
* `azurerm_kubernetes_cluster` - clusters can be updated to use the `cilium` dataplane by setting the value in `ebpf_data_plane` ([#22952](https://github.com/hashicorp/terraform-provider-azurerm/issues/22952))
* `azurerm_linux_virtual_machine_scale_set` - cancel rolling upgrades that are in progress before destroying the resource ([#22991](https://github.com/hashicorp/terraform-provider-azurerm/issues/22991))
* `azurerm_servicebus_namespace` - support for `network_rule_set` block ([#23057](https://github.com/hashicorp/terraform-provider-azurerm/issues/23057))
* `azurerm_windows_virtual_machine_scale_set` - cancel rolling upgrades that are in progress before destroying the resource ([#22991](https://github.com/hashicorp/terraform-provider-azurerm/issues/22991))
* `azurerm_synapse_spark_pool` - support addtional values for the `node_size_family` property ([#23040](https://github.com/hashicorp/terraform-provider-azurerm/issues/23040))

BUG FIXES:

* `azurerm_api_management_policy` - fixes an error caused by a migration ([#23018](https://github.com/hashicorp/terraform-provider-azurerm/issues/23018))
* `azurerm_kubernetes_cluster` - deprecate `public_network_access_enabled` and prevent sending it to the API since it isn't functional ([#22478](https://github.com/hashicorp/terraform-provider-azurerm/issues/22478))

## 3.70.0 (August 17, 2023)

FEATURES:

* **New Resource**: `azurerm_mssql_virtual_machine_availability_group_listener` ([#22808](https://github.com/hashicorp/terraform-provider-azurerm/issues/22808))
* **New Resource**: `azurerm_mssql_virtual_machine_group` ([#22808](https://github.com/hashicorp/terraform-provider-azurerm/issues/22808))

ENHANCEMENTS:

* dependencies: updating to `v0.20230815.1165905` of `github.com/hashicorp/go-azure-sdk` ([#22981](https://github.com/hashicorp/terraform-provider-azurerm/issues/22981))
* `apimanagement`: updating to use `hashicorp/go-azure-sdk` ([#22783](https://github.com/hashicorp/terraform-provider-azurerm/issues/22783))
* `cosmos`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22874](https://github.com/hashicorp/terraform-provider-azurerm/issues/22874))
* `devtestlabs`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22981](https://github.com/hashicorp/terraform-provider-azurerm/issues/22981))
* `policy`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22874](https://github.com/hashicorp/terraform-provider-azurerm/issues/22874))
* `postgresql`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22874](https://github.com/hashicorp/terraform-provider-azurerm/issues/22874))
* `recoveryservices`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22874](https://github.com/hashicorp/terraform-provider-azurerm/issues/22874))
* `resources`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22874](https://github.com/hashicorp/terraform-provider-azurerm/issues/22874))
* `storage`: updating Storage Account and Storage Blob Container to use Common IDs to enable migrating to `hashicorp/go-azure-sdk` in the future ([#22915](https://github.com/hashicorp/terraform-provider-azurerm/issues/22915))
* Data Source: `azurerm_kubernetes_cluster` - add support for the `current_kubernetes_version` property ([#22986](https://github.com/hashicorp/terraform-provider-azurerm/issues/22986))
* `azurerm_mssql_virtual_machine` - add support for the `sql_virtual_machine_group_id` and `wsfc_domain_credential` properties ([#22808](https://github.com/hashicorp/terraform-provider-azurerm/issues/22808))
* `azurerm_netapp_pool` - `size_in_tb` can be sized down to 2 TB ([#22943](https://github.com/hashicorp/terraform-provider-azurerm/issues/22943))
* `azurerm_stack_hci_cluster` - add support for the `automanage_configuration_id` property ([#22857](https://github.com/hashicorp/terraform-provider-azurerm/issues/22857))
* Data Source: `azurerm_disk_encryption_set` - now exports `key_vault_key_url` ([#22893](https://github.com/hashicorp/terraform-provider-azurerm/issues/22893))
* `azurerm_disk_encryption_set` - now exports `key_vault_key_url` ([#22893](https://github.com/hashicorp/terraform-provider-azurerm/issues/22893))

BUG FIXES:

* `azurerm_cognitive_deployment` - add lock on parent resource to prevent errors when deleting the resource ([#22940](https://github.com/hashicorp/terraform-provider-azurerm/issues/22940))
* `azurerm_cost_management_scheduled_action` - fix update for `email_address_sender` ([#22930](https://github.com/hashicorp/terraform-provider-azurerm/issues/22930))
* `azurerm_disk_encryption_set` - now correctly supports key rotation by specifying a versionless Key ID when setting `auto_key_rotation_enabled` to `true` ([#22893](https://github.com/hashicorp/terraform-provider-azurerm/issues/22893))
* `azurerm_iothub_dps` - updating the validation for `target` within the `ip_filter_rule` block to match the values defined in the Azure API Definitions ([#22891](https://github.com/hashicorp/terraform-provider-azurerm/issues/22891))
* `azurerm_postgresql_database` - reworking the validation for database collation ([#22928](https://github.com/hashicorp/terraform-provider-azurerm/issues/22928))
* `azurerm_postgresql_flexible_database` - reworking the validation for database collation ([#22928](https://github.com/hashicorp/terraform-provider-azurerm/issues/22928))
* `azurerm_storage_management_policy` - check for an existing resource to prevent overwriting property values ([#22966](https://github.com/hashicorp/terraform-provider-azurerm/issues/22966))
* `azurerm_virtual_network_gateway_connection` - `custom_bgp_addresses.secondary` is now `Optional` rather than `Required` ([#22912](https://github.com/hashicorp/terraform-provider-azurerm/issues/22912))
* `azurerm_web_application_firewall_policy` - fix handling not found in read ([#22982](https://github.com/hashicorp/terraform-provider-azurerm/issues/22982))

---

For information on changes between the v3.69.0 and v3.0.0 releases, please see [the previous v3.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v3.md).

For information on changes between the v2.99.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v1.44.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
