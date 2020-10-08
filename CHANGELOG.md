## 2.31.0 (Unreleased)

UPGRADE NOTES

* This release updates the `azurerm_security_center_subscription_pricing` resource to use the latest version of the Security API which now allows configuring multiple Resource Types - as such a new field `resource_type` is now available. Configurations default the `resource_type` to `VirtualMachines` which matches the behaviour of the previous release - but your Terraform Configuration may need updating.

FEATURES:

* **New Resource:** `azurerm_service_fabric_mesh_application` [GH-6761]
* **New Resource:** `azurerm_virtual_desktop_workspace` [GH-8605]
* **New Resource:** `azurerm_virtual_desktop_host_pool` [GH-8605]
* **New Resource:** `azurerm_virtual_desktop_application_group` [GH-8605]
* **New Resource:** `azurerm_virtual_desktop_workspace_application_group_association` [GH-8605]

IMPROVEMENTS:

* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v46.4.0` [GH-8642]
* `data.azurerm_application_insights` - Add support for the `connection_string` property [GH-8699]
* `azurerm_app_service` - allow v6 IPs for the `ip_restriction` property [GH-8599]
* `azurerm_application_insights` - Add support for the `connection_string` property [GH-8699]
* `azurerm_backup_policy_vm` - validate daily backups is > `7` [GH-7898]
* `azurerm_dedicated_host` - add support for the `DSv4-Type1` and `sku_name` properties [GH-8718]
* `azurerm_iothub` - Support for the `public_network_access_enabled` property [GH-8586]
* `azurerm_key_vault_certificate_issuer` - the `org_id` property is now optional [GH-8687]
* `azurerm_kubernetes_cluster_node_pool` - `node_count`, `min_node`, and `max_node` can now be set to `0` [GH-8300]
* `azurerm_mssql_database` - support `0` for the `min_capacity` property [GH-8308]
* `azurerm_mssql_database` - add support for `long_term_retention_policy` and `short_term_retention_policy` [GH-8765] 
* `azurerm_mssql_server` - support the `minimum_tls_version` property [GH-8361]
* `azurerm_mssql_virtual_machine` - Add support for `storage_configuration_settings` [GH-8623]
* `azurerm_security_center_subscription_pricing` - now supports per `resource_type` pricing [GH-8549]
* `azurerm_storage_account` - add support for `large_file_share_enabled` [GH-8789]
* `azurerm_storage_share` allow large quotas to be specified.  [GH-8666]

BUG FIXES:

* `azurerm_function_app` - mark the `app_settings` block as computed [GH-8682]
* `azurerm_function_app_slot` - mark the `app_settings` block as computed [GH-8682]
* `azurerm_policy_set_definition` - corrects issue with empty `parameter_values` attribute [GH-8668]
* `azurerm_policy_definition` - `mode` property now enforces correct case [GH-8795]

---

For information on changes between the v2.30.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.44.0 changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
