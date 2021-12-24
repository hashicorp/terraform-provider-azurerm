---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cdn_frontdoor_rule_set_rule"
description: |-
  Manages an Azure Front Door (Standard/Premium) instance. (currently in public preview)
---

# azurerm_cdn_frontdoor_rule_set_rule

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cdn_frontdoor_profile" "example" {
  name                = "afdpremv2"
  location            = "global"
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "Premium_AzureFrontDoor"
}

resource "azurerm_cdn_frontdoor_rule_set" "example" {
  name       = "ruleset1"
  profile_id = azurerm_cdn_frontdoor_profile.example.id
}

resource "azurerm_cdn_frontdoor_rule_set_rule" "example" {
  name        = "rule1"
  rule_set_id = azurerm_cdn_frontdoor_rule_set.example.id

  order = 1

  action {
    type            = "UrlRedirect"
    custom_hostname = "redirected.to.target"
    custom_path     = "/path/xyz/"
    redirect_type   = "Moved"

  }

  condition {

  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the rule.

* `rule_set_id` - (Required) ID of the parent rule set resource.

* `order` - (Required) Order

---

The `action` block supports the following:

* `type` - (Required) Action type. Can be set to `UrlSigning`, `UrlRewrite`, `UrlRedirect`. (list is to be completed) Changing the type attribute forces a re-creation of a rule.

---

The `UrlRedirect` action type supports the following parameters:

* `custom_hostname` - Host to redirect. Leave empty to use the incoming host as the destination host.

* `custom_path` - The full path to redirect. Path cannot be empty and must start with /. Leave empty to use the incoming path as destination path.

* `custom_querystring` - The set of query strings to be placed in the redirect URL. Setting this value would replace any existing query string; leave empty to preserve the incoming query string. Query string must be in `<key>=<value>` format. `?` and `&` will be added automatically so do not include them.

* `custom_fragment` - Fragment to add to the redirect URL. Fragment is the part of the URL that comes after #. Do not include the #.

* `redirect_type` - (Required) Can be set to `Found`, `Moved`, `PermanentRedirect` or `TemporaryRedirect`.

* `destination_protocol` - Destination protocol. Can be set to `Http`, `Https` or `MatchRequest`.

---

The `UrlSigning` action type supports the following parameters:

* `param_name` 

* `param_indicator` can be set to `Expires`, `KeyId` or `Signature`.

---

The `UrlRewrite` action type supports the following parameters:

* `destination` - Destination pattern, needs to start with `/`.

* `source_pattern` - Source pattern, needs to start with `/`.

* `preserve_unmatched_path` - Can be set to `true` or `false`.

---

The `condition` block supports the following:
