---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_account"
sidebar_current: "docs-azurerm-datasource-netapp-account"
description: |-
  Gets information about an existing NetApp Account
---

# Data Source: azurerm_netapp_account

Use this data source to access information about an existing NetApp Account.


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

* `location` - Resource location.

* `active_directories` - One or more `active_directories` block defined below.

* `tags` - Resource tags.

---

The `active_directories` block contains the following:

* `id` - Resource ID.

* `dns` - Comma separated list of DNS server IP addresses for the Active Directory domain.

* `domain` - Name of the Active Directory domain.

* `organizational_unit` - The Organizational Unit (OU) within the Windows Active Directory.

* `password` - Plain text password of Active Directory domain administrator.

* `smb_server_name` - NetBIOS name of the SMB server. This name will be registered as a computer account in the AD and used to mount volumes.

* `status` - Status of the Active Directory.

* `username` - Username of Active Directory domain administrator.
