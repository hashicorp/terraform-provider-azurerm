---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_security_perimeter_access_rule"
description: |-
  Manages a Network Security Perimeter Access Rule.
---

# azurerm_network_security_perimeter_access_rule

Manages a Network Security Perimeter Access Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_network_security_perimeter" "example" {
  name = "example"
  resource_group_name = azurerm_resource_group.example.name
  location = "West Europe"
}

resource "azurerm_network_security_perimeter_profile" "example" {
  name = "example"
  perimeter_id = azurerm_network_security_perimeter.example.id
}

resource "azurerm_network_security_perimeter_access_rule" "example" {
  name = "example"
  profile_id = azurerm_network_security_perimeter_profile.example.id
  direction = "Inbound"

  address_prefixes = [
    "8.8.8.8/32"
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `direction` - (Required) The direction of the rule. Possible values are `Inbound` and `Outbound`. Changing this forces a new Network Security Perimeter Access Rule to be created.

* `name` - (Required) The name which should be used for this Network Security Perimeter Access Rule. Changing this forces a new Network Security Perimeter Access Rule to be created.

* `profile_id` - (Required) The ID of the Network Security Perimeter Profile within which this Access Rule is created. Changing this forces a new Network Security Perimeter Access Rule to be created.

---

* `address_prefixes` - (Optional) Specifies a list of CIDRs. Can only be specified when direction is set to `Inbound`. Conflicts with `fqdns` and `subscription_ids`. 

* `fqdns` - (Optional) Specifies a list of fully qualified domain names. Can only be specified when direction is set to `Outbound`. Conflicts with `address_prefixes` and `subscription_ids`. 

* `subscription_ids` - (Optional) Specifies a list of subscription IDs this rule applies to. Can only be specified when direction is set to `Inbound`. Conflicts with `address_prefixes` and `fqdns`. 


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Network Security Perimeter Access Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Security Perimeter Access Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Security Perimeter Access Rule.
* `update` - (Defaults to 30 minutes) Used when updating the Network Security Perimeter Access Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Security Perimeter Access Rule.

## Import

Network Security Perimeter Access Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_security_perimeter_access_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/networkSecurityPerimeters/example-nsp/profiles/defaultProfile/accessRules/example-accessrule
```