---
subcategory: "Oracle"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_oracle_adbs_character_sets"
description: |-
  This data source provides the list of Autonomous Database Character Sets.
---

# Data Source: azurerm_oracle_adbs_character_sets

Gets a list of supported character sets.

## Example Usage

```hcl
data "azurerm_oracle_adbs_character_sets" "example" {
  location = "West Europe"
}

output "example" {
  value = data.azurerm_oracle_adbs_character_sets.example
}
```

## Arguments Reference

The following arguments are supported:

* `location` - (Required) The Azure Region to query for the character sets in.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `character_sets` - A `character_sets` block as defined below.

---

A `character_sets` block exports the following:

* `character_set` - A valid Oracle character set.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Oracle character set.
