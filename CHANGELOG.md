## 1.15.0 (Unreleased)

IMPROVEMENTS:

* dependencies: upgrading to v20.1.0 of `github.com/Azure/azure-sdk-for-go` [GH-1861]
* dependencies: upgrading to v10.15.4 of `github.com/Azure/go-autorest` [GH-1861] [GH-1909]
* sdk: upgrading to version `2018-06-01` of the Compute API's [GH-1861]
* `azurerm_automation_runbook` - support for specifying the content field [GH-1696]
* `azurerm_app_service` - adding the `virtual_network_name` property [GH-1896]
* `azurerm_app_service_slot` - adding the `virtual_network_name` property [GH-1896]
* `azurerm_key_vault_certificate` - adding the `thumbprint` property [GH-1904]
* `azurerm_sql_database` - adding the `threat_detection_policy` property [GH-1628]
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
* `azurerm_virtual_machine_data_disk_attachment` - adding support for `write_accelerator_enabled` ([#1473](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1473))
* `azurerm_virtual_machine_scale_set` - ensuring we set the `vhd_containers` field to fix a crash ([#1411](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1411))
* `azurerm_virtual_machine_scale_set` - support for specifying user assigned identities ([#1448](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1448))
* `azurerm_virtual_machine_scale_set` - making the `content` field in the `additional_unattend_config`  block (within `os_profile_windows_config`) sensitive ([#1471](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1471))
* `azurerm_virtual_network_gateway` - adding support for the `radius_server_address`, `radius_server_secret` and `vpn_client_protocols` fields to the Data Source ([#1505](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1505))

BUG FIXES:

* `azurerm_key_vault_key` - handling the parent Key Vault being deleted ([#1535](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1535))
* `azurerm_sql_database` - fix `requested_service_objective_name` updates ([#1503](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1503))
* `azurerm_storage_account` - limiting the `tags` field to 128 characters to match the service ([#1524](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1524))
* `azurerm_virtual_network_gateway` - fix `azurerm_virtual_network_gateway` crashing when `vpn_client_configuration` was not supplied ([#1505](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1505))

## 1.8.0 (June 28, 2018)

FEATURES:

* **New Resource:** `azurerm_dns_caa_record` support ([#1450](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1450))
* **New Resource:** `azurerm_virtual_machine_data_disk_attachment` ([#1207](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1207))

IMPROVEMENTS:

* dependencies: upgrading to v10.11.4 of `Azure/go-autorest` ([#1418](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1418))
* dependencies: upgrading to v17.4.0 of `Azure/azure-sdk-for-go` ([#1418](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1418))
* `azurerm_lb` - addtional validation on properties ([#1403](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1403))
* `azurerm_application_gateway` - support for the `match` block for Probes ([#1446](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1446))
* `azurerm_log_analytics_solution` - support for Sovereign Clouds ([#1410](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1410))
* `azurerm_log_analytics_workspace` - support for Sovereign Clouds ([#1410](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1410))
* `azurerm_log_analytics_workspace` - support for the `PerGB2018` SKU ([#1079](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1079))
* `azurerm_mysql_server` -  `GeneralPurpose` and `MemoryOptimized` sku tiers now allow 4tb for the `storage_mb` property ([#1449](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1449))
* `azurerm_network_interface` - addtional validation on properties ([#1403](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1403))
* `azurerm_postgresql_server` -  `GeneralPurpose` and `MemoryOptimized` sku tiers now allow 4tb for the `storage_mb` property ([#1449](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1449))
* `azurerm_postgresql_server` - adding support for version 10.0 ([#1457](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1457))
* `azurerm_route_table` - adding the  disable BGP propagation property ([#1435](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1435))
* `azurerm_sql_database` - support for importing from a bacpac backup ([#972](https://github.com/terraform-providers/terraform-provider-azurerm/issues/972))
* `azurerm_virtual_machine` - support for setting the TimeZone on Windows ([#1265](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1265))

BUG FIXES:

* validation: ensuring IPv4/MAC addresses are detected correctly ([#1431](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1431))

## 1.7.0 (June 16, 2018)

UPGRADE NOTES:

~> **Please Note:** The field `overprovision` on the `azurerm_virtual_machine_scale_set` resource has changed from `false` to `true` to match the behaviour of Azure in this release. ([#1322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1322))

BUG FIXES:

* `azurerm_key_vault` - respecting the proxy environment varibles terraform does and now can create vaults when behind a proxy ([#1393](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1393))
* `azurerm_kubernetes_cluster` - `dns_prefix` is now required ([#1333](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1333))
* `azurerm_network_interface` - ensuring that Public IP's/Private IP Addresses can be removed once assigned ([#1295](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1295))
* `azurerm_public_ip` - setting the `domain_name_label` property into state ([#1287](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1287))
* `azurerm_storage_account` - file and blob encryption is now explicity `true` by default ([#1380](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1380))
* `azurerm_servicebus_namespace` - the `capacity` propety no longer unnecessarily forces a new resource when changed ([#1382](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1382))
* `azurerm_virtual_machine_scale_set` - the field `overprovision` is now `true` by default ([#1322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1322))
* `azurerm_app_service_plan` - the `name` property validation now allows understores ([#1351](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1351))

IMPROVEMENTS:

* `azurerm_automation_schedule` - adding the `interval` property and supporting recurring schedules ([#1384](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1384))
* `azurerm_dns_ns_record` - deprecated `record` properties in favor of a `records` list ([#991](https://github.com/terraform-providers/terraform-provider-azurerm/issues/991))
* `azurerm_function_app` - adding the `identity` property ([#1369](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1369))
* `azurerm_role_definition` - the `role_definition_id` property is now optional. The resource will now generate a random UUID if it is ommited ([#1378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1378))
* `azurerm_storage_account` - adding the `network_rules` property ([#1334](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1334))
* `azurerm_storage_account` - adding the `identity` property ([#1323](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1323))
* `azurerm_storage_blob` - adding the `content_type` property ([#1304](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1304))
* `azurerm_virtual_machine` - support for `write_accelerator_enabled` property on Premium disks attached to MS-series machines ([#964](https://github.com/terraform-providers/terraform-provider-azurerm/issues/964))
* `azurerm_virtual_machine_scale_set` - adding the `dns_settings` and `dns_servers` property ([#1209](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1209))
* `azurerm_virtual_machine_scale_set` - adding the `ip_forwarding` property ([#1209](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1209))
* `azurerm_virtual_network_gateway` - adding the properties `vpn_client_protocols`, `radius_server_address` and `radius_server_secret` ([#946](https://github.com/terraform-providers/terraform-provider-azurerm/issues/946))
* dependencies: migrating to the un-deprecated Preview's for Container Instance, EventGrid, Log Analytics and SQL ([#1322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1322))
* dependencies: upgrading to `2018-01-01` of the EventGrid API ([#1322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1322))
* dependencies: upgrading to `2018-03-01` of the Monitor API ([#1322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1322))

## 1.6.0 (May 24, 2018)

UPGRADE NOTES:

~> **Please Note:** The `azurerm_mysql_server` resource has been updated from the Preview API's to the GA API's - which requires code changes in your Terraform Configuration to use the new Pricing SKU's. Upon updating to v1.6.0 - you'll need to update the configuration from the Preview SKU's to the GA SKU's.

~> **Please Note:** The `azurerm_postgresql_server` resource has been updated from the Preview API's to the GA API's - which requires code changes in your Terraform Configuration to use the new Pricing SKU's. Upon updating to v1.6.0 - you'll need to update the configuration from the Preview SKU's to the GA SKU's.

* `azurerm_scheduler_job_collection` - the property `max_retry_interval` on both the resource and datasource has been deprecated in favour of `max_recurrence_interval` to better match Azure ([#1218](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1218))

FEATURES:

* **New Data Source:** `azurerm_storage_account_sas` ([#1011](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1011))
* **New Resource:** `azurerm_data_lake_store` ([#1219](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1219))
* **New Resource:** `azurerm_relay_namespace` ([#1233](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1233))

BUG FIXES:

* across data-sources and resources: making Connection Strings, Keys and Passwords sensitive fields ([#1242](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1242))
* `azurerm_virtual_machine_scale_set` - an empty `os_profile_windows_config` block no longer causes a panic ([#1224](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1224))

IMPROVEMENTS:

* authorization: upgrading to API version `2018-01-01-preview`
* `azurerm_app_service` - adding support for `ip_restriction`'s ([#1231](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1231))
* `azurerm_app_service_slot` - adding support for `ip_restriction`'s ([#1246](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1246))
* `azurerm_container_registry` - no longer forces a new resource on SKU change ([#1264](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1264))
* `azurerm_dns_zone` - datasource's `resource_group` field is now optional ([#1180](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1180))
* `azurerm_mysql_database` - ignoring casing for the `charset` field ([#1281](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1281))
* `azurerm_mysql_server` - support for the new GA Pricing SKU's ([#1154](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1154))
* `azurerm_postgresql_database` - ignoring the casing on the `collation` field ([#1255](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1255))
* `azurerm_postgresql_server` - support for the new GA Pricing SKU's ([#1190](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1190))
* `azurerm_public_ip` - computed values now default to an empy string ([#1247](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1247))
* `azurerm_role_assignment` - support for roles containing DataActions ([#1284](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1284))
* `azurerm_servicebus_queue` - adding `dead_lettering_on_message_expiration` ([#1235](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1235))
* `azurerm_virtual_machine_scale_set` - adding the `licence_type` property ([#1245](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1245))
* `azurerm_virtual_machine_scale_set` - adding the `priority` property ([#1250](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1250))

## 1.5.0 (May 14, 2018)

UPGRADE NOTES:

~> **Please Note:** Prior to v1.5 Data Sources in the AzureRM Provider returned `nil` rather than an error message when a Resource didn't exist, which was a bug. In order to bring this into line with other Providers - starting in v1.5 the AzureRM Provider will return an error message when a resource doesn't exist.

~> **Please Note:** This release fixes a bug in the `azurerm_redis_cache` resource where changes to fields weren't detected; as such you may see changes in the `redis_configuration` block, particularly with the `rdb_storage_connection_string` field. There's a bug tracking this inconsistency in [the Azure Rest API Specs Repository](https://github.com/Azure/azure-rest-api-specs/issues/3037).

FEATURES:

* **New Data Source:** `azurerm_cosmosdb_account` ([#1056](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1056))
* **New Data Source:** `azurerm_kubernetes_cluster` ([#1204](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1204))
* **New Data Source:** `azurerm_key_vault` ([#1202](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1202))
* **New Data Source:** `azurerm_key_vault_secret` ([#1202](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1202))
* **New Data Source:** `azurerm_route_table` ([#1203](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1203))

BUG FIXES:

* `azurerm_redis_cache` - changes to the `redis_configuration` block are now detected - please see the note above for more information ([#1211](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1211))

IMPROVEMENTS:

* dependencies - upgrading to v16.2.1 of `Azure/azure-sdk-for-go` ([#1198](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1198))
* dependencies - upgrading to v10.8.1 of `Azure/go-autorest` ([#1198](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1198))
* `azurerm_app_service` - support for HTTP2 ([#1188](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1188))
* `azurerm_app_service` - support for Managed Service Identity ([#1130](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1130))
* `azurerm_app_service_slot` - support for HTTP2 ([#1205](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1205))
* `azurerm_cosmosdb_account` - added support for the `connection_strings` property ([#1194](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1194))
* `azurerm_key_vault_certificate` - exposing the `certificate_data` ([#1200](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1200))
* `azurerm_kubernetes_cluster` - making `kube_config_raw` a sensitive field ([#1225](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1225))
* `azurerm_redis_cache` - Redis Caches can now be Imported ([#1211](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1211))
* `azurerm_redis_firewall_rule` - Redis Firewall Rules can now be Imported ([#1211](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1211))
* `azurerm_virtual_network` - guarding against nil-objects in the response ([#1208](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1208))
* `azurerm_virtual_network_gateway` - ignoring the case of the `GatewaySubnet` ([#1141](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1141))

## 1.4.0 (April 26, 2018)

UPGRADE NOTES:

* `azurerm_cosmosdb_account` - the field `failover_policy` has been deprecated in favour of `geo_locations` to better match Azure

FEATURES:

* **New Data Source:** `azurerm_recovery_services_vault` ([#995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/995))
* **New Resource:** `azurerm_recovery_services_vault` ([#995](https://github.com/terraform-providers/terraform-provider-azurerm/issues/995))
* **New Resource:** `azurerm_servicebus_subscription_rule` ([#1124](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1124))

IMPROVEMENTS:

* `azurerm_app_service` - support for updating in-place ([#1125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1125))
* `azurerm_app_service_plan` - support for `kind` being `app` ([#1156](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1156))
* `azurerm_cosmosdb_account` - support for `enable_automatic_failover` ([#1055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1055))
* `azurerm_cosmosdb_account` - support for the `ConsistentPrefix` consistncy level ([#1055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1055))
* `azurerm_cosmosdb_account` - `prefixes` can now be configured for locations ([#1055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1055))
* `azurerm_function_app` - support for updating in-place ([#1125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1125))
* `azurerm_key_vault` - adding cert permissions for `Purge` and `Recover` ([#1132](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1132))
* `azurerm_key_vault` - polling to ensure the Key Vault is resolvable via DNS ([#1081](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1081)] [[#1164](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1164))
* `azurerm_kubernetes_cluster` - only setting the Subnet ID when it's not an empty string ([#1158](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1158))
* `azurerm_kubernetes_cluster` - exposing the clusters credentials as `kube_config` ([#953](https://github.com/terraform-providers/terraform-provider-azurerm/issues/953))
* `azurerm_metric_alertrule` - filtering out tags prefixed with `$type` ([#1107](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1107))
* `azurerm_virtual_machine` - loading managed disk information from Azure when the machine is stopped ([#1100](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1100))
* `azurerm_virtual_machine` - make the `vm_size` property case insensitive ([#1131](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1131))

BUG FIXES:

* `azurerm_cosmosdb_account` - locations can now be modified in-place (without requiring multiple apply's) ([#1055](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1055))


## 1.3.3 (April 17, 2018)

FEATURES:

* **New Data Source:** `azurerm_app_service` ([#1071](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1071))
* **New Resource:** `azurerm_app_service_custom_hostname_binding` ([#1087](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1087))

IMPROVEMENTS:

* dependencies: upgrading to `v15.1.0` of `Azure/azure-sdk-for-go` ([#1099](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1099))
* dependencies: upgrading to `v10.6.0` of `Azure/go-autorest` ([#1077](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1077))
* `azurerm_app_service` - added support for the `https_only` field ([#1080](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1080))
* `azurerm_app_service_slot` - added support for the `https_only` field ([#1080](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1080))
* `azurerm_function_app` - added support for the `https_only` field ([#1080](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1080))
* `azurerm_key_vault_certificate` - exposing the certificate's associated `secret_id` ([#1096](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1096))
* `azurerm_redis_cache` - support for clusters on the internal network ([#1086](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1086))
* `azurerm_servicebus_queue` - support for setting `requires_session` ([#1111](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1111))
* `azurerm_sql_database` - changes to `collation` force a new resource ([#1066](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1066))

## 1.3.2 (April 04, 2018)

FEATURES:

* **New Resource:** `azurerm_packet_capture` ([#1044](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1044))
* **New Resource:** `azurerm_policy_assignment` ([#1051](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1051))

IMPROVEMENTS:

* `azurerm_virtual_machine_scale_set` - adds support for MSI ([#1018](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1018))

## 1.3.1 (March 29, 2018)

FEATURES:

* **New Data Source:** `azurerm_scheduler_job_collection` ([#990](https://github.com/terraform-providers/terraform-provider-azurerm/issues/990))
* **New Data Source:** `azurerm_traffic_manager_geographical_location` ([#987](https://github.com/terraform-providers/terraform-provider-azurerm/issues/987))
* **New Resource:** `azurerm_express_route_circuit_authorization` ([#992](https://github.com/terraform-providers/terraform-provider-azurerm/issues/992))
* **New Resource:** `azurerm_express_route_circuit_peering` ([#1033](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1033))
* **New Resource:** `azurerm_iothub` ([#887](https://github.com/terraform-providers/terraform-provider-azurerm/issues/887))
* **New Resource:** `azurerm_policy_definition` ([#1010](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1010))
* **New Resource:** `azurerm_sql_virtual_network_rule` ([#978](https://github.com/terraform-providers/terraform-provider-azurerm/issues/978))

IMPROVEMENTS:

* `azurerm_app_service` - allow changing `client_affinity_enabled` without requiring a resource recreation ([#993](https://github.com/terraform-providers/terraform-provider-azurerm/issues/993))
* `azurerm_app_service` - support for configuring `LocalSCM` source control ([#826](https://github.com/terraform-providers/terraform-provider-azurerm/issues/826))
* `azurerm_app_service` - returning a clearer error message when the name (which needs to be globally unique) is in use ([#1037](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1037))
* `azurerm_cosmosdb_account` - increasing the maximum value for `max_interval_in_seconds` from 100s to 86400s (1 day) ([#1000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1000))
* `azurerm_function_app` - returning a clearer error message when the name (which needs to be globally unique) is in use ([#1037](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1037))
* `azurerm_network_interface` - support for attaching to Application Gateways ([#1027](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1027))
* `azurerm_traffic_manager_endpoint` - adding support for `geo_mappings` ([#986](https://github.com/terraform-providers/terraform-provider-azurerm/issues/986))
* `azurerm_traffic_manager_profile` - adding support for the `traffic_routing_method` `Geographic` ([#986](https://github.com/terraform-providers/terraform-provider-azurerm/issues/986))
* `azurerm_virtual_machine_scale_sets` - support for attaching to Application Gateways ([#1027](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1027))
* `azurerm_virtual_network_gateway` - changes to `peering_address` now force a new resource ([#1040](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1040))

## 1.3.0 (March 15, 2018)

FEATURES:

* **New Data Source:** `azurerm_cdn_profile` ([#950](https://github.com/terraform-providers/terraform-provider-azurerm/issues/950))
* **New Data Source:** `azurerm_network_interface` ([#854](https://github.com/terraform-providers/terraform-provider-azurerm/issues/854))
* **New Data Source:** `azurerm_public_ips` ([#304](https://github.com/terraform-providers/terraform-provider-azurerm/issues/304))
* **New Data Source:** `azurerm_subscriptions` ([#940](https://github.com/terraform-providers/terraform-provider-azurerm/issues/940))
* **New Resource:** `azurerm_log_analytics_solution` ([#952](https://github.com/terraform-providers/terraform-provider-azurerm/issues/952))
* **New Resource:** `azurerm_sql_active_directory_administrator` ([#765](https://github.com/terraform-providers/terraform-provider-azurerm/issues/765))
* **New Resource:** `azurerm_scheduler_job_collection` ([#963](https://github.com/terraform-providers/terraform-provider-azurerm/issues/963))

BUG FIXES:

* `azurerm_application_gateway` - fixes a crash where `ssl_policy` isn't returned from the Azure API when importing existing resources ([#935](https://github.com/terraform-providers/terraform-provider-azurerm/issues/935))
* `azurerm_app_service` - supporting `client_affinity_enabled` being `false` ([#973](https://github.com/terraform-providers/terraform-provider-azurerm/issues/973))
* `azurerm_kubernetes_cluster` - exporting the FQDN ([#907](https://github.com/terraform-providers/terraform-provider-azurerm/issues/907))
* `azurerm_sql_elasticpool` - fixing a crash where `location` isn't returned for legacy resources ([#982](https://github.com/terraform-providers/terraform-provider-azurerm/issues/982))

IMPROVEMENTS:

* Data Source: `azurerm_builtin_role_definition` - loading available role definitions from Azure ([#770](https://github.com/terraform-providers/terraform-provider-azurerm/issues/770))
* Data Source: `azurerm_managed_disk` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))
* Data Source: `azurerm_network_security_group` - support for security rules including Application Security Groups ([#925](https://github.com/terraform-providers/terraform-provider-azurerm/issues/925))
* `azurerm_app_service_plan` -  support for provisioning Consumption Plans ([#981](https://github.com/terraform-providers/terraform-provider-azurerm/issues/981))
* `azurerm_cdn_endpoint` - adding support for GeoFilters, ProbePaths ([#967](https://github.com/terraform-providers/terraform-provider-azurerm/issues/967))
* `azurerm_cdn_endpoint` - making the `origin` block ForceNew to match Azure ([#967](https://github.com/terraform-providers/terraform-provider-azurerm/issues/967))
* `azurerm_function_app` - adding `client_affinity_enabled`, `use_32_bit_worker_process` and `websockets_enabled` ([#886](https://github.com/terraform-providers/terraform-provider-azurerm/issues/886))
* `azurerm_load_balancer` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))
* `azurerm_managed_disk` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))
* `azurerm_network_interface` - setting `internal_fqdn` if it's not nil ([#977](https://github.com/terraform-providers/terraform-provider-azurerm/issues/977))
* `azurerm_network_security_group` - support for security rules including Application Security Groups ([#925](https://github.com/terraform-providers/terraform-provider-azurerm/issues/925))
* `azurerm_network_security_rule` - support for security rules including Application Security Groups ([#925](https://github.com/terraform-providers/terraform-provider-azurerm/issues/925))
* `azurerm_public_ip` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))
* `azurerm_redis_cache` - add support for `notify-keyspace-events` ([#949](https://github.com/terraform-providers/terraform-provider-azurerm/issues/949))
* `azurerm_template_deployment` - support for specifying parameters via `parameters_body` ([#404](https://github.com/terraform-providers/terraform-provider-azurerm/issues/404))
* `azurerm_virtual_machine` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))
* `azurerm_virtual_machine_scale_set` - adding support for Availability Zones ([#811](https://github.com/terraform-providers/terraform-provider-azurerm/issues/811))

## 1.2.0 (March 02, 2018)

FEATURES:

* **New Data Source:** `azurerm_application_security_group` ([#914](https://github.com/terraform-providers/terraform-provider-azurerm/issues/914))
* **New Resource:** `azurerm_application_security_group` ([#905](https://github.com/terraform-providers/terraform-provider-azurerm/issues/905))
* **New Resource:** `azurerm_servicebus_topic_authorization_rule` ([#736](https://github.com/terraform-providers/terraform-provider-azurerm/issues/736))

BUG FIXES:

* `azurerm_kubernetes_cluster` - an empty `linux_profile.ssh_key.keydata` no longer causes a crash ([#903](https://github.com/terraform-providers/terraform-provider-azurerm/issues/903))
* `azurerm_kubernetes_cluster` - the `linux_profile.admin_username` and `linux_profile.ssh_key.keydata` fields now force a new resource ([#895](https://github.com/terraform-providers/terraform-provider-azurerm/issues/895))
* `azurerm_network_interface` - the `subnet_id` field is now case insensitive ([#866](https://github.com/terraform-providers/terraform-provider-azurerm/issues/866))
* `azurerm_network_security_group` - reverting `security_rules` to a set to fix an ordering issue ([#893](https://github.com/terraform-providers/terraform-provider-azurerm/issues/893))
* `azurerm_virtual_machine_scale_set` - the `computer_name_prefix` field now forces a new resource ([#871](https://github.com/terraform-providers/terraform-provider-azurerm/issues/871))

IMPROVEMENTS:

* authentication: adding support for Managed Service Identity ([#639](https://github.com/terraform-providers/terraform-provider-azurerm/issues/639))
* `azurerm_container_group` - added `dns_name_label` and `FQDN` properties ([#877](https://github.com/terraform-providers/terraform-provider-azurerm/issues/877))
* `azurerm_network_interface` - support for attaching to Application Security Groups ([#911](https://github.com/terraform-providers/terraform-provider-azurerm/issues/911))
* `azurerm_network_security_group` - support for augmented security rules ([#781](https://github.com/terraform-providers/terraform-provider-azurerm/issues/781))
* `azurerm_servicebus_subscription` - added support for the `forward_to` property ([#861](https://github.com/terraform-providers/terraform-provider-azurerm/issues/861))
* `azurerm_storage_account` - adding support for `account_kind` being `StorageV2` ([#851](https://github.com/terraform-providers/terraform-provider-azurerm/issues/851))
* `azurerm_virtual_network_gateway_connection` - support for IPsec/IKE Policies ([#834](https://github.com/terraform-providers/terraform-provider-azurerm/issues/834))

## 1.1.2 (February 19, 2018)

FEATURES:

* **New Resource:** `azurerm_kubernetes_cluster` ([#693](https://github.com/terraform-providers/terraform-provider-azurerm/issues/693))
* **New Resource:** `azurerm_app_service_active_slot` ([#818](https://github.com/terraform-providers/terraform-provider-azurerm/issues/818))
* **New Resource:** `azurerm_app_service_slot` ([#818](https://github.com/terraform-providers/terraform-provider-azurerm/issues/818))

BUG FIXES:

* **Data Source:** `azurerm_app_service_plan`: handling a 404 not being returned as an error ([#849](https://github.com/terraform-providers/terraform-provider-azurerm/issues/849))
* **Data Source:** `azurerm_virtual_network` - Fixing a crash when the DhcpOptions aren't specified ([#803](https://github.com/terraform-providers/terraform-provider-azurerm/issues/803))
* `azurerm_application_gateway` - fixing crashes due to schema mismatches for existing resources ([#848](https://github.com/terraform-providers/terraform-provider-azurerm/issues/848))
* `azurerm_storage_container` - add a retry for creation ([#846](https://github.com/terraform-providers/terraform-provider-azurerm/issues/846))

IMPROVEMENTS:

* authentication: pulling the `Environment` key from the Azure CLI Config ([#842](https://github.com/terraform-providers/terraform-provider-azurerm/issues/842))
* core: upgrading to `v12.5.0-beta` of the Azure SDK for Go ([#830](https://github.com/terraform-providers/terraform-provider-azurerm/issues/830))
* compute: upgrading to use the `2017-12-01` API Version ([#797](https://github.com/terraform-providers/terraform-provider-azurerm/issues/797))
* `azurerm_app_service_plan`: support for attaching to an App Service Environment ([#850](https://github.com/terraform-providers/terraform-provider-azurerm/issues/850))
* `azurerm_container_group` - adding `restart_policy` ([#827](https://github.com/terraform-providers/terraform-provider-azurerm/issues/827))
* `azurerm_managed_disk` - updated the validation on `disk_size_gb` / made it computed ([#800](https://github.com/terraform-providers/terraform-provider-azurerm/issues/800))
* `azurerm_role_assignment` - add `role_definition_name` ([#775](https://github.com/terraform-providers/terraform-provider-azurerm/issues/775))
* `azurerm_subnet` - add support for Service Endpoints ([#786](https://github.com/terraform-providers/terraform-provider-azurerm/issues/786))
* `azurerm_virtual_machine` - changing `managed_disk_id` and `create_option` to be not ForceNew ([#813](https://github.com/terraform-providers/terraform-provider-azurerm/issues/813))


## 1.1.1 (February 06, 2018)

BUG FIXES:

* `azurerm_public_ip` - Setting the `ip_address` field regardless of the DNS Settings ([#772](https://github.com/terraform-providers/terraform-provider-azurerm/issues/772))
* `azurerm_virtual_machine` - ignores the case of the Managed Data Disk ID's to work around an Azure Portal bug ([#792](https://github.com/terraform-providers/terraform-provider-azurerm/issues/792))

FEATURES:

* **New Data Source:** `azurerm_storage_account` ([#794](https://github.com/terraform-providers/terraform-provider-azurerm/issues/794))
* **New Data Source:** `azurerm_virtual_network_gateway` ([#796](https://github.com/terraform-providers/terraform-provider-azurerm/issues/796))

## 1.1.0 (January 26, 2018)

UPGRADE NOTES:

* Data Source: `azurerm_builtin_role_definition` - now returns the correct UUID/GUID for the `Virtual Machines Contributor` role (previously the ID for the `Classic Virtual Machine Contributor` role was returned) ([#762](https://github.com/terraform-providers/terraform-provider-azurerm/issues/762))
* `azurerm_snapshot` - `source_uri` now forces a new resource on changes due to behavioural changes in the Azure API ([#744](https://github.com/terraform-providers/terraform-provider-azurerm/issues/744))

FEATURES:

* **New Data Source:** `azurerm_dns_zone` ([#702](https://github.com/terraform-providers/terraform-provider-azurerm/issues/702))
* **New Resource:** `azurerm_metric_alertrule` ([#478](https://github.com/terraform-providers/terraform-provider-azurerm/issues/478))
* **New Resource:** `azurerm_virtual_network_gateway` ([#133](https://github.com/terraform-providers/terraform-provider-azurerm/issues/133))
* **New Resource:** `azurerm_virtual_network_gateway_connection` ([#133](https://github.com/terraform-providers/terraform-provider-azurerm/issues/133))

IMPROVEMENTS:

* core: upgrading to `v12.2.0-beta` of `Azure/azure-sdk-for-go` ([#684](https://github.com/terraform-providers/terraform-provider-azurerm/issues/684))
* core: upgrading to `v9.7.0` of `Azure/go-autorest` ([#684](https://github.com/terraform-providers/terraform-provider-azurerm/issues/684))
* Data Source: `azurerm_builtin_role_definition` - adding extra role definitions ([#762](https://github.com/terraform-providers/terraform-provider-azurerm/issues/762))
* `azurerm_app_service` - exposing the `outbound_ip_addresses` field ([#700](https://github.com/terraform-providers/terraform-provider-azurerm/issues/700))
* `azurerm_function_app` - exposing the `outbound_ip_addresses` field ([#706](https://github.com/terraform-providers/terraform-provider-azurerm/issues/706))
* `azurerm_function_app` - add support for the `always_on` and `connection_string` fields ([#695](https://github.com/terraform-providers/terraform-provider-azurerm/issues/695))
* `azurerm_image` - add support for filtering images by a regex on the name ([#642](https://github.com/terraform-providers/terraform-provider-azurerm/issues/642))
* `azurerm_lb` - adding support for the `Standard` SKU (in Preview) ([#665](https://github.com/terraform-providers/terraform-provider-azurerm/issues/665))
* `azurerm_public_ip` - adding support for the `Standard` SKU (in Preview) ([#665](https://github.com/terraform-providers/terraform-provider-azurerm/issues/665))
* `azurerm_network_security_rule` - add support for augmented security rules ([#692](https://github.com/terraform-providers/terraform-provider-azurerm/issues/692))
* `azurerm_role_assignment` - generating a name if one isn't specified ([#685](https://github.com/terraform-providers/terraform-provider-azurerm/issues/685))
* `azurerm_traffic_manager_profile` - adding support for setting `protocol` to `TCP` ([#742](https://github.com/terraform-providers/terraform-provider-azurerm/issues/742))

## 1.0.1 (January 12, 2018)

FEATURES:

* **New Data Source:** `azurerm_app_service_plan` ([#668](https://github.com/terraform-providers/terraform-provider-azurerm/issues/668))
* **New Data Source:** `azurerm_eventhub_namespace` ([#673](https://github.com/terraform-providers/terraform-provider-azurerm/issues/673))
* **New Resource:** `azurerm_function_app` ([#647](https://github.com/terraform-providers/terraform-provider-azurerm/issues/647))

IMPROVEMENTS:

* core: adding a cache to the Storage Account Keys ([#634](https://github.com/terraform-providers/terraform-provider-azurerm/issues/634))
* `azurerm_eventhub` - added support for `capture_description` ([#681](https://github.com/terraform-providers/terraform-provider-azurerm/issues/681))
* `azurerm_eventhub_consumer_group` - adding validation for the user metadata field ([#641](https://github.com/terraform-providers/terraform-provider-azurerm/issues/641))
* `azurerm_lb` - adding the computed field `public_ip_addresses` ([#633](https://github.com/terraform-providers/terraform-provider-azurerm/issues/633))
* `azurerm_local_network_gateway` - add support for `tags` ([#638](https://github.com/terraform-providers/terraform-provider-azurerm/issues/638))
* `azurerm_network_interface` - support for Accelerated Networking ([#672](https://github.com/terraform-providers/terraform-provider-azurerm/issues/672))
* `azurerm_storage_account` - expose `primary_connection_string` and `secondary_connection_string` ([#647](https://github.com/terraform-providers/terraform-provider-azurerm/issues/647))


## 1.0.0 (December 15, 2017)

FEATURES:

* **New Data Source:** `azurerm_network_security_group` ([#623](https://github.com/terraform-providers/terraform-provider-azurerm/issues/623))
* **New Data Source:** `azurerm_virtual_network` ([#533](https://github.com/terraform-providers/terraform-provider-azurerm/issues/533))
* **New Resource:** `azurerm_management_lock` ([#575](https://github.com/terraform-providers/terraform-provider-azurerm/issues/575))
* **New Resource:** `azurerm_network_watcher` ([#571](https://github.com/terraform-providers/terraform-provider-azurerm/issues/571))

IMPROVEMENTS:

* authentication - add support for the latest Azure CLI configuration ([#573](https://github.com/terraform-providers/terraform-provider-azurerm/issues/573))
* authentication - conditional loading of the Subscription ID / Tenant ID / Environment ([#574](https://github.com/terraform-providers/terraform-provider-azurerm/issues/574))
* core - appending additions to the User Agent, so we don't overwrite the Go SDK User Agent info ([#587](https://github.com/terraform-providers/terraform-provider-azurerm/issues/587))
* core - Upgrading `Azure/azure-sdk-for-go` to v11.2.2-beta ([#594](https://github.com/terraform-providers/terraform-provider-azurerm/issues/594))
* core - upgrading `Azure/go-autorest` to v9.5.2 ([#617](https://github.com/terraform-providers/terraform-provider-azurerm/issues/617))
* core - skipping Resource Provider Registration in AutoRest when opted-out ([#630](https://github.com/terraform-providers/terraform-provider-azurerm/issues/630))
* `azurerm_app_service` - exposing the Default Hostname as a Computed field

## 0.3.3 (November 14, 2017)

FEATURES:

* **New Resource:** `azurerm_redis_firewall_rule` ([#529](https://github.com/terraform-providers/terraform-provider-azurerm/issues/529))

IMPROVEMENTS:

* authentication: allow using multiple subscriptions for Azure CLI auth ([#445](https://github.com/terraform-providers/terraform-provider-azurerm/issues/445))
* core: appending the CloudShell version to the user agent when running within CloudShell ([#483](https://github.com/terraform-providers/terraform-provider-azurerm/issues/483))
* `azurerm_app_service` / `azurerm_app_service_plan` - adding validation for the `name` fields ([#528](https://github.com/terraform-providers/terraform-provider-azurerm/issues/528))
* `azurerm_container_registry` - Migration: Fixing a crash when the storage_account block is nil ([#551](https://github.com/terraform-providers/terraform-provider-azurerm/issues/551))
* `azurerm_lb_nat_rule`: support for floating IP's ([#542](https://github.com/terraform-providers/terraform-provider-azurerm/issues/542))
* `azurerm_public_ip` - Clarify the error message for the validation of domain name label ([#485](https://github.com/terraform-providers/terraform-provider-azurerm/issues/485))
* `azurerm_network_security_group` - fixing a crash when changes were made outside of Terraform ([#492](https://github.com/terraform-providers/terraform-provider-azurerm/issues/492))
* `azurerm_redis_cache`: support for Patch Schedules ([#540](https://github.com/terraform-providers/terraform-provider-azurerm/issues/540))
* `azurerm_virtual_machine` - ensuring `vhd_uri` is validated ([#470](https://github.com/terraform-providers/terraform-provider-azurerm/issues/470))
* `azurerm_virtual_machine_scale_set`: fixing a crash where accelerated networking isn't returned by the API ([#480](https://github.com/terraform-providers/terraform-provider-azurerm/issues/480))

## 0.3.2 (October 30, 2017)

FEATURES: 

* **New Resource:** `azurerm_application_gateway` ([#413](https://github.com/terraform-providers/terraform-provider-azurerm/issues/413))

IMPROVEMENTS: 

  - `azurerm_virtual_machine_scale_set` - Add nil check to os disk ([#436](https://github.com/terraform-providers/terraform-provider-azurerm/issues/436))

  - `azurerm_key_vault` - Increased timeout on dns availability ([#457](https://github.com/terraform-providers/terraform-provider-azurerm/issues/457))
  
  - `azurerm_route_table` - Fix issue when routes are computed ([#450](https://github.com/terraform-providers/terraform-provider-azurerm/issues/450))

## 0.3.1 (October 21, 2017)

IMPROVEMENTS:

  - `azurerm_virtual_machine_scale_set` - Updating this resource with the v11 of the Azure SDK for Go ([#448](https://github.com/terraform-providers/terraform-provider-azurerm/issues/448))

## 0.3.0 (October 17, 2017)

UPGRADE NOTES:

  - `azurerm_automation_account` - the SKU `Free` has been replaced with `Basic`.
  - `azurerm_container_registry` - Azure has updated the SKU from `Basic` to `Classic`, with new `Basic`, `Standard` and `Premium` SKU's introduced.
  - `azurerm_container_registry` - the `storage_account` block is now `storage_account_id` and is only required for `Classic` SKU's
  - `azurerm_key_vault` - `certificate_permissions`, `key_permissions` and `secret_permissions` have all had the `All` option removed by Azure. Each permission now needs to be specified manually.
  * `azurerm_route_table` - `route` is no longer computed
  - `azurerm_servicebus_namespace` - The `capacity` field can only be set for `Premium` SKU's
  - `azurerm_servicebus_queue` - The `enable_batched_operations` and `support_ordering` fields have been deprecated by Azure.
  - `azurerm_servicebus_subscription` - The `dead_lettering_on_filter_evaluation_exceptions` has been removed by Azure.
  - `azurerm_servicebus_topic` - The `enable_filtering_messages_before_publishing` field has been removed by Azure.

FEATURES:

* **New Data Source:** `azurerm_builtin_role_definition` ([#384](https://github.com/terraform-providers/terraform-provider-azurerm/issues/384))
* **New Data Source:** `azurerm_image` ([#382](https://github.com/terraform-providers/terraform-provider-azurerm/issues/382))
* **New Data Source:** `azurerm_key_vault_access_policy` ([#423](https://github.com/terraform-providers/terraform-provider-azurerm/issues/423))
* **New Data Source:** `azurerm_platform_image` ([#375](https://github.com/terraform-providers/terraform-provider-azurerm/issues/375))
* **New Data Source:** `azurerm_role_definition` ([#414](https://github.com/terraform-providers/terraform-provider-azurerm/issues/414))
* **New Data Source:** `azurerm_snapshot` ([#420](https://github.com/terraform-providers/terraform-provider-azurerm/issues/420))
* **New Data Source:** `azurerm_subnet` ([#411](https://github.com/terraform-providers/terraform-provider-azurerm/issues/411))
* **New Resource:** `azurerm_key_vault_certificate` ([#408](https://github.com/terraform-providers/terraform-provider-azurerm/issues/408))
* **New Resource:** `azurerm_role_assignment` ([#414](https://github.com/terraform-providers/terraform-provider-azurerm/issues/414))
* **New Resource:** `azurerm_role_definition` ([#414](https://github.com/terraform-providers/terraform-provider-azurerm/issues/414))
* **New Resource:** `azurerm_snapshot` ([#420](https://github.com/terraform-providers/terraform-provider-azurerm/issues/420))

IMPROVEMENTS:

* Upgrading to v11 of the Azure SDK for Go ([#367](https://github.com/terraform-providers/terraform-provider-azurerm/issues/367))
* `azurerm_client_config` - updating the data source to work when using AzureCLI auth ([#393](https://github.com/terraform-providers/terraform-provider-azurerm/issues/393))
* `azurerm_container_group` - add support for volume mounts ([#366](https://github.com/terraform-providers/terraform-provider-azurerm/issues/366))
* `azurerm_key_vault` - fix a crash when no certificate_permissions are defined ([#374](https://github.com/terraform-providers/terraform-provider-azurerm/issues/374))
* `azurerm_key_vault` - waiting for the DNS to propagate ([#401](https://github.com/terraform-providers/terraform-provider-azurerm/issues/401))
* `azurerm_managed_disk` - support for creating Managed Disks from Platform Images by supporting "FromImage" ([#399](https://github.com/terraform-providers/terraform-provider-azurerm/issues/399))
* `azurerm_managed_disk` - support for creating Encrypted Managed Disks ([#399](https://github.com/terraform-providers/terraform-provider-azurerm/issues/399))
* `azurerm_mysql_*` - Ensuring we register the MySQL Resource Provider ([#397](https://github.com/terraform-providers/terraform-provider-azurerm/issues/397))
* `azurerm_network_interface` - exposing all of the Private IP Addresses assigned to the NIC ([#409](https://github.com/terraform-providers/terraform-provider-azurerm/issues/409))
* `azurerm_network_security_group` / `azurerm_network_security_rule` - refactoring ([#405](https://github.com/terraform-providers/terraform-provider-azurerm/issues/405))
* `azurerm_route_table` - removing routes when none are specified ([#403](https://github.com/terraform-providers/terraform-provider-azurerm/issues/403))
* `azurerm_route_table` - refactoring `route` from a Set to a List ([#402](https://github.com/terraform-providers/terraform-provider-azurerm/issues/402))
* `azurerm_route` - refactoring `route` from a Set to a List ([#402](https://github.com/terraform-providers/terraform-provider-azurerm/issues/402))
* `azurerm_storage_account` - support for File Encryption ([#363](https://github.com/terraform-providers/terraform-provider-azurerm/issues/363))
* `azurerm_storage_account` - support for Custom Domain ([#363](https://github.com/terraform-providers/terraform-provider-azurerm/issues/363))
* `azurerm_storage_account` - splitting the storage account Tier and Replication out into separate fields ([#363](https://github.com/terraform-providers/terraform-provider-azurerm/issues/363))
- `azurerm_storage_account` - returning a user friendly error when trying to provision a Blob Storage Account with ZRS redundancy ([#421](https://github.com/terraform-providers/terraform-provider-azurerm/issues/421))
* `azurerm_subnet` - making it possible to remove Network Security Groups / Route Tables ([#411](https://github.com/terraform-providers/terraform-provider-azurerm/issues/411))
* `azurerm_virtual_machine` - fixing a bug where `additional_unattend_config.content` was being updated unintentionally ([#377](https://github.com/terraform-providers/terraform-provider-azurerm/issues/377))
* `azurerm_virtual_machine` - switching to use Lists instead of Sets ([#426](https://github.com/terraform-providers/terraform-provider-azurerm/issues/426))
* `azurerm_virtual_machine_scale_set` - fixing a bug where `additional_unattend_config.content` was being updated unintentionally ([#377](https://github.com/terraform-providers/terraform-provider-azurerm/issues/377))
* `azurerm_virtual_machine_scale_set` - support for multiple network profiles ([#378](https://github.com/terraform-providers/terraform-provider-azurerm/issues/378))

## 0.2.2 (September 28, 2017)

FEATURES:

* **New Resource:** `azurerm_key_vault_key` ([#356](https://github.com/terraform-providers/terraform-provider-azurerm/issues/356))
* **New Resource:** `azurerm_log_analytics_workspace` ([#331](https://github.com/terraform-providers/terraform-provider-azurerm/issues/331))
* **New Resource:** `azurerm_mysql_configuration` ([#352](https://github.com/terraform-providers/terraform-provider-azurerm/issues/352))
* **New Resource:** `azurerm_mysql_database` ([#352](https://github.com/terraform-providers/terraform-provider-azurerm/issues/352))
* **New Resource:** `azurerm_mysql_firewall_rule` ([#352](https://github.com/terraform-providers/terraform-provider-azurerm/issues/352))
* **New Resource:** `azurerm_mysql_server` ([#352](https://github.com/terraform-providers/terraform-provider-azurerm/issues/352))

IMPROVEMENTS:

* Updating the provider initialization & adding a `skip_credentials_validation` field to the provider for some advanced scenarios ([#322](https://github.com/terraform-providers/terraform-provider-azurerm/issues/322))

## 0.2.1 (September 25, 2017)

FEATURES:

* **New Resource:** `azurerm_automation_account` ([#257](https://github.com/terraform-providers/terraform-provider-azurerm/issues/257))
* **New Resource:** `azurerm_automation_credential` ([#257](https://github.com/terraform-providers/terraform-provider-azurerm/issues/257))
* **New Resource:** `azurerm_automation_runbook` ([#257](https://github.com/terraform-providers/terraform-provider-azurerm/issues/257))
* **New Resource:** `azurerm_automation_schedule` ([#257](https://github.com/terraform-providers/terraform-provider-azurerm/issues/257))
* **New Resource:** `azurerm_app_service` ([#344](https://github.com/terraform-providers/terraform-provider-azurerm/issues/344))

IMPROVEMENTS:

* `azurerm_client_config` - adding `service_principal_application_id` ([#348](https://github.com/terraform-providers/terraform-provider-azurerm/issues/348))
* `azurerm_key_vault` - adding `application_id` and `certificate_permissions` ([#348](https://github.com/terraform-providers/terraform-provider-azurerm/issues/348))

BUG FIXES:

* `azurerm_virtual_machine_scale_set` - fix panic with `additional_unattend_config` block ([#266](https://github.com/terraform-providers/terraform-provider-azurerm/issues/266))

## 0.2.0 (September 15, 2017)

FEATURES:

* **Support for authenticating using the Azure CLI** ([#316](https://github.com/terraform-providers/terraform-provider-azurerm/issues/316))
* **New Resource:** `azurerm_container_group` ([#333](https://github.com/terraform-providers/terraform-provider-azurerm/issues/333)] [[#311](https://github.com/terraform-providers/terraform-provider-azurerm/issues/311)] [[#338](https://github.com/terraform-providers/terraform-provider-azurerm/issues/338))

IMPROVEMENTS:

* `azurerm_app_service_plan` - support for Linux App Service Plans ([#332](https://github.com/terraform-providers/terraform-provider-azurerm/issues/332))
* `azurerm_postgresql_server` - supporting additional storage sizes ([#239](https://github.com/terraform-providers/terraform-provider-azurerm/issues/239))
* `azurerm_public_ip` - verifying the ID is valid before importing ([#320](https://github.com/terraform-providers/terraform-provider-azurerm/issues/320))
* `azurerm_sql_server` - verifying the name is valid before creating ([#323](https://github.com/terraform-providers/terraform-provider-azurerm/issues/323))
* `resource_group_name` - validation has been added to all resources that use this attribute ([#330](https://github.com/terraform-providers/terraform-provider-azurerm/issues/330))

## 0.1.7 (September 11, 2017)

FEATURES:

* **New Resource:** `azurerm_postgresql_configuration` ([#210](https://github.com/terraform-providers/terraform-provider-azurerm/issues/210))
* **New Resource:** `azurerm_postgresql_database` ([#210](https://github.com/terraform-providers/terraform-provider-azurerm/issues/210))
* **New Resource:** `azurerm_postgresql_firewall_rule` ([#210](https://github.com/terraform-providers/terraform-provider-azurerm/issues/210))
* **New Resource:** `azurerm_postgresql_server` ([#210](https://github.com/terraform-providers/terraform-provider-azurerm/issues/210))

IMPROVEMENTS:

* `azurerm_cdn_endpoint` - defaulting the `http_port` and `https_port` ([#301](https://github.com/terraform-providers/terraform-provider-azurerm/issues/301))
* `azurerm_cosmos_db_account`: allow setting the Kind to MongoDB/GlobalDocumentDB ([#299](https://github.com/terraform-providers/terraform-provider-azurerm/issues/299))

## 0.1.6 (August 31, 2017)

FEATURES:

* **New Data Source**: `azurerm_subscription` ([#285](https://github.com/terraform-providers/terraform-provider-azurerm/issues/285))
* **New Resource:** `azurerm_app_service_plan` ([#1](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1))
* **New Resource:** `azurerm_eventgrid_topic` ([#260](https://github.com/terraform-providers/terraform-provider-azurerm/issues/260))
* **New Resource:** `azurerm_key_vault_secret` ([#269](https://github.com/terraform-providers/terraform-provider-azurerm/issues/269))

IMPROVEMENTS:

* `azurerm_image` - added a default to the `caching` field ([#259](https://github.com/terraform-providers/terraform-provider-azurerm/issues/259))
* `azurerm_key_vault` - validation for the `name` field ([#270](https://github.com/terraform-providers/terraform-provider-azurerm/issues/270))
* `azurerm_network_interface` - support for multiple IP Configurations / setting the Primary IP Configuration ([#245](https://github.com/terraform-providers/terraform-provider-azurerm/issues/245))
* `azurerm_resource_group` - poll until the resource group is created (by migrating to the Azure SDK for Go) ([#289](https://github.com/terraform-providers/terraform-provider-azurerm/issues/289))
* `azurerm_search_service` - migrating to use the Azure SDK for Go ([#283](https://github.com/terraform-providers/terraform-provider-azurerm/issues/283))
* `azurerm_sql_*` - ensuring deleted resources are detected ([#289](https://github.com/terraform-providers/terraform-provider-azurerm/issues/289)] / [[#255](https://github.com/terraform-providers/terraform-provider-azurerm/issues/255))
* `azurerm_sql_database` - Import Support ([#289](https://github.com/terraform-providers/terraform-provider-azurerm/issues/289))
* `azurerm_sql_database` - migrating to using the Azure SDK for Go ([#289](https://github.com/terraform-providers/terraform-provider-azurerm/issues/289))
* `azurerm_sql_firewall_rule` - migrating to using the Azure SDK for Go ([#289](https://github.com/terraform-providers/terraform-provider-azurerm/issues/289))
* `azurerm_sql_server` - added checks to handle `name` not being globally unique ([#189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/189))
* `azurerm_sql_server` - making `administrator_login` `ForceNew` ([#189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/189))
* `azurerm_sql_server` - migrate to using the azure-sdk-for-go ([#189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/189))
* `azurerm_virtual_machine` - Force recreation if `storage_data_disk`.`create_option` changes ([#240](https://github.com/terraform-providers/terraform-provider-azurerm/issues/240))
* `azurerm_virtual_machine_scale_set` - Fix address issue when setting the `winrm` block ([#271](https://github.com/terraform-providers/terraform-provider-azurerm/issues/271))
* updating to `v10.3.0-beta` of the Azure SDK for Go ([#258](https://github.com/terraform-providers/terraform-provider-azurerm/issues/258))
* Removing the (now unused) Riviera SDK ([#289](https://github.com/terraform-providers/terraform-provider-azurerm/issues/289)] [[#291](https://github.com/terraform-providers/terraform-provider-azurerm/issues/291))

BUG FIXES:

* `azurerm_cosmosdb_account` - fixing the validation on the name field ([#263](https://github.com/terraform-providers/terraform-provider-azurerm/issues/263))
* `azurerm_sql_server` - handle deleted servers correctly ([#189](https://github.com/terraform-providers/terraform-provider-azurerm/issues/189))
* Fixing the `Microsoft.Insights` Resource Provider Registration ([#282](https://github.com/terraform-providers/terraform-provider-azurerm/issues/282))

## 0.1.5 (August 09, 2017)

IMPROVEMENTS:

* `azurerm_sql_*` - upgrading to version `2014-04-01` of the SQL API's ([#201](https://github.com/terraform-providers/terraform-provider-azurerm/issues/201))
* `azurerm_virtual_machine` - support for the `Windows_Client` Hybrid Use Benefit type ([#212](https://github.com/terraform-providers/terraform-provider-azurerm/issues/212))
* `azurerm_virtual_machine_scale_set` - support for custom images and managed disks ([#203](https://github.com/terraform-providers/terraform-provider-azurerm/issues/203))

BUG FIXES:

* `azurerm_sql_database` - fixing creating a DB with a PointInTimeRestore ([#197](https://github.com/terraform-providers/terraform-provider-azurerm/issues/197))
* `azurerm_virtual_machine` - fix a crash when the properties for a network inteface aren't returned ([#208](https://github.com/terraform-providers/terraform-provider-azurerm/issues/208))
* `azurerm_virtual_machine` - changes to custom data should force new resource ([#211](https://github.com/terraform-providers/terraform-provider-azurerm/issues/211))
* `azurerm_virtual_machine` - fixes a crash caused by an empty `os_profile_windows_config` block ([#222](https://github.com/terraform-providers/terraform-provider-azurerm/issues/222))
* Checking to ensure the HTTP Response isn't `nil` before accessing it (fixes ([#200](https://github.com/terraform-providers/terraform-provider-azurerm/issues/200)]) [[#204](https://github.com/terraform-providers/terraform-provider-azurerm/issues/204))

## 0.1.4 (July 26, 2017)

BUG FIXES:

* `azurerm_dns_*` - upgrading to version `2016-04-01` of the Azure DNS API by switching from Riviera -> Azure SDK for Go ([#192](https://github.com/terraform-providers/terraform-provider-azurerm/issues/192))

## 0.1.3 (July 21, 2017)

FEATURES:

* **New Resource:** `azurerm_dns_ptr_record` ([#141](https://github.com/terraform-providers/terraform-provider-azurerm/issues/141))
* **New Resource:**`azurerm_image` ([#8](https://github.com/terraform-providers/terraform-provider-azurerm/issues/8))
* **New Resource:** `azurerm_servicebus_queue` ([#151](https://github.com/terraform-providers/terraform-provider-azurerm/issues/151))

IMPROVEMENTS:

* `azurerm_client_config` - added a `service_principal_object_id` attribute to the data source ([#175](https://github.com/terraform-providers/terraform-provider-azurerm/issues/175))
* `azurerm_search_service` - added import support ([#172](https://github.com/terraform-providers/terraform-provider-azurerm/issues/172))
* `azurerm_servicebus_topic` - added a `status` field to allow disabling the topic ([#150](https://github.com/terraform-providers/terraform-provider-azurerm/issues/150))
* `azurerm_storage_account` - Added support for Require secure transfer ([#167](https://github.com/terraform-providers/terraform-provider-azurerm/issues/167))
* `azurerm_storage_table` - updating the name validation ([#143](https://github.com/terraform-providers/terraform-provider-azurerm/issues/143))
* `azurerm_virtual_machine` - making `admin_password` optional for Linux VM's ([#154](https://github.com/terraform-providers/terraform-provider-azurerm/issues/154))
* `azurerm_virtual_machine_scale_set` - adding a `plan` block for Marketplace images ([#161](https://github.com/terraform-providers/terraform-provider-azurerm/issues/161))

## 0.1.2 (June 29, 2017)

FEATURES:

* **New Data Source:** `azurerm_managed_disk` ([#121](https://github.com/terraform-providers/terraform-provider-azurerm/issues/121))
* **New Resource:** `azurerm_application_insights` ([#3](https://github.com/terraform-providers/terraform-provider-azurerm/issues/3))
* **New Resource:** `azurerm_cosmosdb_account` ([#108](https://github.com/terraform-providers/terraform-provider-azurerm/issues/108))
* `azurerm_network_interface` now supports import ([#119](https://github.com/terraform-providers/terraform-provider-azurerm/issues/119))

IMPROVEMENTS:

* Ensuring consistency in when storing the `location` field in the state for the `azurerm_availability_set`, `azurerm_express_route_circuit`, `azurerm_load_balancer`, `azurerm_local_network_gateway`, `azurerm_managed_disk`, `azurerm_network_security_group`
`azurerm_public_ip`, `azurerm_resource_group`, `azurerm_route_table`, `azurerm_storage_account`, `azurerm_virtual_machine` and `azurerm_virtual_network` resources ([#123](https://github.com/terraform-providers/terraform-provider-azurerm/issues/123))
* `azurerm_redis_cache` - now supports backup settings for Premium Redis Cache's ([#130](https://github.com/terraform-providers/terraform-provider-azurerm/issues/130))
* `azurerm_storage_account` - exposing a formatted Connection String for Blob access ([#142](https://github.com/terraform-providers/terraform-provider-azurerm/issues/142))

BUG FIXES:

* `azurerm_cdn_endpoint` - fixing update of the `origin_host_header` ([#134](https://github.com/terraform-providers/terraform-provider-azurerm/issues/134))
* `azurerm_container_service` - exposes the FQDN of the `master_profile` as a computed field ([#125](https://github.com/terraform-providers/terraform-provider-azurerm/issues/125))
* `azurerm_key_vault` - fixing import / the validation on Access Policies ([#124](https://github.com/terraform-providers/terraform-provider-azurerm/issues/124))
* `azurerm_network_interface` - Normalizing the location field in the state ([#122](https://github.com/terraform-providers/terraform-provider-azurerm/issues/122))
* `azurerm_network_interface` - fixing a crash when importing a NIC with a Public IP ([#128](https://github.com/terraform-providers/terraform-provider-azurerm/issues/128))
* `azurerm_network_security_rule`: `network_security_group_name` is now `ForceNew` ([#138](https://github.com/terraform-providers/terraform-provider-azurerm/issues/138))
* `azurerm_subnet` now correctly detects changes to Network Securtiy Groups and Routing Table's ([#113](https://github.com/terraform-providers/terraform-provider-azurerm/issues/113))
* `azurerm_virtual_machine_scale_set` - making `storage_profile_os_disk`.`name` optional ([#129](https://github.com/terraform-providers/terraform-provider-azurerm/issues/129))

## 0.1.1 (June 21, 2017)

BUG FIXES:

* Sort ResourceID.Path keys for consistent output ([#116](https://github.com/terraform-providers/terraform-provider-azurerm/issues/116))

## 0.1.0 (June 20, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

FEATURES:

* **New Data Source:** `azurerm_resource_group` [[#15022](https://github.com/terraform-providers/terraform-provider-azurerm/issues/15022)](https://github.com/hashicorp/terraform/pull/15022)

IMPROVEMENTS:

* Add diff supress func to endpoint_location [[#15094](https://github.com/terraform-providers/terraform-provider-azurerm/issues/15094)](https://github.com/hashicorp/terraform/pull/15094)

BUG FIXES:

* Fixing the Deadlock issue ([#6](https://github.com/terraform-providers/terraform-provider-azurerm/issues/6))
