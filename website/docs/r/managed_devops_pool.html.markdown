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

resource "azurerm_user_assigned_identity" "test" {
  name                = "example-uai"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_dev_center" "test" {
  name                = "example-devcenter"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_dev_center_project" "example" {
  dev_center_id       = azurerm_dev_center.example.id
  location            = azurerm_resource_group.example.location
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_virtual_network" "test" {
  name                = "example-vnet"
  address_space       = ["10.0.0.0/16", "ace:cab:deca::/48"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                            = "example-subnet"
  resource_group_name             = azurerm_resource_group.test.name
  virtual_network_name            = azurerm_virtual_network.test.name
  address_prefixes                = ["10.0.2.0/24"]
  default_outbound_access_enabled = false
}

data "azurerm_platform_image" "test" {
  location  = azurerm_resource_group.test.location
  publisher = "Canonical"
  offer     = "0001-com-ubuntu-server-focal"
  sku       = "20_04-lts-gen2"
}

resource "azurerm_managed_devops_pool" "example" {
  name                           = "example"
  resource_group_name            = azurerm_resource_group.example.name
  location                       = azurerm_resource_group.example.location
  dev_center_project_resource_id = azurerm_dev_center_project.example.id
  maximum_concurrency            = 1

  azure_devops_organization_profile {
    organization {
      parallelism = 1
      url         = "https://dev.azure.com/example"
    }

    permission_profile {
      kind = "SpecificAccounts"
      
      administrator_accounts {
        groups = ["group1@example.com", "group2@example.com"]
        users  = ["user1@example.com", "user2@example.com"]
      }
    }
  }

  stateful_agent_profile {
    manual_resource_predictions_profile {
      time_zone = "Eastern Standard Time"
      
      monday_schedule    = { "09:00:00" = 1, "17:00:00" = 0 }
      tuesday_schedule   = { "09:00:00" = 1, "17:00:00" = 0 }
      wednesday_schedule = { "09:00:00" = 1, "17:00:00" = 0 }
      thursday_schedule  = { "09:00:00" = 1, "17:00:00" = 0 }
      friday_schedule    = { "09:00:00" = 1, "17:00:00" = 0 }
    }
  }

  vmss_fabric_profile {
    sku_name = "Standard_D2ads_v5"

    image {
      resource_id = data.azurerm_platform_image.test.id
      buffer      = "*"
    }

    image {
      well_known_image_name = "ubuntu-24.0"
      buffer                = "*"
      alias                 = "well known image"
    }

    storage_profile {
      os_disk_storage_account_type = "Standard"

      data_disk {
        disk_size_gb = 1
        drive_letter = "B"
      }
    }

    network_profile {
      subnet_id = azurerm_subnet.test.id
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Managed DevOps Pool. Changing this forces a new Managed DevOps Pool to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed DevOps Pool should exist. Changing this forces a new Managed DevOps Pool to be created.

* `location` - (Required) The Azure Region where the Managed DevOps Pool should exist. Changing this forces a new Managed DevOps Pool to be created.

* `dev_center_project_resource_id` - (Required) The ID of the Dev Center project.

* `vmss_fabric_profile` - (Required) A `vmss_fabric_profile` block as defined below.

* `maximum_concurrency` - (Required) Defines how many resources can there be created at any given time.

* `azure_devops_organization_profile` - (Required) An `azure_devops_organization_profile` block as defined below.

* `stateful_agent_profile` - (Optional) A `stateful_agent_profile` block as defined below.

* `stateless_agent_profile` - (Optional) A `stateless_agent_profile` block as defined below.

~> **NOTE:** Exactly one of `stateful_agent_profile` or `stateless_agent_profile` must be specified.

---

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Managed DevOps Pool.

---

A `stateful_agent_profile` block supports the following:

* `grace_period_time_span` - (Optional) Configures the amount of time an agent in a `stateful` pool waits for new jobs before shutting down after all current and queued jobs are complete. The format for Grace Period is `dd.hh:mm:ss`. Defaults to `00:00:00`.

* `max_agent_lifetime` - (Optional) Configures the maximum duration an agent in a `stateful` pool can run before it is shut down and discarded. The format for Max time to live for standby agents is `dd.hh:mm:ss`. Defaults to `7.00:00:00`.

* `manual_resource_predictions_profile` - (Optional) A `manual_resource_predictions_profile` block as defined below.

* `automatic_resource_predictions_profile` - (Optional) An `automatic_resource_predictions_profile` block as defined below.

~> **NOTE:** Exactly one of `manual_resource_predictions_profile` or `automatic_resource_predictions_profile` may be specified.

---

A `stateless_agent_profile` block supports the following:

* `manual_resource_predictions_profile` - (Optional) A `manual_resource_predictions_profile` block as defined below.

* `automatic_resource_predictions_profile` - (Optional) An `automatic_resource_predictions_profile` block as defined below.

~> **NOTE:** Exactly one of `manual_resource_predictions_profile` or `automatic_resource_predictions_profile` may be specified.

---

A `data_disk` block supports the following:

* `caching` - (Optional) The type of caching in a data disk. Possible values are `None`, `ReadOnly` and `ReadWrite`.

* `disk_size_gb` - (Optional) The initial disk size in gigabytes. Possible values are between `1` and `32767`.

* `drive_letter` - (Optional) The drive letter for the empty data disk. If not specified, it will be the first available letter.

* `storage_account_type` - (Optional) The Storage Account type to be used for the data disk. Possible values are `Premium_LRS`, `Premium_ZRS`, `Standard_LRS`, `StandardSSD_LRS` and `StandardSSD_ZRS`. Defaults to `Standard_LRS`.

---

A `vmss_fabric_profile` block supports the following:

* `image` - (Required) One or more `image` blocks as defined below.

* `sku_name` - (Required) The Azure SKU name of the machines in the pool.

* `network_profile` - (Optional) A `network_profile` block as defined below.

* `os_profile` - (Optional) An `os_profile` block as defined below.

* `storage_profile` - (Optional) A `storage_profile` block as defined below.

---

An `identity` block supports the following:

* `type` - (Required) The type of managed service identity. Possible value is `UserAssigned`.

* `identity_ids` - (Required) Specifies a list User assigned managed identity Id's.

---

An `image` block supports the following:

* `aliases` - (Optional) List of aliases to reference the image by.

* `buffer` - (Optional) The percentage of the buffer to be allocated to this image. Possible values are `*` or between `0` and `100`.

* `resource_id` - (Optional) The resource id of the image.

* `well_known_image_name` - (Optional) The image to use from a well-known set of images made available to customers.

~> **Note:** More information about supported images can be found in [Microsoft Learn documentation](https://learn.microsoft.com/en-us/azure/devops/managed-devops-pools/configure-pool-settings?view=azure-devops&tabs=azure-portal#images)

~> **Note:** Exactly one of `resource_id` or `well_known_image_name` are required per `image`

---

A `network_profile` block supports the following:

* `subnet_id` - (Required) The subnet id on which to put all machines created in the pool.

---

An `azure_devops_organization_profile` block supports the following:

* `organization` - (Required) One or more `organization` blocks as defined below.

* `permission_profile` - (Optional) A `permission_profile` block as defined below.

---

A `permission_profile` block supports the following:

* `kind` - (Required) Determines who has admin permissions to the Azure DevOps pool. Possible values are `CreatorOnly`, `Inherit` and `SpecificAccounts`.

* `administrator_accounts` - (Optional) One or more `administrator_accounts` block as defined below. This block is only valid when `kind` is set to `SpecificAccounts`.

---

An `organization` block supports the following:

* `url` - (Required) The Azure DevOps organization URL in which the pool should be created.

* `parallelism` - (Optional) Specifies how many machines can be created at maximum in this organization out of the `maximum_concurrency` of the pool.

* `projects` - (Optional) List of projects in which the pool should be created.

---

An `administrator_accounts` block supports the following:

* `groups` - (Optional) Specifies a list of group email addresses.

* `users` - (Optional) Specifies a list of user email addresses.

---

An `os_profile` block supports the following:

* `logon_type` - (Optional) Determines how the service should be run. Possible values are `Interactive` and `Service`. Defaults to `Service`.

* `secrets_management` - (Optional) A `secrets_management` block as defined below.

---

A `resource_predictions_profile` block supports the following:

* `kind` - (Required) Determines how the stand-by scheme should be provided. Possible values are: `Manual` and `Automatic`.

* `prediction_preference` - (Optional) Specifies the desired balance between cost and performance. Possible values are: `MostCostEffective`, `MoreCostEffective`, `Balanced`, `MorePerformance`, and `BestPerformance`. This is only valid when `kind` is set to `Automatic`. Defaults to `Balanced`.

---

A `manual_resource_predictions_profile` block supports the following:

* `time_zone` - (Optional) Specifies the time zone for the predictions data to be provisioned at. Defaults to `UTC`.

* `all_week_schedule` - (Optional) A number of agents available 24/7 all week. Possible values are between `1` and `maximum_concurrency`.

* `sunday_schedule` - (Optional) A map of time-to-agent-count pairs for Sunday. Agent counts must not exceed `maximum_concurrency`.

* `monday_schedule` - (Optional) A map of time-to-agent-count pairs for Monday. Agent counts must not exceed `maximum_concurrency`.

* `tuesday_schedule` - (Optional) A map of time-to-agent-count pairs for Tuesday. Agent counts must not exceed `maximum_concurrency`.

* `wednesday_schedule` - (Optional) A map of time-to-agent-count pairs for Wednesday. Agent counts must not exceed `maximum_concurrency`.

* `thursday_schedule` - (Optional) A map of time-to-agent-count pairs for Thursday. Agent counts must not exceed `maximum_concurrency`.

* `friday_schedule` - (Optional) A map of time-to-agent-count pairs for Friday. Agent counts must not exceed `maximum_concurrency`.

* `saturday_schedule` - (Optional) A map of time-to-agent-count pairs for Saturday. Agent counts must not exceed `maximum_concurrency`.

~> **NOTE:** Either `all_week_schedule` or individual daily schedules can be specified. Time keys must be in 24-hour format (HH:MM:SS). Agent counts must be non-negative integers and cannot exceed the `maximum_concurrency` value. Please refer to [Microsoft documentation](https://learn.microsoft.com/en-us/azure/devops/managed-devops-pools/configure-scaling?view=azure-devops&tabs=azure-cli#manual) for more information about the manual predictions setup.

---

An `automatic_resource_predictions_profile` block supports the following:

* `prediction_preference` - (Optional) Specifies the desired balance between cost and performance. Possible values are: `MostCostEffective`, `MoreCostEffective`, `Balanced`, `MorePerformance`, and `BestPerformance`. Defaults to `Balanced`.

---

A `secrets_management` block supports the following:

* `certificate_store_location` - (Optional) Specifies where to store certificates on the machine.

* `certificate_store_name` - (Optional) Name of the certificate store to use on the machine. Possible values are: 'My' and 'Root'.

* `key_export_enabled` - (Required) Defines if the key of the certificates should be exportable.

* `observed_certificates` - (Required) Specifies the list of certificates from Azure Key vault to install on all machines in the pool.

---

A `storage_profile` block supports the following:

* `data_disk` - (Optional) One or more `data_disk` blocks as defined above.

* `os_disk_storage_account_type` - (Optional) The storage account type of the OS disk. Possible values are: `Premium`, `Standard` and `StandardSSD`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Managed DevOps Pool.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

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
<!-- This section is generated, changes will be overwritten -~>
This resource uses the following Azure API Providers:

* `Microsoft.DevOpsInfrastructure` - 2025-01-21
