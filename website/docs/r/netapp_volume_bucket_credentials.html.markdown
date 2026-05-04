---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume_bucket_credentials"
description: |-
  Manages credentials (access key / secret key) for a NetApp Files Volume Bucket.
---

# azurerm_netapp_volume_bucket_credentials

Generates the access key and secret key pair used by S3 clients to authenticate against an `azurerm_netapp_volume_bucket`.

~> **Note:** Generating new credentials immediately invalidates any existing credentials for the bucket. To rotate the credentials, replace this resource - `key_pair_expiry_days` is `ForceNew`.

~> **Note:** The Azure NetApp Files Object REST API does not expose an explicit revoke operation. Credentials become invalid on their expiry date or when a new credential pair is generated for the same bucket. Destroying this resource is therefore a no-op; the credentials remain valid until expiry.

~> **Note:** For production workloads set `store_in_key_vault = true` so that Azure NetApp Files writes the generated `access_key` / `secret_key` to the Azure Key Vault configured on the parent bucket (`key_vault.0.credentials_key_vault_uri` / `credentials_secret_name`) instead of returning them inline. The default (`false`) puts the bucket access / secret keys into Terraform state.

## Example Usage

This example generates an inline credential pair for an `azurerm_netapp_volume_bucket`. The keys are returned by the Azure NetApp Files API and stored as sensitive attributes on this resource. Replace this with `store_in_key_vault = true` (and a `key_vault` block on the parent bucket) for production deployments.

```hcl
provider "azurerm" {
  features {}
}

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

resource "tls_private_key" "bucket" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "bucket" {
  private_key_pem = tls_private_key.bucket.private_key_pem

  subject {
    common_name = "example-bucket.example.internal"
  }

  dns_names             = ["example-bucket.example.internal"]
  validity_period_hours = 8760

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
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

  server {
    fqdn            = "example-bucket.example.internal"
    certificate_pem = base64encode("${tls_self_signed_cert.bucket.cert_pem}${tls_private_key.bucket.private_key_pem}")
  }
}

resource "azurerm_netapp_volume_bucket_credentials" "example" {
  bucket_id            = azurerm_netapp_volume_bucket.example.id
  key_pair_expiry_days = 30
}
```

## Arguments Reference

The following arguments are supported:

* `bucket_id` - (Required) The ARM ID of the NetApp Volume Bucket the credentials apply to. Changing this forces a new resource to be created.

* `key_pair_expiry_days` - (Required) Number of days the generated key pair is valid for. Must be at least `1`. Changing this forces a new resource to be created (and rotates the credentials).

* `store_in_key_vault` - (Optional) When `true`, the credentials are written to the Key Vault configured on the parent bucket via its `key_vault.0.credentials_key_vault_uri` / `credentials_secret_name`, and `access_key`/`secret_key` are not returned by the API. Defaults to `false`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the NetApp Volume Bucket the credentials are bound to.

* `access_key` - (Sensitive) The generated S3 access key. Only populated when `store_in_key_vault` is `false`.

* `secret_key` - (Sensitive) The generated S3 secret key. Only populated when `store_in_key_vault` is `false`.

* `key_pair_expiry` - The expiry timestamp of the credential pair, in RFC3339 format.

* `status` - The credentials status of the bucket. Possible values are `NoCredentialsSet`, `CredentialsExpired` and `Active`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the bucket credentials.
* `read` - (Defaults to 5 minutes) Used when retrieving the bucket credentials state.
* `delete` - (Defaults to 5 minutes) Used when deleting the bucket credentials from state.

## Import

NetApp Volume Bucket Credentials are bound to a bucket and can be imported using the bucket `resource id`, e.g.

```shell
terraform import azurerm_netapp_volume_bucket_credentials.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/vol1/buckets/bucket1
```

~> **Note:** The actual `access_key` and `secret_key` are returned only at generation time and are not retrievable on import. After importing, regenerate the credentials by tainting the resource if you need to recover them.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.NetApp` - 2026-01-01
