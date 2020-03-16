---
subcategory: "Advisor"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_advisor_recommendations"
description: |-
  Gets information about an existing Advisor Recommendations.
---

# Data Source: azurerm_advisor_recommendations

Use this data source to access information about an existing Advisor Recommendations.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_advisor_recommendations" "example" {
  categories_filter           = ["security", "cost"]
  resource_group_names_filter = ["example-resgroups"]
}

output "recommendations" {
  value = data.azurerm_advisor_recommendations.example.recommendations
}
```

## Arguments Reference

The following arguments are supported:

* `categories_filter` - (Optional) Specifies a list of categories in which the Advisor Recommendations will be listed. Possible values are 'HighAvailability', 'Security', 'Performance', 'Cost' and 'OperationalExcellence'.

* `resource_group_names_filter` - (Optional) Specifies a list of resource groups about which the Advisor Recommendations will be listed.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Advisor Recommendations.

* `recommendations` - One or more `recommendations` blocks as defined below.

---

A `recommendations` block exports the following:

* `category` - The category of the recommendation.

* `description` - The description of the issue or the opportunity identified by the recommendation.

* `impact` - The business impact of the recommendation.

* `recommendation_name` - The name of the Advisor Recommendation.

* `recommendation_type_id` - The recommendation type id of the Advisor Recommendation.

* `resource_name` - The name of the identified resource of the Advisor Recommendation.

* `resource_type` - The type of the identified resource of the Advisor Recommendation.

* `suppression_names` - A list of Advisor Suppression names of the Advisor Recommendation.

* `updated_time` - The most recent time that Advisor checked the validity of the recommendation..

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 10 minutes) Used when retrieving the Advisor Recommendations.
