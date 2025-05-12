---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_deployment"
description: |-
  Manages a Cognitive Services Account Deployment.
---

# azurerm_cognitive_deployment

Manages a Cognitive Services Account Deployment.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
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
}

```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Services Account Deployment. Changing this forces a new resource to be created.

* `cognitive_account_id` - (Required) The ID of the Cognitive Services Account. Changing this forces a new resource to be created.

* `model` - (Required) A `model` block as defined below. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` block as defined below.

* `dynamic_throttling_enabled` - (Optional) Whether dynamic throttling is enabled.

* `rai_policy_name` - (Optional) The name of RAI policy.

* `version_upgrade_option` - (Optional) Deployment model version upgrade option. Possible values are `OnceNewDefaultVersionAvailable`, `OnceCurrentVersionExpired`, and `NoAutoUpgrade`. Defaults to `OnceNewDefaultVersionAvailable`.

---

A `model` block supports the following:

* `format` - (Required) The format of the Cognitive Services Account Deployment model. Possible values are `OpenAI` and `Cohere`. Changing this forces a new resource to be created.

* `name` - (Required) The name of the Cognitive Services Account Deployment model. Changing this forces a new resource to be created.

* `version` - (Optional) The version of Cognitive Services Account Deployment model. If `version` is not specified, the default version of the model at the time will be assigned.

---

A `sku` block supports the following:

* `name` - (Required) The name of the SKU. Possible values include `Standard`, `DataZoneBatch`, `DataZoneStandard`, `DataZoneProvisionedManaged`, `GlobalBatch`, `GlobalProvisionedManaged`, `GlobalStandard`, and `ProvisionedManaged`.

~> **Note:** `DataZoneProvisionedManaged`, `GlobalProvisionedManaged`, and `ProvisionedManaged` are purchased on-demand at an hourly basis based on the number of deployed PTUs, with substantial term discount available via the purchase of Azure Reservations. Currently, this step cannot be completed using Terraform. For more details, please refer to the [provisioned throughput onboarding documentation](https://learn.microsoft.com/en-us/azure/ai-services/openai/how-to/provisioned-throughput-onboarding).

* `tier` - (Optional) Possible values are `Free`, `Basic`, `Standard`, `Premium`, `Enterprise`. This property is required only when multiple tiers are available with the SKU name. Changing this forces a new resource to be created.

* `size` - (Optional) The SKU size. When the name field is the combination of tier and some other value, this would be the standalone code. Changing this forces a new resource to be created.

* `family` - (Optional) If the service has different generations of hardware, for the same SKU, then that can be captured here. Changing this forces a new resource to be created.

* `capacity` - (Optional) Tokens-per-Minute (TPM). The unit of measure for this field is in the thousands of Tokens-per-Minute. Defaults to `1` which means that the limitation is `1000` tokens per minute. If the resources SKU supports scale in/out then the capacity field should be included in the resources' configuration. If the scale in/out is not supported by the resources SKU then this field can be safely omitted. For more information about TPM please see the [product documentation](https://learn.microsoft.com/azure/ai-services/openai/how-to/quota?tabs=rest).

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Deployment for Azure Cognitive Services Account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Services Account Deployment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Services Account Deployment.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Services Account Deployment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Services Account Deployment.

## Import

Cognitive Services Account Deployment can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_deployment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.CognitiveServices/accounts/account1/deployments/deployment1
```
