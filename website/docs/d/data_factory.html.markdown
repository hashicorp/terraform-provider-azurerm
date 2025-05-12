---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_data_factory"
description: |-
  Gets information about an existing Azure Data Factory (Version 2).
---

# Data Source: azurerm_data_factory

Use this data source to access information about an existing Azure Data Factory (Version 2).

## Example Usage

```hcl
data "azurerm_data_factory" "example" {
  name                = "existing-adf"
  resource_group_name = "existing-rg"
}

output "id" {
  value = data.azurerm_data_factory.example.id
}
```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of this Azure Data Factory.

- `resource_group_name` - (Required) The name of the Resource Group where the Azure Data Factory exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the Azure Data Factory.

- `github_configuration` - A `github_configuration` block as defined below.

- `identity` - An `identity` block as defined below.

- `location` - The Azure Region where the Azure Data Factory exists.

- `purview_id` - (Optional) The ID of purview account.

- `tags` - A mapping of tags assigned to the Azure Data Factory.

- `vsts_configuration` - A `vsts_configuration` block as defined below.

---

A `github_configuration` block exports the following:

- `account_name` - The GitHub account name.

- `branch_name` - The branch of the repository to get code from.

- `git_url` - The GitHub repository url.

- `repository_name` - The name of the git repository.

- `root_folder` - The root folder within the repository.

---

An `identity` block exports the following:

- `type` - The type of Managed Service Identity that is configured on this Data Factory.

- `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Data Factory.

- `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Data Factory.

- `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Data Factory.

---

A `vsts_configuration` block exports the following:

- `account_name` - The VSTS account name.

- `branch_name` - The branch of the repository to get code from.

- `project_name` - The name of the VSTS project.

- `repository_name` - The name of the git repository.

- `root_folder` - The root folder within the repository.

- `tenant_id` - The Tenant ID associated with the VSTS account.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Data Factory.
