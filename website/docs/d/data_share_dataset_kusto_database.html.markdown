---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_share_dataset_kusto_database"
description: |-
  Gets information about an existing Data Share Kusto Database Dataset.
---

# Data Source: azurerm_data_share_dataset_kusto_database

Use this data source to access information about an existing Data Share Kusto Database Dataset.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_data_share_dataset_kusto_database" "example" {
  name     = "example-dskdds"
  share_id = "example-share-id"
}

output "id" {
  value = data.azurerm_data_share_dataset_kusto_database.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Data Share Kusto Database Dataset.

* `share_id` - (Required) The resource ID of the Data Share where this Data Share Kusto Database Dataset should be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The resource ID of the Data Share Kusto Database Dataset.

* `kusto_database_id` - The resource ID of the Kusto Cluster Database to be shared with the receiver.

* `display_name` - The name of the Data Share Dataset.

* `kusto_cluster_location` - The location of the Kusto Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Kusto Database Dataset.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DataShare`: 2019-11-01
