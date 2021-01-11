---
subcategory: "Data Share"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_share_dataset_kusto_cluster"
description: |-
  Gets information about an existing Data Share Kusto Cluster Dataset.
---

# Data Source: azurerm_data_share_dataset_kusto_cluster

Use this data source to access information about an existing Data Share Kusto Cluster Dataset.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_data_share_dataset_kusto_cluster" "example" {
  name     = "example-dskc"
  share_id = "example-share-id"
}

output "id" {
  value = data.azurerm_data_share_dataset_kusto_cluster.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Data Share Kusto Cluster Dataset.

* `share_id` - (Required) The resource ID of the Data Share where this Data Share Kusto Cluster Dataset should be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The resource ID of the Data Share Kusto Cluster Dataset.

* `kusto_cluster_id` - The resource ID of the Kusto Cluster to be shared with the receiver.

* `display_name` - The name of the Data Share Dataset.

* `kusto_cluster_location` - The location of the Kusto Cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Data Share Kusto Cluster Dataset.
