---
subcategory: "Advisor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_advisor_suppression"
description: |-
  Specifies a suppression for an Azure Advisor recommendation.
---

# azurerm_advisor_suppression

Specifies a suppression for an Azure Advisor recommendation.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_advisor_recommendations" "example" {}

resource "azurerm_advisor_suppression" "example" {
  name              = "HardcodedSuppressionName"
  recommendation_id = data.azurerm_advisor_recommendations.test.recommendations[0].recommendation_name
  resource_id       = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  ttl               = "01:00:00:00"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this Advisor suppression. Changing this forces a new Advisor suppression to be created.

* `recommendation_id` - (Required) The ID of the Advisor recommendation to suppress. Changing this forces a new Advisor suppression to be created.

* `resource_id` - (Required) The ID of the Resource to suppress the Advisor recommendation for. Changing this forces a new Advisor suppression to be created.

---

* `ttl` - (Optional) A optional time to live value. If omitted, the suppression will not expire. Changing this forces a new Advisor suppression to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Advisor suppression.

* `suppression_id` - The GUID of the suppression.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Advisor suppression.
* `read` - (Defaults to 5 minutes) Used when retrieving the Advisor suppression.
* `delete` - (Defaults to 30 minutes) Used when deleting the Advisor suppression.

## Import

Advisor suppressions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_advisor_suppression.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Advisor/recommendations/00000000-0000-0000-0000-000000000000/suppressions/name
```
