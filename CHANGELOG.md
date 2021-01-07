## 2.42.0 (Unreleased)

FEATURES:

* **New Data Source:** `azurerm_eventgrid_domain_topic` [GH-10050]
* **New Resource:** `azurerm_data_factory_linked_service_synapse` [GH-9928]
* **New Resource:** `azurerm_disk_access` [GH-9889]
* **New Resource:** `azurerm_media_streaming_locator` [GH-9992]
* **New Resource:** `azurerm_sentinel_alert_rule_fusion` [GH-9829]

IMPROVEMENTS:

* batch: updating to API version `2020-03-01` [GH-10036]
* dependencies: upgrading to `v49.2.0` of `github.com/Azure/azure-sdk-for-go` [GH-10042]
* dependencies: upgrading to `v0.15.1` of `github.com/tombuildsstuff/giovanni` [GH-10035]
* Data Source: `azurerm_hdinsight_cluster` - support for the `kafka_rest_proxy_endpoint` property [GH-8064]
* Data Source: `azurerm_databricks_workspace` - support for the `tags` property [GH-9933]
* Data Source: `azurerm_subscription` - support for the `tags` property [GH-8064]
* `azurerm_batch_pool` support for the `public_address_provisioning_type` property [GH-10036]
* `azurerm_api_management` - support `Consumption_0` for the `sku_name` property [GH-6868]
* `azurerm_cdn_endpoint` - only send `content_types_to_compress` and `geo_filter` to the API when actually set [GH-9902]
* `azurerm_cosmosdb_mongo_collection` - correctly read back the `_id` index when mongo 3.6 [GH-8690]
* `azurerm_container_group` - support for the `volume.empty_dir` property [GH-9836]
* `azurerm_data_factory_linked_service_azure_file_storage` - support for the `file_share` property [GH-9934]
* `azurerm_dedicated_host` - support for addtional `sku_name` values [GH-9951]
* `azurerm_devspace_controller` - deprecating since new DevSpace Controllers can no longer be provisioned, this will be removed in version 3.0 of the Azure Provider [GH-10049]
* `azurerm_function_app` - make `pre_warmed_instance_count` computed to use azure's default [GH-9069]
* `azurerm_hdinsight_hadoop_cluster` - allow the value `Standard_D4a_V4` for the `vm_type` property [GH-10000]
* `azurerm_hdinsight_kafka_cluster` - support for the `rest_proxy` and `kafka_management_node` blocks [GH-8064]
* `azurerm_log_analytics_linked_service` - add validation for resource ID type [GH-9932]
* `azurerm_log_analytics_linked_service` - update validation to use generated validate functions [GH-9950]
* `azurerm_monitor_diagnostic_setting` - validation that `eventhub_authorization_rule_id` is a EventHub Namespace Authorization Rule ID [GH-9914]
* `azurerm_monitor_diagnostic_setting` - validation that `log_analytics_workspace_id` is a Log Analytics Workspace ID [GH-9914]
* `azurerm_monitor_diagnostic_setting` - validation that `storage_account_id` is a Storage Account ID [GH-9914]
* `azurerm_network_security_rule` - increase allowed the number of `application_security_group` blocks allowed [GH-9884]
* `azurerm_sentinel_alert_rule_ms_security_incident` - support the `alert_rule_template_guid` and `display_name_exclude_filter` properties [GH-9797]
* `azurerm_sentinel_alert_rule_scheduled` - support for the `alert_rule_template_guid` property [GH-9712]
* `azurerm_sentinel_alert_rule_scheduled` - support for creating incidents [GH-8564]
* `azurerm_synapse_workspace` - support for the `managed_resource_group_name` property [GH-10017]
* `azurerm_traffic_manager_profile` - support for the `traffic_view_enabled` property [GH-10005]

BUG FIXES:

provider: will not correctly register the `Microsoft.Blueprint` and `Microsoft.HealthcareApis` RPs [GH-10062]
* `azurerm_application_gateway` - allow `750` for `file_upload_limit_mb` when the sku is `WAF_v2` [GH-8753]
* `azurerm_firewall_policy_rule_collection_group` - correctly validate the `network_rule_collection.destination_ports` property [GH-9490]
* `azurerm_cdn_endpoint` - changing many `delivery_rule` condition `match_values` to optional [GH-8850]
* `azurerm_cosmosdb_account` - always include `key_vault_id` in update requests for azure policy enginer compatibility [GH-9966]
* `azurerm_cosmosdb_table` - do not call the throughput api when serverless [GH-9749]
* `azurerm_kubernetes_cluster` - parse oms `log_analytics_workspace_id` to ensure correct casing [GH-9976]
* `azurerm_role_assignment` fix crash in retry logic [GH-10051]
* `azurerm_storage_account` - allow hns when `account_tier` is `Premium` [GH-9548]
* `azurerm_storage_share_file` - allowing files smaller than 4KB to be uploaded [GH-10035]

## 2.41.0 (December 17, 2020)

UPGRADE NOTES:

* `azurerm_key_vault` - Azure will be introducing a breaking change on December 31st, 2020 by force-enabling Soft Delete on all new and existing Key Vaults. To workaround this, this release of the Azure Provider still allows you to configure Soft Delete on before this date (but once this is enabled this cannot be disabled). Since new Key Vaults will automatically be provisioned using Soft Delete in the future, and existing Key Vaults will be upgraded - a future release will deprecate the `soft_delete_enabled` field and default this to true early in 2021. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_certificate` - Terraform will now attempt to `purge` Certificates during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_key` - Terraform will now attempt to `purge` Keys during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_secret` - Terraform will now attempt to `purge` Secrets during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))

FEATURES:

* **New Resource:** `azurerm_eventgrid_system_topic_event_subscription` ([#9852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9852))
* **New Resource:** `azurerm_media_job` ([#9859](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9859))
* **New Resource:** `azurerm_media_streaming_endpoint` ([#9537](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9537))
* **New Resource:** `azurerm_subnet_service_endpoint_storage_policy` ([#8966](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8966))
* **New Resource:** `azurerm_synapse_managed_private_endpoint` ([#9260](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9260))

IMPROVEMENTS:

* `azurerm_app_service` - Add support for `outbound_ip_address_list` and `possible_outbound_ip_address_list` ([#9871](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9871))
* `azurerm_disk_encryption_set` - support for updating `key_vault_key_id` ([#7913](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7913))
* `azurerm_iot_time_series_insights_gen2_environment` - exposing `data_access_fqdn` ([#9848](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9848))
* `azurerm_key_vault_certificate` - performing a "purge" of the Certificate during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_key` - performing a "purge" of the Key during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_secret` - performing a "purge" of the Secret during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_log_analytics_linked_service` - Add new fields `workspace_id`, `read_access_id`, and `write_access_id` ([#9410](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9410))
* `azurerm_linux_virtual_machine` - Normalise SSH keys to cover VM import cases ([#9897](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9897))
* `azurerm_subnet` - support for the `service_endpoint_policy` block ([#8966](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8966))
* `azurerm_traffic_manager_profile` - support for new field `max_return` and support for `traffic_routing_method` to be `MultiValue` ([#9487](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9487))

BUG FIXES:

* `azurerm_key_vault_certificate` - reading `dns_names` and `emails` within the `subject_alternative_names` block from the Certificate if not returned from the API ([#8631](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8631))
* `azurerm_key_vault_certificate` - polling until the Certificate is fully deleted during deletion ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_key` - polling until the Key is fully deleted during deletion ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_key_vault_secret` -  polling until the Secret is fully deleted during deletion ([#9911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9911))
* `azurerm_log_analytics_workspace` - adding a state migration to correctly update the Resource ID ([#9853](https://github.com/terraform-providers/terraform-provider-azurerm/issues/9853))

---

For information on changes between the v2.40.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
