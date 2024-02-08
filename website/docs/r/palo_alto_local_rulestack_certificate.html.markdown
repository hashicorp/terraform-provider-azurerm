---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_local_rulestack_certificate"
description: |-
  Manages a Palo Alto Networks Rulestack Certificate.
---

# azurerm_palo_alto_local_rulestack_certificate

Manages a Palo Alto Networks Rulestack Certificate.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

resource "azurerm_palo_alto_local_rulestack" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_palo_alto_local_rulestack_certificate" "example" {
  name         = "example"
  rulestack_id = azurerm_palo_alto_local_rulestack.example.id
  self_signed  = true
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Palo Alto Networks Rulestack Certificate.

* `rulestack_id` - (Required) The ID of the TODO. Changing this forces a new Palo Alto Networks Rulestack Certificate to be created.

---

* `key_vault_certificate_id` - (Optional) The `versionles_id` of the Key Vault Certificate to use. Changing this forces a new Palo Alto Networks Rulestack Certificate to be created.

* `self_signed` - (Optional) Should a Self Signed Certificate be used. Defaults to `false`. Changing this forces a new Palo Alto Networks Rulestack Certificate to be created.

~> **Note:** One and only one of `self_signed` or `key_vault_certificate_id` must be specified.

* `audit_comment` - (Optional) The comment for Audit purposes.

* `description` - (Optional) The description for the Certificate.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Networks Rulestack Certificate.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Palo Alto Networks Rulestack Certificate.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Networks Rulestack Certificate.
* `update` - (Defaults to 30 minutes) Used when updating the Palo Alto Networks Rulestack Certificate.
* `delete` - (Defaults to 30 minutes) Used when deleting the Palo Alto Networks Rulestack Certificate.

## Import

Palo Alto Networks Rulestack Certificates can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_local_rulestack_certificate.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/myLocalRulestack/certificates/myCertificate
```
