## 2.31.0 (Unreleased)

IMPROVEMENTS:

* `azurerm_app_service` - allow v6 IPs for the `ip_restriction` property [GH-8599]
* `azurerm_dedicated_host` - add support for `DSv4-Type1` `sku_name` [GH-8718]
* `azurerm_iothub` - Support for `public_network_access_enabled` [GH-8586]
* `azurerm_key_vault_certificate_issuer` - `org_id` is now optional [GH-8687]

ENHANCEMENTS:

* dependencies: updating `github.com/Azure/azure-sdk-for-go` to `v46.4.0` [GH-8642]

---

For information on changes between the v2.30.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.44.0 changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
