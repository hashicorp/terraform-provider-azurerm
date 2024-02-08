---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_local_rulestack_outbound_trust_certificate_association"
description: |-
  Manages a Palo Alto Networks Rulestack Outbound Trust Certificate Association.
---

# azurerm_palo_alto_local_rulestack_outbound_trust_certificate_association

Manages a Palo Alto Networks Rulestack Outbound Trust Certificate Association.

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

resource "azurerm_palo_alto_local_rulestack_outbound_trust_certificate_association" "example" {
  certificate_id = azurerm_palo_alto_local_rulestack_certificate.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `certificate_id` - (Required) The ID of the Certificate to use as the Outbound Trust Certificate. Changing this forces a new Palo Alto Networks Rulestack Outbound Trust Certificate Association to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Networks Rulestack Outbound Trust Certificate Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Palo Alto Networks Rulestack Outbound Trust Certificate Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Networks Rulestack Outbound Trust Certificate Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Palo Alto Networks Rulestack Outbound Trust Certificate Association.
