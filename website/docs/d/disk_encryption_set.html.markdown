---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_disk_encryption_set"
description: |-
  Gets information about an existing Disk Encryption Set
---

# Data Source: azurerm_disk_encryption_set

Use this data source to access information about an existing Disk Encryption Set.

## Example Usage

```hcl
data "azurerm_disk_encryption_set" "existing" {
  name                = "example-des"
  resource_group_name = "example-resources"
}

output "id" {
  value = data.azurerm_disk_encryption_set.existing.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the existing Disk Encryption Set.

* `resource_group_name` - The name of the Resource Group where the Disk Encryption Set exists.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Disk Encryption Set.

* `location` - The location where the Disk Encryption Set exists.

* `auto_key_rotation_enabled` - Is the Azure Disk Encryption Set Key automatically rotated to latest version?

* `key_vault_key_url` - The URL for the Key Vault Key or Key Vault Secret that is currently being used by the service.

* `managed_hsm_key_id` - Key ID of a key in a managed HSM.

~> **Note:** Only one of `key_vault_key_url` and `managed_hsm_key_id` will be set, depending on where the encryption key is stored.

* `identity` - An `identity` block as defined below.

* `tags` - A mapping of tags assigned to the Disk Encryption Set.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Disk Encryption Set.

* `identity_ids` - A list of User Assigned Managed Identity IDs assigned to this Disk Encryption Set.

* `principal_id` - The (Client) ID of the Service Principal.

* `tenant_id` - The ID of the Tenant the Service Principal is assigned in.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Disk Encryption Set.
