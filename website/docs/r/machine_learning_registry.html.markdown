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

* `identity` - (Required) An `identity` block as defined below.

* `location` - (Required) The Azure Region where the Machine Learning Registry should exist. Changing this forces a new Machine Learning Registry to be created.

* `name` - (Required) The name which should be used for this Machine Learning Registry. Changing this forces a new Machine Learning Registry to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Machine Learning Registry should exist. Changing this forces a new Machine Learning Registry to be created.

---

* `public_network_access_enabled` - (Optional) Whether to enable public network access for the Machine Learning Registry. Defaults to `true`.

* `replication_region` - (Optional) One or more `replication_region` blocks as defined below.

* `system_created_container_registry_sku` - (Optional) The SKU of the system-created container registry in the primary region. The only supported value is `Premium`. Defaults to `Premium`.

* `system_created_storage_account_blob_public_access_enabled` - (Optional) Whether to allow public blob access for the system-created storage account in the primary region. Defaults to `false`.

* `system_created_storage_account_hns_enabled` - (Optional) Whether to enable the hierarchical namespace feature for the system-created storage account in the primary region. Defaults to `false`.

* `system_created_storage_account_type` - (Optional) The storage account type for the system-created storage account in the primary region. Possible values are `Standard_LRS`, `Standard_GRS`, `Standard_RAGRS`, `Standard_ZRS`, `Standard_GZRS`, `Standard_RAGZRS`, `Premium_LRS` and `Premium_ZRS`. Defaults to `Standard_LRS`.

* `tags` - (Optional) A mapping of tags which should be assigned to the Machine Learning Registry.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Machine Learning Registry. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Machine Learning Registry.

---

A `replication_region` block supports the following:

* `location` - (Required) The Azure Region where the replicated Machine Learning Registry resources should exist.

* `system_created_container_registry_sku` - (Optional) The SKU of the system-created container registry in this region. The only supported value is `Premium`. Defaults to `Premium`.

* `system_created_storage_account_blob_public_access_enabled` - (Optional) Whether to allow public blob access for the system-created storage account in this region. Defaults to `false`.

* `system_created_storage_account_hns_enabled` - (Optional) Whether to enable the hierarchical namespace feature for the system-created storage account in this region. Defaults to `false`.

* `system_created_storage_account_type` - (Optional) The storage account type for the system-created storage account in this region. Possible values are `Standard_LRS`, `Standard_GRS`, `Standard_RAGRS`, `Standard_ZRS`, `Standard_GZRS`, `Standard_RAGZRS`, `Premium_LRS` and `Premium_ZRS`. Defaults to `Standard_LRS`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Machine Learning Registry.

* `discovery_url` - The discovery URL for the Machine Learning Registry.

* `managed_resource_group_id` - The ID of the managed resource group created for the Machine Learning Registry.

* `machine_learning_flow_registry_uri` - The ML Flow registry URI for the Machine Learning Registry.

* `system_created_container_registry_id` - The ID of the system-created container registry in the primary region.

* `system_created_container_registry_name` - The name of the system-created container registry in the primary region.

* `system_created_storage_account_id` - The ID of the system-created storage account in the primary region.

* `system_created_storage_account_name` - The name of the system-created storage account in the primary region.

---

In addition to the above, each `replication_region` block exports the following:

* `system_created_storage_account_id` - The ID of the system-created storage account for this region.

* `system_created_storage_account_name` - The name of the system-created storage account for this region.

* `system_created_container_registry_id` - The ID of the system-created container registry for this region.

* `system_created_container_registry_name` - The name of the system-created container registry for this region.

---

A `private_endpoint_connections` block exports the following:

* `id` - The ID of the private endpoint connection.

* `location` - The Azure Region where the private endpoint connection exists.

* `group_ids` - A list of group IDs for the private endpoint connection.

* `subnet_id` - The ID of the Subnet that the private endpoint is connected to.

* `provisioning_state` - The provisioning state of the private endpoint connection.

* `connection_state` - A `connection_state` block as defined below.

---

A `connection_state` block exports the following:

* `status` - The connection status of the service consumer with the service provider.

* `description` - A user-defined message that may be used for approval-related messages.

* `actions_required` - A message indicating if changes on the service provider require any updates on the consumer.

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
