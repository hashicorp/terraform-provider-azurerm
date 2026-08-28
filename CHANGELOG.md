## 5.7.0 (Unreleased)

ENHANCEMENTS:

FEATURES:

BUG FIXES:

## 5.6.0 (Unreleased)

ENHANCEMENTS:
* Data Source: `azurerm_api_management_workspace` - export the `description` attribute [GH-33205]

FEATURES:

BUG FIXES:

## 5.5.0 (Unreleased)

ENHANCEMENTS:

FEATURES:
* New List Resource: azurerm_availability_set [GH-33241]

BUG FIXES:

## 5.4.0 (Unreleased)

ENHANCEMENTS:
* dependencies: `azurerm_mongo_cluster` - update API version to `2026-06-01` [GH-33195]
* dependencies: `azurerm_mongo_cluster_firewall_rule` - update API version to `2026-06-01` [GH-33195]
* dependencies: `azurerm_mongo_cluster_user` - update API version to `2026-06-01` [GH-33195]
* `azurerm_postgresql_flexible_server` - add new arguments `storage_type`, `storage_iops`, and `storage_throughput` which allows choice of the new "Premium V2 LRS" storage type [GH-32121]

FEATURES:

BUG FIXES:

## 5.3.0 (August 27, 2026)

FEATURES:

