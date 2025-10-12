---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_machine_learning_registry"
description: |-
  Manages a Machine Learning Registry.
---

# azurerm_machine_learning_registry

Manages a Machine Learning Registry.

## Example Usage

```hcl
resource "azurerm_machine_learning_registry" "example" {
  name = "example"
  resource_group_name = "example"
  location = "West Europe"
  identity {
    type = "SystemAssigned"
  }
  main_region {
    location = "West Europe"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `identity` - (Required) A `identity` block as defined below.

* `location` - (Required) The Azure Region where the Machine Learning Registry should exist. Changing this forces a new Machine Learning Registry to be created.

* `main_region` - (Required) A `main_region` block as defined below.

* `name` - (Required) The name which should be used for this Machine Learning Registry. Changing this forces a new Machine Learning Registry to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Machine Learning Registry should exist. Changing this forces a new Machine Learning Registry to be created.

---

* `public_network_access_enabled` - (Optional) Whether to enable the TODO. Defaults to `true`.

* `replication_region` - (Optional) One or more `replication_region` blocks as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning Registry.

---

A `identity` block supports the following:

* `type` - (Required) TODO.

* `identity_ids` - (Optional) Specifies a list of TODO.

---

A `main_region` block supports the following:

* `custom_container_registry_account_id` - (Optional) The ID of the TODO.

* `custom_storage_account_id` - (Optional) The ID of the user supplied storage account.Conflicts with `main_region.0.storage_account_type`.

* `hns_enabled` - (Optional) Whether to enable the hierarchical namespace feature for the blob storage container. Defaults to `false`.

* `location` - (Required) The Azure Region where the Machine Learning Registry should exist. Changing this forces a new Machine Learning Registry to be created. It must be the same location as the Registry.

* `storage_account_type` - (Optional) The type of blob storage to use. Defaults to `Standard_LRS`.

---

A `replication_region` block supports the following:

* `location` - (Required) The Azure Region where the Machine Learning Registry should exist. Changing this forces a new Machine Learning Registry to be created.

* `custom_container_registry_account_id` - (Optional) The ID of the TODO.

* `custom_storage_account_id` - (Optional) The ID of the TODO.

* `hns_enabled` - (Optional) Whether to enable the hierarchical namespace feature for the blob storage container. Defaults to `false`.

* `storage_account_type` - (Optional) TODO.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Machine Learning Registry.

* `discovery_url` - TODO.

* `intellectual_property_publisher` - TODO.

* `managed_resource_group` - TODO.

* `ml_flow_registry_uri` - TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Machine Learning Registry.
* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Registry.
* `update` - (Defaults to 30 minutes) Used when updating the Machine Learning Registry.
* `delete` - (Defaults to 30 minutes) Used when deleting the Machine Learning Registry.

## Import

Machine Learning Registrys can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_machine_learning_registry.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.MachineLearningServices/registries/exampleregistry
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.MachineLearningServices` - 2025-06-01
