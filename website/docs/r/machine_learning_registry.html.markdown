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
  name                = "example"
  resource_group_name = "example"
  location            = "West Europe"
  identity {
    type = "SystemAssigned"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `identity` - (Required) A `identity` block as defined below.

* `location` - (Required) The Azure Region where the Machine Learning Registry should exist. Changing this forces a new Machine Learning Registry to be created.

* `primary_region` - (Optional) A `primary_region` block as defined below.

* `name` - (Required) The name which should be used for this Machine Learning Registry. Changing this forces a new Machine Learning Registry to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Machine Learning Registry should exist. Changing this forces a new Machine Learning Registry to be created.

---

* `public_network_access_enabled` - (Optional) Whether to enable public network access for the Machine Learning Registry. Defaults to `true`.

* `replication_region` - (Optional) One or more `replication_region` blocks as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning Registry.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Machine Learning Registry. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Machine Learning Registry.

---

A `primary_region` block supports the following:

* `system_created_storage_account_type` - (Optional) The type of blob storage to use. Possible values are `Standard_LRS`, `Standard_GRS`, `Standard_RAGRS`, `Standard_ZRS`, `Standard_GZRS`, `Standard_RAGZRS`, `Premium_LRS` and `Premium_ZRS`. Defaults to `Standard_LRS`.

* `hns_enabled` - (Optional) Whether to enable the hierarchical namespace feature for the blob storage container. Defaults to `false`.

---

A `replication_region` block supports the following:

* `location` - (Required) The Azure Region where the Machine Learning Registry should exist.

* `system_created_storage_account_type` - (Optional) The type of blob storage to use. Possible values are `Standard_LRS`, `Standard_GRS`, `Standard_RAGRS`, `Standard_ZRS`, `Standard_GZRS`, `Standard_RAGZRS`, `Premium_LRS` and `Premium_ZRS`. Defaults to `Standard_LRS`.

* `hns_enabled` - (Optional) Whether to enable the hierarchical namespace feature for the blob storage container. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Machine Learning Registry.

* `discovery_url` - The discovery URL for the Machine Learning Registry.

* `intellectual_property_publisher` - The intellectual property publisher for the Machine Learning Registry.

* `managed_resource_group_id` - The ID of the managed resource group created for the Machine Learning Registry.

* `machine_learning_flow_registry_uri` - The ML Flow registry URI for the Machine Learning Registry.

---

In addition to the above, each `primary_region` and `replication_region` block exports the following:

* `system_created_storage_account_id` - The ID of the system-created storage account for this region.

* `system_created_container_registry_id` - The ID of the system-created container registry for this region.

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
