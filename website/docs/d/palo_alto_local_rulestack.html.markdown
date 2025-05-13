---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_palo_alto_local_rulestack"
description: |-
  Gets information about an existing Palo Alto Networks Rulestack.
---

# Data Source: azurerm_palo_alto_local_rulestack

Use this data source to access information about an existing Palo Alto Networks Rulestack.

## Example Usage

```hcl
data "azurerm_palo_alto_local_rulestack" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_palo_alto_local_rulestack.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Palo Alto Networks Rulestack.

* `resource_group_name` - (Required) The name of the Resource Group where the Palo Alto Networks Rulestack exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Networks Rulestack.

* `anti_spyware_profile` - The Anti-Spyware setting used by the Palo Alto Networks Rulestack.

* `anti_virus_profile` - The Anti-Virus setting used by the Palo Alto Networks Rulestack.

* `description` - The description of the Palo Alto Networks Rulestack.

* `dns_subscription` - The DNS Subscription setting used by the Palo Alto Networks Rulestack.

* `file_blocking_profile` - The File Blocking Profile used by the Palo Alto Networks Rulestack.

* `location` - The Azure Region where the Palo Alto Networks Rulestack exists.

* `outbound_trust_certificate` - The trusted egress decryption profile data for the Palo Alto Networks Rulestack.

* `outbound_untrust_certificate` - The untrusted egress decryption profile data for the Palo Alto Networks Rulestack.

* `url_filtering_profile` - The URL Filtering Profile used by the Palo Alto Networks Rulestack.

* `vulnerability_profile` - The Vulnerability Profile used by the Palo Alto Networks Rulestack.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Networks Rulestack.
