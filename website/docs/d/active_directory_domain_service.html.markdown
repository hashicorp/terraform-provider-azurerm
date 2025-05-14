---
subcategory: "Active Directory Domain Services"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_active_directory_domain_service"
description: |-
  Gets information about an Active Directory Domain Service.
---

# Data Source: azurerm_active_directory_domain_service

Gets information about an Active Directory Domain Service.

-> **Note:** At present this data source only supports **User Forest** mode and _not_ **Resource Forest** mode. [Read more](https://docs.microsoft.com/azure/active-directory-domain-services/concepts-resource-forest) about the different operation modes for this service.

## Example Usage

```hcl
data "azurerm_active_directory_domain_service" "example" {
  name                = "example-aadds"
  resource_group_name = "example-aadds-rg"
}
```

## Argument Reference

* `name` - (Required) The display name for your managed Active Directory Domain Service resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Domain Service should exist. Changing this forces a new resource to be created.

## Attributes Reference

* `id` - The ID of the Domain Service.

* `deployment_id` - A unique ID for the managed domain deployment.

* `domain_configuration_type` - The forest type used by the managed domain. One of `ResourceTrusting`, for a _Resource Forest_, or blank, for a _User Forest_.
  
* `domain_name` - The Active Directory domain of the Domain Service. See [official documentation](https://docs.microsoft.com/azure/active-directory-domain-services/tutorial-create-instance#create-a-managed-domain) for constraints and recommendations.

* `filtered_sync_enabled` - Whether group-based filtered sync (also called scoped synchronisation) is enabled.

* `secure_ldap` - A `secure_ldap` block as defined below.

* `location` - The Azure location where the Domain Service exists.

* `notifications` - A `notifications` block as defined below.

* `replica_sets` - One or more `replica_set` blocks as defined below.

* `security` - A `security` block as defined below.

* `sku` - The SKU of the Domain Service resource. One of `Standard`, `Enterprise` or `Premium`.

* `tags` - A mapping of tags assigned to the resource.

---

A `secure_ldap` block exports the following:

* `enabled` - Whether secure LDAP is enabled for the managed domain.

* `external_access_enabled` - Whether external access to LDAPS over the Internet, is enabled.
  
* `external_access_ip_address` - The publicly routable IP address for LDAPS clients to connect to.

* `pfx_certificate` - The certificate to use for LDAPS, as a base64-encoded TripleDES-SHA1 encrypted PKCS#12 bundle (PFX file).

---

A `notifications` block exports the following:

* `additional_recipients` - A list of additional email addresses to notify when there are alerts in the managed domain.

* `notify_dc_admins` - Whethermembers of the _AAD DC Administrators_ group are notified when there are alerts in the managed domain.

* `notify_global_admins` - Whether all Global Administrators are notified when there are alerts in the managed domain.

---

A `replica_set` block exports the following:

* `domain_controller_ip_addresses` - A list of subnet IP addresses for the domain controllers in the replica set, typically two.

* `external_access_ip_address` - The publicly routable IP address for the domain controllers in the replica set.

* `location` - The Azure location in which the replica set resides.

* `replica_set_id` - A unique ID for the replica set.

* `service_status` - The current service status for the replica set.

* `subnet_id` - The ID of the subnet in which the replica set resides.

---

A `security` block exports the following:

* `kerberos_armoring_enabled` - (Optional) Whether the Kerberos Armoring is enabled.

* `kerberos_rc4_encryption_enabled` - (Optional) Whether the Kerberos RC4 Encryption is enabled.

* `ntlm_v1_enabled` - Whether legacy NTLM v1 support is enabled.

* `sync_kerberos_passwords` - Whether Kerberos password hashes are synchronized to the managed domain.

* `sync_ntlm_passwords` - Whether NTLM password hashes are synchronized to the managed domain.

* `sync_on_prem_passwords` - Whether on-premises password hashes are synchronized to the managed domain.

* `tls_v1_enabled` - Whether legacy TLS v1 support is enabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Domain Service.
