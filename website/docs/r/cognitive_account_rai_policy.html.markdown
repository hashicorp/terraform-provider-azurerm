---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_rai_policy"
description: |-
  Manages a Cognitive Services Account RAI Policy.
---

# azurerm_cognitive_account

Manages a Cognitive Services Account RAI Policy.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "East US"
}

resource "azurerm_cognitive_account" "example" {
  name                = "example-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "OpenAI"
  sku_name            = "S0"
}

resource "azurerm_cognitive_account_rai_policy" "example" {
  name                 = "example-rai-policy"
  cognitive_account_id = azurerm_cognitive_account.example.id
  base_policy_name     = "Microsoft.Default"
  content_filter {
    name               = "Hate"
    filter_enabled     = true
    block_enabled      = true
    severity_threshold = "High"
    source             = "Prompt"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Service Account RAI Policy. Changing this forces a new resource to be created.

* `cognitive_account_id` - (Required) The ID of the Cognitive Service Account to which this RAI Policy should be associated. Changing this forces a new resource to be created.

* `base_policy_name` - (Required) The name of the base policy to use for this RAI Policy. Changing this forces a new resource to be created.

* `content_filter` - (Required) A `content_filter` block as defined below.

* `mode` - (Optional) The mode of the RAI Policy. Possible values are `Default`, `Deferred`, `Blocking` or `Asynchronous_filter`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `content_filter` block supports the following:

* `name` - (Required) The name of the content filter.

* `filter_enabled` - (Required) Whether the filter is enabled. Possible values are `true` or `false`.

* `block_enabled` - (Required) Whether the filter should block content. Possible values are `true` or `false`.

* `severity_threshold` - (Required) The severity threshold for the filter. Possible values are `Low`, `Medium` or `High`.

* `source` - (Required) Content source to apply the content filter. Possible values are `Prompt` or `Completion`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Service Account RAI Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Service Account RAI Policy.

* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Service Account RAI Policy.

* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Service Account RAI Policy.

* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Service Account RAI Policy.

## Import

Cognitive Service Account RAI Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_rai_policy.policy1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1/raiPolicies/policy1
```
