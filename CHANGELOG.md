## 2.41.0 (Unreleased)

FEATURES:

* **New Resource:** `azurerm_eventgrid_system_topic_event_subscription` [GH-9852]
* **New Resource:** `azurerm_media_job` [GH-9859]
* **New Resource:** `azurerm_media_streaming_endpoint` [GH-9537]
* **New Resource:** `azurerm_subnet_service_endpoint_storage_policy` [GH-8966]
* **New Resource:** `azurerm_synapse_managed_private_endpoint` [GH-9260]

IMPROVEMENTS:

* `azurerm_app_service` - Add support for `outbound_ip_address_list` and `possible_outbound_ip_address_list` [GH-9871]
* `azurerm_disk_encryption_set` - support for updating `key_vault_key_id` [GH-7913]
* `azurerm_iot_time_series_insights_gen2_environment` - exposing `data_access_fqdn` [GH-9848]
* `azurerm_log_analytics_linked_service` - Add new fields `workspace_id`, `read_access_id`, and `write_access_id` [GH-9410]
* `azurerm_subnet` - support for the `service_endpoint_policy` block [GH-8966]
* `azurerm_traffic_manager_profile` - support for new field `max_return` and support for `traffic_routing_method` to be `MultiValue` [GH-9487]

BUG FIXES:

* `azurerm_key_vault_certificate` - reading `dns_names` and `emails` within the `subject_alternative_names` block from the Certificate if not returned from the API [GH-8631]
* `azurerm_log_analytics_workspace` - adding a state migration to correctly update the Resource ID [GH-9853]

---

For information on changes between the v2.40.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
