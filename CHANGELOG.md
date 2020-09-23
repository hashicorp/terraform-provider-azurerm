## 2.29.0 (Unreleased)

UPGRADE NOTES:

* `azurerm_api_management` - the value `None` has been removed from the `identity` block to match other resources, to specify an API Management Service with no Managed Identity remove the `identity` block [GH-8411]
* `azurerm_container_registry` -  the `storage_account_id` property now forces a new resource as required by the updated API version [GH-8477]

FEATURES: 

* **New Data Source:** `azurerm_data_share_dataset_kusto_cluster` [GH-8464]
* **New Data Source:** `azurerm_databricks_workspace` [GH-8502]
* **New Data Source:** `azurerm_firewall_policy` [GH-7390]
* **New Data Source:** `azurerm_storage_sync_group` [GH-8462]
* **New Data Source:** `azurerm_mssql_server` [GH-7917]
* **New Resource:** `azurerm_data_share_dataset_kusto_cluster` [GH-8464]
* **New resource:** `azurerm_firewall_policy` [GH-7390]
* **New resource:** `azurerm_mysql_server_key` [GH-8125]
* **New Resource:** `azurerm_postgresql_server_key` [GH-8126]

IMPROVEMENTS:

* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v46.3.0` [GH-8592]
* dependencies: updating `containerregistry` to `2019-05-01` [GH-8477]
* Data Source: `azurerm_api_management` - export the `private_ip_addresses` attribute for primary and additional locations [GH-8290]
* `azurerm_api_management` - support the `virtual_network_configuration` block for additional locations [GH-8290]
* `azurerm_api_management` - export the `private_ip_addresses` attribute for additional locations [GH-8290]
* `azurerm_cosmosdb_account` - support the `Serverless` value for the `capabilities` property [GH-8533]
* `azurerm_cosmosdb_sql_container` - support for the `indexing_policy` property [GH-8461]
* `azurerm_mssql_server` - support for the `recover_database_id` and `restore_dropped_database_id` properties [GH-7917]
* `azurerm_policy_set_definition` - support for typed parameter values other then string in `the policy_definition_reference` block deprecating `parameters` in favour of `parameter_vcaluess` [GH-8270]
* `azurerm_search_service` - Add support for `allowed_ips` [GH-8557]
* `azurerm_service_fabric_cluster` - Remove two block limit for `client_certificate_thumbprint` [GH-8521]
* `azurerm_signalr_service` - support for delta updates [GH-8541]
* `azurerm_windows_virtual_machine` - support for updating the `license_type` field [GH-8542]

BUG FIXES:

* `azurerm_api_management` - the value `None` for the field `type` within the `identity` block has been removed - to remove a managed identity remove the `identity` block [GH-8411]
* `azurerm_app_service` - don't try to manage source_control when scm_type is `VSTSRM` [GH-8531]
* `azurerm_function_app` - don't try to manage source_control when scm_type is `VSTSRM` [GH-8531]
* `azurerm_kubernetes_cluster` - picking the first system node pool if the original `default_node_pool` has been removed [GH-8503]

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

ENHANCEMENTS:

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

ENHANCEMENTS:

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

ENHANCEMENTS:

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

* Data Source: `azurerm_app_service` now exports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* Data Source: `azurerm_function_app` now exports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* Data Source: `azurerm_function_app` now exports `site_config` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_app_service` now supports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports full `ip_restriction` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports full `scm_ip_restriction` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_site_recovery_replicated_vm` - support setting `target_network_id` and `network_interface` on failover ([#5688](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5688))
* `azurerm_storage_account` - support `static_website` for `BlockBlobStorage` account type ([#7890](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7890))
* `azurerm_storage_account` - filter `allow_blob_public_access` and `min_tls_version` from Azure US Government ([#8092](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8092))

ENHANCEMENTS:

* dependencies: updating `containerservice` to `2020-04-01` ([#7894](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7894))
* dependencies: updating `mysql` to `2020-01-01` ([#8062](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8062))
* dependencies: updating `postgresql` to `2020-01-01` ([#8045](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8045))
* `azurerm_eventhub_namespace` - support for the `identity` block ([#8065](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8065))
* `azurerm_postgresql_server` - support for the `identity` block ([#8044](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8044))

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

---

For information on changes between the v2.20.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.44.0 changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
