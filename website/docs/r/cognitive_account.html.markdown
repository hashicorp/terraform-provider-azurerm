---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account"
description: |-
  Manages a Cognitive Services Account.
---

# azurerm_cognitive_account

Manages a Cognitive Services Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                = "example-account"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  kind                = "Face"

  sku_name = "S0"

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

* `kind` - (Required) Specifies the type of Cognitive Service Account that should be created. Possible values are `Academic`, `AnomalyDetector`, `Bing.Autosuggest`, `Bing.Autosuggest.v7`, `Bing.CustomSearch`, `Bing.Search`, `Bing.Search.v7`, `Bing.Speech`, `Bing.SpellCheck`, `Bing.SpellCheck.v7`, `CognitiveServices`, `ComputerVision`, `ContentModerator`, `CustomSpeech`, `CustomVision.Prediction`, `CustomVision.Training`, `Emotion`, `Face`,`FormRecognizer`, `ImmersiveReader`, `LUIS`, `LUIS.Authoring`, `Personalizer`, `QnAMaker`, `Recommendations`, `SpeakerRecognition`, `Speech`, `SpeechServices`, `SpeechTranslation`, `TextAnalytics`, `TextTranslation` and `WebLM`. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for this Cognitive Service Account. Possible values are `F0`, `F1`, `S`, `S0`, `S1`, `S2`, `S3`, `S4`, `S5`, `S6`, `P0`, `P1`, and `P2`.

* `qna_runtime_endpoint` - (Optional) A URL to link a QnAMaker cognitive account to a QnA runtime.

-> **NOTE:** This URL is mandatory if the `kind` is set to `QnAMaker`.

* `network_acls` - (Optional) A `network_acls` block as defined below.

* `custom_subdomain_name` - (Optional) The subdomain name used for token-based authentication. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `network_acls` block supports the following:

* `default_action` - (Required) The Default Action to use when no rules match from `ip_rules` / `virtual_network_subnet_ids`. Possible values are `Allow` and `Deny`.

* `ip_rules` - (Optional) One or more IP Addresses, or CIDR Blocks which should be able to access the Cognitive Account.

* `virtual_network_subnet_ids` - (Optional) One or more Subnet ID's which should be able to access this Cognitive Account.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Cognitive Service Account.

* `endpoint` - The endpoint used to connect to the Cognitive Service Account.

* `primary_access_key` - A primary access key which can be used to connect to the Cognitive Service Account.

* `secondary_access_key` - The secondary access key which can be used to connect to the Cognitive Service Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Service Account.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Service Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Service Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Service Account.

## Import

Cognitive Service Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account.account1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1
```
