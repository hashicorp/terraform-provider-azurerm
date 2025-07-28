---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_workspace"
description: |-
  Gets information on an existing Databricks Workspace
---

# Data Source: azurerm_databricks_workspace

Use this data source to access information about an existing Databricks workspace.

## Example Usage

```hcl
data "azurerm_databricks_workspace" "example" {
  name                = "example-workspace"
  resource_group_name = "example-rg"
}

output "databricks_workspace_id" {
  value = data.azurerm_databricks_workspace.example.workspace_id
}
```

## Argument Reference

* `name` - The name of the Databricks Workspace.
* `resource_group_name` - The Name of the Resource Group where the Databricks Workspace exists.

## Attributes Reference

* `id` - The ID of the Databricks Workspace.

* `location` - The Azure location where the Databricks Workspace exists.

* `sku` - SKU of this Databricks Workspace.

* `workspace_id` - Unique ID of this Databricks Workspace in Databricks management plane.

* `workspace_url` - URL this Databricks Workspace is accessible on.

* `managed_disk_identity` - A `managed_disk_identity` block as documented below.

* `storage_account_identity` - A `storage_account_identity` block as documented below.

* `enhanced_security_compliance` - An `enhanced_security_compliance` block as documented below.

* `custom_parameters` - A `custom_parameters` block as documented below.
* 
* `tags` - A mapping of tags to assign to the Databricks Workspace.

---

A `managed_disk_identity` block exports the following:

* `principal_id` - The principal UUID for the internal databricks disks identity needed to provide access to the workspace for enabling Customer Managed Keys.

* `tenant_id` - The UUID of the tenant where the internal databricks disks identity was created.

* `type` - The type of the internal databricks disk identity.

---

A `storage_account_identity` block exports the following:

* `principal_id` - The principal UUID for the internal databricks storage account needed to provide access to the workspace for enabling Customer Managed Keys.

* `tenant_id` - The UUID of the tenant where the internal databricks storage account was created.

* `type` - The type of the internal databricks storage account.

---

An `enhanced_security_compliance` block exports the following:

* `automatic_cluster_update_enabled` - Whether automatic cluster updates for this workspace is enabled.

* `compliance_security_profile_enabled` - Whether compliance security profile for this workspace is enabled.

* `compliance_security_profile_standards` - A list of standards enforced on this workspace.

* `enhanced_security_monitoring_enabled` - Whether enhanced security monitoring for this workspace is enabled.

---

A `custom_parameters` block exports the following:

* `machine_learning_workspace_id` - The ID of a Azure Machine Learning workspace to link with Databricks workspace.

* `nat_gateway_name` - Name of the NAT gateway for Secure Cluster Connectivity (No Public IP) workspace subnets (only for workspace with managed virtual network).

* `no_public_ip` - Are public IP Addresses not allowed?

* `private_subnet_name` - The name of the Private Subnet within the Virtual Network.

* `public_ip_name` - Name of the Public IP for No Public IP workspace with managed virtual network.

* `public_subnet_name` - The name of the Public Subnet within the Virtual Network.

* `storage_account_name` - Default Databricks File Storage account name.

* `storage_account_sku_name` - Storage account SKU name.

* `virtual_network_id` - The ID of a Virtual Network where this Databricks Cluster should be created.

* `vnet_address_prefix` - Address prefix for Managed virtual network.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Workspace.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Databricks`: 2024-05-01
