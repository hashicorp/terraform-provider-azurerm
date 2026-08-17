---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume_bucket_with_server"
description: |-
  Manages the first NetApp Files Volume Bucket on a volume, including its server configuration (Object REST API).
---

# azurerm_netapp_volume_bucket_with_server

Manages a NetApp Files Volume Bucket including its bucket server configuration. Buckets expose the contents of an Azure NetApp Files volume (or a sub-path within it) as an S3-compatible object endpoint via the Azure NetApp Files Object REST API.

Use this resource to create the **first** bucket on a volume. The first bucket establishes the bucket server (FQDN and certificate) that is shared by every bucket on the volume. Create any **subsequent** buckets with the server-less [`azurerm_netapp_volume_bucket`](netapp_volume_bucket.html.markdown) resource, which reuses the server configuration established here.

~> **Note:** Declaring a `server` block on more than one bucket of the same volume overwrites the shared server configuration. Only the first bucket should manage the server, via this resource.

~> **Note:** The Object REST API feature is in preview and must be registered on the subscription via `Microsoft.NetApp / ANFEnableObjectRESTAPI` before buckets can be created. See [Configure access to the Azure NetApp Files Object REST API](https://learn.microsoft.com/en-us/azure/azure-netapp-files/object-rest-api-access-configure) for the registration command and Key Vault prerequisites.

~> **Note:** Buckets are supported on cool-access and large NetApp volumes. Buckets are not supported on cache volumes. Deleting the parent volume cascade-deletes its buckets.

## Example Usage (inline certificate)

This example generates a self-signed server certificate via the [`hashicorp/tls`](https://registry.terraform.io/providers/hashicorp/tls/latest/docs) provider and passes it directly to the bucket through `server.certificate_pem`. Bucket credentials are minted separately by the [`azurerm_netapp_volume_bucket_credentials`](../actions/netapp_volume_bucket_credentials.html.markdown) action and require the bucket to be configured with a `key_vault` block (see the Key Vault example below).

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

resource "azurerm_netapp_volume_bucket_with_server" "example" {
  name      = "example-bucket"
  volume_id = azurerm_netapp_volume.example.id

  file_system_nfs_user {
    group_id = 1000
    user_id  = 1000
  }

  server {
    fqdn            = "example-bucket.example.internal"
    certificate_pem = base64encode("${tls_self_signed_cert.bucket.cert_pem}${tls_private_key.bucket.private_key_pem}")
  }
}
```

## Example Usage (Azure Key Vault)

This example sources the bucket server certificate from Azure Key Vault and persists the generated credentials to a second Azure Key Vault. The NetApp account uses a system-assigned managed identity to access both vaults.

```hcl
provider "azurerm" {
  features {
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = true
    }
  }
}

data "azurerm_client_config" "current" {}

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

  identity {
    type = "SystemAssigned"
  }
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

resource "azurerm_key_vault" "certificate" {
  name                       = "example-cert-kv"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  rbac_authorization_enabled = false
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault" "credentials" {
  name                       = "example-creds-kv"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  rbac_authorization_enabled = false
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  sku_name                   = "standard"
  soft_delete_retention_days = 7
}

resource "azurerm_key_vault_access_policy" "deployer_certificate" {
  key_vault_id = azurerm_key_vault.certificate.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = data.azurerm_client_config.current.object_id

  certificate_permissions = ["Get", "List", "Create", "Import", "Update", "Delete", "Purge", "Recover"]
  secret_permissions      = ["Get", "List", "Set", "Delete", "Purge", "Recover"]
}

resource "azurerm_key_vault_access_policy" "anf_certificate" {
  key_vault_id = azurerm_key_vault.certificate.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_netapp_account.example.identity[0].principal_id

  certificate_permissions = ["Get", "List", "Update", "Create", "Import", "ManageContacts", "GetIssuers", "ListIssuers", "SetIssuers", "DeleteIssuers"]
  secret_permissions      = ["Get", "List", "Set", "Delete"]
}

resource "azurerm_key_vault_access_policy" "anf_credentials" {
  key_vault_id = azurerm_key_vault.credentials.id
  tenant_id    = data.azurerm_client_config.current.tenant_id
  object_id    = azurerm_netapp_account.example.identity[0].principal_id

  secret_permissions = ["Get", "List", "Set", "Delete"]
}

resource "azurerm_key_vault_certificate" "bucket" {
  name         = "example-bucket-cert"
  key_vault_id = azurerm_key_vault.certificate.id

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = false
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      key_usage          = ["digitalSignature", "keyEncipherment"]
      extended_key_usage = ["1.3.6.1.5.5.7.3.1"]
      subject            = "CN=example-bucket.example.internal"

      subject_alternative_names {
        dns_names = ["example-bucket.example.internal"]
      }

      validity_in_months = 12
    }
  }

  depends_on = [
    azurerm_key_vault_access_policy.deployer_certificate,
  ]
}

