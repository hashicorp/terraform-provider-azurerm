---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account"
description: |-
  Manages a Cognitive Services Account.
---

# azurerm_cognitive_account

Manages a Cognitive Services Account.

-> **Note:** The Cognitive Services Account manages the resource type for various Azure AI resource implementations, including Azure AI Foundry, Azure OpenAI, Azure Speech, Azure Vision and others. Each service shares the same control plane but exposes a different subset of developer APIs. Azure AI Foundry (kind = `AIServices`) provides the superset of capabilities. For more information, please see [Azure AI Foundry architecture](https://learn.microsoft.com/en-us/azure/ai-foundry/concepts/architecture).

-> **Note:** The Azure Provider will attempt to Purge the Cognitive Services Account during deletion. This feature can be disabled using the `features` block within the `provider` block, see [the provider documentation on the features block](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/features-block) for more information.

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

* `kind` - (Required) Specifies the type of Cognitive Service Account that should be created. Possible values are `Academic`, `AIServices`, `AnomalyDetector`, `Bing.Autosuggest`, `Bing.Autosuggest.v7`, `Bing.CustomSearch`, `Bing.Search`, `Bing.Search.v7`, `Bing.Speech`, `Bing.SpellCheck`, `Bing.SpellCheck.v7`, `CognitiveServices`, `ComputerVision`, `ContentModerator`, `ContentSafety`, `CustomSpeech`, `CustomVision.Prediction`, `CustomVision.Training`, `Emotion`, `Face`, `FormRecognizer`, `ImmersiveReader`, `LUIS`, `LUIS.Authoring`, `MetricsAdvisor`, `OpenAI`, `Personalizer`, `QnAMaker`, `Recommendations`, `SpeakerRecognition`, `Speech`, `SpeechServices`, `SpeechTranslation`, `TextAnalytics`, `TextTranslation` and `WebLM`. Changing this forces a new resource to be created.

-> **Note:** New Bing Search resources cannot be created as their APIs are moving from Cognitive Services Platform to new surface area under Microsoft.com. Starting from October 30, 2020, existing instances of Bing Search APIs provisioned via Cognitive Services will be continuously supported for next 3 years or till the end of respective Enterprise Agreement, whichever happens first.

-> **Note:** You must create your first Face, Text Analytics, or Computer Vision resources from the Azure portal to review and acknowledge the terms and conditions. In Azure Portal, the checkbox to accept terms and conditions is only displayed when a US region is selected. More information on [Prerequisites](https://docs.microsoft.com/azure/cognitive-services/cognitive-services-apis-create-account-cli?tabs=windows#prerequisites).

* `sku_name` - (Required) Specifies the SKU Name for this Cognitive Service Account. Possible values are `C2`, `C3`, `C4`, `D3`, `DC0`, `E0`, `F0`, `F1`, `P0`, `P1`, `P2`, `S`, `S0`, `S1`, `S2`, `S3`, `S4`, `S5` and `S6`.

-> **Note:** SKU `DC0` is the commitment tier for Cognitive Services containers running in disconnected environments. You must obtain approval from Microsoft by submitting the [request form](https://aka.ms/csdisconnectedcontainers) first, before you can use this SKU. More information on [Purchase a commitment plan to use containers in disconnected environments](https://learn.microsoft.com/en-us/azure/cognitive-services/containers/disconnected-containers?tabs=stt#purchase-a-commitment-plan-to-use-containers-in-disconnected-environments).

* `custom_subdomain_name` - (Optional) The subdomain name used for Entra ID token-based authentication. This attribute is required when `network_acls` is specified. This attribute is also required when using the OpenAI service with libraries which assume the Azure OpenAI endpoint is a subdomain on `https://openai.azure.com/`, eg. `https://<custom_subdomain_name>.openai.azure.com/`. This can be specified during creation or added later, but once set changing this forces a new resource to be created.

~> **Note:** If you do not specify a `custom_subdomain_name` then you will not be able to attach a Private Endpoint to the resource. Moreover, functionality that requires Entra ID authentication, including Agent service, will not be accessible.

* `dynamic_throttling_enabled` - (Optional) Whether to enable the dynamic throttling for this Cognitive Service Account. This attribute cannot be set when the `kind` is `OpenAI` or `AIServices`.

* `customer_managed_key` - (Optional) A `customer_managed_key` block as documented below.

* `fqdns` - (Optional) List of FQDNs allowed for the Cognitive Account.

* `identity` - (Optional) An `identity` block as defined below.

* `local_auth_enabled` - (Optional) Whether local authentication methods is enabled for the Cognitive Account. Defaults to `true`.

* `metrics_advisor_aad_client_id` - (Optional) The Azure AD Client ID (Application ID). This attribute is only set when kind is `MetricsAdvisor`. Changing this forces a new resource to be created.

* `metrics_advisor_aad_tenant_id` - (Optional) The Azure AD Tenant ID. This attribute is only set when kind is `MetricsAdvisor`. Changing this forces a new resource to be created.

* `metrics_advisor_super_user_name` - (Optional) The super user of Metrics Advisor. This attribute is only set when kind is `MetricsAdvisor`. Changing this forces a new resource to be created.

* `metrics_advisor_website_name` - (Optional) The website name of Metrics Advisor. This attribute is only set when kind is `MetricsAdvisor`. Changing this forces a new resource to be created.

-> **Note:** This URL is mandatory if the `kind` is set to `QnAMaker`.

* `network_acls` - (Optional) A `network_acls` block as defined below. When this property is specified, `custom_subdomain_name` is also required to be set.

* `network_injection` - (Optional) A `network_injection` block as defined below. Only applicable if the `kind` is set to `AIServices`.

* `outbound_network_access_restricted` - (Optional) Whether outbound network access is restricted for the Cognitive Account. Defaults to `false`.

* `project_management_enabled` - (Optional) Whether project management is enabled when the `kind` is set to `AIServices`. Once enabled, `project_management_enabled` cannot be disabled. Changing this forces a new resource to be created. Defaults to `false`.

* `public_network_access_enabled` - (Optional) Whether public network access is allowed for the Cognitive Account. Defaults to `true`.

* `qna_runtime_endpoint` - (Optional) A URL to link a QnAMaker cognitive account to a QnA runtime.

* `custom_question_answering_search_service_id` - (Optional) If `kind` is `TextAnalytics` this specifies the ID of the Search service.

* `custom_question_answering_search_service_key` - (Optional) If `kind` is `TextAnalytics` this specifies the key of the Search service.

-> **Note:** `custom_question_answering_search_service_id` and `custom_question_answering_search_service_key` are used for [Custom Question Answering, the renamed version of QnA Maker](https://docs.microsoft.com/azure/cognitive-services/qnamaker/custom-question-answering), while `qna_runtime_endpoint` is used for [the old version of QnA Maker](https://docs.microsoft.com/azure/cognitive-services/qnamaker/overview/overview)

* `storage` - (Optional) A `storage` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `network_acls` block supports the following:

* `bypass` - (Optional) Whether to allow trusted Azure Services to access the service. Possible values are `None` and `AzureServices`.

-> **Note:** `bypass` can only be set when `kind` is set to `OpenAI` or `AIServices`.

* `default_action` - (Required) The Default Action to use when no rules match from `ip_rules` / `virtual_network_rules`. Possible values are `Allow` and `Deny`.

* `ip_rules` - (Optional) One or more IP Addresses, or CIDR Blocks which should be able to access the Cognitive Account.

* `virtual_network_rules` - (Optional) A `virtual_network_rules` block as defined below.

---

A `network_injection` block supports the following:

* `scenario` - (Required) Specifies what features network injection applies to. The only possible value is `agent`.

* `subnet_id` - (Required) The ID of the subnet which the Agent Client is injected into.

~> **Note:** The agent subnet must use an address space in the 172.* or 192.* ranges.

---

A `virtual_network_rules` block supports the following:

* `subnet_id` - (Required) The ID of the subnet which should be able to access this Cognitive Account.

* `ignore_missing_vnet_service_endpoint` - (Optional) Whether ignore missing vnet service endpoint or not. Defaults to `false`.

---

A `customer_managed_key` block supports the following:

* `key_vault_key_id` - (Required) The ID of the Key Vault Key which should be used to Encrypt the data in this Cognitive Account.

* `identity_client_id` - (Optional) The Client ID of the User Assigned Identity that has access to the key. This property only needs to be specified when there're multiple identities attached to the Cognitive Account.

~> **Note:** When `project_management_enabled` is set to `true`, removing this block forces a new resource to be created.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Cognitive Account. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Cognitive Account.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `storage` block supports the following:

* `storage_account_id` - (Required) Full resource id of a Microsoft.Storage resource.

* `identity_client_id` - (Optional) The client ID of the managed identity associated with the storage resource.

~> **Note:** Not all `kind` support a `storage` block. For example the `kind` `OpenAI` does not support it.

## Attribute Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Service Account.

* `endpoint` - The endpoint used to connect to the Cognitive Service Account.

* `identity` - An `identity` block as defined below.

* `primary_access_key` - A primary access key which can be used to connect to the Cognitive Service Account.

* `secondary_access_key` - The secondary access key which can be used to connect to the Cognitive Service Account.

-> **Note:** The `primary_access_key` and `secondary_access_key` properties are only available when `local_auth_enabled` is set to `true`.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Service Account.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Service Account.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Service Account.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Service Account.

## Import

Cognitive Service Accounts can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account.account1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CognitiveServices` - 2025-06-01

* `Microsoft.Network` - 2024-05-01
