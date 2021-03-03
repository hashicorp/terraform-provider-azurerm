## 2.50.0 (Unreleased)

FEATURES:

* **New Data Source:** `azurerm_vmware_private_cloud` [GH-9284]
* **New Resource:** `azurerm_kusto_eventgrid_data_connection` [GH-10712]
* **New Resource:** `azurerm_sentinel_data_connector_aws_cloud_trail` [GH-10664]
* **New Resource:** `azurerm_sentinel_data_connector_azure_active_directory` [GH-10665]
* **New Resource:** `azurerm_sentinel_data_connector_office_365` [GH-10671]
* **New Resource:** `azurerm_sentinel_data_connector_threat_intelligence` [GH-10670]
* **New Resource:** `azurerm_vmware_private_cloud` [GH-9284]

ENHANCEMENTS:

* dependencies: updating to `v52.0.0` of `github.com/Azure/azure-sdk-for-go` [GH-10787]
* dependencies: updating `compute` to API version `2020-12-01` [GH-10650]
* `azurerm_keyvault_secret` - support for the `versionless_id` property [GH-10738]
* `azurerm_kubernetes_cluster` - support `private_dns_zone_id` when using a `service_principal` [GH-10737]
* `azurerm_kusto_cluster` - supports for the `double_encryption_enabled` property [GH-10264]
* `azurerm_mssql_database` - support for the `log_monitoring_enabled` property within the `extended_auditing_policy` block [GH-10324]
* `azurerm_mssql_database_extended_auditing_policy ` - support for the `log_monitoring_enabled` property [GH-10324]
* `azurerm_mssql_server` - support for the `log_monitoring_enabled` property within the `extended_auditing_policy` block [GH-10324]
* `azurerm_mssql_server_extended_auditing_policy ` - support for the `log_monitoring_enabled` property [GH-10324] 
* `azurerm_signalr_service` - support for the `upstream_endpoint` block [GH-10459]
* `azurerm_sql_server` - support for the `log_monitoring_enabled` property within the `extended_auditing_policy` block [GH-10324]
* `azurerm_sql_database` - support for the `log_monitoring_enabled` property within the `extended_auditing_policy` block [GH-10324]
* `azurerm_spring_cloud_java_deployment` - supporting delta updates [GH-10729]
* `azurerm_virtual_network_gateway` - deprecate `peering_address` in favour of `peering_addresses` [GH-10381]

BUG FIXES:

* `azurerm_api_management` - changing the `sku_name` property no longer forces a new resouce to be created [GH-10747]
* `azurerm_api_management` - the field `tenant_access` can only be configured when not using a Consumption SKU [GH-10766]
* `azurerm_cosmosdb_mongo_collection` - ignore throughput if Cosmos DB provisioned in 'serverless' capacity mode [GH-10389]
* `azurerm_linux_virtual_machine` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue [GH-10722]
* `azurerm_linux_virtual_machine_scale_set` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue [GH-10722]
* `azurerm_virtual_machine` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue [GH-10722]
* `azurerm_virtual_machine_scale_set` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue [GH-10722]
* `azurerm_windows_virtual_machine` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue [GH-10722]
* `azurerm_windows_virtual_machine_scale_set` - parsing the User Assigned Identity ID case-insensitively to work around an Azure API issue [GH-10722]

---

For information on changes between the v2.49.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
