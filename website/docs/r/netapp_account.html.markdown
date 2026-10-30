---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_account"
description: |-
  Manages a NetApp Account.
---

# azurerm_netapp_account

Manages a NetApp Account.

~> **Note:** Azure allows only one active directory can be joined to a single subscription at a time for NetApp Account.

## NetApp Account Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

data "azurerm_client_config" "current" {
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "anf-user-assigned-identity"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_account" "example" {
  name                = "netappaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  active_directory {
    username            = "aduser"
    password            = "aduserpwd"
    smb_server_name     = "SMBSERVER"
    dns_servers         = ["1.2.3.4"]
    domain              = "westcentralus.com"
    organizational_unit = "OU=FirstLevel"
  }

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.example.id
    ]
  }
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group where the NetApp Account should be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `active_directory` - (Optional) A `active_directory` block as defined below.

* `identity` - (Optional) The `identity` block where it is used when customer managed keys based encryption will be enabled as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

The `active_directory` block supports the following:

* `dns_servers` - (Required) A list of DNS server IP addresses for the Active Directory domain. Only allows `IPv4` address.

* `domain` - (Required) The name of the Active Directory domain.

* `smb_server_name` - (Required) The NetBIOS name which should be used for the NetApp SMB Server, which will be registered as a computer account in the AD and used to mount volumes.

* `username` - (Required) The Username of Active Directory Domain Administrator.

* `password` - (Required) The password associated with the `username`.

* `administrator_users` - (Optional) A list of Active Directory users to be added to the built-in Administrators active directory group. Maximum of 255 users are allowed.

* `aes_encryption_enabled` - (Optional) If enabled, AES encryption will be enabled for SMB communication. Defaults to `false`.

* `backup_users` - (Optional) A list of Active Directory users to be added to the built-in Backup Operators active directory group. Maximum of 255 users are allowed.

* `dc_connection_encryption_enabled` - (Optional) If enabled, the connection between Active Directory and SMB server will be encrypted. Defaults to `false`.

* `kerberos_ad_name` - (Optional) Name of the active directory machine.

* `kerberos_kdc_ip` - (Optional) kdc server IP addresses for the active directory machine.

~> **Note:** If you plan on using **Kerberos** volumes, both `ad_name` and `kdc_ip` are required in order to create the volume.

* `ldap_client_preferred_server_ips` - (Optional) Specifies the preferred LDAP server IPs for the LDAP client. Maximum of 2 `IPv4` addresses are allowed.

* `ldap_over_tls_enabled` - (Optional) Specifies whether or not the LDAP traffic needs to be secured via TLS. Defaults to `false`.

* `ldap_search_scope` - (Optional) An `ldap_search_scope` block as defined below.

* `ldap_signing_enabled` - (Optional) Specifies whether or not the LDAP traffic needs to be signed. Defaults to `false`.

* `local_nfs_users_with_ldap_allowed` - (Optional) If enabled, NFS client local users can also (in addition to LDAP users) access the NFS volumes. Defaults to `false`.

* `organizational_unit` - (Optional) The Organizational Unit (OU) within Active Directory where machines will be created. If blank, defaults to `CN=Computers`.

* `security_users` - (Optional) A list of Active Directory users to be given SeSecurityPrivilege privilege. Maximum of 255 users are allowed.

* `site_name` - (Optional) The Active Directory site the service will limit Domain Controller discovery to. If blank, defaults to `Default-First-Site-Name`.

* `server_root_ca_certificate` - (Optional) When LDAP over SSL/TLS is enabled, the LDAP client is required to have a *base64 encoded Active Directory Certificate Service's self-signed root CA certificate*, this optional parameter is used only for dual protocol with LDAP user-mapping volumes. Required if `ldap_over_tls_enabled` is set to `true`.

---

The `identity` block supports the following:

* `type` - (Required) The identity type, which can be `SystemAssigned` or `UserAssigned`. Only one type at a time is supported by Azure NetApp Files.

* `identity_ids` - (Optional) The identity id of the user assigned identity to use when type is `UserAssigned`

~> **Note:** Changing identity type from `SystemAssigned` to `UserAssigned` is a supported operation but the reverse is not supported from within Terraform Azure NetApp Files module.

---

The `ldap_search_scope` block supports the following:

* `group_dn` - (Optional) Specifies the group DN, which overrides the base DN for group lookups.

* `group_membership_filter` - (Optional) Specifies the custom LDAP search filter to be used when looking up group membership from LDAP server.

* `user_dn` - (Optional) Specifies the user DN, which overrides the base DN for user lookups.

~> **Note:** When using the `ldap_search_scope` block, at least one of the `group_dn`, `group_membership_filter`, or `user_dn` must be specified.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NetApp Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the NetApp Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Account.
* `update` - (Defaults to 30 minutes) Used when updating the NetApp Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the NetApp Account.

## Import

NetApp Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_account.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.NetApp/netAppAccounts/account1
```

~> **Note:** When importing a NetApp account, the `active_directory.password` and `active_directory.server_root_ca_certificate` values *cannot* be retrieved from the Azure API and will need to be redeclared within the resource.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.NetApp` - 2026-05-01
