## 3.117.1 (February 28, 2025)

SPECIAL NOTES: This 3.x.x patch release is a special, one-off, back-port of an API upgrade for the `azurerm_kubernetes_cluster_trusted_access_role_binding` resource to enable users still on 3.x to continue using this resource.

BUG FIXES:

* dependencies: `azurerm_kubernetes_cluster_trusted_access_role_binding` - upgrade API to `2024-05-01` ([#28910](https://github.com/hashicorp/terraform-provider-azurerm/pull/28910)) 

## 3.117.0 (November 7, 2024)

SPECIAL NOTES: This 3.x release is a special, one-off, back-port of functionality for `azurerm_storage_account` to enable users to deploy this resource in environments which block / are restrictive of Data Plane access, thus preventing the resource being created and/or managed.  This functionality is back-ported from the `v4.9.0` release. Users migrating from this release to the 4.x line, should upgrade directly to `v4.9.0` or later, as these features are not compatible with earlier releases of 4.x.

FEATURES:

* **New Resource:** `azurerm_storage_account_queue_properties` ([#27819](https://github.com/hashicorp/terraform-provider-azurerm/pull/27819))
* **New Resource:** `azurerm_storage_account_static_website`  ([#27819](https://github.com/hashicorp/terraform-provider-azurerm/pull/27819))
* New Provider Feature - storage `data_plane_available` feature flag ([#27819](https://github.com/hashicorp/terraform-provider-azurerm/pull/27819))

ENHANCEMENTS:

* `azurerm_storage_account` - can now be created and managed if Data Plane endpoints are blocked by a firewall ([#27819](https://github.com/hashicorp/terraform-provider-azurerm/pull/27819)) 

## 3.116.0 (August 16, 2024)

DEPRECATIONS:

All Azure Kubernetes Service (AKS) properties related to preview features are deprecated since they will not be available in a stable API. Please see https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/4.0-upgrade-guide#aks-migration-to-stable-api for more details ([#26863](https://github.com/hashicorp/terraform-provider-azurerm/issues/26863))

FEATURES:

* New Resource: `azurerm_ai_services` ([#26008](https://github.com/hashicorp/terraform-provider-azurerm/issues/26008))
* New Resource: `azurerm_communication_service_email_domain_association` ([#26432](https://github.com/hashicorp/terraform-provider-azurerm/issues/26432))
* New Resource: `azurerm_dev_center_project_environment_type` ([#26941](https://github.com/hashicorp/terraform-provider-azurerm/issues/26941))
* New Resource: `azurerm_extended_location_custom_location` ([#24267](https://github.com/hashicorp/terraform-provider-azurerm/issues/24267))
* New Resource: `azurerm_postgresql_flexible_server_virtual_endpoint` ([#26708](https://github.com/hashicorp/terraform-provider-azurerm/issues/26708))

ENHANCEMENTS:

* `notificationhub` - updating to use version `2023-09-01` ([#26528](https://github.com/hashicorp/terraform-provider-azurerm/issues/26528))
* `azurerm_api_management_api` - update validation of `path` to allow single character strings ([#26922](https://github.com/hashicorp/terraform-provider-azurerm/issues/26922))
* `azurerm_cosmosdb_account` - add support for the property `burst_capacity_enabled` ([#26986](https://github.com/hashicorp/terraform-provider-azurerm/issues/26986))
* `azurerm_linux_function_app` - add support for `vnet_image_pull_enabled` property in 4.0 ([#27001](https://github.com/hashicorp/terraform-provider-azurerm/issues/27001))
* `azurerm_linux_function_app_slot` - add support for `vnet_image_pull_enabled` property in 4.0 ([#27001](https://github.com/hashicorp/terraform-provider-azurerm/issues/27001))
* `azurerm_logic_app_standard` - add support for `v8.0` in `site_config.dotnet_framework_version` ([#26983](https://github.com/hashicorp/terraform-provider-azurerm/issues/26983))
* `azurerm_management_group_policy_assignment` - remove length restriction on name ([#27055](https://github.com/hashicorp/terraform-provider-azurerm/issues/27055))
* `azurerm_recovery_services_vault` - add support for the `identity` block ([#26254](https://github.com/hashicorp/terraform-provider-azurerm/issues/26254))
* `azurerm_web_application_firewall_policy` - add support for the `js_challenge_cookie_expiration_in_minutes` property ([#26878](https://github.com/hashicorp/terraform-provider-azurerm/issues/26878))
* `azurerm_windows_function_app` - add support for `vnet_image_pull_enabled` property in 4.0 ([#27001](https://github.com/hashicorp/terraform-provider-azurerm/issues/27001))
* `azurerm_windows_function_app_slot` - add support for `vnet_image_pull_enabled` property in 4.0 ([#27001](https://github.com/hashicorp/terraform-provider-azurerm/issues/27001))

BUG FIXES:

* Data Source: `azurerm_storage_account` - add `default_share_level_permission` to the `azure_files_authentication` to prevent invalid address errors ([#26996](https://github.com/hashicorp/terraform-provider-azurerm/issues/26996))
* Data Source: `azurerm_search_service` - expose the `tags` property ([#26978](https://github.com/hashicorp/terraform-provider-azurerm/issues/26978))
* Data Source: `azurerm_virtual_machine` - populate missing `power_state` ([#26991](https://github.com/hashicorp/terraform-provider-azurerm/issues/26991))
* Data Source: `azurerm_virtual_machine_scale_set` - populate missing `power_state` ([#26991](https://github.com/hashicorp/terraform-provider-azurerm/issues/26991))
* `azurerm_api_management_api_schema` - correctly unmarshal `definition` and `components` ([#26531](https://github.com/hashicorp/terraform-provider-azurerm/issues/26531))
* `azurerm_cdn_frontdoor_secret` - fix issue where `expiration_date` was being set into the parent block ([#26982](https://github.com/hashicorp/terraform-provider-azurerm/issues/26982))
* `azurerm_container_app_environment` - fix diff suppress on `infrastructure_resource_group_name` ([#27007](https://github.com/hashicorp/terraform-provider-azurerm/issues/27007))
* `azurerm_express_route_connection` - prevent sending `private_link_fast_path_enabled` in the payload if it hasn't been explicitly set ([#26928](https://github.com/hashicorp/terraform-provider-azurerm/issues/26928))
* `azurerm_machine_learning_workspace` - `serverless_compute` can now be updated ([#26940](https://github.com/hashicorp/terraform-provider-azurerm/issues/26940))
* `azurerm_mssql_database` - fix issue where the database cannot be upgraded to use serverless due to the behaviour of the `license_type` field ([#26850](https://github.com/hashicorp/terraform-provider-azurerm/issues/26850))
* `azurerm_mssql_database` - prevent error when creating `Free` edition by setting `long_term_retention_policy` and `short_term_retention_policy` as empty ([#26894](https://github.com/hashicorp/terraform-provider-azurerm/issues/26894))
* `azurerm_nginx_deployment` - omit `capacity` when creating deployments with a basic plan ([#26223](https://github.com/hashicorp/terraform-provider-azurerm/issues/26223))
* `azurerm_role_management_policy` - prevent panic when updating `activation_rules.approval_stage` ([#26800](https://github.com/hashicorp/terraform-provider-azurerm/issues/26800))
* `azurerm_sentinel_threat_intelligence_indicator` - prevent panic when importing this resource ([#26976](https://github.com/hashicorp/terraform-provider-azurerm/issues/26976))
* `azurerm_servicebus_namespace` - fix panic reading encryption with versionless ids ([#27060](https://github.com/hashicorp/terraform-provider-azurerm/issues/27060))
* `azurerm_synapse_spark_pool` - prevent plan diff due to API behaviour by setting `node_count` as Computed ([#26953](https://github.com/hashicorp/terraform-provider-azurerm/issues/26953))
* `azurerm_virtual_network_gateway_connection` - fix issue where `ingress_nat_rule_ids` was updating the egress rules on updates ([#27022](https://github.com/hashicorp/terraform-provider-azurerm/issues/27022))

## 3.115.0 (August 09, 2024)

ENHANCEMENTS:

* `cosmosdb` - updating to use version `2024-05-15` ([#26758](https://github.com/hashicorp/terraform-provider-azurerm/issues/26758))
* `healthcare` - updating to use version `2024-03-31` ([#26699](https://github.com/hashicorp/terraform-provider-azurerm/issues/26699))
* `redis` - updating to use version `2024-03-01` ([#26932](https://github.com/hashicorp/terraform-provider-azurerm/issues/26932))
* `azurerm_cosmosdb_account` - avoid infinite diff to `default_identity_type` for legacy resources where an empty string is returned by the RP ([#26525](https://github.com/hashicorp/terraform-provider-azurerm/issues/26525))
* `azurerm_linux_virtual_machine_scale_set` - add support for the `action` property in the `automatic_instance_repair` block ([#26227](https://github.com/hashicorp/terraform-provider-azurerm/issues/26227))
* `azurerm_log_analytics_saved_search` - update the regex for the `function_parameters` property to support more paramters ([#26701](https://github.com/hashicorp/terraform-provider-azurerm/issues/26701))
* `azurerm_monitor_data_collection_rule` - update `performance_counter.x.sampling_frequency_in_seconds` range `1` to `1800` ([#26898](https://github.com/hashicorp/terraform-provider-azurerm/issues/26898))
* `azurerm_orchestrated_virtual_machine_scale_set` - add support for the `action` property in the `automatic_instance_repair` block ([#26227](https://github.com/hashicorp/terraform-provider-azurerm/issues/26227))
* `azurerm_security_center_storage_defender` - add support for the property `scan_results_event_grid_topic_id` ([#26599](https://github.com/hashicorp/terraform-provider-azurerm/issues/26599))
* `azurerm_storage_account` - add support for the property `default_share_level_permission` in the `azure_files_authentication` block ([#26924](https://github.com/hashicorp/terraform-provider-azurerm/issues/26924))
* `azurerm_web_application_firewall_policy` - `excluded_rule_set.0.type` supports `Microsoft_BotManagerRuleSet` ([#26903](https://github.com/hashicorp/terraform-provider-azurerm/issues/26903))
* `azurerm_windows_virtual_machine_scale_set` - add support for the `action` property in the `automatic_instance_repair` block ([#26227](https://github.com/hashicorp/terraform-provider-azurerm/issues/26227))

BUG FIXES:

* `azurerm_container_group` - retrieve and set `storage_account_key` in the payload when updating the resource ([#26640](https://github.com/hashicorp/terraform-provider-azurerm/issues/26640))
* `azurerm_key_vault_managed_hardware_security_module_role_assignment` - fixed a crash in error messages ([#26972](https://github.com/hashicorp/terraform-provider-azurerm/issues/26972))
* `azurerm_kubernetes_cluster` - allow an empty list for `dns_zone_ids` in the `web_app_routing` block ([#26747](https://github.com/hashicorp/terraform-provider-azurerm/issues/26747))
* `azurerm_storage_share_file` - fix a bug when encoding the MD5 hash for the `content_md5` property ([#25715](https://github.com/hashicorp/terraform-provider-azurerm/issues/25715))

## 3.114.0 (August 01, 2024)

UPGRADE NOTES:
* **4.0 Beta:** This release includes a new feature-flag to opt-into the 4.0 Beta - which (when enabled) introduces a number of behavioural changes, field renames and removes some older deprecated resources and data sources. Please read the disclaimers carefully that are outlined in our [guide on how to opt-into the 4.0 Beta](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/4.0-beta) before enabling this, as this will cause irreversible changes to your state. The 4.0 Beta is still a work-in-progress at this time and the changes listed in the [4.0 Upgrade Guide](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/4.0-upgrade-guide) may change. We're interested to hear your feedback which can be provided by following [this link](https://github.com/terraform-providers/terraform-provider-azurerm/issues/new?template=Beta_Feedback.md).

FEATURES:

* **New Resource:** `azurerm_dev_center_network_connection` ([#26718](https://github.com/hashicorp/terraform-provider-azurerm/issues/26718))
* **New Resource:** `azurerm_stack_hci_logical_network` ([#26473](https://github.com/hashicorp/terraform-provider-azurerm/issues/26473))

ENHANCEMENTS:

* dependencies: updating `go-azure-helpers` to `v0.70.1` ([#26757](https://github.com/hashicorp/terraform-provider-azurerm/issues/26757))
* `arckubernetes` - updating to use version `2024-01-01` ([#26761](https://github.com/hashicorp/terraform-provider-azurerm/issues/26761))
* `data.azurerm_storage_account` - the `enable_https_traffic_only` property has been superseded by `https_traffic_only_enabled` ([#26740](https://github.com/hashicorp/terraform-provider-azurerm/issues/26740))
* `azurerm_log_analytics_cluster` - add support for setting `size_gb` to `100` [GH-#26865]
* `azurerm_storage_account` - the `enable_https_traffic_only` property has been superseded by `https_traffic_only_enabled` ([#26740](https://github.com/hashicorp/terraform-provider-azurerm/issues/26740))

BUG FIXES:

* `azurerm_dns_cname_record` - split create and update function to fix lifecycle - ignore ([#26610](https://github.com/hashicorp/terraform-provider-azurerm/issues/26610))
* `azurerm_dns_srv_record` - split create and update function to fix lifecycle - ignore ([#26627](https://github.com/hashicorp/terraform-provider-azurerm/issues/26627))
* `azurerm_kubernetes_cluster` - fix issue that prevented `max_count` from being updated ([#26417](https://github.com/hashicorp/terraform-provider-azurerm/issues/26417))
* `azurerm_linux_web_app` - correctly set `site_config.always_on` as configured during Update ([#25753](https://github.com/hashicorp/terraform-provider-azurerm/issues/25753))
* `azurerm_linux_web_app_slot` - correctly set `site_config.always_on` as configured during Update ([#25753](https://github.com/hashicorp/terraform-provider-azurerm/issues/25753))
* `azurerm_management_group_policy_remediation` - fix panic in deprecated schema change for 4.0 ([#26767](https://github.com/hashicorp/terraform-provider-azurerm/issues/26767))
* `azurerm_network_security_rule` - fix panic when updating `source_port_ranges` ([#26883](https://github.com/hashicorp/terraform-provider-azurerm/issues/26883))
* `azurerm_public_ip` - fix panix when updating `idle_timeout_in_minutes`

DEPRECATIONS:
* `azurerm_redis_cache` - `enable_non_ssl_port` has been superseded by `non_ssl_port_enabled` and `redis_configuration. enable_authentication` has been superseded by `redis_configuration.authentication_enabled` ([#26608](https://github.com/hashicorp/terraform-provider-azurerm/issues/26608))


## 3.113.0 (July 18, 2024)

ENHANCEMENTS:

* dependencies: updating to `v0.20240715.1100358` of `hashicorp/go-azure-sdk` ([#26638](https://github.com/hashicorp/terraform-provider-azurerm/issues/26638))
* `storage` - updating to use `hashicorp/go-azure-sdk` ([#26218](https://github.com/hashicorp/terraform-provider-azurerm/issues/26218))

BUG FIXES:

* `azurerm_storage_account` - fix a validation bug when replacing a StorageV2 account with a StorageV1 account ([#26639](https://github.com/hashicorp/terraform-provider-azurerm/issues/26639))
* `azurerm_storage_account` - resolve an issue refreshing blob or queue properties after recreation ([#26218](https://github.com/hashicorp/terraform-provider-azurerm/issues/26218))
* `azurerm_storage_account` - resolve an issue setting tags for an existing storage account where a policy mandates them ([#26218](https://github.com/hashicorp/terraform-provider-azurerm/issues/26218))
* `azurerm_storage_account` - fix a persistent diff with the `customer_managed_key` block ([#26218](https://github.com/hashicorp/terraform-provider-azurerm/issues/26218))
* `azurerm_storage_account` - resolve several consistency related issues when crreating a new storage account ([#26218](https://github.com/hashicorp/terraform-provider-azurerm/issues/26218))

DEPRECATIONS:

* `azurerm_eventhub_namespace` - deprecate the `zone_redundant` field in v4.0 ([#26611](https://github.com/hashicorp/terraform-provider-azurerm/issues/26611))
* `azurerm_servicebus_namespace` - deprecate the `zone_redundant` field in v4.0 ([#26611](https://github.com/hashicorp/terraform-provider-azurerm/issues/26611))

## 3.112.0 (July 12, 2024)

FEATURES:

* New Data Source: `azurerm_elastic_san_volume_snapshot` ([#26439](https://github.com/hashicorp/terraform-provider-azurerm/issues/26439))
* New Resource: `azurerm_dev_center_dev_box_definition` ([#26307](https://github.com/hashicorp/terraform-provider-azurerm/issues/26307))
* New Resource: `azurerm_dev_center_environment_type` ([#26291](https://github.com/hashicorp/terraform-provider-azurerm/issues/26291))
* New Resource: `azurerm_virtual_machine_restore_point` ([#26526](https://github.com/hashicorp/terraform-provider-azurerm/issues/26526))
* New Resource: `azurerm_virtual_machine_restore_point_collection` ([#26526](https://github.com/hashicorp/terraform-provider-azurerm/issues/26526))

ENHANCEMENTS:

* dependencies: updating to `v0.20240710.1114656` of `github.com/hashicorp/go-azure-sdk` ([#26588](https://github.com/hashicorp/terraform-provider-azurerm/issues/26588))
* dependencies: updating to `v0.70.0` of `go-azure-helpers` ([#26601](https://github.com/hashicorp/terraform-provider-azurerm/issues/26601))
* `containerservice`: updating the Fleet resources to use API Version `2024-04-01` ([#26588](https://github.com/hashicorp/terraform-provider-azurerm/issues/26588))
* Data Source: `azurerm_network_service_tags` - extend validation for `service` to allow `AzureFrontDoor.Backend`, `AzureFrontDoor.Frontend`, and `AzureFrontDoor.FirstParty` ([#26429](https://github.com/hashicorp/terraform-provider-azurerm/issues/26429))
* `azurerm_api_management_identity_provider_aad` - support for the `client_library` property ([#26093](https://github.com/hashicorp/terraform-provider-azurerm/issues/26093))
* `azurerm_api_management_identity_provider_aadb2c` - support for the `client_library` property ([#26093](https://github.com/hashicorp/terraform-provider-azurerm/issues/26093))
* `azurerm_dev_test_virtual_network` - support for the `shared_public_ip_address` property ([#26299](https://github.com/hashicorp/terraform-provider-azurerm/issues/26299))
* `azurerm_kubernetes_cluster` - support for the `certificate_authority` block under the `service_mesh_profile` block ([#26543](https://github.com/hashicorp/terraform-provider-azurerm/issues/26543))
* `azurerm_linux_web_app` - support the value `8.3` for the `php_version` property ([#26194](https://github.com/hashicorp/terraform-provider-azurerm/issues/26194))
* `azurerm_machine_learning_compute_cluster` - the `identity` property can now be updated ([#26404](https://github.com/hashicorp/terraform-provider-azurerm/issues/26404))
* `azurerm_web_application_firewall_policy` - support for the `JSChallenge` value for `managed_rules.managed_rule_set.rule_group_override.rule_action` ([#26561](https://github.com/hashicorp/terraform-provider-azurerm/issues/26561))

BUG FIXES:

* Data Source: `azurerm_communication_service` - `primary_connection_string`, `primary_key`, `secondary_connection_string` and `secondary_key` are marked as Sensitive ([#26560](https://github.com/hashicorp/terraform-provider-azurerm/issues/26560))
* `azurerm_app_configuration_feature` - fix issue when updating the resource without an existing `targeting_filter` ([#26506](https://github.com/hashicorp/terraform-provider-azurerm/issues/26506))
* `azurerm_backup_policy_vm` - split create and update function to fix lifecycle - ignore ([#26591](https://github.com/hashicorp/terraform-provider-azurerm/issues/26591))
* `azurerm_backup_protected_vm` - split create and update function to fix lifecycle - ignore ([#26583](https://github.com/hashicorp/terraform-provider-azurerm/issues/26583))
* `azurerm_communication_service` - the `primary_connection_string`, `primary_key`, `secondary_connection_string`, and `secondary_key` properties are now sensitive ([#26560](https://github.com/hashicorp/terraform-provider-azurerm/issues/26560))
* `azurerm_mysql_flexible_server_configuration` - add locks to prevent conflicts when deleting the resource ([#26289](https://github.com/hashicorp/terraform-provider-azurerm/issues/26289))
* `azurerm_nginx_deployment` - changing the `frontend_public.ip_address`, `frontend_private.ip_address`, `frontend_private.allocation_method`, and `frontend_private.subnet_id` now creates a new resource ([#26298](https://github.com/hashicorp/terraform-provider-azurerm/issues/26298))
* `azurerm_palo_alto_local_rulestack_rule` - correctl read the `protocol` property on read when the `protocol_ports` property is configured ([#26510](https://github.com/hashicorp/terraform-provider-azurerm/issues/26510))
* `azurerm_servicebus_namespace` - parse the identity returned by the API insensitively before setting into state ([#26540](https://github.com/hashicorp/terraform-provider-azurerm/issues/26540))

DEPRECATIONS:

* `azurerm_servicebus_queue` - `enable_batched_operations`, `enable_express` and `enable_partitioning` are superseded by `batched_operations_enabled`, `express_enabled` and `partitioning_enabled` ([#26479](https://github.com/hashicorp/terraform-provider-azurerm/issues/26479))
* `azurerm_servicebus_subscription` - `enable_batched_operations` has been superseded  by `batched_operations_enabled` ([#26479](https://github.com/hashicorp/terraform-provider-azurerm/issues/26479))
* `azurerm_servicebus_topic` - `enable_batched_operations`, `enable_express` and `enable_partitioning` are superseded by `batched_operations_enabled`, `express_enabled` and `partitioning_enabled` ([#26479](https://github.com/hashicorp/terraform-provider-azurerm/issues/26479))

## 3.111.0 (July 04, 2024)

FEATURES:

* **New Resource:** `azurerm_restore_point_collection` ([#26518](https://github.com/hashicorp/terraform-provider-azurerm/issues/26518))

ENHANCEMENTS:

* dependencies: updating to `v0.20240701.1082110` of `github.com/hashicorp/go-azure-sdk` ([#26502](https://github.com/hashicorp/terraform-provider-azurerm/issues/26502))
* `azurerm_disk_encryption_set` - support for the `managed_hsm_key_id` property ([#26201](https://github.com/hashicorp/terraform-provider-azurerm/issues/26201))
* `azurerm_firewall_policy` - remove Computed from the `sku` property and add a default of `Standard` in 4.0 ([#26499](https://github.com/hashicorp/terraform-provider-azurerm/issues/26499))
* `azurerm_kubernetes_cluster` - support updating `default_node_pool.os_sku` between `Ubuntu` and `AzureLinux` ([#26262](https://github.com/hashicorp/terraform-provider-azurerm/issues/26262))
* `azurerm_kubernetes_cluster_node_pool` - support updating `os_sku` between `Ubuntu` and `AzureLinux` ([#26139](https://github.com/hashicorp/terraform-provider-azurerm/issues/26139))
* `azurerm_service_plan` - support for new the Flex Consumption plan ([#26351](https://github.com/hashicorp/terraform-provider-azurerm/issues/26351))

BUG FIXES:

* `azurerm_kubernetes_cluster` - prevent a panic ([#26478](https://github.com/hashicorp/terraform-provider-azurerm/issues/26478))
* `azurerm_kubernetes_cluster` - prevent a diff in `upgrade_settings` when the API returns an empty object ([#26541](https://github.com/hashicorp/terraform-provider-azurerm/issues/26541))
* `azurerm_kubernetes_cluster_node_pool` - prevent a diff in `upgrade_settings` when the API returns an empty object ([#26541](https://github.com/hashicorp/terraform-provider-azurerm/issues/26541))
* `azurerm_virtual_network_gateway` - split create and update function to fix lifecycle - ignore ([#26451](https://github.com/hashicorp/terraform-provider-azurerm/issues/26451))
* `azurerm_virtual_network_gateway_connection` - split create and update function to fix lifecycle - ignore ([#26431](https://github.com/hashicorp/terraform-provider-azurerm/issues/26431))

## 3.110.0 (June 27, 2024)

FEATURES:

* **New Data Source:** `azurerm_load_test` ([#26376](https://github.com/hashicorp/terraform-provider-azurerm/issues/26376))
* **New Resource:** `azurerm_virtual_desktop_scaling_plan_host_pool_association` ([#24670](https://github.com/hashicorp/terraform-provider-azurerm/issues/24670))

ENHANCEMENTS:

* Data Source: `azurerm_monitor_data_collection_endpoint` - support for the `immutable_id` property ([#26380](https://github.com/hashicorp/terraform-provider-azurerm/issues/26380))
* Data Source: `azurerm_nginx_certificate` - export the properties `sha1_thumbprint`, `key_vault_secret_version`, `key_vault_secret_creation_date`, `error_code` and `error_message` ([#26160](https://github.com/hashicorp/terraform-provider-azurerm/issues/26160))
* `azurerm_backup_policy_vm` - support for the `tiering_policy` property ([#26263](https://github.com/hashicorp/terraform-provider-azurerm/issues/26263))
* `azurerm_kubernetes_cluster_node_pool` - Pod Disruption Budgets are now respected when deleting a node pool ([#26471](https://github.com/hashicorp/terraform-provider-azurerm/issues/26471))
* `azurerm_monitor_data_collection_endpoint` - support for the `immutable_id` property ([#26380](https://github.com/hashicorp/terraform-provider-azurerm/issues/26380))
* `azurerm_mssql_managed_instance` - support the value `GZRS` for the `storage_account_type` property ([#26448](https://github.com/hashicorp/terraform-provider-azurerm/issues/26448))
* `azurerm_mssql_managed_instance_transparent_data_encryption` - support for the `managed_hsm_key_id` property ([#26496](https://github.com/hashicorp/terraform-provider-azurerm/issues/26496))
* `azurerm_redis_cache_access_policy` - allow updates to `permissions` ([#26440](https://github.com/hashicorp/terraform-provider-azurerm/issues/26440))
* `azurerm_redhat_openshift_cluster` - support for the `managed_resource_group_name` property ([#25529](https://github.com/hashicorp/terraform-provider-azurerm/issues/25529))
* `azurerm_redhat_openshift_cluster` - support for the `preconfigured_network_security_group_enabled` property ([#26082](https://github.com/hashicorp/terraform-provider-azurerm/issues/26082))
* `azurerm_iotcentral_application` - remove Computed from `template` and set default of `iotc-pnp-preview@1.0.0` in 4.0  ([#26485](https://github.com/hashicorp/terraform-provider-azurerm/issues/26485))
* `azurerm_digital_twins_time_series_database_connection` - remove Computed from `kusto_table_name` and set a default of `AdtPropertyEvents` in 4.0 ([#26484](https://github.com/hashicorp/terraform-provider-azurerm/issues/26484))

BUG FIXES:

* Data Source: `azurerm_express_route_circuit_peering` - fix issue where data source attempts to parse an empty string instead of generating the resource ID ([#26441](https://github.com/hashicorp/terraform-provider-azurerm/issues/26441))
* `azurerm_express_route_gateway` - prevent a panic ([#26467](https://github.com/hashicorp/terraform-provider-azurerm/issues/26467))
* `azurerm_monitor_scheduled_query_rules_alert_v2` - correctly handle the `identity` block if not specified ([#26364](https://github.com/hashicorp/terraform-provider-azurerm/issues/26364))
* `azurerm_security_center_automation` - prevent resource recreation when `tags` are updated ([#26292](https://github.com/hashicorp/terraform-provider-azurerm/issues/26292))
* `azurerm_synapse_workspace` - fix issue where `azure_devops_repo` or `github_repo` configuration could not be removed ([#26421](https://github.com/hashicorp/terraform-provider-azurerm/issues/26421))
* `azurerm_virtual_network_dns_servers` - split create and update function to fix lifecycle - ignore ([#26427](https://github.com/hashicorp/terraform-provider-azurerm/issues/26427))
* `azurerm_linux_function_app` - set `allowed_applications` in the request payload ([#26462](https://github.com/hashicorp/terraform-provider-azurerm/issues/26462))
* `azurerm_linux_function_app_slot` - set `allowed_applications` in the request payload ([#26462](https://github.com/hashicorp/terraform-provider-azurerm/issues/26462))
* `azurerm_windows_function_app` - set `allowed_applications` in the request payload ([#26462](https://github.com/hashicorp/terraform-provider-azurerm/issues/26462))
* `azurerm_windows_function_app_slot` - set `allowed_applications` in the request payload ([#26462](https://github.com/hashicorp/terraform-provider-azurerm/issues/26462))
* `azurerm_linux_web_app` - set `allowed_applications` in the request payload ([#26462](https://github.com/hashicorp/terraform-provider-azurerm/issues/26462))
* `azurerm_linux_web_app_slot` - set `allowed_applications` in the request payload ([#26462](https://github.com/hashicorp/terraform-provider-azurerm/issues/26462))
* `azurerm_windows_web_app` - set `allowed_applications` in the request payload ([#26462](https://github.com/hashicorp/terraform-provider-azurerm/issues/26462))
* `azurerm_windows_web_app_slot` - set `allowed_applications` in the request payload ([#26462](https://github.com/hashicorp/terraform-provider-azurerm/issues/26462))
* `azurerm_api_management` - remove ForceNew from `additional_location.zones` ([#26384](https://github.com/hashicorp/terraform-provider-azurerm/issues/26384))
* `azurerm_logic_app_integration_account_schema` - the `name` property now allows underscores ([#26475](https://github.com/hashicorp/terraform-provider-azurerm/issues/26475))
* `azurerm_palo_alto_local_rulestack_rule` - prevent error when switching between `protocol` and `protocol_ports` ([#26490](https://github.com/hashicorp/terraform-provider-azurerm/issues/26490))

DEPRECATIONS:

* `azurerm_analysis_service_server` - the property `enable_power_bi_service` has been superseded by `power_bi_service_enabled` ([#26456](https://github.com/hashicorp/terraform-provider-azurerm/issues/26456))

## 3.109.0 (June 20, 2024)

FEATURES:

* **New Data Source:** `azurerm_automation_runbook` ([#26359](https://github.com/hashicorp/terraform-provider-azurerm/issues/26359))
* **New Resource:** `azurerm_data_protection_backup_instance_postgresql_flexible_server` ([#26249](https://github.com/hashicorp/terraform-provider-azurerm/issues/26249))
* **New Resource:** `azurerm_email_communication_service_domain` ([#26179](https://github.com/hashicorp/terraform-provider-azurerm/issues/26179))
* **New Resource:** `azurerm_system_center_virtual_machine_manager_cloud` ([#25429](https://github.com/hashicorp/terraform-provider-azurerm/issues/25429))
* **New Resource:** `azurerm_system_center_virtual_machine_manager_virtual_machine_template` ([#25449](https://github.com/hashicorp/terraform-provider-azurerm/issues/25449))
* **New Resource:** `azurerm_system_center_virtual_machine_manager_virtual_network` ([#25451](https://github.com/hashicorp/terraform-provider-azurerm/issues/25451))

ENHANCEMENTS:

* Data Source: `azurerm_hdinsight_cluster` - export the `cluster_id` attribute ([#26228](https://github.com/hashicorp/terraform-provider-azurerm/issues/26228))
* `azurerm_cosmosdb_sql_container` - support for the `partition_key_kind` and `partition_key_paths` properties ([#26372](https://github.com/hashicorp/terraform-provider-azurerm/issues/26372))
* `azurerm_data_protection_backup_instance_blob_storage` - support for the `storage_account_container_names` property ([#26232](https://github.com/hashicorp/terraform-provider-azurerm/issues/26232))
* `azurerm_virtual_network_peering` - support for the `peer_complete_virtual_networks_enabled`, `only_ipv6_peering_enabled`, `local_subnet_names`, and `remote_subnet_names` properties ([#26229](https://github.com/hashicorp/terraform-provider-azurerm/issues/26229))
* `azurerm_virtual_desktop_host_pool` - changing the `preferred_app_group_type` property no longer creates a new resource ([#26333](https://github.com/hashicorp/terraform-provider-azurerm/issues/26333))
* `azurerm_maps_account` - support for the `location`, `identity`, `cors` and `data_store` properties ([#26397](https://github.com/hashicorp/terraform-provider-azurerm/issues/26397))

BUG FIXES:

* `azurerm_automation_job_schedule` - updates `azurerm_automation_job_schedule` to use a composite resource id and allows `azurerm_automation_runbook` to be updated without causing `azurerm_automation_job_schedule` to recreate ([#22164](https://github.com/hashicorp/terraform-provider-azurerm/issues/22164))
* `azurerm_databricks_workspace`- correctly allow disabling the default firewall ([#26339](https://github.com/hashicorp/terraform-provider-azurerm/issues/26339))
* `azurerm_virtual_hub_*` - spliting create and update so lifecycle ignore changes works correctly ([#26310](https://github.com/hashicorp/terraform-provider-azurerm/issues/26310))

DEPRECATIONS:

* Data Source: `azurerm_mariadb_server` - deprecated since the service is retiring. Please use `azurerm_mysql_flexible_server` instead ([#26354](https://github.com/hashicorp/terraform-provider-azurerm/issues/26354))
* `azurerm_mariadb_configuration` - deprecated since the service is retiring. Please use `azurerm_mysql_flexible_server_configuration` instead ([#26354](https://github.com/hashicorp/terraform-provider-azurerm/issues/26354))
* `azurerm_mariadb_database` - deprecated since the service is retiring. Please use `azurerm_mysql_flexible_database` instead ([#26354](https://github.com/hashicorp/terraform-provider-azurerm/issues/26354))
* `azurerm_mariadb_firewall_rule` - deprecated since the service is retiring. Please use `azurerm_mysql_flexible_server_firewall_rule` instead ([#26354](https://github.com/hashicorp/terraform-provider-azurerm/issues/26354))
* `azurerm_mariadb_server` - deprecated since the service is retiring. Please use `azurerm_mysql_flexible_server` instead ([#26354](https://github.com/hashicorp/terraform-provider-azurerm/issues/26354))
* `azurerm_mariadb_virtual_network_rule` - deprecated since the service is retiring ([#26354](https://github.com/hashicorp/terraform-provider-azurerm/issues/26354))

## 3.108.0 (June 13, 2024)

FEATURES:

* **New Data Source:** `azurerm_role_management_policy` ([#25900](https://github.com/hashicorp/terraform-provider-azurerm/issues/25900))
* **New Resource:** `azurerm_role_management_policy` ([#25900](https://github.com/hashicorp/terraform-provider-azurerm/issues/25900))

ENHANCEMENTS:

* provider: support subscription ID hinting when using Azure CLI authentication ([#26282](https://github.com/hashicorp/terraform-provider-azurerm/issues/26282))
* `serviceconnector`: updating to use API Version `2024-04-01` ([#26248](https://github.com/hashicorp/terraform-provider-azurerm/issues/26248))
* `azurerm_container_groups` - can now be created with a User Assigned Identity when running Windows ([#26308](https://github.com/hashicorp/terraform-provider-azurerm/issues/26308))
* `azurerm_kubernetes_cluster` - updating the `network_profile.network_policy` property to `azure` and `calico` when it hasn't been previously set is supported ([#26176](https://github.com/hashicorp/terraform-provider-azurerm/issues/26176))
* `azurerm_kubernetes_cluster` - respect Pod Distruption Budgets when rotating the `default_node_pool` ([#26274](https://github.com/hashicorp/terraform-provider-azurerm/issues/26274))
* `azurerm_lb_backend_address_pool` - support for the `synchronous_mode` property ([#26309](https://github.com/hashicorp/terraform-provider-azurerm/issues/26309))
* `azurerm_private_endpoint` - support symultaneous creation of multiple resources of this type per subnet ([#26006](https://github.com/hashicorp/terraform-provider-azurerm/issues/26006))

BUG FIXES:

* `azurerm_express_route_circuit_peering`, `azurerm_express_route_circuit`, `azurerm_express_route_gateway`, `azurerm_express_route_port` - split create and update ([#26237](https://github.com/hashicorp/terraform-provider-azurerm/issues/26237))
* `azurerm_lb_backend_address_pool_address` - when using this resource, values are no longer reset on `azurerm_lb_backend_address_pool` ([#26264](https://github.com/hashicorp/terraform-provider-azurerm/issues/26264))
* `azurerm_route_filter` - spliting create and update so lifecycle ignore changes works correctly ([#26266](https://github.com/hashicorp/terraform-provider-azurerm/issues/26266))
* `azurerm_route_server` - spliting create and update so lifecycle ignore changes works correctly ([#26266](https://github.com/hashicorp/terraform-provider-azurerm/issues/26266))
* `azurerm_synapse_workspace` - updates the client used in all operations of `azurerm_synapse_workspace_sql_aad_admin` to prevent this resource from modifying the same resource as `azurerm_synapse_workspace_aad_admin` ([#26317](https://github.com/hashicorp/terraform-provider-azurerm/issues/26317))
* `azurerm_virtual_network` - correctly parse network securty group IDs ([#26283](https://github.com/hashicorp/terraform-provider-azurerm/issues/26283))

DEPRECATIONS:

* Data Source: `azurerm_network_interface` - the `enable_ip_forwarding` and `enable_accelerated_networking` properties have been deprecated and superseded by the `ip_forwarding_enabled` and `accelerated_networking_enabled` properties ([#26293](https://github.com/hashicorp/terraform-provider-azurerm/issues/26293))
* `azurerm_api_management` - the `policy` block has been deprecated is superseded by the `azurerm_api_management_policy` resource ([#26305](https://github.com/hashicorp/terraform-provider-azurerm/issues/26305))
* `azurerm_kubernetes_cluster` - the `ebpf_data_plane` property has been deprecated and superseded by the `network_data_plane` property ([#26251](https://github.com/hashicorp/terraform-provider-azurerm/issues/26251))
* `azurerm_network_interface` - the `enable_ip_forwarding` and `enable_accelerated_networking` properties have been deprecated and superseded by the `ip_forwarding_enabled` and `accelerated_networking_enabled` properties ([#26293](https://github.com/hashicorp/terraform-provider-azurerm/issues/26293))
* `azurerm_synapse_workspace` - the `aad_admin` and `sql_aad_admin` blocks have been deprecated and superseded by the `azurerm_synapse_workspace_aad_admin` and `azurerm_synapse_workspace_sql_aad_admin` resources ([#26317](https://github.com/hashicorp/terraform-provider-azurerm/issues/26317))

## 3.107.0 (June 06, 2024)

FEATURES:

* **New Resource:** `azurerm_data_protection_backup_policy_postgresql_flexible_server` ([#26024](https://github.com/hashicorp/terraform-provider-azurerm/issues/26024))

ENHANCEMENTS:

* dependencies: updating to `v0.20240604.1114748` of `github.com/hashicorp/go-azure-sdk` ([#26216](https://github.com/hashicorp/terraform-provider-azurerm/issues/26216))
* `advisor`: update API version to `2023-01-01` ([#26205](https://github.com/hashicorp/terraform-provider-azurerm/issues/26205))
* `keyvault`: handling the Resources API returning Key Vaults that have been deleted when populating the cache ([#26199](https://github.com/hashicorp/terraform-provider-azurerm/issues/26199))
* `machinelearning`: update API version to `2024-04-01` ([#26168](https://github.com/hashicorp/terraform-provider-azurerm/issues/26168))
* `network/privatelinkservices` - update to use `hashicorp/go-azure-sdk` ([#26212](https://github.com/hashicorp/terraform-provider-azurerm/issues/26212))
* `network/serviceendpointpolicies` - update to use `hashicorp/go-azure-sdk` ([#26196](https://github.com/hashicorp/terraform-provider-azurerm/issues/26196))
* `network/virtualnetworks` - update to use `hashicorp/go-azure-sdk` ([#26217](https://github.com/hashicorp/terraform-provider-azurerm/issues/26217))
* `network/virtualwans`: update route resources to use `hashicorp/go-azure-sdk` ([#26189](https://github.com/hashicorp/terraform-provider-azurerm/issues/26189))
* `azurerm_container_app_job` - support for the `key_vault_secret_id` and `identity`  properties in the  `secret` block ([#25969](https://github.com/hashicorp/terraform-provider-azurerm/issues/25969))
* `azurerm_kubernetes_cluster` -  support forthe  `dns_zone_ids` popperty in the `web_app_routing` block ([#26117](https://github.com/hashicorp/terraform-provider-azurerm/issues/26117))
* `azurerm_notification_hub_authorization_rule` - support for the `primary_connection_string` and `secondary_connection_string` properties ([#26188](https://github.com/hashicorp/terraform-provider-azurerm/issues/26188))
* `azurerm_subnet` - support for the `default_outbound_access_enabled` property ([#25259](https://github.com/hashicorp/terraform-provider-azurerm/issues/25259))

BUG FIXES:

* `azurerm_api_management_named_value` - will now enforce setting the `secret` property when setting the `value_from_key_vault` property ([#26150](https://github.com/hashicorp/terraform-provider-azurerm/issues/26150))
* `azurerm_storage_sync_server_endpoint` - improve pooling to work around api inconsistencies ([#26204](https://github.com/hashicorp/terraform-provider-azurerm/issues/26204))
* `azurerm_virtual_network` - split create and update function to fix lifecycle - ignore ([#26246](https://github.com/hashicorp/terraform-provider-azurerm/issues/26246))
* `azurerm_vpn_server_configuration` - split create and update function to fix lifecycle - ignore ([#26175](https://github.com/hashicorp/terraform-provider-azurerm/issues/26175))
* `azurerm_vpn_server_configuration_policy_group` - split create and update function to fix lifecycle - ignore ([#26207](https://github.com/hashicorp/terraform-provider-azurerm/issues/26207))
* `azurerm_vpn_site` -  split create and update function to fix lifecycle - ignore changes ([#26163](https://github.com/hashicorp/terraform-provider-azurerm/issues/26163))

DEPRECATIONS:

* `azurerm_kubernetes_cluster` - the property `dns_zone_id` has been superseded by the property `dns_zone_ids` in the `web_app_routing` block ([#26117](https://github.com/hashicorp/terraform-provider-azurerm/issues/26117))
* `azurerm_nginx_deployment` - the block `configuration` has been deprecated and superseded by the resource `azurerm_nginx_configuration` ([#25773](https://github.com/hashicorp/terraform-provider-azurerm/issues/25773))

## 3.106.1 (May 31, 2024)

BUG FIXES:

* Data Source: `azurerm_kubernetes_cluster` - fix a crash when reading/setting `upgrade_settings` ([#26173](https://github.com/hashicorp/terraform-provider-azurerm/issues/26173))

## 3.106.0 (May 31, 2024)

UPGRADE NOTES:

* This release updates the Key Vault cache to load Key Vaults using both the Key Vaults List API **and** the Resources API to workaround the API returning incomplete/stale data. To achieve this, and provide consistency between tooling, we are intentionally using the same older version of the Resources API as the current version of Azure CLI. ([#26070](https://github.com/hashicorp/terraform-provider-azurerm/issues/26070))

FEATURES:

* **New Data Source:** `azurerm_arc_resource_bridge_appliance` ([#25731](https://github.com/hashicorp/terraform-provider-azurerm/issues/25731))
* **New Data Source:** `azurerm_elastic_san_volume_group` ([#26111](https://github.com/hashicorp/terraform-provider-azurerm/issues/26111))
* **New Data Source:** `azurerm_storage_queue` ([#26087](https://github.com/hashicorp/terraform-provider-azurerm/issues/26087))
* **New Data Source:** `azurerm_storage_table` ([#26126](https://github.com/hashicorp/terraform-provider-azurerm/issues/26126))
* **New Resource:** `azurerm_container_registry_cache_rule` ([#26034](https://github.com/hashicorp/terraform-provider-azurerm/issues/26034))
* **New Resource:** `azurerm_virtual_machine_implicit_data_disk_from_source` ([#25537](https://github.com/hashicorp/terraform-provider-azurerm/issues/25537))

ENHANCEMENTS:

* Data Source: azurerm_kubernetes_cluster - add support for the `drain_timeout_in_minutes` and `node_soak_duration_in_minutes` properties in the `upgrade_settings` block ([#26137](https://github.com/hashicorp/terraform-provider-azurerm/issues/26137))
* dependencies: updating to `v0.20240529.1155048` of `github.com/hashicorp/go-azure-sdk` ([#26148](https://github.com/hashicorp/terraform-provider-azurerm/issues/26148))
* `containerapps`: update API version to `2024-03-01` ([#25993](https://github.com/hashicorp/terraform-provider-azurerm/issues/25993))
* `expressroute`: update to use `hashicorp/go-azure-sdk` ([#26066](https://github.com/hashicorp/terraform-provider-azurerm/issues/26066))
* `keyvault`: populating the cache using both the Key Vault List and Resources API to workaround incomplete/stale data being returned ([#26070](https://github.com/hashicorp/terraform-provider-azurerm/issues/26070))
* `servicenetworking`: updating to API Version `2023-11-01` ([#26148](https://github.com/hashicorp/terraform-provider-azurerm/issues/26148))
* `virtualnetworkpeerings`: update to use `hashicorp/go-azure-sdk` ([#26065](https://github.com/hashicorp/terraform-provider-azurerm/issues/26065))
* `azurerm_automation_powershell72_module` - support for the `tags` property ([#26106](https://github.com/hashicorp/terraform-provider-azurerm/issues/26106))
* `azurerm_bastion_host` - support for `Developer` SKU ([#26068](https://github.com/hashicorp/terraform-provider-azurerm/issues/26068))
* `azurerm_container_app_environment` - support for the `mutual_tls_enabled` property ([#25993](https://github.com/hashicorp/terraform-provider-azurerm/issues/25993))
* `azurerm_container_registry` - validation to fail fast when setting `public_network_access_enabled` with an invalid SKU ([#26054](https://github.com/hashicorp/terraform-provider-azurerm/issues/26054))
* `azurerm_key_vault_managed_hardware_security_module` - the `public_network_access_enabled` property can now be updated ([#26075](https://github.com/hashicorp/terraform-provider-azurerm/issues/26075))
* `azurerm_kubernetes_cluster` - support for the `cost_analysis_enabled` property ([#26052](https://github.com/hashicorp/terraform-provider-azurerm/issues/26052))
* `azurerm_kubernetes_cluster` - support for the `drain_timeout_in_minutes` and `node_soak_duration_in_minutes` properties in the `upgrade_settings` block ([#26137](https://github.com/hashicorp/terraform-provider-azurerm/issues/26137))
* `azurerm_kubernetes_cluster_node_pool` - support for the `drain_timeout_in_minutes` and `node_soak_duration_in_minutes` properties in the `upgrade_settings` block ([#26137](https://github.com/hashicorp/terraform-provider-azurerm/issues/26137))
* `azurerm_linux_virtual_machine` - the `hibernation_enabled` property can now be updated ([#26112](https://github.com/hashicorp/terraform-provider-azurerm/issues/26112))
* `azurerm_logic_app_trigger_custom` - support for the property `callback_url` ([#25979](https://github.com/hashicorp/terraform-provider-azurerm/issues/25979))
* `azurerm_machine_learning_workspace` - support for the `serverless_compute` block ([#25660](https://github.com/hashicorp/terraform-provider-azurerm/issues/25660))
* `azurerm_mssql_elasticpool` - support the sku `HS_PRMS` ([#26161](https://github.com/hashicorp/terraform-provider-azurerm/issues/26161))
* `azurerm_new_relic_monitor` - support for the `identity` block ([#26115](https://github.com/hashicorp/terraform-provider-azurerm/issues/26115))
* `azurerm_route_map` - the `parameter` property is now Optional when the action type is `Drop` ([#26003](https://github.com/hashicorp/terraform-provider-azurerm/issues/26003))
* `azurerm_windows_virtual_machine` - the `hibernation_enabled` property can now be updated ([#26112](https://github.com/hashicorp/terraform-provider-azurerm/issues/26112))

BUG FIXES:

* Data Source: `azurerm_system_center_virtual_machine_manager_inventory_items` - normalise the resource ID for Intentory Items ([#25955](https://github.com/hashicorp/terraform-provider-azurerm/issues/25955))
* `azurerm_app_configuration_feature` - update polling interval to tolerate eventual consistency of the API ([#26025](https://github.com/hashicorp/terraform-provider-azurerm/issues/26025))
* `azurerm_app_configuration_key` - update polling interval to tolerate eventual consistency of the API ([#26025](https://github.com/hashicorp/terraform-provider-azurerm/issues/26025))
* `azurerm_eventhub_namespace_customer_managed_key` - validating that the User Assigned Identity used for accessing the Key Vault is assigned to the EventHub Namespace ([#28509](https://github.com/hashicorp/terraform-provider-azurerm/issues/28509))
* `azurerm_linux_function_app` - fix update handling of `health_check_eviction_time_in_min` and `WEBSITE_HEALTHCHECK_MAXPINGFAILURES` ([#26107](https://github.com/hashicorp/terraform-provider-azurerm/issues/26107))
* `azurerm_linux_function_app_slot` - fix update handling of `health_check_eviction_time_in_min` and `WEBSITE_HEALTHCHECK_MAXPINGFAILURES` ([#26107](https://github.com/hashicorp/terraform-provider-azurerm/issues/26107))
* `azurerm_linux_web_app` - fix update handling of `health_check_eviction_time_in_min` and `WEBSITE_HEALTHCHECK_MAXPINGFAILURES` ([#26107](https://github.com/hashicorp/terraform-provider-azurerm/issues/26107))
* `azurerm_linux_web_app_slot` - fix update handling of `health_check_eviction_time_in_min` and `WEBSITE_HEALTHCHECK_MAXPINGFAILURES` ([#26107](https://github.com/hashicorp/terraform-provider-azurerm/issues/26107))
* `azurerm_postgresql_flexible_server` - prevent premature check on updated `storage_mb` value that prevents the resource from being re-created ([#25986](https://github.com/hashicorp/terraform-provider-azurerm/issues/25986))
* `azurerm_redis_access_cache_policy_assignment` - add locks to stabilize creation of multiple policy assignments ([#26085](https://github.com/hashicorp/terraform-provider-azurerm/issues/26085))
* `azurerm_redis_access_cache_policy` - add locks to stabilize creation of multiple policy assignments ([#26085](https://github.com/hashicorp/terraform-provider-azurerm/issues/26085))
* `azurerm_windows_function_app` - fix update handling of `health_check_eviction_time_in_min` and `WEBSITE_HEALTHCHECK_MAXPINGFAILURES` ([#26107](https://github.com/hashicorp/terraform-provider-azurerm/issues/26107))
* `azurerm_windows_function_app_slot` - fix update handling of `health_check_eviction_time_in_min` and `WEBSITE_HEALTHCHECK_MAXPINGFAILURES` ([#26107](https://github.com/hashicorp/terraform-provider-azurerm/issues/26107))
* `azurerm_windows_web_app` - fix update handling of `health_check_eviction_time_in_min` and `WEBSITE_HEALTHCHECK_MAXPINGFAILURES` ([#26107](https://github.com/hashicorp/terraform-provider-azurerm/issues/26107))
* `azurerm_windows_web_app_slot` - fix update handling of `health_check_eviction_time_in_min` and `WEBSITE_HEALTHCHECK_MAXPINGFAILURES` ([#26107](https://github.com/hashicorp/terraform-provider-azurerm/issues/26107))

## 3.105.0 (May 24, 2024)

BREAKING CHANGE:

* `azurerm_kubernetes_cluster` - the properties `workload_autoscaler_profile.vertical_pod_autoscaler_update_mode` and `workload_autoscaler_profile.vertical_pod_autoscaler_controlled_values` are no longer populated since they're not exported in API version `2023-09-02-preview` ([#25663](https://github.com/hashicorp/terraform-provider-azurerm/issues/25663))

FEATURES:

* New Resource: `azurerm_api_management_policy_fragment` ([#24968](https://github.com/hashicorp/terraform-provider-azurerm/issues/24968))

ENHANCEMENTS:

* dependencies: updating to `v0.20240522.1080424` of `github.com/hashicorp/go-azure-sdk` ([#26069](https://github.com/hashicorp/terraform-provider-azurerm/issues/26069))
* `containerservice`: updating to use API Version `2023-09-02-preview` ([#25663](https://github.com/hashicorp/terraform-provider-azurerm/issues/25663))
* `azurerm_application_insights_standard_web_test` - `http_verb` can now be set to `HEAD` and `OPTIONS` ([#26077](https://github.com/hashicorp/terraform-provider-azurerm/issues/26077))
* `azurerm_cdn_frontdoor_rule` - updating the validation for `match_values` within the `uri_path_condition` block  to support a forward-slash ([#26017](https://github.com/hashicorp/terraform-provider-azurerm/issues/26017))
* `azurerm_linux_web_app` - normalising the value for `virtual_network_subnet_id` ([#25885](https://github.com/hashicorp/terraform-provider-azurerm/issues/25885))
* `azurerm_machine_learning_compute_cluster` - add validation for `name` ([#26060](https://github.com/hashicorp/terraform-provider-azurerm/issues/26060))
* `azurerm_machine_learning_compute_cluster` - improve validation to allow an empty `subnet_resource_id` when the Workspace is using a managed Virtual Network ([#26073](https://github.com/hashicorp/terraform-provider-azurerm/issues/26073))
* `azurerm_postgresql_flexible_server` - the field `public_network_access_enabled` is now configurable (previously this was computed-only/not settable via the API) ([#25812](https://github.com/hashicorp/terraform-provider-azurerm/issues/25812))
* `azurerm_snapshot` - support for `disk_access_id` ([#25996](https://github.com/hashicorp/terraform-provider-azurerm/issues/25996))
* `azurerm_windows_web_app` - normalising the value for `virtual_network_subnet_id` ([#25885](https://github.com/hashicorp/terraform-provider-azurerm/issues/25885))

BUG FIXES:

* `azurerm_container_app_environment_custom_domain`: parsing the Log Analytics Workspace ID insensitively to workaround the API returning this inconsistently ([#26074](https://github.com/hashicorp/terraform-provider-azurerm/issues/26074))
* `azurerm_container_app_job` - updating the validation for the `name` field ([#26049](https://github.com/hashicorp/terraform-provider-azurerm/issues/26049))
* `azurerm_container_app_job` - updating the validation for the `name` field within the `custom_scale_rule` block ([#26049](https://github.com/hashicorp/terraform-provider-azurerm/issues/26049))
* `azurerm_container_app_job` - updating the validation for the `name` field within the `rules` block ([#26049](https://github.com/hashicorp/terraform-provider-azurerm/issues/26049))
* `azurerm_linux_function_app_slot` - fixed panic when planning from a version older than 3.88.0 ([#25838](https://github.com/hashicorp/terraform-provider-azurerm/issues/25838))
* `azurerm_pim_active_role_assignment` - fix a persistent diff when `ticket` is not specified ([#26059](https://github.com/hashicorp/terraform-provider-azurerm/issues/26059))
* `azurerm_pim_eligible_role_assignment` - fix a persistent diff when `ticket` is not specified ([#26059](https://github.com/hashicorp/terraform-provider-azurerm/issues/26059))
* `azurerm_policy_definition` - recreate the resource if the `parameters` property is updated to include fewer items ([#26083](https://github.com/hashicorp/terraform-provider-azurerm/issues/26083))
* `azurerm_windows_function_app_slot` - set Server Farm ID in payload when using a Virtual Network Subnet for the slot ([#25634](https://github.com/hashicorp/terraform-provider-azurerm/issues/25634))
* `azurerm_windows_web_app_slot` - set Server Farm ID in payload when using a Virtual Network Subnet for the slot ([#25634](https://github.com/hashicorp/terraform-provider-azurerm/issues/25634))


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

## 3.69.0 (August 10, 2023)

FEATURES:

* **New Data Source**: `azurerm_palo_alto_local_rulestack` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_graph_services_account` ([#22665](https://github.com/hashicorp/terraform-provider-azurerm/issues/22665))
* **New Resource**: `azurerm_managed_lustre_file_system` ([#22680](https://github.com/hashicorp/terraform-provider-azurerm/issues/22680))
* **New Resource**: `azurerm_palo_alto_local_rulestack` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_palo_alto_local_rulestack_certificate` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_palo_alto_local_rulestack_fqdn_list` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_palo_alto_local_rulestack_outbound_trust_certificate_association` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_palo_alto_local_rulestack_outbound_untrust_certificate_association`  ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_palo_alto_local_rulestack_prefix_list` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_palo_alto_local_rulestack_rule` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_palo_alto_virtual_network_appliance` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_palo_alto_next_generation_firewall_virtual_hub_local_rulestack` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_palo_alto_next_generation_firewall_virtual_hub_panorama` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_palo_alto_next_generation_firewall_virtual_network_local_rulestack` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))
* **New Resource**: `azurerm_palo_alto_next_generation_firewall_virtual_network_panorama` ([#22700](https://github.com/hashicorp/terraform-provider-azurerm/issues/22700))

ENHANCEMENTS:

* dependencies: updating to `v0.58.0` of `github.com/hashicorp/go-azure-helpers` ([#22813](https://github.com/hashicorp/terraform-provider-azurerm/issues/22813))
* dependencies: updating to `v0.20230808.1103829` of `github.com/hashicorp/go-azure-sdk` ([#22860](https://github.com/hashicorp/terraform-provider-azurerm/issues/22860))
* `arckubernetes` - updating to use the `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` as a base layer ([#22815](https://github.com/hashicorp/terraform-provider-azurerm/issues/22815))
* `bot` - updating to use the `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` as a base layer ([#22815](https://github.com/hashicorp/terraform-provider-azurerm/issues/22815))
* `blueprints`: updating to use `hashicorp/go-azure-sdk` ([#21569](https://github.com/hashicorp/terraform-provider-azurerm/issues/21569))
* `compute` - updating to use the `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` as a base layer ([#22860](https://github.com/hashicorp/terraform-provider-azurerm/issues/22860))
* `digitaltwins` - updating to API Version `2023-01-31` ([#22782](https://github.com/hashicorp/terraform-provider-azurerm/issues/22782))
* `hsm` - updating to use the `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` as a base layer ([#22815](https://github.com/hashicorp/terraform-provider-azurerm/issues/22815))
* `hybridcompute` - updating to use the `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` as a base layer ([#22815](https://github.com/hashicorp/terraform-provider-azurerm/issues/22815))
* Data Source: `azurerm_network_service_tags` - updating to use `hashicorp/go-azure-sdk` ([#22873](https://github.com/hashicorp/terraform-provider-azurerm/issues/22873))
* Data Source: `azurerm_network_watcher` - updating to use `hashicorp/go-azure-sdk` ([#22873](https://github.com/hashicorp/terraform-provider-azurerm/issues/22873))
* `azurerm_container_app_environment` - `log_analytics_workspace_id` is now an Optional property ([#22733](https://github.com/hashicorp/terraform-provider-azurerm/issues/22733))
* `azurerm_digital_twins_instance` - support for User Assigned Identities ([#22782](https://github.com/hashicorp/terraform-provider-azurerm/issues/22782))
* `azurerm_function_app_function` - hyphen and underscore are now allows characters for function names ([#22519](https://github.com/hashicorp/terraform-provider-azurerm/issues/22519))
* `azurerm_key_vault_certificate` - Support update of certificates based on `certificate_policy` ([#20627](https://github.com/hashicorp/terraform-provider-azurerm/issues/20627))
* `azurerm_kubernetes_cluster` - export the identity for Web App Routing under `web_app_routing_identity` ([#22809](https://github.com/hashicorp/terraform-provider-azurerm/issues/22809))
* `azurerm_kubernetes_cluster` - add support for the `snapshot_id` property in the `default_node_pool` block ([#22708](https://github.com/hashicorp/terraform-provider-azurerm/issues/22708))
* `azurerm_log_analytics_workspace` - support changing value of `sku` from `CapacityReservation` and `PerGB2018` ([#22597](https://github.com/hashicorp/terraform-provider-azurerm/issues/22597))
* `azurerm_managed_application` - deprecate the `parameters` property in favour of `parameter_values` ([#21541](https://github.com/hashicorp/terraform-provider-azurerm/issues/21541))
* `azurerm_monitor_action_group` - the value `https` is now supported for `aad_auth` ([#22888](https://github.com/hashicorp/terraform-provider-azurerm/issues/22888))
* `azurerm_mssql_server` - `SystemAssigned, UserAssigned` identity is now supported ([#22828](https://github.com/hashicorp/terraform-provider-azurerm/issues/22828))
* `azurerm_network_packet_capture` - updating to use `hashicorp/go-azure-sdk` ([#22873](https://github.com/hashicorp/terraform-provider-azurerm/issues/22873))
* `azurerm_network_profile` - refactoring to use `hashicorp/go-azure-sdk` ([#22850](https://github.com/hashicorp/terraform-provider-azurerm/issues/22850))
* `azurerm_network_watcher_flow_log` - updating to use `hashicorp/go-azure-sdk` ([#22873](https://github.com/hashicorp/terraform-provider-azurerm/issues/22873))
* `azurerm_network_watcher` - updating to use `hashicorp/go-azure-sdk` ([#22873](https://github.com/hashicorp/terraform-provider-azurerm/issues/22873))
* `azurerm_postgresql_database` - updating the validation for `collation` ([#22689](https://github.com/hashicorp/terraform-provider-azurerm/issues/22689))
* `azurerm_postgresql_flexible_server_database` - updating the validation for `collation` ([#22689](https://github.com/hashicorp/terraform-provider-azurerm/issues/22689))
* `azurerm_security_center_subscription_pricing` - support for `extensions` block  ([#22643](https://github.com/hashicorp/terraform-provider-azurerm/issues/22643))
* `azurerm_security_center_subscription_pricing` - support for the `resource_type` `Api` ([#22844](https://github.com/hashicorp/terraform-provider-azurerm/issues/22844))
* `azurerm_spring_cloud_configuration_service` - support for the `ca_certificate_id` property ([#22814](https://github.com/hashicorp/terraform-provider-azurerm/issues/22814))
* `azurerm_virtual_desktop_workspace` - added support for the `public_network_access_enabled` property ([#22542](https://github.com/hashicorp/terraform-provider-azurerm/issues/22542))
* `azurerm_virtual_machine_packet_capture` - updating to use `hashicorp/go-azure-sdk` ([#22873](https://github.com/hashicorp/terraform-provider-azurerm/issues/22873))
* `azurerm_virtual_machine_scale_set_packet_capture` - updating to use `hashicorp/go-azure-sdk` ([#22873](https://github.com/hashicorp/terraform-provider-azurerm/issues/22873))
* `azurerm_vpn_gateway_connection` - updating to use `hashicorp/go-azure-sdk` ([#22873](https://github.com/hashicorp/terraform-provider-azurerm/issues/22873))
* `azurerm_vpn_server_configuration` - refactoring to use `hashicorp/go-azure-sdk` ([#22850](https://github.com/hashicorp/terraform-provider-azurerm/issues/22850))
* `azurerm_vpn_server_configuration_policy_group` - refactoring to use `hashicorp/go-azure-sdk` ([#22850](https://github.com/hashicorp/terraform-provider-azurerm/issues/22850))
* `azurerm_vpn_site` - refactoring to use `hashicorp/go-azure-sdk` ([#22850](https://github.com/hashicorp/terraform-provider-azurerm/issues/22850))

BUG FIXES:

* Data Source: `azurerm_virutal_machine` - correctly retrieve and set value for `power_state` ([#22851](https://github.com/hashicorp/terraform-provider-azurerm/issues/22851))
* `azurerm_cdn_endpoint` - conditionally using `PUT` in place of `PATCH` when a field other than `tags` has changed ([#22662](https://github.com/hashicorp/terraform-provider-azurerm/issues/22662))
* `azurerm_cdn_frontdoor_security_policy` - normalizing the value returned from the API for `cdn_frontdoor_domain_id` ([#22841](https://github.com/hashicorp/terraform-provider-azurerm/issues/22841))
* `azurerm_container_group` - set `init_container.secure_environment_variables` into state correctly ([#22832](https://github.com/hashicorp/terraform-provider-azurerm/issues/22832))
* `azurerm_custom_ip_prefix` - support for environments other than Azure Public ([#22812](https://github.com/hashicorp/terraform-provider-azurerm/issues/22812))
* `azurerm_databricks_workspace` - update parse function for `machine_learning_workspace_id` field validation ([#22865](https://github.com/hashicorp/terraform-provider-azurerm/issues/22865))
* `azurerm_key_vault` - fixing support for the `storage` Nested Item type ([#22707](https://github.com/hashicorp/terraform-provider-azurerm/issues/22707))
* `azurerm_kusto_cosmosdb_data_connection_resource` - ensure the `subscriptionId` and `ResourceGroupName` align with the CosmosDB container ([#22663](https://github.com/hashicorp/terraform-provider-azurerm/issues/22663))
* `azurerm_managed_application` - fix an issue where `secureString` parameters were not persisted to state ([#21541](https://github.com/hashicorp/terraform-provider-azurerm/issues/21541))
* `azurerm_managed_application` - the `plan` block is now marked ForceNew to comply with service limitations ([#21541](https://github.com/hashicorp/terraform-provider-azurerm/issues/21541))
* `azurerm_monitor_data_collection_rule` - recreate resource when attempting to remove `kind` ([#22811](https://github.com/hashicorp/terraform-provider-azurerm/issues/22811))
* `azurerm_static_site_custom_domain` - prevent overwriting `validation_token` with an empty value by setting it into state when creating the resource ([#22848](https://github.com/hashicorp/terraform-provider-azurerm/issues/22848))

## 3.68.0 (August 03, 2023)

FEATURES:

* **New Resource:** `azurerm_custom_ip_prefix` ([#21322](https://github.com/hashicorp/terraform-provider-azurerm/issues/21322))
* **New Resource:**: `azurerm_mobile_network_sim` ([#22628](https://github.com/hashicorp/terraform-provider-azurerm/issues/22628))
* **New Data Source:** `azurerm_mobile_network_sim` ([#22628](https://github.com/hashicorp/terraform-provider-azurerm/issues/22628))
* **New Resource:** `azurerm_automation_variable_object` ([#22644](https://github.com/hashicorp/terraform-provider-azurerm/issues/22644))
* **New Data Source:** `azurerm_automation_variable_object` ([#22644](https://github.com/hashicorp/terraform-provider-azurerm/issues/22644))

ENHANCEMENTS

* dependencies: updating to `v0.20230803.1095722` of `github.com/hashicorp/go-azure-sdk` ([#22803](https://github.com/hashicorp/terraform-provider-azurerm/issues/22803))
* dependencies: migrate mysql resources to `hashicorp/go-azure-sdk` ([#22795](https://github.com/hashicorp/terraform-provider-azurerm/issues/22795))
* `advisor`: updating the base layer to use `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22750](https://github.com/hashicorp/terraform-provider-azurerm/issues/22750))
* `apimanagement`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22759](https://github.com/hashicorp/terraform-provider-azurerm/issues/22759))
* `analysisservices`: updating the base layer to use `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22750](https://github.com/hashicorp/terraform-provider-azurerm/issues/22750))
* `automation`: updating `dscnodeconfiguration` and `sourcecontrol` to use API Version `2022-08-08` ([#22781](https://github.com/hashicorp/terraform-provider-azurerm/issues/22781))
* `azurestackhci`: updating the base layer to use `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22750](https://github.com/hashicorp/terraform-provider-azurerm/issues/22750))
* `domainservices`: updating the base layer to use `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22750](https://github.com/hashicorp/terraform-provider-azurerm/issues/22750))
* `eventgrid`: refactoring to use `hashicorp/go-azure-sdk` ([#22673](https://github.com/hashicorp/terraform-provider-azurerm/issues/22673))
* `machinelearningservice`: updating to use API Version `2023-04-01` ([#22729](https://github.com/hashicorp/terraform-provider-azurerm/issues/22729))
* `monitor`: updating the base layer to use `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22750](https://github.com/hashicorp/terraform-provider-azurerm/issues/22750))
* `network`: updating to use API Version `2023-04-01` ([#22727](https://github.com/hashicorp/terraform-provider-azurerm/issues/22727))
* `relay`: updating to use API Version `2021-11-01` ([#22725](https://github.com/hashicorp/terraform-provider-azurerm/issues/22725))
* Data Source: `azurerm_images` - support for `disk_encryption_set_id` ([#22690](https://github.com/hashicorp/terraform-provider-azurerm/issues/22690))
* `azurerm_eventhub_namespace_customer_managed_key` - support for the `infrastructure_encryption_enabled` property ([#22718](https://github.com/hashicorp/terraform-provider-azurerm/issues/22718))
* `azurerm_hpc_cache_blob_nfs_target` - support for setting the `usage_model` property to `READ_ONLY` and `READ_WRITE` ([#22798](https://github.com/hashicorp/terraform-provider-azurerm/issues/22798))
* `azurerm_hpc_cache_nfs_target` - support for setting the `usage_model` property to `READ_ONLY` and `READ_WRITE` ([#22798](https://github.com/hashicorp/terraform-provider-azurerm/issues/22798))
* `azurerm_monitor_aad_diagnostic_setting` - updating to use `hashicorp/go-azure-sdk` ([#22778](https://github.com/hashicorp/terraform-provider-azurerm/issues/22778))
* `azurerm_web_application_firewall_policy` - updating to use API Version `2023-02-01` ([#22455](https://github.com/hashicorp/terraform-provider-azurerm/issues/22455))
* `azurerm_web_application_firewall_policy` - support for `log_scrubbing` property ([#22522](https://github.com/hashicorp/terraform-provider-azurerm/issues/22522))
* `azurerm_shared_image_gallery` - support for the `sharing` block ([#22221](https://github.com/hashicorp/terraform-provider-azurerm/issues/22221))
* `azurerm_virtual_network` - support for the `encryption` block ([#22745](https://github.com/hashicorp/terraform-provider-azurerm/issues/22745))
  
BUG FIXES

* provider: only obtaining an authentication token for Managed HSM in environments where Managed HSM is available ([#22400](https://github.com/hashicorp/terraform-provider-azurerm/issues/22400))
* `azurerm_api_management` - retrieving the `location` from the API rather than the config prior to deletion ([#22752](https://github.com/hashicorp/terraform-provider-azurerm/issues/22752))
* `azurerm_cognitive_deployment` - add locks to parent resource to prevent 409 error ([#22711](https://github.com/hashicorp/terraform-provider-azurerm/issues/22711))
* `azurerm_pim_eligible_role_assignment` - fixing a bug where the context deadline was checked incorrectly during deletion ([#22756](https://github.com/hashicorp/terraform-provider-azurerm/issues/22756))
* `azurerm_private_endpoint` - loading the subnet to lock from the API rather than the config during deletion ([#22676](https://github.com/hashicorp/terraform-provider-azurerm/issues/22676))
* `azurerm_netapp_volume` - updating the validation of `security_style` to match the casing defined in the Azure API Definitions ([#22721](https://github.com/hashicorp/terraform-provider-azurerm/issues/22721))
* `azurerm_netapp_volume_group_sap_hana` - update the validation of `security_style` to match the casing defined in the Azure API Definitions ([#22615](https://github.com/hashicorp/terraform-provider-azurerm/issues/22615))
* `azurerm_site_recovery_replication_recovery_plan` - fix update for `boot_recovery_group`,`failover_recovery_group` and `shutdown_recovery_group` ([#22687](https://github.com/hashicorp/terraform-provider-azurerm/issues/22687))

## 3.67.0 (July 27, 2023)

FEATURES:

* **New Data Source:** `azurerm_eventhub_sas` ([#22215](https://github.com/hashicorp/terraform-provider-azurerm/issues/22215))
* **New Resource**: `azurerm_kubernetes_cluster_trusted_access_role_binding` ([#22647](https://github.com/hashicorp/terraform-provider-azurerm/issues/22647))
* **New Resource:** `azurerm_marketplace_role_assignment` ([#22398](https://github.com/hashicorp/terraform-provider-azurerm/issues/22398))
* **New Resource:** `azurerm_network_function_azure_traffic_collector` ([#22274](https://github.com/hashicorp/terraform-provider-azurerm/issues/22274))

ENHANCEMENTS:

* dependencies: updating to `v0.20230726.1135558` of `github.com/hashicorp/go-azure-sdk` ([#22698](https://github.com/hashicorp/terraform-provider-azurerm/issues/22698))
* `connections`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22681](https://github.com/hashicorp/terraform-provider-azurerm/issues/22681))
* `iothub`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22681](https://github.com/hashicorp/terraform-provider-azurerm/issues/22681))
* `mysql`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22681](https://github.com/hashicorp/terraform-provider-azurerm/issues/22681))
* `orbital`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22681](https://github.com/hashicorp/terraform-provider-azurerm/issues/22681))
* `powerbi`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22681](https://github.com/hashicorp/terraform-provider-azurerm/issues/22681))
* `privatedns`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22681](https://github.com/hashicorp/terraform-provider-azurerm/issues/22681))
* `purview`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22681](https://github.com/hashicorp/terraform-provider-azurerm/issues/22681))
* `relay`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` ([#22681](https://github.com/hashicorp/terraform-provider-azurerm/issues/22681))
* `azurerm_cdn_endpoint_custom_domain` - pass nil as version when Certificate/Secret version is set to Latest ([#22683](https://github.com/hashicorp/terraform-provider-azurerm/issues/22683))
* `azurerm_image` - support for the field `disk_encryption_set_id` within the `os_disk` block ([#22642](https://github.com/hashicorp/terraform-provider-azurerm/issues/22642))
* `azurerm_linux_virtual_machine` - add support for the `bypass_platform_safety_checks_on_user_schedule_enabled` and `reboot_setting` properties ([#22349](https://github.com/hashicorp/terraform-provider-azurerm/issues/22349))
* `azurerm_network_interface` - updating to use `hashicorp/go-azure-sdk` and API Version `2023-02-01` ([#22479](https://github.com/hashicorp/terraform-provider-azurerm/issues/22479))
* `azurerm_redis_enterprise_database` - support `redisSON` module for geo-replication ([#22627](https://github.com/hashicorp/terraform-provider-azurerm/issues/22627))
* `azurerm_windows_virtual_machine` - add support for the `bypass_platform_safety_checks_on_user_schedule_enabled` and `reboot_setting` properties ([#22349](https://github.com/hashicorp/terraform-provider-azurerm/issues/22349))

BUG FIXES:

* `azurerm_cosmosdb_account` - `type` within the `backup` block is updated separately when set to `Continuous` ([#22638](https://github.com/hashicorp/terraform-provider-azurerm/issues/22638))
* `azurerm_cosmosdb_account` - `max_age_in_seconds` within the `cors_rule` block is now Optional and can now be configured up to `2147483647` ([#22552](https://github.com/hashicorp/terraform-provider-azurerm/issues/22552))
* `azurerm_maintenance_configuration` - fixing a bug where include and exclude were set incorrectly ([#22671](https://github.com/hashicorp/terraform-provider-azurerm/issues/22671))
* `azurerm_pim_eligible_role_assignment` - polling for the duration of the timeout, rather than using a hard-coded value ([#22682](https://github.com/hashicorp/terraform-provider-azurerm/issues/22682))
* `azurerm_redis_cache` - only updating `patch_schedule` when it has changed in the config file ([#22661](https://github.com/hashicorp/terraform-provider-azurerm/issues/22661))
* `azurerm_logic_app_standard` - attribute `auto_swap_slot_name` is now under correct block `site_config` ([#22712](https://github.com/hashicorp/terraform-provider-azurerm/issues/22712))
* `azurerm_postgresql_flexible_server` - update the validation of `storage_mb` replacing `33554432` with `33553408` ([#22706](https://github.com/hashicorp/terraform-provider-azurerm/issues/22706))

## 3.66.0 (July 20, 2023)

FEATURES:

* **New Data Source:** `azurerm_mobile_network_attached_data_network` ([#22168](https://github.com/hashicorp/terraform-provider-azurerm/issues/22168))
* **New Resource:** `azurerm_graph_account` ([#22334](https://github.com/hashicorp/terraform-provider-azurerm/issues/22334))
* **New Resource:** `azurerm_mobile_network_attached_data_network` ([#22168](https://github.com/hashicorp/terraform-provider-azurerm/issues/22168))

ENHANCEMENTS:

* dependencies: bump `go-azure-sdk` to `v0.20230720.1190320` and switch `machinelearning`, `mixedreality`, `mariadb`, `storagecache`, `storagepool`, `vmware`, `videoanalyzer`, `voiceServices` and `mobilenetwork` to new base layer ([#22538](https://github.com/hashicorp/terraform-provider-azurerm/issues/22538))
* dependencies: move `azurerm_bastion_host` and `azurerm_network_connection_monitor` over to `hashicorp/go-azure-sdk` ([#22425](https://github.com/hashicorp/terraform-provider-azurerm/issues/22425))
* dependencies: move `azurerm_network_watcher_flow_log` to `hashicorp/go-azure-sdk` ([#22575](https://github.com/hashicorp/terraform-provider-azurerm/issues/22575))
* dependencies: move `mysql` resources over to `hashicorp/go-azure-sdk` ([#22528](https://github.com/hashicorp/terraform-provider-azurerm/issues/22528))
* dependencies: move `storage_sync` resources over to `hashicorp/go-azure-sdk` ([#21928](https://github.com/hashicorp/terraform-provider-azurerm/issues/21928))
* dependencies: updating to API Version `2022-08-08` ([#22440](https://github.com/hashicorp/terraform-provider-azurerm/issues/22440))
* `postgres` - updating to API Version `2023-03-01-preview` ([#22577](https://github.com/hashicorp/terraform-provider-azurerm/issues/22577))
* `data.azurerm_route_table` - support for the `bgp_route_propagation_enabled` property ([#21940](https://github.com/hashicorp/terraform-provider-azurerm/issues/21940))
* `data.azurerm_servicebus_*` - add deprecation messages for the `resource_group_name` and `namespace_name` properties ([#22521](https://github.com/hashicorp/terraform-provider-azurerm/issues/22521))
* `azurerm_cdn_frontdoor_rule` - allow the `conditions.x.url_path_condition.x.match_values` property to be set to `/` ([#22610](https://github.com/hashicorp/terraform-provider-azurerm/issues/22610))
* `azurerm_eventhub_namespace` - updates properly when encryption is enabled ([#22625](https://github.com/hashicorp/terraform-provider-azurerm/issues/22625))
* `azurerm_logic_app_standard` - now exports the `auto_swap_slot_name` attribute ([#22525](https://github.com/hashicorp/terraform-provider-azurerm/issues/22525))
* `azurerm_mysql_flexible_server_configuration` - the `value` property can now be changed without creating a new resource ([#22557](https://github.com/hashicorp/terraform-provider-azurerm/issues/22557))
* `azurerm_postgresql_flexible_server` - support for `33554432` storage ([#22574](https://github.com/hashicorp/terraform-provider-azurerm/issues/22574))
* `azurerm_postgresql_flexible_server` - support for the `geo_backup_key_vault_key_id` and `geo_backup_user_assigned_identity_id` properties ([#22612](https://github.com/hashicorp/terraform-provider-azurerm/issues/22612))
* `azurerm_spring_cloud_service` - support for the `marketplace` block ([#22553](https://github.com/hashicorp/terraform-provider-azurerm/issues/22553))
* `azurerm_spring_cloud_service` - support for the `outbound_type` property ([#22596](https://github.com/hashicorp/terraform-provider-azurerm/issues/22596))

BUG FIXES:

* provider: the Resource Providers `Microsoft.Kubernetes` and `Microsoft.KubernetesConfiguration` are no longer automatically registered ([#22580](https://github.com/hashicorp/terraform-provider-azurerm/issues/22580))
* `data.automation_account_variables` - correctly populate missing variable attributes ([#22611](https://github.com/hashicorp/terraform-provider-azurerm/issues/22611))
* `data.azurerm_virtual_machine_scale_set` - fix an issue where `computer_name`, `latest_model_applied`, `power_state` and `virtual_machine_id` attributes were not correctly set ([#22566](https://github.com/hashicorp/terraform-provider-azurerm/issues/22566))
* `azurerm_app_service_public_certificate` - poll for certificate during read to get around an eventual consistency bug ([#22587](https://github.com/hashicorp/terraform-provider-azurerm/issues/22587))
* `azurerm_application_gateway` - send `min_protocol_version` and correct `policy_type` when using `CustomV2` ([#22535](https://github.com/hashicorp/terraform-provider-azurerm/issues/22535))
* `azurerm_cognitive_deployment` - remove upper limit on validation for the `capacity` property in the `scale` block ([#22502](https://github.com/hashicorp/terraform-provider-azurerm/issues/22502))
* `azurerm_cosmosdb_account` - fixed regression to `default_identity_type` being switched to `FirstPartyIdentity` on update ([#22609](https://github.com/hashicorp/terraform-provider-azurerm/issues/22609))
* `azurerm_kubernetes_cluster` - the `windows_profile.admin_password` property will become Required in `v4.0` ([#22554](https://github.com/hashicorp/terraform-provider-azurerm/issues/22554))
* `azurerm_kusto_cluster` - the `engine` property has been deprecataed and is now non functional as the service team intends to remove it from the API ([#22497](https://github.com/hashicorp/terraform-provider-azurerm/issues/22497))
* `azurerm_maintenance_configuration` - tge `package_names_mask_to_exclude` and `package_names_mask_to_exclude` properties are not set properly ([#22555](https://github.com/hashicorp/terraform-provider-azurerm/issues/22555))
* `azurerm_redis_cache` - only set the `rdb_backup_enabled` property when using a premium SKU ([#22309](https://github.com/hashicorp/terraform-provider-azurerm/issues/22309))
* `azurerm_site_recovery_replication_recovery_plan` - fix an issue where the order of boot recovery groups was not correctly maintained ([#22348](https://github.com/hashicorp/terraform-provider-azurerm/issues/22348))
* `azurerm_synapse_firewall_rule` - correct an overly strict validation for the `name` property ([#22571](https://github.com/hashicorp/terraform-provider-azurerm/issues/22571))

## 3.65.0 (July 13, 2023)

FEATURES:

* **New Data Source**: `azurerm_communication_service` ([#22426](https://github.com/hashicorp/terraform-provider-azurerm/issues/22426))

ENHANCEMENTS:

* dependencies: updating to `v0.20230712.1084117` of `github.com/hashicorp/go-azure-sdk` ([#22491](https://github.com/hashicorp/terraform-provider-azurerm/issues/22491))
* dependencies: updating to `v0.20230703.1101016` of `github.com/tombuildsstuff/kermit` ([#22390](https://github.com/hashicorp/terraform-provider-azurerm/issues/22390))
* provider: the Resource Providers `Microsoft.Kubernetes` and `Microsoft.KubernetesConfiguration` are now automatically registered ([#22463](https://github.com/hashicorp/terraform-provider-azurerm/issues/22463))
* `automation/dscconfiguration` - updating to API Version `2022-08-08` ([#22403](https://github.com/hashicorp/terraform-provider-azurerm/issues/22403))
* `azurestackhcl` - updating to API Version `2023-03-01` ([#22411](https://github.com/hashicorp/terraform-provider-azurerm/issues/22411))
* `batch` - updating to use API Version `2023-05-01` ([#22412](https://github.com/hashicorp/terraform-provider-azurerm/issues/22412))
* `datafactory` - moving `azurerm_data_factory` and `azurerm_data_factory_managed_private_endpoint` over to `hashicorp/go-azure-sdk` ([#22409](https://github.com/hashicorp/terraform-provider-azurerm/issues/22409))
* `elastic` - updating to API Version `2023-06-01` ([#22451](https://github.com/hashicorp/terraform-provider-azurerm/issues/22451))
* `kusto` - updating to API Version `2023-05-02` [GH-22410
* `managedapplications` - migrate to `hashicorp/go-azure-sdk` ([#21571](https://github.com/hashicorp/terraform-provider-azurerm/issues/21571))
* `privatedns`: updating to API Version `2020-06-01` ([#22470](https://github.com/hashicorp/terraform-provider-azurerm/issues/22470))
* `storage` - updating to Data Plane API Version `2020-08-04` ([#22405](https://github.com/hashicorp/terraform-provider-azurerm/issues/22405))
* `network` - `application_security_group` and `private_endpoint` now use `hashicorp/go-azure-sdk` ([#22396](https://github.com/hashicorp/terraform-provider-azurerm/issues/22396))
* `voiceservices`: updating to use API Version `2023-04-03` ([#22469](https://github.com/hashicorp/terraform-provider-azurerm/issues/22469))
* Data Source: `azurerm_kubernetes_cluster` - add support for the `internal_ingress_gateway_enabled` and `external_ingress_gateway_enabled` properties ([#22393](https://github.com/hashicorp/terraform-provider-azurerm/issues/22393))
* `azurerm_batch_account` - support for the `network_profile` block ([#22356](https://github.com/hashicorp/terraform-provider-azurerm/issues/22356))
* `azurerm_container_app` - the `min_replicas` and `max_replicas` propertiesnow support a maximum value of `300` ([#22511](https://github.com/hashicorp/terraform-provider-azurerm/issues/22511))
* `azurerm_dns_zone` - can now use the `host_name` property with `dns_zone` for `soa_record` creation ([#22312](https://github.com/hashicorp/terraform-provider-azurerm/issues/22312))
* `azurerm_kubernetes_cluster` - add support for the `internal_ingress_gateway_enabled` and `external_ingress_gateway_enabled` properties ([#22393](https://github.com/hashicorp/terraform-provider-azurerm/issues/22393))
* `azurerm_site_recovery_vmware_replication_policy_association` - update validation to correctly handle case ([#22443](https://github.com/hashicorp/terraform-provider-azurerm/issues/22443))

BUG FIXES:

* `azurerm_automation_dsc_configuration` - fixing an issue where `content_embedded` couldn't be deserialized ([#22403](https://github.com/hashicorp/terraform-provider-azurerm/issues/22403))
* `azurerm_data_factory_dataset_cosmosdb_sqlapi` - fix incorrect type/error message during read ([#22438](https://github.com/hashicorp/terraform-provider-azurerm/issues/22438))
* `azurerm_data_factory_dataset_mysql` - fix incorrect type/error message during read ([#22438](https://github.com/hashicorp/terraform-provider-azurerm/issues/22438))
* `azurerm_data_factory_dataset_postgresql` - fix incorrect type/error message during read ([#22438](https://github.com/hashicorp/terraform-provider-azurerm/issues/22438))
* `azurerm_logic_app_workflow` - prevent crash when `access_control` is empty block ([#22486](https://github.com/hashicorp/terraform-provider-azurerm/issues/22486))
* `azurerm_vpn_server_configuration` - prevent a potential panic when setting deprecated variables ([#22437](https://github.com/hashicorp/terraform-provider-azurerm/issues/22437))

## 3.64.0 (July 06, 2023)

FEATURES:

* **New Data Source:** `azurerm_automation_variables` ([#22216](https://github.com/hashicorp/terraform-provider-azurerm/issues/22216))
* **New Resource:** `azurerm_arc_private_link_scope` ([#22314](https://github.com/hashicorp/terraform-provider-azurerm/issues/22314))
* **New Resource:** `azurerm_kusto_cosmosdb_data_connection` ([#22295](https://github.com/hashicorp/terraform-provider-azurerm/issues/22295))
* **New Resource:** `azurerm_pim_active_role_assignment` ([#20731](https://github.com/hashicorp/terraform-provider-azurerm/issues/20731))
* **New Resource:** `azurerm_pim_eligible_role_assignment` ([#20731](https://github.com/hashicorp/terraform-provider-azurerm/issues/20731))

ENHANCEMENTS:

* dependencies: `web`: updating to API Version `2022-09-01` ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* dependencies: `cognitive`: updating to API Version `2023-05-01` ([#22223](https://github.com/hashicorp/terraform-provider-azurerm/issues/22223))
* dependencies: updating to `v1.53.0` of `google.golang.org/grpc` ([#22383](https://github.com/hashicorp/terraform-provider-azurerm/issues/22383))
* `azurerm_cognitive_deployment` - suppot for the `scale` block propeties `tier`, `size`, `family`, and `capacity` ([#22223](https://github.com/hashicorp/terraform-provider-azurerm/issues/22223))
* `azurerm_linux_function_app` - added support for the `public_network_access_enabled` property ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_linux_function_app_slot` - added support for the `public_network_access_enabled` property ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_linux_web_app` - added support for the `public_network_access_enabled` property ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_linux_web_app_slot`  - added support for the `public_network_access_enabled` property ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_windows_function_app` - added support for the `public_network_access_enabled` property ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_windows_function_app_slot` - added support for the `public_network_access_enabled` property
* `azurerm_windows_web_app` - added support for the `public_network_access_enabled` property ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_windows_web_app_slot` - added support for the `public_network_access_enabled` property ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_stream_analytics_output_blob` - increase the `batch_min_rows` property allowed values to `1000000` ([#22331](https://github.com/hashicorp/terraform-provider-azurerm/issues/22331))
* `azurerm_spring_cloud_gateway` - support for the the `allowed_origin_patterns` property ([#22317](https://github.com/hashicorp/terraform-provider-azurerm/issues/22317))

BUG FIXES:

* Data Source `azurerm_virtual_machine_scale_set` - prevent a nil pointer panic during reads ([#22335](https://github.com/hashicorp/terraform-provider-azurerm/issues/22335))
* `azurerm_application_insights_api_key` - prevent a nil pointer panic ([#22388](https://github.com/hashicorp/terraform-provider-azurerm/issues/22388))
* `azurerm_linux_function_app` - the `allowed_origins` property in the `cors` block now has a minimum entry count of `1` ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_linux_function_app_slot` - the `allowed_origins` property in the `cors` block now has a minimum entry count of `1` ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_linux_web_app` - the `allowed_origins` property in the `cors` block now has a minimum entry count of `1` ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_linux_web_app` - prevent a nil pointer panic in docker settings processing ([#22347](https://github.com/hashicorp/terraform-provider-azurerm/issues/22347))
* `azurerm_linux_web_app_slot`  - the `allowed_origins` property in the `cors` block now has a minimum entry count of `1` ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_private_dns_resolver_forwarding_rule_resource` - changing the `domain_name` property now creates a new resource ([#22375](https://github.com/hashicorp/terraform-provider-azurerm/issues/22375))
* `azurerm_windows_function_app` - the `allowed_origins` property in the `cors` block now has a minimum entry count of `1` ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_windows_function_app_slot` - the `allowed_origins` property in the `cors` block now has a minimum entry count of `1` ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_windows_web_app` - the `allowed_origins` property in the `cors` block now has a minimum entry count of `1` ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_windows_web_app_slot` - the `allowed_origins` property in the `cors` block now has a minimum entry count of `1` ([#22352](https://github.com/hashicorp/terraform-provider-azurerm/issues/22352))
* `azurerm_network_security_rule` - improve validation of the `name` property and prevent creation of resources that are broken ([#22336](https://github.com/hashicorp/terraform-provider-azurerm/issues/22336))

DEPRECATION:

* `media` - all resources and data sources are deprecated ahead of service being retired ([#22350](https://github.com/hashicorp/terraform-provider-azurerm/issues/22350))

## 3.63.0 (June 29, 2023)

FEATURES:

* **New Data Source:** `azurerm_network_manager_network_group` ([#22277](https://github.com/hashicorp/terraform-provider-azurerm/issues/22277))

BREAKING CHANGES:

* `azurerm_linux_web_app` - the `win32_status` property of the `status_code` block in `auto_heal` has changed from `string` to `int`.  ([#22003](https://github.com/hashicorp/terraform-provider-azurerm/issues/22003))
* `azurerm_linux_web_app_slot` -the `win32_status` property of the `status_code` block in `auto_heal` has changed from `string` to `int`.  ([#22003](https://github.com/hashicorp/terraform-provider-azurerm/issues/22003))
* `azurerm_windows_web_app` - the `win32_status` property of the `status_code` block in `auto_heal` has changed from `string` to `int`.  ([#22003](https://github.com/hashicorp/terraform-provider-azurerm/issues/22003))
* `azurerm_windows_web_app_slot` - the `win32_status` property of the `status_code` block in `auto_heal` has changed from `string` to `int`.  ([#22003](https://github.com/hashicorp/terraform-provider-azurerm/issues/22003))

ENHANCEMENTS:

* dependencies: updating to `v0.20230623.1103505` of `github.com/hashicorp/go-azure-sdk` ([#22263](https://github.com/hashicorp/terraform-provider-azurerm/issues/22263))
* dependencies: updating to `v0.57.0` of `github.com/hashicorp/go-azure-helpers` ([#22247](https://github.com/hashicorp/terraform-provider-azurerm/issues/22247))
* dependencies: `containers/containerinstance`: updating to API Version `2023-05-01` ([#22276](https://github.com/hashicorp/terraform-provider-azurerm/issues/22276))
* dependencies: `network/securityrules`: migrate to `go-azure-sdk` ([#22242](https://github.com/hashicorp/terraform-provider-azurerm/issues/22242))
* dependencies: `redis`: updating to API Version `2023-04-01` ([#22285](https://github.com/hashicorp/terraform-provider-azurerm/issues/22285))
* Data Source: `azurerm_kubernetes_cluster` - add support for the `custom_ca_trust_certificates_base64` property ([#22032](https://github.com/hashicorp/terraform-provider-azurerm/issues/22032))
* `azurerm_automation_software_update_configuration` - the `duration` property now defaults to `PT2H` as per the service. ([#22204](https://github.com/hashicorp/terraform-provider-azurerm/issues/22204))
* `azurerm_automation_software_update_configuration` - the `schedule` block is now limited to `1`, to match the API limit. ([#22204](https://github.com/hashicorp/terraform-provider-azurerm/issues/22204))
* `azurerm_automation_software_update_configuration` - the `schedule` block is now `Required` to match the API specification. The API  rejects requests that do not specify this block, with at least a `frequency` value. ([#22204](https://github.com/hashicorp/terraform-provider-azurerm/issues/22204))
* `azurerm_automation_software_update_configuration` - the `frequency` property is now a `Required` property of the `schedule` block. This is to match the minimum requirements of the API. ([#22204](https://github.com/hashicorp/terraform-provider-azurerm/issues/22204))
* `azurerm_automation_software_update_configuration` - the `pre_task` blocks are now limited to `1` to match the API. ([#22204](https://github.com/hashicorp/terraform-provider-azurerm/issues/22204))
* `azurerm_automation_software_update_configuration` - the `post_task` blocks are now limited to `1` to match the API. ([#22204](https://github.com/hashicorp/terraform-provider-azurerm/issues/22204))
* `azurerm_automation_software_update_configuration` - the `operating_system` property has been deprecated and is now controlled by the presence of either a `linux` or `windows` block. ([#22204](https://github.com/hashicorp/terraform-provider-azurerm/issues/22204))
* `azurerm_automation_software_update_configuration` - one of the `linux` or `windows` blocks must now be present. This is a requirement of the API, so is a non-breaking `Optional` to `Required` change. ([#22204](https://github.com/hashicorp/terraform-provider-azurerm/issues/22204))
* `azurerm_automation_software_update_configuration` - the `monthly_occurrence` blocks are now limited to `1` to match the API. ([#22204](https://github.com/hashicorp/terraform-provider-azurerm/issues/22204))
* `azurerm_container_app` - support for both system and user assigned identities at the same time ([#21149](https://github.com/hashicorp/terraform-provider-azurerm/issues/21149))
* `azurerm_key_vault_managed_hardware_security_module` - support for activating an HSM through `security_domain_key_vault_certificate_ids` ([#22162](https://github.com/hashicorp/terraform-provider-azurerm/issues/22162))
* `azurerm_kubernetes_cluster` - support for the `custom_ca_trust_certificates_base64` property ([#22032](https://github.com/hashicorp/terraform-provider-azurerm/issues/22032))
* `azurerm_kubernetes_cluster` - support for the `maintenance_window_auto_upgrade` block ([#21760](https://github.com/hashicorp/terraform-provider-azurerm/issues/21760))
* `azurerm_kubernetes_cluster` - support for the `maintenance_window_node_os` block ([#21760](https://github.com/hashicorp/terraform-provider-azurerm/issues/21760))
* `azurerm_monitor_aad_diagnostic_setting` - deprecate `log` in favour of `enabled_log` ([#21390](https://github.com/hashicorp/terraform-provider-azurerm/issues/21390))
* `azurerm_resource_group` - support for the `managed_by` property ([#22012](https://github.com/hashicorp/terraform-provider-azurerm/issues/22012))

BUG FIXES:

* `azurerm_automation_schedule` - prevent diffs for the `expiry_time` property when it hasn't been set in the user's configuration ([#21886](https://github.com/hashicorp/terraform-provider-azurerm/issues/21886))
* `azurerm_frontdoor` - throw an error if the resource cannot be found during an update ([#21975](https://github.com/hashicorp/terraform-provider-azurerm/issues/21975))
* `azurerm_image` - changing the `os_disk.size_gb` propety now creates a new resource ([#22272](https://github.com/hashicorp/terraform-provider-azurerm/issues/22272))
* `azurerm_kubernetes_cluster` - fix the validation for `node_os_channel_upgrade` block ([#22284](https://github.com/hashicorp/terraform-provider-azurerm/issues/22284))
* `azurerm_linux_virtual_machine` - raise an error if the resource cannot be found during an update ([#21975](https://github.com/hashicorp/terraform-provider-azurerm/issues/21975))
* `azurerm_linux_web_app` - deprecated the `docker_image` and `docker_image_tag` properties in favour of `docker_image_name`, `docker_registry_url`, `docker_registry_username`, and `docker_registry_password`. These settings now manage the respective `app_settings` values of the same name. ([#22003](https://github.com/hashicorp/terraform-provider-azurerm/issues/22003))
* `azurerm_linux_web_app_slot` - deprecated the `docker_image` and `docker_image_tag` properties in favour of `docker_image_name`, `docker_registry_url`, `docker_registry_username`, and `docker_registry_password`. These settings now manage the respective `app_settings` values of the same name.  ([#22003](https://github.com/hashicorp/terraform-provider-azurerm/issues/22003))
* `azurerm_site_recovery_replicated_vm` - set the `network_interface.failover_test_subnet_name`, `network_interface.failover_test_public_ip_address_id` and `network_interface.failover_test_static_ip` properties correctly ([#22217](https://github.com/hashicorp/terraform-provider-azurerm/issues/22217))
* `azurerm_ssh_public_key` - throw an error if the resource cannot be found during an update ([#21975](https://github.com/hashicorp/terraform-provider-azurerm/issues/21975))
* `azurerm_storage_share` - revert the resource ID format back to what it was previously due to a discrepancy in the API and Portal ([#22271](https://github.com/hashicorp/terraform-provider-azurerm/issues/22271))
* `azurerm_storage_account` - the `last_access_time_enabled` and `container_delete_retention_policy` properties are now supported in usgovernment ([#22273](https://github.com/hashicorp/terraform-provider-azurerm/issues/22273))
* `azurerm_windows_virtual_machine` - reaise an error if the resource cannot be found during an update ([#21975](https://github.com/hashicorp/terraform-provider-azurerm/issues/21975))
* `azurerm_windows_web_app` - deprecated the `docker_container_registry`, `docker_container_name`, and `docker_container_tag` properties in favour of `docker_image_name`, `docker_registry_url`, `docker_registry_username`, and `docker_registry_password`. These settings now manage the respective `app_settings` values of the same name.  ([#22003](https://github.com/hashicorp/terraform-provider-azurerm/issues/22003))
* `azurerm_windows_web_app_slot` - deprecated the `docker_container_registry`, `docker_container_name`, and `docker_container_tag` properties in favour of `docker_image_name`, `docker_registry_url`, `docker_registry_username`, and `docker_registry_password`. These settings now manage the respective `app_settings` values of the same name.  ([#22003](https://github.com/hashicorp/terraform-provider-azurerm/issues/22003))

## 3.62.1 (June 22, 2023)

BUG FIXES:

dependencies: `compute/marketplace_agreement` - Downgrade API version to 2015-06-01 ([#22264](https://github.com/hashicorp/terraform-provider-azurerm/issues/22264))

## 3.62.0 (June 22, 2023)

FEATURES:

* **New Resource:** `azurerm_new_relic_monitor` ([#21958](https://github.com/hashicorp/terraform-provider-azurerm/issues/21958))

ENHANCEMENTS:

* dependencies: updating to `v0.20230614.1151152` of `github.com/hashicorp/go-azure-sdk` ([#22176](https://github.com/hashicorp/terraform-provider-azurerm/issues/22176))
* dependencies: `compute/marketplace_agreement` - swap to use `hashicorp/go-azure-sdk` ([#21938](https://github.com/hashicorp/terraform-provider-azurerm/issues/21938))
* dependencies: `network/manager` - swap to use `hashicorp/go-azure-sdk` ([#22119](https://github.com/hashicorp/terraform-provider-azurerm/issues/22119))
* dependencies: `network/route` - swap to use `hashicorp/go-azure-sdk` ([#22227](https://github.com/hashicorp/terraform-provider-azurerm/issues/22227))
* `azurerm_cosmosdb_gremlin_graph` - support for the `analytical_storage_ttl` property ([#22179](https://github.com/hashicorp/terraform-provider-azurerm/issues/22179))
* `azurerm_kubernetes_cluster` - support for the value `AzureLinux` for the field `os_sku` within the `default_node_pool` block ([#22139](https://github.com/hashicorp/terraform-provider-azurerm/issues/22139))
* `azurerm_kubernetes_cluster` - support for the property `node_os_channel_upgrade` ([#22187](https://github.com/hashicorp/terraform-provider-azurerm/issues/22187))
* `azurerm_kubernetes_cluster_node_pool` - support for the value `AzureLinux` for the field `os_sku` ([#22139](https://github.com/hashicorp/terraform-provider-azurerm/issues/22139))
* `azurerm_monitor_workspace` - support for `public_network_access_enabled` ([#22197](https://github.com/hashicorp/terraform-provider-azurerm/issues/22197))
* `azurerm_virtual_hub` - support for `virtual_router_auto_scale_min_capacity` ([#21614](https://github.com/hashicorp/terraform-provider-azurerm/issues/21614))

BUG FIXES:

* `azurerm_application_insights_workbook` - the `display_name` property can now be updated ([#22148](https://github.com/hashicorp/terraform-provider-azurerm/issues/22148))
* `azurerm_bastion_host` - will now create a new resource when the `sku` property is downgraded ([#22147](https://github.com/hashicorp/terraform-provider-azurerm/issues/22147))
* `azurerm_container_app` - the `EmptyDir` property now functions ([#22196](https://github.com/hashicorp/terraform-provider-azurerm/issues/22196))
* `azurerm_kubernetes_cluster` - fix the validation preventing cluster's with `network_plugin_mode` set to `Overlay` due to a case change in the upstream API ([#22153](https://github.com/hashicorp/terraform-provider-azurerm/issues/22153))
* `azurerm_resource_deployment_script_*` - fix issue where `identity` wasn't specified but was being sent as `TypeNone` to the api ([#22165](https://github.com/hashicorp/terraform-provider-azurerm/issues/22165))
* `azurerm_bastion_host` - the `ip_configuration` propery is now required ([#22154](https://github.com/hashicorp/terraform-provider-azurerm/issues/22154))

## 3.61.0 (June 12, 2023)

FEATURES:

* **New Data Source:** `azurerm_mobile_network_packet_core_data_plane` ([#21053](https://github.com/hashicorp/terraform-provider-azurerm/issues/21053))
* **New Resource:** `azurerm_arc_machine_extension` ([#22051](https://github.com/hashicorp/terraform-provider-azurerm/issues/22051))
* **New Resource:** `azurerm_arc_kubernetes_flux_configuration` ([#21579](https://github.com/hashicorp/terraform-provider-azurerm/issues/21579))
* **New Resource:** `azurerm_kubernetes_flux_configuration` ([#21579](https://github.com/hashicorp/terraform-provider-azurerm/issues/21579))
* **New Resource:** `azurerm_mobile_network_packet_core_data_plane` ([#21053](https://github.com/hashicorp/terraform-provider-azurerm/issues/21053))

ENHANCEMENTS:

* dependencies: updating to `v0.20230530.1150329` of `github.com/tombuildsstuff/kermit` ([#21980](https://github.com/hashicorp/terraform-provider-azurerm/issues/21980))
* dependencies: `compute/gallery`: updating to API Version `2022-03-03` ([#21999](https://github.com/hashicorp/terraform-provider-azurerm/issues/21999))
* dependencies: `kusto`: updating to API Version `2022-12-29` ([#21961](https://github.com/hashicorp/terraform-provider-azurerm/issues/21961))
* Data Source `azurerm_site_recovery_replication_recovery_plan` - add support for `azure_to_azure_settings` block ([#22098](https://github.com/hashicorp/terraform-provider-azurerm/issues/22098))
* `compute`: updating to use API Version `2023-03-01` ([#21980](https://github.com/hashicorp/terraform-provider-azurerm/issues/21980))
* `containers`: updating to use API version `2023-04-02-preview` [22048]
* `managedidentity`: updating to use API Version `2023-01-31` ([#22102](https://github.com/hashicorp/terraform-provider-azurerm/issues/22102))
* `azurerm_backup_protected_vm` - support for the `protection_state` property ([#20608](https://github.com/hashicorp/terraform-provider-azurerm/issues/20608))
* `azurerm_batch_account` - the `public_network_access_enabled` property can now be updated ([#22095](https://github.com/hashicorp/terraform-provider-azurerm/issues/22095))
* `azurerm_batch_pool` - support for the `target_node_communication_mode` property ([#22094](https://github.com/hashicorp/terraform-provider-azurerm/issues/22094))
* `azurerm_automanage_configuration` - support for the `log_analytics_enabled` property ([#22121](https://github.com/hashicorp/terraform-provider-azurerm/issues/22121))
* `azurerm_nginx_certificate` - the `key_virtual_path`, `certificate_virtual_path`, and `key_vault_secret_id` proeprties can now be updated ([#22100](https://github.com/hashicorp/terraform-provider-azurerm/issues/22100))
* `azurerm_spring_cloud_gateway` - support for the `client_authentication` property ([#22016](https://github.com/hashicorp/terraform-provider-azurerm/issues/22016))

BUG FIXES:

* `azurerm_databricks_workspace_data_source` - correctly set the `managed_idnetity_id` attribute ([#22021](https://github.com/hashicorp/terraform-provider-azurerm/issues/22021))

## 3.60.0 (June 08, 2023)

NOTES:

* `azurerm_security_center_subscription_pricing` - upon deletion the pricing tier will now reset to `Free` tier ([#21437](https://github.com/hashicorp/terraform-provider-azurerm/issues/21437))

ENHANCEMENTS:

* dependencies: `batch`: updating to API Version `2022-10-01` ([#21962](https://github.com/hashicorp/terraform-provider-azurerm/issues/21962))
* dependencies: `loadtest`: updating to API Version `2022-12-01` ([#22091](https://github.com/hashicorp/terraform-provider-azurerm/issues/22091))
* provider: adding the `client_id_file_path` and `client_secret_file_path` provider properties ([#21764](https://github.com/hashicorp/terraform-provider-azurerm/issues/21764))
* `data.azurerm_key_vault_encrypted_value` - now exports the `decoded_plain_text_value` [attribute GH-21682]
* `azurerm_automanage_configuration` - support for the `backup` and `azure_security_baseline` blocks ([#22081](https://github.com/hashicorp/terraform-provider-azurerm/issues/22081))
* `azurerm_app_configuration` - support toggling of user permission error on soft deleted stores through `app_configuration.recover_soft_deleted` feature flag ([#19661](https://github.com/hashicorp/terraform-provider-azurerm/issues/19661))
* `azurerm_backup_policy_file_share` - support for day-based retention policies and hourly backups ([#21529](https://github.com/hashicorp/terraform-provider-azurerm/issues/21529))
* `azurerm_linux_function_app` - support for Python `3.11` for Linux function app ([#21956](https://github.com/hashicorp/terraform-provider-azurerm/issues/21956))
* `azurerm_linux_function_app_slot` - support for Python `3.11` for Linux function app ([#21956](https://github.com/hashicorp/terraform-provider-azurerm/issues/21956))
* `azurerm_monitor_autoscale_setting` - support for the `predictive` block ([#22038](https://github.com/hashicorp/terraform-provider-azurerm/issues/22038))
* `azurerm_machine_learning_compute_instance` - support for the `node_public_ip_enabled` property ([#22063](https://github.com/hashicorp/terraform-provider-azurerm/issues/22063))
* `azurerm_spring_cloud_service` - support for the `container_registry` block ([#22017](https://github.com/hashicorp/terraform-provider-azurerm/issues/22017))
* `azurerm_site_recovery_replication_recovery_plan` - the order of the `pre_action` and `post_action` properties is now respected ([#22019](https://github.com/hashicorp/terraform-provider-azurerm/issues/22019))

BUG FIXES:

* `azurerm_hdinsight_interactive_query_cluster` - deprecating the `*_node.0.autoscale.0.capacity` property ([#21981](https://github.com/hashicorp/terraform-provider-azurerm/issues/21981))
* `azurerm_key_vault_key` - allow the `rotation_policy` property to be removed ([#21935](https://github.com/hashicorp/terraform-provider-azurerm/issues/21935))
* `azurerm_mssql_server` - fix issue where the `minimum_tls_version` property is being returned as `None` instead of `Disabled` ([#22067](https://github.com/hashicorp/terraform-provider-azurerm/issues/22067))
* `azurerm_sentinel_data_connector_microsoft_threat_intelligence` - the `bing_safety_phishing_url_lookback_date` property has been deprecated ([#21954](https://github.com/hashicorp/terraform-provider-azurerm/issues/21954))

## 3.59.0 (June 01, 2023)

FEATURES:

* **New Data Source:** `azurerm_arc_machine` ([#21796](https://github.com/hashicorp/terraform-provider-azurerm/issues/21796))
* **New Resource:** `azurerm_automanage_configuration` ([#21490](https://github.com/hashicorp/terraform-provider-azurerm/issues/21490))

ENHANCEMENTS:

* dependencies: updating to `v0.20230523.1140858` of `github.com/hashicorp/go-azure-sdk` ([#21910](https://github.com/hashicorp/terraform-provider-azurerm/issues/21910))
* dependencies: `azurem_monitor_action_group` - upgrading `actiongroupsapis` from `2021-09-01` to `2023-01-01` ([#21948](https://github.com/hashicorp/terraform-provider-azurerm/issues/21948))
* dependencies: `policy.guestconfigurationassignments`: migrate to `hashicorp/go-azure-sdk` ([#21927](https://github.com/hashicorp/terraform-provider-azurerm/issues/21927))
* dependencies: `azurerm_monitor_autoscale_setting`  upgrade API version from to `2023-05-01-preview` ([#21953](https://github.com/hashicorp/terraform-provider-azurerm/issues/21953))
* `data.azurerm_linux_web_app` - now exports the `availability` and `usage` attributes ([#21945](https://github.com/hashicorp/terraform-provider-azurerm/issues/21945))
* `data.azurerm_linux_function_app` - now exports the `availability` and `usage` attributes ([#21945](https://github.com/hashicorp/terraform-provider-azurerm/issues/21945))
* `data.azurerm_cdn_frontdoor_secret` - now exports the `expiration_date` attribute ([#21945](https://github.com/hashicorp/terraform-provider-azurerm/issues/21945))
* `data.azurerm_virtual_machine` - now exports the `power_state` ([#21945](https://github.com/hashicorp/terraform-provider-azurerm/issues/21945))
* `data.azurerm_virtual_machine_scale_set` -  now exports the `power_state` attribute ([#21945](https://github.com/hashicorp/terraform-provider-azurerm/issues/21945))
* `data.azurerm_azurerm_resources` - now exports the `resource_group_name` attribute for each resource ([#21676](https://github.com/hashicorp/terraform-provider-azurerm/issues/21676))
* `security.watchitems` - updating to use `hashicorp/go-azure-sdk` ([#21944](https://github.com/hashicorp/terraform-provider-azurerm/issues/21944))
* `azurerm_cosmosdb_account` - support new capabilities for `MongoDB` ([#21974](https://github.com/hashicorp/terraform-provider-azurerm/issues/21974))
* `azurerm_kubernetes_cluster` - the properties `enable_host_encryption`, `enable_node_public_ip`, `kubelet_config`, `linux_os_config`, `max_pods`, `node_taints`, `only_critical_addons_enabled`, `os_disk_size_gb`, `os_disk_type`, `os_sku`, `pod_subnet_id`, `ultra_ssd_enabled`, `vnet_subnet_id` and `zones` are now updateable through cycling of the system node pool ([#21719](https://github.com/hashicorp/terraform-provider-azurerm/issues/21719))
* `azurerm_machine_learning_compute_cluster` - add support for the `node_public_ip_enabled` property ([#21377](https://github.com/hashicorp/terraform-provider-azurerm/issues/21377))
* `azurerm_nginx_certificate` - `key_vault_secret_id` now accepts version-less key vault secret ids ([#21949](https://github.com/hashicorp/terraform-provider-azurerm/issues/21949))
* `azurerm_postgresql_flexible_server` - add support for `version` value `15` ([#21934](https://github.com/hashicorp/terraform-provider-azurerm/issues/21934))
* `azurerm_shared_image_version` - now exports the `id` property ([#22006](https://github.com/hashicorp/terraform-provider-azurerm/issues/22006))
* `azurerm_spring_cloud_certificate` - support for the `exclude_private_key` property ([#21942](https://github.com/hashicorp/terraform-provider-azurerm/issues/21942))
* `azurerm_spring_cloud_customized_accelerator` - support for the `ca_certificate_id` property ([#21943](https://github.com/hashicorp/terraform-provider-azurerm/issues/21943))

BUG FIXES:

* `azurerm_app_configuration` - prevent errors when deleting by checking that the name of the app configuration store has been released ([#21750](https://github.com/hashicorp/terraform-provider-azurerm/issues/21750))
* `azurerm_express_route_port_authorization` - add a lock when create/update/delete authorization of express route port ([#21959](https://github.com/hashicorp/terraform-provider-azurerm/issues/21959))
* `azurerm_kubernetes_cluster` - recompute the field `oidc_issuer_url` if the value of `oidc_issuer_enabled` has changed ([#21911](https://github.com/hashicorp/terraform-provider-azurerm/issues/21911))
* `azurerm_kubernetes_cluster` - set correct value for `default_node_pool.os_sku` when resizing the `default_node_pool` ([#21976](https://github.com/hashicorp/terraform-provider-azurerm/issues/21976))
* `azurerm_postgresql_flexible_server` - fix issue updating `storage_mb` and `backup_retention_days` together ([#21987](https://github.com/hashicorp/terraform-provider-azurerm/issues/21987))

## 3.58.0 (May 25, 2023)

FEATURES:

* **New data Source:** `azurerm_mobile_network_packet_core_control_plane` ([#21071](https://github.com/hashicorp/terraform-provider-azurerm/issues/21071))
* **New Resource:** `azurerm_cosmosdb_mongo_role_definition` ([#21754](https://github.com/hashicorp/terraform-provider-azurerm/issues/21754))
* **New Resource:** `azurerm_cosmosdb_mongo_user_definition` ([#21914](https://github.com/hashicorp/terraform-provider-azurerm/issues/21914))
* **New Resource:** `azurerm_iothub_file_upload` ([#20668](https://github.com/hashicorp/terraform-provider-azurerm/issues/20668))
* **New Resource:** `azurerm_mobile_network_packet_core_control_plane` ([#21071](https://github.com/hashicorp/terraform-provider-azurerm/issues/21071))
* **New Resource:** `azurerm_mysql_flexible_server_active_directory_administrator` ([#21786](https://github.com/hashicorp/terraform-provider-azurerm/issues/21786))
* **New Resource:** `azurerm_monitor_alert_prometheus_rule_group` ([#21751](https://github.com/hashicorp/terraform-provider-azurerm/issues/21751))
* **New Resource:** `azurerm_recovery_services_vault_resource_guard_association` ([#21712](https://github.com/hashicorp/terraform-provider-azurerm/issues/21712))
* **New Resource:** `azurerm_site_recovery_hyperv_network_mapping` ([#21788](https://github.com/hashicorp/terraform-provider-azurerm/issues/21788))
* **New Resource:** `azurerm_site_recovery_vmware_replication_policy_association` ([#21389](https://github.com/hashicorp/terraform-provider-azurerm/issues/21389))

ENHANCEMENTS:

* dependencies: updating to `v0.20230523.1080931` of `github.com/hashicorp/go-azure-sdk` ([#21898](https://github.com/hashicorp/terraform-provider-azurerm/issues/21898))
* dependencies: updating to `v0.20230518.1143920` of `github.com/tombuildsstuff/kermit` ([#21899](https://github.com/hashicorp/terraform-provider-azurerm/issues/21899))
* dependencies: `azurerm_monitor_autoscale_setting`  upgrade API version from `2015-04-01` to `2022-10-01` ([#21887](https://github.com/hashicorp/terraform-provider-azurerm/issues/21887))
* `cosmosdb.gremlin`: updating to use `hashicorp/go-azure-sdk` and api version `2023-04-15` ([#21813](https://github.com/hashicorp/terraform-provider-azurerm/issues/21813))
* `cosmosdb.sql_container`: updating to use `hashicorp/go-azure-sdk` and api version `2023-04-15` ([#21813](https://github.com/hashicorp/terraform-provider-azurerm/issues/21813))
* `nginx`: updating to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#21810](https://github.com/hashicorp/terraform-provider-azurerm/issues/21810))
* `portal`: updating to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#21810](https://github.com/hashicorp/terraform-provider-azurerm/issues/21810))
* `redis`: updating to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#21810](https://github.com/hashicorp/terraform-provider-azurerm/issues/21810))
* `appplatform`: updating to API Version `2023-03-01-preview` ([#21404](https://github.com/hashicorp/terraform-provider-azurerm/issues/21404))
* `redisenterprise`: updating to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#21810](https://github.com/hashicorp/terraform-provider-azurerm/issues/21810))
* `azurerm_cosmosdb_account` - fix for upstream Microsoft API issue where updating `identity` and `default_identity` at the same time silently fails ([#21780](https://github.com/hashicorp/terraform-provider-azurerm/issues/21780))
* `azurerm_monitor_activity_log_alert` - support for the `levels`, `resource_providers`, `resource_types`, `resource_groups`, `resource_ids`, `statuses`, and `sub_statuses` properties ([#21367](https://github.com/hashicorp/terraform-provider-azurerm/issues/21367))
* `azurerm_media_transform` - support for the `experimental_options` property ([#21873](https://github.com/hashicorp/terraform-provider-azurerm/issues/21873))
* `azurerm_backup_policy_vm` - support for the `days` and `include_last_days` properties ([#21434](https://github.com/hashicorp/terraform-provider-azurerm/issues/21434))
* `azurerm_subnet` - the `name` property within the `subnet_delegation` block can now be set to `Microsoft.App/environments` ([#21893](https://github.com/hashicorp/terraform-provider-azurerm/issues/21893))
* `azurerm_subnet_service_endpoint_policy` - support for the `service` property ([#21865](https://github.com/hashicorp/terraform-provider-azurerm/issues/21865))
* `azurerm_signalr_service` - support for the `user_assigned_identity_id` property ([#21055](https://github.com/hashicorp/terraform-provider-azurerm/issues/21055))
* `azurerm_site_recovery_replication_recovery_plan` - support for the `azure_to_azure_settings` block ([#21666](https://github.com/hashicorp/terraform-provider-azurerm/issues/21666))
* `azurerm_cosmosdb_postgresql_cluster` - the `citus_version` property now supports `11.3` ([#21916](https://github.com/hashicorp/terraform-provider-azurerm/issues/21916))

BUG FIXES:

* Data Source: `azurerm_kubernetes_cluster` - prevent a panic when some values returned are nil ([#21867](https://github.com/hashicorp/terraform-provider-azurerm/issues/21867))
* `azurerm_application_insights_web_test` - normalizing the value for the `application_insights_id` property ([#21837](https://github.com/hashicorp/terraform-provider-azurerm/issues/21837))
* `azurerm_api_management` - correctly configure the `triple_des_ciphers_enabled` value ([#21789](https://github.com/hashicorp/terraform-provider-azurerm/issues/21789))
* `azurerm_key_vault` - during creation the`createMode` will now be set to `default` instead of `nil` ([#21668](https://github.com/hashicorp/terraform-provider-azurerm/issues/21668))
* `azurerm_spring_cloud_gateway_route_config` -  the `filters` and `predicates` properties will now be omitted when not specified ([#21745](https://github.com/hashicorp/terraform-provider-azurerm/issues/21745))
* `azurerm_subnet` - permit `Microsoft.BareMetal/AzureHostedService` as an option for the `service_delegation` property ([#21871](https://github.com/hashicorp/terraform-provider-azurerm/issues/21871))

## 3.57.0 (May 19, 2023)

FEATURES:

* **New Data Source:** `azurerm_virtual_hub_connection` ([#21681](https://github.com/hashicorp/terraform-provider-azurerm/issues/21681))

ENHANCEMENTS:

* `synapse`: refactoring to use `tombuildsstuff/kermit` rather than `Azure/azure-sdk-for-go` for Data Plane ([#21792](https://github.com/hashicorp/terraform-provider-azurerm/issues/21792))
* `azurerm_batch_account` - support versionless keys for CMK ([#21677](https://github.com/hashicorp/terraform-provider-azurerm/issues/21677))
* `azurerm_kubernetes_cluster` - changing the `http_proxy_config.no_proxy` no longer creates a new resource ([#21793](https://github.com/hashicorp/terraform-provider-azurerm/issues/21793))
* `azurerm_media_transform` - support for the `jpg_image` and `png_image` blocks within the `custom_preset` block ([#21709](https://github.com/hashicorp/terraform-provider-azurerm/issues/21709))
* `azurerm_recovery_services_vault` - support the `monitoring` block ([#21691](https://github.com/hashicorp/terraform-provider-azurerm/issues/21691))

BUG FIXES:

* `data.azurerm_kubernetes_cluster` - prevent a panic when some values returned are nil ([#21850](https://github.com/hashicorp/terraform-provider-azurerm/issues/21850))

## 3.56.0 (May 11, 2023)

FEATURES:

* **New Resource:** `azurerm_cosmosdb_postgresql_coordinator_configuration` ([#21595](https://github.com/hashicorp/terraform-provider-azurerm/issues/21595))
* **New Resource:** `azurerm_cosmosdb_postgresql_node_configuration` ([#21596](https://github.com/hashicorp/terraform-provider-azurerm/issues/21596))
* **New Resource:** `azurerm_cosmosdb_postgresql_role` ([#21597](https://github.com/hashicorp/terraform-provider-azurerm/issues/21597))
* **New Resource:** `azurerm_monitor_workspace` ([#21598](https://github.com/hashicorp/terraform-provider-azurerm/issues/21598))
* **New Resource:** `azurerm_network_manager_deployment` ([#20451](https://github.com/hashicorp/terraform-provider-azurerm/issues/20451))

ENHANCEMENTS:

* dependencies: updating to `v0.56.0` of `github.com/hashicorp/go-azure-helpers` ([#21725](https://github.com/hashicorp/terraform-provider-azurerm/issues/21725))
* dependencies: updating to `v0.20230511.1094507` of `github.com/hashicorp/go-azure-sdk` ([#21759](https://github.com/hashicorp/terraform-provider-azurerm/issues/21759))
* provider: improving the error messages when parsing a Resource ID and the ID doesn't match what's expected ([#21725](https://github.com/hashicorp/terraform-provider-azurerm/issues/21725))
* provider: Resource Provider Registration now uses API Version `2022-09-01` ([#21695](https://github.com/hashicorp/terraform-provider-azurerm/issues/21695))
* provider: updating the `IsAzureStack` check to use `hashicorp/go-azure-sdk` rather than relying on the environment from `Azure/go-autorest` ([#21697](https://github.com/hashicorp/terraform-provider-azurerm/issues/21697))
* `appconfiguration`: updating to API Version `2023-03-01` ([#21660](https://github.com/hashicorp/terraform-provider-azurerm/issues/21660))
* `keyvault`: refactoring to use `hashicorp/go-azure-sdk` ([#21621](https://github.com/hashicorp/terraform-provider-azurerm/issues/21621))
* `azurerm_machine_learning_workspace` - exporting `workspace_id` ([#21746](https://github.com/hashicorp/terraform-provider-azurerm/issues/21746))
* `azurerm_mssql_server` - expose the ability to enable `Transparent Data Encryption` using a `Customer Managed Key` during server deployment ([#21704](https://github.com/hashicorp/terraform-provider-azurerm/issues/21704))
* `azurerm_orbital_contact_profile` - `ip_address` is now optional ([#21721](https://github.com/hashicorp/terraform-provider-azurerm/issues/21721))

BUG FIXES:

* provider: fixing a bug where we would invoke but not poll for the Registration State during automatic Resource Provider Registration ([#21695](https://github.com/hashicorp/terraform-provider-azurerm/issues/21695))
* `azurerm_app_configuration`: handling an API bug where when polling for `PurgeDeleted` returns a 404 rather the payload for a long-running operation ([#21665](https://github.com/hashicorp/terraform-provider-azurerm/issues/21665))
* `azurerm_api_management_api` - fixing a bug where an empty `contact` bug would cause a crash ([#21740](https://github.com/hashicorp/terraform-provider-azurerm/issues/21740))
* `azurerm_eventhub_namespace` - add locks and remove unneeded WaitForState functions ([#21656](https://github.com/hashicorp/terraform-provider-azurerm/issues/21656))
* `azurerm_machine_learning_workspace` - parse `key_vault_id` insensitively ([#21684](https://github.com/hashicorp/terraform-provider-azurerm/issues/21684))
* `azurerm_monitor_action_group` - further expand ExactlyOneOf logic for `event_hub_receiver` attributes ([#21735](https://github.com/hashicorp/terraform-provider-azurerm/issues/21735))
* `azurerm_monitor_metric_alert` - fix regression by using `SingleResourceMultiMetricCriteria` for new metric alerts  ([#21658](https://github.com/hashicorp/terraform-provider-azurerm/issues/21658))
* `azurerm_service_fabric_managed_cluster` - fixing a bug where `certificates` within the `vm_secrets` block wouldn't be set into the state ([#21680](https://github.com/hashicorp/terraform-provider-azurerm/issues/21680))
* `azurerm_storage_share` - correct resource ID segment from `fileshares` to `shares` ([#21645](https://github.com/hashicorp/terraform-provider-azurerm/issues/21645))
* `azurerm_virtual_machine_scale_set`,  - - support specifying `ultra_ssd_disk_iops_read_write` and `ultra_ssd_disk_mbps_read_write` for `PremiumV2_LRS` ([#21530](https://github.com/hashicorp/terraform-provider-azurerm/issues/21530)) 


## 3.55.0 (May 04, 2023)

FEATURES:

* **New Data Source:** `azurerm_kubernetes_node_pool_snapshot` ([#21511](https://github.com/hashicorp/terraform-provider-azurerm/issues/21511))
* **New Resource:** `azurerm_cosmosdb_postgresql_firewall_rule` ([#21599](https://github.com/hashicorp/terraform-provider-azurerm/issues/21599))

ENHANCEMENTS:

* `appconfiguration`: refactoring to use `tombuildsstuff/kermit` rather than an embedded SDK ([#21623](https://github.com/hashicorp/terraform-provider-azurerm/issues/21623))
* `recoveryservicesbackup` - updating to use API Version `2023-02-01` ([#21575](https://github.com/hashicorp/terraform-provider-azurerm/issues/21575))
* `azurerm_kubernetes_cluster_node_pool` - support for the `snapshot_id` property ([#21511](https://github.com/hashicorp/terraform-provider-azurerm/issues/21511))

BUG FIXES:

* Data Source: `azurerm_healthcare_fhir_service` - `identity` now exports both `SystemAssigned` and `UserAssigned` identities ([#21594](https://github.com/hashicorp/terraform-provider-azurerm/issues/21594))
* `azurerm_local_network_gateway` - validating that `address_space` isn't set to an empty string ([#21566](https://github.com/hashicorp/terraform-provider-azurerm/issues/21566))
* `azurerm_log_analytics_cluster` -  Add locks and remove unneeded WaitForState checks ([#21631](https://github.com/hashicorp/terraform-provider-azurerm/issues/21631))
* `azurerm_log_analytics_cluster_customer_managed_key` - Add locks and remove unneeded WaitForState checks ([#21631](https://github.com/hashicorp/terraform-provider-azurerm/issues/21631))
* `azurerm_managed_disk` - now detaches when `disk_size_gb` increases from below `4095` to above `4095` ([#21620](https://github.com/hashicorp/terraform-provider-azurerm/issues/21620))
* Service `mssqlmanagedinstance` - add initialize of `client.MSSQLManagedInstance` to fix panic ([#21657](https://github.com/hashicorp/terraform-provider-azurerm/issues/21657))
* `azurerm_virtual_machine` - fixing a regression when parsing the OS Disk ID from the Azure API ([#21606](https://github.com/hashicorp/terraform-provider-azurerm/issues/21606))
* `azurerm_virtual_machine` - fixing a regression when parsing the Data Disk ID from the Azure API ([#21606](https://github.com/hashicorp/terraform-provider-azurerm/issues/21606))

## 3.54.0 (April 27, 2023)

BREAKING CHANGES:

* `azurerm_attestation_provider` - the field `policy` is deprecated and non-functional due to a design issue with the original resource (where this wasn't retrieved from the Azure API and thus wasn't exposed correctly) - this has been superseded by the fields `open_enclave_policy_base64`, `sgx_enclave_policy_base64` and `tpm_policy_base64`. ([#21524](https://github.com/hashicorp/terraform-provider-azurerm/issues/21524))

FEATURES:

* **New Resource:** `azurerm_arc_kubernetes_cluster_extension` ([#21310](https://github.com/hashicorp/terraform-provider-azurerm/issues/21310))
* **New Resource:** `azurerm_cosmosdb_postgresql_cluster` ([#21090](https://github.com/hashicorp/terraform-provider-azurerm/issues/21090))
* **New Resource:** `azurerm_email_communication_service` ([#21526](https://github.com/hashicorp/terraform-provider-azurerm/issues/21526))
* **New Resource:** `azurerm_kubernetes_cluster_extension` ([#21310](https://github.com/hashicorp/terraform-provider-azurerm/issues/21310))
* **New Resource:** `azurerm_netapp_volume_group_sap_hana` ([#21290](https://github.com/hashicorp/terraform-provider-azurerm/issues/21290))
* **New Resource:** `azurerm_storage_mover_project` ([#21477](https://github.com/hashicorp/terraform-provider-azurerm/issues/21477))
* **New Resource:** `azurerm_storage_mover_job_definition` ([#21514](https://github.com/hashicorp/terraform-provider-azurerm/issues/21514))

ENHANCEMENTS:

* dependencies: updating to `v0.20230427.1112058` of `github.com/hashicorp/go-azure-sdk` ([#21583](https://github.com/hashicorp/terraform-provider-azurerm/issues/21583))
*  `security`: updating to API Version `2023-01-01` ([#21531](https://github.com/hashicorp/terraform-provider-azurerm/issues/21531))
* Data Source: `azurerm_virtual_network_gateway` - add support for the field `private_ip_address` ([#21432](https://github.com/hashicorp/terraform-provider-azurerm/issues/21432))
* `azurerm_active_directory_domain_service` - `domain_name` now supports a length up to 30 characters ([#21555](https://github.com/hashicorp/terraform-provider-azurerm/issues/21555))
* `azurerm_attestation_provider` - adding support for the field `open_enclave_policy_base64`, `sgx_enclave_policy_base64` and `tpm_policy_base64` ([#21524](https://github.com/hashicorp/terraform-provider-azurerm/issues/21524))
* `azurerm_attestation_provider` - adding support for the field `sgx_enclave_policy_base64` ([#21524](https://github.com/hashicorp/terraform-provider-azurerm/issues/21524))
* `azurerm_attestation_provider` - adding support for the field `tpm_policy_base64` ([#21524](https://github.com/hashicorp/terraform-provider-azurerm/issues/21524))
* `azurerm_billing_account_cost_management_export` - the field `time_frame` can now be set to `TheLast7Days` ([#21528](https://github.com/hashicorp/terraform-provider-azurerm/issues/21528))
* `azurerm_firewall_policy_rule_collection_group` - the fields `source_addresses` and `destination_addresses` now accepts an IPv4 range ([#21542](https://github.com/hashicorp/terraform-provider-azurerm/issues/21542))
* `azurerm_kubernetes_cluster` - add support for the `service_mesh_profile` block ([#21516](https://github.com/hashicorp/terraform-provider-azurerm/issues/21516))
* `azurerm_resource_group_cost_management_export` - the field `time_frame` can now be set to `TheLast7Days` ([#21528](https://github.com/hashicorp/terraform-provider-azurerm/issues/21528))
* `azurerm_search_service` - adding support for `authentication_failure_mode ` ([#21323](https://github.com/hashicorp/terraform-provider-azurerm/issues/21323))
* `azurerm_search_service` - adding support for `customer_managed_key_enforcement_enabled ` ([#21323](https://github.com/hashicorp/terraform-provider-azurerm/issues/21323))
* `azurerm_search_service` - adding support for `hosting_mode ` ([#21323](https://github.com/hashicorp/terraform-provider-azurerm/issues/21323))
* `azurerm_search_service` - adding support for `local_authentication_enabled` ([#21323](https://github.com/hashicorp/terraform-provider-azurerm/issues/21323))
* `azurerm_search_service` - support for setting `sku` to `StorageOptimizedL2` ([#21323](https://github.com/hashicorp/terraform-provider-azurerm/issues/21323))
* `azurerm_subscription_cost_management_export` - the field `time_frame` can now be set to `TheLast7Days` ([#21528](https://github.com/hashicorp/terraform-provider-azurerm/issues/21528))

BUG FIXES:

* **Provider:** fix an authentication bug when specifying `auxiliary_tenant_ids` whilst authenticating using Azure CLI ([#21583](https://github.com/hashicorp/terraform-provider-azurerm/issues/21583))
* `azurerm_attestation_provider` - the field `policy` is deprecated and non-functional - instead please use the fields  `open_enclave_policy_base64`, `sgx_enclave_policy_base64` and `tpm_policy_base64` ([#21524](https://github.com/hashicorp/terraform-provider-azurerm/issues/21524))
* `azurerm_mysql_flexible_server` - fix issue where `identity` was not being removed properly on updates ([#21533](https://github.com/hashicorp/terraform-provider-azurerm/issues/21533))
* `azurerm_search_service` - updating the default value for `partition_count` to `1` to match the API ([#21323](https://github.com/hashicorp/terraform-provider-azurerm/issues/21323))
* `azurerm_search_service` - updating the default value for `replica_count` to `1` to match the API ([#21323](https://github.com/hashicorp/terraform-provider-azurerm/issues/21323))
* `azurerm_search_service` - the field `allowed_ips` is now a Set rather than a List ([#21323](https://github.com/hashicorp/terraform-provider-azurerm/issues/21323))

## 3.53.0 (April 20, 2023)

FEATURES:

* **New Resource:** `azurerm_cost_management_scheduled_action` ([#21325](https://github.com/hashicorp/terraform-provider-azurerm/issues/21325))
* **New Resource:** `azurerm_storage_mover_agent` ([#21273](https://github.com/hashicorp/terraform-provider-azurerm/issues/21273))
* **New Resource:** `azurerm_storage_mover_source_endpoint` ([#21449](https://github.com/hashicorp/terraform-provider-azurerm/issues/21449))
* **New Resource:** `azurerm_storage_mover_target_endpoint` ([#21449](https://github.com/hashicorp/terraform-provider-azurerm/issues/21449))

ENHANCEMENTS:

* `advisor` - refactoring to use `hashicorp/go-azure-sdk` ([#21307](https://github.com/hashicorp/terraform-provider-azurerm/issues/21307))
* `healthcare`: refactoring to use `hashicorp/go-azure-sdk` ([#21327](https://github.com/hashicorp/terraform-provider-azurerm/issues/21327))
* `hpccache` - refactoring to use `hashicorp/go-azure-sdk` ([#21303](https://github.com/hashicorp/terraform-provider-azurerm/issues/21303))
* `logz` - refactoring to use `hashicorp/go-azure-sdk` ([#21321](https://github.com/hashicorp/terraform-provider-azurerm/issues/21321))
* `hpccache`: updating to API Version `2023-01-01` ([#21459](https://github.com/hashicorp/terraform-provider-azurerm/issues/21459))
* `orbital`: updating to API Version `2022-11-01` ([#21405](https://github.com/hashicorp/terraform-provider-azurerm/issues/21405))
* `vmware`: updating to API Version `2022-05-01` ([#21458](https://github.com/hashicorp/terraform-provider-azurerm/issues/21458))
* `azurerm_attestation_provider` - support for the `policy` block ([#20972](https://github.com/hashicorp/terraform-provider-azurerm/issues/20972))
* `azurerm_linux_function_app` - support for the `hosting_environment_id` property ([#20471](https://github.com/hashicorp/terraform-provider-azurerm/issues/20471))
* `azurerm_linux_function_app_slot` - support for the `hosting_environment_id` property ([#20471](https://github.com/hashicorp/terraform-provider-azurerm/issues/20471))
* `azurerm_linux_web_app` - support for the `hosting_environment_id` property ([#20471](https://github.com/hashicorp/terraform-provider-azurerm/issues/20471))
* `azurerm_linux_web_app` - support `PHP 8.2` for the `application_stack` property ([#21420](https://github.com/hashicorp/terraform-provider-azurerm/issues/21420))
* `azurerm_linux_web_app_slot`  - support for the `hosting_environment_id` property ([#20471](https://github.com/hashicorp/terraform-provider-azurerm/issues/20471))
* `azurerm_linux_web_app_slot` support `PHP 8.2` for the `application_stack` property ([#21420](https://github.com/hashicorp/terraform-provider-azurerm/issues/21420))
* `azurerm_signalr_service` - add addtional valid values for `sku.0.capacity` ([#21494](https://github.com/hashicorp/terraform-provider-azurerm/issues/21494))
* `azurerm_windows_function_app` - support for the `hosting_environment_id` property ([#20471](https://github.com/hashicorp/terraform-provider-azurerm/issues/20471))
* `azurerm_windows_function_app_slot` - support for the `hosting_environment_id` property ([#20471](https://github.com/hashicorp/terraform-provider-azurerm/issues/20471))
* `azurerm_windows_web_app` - support for the `hosting_environment_id` property ([#20471](https://github.com/hashicorp/terraform-provider-azurerm/issues/20471))
* `azurerm_windows_web_app_slot` - support for the `hosting_environment_id` property ([#20471](https://github.com/hashicorp/terraform-provider-azurerm/issues/20471))

BUG FIXES: 

* `azurerm_cdn_endpoint` - remove the length limit for the `query_string` property ([#21474](https://github.com/hashicorp/terraform-provider-azurerm/issues/21474))
* `azurerm_cognitive_account` - mark the `custom_question_answering_search_service_key` property as sensitive ([#21469](https://github.com/hashicorp/terraform-provider-azurerm/issues/21469))
* `azurerm_monitor_metric_alert` - fix crash when the `dynamic_criteria.0.ignore_data_before` property  isn't set ([#21446](https://github.com/hashicorp/terraform-provider-azurerm/issues/21446))
* `azurerm_postgresql_flexible_server` - correctly set the `point_in_time_restore_time_in_utc` property ([#21501](https://github.com/hashicorp/terraform-provider-azurerm/issues/21501))
* `azurerm_search_service` - mark the `primary_key` and `secondary_key` properties as sensitive ([#21469](https://github.com/hashicorp/terraform-provider-azurerm/issues/21469))

## 3.52.0 (April 13, 2023)

ENHANCEMENTS:

* `containerRegistry` - refactoring to use `hashicorp/go-azure-sdk` ([#21344](https://github.com/hashicorp/terraform-provider-azurerm/issues/21344))
* `monitor` - refactoring to use `hashicorp/go-azure-sdk` ([#21392](https://github.com/hashicorp/terraform-provider-azurerm/issues/21392))
* `recoveryServices` - refactoring to use `hashicorp/go-azure-sdk` ([#21344](https://github.com/hashicorp/terraform-provider-azurerm/issues/21344))
* Data Source: `azurerm_key_vault_certificate` - add support for `resource_manager_id` and `resource_manager_versionless_id` ([#21314](https://github.com/hashicorp/terraform-provider-azurerm/issues/21314))
* Data Source: `azurerm_key_vault_secret` - support for `not_before_date` and `expiration_date` ([#21359](https://github.com/hashicorp/terraform-provider-azurerm/issues/21359))
* Data Source: `azurerm_key_vault_secret` - support specifying the keyvault secret version ([#21336](https://github.com/hashicorp/terraform-provider-azurerm/issues/21336))
* `azurerm_dashboard_grafana`- support for `UserAssigned` identitiues ([#21394](https://github.com/hashicorp/terraform-provider-azurerm/issues/21394))
* `azurerm_key_vault_certificate` - add support for `resource_manager_id` and `resource_manager_versionless_id` ([#21314](https://github.com/hashicorp/terraform-provider-azurerm/issues/21314))
* `azurerm_linux_function_app` - mark the `site_credential` block as `Sensitive` ([#21393](https://github.com/hashicorp/terraform-provider-azurerm/issues/21393))
* `azurerm_linux_function_app_slot` - mark the `site_credential` block as `Sensitive` ([#21393](https://github.com/hashicorp/terraform-provider-azurerm/issues/21393))
* `azurerm_linux_web_app` - mark the `site_credential` block as `Sensitive` ([#21393](https://github.com/hashicorp/terraform-provider-azurerm/issues/21393))
* `azurerm_linux_web_app_slot`  - mark the `site_credential` block as `Sensitive` ([#21393](https://github.com/hashicorp/terraform-provider-azurerm/issues/21393))
* `azurerm_windows_function_app` - mark the `site_credential` block as `Sensitive` ([#21393](https://github.com/hashicorp/terraform-provider-azurerm/issues/21393))
* `azurerm_windows_function_app_slot` - mark the `site_credential` block as `Sensitive` ([#21393](https://github.com/hashicorp/terraform-provider-azurerm/issues/21393))
* `azurerm_windows_web_app` - mark the `site_credential` block as `Sensitive` ([#21393](https://github.com/hashicorp/terraform-provider-azurerm/issues/21393))
* `azurerm_windows_web_app_slot` - mark the  `site_credential` block as `Sensitive` ([#21393](https://github.com/hashicorp/terraform-provider-azurerm/issues/21393))


BUG FIXES:

* `azurerm_app_configuration_key` - extend timeout for polling resource to allow propagation of read permission ([#21337](https://github.com/hashicorp/terraform-provider-azurerm/issues/21337))
* `azurerm_app_configuration_feature` - extend timeout for polling resource to allow propagation of read permission ([#21337](https://github.com/hashicorp/terraform-provider-azurerm/issues/21337))
* `azurerm_cdn_endpoint` - the `global_delivery_rule` property must have at least one action specified ([#21403](https://github.com/hashicorp/terraform-provider-azurerm/issues/21403))
* `azurerm_kubernetes_cluster` - the `enable_host_encryption` properly is not set when when resizing the `default_node_pool` ([#21379](https://github.com/hashicorp/terraform-provider-azurerm/issues/21379))
* `azurerm_linux_function_app` - fix a crash in `auth_v2` in `active_directory_v2` ([#21381](https://github.com/hashicorp/terraform-provider-azurerm/issues/21381))
* `azurerm_linux_function_app_slot` - fix a crash in `auth_v2` in `active_directory_v2` ([#21381](https://github.com/hashicorp/terraform-provider-azurerm/issues/21381))
* `azurerm_linux_web_app` - fix a crash in `auth_v2` in `active_directory_v2` ([#21381](https://github.com/hashicorp/terraform-provider-azurerm/issues/21381))
* `azurerm_linux_web_app_slot` - fix a crash in `auth_v2` in `active_directory_v2` ([#21381](https://github.com/hashicorp/terraform-provider-azurerm/issues/21381))
* `azurerm_service_plan` - support for new Premium V3 and Memory Optimised SKUs ([#21371](https://github.com/hashicorp/terraform-provider-azurerm/issues/21371))
* `azurerm_storage_account_local_user` - the `ssh_authorized_key` property can now be updated ([#21362](https://github.com/hashicorp/terraform-provider-azurerm/issues/21362))
* `azurerm_storage_mover` - remove `Microsoft.StorageMover` from required list of Resource Providers ([#21370](https://github.com/hashicorp/terraform-provider-azurerm/issues/21370))
* `azurerm_subscription` - fix an error during update ([#21255](https://github.com/hashicorp/terraform-provider-azurerm/issues/21255))
* `azurerm_windows_function_app` - fix acrash in `auth_v2` in `active_directory_v2` ([#21381](https://github.com/hashicorp/terraform-provider-azurerm/issues/21381))
* `azurerm_windows_function_app_slot` - fix a crash in `auth_v2` in `active_directory_v2` ([#21381](https://github.com/hashicorp/terraform-provider-azurerm/issues/21381))
* `azurerm_windows_web_app` - fix a crash in `auth_v2` in `active_directory_v2` ([#21381](https://github.com/hashicorp/terraform-provider-azurerm/issues/21381))
* `azurerm_windows_web_app_slot` - fix a crash in `auth_v2` in `active_directory_v2` ([#21381](https://github.com/hashicorp/terraform-provider-azurerm/issues/21381))

## 3.51.0 (April 06, 2023)

BREAKING CHANGES:

* `azurerm_kubernetes_cluster` - the `sku_tier` property no longer accepts the value `Paid`, it must be updated to `Standard` ([#21256](https://github.com/hashicorp/terraform-provider-azurerm/issues/21256))

FEATURES:

* **New Resource:** `azurerm_arc_kubernetes_cluster` ([#15401](https://github.com/hashicorp/terraform-provider-azurerm/issues/15401))
* **New Resource:** `azurerm_resource_group_cost_management_view` ([#21112](https://github.com/hashicorp/terraform-provider-azurerm/issues/21112))
* **New Resource:** `azurerm_signalr_service_custom_certificate` ([#21112](https://github.com/hashicorp/terraform-provider-azurerm/issues/21112))
* **New Resource:** `azurerm_storage_mover` ([#21000](https://github.com/hashicorp/terraform-provider-azurerm/issues/21000))
* **New Resource:** `azurerm_subscription_cost_management_view` ([#21112](https://github.com/hashicorp/terraform-provider-azurerm/issues/21112))
* **New Resource:** `azurerm_voice_services_communications_gateway_test_line` ([#21111](https://github.com/hashicorp/terraform-provider-azurerm/issues/21111))

ENHANCEMENTS:

* dependencies: updating to `v0.20230405.1143248` of `github.com/hashicorp/go-azure-sdk` ([#21312](https://github.com/hashicorp/terraform-provider-azurerm/issues/21312))
* dependencies: updating to `v0.20230331.1120327` of `github.com/tombuildsstuff/kermit` ([#21235](https://github.com/hashicorp/terraform-provider-azurerm/issues/21235))
* dependencies: updating `containerservice/2022-09-02-preview` to `2023-02-02-preview` ([#21256](https://github.com/hashicorp/terraform-provider-azurerm/issues/21256))
* dependencies: updating `search/2020-03-13` to `search/2022-09-01` ([#21250](https://github.com/hashicorp/terraform-provider-azurerm/issues/21250))
* `batch`: updating to API Version `2022-01-01.15.0` (from `github.com/tombuildsstuff/kermit`) ([#21234](https://github.com/hashicorp/terraform-provider-azurerm/issues/21234))
* Data Source: `azurerm_monitor_data_collection_rule` - support for the `data_collection_endpoint_id` property ([#21159](https://github.com/hashicorp/terraform-provider-azurerm/issues/21159))
* Data Source: `azurerm_monitor_data_collection_rule` - support for the `identity` and `stream_declaration` blocks ([#21159](https://github.com/hashicorp/terraform-provider-azurerm/issues/21159))
* Data Source: `azurerm_monitor_data_collection_rule` - support for additional `destinations`, `data_sources` and `data_flow` transformations ([#21159](https://github.com/hashicorp/terraform-provider-azurerm/issues/21159))
* `azurerm_app_configuration_feature` - support for the `key` property ([#21252](https://github.com/hashicorp/terraform-provider-azurerm/issues/21252))
* `azurerm_container_app` - the `app_port` property is now optional ([#20567](https://github.com/hashicorp/terraform-provider-azurerm/issues/20567))
* `azurerm_healthcare_fhir_service` - support for `PATCH` as an available value for `cors` ([#21222](https://github.com/hashicorp/terraform-provider-azurerm/issues/21222))
* `azurerm_healthcare_service` - upport for `PATCH` as an available value for `cors` ([#21222](https://github.com/hashicorp/terraform-provider-azurerm/issues/21222))
* `azurerm_kubernetes_cluster` - support `KataMshvVmIsolation` as a option for the `workload_runtime` property ([#21176](https://github.com/hashicorp/terraform-provider-azurerm/issues/21176))
* `azurerm_kubernetes_cluster_node_pool` - support `KataMshvVmIsolation` as a option for the `workload_runtime` property ([#21176](https://github.com/hashicorp/terraform-provider-azurerm/issues/21176))
* `azurerm_monitor_data_collection_rule` - support for the `data_collection_endpoint_id` property ([#21159](https://github.com/hashicorp/terraform-provider-azurerm/issues/21159))
* `azurerm_monitor_data_collection_rule` - support for the `identity` and `stream_declaration` blocks ([#21159](https://github.com/hashicorp/terraform-provider-azurerm/issues/21159))
* `azurerm_monitor_data_collection_rule` - support for additional `destinations`, `data_sources` and `data_flow` transformations ([#21159](https://github.com/hashicorp/terraform-provider-azurerm/issues/21159))
* `azurerm_signalr_service` - support for the `http_request_logs_enabled` property ([#21032](https://github.com/hashicorp/terraform-provider-azurerm/issues/21032))
* `azurerm_snapshot` - support for the `incremental_enabled` property ([#21263](https://github.com/hashicorp/terraform-provider-azurerm/issues/21263))
* `azurerm_web_pubsub_hub` - support for the `event_listener` block ([#21145](https://github.com/hashicorp/terraform-provider-azurerm/issues/21145))

BUG FIXES:

* Data Source: `azurerm_app_configuration_keys` - fixing a regression where the API doesn't return the http endpoint when listing items ([#21208](https://github.com/hashicorp/terraform-provider-azurerm/issues/21208))
* Data Source: `azurerm_kubernetes_cluster` - prevent errors when used with limited permissions ([#21229](https://github.com/hashicorp/terraform-provider-azurerm/issues/21229))
* `azurerm_api_management` - prevent error from empty response body when updating the resource ([#21221](https://github.com/hashicorp/terraform-provider-azurerm/issues/21221))
* `azurerm_application_gateway` - correctly validate the `firewall_policy_id` property ([#21238](https://github.com/hashicorp/terraform-provider-azurerm/issues/21238))
* `azurerm_automation_software_update_configuration` - `time_zone` correctly defaults to `Etc/UTC` ([#21254](https://github.com/hashicorp/terraform-provider-azurerm/issues/21254))
* `azurerm_digital_twins_time_series_database_connection` - insensitively parse `kusto_cluster_uri` ([#21243](https://github.com/hashicorp/terraform-provider-azurerm/issues/21243))
* `azurerm_express_route_circuit` - can now set `authorization_key` during creation ([#21132](https://github.com/hashicorp/terraform-provider-azurerm/issues/21132))
* `azurerm_kusto_eventhub_data_connection` - insensitively parse `identity_id` if it applies to a Kusto Cluster ([#21243](https://github.com/hashicorp/terraform-provider-azurerm/issues/21243))
* `azurerm_linux_function_app`  - fix a crash in `auth_v2` in `active_directory_v2` ([#21219](https://github.com/hashicorp/terraform-provider-azurerm/issues/21219))
* `azurerm_linux_function_app_slot`  - fix a crash in `auth_v2` in `active_directory_v2` ([#21219](https://github.com/hashicorp/terraform-provider-azurerm/issues/21219))
* `azurerm_linux_web_app`  - fix a crash in `auth_v2` in `active_directory_v2` ([#21219](https://github.com/hashicorp/terraform-provider-azurerm/issues/21219))
* `azurerm_linux_web_app_slot`  - fix a crash in `auth_v2` in `active_directory_v2` ([#21219](https://github.com/hashicorp/terraform-provider-azurerm/issues/21219))
* `azurerm_monitor_diagnostic_setting` - insensitively parse the resource's ID if it has been created for a Kusto Cluster ([#21243](https://github.com/hashicorp/terraform-provider-azurerm/issues/21243))
* `azurerm_mssql_database` - fix a issue with `short_term_retention_policy` preventing creation ([#21268](https://github.com/hashicorp/terraform-provider-azurerm/issues/21268))
* `azurerm_windows_function_app`  - fix a crash in `auth_v2` in `active_directory_v2` ([#21219](https://github.com/hashicorp/terraform-provider-azurerm/issues/21219))
* `azurerm_windows_function_app_slot`  - fix a crash in `auth_v2` in `active_directory_v2` ([#21219](https://github.com/hashicorp/terraform-provider-azurerm/issues/21219))
* `azurerm_windows_web_app`  - fix a crash in `auth_v2` in `active_directory_v2` ([#21219](https://github.com/hashicorp/terraform-provider-azurerm/issues/21219))
* `azurerm_windows_web_app_slot`  - fix a crash in `auth_v2` in `active_directory_v2` ([#21219](https://github.com/hashicorp/terraform-provider-azurerm/issues/21219))

## 3.50.0 (March 30, 2023)

FEATURES:

* **New DataSource:** `azurerm_container_app` ([#21199](https://github.com/hashicorp/terraform-provider-azurerm/issues/21199))
* **New Resource:** `azurerm_web_pubsub_custom_certificate` ([#21114](https://github.com/hashicorp/terraform-provider-azurerm/issues/21114))

ENHANCEMENTS:

* dependencies: updating to `v0.20230329.1052505` of `github.com/hashicorp/go-azure-sdk` ([#21175](https://github.com/hashicorp/terraform-provider-azurerm/issues/21175))
* dependencies: updated `azurerm_subscription` to use new SDK ([#18813](https://github.com/hashicorp/terraform-provider-azurerm/issues/18813))
* `azurerm_databricks_access_connector` - support for user assigned identities ([#21059](https://github.com/hashicorp/terraform-provider-azurerm/issues/21059))
* `azurerm_linux_function_app`  - add support for `zip_deploy_file` ([#20544](https://github.com/hashicorp/terraform-provider-azurerm/issues/20544))
* `azurerm_monitor_scheduled_query_rules_alert` - `trigger.x.metric_column` is now optional ([#21203](https://github.com/hashicorp/terraform-provider-azurerm/issues/21203))
* `azurerm_mssql_database` - HyperScale Skus now support `long_term_retention_policy` and `short_term_retention_policy` ([#21166](https://github.com/hashicorp/terraform-provider-azurerm/issues/21166))
* `azurerm_windows_function_app` - add support for `zip_deploy_file` ([#20544](https://github.com/hashicorp/terraform-provider-azurerm/issues/20544))

BUG FIXES:

* Data Source: `azurerm_databricks_workspace_private_endpoint_connection`: validating `private_endpoint_id` and `workspace_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* Data Source: `azurerm_healthcare_medtech_service` - the `workspace_id` field is no longer marked as ForceNew ([#21077](https://github.com/hashicorp/terraform-provider-azurerm/issues/21077))
* Data Source: `azurerm_healthcare_medtech_service` - support for Azure Environments other then Azure Public ([#21077](https://github.com/hashicorp/terraform-provider-azurerm/issues/21077))
* `azurerm_api_management` - validating `public_ip_address_id` and `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_api_management_custom_domain` - validating `api_management_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_api_management_policy` - validating `api_management_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_api_management_gateway_api` - validating `api_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_application_gateway` - validating `firewall_policy_id` and `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_application_gateway` - validating that `data` within the `ssl_certificate` block is a base64-encoded value ([#21191](https://github.com/hashicorp/terraform-provider-azurerm/issues/21191))
* `azurerm_application_insights_analytics_item` - validating `application_insights_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_application_insights_api_key` - validating `application_insights_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_application_insights_smart_detection_rule` - validating `application_insights_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_application_insights_standard_webtests` - validating `application_insights_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_application_insights_webtests` - validating `application_insights_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_app_service_virtual_network_swift_connection` - validating `app_service_id` and `virtual_network_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_bastion_host` - validating `public_ip_address_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_container_registry` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_database_migration_service` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_databricks_workspace` - validating `load_balancer_backend_address_pool_id`, `machine_learning_workspace_id` and `virtual_network_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_data_factory_linked_service_key_vault` - validating `key_vault_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_data_factory_integration_runtime_managed` - validating `vnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_data_share_dataset_kusto_cluster` - validating `kusto_cluster_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_data_share_dataset_kusto_database` - validating `kusto_database_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_eventhub_namespace` - validating the `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_eventhub_namespace_disaster_recovery_config` - fixing a bug where `partner_namespace_id` would validate with an empty string when the field should instead be omitted ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_express_route_circuit_peering` - validating `route_filter_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_express_route_gateway` - validating `virtual_hub_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_eventhub` - validating `storage_account_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_eventgrid_event_subscription` - validating `eventhub_resource_id`, `servicebus_queue_endpoint_id`, `servicebus_topic_endpoint_id` and `storage_account_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_frontdoor` - validating `web_application_firewall_policy_link_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_hdinsight_hadoop_cluster` - validating `storage_resource_id`, `subnet_id` and `virtual_network_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_hdinsight_hbase_cluster` - validating `storage_resource_id`, `subnet_id` and `virtual_network_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_hdinsight_interactive_query_cluster` - validating `storage_resource_id`, `subnet_id` and `virtual_network_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_hdinsight_kafka_cluster` - validating `storage_resource_id`, `subnet_id` and `virtual_network_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_hdinsight_spark_cluster` - validating `storage_resource_id`, `subnet_id` and `virtual_network_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_healthcare_medtech_service` - support for Azure Environments other then Azure Public ([#21077](https://github.com/hashicorp/terraform-provider-azurerm/issues/21077))
* `azurerm_hpc_cache` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_image` - validating `managed_disk_id` and `source_virtual_machine_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_iothub_certificate` - certificate content now updates correctly ([#21163](https://github.com/hashicorp/terraform-provider-azurerm/issues/21163))
* `azurerm_iothub_dps_certificate` - certificate content now updates correctly ([#21163](https://github.com/hashicorp/terraform-provider-azurerm/issues/21163))
* `azurerm_key_vault_access_policy` - validating `key_vault_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_key_vault_certificate_issuer` - validating `key_vault_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_kubernetes_cluster` - validating `vnet_subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_kubernetes_cluster_node_pool` - validating `vnet_subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_kusto_attached_database_configuration` - validating the `cluster_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_kusto_cluster` - validating `subnet_id`, `engine_public_ip_id` and `data_management_public_ip_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_kusto_eventgrid_data_connection` - validating `eventgrid_resource_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_lb` - validating `public_ip_address_id`, `public_ip_prefix_id` and `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_lb_nat_rule` - validating the `backend_address_pool_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_linux_function_app`  - fix crash in `auth_v2` in `active_directory_v2` ([#21113](https://github.com/hashicorp/terraform-provider-azurerm/issues/21113))
* `azurerm_linux_function_app_slot`  - fix crash in `auth_v2` in `active_directory_v2` ([#21113](https://github.com/hashicorp/terraform-provider-azurerm/issues/21113))
* `azurerm_linux_web_app`  - fix crash in `auth_v2` in `active_directory_v2` ([#21113](https://github.com/hashicorp/terraform-provider-azurerm/issues/21113))
* `azurerm_linux_web_app_slot`  - fix crash in `auth_v2` in `active_directory_v2` ([#21113](https://github.com/hashicorp/terraform-provider-azurerm/issues/21113))
* `azurerm_linux_virtual_machine` - validating `application_security_group_ids` and `key_vault_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_linux_virtual_machine_scale_set` - validating `key_vault_id`, `network_security_group_id`, `public_ip_prefix_id` and `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_log_analytics_linked_service` - validating the workspace id ([#21170](https://github.com/hashicorp/terraform-provider-azurerm/issues/21170))
* `azurerm_log_analytics_linked_storage_account` - validating the `storage_account_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_logic_app_action_custom` - validating `logic_app_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_logic_app_action_http` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_logic_app_trigger_custom` - validating `logic_app_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_logic_app_trigger_http_request` - validating `logic_app_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_logic_app_trigger_recurrence` - validating `logic_app_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_mssql_virtual_machine` - the `sql_license_type` property is now optional ([#21138](https://github.com/hashicorp/terraform-provider-azurerm/issues/21138))
* `azurerm_managed_disk` - validating `disk_access_id` and `storage_account_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_mariadb_virtual_network_rule` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_monitor_action_group` - validating `automation_account_id` and `function_app_resource_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_monitor_log_profile` - validating `storage_account_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_mssql_database` - fixing an int64 overflow for `max_size_gb` on 32-bit platforms ([#21155](https://github.com/hashicorp/terraform-provider-azurerm/issues/21155))
* `azurerm_mssql_database` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_mysql_virtual_network_rule` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_netapp_volume` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_network_interface` - validating `public_ip_address_id` and `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_network_interface_application_gateway_association` - validating `backend_address_pool_id` and `network_interface_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_network_interface_application_security_group_association` - validate `application_security_group_id` and `network_interface_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_network_interface_backend_address_pool_association` - validating the `backend_address_pool_id` and `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_network_interface_network_security_group_association` - validating `network_security_group_id` and `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_network_interface_nat_rule_association` - validating `network_interface_id` and `nat_rule_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_network_profile` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_network_watcher_flow_log` - fixing the delete function to work reliably during deletion ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_network_watcher_flow_log` - validating `storage_account_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_orchestrated_virtual_machine_scale_set` - validating `application_security_group_ids`, `key_vault_id`,  `proximity_placement_group_id`, `public_ip_prefix_id` and `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_private_link_service` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_public_ip` - validating `public_ip_prefix_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_postgresql_virtual_network_rule` - validating `subnet_id` is a subnet ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_private_dns_zone_virtual_network_link` - validating `virtual_network_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_role_definition` - polling for longer during deletion ([#21151](https://github.com/hashicorp/terraform-provider-azurerm/issues/21151))
* `azurerm_sentinel_automation_rule` - validating `logic_app_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_security_center_workspace` - validating `log_analytics_workspace_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_security_center_automation` - validating that a Scope is specified ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_sql_managed_database` - validating `managed_instance_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_sql_managed_instance` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_static_site_custom_domain` - validating `static_site_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_storage_account` - updating the validation for `ip_rules` to highlight the IP Range that's invalid when the validation fails ([#21178](https://github.com/hashicorp/terraform-provider-azurerm/issues/21178))
* `azurerm_storage_account_network_rules` - validating `ip_rules` ([#21178](https://github.com/hashicorp/terraform-provider-azurerm/issues/21178))
* `azurerm_storage_management_policy` - validating `storage_account_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_subnet_nat_gateway_association` - validating `nat_gateway_id` and `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_subnet_network_security_group_association` - validating `network_security_group_id` and `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_subnet_route_table_association` - validating `subnet_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_virtual_hub` - validating `virtual_wan_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_virtual_machine_data_disk_attachment` - validating `managed_disk_id` and `virtual_machine_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_virtual_network` - validating `ddos_protection_plan_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_virtual_network_gateway` - validating `default_local_network_gateway_id` and `public_ip_address_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_virtual_network_gateway_connection` - validating `express_route_circuit_id`, `local_network_gateway_id` and `peer_virtual_network_gateway_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_web_application_firewall_policy` - the `match_values` property is now optional ([#21125](https://github.com/hashicorp/terraform-provider-azurerm/issues/21125))
* `azurerm_windows_function_app`  - fix crash in `auth_v2` in `active_directory_v2` ([#21113](https://github.com/hashicorp/terraform-provider-azurerm/issues/21113))
* `azurerm_windows_function_app_slot`  - fix crash in `auth_v2` in `active_directory_v2` ([#21113](https://github.com/hashicorp/terraform-provider-azurerm/issues/21113))
* `azurerm_windows_virtual_machine_scale_set` - validating `application_security_group_ids`, `network_security_group_id`, `proximity_placement_group_id`, `public_ip_prefix_id`, `subnet_id` and `virtual_network_gateway_id` ([#21129](https://github.com/hashicorp/terraform-provider-azurerm/issues/21129))
* `azurerm_windows_web_app`  - fix crash in `auth_v2` in `active_directory_v2` ([#21113](https://github.com/hashicorp/terraform-provider-azurerm/issues/21113))
* `azurerm_windows_web_app_slot`  - fix crash in `auth_v2` in `active_directory_v2` ([#21113](https://github.com/hashicorp/terraform-provider-azurerm/issues/21113))

## 3.49.0 (March 23, 2023)

BREAKING CHANGES: 
App Service `site_config`
* `ip_restriction` blocks are no longer computed - changes to IP restrictions outside of Terraform will now present a diff
* `scm_ip_restriction` blocks are no longer computed - changes to SCM IP restrictions outside of Terraform will now present a diff
* `cors` blocks no longer require `allowed_origins`, however, if the property is supplied it must contain at least one item. Omitting this property will set the array empty

FEATURES: 

* **New Datasource:** `azurerm_orchestrated_virtual_machine_scale_set` ([#21050](https://github.com/hashicorp/terraform-provider-azurerm/issues/21050))
* **New Resource:** `azurerm_databricks_virtual_network_peering #20728` ([#20728](https://github.com/hashicorp/terraform-provider-azurerm/issues/20728))
* **New Resource:** `azurerm_sentinel_threat_intelligence_indicator` ([#20771](https://github.com/hashicorp/terraform-provider-azurerm/issues/20771))
* **New Resource:** `azurerm_voice_services_communications_gateway` ([#20607](https://github.com/hashicorp/terraform-provider-azurerm/issues/20607))

ENHANCEMENTS:

* dependencies: updating to `v0.20230322.1105901` of `hashicorp/go-azure-sdk` ([#21079](https://github.com/hashicorp/terraform-provider-azurerm/issues/21079))
* `databricks`: refactoring to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#21004](https://github.com/hashicorp/terraform-provider-azurerm/issues/21004))
* `azurerm_app_configuration_key` - the resource's ID has been changed to match the Data Plane URL format to work around a number of bugs in the previous parsing logic ([#20082](https://github.com/hashicorp/terraform-provider-azurerm/issues/20082))
* `azurerm_app_configuration_feature` - the resource's ID has been changed to match the Data Plane URL format to work around a number of bugs in the previous parsing logic ([#20082](https://github.com/hashicorp/terraform-provider-azurerm/issues/20082))
* `azurerm_express_route_circuit` - add support for `authorization_key` ([#21104](https://github.com/hashicorp/terraform-provider-azurerm/issues/21104))
* `azurerm_media_job` - updating to use API Version `2022-07-01` ([#20956](https://github.com/hashicorp/terraform-provider-azurerm/issues/20956))
* `azurerm_media_transform` - updating to use API Version `2022-07-01` ([#20956](https://github.com/hashicorp/terraform-provider-azurerm/issues/20956))
* `azurerm_virtual_network_gateway` - support for conditional/patch updates ([#21009](https://github.com/hashicorp/terraform-provider-azurerm/issues/21009))
* `azurerm_web_application_firewall_policy` - the field `operator` within the `match_conditions` block can now be set to `Any` ([#20971](https://github.com/hashicorp/terraform-provider-azurerm/issues/20971))
* `azurerm_kubernetes_cluster` - add missing property to `oms_agent` schema([#21046](https://github.com/hashicorp/terraform-provider-azurerm/issues/21046))
* `azurerm_kubernetes_cluster` - deprecate `docker_bridge_cidr` which is no longer supported by the API since docker is no longer a valid container runtime ([#20952](https://github.com/hashicorp/terraform-provider-azurerm/issues/20952))
* `azurerm_management_group_policy_assignment` - support for the `overrides` and `resource_selectors` blocks ([#20686](https://github.com/hashicorp/terraform-provider-azurerm/issues/20686))
* `azurerm_mysql_flexible_server` - support for the `geo_backup_key_vault_key_id` and `geo_backup_user_assigned_identity_id` properties ([#20796](https://github.com/hashicorp/terraform-provider-azurerm/issues/20796))
* `azurerm_resource_group_policy_assignment` - support for the `overrides` and `resource_selectors` blocks ([#20686](https://github.com/hashicorp/terraform-provider-azurerm/issues/20686))
* `azurerm_resource_policy_assignment` - support for the `overrides` and `resource_selectors` blocks ([#20686](https://github.com/hashicorp/terraform-provider-azurerm/issues/20686))
* `azurerm_role_assignment` - support subscription aliases scopes ([#20895](https://github.com/hashicorp/terraform-provider-azurerm/issues/20895))
* `azurerm_signalr_service` - support for `public_network_access_enabled`, `local_auth_enabled`, `aad_auth_enabled`, `tls_client_cert_enabled`, and `serverless_connection_timeout_in_seconds` properties ([#20975](https://github.com/hashicorp/terraform-provider-azurerm/issues/20975))
* `azurerm_subscription_policy_assignment` - support for the `overrides` and `resource_selectors` blocks ([#20686](https://github.com/hashicorp/terraform-provider-azurerm/issues/20686))
* `azurerm_sentinel_log_analytics_workspace_onboarding` - the `resource_group_name` and `workspace_name` properties have been deprecated in favour of workspace_id ([#20661](https://github.com/hashicorp/terraform-provider-azurerm/issues/20661))
* `azurerm_virtual_network_peering` - adding an explicit default value for `allow_forwarded_traffic`, `allow_gateway_transit` and `use_remote_gateways` ([#21009](https://github.com/hashicorp/terraform-provider-azurerm/issues/21009))
* `azurerm_virtual_hub` - support for the `hub_routing_preference` property ([#21028](https://github.com/hashicorp/terraform-provider-azurerm/issues/21028))

BUG FIXES:

* `azurerm_automation_account` - the `key_source` property has been deprecated ([#21041](https://github.com/hashicorp/terraform-provider-azurerm/issues/21041))
* `azurerm_application_insights` - the `workspace_id` can now be updated without creating a new resource ([#21029](https://github.com/hashicorp/terraform-provider-azurerm/issues/21029))
* `azurerm_firewall` - Prevent duplicate name from being used for `ip_configuration` and `management_ip_configuration` ([#21068](https://github.com/hashicorp/terraform-provider-azurerm/issues/21068))
* `azurerm_kubernetes_cluster` - replace calls to the deprecated accessProfiles endpoint with listUserCredentials ([#20927](https://github.com/hashicorp/terraform-provider-azurerm/issues/20927))
* `azurerm_kusto_cluster` - `language_extensions` is now a Set rather than a List ([#20951](https://github.com/hashicorp/terraform-provider-azurerm/issues/20951))
* `azurerm_linux_function_app`  - fixan update bug with the `health_check_eviction_time_in_min` property ([#21095](https://github.com/hashicorp/terraform-provider-azurerm/issues/21095))
* `azurerm_linux_function_app` - fixed processing of `cors` block ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_linux_function_app`  - fixauth_v2 `active_directory_v2` sending empty data ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_linux_function_app`  - fixread for `token_store_enabled` to correctly set returned value in state ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_linux_function_app` - the `ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_linux_function_app` - the `scm_ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_linux_function_app_slot` - fixed processing of `cors` block ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_linux_function_app_slot` - the `ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_linux_function_app_slot` - the `scm_ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_linux_function_app_slot`  - fixauth_v2 `active_directory_v2` sending empty data ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_linux_function_app_slot`  - fixread for `token_store_enabled` to correctly set returned value in state ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_linux_web_app`  - fixauth_v2 `active_directory_v2` sending empty data ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_linux_web_app`  - fixread for `token_store_enabled` to correctly set returned value in state ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_linux_web_app` - fixed processing of `cors` block ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_linux_web_app` - the `ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_linux_web_app` - the `scm_ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_linux_web_app_slot`  - fixauth_v2 `active_directory_v2` sending empty data ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_linux_web_app_slot`  - fixread for `token_store_enabled` to correctly set returned value in state ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_linux_web_app_slot` - fixed processing of `cors` block ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_linux_web_app_slot` - the `ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_linux_web_app_slot` - the `scm_ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_machine_learning_datastore_blobstorage`  - fixan issue creating this resource in clouds other than public ([#21016](https://github.com/hashicorp/terraform-provider-azurerm/issues/21016))
* `azurerm_virtual_desktop_host_pool` - changing the `load_balancer_type` property no longer creates a new resource ([#20947](https://github.com/hashicorp/terraform-provider-azurerm/issues/20947))
* `azurerm_windows_function_app`  - fixan update bug with the `health_check_eviction_time_in_min` property ([#21095](https://github.com/hashicorp/terraform-provider-azurerm/issues/21095))
* `azurerm_windows_function_app` - fixed processing of `cors` block ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_function_app`  - fixauth_v2 `active_directory_v2` sending empty data ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_windows_function_app`  - fixread for `token_store_enabled` to correctly set returned value in state ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_windows_function_app` - the `ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_function_app` - the `scm_ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_function_app_slot` - fixed processing of `cors` block ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_function_app_slot`  - fixauth_v2 `active_directory_v2` sending empty data ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_windows_function_app_slot`  - fixread for `token_store_enabled` to correctly set returned value in state ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_windows_function_app_slot` - the `ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_function_app_slot` - the `scm_ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_web_app` - the `ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_web_app` - fixed processing of `cors` block ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_web_app` - the `scm_ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_web_app`  - fixauth_v2 `active_directory_v2` sending empty data ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_windows_web_app`  - fixread for `token_store_enabled` to correctly set returned value in state ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_windows_web_app_slot` - the `ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_web_app_slot` - fixed processing of `cors` block ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_web_app_slot` - the `scm_ip_restriction` block can is now successfully removed by removing from config ([#20987](https://github.com/hashicorp/terraform-provider-azurerm/issues/20987))
* `azurerm_windows_web_app_slot`  - fixauth_v2 `active_directory_v2` sending empty data ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))
* `azurerm_windows_web_app_slot`  - fixread for `token_store_enabled` to correctly set returned value in state ([#21091](https://github.com/hashicorp/terraform-provider-azurerm/issues/21091))


## 3.48.0 (March 16, 2023)

FEATURES: 

* **New Data Source:** `azurerm_mobile_network_sim_policy` [FGH-20732]
* **New Resource:** `azurerm_express_route_port_authorization` ([#20736](https://github.com/hashicorp/terraform-provider-azurerm/issues/20736))
* **New Resource:** `azurerm_mobile_network_sim_policy` ([#20732](https://github.com/hashicorp/terraform-provider-azurerm/issues/20732))
* **New Resource:** `azurerm_site_recovery_vmware_replication_policy` ([#20881](https://github.com/hashicorp/terraform-provider-azurerm/issues/20881))
* **New Resource:** `azurerm_sentinel_alert_rule_anomaly_duplicate` ([#20760](https://github.com/hashicorp/terraform-provider-azurerm/issues/20760))

ENHANCEMENTS:

* dependencies: updating to `v0.20230316.1132628` of `github.com/hashicorp/go-azure-sdk` ([#20986](https://github.com/hashicorp/terraform-provider-azurerm/issues/20986))
* `signalr`: updating to API Version `2023-02-01` ([#20910](https://github.com/hashicorp/terraform-provider-azurerm/issues/20910))
* `webpubsub`: updating to API Version `2023-02-01` ([#20910](https://github.com/hashicorp/terraform-provider-azurerm/issues/20910))
* `azurerm_express_route_gateway` - support for the `allow_non_virtual_wan_traffic` property ([#20667](https://github.com/hashicorp/terraform-provider-azurerm/issues/20667))
* `azurerm_ssh_public_key` -  allow `.` for `name` validation ([#20955](https://github.com/hashicorp/terraform-provider-azurerm/issues/20955))

BUG FIXES:

* provider: fix an authentication bug which sometimes caused access tokens to be refreshed too late ([#20894](https://github.com/hashicorp/terraform-provider-azurerm/issues/20894))
* `azurerm_bot_channel_directline` - fixing an issue where an empty `site` was passed to the API ([#20890](https://github.com/hashicorp/terraform-provider-azurerm/issues/20890))
* `azurerm_healthcare_dicom_service` - extending the `create` and `update` timeouts to `90` minutes ([#20932](https://github.com/hashicorp/terraform-provider-azurerm/issues/20932))
* `azurerm_kusto_eventhub_data_connection` - fixing an issue where an existing resource wouldn't be flagged during creation ([#20926](https://github.com/hashicorp/terraform-provider-azurerm/issues/20926))
* `azurerm_linux_function_app` - Fixed apply time validation when using `WEBSITE_CONTENTOVERVNET`  ([#18258](https://github.com/hashicorp/terraform-provider-azurerm/issues/18258))
* `azurerm_windows_function_app` - Fixed apply time validation when using `WEBSITE_CONTENTOVERVNET` ([#18258](https://github.com/hashicorp/terraform-provider-azurerm/issues/18258))


## 3.47.0 (March 09, 2023)

FEATURES: 

* **New Resource:** `azurerm_sentinel_metadata` ([#20801](https://github.com/hashicorp/terraform-provider-azurerm/issues/20801))

ENHANCEMENTS

* dependencies: updating to `v4.4.0+incompatible` of `github.com/gofrs/uuid` ([#20821](https://github.com/hashicorp/terraform-provider-azurerm/issues/20821))
* dependencies: updating to `v0.55.0` of `github.com/hashicorp/go-azure-helpers` ([#20807](https://github.com/hashicorp/terraform-provider-azurerm/issues/20807))
* dependencies: updating to version `v0.20230309.1123256` of `github.com/hashicorp/go-azure-sdk` ([#20810](https://github.com/hashicorp/terraform-provider-azurerm/issues/20810))
* dependencies: updating to `v0.20230307.1105329` of `github.com/tombuildsstuff/kermit` ([#20821](https://github.com/hashicorp/terraform-provider-azurerm/issues/20821))
* dependencies: updating `redis/2021-06-01` to `redis/2022-06-01` ([#20839](https://github.com/hashicorp/terraform-provider-azurerm/issues/20839))
* `dashboard`: refactoring to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#20810](https://github.com/hashicorp/terraform-provider-azurerm/issues/20810))
* `media`: refactoring to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#20810](https://github.com/hashicorp/terraform-provider-azurerm/issues/20810))
* `servicebus`: refactoring to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#20810](https://github.com/hashicorp/terraform-provider-azurerm/issues/20810))
* Data Source: `azurerm_function_app_host_keys` - exporting `blobs_extension_key` ([#20837](https://github.com/hashicorp/terraform-provider-azurerm/issues/20837))
* Data Source: `azurerm_servicebus_namespace` - exporting `endpoint` ([#20790](https://github.com/hashicorp/terraform-provider-azurerm/issues/20790))
* Data Source: `azurerm_kubernetes_cluster` - generate and export `node_resource_group_id` ([#20830](https://github.com/hashicorp/terraform-provider-azurerm/issues/20830))
* `azurerm_kubernetes_cluster` - generate and export `node_resource_group_id` ([#20830](https://github.com/hashicorp/terraform-provider-azurerm/issues/20830))
* `azurerm_kubernetes_cluster` - support for the`vertical_pod_autoscaler_enabled` property ([#20751](https://github.com/hashicorp/terraform-provider-azurerm/issues/20751))
* `azurerm_kubernetes_cluster` - support for the `msi_auth_for_monitoring_enabled` property ([#20757](https://github.com/hashicorp/terraform-provider-azurerm/issues/20757))
* `azurerm_kubernetes_cluster` - the `vm_size` property of the `default_node_pool` is no longer ForceNew and can be resized by specifying `temporary_name_for_rotation` ([#20628](https://github.com/hashicorp/terraform-provider-azurerm/issues/20628))
* `azurerm_mariadb_server` - support for the `ssl_minimal_tls_version_enforced` property ([#20782](https://github.com/hashicorp/terraform-provider-azurerm/issues/20782))
* `azurerm_monitor_action_group` - support for the `location` property ([#20603](https://github.com/hashicorp/terraform-provider-azurerm/issues/20603))
* `azurerm_mssql_database` - support for `ServerlessGen5` Hyperscale ([#20875](https://github.com/hashicorp/terraform-provider-azurerm/issues/20875))
* `azurerm_mssql_managed_database` - support for retention policies ([#20845](https://github.com/hashicorp/terraform-provider-azurerm/issues/20845))
* `azurerm_servicebus_namespace` - exports the `endpoint` attribute ([#20790](https://github.com/hashicorp/terraform-provider-azurerm/issues/20790))
* `azurerm_virtual_network_peering` - support for  the `triggers` property to allow `address_space` synchronization ([#20877](https://github.com/hashicorp/terraform-provider-azurerm/issues/20877))

BUG FIXES:

* provider: fix an issue with authentication using `oidc_token_file_path` ([#20824](https://github.com/hashicorp/terraform-provider-azurerm/issues/20824))
* provider: fix an issue with Azure CLI authentication when running in Azure Cloud Shell ([#20824](https://github.com/hashicorp/terraform-provider-azurerm/issues/20824))
* `azurerm_application_insights_analytics_item` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_automated_connection_type` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_automation_software_update_configuration` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_automation_source_control` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_automation_watcher` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_cdn_frontdoor_origin`  - fixregression where `origin_host_header` value would be inadvertently removed ([#20874](https://github.com/hashicorp/terraform-provider-azurerm/issues/20874))
* `azurerm_cdn_frontdoor_route_disable_link_to_default_domain` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_container_registry_task`  - fixupdating failed due to incomplete `registry_credential` ([#20841](https://github.com/hashicorp/terraform-provider-azurerm/issues/20841))
* `azurerm_digital_twins_time_series_database_connection` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_fluid_relay_server` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_function_app_active_slot` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_iothub_endpoint_eventhub` - marking the resource as gone when it's been deleted outside of Terraform ([#20798](https://github.com/hashicorp/terraform-provider-azurerm/issues/20798))
* `azurerm_iothub`  - fixwrong default value of `file_upload.sas_ttl` when not specified ([#20854](https://github.com/hashicorp/terraform-provider-azurerm/issues/20854))
* `azurerm_iothub_endpoint_servicebus_queue` - marking the resource as gone when it's been deleted outside of Terraform ([#20798](https://github.com/hashicorp/terraform-provider-azurerm/issues/20798))
* `azurerm_iothub_endpoint_servicebus_topic` - marking the resource as gone when it's been deleted outside of Terraform ([#20798](https://github.com/hashicorp/terraform-provider-azurerm/issues/20798))
* `azurerm_iothub_endpoint_servicebus_queue` - marking the resource as gone when it's been deleted outside of Terraform ([#20798](https://github.com/hashicorp/terraform-provider-azurerm/issues/20798))
* `azurerm_iothub_endpoint_storage_container` - marking the resource as gone when it's been deleted outside of Terraform ([#20798](https://github.com/hashicorp/terraform-provider-azurerm/issues/20798))
* `azurerm_iothub_fallback_route` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_iothub_route` - marking the resource as gone when it's been deleted outside of Terraform ([#20798](https://github.com/hashicorp/terraform-provider-azurerm/issues/20798))
* `azurerm_kubernetes_cluster`  - fixvalidation logic for `dns_prefix` ([#20813](https://github.com/hashicorp/terraform-provider-azurerm/issues/20813))
* `azurerm_linux_function_app_slot`  - fixhealth_check_eviction_time_in_min ([#20816](https://github.com/hashicorp/terraform-provider-azurerm/issues/20816))
* `azurerm_logic_app_integration_account` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_maintenance_assignment_virtual_machine` - prevent a potential panic from a nil value ([#20781](https://github.com/hashicorp/terraform-provider-azurerm/issues/20781))
* `azurerm_maintenance_assignment_virtual_machine` - maintenance configuration is now obtained by name rather than using the first in the list ([#20766](https://github.com/hashicorp/terraform-provider-azurerm/issues/20766))
* `azurerm_nginx_certificate` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_nginx_configuration` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_nginx_deployment` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_synapse_workspace_aad_admin` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_synapse_workspace_key` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_synapse_workspace_sql_aad_admin` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_web_app_active_slot` - marking the resource as gone when it's been deleted outside of Terraform ([#20797](https://github.com/hashicorp/terraform-provider-azurerm/issues/20797))
* `azurerm_windows_function_app_slot`  - fixhealth_check_eviction_time_in_min ([#20816](https://github.com/hashicorp/terraform-provider-azurerm/issues/20816))


## 3.46.0 (March 02, 2023)

FEATURES

* **New Data Source:** `azurerm_mobile_network_data_network` ([#20338](https://github.com/hashicorp/terraform-provider-azurerm/issues/20338))
* **New Data Source:** `azurerm_sentinel_alert_rule_anomaly_built_in` ([#20368](https://github.com/hashicorp/terraform-provider-azurerm/issues/20368))
* **New Resource:** `azurerm_mobile_network_data_network` ([#20338](https://github.com/hashicorp/terraform-provider-azurerm/issues/20338))
* **New Resource:** `azurerm_sentinel_alert_rule_anomaly_built_in` ([#20368](https://github.com/hashicorp/terraform-provider-azurerm/issues/20368))
* **New Resource:** `azurerm_sentinel_alert_rule_threat_intelligence` ([#20739](https://github.com/hashicorp/terraform-provider-azurerm/issues/20739))

ENHANCEMENTS

* dependencies: updating to `v0.20230228.1160358` of `github.com/hashicorp/go-azure-sdk` ([#20688](https://github.com/hashicorp/terraform-provider-azurerm/issues/20688))
* dependencies: updating to `v0.20230224.1120200` of `github.com/tombuildsstuff/kermit` ([#20649](https://github.com/hashicorp/terraform-provider-azurerm/issues/20649))
* dependencies: updating `containerservice/2022-09-02-preview` to `2023-01-02-preview` ([#20734](https://github.com/hashicorp/terraform-provider-azurerm/issues/20734))
* dependencies: updating `hybridCompute/2022-03-10` to `2022-11-10` ([#20733](https://github.com/hashicorp/terraform-provider-azurerm/issues/20733))
* `aadb2c`: refactoring to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#20715](https://github.com/hashicorp/terraform-provider-azurerm/issues/20715))
* `databoxedge` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20638](https://github.com/hashicorp/terraform-provider-azurerm/issues/20638))
* `dns`: refactoring to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#20688](https://github.com/hashicorp/terraform-provider-azurerm/issues/20688))
* `maps`: refactoring to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#20688](https://github.com/hashicorp/terraform-provider-azurerm/issues/20688))
* `signalr`: refactoring to use `hashicorp/go-azure-sdk` as a base layer rather than `Azure/go-autorest` ([#20688](https://github.com/hashicorp/terraform-provider-azurerm/issues/20688))
* `compute/shared_image_gallery` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20599](https://github.com/hashicorp/terraform-provider-azurerm/issues/20599))
* `compute/gallery_application` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20599](https://github.com/hashicorp/terraform-provider-azurerm/issues/20599))
* `compute/gallery_application_version` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20599](https://github.com/hashicorp/terraform-provider-azurerm/issues/20599))
* `iottimeseriesinsights` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20416](https://github.com/hashicorp/terraform-provider-azurerm/issues/20416))
* `policy/assignment` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20638](https://github.com/hashicorp/terraform-provider-azurerm/issues/20638))
* `sentinel/alert_rule` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20680](https://github.com/hashicorp/terraform-provider-azurerm/issues/20680))
* `sentinel/automation_rule` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20726](https://github.com/hashicorp/terraform-provider-azurerm/issues/20726))
* **Data Source:** `azurerm_linux_function_app` - support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20722](https://github.com/hashicorp/terraform-provider-azurerm/issues/20722))
* **Data Source:**`azurerm_windows_function_app` -support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20722](https://github.com/hashicorp/terraform-provider-azurerm/issues/20722))
* `azurerm_app_service_connection` - support for the `secret_store` block ([#20613](https://github.com/hashicorp/terraform-provider-azurerm/issues/20613))
* `express_route_circuit_peering_resource` - support for the `advertised_communities` property ([#20708](https://github.com/hashicorp/terraform-provider-azurerm/issues/20708))
* `azurerm_healthcare_service` - extend range of the cosmosdb_throughput to a maximum of `100000` ([#20755](https://github.com/hashicorp/terraform-provider-azurerm/issues/20755))
* `azurerm_key_vault_key` - support for the `rotation_policy` block ([#19113](https://github.com/hashicorp/terraform-provider-azurerm/issues/19113))
* `azurerm_kubernetes_cluster` - support for `Standard` with the `sku_tier` ([#20734](https://github.com/hashicorp/terraform-provider-azurerm/issues/20734))
* `azurerm_linux_function_app` - support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20722](https://github.com/hashicorp/terraform-provider-azurerm/issues/20722))
* `azurerm_linux_function_app_slot` - support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20722](https://github.com/hashicorp/terraform-provider-azurerm/issues/20722))
* `azurerm_media_streaming_policy` - support for the `common_encryption_cbcs.clear_key_encryption`, `common_encryption_cenc.clear_key_encryption`, `common_encryption_cenc.clear_track`, `common_encryption_cenc.content_key_to_track_mapping` and `envelope_encryption` properties ([#20524](https://github.com/hashicorp/terraform-provider-azurerm/issues/20524))
* `azurerm_orchestrated_virtual_machine_scale_set` - support for the  `priority_mix` property ([#20618](https://github.com/hashicorp/terraform-provider-azurerm/issues/20618))
* `azurerm_storage_management_policy` - support for `auto_tier_to_hot_from_cool_enabled` ([#20641](https://github.com/hashicorp/terraform-provider-azurerm/issues/20641))
* `azurerm_spring_cloud_connection` - support for the `secret_store` block ([#20613](https://github.com/hashicorp/terraform-provider-azurerm/issues/20613))
* `azurerm_windows_function_app` - support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20722](https://github.com/hashicorp/terraform-provider-azurerm/issues/20722))
* `azurerm_windows_function_app_slot` - support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20722](https://github.com/hashicorp/terraform-provider-azurerm/issues/20722))



BUG FIXES

* Data Source: `azurerm_automation_variable_bool` - fixed a regression in read ([#20665](https://github.com/hashicorp/terraform-provider-azurerm/issues/20665))
* Data Source: `azurerm_automation_variable_datetime` - fixed a regression in read ([#20665](https://github.com/hashicorp/terraform-provider-azurerm/issues/20665))
* Data Source: `azurerm_automation_variable_int` - fixed a regression in read ([#20665](https://github.com/hashicorp/terraform-provider-azurerm/issues/20665))
* Data Source: `azurerm_automation_variable_string` - fixed a regression in read ([#20665](https://github.com/hashicorp/terraform-provider-azurerm/issues/20665))
* `azurerm_aadb2c_directory` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_cdn_frontdoor_origin` - `origin_host_header` can now be cleared once it has been set ([#20679](https://github.com/hashicorp/terraform-provider-azurerm/issues/20679))
* `azurerm_container_app` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_communication_service` - changing the `data_location` property now creates a new resource ([#20711](https://github.com/hashicorp/terraform-provider-azurerm/issues/20711))
* `azurerm_eventhub_cluster` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_eventhub_namespace` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_eventhub_namespace_disaster_recovery_config` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_kubernetes_cluster_node_pool` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_iothub_dps` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_media_services_account`: fix crash around `key_delivery_access_control` ([#20685](https://github.com/hashicorp/terraform-provider-azurerm/issues/20685))
* `azurerm_netapp_account` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_netapp_pool` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_netapp_snapshot` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_netapp_snapshot_policy` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_netapp_volume` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_netapp_volume`  - fixpotential nil panic in resource read ([#20662](https://github.com/hashicorp/terraform-provider-azurerm/issues/20662))
* `azurerm_notification_hub` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_notification_hub_namespace` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_proximity_placement_group` - will now correctly update when a vm is attached ([#20131](https://github.com/hashicorp/terraform-provider-azurerm/issues/20131))
* `azurerm_sentinel_log_analytics_workspace_onboard` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_servicebus_namespace_disaster_recovery_config` - fixing a crash when the connection dropped ([#20670](https://github.com/hashicorp/terraform-provider-azurerm/issues/20670))
* `azurerm_storage_object_replication` - now functions when cross tenant replication is disabled ([#20132](https://github.com/hashicorp/terraform-provider-azurerm/issues/20132))

## 3.45.0 (February 23, 2023)

FEATURES

* `App Service` - Add authV2 to Web Apps ([#20449](https://github.com/hashicorp/terraform-provider-azurerm/issues/20449))
* **New Resource:** `azurerm_site_recovery_hyperv_replication_policy` ([#20454](https://github.com/hashicorp/terraform-provider-azurerm/issues/20454))
* **New Resource:** `azurerm_site_recovery_hyperv_replication_policy_association` ([#20630](https://github.com/hashicorp/terraform-provider-azurerm/issues/20630))

ENHANCEMENTS

* dependencies: updating to `v0.20230222.1094703` of `github.com/hashicorp/go-azure-sdk` ([#20610](https://github.com/hashicorp/terraform-provider-azurerm/issues/20610))
* dependencies: updating to `v0.7.0` of `golang.org/x/net` ([#20541](https://github.com/hashicorp/terraform-provider-azurerm/issues/20541))
* `automation` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20568](https://github.com/hashicorp/terraform-provider-azurerm/issues/20568))
* `compute/capacityreservations` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20580](https://github.com/hashicorp/terraform-provider-azurerm/issues/20580))
* `compute/capacityreservationgroups` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20580](https://github.com/hashicorp/terraform-provider-azurerm/issues/20580))
* `kusto` - switching to use `github.com/hashicorp/go-azure-sdk` ([#20563](https://github.com/hashicorp/terraform-provider-azurerm/issues/20563))
* `azurerm_backup_policy_vm` - add support for `instant_restore_resource_group` ([#20562](https://github.com/hashicorp/terraform-provider-azurerm/issues/20562))
* `azurerm_express_route_connection` - support for the `inbound_route_map_id`, `outbound_route_map_id`, and `enabled_private_link_fast_path` properties ([#20619](https://github.com/hashicorp/terraform-provider-azurerm/issues/20619))
* `azurerm_kusto_cluster_customer_managed_key` - `key_version` is now Optional to allow for auto-rotation of key ([#20583](https://github.com/hashicorp/terraform-provider-azurerm/issues/20583))
* `azurerm_linux_virtual_machine` - strengthen validation for `admin_password` ([#20558](https://github.com/hashicorp/terraform-provider-azurerm/issues/20558))
* `azurerm_linux_web_app` - add support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20449](https://github.com/hashicorp/terraform-provider-azurerm/issues/20449))
* `azurerm_linux_web_app_slot` - add support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20449](https://github.com/hashicorp/terraform-provider-azurerm/issues/20449))
* `azurerm_postgresql_flexible_server` - a server can now be created without enabling password authtication ([#20578](https://github.com/hashicorp/terraform-provider-azurerm/issues/20578))
* `azurerm_media_streaming_endpoint` - add support for reading `sku` and increase limit for `scale_units` ([#20585](https://github.com/hashicorp/terraform-provider-azurerm/issues/20585))
* `azurerm_recovery_services_vault` - add support for `classic_vmware_replication_enabled` ([#20473](https://github.com/hashicorp/terraform-provider-azurerm/issues/20473))
* `azurerm_windows_virtual_machine` - strengthen validation for `admin_password` ([#20558](https://github.com/hashicorp/terraform-provider-azurerm/issues/20558))
* `azurerm_windows_web_app` - add support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20449](https://github.com/hashicorp/terraform-provider-azurerm/issues/20449))
* `azurerm_windows_web_app_slot` - add support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20449](https://github.com/hashicorp/terraform-provider-azurerm/issues/20449))
* **Data Source:** `azurerm_linux_web_app` - add support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20449](https://github.com/hashicorp/terraform-provider-azurerm/issues/20449))
* **Data Source:**`azurerm_windows_web_app` - add support for AuthV2 (EasyAuthV2) `auth_settings_v2` ([#20449](https://github.com/hashicorp/terraform-provider-azurerm/issues/20449))

BUG FIXES

* Data Source: `azurerm_linux_web_app` - set `virtual_network_subnet_id` correctly ([#20577](https://github.com/hashicorp/terraform-provider-azurerm/issues/20577))
* Data Source: `azurerm_redis_cache`  - fixissue when no patch schedules can be found ([#20516](https://github.com/hashicorp/terraform-provider-azurerm/issues/20516))
* Data Source: `azurerm_windows_web_app` - set `virtual_network_subnet_id` correctly ([#20577](https://github.com/hashicorp/terraform-provider-azurerm/issues/20577))
* `azurerm_batch_pool` - set user assigned id for `azure_blob_file_system` correctly ([#20560](https://github.com/hashicorp/terraform-provider-azurerm/issues/20560))
* `azurerm_iot_dps` - allow older resources to update without having set `data_residency_enabled` ([#20632](https://github.com/hashicorp/terraform-provider-azurerm/issues/20632))
* `azurerm_kubernetes_cluster` - prevent crash when `SecurityProfile` is nil ([#20584](https://github.com/hashicorp/terraform-provider-azurerm/issues/20584))
* `azurerm_log_analytics_workspace` - prevent ForceNew when `sku` is `LACluster` ([#19608](https://github.com/hashicorp/terraform-provider-azurerm/issues/19608))
* `azurerm_media_streaming_endpoint` - set and update `tags` properly ([#20585](https://github.com/hashicorp/terraform-provider-azurerm/issues/20585))
* `azurerm_mobile_network_sim_group` - update `identity` to only support User Assigned Identities ([#20474](https://github.com/hashicorp/terraform-provider-azurerm/issues/20474))
* `azurerm_monitor_diagnostic_setting` - the `log_analytics_destination_type` property is nto computer rather then defaulting to `AzureDiagnostics` on new resources ([#20203](https://github.com/hashicorp/terraform-provider-azurerm/issues/20203))

## 3.44.1 (February 17, 2023)

ENHANCEMENTS

* dependencies: updating to `v0.20230217.1150808` of `github.com/hashicorp/go-azure-sdk` ([#20539](https://github.com/hashicorp/terraform-provider-azurerm/issues/20539))

BUG FIXES

* authentication: fixing an issue when obtaining the auth token for Resource Manager in Azure Government ([#20523](https://github.com/hashicorp/terraform-provider-azurerm/issues/20523))
* authentication: fixing an issue where the default subscription ID was not detected when authenticating using Azure CLI ([#20526](https://github.com/hashicorp/terraform-provider-azurerm/issues/20526))
* authentication: fixing an issue where Managed Identity authentication would fail ([#20523](https://github.com/hashicorp/terraform-provider-azurerm/issues/20523))
* Data Source: `azurerm_app_configuration_key` - fixing an issue where the App Configuration was misleadingly marked as gone when the data plane client couldn't be build ([#20533](https://github.com/hashicorp/terraform-provider-azurerm/issues/20533))
* Data Source: `azurerm_app_configuration_key` - surfacing the error when a data plane client can't be built ([#20533](https://github.com/hashicorp/terraform-provider-azurerm/issues/20533))
* Data Source: `azurerm_app_configuration_keys` - fixing an issue where the App Configuration was misleadingly marked as gone when the data plane client couldn't be build ([#20533](https://github.com/hashicorp/terraform-provider-azurerm/issues/20533))
* Data Source: `azurerm_app_configuration_keys` - surfacing the error when a data plane client can't be built ([#20533](https://github.com/hashicorp/terraform-provider-azurerm/issues/20533))
* `azurerm_app_configuration_feature` - fixing an issue where the App Configuration was misleadingly marked as gone when the data plane client couldn't be build ([#20533](https://github.com/hashicorp/terraform-provider-azurerm/issues/20533))
* `azurerm_app_configuration_feature` - surfacing the error when a data plane client can't be built ([#20533](https://github.com/hashicorp/terraform-provider-azurerm/issues/20533))
* `azurerm_app_configuration_key` - fixing an issue where the App Configuration was misleadingly marked as gone when the data plane client couldn't be build ([#20533](https://github.com/hashicorp/terraform-provider-azurerm/issues/20533))
* `azurerm_app_configuration_key` - surfacing the error when a data plane client can't be built ([#20533](https://github.com/hashicorp/terraform-provider-azurerm/issues/20533))
* `azurerm_kubernetes_cluster`  - fixa crash when `securityProfile` is nil in the API Response ([#20517](https://github.com/hashicorp/terraform-provider-azurerm/issues/20517))
* `azurerm_logic_app_standard` - fixing an issue where the `storage endpoint suffix` couldn't be found ([#20536](https://github.com/hashicorp/terraform-provider-azurerm/issues/20536))
* `azurerm_synapse_role_assignment` - fixing an issue where the `Synapse domain suffix` couldn't be found ([#20536](https://github.com/hashicorp/terraform-provider-azurerm/issues/20536))

## 3.44.0 (February 16, 2023)

FEATURES:

* **New Data Source:** `azurerm_hybrid_compute_machine` ([#20211](https://github.com/hashicorp/terraform-provider-azurerm/issues/20211))
* **New Data Source:** `azurerm_policy_definition_built_in` ([#19933](https://github.com/hashicorp/terraform-provider-azurerm/issues/19933))
* **New Data Source:** `azurerm_mobile_network_service` ([#20337](https://github.com/hashicorp/terraform-provider-azurerm/issues/20337))
* **New Data Source:** `azurerm_mobile_network_site` ([#20334](https://github.com/hashicorp/terraform-provider-azurerm/issues/20334))
* **New Data Source:** `azurerm_mobile_network_slice` ([#20336](https://github.com/hashicorp/terraform-provider-azurerm/issues/20336))
* **New Data Source:** `azurerm_mobile_network_sim_group` ([#20339](https://github.com/hashicorp/terraform-provider-azurerm/issues/20339))
* **New Data Source:** `azurerm_virtual_desktop_host_pool` ([#20505](https://github.com/hashicorp/terraform-provider-azurerm/issues/20505))
* **New Resource:** `azurerm_network_manager_security_admin_configuration` ([#20233](https://github.com/hashicorp/terraform-provider-azurerm/issues/20233))
* **New Resource:** `azurerm_network_manager_admin_rule_collection` ([#20233](https://github.com/hashicorp/terraform-provider-azurerm/issues/20233))
* **New Resource:** `azurerm_network_manager_admin_rule` ([#20233](https://github.com/hashicorp/terraform-provider-azurerm/issues/20233))
* **New Resource:** `azurerm_mobile_network_service` ([#20337](https://github.com/hashicorp/terraform-provider-azurerm/issues/20337))
* **New Resource:** `azurerm_mobile_network_site` ([#20334](https://github.com/hashicorp/terraform-provider-azurerm/issues/20334))
* **New Resource:** `azurerm_mobile_network_slice` ([#20336](https://github.com/hashicorp/terraform-provider-azurerm/issues/20336))
* **New Resource:** `azurerm_mobile_network_sim_group` [GH-20339
* **New Resource:** `azurerm_site_recovery_services_vault_hyperv_site` [GH-204309

ENHANCEMENTS:

* dependencies: updating to `v0.20230216.1112535` of `github.com/hashicorp/go-azure-sdk` ([#20465](https://github.com/hashicorp/terraform-provider-azurerm/issues/20465))
* dependencies: no longer utilizing `github.com/manicminer/hamilton` ([#20320](https://github.com/hashicorp/terraform-provider-azurerm/issues/20320))
* provider: support for the `client_certificate` provider property ([#20320](https://github.com/hashicorp/terraform-provider-azurerm/issues/20320))
* provider: support for the `use_cli` provider property ([#20320](https://github.com/hashicorp/terraform-provider-azurerm/issues/20320))
* provider: authentication now uses the `github.com/hashicorp/go-azure-sdk/sdk/auth` package ([#20320](https://github.com/hashicorp/terraform-provider-azurerm/issues/20320))
* provider: cloud configuration now uses the `github.com/hashicorp/go-azure-sdk/sdk/environments` package ([#20320](https://github.com/hashicorp/terraform-provider-azurerm/issues/20320))
* `datashare`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20501](https://github.com/hashicorp/terraform-provider-azurerm/issues/20501))
* `managementlocks`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20387](https://github.com/hashicorp/terraform-provider-azurerm/issues/20387))
* `media`: refactoring `StreamingEndpoints` to use API Version `2022-08-01` ([#20457](https://github.com/hashicorp/terraform-provider-azurerm/issues/20457))
* `postgres` - updating API to `2022-12-01` ([#20370](https://github.com/hashicorp/terraform-provider-azurerm/issues/20370))
* Data Source: `azurerm_policy_definition` - support for the `mode` property ([#20420](https://github.com/hashicorp/terraform-provider-azurerm/issues/20420))
* Data Source: `azurerm_key_vault_certificates` - now exports the `certificates` block ([#20498](https://github.com/hashicorp/terraform-provider-azurerm/issues/20498))
* Data Source: `azurerm_key_vault_secrets` - now exports the `secrets` block ([#20498](https://github.com/hashicorp/terraform-provider-azurerm/issues/20498))
* `azurerm_api_management` - support for the `delegation` block ([#20399](https://github.com/hashicorp/terraform-provider-azurerm/issues/20399))
* `azurerm_container_app` - now supports multiple `container` blocks ([#20423](https://github.com/hashicorp/terraform-provider-azurerm/issues/20423))
* `azurerm_cognitive_account` - the field `sku_name` can now be set to `DC0` ([#20426](https://github.com/hashicorp/terraform-provider-azurerm/issues/20426))
* `azurerm_container_app` - support for the `registry.identity` property ([#20466](https://github.com/hashicorp/terraform-provider-azurerm/issues/20466))
* `azurerm_data_factory_linked_service_azure_blob_storage` - Add support for `connection_string_insecure`  [Gh-20494]
* `azurerm_express_route_port` - support for the `billing_type` property ([#20361](https://github.com/hashicorp/terraform-provider-azurerm/issues/20361))
* `azurerm_kubernetes_cluster` - the `web_app_routing.dns_zone_id` property now accepts an empty string for BYO DNS ([#20341](https://github.com/hashicorp/terraform-provider-azurerm/issues/20341))
* `azurerm_linux_virtual_machine` - validating that the value for the `admin_username` property isn't a disallowed username ([#20424](https://github.com/hashicorp/terraform-provider-azurerm/issues/20424))
* `azurerm_windows_virtual_machine` - validating that the value for the `admin_username` property isn't a disallowed username ([#20424](https://github.com/hashicorp/terraform-provider-azurerm/issues/20424))

BUG FIXES:

* Data Source: `azurerm_aadb2c_directory` - fixing a bug where the Data Source didn't return an error when the AAD B2C was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_app_service_environment_v3` - fixing a bug where the Data Source didn't return an error when the App Service Environment was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_consumption_budget_resource_group` - using the correct timeout value ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_consumption_budget_resource_group` - fixing a bug where the Data Source didn't return an error when the Consumption Budget Resource Group was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_data_protection_backup_vault` - fixing a bug where the Data Source didn't return an error when the Data Protection Backup Vault was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_databox_edge_device` - fixing a bug where the Data Source didn't return an error when the DataBox Edge Device was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_healthcare_dicom` - fixing a bug where the Data Source didn't return an error when the HealthCare DICOM was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_healthcare_fhir` - fixing a bug where the Data Source didn't return an error when the HealthCare FHIR was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_healthcare_medtech_service` - fixing a bug where the Data Source didn't return an error when the HealthCare MedTech Service was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_key_vault_certificate_data` - fixing a bug where the Data Source didn't return an error when the KeyVault Certificate was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_key_vault_certificate` - fixing a bug where the Data Source didn't return an error when the KeyVault Certificate was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_lb_outbound_rule` - fixing a bug where the Data Source didn't return an error when the Load Balancer Outbound Rule was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_lb_rule` - fixing a bug where the Data Source didn't return an error when the Load Balancer Rule was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_local_network_gateway` - fixing a bug where the Data Source didn't return an error when the Local Network Gateway was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_mobile_network` - fixing a bug where the Data Source didn't return an error when the Mobile Network was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_monitor_data_collection_endpoint` - fixing a bug where the Data Source didn't return an error when the Monitor Data Collection Endpoint was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_mssql_managed_instance` - fixing a bug where the Data Source didn't return an error when the MSSQL Managed Instance was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_policy_assignment` - fixing a bug where the Data Source didn't return an error when the Policy Assignment was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_redis_enterprise_database` - fixing a bug where the Data Source didn't return an error when the Redis Enterprise Database was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_servicebus_namespace_disaster_recovery_config` - fixing a bug where the Data Source didn't return an error when the ServiceBus Namespace Disaster Recovery Config was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_site_recovery_replication_recovery_plan` - fixing a bug where the Data Source didn't return an error when the Site Recovery Replication Recovery Plan was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_storage_blob` - fixing a bug where the Data Source didn't return an error when the Blob was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_storage_table_entity` - fixing a bug where the Data Source didn't return an error when the Table Entity was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_vpn_gateway` - fixing a bug where the Data Source didn't return an error when the VPN Gateway was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* Data Source: `azurerm_web_pubsub` - fixing a bug where the Data Source didn't return an error when the Web PubSub was not found ([#20479](https://github.com/hashicorp/terraform-provider-azurerm/issues/20479))
* `azurerm_backup_protected_vm` - will now correctly delete ([#20469](https://github.com/hashicorp/terraform-provider-azurerm/issues/20469))
* `azurerm_eventhub` - changing the `partition_count` property now works by creating a new resource ([#20480](https://github.com/hashicorp/terraform-provider-azurerm/issues/20480))
* `azurerm_eventgrid_domain_topic` - the `name` property can now be up to 128 characters ([#20407](https://github.com/hashicorp/terraform-provider-azurerm/issues/20407))
* `azurerm_kubernetes_cluster` - parsing the API response for the `log_analytics_workspace_id` field case-insensitively ([#20484](https://github.com/hashicorp/terraform-provider-azurerm/issues/20484))
* `azurerm_private_endpoint` - normalizing the `private_connection_resource_id` propety for a redis cache ([#20418](https://github.com/hashicorp/terraform-provider-azurerm/issues/20418))
* `azurerm_private_endpoint` - consistently normalizing the value returned from the API for `private_connection_resource_id` ([#20452](https://github.com/hashicorp/terraform-provider-azurerm/issues/20452))
* `azurerm_recovery_services_vault` - updating `cross_region_restore_enabled` to `false` recreates the resource since this operation isn't supported by the API ([#20406](https://github.com/hashicorp/terraform-provider-azurerm/issues/20406))
* `azurerm_storage_management_policy` - the `rule.filters` property is now Required since storage management policies fail if it's unspecified ([#20448](https://github.com/hashicorp/terraform-provider-azurerm/issues/20448))

## 3.43.0 (February 09, 2023)

FEATURES

* **New Data Source:** `azurerm_container_app_environment` ([#18008](https://github.com/hashicorp/terraform-provider-azurerm/issues/18008))
* **New Data Source:** `azurerm_container_app_environment_certificate` ([#18008](https://github.com/hashicorp/terraform-provider-azurerm/issues/18008))
* **New Data Source:** `azurerm_mobile_network` ([#20128](https://github.com/hashicorp/terraform-provider-azurerm/issues/20128))
* **New Resource:** `azurerm_container_app_environment` ([#18008](https://github.com/hashicorp/terraform-provider-azurerm/issues/18008))
* **New Resource:** `azurerm_container_app_environment_storage` ([#18008](https://github.com/hashicorp/terraform-provider-azurerm/issues/18008))
* **New Resource:** `azurerm_container_app_environment_dapr_component` ([#18008](https://github.com/hashicorp/terraform-provider-azurerm/issues/18008))
* **New Resource:** `azurerm_container_app_environment_certificate` ([#18008](https://github.com/hashicorp/terraform-provider-azurerm/issues/18008))
* **New Resource:** `azurerm_container_app` ([#18008](https://github.com/hashicorp/terraform-provider-azurerm/issues/18008))
* **New Resource:** `azurerm_machine_learning_datastore_fileshare` ([#19934](https://github.com/hashicorp/terraform-provider-azurerm/issues/19934))
* **New Resource:** `azurerm_machine_learning_datastore_datalake_gen2` ([#20045](https://github.com/hashicorp/terraform-provider-azurerm/issues/20045))
* **New Resource:** `azurerm_mobile_network` ([#20128](https://github.com/hashicorp/terraform-provider-azurerm/issues/20128))
* **New Resource:** `azurerm_sentinel_data_connector_microsoft_threat_intelligence` ([#20273](https://github.com/hashicorp/terraform-provider-azurerm/issues/20273))

ENHANCEMENTS:

* dependencies: updating to `v0.11.28` of `github.com/Azure/go-autorest/autorest` ([#20272](https://github.com/hashicorp/terraform-provider-azurerm/issues/20272))
* dependencies: updating to `v0.50.0` of `github.com/hashicorp/go-azure-helpers` ([#20272](https://github.com/hashicorp/terraform-provider-azurerm/issues/20272))
* dependencies: updating to `v0.20230208.1165725` of `github.com/hashicorp/go-azure-sdk` ([#20381](https://github.com/hashicorp/terraform-provider-azurerm/issues/20381))
* dependencies: updating to `v0.55.0` of `github.com/manicminer/hamilton` ([#20272](https://github.com/hashicorp/terraform-provider-azurerm/issues/20272))
* dependencies: updating to `v0.20230208.1135849` of `github.com/tombuildsstuff/kermit` ([#20381](https://github.com/hashicorp/terraform-provider-azurerm/issues/20381))
* dependences: updating `postgresql/2021-06-01/databases` to 2022-12-01 ([#20369](https://github.com/hashicorp/terraform-provider-azurerm/issues/20369))
* `appservice`: updating to API Version `2021-03-01` ([#20349](https://github.com/hashicorp/terraform-provider-azurerm/issues/20349))
* `azurestackhci`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20318](https://github.com/hashicorp/terraform-provider-azurerm/issues/20318))
* `batch`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20375](https://github.com/hashicorp/terraform-provider-azurerm/issues/20375))
* `databricks`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20309](https://github.com/hashicorp/terraform-provider-azurerm/issues/20309))
* `datadog`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20311](https://github.com/hashicorp/terraform-provider-azurerm/issues/20311))
* `databoxedge`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20236](https://github.com/hashicorp/terraform-provider-azurerm/issues/20236))
* `digitaltwins`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20318](https://github.com/hashicorp/terraform-provider-azurerm/issues/20318))
* `postgresql`: updating to API Version `2022-12-01` ([#20367](https://github.com/hashicorp/terraform-provider-azurerm/issues/20367))
* `redis`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20313](https://github.com/hashicorp/terraform-provider-azurerm/issues/20313))
* `azurerm_media_streaming_locator` - support for the `filter_names` property ([#20274](https://github.com/hashicorp/terraform-provider-azurerm/issues/20274))
* `azurerm_media_live_event_output` - support for the `rewind_window_duration` property ([#20271](https://github.com/hashicorp/terraform-provider-azurerm/issues/20271))
* `azurerm_media_streaming_live_event` - support for the `stream_options` property ([#20254](https://github.com/hashicorp/terraform-provider-azurerm/issues/20254))
* `azurerm_storage_blob_inventory_policy` - support for the `exclude_prefixes` property ([#20281](https://github.com/hashicorp/terraform-provider-azurerm/issues/20281))
* `azurerm_sentinel_alert_rule_nrt` - support for the `dynamic_property` block ([#20212](https://github.com/hashicorp/terraform-provider-azurerm/issues/20212))
* `azurerm_sentinel_alert_rule_nrt` - support for the `sentinel_entity_mapping` block ([#20230](https://github.com/hashicorp/terraform-provider-azurerm/issues/20230))
* `azurerm_sentinel_alert_rule_nrt` - support for the `event_grouping` block ([#20231](https://github.com/hashicorp/terraform-provider-azurerm/issues/20231))
* `azurerm_sentinel_alert_rule_scheduled` - support for the `dynamic_property` block ([#20212](https://github.com/hashicorp/terraform-provider-azurerm/issues/20212))
* `azurerm_sentinel_alert_rule_scheduled` - support for the `sentinel_entity_mapping` block ([#20230](https://github.com/hashicorp/terraform-provider-azurerm/issues/20230))
* `azurerm_shared_image` - support for the `confidential_vm_supported` and `confidential_vm_enabled` properties ([#20249](https://github.com/hashicorp/terraform-provider-azurerm/issues/20249))
* `azurerm_postgresql_flexible_server` - support for `replication_role` and new enum value `Replica` for `create_mode` ([#20364](https://github.com/hashicorp/terraform-provider-azurerm/issues/20364))

BUG FIXES:

* `azurerm_custom_provider` - switching a spurious usage of `Azure/azure-sdk-for-go` to `hashicorp/go-azure-sdk` ([#20315](https://github.com/hashicorp/terraform-provider-azurerm/issues/20315))
* `azurerm_function_app_function` - prevent a bug with multiple file blocks resulting in last file being used for all entries ([#20198](https://github.com/hashicorp/terraform-provider-azurerm/issues/20198))
* `azurerm_monitor_diagnostic_setting` - changing the `storage_account_id`, `eventhub_authorization_rule_id`, and `eventhub_name` properties no longer creates a new resource ([#20307](https://github.com/hashicorp/terraform-provider-azurerm/issues/20307))
* `azurerm_redis_enterprise_cluster` - switching a spurious usage of `Azure/azure-sdk-for-go` to `hashicorp/go-azure-sdk` ([#20314](https://github.com/hashicorp/terraform-provider-azurerm/issues/20314))
* `azurerm_service_fabric_managed_cluster`  - fixpotential panic when setting `node_type` ([#20345](https://github.com/hashicorp/terraform-provider-azurerm/issues/20345))
* `azurerm_web_application_firewall_policy` - prevent a failure caused by changing the order of the `disabled_rules` properties ([#20285](https://github.com/hashicorp/terraform-provider-azurerm/issues/20285))
* `azurerm_databricks_access_connector` - `name` can now be up to 64 character in length ([#20353](https://github.com/hashicorp/terraform-provider-azurerm/issues/20353))

## 3.42.0 (February 02, 2023)

FEATURES

* **New Resource:** `azurerm_ip_group_cidr` ([#20225](https://github.com/hashicorp/terraform-provider-azurerm/issues/20225))
* **New Resource:** `azurerm_network_manager_connectivity_configuration` ([#20133](https://github.com/hashicorp/terraform-provider-azurerm/issues/20133))

ENHANCEMENTS:

* dependencies: updating to `v0.20230130.1140358 ` of `github.com/hashicorp/go-azure-sdk` ([#20293](https://github.com/hashicorp/terraform-provider-azurerm/issues/20293))
* `databasemigration`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20214](https://github.com/hashicorp/terraform-provider-azurerm/issues/20214))
* `servicefabric`: refactoring to use github.com/hashicorp/go-azure-sdk ([#20202](https://github.com/hashicorp/terraform-provider-azurerm/issues/20202))
* `azurerm_kubernetes_cluster` - add support for the `confidential_computing` add-on ([#20194](https://github.com/hashicorp/terraform-provider-azurerm/issues/20194))
* `azurerm_kubernetes_cluster` - export the identity for the `aci_connector_linux` add-on ([#20194](https://github.com/hashicorp/terraform-provider-azurerm/issues/20194))
* `azurerm_lb_backend_address_pool` - support for the `virtual_network_id` property ([#20205](https://github.com/hashicorp/terraform-provider-azurerm/issues/20205))
* `azurerm_postgresql_flexible_server`: add default value for `authentication.active_directory_auth_enabled` and `authentication.password_auth_enabled` ([#20054](https://github.com/hashicorp/terraform-provider-azurerm/issues/20054))
* `azurerm_site_recovery_protection_container_mapping` - support for the `automatic_update` block ([#19710](https://github.com/hashicorp/terraform-provider-azurerm/issues/19710))
* `azurerm_site_recovery_replicated_vm` - support for the `unmanaged_disk`, `target_proximity_placement_group_id`, `target_boot_diag_storage_account_id`,  `target_capacity_reservation_group_id`, `target_virtual_machine_scale_set_id`, `multi_vm_group_name`, `target_edge_zone`, and `test_network_id` properties ([#19939](https://github.com/hashicorp/terraform-provider-azurerm/issues/19939))

BUG FIXES:

* `data.azurerm_monitor_data_collection_rule` - raises an error when the specified data collection rule can't be found ([#20282](https://github.com/hashicorp/terraform-provider-azurerm/issues/20282))
* `azurerm_federated_identity_credential` - prevent a perpetual diff ([#20219](https://github.com/hashicorp/terraform-provider-azurerm/issues/20219))
* `azurerm_linux_function_app`  - fix`linuxFxVersion` for docker `registry_url` processing ([#18194](https://github.com/hashicorp/terraform-provider-azurerm/issues/18194))
* `azurerm_monitor_aad_diagnostic_setting` - the field `log_analytics_workspace_id` is now parsed case-insensitively from the API Response ([#20206](https://github.com/hashicorp/terraform-provider-azurerm/issues/20206))

## 3.41.0 (January 26, 2023)

FEATURES

* **New Data Source:** `azurerm_key_vault_certificates` ([#19498](https://github.com/hashicorp/terraform-provider-azurerm/issues/19498))
* **New Data Source:** `azurerm_site_recovery_replication_recovery_plan` ([#19940](https://github.com/hashicorp/terraform-provider-azurerm/issues/19940))
* **New Resource:** `azurerm_orbital_contact` ([#19036](https://github.com/hashicorp/terraform-provider-azurerm/issues/19036))
* **New Resource:** `azurerm_site_recovery_replication_recovery_plan` ([#19940](https://github.com/hashicorp/terraform-provider-azurerm/issues/19940))

ENHANCEMENTS:

* dependencies: updating to `v0.20230124.1111819` of `github.com/hashicorp/go-azure-sdk` ([#20160](https://github.com/hashicorp/terraform-provider-azurerm/issues/20160))
* resourceproviders: no longer registering `Microsoft.ServiceFabricMesh` by default ([#20165](https://github.com/hashicorp/terraform-provider-azurerm/issues/20165))
* testing: refactoring to use `hashicorp/terraform-plugin-testing` ([#20114](https://github.com/hashicorp/terraform-provider-azurerm/issues/20114))
* `devtestlabs`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20139](https://github.com/hashicorp/terraform-provider-azurerm/issues/20139))
* `logic`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#20144](https://github.com/hashicorp/terraform-provider-azurerm/issues/20144))
* `network`: updating to API version `2022-07-01` ([#20097](https://github.com/hashicorp/terraform-provider-azurerm/issues/20097))
* `postgresql`: updating to API version `2022-03-08-preview` ([#20073](https://github.com/hashicorp/terraform-provider-azurerm/issues/20073))
* `streamanalytics`: updating to API Version `2021-10-01-preview` ([#20145](https://github.com/hashicorp/terraform-provider-azurerm/issues/20145))
* `azurerm_*_app_slot` - support for slots to be placed in different service plans ([#20184](https://github.com/hashicorp/terraform-provider-azurerm/issues/20184))
* `azurerm_databricks_workspace` - support for customer managed keys for managed disks attached to the workspace ([#19992](https://github.com/hashicorp/terraform-provider-azurerm/issues/19992))
* `azurerm_databricks_workspace` - support for updating the properties `public_network_access_enabled`, `network_security_group_rules_required` and ` managed_services_cmk_key_vault_key_id` ([#19992](https://github.com/hashicorp/terraform-provider-azurerm/issues/19992))
* `azurerm_kubernetes_cluster` - support for `node_public_ip_tags` ([#19731](https://github.com/hashicorp/terraform-provider-azurerm/issues/19731))
* `azurerm_kubernetes_cluster_node_pool` - support for `node_public_ip_tags` ([#19731](https://github.com/hashicorp/terraform-provider-azurerm/issues/19731))
* `azurerm_log_analytics_workspace` - support for the `local_authentication_disabled` property ([#20092](https://github.com/hashicorp/terraform-provider-azurerm/issues/20092))
* `azurerm_postgresql_flexible_server` - support for customer managed keys ([#20086](https://github.com/hashicorp/terraform-provider-azurerm/issues/20086))
* `azurerm_storage_account` - support for `AADKERB` to `azure_files_authentication.0.directory_type` ([#20168](https://github.com/hashicorp/terraform-provider-azurerm/issues/20168))

BUG FIXES:

* `azurerm_stream_analytics_output_servicebus_queue` - shared access policy name and key are now optional for MSI authentication ([#19712](https://github.com/hashicorp/terraform-provider-azurerm/issues/19712))
* `azurerm_stream_analytics_output_servicebus_topic` - shared access policy name and key are now optional for MSI authentication ([#19708](https://github.com/hashicorp/terraform-provider-azurerm/issues/19708))

## 3.40.0 (January 19, 2023)

FEATURES

* **New Data Source:** `azurerm_bastion_host` ([#20062](https://github.com/hashicorp/terraform-provider-azurerm/issues/20062))
* **New Resource:** `azurerm_lab_service_schedule` ([#19977](https://github.com/hashicorp/terraform-provider-azurerm/issues/19977))
* **New Resource:** `azurerm_machine_learning_datastore_blobstorage` ([#19909](https://github.com/hashicorp/terraform-provider-azurerm/issues/19909))
* **New Resource:** `azurerm_network_manager_scope_connection` ([#19610](https://github.com/hashicorp/terraform-provider-azurerm/issues/19610))
* **New Resource:** `azurerm_network_manager_static_member` ([#20077](https://github.com/hashicorp/terraform-provider-azurerm/issues/20077))
* **New Resource:** `azurerm_sentinel_log_analytics_workspace_onboarding` ([#19692](https://github.com/hashicorp/terraform-provider-azurerm/issues/19692))

ENHANCEMENTS:

* dependencies: updating to `v0.20230117.1125206` of `github.com/hashicorp/go-azure-sdk` ([#20081](https://github.com/hashicorp/terraform-provider-azurerm/issues/20081))
* `azurerm_application_gateway` - support for TLS 1.3 and CustomV2 ([#20029](https://github.com/hashicorp/terraform-provider-azurerm/issues/20029))
* `azurerm_kubernetes_cluster` - support for the `key_management_service` block ([#19893](https://github.com/hashicorp/terraform-provider-azurerm/issues/19893))
* `azurerm_linux_web_app` - support for Python 3.11  ([#20001](https://github.com/hashicorp/terraform-provider-azurerm/issues/20001))
* `azurerm_linux_web_app_slot` - support for Python 3.11 ([#20001](https://github.com/hashicorp/terraform-provider-azurerm/issues/20001))
* `azurerm_ip_group` - support for the `firewall_ids` and `firewall_policy_ids` properties ([#19845](https://github.com/hashicorp/terraform-provider-azurerm/issues/19845))
* `azurerm_recovery_services_vault` - support for the `immutability`, user assigned `identity` and `use_system_assigned_identity` properties ([#20109](https://github.com/hashicorp/terraform-provider-azurerm/issues/20109))
* `azurerm_synapse_sql_pool` - support for `geo_backup_policy_enabled` and fix `recovery_database_id` [([#20010](https://github.com/hashicorp/terraform-provider-azurerm/issues/20010))

BUG FIXES: 

* Data Source: `azurerm_batch_pool` - the field `password` is now correctly marked as a sensitive value ([#20061](https://github.com/hashicorp/terraform-provider-azurerm/issues/20061))
* Data Source: `azurerm_batch_pool` - the field `ssh_private_key ` is now correctly marked as a sensitive value ([#20061](https://github.com/hashicorp/terraform-provider-azurerm/issues/20061))
* `azurerm_api_management_identity_provider_twitter` - the field `api_key` is now correctly marked as a sensitive value ([#20061](https://github.com/hashicorp/terraform-provider-azurerm/issues/20061))
* `azurerm_cdn_frontdoor_origin_group` - shim SDK to allow `health_probe` to be passed as `null` ([#20015](https://github.com/hashicorp/terraform-provider-azurerm/issues/20015))
* `azurerm_container_group`  - fix dynamic setting `dns_config` crash issue ([#20002](https://github.com/hashicorp/terraform-provider-azurerm/issues/20002))
* `azurerm_container_registry_task` - the field `password` is now correctly marked as a sensitive value ([#20061](https://github.com/hashicorp/terraform-provider-azurerm/issues/20061))
* `azurerm_dev_test_windows_virtual_machine` - the `password` field is now correctly marked as a sensitive value ([#20061](https://github.com/hashicorp/terraform-provider-azurerm/issues/20061))
* `azurerm_federated_identity_credential` - preent concurrent write to parent resource ([#20003](https://github.com/hashicorp/terraform-provider-azurerm/issues/20003))
* `azurerm_linux_web_app_slot`  - fixa bug where `use_32_bit_worker` would not be set correctly ([#20051](https://github.com/hashicorp/terraform-provider-azurerm/issues/20051))
* `azurerm_postgresql_flexible_server_configuration` - restart server when required ([#20044](https://github.com/hashicorp/terraform-provider-azurerm/issues/20044))
* `azurerm_kubernetes_cluster` - prevent a possible panic while importing ([#20107](https://github.com/hashicorp/terraform-provider-azurerm/issues/20107))
* `azurerm_service_fabric_managed_cluster` - the `password` field is now correctly marked as a sensitive value ([#20061](https://github.com/hashicorp/terraform-provider-azurerm/issues/20061))
* `azurerm_service_fabric_managed_cluster` - the `resource_group_name` field is now correctly marked as ForceNew ([#20061](https://github.com/hashicorp/terraform-provider-azurerm/issues/20061))
* `azurerm_spring_cloud_configuration_service ` - the field `password` is now correctly marked as a sensitive value ([#20061](https://github.com/hashicorp/terraform-provider-azurerm/issues/20061))
* `azurerm_spring_cloud_configuration_service ` - the field `private_key` is now correctly marked as a sensitive value ([#20061](https://github.com/hashicorp/terraform-provider-azurerm/issues/20061))
* `azurerm_static_site` - the field `api_key` is now correctly marked as a sensitive value ([#20061](https://github.com/hashicorp/terraform-provider-azurerm/issues/20061))
* `azurerm_storage_account` - will no longer silently ignore `404` error while reading service properties ([#19062](https://github.com/hashicorp/terraform-provider-azurerm/issues/19062))
* `azurerm_storage_account` - the `infrastructure_encryption_enabled` is now supportted for premium accounts ([#20028](https://github.com/hashicorp/terraform-provider-azurerm/issues/20028))
* `azurerm_windows_web_app_slot`  - fixa bug where `use_32_bit_worker` would not be set correctly ([#20051](https://github.com/hashicorp/terraform-provider-azurerm/issues/20051))

## 3.39.1 (January 13, 2023)

BUG FIXES:

* `azurerm_cosmosdb_sql_container`  - fixproperty `included_path` can not be removed issue ([#19998](https://github.com/hashicorp/terraform-provider-azurerm/issues/19998))
* `azurerm_log_analytics `- fixing crash during read ([#20011](https://github.com/hashicorp/terraform-provider-azurerm/issues/20011))

## 3.39.0 (January 12, 2023)

BREAKING CHANGES:

* **App Service App Stack Re-alignment** - due to a number of changes in how the Service manages App and Stack settings, the Terraform resource schema and validation needs to be updated to re-align with the service. Whist we ordinarily avoid breaking changes outside a major release, the drift has made several aspects of these resources in an unworkable position resulting in a poor experience for many users ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))

* `azurerm_windows_web_app`
    * `node_version` Valid values are now `~12`, `~14`, `~16`, and  `~18`. This is due to an underlying change to where the Service reads the Node value from in the API request.
    * `dotnet_version` valid values are now `v2.0`, `v3.0`, `v4.0`, `v5.0`, `v6.0`, and `v7.0`
    * New setting `dotnet_core_version` - Valid values are `v4.0`. This setting replaces the hybrid setting of `core3.1` in `dotnet_version` since the removal of core3.1 from the supported versions.
    * `tomcat_version` - Configured the Web App to use Tomcat as the JWS at the specified version. See the official docs for supported versions. Examples include `10.0`, and `10.0.20`
    * `java_embedded_server_enabled` - configures the JWS to be the Embedded server at the version specified by `java_version`. Defaults to `false`. Note: One of `java_embedded_server_enabled` or `tomcat_version` is required when `java_version` is set.

* `azurerm_windows_web_app_slot`
    * `node_version` Valid values are now `~12`, `~14`, `~16`, and  `~18`. This is due to an underlying change to where the Service reads the Node value from in the API request.
    * `dotnet_version` valid values are now `v2.0`, `v3.0`, `v4.0`, `v5.0`, `v6.0`, and `v7.0`
    * New setting `dotnet_core_version` - Valid values are `v4.0`. This setting replaces the hybrid setting of `core3.1` in `dotnet_version` since the removal of core3.1 from the supported versions.
    * `tomcat_version` - Configured the Web App to use Tomcat as the JWS at the specified version. See the official docs for supported versions. Examples include `10.0`, and `10.0.20`
    * `java_embedded_server_enabled` - configures the JWS to be the Embedded server at the version specified by `java_version`. Defaults to `false`. Note: One of `java_embedded_server_enabled` or `tomcat_version` is required when `java_version` is set.

* `azurerm_windows_function_app`
    * `dotnet_version` - Valid values are now `v3.0`, `v4.0`, `v6.0`, and `v7.0`, defaulting to `v4.0`
    * `java_version` - Valid values are now `1.8`, `11`, and `17`

* `azurerm_windows_function_app_slot`
    * `dotnet_version` - Valid values are now `v3.0`, `v4.0`, `v6.0`, and `v7.0`, defaulting to `v4.0`
    * `java_version` - Valid values are now `1.8`, `11`, and `17`

* `azurerm_linux_web_app`
    * `java_version` - input validation has been introduced based on supported values within the service. Valid values are now: `8`,`11`, and `17`. 

FEATURES:

* **New Data Source:** `azurerm_private_dns_resolver` ([#19885](https://github.com/hashicorp/terraform-provider-azurerm/issues/19885))
* **New Data Source:** `azurerm_private_dns_resolver_dns_forwarding_ruleset` ([#19941](https://github.com/hashicorp/terraform-provider-azurerm/issues/19941))
* **New Data Source:** `azurerm_private_dns_resolver_forwarding_rule` ([#19947](https://github.com/hashicorp/terraform-provider-azurerm/issues/19947))
* **New Data Source:** `azurerm_private_dns_resolver_inbound_endpoint` ([#19948](https://github.com/hashicorp/terraform-provider-azurerm/issues/19948))
* **New Data Source:** `azurerm_private_dns_resolver_outbound_endpoint` ([#19950](https://github.com/hashicorp/terraform-provider-azurerm/issues/19950))
* **New Data Source:** `azurerm_private_dns_resolver_virtual_network_link` ([#19951](https://github.com/hashicorp/terraform-provider-azurerm/issues/19951))
* **New Resource:** `azurerm_application_insights_standard_web_test` ([#19954](https://github.com/hashicorp/terraform-provider-azurerm/issues/19954))
* **New Resource:** `azurerm_cost_anomaly_alert` ([#19899](https://github.com/hashicorp/terraform-provider-azurerm/issues/19899))
* **New Resource:** `azurerm_lab_service_lab` ([#19852](https://github.com/hashicorp/terraform-provider-azurerm/issues/19852))
* **New Resource:** `azurerm_lab_service_user` ([#19957](https://github.com/hashicorp/terraform-provider-azurerm/issues/19957))
* **New Resource:** `azurerm_network_manager_subscription_connection` ([#19617](https://github.com/hashicorp/terraform-provider-azurerm/issues/19617))
* **New Resource:** `azurerm_network_manager_management_group_connection` ([#19621](https://github.com/hashicorp/terraform-provider-azurerm/issues/19621))
* **New Resource:** `azurerm_media_services_account_filter` ([#19964](https://github.com/hashicorp/terraform-provider-azurerm/issues/19964))
* **New Resource:** `azurerm_private_endpoint_application_security_group_association` ([#19825](https://github.com/hashicorp/terraform-provider-azurerm/issues/19825))
* **New Resource:** `azurerm_sentinel_data_connector_threat_intelligence_taxii` ([#19209](https://github.com/hashicorp/terraform-provider-azurerm/issues/19209))
* **New Resource:** `azurerm_storage_account_local_user` ([#19592](https://github.com/hashicorp/terraform-provider-azurerm/issues/19592))

ENHANCEMENTS:

* `siterecovery`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#19571](https://github.com/hashicorp/terraform-provider-azurerm/issues/19571))
* `siterecovery`: updating to API version `2021-11-01` ([#19571](https://github.com/hashicorp/terraform-provider-azurerm/issues/19571))
* Data Source: `azurerm_shared_image` - add support for the `purchase_plan` block ([#19873](https://github.com/hashicorp/terraform-provider-azurerm/issues/19873))
* `azurerm_kubernetes_cluster` - add support for the `vnet_integration_enabled` and `subnet_id` properties ([#19438](https://github.com/hashicorp/terraform-provider-azurerm/issues/19438))
* `azurerm_log_analytics_data_export_rule` - `destination_resource_id` accepts an Event Hub Namespace ID ([#19868](https://github.com/hashicorp/terraform-provider-azurerm/issues/19868))
* `azurerm_linux_web_app`- support for the `application_stack.go_version` property ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_linux_web_app_slot` -support for the `application_stack.go_version` property ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_logic_app_action_http` - add support for `@` in the `body` property ([#19754](https://github.com/hashicorp/terraform-provider-azurerm/issues/19754))
* `azurerm_maintenance_configuration` - support for the `in_guest_user_patch_mode` and `install_patches` properties ([#19865](https://github.com/hashicorp/terraform-provider-azurerm/issues/19865))
* `azurerm_monitor_diagnostic_setting` - deprecate `log` in favour of `enabled_log` ([#19504](https://github.com/hashicorp/terraform-provider-azurerm/issues/19504))
* `azurerm_media_services_account` - support for the `encryption` and `public_network_access_enabled` properties ([#19891](https://github.com/hashicorp/terraform-provider-azurerm/issues/19891))
* `azurerm_mysql_flexible_server` - support for the `customer_managed_key` properties ([#19905](https://github.com/hashicorp/terraform-provider-azurerm/issues/19905))
* `azurerm_sentinel_automation_rule` - support for the `triggers_on`, `triggers_when`, and `condition_json` properties ([#19309](https://github.com/hashicorp/terraform-provider-azurerm/issues/19309))
* `azurerm_spring_cloud_gateway` - support for the `application_performance_monitoring_types`, `environment_variables`, and `sensitive_environment_variables` properties ([#19884](https://github.com/hashicorp/terraform-provider-azurerm/issues/19884))
* `azurerm_storage_account` - support for the `allowed_copy_scope` property ([#19906](https://github.com/hashicorp/terraform-provider-azurerm/issues/19906))
* `azurerm_storage_queue` - exporting `resource_manager_id` ([#19969](https://github.com/hashicorp/terraform-provider-azurerm/issues/19969))
* `azurerm_synapse_spark_pool` - add support for Spark 3.3 ([#19866](https://github.com/hashicorp/terraform-provider-azurerm/issues/19866))
* `azurerm_windows_web_app` - the `php_version` property supported values now include: `7.1`, `7.4`, and `Off`. Note: `7.1` is currently deprecated. `Off` will configure the system to use the latest available to the App service image ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app` - the `python_version` property has been deprecated and is no longer used by the service  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app` - the `python` property supersedes `python_version`. Defaults to `false`. When true uses the latest Python version supported by the Windows App image  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app` - the `java_container` property has been deprecated in favour of `tomcat_version` and `java_embedded_server_enabled`  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app` - the `java_container_version` property This property has been deprecated in favour of `tomcat_version` and `java_embedded_server_enabled`  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app` - the `current_stack` property will now be computed if only one stack is configured on the Windows Web App. This will ensure the portal displays the appropriate metadata and configuration for this stack  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app` - Added input validation for `interval` values in the `auto_heal` block. These properties now enforce HH:MM:SS values up to `99:59:59` ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app_slot` - the `php_version` property supported values now include: `7.1`, `7.4`, and `Off`. Note: `7.1` is currently deprecated. `Off` will configure the system to use the latest available to the App service image ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app_slot` - the `python_version` property has been deprecated and is no longer used by the service  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app_slot` - the `python` property supersedes `python_version`. Defaults to `false`. When true uses the latest Python version supported by the Windows App image  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app_slot` - the `java_container` property has been deprecated in favour of `tomcat_version` and `java_embedded_server_enabled`  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app_slot` - the `java_container_version` property This property has been deprecated in favour of `tomcat_version` and `java_embedded_server_enabled`  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app_slot` - the `current_stack` property will now be computed if only one stack is configured on the Windows Web App. This will ensure the portal displays the appropriate metadata and configuration for this stack  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app_slot` - Added input validation for `interval` values in the `auto_heal` block. These properties now enforce HH:MM:SS values up to `99:59:59` ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))

BUG FIXES:

* `azurerm_app_configuration_feature` - handle updates correctly where the label ID is omitted ([#19900](https://github.com/hashicorp/terraform-provider-azurerm/issues/19900))
* `azurerm_cdn_frontdoor_rule` - handle empty string value for `query_string` ([#19927](https://github.com/hashicorp/terraform-provider-azurerm/issues/19927))
* `azurerm_cosmosdb_account` - `default_identity_type` is now computed to allow for restores ([#19956](https://github.com/hashicorp/terraform-provider-azurerm/issues/19956))
* `azurerm_linux_web_app`- prevent a bug where `backup_config` could silently fail to expand resulting in the config not being sent ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_linux_web_app`- prevent a bug where `health_check_eviction_time_in_min` would not be correctly read back from the service ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_linux_web_app_slot`- prevent a bug where `backup_config` could silently fail to expand resulting in the config not being sent ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_linux_web_app_slot`- prevent a bug where `health_check_eviction_time_in_min` would not be correctly read back from the service ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_policy_set_definition`  - fixupdate of for empty group names in `policy_definition_reference.policy_group_names` ([#19890](https://github.com/hashicorp/terraform-provider-azurerm/issues/19890))
* `azurerm_storage_account` -  `403` is now a valid status code for when permissions to list keys is missing ([#19645](https://github.com/hashicorp/terraform-provider-azurerm/issues/19645))
* `azurerm_storage_data_lake_gen2_path` - `ace` generated by default are no longer stored in state to prevent perpetual state diffs ([#18494](https://github.com/hashicorp/terraform-provider-azurerm/issues/18494))
* `azurerm_storage_data_lake_gen2_filesystem` - `ace` generated by default are no longer stored in state to prevent perpetual state diffs ([#18494](https://github.com/hashicorp/terraform-provider-azurerm/issues/18494))
* `azurerm_web_pubsub_hub` - the `event_handler` property is now a list instead of set to respect the input order ([#19886](https://github.com/hashicorp/terraform-provider-azurerm/issues/19886))
* `azurerm_windows_web_app` - prevent a bug where `backup_config` could silently fail to expand resulting in the config not being sent  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app` - prevent a bug where `health_check_eviction_time_in_min` would not be correctly set on Crete or Update  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app_slot` - prevent a bug where `backup_config` could silently fail to expand resulting in the config not being sent  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))
* `azurerm_windows_web_app_slot` - prevent a bug where `health_check_eviction_time_in_min` would not be correctly set on Crete or Update  ([#19685](https://github.com/hashicorp/terraform-provider-azurerm/issues/19685))

## 3.38.0 (January 05, 2023)

FEATURES:

* **New Data Source:** `azurerm_marketplace_agreement` ([#19628](https://github.com/hashicorp/terraform-provider-azurerm/issues/19628))
* **New Data Source:** `azurerm_network_manager_network_group` ([#19593](https://github.com/hashicorp/terraform-provider-azurerm/issues/19593))
* **New Data Source:** `azurerm_virtual_hub_route_table` ([#19628](https://github.com/hashicorp/terraform-provider-azurerm/issues/19628))

ENHANCEMENTS:

* dependencies: updating to `v0.20230105.1121404` of `github.com/hashicorp/go-azure-sdk` ([#19872](https://github.com/hashicorp/terraform-provider-azurerm/issues/19872))
* dependencies: updating to `v0.20221207.1110610` of `github.com/tombuildsstuff/kermit` ([#19698](https://github.com/hashicorp/terraform-provider-azurerm/issues/19698))
* `azurerm_dedicated_host` - add support for`LSv3-Type1` type ([#19875](https://github.com/hashicorp/terraform-provider-azurerm/issues/19875))
* `azurerm_proximity_placement_group` - support for the `allowed_vm_sizes` and `zone` properties ([#19675](https://github.com/hashicorp/terraform-provider-azurerm/issues/19675))

BUG FIXES:

* `azurerm_automation_software_update_configuration` - correctly handle empty `expiry_time` api values ([#19774](https://github.com/hashicorp/terraform-provider-azurerm/issues/19774))
* `azurerm_app_service_connection` - polling until the resource is fully created, updated and deleted ([#19792](https://github.com/hashicorp/terraform-provider-azurerm/issues/19792))
* `azurerm_batch_pool` - correctly handle the resource being deleted outside of terraform ([#19780](https://github.com/hashicorp/terraform-provider-azurerm/issues/19780))
* `azurerm_billing_account_cost_management_export` - marking the resource as gone when it's no longer present in Azure ([#19871](https://github.com/hashicorp/terraform-provider-azurerm/issues/19871))
* `azurerm_bot_service_azure_bot` - marking the resource as gone when it's no longer present in Azure ([#19871](https://github.com/hashicorp/terraform-provider-azurerm/issues/19871))
* `azurerm_databricks_access_connector` - polling until the resource is fully created, updated and deleted ([#19792](https://github.com/hashicorp/terraform-provider-azurerm/issues/19792))
* `azurerm_databricks_access_connector` - marking the resource as gone when it's no longer present in Azure ([#19871](https://github.com/hashicorp/terraform-provider-azurerm/issues/19871))
* `azurerm_datadog_monitor_sso_configuration` - polling until the resource is fully created and deleted ([#19792](https://github.com/hashicorp/terraform-provider-azurerm/issues/19792))
* `azurerm_hdinsight_kafka_cluster` - the `kafka_management_node` property has been deprecated and will be removed in `v4.0` ([#19423](https://github.com/hashicorp/terraform-provider-azurerm/issues/19423))
* `azurerm_kubernetes_cluster` - `scale_down_mode` of the default node pool can now be updated without rebuilding the entire cluster ([#19823](https://github.com/hashicorp/terraform-provider-azurerm/issues/19823))
* `azurerm_orbital_contact_profile` - polling until the resource is fully created, updated and deleted ([#19792](https://github.com/hashicorp/terraform-provider-azurerm/issues/19792))
* `azurerm_orbital_spacecraft` - polling until the resource is fully created, updated and deleted ([#19792](https://github.com/hashicorp/terraform-provider-azurerm/issues/19792))
* `azurerm_postgresql_flexible_server` - correctly handle password authentication ([#19800](https://github.com/hashicorp/terraform-provider-azurerm/issues/19800))
* `azurerm_resource_group_cost_management_export` - marking the resource as gone when it's no longer present in Azure ([#19871](https://github.com/hashicorp/terraform-provider-azurerm/issues/19871))
* `azurerm_spring_cloud_connection` - polling until the resource is fully updated and deleted ([#19792](https://github.com/hashicorp/terraform-provider-azurerm/issues/19792))
* `azurerm_stack_hci_cluster` - polling until the resource is fully deleted ([#19792](https://github.com/hashicorp/terraform-provider-azurerm/issues/19792))
* `azurerm_stream_analytics_cluster` - polling until the resource is fully deleted ([#19792](https://github.com/hashicorp/terraform-provider-azurerm/issues/19792))
* `azurerm_stream_analytics_reference_input_blob` - the `storage_account_key` property is now optional when MSI auth is used ([#19676](https://github.com/hashicorp/terraform-provider-azurerm/issues/19676))
* `azurerm_storage_account_network_rules` - the requires import check no longer checks the `bypass` field to workaround an issue within the Azure API ([#19719](https://github.com/hashicorp/terraform-provider-azurerm/issues/19719))
* `azurerm_subscription_cost_management_export` - marking the resource as gone when it's no longer present in Azure ([#19871](https://github.com/hashicorp/terraform-provider-azurerm/issues/19871))
* `azurerm_synapse_linked_service` - report error during create/update ([#19849](https://github.com/hashicorp/terraform-provider-azurerm/issues/19849))
* `azurerm_virtual_desktop_application_group` - changing the `host_pool_id` now creates a new resource ([#19689](https://github.com/hashicorp/terraform-provider-azurerm/issues/19689))
* `azurerm_vmware_express_route_authorization` - marking the resource as gone when it's no longer present in Azure ([#19871](https://github.com/hashicorp/terraform-provider-azurerm/issues/19871))

## 3.37.0 (December 21, 2022)

FEATURES:

* **New Resource:** `azurerm_cognitive_deployment` ([#19526](https://github.com/hashicorp/terraform-provider-azurerm/issues/19526))
* **New Resource:** `azurerm_billing_account_cost_management_export` ([#19723](https://github.com/hashicorp/terraform-provider-azurerm/issues/19723))
* **New resource:** `azurerm_key_vault_certificate_contacts` ([#19743](https://github.com/hashicorp/terraform-provider-azurerm/issues/19743))
* **New Resource:** `azurerm_lab_service_plan` ([#19312](https://github.com/hashicorp/terraform-provider-azurerm/issues/19312))
* **New Resource:** `azurerm_resource_deployment_script` ([#19436](https://github.com/hashicorp/terraform-provider-azurerm/issues/19436))
* **New Resource:** `azurerm_spring_cloud_customized_accelerator` ([#19736](https://github.com/hashicorp/terraform-provider-azurerm/issues/19736))

ENHANCEMENTS:

* `azurerm_netapp_volume` - support for the `zone` property ([#19669](https://github.com/hashicorp/terraform-provider-azurerm/issues/19669))

BUG FIXES: 

* `azurerm_app_configuration_key`  - fixa regression when handling IDs containing a `:` ([#19722](https://github.com/hashicorp/terraform-provider-azurerm/issues/19722))
* `azurerm_virtual_network_gateway_connection` -  can now be created with a `azurerm_virtual_network_gateway` in another resource group ([#19699](https://github.com/hashicorp/terraform-provider-azurerm/issues/19699))

## 3.36.0 (December 15, 2022)

FEATURES:

* **New Resource:** `azurerm_virtual_machine_packet_capture` ([#19385](https://github.com/hashicorp/terraform-provider-azurerm/issues/19385))
* **New Resource:** `azurerm_virtual_machine_scale_set_packet_capture` ([#19385](https://github.com/hashicorp/terraform-provider-azurerm/issues/19385))
* **New Resource:** `azurerm_spring_cloud_accelerator` ([#19572](https://github.com/hashicorp/terraform-provider-azurerm/issues/19572))
* **New Resource:** `azurerm_spring_cloud_dev_tool_portal` ([#19568](https://github.com/hashicorp/terraform-provider-azurerm/issues/19568))
* **New Resource:** `azurerm_route_map` ([#19402](https://github.com/hashicorp/terraform-provider-azurerm/issues/19402))
* **New Data Source:** `azurerm_lb_outbound_rule` ([#19345](https://github.com/hashicorp/terraform-provider-azurerm/issues/19345))

ENHANCEMENTS:

* `healthbot`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#19433](https://github.com/hashicorp/terraform-provider-azurerm/issues/19433))
* `media`: updating to API version `2021-11-01 `and `2022-08-01` ([#19623](https://github.com/hashicorp/terraform-provider-azurerm/issues/19623))
* `azurerm_cosmosdb_account` - support  for updating some `capabilities`  ([#14991](https://github.com/hashicorp/terraform-provider-azurerm/issues/14991))
* `azurerm_key_vault_managed_hardware_security_module` - support for the `public_network_access_enabled` and `network_acls` properties ([#19640](https://github.com/hashicorp/terraform-provider-azurerm/issues/19640))
* `azurerm_kubernetes_cluster` - support for the `monitor_metrics` block ([#19530](https://github.com/hashicorp/terraform-provider-azurerm/issues/19530))
* `azurerm_kubernetes_cluster` - the `ssh_key` property can now be updated ([#19634](https://github.com/hashicorp/terraform-provider-azurerm/issues/19634))
* `azurerm_kubernetes_cluster_node_pool` - support for the `outbound_nat_enabled` property ([#19663](https://github.com/hashicorp/terraform-provider-azurerm/issues/19663))
* `azurerm_lighthouse_definition` - support for the `eligible_authorization` property ([#19569](https://github.com/hashicorp/terraform-provider-azurerm/issues/19569))
* `azurerm_log_analytics_workspace` - support for the `allow_resource_only_permissions` property ([#19346](https://github.com/hashicorp/terraform-provider-azurerm/issues/19346))
* `azurerm_private_endpoint` - support for the `member_name` property in the `ip_configuration` block and support for multiple `ip_configuration` blocks ([#19389](https://github.com/hashicorp/terraform-provider-azurerm/issues/19389))
* `azurerm_storage_account` - support for the `blob_properties.restore_policy` property ([#19644](https://github.com/hashicorp/terraform-provider-azurerm/issues/19644))
* `azurerm_vpn_gateway_connection` - support for the `inbound_route_map_id` and `outbound_route_map_id` properties ([#19681](https://github.com/hashicorp/terraform-provider-azurerm/issues/19681))
* `azurerm_point_to_site_vpn_gateway` - support for the `routing_preference_internet_enabled`, `inbound_route_map_id`, and `outbound_route_map_id` properties ([#19672](https://github.com/hashicorp/terraform-provider-azurerm/issues/19672))
* `azurerm_web_application_firewall_policy` - support the `rule` property in the `rule_group_override` block ([#19497](https://github.com/hashicorp/terraform-provider-azurerm/issues/19497))

BUG FIXES:

* Data Source: `azurerm_api_management` - prevent failure when retrieving tenant access properties when permissions are missing ([#19626](https://github.com/hashicorp/terraform-provider-azurerm/issues/19626))
* `azurerm_cdn_frontdoor_firewall_policy` - allow `Log` as a valid value for managed rule override `action` in DRS 2.0 and above ([#19637](https://github.com/hashicorp/terraform-provider-azurerm/issues/19637))
* `azurerm_cosmosdb_account` - enabling `analytical_storage_enabled` no longer forces recreation ([#19659](https://github.com/hashicorp/terraform-provider-azurerm/issues/19659))
* `azurerm_monitor_scheduled_query_rules_alert_v2` - use the correct alue `Equals` for operator ([#19594](https://github.com/hashicorp/terraform-provider-azurerm/issues/19594))
* `azurerm_mssql_database` - the `threat_detection_policy.storage_*` properties can now be correctly set as empty ([#19670](https://github.com/hashicorp/terraform-provider-azurerm/issues/19670))
* `azurerm_synapse_linked_service` - add validation for `type` ([#19636](https://github.com/hashicorp/terraform-provider-azurerm/issues/19636))
* `azurerm_resource_policy_exemption` - changing the `policy_assignment_id` property not created a new resource ([#19674](https://github.com/hashicorp/terraform-provider-azurerm/issues/19674))
* `azurerm_resource_group_policy_exemption` - changing the `policy_assignment_id` property not created a new resource ([#19674](https://github.com/hashicorp/terraform-provider-azurerm/issues/19674))
* `azurerm_subscription_policy_exemption` - changing the `policy_assignment_id` property not created a new resource ([#19674](https://github.com/hashicorp/terraform-provider-azurerm/issues/19674))
* `azurerm_stream_analytics_output_mssql` - the `user` and `password` properties are not optional when using MSI authentication ([#19696](https://github.com/hashicorp/terraform-provider-azurerm/issues/19696))

## 3.35.0 (December 09, 2022)

BREAKING CHANGES:

* `azurerm_stream_analytics_output_blob` - the field `batch_min_rows` is now an integer rather than a float due to a [breaking change in the API Specifications](https://github.com/Azure/azure-rest-api-specs) - we believe this was only previously valid as an integer, as such whilst this is a breaking change we believe this shouldn't cause an issue for most users (since the API required that this was an integer) ([#19602](https://github.com/hashicorp/terraform-provider-azurerm/issues/19602))

FEATURES:

* **New Resource:** `azurerm_digital_twins_time_series_database_connection` ([#19576](https://github.com/hashicorp/terraform-provider-azurerm/issues/19576))
* **New Resource:** `azurerm_network_manager` ([#19334](https://github.com/hashicorp/terraform-provider-azurerm/issues/19334))
* **New Resource:** `azurerm_spring_cloud_application_live_view` ([#19495](https://github.com/hashicorp/terraform-provider-azurerm/issues/19495))
* **New Resource:** `azurerm_sentinel_data_connector_microsoft_threat_protection` ([#19162](https://github.com/hashicorp/terraform-provider-azurerm/issues/19162))
* **New Resource:** `azurerm_vmware_netapp_volume_attachment` ([#19082](https://github.com/hashicorp/terraform-provider-azurerm/issues/19082))

ENHANCEMENTS:

* dependencies: updating to `v0.20221207.1121859` of `github.com/hashicorp/go-azure-sdk` ([#19602](https://github.com/hashicorp/terraform-provider-azurerm/issues/19602))
* `lighthouse`: updating to API version `2022-10-01` ([#19499](https://github.com/hashicorp/terraform-provider-azurerm/issues/19499))
* `proximityplacementgroups`: updating to API Version `2022-03-01` ([#19537](https://github.com/hashicorp/terraform-provider-azurerm/issues/19537))
* Data Source: `azurerm_kubernetes_cluster` - support for the `storage_profile` block ([#19396](https://github.com/hashicorp/terraform-provider-azurerm/issues/19396))
* `azurerm_firewall_policy` - support for the `explicit_proxy` block and `auto_learn_private_ranges_mode` property ([#19313](https://github.com/hashicorp/terraform-provider-azurerm/issues/19313))
* `azurerm_kubernetes_cluster` - support for the `custom_ca_trust_enabled` property ([#19546](https://github.com/hashicorp/terraform-provider-azurerm/issues/19546))
* `azurerm_kubernetes_cluster` - support for the `storage_profile` block ([#19396](https://github.com/hashicorp/terraform-provider-azurerm/issues/19396))
* `azurerm_kubernetes_cluster` - support for the `image_cleaner` block ([#19368](https://github.com/hashicorp/terraform-provider-azurerm/issues/19368))
* `azurerm_kubernetes_cluster` - support for the `network_plugin_mode` and `ebpf_data_plane` properties ([#19527](https://github.com/hashicorp/terraform-provider-azurerm/issues/19527))
* `azurerm_kubernetes_cluster_node_pool` - support for the `custom_ca_trust_enabled` property ([#19546](https://github.com/hashicorp/terraform-provider-azurerm/issues/19546))
* `azurerm_lb_probe` - support for the `probe_threshold` property ([#19573](https://github.com/hashicorp/terraform-provider-azurerm/issues/19573))
* `azurerm_mssql_virtual_machine` - support for the `days_of_week` property ([#19553](https://github.com/hashicorp/terraform-provider-azurerm/issues/19553))
* `azurerm_spring_cloud_gateway_route_config` - support for the `filters`, `predicates`, and `sso_validation_enabled` properties ([#19493](https://github.com/hashicorp/terraform-provider-azurerm/issues/19493))

BUG FIXES:

* Data Source: `azurerm_sentinel_alert_rule_template`: Set custom ID rather than using ID returned from API ([#19580](https://github.com/hashicorp/terraform-provider-azurerm/issues/19580))
* `azurerm_app_service_connection` - correctly pass the secret to the service ([#19519](https://github.com/hashicorp/terraform-provider-azurerm/issues/19519))
* `azurerm_automation_software_update_configuration`  - fixissue where omitting `tags`and `tag_filter` result in an error ([#19516](https://github.com/hashicorp/terraform-provider-azurerm/issues/19516))
* `azurerm_automation_source_control` - a state migration to work around the previously incorrect id casing ([#19506](https://github.com/hashicorp/terraform-provider-azurerm/issues/19506))
* `azurerm_automation_webhook` - a state migration to work around the previously incorrect id casing ([#19506](https://github.com/hashicorp/terraform-provider-azurerm/issues/19506))
* `azurerm_container_registry_webhook` - added a state migration to work around the previously incorrect id casing ([#19507](https://github.com/hashicorp/terraform-provider-azurerm/issues/19507))
* `azurerm_frontdoor_rules_engine` - a state migration to work around the previously incorrect id casing ([#19512](https://github.com/hashicorp/terraform-provider-azurerm/issues/19512))
* `azurerm_healthcare_*` - added a state migration to work around the previously incorrect id casing ([#19511](https://github.com/hashicorp/terraform-provider-azurerm/issues/19511))
* `azurerm_iothub_*` - added a state migration to work around the previously incorrect id casing ([#19524](https://github.com/hashicorp/terraform-provider-azurerm/issues/19524))
* `azurerm_key_vault` - allow for keyvaults in two different subscriptions ([#19531](https://github.com/hashicorp/terraform-provider-azurerm/issues/19531))
* `azurerm_key_vault_certificate` - skip purging  during deletion if the parent key vault has purge protection enabled ([#19528](https://github.com/hashicorp/terraform-provider-azurerm/issues/19528))
* `azurerm_key_vault_key` - skip purging  during deletion if the parent key vault has purge protection enabled ([#19528](https://github.com/hashicorp/terraform-provider-azurerm/issues/19528))
* `azurerm_key_vault_managed_hardware_security_module` - skip purging  during deletion if the parent key vault has purge protection enabled ([#19528](https://github.com/hashicorp/terraform-provider-azurerm/issues/19528))
* `azurerm_key_vault_secret` - skip purging  during deletion if the parent key vault has purge protection enabled ([#19528](https://github.com/hashicorp/terraform-provider-azurerm/issues/19528))
* `azurerm_lb` - adding/removing a frontend configuration will no longer force recreation a new resource to be created ([#19548](https://github.com/hashicorp/terraform-provider-azurerm/issues/19548))
* `azurerm_kusto_*` - added a state migration to work around the previously incorrect id casing ([#19525](https://github.com/hashicorp/terraform-provider-azurerm/issues/19525))
* `azurerm_media_services_account` - fixing an issue in the state upgrade where the Resource ID was being parsed incorrectly ([#19578](https://github.com/hashicorp/terraform-provider-azurerm/issues/19578))
* `azurerm_mssql_elasticpool` - Prevent `license_type` from being configured in specific scenarios ([#19586](https://github.com/hashicorp/terraform-provider-azurerm/issues/19586))
* `azurerm_monitor_smart_detector_alert_rule` - added a state migration to work around the previously incorrect id casing ([#19513](https://github.com/hashicorp/terraform-provider-azurerm/issues/19513))
* `azurerm_spring_cloud_*` - added a state migration to work around the previously incorrect id casing ([#19564](https://github.com/hashicorp/terraform-provider-azurerm/issues/19564))
* `azurerm_stream_analytics_output_blob` - the field `batch_min_rows` is now an integer rather than a float due to a [breaking change in the API Specifications](https://github.com/Azure/azure-rest-api-specs) - we believe this was only previously valid as an integer, as such whilst this is a breaking change we believe this shouldn't cause an issue for most users (since the API required that this was an integer) ([#19602](https://github.com/hashicorp/terraform-provider-azurerm/issues/19602))
* `azurerm_virtual_desktop_workspace_application_group_association` - set `tags` properly ([#19574](https://github.com/hashicorp/terraform-provider-azurerm/issues/19574))

## 3.34.0 (December 01, 2022)

ENHANCEMENTS:

* dependencies: updating to `v0.20221129.1175354` of `github.com/hashicorp/go-azure-sdk` ([#19483](https://github.com/hashicorp/terraform-provider-azurerm/issues/19483))
* `media`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#19285](https://github.com/hashicorp/terraform-provider-azurerm/issues/19285))
* `springcloud`: updating to use API Version `2022-11-01-preview` ([#19445](https://github.com/hashicorp/terraform-provider-azurerm/issues/19445))
* `streamanalytics`: refactoring to use `github.com/hashicorp/go-azure-sdk` ([#19395](https://github.com/hashicorp/terraform-provider-azurerm/issues/19395))
* `synapse`: refactoring to use `github.com/tombuildstuff/kermit` rather than the embedded sdk ([#19484](https://github.com/hashicorp/terraform-provider-azurerm/issues/19484))
* Data Source: `azurerm_api_management` - support for `tenant_access` property ([#19422](https://github.com/hashicorp/terraform-provider-azurerm/issues/19422))
* `azurerm_kusto_database` - supports underscores in the name ([#19466](https://github.com/hashicorp/terraform-provider-azurerm/issues/19466))
* `azurerm_managed_disk` - support for `upload_size_bytes` property ([#19458](https://github.com/hashicorp/terraform-provider-azurerm/issues/19458))
* `azurerm_monitor_activity_log_alert` - `action` is now supplied as a list instead of a set ([#19425](https://github.com/hashicorp/terraform-provider-azurerm/issues/19425))
* `azurerm_spring_cloud_gateway_route_config` - support for `protocol` property ([#19382](https://github.com/hashicorp/terraform-provider-azurerm/issues/19382))
* `azurerm_storage_account` - support for `sftp_enabled` ([#19428](https://github.com/hashicorp/terraform-provider-azurerm/issues/19428))
* `azurerm_storage_management_policy` - `tier_to_cool_after_days_since_creation_greater_than` - support for the `tier_to_cool_after_days_since_creation_greater_than`, `tier_to_archive_after_days_since_creation_greater_than`, `delete_after_days_since_creation_greater_than` properties ([#19446](https://github.com/hashicorp/terraform-provider-azurerm/issues/19446))

BUG FIXES:

* `data.azurerm_sentinel_alert_rule_template` - a  state migration to work around the previously incorrect id casing ([#19487](https://github.com/hashicorp/terraform-provider-azurerm/issues/19487))
* `azurerm_app_configuration_key` - prevent crash when retrieving the key value ([#19464](https://github.com/hashicorp/terraform-provider-azurerm/issues/19464))
* `azurerm_data_factory_linked_service_azure_file_storage` - send `host` and and `user_id` in the payload only when it's been set ([#19468](https://github.com/hashicorp/terraform-provider-azurerm/issues/19468))
* `azurerm_eventgrid_topic`  - fixsetting of fields in `input_mapping_fields` during read ([#19494](https://github.com/hashicorp/terraform-provider-azurerm/issues/19494)) 
* `azurerm_iot_security_solution` - a  state migration to work around the previously incorrect id casing ([#19489](https://github.com/hashicorp/terraform-provider-azurerm/issues/19489))
* `azurerm_monitor_autoscale_setting` - a  state migration to work around the previously incorrect id casing ([#19492](https://github.com/hashicorp/terraform-provider-azurerm/issues/19492))
* `azurerm_sentinel_automation_rule` - a  state migration to work around the previously incorrect id casing ([#19487](https://github.com/hashicorp/terraform-provider-azurerm/issues/19487))
* `azurerm_sql_active_directory_administrator` - a  state migration to work around the previously incorrect id casing ([#19486](https://github.com/hashicorp/terraform-provider-azurerm/issues/19486))
* `azurerm_stream_analytics_output_eventhub` - `shared_access_policy_key` and `shared_access_policy_name` are now optional  ([#19447](https://github.com/hashicorp/terraform-provider-azurerm/issues/19447))
* `azurerm_synapse_integration_runtime_azure` - a state migration to work around the previously incorrect id casing ([#19485](https://github.com/hashicorp/terraform-provider-azurerm/issues/19485))
* `azurerm_synapse_integration_runtime_self_hosted` - a state migration to work around the previously incorrect id casing ([#19485](https://github.com/hashicorp/terraform-provider-azurerm/issues/19485))
* `azurerm_synapse_linked_service` - a state migration to work around the previously incorrect id casing  ([#19477](https://github.com/hashicorp/terraform-provider-azurerm/issues/19477))
* `azurerm_windows_web_app`  - fixcurrentStack is being reset when other `site_config` values are changed. ([#18568](https://github.com/hashicorp/terraform-provider-azurerm/issues/18568))
* `azurerm_windows_web_app_slot`  - fixcurrentStack is being reset when other `site_config` values are changed. ([#18568](https://github.com/hashicorp/terraform-provider-azurerm/issues/18568))
* `azurerm_windows_virtual_machine_scale_set` Fix crash when upgrading `automatic_os_upgrade_policy` ([#19465](https://github.com/hashicorp/terraform-provider-azurerm/issues/19465))

## 3.33.0 (November 24, 2022)

FEATURES:

* **New Data Source:** `azurerm_cdn_frontdoor_custom_domain` ([#19357](https://github.com/hashicorp/terraform-provider-azurerm/issues/19357))
* **New Resource:** `azurerm_mssql_managed_instance_transparent_data_encryption` ([#18918](https://github.com/hashicorp/terraform-provider-azurerm/issues/18918))
* **New Resource:** `azurerm_postgresql_flexible_server_active_directory_administrator` ([#19269](https://github.com/hashicorp/terraform-provider-azurerm/issues/19269))

ENHANCEMENTS:

* build: updating to use Go `1.19.3` ([#19362](https://github.com/hashicorp/terraform-provider-azurerm/issues/19362))
* dependencies: updating to `v0.20221122.1115312` of `github.com/hashicorp/go-azure-sdk` ([#19412](https://github.com/hashicorp/terraform-provider-azurerm/issues/19412))
* dependencies: upgrading to `v2.24.1` of `github.com/hashicorp/terraform-plugin-sdk` ([#19303](https://github.com/hashicorp/terraform-provider-azurerm/issues/19303))
* `cognitive`: updating to API Version `2022-10-01` ([#19344](https://github.com/hashicorp/terraform-provider-azurerm/issues/19344))
* `springcloud`: updating to API Version `2022-09-01-preview` ([#19340](https://github.com/hashicorp/terraform-provider-azurerm/issues/19340))
* Data Source: `azurerm_mssql_managed_instance` - support for `customer_managed_key_id` attribute and user-assigned identity ([#18918](https://github.com/hashicorp/terraform-provider-azurerm/issues/18918))
* `azurerm_cognitive_account` - support for `dynamic_throttling_enabled` ([#19371](https://github.com/hashicorp/terraform-provider-azurerm/issues/19371))
* `azurerm_databricks_workspace` - support for `storage_account_identity` property in datasource ([#19336](https://github.com/hashicorp/terraform-provider-azurerm/issues/19336))
* `azurerm_mssql_managed_instance` - support for user-assigned identity ([#18918](https://github.com/hashicorp/terraform-provider-azurerm/issues/18918))
* `azurerm_postgresql_flexible_server` - support for `authentication` ([#19269](https://github.com/hashicorp/terraform-provider-azurerm/issues/19269))
* `azurerm_spring_cloud_app` - support for the `ingress_settings` block ([#19386](https://github.com/hashicorp/terraform-provider-azurerm/issues/19386))

BUG FIXES:

* `azurerm_application_insights` - validating/normalizing the `workspace_id` as a Workspace ID ([#19325](https://github.com/hashicorp/terraform-provider-azurerm/issues/19325))
* `azurerm_cdn_frontdoor_rule` - allow `cache_duration` to be `null` if `cache_behavior` is set to `HonorOrigin` ([#19378](https://github.com/hashicorp/terraform-provider-azurerm/issues/19378))
* `azurerm_monitor_alert_processing_rule_action_group` - `condition.x.monitor_condition` can be correctly specified alone ([#19338](https://github.com/hashicorp/terraform-provider-azurerm/issues/19338))
* `azurerm_monitor_alert_processing_rule_suppression` - `condition.x.monitor_condition` can be correctly specified alone ([#19338](https://github.com/hashicorp/terraform-provider-azurerm/issues/19338))
* `azurerm_mysql_flexible_server` - increase validation max value for the `iops` property ([#19419](https://github.com/hashicorp/terraform-provider-azurerm/issues/19419))
* `azurerm_servicebus_subscription_rule` - `correlation_filter` with empty attributes no longer crashes ([#19352](https://github.com/hashicorp/terraform-provider-azurerm/issues/19352))
* `azurerm_storage_account`  - fix crash in multichannel checking ([#19298](https://github.com/hashicorp/terraform-provider-azurerm/issues/19298))
* `azurerm_storage_account` - prevent both `blob_properties.0.versioning_enabled` and `is_hns_enabled` being set to true ([#19418](https://github.com/hashicorp/terraform-provider-azurerm/issues/19418))

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

* Data Source: `azurerm_app_configuration_keys`  - fixa crash when `label` is not set ([#19032](https://github.com/hashicorp/terraform-provider-azurerm/issues/19032))
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
* `azurerum_cdn_frontdoor_route`  - fixa panic on import ([#18824](https://github.com/hashicorp/terraform-provider-azurerm/issues/18824))
* `azurerm_eventhub_namespace` - ignore case for `network_rulesets.x.virtual_network_rule.x.subnet_id` ([#18818](https://github.com/hashicorp/terraform-provider-azurerm/issues/18818))
* `azurerm_firewall_policy_rule_collection_group` - limit the number of destination ports in a NAT rule to one ([#18766](https://github.com/hashicorp/terraform-provider-azurerm/issues/18766))
* Data Source: `azurerm_linux_function_app`  - fixmissing error on data source not found ([#18876](https://github.com/hashicorp/terraform-provider-azurerm/issues/18876))
* `azurerm_linux_function_app`  - fixan issue where `app_settings` would show a diff when setting `vnet_route_all_enabled` to true ([#18836](https://github.com/hashicorp/terraform-provider-azurerm/issues/18836))
* `azurerm_linux_function_app_slot`  - fixan issue where `app_settings` would show a diff when setting `vnet_route_all_enabled` to true ([#18836](https://github.com/hashicorp/terraform-provider-azurerm/issues/18836))
* `azurerm_linux_virtual_machine` - allow disabling secure boot when creating a virtual machine with disk encryption type `VMGuestStateOnly` ([#18749](https://github.com/hashicorp/terraform-provider-azurerm/issues/18749))
* `azurerm_linux_virtual_machine_scale_set` - allow disabling secure boot when creating a virtual machine scale set with disk encryption type `VMGuestStateOnly` ([#18749](https://github.com/hashicorp/terraform-provider-azurerm/issues/18749))
* `azurerm_network_security_group` - correct the casing of the `protocol` property ([#18799](https://github.com/hashicorp/terraform-provider-azurerm/issues/18799))
* `azurerm_network_security_rule` - correct the casing of the `protocol` property ([#18799](https://github.com/hashicorp/terraform-provider-azurerm/issues/18799))
* `azurerm_recovery_services_vault`  - fixissue where `soft_delete_enabled` is reset to the default value when the `identity` block is updated ([#18871](https://github.com/hashicorp/terraform-provider-azurerm/issues/18871))
* `azurerm_windows_virtual_machine` - allow disabling secure boot when creating a virtual machine with disk encryption type `VMGuestStateOnly` ([#18749](https://github.com/hashicorp/terraform-provider-azurerm/issues/18749))
* `azurerm_windows_function_app`  - fixan issue where `app_settings` would show a diff when setting `vnet_route_all_enabled` to true ([#18836](https://github.com/hashicorp/terraform-provider-azurerm/issues/18836))
* `azurerm_windows_function_app_slot`  - fixan issue where `app_settings` would show a diff when setting `vnet_route_all_enabled` to true ([#18836](https://github.com/hashicorp/terraform-provider-azurerm/issues/18836))
* `azurerm_windows_web_app`  - fixparsing of `docker_container_name` and `docker_container_registry` on read ([#18251](https://github.com/hashicorp/terraform-provider-azurerm/issues/18251))

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
* `azurerm_kubernetes_cluster_node_pool`  - fixa crash in expanding upgrade settings ([#18074](https://github.com/hashicorp/terraform-provider-azurerm/issues/18074))
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
* `azurerm_windows_function_app_slot`  - fixa bug with setting values for `WindowsFxString` ([#18014](https://github.com/hashicorp/terraform-provider-azurerm/issues/18014))
* `azurerm_linux_function_app` - correctly send `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_linux_function_app`  - fixcontent settings when `storage_uses_managed_identity` is set to `true` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_linux_function_app_slot` - correctly send `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_linux_function_app_slot`  - fixcontent settings when `storage_uses_managed_identity` is set to `true` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_windows_function_app` - correctly send `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_windows_function_app`  - fixcontent settings when `storage_uses_managed_identity` is set to `true` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_windows_function_app_slot` - correctly send `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))
* `azurerm_windows_function_app_slot`  - fixcontent settings when `storage_uses_managed_identity` is set to `true` ([#18035](https://github.com/hashicorp/terraform-provider-azurerm/issues/18035))

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
* Data Source: `azurerm_blueprint_definition`  - fix`version` property output ([#16299](https://github.com/hashicorp/terraform-provider-azurerm/issues/16299))

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

* `azurerm_linux_function_app`  - fixcasing bug with the `linux_fx_string` property for Node apps ([#17789](https://github.com/hashicorp/terraform-provider-azurerm/issues/17789))
* `azurerm_linux_function_app_slot`  - fixcasing bug with the `linux_fx_string` property for Node apps ([#17789](https://github.com/hashicorp/terraform-provider-azurerm/issues/17789))
* `azurerm_resource_group_template_deployment` - fixing a bug where the same Resource Provider defined in different casings would cause the API Version to not be identified ([#17707](https://github.com/hashicorp/terraform-provider-azurerm/issues/17707))

## 3.15.1 (July 25, 2022)

BUG FIXES: 

* `data.azurerm_servicebus_queue`  - fixa regression around `namespace_id` ([#17755](https://github.com/hashicorp/terraform-provider-azurerm/issues/17755))
* `azurerm_postgresql_aad_administrator`  - fixthe state migration ([#17732](https://github.com/hashicorp/terraform-provider-azurerm/issues/17732))
* `azurerm_postgresql_server`  - fixa regression around `id` ([#17755](https://github.com/hashicorp/terraform-provider-azurerm/issues/17755))

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
* `azurerm_linux_function_app_slot`  - fix`app_settings.WEBSITE_RUN_FROM_PACKAGE` handling from external sources ([#16641](https://github.com/hashicorp/terraform-provider-azurerm/issues/16641))
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

* `azurerm_datafactory_dataset_x`  - fix crash around `azure_blob_storage_location.0.dynamic_container_enabled` ([#16514](https://github.com/hashicorp/terraform-provider-azurerm/issues/16514))
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
* `azurerm_linux_function_app`  - fixa bug in updates to `app_settings` where settings could be lost ([#16442](https://github.com/hashicorp/terraform-provider-azurerm/issues/16442))
* `azurerm_linux_function_app_slot` -  this `ip_address` property is now correctly set into state when the `service_tag` property is specified ([#16426](https://github.com/hashicorp/terraform-provider-azurerm/issues/16426))
* `azurerm_linux_web_app` -  the `ip_address` property is correctly set into state when the `service_tag` property is specified ([#16426](https://github.com/hashicorp/terraform-provider-azurerm/issues/16426))
* `azurerm_linux_web_app`  - fixa potential crash when an empty `app_stack` block is used ([#16446](https://github.com/hashicorp/terraform-provider-azurerm/issues/16446))
* `azurerm_linux_web_app_slot` -  the `ip_address` property is now correctly set into state when the `service_tag` property is specified ([#16426](https://github.com/hashicorp/terraform-provider-azurerm/issues/16426))
* `azurerm_linux_web_app_slot`  - fixa potential crash when an empty `app_stack` block is used ([#16446](https://github.com/hashicorp/terraform-provider-azurerm/issues/16446))
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
* `azurerm_windows_function_app`  - fixa bug in updates to `app_settings` where settings could be lost ([#16442](https://github.com/hashicorp/terraform-provider-azurerm/issues/16442))
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
* `azurerm_windows_function_app`  - fixthe import check for Service Plan OS type ([#16164](https://github.com/hashicorp/terraform-provider-azurerm/issues/16164))
* `azurerm_linux_web_app_slot `  - fix`container_registry_managed_identity_client_id` property validation ([#16149](https://github.com/hashicorp/terraform-provider-azurerm/issues/16149))
* `azurerm_windows_web_app` - add support for `dotnetcore` in site metadata property `current_stack` ([#16129](https://github.com/hashicorp/terraform-provider-azurerm/issues/16129))
* `azurerm_windows_web_app`  - fixdocker `windowsFXVersion` when `docker_container_registry` is specified ([#16192](https://github.com/hashicorp/terraform-provider-azurerm/issues/16192))
* `azurerm_windows_web_app_slot` - add support for `dotnetcore` in site metadata property `current_stack` ([#16129](https://github.com/hashicorp/terraform-provider-azurerm/issues/16129))
* `azurerm_windows_web_app_slot`  - fixdocker `windowsFXVersion` when `docker_container_registry` is specified ([#16192](https://github.com/hashicorp/terraform-provider-azurerm/issues/16192))
* `azurerm_storage_data_lake_gen2_filesystem` - add support for `$superuser` in `group` and `owner` properties ([#16215](https://github.com/hashicorp/terraform-provider-azurerm/issues/16215))

## 3.0.2 (March 26, 2022)

BUG FIXES:

* `azurerm_cosmosdb_account` - prevent a panic when the API returns an empty list of read or write locations ([#16031](https://github.com/hashicorp/terraform-provider-azurerm/issues/16031))
* `azurerm_cdn_endpoint` - prevent a panic when there is an empty `country_codes` property ([#16066](https://github.com/hashicorp/terraform-provider-azurerm/issues/16066))
* `azurerm_key_vault`  - fixthe `authorizer was not an auth.CachedAuthorizer ` error ([#16078](https://github.com/hashicorp/terraform-provider-azurerm/issues/16078))
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
* `azurerm_local_network_gateway`  - fixfor `address_space` cannot be updated ([#15159](https://github.com/hashicorp/terraform-provider-azurerm/issues/15159))
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
