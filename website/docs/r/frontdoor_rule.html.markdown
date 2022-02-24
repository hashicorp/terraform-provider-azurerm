---
subcategory: "Cdn"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_frontdoor_rule"
description: |-
  Manages a Frontdoor Rule.
---

# azurerm_frontdoor_rule

Manages a Frontdoor Rule.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-cdn"
  location = "West Europe"
}

resource "azurerm_frontdoor_profile" "test" {
  name                = "example-profile"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_frontdoor_rule_set" "test" {
  name                 = "exampleruleset"
  frontdoor_profile_id = azurerm_frontdoor_profile.test.id
}

resource "azurerm_frontdoor_rule" "test" {
  name                  = "example-rule"
  frontdoor_rule_set_id = azurerm_frontdoor_rule_set.test.id
  order                 = 1

  actions    = ["CacheExpiration", "UrlRedirect", "OriginGroupOverride"]
  conditions = ["HostName", "IsDevice", "PostArgs", "RequestMethod"]

  match_processing_behavior = "Continue"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Frontdoor Rule. Changing this forces a new Frontdoor Rule to be created.

* `frontdoor_rule_set_id` - (Required) The ID of the Frontdoor Rule. Changing this forces a new Frontdoor Rule to be created.

* `order` - (Required) The order in which the rules will be applied for the Frontdoor Endpoint. Possible values include `0`, `1`, `2`, `3` etc. A Frontdoor Rule with a lesser order value will be applied before a rule with a greater order value. 

~> **NOTE:** If the Frontdoor Rule has an order value of `0` they do not require any conditions or actions and they will always be applied.

* `actions` - (Required) A list of upto 5 actions for this Frontdoor Rule. Possible values include `CacheExpiration`, `CacheKeyQueryString`, `ModifyRequestHeader`, `ModifyResponseHeader`, `OriginGroupOverride`, `RouteConfigurationOverride`, `UrlRedirect`, `UrlRewrite` or `UrlSigning`.

* `conditions` - (Optional) A list of upto 10 conditions for the Frontdoor Rule. Possible values are `ClientPort`, `Cookies`, `HostName`, `HttpVersion`, `IsDevice`, `PostArgs`, `QueryString`, `RemoteAddress`, `RequestBody`, `RequestHeader`, `RequestMethod`, `RequestScheme`, `RequestUri`, `ServerPort`, `SocketAddr`, `SslProtocol`, `UrlFileExtension`, `UrlFileName` or `UrlPath`.

* `match_processing_behavior` - (Optional) If this rule is a match should the rules engine continue processing the remaining rules or stop? Possible values are `Continue` and `Stop`. Defaults to `Continue`.

---

An `actions` block supports the following:

* `name` - (Required) The name of the action for the delivery rule.

---

A `conditions` block supports the following:

* `name` - (Required) The name of the condition for the delivery rule.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Frontdoor Rule.

* `rule_set_name` - The name of the Frontdoor Rule Set containing this Frontdoor Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the cdn Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the cdn Rule.
* `update` - (Defaults to 30 minutes) Used when updating the cdn Rule.
* `delete` - (Defaults to 30 minutes) Used when deleting the cdn Rule.

## Import

cdn Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_frontdoor_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Cdn/profiles/profile1/ruleSets/ruleSet1/rules/rule1
```
