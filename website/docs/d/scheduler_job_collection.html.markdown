---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_scheduler_job_collection"
sidebar_current: "docs-azurerm-datasource-scheduler_job_collection"
description: |-
  Get information about the specified scheduler job collection.
---

# Data Source: azurerm_scheduler_job_collection

Use this data source to access the properties of an Azure scheduler job collection.

## Example Usage

```hcl
data "azurerm_scheduler_job_collection" "test" {
  name                = "tfex-job-collection"
  resource_group_name = "tfex-job-collection-rg"
}

output "job_collection_state" {
  value = "${data.azurerm_scheduler_job_collection.jobs.state}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Scheduler Job Collection. 

* `resource_group_name` - (Required) Specifies the name of the resource group in which the Scheduler Job Collection resides. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Scheduler Job Collection.

* `location` - The Azure location where the resource exists. 

* `tags` - A mapping of tags assigned to the resource.

* `sku` - The Job Collection's pricing level's SKU. 

* `state` - The Job Collection's state. 

* `quota` - The Job collection quotas as documented in the `quota` block below. 

The `quota` block supports:

* `max_job_count` - Sets the maximum number of jobs in the collection. 

* `max_recurrence_frequency` - The maximum frequency of recurrence. 

* `max_retry_interval` - The maximum interval between retries.
