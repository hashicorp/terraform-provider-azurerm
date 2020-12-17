## 2.41.0 (Unreleased)

UPGRADE NOTES:

* `azurerm_key_vault` - Azure will be introducing a breaking change on December 31st, 2020 by force-enabling Soft Delete on all new and existing Key Vaults. To workaround this, this release of the Azure Provider still allows you to configure Soft Delete on before this date (but once this is enabled this cannot be disabled). Since new Key Vaults will automatically be provisioned using Soft Delete in the future, and existing Key Vaults will be upgraded - a future release will deprecate the `soft_delete_enabled` field and default this to true early in 2021. [GH-9911]
* `azurerm_key_vault_certificate` - Terraform will now attempt to `purge` Certificates during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. [GH-9911]
* `azurerm_key_vault_key` - Terraform will now attempt to `purge` Keys during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. [GH-9911]
* `azurerm_key_vault_secret` - Terraform will now attempt to `purge` Secrets during deletion due to the upcoming breaking change in the Azure API where Key Vaults will have soft-delete force-enabled. This can be disabled by setting the `purge_soft_delete_on_destroy` field within the `features -> keyvault` block to `false`. [GH-9911]

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
* `azurerm_key_vault_certificate` - performing a "purge" of the Certificate during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information [GH-9911]
* `azurerm_key_vault_key` - performing a "purge" of the Key during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information [GH-9911]
* `azurerm_key_vault_secret` - performing a "purge" of the Secret during deletion if the feature is opted-in within the `features` block, see the "Upgrade Notes" for more information [GH-9911]
* `azurerm_log_analytics_linked_service` - Add new fields `workspace_id`, `read_access_id`, and `write_access_id` [GH-9410]
* `azurerm_linux_virtual_machine` - Normalise SSH keys to cover VM import cases [GH-9897]
* `azurerm_subnet` - support for the `service_endpoint_policy` block [GH-8966]
* `azurerm_traffic_manager_profile` - support for new field `max_return` and support for `traffic_routing_method` to be `MultiValue` [GH-9487]

BUG FIXES:

* `azurerm_key_vault_certificate` - reading `dns_names` and `emails` within the `subject_alternative_names` block from the Certificate if not returned from the API [GH-8631]
* `azurerm_key_vault_certificate` - polling until the Certificate is fully deleted during deletion [GH-9911]
* `azurerm_key_vault_key` - polling until the Key is fully deleted during deletion [GH-9911]
* `azurerm_key_vault_secret` -  polling until the Secret is fully deleted during deletion [GH-9911]
* `azurerm_log_analytics_workspace` - adding a state migration to correctly update the Resource ID [GH-9853]

---

For information on changes between the v2.40.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.x changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
