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

* `agent_profile` - A `agent_profile` block as defined below.

* `dev_center_project_resource_id` - The ID of the Dev Center project.

* `fabric_profile` - A `fabric_profile` block as defined below.

* `identity` - An `identity` block as defined below.

* `location` - The Azure Region where the Managed DevOps Pool exists.

* `maximum_concurrency` - Defines how many resources can there be created at any given time.

* `organization_profile` - A `organization_profile` block as defined below.

* `provisioning_state` - The status of the pool operation.

* `tags` - A mapping of tags assigned to the Managed DevOps Pool.

* `type` - This is the resource type `Microsoft.DevOpsInfrastructure/pools.

---

A `agent_profile` block supports the following:

* `kind` - Defines the type of agent profile. Possible values are: `Stateful` and `Stateless`.

* `grace_period_time_span` - How long should the machine be kept around after it ran a workload when there are no stand-by agents. The maximum is one week. This is needed when kind is `Stateful`.

* `max_agent_lifetime` - How long should stateful machines be kept around. The maximum is one week. This is needed when kind is `Stateful`.

* `resource_predictions` - A `resource_predictions` block as defined below.

* `resource_predictions_profile` - One or more `resource_predictions_profile` blocks as defined below.

---

A `resource_predictions` block supports the following:

* `days_data` - A JSON string containing a list of maps, where each map represents a day of the week. Each map includes keys for start and end times, with corresponding values indicating the number of agents available during those times. For example, jsonencode([{},{"09:00:00": 1, "17:00:00": 0},{},{},{},{},{}]) specifies 1 standby agent available every Monday from 9:00 AM to 5:00 PM.

* `time_zone` - Specifies the time zone for the predictions data to be provisioned at. Default is UTC.

---

A `resource_predictions_profile` block supports the following:

* `kind` - Determines how the stand-by scheme should be provided. Possible values are: `Manual` and `Automatic`.

* `prediction_preference` - Specifies the desired balance between cost and performance. Accepted values are: `MostCostEffective`, `MoreCostEffective`, `Balanced`, `MorePerformance`, and `BestPerformance`.

---

A `fabric_profile` block supports the following:

* `images` - One or more `images` blocks as defined below.

* `kind` - Discriminator property for FabricProfile. Possible value is: `Vmss`.

* `sku` - A `sku` block as defined below.

* `network_profile` - A `network_profile` block as defined below.

* `os_profile` - A `os_profile` block as defined below.

* `storage_profile` - A `storage_profile` block as defined below.

---

A `images` block supports the following:

* `aliases` - List of aliases to reference the image by.

* `buffer` - The percentage of the buffer to be allocated to this image.

* `resource_id` - The resource id of the image.

* `well_known_image_name` - The image to use from a well-known set of images made available to customers.

---

A `sku` block supports the following:

* `name` - (Required) The Azure SKU of the machines in the pool.

---

A `network_profile` block supports the following:

* `subnet_id` - The subnet id on which to put all machines created in the pool.

---

A `os_profile` block supports the following:

* `logon_type` - Determines how the service should be run. Accepted values are: `Interactive` and `Service`.

* `secrets_management_settings` - A `secrets_management_settings` block as defined below.

---

A `secrets_management_settings` block supports the following:

* `certificate_store_location` - Specified where to store certificates on the machine.

* `certificate_store_name` - Name of the certificate store to use on the machine. Possible values are: 'My' and 'Root'.

* `key_exportable` - Defines if the key of the certificates should be exportable.

* `observed_certificates` - Specifies the list of certificates from Azure Key vault to install on all machines in the pool.

---

A `storage_profile` block supports the following:

* `data_disks` - One or more `data_disks` blocks as defined above.

* `os_disk_storage_account_type` - The storage account type of the OS disk. Possible values are: `Premium`, `Standard` and `StandardSSD`.

---

A `data_disks` block supports the following:

* `caching` - The type of caching in a data disk. Possible values are: `None`, `ReadOnly` and `ReadWrite`.

* `disk_size_gb` - The initial disk size in gigabytes.

* `drive_letter` - The drive letter for the empty data disk. If not specified, it will be the first available letter.

* `storage_account_type` - The storage Account type to be used for the data disk. If omitted, the default is "standard_lrs".

---

A `identity` block exports the following:

* `identity_ids` - Specifies a list User assigned managed identity Id's.

* `principal_id` - The principal ID for the identity.

* `tenant_id` - The tenant ID for the identity.

* `type` - Type of managed service identity. Possible value is `UserAssigned`.

---

A `organization_profile` block supports the following:

* `kind` - Discriminator property for OrganizationProfile. Possible value is`AzureDevOps` currently.

* `organizations` - One or more `organizations` blocks as defined below.

* `permission_profile` - One or more `permission_profile` blocks as defined below.

---

A `organizations` block supports the following:

* `url` - (Required) The Azure DevOps organization URL in which the pool should be created.

* `parallelism` - (Optional) Specifies how many machines can be created at maximum in this organization out of the `maximum_concurrency` of the pool.

* `projects` - (Optional) List of projects in which the pool should be created.

---

A `permission_profile` block supports the following:

* `kind` - (Required) Determines who has admin permissions to the Azure DevOps pool. Accepted values are: `CreatorOnly`, `Inherit` and `SpecificAccounts`.

* `groups` - (Optional) Specifies a list of group email addresses.

* `users` - (Optional) Specifies a list of user email addresses.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Managed DevOps Pool.
