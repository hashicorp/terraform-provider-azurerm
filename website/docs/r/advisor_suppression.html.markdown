---
subcategory: "Advisor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_advisor_suppression"
description: |-
  Manages an Advisor Suppression.
---

# azurerm_advisor_suppression

Manages an Advisor Suppression.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_advisor_recommendations" "example" {}

resource "azurerm_advisor_suppression" "example" {
  name              = "example-as"
  recommendation_id = data.azurerm_advisor_recommendations.example.recommendations.0.recommendation_id
  duration_in_days  = "3"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this Advisor Suppression. Changing this forces a new Advisor Suppression to be created.

* `recommendation_id` - (Required) The ID of the Advisor Recommendation to be postponed for a specific duration or infinitely. Changing this forces a new Advisor Suppression to be created.

---

* `duration_in_days` - (Optional) The interval in seconds during which the Advisor Recommendation will be suppressed. Specifying `-1` will suppress forever.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Advisor Suppression.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Advisor Suppression.
* `read` - (Defaults to 5 minutes) Used when retrieving the Advisor Suppression.
* `update` - (Defaults to 30 minutes) Used when updating the Advisor Suppression.
* `delete` - (Defaults to 30 minutes) Used when deleting the Advisor Suppression.

## Import

Advisor Suppressions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_advisor_suppression.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/microsoft.web/sites/site1/providers/Microsoft.Advisor/recommendations/recommendation1/suppressions/suppression1
```
