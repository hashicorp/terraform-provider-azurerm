## 0.1.2 (Unreleased)

FEATURES:

* **New Data Source:** `azurerm_managed_disk` [GH-121]
* `azurerm_network_interface` now supports import [GH-119]

IMPROVEMENTS:

* Ensuring consistency in when storing the `location` field in the state for the `azurerm_availability_set`, `azurerm_express_route_circuit`, `azurerm_load_balancer`, `azurerm_local_network_gateway`, `azurerm_managed_disk`, `azurerm_network_security_group`
`azurerm_public_ip`, `azurerm_resource_group`, `azurerm_route_table`, `azurerm_storage_account`, `azurerm_virtual_machine` and `azurerm_virtual_network` resources [GH-123]

BUG FIXES:

* `azurerm_container_service` - exposes the FQDN of the master_profile as a computed field [GH-125]
* `azurerm_key_vault` - fixing import / the validation on Access Policies [GH-124]
* `azurerm_network_interface`: Normalizing the location field in the state [GH-122]
* `azurerm_subnet` now correctly detects changes to Network Securtiy Groups and Routing Table's [GH-113]

## 0.1.1 (June 21, 2017)

BUG FIXES:

* Sort ResourceID.Path keys for consistent output ([#116](https://github.com/terraform-providers/terraform-provider-azurerm/116))

## 0.1.0 (June 20, 2017)

BACKWARDS INCOMPATIBILITIES / NOTES:

FEATURES:

* **New Data Source:** `azurerm_resource_group` [[#15022](https://github.com/terraform-providers/terraform-provider-azurerm/15022)](https://github.com/hashicorp/terraform/pull/15022)

IMPROVEMENTS:

* Add diff supress func to endpoint_location [[#15094](https://github.com/terraform-providers/terraform-provider-azurerm/15094)](https://github.com/hashicorp/terraform/pull/15094)

BUG FIXES:

* Fixing the Deadlock issue ([#6](https://github.com/terraform-providers/terraform-provider-azurerm/6))
