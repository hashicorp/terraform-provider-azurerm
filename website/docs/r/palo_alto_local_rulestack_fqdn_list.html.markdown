---
subcategory: "Palo Alto"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_palo_alto_local_rulestack_fqdn_list"
description: |-
  Manages a Palo Alto Local Rulestack FQDN List.
---

# azurerm_palo_alto_local_rulestack_fqdn_list

Manages a Palo Alto Local Rulestack FQDN List.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

resource "azurerm_palo_alto_local_rulestack" "example" {
  name                = "example"
  resource_group_name = azurerm_resrouce_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_palo_alto_local_rulestack_fqdn_list" "example" {
  name         = "example"
  rulestack_id = azurerm_palo_alto_local_rulestack.example.id

  fully_qualified_domain_names = ["contoso.com"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Palo Alto Local Rulestack FQDN List.

* `rulestack_id` - (Required) The ID of the TODO. Changing this forces a new Palo Alto Local Rulestack FQDN List to be created.

* `fully_qualified_domain_names` - (Required) Specifies a list of Fully Qualified Domain Names.

---

* `audit_comment` - (Optional) The comment for Audit purposes.

* `description` - (Optional) The description for the FQDN List.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Palo Alto Local Rulestack FQDN List.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Palo Alto Local Rulestack FQDN List.
* `read` - (Defaults to 5 minutes) Used when retrieving the Palo Alto Local Rulestack FQDN List.
* `update` - (Defaults to 30 minutes) Used when updating the Palo Alto Local Rulestack FQDN List.
* `delete` - (Defaults to 30 minutes) Used when deleting the Palo Alto Local Rulestack FQDN List.

## Import

Palo Alto Local Rulestack FQDN Lists can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_palo_alto_local_rulestack_fqdn_list.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/PaloAltoNetworks.Cloudngfw/localRulestacks/myLocalRulestack/fqdnLists/myFQDNList1
```
