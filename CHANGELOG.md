## 0.2.1 (Unreleased)
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
