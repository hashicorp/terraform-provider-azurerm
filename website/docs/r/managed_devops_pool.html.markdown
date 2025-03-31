---
subcategory: "Managed DevOps Pool"
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

resource "azurerm_managed_devops_pool" "example" {
  name = "example"
  resource_group_name = azurerm_resource_group.example.name
  location = azurerm_resource_group.example.location
  dev_center_project_resource_id = azurerm_dev_center_project.example.id
  maximum_concurrency = 1

  organization_profile {
    kind = "AzureDevOps"

    organizations {
      parallelism = 1
      url = "https://dev.azure.com/example""      
    }    
  }

  agent_profile {
    kind = "Stateless"

    resource_predictions_profile {
      kind = "Automatic"
      prediction_preference = "Balanced"
    }
  }

  fabric_profile {
    kind = "Vmss" 

    sku {
      name = "Standard_D2ads_v5"     
    }

    images {
      resource_id = "/Subscriptions/00000000-0000-0000-0000-000000000000/Providers/Microsoft.Compute/Locations/australiaeast/publishers/canonical/artifacttypes/vmimage/offers/0001-com-ubuntu-server-focal/skus/20_04-lts-gen2/versions/latest"
      buffer = "*"
    }

    storage_profile {
      os_disk_storage_account_type = "Standard"

      data_disks {
        disk_size_gib = 1
        drive_letter = "B"
      }
    }

    network_profile {
      subnet_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `agent_profile` - (Required) A `agent_profile` block as defined below.

* `dev_center_project_resource_id` - (Required) The ID of the Dev Center project.

* `fabric_profile` - (Required) A `fabric_profile` block as defined below.

* `location` - (Required) The Azure Region where the Managed DevOps Pool should exist. Changing this forces a new Managed DevOps Pool to be created.

* `maximum_concurrency` - (Required) Defines how many resources can there be created at any given time.

* `name` - (Required) The name which should be used for this Managed DevOps Pool. Changing this forces a new Managed DevOps Pool to be created.

* `organization_profile` - (Required) A `organization_profile` block as defined below.

* `resource_group_name` - (Required) The name of the Resource Group where the Managed DevOps Pool should exist. Changing this forces a new Managed DevOps Pool to be created.

---

* `identity` - (Optional) A `identity` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Managed DevOps Pool.

---

A `agent_profile` block supports the following:

* `kind` - (Required) Defines the type of agent profile. This can be `Stateful` or `Stateless`.

* `grace_period_time_span` - (Optional) How long should the machine be kept around after it ran a workload when there are no stand-by agents. The maximum is one week. This is needed when kind is Stateful.

* `max_agent_lifetime` - (Optional) How long should stateful machines be kept around. The maximum is one week. This is needed when kind is Stateful.

* `resource_predictions` - (Optional) A `resource_predictions` block as defined below.

* `resource_predictions_profile` - (Optional) One or more `resource_predictions_profile` blocks as defined below.

---

A `data_disks` block supports the following:

* `caching` - (Optional) The type of caching in a data disk. Accepted values are: `None`, `ReadOnly` and `ReadWrite`.

* `disk_size_gb` - (Optional) The initial disk size in gigabytes.

* `drive_letter` - (Optional) The drive letter for the empty data disk. If not specified, it will be the first available letter.

* `storage_account_type` - (Optional) The storage Account type to be used for the data disk. If omitted, the default is "standard_lrs".

---

A `fabric_profile` block supports the following:

* `images` - (Required) One or more `images` blocks as defined below.

* `kind` - (Required) Discriminator property for FabricProfile. Accepted value is: `Vmss`.

* `sku` - (Required) A `sku` block as defined below.

* `network_profile` - (Optional) A `network_profile` block as defined below.

* `os_profile` - (Optional) A `os_profile` block as defined below.

* `storage_profile` - (Optional) A `storage_profile` block as defined below.

---

A `identity` block supports the following:

* `type` - (Required) Type of managed service identity. Accepted value is `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list User assigned managed identity Id's.

---

A `images` block supports the following:

* `aliases` - (Optional) List of aliases to reference the image by.

* `buffer` - (Optional) The percentage of the buffer to be allocated to this image.

* `resource_id` - (Optional) The resource id of the image.

* `well_known_image_name` - (Optional) The image to use from a well-known set of images made available to customers.

---

A `network_profile` block supports the following:

* `subnet_id` - (Required) The subnet id on which to put all machines created in the pool.

---

A `organization_profile` block supports the following:

* `kind` - (Required) Discriminator property for OrganizationProfile. Accepted value is`AzureDevOps` currently.

* `organizations` - (Required) One or more `organizations` blocks as defined below.

* `permission_profile` - (Optional) One or more `permission_profile` blocks as defined below.

---

A `organizations` block supports the following:

* `url` - (Required) The Azure DevOps organization URL in which the pool should be created.

* `parallelism` - (Optional) Specifies how many machines can be created at maximum in this organization out of the maximum_Concurrency of the pool.

* `projects` - (Optional) List of projects in which the pool should be created.

---

A `os_profile` block supports the following:

* `logon_type` - (Required) Determines how the service should be run. Accepted values are: `Interactive` and `Service`.

* `secrets_management_settings` - (Optional) A `secrets_management_settings` block as defined below.

---

A `permission_profile` block supports the following:

* `kind` - (Required) Determines who has admin permissions to the Azure DevOps pool. Accepted values are: `CreatorOnly`, `Inherit` and `SpecificAccounts`.

* `groups` - (Optional) Specifies a list of group email addresses.

* `users` - (Optional) Specifies a list of user email addresses.

---

A `resource_predictions` block supports the following:

* `days_data` - (Optional) A JSON string containing a list of maps, where each map represents a day of the week. Each map includes keys for start and end times, with corresponding values indicating the number of agents available during those times. For example, jsonencode([{},{"09:00:00": 1, "17:00:00": 0},{},{},{},{},{}]) specifies 1 standby agent available every Monday from 9:00 AM to 5:00 PM.

* `time_zone` - (Optional) Specifies the time zone for the predictions data to be provisioned at. Default is UTC.

---

A `resource_predictions_profile` block supports the following:

* `kind` - (Required) Determines how the stand-by scheme should be provided.

* `prediction_preference` - (Optional) Specifies the desired balance between cost and performance. Accepted values are: `MostCostEffective`, `MoreCostEffective`, `Balanced`, `MorePerformance`, and `BestPerformance`.

---

A `secrets_management_settings` block supports the following:

* `certificate_store_location` - (Optional) Specified where to store certificates on the machine.

* `certificate_store_name` - (Optional) Name of the certificate store to use on the machine, Accepted values are: 'My' and 'Root'.

* `key_exportable` - (Optional) Defines if the key of the certificates should be exportable.

* `observed_certificates` - (Optional) Specifies the list of certificates from Azure Key vault to install on all machines in the pool.

---

A `sku` block supports the following:

* `name` - (Required) The Azure SKU of the machines in the pool.

---

A `storage_profile` block supports the following:

* `data_disks` - (Optional) One or more `data_disks` blocks as defined above.

* `os_disk_storage_account_type` - (Optional) The storage account type of the OS disk. Accepted values are: `Premium`, `Standard` and `StandardSSD`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Managed DevOps Pool.

* `provisioning_state` - The status of the pool.

* `type` - This is the resource type `Microsoft.DevOpsInfrastructure/pools`.

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