* **New Data Source**: `azurerm_playwright_workspace` ([#31954](https://github.com/hashicorp/terraform-provider-azurerm/issues/31954))
* **New List Resource**: `azurerm_cognitive_deployment` ([#33149](https://github.com/hashicorp/terraform-provider-azurerm/issues/33149))
* **New List Resource**: `azurerm_playwright_workspace` ([#31954](https://github.com/hashicorp/terraform-provider-azurerm/issues/31954))
* **New Resource**: `azurerm_playwright_workspace` ([#31954](https://github.com/hashicorp/terraform-provider-azurerm/issues/31954))

ENHANCEMENTS:

* dependencies: `go-azure-helpers` - update version to `0.82.0` ([#33142](https://github.com/hashicorp/terraform-provider-azurerm/issues/33142))
* dependencies: `sql` - update API version to `2025-01-01` ([#33201](https://github.com/hashicorp/terraform-provider-azurerm/issues/33201))
* Data Source: `azurerm_role_definition` - export the `role_definition_resource_id` property ([#33126](https://github.com/hashicorp/terraform-provider-azurerm/issues/33126))
* `azurerm_cognitive_deployment` - add Resource Identity support ([#33149](https://github.com/hashicorp/terraform-provider-azurerm/issues/33149))
* `azurerm_federated_identity_credential` - add additional polling to account for Azure's eventual consistency ([#32935](https://github.com/hashicorp/terraform-provider-azurerm/issues/32935))
* `azurerm_kubernetes_cluster` - add support for the `oms_agent.retina_flow_logs_enabled` property ([#33222](https://github.com/hashicorp/terraform-provider-azurerm/issues/33222))
* `azurerm_managed_application` - add support for the `identity` block ([#30725](https://github.com/hashicorp/terraform-provider-azurerm/issues/30725))
* `azurerm_private_endpoint` - extend validation for the `private_service_connection.subresource_names` property to allow names containing spaces ([#32887](https://github.com/hashicorp/terraform-provider-azurerm/issues/32887))
* `azurerm_search_service` - allow in-place downgrades of the `sku` property between Basic and Standard tiers ([#33069](https://github.com/hashicorp/terraform-provider-azurerm/issues/33069))
* `azurerm_site_recovery_replicated_vm` - add update support to the `managed_disk` block without requiring resource recreation ([#33140](https://github.com/hashicorp/terraform-provider-azurerm/issues/33140))
* `azurerm_user_assigned_identity` - add additional polling to account for Azure's eventual consistency ([#33142](https://github.com/hashicorp/terraform-provider-azurerm/issues/33142))

BUG FIXES:

* Data Source: `azurerm_app_configuration_key` - now correctly sets `tags` into state ([#33182](https://github.com/hashicorp/terraform-provider-azurerm/issues/33182))
* `azurerm_eventhub_namespace` - prevent `network_rulesets.x.default_action` being set to `Deny` if `ip_rule` or `virtual_network_rule` is not specified ([#33216](https://github.com/hashicorp/terraform-provider-azurerm/issues/33216))

## 5.2.0 (August 20, 2026)

FEATURES:

* **New List Resource**: `azurerm_user_assigned_identity` ([#32667](https://github.com/hashicorp/terraform-provider-azurerm/issues/32667))

ENHANCEMENTS:

* dependencies: `go` - update to `1.26.6` ([#33141](https://github.com/hashicorp/terraform-provider-azurerm/issues/33141))
* dependencies: `go-azure-sdk` - update to `v0.20260811.1225050` ([#33079](https://github.com/hashicorp/terraform-provider-azurerm/issues/33079))
* `azurerm_cdn_frontdoor_batch_rule_set` - allow `/` as an input to `rule.conditions.request_path.values` ([#33023](https://github.com/hashicorp/terraform-provider-azurerm/issues/33023))
* `azurerm_databricks_workspace` - remove a redundant key vault existence check ([#33136](https://github.com/hashicorp/terraform-provider-azurerm/issues/33136))
* `azurerm_databricks_workspace_root_dbfs_customer_managed_key` - remove a redundant key vault existence check ([#33136](https://github.com/hashicorp/terraform-provider-azurerm/issues/33136))
* `azurerm_logic_app_standard` - add support for `v10.0` to `site_config.dotnet_framework_version` ([#33116](https://github.com/hashicorp/terraform-provider-azurerm/issues/33116))
* `azurerm_mongo_cluster` - `administrator_password` is no longer required when `create_mode` is `Default` to support Entra ID-only authentication ([#32092](https://github.com/hashicorp/terraform-provider-azurerm/issues/32092))
* `azurerm_redhat_openshift_cluster` - add support for the `network_profile.load_balancer_profile` block ([#32473](https://github.com/hashicorp/terraform-provider-azurerm/issues/32473))
* `azurerm_redhat_openshift_cluster` - add support for the `platform_workload_identity_profile` block ([#32473](https://github.com/hashicorp/terraform-provider-azurerm/issues/32473))
* `azurerm_role_assignment` - the `condition`, `condition_version`, and `description` properties can now be updated in-place ([#32714](https://github.com/hashicorp/terraform-provider-azurerm/issues/32714))
* `azurerm_snapshot` - `create_option` now supports `CopyStart` ([#32834](https://github.com/hashicorp/terraform-provider-azurerm/issues/32834))

BUG FIXES:

* `azurerm_cognitive_account_project` - added create/update/delete lock on parent AccountID to make sure operations on parent account are processed in serial (required by Cognitive service) ([#33151](https://github.com/hashicorp/terraform-provider-azurerm/issues/33151))
* `azurerm_databricks_workspace` - fix a persistent diff on removal of `managed_disk_cmk_key_vault_key_id` or `managed_services_cmk_key_vault_key_id` ([#33136](https://github.com/hashicorp/terraform-provider-azurerm/issues/33136))
* `azurerm_oracle_exadata_infrastructure` - fix an issue that prevented users from deploying with no `zones` set ([#33011](https://github.com/hashicorp/terraform-provider-azurerm/issues/33011))

## 5.1.0 (August 13, 2026)

ENHANCEMENTS:

* dependencies: `azurerm_linux_virtual_machine_scale_set` - update to API version `2025-04-01` ([#31586](https://github.com/hashicorp/terraform-provider-azurerm/issues/31586))
* dependencies: `azurerm_orchestrated_virtual_machine_scale_set` - update to API version `2025-04-01` ([#31586](https://github.com/hashicorp/terraform-provider-azurerm/issues/31586))
* dependencies: `azurerm_virtual_machine_scale_set` - update to API version `2025-04-01` ([#31586](https://github.com/hashicorp/terraform-provider-azurerm/issues/31586))
* dependencies: `azurerm_virtual_machine_scale_set_extension` - update to API version `2025-04-01` ([#31586](https://github.com/hashicorp/terraform-provider-azurerm/issues/31586))
* dependencies: `azurerm_windows_virtual_machine_scale_set` - update to API version `2025-04-01` ([#31586](https://github.com/hashicorp/terraform-provider-azurerm/issues/31586))
* dependencies: `codesigning` - update to API version `2025-10-13` ([#31714](https://github.com/hashicorp/terraform-provider-azurerm/issues/31714))
* `azurerm_linux_virtual_machine` - `encryption_at_host_enabled` can now be set to `true` when `os_disk.security_encryption_type` is set to `DiskWithVMGuestState` ([#32885](https://github.com/hashicorp/terraform-provider-azurerm/issues/32885))
* `azurerm_linux_virtual_machine_scale_set` - `encryption_at_host_enabled` can now be set to `true` when `os_disk.security_encryption_type` is set to `DiskWithVMGuestState` ([#32885](https://github.com/hashicorp/terraform-provider-azurerm/issues/32885))
* `azurerm_managed_devops_pool` - add support for the `CreatorOnly` value to `azure_devops_organization.permission.kind` property ([#32753](https://github.com/hashicorp/terraform-provider-azurerm/issues/32753))
* `azurerm_windows_virtual_machine` - `encryption_at_host_enabled` can now be set to `true` when `os_disk.security_encryption_type` is set to `DiskWithVMGuestState` ([#32885](https://github.com/hashicorp/terraform-provider-azurerm/issues/32885))
* `azurerm_windows_virtual_machine_scale_set` - `encryption_at_host_enabled` can now be set to `true` when `os_disk.security_encryption_type` is set to `DiskWithVMGuestState` ([#32885](https://github.com/hashicorp/terraform-provider-azurerm/issues/32885))

BUG FIXES:

* `azurerm_cdn_frontdoor_batch_ruleset` - parse `rule.actions.route_configuration_override.origin_group.cdn_frontdoor_origin_group_id` case-insensitively and normalize the resulting value to prevent diffs ([#32980](https://github.com/hashicorp/terraform-provider-azurerm/issues/32980))
* `azurerm_cdn_frontdoor_route` - parse `cdn_frontdoor_origin_group_id` case-insensitively and normalize the resulting value to prevent diffs ([#32980](https://github.com/hashicorp/terraform-provider-azurerm/issues/32980))
* `azurerm_cdn_frontdoor_secret` - fix an incorrect type assertion ([#32982](https://github.com/hashicorp/terraform-provider-azurerm/issues/32982))
* `azurerm_dev_center_project` - parse `dev_center_id` case-insensitively and normalize the resulting value to prevent diffs ([#32798](https://github.com/hashicorp/terraform-provider-azurerm/issues/32798))
* `azurerm_eventhub` - now prevents the `status` property from being set to `SendDisabled` on create  ([#33071](https://github.com/hashicorp/terraform-provider-azurerm/issues/33071))
* `azurerm_storage_container` - add a state migration for the `id` field, fixing the upgrade path from 4.x to 5.x ([#32978](https://github.com/hashicorp/terraform-provider-azurerm/issues/32978))
* `azurerm_storage_queue` - extend state migration to handle a malformed `resource_manager_id` ([#32979](https://github.com/hashicorp/terraform-provider-azurerm/issues/32979))
* `azurerm_storage_share` - add a state migration for the `id` field, fixing the upgrade path from 4.x to 5.x ([#33075](https://github.com/hashicorp/terraform-provider-azurerm/issues/33075))

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
