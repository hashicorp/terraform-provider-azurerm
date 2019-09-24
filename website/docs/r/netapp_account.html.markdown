---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_account"
sidebar_current: "docs-azurerm-resource-netapp-account"
description: |-
  Manage Azure NetApp Account instance.
---

# azurerm_netapp_account

Manage Azure NetApp Account instance.


## NetApp Account Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "acctestRG"
  location = "Eastus2"
}

resource "azurerm_netapp_account" "example" {
  name                = "acctestnetappaccount"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  active_directories {
    username            = "aduser"
    password            = "aduser"
    smb_server_name     = "SMBSERVER"
    dns                 = "1.2.3.4"
    domain              = "westcentralus.com"
    organizational_unit = "OU=FirstLevel"
  }

  tags = {
    env = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The resource group name of the NetApp Account. Changing this forces a new resource to be created.

* `location` - (Optional) Resource location. Changing this forces a new resource to be created.

* `active_directories` - (Optional) One or more `active_directories` block defined below.

* `tags` - (Optional) Resource tags. Changing this forces a new resource to be created.

---

The `active_directories` block supports the following:

* `dns` - (Optional) Comma separated list of DNS server IP addresses for the Active Directory domain.

* `domain` - (Optional) Name of the Active Directory domain.

* `organizational_unit` - (Optional) The Organizational Unit (OU) within the Windows Active Directory.

* `password` - (Optional) Plain text password of Active Directory domain administrator.

* `smb_server_name` - (Optional) NetBIOS name of the SMB server. This name will be registered as a computer account in the AD and used to mount volumes.

* `username` - (Optional) Username of Active Directory domain administrator.

---

## Attributes Reference

The following attributes are exported:

* `id` - Resource id.

## Import

NetApp Account can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG/providers/Microsoft.NetApp/netAppAccounts/
```