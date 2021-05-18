## 2.60.0 (Unreleased)

FEATURES:

* **New Data Source:** `azurerm_redis_enterprise_database` [GH-11734]

ENHANCEMENTS:

* dependencies: updating `aks` to use API Version `2021-03-01` [GH-11708]
* dependencies: updating `eventgrid` to use API Version `2020-10-15-preview` [GH-11746]
* `azurerm_cosmosdb_mongo_collection` - support for the `analytical_storage_ttl` property [GH-11735]
* `azurerm_cosmosdb_cassandra_table` - support for the `analytical_storage_ttl` property [GH-11755]
* `azurerm_hdinsight_kafka_cluster` - support for the `encryption_in_transit_enabled` property [GH-11737]
* `azurerm_media_services_account` - support for the `key_delivery_access_control` block [GH-11726]
* `azurerm_redis_cache` - suppot for the `replicas_per_master` peoperty [GH-11714]

BUG FIXES:

* `azurerm_linux_virtual_machine_scale_set` - the `extension` blocks are now a set [GH-11425]
* `azurerm_virtual_network_gateway_connection` - fix a bug where `shared_key` was not being updated [GH-11742]
* `azurerm_windows_virtual_machine_scale_set` - the `extension` blocks are now a set [GH-11425]
* `azurerm_windows_virtual_machine_scale_set` - changing the `license_type` will no longer create a new resource [GH-11731]

---

For information on changes between the v2.59.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
