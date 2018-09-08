## 1.15.0 (Unreleased)

IMPROVEMENTS:

* dependencies: upgrading to v20.1.0 of `github.com/Azure/azure-sdk-for-go` [GH-1861]
* dependencies: upgrading to v10.15.3 of `github.com/Azure/go-autorest` [GH-1861]
* sdk: upgrading to version `2018-06-01` of the Compute API's [GH-1861]
* `azurerm_automation_runbook` - support for specifying the content field [GH-1696]
* `azurerm_virtual_network` - adding validation to `name` preventing empty values [GH-1898]

BUG FIXES:


* Data Source: `azurerm_azuread_service_principal` - passing a filter containing the name to Azure rather than querying locally [GH-1862]
* Data Source: `azurerm_azuread_service_principal` - passing a filter containing the name to Azure rather than querying locally [GH-1862]
* `azurerm_role_assignment` - parsing the Resource ID during deletion [GH-1887]
* `azurerm_role_definition` - parsing the Resource ID during deletion [GH-1887]

## 1.14.0 (September 06, 2018)

FEATURES:

* **New Data Source:** `azurerm_management_group` ([#1877](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1877))
* **New Resource:** `azurerm_management_group` ([#1788](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1788))
* **New Resource:** `azurerm_postgresql_virtual_network_rule` ([#1774](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1774))

IMPROVEMENTS:

* authentication: making the client registration consistent ([#1845](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1845))
* `azurerm_application_insights` - support for the `MobileCenter` kind ([#1878](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1878))
* `azurerm_function_app` - removing validation from the `version` field ([#1872](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1872))
* `azurerm_iothub` - exporting the `event_hub_events_endpoint`, `event_hub_events_path`, `event_hub_operations_endpoint` and `event_hub_operations_path` fields ([#1789](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1789))
* `azurerm_iothub` - support for `endpoint` and `route` blocks ([#1693](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1693))
* `azurerm_kubernetes_cluster` - making `linux_profile` optional ([#1821](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1821))
* `azurerm_storage_blob` - support for import ([#1816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1816))
* `azurerm_storage_container` - support for import ([#1816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1816))
* `azurerm_storage_queue` - support for import ([#1816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1816))
* `azurerm_storage_table` - support for import ([#1816](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1816))

BUG FIXES:

* `azurerm_data_lake_store_file` - updating the Resource ID to match the file path ([#1856](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1856))
* `azurerm_eventhub` - updating the validation to support periods, hyphens and underscores ([#1795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1795))
* `azurerm_eventhub_authorization_rule` - updating the validation error ([#1795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1795))
* `azurerm_eventhub_consumer_group` - updating the validation to support periods, hyphens and underscores ([#1795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1795))
* `azurerm_eventhub_namespace` - updating the validation error ([#1795](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1795))
* `azurerm_function_app` - support for names in upper-case ([#1835](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1835))
* `azurerm_kubernetes_cluster` - removing validation for the `pod_cidr` field when `network_plugin` is set to `azure` ([#1798](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1798))
* `azurerm_logic_app_workflow` - ensuring parameters are strings ([#1843](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1843))
* `azurerm_virtual_machine` - setting the `image_uri` property within the `storage_os_disk` block ([#1799](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1799))
* `azurerm_virtual_machine_data_disk_attachment` - obtaining a basic view, rather than the entire instance view of the Virtual Machine to work around an issue in the API ([#1855](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1855))

## 1.13.0 (August 15, 2018)

FEATURES:

* **New Data Source:** `azurerm_log_analytics_workspace` ([#1755](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1755))
* **New Resource:** `azurerm_monitor_action_group` ([#1725](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1725))

IMPROVEMENTS:

* dependencies: upgrading to `2018-04-01` of the IoTHub SDK ([#1717](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1717))
* Azure CLI Auth - using the `USERPROFILE` environment variable to locate the users home directory, if set ([#1718](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1718))
* Data Source `azurerm_kubernetes_cluster` - exposing the `max_pods` field within the `agent_pool_profile` block ([#1753](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1753))
* Data Source: `azurerm_kubernetes_cluster` - exposing the `add_on_profile` block ([#1751](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1751))
* `azurerm_automation_schedule` - adding the `week_days`, `month_days` and `monthly_occurrence` properties ([#1626](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1626))
* `azurerm_container_group` - adding a new `commands` field / deprecating the `command` field ([#1740](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1740))
* `azurerm_iothub` - support for the `Basic` SKU ([#1717](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1717))
* `azurerm_kubernetes_cluster` - support for `max_pods` within the `agent_pool_profile` block ([#1753](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1753))
* `azurerm_kubernetes_cluster` - support for the `add_on_profile` block ([#1751](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1751))
* `azurerm_kubernetes_cluster` - validation for when `pod_cidr` is set with a `network_plugin` set to `azure` ([#1763](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1763))
* `azurerm_kubernetes_cluster` - `client_id` and `client_secret` in the `service_principal` block are now ForceNew ([#1737](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1737))
* `azurerm_kubernetes_cluster` - `docker_bridge_cidr`, `dns_service_ip` and `service_cidr` are now conditionally set ([#1715](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1715))
* `azurerm_lb_nat_rule` - `protocol` property now supports `All` ([#1736](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1736))
* `azurerm_lb_nat_pool` - `protocol` property now supports `All` ([#1748](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1748))
* `azurerm_lb_probe` - `protocol` property now supports `Https` ([#1742](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1742))
* `azurerm_lb_rule` - support for the `All` protocol / adding validation ([#1754](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1754))

BUG FIXES:

* `azurerm_application_insights` - handling a `HTTP 201` being returned from the Create API which working around a breaking change in the API ([#1769](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1769))
* `azurerm_autoscale_setting` - filtering out the `$tags` tag ([#1770](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1770))
* `azurerm_eventhub` - allowing underscores in the name field ([#1768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1768))
* `azurerm_eventhub_authorization_rule` - allowing underscores in the name field ([#1768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1768))
* `azurerm_eventhub_consumer_group` - allowing underscores in the name field ([#1768](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1768))

## 1.12.0 (August 03, 2018)

UPGRADE NOTES:

* **Please Note:** When upgrading to v1.12.0 of the Azure Provider, you may need to specify the `priority` of any VM Scale Sets created between v1.6 of the Provider and v1.12. ([#1586](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1586))

FEATURES:

* **New Data Source:** `azurerm_container_registry` ([#1642](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1642))
* **New Resource:** `azurerm_service_fabric_cluster` ([#4](https://github.com/terraform-providers/terraform-provider-azurerm/issues/4))

IMPROVEMENTS:

* sdk: switching from `WaitForCompletion` -> `WaitForCompletionRef` when polling Future's ([#1660](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1660))
* Data Source: `azurerm_kubernetes_cluster` - support for specifying the `network_profile` block ([#1479](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1479))
* Data Source: `azurerm_kubernetes_cluster` - outputting the `node_resource_group` field ([#1649](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1649))
* `azurerm_kubernetes_cluster` - support for specifying the `network_profile` block ([#1479](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1479))
* `azurerm_kubernetes_cluster` - outputting the `node_resource_group` field ([#1649](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1649))
* `azurerm_role_assignment` - retrying resource creation to match the Azure CLI's behaviour ([#1647](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1647))
* `azurerm_virtual_machine` - setting the connection information for Provisioners ([#1646](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1646))


BUG FIXES:

* `azurerm_virtual_machine_scale_set` - removing the default of `priority`, since this isn't set on older instances. ([#1586](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1586))

## 1.11.0 (July 25, 2018)

FEATURES:

* **New Resource:** `azurerm_data_lake_store_file` ([#1261](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1261))

IMPROVEMENTS:

* `azurerm_app_service` - support for `min_tls_version` in the `site_config` block ([#1601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1601))
* `azurerm_app_service_slot` - support for `min_tls_version` in the `site_config` block ([#1601](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1601))
* `azurerm_data_lake_store` - support for enabling/disabling encryption ([#1623](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1623))
* `azurerm_data_lake_store` - support for managing the firewall state ([#1623](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1623))

BUG FIXES:

* `azurerm_servicebus_topic` - the `name` property now allows the ~ character ([#1640](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1640))

## 1.10.0 (July 21, 2018)

FEATURES:

* **New Data Source:** `azurerm_azuread_application` ([#1552](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1552))
* **New Data Source:** `azurerm_logic_app_workflow` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Data Source:** `azurerm_notification_hub` ([#1589](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1589))
* **New Data Source:** `azurerm_notification_hub_namespace` ([#1589](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1589))
* **New Data Source:** `azurerm_service_principal` ([#1564](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1564))
* **New Resource:** `azurerm_autoscale_setting` ([#1140](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1140))
* **New Resource:** `azurerm_data_lake_analytics_account` ([#1618](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1618))
* **New Resource:** `azurerm_data_lake_analytics_firewall_rule` ([#1618](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1618))
* **New Resource:** `azurerm_eventhub_namespace_authorization_rule` ([#1572](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1572))
* **New Resource:** `azurerm_logic_app_action_custom` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_logic_app_action_http` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_logic_app_trigger_custom` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_logic_app_trigger_http_request` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_logic_app_trigger_recurrence` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_logic_app_workflow` ([#1266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1266))
* **New Resource:** `azurerm_notification_hub` ([#1589](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1589))
* **New Resource:** `azurerm_notification_hub_authorization_rule` ([#1589](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1589))
* **New Resource:** `azurerm_notification_hub_namespace ` ([#1589](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1589))
* **New Resource:** `azurerm_servicebus_queue_authorization_rule` ([#1543](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1543))
* **New Resource:** `azurerm_service_principal` ([#1564](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1564))
* **New Resource:** `azurerm_service_principal_password` ([#1564](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1564))

IMPROVEMENTS:

* authentication: Refreshing the Service Principal Token before using it ([#1544](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1544))
* dependencies: updating to`2018-02-01` of the App Service SDK ([#1436](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1436))
* `azurerm_app_service` - support for setting `ftps_settings` in the `site_config` block ([#1577](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1577))
* `azurerm_app_service` - support for running containers ([#1578](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1578))
* `azurerm_app_service_slot` - support for Managed Service Identity ([#1579](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1579))
* `azurerm_app_service_slot` - Slots can now be updated in-place ([#1436](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1436))
* `azurerm_container_group` - support for images hosted in a private registry ([#1529](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1529))
* `azurerm_function_app` - adding support for the `site_credential` block ([#1567](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1567))
* `azurerm_function_app` - only setting `WEBSITE_CONTENTSHARE` and `WEBSITE_CONTENTAZUREFILECONNECTIONSTRING` for Consumption Apps ([#1515](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1515))
* `azurerm_mysql_server` - changing `tier` or `family` in `sku` property no longer destroys existing resource ([#1598](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1598))
* `azurerm_network_security_rule` - a maximum of 1 Application Security Group can be set per Security Rule  ([#1587](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1587))
* `azurerm_postgresql_server` - changing `tier` or `family` in `sku` property no longer destroys existing resource ([#1598](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1598))
* `azurerm_virtual_machine_scale_set` - `sku` property is now a list #1558 ([#1558](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1558))

BUG FIXES:

* `azurerm_application_insights` - fixing a bug where `application_type` was set to `other` ([#1563](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1563))
* `azurerm_lb` - allow `subnet_id` to be set to an empty value ([#1588](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1588))
* `azurerm_servicebus_subscription` - only sending `correlation_filter` values if they're set ([#1565](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1565))
* `azurerm_servicebus_subscription` - setting the `default_message_ttl` field ([#1568](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1568))
* `azurerm_snapshot` - allowing dashes in the `name` field ([#1574](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1574))
* `azurerm_traffic_manager_endpoint` - working around a bug in the API by setting `target` to nil when a `target_resource_id` is specified ([#1546](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1546))

## 1.9.0 (July 11, 2018)

FEATURES:

* **New Resource:** `azurerm_azuread_application` ([#1269](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1269))
* **New Resource:** `azurerm_data_lake_store_firewall_rule` ([#1499](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1499))
* **New Resource:** `azurerm_key_vault_access_policy` ([#1149](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1149))
* **New Resource:** `azurerm_scheduler_job` ([#1172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1172))
* **New Resource:** `azurerm_servicebus_namespace_authorization_rule` ([#1498](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1498))
* **New Resource:** `azurerm_user_assigned_identity` ([#1448](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1448))

IMPROVEMENTS:

* dependencies: updating the `containerservice` SDK to `2018-03-31` to support AKS GA ([#1474](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1474))
* dependencies: updating to `v18.0.0` of `Azure/azure-sdk-for-go` ([#1487](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1487))
* dependencies: updating to `v10.12.0` of `Azure/go-autorest` ([#1487](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1487))
* `azurerm_application_gateway` - adding `minimum_servers` to the probe resource ([#1510](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1510))
* `azurerm_cdn_profile` - support for `Standard_ChinaCdn` and `Standard_Microsoft` SKU's ([#1465](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1465))
* `azurerm_cosmosdb_account` - checking to see if the name is in use before creating ([#1464](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1464))
* `azurerm_cosmosdb_account` - fixing the validation on the `ip_range_filter` field ([#1463](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1463))
* `azurerm_dns_zone` - support for Private DNS Zones ([#1404](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1404))
* `azurerm_image` - change os_disk property to a list and add addtional property validation ([#1443](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1443))
* `azurerm_lb` - allow `private_ip_address` to be set to an empty value ([#1481](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1481))
* `azurerm_mysql_server` - changing the `storage_mb` property no longer forces a new resource ([#1532](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1532))
* `azurerm_postgresql_server` - changing the `storage_mb` property no longer forces a new resource ([#1532](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1532))
* `azurerm_servicebus_queue` - `enable_partitioning` can now be enabled for `Basic` and `Standard` tiers ([#1391](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1391))
* `azurerm_virtual_machine` - support for specifying user assigned identities ([#1448](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1448))
* `azurerm_virtual_machine` - making the `content` field in the `additional_unattend_config`  block (within `os_profile_windows_config`) sensitive ([#1471](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1471))
* `azurerm_virtual_machine_data_disk_attachment` - adding support for `write_accelerator_enabled` ([#147
