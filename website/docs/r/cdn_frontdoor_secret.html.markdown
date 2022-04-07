---
subcategory: "CDN"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_secret"
description: |-
  Manages a Frontdoor Secret.
---

# azurerm_cdn_frontdoor_secret

Manages a Frontdoor Secret.

## Example Usage

```hcl
resource "azurerm_cdn_frontdoor_secret" "example" {
  name                     = "exampleSecret"
  cdn_frontdoor_profile_id = cdn_frontdoor_profile.example.id

  secret_parameters {
    customer_certificate {
      secret_source_id = azurerm_key_vault_secret.example.id
      use_latest       = true
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Secret. Changing this forces a new Frontdoor Secret to be created.

* `cdn_frontdoor_profile_id` - (Required) The ID of the Frontdoor Profile. Changing this forces a new Frontdoor Secret to be created.

* `secret_parameters` - (Required) A `secret_parameters` block as defined below. Changing this forces a new Frontdoor Secret to be created.

---

A `secret_parameters` block supports the following:

* `customer_certificate` - (Required) A `customer_certificate` block as defined below. Changing this forces a new Frontdoor Secret to be created.

---

A `customer_certificate` - (Required)  block supports the following:

* `secret_source_id` - (Required) The Resource ID of the Azure Key Vault secret. Expected to be in format of /subscriptions/00000000-0000-0000-0000-000000000000​​​​​​​​​/resourceGroups/resourceGroup1​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​/providers/Microsoft.KeyVault/vaults/vault1​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​/secrets/secret1.​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​​
* `secret_version` - (Optional) The version of the secret to be used.

* `use_latest` - (Optional) Should the latest version for the certificate be used? Defaults to `true`.

* `subject_alternative_names` - (Optional) One or more, up to 100, subject alternative names.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Secret.

* `cdn_frontdoor_profile_name` - The name of the Frontdoor Profile containing this Frontdoor Secret.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Frontdoor Secret.
* `read` - (Defaults to 5 minutes) Used when retrieving the Frontdoor Secret.
* `delete` - (Defaults to 30 minutes) Used when deleting the Frontdoor Secret.

## Import

Frontdoor Secrets can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cdn_frontdoor_secret.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/secrets/secrets1
```
