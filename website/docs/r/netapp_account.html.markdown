---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_account"
sidebar_current: "docs-azurerm-resource-netapp-account"
description: |-
  Manages a NetApp Account.
---

# azurerm_netapp_account

Manages a NetApp Account.

~> **NOTE:** Azure allows only one active directory can be joined to a single subscription at a time for NetApp Account.

## NetApp Account Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_netapp_account" "example" {
  name                = "example-netapp"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  active_directory {
    username            = "aduser"
    password            = "aduserpwd"
    smb_server_name     = "SMBSERVER"
    dns                 = ["1.2.3.4"]
    domain              = "westcentralus.com"
    organizational_unit = "OU=FirstLevel"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Account should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `active_directory` - (Optional) A `active_directory` block as defined below.

---

The `active_directory` block supports the following:

* `dns_servers` - (Required) A list of DNS server IP addresses for the Active Directory domain. Only allows `IPv4` address.

* `domain` - (Required) The name of the Active Directory domain.

* `smb_server_name` - (Required) The NetBIOS name which should be used for the NetApp SMB Server, which will be registered as a computer account in the AD and used to mount volumes.

* `username` - (Required) The Username of Active Directory Domain Administrator.

* `password` - (Required) The password associated with the `username`.

* `organizational_unit` - (Optional) The Organizational Unit (OU) within the Active Directory Domain.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NetApp Account.

---

## Import

NetApp Accounts can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1
```
