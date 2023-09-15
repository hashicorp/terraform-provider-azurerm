## 3.74.0 (Unreleased)

ENHANCEMENTS:

* dependencies - swap Virtual Machine Scale Set Ids to `hashicorp/go-azure-helpers/commonids` [GH-23272]

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
