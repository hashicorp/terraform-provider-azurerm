---
subcategory: "Machine Learning"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_machine_learning_registry"
description: |-
  Gets information about an existing Machine Learning Registry.
---

# Data Source: azurerm_machine_learning_registry

Use this data source to access information about an existing Machine Learning Registry.

## Example Usage

```hcl
data "azurerm_machine_learning_registry" "example" {
  name                = "existing-mlregistry"
  resource_group_name = "example-resources"
}

output "registry_id" {
  value = data.azurerm_machine_learning_registry.example.id
}

output "discovery_url" {
  value = data.azurerm_machine_learning_registry.example.discovery_url
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Machine Learning Registry.

* `resource_group_name` - (Required) The name of the Resource Group where the Machine Learning Registry exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Machine Learning Registry.

* `discovery_url` - The discovery URL for the Machine Learning Registry.

* `identity` - An `identity` block as defined below.

* `location` - The Azure Region where the Machine Learning Registry exists.

* `machine_learning_flow_registry_uri` - The ML Flow registry URI for the Machine Learning Registry.

* `managed_resource_group_id` - The ID of the managed resource group created for the Machine Learning Registry.

* `public_network_access_enabled` - Whether public network access is enabled for the Machine Learning Registry.

* `replication_regions` - One or more `replication_regions` blocks as defined below.

* `system_created_container_registry_id` - The ID of the system-created container registry in the primary region.

* `system_created_container_registry_name` - The name of the system-created container registry in the primary region.

* `system_created_container_registry_sku` - The SKU of the system-created container registry in the primary region. This is always `Premium`.

* `system_created_storage_account_hierarchical_namespace_enabled` - Whether hierarchical namespace is enabled for the system-created storage account in the primary region.

* `system_created_storage_account_id` - The ID of the system-created storage account in the primary region.

* `system_created_storage_account_name` - The name of the system-created storage account in the primary region.

* `system_created_storage_account_type` - The storage account type for the system-created storage account in the primary region.

* `tags` - A mapping of tags assigned to the Machine Learning Registry.

---

An `identity` block exports the following:

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Machine Learning Registry.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Machine Learning Registry.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Machine Learning Registry.

* `type` - The type of Managed Service Identity that is configured on this Machine Learning Registry.

---

A `replication_regions` block exports the following:

* `location` - The Azure Region for the replication region.

* `system_created_container_registry_id` - The ID of the system-created container registry for this region.

* `system_created_container_registry_name` - The name of the system-created container registry for this region.

* `system_created_container_registry_sku` - The SKU of the system-created container registry for this region. This is always `Premium`.

* `system_created_storage_account_hierarchical_namespace_enabled` - Whether hierarchical namespace is enabled for the system-created storage account for this region.

* `system_created_storage_account_id` - The ID of the system-created storage account for this region.

* `system_created_storage_account_name` - The name of the system-created storage account for this region.

* `system_created_storage_account_type` - The storage account type for the system-created storage account for this region.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Registry.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.MachineLearningServices` - 2025-06-01
