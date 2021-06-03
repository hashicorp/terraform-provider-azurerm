## 2.62.0 (Unreleased)

FEATURES:

* **New Resource** `azurerm_data_protection_backup_vault` [GH-11955]
* **New Resource** `azurerm_postgresql_flexible_server_firewall_rule` [GH-11834]
* **New Resource** `azurerm_vmware_express_route_authorization` [GH-11812]
* **New Resource** `azurerm_storage_object_replication_policy` [GH-11744]

ENHANCEMENTS:

* `azurerm_app_service_environment` - support for the `internal_ip_address`, `service_ip_address`, and `outbound_ip_addresses`properties [GH-12026]
* `azurerm_api_management_api_subscription` - support for the `api_id` property [GH-12025]
* `azurerm_container_registry` - support for  versionless encryption keys for ACR [GH-11856]
* `azurerm_kubernetes_cluster` -  support for `gateway_name` for Application Gateway add-on [GH-11984]
* `azurerm_kubernetes_cluster` - support update of `azure_rbac_enabled` [GH-12029]
* `azurerm_kubernetes_cluster` - support for `node_public_ip_prefix_id` [GH-11635]
* `azurerm_kubernetes_cluster_node_pool` - support for `node_public_ip_prefix_id` [GH-11635]
* `azurerm_machine_learning_inference_cluster` - support for the `ssl.leaf_domain_label` and `ssl.overwrite_existing_domain` properties [GH-11830]
* `azurerm_network_watcher_flow_log` - support for the `location` and `tags` properties [GH-11670]
* `azurerm_role_assignment` - support the `delegated_managed_identity_resource_id` property [GH-11848]

BUG FIXES:

* `azuerrm_postgres_server` - do no update `password` unless its changed [GH-12008]
* `azuerrm_storage_acount` - prevent `containerDeleteRetentionPolicy` and `lastAccessTimeTrackingPolicy` not supported in `AzureUSGovernment` errors [GH-11960]

## 2.61.0 (May 27, 2021)

FEATURES:

* **New Data Source:** `azurerm_spatial_anchors_account` ([#11824](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11824))

ENHANCEMENTS:

* dependencies: updating to `v54.3.0` of `github.com/Azure/azure-sdk-for-go` ([#11813](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11813))
* dependencies: updating `mixedreality` to use API Version `2021-01-01` ([#11824](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11824))
* refactor: switching to use an embedded SDK for `appconfiguration` ([#11959](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11959))
* refactor: switching to use an embedded SDK for `eventhub` ([#11973](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11973))
* provider: support for the Virtual Machine `skip_shutdown_and_force_delete` feature ([#11216](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11216))
* provider: support for the Virtual Machine Scale Set `force_delete` feature ([#11216](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11216))
* provider: no longer auto register the Microsoft.DevSpaces RP ([#11822](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11822))
* Data Source: `azurerm_key_vault_certificate_data` - support certificate bundles and add support for ECDSA keys ([#11974](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11974))
* `azurerm_data_factory_linked_service_sftp` - support for hostkey related properties ([#11825](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11825))
* `azurerm_spatial_anchors_account` - support for `account_domain` and `account_id` ([#11824](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11824))
* `azurerm_static_site` - Add support for `tags` attribute ([#11849](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11849))
* `azurerm_storage_account` - `private_link_access` supports more values ([#11957](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11957))
* `azurerm_storage_account_network_rules`: `private_link_access` supports more values ([#11957](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11957))
* `azurerm_synapse_spark_pool` - `spark_version` now supports `3.0` ([#11972](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11972))

BUG FIXES:

* `azurerm_cdn_endpoint` - do not send an empty `origin_host_header` to the api ([#11852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11852))
* `azurerm_linux_virtual_machine_scale_set`: changing the `disable_automatic_rollback` and `enable_automatic_os_upgrade` properties no longer created a new resource ([#11723](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11723))
* `azurerm_storage_share`: Fix ID for `resource_manager_id` ([#11828](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11828))
* `azurerm_windows_virtual_machine_scale_set`: changing the `disable_automatic_rollback` and `enable_automatic_os_upgrade` properties no longer created a new resource ([#11723](https://github.com/terraform-providers/terraform-provider-azurerm/issues/11723))

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
