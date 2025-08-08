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

* `name` - (Required) The name of this Managed DevOps Pool. Changing this forces a new Managed DevOps Pool to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed DevOps Pool exists. Changing this forces a new Managed DevOps Pool to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed DevOps Pool.

* `agent_profile` - An `agent_profile` block as defined below.

* `dev_center_project_resource_id` - The ID of the Dev Center project.

* `fabric_profile` - A `fabric_profile` block as defined below.

* `identity` - An `identity` block as defined below.

* `location` - The Azure Region where the Managed DevOps Pool exists.

* `maximum_concurrency` - The maximum number of agents that can be created.

* `organization_profile` - An `organization_profile` block as defined below.

* `provisioning_state` - The status of the pool operation.

* `tags` - A mapping of tags assigned to the Managed DevOps Pool.

---

An `agent_profile` block exports the following:

* `kind` - Defines the type of agent profile. Possible values are: `Stateful` and `Stateless`.

* `grace_period_time_span` - The amount of time an agent in a `stateful` pool waits for new jobs before shutting down after all current and queued jobs are complete.

* `max_agent_lifetime` - The maximum duration an agent in a `stateful` pool can run before it is shut down and discarded.

* `resource_predictions` - A `resource_predictions` block as defined below.

* `resource_predictions_profile` - A `resource_predictions_profile` block as defined below.

---

A `resource_predictions` block exports the following:

* `days_data` - A JSON string containing a list of maps, where each map represents a day of the week. Each map includes keys for start and end times, with corresponding values indicating the number of agents available during those times.

* `time_zone` - The time zone for the predictions data to be provisioned at.

---

A `resource_predictions_profile` block exports the following:

* `kind` - Determines how the stand-by scheme should be provided. Possible values are: `Manual` and `Automatic`.

* `prediction_preference` - The desired balance between cost and performance.

---

A `fabric_profile` block exports the following:

* `image` - One or more `image` blocks as defined below.

* `kind` - Discriminator property for FabricProfile. Possible value is: `Vmss`.

* `sku` - A `sku` block as defined below.

* `network_profile` - A `network_profile` block as defined below.

* `os_profile` - An `os_profile` block as defined below.

* `storage_profile` - A `storage_profile` block as defined below.

---

An `image` block exports the following:

* `aliases` - A list of image aliases.

* `buffer` - The percentage of the buffer allocated to this image.

* `resource_id` - The resource id of the image.

* `well_known_image_name` - The image name from a well-known set of images made available to customers.

---

A `sku` block exports the following:

* `name` - The Azure SKU of the machines in the pool.

---

A `network_profile` block exports the following:

* `subnet_id` - The ID of the subnet associated with the Managed DevOps Pool.

---

An `os_profile` block exports the following:

* `logon_type` - The logon type.

* `secrets_management_settings` - A `secrets_management_settings` block as defined below.

---

A `secrets_management_settings` block exports the following:

* `certificate_store_location` -  The location where the certificates are stored.

* `certificate_store_name` - The certificate store name.

* `key_export_enabled` - Whether the keys of the certificates are exportable.

* `observed_certificates` - A list of certificates installed on the machines in the Managed DevOps Pool.

---

A `storage_profile` block exports the following:

* `data_disk` - A `data_disk` block as defined above.

* `os_disk_storage_account_type` - The storage account type of the OS disk.

---

A `data_disk` block exports the following:

* `caching` - The type of caching on the data disk.

* `disk_size_gb` - The initial disk size in gigabytes.

* `drive_letter` - The drive letter for the empty data disk.

* `storage_account_type` - The storage account type of the data disk.

---

An `identity` block exports the following:

* `identity_ids` - Specifies a list User assigned managed identity Id's.

* `principal_id` - The principal ID for the identity.

* `tenant_id` - The tenant ID for the identity.

* `type` - The type of managed service identity. The only possible value is `UserAssigned`.

---

An `organization_profile` block exports the following:

* `kind` - Discriminator property for OrganizationProfile. Possible value is`AzureDevOps` currently.

* `organizations` - One or more `organization` blocks as defined below.

* `permission_profile` - One or more `permission_profile` blocks as defined below.

---

An `organization` block exports the following:

* `url` - The Azure DevOps organization URL in which the pool should be created.

* `parallelism` - (Optional) Specifies how many machines can be created at maximum in this organization out of the `maximum_concurrency` of the pool.

* `projects` - (Optional) List of projects in which the pool should be created.

---

A `permission_profile` block exports the following:

* `kind` - Determines who has admin permissions to the Azure DevOps pool. Possible values are: `CreatorOnly`, `Inherit` and `SpecificAccounts`.

* `groups` - (Optional) Specifies a list of group email addresses.

* `users` - (Optional) Specifies a list of user email addresses.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed DevOps Pool.
