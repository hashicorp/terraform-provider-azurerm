## 2.30.0 (October 01, 2020)

UPGRADE NOTES

* This release renames certain fields within the `azurerm_cosmosdb_account` (data source & resource) and `azurerm_function_app_host_keys` data source to follow HashiCorp's [inclusive language guidelines](https://discuss.hashicorp.com/t/inclusive-language-changes) - where fields have been renamed, existing fields will continue to remain available until the next major version of the Azure Provider (`v3.0`)

FEATURES: 

* **New Data Source:** `azurerm_cosmosdb_sql_storedprocedure` ([#6189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6189))
* **New Data Source:** `azurerm_ip_groups` ([#8556](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8556))
* **New Resource:** `azurerm_ip_groups` ([#8556](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8556))
* **New Resource:** `azurerm_resource_group_template_deployment` ([#8672](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8672))
* **New Resource:** `azurerm_subscription_template_deployment` ([#8672](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8672))

IMPROVEMENTS:

* dependencies: updating `iothub` to `2020-03-01` ([#8688](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8688))
* dependencies: updating `storagecache` to `2020-03-01` ([#8078](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8078))
* dependencies: updating `resources` to API Version `2020-06-01` ([#8672](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8672))
* `azurerm_analysis_services_server` - support for the `S8v2` and `S9v2` SKU's ([#8707](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8707))
* `azurerm_cognitive_account` - support for the `S` `sku` ([#8639](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8639))
* `azurerm_container_group` - support for the `dns_config` block ([#7912](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7912))
* `azurerm_cosmosdb_account` - support the `zone_reduntant` property ([#8295](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8295))
* `azurerm_cosmosdb_mongo_collection` - will now respect the order of the `keys` property in the `index` block ([#8602](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8602))
* `azurerm_hpc_cache` -  support the `mtu` and `root_squash_enabled` properties ([#8078](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8078))
* `azurerm_key_vault` - add support for `enable_rbac_authorization` ([#8670](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8670))
* `azurerm_lighthouse_assignment` - limit the `scope` property to subsriptions ([#8601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8601))
* `azurerm_logic_app_workflow` - support for the `integration_service_environment_id` property ([#8504](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8504))
* `azurerm_servicebus_topic` - validate the `max_size_in_megabytes` property ([#8648](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8648))
* `azurerm_servicebus_queue` - validate the `max_size_in_megabytes` property ([#8648](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8648))
* `azurerm_servicebus_subscription_rule` - support the `correlation_filter.properties` property ([#8646](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8646))
* `azurerm_storage_management_policy` - support the `appendBlob` value for `blob_types` ([#8659](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8659))


BUG FIXES:

* `azurerm_monitor_metric_alert` - property wait when creating/updating multiple monitor metric alerts ([#8667](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8667))
* `azurerm_linux_virtual_machine_scale_set` - fix empty JSON error in `settings` and `protected_settings` when these values are not used ([#8627](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8627))

## 2.29.0 (September 24, 2020)

UPGRADE NOTES:

* `azurerm_api_management` - the value `None` has been removed from the `identity` block to match other resources, to specify an API Management Service with no Managed Identity remove the `identity` block ([#8411](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8411))
* `azurerm_container_registry` -  the `storage_account_id` property now forces a new resource as required by the updated API version ([#8477](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8477))
* `azurerm_virtual_hub_connection` - deprecating the field `vitual_network_to_hub_gateways_traffic_allowed` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_virtual_hub_connection` - deprecating the field `hub_to_vitual_network_traffic_allowed` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_virtual_wan` - deprecating the field `allow_vnet_to_vnet_traffic` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))

FEATURES: 

* **New Data Source:** `azurerm_data_share_dataset_kusto_cluster` ([#8464](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8464))
* **New Data Source:** `azurerm_databricks_workspace` ([#8502](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8502))
* **New Data Source:** `azurerm_firewall_policy` ([#7390](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7390))
* **New Data Source:** `azurerm_storage_sync_group` ([#8462](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8462))
* **New Data Source:** `azurerm_mssql_server` ([#7917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7917))
* **New Resource:** `azurerm_data_share_dataset_kusto_cluster` ([#8464](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8464))
* **New Resource:** `azurerm_firewall_policy` ([#7390](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7390))
* **New Resource:** `azurerm_mysql_server_key` ([#8125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8125))
* **New Resource:** `azurerm_postgresql_server_key` ([#8126](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8126))

IMPROVEMENTS:

* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v46.3.0` ([#8592](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8592))
* dependencies: updating `containerregistry` to `2019-05-01` ([#8477](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8477))
* Data Source: `azurerm_api_management` - export the `private_ip_addresses` attribute for primary and additional locations ([#8290](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8290))
* `azurerm_api_management` - support the `virtual_network_configuration` block for additional locations ([#8290](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8290))
* `azurerm_api_management` - export the `private_ip_addresses` attribute for additional locations ([#8290](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8290))
* `azurerm_cosmosdb_account` - support the `Serverless` value for the `capabilities` property ([#8533](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8533))
* `azurerm_cosmosdb_sql_container` - support for the `indexing_policy` property ([#8461](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8461))
* `azurerm_mssql_server` - support for the `recover_database_id` and `restore_dropped_database_id` properties ([#7917](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7917))
* `azurerm_policy_set_definition` - support for typed parameter values other then string in `the policy_definition_reference` block deprecating `parameters` in favour of `parameter_vcaluess` ([#8270](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8270))
* `azurerm_search_service` - Add support for `allowed_ips` ([#8557](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8557))
* `azurerm_service_fabric_cluster` - Remove two block limit for `client_certificate_thumbprint` ([#8521](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8521))
* `azurerm_signalr_service` - support for delta updates ([#8541](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8541))
* `azurerm_spring_cloud_service` - support for configuring the `network` block ([#8568](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8568))
* `azurerm_virtual_hub_connection` - deprecating the field `vitual_network_to_hub_gateways_traffic_allowed` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_virtual_hub_connection` - deprecating the field `hub_to_vitual_network_traffic_allowed` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_virtual_hub_connection` - switching to use the now separate API for provisioning these resources ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_virtual_wan` - deprecating the field `allow_vnet_to_vnet_traffic` since due to a breaking change in the API behaviour this is no longer used ([#7601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7601))
* `azurerm_windows_virtual_machine` - support for updating the `license_type` field ([#8542](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8542))

BUG FIXES:

* `azurerm_api_management` - the value `None` for the field `type` within the `identity` block has been removed - to remove a managed identity remove the `identity` block ([#8411](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8411))
* `azurerm_app_service` - don't try to manage source_control when scm_type is `VSTSRM` ([#8531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8531))
* `azurerm_function_app` - don't try to manage source_control when scm_type is `VSTSRM` ([#8531](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8531))
* `azurerm_kubernetes_cluster` - picking the first system node pool if the original `default_node_pool` has been removed ([#8503](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8503))

## 2.28.0 (September 17, 2020)

UPGRADE NOTES

* The `id` field for the `azurerm_role_definition` changed in release 2.27.0 to work around a bug in the Azure API when using management groups, where the Scope isn't returned - the existing `id` field is available as `role_definition_resource_id` from this version of the Azure Provider.

FEATURES:

* **New Data Source:** `azurerm_data_share_dataset_data_lake_gen2` [[#7907](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7907)] 
* **New Data Source:** `azurerm_servicebus_queue_authorization_rule` ([#8438](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8438))
* **New Data Source:** `azurerm_storage_sync` [[#7843](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7843)] 
* **New Resource:** `azurerm_data_share_dataset_data_lake_gen2` ([#7907](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7907))
* **New Resource:** `azurerm_lighthouse_definition` ([#6560](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6560))
* **New Resource:** `azurerm_lighthouse_assignment` ([#6560](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6560))
* **New Resource:** `azurerm_mssql_server_extended_auditing_policy`  ([#8447](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8447))
* **New Resource:** `azurerm_storage_sync` ([#7843](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7843))
* **New Resource:** `azurerm_synapse_sql_pool` ([#8095](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8095))

IMPROVEMENTS:

* Data Source: `azurerm_app_service_environment` - Expose vip information of an app service environment ([#8487](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8487))
* Data Source: `azurerm_function_app` - export the `identity` block ([#8389](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8389))
* `azurerm_app_service_hybrid_connection` - support relays in different namespaces ([#8370](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8370))
* `azurerm_cosmosdb_cassandra_keyspace` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_gremlin_database` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_gremlin_graph` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_mongo_collection` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_mongo_database` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_sql_container` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_sql_database` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_cosmosdb_table` - support the `autoscale_settings` block ([#7773](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7773))
* `azurerm_firewall` - support the `management_ip_configuration` block ([#8235](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8235))
* `azurerm_storage_account_customer_managed_key` - support for key rotation ([#7836](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7836))

BUG FIXES:

* Data Source: `azurerm_function_app_host_keys` - Fix a crash when null ID sometimes returned by API ([#8430](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8430))
* `azurerm_cognitive_account` - correctly wait on update logic ([#8386](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8386))
* `azurerm_eventhub_consumer_group` - allow the `name` property to be set to `$Default` ([#8388](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8388))
* `azurerm_kubernetes_cluster` - ensure the OMS Agent Log Analytics Workspace case is preserved after disabling/enabling ([#8374](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8374))
* `azurerm_management_group_id` - loosen case restritions during parsing of management group ID ([#8024](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8024))
* `azurerm_packet_capture` - fix to ID path to match change in API ([#8167](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8167))
* `azurerm_role_definition` - expose `role_definition_resource_id` ([#8492](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8492))

## 2.27.0 (September 10, 2020)

UPGRADE NOTES

* The `id` field for the `azurerm_role_definition` has changed in this release to work around a bug in the Azure API when using management groups, where the Scope isn't returned - the existing `id` field is available as `role_definition_resource_id` on the new resource from version 2.28.0 of the Azure Provider.

FEATURES:

* **New Data Source:** `azurerm_attestation_provider` ([#7885](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7885))
* **New Data Source:** `azurerm_function_app_host_keys` ([#7902](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7902))
* **New Data Source:** `azurerm_lb_rule` ([#8365](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8365))
* **New Resource:** `azurerm_mssql_database_extended_auditing_policy` ([#7793](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7793))
* **New Resource:** `azurerm_attestation_provider` ([#7885](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7885))
* **New Resource:** `azurerm_api_management_api_diagnostic` ([#7873](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7873))
* **New Resource:** `azurerm_data_factory_linked_service_azure_sql_database` ([#8349](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8349))

IMPROVEMENTS:

* Data Source: `azurerm_virtual_network_gateway` - exposing `aad_audience`, `aad_issuer` and `aad_tenant` within the `vpn_client_configuration` block ([#8294](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8294))
* `azurerm_cosmosdb_account` - supporting the value `AllowSelfServeUpgradeToMongo36` for the `name` field within the `capabilities` block ([#8335](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8335))
* `azurerm_linux_virtual_machine` - Add support for `encryption_at_host_enabled` ([#8322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8322))
* `azurerm_linux_virtual_machine_scale_set` - Add support for `encryption_at_host_enabled` ([#8322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8322))
* `azurerm_servicebus_subscription` - add support for `dead_lettering_on_filter_evaluation_error` ([#8412](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8412))
* `azurerm_spring_cloud_app` - support for the `identity` block ([#8336](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8336))
* `azurerm_storage_share_directory` - Update name validation ([#8366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8366))
* `azurerm_virtual_network_gateway` - support for `aad_audience`, `aad_issuer` and `aad_tenant` within the `vpn_client_configuration` block ([#8294](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8294))
* `azurerm_windows_virtual_machine` - Add support for `encryption_at_host_enabled` ([#8322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8322))
* `azurerm_windows_virtual_machine_scale_set` - Add support for `encryption_at_host_enabled` ([#8322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8322))

BUG FIXES:

* `azurerm_api_management_x.y.api_name` - validation fix ([#8409](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8409))
* `azurerm_application_insights_webtests` - Fix an issue where the `kind` property is sometimes set to `null` ([#8372](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8372))
* `azurerm_cognitive_account` - Fixes a crash when provisioning a QnAMaker and supports AnomalyDetector ([#8357](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8357))
* `azurerm_linux_virtual_machine` - Add WaitForState on VM delete ([#8383](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8383))
* `azurerm_network_security_group` - fixed issue where updates would fail for resource ([#8384](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8384))
* `azurerm_role_definition` - fixed delete operation when role is scoped to Management Group ([#6107](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6107))
* `azurerm_windows_virtual_machine` - Add WaitForState on VM delete ([#8383](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8383))

## 2.26.0 (September 04, 2020)

UPGRADE NOTES:

* **Opt-In Beta:** This release introduces an opt-in beta for in-line Virtual Machine Scale Set Extensions. This functionality enables the resource to be used with Azure Service Fabric and other extensions that may require creation time inclusion on Scale Set members. Please see the documentation for `azurerm_linux_virtual_machine_scale_set` and `azurerm_windows_virtual_machine_scale_set` for information.

FEATURES:

* **New Resource:** `azurerm_log_analytics_saved_search` ([#8253](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8253))

IMPROVEMENTS:

* dependencies: updating `loganalytics` to `2020-03-01-preview` ([#8234](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8234))
* `azurerm_api_management_subscription` - Support `allow_tracing property` ([#7969](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7969))
* `azurerm_application_gateway ` - Add support for `probe.properties.port` ([#8278](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8278))
* `azurerm_linux_virtual_machine_scale_set` - Beta support for `extension` blocks ([#8222](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8222))
* `azurerm_log_analytics_workspace`- the `sku` value is now optional and defaults to `PerGB2018` ([#8272](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8272))
* `azurerm_windows_virtual_machine_scale_set` - Beta support for `extension` blocks ([#8222](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8222))

BUG FIXES:

* `azurerm_cdn_endpoint` - fixing the casing of the Resource ID to be consistent ([#8237](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8237))
* `azurerm_cdn_profile` - fixing the casing of the Resource ID to be consistent ([#8237](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8237))
* `azurerm_key_vault_key` - updating the latest version of the key when updating metadata ([#8304](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8304))
* `azurerm_key_vault_secret` - updating the latest version of the secret when updating metadata ([#8304](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8304))
* `azurerm_linux_virtual_machine` - allow updating `allow_extension_operations` regardless of the value of `provision_vm_agent` (for when the VM Agent has been installed manually) ([#8001](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8001))
* `azurerm_linux_virtual_machine_scale_set` - working around a bug in the Azure API by always sending the existing Storage Image Reference during updates ([#7983](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7983))
* `azurerm_network_interface_application_gateway_association` - handling the Network Interface being deleted during a refresh ([#8267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8267))
* `azurerm_network_interface_application_security_group_association` - handling the Network Interface being deleted during a refresh ([#8267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8267))
* `azurerm_network_interface_backend_address_pool_association` - handling the Network Interface being deleted during a refresh ([#8267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8267))
* `azurerm_network_interface_nat_rule_association_resource` - handling the Network Interface being deleted during a refresh ([#8267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8267))
* `azurerm_network_interface_network_security_group_association` - handling the Network Interface being deleted during a refresh ([#8267](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8267))
* `azurerm_windows_virtual_machine` - allow updating `allow_extension_operations` regardless of the value of `provision_vm_agent` (for when the VM Agent has been installed manually) ([#8001](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8001))
* `azurerm_windows_virtual_machine_scale_set` - working around a bug in the Azure API by always sending the existing Storage Image Reference during updates ([#7983](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7983))

## 2.25.0 (August 27, 2020)

UPGRADE NOTES:

* `azurerm_container_group` - The `secure_environment_variables` field within the `container` now maps keys with empty values, which differs from previous versions of this provider which ignored empty values ([#8151](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8151))

FEATURES:

* **New Resource** `azurerm_spring_cloud_certificate` ([#8067](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8067))

IMPROVEMENTS:

* dependencies: updating `keyvault` to `2019-09-01` ([#7822](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7822))
* `azurerm_app_service_slot_virtual_network_swift_connection` - adding validation that the `app_service_id` is an App Service / Function App ID ([#8111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8111))
* `azurerm_app_service_slot_virtual_network_swift_connection` - adding validation that the `subnet` is a Subnet ID ([#8111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8111))
* `azurerm_batch_pool` - Remove `network_configuration` from update payload ([#8189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8189))
* `azurerm_frontdoor_firewall_policy` - `match_variable` within the `match_condition` block can now be set to `SocketAddr` ([#8244](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8244))
* `azurerm_linux_virtual_machine_scale_set` - `upgrade_mode="Automatic"` no longer requires health probe ([#6667](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6667))
* `azurerm_key_vault` - support for `soft_delete_retention_days` ([#7822](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7822))
* `azurerm_shared_image` - Support for `purchase_plan` ([#8124](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8124))
* `azurerm_shared_image_gallery` - validating at import time that the ID is for a Shared Image Gallery ([#8240](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8240))
* `azurerm_windows_virtual_machine_scale_set` - `upgrade_mode="Automatic"` no longer requires health probe ([#6667](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6667))

BUG FIXES:

* Data Source: `azurerm_app_service` - ensuring the `site_config` block is correctly set into the state ([#8212](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8212))
* Enhanced Validation: supporting "centralindia", "southindia" and "westindia" as valid regions in Azure Public (working around invalid data from the Azure API) ([#8217](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8217))
* `azurerm_application_gateway` - allow setting `ip_addresses` within the `backend_address_pool` block to an empty list ([#8210](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8210))
* `azurerm_application_gateway` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_container_group` - the `secure_environment_variables` field within the `container` now maps keys with empty values ([#8151](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8151))
* `azurerm_dedicated_host` - waiting for the resource to be gone 20 times rather than 10 to work around an API issue ([#8221](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8221))
* `azurerm_dedicated_host_group` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_firewall` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_hardware_security_module` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_lb` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_linux_virtual_machine` - support for updating `ultra_ssd_enabled` within the `additional_capabilities` block without recreating the virtual machine ([#8015](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8015))
* `azurerm_linux_virtual_machine_scale_set` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_managed_disk` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_nat_gateway` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_orchestrated_virtual_machine_scale_set` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_public_ip_prefix` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_public_ip` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_redis_cache` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_virtual_machine` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_virtual_machine_scale_set` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))
* `azurerm_windows_virtual_machine` - support for updating `ultra_ssd_enabled` within the `additional_capabilities` block without recreating the virtual machine ([#8015](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8015))
* `azurerm_windows_virtual_machine_scale_set` - adding validation to the `zone` field ([#8233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8233))

## 2.24.0 (August 20, 2020)

FEATURES:

* **New Resource:** `azurerm_synapse_spark_pool` ([#7886](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7886))

IMPROVEMENTS:

* dependencies: update `containerinstance` to API version `2019-12-01` ([#8110](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8110))
* `azurerm_api_management_api` - now supports `oauth2_authorization` and `openid_authentication` ([#7617](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7617))
* `azurerm_policy_definition` - `mode` can now be updated without recreating the resource ([#7976](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7976))

BUG FIXES:

* `azurerm_frontdoor` - ensuring all fields are set into the state ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor` - rewriting case-inconsistent Resource ID's to ensure they're reliable ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor_firewall_policy` - ensuring all fields are set into the state ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor_firewall_policy` - rewriting case-inconsistent Resource ID's to ensure they're reliable ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor_custom_https_configuration` - ensuring all fields are set into the state ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor_custom_https_configuration` - ensuring the `resource_group_name` field is set into the state ([#8173](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8173))
* `azurerm_frontdoor_custom_https_configuration` - rewriting case-inconsistent Resource ID's to ensure they're reliable ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_frontdoor_custom_https_configuration` - updating the ID to use the frontendEndpoint's Resource ID rather than a custom Resource ID ([#8146](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8146))
* `azurerm_lb` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_backend_address_pool` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_backend_address_pool` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_backend_address_pool` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_nat_pool` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_nat_pool` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_nat_pool` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_nat_rule` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_nat_rule` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_nat_rule` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_outbound_rule` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_outbound_rule` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_outbound_rule` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_probe` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_probe` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_probe` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_lb_rule` - adding more specific validation for the Load Balancer ID field ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_rule` - ensuring all fields are always set into the state ([#8172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8172))
* `azurerm_lb_rule` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` ([#8006](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8006))
* `azurerm_storage_account` - only sending `allow_blob_public_access` and `min_tls_version` in Azure Public since these are currently not supported in other regions ([#8148](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8148))

## 2.23.0 (August 13, 2020)

FEATURES:

* **New Resource:** `azurerm_integration_service_environment` ([#7763](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7763))
* **New Resource:** `azurerm_redis_linked_server` ([#8026](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8026))
* **New Resource:** `azurerm_synapse_firewall_rule` ([#7904](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7904))

IMPROVEMENTS:

* dependencies: updating `containerservice` to `2020-04-01` ([#7894](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7894))
* dependencies: updating `mysql` to `2020-01-01` ([#8062](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8062))
* dependencies: updating `postgresql` to `2020-01-01` ([#8045](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8045))
* Data Source: `azurerm_app_service` now exports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* Data Source: `azurerm_function_app` now exports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* Data Source: `azurerm_function_app` now exports `site_config` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_app_service` now supports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports full `ip_restriction` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports full `scm_ip_restriction` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_eventhub_namespace` - support for the `identity` block ([#8065](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8065))
* `azurerm_postgresql_server` - support for the `identity` block ([#8044](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8044))
* `azurerm_site_recovery_replicated_vm` - support setting `target_network_id` and `network_interface` on failover ([#5688](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5688))
* `azurerm_storage_account` - support `static_website` for `BlockBlobStorage` account type ([#7890](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7890))
* `azurerm_storage_account` - filter `allow_blob_public_access` and `min_tls_version` from Azure US Government ([#8092](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8092))

BUG FIXES:

* All resources using a `location` field - allowing the value `global` when using enhanced validation ([#8042](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8042))
* Data Source: `azurerm_api_management_user` - `user_id` now accepts single characters ([#7975](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7975))
* `azurerm_application_gateway` - enforce case for the `rule_type` property ([#8061](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8061))
* `azurerm_iothub_consumer_group` - lock during creation and deletion to workaround an API issue ([#8041](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8041))
* `azurerm_iothub` - the `endpoint` and `route` lists can now be cleared by setting them to `[]` ([#8028](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8028))
* `azurerm_linux_virtual_machine` - handling machines which are already stopped/deallocated ([#8000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8000))
* `azurerm_mariadb_virtual_network_rule` will now work across subscriptions ([#8100](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8100))
* `azurerm_monitor_metric_alert_resource` - continue using `SingleResourceMultiMetricCriteria` for existing alerts ([#7995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7995))
* `azurerm_mysql_server` - prevent a non empty plan when using `threat_detection_policy` ([#7981](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7981))
* `azurerm_orchestrated_virtual_machine_scale_set` - allow `single_placement_group` to be `true` ([#7821](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7821))
* `azurerm_mysql_server` - support for the `identity` block ([#8059](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8059))
* `azurerm_storage_account` - set default for `min_tls_version` to `TLS_10` ([#8152](https://github.com/terraform-providers/terraform-provider-azurerm/pull/8152))
* `azurerm_traffic_manager_profile` - updating no longer clears all endpoints ([#7846](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7846))
* `azurerm_windows_virtual_machine` - handling machines which are already stopped/deallocated [[#8000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8000)]'
* `azurerm_data_factory_dataset_delimited_text` - fix issue with property `azure_blob_storage_account` ([#7953](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7953))

## 2.22.0 (August 07, 2020)

DEPENDENCIES:

* updating `github.com/Azure/azure-sdk-for-go` to `v44.2.0` ([#7933](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7933))

IMPROVEMENTS:

* `azurerm_cosmosdb_account` - support `DisableRateLimitingResponses` with the `capabilities` property ([#8016](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8016))
* `azurerm_storage_account` - support for the `min_tls_version` property ([#7879](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7879))
* `azurerm_storage_account_sas` - support for the `signed_version attribute` property ([#8020](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8020))
* `azurerm_servicebus_queue` - support for the `enable_batched_operations`, `status`, `forward_to`, and `forward_dead_lettered_messages_to` ([#7990](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7990))

BUG FIXES:

* Data Source: `azurerm_key_vault_certificate` - fixing a crash when using acmebot certificates ([#8029](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8029))
* `azurerm_iothub_shared_access_policy` - prevent `primary_connection_string` & `secondary_connection_string` from regenerating during every apply ([#8017](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8017))

## 2.21.0 (July 31, 2020)

DEPENDENCIES:

* updating `search` to `2020-03-13` ([#7867](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7867))
* updating `go-azure-helpers` to `v0.11.2` ([#7911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7911))

FEATURES:

* **New Data Source:** `azurerm_data_share_dataset_data_lake_gen1` ([#7840](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7840))
* **New Resource:** `azurerm_dedicated_hardware_security_module` ([#7727](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7727))

IMPROVEMENTS:
* `azurerm_api_management_identity_provider_aad` - Support for `signin_tenant` ([#7901](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7901))
* `azurerm_app_service_plan` - update the relation between `kind` and `reserved` ([#7943](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7943))
* `azurerm_automation_runbook` - recreate `azurerm_automation_job_schedule` after an update ([#7555](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7555))
* `azurerm_app_service_slot` - support for the `application_logs.file_system` ([#7311](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7311))
* `azurerm_firewall` - no longer requires a `zone` ([#7817](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7817))
* `azurerm_function_app_slot` - support for the `site_config.auto_swap_slot_name` property ([#7859](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7859))
* `azurerm_kubernetes_cluster` - support for in-place upgrade from `Free` to `Paid` for `sku_tier` ([#7927](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7927))
* `azurerm_monitor_scheduled_query_rules_alert` - `action.0.custom_webhook_payload` is now sent as empty to allow for Azure's default to take effect([#7838](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7838))
* `azurerm_search_service` - support for the `public_network_access_enabled` property ([#7867](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7867))
* `azurerm_servicebus_subscription` - support for the `status` property ([#7852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7852))

BUG FIXES:

* `azurerm_automation_runbook` - allow `publish_content_link` resource to not be set ([#7824](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7824))
* `azurerm_api_management_named_value` - the `value` has been marked as sensitive to hide secret values ([#7819](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7819))
* `azurerm_cognitive_account` - allow `qname_runtime_endpoint` to not be set ([#7916](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7916))
* `azurerm_iothub_dps` - the only valid value for the `sku` property for the API is now `S1` ([#7847](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7847))
* `azurerm_eventgrid_event_subscription` - deprecate the `topic_name` as it is now readonly in the API ([#7871](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7871))
* `azurerm_kubernetes_cluster` - updates will no longer fail when using managed AAD integration ([#7874](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7874))

## 2.20.0 (July 23, 2020)

UPGRADE NOTES

* **Enhanced Validation for Locations** - the Azure Provider now validates that the value for the `location` argument is a supported Azure Region within the Azure Environment being used (from the Azure Metadata Service) - which allows us to catch configuration errors for this field at `terraform plan` time, rather than during a `terraform apply`. This functionality is now enabled by default, and can be opted-out of by setting the Environment Variable `ARM_PROVIDER_ENHANCED_VALIDATION` to `false`
* `azurerm_storage_account` - will now default `allow_blob_public_access` to false to align with the portal and be secure by default ([#7784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7784))

DEPENDENCIES:

* updating `github.com/Azure/azure-sdk-for-go` to `v44.1.0` ([#7774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7774))
* updating `cosmos` to `2020-04-01` ([#7597](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7597))

FEATURES: 

* **New Data Source:** `azurerm_synapse_workspace` ([#7517](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7517))
* **New Resource:** `azurerm_data_share_dataset_data_lake_gen1` - add `dataset_data_lake_gen1` suppport for `azurerm_data_share` ([#7511](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7511))
* **New Resource:** `azurerm_frontdoor_custom_https_configuration` - move the front door `custom_https_configuration` to its own resource to allow for parallel creation/update of custom https certificates. ([#7498](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7498))
* **New Resource:** `azurerm_kusto_cluster_customer_managed_key` ([#7520](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7520))
* **New Resource:** `azurerm_synapse_workspace` ([#7517](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7517))

IMPROVEMENTS:

* `azurerm_cosmos_db_account` - add support for the `enable_free_tier` property ([#7814](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7814))

BUG FIXES:

* Data Source: `azurerm_private_dns_zone` - fix a crash when the zone does not exist ([#7783](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7783))
* `azurerm_application_gateway` - fix crash with `gateway_ip_configuration` ([#7789](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7789))
* `azurerm_cosmos_account` - the `geo_location.prefix` property has been deprecated as service no longer accepts it as an input since Apr 25, 2019 ([#7597](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7597))
* `azurerm_monitor_autoscale_setting` - fix crash in `notification` ([#7835](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7835))
* `azurerm_storage_account` - will now default `allow_blob_public_access` to false to align with the portal and be secure by default ([#7784](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7784))

## 2.19.0 (July 16, 2020)

UPGRADE NOTES:

* HDInsight 3.6 will be retired (in Azure Public) on 2020-12-30 - HDInsight 4.0 does not support ML Services, RServer or Storm Clusters - as such the `azurerm_hdinsight_ml_services_cluster`, `azurerm_hdinsight_rserver_cluster` and `azurerm_hdinsight_storm_cluster` resources are deprecated and will be removed in the next major version of the Azure Provider. ([#7706](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7706))
* provider: no longer auto register the Microsoft.StorageCache RP ([#7768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7768))

FEATURES:

* **New Data source:** `azurerm_route_filter` ([#6341](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6341))
* **New Resource:** `azurerm_route_filter` ([#6341](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6341))

IMPROVEMENTS:

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
* `azurerm_storage_container` - container creation will retry if a container of the same name has not completed its delete operation ([#7179](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7179))
* `azurerm_storage_share` - share creation will retry if a share of the same name has not completed its previous delete operation ([#7179](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7179))
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

IMPROVEMENTS:

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

IMPROVEMENTS:

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

IMPROVEMENTS:

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

IMPROVEMENTS:

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
