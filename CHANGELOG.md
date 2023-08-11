## 3.70.0 (Unreleased)

ENHANCEMENTS:

* dependencies: updating to `v0.20230810.1125717` of `github.com/hashicorp/go-azure-sdk` [GH-22874]
* `cosmos`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` [GH-22874]
* `policy`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` [GH-22874]
* `postgres`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` [GH-22874]
* `recoveryservices`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` [GH-22874]
* `resources`: updating to use the base layer from `hashicorp/go-azure-sdk` rather than `Azure/go-autorest` [GH-22874]

BUG FIXES:

* `azurerm_iothub_dps` - updating the validation for `target` within the `ip_filter_rule` block to match the values defined in the Azure API Definitions [GH-22891]
* `azurerm_postgresql_database` - updating the validation for `collation` to include `en-GB` [GH-22907]

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

---

For information on changes between the v3.59.0 and v3.0.0 releases, please see [the previous v3.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v3.md).

For information on changes between the v2.99.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v1.44.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
