---
subcategory: "Key Vault"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_key_vault_managed_hardware_security_module_key_rotation_policy"
description: |-
  Manages a Managed HSM Key rotation policy.
---

# azurerm_key_vault_managed_hardware_security_module_key_rotation_policy

Manages a Managed HSM Key rotation policy.

## Example Usage

```hcl
resource "azurerm_key_vault_managed_hardware_security_module_key" "example" {
  name           = "example-key"
  managed_hsm_id = azurerm_key_vault_managed_hardware_security_module.example.id
  key_type       = "EC-HSM"
  curve          = "P-521"
  key_opts       = ["sign"]
}

resource "azurerm_key_vault_managed_hardware_security_module_key_rotation_policy" "example" {
  managed_hsm_key_id = azurerm_key_vault_managed_hardware_security_module_key.example.id
  expire_after       = "P60D"
  time_before_expiry = "P30D"
}
```

## Arguments Reference

The following arguments are supported:

* `managed_hsm_key_id` - (Required) The ID of the Managed HSM Key. Changing this forces a new Managed HSM Key rotation policy to be created.

* `expire_after` - (Required) Specify the expiration duration on a newly rotated key as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). The minimum duration is `P28D`.

* `time_after_creation` - (Optional) Rotate automatically at a duration after key creation as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). Exactly one of `time_after_creation` or `time_before_expiry` should be specified.

* `time_before_expiry` - (Optional) Rotate automatically at a duration before key expiry as an [ISO 8601 duration](https://en.wikipedia.org/wiki/ISO_8601#Durations). Exactly one of `time_after_creation` or `time_before_expiry` should be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed HSM Key Rotation policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed HSM Key rotation policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed HSM Key rotation policy.
* `update` - (Defaults to 30 minutes) Used when updating the Managed HSM Key rotation policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed HSM Key rotation policy.

## Import

Managed HSM Key rotation policy can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_key_vault_managed_hardware_security_module_key_rotation_policy.example https://example-hsm.managedhsm.azure.net/keys/example
```
