---
subcategory: "Synapse"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_synapse_workspace"
description: |-
  Gets information about an existing Synapse Workspace.
---

# Data Source: azurerm_synapse_workspace

Use this data source to access information about an existing Synapse Workspace.

## Example Usage

```hcl
data "azurerm_synapse_workspace" "example" {
  name                = "existing"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_synapse_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Synapse Workspace.

* `resource_group_name` - (Required) The name of the Resource Group where the Synapse Workspace exists.

## Attributes Reference

the following Attributes are exported:

* `id` - The ID of the synapse Workspace.

* `location` - The Azure location where the Synapse Workspace exists.

* `storage_data_lake_gen2_filesystem_id` - The ID of the Data Lake Storage Gen2 Filesystem associated with this Synapse Workspace.

* `sql_administrator_login` - The login name of the SQL administrator.

* `linking_allowed_for_aad_tenant_ids` - The allowed Azure AD tenant IDs for linking.

* `compute_subnet_id` - The ID of the Subnet used for computes in this Synapse Workspace.

* `data_exfiltration_protection_enabled` - Whether data exfiltration protection is enabled in this Synapse Workspace.

* `managed_virtual_network_enabled` - Whether managed virtual network is enabled for all computes in this Synapse Workspace.

* `connectivity_endpoints` - A map of Connectivity endpoints for this Synapse Workspace.

* `managed_resource_group_name` - The managed Resource Group name of this Synapse Workspace.

* `azure_devops_repo` - An `azure_devops_repo` block as defined below.

* `github_repo` - A `github_repo` block as defined below.

* `public_network_access_enabled` - Whether public network access is enabled for this Synapse Workspace.

* `purview_id` - The ID of the Purview Account associated with this Synapse Workspace.

* `sql_identity_control_enabled` - Whether pipelines running as the workspace's system assigned identity are allowed to access SQL pools.

* `customer_managed_key` - A `customer_managed_key` block as defined below.

* `azuread_authentication_only` - Whether Azure AD authentication is the only way to authenticate with resources inside this Synapse Workspace.

* `tags` - A mapping of tags assigned to the resource.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this Synapse Workspace.

---

An `azure_devops_repo` block exports the following:

* `account_name` - The Azure DevOps account name.

* `branch_name` - The collaboration branch of the repository.

* `last_commit_id` - The last commit ID.

* `project_name` - The name of the Azure DevOps project.

* `repository_name` - The name of the Git repository.

* `root_folder` - The root folder within the repository.

* `tenant_id` - The ID of the tenant for the Azure DevOps account.

---

A `customer_managed_key` block exports the following:

* `key_versionless_id` - The Azure Key Vault Key Versionless ID used as the Customer Managed Key.

* `key_name` - The identifier for the key.

* `user_assigned_identity_id` - The User Assigned Identity ID used for accessing the Customer Managed Key.

---

The `identity` block exports the following:

* `type` - The Identity Type for the Service Principal associated with the Managed Service Identity of this Synapse Workspace.

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Synapse Workspace.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Synapse Workspace.

---

A `github_repo` block exports the following:

* `account_name` - The GitHub account name.

* `branch_name` - The collaboration branch of the repository.

* `git_url` - The GitHub Enterprise host name.

* `last_commit_id` - The last commit ID.

* `repository_name` - The name of the Git repository.

* `root_folder` - The root folder within the repository.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Synapse Workspace.
