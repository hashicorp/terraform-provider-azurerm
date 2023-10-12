---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_local_rulestack"
description: |-
  Manages a Palo Alto Networks Local Rulestack.
---

# azurerm_palo_alto_local_rulestack

Manages a Palo Alto Networks Rulestack.

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
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Palo Alto Networks Rulestack. Changing this forces a new Palo Alto Networks Rulestack to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Palo Alto Networks Rulestack should exist. Changing this forces a new Palo Alto Networks Rulestack to be created.

* `location` - (Required) The Azure Region where the Palo Alto Networks Rulestack should exist. Changing this forces a new Palo Alto Networks Rulestack to be created.

---

* `anti_spyware_profile` - (Optional) The setting to use for Anti-Spyware.  Possible values include `BestPractice`, and `Custom`.

* `anti_virus_profile` - (Optional) The setting to use for Anti-Virus. Possible values include `BestPractice`, and `Custom`.

* `description` - (Optional) The description for this Local Rulestack.

* `dns_subscription` - (Optional) TThe setting to use for DNS Subscription. Possible values include `BestPractice`, and `Custom`.

* `file_blocking_profile` - (Optional) The setting to use for the File Blocking Profile. Possible values include `BestPractice`, and `Custom`.

* `url_filtering_profile` - (Optional) The setting to use for the URL Filtering Profile. Possible values include `BestPractice`, and `Custom`.

* `vulnerability_profile` - (Optional) The setting to use for the Vulnerability Profile. Possible values include `BestPractice`, and `Custom`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Networks Rulestack.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Palo Alto Networks Rulestack.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Networks Rulestack.
* `update` - (Defaults to 30 minutes) Used when updating the Palo Alto Networks Rulestack.
* `delete` - (Defaults to 30 minutes) Used when deleting the Palo Alto Networks Rulestack.

## Import

Palo Alto Networks Rulestacks can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_local_rulestack.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/myLocalRulestack
```
