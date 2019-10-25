---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_account"
sidebar_current: "docs-azurerm-datasource-netapp-account"
description: |-
  Gets information about an existing NetApp Account
---

# Data Source: azurerm_netapp_account

Uses this data source to access information about an existing NetApp Account.

## NetApp Account Usage

```hcl
data "azurerm_netapp_account" "example" {
  resource_group_name = "acctestRG"
  name                = "acctestnetappaccount"
}

output "netapp_account_id" {
  value = "${data.azurerm_netapp_account.example.id}"
}
```


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Account.

* `resource_group_name` - (Required) The Name of the Resource Group where the NetApp Account exists.


## Attributes Reference

The following attributes are exported:

* `location` - The Azure Region where the NetApp Account exists.

* `active_directory` - An `active_directory` block defined below.

---

The `active_directory` block exports the following:

* `dns_servers` - A list of IP Addresses used as DNS Servers for the Active Directory domain.

* `domain` - The Active Directory Domain Name.

* `smb_server_name` - The NetBIOS name of the SMB Server.

* `username` - The Username of Active Directory domain administrator.

* `organizational_unit` - The Organizational Unit (OU) within the Active Directory Domain where the NetApp is located.
