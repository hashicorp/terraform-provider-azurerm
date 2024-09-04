## 4.1.0 (Unreleased)

BUG FIXES:

* `azurerm_cosmosdb_account` - the `ip_range_filter` property now supports IPV4 addresses [GH-27208]
* `azurerm_linux_virtual_machine` - the `admin_ssh_key.public_key` property now supports ed25519 ssh keys [GH-27202]

ENHANCEMENTS:

* `azurerm_*_virtual_machine_scale_set` - upgrade api version from `2024-03-01` to `2024-07-01` [GH-27230]
* `azurerm_api_management_logger` - support for the `application_insights.connection_string` property [GH-27137]

## 4.0.1 (August 23, 2024)

BUG FIXES:


* provider: fix a validation bug that prevents `terraform validate` from working when `subscription_id` is not specified ([#27178](https://github.com/hashicorp/terraform-provider-azurerm/issues/27178))
* `azurerm_cognitive_deployment` - fixed replacement of `scale` block with `sku` ([#27173](https://github.com/hashicorp/terraform-provider-azurerm/issues/27173))
* `azurerm_kubernetes_cluster` - prevent a panic ([#27183](https://github.com/hashicorp/terraform-provider-azurerm/issues/27183))
* `azurerm_kubernetes_cluster_node_pool` - prevent a panic caused by renamed `enable_*` properties ([#27164](https://github.com/hashicorp/terraform-provider-azurerm/issues/27164))
* `azurerm_sentinel_data_connector_microsoft_threat_intelligence` - prevent error by removing deprecated property `bing_safety_phishing_url_lookback_date` ([#27171](https://github.com/hashicorp/terraform-provider-azurerm/issues/27171))

## 4.0.0 (August 22, 2024)

NOTES:

* **Major Version**: Version 4.0 of the Azure Provider is a major version - some behaviours have changed and some deprecated fields/resources have been removed - please refer to [the 4.0 upgrade guide for more information](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/4.0-upgrade-guide).
* When upgrading to v4.0 of the AzureRM Provider, we recommend upgrading to the latest version of Terraform Core ([which can be found here](https://www.terraform.io/downloads)).

ENHANCEMENTS:

* Data Source: `azurerm_shared_image` - add support for the `trusted_launch_supported`, `trusted_launch_enabled`, `confidential_vm_supported`, `confidential_vm_enabled`, `accelerated_network_support_enabled` and `hibernation_enabled` properties ([#26975](https://github.com/hashicorp/terraform-provider-azurerm/issues/26975))
* dependencies: updating `hashicorp/go-azure-sdk` to `v0.20240819.1075239` ([#27107](https://github.com/hashicorp/terraform-provider-azurerm/issues/27107))
* `applicationgateways` - updating to use `2023-11-01` ([#26776](https://github.com/hashicorp/terraform-provider-azurerm/issues/26776))
* `containerregistry` - updating to use `2023-06-01-preview` ([#23393](https://github.com/hashicorp/terraform-provider-azurerm/issues/23393))
* `containerservice` - updating to `2024-05-01` ([#27105](https://github.com/hashicorp/terraform-provider-azurerm/issues/27105))
* `mssql` - updating to use `hashicorp/go-azure-sdk` and `023-08-01-preview` ([#27073](https://github.com/hashicorp/terraform-provider-azurerm/issues/27073))
* `mssqlmanagedinstance` - updating to use `hashicorp/go-azure-sdk` and `2023-08-01-preview` ([#26872](https://github.com/hashicorp/terraform-provider-azurerm/issues/26872))
* `azurerm_image` - add support for the `disk_encryption_set_id` property to the `data_disk` block ([#27015](https://github.com/hashicorp/terraform-provider-azurerm/issues/27015))
* `azurerm_log_analytics_workspace_table` - add support for more `total_retention_in_days` and `retention_in_days` values ([#27053](https://github.com/hashicorp/terraform-provider-azurerm/issues/27053))
* `azurerm_mssql_elasticpool` - add support for the `HS_MOPRMS` and `MOPRMS` skus ([#27085](https://github.com/hashicorp/terraform-provider-azurerm/issues/27085))
* `azurerm_netapp_pool` - allow `1` as a valid value for `size_in_tb` ([#27095](https://github.com/hashicorp/terraform-provider-azurerm/issues/27095))
* `azurerm_notification_hub` - add support for the `browser_credential` property ([#27058](https://github.com/hashicorp/terraform-provider-azurerm/issues/27058))
* `azurerm_redis_cache` - add support for the `access_keys_authentication_enabled` property ([#27039](https://github.com/hashicorp/terraform-provider-azurerm/issues/27039))
* `azurerm_role_assignment` - add support for the `/`, `/providers/Microsoft.Capacity` and `/providers/Microsoft.BillingBenefits` scopes ([#26663](https://github.com/hashicorp/terraform-provider-azurerm/issues/26663))
* `azurerm_shared_image` - add support for the `hibernation_enabled` property ([#26975](https://github.com/hashicorp/terraform-provider-azurerm/issues/26975))
* `azurerm_storage_account` - support `queue_encryption_key_type` and `table_encryption_key_type` for more storage account kinds ([#27112](https://github.com/hashicorp/terraform-provider-azurerm/issues/27112))
* `azurerm_web_application_firewall_policy` - add support for the `request_body_enforcement` property ([#27094](https://github.com/hashicorp/terraform-provider-azurerm/issues/27094))

BUG FIXES:

* `azurerm_ip_group_cidr` - fixed the position of the CIDR check to correctly refresh the resource when it's no longer present ([#27103](https://github.com/hashicorp/terraform-provider-azurerm/issues/27103))
* `azurerm_monitor_diagnostic_setting` - add further polling to work around an eventual consistency issue when creating the resource ([#27088](https://github.com/hashicorp/terraform-provider-azurerm/issues/27088))
* `azurerm_storage_account` - prevent API error by populating `infrastructure_encryption_enabled` when updating `customer_managed_key` ([#26971](https://github.com/hashicorp/terraform-provider-azurerm/issues/26971))
* `azurerm_storage_blob_inventory_policy` - the `filter` property can now be set when `scope` is `container` ([#27113](https://github.com/hashicorp/terraform-provider-azurerm/issues/27113))
* `azurerm_virtual_network_dns_servers` - moved locks to prevent the creation of subnets with stale data ([#27036](https://github.com/hashicorp/terraform-provider-azurerm/issues/27036))
* `azurerm_virtual_network_gateway_connection` - allow `0` as a valid value for `ipsec_policy.sa_datasize` ([#27056](https://github.com/hashicorp/terraform-provider-azurerm/issues/27056))

---

For information on changes between the v3.116.0 and v3.0.0 releases, please see [the previous v3.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v3.md).

For information on changes between the v2.99.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v2.md).

For information on changes between the v1.44.0 and v1.0.0 releases, please see [the previous v1.x changelog entries](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v1.md).

For information on changes prior to the v1.0.0 release, please see [the v0.x changelog](https://github.com/hashicorp/terraform-provider-azurerm/blob/main/CHANGELOG-v0.md).
