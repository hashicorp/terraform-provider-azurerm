---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_account"
sidebar_current: "docs-azurerm-resource-netapp-account"
description: |-
  Manages Azure NetApp Account instance.
---

# azurerm_netapp_account

Manages Azure NetApp Account instance.


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

  active_directory {
    username            = "aduser"
    password            = "aduserpwd"
    smb_server_name     = "SMBSERVER"
    dns                 = ["1.2.3.4"]
    domain              = "westcentralus.com"
    organizational_unit = "OU=FirstLevel"
  }

  tags = {
    foo = "bar"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Account resides. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `active_directory` - (Optional) One or more `active_directory` block defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `active_directory` block supports the following:

* `dns_servers` - (Required) A list of DNS server IP addresses for the Active Directory domain.

* `domain` - (Required) Name of the Active Directory domain.

* `smb_server_name` - (Required) NetBIOS name of the SMB server. This name will be registered as a computer account in the AD and used to mount volumes.

* `username` - (Required) Username of Active Directory domain administrator, which have permissions to create a SMB machine account in the AD domain.

* `password` - (Required) Plain text password of Active Directory domain administrator.

* `organizational_unit` - (Optional) The Organizational Unit (OU) within the Windows Active Directory.

---

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the NetApp resource.

* `active_directory` - One or more `active_directory` block defined below.

---

The `active_directory` block supports the following:

* `id` - The resource id of Active Directory.

---

## Import

NetApp Account can be imported using the `resource id`, e.g.

```shell
$ terraform import azurerm_netapp_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctestRG/providers/Microsoft.NetApp/netAppAccounts/acctestnetappaccount
```