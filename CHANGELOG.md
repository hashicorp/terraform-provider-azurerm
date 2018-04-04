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
* `azurerm_cosmosdb_account` - increasing the maximum value for `max_interval_in_seconds` from 100s to 86400s (1 day) [[#1000](https://github.com/terraform-providers/terraform-provider-azurerm/issues/1000)] 
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
