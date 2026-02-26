---
subcategory: "Managed DevOps Pools"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_managed_devops_pool"
description: |-
  Gets information about an existing Managed DevOps Pool.
---

# Data Source: azurerm_managed_devops_pool

Use this data source to access information about an existing Managed DevOps Pool.

## Example Usage

```hcl
data "azurerm_managed_devops_pool" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_managed_devops_pool.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Managed DevOps Pool.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed DevOps Pool exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed DevOps Pool.

* `azure_devops_organization` - An `azure_devops_organization` block as defined below.

* `dev_center_project_id` - The ID of the Dev Center project.

* `identity` - An `identity` block as defined below.

* `location` - The Azure Region where the Managed DevOps Pool exists.

* `maximum_concurrency` - The maximum number of agents that can be created.

* `stateful_agent` - A `stateful_agent` block as defined below.

* `stateless_agent` - A `stateless_agent` block as defined below.

* `tags` - A mapping of tags assigned to the Managed DevOps Pool.

* `vmss_fabric` - A `vmss_fabric` block as defined below.

---

An `administrator_account` block exports the following:

* `groups` - A list of group email addresses.

* `users` - A list of user email addresses.

---

An `automatic_resource_predictions` block exports the following:

* `prediction_preference` - The desired balance between cost and performance.

---

An `azure_devops_organization` block exports the following:

* `organization` - One or more `organization` blocks as defined below.

* `permission` - A `permission` block as defined below.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Managed DevOps Pool

* `identity_ids` - A list of User Assigned Identity IDs assigned to this Managed DevOps Pool.

---

An `image` block exports the following:

* `aliases` - A list of image aliases.

* `buffer` - The percentage of the buffer allocated to this image.

* `id` - The resource id of the image.

* `well_known_image_name` - The image name from a well-known set of images made available to customers.

---

A `manual_resource_prediction` block exports the following:

* `time_zone_name` - The time zone for the predictions data to be provisioned at.

* `all_week_schedule` - A number of agents available 24/7 all week.

* `sunday_schedule` - One or more `sunday_schedule` blocks as defined below.

* `monday_schedule` - One or more `monday_schedule` blocks as defined below.

* `tuesday_schedule` - One or more `tuesday_schedule` blocks as defined below.

* `wednesday_schedule` - One or more `wednesday_schedule` blocks as defined below.

* `thursday_schedule` - One or more `thursday_schedule` blocks as defined below.

* `friday_schedule` - One or more `friday_schedule` blocks as defined below.

* `saturday_schedule` - One or more `saturday_schedule` blocks as defined below.

---

Each schedule block (`sunday_schedule`, `monday_schedule`, `tuesday_schedule`, `wednesday_schedule`, `thursday_schedule`, `friday_schedule`, `saturday_schedule`) exports the following:

* `time` - The time of day at which the agent count changes, in 24-hour format `HH:MM:SS`.

* `count` - The number of standby agents provisioned at this time.

---

An `organization` block exports the following:

* `url` - The URL of  The Azure DevOps organization.

* `parallelism` - Maximum numbers of machines in this organization out of the `maximum_concurrency` of the pool.

* `projects` - A list of projects in which the pool should be created.

---

A `security` block exports the following:

* `interactive_logon_enabled` - Whether the agent runs in interactive mode.

* `key_vault_management` - A `key_vault_management` block as defined below.

---

A `permission` block exports the following:

* `kind` - The type of Azure DevOps pool permission.

* `administrator_account` - One or more `administrator_account` block as defined below.

---

A `key_vault_management` block exports the following:

* `certificate_store_location` -  The location where the certificates are stored.

* `certificate_store_name` - The certificate store name.

* `key_export_enabled` - Whether the keys of the certificates are exportable.

* `key_vault_certificate_ids` - A list of certificates installed on the machines in the Managed DevOps Pool.

---

A `stateful_agent` block exports the following:

* `grace_period_time_span` - The amount of time an agent in a `stateful` pool waits for new jobs before shutting down after all current and queued jobs are complete.

* `max_agent_lifetime` - The maximum duration an agent in a `stateful` pool can run before it is shut down and discarded.

* `manual_resource_prediction` - A `manual_resource_prediction` block as defined below.

* `automatic_resource_prediction` - An `automatic_resource_prediction` block as defined below.

---

A `stateless_agent` block exports the following:

* `manual_resource_prediction` - A `manual_resource_prediction` block as defined below.

* `automatic_resource_prediction` - An `automatic_resource_prediction` block as defined below.

---

A `storage` block exports the following:

* `caching` - The type of caching for the data disk.

* `disk_size_in_gb` - The initial disk size in gigabytes.

* `drive_letter` - The drive letter for the data disk.

* `storage_account_type` - The storage account type of the data disk.

---

A `vmss_fabric` block exports the following:

* `image` - One or more `image` blocks as defined below.

* `sku_name` - The Azure SKU of the machines in the pool.

* `subnet_id` - The ID of the subnet associated with the Managed DevOps Pool.

* `security` - A `security` block as defined below.

* `os_disk_storage_account_type` - The storage account type for the OS disk.

* `storage` - A `storage` block as defined below.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed DevOps Pool.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DevOpsInfrastructure` - 2025-01-21
