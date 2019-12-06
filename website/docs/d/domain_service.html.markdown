---
subcategory: "Domain Service"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_domain_service"
sidebar_current: "docs-azurerm-resource-domain-service"
description: |-
  Gets information about an existing Azure DomainService instance.
---

# azurerm_domain_service

Use this data source to access information about an existing DomainService instance.

## Example Usage

```hcl
data "azurerm_domain_service" "example" {
	name                  = "example.onmicrosoft.com"
	resource_group_name   = "example"
}

output "domain_service_id" {
  value = "${data.azurerm_domain_service.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the domain service. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group within the user's subscription. The name is case insensitive. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - Resource Id

* `name` - Resource name

* `domain_controller_ip_address` - List of Domain Controller IP Address

* `location` -  Resource location Changing this forces a new resource to be created.

* `domain_security_settings` -  One `domain_security_setting` block defined below.

* `filtered_sync` -  Whether turn on Group-based filtered sync.

* `ldaps_settings` -  One `ldaps_setting` block defined below.

* `notification_settings` -  One `notification_setting` block defined below.

* `subnet_id` - The id of the subnet that Domain Services will be deployed on.

* `tags` -  Resource tags Changing this forces a new resource to be created.

---

The `domain_security_setting` block supports the following:

* `ntlm_v1` -  A flag to determine whether or not NtlmV1 is enabled or disabled.

* `tls_v1` -  A flag to determine whether or not TlsV1 is enabled or disabled.

* `sync_ntlm_passwords` -  A flag to determine whether or not SyncNtlmPasswords is enabled or disabled.

---

The `ldaps_setting` block supports the following:

* `ldaps` -  A flag to determine whether or not Secure LDAP is enabled or disabled. Defaults to `Enabled`.

* `pfx_certificate` -  The certificate required to configure Secure LDAP. The parameter passed here should be a base64encoded representation of the certificate pfx file.

* `pfx_certificate_password` -  The password to decrypt the provided Secure LDAP certificate pfx file.

* `external_access` -  A flag to determine whether or not Secure LDAP access over the internet is enabled or disabled. Defaults to `Enabled`.

* `external_access_ip_address` - the ip address of Secure LDAP access over the internet
---

The `notification_setting` block supports the following:

* `notify_global_admins` -  Should global admins be notified Defaults to `Enabled`.

* `notify_dc_admins` -  Should domain controller admins be notified Defaults to `Enabled`.

* `additional_recipients` -  The list of additional recipients
