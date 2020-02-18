## 2.0.0 (Unreleased)

NOTES:

* **Major Version:** Version 2.0 of the Azure Provider is a major version - some deprecated fields/resources have been removed - please [refer to the 2.0 upgrade guide for more information](https://www.terraform.io/docs/providers/azurerm/guides/2.0-upgrade-guide.html).
* **Provider Block:** The Azure Provider now requires that a `features` block is specified within the Provider block, which can be used to alter the behaviour of certain resources - [more information on the `features` block can be found in the documentation](https://www.terraform.io/docs/providers/azurerm/index.html#features).
* **Terraform 0.10/0.11:** Version 2.0 of the Azure Provider no longer supports Terraform 0.10 or 0.11 - you must upgrade to Terraform 0.12 to use version 2.0 of the Azure Provider.

FEATURES:

* **Custom Timeouts:** - all resources within the Azure Provider now allow configuring custom timeouts - please [see Terraform's Timeout documentation](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) and the documentation in each data source resource for more information.
* **Requires Import:** The Azure Provider now checks for the presence of an existing resource prior to creating it - which means that if you try and create a resource which already exists (without importing it) you'll be prompted to import this into the state.
* **New Resource:** `azurerm_linux_virtual_machine` [GH-5705]
* **New Resource:** `azurerm_linux_virtual_machine_scale_set` [GH-5705]
* **New Resource:** `azurerm_virtual_machine_scale_set_extension` [GH-5705]
* **New Resource:** `azurerm_windows_virtual_machine` [GH-5705]
* **New Resource:** `azurerm_windows_virtual_machine_scale_set` [GH-5705]

BREAKING CHANGES:

* Data Source: `azurerm_app_service_plan` - the deprecated `properties` block has been removed since these properties have been moved to the top level [GH-5717]
* Data Source: `azurerm_azuread_application` - This data source has been removed since it was deprecated [GH-5748]
* Data Source: `azurerm_azuread_service_principal` - This data source has been removed since it was deprecated [GH-5748]
* Data Source: `azurerm_key_vault` - removing the `sku` block since this has been deprecated in favour of the `sku_name` field [GH-5774]
* Data Source: `azurerm_key_vault_key` - removing the deprecated `vault_uri` field [GH-5774]
* Data Source: `azurerm_key_vault_secret` - removing the deprecated `vault_uri` field [GH-5774]
* Data Source: `azurerm_role_definition` - removing the alias `VirtualMachineContributor` which has been deprecated in favour of the full name `Virtual Machine Contributor` [GH-5733]
* Data Source: `azurerm_scheduler_job_collection` - This data source has been removed since it was deprecated [GH-5712]
* `azurerm_app_service_plan` - the deprecated `properties` block has been removed since these properties have been moved to the top level [GH-5717]
* `azurerm_application_gateway` - updating the default value for the `body` field within the `match` block from `*` to an empty string [GH-5752]
* `azurerm_availability_set` - updating the default value for `managed` from `false` to `true` [GH-5724]
* `azurerm_azuread_application` - This resource has been removed since it was deprecated [GH-5748]
* `azurerm_azuread_service_principal_password` - This resource has been removed since it was deprecated [GH-5748]
* `azurerm_azuread_service_principal` - This resource has been removed since it was deprecated [GH-5748]
* `azurerm_cognitive_account` - removing the deprecated `sku_name` block [GH-5797]
* `azurerm_container_service` - This resource has been removed since it was deprecated [GH-5709]
* `azurerm_iot_dps` - This resource has been removed since it was deprecated [GH-5753]
* `azurerm_iot_dps_certificate` - This resource has been removed since it was deprecated [GH-5753]
* `azurerm_key_vault` - removing the `sku` block since this has been deprecated in favour of the `sku_name` field [GH-5774]
* `azurerm_key_vault_access_policy` - removing the deprecated field `vault_name` which has been superseded by the `key_vault_id` field [GH-5774]
* `azurerm_key_vault_access_policy` - removing the deprecated field `resource_group_name ` which has been superseded by the `key_vault_id` field [GH-5774]
* `azurerm_key_vault_certificate` - removing the deprecated `vault_uri` field [GH-5774]
* `azurerm_key_vault_key` - removing the deprecated `vault_uri` field [GH-5774]
* `azurerm_key_vault_secret` - removing the deprecated `vault_uri` field [GH-5774]
* `azurerm_kubernetes_cluster` - updating the default value for `load_balancer_sku` to `Standard` from `Basic` [GH-5747]
* `azurerm_log_analytics_workspace_linked_service` - This resource has been removed since it was deprecated [GH-5754]
* `azurerm_maps_account` - the `sku_name` field is now case-sensitive [GH-5776]
* `azurerm_mariadb_server` - removing the `sku` block since it's been deprecated in favour of the `sku_name` field [GH-5777]
* `azurerm_mssql_elasticpool` - removing the deprecated `elastic_pool_properties` block [GH-5744]
* `azurerm_notification_hub_namesapce` - removing the `sku` block in favour of the `sku_name` argument [GH-5722]
* `azurerm_postgresql_server` - removing the `sku` block which has been deprecated in favour of the `sku_name` field [GH-5721]
* `azurerm_relay_namespace` - removing the `sku` block in favour of the `sku_name` field [GH-5719]
* `azurerm_scheduler_job` - This resource has been removed since it was deprecated [GH-5712]
* `azurerm_scheduler_job_collection` - This resource has been removed since it was deprecated [GH-5712]
* `azurerm_storage_account` - removing the deprecated `account_type` field [GH-5710]
* `azurerm_storage_account` - removing the deprecated `enable_advanced_threat_protection` field [GH-5710]
* `azurerm_storage_blob` - making the `type` field case-sensitive [GH-5710]
* `azurerm_storage_blob` - removing the deprecated `attempts` field [GH-5710]
* `azurerm_storage_blob` - removing the deprecated `resource_group_name` field [GH-5710]
* `azurerm_storage_container` - removing the deprecated `resource_group_name` field [GH-5710]
* `azurerm_storage_container` - removing the deprecated `properties` block [GH-5710]
* `azurerm_storage_queue` - removing the deprecated `resource_group_name` field [GH-5710]
* `azurerm_storage_share` - removing the deprecated `resource_group_name` field [GH-5710]
* `azurerm_storage_table` - removing the deprecated `resource_group_name` field [GH-5710]

IMPROVEMENTS:

* Data Source: `azurerm_kubernetes_service_version` - support for filtering of preview releases [GH-5662]
* `azurerm_dedicated_host` - support for setting `sku_name` to `DSv3-Type2` and `ESv3-Type2` [GH-5768]
* `azurerm_storage_account` - support for configuring the `static_website` block [GH-5649]
* `azurerm_storage_account` - support for configuring `cors_rules` within the `blob_properties` block [GH-5425]
* `azurerm_windows_virtual_machine` - fixing a bug when provisioning from a Shared Gallery image [GH-5661]

BUG FIXES:

* `azurerm_linux_virtual_machine` - using the delete custom timeout during deletion [GH-5764]
* `azurerm_public_ip_prefix` - fixing the validation for the `prefix_length` to match the Azure API [GH-5693]
* `azurerm_role_assignment` - validating that the `name` is a UUID [GH-5624]
* `azurerm_signalr_service` - ensuring the SignalR segment is parsed in the correct case [GH-5737]
* `azurerm_windows_virtual_machine` - using the delete custom timeout during deletion [GH-5764]

---

For information on v1.44.0 and prior releases, please see [the v1.44.0 changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/v1.44.0/CHANGELOG.md).
