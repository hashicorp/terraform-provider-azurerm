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
