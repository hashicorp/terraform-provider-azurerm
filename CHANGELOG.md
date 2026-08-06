* `azurerm_cdn_frontdoor_route` - parse `cdn_frontdoor_origin_group_id` case-insensitively and normalize the resulting value to prevent diffs [GH-32980]
* `azurerm_cdn_frontdoor_secret` - fix an incorrect type assertion [GH-32982]
* `azurerm_storage_queue` - add a state migration to fix an ID parsing error when upgrading from 4.x to 5.x [GH-32914]
## 5.0.1 (July 30, 2026)

NOTES:

In addition to the bug fixes below, a number of resource documentation pages and the 5.0-upgrade-guide have been updated.

BUG FIXES:

* `azurerm_cdn_frontdoor_origin` - fix a regression that prevented valid values as input to `private_link.private_link_target_id` ([#32912](https://github.com/hashicorp/terraform-provider-azurerm/issues/32912))
* `azurerm_storage_queue` - add a state migration for the `id` field, fixing the upgrade path from 4.x to 5.x ([#32914](https://github.com/hashicorp/terraform-provider-azurerm/issues/32914))
* `azurerm_storage_table_entity` - add a state migration for the `storage_table_id` field, fixing the upgrade path from 4.x to 5.x ([#32929](https://github.com/hashicorp/terraform-provider-azurerm/issues/32929))

## 5.0.0 (July 27, 2026)

NOTES:

* **Major Version**: Version 5.0 of the Azure Provider is a major version - some behaviours have changed and some deprecated fields/resources have been removed - please refer to [the 5.0 upgrade guide for more information](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/5.0-upgrade-guide).
* When upgrading to v5.0 of the AzureRM Provider, we recommend upgrading to the latest version of Terraform Core ([which can be found here](https://developer.hashicorp.com/terraform/install)).

FEATURES:

* **New Action**: `azurerm_web_app_set_slot_distribution` ([#32364](https://github.com/hashicorp/terraform-provider-azurerm/issues/32364))
* **New Datasource** adds `azurerm_kubernetes_automatic_cluster_datasource` ([#32881](https://github.com/hashicorp/terraform-provider-azurerm/issues/32881))

ENHANCEMENTS:

* dependencies: `grpc` update to `1.82.1` ([#32852](https://github.com/hashicorp/terraform-provider-azurerm/issues/32852))
* dependencies: `loadbalancers` - update to API version `2025-01-01` ([#32644](https://github.com/hashicorp/terraform-provider-azurerm/issues/32644))
* `azurerm_cognitive_account_rai_policy` - the `content_filter.severity_threshold` property is now optional ([#32100](https://github.com/hashicorp/terraform-provider-azurerm/issues/32100))
* `azurerm_container_registry` - the `trust_policy_enabled` property has been deprecated and removed from the provider ([#32752](https://github.com/hashicorp/terraform-provider-azurerm/issues/32752))
* `azurerm_dashboard_grafana` - the `11` value for the `grafana_major_version` property has been deprecated and the property now supports `13` ([#32777](https://github.com/hashicorp/terraform-provider-azurerm/issues/32777))
* `azurerm_log_analytics_workspace` - add support for the `internet_ingestion_access_type` and `internet_query_access_type` properties ([#32562](https://github.com/hashicorp/terraform-provider-azurerm/issues/32562))
* `azurerm_subnet` - add support for the `network_security_group_id_wo` and `network_security_group_id_wo_version` properties ([#32847](https://github.com/hashicorp/terraform-provider-azurerm/issues/32847))
* `azurerm_subnet` - add support for the `route_table_id_wo` and `route_table_id_wo_version` properties ([#32847](https://github.com/hashicorp/terraform-provider-azurerm/issues/32847))
* `azurerm_subnet` - export the `network_security_group_id` property ([#32847](https://github.com/hashicorp/terraform-provider-azurerm/issues/32847))
* `azurerm_subnet` - export the `route_table_id` property ([#32847](https://github.com/hashicorp/terraform-provider-azurerm/issues/32847))
* `azurerm_windows_web_app` - add support for `~24` to `site_config.application_stack.node_version` ([#32840](https://github.com/hashicorp/terraform-provider-azurerm/issues/32840))
* `azurerm_windows_web_app_slot` - add support for `~24` to `site_config.application_stack.node_version` ([#32840](https://github.com/hashicorp/terraform-provider-azurerm/issues/32840))
* `cdn` - migrate to `go-azure-sdk` ([#32849](https://github.com/hashicorp/terraform-provider-azurerm/issues/32849))
* `sentinel` - migrate to `go-azure-sdk` ([#32759](https://github.com/hashicorp/terraform-provider-azurerm/issues/32759))
