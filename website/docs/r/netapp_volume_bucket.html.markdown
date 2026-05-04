---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume_bucket"
description: |-
  Manages a NetApp Files Volume Bucket (Object REST API).
---

# azurerm_netapp_volume_bucket

Manages a NetApp Files Volume Bucket. Buckets expose the contents of an Azure NetApp Files volume (or a sub-path within it) as an S3-compatible object endpoint via the Azure NetApp Files Object REST API.

~> **Note:** The Object REST API feature is in preview and must be registered on the subscription via `Microsoft.NetApp / ANFEnableObjectRESTAPI` before buckets can be created. See [Configure access to the Azure NetApp Files Object REST API](https://learn.microsoft.com/en-us/azure/azure-netapp-files/object-rest-api-access-configure) for the registration command and Key Vault prerequisites.

~> **Note:** Buckets are supported on cool-access and large NetApp volumes. Buckets are not supported on cache volumes. Deleting the parent volume cascade-deletes its buckets.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_virtual_network" "example" {
  name                = "example-vnet"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "example" {
  name                 = "example-delegated"
  resource_group_name  = azurerm_resource_group.example.name
  virtual_network_name = azurerm_virtual_network.example.name
  address_prefixes     = ["10.0.2.0/24"]

  delegation {
    name = "netapp"

    service_delegation {
      name    = "Microsoft.Netapp/volumes"
      actions = ["Microsoft.Network/networkinterfaces/*", "Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azurerm_netapp_account" "example" {
  name                = "example-anfaccount"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_netapp_pool" "example" {
  name                = "example-anfpool"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  service_level       = "Standard"
  size_in_tb          = 4
}

resource "azurerm_netapp_volume" "example" {
  name                = "example-anfvolume"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  account_name        = azurerm_netapp_account.example.name
  pool_name           = azurerm_netapp_pool.example.name
  volume_path         = "example-vol"
  service_level       = "Standard"
  subnet_id           = azurerm_subnet.example.id
  storage_quota_in_gb = 100
  protocols           = ["NFSv3"]
}

resource "azurerm_netapp_volume_bucket" "example" {
  name      = "example-bucket"
  volume_id = azurerm_netapp_volume.example.id

  file_system_user {
    nfs_user {
      group_id = 1000
      user_id  = 1000
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The S3-compatible name of the bucket. Must be 3-63 characters long, DNS-compliant (lowercase letters, digits, hyphens or periods), must start and end with a letter or number and must not look like an IPv4 address. Changing this forces a new resource to be created.

* `volume_id` - (Required) The ARM ID of the parent NetApp Volume the bucket attaches to. Changing this forces a new resource to be created.

* `file_system_user` - (Required) A `file_system_user` block as defined below. Exactly one of `nfs_user` or `cifs_user` must be specified.

---

* `path` - (Optional) The volume sub-path mounted inside the bucket. Defaults to `/`. Changing this forces a new resource to be created.

* `permissions` - (Optional) The bucket permission level. Possible values are `ReadOnly` and `ReadWrite`. Defaults to `ReadOnly`.

* `server` - (Optional) A `server` block as defined below. Used to provide the bucket server FQDN and a directly uploaded PEM certificate. Mutually exclusive with `key_vault`.

* `key_vault` - (Optional) A `key_vault` block as defined below. Used to source the server certificate and to store generated credentials in Azure Key Vault. Mutually exclusive with `server.0.certificate_pem`.

---

A `file_system_user` block supports the following:

* `nfs_user` - (Optional) A `nfs_user` block as defined below.

* `cifs_user` - (Optional) A `cifs_user` block as defined below.

---

A `nfs_user` block supports the following:

* `group_id` - (Required) The POSIX group ID used by the bucket when accessing volume data over NFS.

* `user_id` - (Required) The POSIX user ID used by the bucket when accessing volume data over NFS.

---

A `cifs_user` block supports the following:

* `username` - (Required) The CIFS username used by the bucket when accessing volume data over SMB.

---

A `server` block supports the following:

* `fqdn` - (Optional) The DNS name that resolves to the bucket endpoint IP address.

* `certificate_pem` - (Optional, Sensitive) Base64-encoded PEM blob containing the server certificate and the private key. Used when the certificate is supplied directly instead of via Key Vault. Mutually exclusive with `key_vault`.

* `on_certificate_conflict_action` - (Optional) Behaviour when an existing certificate already matches during a certificate rotation. Possible values are `Update` and `Fail`. Only used during update.

---

A `key_vault` block supports the following:

* `certificate_key_vault_uri` - (Required) The URI of the Azure Key Vault that stores the bucket server certificate.

* `certificate_name` - (Required) The name of the certificate object inside the Key Vault.

* `credentials_key_vault_uri` - (Required) The URI of the Azure Key Vault used to store the generated bucket access and secret keys. May be the same vault as `certificate_key_vault_uri` but the documentation recommends using two separate vaults.

* `credentials_secret_name` - (Required) The name of the secret in `credentials_key_vault_uri` that stores the generated bucket credentials. The Key Vault secret value is a JSON document with `access_key_id` and `secret_access_key` properties.

~> **Note:** The Azure NetApp Files service principal needs the following Key Vault permissions on `certificate_key_vault_uri`: `Get, List, Update, Create, Import, Manage Certificate Authorities, Get Certificate Authorities, List Certificate Authorities, Set Certificate Authorities, Delete Certificate Authorities` for certificates, and `Get, List, Set, Delete` for secrets on `credentials_key_vault_uri`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NetApp Volume Bucket.

* `status` - The credentials status of the bucket. Possible values are `NoCredentialsSet`, `CredentialsExpired` and `Active`.

* `server_ip_address` - The IP address that backs the bucket endpoint.

* `server_certificate_common_name` - The Common Name (CN) of the bucket server certificate.

* `server_certificate_expiry_date` - The expiry date of the bucket server certificate, in RFC3339 format.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the NetApp Volume Bucket.
* `read` - (Defaults to 5 minutes) Used when retrieving the NetApp Volume Bucket.
* `update` - (Defaults to 1 hour) Used when updating the NetApp Volume Bucket.
* `delete` - (Defaults to 1 hour) Used when deleting the NetApp Volume Bucket.

## Import

NetApp Volume Buckets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_netapp_volume_bucket.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/vol1/buckets/bucket1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.NetApp` - 2026-01-01
