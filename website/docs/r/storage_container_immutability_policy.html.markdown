---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_container_immutability_policy"
description: |-
  Manages an Immutability Policy for a Container within an Azure Storage Account.
---

# azurerm_storage_container_immutability_policy

Manages an Immutability Policy for a Container within an Azure Storage Account.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoraccount"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_container" "example" {
  name                  = "example"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_storage_container_immutability_policy" "example" {
  storage_container_resource_manager_id = azurerm_storage_container.example.resource_manager_id
  immutability_period_in_days           = 14
  protected_append_writes_all_enabled   = false
  protected_append_writes_enabled       = true
}
```

## Argument Reference

The following arguments are supported:

* `storage_container_resource_manager_id` - (Required) The Resource Manager ID of the Storage Container where this Immutability Policy should be applied. Changing this forces a new resource to be created.

* `immutability_period_in_days` - (Required) The time interval in days that the data needs to be kept in a non-erasable and non-modifiable state.

* `locked` - (Optional) Whether to lock this immutability policy. Cannot be set to `false` once the policy has been locked.

!> **Note:** Once an Immutability Policy has been locked, it cannot be unlocked. After locking, it will only be possible to increase the value for `retention_period_in_days` up to 5 times for the lifetime of the policy. No other properties will be updateable. Furthermore, the Storage Container and the Storage Account in which it resides will become protected by the policy. It will no longer be possible to delete the Storage Container or the Storage Account. Please refer to [official documentation](https://learn.microsoft.com/en-us/azure/storage/blobs/immutable-policy-configure-container-scope?tabs=azure-portal#lock-a-time-based-retention-policy) for more information.

* `protected_append_writes_all_enabled` - (Optional) Whether to allow protected append writes to block and append blobs to the container. Defaults to `false`. Cannot be set with `protected_append_writes_enabled`.

* `protected_append_writes_enabled` - (Optional) Whether to allow protected append writes to append blobs to the container. Defaults to `false`. Cannot be set with `protected_append_writes_all_enabled`.

## Attributes Reference

No additional attributes are exported.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Storage Container Immutability Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Container Immutability Policy.
* `update` - (Defaults to 10 minutes) Used when updating the Storage Container Immutability Policy.
* `delete` - (Defaults to 10 minutes) Used when deleting the Storage Container Immutability Policy.

## Import

Storage Container Immutability Policies can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_container_immutability_policy.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myresourcegroup/providers/Microsoft.Storage/storageAccounts/myaccount/blobServices/default/containers/mycontainer/immutabilityPolicies/default
```
