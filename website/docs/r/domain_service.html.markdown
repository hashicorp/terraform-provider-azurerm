---
subcategory: "Domain Service"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_domain_service"
sidebar_current: "docs-azurerm-resource-domain-service"
description: |-
  Manage Azure DomainService instance.
---

# azurerm_domain_service

Manage Azure DomainService instance.


## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the domain service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group within the user's subscription. The name is case insensitive. Changing this forces a new resource to be created.

* `location` - (Optional) Resource location Changing this forces a new resource to be created.

* `domain_security_settings` - (Optional) One `domain_security_setting` block defined below.

* `filtered_sync` - (Optional) Enabled or Disabled flag to turn on Group-based filtered sync Defaults to `Enabled`.

* `ldaps_settings` - (Optional) One `ldaps_setting` block defined below.

* `notification_settings` - (Optional) One `notification_setting` block defined below.

* `subnet_id` - (Required) The id of the subnet that Domain Services will be deployed on. Changing this forces a new resource to be created.

* `tags` - (Optional) Resource tags Changing this forces a new resource to be created.

---

The `domain_security_setting` block supports the following:

* `ntlm_v1` - (Optional) A flag to determine whether or not NtlmV1 is enabled or disabled. Defaults to `Enabled`.

* `tls_v1` - (Optional) A flag to determine whether or not TlsV1 is enabled or disabled. Defaults to `Enabled`.

* `sync_ntlm_passwords` - (Optional) A flag to determine whether or not SyncNtlmPasswords is enabled or disabled. Defaults to `Enabled`.

---

The `ldaps_setting` block supports the following:

* `ldaps` - (Optional) A flag to determine whether or not Secure LDAP is enabled or disabled. Defaults to `Enabled`.

* `pfx_certificate` - (Optional) The certificate required to configure Secure LDAP. The parameter passed here should be a base64encoded representation of the certificate pfx file.

* `pfx_certificate_password` - (Optional) The password to decrypt the provided Secure LDAP certificate pfx file.

* `external_access` - (Optional) A flag to determine whether or not Secure LDAP access over the internet is enabled or disabled. Defaults to `Enabled`.

* `external_access_ip_address` - (Computed) the ip address of Secure LDAP access over the internet
---

The `notification_setting` block supports the following:

* `notify_global_admins` - (Optional) Should global admins be notified Defaults to `Enabled`.

* `notify_dc_admins` - (Optional) Should domain controller admins be notified Defaults to `Enabled`.

* `additional_recipients` - (Optional) The list of additional recipients

## Attributes Reference

The following attributes are exported:

* `domain_controller_ip_address` - List of Domain Controller IP Address

* `id` - Resource Id

* `name` - Resource name
