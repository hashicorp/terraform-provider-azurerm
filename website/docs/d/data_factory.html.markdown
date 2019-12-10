---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory"
sidebar_current: "docs-azurerm-datasource-data-factory-x"
description: |-
  Manages an Azure Data Factory (Version 2).
---

# Data Source: azurerm_data_factory

Use this data source to access information about an existing Azure Data Factory (Version 2).

## Example Usage

```hcl

data "azurerm_data_factory" "example" {
  name                = "${azurerm_data_factory.example.name}"
  resource_group_name = "${azurerm_data_factory.example.resource_group_name}"
}

output "data_factory_id" {
  value = "${azurerm_data_factory.example.id}"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory to retrieve information about. 

* `resource_group_name` - (Required) The name of the resource group where the Data Factory exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Data Factory ID.

* `location` - The Azure location where the resource exists.

* `github_configuration` - A `github_configuration` block as defined below.

* `identity` - An `identity` block as defined below.

* `vsts_configuration` - A `vsts_configuration` block as defined below.

* `tags` - A mapping of tags assigned to the resource.
---

A `github_configuration` block exports the following:

* `account_name` - The GitHub account name.

* `branch_name` - The branch of the repository to get code from.

* `git_url` - The GitHub Enterprise host name. 

* `repository_name` - The name of the git repository.

* `root_folder` - The root folder within the repository.

---

An `identity` block exports the following:

* `principal_id` - The ID of the Principal (Client) in Azure Active Directory.

* `tenant_id` - The ID of the Azure Active Directory Tenant.

* `type` - The identity type of the Data Factory.

---

A `vsts_configuration` block exports the following:

* `account_name` - The VSTS account name.

* `branch_name` - The branch of the repository to get code from.

* `project_name` - The name of the VSTS project.

* `repository_name` - The name of the git repository.

* `root_folder` - The root folder within the repository.

* `tenant_id` - The Tenant ID associated with the VSTS account.
