## 2.71.0 (Unreleased)

FEATURES:

* **New Data Source:** `azurerm_databricks_workspace_private_endpoint_connection` [GH-12543]
* **New Resource:** `azurerm_bot_channel_line` [GH-12746]
* **New Resource:** `azurerm_cdn_endpoint_custom_domain` [GH-12496]
* **New Resource:** `azurerm_data_factory_data_flow` [GH-12588]

ENHANCEMENTS:

* dependencies: upgrading to `v56.0.0` of `github.com/Azure/azure-sdk-for-go` [GH-12781]
* dependencies: updating `containerservice` to use API Version `2021-05-1` [GH-12747]
* dependencies: updating `machinelearning` to use API Version `2021-04-01` [GH-12804]
* dependencies: updating `databricks` to use API Version `2021-04-01-preview` [GH-12543]
* PowerBI: refactoring to use an Embedded SDK [GH-12787]
* SignalR: refactoring to use an Embedded SDK [GH-12785]
* `azurerm_api_management_api_diagnostic` - support for the `operation_name_format` property [GH-12782]
* `azurerm_databricks_workspace` - support for private link endpoint [GH-12543]
* `azurerm_databricks_workspace` - add support for `Customer Managed Keys for Managed Services` [GH-12799]
* `azurerm_kubernetes_cluster` - support for the `local_account_disabled` property [GH-12386]
* `azurerm_logic_app_workflow` - support for the `workflow_parameters` [GH-12314]
* `azurerm_mssql_database` - support for the `Free` and `FSV2` SKU's [GH-12835]
* `azurerm_public_ip_resource` - support for sku_tier property [GH-12775]

BUG FIXES:

* `azurerm_api_management_api_subscription` - fixing the default scope to be `/apis` rather than `all_apis` as required by the latest API [GH-12829]
* `azurerm_app_service_active_slot` - fix 404 not found on read for slot [GH-12792]
* `azurerm_linux_virtual_machine_scale_set` - fix crash in checking for latest image [GH-12808]
* `azurerm_local_network_gateway` - fixing a crash where the `LocalNetworkAddressSpace` block was nil [GH-12822]
* `azurerm_postgresql_database` - fixing a crash in the Azure SDK [GH-12823]
* `azurerm_private_dns_zone` - fixing a crash during deletion [GH-12824]
* `azurerm_virtual_machine_extension` - changing the `publisher` property now creates a new resource [GH-12790]


## 2.70.0 (July 30, 2021)

FEATURES:

* **New Data Source** `azurerm_storage_share` ([#12693](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12693))
* **New Resource** `azurerm_bot_channel_alexa` ([#12682](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12682))
* **New Resource** `azurerm_bot_channel_direct_line_speech` ([#12735](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12735))
* **New Resource** `azurerm_bot_channel_facebook` ([#12709](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12709))
* **New Resource** `azurerm_bot_channel_sms` ([#12713](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12713))
* **New Resource** `azurerm_data_factory_trigger_custom_event` ([#12448](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12448))
* **New Resource** `azurerm_data_factory_trigger_tumbling_window` ([#12437](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12437))
* **New Resource** `azurerm_data_protection_backup_instance_disk` ([#12617](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12617))

ENHANCEMENTS:

* dependencies: Upgrade `web` (App Service) API to `2021-01-15` ([#12635](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12635))
* analysisservices: refactoring to use an Embedded SDK ([#12771](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12771))
* maps: refactoring to use an Embedded SDK ([#12716](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12716))
* msi: refactoring to use an Embedded SDK ([#12715](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12715))
* relay: refactoring to use an Embedded SDK ([#12772](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12772))
* vmware: refactoring to use an Embedded SDK ([#12751](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12751))
* Data Source: `azurerm_storage_account_sas` - support for the property `ip_addresses` ([#12705](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12705))
* `azurerm_api_management_diagnostic` - support for the property `operation_name_format` ([#12736](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12736))
* `azurerm_automation_certificate` - the `exportable` property can now be set ([#12738](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12738))
* `azurerm_data_factory_dataset_binary` - the blob `path` and `filename` propeties are now optional ([#12676](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12676))
* `azurerm_data_factory_trigger_blob_event` - support for the `activation` property ([#12644](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12644))
* `azurerm_data_factory_pipeline` - support for the `concurrency` and `moniter_metrics_after_duration` properties ([#12685](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12685))
* `azurerm_hdinsight_interactive_query_cluster` - support for the `encryption_in_transit_enabled` property ([#12767](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12767))
* `azurerm_hdinsight_spark_cluster` - support for the `encryption_in_transit_enabled` property ([#12767](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12767))
* `azurerm_firewall_polcy` - support for property `private_ip_ranges` ([#12696](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12696))

BUG FIXES:

* `azurerm_cdn_endpoint` - fixing a crash when the future is nil ([#12743](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12743))
* `azurerm_private_endpoint` - working around a casing issue in `private_connection_resource_id` for MariaDB, MySQL and PostgreSQL resources ([#12761](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12761))

---

For information on changes between the v2.69.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