resource "azurerm_netapp_volume_bucket_with_server" "example" {
  name      = "example-bucket"
  volume_id = azurerm_netapp_volume.example.id

  file_system_nfs_user {
    group_id = 1000
    user_id  = 1000
  }

  server {
    fqdn = "example-bucket.example.internal"
  }

  key_vault {
    certificate_key_vault_uri = azurerm_key_vault.certificate.vault_uri
    certificate_name          = azurerm_key_vault_certificate.bucket.name
    credentials_key_vault_uri = azurerm_key_vault.credentials.vault_uri
    credentials_secret_name   = "example-bucket-creds"
  }

  depends_on = [
    azurerm_key_vault_access_policy.anf_certificate,
    azurerm_key_vault_access_policy.anf_credentials,
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The S3-compatible name of the bucket. Must be 3-63 characters long, DNS-compliant (lowercase letters, digits, hyphens or periods), must start and end with a letter or number and must not look like an IPv4 address. Changing this forces a new resource to be created.

* `volume_id` - (Required) The ARM ID of the parent NetApp Volume the bucket attaches to. Changing this forces a new resource to be created.

* `file_system_nfs_user` - (Optional) A `file_system_nfs_user` block as defined below. Exactly one of `file_system_nfs_user` or `file_system_cifs_username` must be specified.

* `file_system_cifs_username` - (Optional) The CIFS username used by the bucket when accessing volume data over SMB. Exactly one of `file_system_nfs_user` or `file_system_cifs_username` must be specified.

* `server` - (Required) A `server` block as defined below. Used to provide the bucket server FQDN and a directly uploaded PEM certificate. The certificate source (`server.0.certificate_pem`) is mutually exclusive with `key_vault`.

* `path` - (Optional) The volume sub-path mounted inside the bucket. Defaults to `/`. Changing this forces a new resource to be created.

* `permissions` - (Optional) The bucket permission level. Possible values are `ReadOnly` and `ReadWrite`. Defaults to `ReadOnly`.

* `key_vault` - (Optional) A `key_vault` block as defined below. Used to source the server certificate and to store generated credentials in Azure Key Vault. Mutually exclusive with `server.0.certificate_pem`.

---

A `file_system_nfs_user` block supports the following:

* `group_id` - (Required) The POSIX group ID used by the bucket when accessing volume data over NFS.

* `user_id` - (Required) The POSIX user ID used by the bucket when accessing volume data over NFS.

---

A `server` block supports the following:

* `fqdn` - (Required) The DNS name that resolves to the bucket endpoint IP address.

* `certificate_pem` - (Optional, Sensitive) Base64-encoded PEM blob containing the server certificate concatenated with the private key. Used when the certificate is supplied directly instead of via Key Vault. Mutually exclusive with `key_vault`.

* `on_certificate_conflict_action` - (Optional) Behaviour when an existing certificate already matches during a certificate rotation. Possible values are `Update` and `Fail`. Defaults to `Fail`.

---

A `key_vault` block supports the following:

* `certificate_key_vault_uri` - (Required) The URI of the Azure Key Vault that stores the bucket server certificate.

* `certificate_name` - (Required) The name of the certificate object inside the Key Vault.

* `credentials_key_vault_uri` - (Required) The URI of the Azure Key Vault used to store the generated bucket access and secret keys. May be the same vault as `certificate_key_vault_uri` but the documentation recommends using two separate vaults.

* `credentials_secret_name` - (Required) The name of the secret in `credentials_key_vault_uri` that stores the generated bucket credentials. The Key Vault secret value is a JSON document with `access_key_id` and `secret_access_key` properties.

~> **Note:** When `key_vault` is used, the parent NetApp account must have a system-assigned managed identity (`identity { type = "SystemAssigned" }` on `azurerm_netapp_account`). That identity is the principal that needs Key Vault access. Grant it `Get, List, Update, Create, Import, ManageContacts, GetIssuers, ListIssuers, SetIssuers, DeleteIssuers` certificate permissions on `certificate_key_vault_uri` and `Get, List, Set, Delete` secret permissions on `credentials_key_vault_uri`.

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
terraform import azurerm_netapp_volume_bucket_with_server.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.NetApp/netAppAccounts/account1/capacityPools/pool1/volumes/vol1/buckets/bucket1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.NetApp` - 2026-01-01
