---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_instance"
description: |-
  Manages a Microsoft SQL Azure Managed Instance.
---

# azurerm_mssql_managed_instance

Manages a Microsoft SQL Azure Managed Instance.

~> **Note:** All arguments including the administrator login and password will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

~> **Note:** SQL Managed Instance needs permission to read Azure Active Directory when configuring the AAD administrator. [Read more about provisioning AAD administrators](https://learn.microsoft.com/en-us/azure/azure-sql/database/authentication-aad-configure?view=azuresql#provision-azure-ad-admin-sql-managed-instance).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "database-rg"
  location = "West Europe"
}

resource "azurerm_network_security_group" "example" {
  name                = "mi-security-group"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}


resource "azurerm_network_security_rule" "allow_management_inbound" {
  name                        = "allow_management_inbound"
  priority                    = 106
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges     = ["9000", "9003", "1438", "1440", "1452"]
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "allow_misubnet_inbound" {
  name                        = "allow_misubnet_inbound"
  priority                    = 200
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "10.0.0.0/24"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "allow_health_probe_inbound" {
  name                        = "allow_health_probe_inbound"
  priority                    = 300
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "AzureLoadBalancer"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "allow_tds_inbound" {
  name                        = "allow_tds_inbound"
  priority                    = 1000
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "1433"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "deny_all_inbound" {
  name                        = "deny_all_inbound"
  priority                    = 4096
  direction                   = "Inbound"
  access                      = "Deny"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "allow_management_outbound" {
  name                        = "allow_management_outbound"
  priority                    = 102
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_ranges     = ["80", "443", "12000"]
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "allow_misubnet_outbound" {
  name                        = "allow_misubnet_outbound"
  priority                    = 200
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "10.0.0.0/24"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_network_security_rule" "deny_all_outbound" {
  name                        = "deny_all_outbound"
  priority                    = 4096
  direction                   = "Outbound"
  access                      = "Deny"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = azurerm_resource_group.example.name
  network_security_group_name = azurerm_network_security_group.example.name
}

resource "azurerm_virtual_network" "example" {
  name                = "vnet-mi"
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.example.location
}

resource "azurerm_subnet" "example" {
  name                 = "subnet-mi"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.0.0/24"]

  delegation {
    name = "managedinstancedelegation"

    service_delegation {
      name    = "Microsoft.Sql/managedInstances"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action", "Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action", "Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action"]
    }
  }
}

resource "azurerm_subnet_network_security_group_association" "example" {
  subnet_id                 = azurerm_subnet.example.id
  network_security_group_id = azurerm_network_security_group.example.id
}

resource "azurerm_route_table" "example" {
  name                          = "routetable-mi"
  location                      = azurerm_resource_group.example.location
  resource_group_name           = azurerm_resource_group.example.name
  bgp_route_propagation_enabled = true
  depends_on = [
    azurerm_subnet.example,
  ]
}

resource "azurerm_subnet_route_table_association" "example" {
  subnet_id      = azurerm_subnet.example.id
  route_table_id = azurerm_route_table.example.id
}

resource "azurerm_mssql_managed_instance" "example" {
  name                = "managedsqlinstance"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  license_type       = "BasePrice"
  sku_name           = "GP_Gen5"
  storage_size_in_gb = 32
  subnet_id          = azurerm_subnet.example.id
  vcores             = 4

  administrator_login          = "mradministrator"
  administrator_login_password = "thisIsDog11"

  depends_on = [
    azurerm_subnet_network_security_group_association.example,
    azurerm_subnet_route_table_association.example,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `license_type` - (Required) What type of license the Managed Instance will use. Possible values are `LicenseIncluded` and `BasePrice`.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `name` - (Required) The name of the SQL Managed Instance. This needs to be globally unique within Azure. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the SQL Managed Instance. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for the SQL Managed Instance. Valid values include `GP_Gen4`, `GP_Gen5`, `GP_Gen8IM`, `GP_Gen8IH`, `BC_Gen4`, `BC_Gen5`, `BC_Gen8IM` or `BC_Gen8IH`.

* `storage_size_in_gb` - (Required) Maximum storage space for the SQL Managed instance. This should be a multiple of 32 (GB).

* `subnet_id` - (Required) The subnet resource id that the SQL Managed Instance will be associated with.

* `vcores` - (Required) Number of cores that should be assigned to the SQL Managed Instance. Values can be `8`, `16`, or `24` for Gen4 SKUs, or `4`, `6`, `8`, `10`, `12`, `16`, `20`, `24`, `32`, `40`, `48`, `56`, `64`, `80`, `96` or `128` for Gen5 SKUs.

* `administrator_login` - (Optional) The administrator login name for the new SQL Managed Instance. Changing this forces a new resource to be created.

* `administrator_login_password` - (Optional) The password associated with the `administrator_login` user. Needs to comply with Azure's [Password Policy](https://msdn.microsoft.com/library/ms161959.aspx)

* `azure_active_directory_administrator` - (Optional) An `azure_active_directory_administrator` block as defined below.

* `collation` - (Optional) Specifies how the SQL Managed Instance will be collated. Default value is `SQL_Latin1_General_CP1_CI_AS`. Changing this forces a new resource to be created.

* `database_format` - (Optional) Specifies the internal format of the SQL Managed Instance databases specific to the SQL engine version. Possible values are `AlwaysUpToDate` and `SQLServer2022`. Defaults to `SQLServer2022`.

~> **Note:** Changing `database_format` from `AlwaysUpToDate` to `SQLServer2022` forces a new SQL Managed Instance to be created.

* `dns_zone_partner_id` - (Optional) The ID of the SQL Managed Instance which will share the DNS zone. This is a prerequisite for creating an `azurerm_sql_managed_instance_failover_group`. Setting this after creation forces a new resource to be created.

* `hybrid_secondary_usage` - (Optional) Specifies the hybrid secondary usage for disaster recovery of the SQL Managed Instance. Possible values are `Active` and `Passive`. Defaults to `Active`.

* `identity` - (Optional) An `identity` block as defined below.

* `maintenance_configuration_name` - (Optional) The name of the Public Maintenance Configuration window to apply to the SQL Managed Instance. Valid values include `SQL_Default` or an Azure Location in the format `SQL_{Location}_MI_{Size}`(for example `SQL_EastUS_MI_1`). Defaults to `SQL_Default`.

* `minimum_tls_version` - (Optional) The Minimum TLS Version. Default value is `1.2` Valid values include `1.0`, `1.1`, `1.2`.

~> **Note:** Azure Services will require TLS 1.2+ by August 2025, please see this [announcement](https://azure.microsoft.com/en-us/updates/v2/update-retirement-tls1-0-tls1-1-versions-azure-services/) for more.

* `proxy_override` - (Optional) Specifies how the SQL Managed Instance will be accessed. Default value is `Default`. Valid values include `Default`, `Proxy`, and `Redirect`.

* `public_data_endpoint_enabled` - (Optional) Is the public data endpoint enabled? Default value is `false`.

* `service_principal_type` - (Optional) The service principal type. The only possible value is `SystemAssigned`.

* `storage_account_type` - (Optional) Specifies the storage account type used to store backups for this database. Possible values are `GRS`, `GZRS`, `LRS`, and `ZRS`. Defaults to `GRS`.

* `zone_redundant_enabled` - (Optional) Specifies whether or not the SQL Managed Instance is zone redundant. Defaults to `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

* `timezone_id` - (Optional) The TimeZone ID that the SQL Managed Instance will be operating in. Default value is `UTC`. Changing this forces a new resource to be created.

---

An `azure_active_directory_administrator` block supports the following:

* `login_username` - (Required) The login username of the Azure AD Administrator of this SQL Managed Instance.

* `object_id` - (Required) The object id of the Azure AD Administrator of this SQL Managed Instance.

* `principal_type` - (Required) The principal type of the Azure AD Administrator of this SQL Managed Instance. Possible values are `Application`, `Group`, `User`.

* `azuread_authentication_only_enabled` - (Optional) Specifies whether only Azure AD authentication can be used to log in to this SQL Managed Instance. When `true`, the `administrator_login` and `administrator_login_password` properties can be omitted. Defaults to `false`.

* `tenant_id` - (Optional) The tenant id of the Azure AD Administrator of this SQL Managed Instance. Should be specified if the Azure AD Administrator is homed in a different tenant to the SQL Managed Instance.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this SQL Managed Instance. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this SQL Managed Instance. Required when `type` includes `UserAssigned`.

~> **Note:** The assigned `principal_id` and `tenant_id` can be retrieved after the identity `type` has been set to `SystemAssigned` and SQL Managed Instance has been created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The SQL Managed Instance ID.

* `dns_zone` - The Dns Zone where the SQL Managed Instance is located.

* `fqdn` - The fully qualified domain name of the Azure Managed SQL Instance

---

An `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Identity of this SQL Managed Instance.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Identity of this SQL Managed Instance.

-> **Note:** You can access the Principal ID via `azurerm_mssql_managed_instance.example.identity[0].principal_id` and the Tenant ID via `azurerm_mssql_managed_instance.example.identity[0].tenant_id`

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 24 hours) Used when creating the Microsoft SQL Managed Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the Microsoft SQL Managed Instance.
* `update` - (Defaults to 24 hours) Used when updating the Microsoft SQL Managed Instance.
* `delete` - (Defaults to 24 hours) Used when deleting the Microsoft SQL Managed Instance.

## Import

Microsoft SQL Managed Instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_managed_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Sql/managedInstances/myserver
```
