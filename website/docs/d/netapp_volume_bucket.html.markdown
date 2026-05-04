---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume_bucket"
description: |-
  Gets information about an existing NetApp Files Volume Bucket (Object REST API).
---

# Data Source: azurerm_netapp_volume_bucket

Use this data source to access information about an existing NetApp Files Volume Bucket.

## Example Usage

```hcl
data "azurerm_netapp_volume_bucket" "example" {
  name                = "example-bucket"
  resource_group_name = "example-resources"
  account_name        = "example-anfaccount"
  pool_name           = "example-anfpool"
  volume_name         = "example-anfvolume"
}

output "bucket_status" {
  value = data.azurerm_netapp_volume_bucket.example.status
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the NetApp Volume Bucket.

* `resource_group_name` - (Required) The name of the Resource Group where the parent NetApp account exists.

* `account_name` - (Required) The name of the parent NetApp Account.

* `pool_name` - (Required) The name of the parent NetApp Capacity Pool.

* `volume_name` - (Required) The name of the parent NetApp Volume.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NetApp Volume Bucket.

* `volume_id` - The ARM ID of the parent NetApp Volume.

* `path` - The volume sub-path mounted inside the bucket.

* `permissions` - The bucket permission level (`ReadOnly` or `ReadWrite`).

* `file_system_user` - A `file_system_user` block as defined below.

* `server` - A `server` block as defined below.

* `key_vault` - A `key_vault` block as defined below (populated only when the bucket is configured against Azure Key Vault).

* `status` - The credentials status of the bucket. Possible values are `NoCredentialsSet`, `CredentialsExpired` and `Active`.

* `server_ip_address` - The IP address that backs the bucket endpoint.

* `server_certificate_common_name` - The Common Name (CN) of the bucket server certificate.

* `server_certificate_expiry_date` - The expiry date of the bucket server certificate, in RFC3339 format.

---

A `file_system_user` block exports the following:

* `nfs_user` - A `nfs_user` block as defined below (only set when the bucket is configured for NFS).

* `cifs_user` - A `cifs_user` block as defined below (only set when the bucket is configured for CIFS).

---

A `nfs_user` block exports the following:

* `group_id` - The POSIX group ID used by the bucket.

* `user_id` - The POSIX user ID used by the bucket.

---

A `cifs_user` block exports the following:

* `username` - The CIFS username used by the bucket.

---

A `server` block exports the following:

* `fqdn` - The DNS name that resolves to the bucket endpoint IP address.

* `on_certificate_conflict_action` - The action that runs when a certificate rotation conflicts with an existing certificate.

---

A `key_vault` block exports the following:

* `certificate_key_vault_uri` - The URI of the Azure Key Vault that stores the bucket server certificate.

* `certificate_name` - The name of the certificate object inside the Key Vault.

* `credentials_key_vault_uri` - The URI of the Azure Key Vault used to store the generated bucket access and secret keys.

* `credentials_secret_name` - The name of the secret in `credentials_key_vault_uri` that stores the generated bucket credentials.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Volume Bucket.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.NetApp` - 2026-01-01
