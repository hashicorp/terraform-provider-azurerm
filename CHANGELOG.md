## 1.22.0 (Unreleased)

FEATURES:

* **New Resource:**  `azurerm_ddos_protection_plan` [GH-2654]

IMPROVEMENTS:

* dependencies: switching to Go Modules [GH-2705]
* dependencies: upgrading to v11.3.2 of github.com/Azure/go-autorest [GH-2744]
* `azurerm_application_gateway` - support for the `http2` property [GH-2735]
* `azurerm_application_gateway` - support for the `file_upload_limit_mb` property [GH-2666]
* `azurerm_application_gateway` - Support for `pick_host_name_from_backend_address` and `pick_host_name_from_backend_http_settings` properties [GH-2658]
* `azurerm_cosmosdb_account` - support for the `EnableAggregationPipeline`, `MongoDBv3.4` and ` mongoEnableDocLevelTTL` capabilities [GH-2715]
* `azurerm_data_lake_store_file` - support file uploads greater then 4 megabytes [GH-2633]
* `azurerm_mssql_elasticpool` - support for setting `max_size_bytes` [GH-2346]
* `azurerm_signalr_service` - exporting `primary_access_key`, `secondary_access_key`, `primary_connection_string` and `secondary_connection_string` and secondary access keys and connection strings [GH-2655]
* `azurerm_subnet` - support for additional subnet delegation types [GH-2667]

BUG FIXES:

* `azurerm_azuread_application` - fixing a bug where `reply_uris` was set incorrectly [GH-2729]
* `azurerm_batch_pool` - can now set multiple environment variables [GH-2685]
* `azurerm_cosmosdb_account` - prevent occasional error when deleting the resource [GH-2702]
* `azurerm_cosmosdb_account` - allow empty values for the `ip_range_filter` property [GH-2713]
* `azurerm_express_route_circuit` - added the `premium` SKU back to validation logic [GH-2692]
* `azurerm_firewall` - ensuring rules aren't removed during an update [GH-2663]
* `azurerm_notification_hub_namespace` - now polls on creation to handle eventual consistency [GH-2701]
* `azurerm_redis_cache` - locking on the Virtual Network/Subnet name to avoid a race condition [GH-2725]
* `azurerm_service_bus_subscription` - name's can now start with a digit [GH-2672]
* `azurerm_security_center` - increase the creation timeout to `30m` [GH-2724]
* `azurerm_subnet` - fixing a crash when service endpoints was nil [GH-2742]

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
* `azurerm_automation_account` - exposing `dsc_server_endpoint`, `dsc_primary_access_key`, `dsc_secondary_access_key` properties [[#2166](https://github.com/terraform-providers/terraform-provider-azurerm/issues/2166)] 
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
* dependencies: upgrading to v10.15.4 of `github.com/Azure/go-autorest` ([#1861](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1861)] [[#1909](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1909))
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
* `azurerm_key_vault` - polling to ensure the Key Vault is resolvable via DNS ([#1081](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1081)] [[#1164](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1164))
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
