---
subcategory: "Elastic"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_elastic_cloud_serverless"
description: |-
  Lists Elastic Cloud Serverless resources.
---

# List resource: azurerm_elastic_cloud_serverless

Lists Elastic Cloud Serverless resources.

## Example Usage

### List all Elastic Cloud Serverless projects in the subscription

```hcl
list "azurerm_elastic_cloud_serverless" "example" {
  provider = azurerm
  config {}
}
```

### List all Elastic Cloud Serverless projects in a specific Resource Group

```hcl
list "azurerm_elastic_cloud_serverless" "example" {
  provider = azurerm
  config {
    resource_group_name = "example-resource-group"
  }
}
```

## Argument Reference

This list resource supports the following arguments:

* `resource_group_name` - (Optional) The name of the Resource Group to query.

* `subscription_id` - (Optional) The Subscription ID to query. Defaults to the value specified in the Provider Configuration.