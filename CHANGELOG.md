## 2.61.0 (Unreleased)

FEATURES:

* **New Data Source:** `azurerm_spatial_anchors_account` [GH-11824]

ENHANCEMENTS:

* dependencies: updating to `v54.3.0` of `github.com/Azure/azure-sdk-for-go` [GH-11813]
* dependencies: updating `mixedreality` to use API Version `2021-01-01` [GH-11824]
* refactor: switching to use an embedded SDK for `appconfiguration` [GH-11959]
* refactor: switching to use an embedded SDK for `eventhub` [GH-11973]
* provider: support for the Virtual Machine `skip_shutdown_and_force_delete` feature [GH-11216]
* provider: support for the Virtual Machine Scale Set `force_delete` feature [GH-11216]
* provider: no longer auto register the Microsoft.DevSpaces RP [GH-11822]
* Data Source: `azurerm_key_vault_certificate_data` - support certificate bundles and add support for ECDSA keys [GH-11974]
* `azurerm_data_factory_linked_service_sftp` - support for hostkey related properties [GH-11825]
* `azurerm_spatial_anchors_account` - support for `account_domain` and `account_id` [GH-11824]
* `azurerm_static_site`: Add support for `tags` attribute [GH-11849]
* `azurerm_storage_account`: `private_link_access` supports more values [GH-11957]
* `azurerm_storage_account_network_rules`: `private_link_access` supports more values [GH-11957]
* `azurerm_synapse_spark_pool` - `spark_version` now supports `3.0` [GH-11972]

BUG FIXES:

* `azurerm_cdn_endpoint` - do not send an empty `origin_host_header` to the api [GH-11852]
* `azurerm_linux_virtual_machine_scale_set`: changing the `disable_automatic_rollback` and `enable_automatic_os_upgrade` properties no longer created a new resource [GH-11723]
* `azurerm_storage_share`: Fix ID for `resource_manager_id` [GH-11828]
* `azurerm_windows_virtual_machine_scale_set`: changing the `disable_automatic_rollback` and `enable_automatic_os_upgrade` properties no longer created a new resource [GH-11723]

## 2.60.0 (May 20, 2021)

FEATURES:

* **New Data Source:** `azurerm_eventhub_cluster` ([#11763](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11763))
* **New Data Source:** `azurerm_redis_enterprise_database` ([#11734](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11734))
* **New Resource:** `azurerm_static_site` ([#7150](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7150))
* **New Resource:** `azurerm_machine_learning_inference_cluster` ([#11550](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11550))

ENHANCEMENTS:

* dependencies: updating `aks` to use API Version `2021-03-01` ([#11708](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11708))
* dependencies: updating `eventgrid` to use API Version `2020-10-15-preview` ([#11746](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11746))
* `azurerm_cosmosdb_mongo_collection` - support for the `analytical_storage_ttl` property ([#11735](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11735))
* `azurerm_cosmosdb_cassandra_table` - support for the `analytical_storage_ttl` property ([#11755](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11755))
* `azurerm_healthcare_service` - support for the `public_network_access_enabled` property ([#11736](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11736))
* `azurerm_hdinsight_kafka_cluster` - support for the `encryption_in_transit_enabled` property ([#11737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11737))
* `azurerm_media_services_account` - support for the `key_delivery_access_control` block ([#11726](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11726))
* `azurerm_monitor_activity_log_alert` - support for `Security` event type for Azure Service Health alerts ([#11802](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11802))
* `azurerm_netapp_volume` - support for the `security_style` property - ([#11684](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11684))
* `azurerm_redis_cache` - suppot for the `replicas_per_master` peoperty ([#11714](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11714))
* `azurerm_spring_cloud_service` - support for the `required_network_traffic_rules` block ([#11633](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11633))
* `azurerm_storage_account_management_policy` - the `name` property can now contain `-` ([#11792](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11792))

BUG FIXES:

* `azurerm_frontdoor` - added a check for `nil` to avoid panic on destroy ([#11720](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11720))
* `azurerm_linux_virtual_machine_scale_set` - the `extension` blocks are now a set ([#11425](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11425))
* `azurerm_virtual_network_gateway_connection` - fix a bug where `shared_key` was not being updated ([#11742](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11742))
* `azurerm_windows_virtual_machine_scale_set` - the `extension` blocks are now a set ([#11425](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11425))
* `azurerm_windows_virtual_machine_scale_set` - changing the `license_type` will no longer create a new resource ([#11731](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11731))

---

For information on changes between the v2.59.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
