## 0.1.6 (Unreleased)

FEATURES:

* **New Resource:** `azurerm_app_service_plan` [GH-1]
* **New Resource:** `azurerm_eventgrid_topic` [GH-260]
* **New Resource:** `azurerm_key_vault_secret` [GH-269]

IMPROVEMENTS:

* `azurerm_image` - added a default to the `caching` field [GH-259]
* `azurerm_key_vault` - validation for the `name` field [GH-270]
* `azurerm_network_interface` - support for multiple IP Configurations / setting the Primary IP Configuration [GH-245]
* `azurerm_sql_server` - added checks to handle `name` not being globally unique [GH-189]
* `azurerm_sql_server` - making `administrator_login` `ForceNew` [GH-189]
* `azurerm_sql_server` - migrate to using the azure-sdk-for-go [GH-189]
* `azurerm_virtual_machine` - Force recreation if `storage_data_disk`.`create_option` changes [GH-240]
* `azurerm_virtual_machine_scale_set` - Fix address issue when setting the `winrm` block [GH-271]
* updating to `v10.3.0-beta` of the Azure SDK for Go [GH-258]

BUG FIXES:

* `azurerm_cosmosdb_account` - fixing the validation on the name field [GH-263]
* `azurerm_sql_server` - handle deleted servers correctly [GH-189]

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
