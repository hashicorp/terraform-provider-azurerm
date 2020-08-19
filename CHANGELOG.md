## 2.24.0 (Unreleased)

IMPROVEMENTS:

* dependencies: update `containerinstance` to API version `2019-12-01` [GH-8110]
* `azurerm_api_management_api` - now supports `oauth2_authorization` and `openid_authentication` [GH-7617]

BUG FIXES:

* `azurerm_frontdoor` - ensuring all fields are set into the state [GH-8146]
* `azurerm_frontdoor` - rewriting case-inconsistent Resource ID's to ensure they're reliable [GH-8146]
* `azurerm_frontdoor_firewall_policy` - ensuring all fields are set into the state [GH-8146]
* `azurerm_frontdoor_firewall_policy` - rewriting case-inconsistent Resource ID's to ensure they're reliable [GH-8146]
* `azurerm_frontdoor_custom_https_configuration` - ensuring all fields are set into the state [GH-8146]
* `azurerm_frontdoor_custom_https_configuration` - rewriting case-inconsistent Resource ID's to ensure they're reliable [GH-8146]
* `azurerm_frontdoor_custom_https_configuration` - updating the ID to use the frontendEndpoint's Resource ID rather than a custom Resource ID [GH-8146]
* `azurerm_lb` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` [GH-8006]
* `azurerm_lb_backend_address_pool` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` [GH-8006]
* `azurerm_lb_nat_pool` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` [GH-8006]
* `azurerm_lb_nat_rule` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` [GH-8006]
* `azurerm_lb_probe` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` [GH-8006]
* `azurerm_lb_outbound_rule` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` [GH-8006]
* `azurerm_lb_rule` - switching to use API version `2020-03-01` to workaround a bug in API version `2020-05-01` [GH-8006]
* `azurerm_storage_account` - only sending `allow_blob_public_access` and `min_tls_version` in Azure Public since these are currently not supported in other regions [GH-8148]

## 2.23.0 (August 13, 2020)

FEATURES:

