---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account"
sidebar_current: "docs-azurerm-resource-cognitive-account"
description: |-
  Manages a Cognitive Services Account.
---

# azurerm_cognitive_account

Manages a Cognitive Services Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "test" {
  name                = "example-account"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  kind                = "Face"

  sku {
    name = "S0"
    tier = "Standard"
  }

  tags = {
    Acceptance = "Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Cognitive Service Account. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which the Cognitive Service Account is created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `kind` - (Required) Specifies the type of Cognitive Service Account that should be created. Possible values are `Academic`, `Bing.Autosuggest`, `Bing.Autosuggest.v7`, `Bing.CustomSearch`, `Bing.Search`, `Bing.Search.v7`, `Bing.Speech`, `Bing.SpellCheck`, `Bing.SpellCheck.v7`, `ComputerVision`, `ContentModerator`, `CustomSpeech`, `CustomVision.Prediction`, `CustomVision.Training`, `Emotion`, `Face`, `LUIS`, `QnAMaker`, `Recommendations`, `SpeakerRecognition`, `Speech`, `SpeechServices`, `SpeechTranslation`, `TextAnalytics`, `TextTranslation` and `WebLM`. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `sku` block supports the following:

* `name` - (Required) Specifies the Name of the Sku. Possible values are `F0`, `S0`, `S1`, `S2`, `S3`, `S4`, `S5`, `S6`, `P0`, `P1` and `P2`.

* `tier` - (Required) Specifies the Tier of the Sku. Possible values include `Free`, `Standard` and `Premium`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cognitive Service Account.

* `endpoint` - The endpoint used to connect to the Cognitive Service Account.

* `primary_access_key` - A primary access key which can be used to connect to the Cognitive Service Account.

* `secondary_access_key` - The secondary access key which can be used to connect to the Cognitive Service Account.

## Import

Cognitive Service Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account.account1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1
```
