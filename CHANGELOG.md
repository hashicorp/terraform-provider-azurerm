## 0.1.2 (Unreleased)

FEATURES:

* **New Data Source:** `azurerm_managed_disk` [GH-121]
* `azurerm_network_interface` now supports import [GH-119]

BUG FIXES:

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
