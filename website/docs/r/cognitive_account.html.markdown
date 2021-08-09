---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account"
description: |-
  Manages a Cognitive Services Account.
---

# azurerm_cognitive_account

Manages a Cognitive Services Account.

-> **Note:** Version v2.65.0 of the Azure Provider and later will attempt to Purge the Cognitive Account during deletion. This feature can be disabled using the `features` block within the `provider` block, see [the provider documentation on the features block](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs#features) for more information.

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

* `kind` - (Required) Specifies the type of Cognitive Service Account that should be created. Possible values are `Academic`, `AnomalyDetector`, `Bing.Autosuggest`, `Bing.Autosuggest.v7`, `Bing.CustomSearch`, `Bing.Search`, `Bing.Search.v7`, `Bing.Speech`, `Bing.SpellCheck`, `Bing.SpellCheck.v7`, `CognitiveServices`, `ComputerVision`, `ContentModerator`, `CustomSpeech`, `CustomVision.Prediction`, `CustomVision.Training`, `Emotion`, `Face`,`FormRecognizer`, `ImmersiveReader`, `LUIS`, `LUIS.Authoring`, `MetricsAdvisor`, `Personalizer`, `QnAMaker`, `Recommendations`, `SpeakerRecognition`, `Speech`, `SpeechServices`, `SpeechTranslation`, `TextAnalytics`, `TextTranslation` and `WebLM`. Changing this forces a new resource to be created.

* `sku_name` - (Required) Specifies the SKU Name for this Cognitive Service Account. Possible values are `F0`, `F1`, `S`, `S0`, `S1`, `S2`, `S3`, `S4`, `S5`, `S6`, `P0`, `P1`, and `P2`.

* `custom_subdomain_name` - (Optional) The subdomain name used for token-based authentication. Changing this forces a new resource to be created.

* `fqdns` - (Optional) List of FQDNs allowed for the Cognitive Account.

* `identity` - (Optional) An `identity` block is documented below.

* `local_auth_enabled` - (Optional) Whether local authentication methods is enabled for the Cognitive Account. Defaults to `true`.

* `metrics_advisor_aad_client_id` - (Optional) The Azure AD Client ID (Application ID). This attribute is only set when kind is `MetricsAdvisor`. Changing this forces a new resource to be created.

* `metrics_advisor_aad_tenant_id` - (Optional) The Azure AD Tenant ID. This attribute is only set when kind is `MetricsAdvisor`. Changing this forces a new resource to be created.

* `metrics_advisor_super_user_name` - (Optional) The super user of Metrics Advisor. This attribute is only set when kind is `MetricsAdvisor`. Changing this forces a new resource to be created.

* `metrics_advisor_website_name` - (Optional) The website name of Metrics Advisor. This attribute is only set when kind is `MetricsAdvisor`. Changing this forces a new resource to be created.

-> **NOTE:** This URL is mandatory if the `kind` is set to `QnAMaker`.

* `network_acls` - (Optional) A `network_acls` block as defined below.

* `outbound_network_access_restrited` - (Optional) Whether outbound network access is restricted for the Cognitive Account. Defaults to `false`.

* `public_network_access_enabled` - (Optional) Whether public network access is allowed for the Cognitive Account. Defaults to `true`.

* `qna_runtime_endpoint` - (Optional) A URL to link a QnAMaker cognitive account to a QnA runtime.

* `storage` - (Optional) An `identity` block is documented below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `network_acls` block supports the following:

* `default_action` - (Required) The Default Action to use when no rules match from `ip_rules` / `virtual_network_subnet_ids`. Possible values are `Allow` and `Deny`.

* `ip_rules` - (Optional) One or more IP Addresses, or CIDR Blocks which should be able to access the Cognitive Account.

* `virtual_network_rules` - (Optional) A `virtual_network_rules` block as defined below.

A `virtual_network_rules` block supports the following:

* `subnet_id` - (Required) The ID of the subnet which should be able to access this Cognitive Account.

* `ignore_missing_vnet_service_endpoint` - (Optional) Whether ignore missing vnet service endpoint or not. Default to `false`.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on the Cognitive Account. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of IDs for User Assigned Managed Identity resources to be assigned.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `storage` block supports the following:

* `storage_account_id` - (Required) Full resource id of a Microsoft.Storage resource.

* `identity_client_id` - (Optional) The client ID of the managed identity associated with the storage resource.

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