* **New Resource:** `azurerm_integration_service_environment` ([#7763](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7763))
* **New Resource:** `azurerm_redis_linked_server` ([#8026](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8026))
* **New Resource:** `azurerm_synapse_firewall_rule` ([#7904](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7904))

IMPROVEMENTS:

* Data Source: `azurerm_app_service` now exports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* Data Source: `azurerm_function_app` now exports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* Data Source: `azurerm_function_app` now exports `site_config` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_app_service` now supports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports `source_control` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports full `ip_restriction` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_function_app` now supports full `scm_ip_restriction` configuration ([#7945](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7945))
* `azurerm_site_recovery_replicated_vm` - support setting `target_network_id` and `network_interface` on failover ([#5688](https://github.com/terraform-providers/terraform-provider-azurerm/issues/5688))
* `azurerm_storage_account` - support `static_website` for `BlockBlobStorage` account type ([#7890](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7890))
* `azurerm_storage_account` - filter `allow_blob_public_access` and `min_tls_version` from Azure US Government ([#8092](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8092))

ENHANCEMENTS:

* dependencies: updating `containerservice` to `2020-04-01` ([#7894](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7894))
* dependencies: updating `mysql` to `2020-01-01` ([#8062](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8062))
* dependencies: updating `postgresql` to `2020-01-01` ([#8045](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8045))
* `azurerm_eventhub_namespace` - support for the `identity` block ([#8065](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8065))
* `azurerm_postgresql_server` - support for the `identity` block ([#8044](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8044))

BUG FIXES:

* All resources using a `location` field - allowing the value `global` when using enhanced validation ([#8042](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8042))
* Data Source: `azurerm_api_management_user` - `user_id` now accepts single characters ([#7975](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7975))
* `azurerm_application_gateway` - enforce case for the `rule_type` property ([#8061](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8061))
* `azurerm_iothub_consumer_group` - lock during creation and deletion to workaround an API issue ([#8041](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8041))
* `azurerm_iothub` - the `endpoint` and `route` lists can now be cleared by setting them to `[]` ([#8028](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8028))
* `azurerm_linux_virtual_machine` - handling machines which are already stopped/deallocated ([#8000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8000))
* `azurerm_mariadb_virtual_network_rule` will now work across subscriptions ([#8100](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8100))
* `azurerm_monitor_metric_alert_resource` - continue using `SingleResourceMultiMetricCriteria` for existing alerts ([#7995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7995))
* `azurerm_mysql_server` - prevent a non empty plan when using `threat_detection_policy` ([#7981](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7981))
* `azurerm_orchestrated_virtual_machine_scale_set` - allow `single_placement_group` to be `true` ([#7821](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7821))
* `azurerm_mysql_server` - support for the `identity` block ([#8059](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8059))
* `azurerm_storage_account` - set default for `min_tls_version` to `TLS_10` ([#8152](https://github.com/terraform-providers/terraform-provider-azurerm/pull/8152))
* `azurerm_traffic_manager_profile` - updating no longer clears all endpoints ([#7846](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7846))
* `azurerm_windows_virtual_machine` - handling machines which are already stopped/deallocated [[#8000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8000)]'
* `azurerm_data_factory_dataset_delimited_text` - fix issue with property `azure_blob_storage_account` ([#7953](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7953))

## 2.22.0 (August 07, 2020)

DEPENDENCIES:

* updating `github.com/Azure/azure-sdk-for-go` to `v44.2.0` ([#7933](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7933))

IMPROVEMENTS:

* `azurerm_cosmosdb_account` - support `DisableRateLimitingResponses` with the `capabilities` property ([#8016](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8016))
* `azurerm_storage_account` - support for the `min_tls_version` property ([#7879](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7879))
* `azurerm_storage_account_sas` - support for the `signed_version attribute` property ([#8020](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8020))
* `azurerm_servicebus_queue` - support for the `enable_batched_operations`, `status`, `forward_to`, and `forward_dead_lettered_messages_to` ([#7990](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7990))

BUG FIXES:

* Data Source: `azurerm_key_vault_certificate` - fixing a crash when using acmebot certificates ([#8029](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8029))
* `azurerm_iothub_shared_access_policy` - prevent `primary_connection_string` & `secondary_connection_string` from regenerating during every apply ([#8017](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8017))

## 2.21.0 (July 31, 2020)

DEPENDENCIES:

* updating `search` to `2020-03-13` ([#7867](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7867))
* updating `go-azure-helpers` to `v0.11.2` ([#7911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7911))

FEATURES:

* **New Data Source:** `azurerm_data_share_dataset_data_lake_gen1` ([#7840](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7840))
* **New Resource:** `azurerm_dedicated_hardware_security_module` ([#7727](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7727))

IMPROVEMENTS:
* `azurerm_api_management_identity_provider_aad` - Support for `signin_tenant` ([#7901](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7901))
* `azurerm_app_service_plan` - update the relation between `kind` and `reserved` ([#7943](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7943))
* `azurerm_automation_runbook` - recreate `azurerm_automation_job_schedule` after an update ([#7555](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7555))
* `azurerm_app_service_slot` - support for the `application_logs.file_system` ([#7311](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7311))
* `azurerm_firewall` - no longer requires a `zone` ([#7817](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7817))
* `azurerm_function_app_slot` - support for the `site_config.auto_swap_slot_name` property ([#7859](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7859))
* `azurerm_kubernetes_cluster` - support for in-place upgrade from `Free` to `Paid` for `sku_tier` ([#7927](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7927))
* `azurerm_monitor_scheduled_query_rules_alert` - `action.0.custom_webhook_payload` is now sent as empty to allow for Azure's default to take effect([#7838](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7838))
* `azurerm_search_service` - support for the `public_network_access_enabled` property ([#7867](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7867))
* `azurerm_servicebus_subscription` - support for the `status` property ([#7852](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7852))

BUG FIXES:

* `azurerm_automation_runbook` - allow `publish_content_link` resource to not be set ([#7824](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7824))
* `azurerm_api_management_named_value` - the `value` has been marked as sensitive to hide secret values ([#7819](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7819))
* `azurerm_cognitive_account` - allow `qname_runtime_endpoint` to not be set ([#7916](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7916))
* `azurerm_iothub_dps` - the only valid value for the `sku` property for the API is now `S1` ([#7847](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7847))
* `azurerm_eventgrid_event_subscription` - deprecate the `topic_name` as it is now readonly in the API ([#7871](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7871))
* `azurerm_kubernetes_cluster` - updates will no longer fail when using managed AAD integration ([#7874](https://github.com/terraform-providers/terraform-provider-azurerm/issues/7874))

---

For information on changes between the v2.20.0 and v2.0.0 releases, please see [the previous v2.x changelog entries](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v2.md).

For information on changes in version v1.44.0 and prior releases, please see [the v1.44.0 changelog](https://github.com/terraform-providers/terraform-provider-azurerm/blob/master/CHANGELOG-v1.md).
