---
subcategory: "Managed DevOps Pools"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_devops_pool"
description: |-
  Manages a Managed DevOps Pool.
---

# azurerm_managed_devops_pool

Manages a Managed DevOps Pool.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_dev_center" "example" {
  name                = "example-devcenter"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_dev_center_project" "example" {
  dev_center_id       = azurerm_dev_center.example.id
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_managed_devops_pool" "example" {
  name                  = "example-manageddevopspools"
  resource_group_name   = azurerm_resource_group.example.name
  location              = azurerm_resource_group.example.location
  dev_center_project_id = azurerm_dev_center_project.example.id
  maximum_concurrency   = 1

  azure_devops_organization {
    organization {
      parallelism = 1
      url         = "https://dev.azure.com/example"
    }
  }

  stateless_agent {}

  virtual_machine_scale_set_fabric {
    sku_name = "Standard_D2ads_v5"

    image {
      well_known_image_name = "ubuntu-24.04/buffer"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Managed DevOps Pool. The name must be between 3 and 44 characters, can only include alphanumeric characters, periods (`.`) and hyphens (`-`), must start with an alphanumeric character and cannot end with a period (`.`). Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed DevOps Pool should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Managed DevOps Pool should exist. Changing this forces a new resource to be created.

* `azure_devops_organization` - (Required) An `azure_devops_organization` block as defined below.

* `dev_center_project_id` - (Required) The ID of the Dev Center project.

* `maximum_concurrency` - (Required) Defines how many resources can there be created at any given time. Possible values range between `1` and `10000`.

* `virtual_machine_scale_set_fabric` - (Required) A `virtual_machine_scale_set_fabric` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `stateful_agent` - (Optional) A `stateful_agent` block as defined below.

* `stateless_agent` - (Optional) A `stateless_agent` block as defined below.

~> **Note:** Exactly one of `stateful_agent` or `stateless_agent` must be specified.

* `work_folder` - (Optional) Specifies the work folder for every agent in the pool.

* `tags` - (Optional) A mapping of tags which should be assigned to the Managed DevOps Pool.

---

An `administrator_account` block supports the following:

* `groups` - (Optional) Specifies a list of group email addresses. Changing this forces a new resource to be created.

* `users` - (Optional) Specifies a list of user email addresses. Changing this forces a new resource to be created.

~> **Note:** At least one of `groups` and `users` must be specified.

---

An `automatic_resource_prediction` block supports the following:

* `prediction_preference` - (Optional) Specifies the desired balance between cost and performance. Possible values are `MostCostEffective`, `MoreCostEffective`, `Balanced`, `MorePerformance`, and `BestPerformance`. Defaults to `Balanced`.

---

An `azure_devops_organization` block supports the following:

* `organization` - (Required) One or more `organization` blocks as defined below.

* `permission` - (Optional) A `permission` block as defined below. Changing this forces a new resource to be created.

---

A `daily_schedule` block supports the following:

* `count` - (Required) The number of standby agents to provision at this time. Possible values range between `0` and `maximum_concurrency`.

* `time` - (Required) The time of day at which the agent count changes, in 24-hour format `HH:MM:SS`.

---

An `identity` block supports the following:

* `identity_ids` - (Required) Specifies a list of User Assigned Managed Identity IDs.

* `type` - (Required) The type of managed service identity. The only possible value is `UserAssigned`.

---

An `image` block supports the following:

* `aliases` - (Optional) List of aliases to reference the image by.

* `buffer` - (Optional) The percentage of the buffer to be allocated to this image. Possible values are `*` or between `0` and `100`. Defaults to `*`.

* `id` - (Optional) The resource id of the image.

* `well_known_image_name` - (Optional) The image to use from a well-known set of images made available to customers.

-> **Note:** More information about supported images can be found in [list of Azure Pipelines image predefined aliases](https://learn.microsoft.com/azure/devops/managed-devops-pools/configure-images?view=azure-devops&tabs=arm#azure-pipelines-images). You can optionally specify a version in your `well_known_image_name`, for example `windows-2022/latest` or `windows-2022/20250427.1.0`. If you don't specify a version, latest is used.

~> **Note:** Exactly one of `id` or `well_known_image_name` are required per `image`

---

A `key_vault_management` block supports the following:

* `key_vault_certificate_ids` - (Required) A list of `versionless_id` from Azure Key vault certificates to install on all machines in the pool.

* `certificate_store_location` - (Optional) Specifies where to store certificates on the machine.

* `certificate_store_name` - (Optional) Name of the certificate store to use on the machine. Possible values are `My` and `Root`.

* `key_export_enabled` - (Optional) Defines if the key of the certificates should be exportable. Defaults to `false`.

---

A `manual_resource_prediction` block supports the following:

* `all_week_schedule` - (Optional) A number of agents available 24/7 all week. Possible values range between `1` and `maximum_concurrency`.

* `friday_schedule` - (Optional) One or more `daily_schedule` blocks as defined below.

* `monday_schedule` - (Optional) One or more `daily_schedule` blocks as defined below.

* `saturday_schedule` - (Optional) One or more `daily_schedule` blocks as defined below.

* `sunday_schedule` - (Optional) One or more `daily_schedule` blocks as defined below.

* `thursday_schedule` - (Optional) One or more `daily_schedule` blocks as defined below.

* `time_zone_name` - (Optional) Specifies the time zone for the predictions data to be provisioned at. Defaults to `UTC`.

-> **Note:** A list of possible values for `time_zone_name` are available by executing `[System.TimeZoneInfo]::GetSystemTimeZones()` in PowerShell.

* `tuesday_schedule` - (Optional) One or more `daily_schedule` blocks as defined below.

* `wednesday_schedule` - (Optional) One or more `daily_schedule` blocks as defined below.

~> **Note:** Exactly one of `all_week_schedule` or at least one individual daily schedule block must be specified.

-> **Note:** Please refer to [Microsoft documentation](https://learn.microsoft.com/azure/devops/managed-devops-pools/configure-scaling?view=azure-devops&tabs=azure-cli#manual) for more information about the manual predictions setup.

---

An `organization` block supports the following:

* `parallelism` - (Required) Specifies how many machines can be created at maximum in this organization out of the `maximum_concurrency` of the pool. Possible values range between `1` and `10000`.

~> **Note:** The sum of `parallelism` across orgs should be equal to `maximum_concurrency`.

* `url` - (Required) The Azure DevOps organization URL in which the pool should be created. It must end with a letter or number.

* `projects` - (Optional) List of projects in which the pool should be created.

-> **Note:** Please refer to [Azure DevOps Project Names](https://learn.microsoft.com/azure/devops/organizations/settings/naming-restrictions?view=azure-devops#project-names) for more information on project naming restrictions.

---

A `permission` block supports the following:

* `kind` - (Required) Determines who has admin permissions to the Azure DevOps pool. Possible values are `Inherit` and `SpecificAccounts`. Changing this forces a new resource to be created.

* `administrator_account` - (Optional) An `administrator_account` block as defined below. This block is only valid when `kind` is set to `SpecificAccounts`. Changing this forces a new resource to be created.

---

A `security` block supports the following:

* `interactive_logon_enabled` - (Optional) Specifies whether the agent should run in interactive mode. Defaults to `false`.

* `key_vault_management` - (Optional) A `key_vault_management` block as defined below.

---

A `stateful_agent` block supports the following:

* `automatic_resource_prediction` - (Optional) An `automatic_resource_prediction` block as defined below.

* `grace_period_time_span` - (Optional) Configures the amount of time an agent in a `stateful` pool waits for new jobs before shutting down after all current and queued jobs are complete. The format for Grace Period is `dd.hh:mm:ss` or `hh:mm:ss`. Defaults to `00:00:00`.

* `manual_resource_prediction` - (Optional) A `manual_resource_prediction` block as defined below.

* `maximum_agent_lifetime` - (Optional) Configures the maximum duration an agent in a `stateful` pool can run before it is shut down and discarded. The format for Max time to live for standby agents is `dd.hh:mm:ss` or `hh:mm:ss`. Defaults to `7.00:00:00`.

~> **Note:** Exactly one of `manual_resource_prediction` or `automatic_resource_prediction` may be specified.

---

A `stateless_agent` block supports the following:

* `automatic_resource_prediction` - (Optional) An `automatic_resource_prediction` block as defined below.

* `manual_resource_prediction` - (Optional) A `manual_resource_prediction` block as defined below.

~> **Note:** Exactly one of `manual_resource_prediction` or `automatic_resource_prediction` may be specified.

---

A `storage` block supports the following:

* `disk_size_in_gb` - (Required) The initial disk size in gigabytes. Possible values range between `1` and `32767`.

* `caching` - (Optional) The type of caching for the data disk. Possible values are `ReadOnly` and `ReadWrite`.

* `drive_letter` - (Optional) The drive letter for the data disk.

* `storage_account_type` - (Optional) The storage account type of the data disk. Possible values are `Premium_LRS`, `Premium_ZRS`, `Standard_LRS`, `StandardSSD_LRS`, and `StandardSSD_ZRS`. Defaults to `Standard_LRS`.

---

A `virtual_machine_scale_set_fabric` block supports the following:

* `image` - (Required) One or more `image` blocks as defined below.

* `sku_name` - (Required) The Azure SKU name of the machines in the pool.

-> **Note:** Please refer to the [Microsoft Documentation](https://learn.microsoft.com/azure/devops/managed-devops-pools/configure-pool-settings?view=azure-devops&tabs=azure-portal#agent-size) for more information about available SKUs.

* `os_disk_storage_account_type` - (Optional) The storage account type for the OS disk. Possible values are `Premium`, `Standard`, and `StandardSSD`. Defaults to `Standard`.

* `security` - (Optional) A `security` block as defined below.

* `storage` - (Optional) A `storage` block as defined below.

* `subnet_id` - (Optional) The subnet ID on which to put all machines created in the pool.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed DevOps Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Managed DevOps Pool.
* `read` - (Defaults to 5 minutes) Used when retrieving the Managed DevOps Pool.
* `update` - (Defaults to 30 minutes) Used when updating the Managed DevOps Pool.
* `delete` - (Defaults to 30 minutes) Used when deleting the Managed DevOps Pool.

## Import

Managed DevOps Pool can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_devops_pool.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.DevOpsInfrastructure/pools/pool1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.DevOpsInfrastructure` - 2025-09-20
