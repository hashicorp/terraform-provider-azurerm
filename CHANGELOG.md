## 2.70.0 (Unreleased)

FEATURES:

* **New Resource** `azurerm_bot_channel_alexa` [GH-12682]
* **New Resource** `azurerm_bot_channel_facebook` [GH-12709]
* **New Resource** `azurerm_bot_channel_sms` [GH-12713]
* **New Resource** `azurerm_data_factory_trigger_custom_event` [GH-12448]
* **New Resource** `azurerm_data_factory_trigger_tumbling_window` [GH-12437]
* **New Resource** `azurerm_data_protection_backup_instance_disk` [GH-12617]

ENHANCEMENTS:

* dependencies: Upgrade `web` (App Service) API to `2021-01-15` [GH-12635]
* analysisservices: refactoring to use an Embedded SDK [GH-12771]
* maps: refactoring to use an Embedded SDK [GH-12716]
* msi: refactoring to use an Embedded SDK [GH-12715]
* relay: refactoring to use an Embedded SDK [GH-12772]
* vmware: refactoring to use an Embedded SDK [GH-12751]
* Data Source: `azurerm_storage_account_sas` - support for the property `ip_addresses` [GH-12705]
* `azurerm_api_management_diagnostic` - support for the property `operation_name_format` [GH-12736]
* `azurerm_automation_certificate` - the `exportable` property can now be set [GH-12738]
* `azurerm_data_factory_dataset_binary` - the blob `path` and `filename` propeties are now optional [GH-12676]
* `azurerm_data_factory_trigger_blob_event` - support for the `activation` property [GH-12644]
* `azurerm_data_factory_pipeline` - support for the `concurrency` and `moniter_metrics_after_duration` properties [GH-12685]
* `azurerm_hdinsight_interactive_query_cluster` - support for the `encryption_in_transit_enabled` property [GH-12767]
* `azurerm_hdinsight_spark_cluster` - support for the `encryption_in_transit_enabled` property [GH-12767]

BUG FIXES:

* `azurerm_cdn_endpoint` - fixing a crash when the future is nil [GH-12743]
* `azurerm_private_endpoint` - working around a casing issue in `private_connection_resource_id` for MariaDB, MySQL and PostgreSQL resources [GH-12761]

---

For information on changes between the v2.69.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
