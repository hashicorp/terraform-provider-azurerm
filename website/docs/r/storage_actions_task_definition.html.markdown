---
subcategory: "Storage"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_storage_actions_task_definition"
description: |-
  Manages a Storage Actions Task Definition.
---

# azurerm_storage_actions_task_definition

Manages a Storage Actions Task Definition.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_storage_actions_task_definition" "example" {
  name                = "examplestoragetask"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  description         = "Example Storage Actions Task Definition"
  enabled             = true

  identity {
    type = "SystemAssigned"
  }

  action {
    if {
      condition = "[[equals(AccessTier, 'Cool')]]"

      operation {
        name       = "SetBlobTier"
        on_failure = "break"
        on_success = "continue"

        parameters = {
          tier = "Hot"
        }
      }
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Storage Actions Task Definition. The name must be between 3 and 18 characters in length and may contain lowercase letters and numbers only. Changing this forces a new Storage Actions Task Definition to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Storage Actions Task Definition should exist. Changing this forces a new Storage Actions Task Definition to be created.

* `location` - (Required) The Azure Region where the Storage Actions Task Definition should exist. Changing this forces a new Storage Actions Task Definition to be created.

* `action` - (Required) An `action` block as defined below.

* `description` - (Required) A description for this Storage Actions Task Definition.

* `enabled` - (Required) Whether this Storage Actions Task Definition is enabled.

* `identity` - (Required) An `identity` block as defined below.

---

* `tags` - (Optional) A mapping of tags which should be assigned to the Storage Actions Task Definition.

---

An `action` block supports the following:

* `if` - (Required) An `if` block as defined below.

* `else` - (Optional) An `else` block as defined below.

---

An `else` block supports the following:

* `operation` - (Required) One or more `operation` blocks as defined below.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Storage Actions Task Definition. Possible values are `SystemAssigned` and `UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Storage Actions Task Definition.

~> **Note:** This is required when `type` is set to `UserAssigned`.

---

An `if` block supports the following:

* `condition` - (Required) The condition predicate that is evaluated against each blob to determine whether the operations defined under this block should run. See the [Azure Storage Actions condition documentation](https://learn.microsoft.com/azure/storage-actions/storage-tasks/storage-task-conditions) for the supported syntax.

* `operation` - (Required) One or more `operation` blocks as defined below.

---

An `operation` block supports the following:

* `name` - (Required) The name of the operation to perform on each blob matched by the condition. Possible values are `DeleteBlob`, `SetBlobExpiry`, `SetBlobImmutabilityPolicy`, `SetBlobLegalHold`, `SetBlobTags`, `SetBlobTier`, and `UndeleteBlob`.

~> **Note:** When `name` is set to `DeleteBlob` it must be the only `operation` defined within its parent `if` or `else` block.

* `on_failure` - (Optional) The action to take when the operation fails on a blob. The only possible value is `break`.

* `on_success` - (Optional) The action to take when the operation succeeds on a blob. The only possible value is `continue`.

* `parameters` - (Optional) A mapping of parameters used by the operation.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Actions Task Definition.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Storage Actions Task Definition.
* `read` - (Defaults to 5 minutes) Used when retrieving the Storage Actions Task Definition.
* `update` - (Defaults to 30 minutes) Used when updating the Storage Actions Task Definition.
* `delete` - (Defaults to 30 minutes) Used when deleting the Storage Actions Task Definition.

## Import

Storage Actions Task Definitions can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_storage_actions_task_definition.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.StorageActions/storageTasks/storageTask1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.StorageActions` - 2023-01-01
