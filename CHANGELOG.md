## 2.91.0 (Unreleased)

FEATURES:

ENHANCEMENTS:

* `azurerm_application_gateway` - support for the `key_vault_secret_id` and `force_firewall_policy_association` property [GH-14413]
* `azurerm_managed_disk` - support for the `gallery_image_reference_id` property [GH-14121]

BUG FIXES:

## 2.90.0 (December 17, 2021)

FEATURES:

* **New Datasource:** `azurerm_app_configuration_key` ([#14484](https://github.com/hashicorp/terraform-provider-azurerm/issues/14484))
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
