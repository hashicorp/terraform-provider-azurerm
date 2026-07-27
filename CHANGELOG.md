## 5.0.0 (Unreleased)

NOTES:

* **Major Version** TODO

ENHANCEMENTS:
* `azurerm_subnet` - add support for the `network_security_group_id_wo` and `network_security_group_id_wo_version` properties [GH-32847]
* `azurerm_subnet` - add support for the `route_table_id_wo` and `route_table_id_wo_version` properties [GH-32847]
* `azurerm_subnet` - export the `network_security_group_id` property [GH-32847]
* `azurerm_subnet` - export the `route_table_id` property [GH-32847]
* dependencies: `loadbalancers` - update to API version `2025-01-01` [GH-32644]
* `cdn` - migrate to `go-azure-sdk` [GH-32849]
* `azurerm_windows_web_app` - add support for `~24` to `site_config.application_stack.node_version` [GH-32840]
* `azurerm_windows_web_app_slot` - add support for `~24` to `site_config.application_stack.node_version` [GH-32840]
* `azurerm_cognitive_account_rai_policy` - the `content_filter.severity_threshold` property is now optional [GH-32100]
* dependencies: `grpc` update to `1.82.1` [GH-32852]
* `azurerm_dashboard_grafana` - the `11` value for the `grafana_major_version` property has been deprecated and the property now supports `13` [GH-32777]
* `azurerm_container_registry` - the `trust_policy_enabled` property has been deprecated and removed from the provider [GH-32752]
* `sentinel` - migrate to `go-azure-sdk` [GH-32759]

FEATURES:
* New Action: `azurerm_web_app_set_slot_distribution` [GH-32364]

BUG FIXES:
