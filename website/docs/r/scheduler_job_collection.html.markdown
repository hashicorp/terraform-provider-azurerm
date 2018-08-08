---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_scheduler_job_collection"
sidebar_current: "docs-azurerm-resource-scheduler-job-collection"
description: |-
  Manages a Scheduler Job Collection.
---

# azurerm_scheduler_job_collection

Manages a Scheduler Job Collection.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_scheduler_job_collection" "example" {
  name                = "example-collection"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "free"
  state               = "enabled"

  quota {
    max_job_count            = 5
    max_recurrence_interval  = 24
    max_recurrence_frequency = "hour"
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Scheduler Job Collection. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Scheduler Job Collection. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `sku` - (Required) Sets the Job Collection's pricing level's SKU. Possible values include: `Standard`, `Free`, `P10Premium`, `P20Premium`.

* `quota` - (Optional) A `quota` block as defined below.

* `state` - (Optional) Sets Job Collection's state. Possible values include: `Enabled`, `Disabled`, `Suspended`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `quota` block supports:

* `max_job_count` - (Optional) Sets the maximum number of jobs in the collection.

* `max_recurrence_frequency` - (Required) The maximum frequency of recurrence. Possible values include: `Minute`, `Hour`, `Day`, `Week` and `Month`.

* `max_recurrence_interval` - (Optional) The maximum interval between recurrence.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Scheduler Job Collection.

## Import

Scheduler Job Collections can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_scheduler_job_collection.jobcollection1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Scheduler/jobCollections/jobcollection1
```
