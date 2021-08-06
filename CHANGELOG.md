## 2.71.0 (August 06, 2021)

FEATURES:

* **New Data Source:** `azurerm_databricks_workspace_private_endpoint_connection` ([#12543](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12543))
* **New Resource:** `azurerm_api_management_tag` ([#12535](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12535))
* **New Resource:** `azurerm_bot_channel_line` ([#12746](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12746))
* **New Resource:** `azurerm_cdn_endpoint_custom_domain` ([#12496](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12496))
* **New Resource:** `azurerm_data_factory_data_flow` ([#12588](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12588))
* **New Resource:** `azurerm_postgresql_flexible_server_database` ([#12550](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12550))

ENHANCEMENTS:

* dependencies: upgrading to `v56.0.0` of `github.com/Azure/azure-sdk-for-go` ([#12781](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12781))
* dependencies: updating `appinsights` to use API Version `2020-02-02` ([#12818](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12818))
* dependencies: updating `containerservice` to use API Version `2021-05-1` ([#12747](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12747))
* dependencies: updating `machinelearning` to use API Version `2021-04-01` ([#12804](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12804))
* dependencies: updating `databricks` to use API Version `2021-04-01-preview` ([#12543](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12543))
* PowerBI: refactoring to use an Embedded SDK ([#12787](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12787))
* SignalR: refactoring to use an Embedded SDK ([#12785](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12785))
* `azurerm_api_management_api_diagnostic` - support for the `operation_name_format` property ([#12782](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12782))
* `azurerm_app_service` - support for the acr_use_managed_identity_credentials and acr_user_managed_identity_client_id properties ([#12745](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12745))
* `azurerm_app_service` - support `v6.0` for the `dotnet_framework_version` property ([#12788](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12788))
* `azurerm_application_insights` - support for the `workspace_id` property ([#12818](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12818))
* `azurerm_databricks_workspace` - support for private link endpoint ([#12543](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12543))
* `azurerm_databricks_workspace` - add support for `Customer Managed Keys for Managed Services` ([#12799](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12799))
* `azurerm_data_factory_linked_service_data_lake_storage_gen2` - don't send a secure connection string when using a managed identity ([#12359](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12359))
* `azurerm_function_app` - support for the `elastic_instance_minimum`, `app_scale_limit`, and `runtime_scale_monitoring_enabled` properties ([#12741](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12741))
* `azurerm_kubernetes_cluster` - support for the `local_account_disabled` property ([#12386](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12386))
* `azurerm_kubernetes_cluster` - support for the `maintenance_window` block ([#12762](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12762))
* `azurerm_kubernetes_cluster` - the field `automatic_channel_upgrade` can now be set to `node-image` ([#12667](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12667))
* `azurerm_logic_app_workflow` - support for the `workflow_parameters` ([#12314](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12314))
* `azurerm_mssql_database` - support for the `Free` and `FSV2` SKU's ([#12835](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12835))
* `azurerm_network_security_group` - the `protocol` property now supports `Ah` and `Esp` values ([#12865](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12865))
* `azurerm_public_ip_resource` - support for sku_tier property ([#12775](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12775))
* `azurerm_redis_cache` - support for the `replicas_per_primary`, `redis_version`, and `tenant_settings` properties and blocks ([#12820](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12820))
* `azurerm_redis_enterprise_cluster` - this can now be provisioned in `Canada Central` ([#12842](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12842))
* `azurerm_static_site` - support `Standard` SKU ([#12510](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12510))

BUG FIXES:

* Data Source `azurerm_ssh_public_key` - normalising the SSH Public Key ([#12800](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12800))
* `azurerm_api_management_api_subscription` - fixing the default scope to be `/apis` rather than `all_apis` as required by the latest API ([#12829](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12829))
* `azurerm_app_service_active_slot` - fix 404 not found on read for slot ([#12792](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12792))
* `azurerm_linux_virtual_machine_scale_set` - fix crash in checking for latest image ([#12808](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12808))
* `azurerm_kubernetes_cluster` - corrently valudate the `net_ipv4_ip_local_port_range_max` property ([#12859](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12859))
* `azurerm_local_network_gateway` - fixing a crash where the `LocalNetworkAddressSpace` block was nil ([#12822](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12822))
* `azurerm_notification_hub_authorization_rule` - switching to use an ID Formatter ([#12845](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12845))
* `azurerm_notification_hub` - switching to use an ID Formatter ([#12845](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12845))
* `azurerm_notification_hub_namespace` - switching to use an ID Formatter ([#12845](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12845))
* `azurerm_postgresql_database` - fixing a crash in the Azure SDK ([#12823](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12823))
* `azurerm_private_dns_zone` - fixing a crash during deletion ([#12824](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12824))
* `azurerm_resource_group_template_deployment` - fixing deletion of nested items when using non-top level items ([#12421](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12421))
* `azurerm_subscription_template_deployment` - fixing deletion of nested items when using non-top level items ([#12421](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12421))
* `azurerm_virtual_machine_extension` - changing the `publisher` property now creates a new resource ([#12790](https://github.com/terraform-providers/terraform-provider-azurerm/issues/12790))

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
