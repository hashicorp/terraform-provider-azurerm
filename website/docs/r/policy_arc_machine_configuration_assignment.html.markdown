---
subcategory: "Policy"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_policy_arc_machine_configuration_assignment"
description: |-
  Manages a Policy.
---

# azurerm_policy_arc_machine_configuration_assignment

Applies a Guest Configuration Policy to an Arc machine.

## Example Usage

```hcl
data "azurerm_arc_machine" "machine" {
  name                = "example-machine"
  resource_group_name = "example-rg"
}

data "azurerm_storage_blob" "configuration" {
  name                   = "example-configuration"
  storage_account_name   = "example-storage"
  storage_container_name = "example-container"
}

resource "azurerm_policy_arc_machine_configuration_assignment" "import" {
  name       = "example"
  location   = "West Europe"
  machine_id = data.azurerm_arc_machine.machine.id

  configuration {
    version         = "1.0.0"
    content_uri     = data.azurerm_storage_blob.configuration.url
    content_hash    = "315F5BDB76D078C43B8AC0064E4A0164612B1FCE77C869345BFC94C75894EDD3"
    assignment_type = "Audit"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `configuration` - (Required) A `configuration` block as defined below.

* `location` - (Required) The Azure Region where the Guest Configuration Assignment should exist. Changing this forces a new resource to be created.

* `machine_id` - (Required) The resource ID of the HybridCompute machine which the Guest Configuration Assignment should apply to. Changing this forces a new resource to be created.

* `name` - (Required) The name which should be used for this Guest Configuration Assignment. Changing this forces a new resource to be created.

---

A `configuration` block supports the following:

* `assignment_type` - (Optional) The assignment type for the Guest Configuration Assignment. Possible values are `Audit`, `ApplyAndAutoCorrect`, `ApplyAndMonitor` and `DeployAndAutoCorrect`.

* `content_hash` - (Optional) The sha256 content hash for the Guest Configuration package.

* `content_uri` - (Optional) The content URI where the Guest Configuration package is stored.

* `parameter` - (Optional) One or more `parameter` blocks as defined below.

* `version` - (Optional) The assignment version of the Guest Configuration Assignment being applied.

---

A `parameter` block supports the following:

* `name` - (Required) The name of the configuration parameter to apply for the Guest Configuration package.

* `value` - (Required) The value of the configuration parameter to apply for the Guest Configuration package.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Guest Configuration Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Policy.
* `read` - (Defaults to 5 minutes) Used when retrieving the Policy.
* `update` - (Defaults to 30 minutes) Used when updating the Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the Policy.

## Import

Policy Arc Machine Configuration Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_policy_arc_machine_configuration_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.HybridCompute/machines/arcMachine1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/configuration1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.HybridCompute` - 2024-04-05
