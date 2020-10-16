---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_instance"
description: |-
  Manages a MS SQL Managed instance.
---

# azurerm_mssql_managed_instance

Manages a MS SQL Managed instance.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "westus"
}

resource "azurerm_network_security_group" "example" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "sql_mi-network"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "internal"
  virtual_network_name = azurerm_virtual_network.example.name
  resource_group_name  = azurerm_resource_group.example.name
  address_prefixes     = ["10.0.1.0/24"]
  delegation {
    name = "miDelegation"

    service_delegation {
      name = "Microsoft.Sql/managedInstances"
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_route_table" "example" {
  name                = "example-routetable"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  route {
    name                   = "example"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_subnet_route_table_association" "example" {
  subnet_id      = azurerm_subnet.example.id
  route_table_id = azurerm_route_table.example.id
}

resource "azurerm_mssql_managed_instance" "dns" {
  name                         = "sql-dns-partner-mi"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  administrator_login          = "demoReadUser"
  administrator_login_password = "ReadUserDemo@12345"
  subnet_id                    = azurerm_subnet.example.id
  identity {
    type = "SystemAssigned"
  }
  sku {
    capacity = 8
    family   = "Gen5"
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
  }
  license_type                 = "BasePrice"
  collation                    = "SQL_Latin1_General_CP1_CI_AS"
  proxy_override               = "Redirect"
  storage_size_gb              = 64
  vcores                       = 8
  public_data_endpoint_enabled = true
  timezone_id                  = "UTC"
  minimal_tls_version          = "1.1"
}


resource "azurerm_mssql_managed_instance" "example" {
  name                         = "sql-mi"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  administrator_login          = "demoReadUser"
  administrator_login_password = "ReadUser@123456"
  subnet_id                    = azurerm_subnet.example.id
  identity {
    type = "SystemAssigned"
  }
  sku {
    capacity = 8
    family   = "Gen5"
    name     = "GP_Gen5"
    tier     = "GeneralPurpose"
  }
  license_type                 = "LicenseIncluded"
  collation                    = "SQL_Latin1_General_CP1_CI_AS"
  proxy_override               = "Redirect"
  storage_size_gb              = 64
  vcores                       = 8
  public_data_endpoint_enabled = false
  timezone_id                  = "Central America Standard Time"
  minimal_tls_version          = "1.2"
  dns_zone_partner             = azurerm_mssql_managed_instance.dns.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the MS SQL Managed Instance. Changing this forces a new resource to be created.

* `location` - (Required) The location of the Managed Instance.

* `resource_group_name` - (Required) The resource group name of the Managed instance.

* `administrator_login` - (Required) The admin login user name of the Managed instance. Changing this forces a new resource to be created.

* `administrator_login_password` - (Required) The admin login user password(sensitive) of the Managed instance.

* `collation` - (Optional) The collation of the Managed instance. Changing this forces a new resource to be created.

* `dns_zone_partner` - (Optional) The resource id of another managed instance whose DNS zone this managed instance will share after creation. Changing this forces a new resource to be created.

* `instance_pool_id` - (Optional) The Id of the instance pool this managed instance belongs to. Changing this forces a new resource to be created.

* `license_type` - (Optional) The license type. Possible values are `LicenseIncluded` (regular price inclusive of a new SQL license) and `BasePrice` (discounted AHB price for bringing your own SQL licenses),

* `maintenance_configuration_id` - (Optional) Specifies maintenance configuration id to apply to this managed instance. Changing this forces a new resource to be created.

* `create_mode` - (Optional) Specifies the mode of database creation. Possible values are `Default` (regular instance creation) and `Restore` (creates an instance by restoring a set of backups to specific point in time. `restore_point_in_time` and `source_managed_instance_id` must also be specified). Changing this forces a new resource to be created.

* `minimal_tls_version` - (Optional) Minimal TLS version. Allowed values: `None`, `1.0`, `1.1` and `1.2`.

* `proxy_override` - (Optional) Connection type used for connecting to the instance. Allowed values: `Default`, `Proxy` and `Redirect`.

* `public_data_endpoint_enabled` - (Optional) Boolean value indicating whether or not the public data endpoint is enabled. 

* `restore_point_in_time` - (Optional) Specifies the point in time (ISO8601 format) of the source database that will be restored to create the new database. Specified for `Restore` create mode. `source_managed_instance_id` must be specified too. Changing this forces a new resource to be created.

* `source_managed_instance_id` - (Optional) The resource id of the source managed instance associated with create operation of this instance. Specified for `Restore` create mode. `restore_point_in_time` must be specified too. Changing this forces a new resource to be created.

* `storage_size_gb` - (Optional) Storage size in GB. Minimum value: 32. Maximum value: 16384. Increments of 32 GB allowed only.

~> **NOTE:** The storage account type (used to store backups) is not supported yet in Go-SDK. 

* `subnet_id` - (Required) Subnet resource ID for the managed instance.

* `timezone_id` - (Optional) Id of the timezone. Allowed values are timezones supported by Windows. Changing this forces a new resource to be created.

~> **NOTE:** The supported timezones can be found from the link [timezones](https://docs.microsoft.com/en-us/azure/azure-sql/managed-instance/timezones-overview)

* `vcores` - (Optional) The number of vCores. Allowed values: 4, 8, 16, 24, 32, 40, 64 and 80.

* `sku` - (Optional) Managed instance SKU. The `sku` block is defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) resource tags.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Microsoft SQL Managed instance. At this time the only allowed value is `SystemAssigned`.

~> **NOTE:** The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned` and the managed instance has been created.

---

a `sku` block supports the following:

* `capacity` - (Optional) Capacity of the particular SKU.

* `family` - (Optional) If the service has different generations of hardware, for the same SKU, then that can be captured here.

* `name` - (Required) The name of the SKU, typically, a letter + Number code, e.g. P3.

* `size` - (Optional) Size of the particular SKU.

* `tier` - (Optional) The tier or edition of the particular SKU, e.g. Basic, Premium.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Managed instance.

* `fully_qualified_domain_name` - The FQDN of the Managed instance.

* `state` - The state of the Managed instance.

* `type` - The type of the Managed instance.

* `dns_zone` - The DNS zone of the Managed instance.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 600 minutes) Used when creating the Managed instance. 
* `update` - (Defaults to 1200 minutes) Used when updating the Managed instance(When SKU details are changed, it may take upto 12+ hours for the changes to take affect).
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed instance.
* `delete` - (Defaults to 600 minutes) Used when deleting the Managed instance(If the managed instance virtual cluster has more than 1 managed instances in it, first instance deletion may take upto 10 minutes. The deletion of last managed instance may take 3+ hours which deletes both the last managed instance and the empty virtual cluster).

## Import

SQL managed instance can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_managed_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Sql/managedInstances/sql-mi
```
