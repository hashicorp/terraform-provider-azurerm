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
  name                = "example-mlregistry"
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

* `location` - The Azure Region where the Machine Learning Registry exists.

* `public_network_access_enabled` - Whether public network access is enabled for the Machine Learning Registry.

* `main_region` - A `main_region` block as defined below.

* `replication_region` - One or more `replication_region` blocks as defined below.

* `identity` - An `identity` block as defined below.

* `discovery_url` - The discovery URL for the Machine Learning Registry.

* `intellectual_property_publisher` - The intellectual property publisher for the Machine Learning Registry.

* `ml_flow_registry_uri` - The ML Flow registry URI for the Machine Learning Registry.

* `managed_resource_group` - The ID of the managed resource group created for the Machine Learning Registry.

* `tags` - A mapping of tags assigned to the Machine Learning Registry.

---

A `main_region` block exports the following:

* `location` - The Azure Region for the main region.

* `storage_account_type` - The storage account type for the main region.

* `hns_enabled` - Whether hierarchical namespace is enabled for the main region storage account.

* `system_created_storage_account_id` - The ID of the system-created storage account for the main region.

* `system_created_container_registry_id` - The ID of the system-created container registry for the main region.

---

A `replication_region` block exports the following:

* `location` - The Azure Region for the replication region.

* `storage_account_type` - The storage account type for the replication region.

* `hns_enabled` - Whether hierarchical namespace is enabled for the replication region storage account.

* `system_created_storage_account_id` - The ID of the system-created storage account for the replication region.

* `system_created_container_registry_id` - The ID of the system-created container registry for the replication region.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Machine Learning Registry.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Machine Learning Registry.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Machine Learning Registry.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Machine Learning Registry.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Machine Learning Registry.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.MachineLearningServices` - 2025-06-01