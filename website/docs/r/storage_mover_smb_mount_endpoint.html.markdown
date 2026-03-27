---
subcategory: "Storage Mover"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_mover_smb_mount_endpoint"
description: |-
  Manages a Storage Mover SMB Mount Endpoint.
---

# azurerm_storage_mover_smb_mount_endpoint

Manages a Storage Mover SMB Mount Endpoint for migrating from on-premises Windows file servers to Azure.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
}

resource "azurerm_storage_mover_smb_mount_endpoint" "example" {
  name             = "example-smbme"
  storage_mover_id = azurerm_storage_mover.example.id
  host             = "server.contoso.com"
  share_name       = "data"
  description      = "Example SMB Mount Endpoint"
}
```

## Example Usage with Credentials

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_mover" "example" {
  name                = "example-ssm"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
}

resource "azurerm_storage_mover_smb_mount_endpoint" "example" {
  name                         = "example-smbme"
  storage_mover_id             = azurerm_storage_mover.example.id
  host                         = "server.contoso.com"
  share_name                   = "data"
  username_key_vault_secret_id = azurerm_key_vault_secret.username.versionless_id
  password_key_vault_secret_id = azurerm_key_vault_secret.password.versionless_id
  description                  = "Example SMB Mount Endpoint with credentials"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Storage Mover SMB Mount Endpoint. Changing this forces a new resource to be created.

* `storage_mover_id` - (Required) Specifies the ID of the Storage Mover for this SMB Mount Endpoint. Changing this forces a new resource to be created.

* `host` - (Required) Specifies the host name or IP address of the SMB server. Must be a valid IPv4 address or hostname/FQDN (letters, numbers, dots, hyphens only). Changing this forces a new resource to be created.

* `share_name` - (Required) Specifies the name of the SMB share. Must be 1-80 characters and cannot contain: `\` `/` `[` `]` `:` `<` `>` `+` `=` `;` `,` `*` `?` or control characters. Changing this forces a new resource to be created.

* `username_key_vault_secret_id` - (Optional) Specifies the Azure Key Vault secret ID for the username to use for authentication. Must be specified together with `password_key_vault_secret_id`. You can use `azurerm_key_vault_secret.<name>.id` or `azurerm_key_vault_secret.<name>.versionless_id` (recommended to avoid pinning to a specific secret version).

* `password_key_vault_secret_id` - (Optional) Specifies the Azure Key Vault secret ID for the password to use for authentication. Must be specified together with `username_key_vault_secret_id`. You can use `azurerm_key_vault_secret.<name>.id` or `azurerm_key_vault_secret.<name>.versionless_id` (recommended).

* `description` - (Optional) Specifies a description for the Storage Mover SMB Mount Endpoint. Maximum length 1024 characters.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Mover SMB Mount Endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Mover SMB Mount Endpoint.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Mover SMB Mount Endpoint.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Mover SMB Mount Endpoint.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Mover SMB Mount Endpoint.

## Import

Storage Mover SMB Mount Endpoint can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_mover_smb_mount_endpoint.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.StorageMover/storageMovers/storageMover1/endpoints/endpoint1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.StorageMover` - 2025-07-01
