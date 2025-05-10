---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_managed_instance_security_alert_policy"
description: |-
  Manages a Security Alert Policy for an MS SQL Managed Instance.

---

# azurerm_mssql_managed_instance_security_alert_policy

Manages a Security Alert Policy for an MS SQL Managed Instance.

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
  disable_bgp_route_propagation = false
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

resource "azurerm_mssql_managed_instance_security_alert_policy" "example" {
  resource_group_name        = azurerm_resource_group.example.name
  managed_instance_name      = azurerm_mssql_managed_instance.example.name
  enabled                    = true
  storage_endpoint           = azurerm_storage_account.example.primary_blob_endpoint
  storage_account_access_key = azurerm_storage_account.example.primary_access_key
  disabled_alerts = [
    "Sql_Injection",
    "Data_Exfiltration"
  ]
  retention_days = 20
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_name` - (Required) The name of the resource group that contains the MS SQL Managed Instance. Changing this forces a new resource to be created.

* `managed_instance_name` - (Required) Specifies the name of the MS SQL Managed Instance. Changing this forces a new resource to be created.

* `disabled_alerts` - (Optional) Specifies an array of alerts that are disabled. Possible values are `Sql_Injection`, `Sql_Injection_Vulnerability`, `Access_Anomaly`, `Data_Exfiltration`, `Unsafe_Action` and `Brute_Force`.

* `enabled` - (Optional) Specifies the state of the Security Alert Policy, whether it is enabled or disabled. Possible values are `true`, `false`.

* `email_account_admins_enabled` - (Optional) Boolean flag which specifies if the alert is sent to the account administrators or not. Defaults to `false`.

* `email_addresses` - (Optional) Specifies an array of email addresses to which the alert is sent.

* `retention_days` - (Optional) Specifies the number of days to keep in the Threat Detection audit logs. Defaults to `0`.

* `storage_endpoint` - (Optional) Specifies the blob storage endpoint (e.g. https://example.blob.core.windows.net). This blob storage will hold all Threat Detection audit logs.

* `storage_account_access_key` - (Optional) Specifies the identifier key of the Threat Detection audit storage account. This is mandatory when you use `storage_endpoint` to specify a storage account blob endpoint.

-> **Note:** Please note that storage accounts configured with `shared_access_key_enabled = false` cannot be used to configure `azurerm_mssql_managed_instance_security_alert_policy` with `storage_endpoint` for now.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the MS SQL Managed Instance Security Alert Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the MS SQL Managed Instance Security Alert Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the MS SQL Managed Instance Security Alert Policy.
* `update` - (Defaults to 30 minutes) Used when updating the MS SQL Managed Instance Security Alert Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the MS SQL Managed Instance Security Alert Policy.

## Import

MS SQL Managed Instance Security Alert Policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mssql_managed_instance_security_alert_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acceptanceTestResourceGroup1/providers/Microsoft.Sql/managedInstances/instance1/securityAlertPolicies/Default
```
