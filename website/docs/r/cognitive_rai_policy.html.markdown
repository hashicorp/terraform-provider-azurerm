---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_rai_policy"
description: |-
  Manages a Cognitive Rai Policy Example.
---

# azurerm_cognitive_rai_policy

Manages a Cognitive Rai Policy Example.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "Brazil South"
}

resource "azurerm_cognitive_account" "example" {
  name                = "example-ca"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "OpenAI"
  sku_name            = "S0"
}

resource "azurerm_cognitive_deployment" "example" {
  name                 = "example-cd"
  cognitive_account_id = azurerm_cognitive_account.example.id
  model {
    format  = "OpenAI"
    name    = "text-curie-001"
    version = "1"
  }

  sku {
    name = "Standard"
  }
  rai_policy_name      = azurerm_cognitive_rai_policy.example.name
}

resource "azurerm_cognitive_rai_policy" "example" {
  name                 = "example-rp"
  cognitive_account_id = "azurerm_cognitive_account.example.id"
  mode                 = "Default"

  content_filters = [
    {
      name               = "Hate"
      blocking           = true
      enabled            = true
      severity_threshold = "High"
      source             = "Prompt"
    },
    {
      name               = "Hate"
      blocking           = true
      enabled            = true
      severity_threshold = "High"
      source             = "Completion"
    },
  ]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Cognitive Rai Policy. Changing this forces a new Cognitive Rai Policy Example to be created.
  
* `cognitive_account_id` - (Required) The ID of the Cognitive Services Account. Changing this forces a new resource to be created.

* `mode` - (Required) Rai Policy mode option. Possible values are `Default`, `Asynchronous_filter`,	`Blocking`, and `Deferred`.

---

* `type` - (Optional) Rai Policy type option. Possible values are `SystemManaged` and `UserManaged`.

* `base_policy_name` - (Optional) Rai Policy Base Policy Name option. Defaults to `Microsoft.Default`.

* `content_filters` - (Optional) One or more `content_filters` blocks as defined below.

* `custom_blocklists` - (Optional) A `custom_blocklists` block as defined below.

---

A `content_filters` block supports the following:

* `name` - (Required) The Name which should be used for this Content Filter. Possible values are `Hate`, `Sexual`, `Selfharm`, `Violence`, `Jailbreak`, `Protected Material Text`, `Protected Material Code`, `Indirect Attack` and `Profanity`.

* `blocking` - (Required) Should the Content Filter be blocked?

* `enabled` - (Required) Should the Content Filter be enabled?

* `source` - (Required) The Source which should be used for this Content Filter. Possible values are `Prompt` and `Completion`.

* `severity_threshold` - (Optional) The Severity Threshold which should be used for this Content Filter. Possible values are `High`, `Medium` and `Low`. Defaults to `Medium`.

---

A `custom_blocklists` block supports the following:

* `blocklist_name` - (Required) The Custom Blocklist's Name that must already be created.

* `blocking` - (Required) Should the Custom Blocklist be enabled?.

* `source` - (Required) The Source which should be used for this Custom Blocklist. Possible values are `Prompt` and `Completion`..

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Cognitive Rai Policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Rai Policy Example.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Rai Policy Example.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Rai Policy Example.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Rai Policy Example.

## Import

Cognitive Rai Policy Examples can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_rai_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.CognitiveServices/accounts/account1/raipolicies/raipolicy1
```