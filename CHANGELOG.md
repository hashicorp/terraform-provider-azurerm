## 1.8.0 (Unreleased)

IMPROVEMENTS:

* `azurerm_log_analytics_solution` - support for Sovereign Clouds [GH-1410]
* `azurerm_log_analytics_workspace` - support for Sovereign Clouds [GH-1410]
* `azurerm_log_analytics_workspace` - support for the `PerGB2018` SKU [GH-1079]
* `azurerm_sql_database` - support for importing from a bacpac backup [GH-972]
* `azurerm_network_interface` - addtional validation on properties [GH-1403]
* `azurerm_arm_loadbalancer` - addtional validation on properties [GH-1403]


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
