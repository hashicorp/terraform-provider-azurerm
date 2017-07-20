## 0.1.3 (Unreleased)

FEATURES:

* **New Resource:** `azurerm_dns_ptr_record` [GH-141]
* **New Resource:**`azurerm_image` [GH-8]
* **New Resource:** `azurerm_servicebus_queue` [GH-151]

IMPROVEMENTS:

* `azurerm_client_config` - added a `service_principal_object_id` attribute to the data source [GH-175]
* `azurerm_search_service` - added import support [GH-172]
* `azurerm_servicebus_topic` - added a `status` field to allow disabling the topic [GH-150]
* `azurerm_storage_account` - Added support for Require secure transfer [GH-167]
* `azurerm_storage_table` - updating the name validation [GH-143]
* `azurerm_virtual_machine` - making `admin_password` optional for Linux VM's [GH-154]
* `azurerm_virtual_machine_scale_set` - adding a `plan` block for Marketplace images [GH-161]

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
