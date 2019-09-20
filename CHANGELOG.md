## 1.35.0 (Unreleased)

FEATURES:

* **New Resource:** `azurerm_bot_channel_slack` [GH-4367]

IMPROVEMENTS:

* dependencies: upgrading `github.com/Azure/azure-sdk-for-go` to `v33.2.0` [GH-4334]
* kusto: updating to API version `2019-05-15` [GH-4376]

BUG FIXES:

* `azurerm_key_vault_certificate` - storing the certificate data as hex [GH-4335]
* `azurerm_public_ip` - ensuring that `public_ip_prefix_id` is read [GH-4344]

## 1.34.0 (September 18, 2019)

FEATURES:

* **New Data Source:** `azurerm_network_ddos_protection_plan` ([#4228](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4228))
* **New Data Source:** `azurerm_proximity_placement_group` ([#4020](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4020))
* **New Data Source:** `azurerm_servicebus_namespace_authorization_rule` ([#4294](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4294))
* **New Data Source:** `azurerm_sql_database` ([#4210](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4210))
* **New Data Source:** `azurerm_storage_account_blob_container_sas` ([#4195](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4195))
* **New Resource:** `azurerm_app_service_certificate` ([#4192](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4192))
* **New Resource:** `azurerm_app_service_source_control_token` ([#4214](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4214))
* **New Resource:** `azurerm_bot_channels_registration` ([#4245](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4245))
* **New Resource:** `azurerm_bot_connection` ([#4311](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4311))
* **New Resource:** `azurerm_frontdoor` ([#3933](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3933))
* **New Resource:** `azurerm_frontdoor_firewall_policy` ([#4125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4125))
* **New Resource:** `azurerm_kusto_cluster` ([#4129](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4129))
* **New Resource:** `azurerm_kusto_database` ([#4149](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4149))
* **New Resource:** `azurerm_marketplace_agreement` ([#4305](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4305))
* **New Resource:** `azurerm_private_dns_zone_virtual_network_link` ([#3789](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3789))
* **New Resource:** `azurerm_proximity_placement_group` ([#4020](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4020))
* **New Resource:** `azurerm_stream_analytics_output_servicebus_topic` ([#4164](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4164))
* **New Resource:** `azurerm_web_application_firewall_policy` ([#4119](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4119))

IMPROVEMENTS:

* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v32.5.0` ([#4166](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4166))
* dependencies: updating `github.com/Azure/go-autorest` to `v0.9.0` ([#4166](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4166))
* dependencies: updating `github.com/hashicorp/go-azure-helpers` to `v0.7.0` ([#4166](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4166))
* dependencies: updating `github.com/terraform-providers/terraform-provider-azuread` to `v0.6.0` ([#4166](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4166))
* dependencies: updating `github.com/hashicorp/terraform` to `v0.12.8` ([#4341](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4341))
* compute: updating the API Version to `2019-07-01` ([#4331](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4331))
* network: updating to API version `2019-06-01` ([#4291](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4291))
* network: reverting the locking changes from #3673 ([#3673](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3673))
* storage: caching the Resource Group Name / Account Key ([#4205](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4205))
* storage: switching to use SharedKey for authentication with Blobs/Containers rather than SharedKeyLite ([#4235](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4235))
* Data Source: `azurerm_storage_account` - gracefully degrading when there's a ReadOnly lock/the user doesn't have permissions to list the Keys for the storage account ([#4248](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4248))
* Data Source: `azurerm_storage_account_sas` - adding an `ISO8601` validator to the `start` and `end` dates ([#4064](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4064))
* Data Source: `azurerm_virtual_network` - support for the `location` property ([#4281](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4281))
* `azurerm_api_management` - support for multiple `additional_location` blocks ([#4175](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4175))
* `azurerm_application_gateway` - allowing `capacity` to be set to `32` ([#4189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4189))
* `azurerm_application_gateway` - support OWASP version `3.1` for the `rule_set_version` property ([#4263](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4263))
* `azurerm_application_gateway` - support for the `trusted_root _certificate` property ([#4206](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4206))
* `azurerm_app_service` - fixing a bug where the Application `logs` block would get reset when `app_settings` were configured ([#4243](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4243))
* `azurerm_app_service` - support for sending HTTP Logs to Blob Storage ([#4249](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4249))
* `azurerm_app_service` - the `ip_restriction.ip_address` property is now optional ([#4184](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4184))
* `azurerm_app_service_slot` - the `ip_restriction.ip_address` property is now optional ([#4184](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4184))
* `azurerm_availability_set` - support for the `proximity_placement_group_id` property ([#4020](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4020))
* `azurerm_cognitive_account` - supporting `CognitiveServices` as a `kind` ([#4209](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4209))
* `azurerm_container_registry` - support for configuring Virtual Network Rules to Subnets ([#4293](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4293))
* `azurerm_cosmosdb_account` - correctly validate `max_interval_in_seconds` & `max_staleness_prefix` for geo replicated accounts ([#4273](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4273))
* `azurerm_cosmosdb_account` - increase creation & deletion wait timeout to `3` hours ([#4271](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4271))
* `azurerm_cosmosdb_sql_container` - changing the `unique_key.paths` property now forces a new resource ([#4163](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4163))
* `azurerm_eventhub_namespace` - changing the `kafka_enabled` property now forces a new resource ([#4264](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4264))
* `azurerm_kubernetes_cluster` - support for configuring the `kube_dashboard` within the `addon_profile` block ([#4139](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4139))
* `azurerm_kubernetes_cluster` - prevent `pod_cidr` and azure `network_plugin` from being set at the same time causing a new resource to be created ([#4286](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4286))
* `azurerm_mariadb_server` - support for version `10.3` ([#4170](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4170))
* `azurerm_mariadb_server` - support for configuring `auto_grow` ([#4302](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4302))
* `azurerm_managed_disk` - add support for the Ultra SSD `disk_iops_read_write` & `disk_mbps_read_write` properties ([#4102](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4102))
* `azurerm_mysql_server` - support for configuring `auto_grow` ([#4303](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4303))
* `azurerm_private_dns_zone` - polling until the dns zone is marked as fully provisioned ([#4307](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4307))
* `azurerm_postgresql_server` - support for configuring `auto_grow` ([#4220](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4220))
* `azurerm_resource_group` - the `name` field can now be up to 90 characters ([#4233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4233))
* `azurerm_role_assignment` - add `principal_type` and `skip_service_principal_aad_check` properties ([#4168](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4168))
* `azurerm_storage_account` - gracefully degrading when there's a ReadOnly lock/the user doesn't have permissions to list the Keys for the storage account ([#4248](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4248))
* `azurerm_storage_blob` - switching over to use the new Storage SDK ([#4179](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4179))
* `azurerm_storage_blob` - support for Append Blobs ([#4238](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4238))
* `azurerm_storage_blob` - support for configuring the `access_tier` ([#4238](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4238))
* `azurerm_storage_blob` - support for specifying Block Blob content via `source_content` ([#4238](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4238))
* `azurerm_storage_blob` - the `type` field is now Required, since it had to be set anyway ([#4238](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4238))
* `azurerm_storage_share_directory` - support for upper-case characters in the `name` field ([#4178](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4178))
* `azurerm_storage_table` - using the correct storage account name when checking for the presence of an existing storage table ([#4234](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4234))
* `azurerm_stream_analytics_job` - the field `data_locale` is now optional ([#4190](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4190))
* `azurerm_stream_analytics_job` - the field `events_late_arrival_max_delay_in_seconds` is now optional ([#4190](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4190))
* `azurerm_stream_analytics_job` - the field `events_out_of_order_policy` is now optional ([#4190](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4190))
* `azurerm_stream_analytics_job` - the field `output_error_policy` is now optional ([#4190](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4190))
* `azurerm_subnet` - support for the actions `Microsoft.Network/virtualNetworks/subnets/join/action` and `Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action` ([#4137](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4137))
* `azurerm_virtual_machine` - support for `UltraSSD_LRS` managed disks ([#3860](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3860))
* `azurerm_virtual_machine` - support for the `proximity_placement_group_id` property ([#4020](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4020))
* `azurerm_virtual_machine_scale_set` - support for the `proximity_placement_group_id` property ([#4020](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4020))

BUG FIXES:

* `azurerm_app_service` - will no longer panic from when an access restriction rule involves a virtual network ([#4184](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4184))
* `azurerm_app_service_slot` - will no longer panic from when an access restriction rule involves a virtual network ([#4184](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4184))
* `azurerm_app_service_plan` and `azurerm_app_service_slot` crash fixes ([#4184](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4184))
* `azurerm_container_group` - make `storage_account_key` field in `volume` block sensitive ([#4201](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4201))
* `azurerm_key_vault_certificate` - prevented a panic caused by an empty element in `extended_key_usage` ([#4272](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4272))
* `azurerm_log_analytics_linked_service` - will no longer panic if no items are passed into the property `linked_service_properties` ([#4142](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4142))
* `azurerm_log_analytics_workspace_linked_service` - will no longer panic if no items are passed into the property `linked_service_properties` ([#4152](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4152))
* `azurerm_network_interface` - changing the `ip_configuration` property to no longer force new resource ([#4155](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4155))
* `azurerm_virtual_network_peering` - prevent nil object from being read ([#4180](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4180))

## 1.33.1 (August 27, 2019)

* networking: reducing the number of locks to avoid deadlock when creating 3 or more subnets with Network Security Group/Route Table Associations ([#3673](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3673))

## 1.33.0 (August 22, 2019)

FEATURES:

* **New Data Source:** `azurerm_dev_test_virtual_network` ([#3746](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3746))
* **New Resource:** `azurerm_cosmosdb_sql_container` ([#3871](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3871))
* **New Resource:** `azurerm_container_registry_webhook` ([#4112](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4112))
* **New Resource:** `azurerm_dev_test_lab_schedule` ([#3554](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3554))
* **New Resource:** `azurerm_mariadb_virtual_network_rule` ([#4048](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4048))
* **New Resource:** `azurerm_mariadb_configuration` ([#4060](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4060))
* **New Resource:** `azurerm_private_dns_cname_record` ([#4028](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4028))
* **New Resource:** `azurerm_recovery_services_fabric` ([#4003](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4003))
* **New Resource:** `azurerm_recovery_services_protection_container` ([#4003](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4003))
* **New Resource:** `azurerm_recovery_services_replication_policy` ([#4003](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4003))
* **New Resource:** `azurerm_recovery_services_protection_container_mapping` ([#4003](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4003))
* **New Resource:** `azurerm_recovery_network_mapping` ([#4003](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4003))
* **New Resource:** `azurerm_recovery_replicated_vm`  ([#4003](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4003))
* **New Resource:** `azurerm_sql_failover_group` ([#3901](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3901))
* **New Resource:** `azurerm_virtual_wan` ([#4089](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4089))

IMPROVEMENTS:

* all resources: increasing the maximum number of tags from `15` to `50` ([#4071](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4071))
* dependencies: upgrading `github.com/tombuildsstuff/giovanni` to `v0.3.2` ([#4122](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4122))
* dependencies: upgrading the `authorization` SDK to `2018-09-01` ([#4063](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4063))
* dependencies: upgrading `github.com/hashicorp/terraform` to `0.12.6` ([#4041](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4041))
* internal: removing a duplicate Date/Time from the debug logs ([#4024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4024))
* Data Source `azurerm_dns_zone`: deprecating the `zone_type` field ([#4033](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4033))
* `azurerm_app_service` - `filesystem` logging can now be set. ([#4025](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4025))
* `azurerm_batch_pool` - Support for Container Registry configurations ([#4072](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4072))
* `azurerm_container_group` - support for attaching to a (Private) Virtual Network ([#3716](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3716))
* `azurerm_container_group` - `log_type` can now be an empty string ([#4013](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4013))
* `azurerm_cognitive_account` - Adding 'QnAMaker' as Kind ([#4126](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4126))
* `azurerm_dns_zone` - deprecating the `zone_type` field ([#4033](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4033))
* `azurerm_function_app` - support for cors ([#3949](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3949))
* `azurerm_function_app` - support for the `virtual_network_name` property ([#4078](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4078))
* `azurerm_iot_dps` - add support for the `linked_hub` property ([#3922](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3922))
* `azurerm_kubernetes_cluster` - support for the `enable_pod_security_policy` property ([#4098](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4098))
* `azurerm_monitor_diagnostic_setting` - support for `log_analytics_destination_type` ([#3987](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3987))
* `azurerm_role_assignment` - now supports management groups ([#4063](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4063))
* `azurerm_storage_account` - requesting an access token using the ARM Authorizer ([#4099](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4099))
* `azurerm_storage_account` - support for `BlockBlobStorage` ([#4131](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4131))
* `azurerm_subnet` - support for the Service Endpoints `Microsoft.BareMetal/AzureVMware`, `Microsoft.BareMetal/CrayServers`, `Microsoft.Databricks/workspaces` and `Microsoft.Web/hostingEnvironments` ([#4115](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4115))
* `azurerm_traffic_manager_profile` - support for the `interval_in_seconds`, `timeout_in_seconds`, and `tolerated_number_of_failures` properties ([#3473](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3473))
* `azurerm_user_assigned_identity` - the `name` field can now be up to 128 characters ([#4094](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4094))

BUG FIXES: 

* `azurerm_app_service_plan` - workaround for missing error on 404 ([#3990](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3990))
* `azurerm_batch_certificate` - the `thumbprint_algorithm` property is now case insensitive ([#3977](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3977))
* `azurerm_notification_hub_authorization_rule - fixing an issue when creating multiple authorization rules at the same time ([#4087](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4087))
* `azurerm_postgresql_server` - removal of unsupported version `10.2` ([#3915](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3915))
* `azurerm_role_definition` - enture `role_definition_id` is correctly set if left empty during creation ([#3913](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3913))
* `azurerm_storage_account` - making `default_action` within the `network_rules` block required ([#4037](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4037))
* `azurerm_storage_account` - making the `network_rules` block computed ([#4037](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4037))
* `azurerm_storage_queue` - switching to using SharedKey for authentication ([#4122](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4122))
* `azurerm_storage_share` - allow up to 100TB for the `quota` property ([#4054](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4054))
* `azurerm_storage_share_directory` - handling the share being eventually consistent ([#4122](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4122))
* `azurerm_storage_share_directory` - allowing nested directories ([#4122](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4122))

## 1.32.1 (July 31, 2019)

BUG FIXES: 

* `azurerm_application_gateway` fix an index out of range crash ([#3966](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3966))
* `azurerm_api_management_backend` - ensuring a nil `certificates` object is sent to the API instead of an empty one ([#3931](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3931))
* `azurerm_api_managment_product` - additional validation for `approval_required` ([#3945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3945))
* `azurerm_network_ddos_protection_plan` - correctly decodes the resource ID on read/delete ([#3975](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3975))
* `azurerm_dev_test_virtual_network` - generate subnet IDs in the correct format ([#3717](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3717))
* `azurerm_iot_dps` fixed deletion issue when using a service principal ([#3973](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3973))
* `azurerm_kubernetes_cluster` - the `load_balancer_sku` property is now case insensitive ([#3958](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3958))
* `azurerm_postgresql_server` - add missing support for version `11.0` ([#3970](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3970))
* `azurerm_storage_*` - prevent multiple panics when a storage account/resource group cannot be found ([#3986](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3986))
* `azurerm_storage_account` - fix `enable_advanced_threat_protection` create/read for unsupported regions ([#3947](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3947))
* `azurerm_storage_table` - now migrates older versions of the resource id to the new format ([#3932](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3932))
* `azurerm_virtual_machine_scale_set` - the `ssh_keys` property of the `os_profile_linux_config` block now recognizes updates ([#3837](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3837))
* `azurerm_virtual_machine_scale_set` - changes made to the `network_profile` property should now be correctly reflected during updates ([#3821](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3821))


## 1.32.0 (July 24, 2019)

FEATURES:

* **New Data Source:** `azurerm_maps_account` ([#3698](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3698))
* **New Data Source:** `azurerm_mssql_elasticpool` ([#3824](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3824))
* **New Resource:** `azurerm_analysis_services_server` ([#3721](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3721))
* **New Resource:** `azurerm_api_management_backend` ([#3676](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3676))
* **New Resource:** `azurerm_batch_application` ([#3825](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3825))
* **New Resource:** `azurerm_maps_account` ([#3698](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3698))
* **New Resource:** `azurerm_private_dns_zone_a_record` ([#3849](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3849))
* **New Resource:** `azurerm_storage_table_entity` ([#3831](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3831))
* **New Resource:** `azurerm_storage_share_directory` ([#3802](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3802))

IMPROVEMENTS:

* dependencies: upgrading to `v31.0.0` of `github.com/Azure/azure-sdk-for-go` ([#3786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3786))
* dependencies: upgrading to `v0.5.0` of `github.com/hashicorp/go-azure-helpers` ([#3850](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3850))
* dependencies: upgrading the `containerservice` SDK to `2019-02-01` ([#3787](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3787))
* dependencies: upgrading the `subscription` SDK to `2018-06-01` ([#3811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3811))
* authentication: showing a more helpful error when attempting to use the Azure CLI authentication when logged in as a Service Principal ([#3850](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3850))
* Data Source `azurerm_function_app` - support for `auth_settings` ([#3893](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3893))
* Data Source `azurerm_subscription` - support the `tenant_id` property ([#3811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3811))
* `azurerm_app_service` - support for backups ([#3804](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3804))
* `azurerm_app_service` - support for storage mounts ([#3792](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3792))
* `azurerm_app_service` - support for user assigned identities ([#3637](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3637))
* `azurerm_app_service_slot` - support for `auth_settings` ([#3897](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3897))
* `azurerm_app_service_slot` - support for user assigned identities ([#3637](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3637))
* `azurerm_application_gateway` - Support for Managed Identities ([#3648](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3648))
* `azurerm_batch_pool` - support for custom images with the `storage_image_reference` property ([#3530](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3530))
* `azurerm_batch_account` - expose required properties for when `pool_allocation_mode` is `UserSubscription` ([#3535](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3535))
* `azurerm_cognitive_account` - add support for `CustomVision.Training` and `CustomVision.Prediction` to the `kind` property ([#3817](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3817))
* `azurerm_container_registry` - support for `network_rule_set` property ([#3194](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3194))
* `azurerm_cosmosdb_account` - validate `max_interval_in_seconds` and `max_staleness_prefix` correctly when using more then 1 geo_location ([#3906](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3906))
* `azurerm_function_app` - support for `auth_settings` ([#3893](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3893))
* `azurerm_iothub` - support for the `file_upload` property ([#3735](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3735))
* `azurerm_kubernetes_cluster` - support for auto scaling ([#3361](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3361))
* `azurerm_kubernetes_cluster` - support for `custom_resource_group_name` ([#3785](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3785))
* `azurerm_kubernetes_cluster` - support for the `node_taints` property ([#3787](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3787))
* `azurerm_kubernetes_cluster`  - support for the `windows_profile` property ([#3519](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3519))
* `kubernetes_cluster` - support for specifying the `load_balancer_sku` property ([#3890](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3890))
* `azurerm_recovery_services_protected_vm` - changing `backup_policy_id` no longer forces a new resource ([#3822](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3822))
* `azurerm_security_center_contact` - the `phone` property is now optional ([#3761](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3761))
* `azurerm_storage_account` - the `account_kind` property now supports `FileStorage` ([#3750](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3750))
* `azurerm_storage_account` - support for the `enable_advanced_threat_protection` property ([#3782](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3782))
* `azurerm_storage_account` - support for `queue_properties` ([#3859](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3859))
* `azurerm_storage_blob` - making `metadata` a computed field ([#3842](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3842))
* `azurerm_storage_container` - switching to use github.com/tombuildsstuff/giovanni ([#3857](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3857))
* `azurerm_storage_container` - adding support for `metadata` ([#3857](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3857))
* `azurerm_storage_container` - can now create containers with the name `$web` ([#3896](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3896))
* `azurerm_storage_queue` - switching to use github.com/tombuildsstuff/giovanni ([#3832](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3832))
* `azurerm_storage_share` - switching to use github.com/tombuildsstuff/giovanni ([#3828](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3828))
* `azurerm_storage_share` - support for configuring ACL's ([#3830](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3830))
* `azurerm_storage_share` - support for configuring MetaData ([#3830](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3830))
* `azurerm_storage_table` - switching to use github.com/tombuildsstuff/giovanni ([#3834](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3834))
* `azurerm_storage_table` - support for configuring ACL's ([#3847](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3847))
* `azurerm_traffic_manager_endpoint` - supper for `custom_header` and `subnet` properties ([#3655](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3655))
* `azurerm_virtual_machine` - switching over to use the github.com/tombuildsstuff/giovanni Storage SDK ([#3838](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3838))
* `azurerm_virtual_machine` - looking up the data disks attached to the Virtual Machine when optionally deleting them upon deletion rather than parsing them from the config ([#3838](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3838))
* `azurerm_virtual_machine_scale_set` - prevent `public_ip_address_configuration` from being lost during update ([#3767](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3767))

BUG FIXES:

* `azurerm_image` - prevent crash when using `data_disk` ([#3797](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3797))
* `azurerm_role_assignment` - now correctly uses `scope` when looking up the role definition by name ([#3768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3768))

## 1.31.0 (June 28, 2019)

FEATURES:

* increase the default timeout to `3 hours` ([#3737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3737))
* **New Resource:** `azurerm_iot_dps` ([#3618](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3618))
* **New Resource:** `azurerm_iot_dps_certificate` ([#3567](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3645))
* **New Resource:** `azurerm_mariadb_firewall_rule` ([#3720](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3720))
* **New Resource:** `azurerm_private_dns_zone` ([#3718](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3718))
* **New Resource:** `azurerm_stream_analytics_output_mssql` ([#3567](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3567))

IMPROVEMENTS:

* Data Source `azurerm_key_vault` - deprecated `sku` in favour of `sku_name` ([#3119](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3119))
* `azurerm_app_service` - support for shipping the application logs to blob storage ([#3520](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3520))
* `azurerm_app_service_plan` - prevent a panic during import ([#3657](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3657))
* `azurerm_app_service_slot` - updating `identity` no longer forces a new resource ([#3702](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3702))
* `azurerm_automation_account` - deprecated `sku` in favour of `sku_name` ([#3119](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3119))
* `azurerm_key_vault` - deprecated `sku` in favour of `sku_name` ([#3119](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3119))
* `azurerm_key_vault_key` - add support for Elliptic Curve based keys ([#1814](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1814))
* `azurerm_traffic_manager_profile` - `ttl` can now be 1 second ([#3632](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3632))
* `azurerm_eventgrid_event_subscription` - now retrieves the full URL for event webhooks ([#3630](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3630))
* `azurerm_lb` - support for the `public_ip_prefix_id` property ([#3675](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3675))
* `azurerm_mysql_server` - add validation to the `name` property ([#3695](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3695))
* `azurerm_notification_hub_namespace` - deprecated `sku` in favour of `sku_name` ([#3119](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3119))
* `azurerm_redis_firewall_rule` - no longer fails with multiple rules ([#3731](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3731))
* `azurerm_relay_namespace` - deprecated `sku` in favour of `sku_name` ([#3119](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3119))
* `azurerm_service_fabric_cluster` - `tenant_id`, `cluster_application_id`, and `client_application_id` are now updateable ([#3654](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3654))
* `azurerm_service_fabric_cluster` - ability to set `certificate_common_names` ([#3652](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3652))
* `azurerm_storage_account` - ability to set `default_action` oi the `network_rules` block ([#3255](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3255))

BUG FIXES:

* `azurerm_cosmosdb_account` - will ignore `500` responses from `documentdb.DatabaseAccountsClient#CheckNameExists` requests to work around a broken API ([#3747](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3747))

## 1.30.1 (June 07, 2019)

BUG FIXES:

* Ensuring the authorization header is set for calls to the User Assigned Identity API's ([#3613](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3613))

## 1.30.0 (June 07, 2019)

FEATURES:

* **New Data Source:** `azurerm_redis_cache` ([#3481](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3481))
* **New Data Source:** `azurerm_sql_server` ([#3513](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3513))
* **New Data Source:** `azurerm_virtual_network_gateway_connection` ([#3571](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3571))

IMPROVEMENTS:

* dependencies: upgrading to Go 1.12 ([#3525](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3525))
* dependencies: upgrading the `storage` SDK to `2019-04-01` ([#3578](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3578))
* Data Source `azurerm_app_service` - support windows containers ([#3566](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3566))
* Data Source `azurerm_app_service_plan` - support windows containers ([#3566](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3566))
* `azurerm_api_management` - rename `disable_triple_des_chipers` to `disable_triple_des_ciphers` ([#3539](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3539))
* `azurerm_application_gateway` - support for the value `General` in the `rule_group_name` field within the `disabled_rule_group` block ([#3533](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3533))
* `azurerm_app_service` - support for windows containers ([#3566](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3566))
* `azurerm_app_service_plan` - support for the `maximum_elastic_worker_count` property ([#3547](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3547))
* `azurerm_managed_disk` - support for the `create_option` of `Restore` ([#3598](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3598))
* `azurerm_app_service_plan` - support for windows containers ([#3566](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3566))


## 1.29.0 (May 25, 2019)

FEATURES:

* **New Resource:** `azurerm_application_insights_web_test` ([#3331](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3331))

IMPROVEMENTS:

* dependencies: upgrading to `v0.12.0` of `github.com/hashicorp/terraform` ([#3417](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3417))
* sdk: configuring the Correlation Request ID ([#3253](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3253))
* `azurerm_application_gateway` - support for rewrite rules ([#3423](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3423))
* `azurerm_application_gateway` - support for `ssl_policy` blocks and deprecating `disabled_ssl_protocols` ([#3360](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3360))
* `azurerm_app_service` - support for configuring authentication settings ([#2831](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2831))
* `azurerm_kubernetes_cluster` - updating the casing on the `SubnetName` field to match a change in the AKS API ([#3484](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3484))
* `azurerm_kubernetes_cluster` - support for multiple agent pools ([#3491](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3491))

BUG FIXES:

* Data Source `azurerm_virtual_network`: add `network_space` property to match resource while deprecating `network_spaces` ([#3494](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3494))
* `azurerm_automation_module` - now polls to wait until the module's finished provisioning ([#3482](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3482))
* `azurerm_api_management_api` - correct validation to allow empty and strings 400 characters long ([#3475](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3475))
* `azurerm_dev_test_virtual_network` - correctly manages `subnets` on the initial creation ([#3501](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3501))
* `azurerm_express_route_circuit` - no longer removes circuit subresources on update ([#3496](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3496))
* `azurerm_role_assignment` - making the `role_definition_name` field case-insensitive ([#3499](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3499))

## 1.28.0 (May 17, 2019)

FEATURES:

* **New Data Source:** `azurerm_automation_variable_bool` ([#3310](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3310))
* **New Data Source:** `azurerm_automation_variable_datetime` ([#3310](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3310))
* **New Data Source:** `azurerm_automation_variable_int` ([#3310](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3310))
* **New Data Source:** `azurerm_automation_variable_string` ([#3310](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3310))
* **New Data Source:** `zurerm_kubernetes_service_versions` ([#3382](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3382))
* **New Data Source:** `azurerm_user_assigned_identity` ([#3343](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3343))
* **New Resource:** `azurerm_automation_variable_bool` ([#3310](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3310))
* **New Resource:** `azurerm_automation_variable_datetime` ([#3310](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3310))
* **New Resource:** `azurerm_automation_variable_int` ([#3310](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3310))
* **New Resource:** `azurerm_automation_variable_string` ([#3310](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3310))
* **New Resource:** `azurerm_api_management_api_operation_policy` ([#3374](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3374))
* **New Resource:** `azurerm_api_management_api_policy` ([#3367](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3367))
* **New Resource:** `azurerm_api_management_product_policy` ([#3325](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3325))
* **New Resource:** `azurerm_api_management_schema` ([#3357](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3357))
* **New Resource:** `azurerm_cosmosdb_table` ([#3442](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3442))
* **New Resource:** `azurerm_cosmosdb_cassandra_keyspace` ([#3442](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3442))
* **New Resource:** `azurerm_cosmosdb_mongo_collection` ([#3459](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3459))
* **New Resource:** `azurerm_cosmosdb_mongo_database` ([#3442](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3442))
* **New Resource:** `azurerm_cosmosdb_sql_database` ([#3442](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3442))
* **New Resource:** `azurerm_firewall_nat_rule_collection` ([#3218](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3218))
* **New Resource:** `azurerm_data_factory_linked_service_data_lake_storage_gen2` ([#3425](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3425))
* **New Resource:** `azurerm_network_profile` ([#2636](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2636))

IMPROVEMENTS:

* dependencies: updating `github.com/Azure/azure-sdk-for-go` to v29.0.0 ([#3335](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3335))
* Data Source `azurerm_kubernetes_cluster` - exposing the `type` field within the `agent_pool_profile ` block ([#3424](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3424))
* `azurerm_application_gateway` - support for the `autoscale_configuration` property ([#3353](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3353))
* `azurerm_application_gateway` added validation to ensure `redirect_configuration_name` must not be set if either `backend_address_pool_name` or `backend_http_settings_name` is set ([#3340](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3340))
* `azurerm_application_gateway` - support for `affinity_cookie_name` ([#3434](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3434))
* `azurerm_application_gateway` - support for `disabled_rule_groups` ([#3394](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3394))
* `azurerm_app_service_slot` - exporting the `site_credential` block ([#3444](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3444))
* `azurerm_batch_pool` support for the `container_configuration` property ([#3311](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3311))
* `azurerm_kubernetes_cluster` - support for the `api_server_authorized_ip_ranges` property ([#3262](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3262))
* `azurerm_kubernetes_cluster` - support for setting `type` within the `agent_pool_profile` block (Agent Pools via Virtual Machine Scale Sets) ([#3424](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3424))
* `azurerm_redis_cache` - support for disabling authentication ([#3389](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3389))
* `azurerm_redis_cache` - make the `redis_configuration` block optional ([#3397](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3397))
* `azurerm_sql_database` - support for the `read_scale` property ([#3377](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3377))
* `azurerm_stream_analytics_job` - `tags` can now be set on the property ([#3329](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3329))
* `azurerm_virtual_network_peering` - retrying provisioning the peering of the virtual network ([#3392](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3392))
* `azurerm_virtual_machine_scale_set` - support for the `provision_after_extensions` property to chain multiple extensions togeather ([#2937](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2937))

BUG FIXES:

* Data Source: `azurerm_api_management` - correctly returning the hostname `portal` and `proxy` values ([#3385](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3385))
* `azurerm_application_gateway` - will no longer prevent `default_backend_address_pool_name` and `redirect_configuration_name` from being set at the same time ([#3286](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3286))
* `azurerm_application_gateway` prevent a potential panic in backend and probe validation ([#3438](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3438))
* `azurerm_eventhub` - decrease minimum `partition_count` to correct value of `1` ([#3439](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3439))
* `azurerm_eventhub_namespace` - decrease maximum `maximum_throughput_units` to correct value of `20` ([#3440](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3440))
* `azurerm_firewall` - ensuring that the value for `subnet_id` within the `ip_configuration` block has the name `AzureFirewallSubnet` ([#3406](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3406))
* `azurerm_managed_disk` - can now actually create `UltraSSD_LRS` disks ([#3453](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3453))
* `azurerm_redis_configuration` - correctly display http errors encoutered during creation ([#3397](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3397))
* `azurerm_sql_database` - making the `collation` field case insensitive to work around a bug in the API ([#3137](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3137))
* `azurerm_stream_analytics_output_eventhub` will now correctly set `format` for JSON output ([#3318](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3318))
* `azurerm_app_service_plan` - supports `elastic` for the sku tier ([#3402](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3402))
* `azurerm_application_gateway` - supports `disabled_rule_group` for waf configurations ([#3394](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3394))
* `azurerm_application_gateway` - supports `exclusion` for waf configurations ([#3407](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3407))
* `azurerm_application_gateway` - supports updating a `gateway_ip_configuration.x.subnet_id` ([#3437](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3437))

## 1.27.1 (April 26, 2019)

BUG FIXES:

* provider will now only register available resource providers ([#3313](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3313))

## 1.27.0 (April 26, 2019)

NOTES:

* This release includes a Terraform SDK upgrade with compatibility for Terraform v0.12. The provider remains backwards compatible with Terraform v0.11 and there should not be any significant behavioural changes. ([#2968](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2968))

## 1.26.0 (April 25, 2019)

IMPROVEMENTS:

* `azurerm_app_service` - support for Java 11 ([#3270](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3270))
* `azurerm_app_service_slot` - support for Java 11 ([#3270](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3270))
* `azurerm_container_group` - support for the `identity` block ([#3243](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3243))

BUG FIXES:

* provider will work through proxies again ([#3301](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3301))

## 1.25.0 (April 17, 2019)

FEATURES:

* **New Data Source:** `azurerm_batch_certificate` ([#3097](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3097))
* **New Data Source:** `azurerm_express_route_circuit` ([#3158](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3158))
* **New Data Source:** `azurerm_firewall` ([#3235](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3235))
* **New Data Source:** `azurerm_hdinsight_cluster` ([#3196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3196))
* **New Data Source:** `azurerm_stream_analytics_job` ([#3227](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3227))
* **New Resource:** `azurerm_batch_certificate` ([#3097](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3097))
* **New Resource:** `azurerm_data_factory` ([#3159](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3159))
* **New Resource:** `azurerm_data_factory_dataset_mysql` ([#3267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3267))
* **New Resource:** `azurerm_data_factory_dataset_postgresql` ([#3267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3267))
* **New Resource:** `azurerm_data_factory_dataset_sql_server_table` ([#3236](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3236))
* **New Resource:** `azurerm_data_factory_linked_service_sql_server` ([#3205](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3205))
* **New Resource:** `azurerm_data_factory_linked_service_mysql` ([#3265](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3265))
* **New Resource:** `azurerm_data_factory_linked_service_postgresql` ([#3266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3266))
* **New Resource:** `azurerm_data_factory_pipeline` ([#3244](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3244))
* **New Resource:** `azurerm_hdinsight_kafka_cluster` ([#3196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3196))
* **New Resource:** `azurerm_hdinsight_kbase_cluster` ([#3196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3196))
* **New Resource:** `azurerm_hdinsight_hadoop_cluster` ([#3196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3196))
* **New Resource:** `azurerm_hdinsight_interactive_query_cluster` ([#3196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3196))
* **New Resource:** `azurerm_hdinsight_ml_services_cluster` ([#3196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3196))
* **New Resource:** `azurerm_hdinsight_rserver_cluster` ([#3196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3196))
* **New Resource:** `azurerm_hdinsight_spark_cluster` ([#3196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3196))
* **New Resource:** `azurerm_hdinsight_storm_cluster` ([#3196](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3196))
* **New Resource:** `azurerm_iothub_shared_access_policy` ([#3009](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3009))
* **New Resource:** `azurerm_public_ip_prefix` ([#3139](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3139))
* **New Resource:** `azurerm_stream_analytics_job` ([#3227](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3227))
* **New Resource:** `azurerm_stream_analytics_function_javascript_udf` ([#3249](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3249))
* **New Resource:** `azurerm_stream_analytics_stream_input_blob` ([#3250](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3250))
* **New Resource:** `azurerm_stream_analytics_stream_input_eventhub` ([#3250](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3250))
* **New Resource:** `azurerm_stream_analytics_stream_input_iothub` ([#3250](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3250))
* **New Resource:** `azurerm_stream_analytics_output_blob` ([#3250](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3250))
* **New Resource:** `azurerm_stream_analytics_output_eventhub` ([#3250](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3250))
* **New Resource:** `azurerm_stream_analytics_output_servicebus_queue` ([#3250](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3250))

IMPROVEMENTS:


* dependencies: updating `github.com/Azure/azure-sdk-for-go` to v26.7.0 ([#3126](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3126))
* dependencies: updating `github.com/Azure/go-autorest` to v11.7.0 ([#3126](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3126))
* dependencies: updating `github.com/hashicorp/terraform` to `44702fa6c163` ([#3181](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3181))
* Data Source: `azurerm_batch_pool` - adding the `resource_file` block to the `start_task` block ([#3192](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3192))
* Data Source: `azurerm_subnet` - exposing the `service_endpoint` field ([#3184](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3184))
* `azurerm_batch_pool` - adding the `resource_file` block to the `start_task` block ([#3192](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3192))
* `azurerm_container_group` - support for specifying `liveness_probe` and `readiness_probe` blocks ([#3118](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3118))
* `azurerm_key_vault_access_policy` - support for setting `storage_permissions` ([#3153](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3153))
* `azurerm_kubernetes_cluster` - `network_policy` now supports `azure` ([#3213](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3213))
* `azurerm_iothub` - support for configuring `ip_filter_rule` ([#3173](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3173))
* `azurerm_public_ip` - support for attaching a `azurerm_public_ip_prefix` ([#3139](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3139))
* `azurerm_redis_cache` - support for setting `aof_backup_enabled`, `aof_storage_connection_string_0` and `aof_storage_connection_string_1` ([#3155](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3155))
* `azurerm_storage_blob` - support for the `metadata` property ([#3206](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3206))
* `azurerm_traffic_manager_profile` - support the `MultiValue` and `Weighted` values for the `traffic_routing_method` property ([#3207](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3207))
* `azurerm_virtual_network_gateway` - support for the `VpnGw1AZ`, `VpnGw2AZ`, and `VpnGw3AZ` SKU's ([#3171](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3171))

BUG FIXES:

* dependencies: downgrading the Security API to `2017-08-01-preview` to work around a breaking API change ([#3269](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3269))
* `azurerm_app_service` - removing Computed from the `use_32_bit_worker_process` property in the `site_config` block ([#3219](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3219))
* `azurerm_app_service_slot` - removing Computed from the `use_32_bit_worker_process` property in the `site_config` block ([#3219](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3219))
* `azurerm_batch_account` - temporarily treating the Resource Group Name as case insensitive to work around an API bug ([#3260](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3260))
* `azurerm_batch_pool` - temporarily treating the Resource Group Name as case insensitive to work around an API bug ([#3260](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3260))
* `azurerm_app_service` - ensuring deleted App Services are detected correctly ([#3198](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3198))
* `azurerm_function_app` - ensuring deleted Function Apps are detected correctly ([#3198](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3198))
* `azurerm_virtual_machine` - adding validation for the `identity_ids` field ([#3183](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3183))

## 1.24.0 (April 03, 2019)

UPGRADE NOTES:

* `azurerm_kubernetes_cluster` - `ssh_key` is now limited to a single element to reflect what the API expects ([#3099](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3099))

FEATURES:

* **New Data Source:** `azurerm_api_management_api` ([#3010](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3010))
* **New Resource:** `azurerm_api_management_api` ([#3010](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3010))
* **New Resource:** `azurerm_api_management_api_operation` ([#3121](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3121))
* **New Resource:** `azurerm_api_management_api_version_set` ([#3073](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3073))
* **New Resource:** `azurerm_api_management_authorization_server` ([#3123](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3123))
* **New Resource:** `azurerm_api_management_certificate` ([#3141](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3141))
* **New Resource:** `azurerm_api_management_logger` ([#2994](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2994))
* **New Resource:** `azurerm_api_management_openid_connect_provider` ([#3143](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3143))
* **New Resource:** `azurerm_api_management_product_api` ([#3066](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3066))
* **New Resource:** `azurerm_api_management_subscription` ([#3103](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3103))

IMPROVEMENTS:

* Data Source: `azurerm_app_service` - exporting the `cors` headers ([#2870](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2870))
* Data Source: `azurerm_storage_account` - exposing the Hierarchical Namespace state ([#3032](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3032))
* `azurerm_api_management` - support for `sign_in`, `sign_up` and `policy` blocks ([#3151](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3151))
* `azurerm_app_service` - support for migrating between App Service Plans ([#3048](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3048))
* `azurerm_app_service` - support for additional types for the `scm_type` field in the `site_config` block ([#3019](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3019))
* `azurerm_app_service` - support for specifying `cors` headers ([#2870](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2870))
* `azurerm_app_service_slot` - support for specifying `cors` headers ([#2870](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2870))
* `azurerm_app_service_slot` - support for additional types for the `scm_type` field in the `site_config` block ([#3019](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3019))
* `azurerm_application_gateway` - support for WAF configuration properties `request_body_check` and `max_request_body_size_kb` ([#3093](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3093))
* `azurerm_application_gateway` - support for the `hostname` property ([#2990](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2990))
* `azurerm_application_gateway` - support for redirect rules ([#2908](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2908))
* `azurerm_application_gateway` - support for `zones` ([#3144](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3144))
* `azurerm_batch_account` - now exports the `primary_access_key`, `secondary_access_key`, and `account_endpoint` properties ([#3071](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3071))
* `azurerm_container_group` - support for attaching GPU's ([#3053](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3053))
* `azurerm_eventhub` - support for the `skip_empty_archives` property ([#3074](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3074))
* `azurerm_eventhub_namespace` - increase maximum `maximum_throughput_units` to 100 ([#3049](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3049))
* `azurerm_function_app` - exporting `possible_outbound_ip_addresses` ([#3043](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3043))
* `azurerm_iothub` - properties `batch_frequency_in_seconds`, `max_chunk_size_in_bytes`, `encoding`, `container_name`, `file_name_format` are now correctly diff'd depending on the type ([#2951](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2951))
* `azurerm_image` - support for the `zone_resilient` property ([#3100](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3100))
* `azurerm_kubernetes_cluster` - support for the `network_profile` property ([#2987](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2987))
* `azurerm_key_vault` - support for the `storage_permissions` property ([#3081](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3081))
* `azurerm_managed_disk` - support for managed disks up to 32TB ([#3062](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3062))
* `azurerm_mssql_elasticpool` - support setting the `zone_redundant` property ([#3104](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3104))
* `azurerm_redis_cache` - support for the `minimum_tls_version` property ([#3111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3111))
* `azurerm_storage_account` - support for configuring the Hierarchical Namespace state ([#3032](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3032))
* `azurerm_storage_account` - exposing the DFS File Secondary and Web endpoints ([#3110](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3110))
* `azurerm_virtual_machine` - support for managed disks up to 32TB ([#3062](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3062))
* `azurerm_virtual_machine_scale_set` - support for managed disks up to 32TB ([#3062](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3062))

BUG FIXES:

* `azurerm_application_gateway` - correctly populating backend addresses from both new and deprecated properties `fqdns`/`fqdn_list` ([#3085](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3085))
* `azurerm_key_vault_certificate` - making `contents` and `password` within the `certificate` block sensitive ([#3064](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3064))
* `monitor_metric_alert` - support for setting `aggregation` to `count`  ([#3047](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3047))
* `azurerm_virtual_network_gateway` - fixing a crash when `bgp_settings` had no elements ([#3038](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3038))
* `azurerm_virtual_machine_scale_set` - support setting `zones` to an empty list ([#3142](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3142))

## 1.23.0 (March 08, 2019)

FEATURES:

* **New Data Source:** `azurerm_api_management_group` ([#2809](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2809))
* **New Data Source:** `azurerm_api_management_product` ([#2953](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2953))
* **New Data Source:** `azurerm_api_management_user` ([#2954](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2954))
* **New Data Source:** `azurerm_availability_set` ([#2850](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2850))
* **New Data Source:** `azurerm_network_watcher` ([#2791](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2791))
* **New Data Source:** `azurerm_recovery_services_protection_policy_vm` ([#2974](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2974))
* **New Resource:** `azurerm_api_management_group` ([#2809](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2809))
* **New Resource:** `azurerm_api_management_group_user` ([#2972](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2972))
* **New Resource:** `azurerm_api_management_product` ([#2953](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2953))
* **New Resource:** `azurerm_api_management_product_group` ([#2984](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2984))
* **New Resource:** `azurerm_api_management_property` ([#2986](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2986))
* **New Resource:** `azurerm_api_management_user` ([#2954](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2954))
* **New Resource:** `azurerm_connection_monitor` ([#2791](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2791))
* **New Resource:** `azurerm_eventgrid_domain` ([#2884](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2884))
* **New Resource:** `azurerm_eventgrid_event_subscription` ([#2967](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2967))
* **New Resource:** `azurerm_lb_outbound_rule` ([#2912](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2912))
* **New Resource:** `azurerm_media_service_account` ([#2711](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2711))

IMPROVEMENTS:

* dependencies: upgrading to v25.1.0 of `github.com/Azure/azure-sdk-for-go` ([#2886](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2886))
* dependencies: upgrading to v11.4.0 of `github.com/Azure/go-autorest` ([#2886](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2886))
* `azurerm_application_gateway` - support for setting `path` within the `backend_http_settings` block ([#2879](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2879))
* `azurerm_application_gateway` - support for setting `connection_draining` to the `backend_http_settings` ([#2778](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2778))
* `azurerm_container_group` - support for specifying the `diagnostics` block ([#2763](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2763))
* `azurerm_iothub` - support for the `fallback_route` property ([#2764](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2764))
* `azurerm_key_vault` - support for 1024 `access_policy` blocks ([#2866](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2866))
* `azurerm_redis_cache` - support for configuring the `maxfragmentationmemory_reserved` in the `redis_configuration` block ([#2887](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2887))
* `azurerm_servicebus_namespace` - allowing `capacity` to be set to `0` for non-Premium SKU's ([#2920](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2920))
* `azurerm_service_fabric_cluster` - support for setting `capacities` and `placement_properties` ([#2936](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2936))
* `azurerm_storage_account` - exposing primary/secondary `_host` attributes ([#2792](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2792))

BUG FIXES:

* `azurerm_api_management` - switching to use API version `2018-01-01` rather than `2018-06-01-preview` ([#2958](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2958))
* `azurerm_application_gateway` - updating the default value for `file_upload_limit_mb` within the `waf_configuration` block to be `100` to match the documentation ([#3012](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3012))
* `azurerm_batch_pool` - updating `max_tasks_per_node` to be ForceNew ([#2856](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2856))
* `azurerm_key_vault_access_policy` - no longer silenty fails on creation of the `key_vault_id` property is invalid/doesn't exist ([#2922](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2922))
* `azurerm_policy_definition` - making the `metadata` field to computed ([#2939](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2939))
* `azurerm_redis_firewall_rule` - allowing underscores in the `name` field ([#2906](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2906))
* `azurerm_iothub` - marking the `connection_string` property as sensitive ([#3007](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3007))
* `azurerm_iothub` - ensuring the `type` property is alwaysa set ([#3007](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3007))

## 1.22.1 (February 14, 2019)

BUG FIXES:

* `azurerm_key_vault_access_policy` - will no longer fail to find the Key Vault if `key_vault_id` is empty ([#2874](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2874))
* `azurerm_key_vault_certificate` - will no longer fail to find the Key Vault if `key_vault_id` is ([#2874](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2874))
* `azurerm_key_vault_key` - will no longer fail to find the Key Vault if `key_vault_id` is ([#2874](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2874))
* `azurerm_key_vault_secret` - will no longer fail to find the Key Vault if `key_vault_id` is ([#2874](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2874))
* `azurerm_storage_container` - support for large numbers of containers within a storage account ([#2873](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2873))

## 1.22.0 (February 11, 2019)

UPGRADE NOTES:

* The v1.22 release includes a few new resources which are duplicates of existing resources, the purpose of this is to correct some invalid naming so that we can remove the mis-named resources in the next major version of the Provider. Please see [the upgrade guide](https://www.terraform.io/docs/providers/azurerm/guides/migrating-between-renamed-resources.html) for more information on how to migrate between these resources.
* The `azurerm_builtin_role_definition` Data Source has been deprecated in favour of the `azurerm_role_definition` Data Source, which now provides the same functionality and will be removed in the next major version of the AzureRM Provider (2.0) ([#2798](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2798))
* The `azurerm_log_analytics_workspace_linked_service` resource has been deprecated in favour of the (new) `azurerm_log_analytics_linked_service` resource and will be removed in the next major version of the AzureRM Provider (2.0) ([#2768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2768))
* The `azurerm_autoscale_setting` resource has been deprecated in favour of the (new) `azurerm_monitor_autoscale_setting` resource and will be removed in the next major version of the AzureRM Provider (2.0) ([#2768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2768))
* The `azurerm_metric_alertrule` resource has been deprecated in favour of the (new) `azurerm_monitor_metric_alertrule` resource and will be removed in the next major version of the AzureRM Provider (2.0) ([#2762](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2762))

FEATURES:

* **New Data Source:** `azurerm_policy_definition` ([#2788](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2788))
* **New Data Source:** `azurerm_servicebus_namespace` ([#2841](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2841))
* **New Resource:** `azurerm_ddos_protection_plan` ([#2654](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2654))
* **New Resource:** `azurerm_log_analytics_linked_service ` ([#2768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2768))
* **New Resource:** `azurerm_monitor_autoscale_setting` ([#2768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2768))
* **New Resource:** `azurerm_monitor_metric_alertrule` ([#2762](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2762))
* **New Resource:** `azurerm_network_interface_application_security_group_association` ([#2789](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2789))

DEPRECATIONS:

* Data Source `azurerm_key_vault_key` - deprecating the `vault_uri` property in favour of `key_vault_id` ([#2820](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2820))
* Data Source `azurerm_key_vault_secret` - deprecating the `vault_uri` property in favour of `key_vault_id` ([#2820](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2820))
* `azurerm_key_vault_certificate` - deprecating the `vault_uri` property in favour of `key_vault_id` ([#2820](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2820))
* `azurerm_key_vault_key` - deprecating the `vault_uri` property in favour of `key_vault_id` ([#2820](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2820))
* `azurerm_key_vault_access_policy` - deprecating the `vault_name` and `resource_group_name` properties in favour of `key_vault_id` ([#2820](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2820))
* `azurerm_key_vault_secret` - deprecating the `vault_uri` property in favour of `key_vault_id` ([#2820](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2820))
* `azurerm_application_gateway` - deprecating the `fqdn_list` field in favour of `fqdns` ([#2768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2768))
* `azurerm_application_gateway` - deprecating the `ip_address_list` field in favour of `ip_addresses` ([#2768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2768))
* `azurerm_builtin_role_definition` - deprecating in favour of the `azurerm_role_definition` data source, which now provides the same functionality ([#2798](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2798))
* `azurerm_log_analytics_workspace_linked_service` - deprecating in favour of the (renamed) `azurerm_log_analytics_linked_service` resource ([#2768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2768))
* `azurerm_monitor_autoscale_setting` - deprecating in favour of the (renamed) `azurerm_autoscale_setting` resource ([#2768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2768))
* `azurerm_network_interface ` - deprecating the `application_security_group_ids` field in favour of the new `azurerm_network_interface_application_security_group_association` resource ([#2789](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2789))

IMPROVEMENTS:

* dependencies: switching to Go Modules ([#2705](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2705))
* dependencies: upgrading to v11.3.2 of github.com/Azure/go-autorest ([#2744](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2744))
* Data Source: `azurerm_role_definition` - support for finding roles by name ([#2798](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2798))
* `azurerm_application_gateway` - support for the `http2` property ([#2735](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2735))
* `azurerm_application_gateway` - support for the `file_upload_limit_mb` property ([#2666](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2666))
* `azurerm_application_gateway` - support for the `custom_error_configuration` property ([#2783](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2783))
* `azurerm_application_gateway` - Support for `pick_host_name_from_backend_address` and `pick_host_name_from_backend_http_settings` properties ([#2658](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2658))
* `azurerm_app_service` - support for the `client_cert_enabled` property ([#2765](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2765))
* `azurerm_autoscale_setting` - support values from `0` to `1000` for the `minimum`, `maximum` and `default` properties ([#2815](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2815))
* `azurerm_batch_pool` - support for the `max_tasks_per_node` property ([#2805](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2805))
* `azurerm_cognitive_account` - exporting `primary_access_key` and `secondary_access_key` ([#2825](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2825))
* `azurerm_cosmosdb_account` - support for the `EnableAggregationPipeline`, `MongoDBv3.4` and ` mongoEnableDocLevelTTL` capabilities ([#2715](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2715))
* `azurerm_data_lake_store_file` - support file uploads greater then 4 megabytes ([#2633](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2633))
* `azurerm_function_app` - support for linux via the `linux_fx_version` property ([#2767](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2767))
* `azurerm_mssql_elasticpool` - support for setting `max_size_bytes` ([#2346](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2346))
* `azurerm_mssql_elasticpool` - support for setting `max_size_gb` ([#2695](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2695))
* `azurerm_postgresql_server` - support for version `10` and `10.2` ([#2768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2768))
* `azurerm_kubernetes_cluster` - add addtional validation ([#2772](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2772))
* `azurerm_signalr_service` - exporting `primary_access_key`, `secondary_access_key`, `primary_connection_string` and `secondary_connection_string` and secondary access keys and connection strings ([#2655](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2655))
* `azurerm_subnet` - support for additional subnet delegation types ([#2667](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2667))

BUG FIXES:

* `azurerm_azuread_application` - fixing a bug where `reply_uris` was set incorrectly ([#2729](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2729))
* `azurerm_batch_pool` - can now set multiple environment variables ([#2685](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2685))
* `azurerm_cosmosdb_account` - prevent occasional error when deleting the resource ([#2702](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2702))
* `azurerm_cosmosdb_account` - allow empty values for the `ip_range_filter` property ([#2713](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2713))
* `azurerm_express_route_circuit` - added the `premium` SKU back to validation logic ([#2692](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2692))
* `azurerm_firewall` - ensuring rules aren't removed during an update ([#2663](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2663))
* `azurerm_notification_hub_namespace` - now polls on creation to handle eventual consistency ([#2701](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2701))
* `azurerm_redis_cache` - locking on the Virtual Network/Subnet name to avoid a race condition ([#2725](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2725))
* `azurerm_service_bus_subscription` - name's can now start with a digit ([#2672](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2672))
* `azurerm_security_center` - increase the creation timeout to `30m` ([#2724](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2724))
* `azurerm_service_fabric_cluster` - no longer pass `reverse_proxy_endpoint_port` to the API when not specified ([#2747](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2747))
* `azurerm_subnet` - fixing a crash when service endpoints was nil ([#2742](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2742))
* `azurerm_subnet` - will no longer lose service endpoints during a virtual network update ([#2738](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2738))

## 1.21.0 (January 11, 2019)

FEATURES:

* **New Data Source:** `azurerm_application_insights` ([#2625](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2625))
* **New Data Source:** `azurerm_batch_account` ([#2428](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2428))
* **New Data Source:** `azurerm_batch_pool` ([#2461](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2461))
* **New Data Source:** `azurerm_lb` ([#2354](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2354))
* **New Data Source:** `azurerm_lb_backend_address_pool` ([#2354](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2354))
* **New Data Source:** `azurerm_virtual_machine` ([#2463](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2463))
* **New Resource:** `azurerm_application_insights_api_key` ([#2556](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2556))
* **New Resource:** `azurerm_batch_account` ([#2428](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2428))
* **New Resource:** `azurerm_batch_pool` ([#2461](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2461))
* **New Resource:** `azurerm_firewall_application_rule_collection` ([#2532](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2532))
* **New Resource:** `azurerm_policy_set_definition` ([#2535](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2535))

IMPROVEMENTS:

* config: support for specifying the `partner_id` for partner resource attribution ([#2643](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2643))
* dependencies: updating to `v24.0.0` of `Azure/azure-sdk-for-go` ([#2572](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2572))
* dependencies: upgrading the `network` SDK to `2018-08-01` ([#2433](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2433))
* Data Source: `azurerm_app_service` - exporting the `possible_outbound_ip_addresses` ([#2513](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2513))
* Data Source: `azurerm_azuread_application` - deprecating in favour of the split-out AzureAD Provider ([#2632](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2632))
* Data Source: `azurerm_azuread_service_principal` - deprecating in favour of the split-out AzureAD Provider ([#2632](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2632))
* Data Source: `azurerm_container_registry` - now exports `tags` ([#2607](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2607))
* Data Source: `azurerm_network_interface` - now exports `ip_configuration.private_ip_address_version` ([#2646](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2646))
* Data Source: `azurerm_public_ip` - now exports `location`, `sku`, `allocation_method`, `reverse_fqdn` and `zones` ([#2576](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2576))
* `azurerm_app_service` - exporting the `possible_outbound_ip_addresses` ([#2513](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2513))
* `azurerm_azuread_application` - deprecating in favour of the split-out AzureAD Provider ([#2632](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2632))
* `azurerm_azuread_service_principal` - deprecating in favour of the split-out AzureAD Provider ([#2632](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2632))
* `azurerm_azuread_service_principal_password` - deprecating in favour of the split-out AzureAD Provider ([#2632](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2632))
* `azurerm_cognitive_account` - support for the `SpeechServices` kind ([#2583](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2583))
* `azurerm_container_group` - deprecated container properties `port` and `protocol` for ports allowing for multiple ports ([#1930](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1930))
* `azurerm_eventhub_namespace` - support for `kafka_enabled` ([#2395](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2395))
* `azurerm_firewall` - renaming the `public_ip_address_id` property to `ip_address_id` ([#2433](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2433))
* `azurerm_kubernetes_cluster` - support for Virtual Nodes ([#2641](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2641))
* `azurerm_kubernetes_cluster` - the `dns_prefix` now forces a new resource and is properly validated ([#2611](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2611))
* `azurerm_log_analytics_workspace_linked_service` - now correctly handels uppcase `workspace_name` values  ([#2594](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2594))
* `azurerm_network_interface` - support for IPv6 addresses ([#2548](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2548))
* `azurerm_policy_assignment` - support for Managed Service Identity ([#2549](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2549))
* `azurerm_policy_assignment` - support exclusions with the `not_scopes` property ([#2620](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2620))
* `azurerm_policy_definition` - polices can now be assigned to a management group ([#2490](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2490))
* `azurerm_policy_set_definition` - policy sets can now be assigned to a management group ([#2618](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2618))
* `azurerm_public_ip` - deprecated `public_ip_address_allocation` in favor of `allocation_method` to better match the SDK ([#2576](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2576))
* `azurerm_redis_cache` - add availability zone support ([#2580](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2580))
* `azurerm_service_fabric_cluster` - support for `azure_active_directory` ([#2553](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2553))
* `azurerm_service_fabric_cluster` - support for `reverse_proxy_certificate` ([#2544](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2544))
* `azurerm_service_fabric_cluster` - support for `reverse_proxy_endpoint_port` ([#2544](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2544))
* `azurerm_subnet` - support for delegation ([#2042](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2042))

BUG FIXES:

* Data Source: `azurerm_managed_disk` - exposing the `create_option` field ([#2597](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2597))
* Data Source: `azurerm_network_interface` - exposing `application_security_group_ids` within the `ip_configuration` block ([#2599](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2599))
* Data Source: `azurerm_snapshot` - ensuring `disk_size_gb` is set ([#2596](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2596))
* Data Source: `azurerm_storage_account` - ensuring the `account_replication_type` field is set correctly ([#2595](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2595))
* `azurerm_app_service` - handling connection strings being in any order ([#2609](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2609))
* `azurerm_app_service_slot` - handling connection strings being in any order ([#2609](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2609))
* `azurerm_network_security_rule` - the properties `source_application_security_group_ids` and `destination_application_security_group_ids` are now correctly read & imported ([#2558](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2558))
* `azurerm_role_assignment` - retrieving the role definition name during import ([#2565](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2565))
* `azurerm_template_deployment` - fixing regression and supportting nested template deployments ([#2514](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2514))

## 1.20.0 (December 12, 2018)

FEATURES:

* **New Data Source:** `azurerm_monitor_action_group` ([#2430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2430))
* **New Resource:** `azurerm_mariadb_database` ([#2445](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2445))
* **New Resource:** `azurerm_mariadb_server` ([#2406](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2406))
* **New Resource:** `azurerm_signalr_service` ([#2410](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2410))

IMPROVEMENTS:

* authentication: switching to use the shared Azure authentication library ([#2355](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2355))
* authentication: support for authenticating using a Service Principal with a Client Certificate ([#2471](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2471))
* authentication: requesting a token using the audience address ([#2381](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2381))
* authentication: switching to request tokens from the Azure CLI ([#2387](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2387))
* sdk: upgrading to version `2018-05-01` of the Policy API ([#2386](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2386))
* Data Source: `azurerm_kubernetes_cluster` - support for Role Based Access Control without Azure AD ([#2495](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2495))
* Data Source: `azurerm_kubernetes_cluster` - exposing the `clusterAdmin` credentials ([#2495](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2495))
* Data Source: `azurerm_subscriptions` - ability to filtering by prefix/contains on the Display Name ([#2429](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2429))
* `azurerm_app_service` - support for configuring `app_command_line` in the `site_config` block ([#2350](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2350))
* `azurerm_app_service_plan` - deprecated the `properties` and moved `app_service_environment_id`, `per_site_scaling` and `reserved` to the top level  ([#2442](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2442))
* `azurerm_app_service_slot` - support for configuring `app_command_line` in the `site_config` block ([#2350](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2350))
* `azurerm_application_insights` - added `Node.JS` application type ([#2407](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2407))
* `azurerm_container_registry` - support for geo-replication via the `georeplication_locations` property ([#2055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2055))
* `azurerm_key_vault` - exposed `backup` and `restore` permissions made `key_permissions` and `secret_permissions` optional ([#2363](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2363))
* `azurerm_kubernetes_cluster` - support for Role Based Access Control without Azure AD ([#2495](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2495))
* `azurerm_kubernetes_cluster` - exposing the `clusterAdmin` credentials ([#2495](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2495))
* `azurerm_mssql_elasticpool` - deprecated the `elastic_pool_properties` property and moved `max_size_bytes` and `zone_redundant` to the top level ([#2378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2378))
* `azurerm_mysql_server` - support for new skus `GP_Gen5_64` and `MO_Gen5_32` ([#2446](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2446))
* `azurerm_postgresql_server` support for new skus `GP_Gen5_64` and `MO_Gen5_32` - ([#2447](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2447))

BUG FIXES:

* Data Source: `azurerm_logic_app_workflow` - ensuing the parameters are a string prior to flattening ([#2348](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2348))
* Data Source: `azurerm_public_ip` - ensuing properties always exist ([#2448](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2448))
* Data Source: `azurerm_route_table` - validation updated to prevent empty and blank `property` values from causing a panic ([#2467](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2467))
* `azurerm_key_vault` - fixing a deadlock situation where multiple subnets are used from the same virtual network ([#2324](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2324))
* `azurerm_eventhub` - making the `partition_count` field ForceNew ([#2400](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2400))
* `azurerm_eventhub` - now validates that the `storage_account_id` is a proper resource ID  ([#2374](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2374))
* `azurerm_mssql_elasticpool` - relaxed validation of the `name` property ([#2398](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2398))
* `azurerm_recovery_services_protection_policy_vm` - added the `timezone` property ([#2404](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2404))
* `azurerm_route_table` - validation updated to prevent empty and blank `property` values from causing a panic ([#2467](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2467))
* `azurerm_sql_server` - only updating the `admin_login_password` when it's changed, allowing this to be managed outside of Terraform ([#2263](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2263))
* `azurerm_virtual_machine` - nil-checking properties prior to accessing ([#2365](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2365))

## 1.19.0 (November 15, 2018)

FEATURES:

* **New Data Source:** `azurerm_key_vault_key` ([#2231](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2231))
* **New Data Source:** `azurerm_monitor_diagnostic_setting` ([#1291](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1291))
* **New Resource:** `azurerm_iothub_consumer_group` ([#2243](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2243))
* **New Resource:** `azurerm_monitor_diagnostic_setting` ([#1291](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1291))
* **New Resource:** `azurerm_mssql_elasticpool` ([#2071](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2071))

IMPROVEMENTS:

* dependencies: switching to Go 1.11 ([#2229](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2229))
* authentication: refactoring to allow authentication modes to be feature-toggled ([#2199](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2199))
* Data Source: `azurerm_kubernetes_cluster` - support for `role_based_access_control` ([#1820](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1820))
* `azurerm_app_service` - support for PHP 7.2 ([#2308](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2308))
* `azurerm_app_service_slot` - support for PHP 7.2 ([#2308](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2308))
* `azurerm_databricks_workspace` - fixing validation on the `name` field ([#2221](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2221))
* `azurerm_function_app` - support for the `enable_builtin_logging` property ([#2268](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2268))
* `azurerm_kubernetes_cluster` - support for `role_based_access_control` ([#1820](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1820))
* `azurerm_network_interface` - deprecating `internal_fqdn` since it's no longer setable/returned by Azure ([#2253](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2253))
* `azurerm_shared_image_version` - allowing larger numbers for versions ([#2301](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2301))
* `azurerm_virtual_machine` - support for assigning both a system and a user managed identity ([#2188](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2188))
* `azurerm_virtual_machine_scale_set` - support for assigning both a system and a user managed identity ([#2188](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2188))
* `azurerm_virtual_machine_scale_set` - support for setting `eviction_policy` ([#2226](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2226))
* `azurerm_virtual_network_gateway` - support for Zone Redundant Gateways ([#2260](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2260))

BUG FIXES:

* Data Source: `azurerm_api_management` - ensuring the `public_ip_addresses` field is set ([#2310](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2310))
* `azurerm_api_management` - ensuring the `public_ip_addresses` field is set ([#2310](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2310))
* `azurerm_application_gateway` - refactoring to ensure all fields are set ([#2054](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2054))
* `azurerm_application_gateway` - SSL certificates no longer continually diff ([#2054](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2054))
* `azurerm_azuread_application` - fix regression and allow `http` for `identifier_uris` and `reply_urls` properties ([#2320](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2320))
* `azurerm_cosmosdb_account` - the `ip_range_filter` range filter now allows /32 ip addresses  ([#2222](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2222))
* `azurerm_public_ip` - fixing the casing of the `ip_version` / `public_ip_address_allocation` fields ([#2296](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2296))
* `azurerm_recovery_services_protected_vm` - VM can now be in a different resource group then the vault ([#2287](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2287))
* `azurerm_role_assignment` - will now wait after a Service Principal is created ([#2204](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2204))
* `azurerm_route` - allowing setting `next_hop_in_ip_address` to an empty value ([#2184](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2184))
* `azurerm_route_table` - allowing setting `next_hop_in_ip_address` to an empty value ([#2184](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2184))
* `azurerm_virtual_network_gateway` - plan is now empty when `bgp_settings` is omitted ([#2304](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2304))
* `azurerm_virtual_network` - add valdiation to prevent panics ([#2305](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2305))

## 1.18.0 (November 02, 2018)

FEATURES:

* **New Resource:** `azurerm_devspace_controller` ([#2086](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2086))
* **New Resource:** `azurerm_log_analytics_workspace_linked_service` ([#2139](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2139))

IMPROVEMENTS:

* authentication: decoupling the authentication methods from the provider to enable splitting out the authentication library ([#2197](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2197))
* authentication: using the Proxy from the Environment, if set ([#2133](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2133))
* dependencies: upgrading to v21.3.0 of `github.com/Azure/azure-sdk-for-go` ([#2163](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2163))
* refactoring:  decoupling Resource Provider Registration to enable splitting out the authentication library ([#2197](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2197))
* sdk: upgrading to `2018-10-01` of the `containerinstance` sdk ([#2174](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2174))
* `azurerm_automation_account` - exposing `dsc_server_endpoint`, `dsc_primary_access_key`, `dsc_secondary_access_key` properties ([#2166](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2166))
* `azurerm_automation_account` - support for the `free` SKU ([#2166](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2166))
* `azurerm_client_config` - ensuring the `service_principal_application_id` and `service_principal_object_id` are always set ([#2120](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2120))
* `azurerm_cosmosdb_account` - support for the `enable_multiple_write_locations` property ([#2109](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2109))
* `azurerm_eventhub_namespace` - allow `maximum_throughput_units` to be zero ([#2124](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2124))
* `azurerm_key_vault_certificate` - support for setting `extended_key_usage` ([#2128](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2128))
* `azurerm_key_vault_certificate` - support for setting `subject_alternative_names` ([#2123](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2123))
* `azurerm_managed_disk` - support for the `UltraSSD_LRS` storage account type ([#2118](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2118))
* `azurerm_monitor_activity_log_alert` - support the criteria fields `resource_provider`, `resource_type`, `resource_group` ([#2150](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2150))
* `azurerm_recovery_services_protected_vm` - `backup_policy_id` is now required ([#2154](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2154))
* `azurerm_sql_database` - adding validation to `requested_service_objective_name` ([#2125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2125))
* `azurerm_virtual_network_gateway` - support for `OpenVPN` as a client protocol option ([#2126](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2126))
* `azurerm_virtual_machine_scale_set` - support for the `application_security_group_ids` property of `ip_configuration`  ([#2009](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2009))
* `azurerm_virtual_machine_scale_set` - support for a Rolling Upgrade Policy with Automatic OS upgrades ([#922](https://github.com/terraform-providers/terraform-provider-azurerm/issues/922))

BUG FIXES:

* security: removing the `Authorization` header from the debug logs ([#2131](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2131))
* `azurerm_api_management` - validating the Key Vault Secret ID for the `key_vault_id` field in the `hostname_configuration` block ([#2189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2189))
* `azurerm_function_app` - correctly marking the resource as missing upon manual deletion ([#2111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2111))
* `azurerm_kubernetes_cluster` - changing `os_disk_size_gb` to computed as the API now returns a valid default ([#2117](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2117))
* `azurerm_public_ip` - `domain_name_label` validation now allows 63 characters ([#2122](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2122))
* `azurerm_virtual_machine` - making `availability_set_id` conflict with `zones` ([#2185](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2185))


## 1.17.0 (October 18, 2018)

UPGRADE NOTES:

* `azurerm_virtual_machine_scale_set` - the field `primary` within the `ip_configuration` block within the `network_profile` block is now Required, to match behavioural changes in the Azure API. ([#2035](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2035))

FEATURES:

* **New Data Source:** `azurerm_monitor_log_profile` ([#1792](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1792))
* **New Resource:** `azurerm_api_management` ([#1516](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1516))
* **New Resource:** `azurerm_automation_dsc_configuration` ([#1512](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1512))
* **New Resource:** `azurerm_automation_dsc_nodeconfiguration` ([#1512](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1512))
* **New Resource:** `azurerm_automation_module` ([#1512](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1512))
* **New Resource:** `azurerm_cognitive_account` ([#962](https://github.com/terraform-providers/terraform-provider-azurerm/issues/962))
* **New Resource:** `azurerm_databricks_workspace` ([#1134](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1134))
* **New Resource:** `azurerm_dev_test_policy` ([#2070](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2070))
* **New Resource:** `azurerm_dev_test_linux_virtual_machine` ([#2058](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2058))
* **New Resource:** `azurerm_dev_test_windows_virtual_machine` ([#2058](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2058))
* **New Resource:** `azurerm_monitor_activitylog_alert` ([#1989](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1989))
* **New Resource:** `azurerm_monitor_metric_alert` ([#2026](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2026))
* **New Resource:** `azurerm_monitor_log_profile` ([#1792](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1792))
* **New Resource:** `azurerm_network_interface_application_gateway_backend_address_pool_association` ([#2079](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2079))
* **New Resource:** `azurerm_network_interface_backend_address_pool_association` ([#2079](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2079))
* **New Resource:** `azurerm_network_interface_nat_rule_association` ([#2079](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2079))
* **New Resource:** `azurerm_recovery_services_protection_policy_vm` ([#1978](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1978))
* **New Resource:** ` azurerm_recovery_services_protected_vm` ([#1637](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1637))
* **New Resource:** `azurerm_security_center_contact` ([#2045](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2045))
* **New Resource:** `azurerm_security_center_subscription_pricing` ([#2043](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2043))
* **New Resource:** `azurerm_security_center_workspace` ([#2072](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2072))
* **New Resource:** `azurerm_subnet_network_security_group_association` ([#1933](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1933))
* **New Resource:** `azurerm_subnet_route_table_association ` ([#1933](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1933))

BUG FIXES:

* Data Source `azurerm_subnet` - fixing the ordering of the resource group name and network name in the error message ([#2017](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2017))
* `azurerm_kubernetes_cluster` - using the correct casing for the `addon_profile` `oms_agent` property ([#1995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1995))
* `azurerm_service_bus_queue` - support for `max_delivery_count` ([#2028](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2028))
* `azurerm_redis_cache` - `capcity` can now be successfully changed ([#2088](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2088))
* `azurerm_virtual_machine_scale_set` - `primary` is now required within the `ip_configuration` block within `network_profile` (matching a behavioural change with the Azure API) ([#2035](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2035))

IMPROVEMENTS:

* `azurerm_application_gateway` - support for the `StandardV2` and `WAFV2` skus and tiers ([#2015](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2015))
* `azurerm_container_group` - adding the `secure_environment_variables` property ([#2024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2024))
* `azurerm_dev_test_virtual_network` - support for managing the Subnet ([#2041](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2041))
* `azurerm_key_vault` - support for Virtual Network Rules ([#2027](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2027))
* `azurerm_kubernetes_cluster` - changing the `oms_agent` property no longer forces a new resource ([#2021](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2021))
* `azurerm_postgresql_virtual_network_rule` - support for the `ignore_missing_vnet_service_endpoint` ([#2056](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2056))
* `azurerm_public_ip` - support for IPv6 addresses ([#2019](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2019))
* `azurerm_search_service` - adding the administrative `primary_key` and `secondary_key` propeties ([#2074](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2074))
* `azurerm_role_definition` - adding the `data_actions` and `not_data_actions` to the data source ([#2110](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2110))
* `azurerm_storage_container` - changing `container_access_type` no longer forces a new resource ([#2075](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2075))
* `azurerm_user_assigned_identity` - now exports the `client_id` property ([#2078](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2078))


## 1.16.0 (October 01, 2018)

UPGRADE NOTES:

* `azurerm_azuread_application` - the properties `homepage`, `identifier_uris` and `reply_urls` are now required to be `https` as required by Azure ([#1960](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1960))

FEATURES:

* **New Data Source:** `azurerm_dev_test_lab` ([#1944](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1944))
* **New Data Source:** `azurerm_shared_image` ([#1987](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1987))
* **New Data Source**: `azurerm_shared_image_gallery` ([#1987](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1987))
* **New Data Source:** `azurerm_shared_image_version` ([#1987](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1987))
* **New Resource:** `azurerm_dev_test_lab` ([#1944](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1944))
* **New Resource:** `azurerm_dev_test_virtual_network` ([#1944](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1944))
* **New Resource:** `azurerm_shared_image` ([#1987](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1987))
* **New Resource**: `azurerm_shared_image_gallery` ([#1987](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1987))
* **New Resource:** `azurerm_shared_image_version` ([#1987](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1987))

IMPROVEMENTS:

* dependencies: upgrading to v21.0.0 of `github.com/Azure/azure-sdk-for-go` ([#1996](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1996))
* `azurerm_cosmosdb_account` - adding the `is_virtual_network_filter_enabled` and `virtual_network_rule` propeties ([#1961](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1961))

BUG FIXES:

* Data Source `azurerm_builtin_role_definition`: support for `data_actions` and `not_data_actions` ([#2000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2000))
* `azurerm_app_service_plan` - exposing additional information on failure ([#1926](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1926))
* `azurerm_app_service_custom_hostname_binding` - handling multiple bindings being created in parallel ([#1970](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1970))
* `azurerm_lb_rule` - allow `0` for `frontend_port` and `backend_port` again ([#1951](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1951))
* `azurerm_public_ip` - correctly reading and importing the `idle_timeout_in_minutes` property ([#1925](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1925))
* `azurerm_role_assignment` - only retry on errors when they are retryable ([#1934](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1934))
* `azurerm_role_definition` - support for the `data_actions` and `not_data_action` blocks ([#1971](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1971))
* `azurerm_service_fabric_cluster` - allow two `client_certificate_thumbprint` blocks ([#1938](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1938))
* `azurerm_service_fabric_cluster` - support for specifying the `cluster_code_version` field ([#1945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1945))
* `azurerm_virtual_network` - exposing the `id` of each subnet ([#1913](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1913))
* `azurerm_virtual_machine` - handling the Managed Disk ID being nil ([#1947](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1947))
* `azurerm_virtual_machine_data_disk_attachment` - supporting data disk attachments when a VM Extension is installed ([#1950](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1950))
* `azurerm_virtual_machine_scale_set` - making `admin_password` in the `os_profile` block optional again ([#1958](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1958))

## 1.15.0 (September 14, 2018)

FEATURES:

* **New Resource:** `azurerm_firewall` ([#1627](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1627))
* **New Resource:** `azurerm_firewall_network_rule_collection` ([#1627](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1627))
* **New Resource:** `azurerm_mysql_virtual_network_rule` ([#1879](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1879))

IMPROVEMENTS:

* dependencies: upgrading to v20.1.0 of `github.com/Azure/azure-sdk-for-go` ([#1861](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1861))
* dependencies: upgrading to v10.15.4 of `github.com/Azure/go-autorest` ([#1861](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1861)) ([#1909](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1909))
* sdk: upgrading to version `2018-06-01` of the Compute API's ([#1861](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1861))
* `azurerm_automation_runbook` - support for specifying the content field ([#1696](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1696))
* `azurerm_app_service` - adding the `virtual_network_name` property ([#1896](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1896))
* `azurerm_app_service_slot` - adding the `virtual_network_name` property ([#1896](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1896))
* `azurerm_key_vault_certificate` - adding the `thumbprint` property ([#1904](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1904))
* `azurerm_servicebus_queue` - adding validation for ISO8601 Durations ([#1921](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1921))
* `azurerm_servicebus_topic` - adding validation for ISO8601 Durations ([#1921](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1921))
* `azurerm_sql_database` - adding the `threat_detection_policy` property ([#1628](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1628))
* `azurerm_virtual_network` - adding validation to `name` preventing empty values ([#1898](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1898))
* `azurerm_virtual_machine` - support for the `managed_disk_type` of `StandardSSD_LRS` ([#1901](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1901))
* `azurerm_virtual_machine_scale_set` - support for the `managed_disk_type` of `StandardSSD_LRS` ([#1901](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1901))
* `azurerm_virtual_network_gateway` - additional validation ([#1899](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1899))

BUG FIXES:

* Data Source: `azurerm_azuread_service_principal` - passing a filter containing the name to Azure rather than querying locally ([#1862](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1862))
* Data Source: `azurerm_azuread_service_principal` - passing a filter containing the name to Azure rather than querying locally ([#1862](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1862))
* `azurerm_logic_app_trigger_http_request` - `relative_path` property now allows `/`s and `{}`s ([#1918](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1918))
* `azurerm_role_assignment` - parsing the Resource ID during deletion ([#1887](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1887))
* `azurerm_role_definition` - parsing the Resource ID during deletion ([#1887](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1887))
* `azurerm_servicebus_namespace` - polling for the deletion of the namespace ([#1908](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1908))

## 1.14.0 (September 06, 2018)

FEATURES:

* **New Data Source:** `azurerm_management_group` ([#1877](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1877))
* **New Resource:** `azurerm_management_group` ([#1788](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1788))
* **New Resource:** `azurerm_postgresql_virtual_network_rule` ([#1774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1774))

IMPROVEMENTS:

* authentication: making the client registration consistent ([#1845](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1845))
* `azurerm_application_insights` - support for the `MobileCenter` kind ([#1878](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1878))
* `azurerm_function_app` - removing validation from the `version` field ([#1872](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1872))
* `azurerm_iothub` - exporting the `event_hub_events_endpoint`, `event_hub_events_path`, `event_hub_operations_endpoint` and `event_hub_operations_path` fields ([#1789](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1789))
* `azurerm_iothub` - support for `endpoint` and `route` blocks ([#1693](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1693))
* `azurerm_kubernetes_cluster` - making `linux_profile` optional ([#1821](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1821))
* `azurerm_storage_blob` - support for import ([#1816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1816))
* `azurerm_storage_container` - support for import ([#1816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1816))
* `azurerm_storage_queue` - support for import ([#1816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1816))
* `azurerm_storage_table` - support for import ([#1816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1816))

BUG FIXES:

* `azurerm_data_lake_store_file` - updating the Resource ID to match the file path ([#1856](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1856))
* `azurerm_eventhub` - updating the validation to support periods, hyphens and underscores ([#1795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1795))
* `azurerm_eventhub_authorization_rule` - updating the validation error ([#1795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1795))
* `azurerm_eventhub_consumer_group` - updating the validation to support periods, hyphens and underscores ([#1795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1795))
* `azurerm_eventhub_namespace` - updating the validation error ([#1795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1795))
* `azurerm_function_app` - support for names in upper-case ([#1835](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1835))
* `azurerm_kubernetes_cluster` - removing validation for the `pod_cidr` field when `network_plugin` is set to `azure` ([#1798](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1798))
* `azurerm_logic_app_workflow` - ensuring parameters are strings ([#1843](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1843))
* `azurerm_virtual_machine` - setting the `image_uri` property within the `storage_os_disk` block ([#1799](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1799))
* `azurerm_virtual_machine_data_disk_attachment` - obtaining a basic view, rather than the entire instance view of the Virtual Machine to work around an issue in the API ([#1855](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1855))

## 1.13.0 (August 15, 2018)

FEATURES:

* **New Data Source:** `azurerm_log_analytics_workspace` ([#1755](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1755))
* **New Resource:** `azurerm_monitor_action_group` ([#1725](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1725))

IMPROVEMENTS:

* dependencies: upgrading to `2018-04-01` of the IoTHub SDK ([#1717](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1717))
* Azure CLI Auth - using the `USERPROFILE` environment variable to locate the users home directory, if set ([#1718](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1718))
* Data Source `azurerm_kubernetes_cluster` - exposing the `max_pods` field within the `agent_pool_profile` block ([#1753](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1753))
* Data Source: `azurerm_kubernetes_cluster` - exposing the `add_on_profile` block ([#1751](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1751))
* `azurerm_automation_schedule` - adding the `week_days`, `month_days` and `monthly_occurrence` properties ([#1626](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1626))
* `azurerm_container_group` - adding a new `commands` field / deprecating the `command` field ([#1740](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1740))
* `azurerm_iothub` - support for the `Basic` SKU ([#1717](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1717))
* `azurerm_kubernetes_cluster` - support for `max_pods` within the `agent_pool_profile` block ([#1753](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1753))
* `azurerm_kubernetes_cluster` - support for the `add_on_profile` block ([#1751](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1751))
* `azurerm_kubernetes_cluster` - validation for when `pod_cidr` is set with a `network_plugin` set to `azure` ([#1763](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1763))
* `azurerm_kubernetes_cluster` - `client_id` and `client_secret` in the `service_principal` block are now ForceNew ([#1737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1737))
* `azurerm_kubernetes_cluster` - `docker_bridge_cidr`, `dns_service_ip` and `service_cidr` are now conditionally set ([#1715](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1715))
* `azurerm_lb_nat_rule` - `protocol` property now supports `All` ([#1736](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1736))
* `azurerm_lb_nat_pool` - `protocol` property now supports `All` ([#1748](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1748))
* `azurerm_lb_probe` - `protocol` property now supports `Https` ([#1742](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1742))
* `azurerm_lb_rule` - support for the `All` protocol / adding validation ([#1754](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1754))

BUG FIXES:

* `azurerm_application_insights` - handling a `HTTP 201` being returned from the Create API which working around a breaking change in the API ([#1769](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1769))
* `azurerm_autoscale_setting` - filtering out the `$tags` tag ([#1770](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1770))
* `azurerm_eventhub` - allowing underscores in the name field ([#1768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1768))
* `azurerm_eventhub_authorization_rule` - allowing underscores in the name field ([#1768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1768))
* `azurerm_eventhub_consumer_group` - allowing underscores in the name field ([#1768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1768))

## 1.12.0 (August 03, 2018)

UPGRADE NOTES:

* **Please Note:** When upgrading to v1.12.0 of the Azure Provider, you may need to specify the `priority` of any VM Scale Sets created between v1.6 of the Provider and v1.12. ([#1586](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1586))

FEATURES:

* **New Data Source:** `azurerm_container_registry` ([#1642](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1642))
* **New Resource:** `azurerm_service_fabric_cluster` ([#4](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4))

IMPROVEMENTS:

* sdk: switching from `WaitForCompletion` -> `WaitForCompletionRef` when polling Future's ([#1660](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1660))
* Data Source: `azurerm_kubernetes_cluster` - support for specifying the `network_profile` block ([#1479](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1479))
* Data Source: `azurerm_kubernetes_cluster` - outputting the `node_resource_group` field ([#1649](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1649))
* `azurerm_kubernetes_cluster` - support for specifying the `network_profile` block ([#1479](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1479))
* `azurerm_kubernetes_cluster` - outputting the `node_resource_group` field ([#1649](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1649))
* `azurerm_role_assignment` - retrying resource creation to match the Azure CLI's behaviour ([#1647](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1647))
* `azurerm_virtual_machine` - setting the connection information for Provisioners ([#1646](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1646))


BUG FIXES:

* `azurerm_virtual_machine_scale_set` - removing the default of `priority`, since this isn't set on older instances. ([#1586](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1586))

## 1.11.0 (July 25, 2018)

FEATURES:

* **New Resource:** `azurerm_data_lake_store_file` ([#1261](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1261))

IMPROVEMENTS:

* `azurerm_app_service` - support for `min_tls_version` in the `site_config` block ([#1601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1601))
* `azurerm_app_service_slot` - support for `min_tls_version` in the `site_config` block ([#1601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1601))
* `azurerm_data_lake_store` - support for enabling/disabling encryption ([#1623](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1623))
* `azurerm_data_lake_store` - support for managing the firewall state ([#1623](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1623))

BUG FIXES:

* `azurerm_servicebus_topic` - the `name` property now allows the ~ character ([#1640](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1640))

## 1.10.0 (July 21, 2018)

FEATURES:

* **New Data Source:** `azurerm_azuread_application` ([#1552](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1552))
* **New Data Source:** `azurerm_logic_app_workflow` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Data Source:** `azurerm_notification_hub` ([#1589](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1589))
* **New Data Source:** `azurerm_notification_hub_namespace` ([#1589](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1589))
* **New Data Source:** `azurerm_service_principal` ([#1564](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1564))
* **New Resource:** `azurerm_autoscale_setting` ([#1140](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1140))
* **New Resource:** `azurerm_data_lake_analytics_account` ([#1618](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1618))
* **New Resource:** `azurerm_data_lake_analytics_firewall_rule` ([#1618](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1618))
* **New Resource:** `azurerm_eventhub_namespace_authorization_rule` ([#1572](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1572))
* **New Resource:** `azurerm_logic_app_action_custom` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_logic_app_action_http` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_logic_app_trigger_custom` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_logic_app_trigger_http_request` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_logic_app_trigger_recurrence` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_logic_app_workflow` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_notification_hub` ([#1589](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1589))
* **New Resource:** `azurerm_notification_hub_authorization_rule` ([#1589](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1589))
* **New Resource:** `azurerm_notification_hub_namespace ` ([#1589](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1589))
* **New Resource:** `azurerm_servicebus_queue_authorization_rule` ([#1543](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1543))
* **New Resource:** `azurerm_service_principal` ([#1564](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1564))
* **New Resource:** `azurerm_service_principal_password` ([#1564](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1564))

IMPROVEMENTS:

* authentication: Refreshing the Service Principal Token before using it ([#1544](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1544))
* dependencies: updating to`2018-02-01` of the App Service SDK ([#1436](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1436))
* `azurerm_app_service` - support for setting `ftps_settings` in the `site_config` block ([#1577](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1577))
* `azurerm_app_service` - support for running containers ([#1578](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1578))
* `azurerm_app_service_slot` - support for Managed Service Identity ([#1579](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1579))
* `azurerm_app_service_slot` - Slots can now be updated in-place ([#1436](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1436))
* `azurerm_container_group` - support for images hosted in a private registry ([#1529](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1529))
* `azurerm_function_app` - adding support for the `site_credential` block ([#1567](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1567))
* `azurerm_function_app` - only setting `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` for Consumption Apps ([#1515](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1515))
* `azurerm_mysql_server` - changing `tier` or `family` in `sku` property no longer destroys existing resource ([#1598](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1598))
* `azurerm_network_security_rule` - a maximum of 1 Application Security Group can be set per Security Rule  ([#1587](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1587))
* `azurerm_postgresql_server` - changing `tier` or `family` in `sku` property no longer destroys existing resource ([#1598](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1598))
* `azurerm_virtual_machine_scale_set` - `sku` property is now a list #1558 ([#1558](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1558))

BUG FIXES:

* `azurerm_application_insights` - fixing a bug where `application_type` was set to `other` ([#1563](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1563))
* `azurerm_lb` - allow `subnet_id` to be set to an empty value ([#1588](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1588))
* `azurerm_servicebus_subscription` - only sending `correlation_filter` values if they're set ([#1565](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1565))
* `azurerm_servicebus_subscription` - setting the `default_message_ttl` field ([#1568](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1568))
* `azurerm_snapshot` - allowing dashes in the `name` field ([#1574](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1574))
* `azurerm_traffic_manager_endpoint` - working around a bug in the API by setting `target` to nil when a `target_resource_id` is specified ([#1546](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1546))

## 1.9.0 (July 11, 2018)

FEATURES:

* **New Resource:** `azurerm_azuread_application` ([#1269](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1269))
* **New Resource:** `azurerm_data_lake_store_firewall_rule` ([#1499](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1499))
* **New Resource:** `azurerm_key_vault_access_policy` ([#1149](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1149))
* **New Resource:** `azurerm_scheduler_job` ([#1172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1172))
* **New Resource:** `azurerm_servicebus_namespace_authorization_rule` ([#1498](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1498))
* **New Resource:** `azurerm_user_assigned_identity` ([#1448](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1448))

IMPROVEMENTS:

* dependencies: updating the `containerservice` SDK to `2018-03-31` to support AKS GA ([#1474](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1474))
* dependencies: updating to `v18.0.0` of `Azure/azure-sdk-for-go` ([#1487](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1487))
* dependencies: updating to `v10.12.0` of `Azure/go-autorest` ([#1487](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1487))
* `azurerm_application_gateway` - adding `minimum_servers` to the probe resource ([#1510](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1510))
* `azurerm_cdn_profile` - support for `Standard_ChinaCdn` and `Standard_Microsoft` SKU's ([#1465](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1465))
* `azurerm_cosmosdb_account` - checking to see if the name is in use before creating ([#1464](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1464))
* `azurerm_cosmosdb_account` - fixing the validation on the `ip_range_filter` field ([#1463](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1463))
* `azurerm_dns_zone` - support for Private DNS Zones ([#1404](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1404))
* `azurerm_image` - change os_disk property to a list and add additional property validation ([#1443](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1443))
* `azurerm_lb` - allow `private_ip_address` to be set to an empty value ([#1481](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1481))
* `azurerm_mysql_server` - changing the `storage_mb` property no longer forces a new resource ([#1532](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1532))
* `azurerm_postgresql_server` - changing the `storage_mb` property no longer forces a new resource ([#1532](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1532))
* `azurerm_servicebus_queue` - `enable_partitioning` can now be enabled for `Basic` and `Standard` tiers ([#1391](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1391))
* `azurerm_virtual_machine` - support for specifying user assigned identities ([#1448](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1448))
* `azurerm_virtual_machine` - making the `content` field in the `additional_unattend_config`  block (within `os_profile_windows_config`) sensitive ([#1471](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1471))
* `azurerm_virtual_machine_data_disk_attachment` - adding support for `write_accelerator_enabled` ([#1473](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1473))
* `azurerm_virtual_machine_scale_set` - ensuring we set the `vhd_containers` field to fix a crash ([#1411](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1411))
* `azurerm_virtual_machine_scale_set` - support for specifying user assigned identities ([#1448](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1448))
* `azurerm_virtual_machine_scale_set` - making the `content` field in the `additional_unattend_config`  block (within `os_profile_windows_config`) sensitive ([#1471](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1471))
* `azurerm_virtual_network_gateway` - adding support for the `radius_server_address`, `radius_server_secret` and `vpn_client_protocols` fields to the Data Source ([#1505](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1505))

BUG FIXES:

* `azurerm_key_vault_key` - handling the parent Key Vault being deleted ([#1535](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1535))
* `azurerm_sql_database` - fix `requested_service_objective_name` updates ([#1503](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1503))
* `azurerm_storage_account` - limiting the `tags` field to 128 characters to match the service ([#1524](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1524))
* `azurerm_virtual_network_gateway` - fix `azurerm_virtual_network_gateway` crashing when `vpn_client_configuration` was not supplied ([#1505](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1505))

## 1.8.0 (June 28, 2018)

FEATURES:

* **New Resource:** `azurerm_dns_caa_record` support ([#1450](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1450))
* **New Resource:** `azurerm_virtual_machine_data_disk_attachment` ([#1207](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1207))

IMPROVEMENTS:

* dependencies: upgrading to v10.11.4 of `Azure/go-autorest` ([#1418](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1418))
* dependencies: upgrading to v17.4.0 of `Azure/azure-sdk-for-go` ([#1418](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1418))
* `azurerm_lb` - additional validation on properties ([#1403](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1403))
* `azurerm_application_gateway` - support for the `match` block for Probes ([#1446](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1446))
* `azurerm_log_analytics_solution` - support for Sovereign Clouds ([#1410](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1410))
* `azurerm_log_analytics_workspace` - support for Sovereign Clouds ([#1410](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1410))
* `azurerm_log_analytics_workspace` - support for the `PerGB2018` SKU ([#1079](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1079))
* `azurerm_mysql_server` -  `GeneralPurpose` and `MemoryOptimized` sku tiers now allow 4tb for the `storage_mb` property ([#1449](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1449))
* `azurerm_network_interface` - additional validation on properties ([#1403](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1403))
* `azurerm_postgresql_server` -  `GeneralPurpose` and `MemoryOptimized` sku tiers now allow 4tb for the `storage_mb` property ([#1449](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1449))
* `azurerm_postgresql_server` - adding support for version 10.0 ([#1457](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1457))
* `azurerm_route_table` - adding the  disable BGP propagation property ([#1435](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1435))
* `azurerm_sql_database` - support for importing from a bacpac backup ([#972](https://github.com/terraform-providers/terraform-provider-azurerm/issues/972))
* `azurerm_virtual_machine` - support for setting the TimeZone on Windows ([#1265](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1265))

BUG FIXES:

* validation: ensuring IPv4/MAC addresses are detected correctly ([#1431](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1431))

## 1.7.0 (June 16, 2018)

UPGRADE NOTES:

~> **Please Note:** The field `overprovision` on the `azurerm_virtual_machine_scale_set` resource has changed from `false` to `true` to match the behaviour of Azure in this release. ([#1322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1322))

BUG FIXES:

* `azurerm_key_vault` - respecting the proxy environment varibles terraform does and now can create vaults when behind a proxy ([#1393](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1393))
* `azurerm_kubernetes_cluster` - `dns_prefix` is now required ([#1333](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1333))
* `azurerm_network_interface` - ensuring that Public IP's/Private IP Addresses can be removed once assigned ([#1295](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1295))
* `azurerm_public_ip` - setting the `domain_name_label` property into state ([#1287](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1287))
* `azurerm_storage_account` - file and blob encryption is now explicity `true` by default ([#1380](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1380))
* `azurerm_servicebus_namespace` - the `capacity` propety no longer unnecessarily forces a new resource when changed ([#1382](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1382))
* `azurerm_virtual_machine_scale_set` - the field `overprovision` is now `true` by default ([#1322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1322))
* `azurerm_app_service_plan` - the `name` property validation now allows understores ([#1351](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1351))

IMPROVEMENTS:

* `azurerm_automation_schedule` - adding the `interval` property and supporting recurring schedules ([#1384](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1384))
* `azurerm_dns_ns_record` - deprecated `record` properties in favor of a `records` list ([#991](https://github.com/terraform-providers/terraform-provider-azurerm/issues/991))
* `azurerm_function_app` - adding the `identity` property ([#1369](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1369))
* `azurerm_role_definition` - the `role_definition_id` property is now optional. The resource will now generate a random UUID if it is ommited ([#1378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1378))
* `azurerm_storage_account` - adding the `network_rules` property ([#1334](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1334))
* `azurerm_storage_account` - adding the `identity` property ([#1323](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1323))
* `azurerm_storage_blob` - adding the `content_type` property ([#1304](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1304))
* `azurerm_virtual_machine` - support for `write_accelerator_enabled` property on Premium disks attached to MS-series machines ([#964](https://github.com/terraform-providers/terraform-provider-azurerm/issues/964))
* `azurerm_virtual_machine_scale_set` - adding the `dns_settings` and `dns_servers` property ([#1209](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1209))
* `azurerm_virtual_machine_scale_set` - adding the `ip_forwarding` property ([#1209](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1209))
* `azurerm_virtual_network_gateway` - adding the properties `vpn_client_protocols`, `radius_server_address` and `radius_server_secret` ([#946](https://github.com/terraform-providers/terraform-provider-azurerm/issues/946))
* dependencies: migrating to the un-deprecated Preview's for Container Instance, EventGrid, Log Analytics and SQL ([#1322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1322))
* dependencies: upgrading to `2018-01-01` of the EventGrid API ([#1322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1322))
* dependencies: upgrading to `2018-03-01` of the Monitor API ([#1322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1322))

## 1.6.0 (May 24, 2018)

UPGRADE NOTES:

~> **Please Note:** The `azurerm_mysql_server` resource has been updated from the Preview API's to the GA API's - which requires code changes in your Terraform Configuration to use the new Pricing SKU's. Upon updating to v1.6.0 - you'll need to update the configuration from the Preview SKU's to the GA SKU's.

~> **Please Note:** The `azurerm_postgresql_server` resource has been updated from the Preview API's to the GA API's - which requires code changes in your Terraform Configuration to use the new Pricing SKU's. Upon updating to v1.6.0 - you'll need to update the configuration from the Preview SKU's to the GA SKU's.

* `azurerm_scheduler_job_collection` - the property `max_retry_interval` on both the resource and datasource has been deprecated in favour of `max_recurrence_interval` to better match Azure ([#1218](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1218))

FEATURES:

* **New Data Source:** `azurerm_storage_account_sas` ([#1011](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1011))
* **New Resource:** `azurerm_data_lake_store` ([#1219](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1219))
* **New Resource:** `azurerm_relay_namespace` ([#1233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1233))

BUG FIXES:

* across data-sources and resources: making Connection Strings, Keys and Passwords sensitive fields ([#1242](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1242))
* `azurerm_virtual_machine_scale_set` - an empty `os_profile_windows_config` block no longer causes a panic ([#1224](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1224))

IMPROVEMENTS:

* authorization: upgrading to API version `2018-01-01-preview`
* `azurerm_app_service` - adding support for `ip_restriction`'s ([#1231](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1231))
* `azurerm_app_service_slot` - adding support for `ip_restriction`'s ([#1246](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1246))
* `azurerm_container_registry` - no longer forces a new resource on SKU change ([#1264](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1264))
* `azurerm_dns_zone` - datasource's `resource_group` field is now optional ([#1180](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1180))
* `azurerm_mysql_database` - ignoring casing for the `charset` field ([#1281](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1281))
* `azurerm_mysql_server` - support for the new GA Pricing SKU's ([#1154](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1154))
* `azurerm_postgresql_database` - ignoring the casing on the `collation` field ([#1255](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1255))
* `azurerm_postgresql_server` - support for the new GA Pricing SKU's ([#1190](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1190))
* `azurerm_public_ip` - computed values now default to an empy string ([#1247](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1247))
* `azurerm_role_assignment` - support for roles containing DataActions ([#1284](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1284))
* `azurerm_servicebus_queue` - adding `dead_lettering_on_message_expiration` ([#1235](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1235))
* `azurerm_virtual_machine_scale_set` - adding the `licence_type` property ([#1245](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1245))
* `azurerm_virtual_machine_scale_set` - adding the `priority` property ([#1250](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1250))

## 1.5.0 (May 14, 2018)

UPGRADE NOTES:

~> **Please Note:** Prior to v1.5 Data Sources in the AzureRM Provider returned `nil` rather than an error message when a Resource didn't exist, which was a bug. In order to bring this into line with other Providers - starting in v1.5 the AzureRM Provider will return an error message when a resource doesn't exist.

~> **Please Note:** This release fixes a bug in the `azurerm_redis_cache` resource where changes to fields weren't detected; as such you may see changes in the `redis_configuration` block, particularly with the `rdb_storage_connection_string` field. There's a bug tracking this inconsistency in [the Azure Rest API Specs Repository](https://github.com/Azure/azure-rest-api-specs/issues/3037).

FEATURES:

* **New Data Source:** `azurerm_cosmosdb_account` ([#1056](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1056))
* **New Data Source:** `azurerm_kubernetes_cluster` ([#1204](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1204))
* **New Data Source:** `azurerm_key_vault` ([#1202](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1202))
* **New Data Source:** `azurerm_key_vault_secret` ([#1202](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1202))
* **New Data Source:** `azurerm_route_table` ([#1203](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1203))

BUG FIXES:

* `azurerm_redis_cache` - changes to the `redis_configuration` block are now detected - please see the note above for more information ([#1211](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1211))

IMPROVEMENTS:

* dependencies - upgrading to v16.2.1 of `Azure/azure-sdk-for-go` ([#1198](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1198))
* dependencies - upgrading to v10.8.1 of `Azure/go-autorest` ([#1198](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1198))
* `azurerm_app_service` - support for HTTP2 ([#1188](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1188))
* `azurerm_app_service` - support for Managed Service Identity ([#1130](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1130))
* `azurerm_app_service_slot` - support for HTTP2 ([#1205](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1205))
* `azurerm_cosmosdb_account` - added support for the `connection_strings` property ([#1194](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1194))
* `azurerm_key_vault_certificate` - exposing the `certificate_data` ([#1200](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1200))
* `azurerm_kubernetes_cluster` - making `kube_config_raw` a sensitive field ([#1225](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1225))
* `azurerm_redis_cache` - Redis Caches can now be Imported ([#1211](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1211))
* `azurerm_redis_firewall_rule` - Redis Firewall Rules can now be Imported ([#1211](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1211))
* `azurerm_virtual_network` - guarding against nil-objects in the response ([#1208](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1208))
* `azurerm_virtual_network_gateway` - ignoring the case of the `GatewaySubnet` ([#1141](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1141))

## 1.4.0 (April 26, 2018)

UPGRADE NOTES:

* `azurerm_cosmosdb_account` - the field `failover_policy` has been deprecated in favour of `geo_locations` to better match Azure

FEATURES:

* **New Data Source:** `azurerm_recovery_services_vault` ([#995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/995))
* **New Resource:** `azurerm_recovery_services_vault` ([#995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/995))
* **New Resource:** `azurerm_servicebus_subscription_rule` ([#1124](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1124))

IMPROVEMENTS:

* `azurerm_app_service` - support for updating in-place ([#1125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1125))
* `azurerm_app_service_plan` - support for `kind` being `app` ([#1156](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1156))
* `azurerm_cosmosdb_account` - support for `enable_automatic_failover` ([#1055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1055))
* `azurerm_cosmosdb_account` - support for the `ConsistentPrefix` consistncy level ([#1055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1055))
* `azurerm_cosmosdb_account` - `prefixes` can now be configured for locations ([#1055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1055))
* `azurerm_function_app` - support for updating in-place ([#1125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1125))
* `azurerm_key_vault` - adding cert permissions for `Purge` and `Recover` ([#1132](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1132))
* `azurerm_key_vault` - polling to ensure the Key Vault is resolvable via DNS ([#1081](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1081)) ([#1164](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1164))
* `azurerm_kubernetes_cluster` - only setting the Subnet ID when it's not an empty string ([#1158](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1158))
* `azurerm_kubernetes_cluster` - exposing the clusters credentials as `kube_config` ([#953](https://github.com/terraform-providers/terraform-provider-azurerm/issues/953))
* `azurerm_metric_alertrule` - filtering out tags prefixed with `$type` ([#1107](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1107))
* `azurerm_virtual_machine` - loading managed disk information from Azure when the machine is stopped ([#1100](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1100))
* `azurerm_virtual_machine` - make the `vm_size` property case insensitive ([#1131](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1131))

BUG FIXES:

* `azurerm_cosmosdb_account` - locations can now be modified in-place (without requiring multiple apply's) ([#1055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1055))


## 1.3.3 (April 17, 2018)

FEATURES:

* **New Data Source:** `azurerm_app_service` ([#1071](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1071))
* **New Resource:** `azurerm_app_service_custom_hostname_binding` ([#1087](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1087))

IMPROVEMENTS:

* dependencies: upgrading to `v15.1.0` of `Azure/azure-sdk-for-go` ([#1099](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1099))
* dependencies: upgrading to `v10.6.0` of `Azure/go-autorest` ([#1077](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1077))
* `azurerm_app_service` - added support for the `https_only` field ([#1080](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1080))
* `azurerm_app_service_slot` - added support for the `https_only` field ([#1080](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1080))
* `azurerm_function_app` - added support for the `https_only` field ([#1080](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1080))
* `azurerm_key_vault_certificate` - exposing the certificate's associated `secret_id` ([#1096](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1096))
* `azurerm_redis_cache` - support for clusters on the internal network ([#1086](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1086))
* `azurerm_servicebus_queue` - support for setting `requires_session` ([#1111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1111))
* `azurerm_sql_database` - changes to `collation` force a new resource ([#1066](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1066))

## 1.3.2 (April 04, 2018)

FEATURES:

* **New Resource:** `azurerm_packet_capture` ([#1044](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1044))
* **New Resource:** `azurerm_policy_assignment` ([#1051](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1051))

IMPROVEMENTS:

* `azurerm_virtual_machine_scale_set` - adds support for MSI ([#1018](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1018))

## 1.3.1 (March 29, 2018)

FEATURES:

* **New Data Source:** `azurerm_scheduler_job_collection` ([#990](https://github.com/terraform-providers/terraform-provider-azurerm/issues/990))
* **New Data Source:** `azurerm_traffic_manager_geographical_location` ([#987](https://github.com/terraform-providers/terraform-provider-azurerm/issues/987))
* **New Resource:** `azurerm_express_route_circuit_authorization` ([#992](https://github.com/terraform-providers/terraform-provider-azurerm/issues/992))
* **New Resource:** `azurerm_express_route_circuit_peering` ([#1033](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1033))
* **New Resource:** `azurerm_iothub` ([#887](https://github.com/terraform-providers/terraform-provider-azurerm/issues/887))
* **New Resource:** `azurerm_policy_definition` ([#1010](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1010))
* **New Resource:** `azurerm_sql_virtual_network_rule` ([#978](https://github.com/terraform-providers/terraform-provider-azurerm/issues/978))

IMPROVEMENTS:

* `azurerm_app_service` - allow changing `client_affinity_enabled` without requiring a resource recreation ([#993](https://github.com/terraform-providers/terraform-provider-azurerm/issues/993))
* `azurerm_app_service` - support for configuring `LocalSCM` source control ([#826](https://github.com/terraform-providers/terraform-provider-azurerm/issues/826))
* `azurerm_app_service` - returning a clearer error message when the name (which needs to be globally unique) is in use ([#1037](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1037))
* `azurerm_cosmosdb_account` - increasing the maximum value for `max_interval_in_seconds` from 100s to 86400s (1 day) ([#1000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1000))
* `azurerm_function_app` - returning a clearer error message when the name (which needs to be globally unique) is in use ([#1037](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1037))
* `azurerm_network_interface` - support for attaching to Application Gateways ([#1027](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1027))
* `azurerm_traffic_manager_endpoint` - adding support for `geo_mappings` ([#986](https://github.com/terraform-providers/terraform-provider-azurerm/issues/986))
* `azurerm_traffic_manager_profile` - adding support for the `traffic_routing_method` `Geographic` ([#986](https://github.com/terraform-providers/terraform-provider-azurerm/issues/986))
* `azurerm_virtual_machine_scale_sets` - support for attaching to Application Gateways ([#1027](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1027))
* `azurerm_virtual_network_gateway` - changes to `peering_address` now force a new resource ([#1040](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1040))

## 1.3.0 (March 15, 2018)

FEATURES:

* **New Data Source:** `azurerm_cdn_profile` ([#950](https://github.com/terraform-providers/terraform-provider-azurerm/issues/950))
* **New Data Source:** `azurerm_network_interface` ([#854](https://github.com/terraform-providers/terraform-provider-azurerm/issues/854))
* **New Data Source:** `azurerm_public_ips` ([#304](https://github.com/terraform-providers/terraform-provider-azurerm/issues/304))
* **New Data Source:** `azurerm_subscriptions` ([#940](https://github.com/terraform-providers/terraform-provider-azurerm/issues/940))
* **New Resource:** `azurerm_log_analytics_solution` ([#952](https://github.com/terraform-providers/terraform-provider-azurerm/issues/952))
* **New Resource:** `azurerm_sql_active_directory_administrator` ([#765](https://github.com/terraform-providers/terraform-provider-azurerm/issues/765))
* **New Resource:** `azurerm_scheduler_job_collection` ([#963](https://github.com/terraform-providers/terraform-provider-azurerm/issues/963))

BUG FIXES:

* `azurerm_application_gateway` - fixes a crash where `ssl_policy` isn't returned from the Azure API when importing existing resources ([#935](https://github.com/terraform-providers/terraform-provider-azurerm/issues/935))
* `azurerm_app_service` - supporting `client_affinity_enabled` being `false` ([#973](https://github.com/terraform-providers/terraform-provider-azurerm/issues/973))
* `azurerm_kubernetes_cluster` - exporting the FQDN ([#907](https://github.com/terraform-providers/terraform-provider-azurerm/issues/907))
* `azurerm_sql_elasticpool` - fixing a crash where `location` isn't returned for legacy resources ([#982](https://github.com/terraform-providers/terraform-provider-azurerm/issues/982))

IMPROVEMENTS:

* Data Source: `azurerm_builtin_role_definition` - loading available role definitions from Azure ([#770](https://github.com/terraform-providers/terraform-provider-azurerm/issues/770))
* Data Source: `azurerm_managed_disk` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))
* Data Source: `azurerm_network_security_group` - support for security rules including Application Security Groups ([#925](https://github.com/terraform-providers/terraform-provider-azurerm/issues/925))
* `azurerm_app_service_plan` -  support for provisioning Consumption Plans ([#981](https://github.com/terraform-providers/terraform-provider-azurerm/issues/981))
* `azurerm_cdn_endpoint` - adding support for GeoFilters, ProbePaths ([#967](https://github.com/terraform-providers/terraform-provider-azurerm/issues/967))
* `azurerm_cdn_endpoint` - making the `origin` block ForceNew to match Azure ([#967](https://github.com/terraform-providers/terraform-provider-azurerm/issues/967))
* `azurerm_function_app` - adding `client_affinity_enabled`, `use_32_bit_worker_process` and `websockets_enabled` ([#886](https://github.com/terraform-providers/terraform-provider-azurerm/issues/886))
* `azurerm_load_balancer` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))
* `azurerm_managed_disk` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))
* `azurerm_network_interface` - setting `internal_fqdn` if it's not nil ([#977](https://github.com/terraform-providers/terraform-provider-azurerm/issues/977))
* `azurerm_network_security_group` - support for security rules including Application Security Groups ([#925](https://github.com/terraform-providers/terraform-provider-azurerm/issues/925))
* `azurerm_network_security_rule` - support for security rules including Application Security Groups ([#925](https://github.com/terraform-providers/terraform-provider-azurerm/issues/925))
* `azurerm_public_ip` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))
* `azurerm_redis_cache` - add support for `notify-keyspace-events` ([#949](https://github.com/terraform-providers/terraform-provider-azurerm/issues/949))
* `azurerm_template_deployment` - support for specifying parameters via `parameters_body` ([#404](https://github.com/terraform-providers/terraform-provider-azurerm/issues/404))
* `azurerm_virtual_machine` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))
* `azurerm_virtual_machine_scale_set` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))

## 1.2.0 (March 02, 2018)

FEATURES:

* **New Data Source:** `azurerm_application_security_group` ([#914](https://github.com/terraform-providers/terraform-provider-azurerm/issues/914))
* **New Resource:** `azurerm_application_security_group` ([#905](https://github.com/terraform-providers/terraform-provider-azurerm/issues/905))
* **New Resource:** `azurerm_servicebus_topic_authorization_rule` ([#736](https://github.com/terraform-providers/terraform-provider-azurerm/issues/736))

BUG FIXES:

* `azurerm_kubernetes_cluster` - an empty `linux_profile.ssh_key.keydata` no longer causes a crash ([#903](https://github.com/terraform-providers/terraform-provider-azurerm/issues/903))
* `azurerm_kubernetes_cluster` - the `linux_profile.admin_username` and `linux_profile.ssh_key.keydata` fields now force a new resource ([#895](https://github.com/terraform-providers/terraform-provider-azurerm/issues/895))
* `azurerm_network_interface` - the `subnet_id` field is now case insensitive ([#866](https://github.com/terraform-providers/terraform-provider-azurerm/issues/866))
* `azurerm_network_security_group` - reverting `security_rules` to a set to fix an ordering issue ([#893](https://github.com/terraform-providers/terraform-provider-azurerm/issues/893))
* `azurerm_virtual_machine_scale_set` - the `computer_name_prefix` field now forces a new resource ([#871](https://github.com/terraform-providers/terraform-provider-azurerm/issues/871))

IMPROVEMENTS:

* authentication: adding support for Managed Service Identity ([#639](https://github.com/terraform-providers/terraform-provider-azurerm/issues/639))
* `azurerm_container_group` - added `dns_name_label` and `FQDN` properties ([#877](https://github.com/terraform-providers/terraform-provider-azurerm/issues/877))
* `azurerm_network_interface` - support for attaching to Application Security Groups ([#911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/911))
* `azurerm_network_security_group` - support for augmented security rules ([#781](https://github.com/terraform-providers/terraform-provider-azurerm/issues/781))
* `azurerm_servicebus_subscription` - added support for the `forward_to` property ([#861](https://github.com/terraform-providers/terraform-provider-azurerm/issues/861))
* `azurerm_storage_account` - adding support for `account_kind` being `StorageV2` ([#851](https://github.com/terraform-providers/terraform-provider-azurerm/issues/851))
* `azurerm_virtual_network_gateway_connection` - support for IPsec/IKE Policies ([#834](https://github.com/terraform-providers/terraform-provider-azurerm/issues/834))

## 1.1.2 (February 19, 2018)

FEATURES:

* **New Resource:** `azurerm_kubernetes_cluster` ([#693](https://github.com/terraform-providers/terraform-provider-azurerm/issues/693))
* **New Resource:** `azurerm_app_service_active_slot` ([#818](https://github.com/terraform-providers/terraform-provider-azurerm/issues/818))
* **New Resource:** `azurerm_app_service_slot` ([#818](https://github.com/terraform-providers/terraform-provider-azurerm/issues/818))

BUG FIXES:

* **Data Source:** `azurerm_app_service_plan`: handling a 404 not being returned as an error ([#849](https://github.com/terraform-providers/terraform-provider-azurerm/issues/849))
* **Data Source:** `azurerm_virtual_network` - Fixing a crash when the DhcpOptions aren't specified ([#803](https://github.com/terraform-providers/terraform-provider-azurerm/issues/803))
* `azurerm_application_gateway` - fixing crashes due to schema mismatches for existing resources ([#848](https://github.com/terraform-providers/terraform-provider-azurerm/issues/848))
* `azurerm_storage_container` - add a retry for creation ([#846](https://github.com/terraform-providers/terraform-provider-azurerm/issues/846))

IMPROVEMENTS:

* authentication: pulling the `Environment` key from the Azure CLI Config ([#842](https://github.com/terraform-providers/terraform-provider-azurerm/issues/842))
* core: upgrading to `v12.5.0-beta` of the Azure SDK for Go ([#830](https://github.com/terraform-providers/terraform-provider-azurerm/issues/830))
* compute: upgrading to use the `2017-12-01` API Version ([#797](https://github.com/terraform-providers/terraform-provider-azurerm/issues/797))
* `azurerm_app_service_plan`: support for attaching to an App Service Environment ([#850](https://github.com/terraform-providers/terraform-provider-azurerm/issues/850))
* `azurerm_container_group` - adding `restart_policy` ([#827](https://github.com/terraform-providers/terraform-provider-azurerm/issues/827))
* `azurerm_managed_disk` - updated the validation on `disk_size_gb` / made it computed ([#800](https://github.com/terraform-providers/terraform-provider-azurerm/issues/800))
* `azurerm_role_assignment` - add `role_definition_name` ([#775](https://github.com/terraform-providers/terraform-provider-azurerm/issues/775))
* `azurerm_subnet` - add support for Service Endpoints ([#786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/786))
* `azurerm_virtual_machine` - changing `managed_disk_id` and `create_option` to be not ForceNew ([#813](https://github.com/terraform-providers/terraform-provider-azurerm/issues/813))


## 1.1.1 (February 06, 2018)

BUG FIXES:

* `azurerm_public_ip` - Setting the `ip_address` field regardless of the DNS Settings ([#772](https://github.com/terraform-providers/terraform-provider-azurerm/issues/772))
* `azurerm_virtual_machine` - ignores the case of the Managed Data Disk ID's to work around an Azure Portal bug ([#792](https://github.com/terraform-providers/terraform-provider-azurerm/issues/792))

FEATURES:

* **New Data Source:** `azurerm_storage_account` ([#794](https://github.com/terraform-providers/terraform-provider-azurerm/issues/794))
* **New Data Source:** `azurerm_virtual_network_gateway` ([#796](https://github.com/terraform-providers/terraform-provider-azurerm/issues/796))

## 1.1.0 (January 26, 2018)

UPGRADE NOTES:

* Data Source: `azurerm_builtin_role_definition` - now returns the correct UUID/GUID for the `Virtual Machines Contributor` role (previously the ID for the `Classic Virtual Machine Contributor` role was returned) ([#762](https://github.com/terraform-providers/terraform-provider-azurerm/issues/762))
* `azurerm_snapshot` - `source_uri` now forces a new resource on changes due to behavioural changes in the Azure API ([#744](https://github.com/terraform-providers/terraform-provider-azurerm/issues/744))

FEATURES:

* **New Data Source:** `azurerm_dns_zone` ([#702](https://github.com/terraform-providers/terraform-provider-azurerm/issues/702))
* **New Resource:** `azurerm_metric_alertrule` ([#478](https://github.com/terraform-providers/terraform-provider-azurerm/issues/478))
* **New Resource:** `azurerm_virtual_network_gateway` ([#133](https://github.com/terraform-providers/terraform-provider-azurerm/issues/133))
* **New Resource:** `azurerm_virtual_network_gateway_connection` ([#133](https://github.com/terraform-providers/terraform-provider-azurerm/issues/133))

IMPROVEMENTS:

* core: upgrading to `v12.2.0-beta` of `Azure/azure-sdk-for-go` ([#684](https://github.com/terraform-providers/terraform-provider-azurerm/issues/684))
* core: upgrading to `v9.7.0` of `Azure/go-autorest` ([#684](https://github.com/terraform-providers/terraform-provider-azurerm/issues/684))
* Data Source: `azurerm_builtin_role_definition` - adding extra role definitions ([#762](https://github.com/terraform-providers/terraform-provider-azurerm/issues/762))
* `azurerm_app_service` - exposing the `outbound_ip_addresses` field ([#700](https://github.com/terraform-providers/terraform-provider-azurerm/issues/700))
* `azurerm_function_app` - exposing the `outbound_ip_addresses` field ([#706](https://github.com/terraform-providers/terraform-provider-azurerm/issues/706))
* `azurerm_function_app` - add support for the `always_on` and `connection_string` fields ([#695](https://github.com/terraform-providers/terraform-provider-azurerm/issues/695))
* `azurerm_image` - add support for filtering images by a regex on the name ([#642](https://github.com/terraform-providers/terraform-provider-azurerm/issues/642))
* `azurerm_lb` - adding support for the `Standard` SKU (in Preview) ([#665](https://github.com/terraform-providers/terraform-provider-azurerm/issues/665))
* `azurerm_public_ip` - adding support for the `Standard` SKU (in Preview) ([#665](https://github.com/terraform-providers/terraform-provider-azurerm/issues/665))
* `azurerm_network_security_rule` - add support for augmented security rules ([#692](https://github.com/terraform-providers/terraform-provider-azurerm/issues/692))
* `azurerm_role_assignment` - generating a name if one isn't specified ([#685](https://github.com/terraform-providers/terraform-provider-azurerm/issues/685))
* `azurerm_traffic_manager_profile` - adding support for setting `protocol` to `TCP` ([#742](https://github.com/terraform-providers/terraform-provider-azurerm/issues/742))

## 1.0.1 (January 12, 2018)

FEATURES:

* **New Data Source:** `azurerm_app_service_plan` ([#668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/668))
* **New Data Source:** `azurerm_eventhub_namespace` ([#673](https://github.com/terraform-providers/terraform-provider-azurerm/issues/673))
* **New Resource:** `azurerm_function_app` ([#647](https://github.com/terraform-providers/terraform-provider-azurerm/issues/647))

IMPROVEMENTS:

* core: adding a cache to the Storage Account Keys ([#634](https://github.com/terraform-providers/terraform-provider-azurerm/issues/634))
* `azurerm_eventhub` - added support for `capture_description` ([#681](https://github.com/terraform-providers/terraform-provider-azurerm/issues/681))
* `azurerm_eventhub_consumer_group` - adding validation for the user metadata field ([#641](https://github.com/terraform-providers/terraform-provider-azurerm/issues/641))
* `azurerm_lb` - adding the computed field `public_ip_addresses` ([#633](https://github.com/terraform-providers/terraform-provider-azurerm/issues/633))
* `azurerm_local_network_gateway` - add support for `tags` ([#638](https://github.com/terraform-providers/terraform-provider-azurerm/issues/638))
* `azurerm_network_interface` - support for Accelerated Networking ([#672](https://github.com/terraform-providers/terraform-provider-azurerm/issues/672))
* `azurerm_storage_account` - expose `primary_connection_string` and `secondary_connection_string` ([#647](https://github.com/terraform-providers/terraform-provider-azurerm/issues/647))


## 1.0.0 (December 15, 2017)

FEATURES:

* **New Data Source:** `azurerm_network_security_group` ([#623](https://github.com/terraform-providers/terraform-provider-azurerm/issues/623))
* **New Data Source:** `azurerm_virtual_network` ([#533](https://github.com/terraform-providers/terraform-provider-azurerm/issues/533))
* **New Resource:** `azurerm_management_lock` ([#575](https://github.com/terraform-providers/terraform-provider-azurerm/issues/575))
* **New Resource:** `azurerm_network_watcher` ([#571](https://github.com/terraform-providers/terraform-provider-azurerm/issues/571))

IMPROVEMENTS:

* authentication - add support for the latest Azure CLI configuration ([#573](https://github.com/terraform-providers/terraform-provider-azurerm/issues/573))
* authentication - conditional loading of the Subscription ID / Tenant ID / Environment ([#574](https://github.com/terraform-providers/terraform-provider-azurerm/issues/574))
* core - appending additions to the User Agent, so we don't overwrite the Go SDK User Agent info ([#587](https://github.com/terraform-providers/terraform-provider-azurerm/issues/587))
* core - Upgrading `Azure/azure-sdk-for-go` to v11.2.2-beta ([#594](https://github.com/terraform-providers/terraform-provider-azurerm/issues/594))
* core - upgrading `Azure/go-autorest` to v9.5.2 ([#617](https://github.com/terraform-providers/terraform-provider-azurerm/issues/617))
* core - skipping Resource Provider Registration in AutoRest when opted-out ([#630](https://github.com/terraform-providers/terraform-provider-azurerm/issues/630))
* `azurerm_app_service` - exposing the Default Hostname as a Computed field

## 0.3.3 (November 14, 2017)

FEATURES:

* **New Resource:** `azurerm_redis_firewall_rule` ([#529](https://github.com/terraform-providers/terraform-provider-azurerm/issues/529))

IMPROVEMENTS:

* authentication: allow using multiple subscriptions for Azure CLI auth ([#445](https://github.com/terraform-providers/terraform-provider-azurerm/issues/445))
* core: appending the CloudShell version to the user agent when running within CloudShell ([#483](https://github.com/terraform-providers/terraform-provider-azurerm/issues/483))
* `azurerm_app_service` / `azurerm_app_service_plan` - adding validation for the `name` fields ([#528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/528))
* `azurerm_container_registry` - Migration: Fixing a crash when the storage_account block is nil ([#551](https://github.com/terraform-providers/terraform-provider-azurerm/issues/551))
* `azurerm_lb_nat_rule`: support for floating IP's ([#542](https://github.com/terraform-providers/terraform-provider-azurerm/issues/542))
* `azurerm_public_ip` - Clarify the error message for the validation of domain name label ([#485](https://github.com/terraform-providers/terraform-provider-azurerm/issues/485))
* `azurerm_network_security_group` - fixing a crash when changes were made outside of Terraform ([#492](https://github.com/terraform-providers/terraform-provider-azurerm/issues/492))
* `azurerm_redis_cache`: support for Patch Schedules ([#540](https://github.com/terraform-providers/terraform-provider-azurerm/issues/540))
* `azurerm_virtual_machine` - ensuring `vhd_uri` is validated ([#470](https://github.com/terraform-providers/terraform-provider-azurerm/issues/470))
* `azurerm_virtual_machine_scale_set`: fixing a crash where accelerated networking isn't returned by the API ([#480](https://github.com/terraform-providers/terraform-provider-azurerm/issues/480))

## 0.3.2 (October 30, 2017)

FEATURES: 

* **New Resource:** `azurerm_application_gateway` ([#413](https://github.com/terraform-providers/terraform-provider-azurerm/issues/413))

IMPROVEMENTS: 

  - `azurerm_virtual_machine_scale_set` - Add nil check to os disk ([#436](https://github.com/terraform-providers/terraform-provider-azurerm/issues/436))

  - `azurerm_key_vault` - Increased timeout on dns availability ([#457](https://github.com/terraform-providers/terraform-provider-azurerm/issues/457))
  
  - `azurerm_route_table` - Fix issue when routes are computed ([#450](https://github.com/terraform-providers/terraform-provider-azurerm/issues/450))

## 0.3.1 (October 21, 2017)

IMPROVEMENTS:

  - `azurerm_virtual_machine_scale_set` - Updating this resource with the v11 of the Azure SDK for Go ([#448](https://github.com/terraform-providers/terraform-provider-azurerm/issues/448))

## 0.3.0 (October 17, 2017)

UPGRADE NOTES:

  - `azurerm_automation_account` - the SKU `Free` has been replaced with `Basic`.
  - `azurerm_container_registry` - Azure has updated the SKU from `Basic` to `Classic`, with new `Basic`, `Standard` and `Premium` SKU's introduced.
  - `azurerm_container_registry` - the `storage_account` block is now `storage_account_id` and is only required for `Classic` SKU's
  - `azurerm_key_vault` - `certificate_permissions`, `key_permissions` and `secret_permissions` have all had the `All` option removed by Azure. Each permission now needs to be specified manually.
  * `azurerm_route_table` - `route` is no longer computed
  - `azurerm_servicebus_namespace` - The `capacity` field can only be set for `Premium` SKU's
  - `azurerm_servicebus_queue` - The `enable_batched_operations` and `support_ordering` fields have been deprecated by Azure.
  - `azurerm_servicebus_subscription` - The `dead_lettering_on_filter_evaluation_exceptions` has been removed by Azure.
  - `azurerm_servicebus_topic` - The `enable_filtering_messages_before_publishing` field has been removed by Azure.

FEATURES:

* **New Data Source:** `azurerm_builtin_role_definition` ([#384](https://github.com/terraform-providers/terraform-provider-azurerm/issues/384))
* **New Data Source:** `azurerm_image` ([#382](https://github.com/terraform-providers/terraform-provider-azurerm/issues/382))
* **New Data Source:** `azurerm_key_vault_access_policy` ([#423](https://github.com/terraform-providers/terraform-provider-azurerm/issues/423))
* **New Data Source:** `azurerm_platform_image` ([#375](https://github.com/terraform-providers/terraform-provider-azurerm/issues/375))
* **New Data Source:** `azurerm_role_definition` ([#414](https://github.com/terraform-providers/terraform-provider-azurerm/issues/414))
* **New Data Source:** `azurerm_snapshot` ([#420](https://github.com/terraform-providers/terraform-provider-azurerm/issues/420))
* **New Data Source:** `azurerm_subnet` ([#411](https://github.com/terraform-providers/terraform-provider-azurerm/issues/411))
* **New Resource:** `azurerm_key_vault_certificate` ([#408](https://github.com/terraform-providers/terraform-provider-azurerm/issues/408))
* **New Resource:** `azurerm_role_assignment` ([#414](https://github.com/terraform-providers/terraform-provider-azurerm/issues/414))
* **New Resource:** `azurerm_role_definition` ([#414](https://github.com/terraform-providers/terraform-provider-azurerm/issues/414))
* **New Resource:** `azurerm_snapshot` ([#420](https://github.com/terraform-providers/terraform-provider-azurerm/issues/420))

IMPROVEMENTS:

* Upgrading to v11 of the Azure SDK for Go ([#367](https://github.com/terraform-providers/terraform-provider-azurerm/issues/367))
* `azurerm_client_config` - updating the data source to work when using AzureCLI auth ([#393](https://github.com/terraform-providers/terraform-provider-azurerm/issues/393))
* `azurerm_container_group` - add support for volume mounts ([#366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/366))
* `azurerm_key_vault` - fix a crash when no certificate_permissions are defined ([#374](https://github.com/terraform-providers/terraform-provider-azurerm/issues/374))
* `azurerm_key_vault` - waiting for the DNS to propagate ([#401](https://github.com/terraform-providers/terraform-provider-azurerm/issues/401))
* `azurerm_managed_disk` - support for creating Managed Disks from Platform Images by supporting "FromImage" ([#399](https://github.com/terraform-providers/terraform-provider-azurerm/issues/399))
* `azurerm_managed_disk` - support for creating Encrypted Managed Disks ([#399](https://github.com/terraform-providers/terraform-provider-azurerm/issues/399))
* `azurerm_mysql_*` - Ensuring we register the MySQL Resource Provider ([#397](https://github.com/terraform-providers/terraform-provider-azurerm/issues/397))
* `azurerm_network_interface` - exposing all of the Private IP Addresses assigned to the NIC ([#409](https://github.com/terraform-providers/terraform-provider-azurerm/issues/409))
* `azurerm_network_security_group` / `azurerm_network_security_rule` - refactoring ([#405](https://github.com/terraform-providers/terraform-provider-azurerm/issues/405))
* `azurerm_route_table` - removing routes when none are specified ([#403](https://github.com/terraform-providers/terraform-provider-azurerm/issues/403))
* `azurerm_route_table` - refactoring `route` from a Set to a List ([#402](https://github.com/terraform-providers/terraform-provider-azurerm/issues/402))
* `azurerm_route` - refactoring `route` from a Set to a List ([#402](https://github.com/terraform-providers/terraform-provider-azurerm/issues/402))
* `azurerm_storage_account` - support for File Encryption ([#363](https://github.com/terraform-providers/terraform-provider-azurerm/issues/363))
* `azurerm_storage_account` - support for Custom Domain ([#363](https://github.com/terraform-providers/terraform-provider-azurerm/issues/363))
* `azurerm_storage_account` - splitting the storage account Tier and Replication out into separate fields ([#363](https://github.com/terraform-providers/terraform-provider-azurerm/issues/363))
- `azurerm_storage_account` - returning a user friendly error when trying to provision a Blob Storage Account with ZRS redundancy ([#421](https://github.com/terraform-providers/terraform-provider-azurerm/issues/421))
* `azurerm_subnet` - making it possible to remove Network Security Groups / Route Tables ([#411](https://github.com/terraform-providers/terraform-provider-azurerm/issues/411))
* `azurerm_virtual_machine` - fixing a bug where `additional_unattend_config.content` was being updated unintentionally ([#377](https://github.com/terraform-providers/terraform-provider-azurerm/issues/377))
* `azurerm_virtual_machine` - switching to use Lists instead of Sets ([#426](https://github.com/terraform-providers/terraform-provider-azurerm/issues/426))
* `azurerm_virtual_machine_scale_set` - fixing a bug where `additional_unattend_config.content` was being updated unintentionally ([#377](https://github.com/terraform-providers/terraform-provider-azurerm/issues/377))
* `azurerm_virtual_machine_scale_set` - support for multiple network profiles ([#378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/378))

## 0.2.2 (September 28, 2017)

FEATURES:

* **New Resource:** `azurerm_key_vault_key` ([#356](https://github.com/terraform-providers/terraform-provider-azurerm/issues/356))
* **New Resource:** `azurerm_log_analytics_workspace` ([#331](https://github.com/terraform-providers/terraform-provider-azurerm/issues/331))
* **New Resource:** `azurerm_mysql_configuration` ([#352](https://github.com/terraform-providers/terraform-provider-azurerm/issues/352))
* **New Resource:** `azurerm_mysql_database` ([#352](https://github.com/terraform-providers/terraform-provider-azurerm/issues/352))
* **New Resource:** `azurerm_mysql_firewall_rule` ([#352](https://github.com/terraform-providers/terraform-provider-azurerm/issues/352))
* **New Resource:** `azurerm_mysql_server` ([#352](https://github.com/terraform-providers/terraform-provider-azurerm/issues/352))

IMPROVEMENTS:

* Updating the provider initialization & adding a `skip_credentials_validation` field to the provider for some advanced scenarios ([#322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/322))

## 0.2.1 (September 25, 2017)

FEATURES:

* **New Resource:** `azurerm_automation_account` ([#257](https://github.com/terraform-providers/terraform-provider-azurerm/issues/257))
* **New Resource:** `azurerm_automation_credential` ([#257](https://github.com/terraform-providers/terraform-provider-azurerm/issues/257))
* **New Resource:** `azurerm_automation_runbook` ([#257](https://github.com/terraform-providers/terraform-provider-azurerm/issues/257))
* **New Resource:** `azurerm_automation_schedule` ([#257](https://github.com/terraform-providers/terraform-provider-azurerm/issues/257))
* **New Resource:** `azurerm_app_service` ([#344](https://github.com/terraform-providers/terraform-provider-azurerm/issues/344))

IMPROVEMENTS:

* `azurerm_client_config` - adding `service_principal_application_id` ([#348](https://github.com/terraform-providers/terraform-provider-azurerm/issues/348))
* `azurerm_key_vault` - adding `application_id` and `certificate_permissions` ([#348](https://github.com/terraform-providers/terraform-provider-azurerm/issues/348))

BUG FIXES:

* `azurerm_virtual_machine_scale_set` - fix panic with `additional_unattend_config` block ([#266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/266))
