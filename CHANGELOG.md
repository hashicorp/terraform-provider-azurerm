## 2.31.0 (Unreleased)

FEATURES:

* **New Resource:** `azurerm_service_fabric_mesh_application` [GH-6761]

IMPROVEMENTS:

* `data.azurerm_application_insights` - Add support for `connection_string` [GH-8699]
* `azurerm_app_service` - allow v6 IPs for the `ip_restriction` property [GH-8599]
* `azurerm_application_insights` - Add support for `connection_string` [GH-8699]
* `azurerm_dedicated_host` - add support for `DSv4-Type1` `sku_name` [GH-8718]
* `azurerm_iothub` - Support for `public_network_access_enabled` [GH-8586]
* `azurerm_key_vault_certificate_issuer` - `org_id` is now optional [GH-8687]
* `azurerm_mssql_server` - support the `minimum_tls_version` property [GH-8361]
* `azurerm_mssql_virtual_machine` - Add support for `storage_configuration_settings` [GH-8623]

ENHANCEMENTS:

* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v46.4.0` [GH-8642]

---

For information on changes between the v2.30.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.44.0 changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
