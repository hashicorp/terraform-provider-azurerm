---
subcategory: "NetApp"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_netapp_volume_bucket_credentials"
description: |-
  Generates the S3 access key and secret key pair for a NetApp Files Volume Bucket and writes them to the bucket's configured Azure Key Vault.
---

# Action: azurerm_netapp_volume_bucket_credentials

Generates the S3 access key and secret key pair used to authenticate against an Azure NetApp Files Volume Bucket (Object REST API). Azure NetApp Files writes the generated keys as a JSON secret (`{"access_key_id": "...", "secret_access_key": "..."}`) into the credentials Key Vault configured on the parent bucket via its `key_vault` block. Terraform never sees the key material.

Each invocation generates a new key pair and immediately invalidates the previous one, so use this action as a rotation primitive: gate it behind a `terraform_data` lifecycle `action_trigger` and re-trigger when you want to rotate.

~> **Note:** The bucket must be configured with a `key_vault` block (specifying both `credentials_key_vault_uri` and `credentials_secret_name`) for the action to succeed. The NetApp account's system-assigned managed identity must have `Get, List, Set, Delete` secret permissions on the credentials Key Vault.

## Example Usage

```terraform
action "azurerm_netapp_volume_bucket_credentials" "example" {
  config {
    bucket_id            = azurerm_netapp_volume_bucket.example.id
    key_pair_expiry_days = 30
  }
}

resource "terraform_data" "bucket_credentials" {
  input = azurerm_netapp_volume_bucket.example.id

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_netapp_volume_bucket_credentials.example]
    }
  }
}
```

See the [`examples/netapp/volume_bucket_akv`](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/netapp/volume_bucket_akv) example for a complete end-to-end walkthrough including the Key Vaults, managed identity and access policies the bucket requires.

## Argument Reference

This action supports the following arguments:

* `bucket_id` - (Required) The ID of the NetApp Volume Bucket to generate credentials for. The bucket must have a `key_vault` block configured.

* `key_pair_expiry_days` - (Required) The lifetime of the generated key pair, in days. Must be greater than or equal to `1`.

---

* `timeout` - (Optional) Timeout duration for the action to complete. Defaults to `60m`.
