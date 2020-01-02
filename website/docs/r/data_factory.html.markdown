---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory"
sidebar_current: "docs-azurerm-resource-data-factory-x"
description: |-
  Manages an Azure Data Factory (Version 2).
---

# azurerm_data_factory

Manages an Azure Data Factory (Version 2).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "northeurope"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `github_configuration` - (Optional) A `github_configuration` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `vsts_configuration` - (Optional) A `vsts_configuration` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `github_configuration` block supports the following:

* `account_name` - (Required) Specifies the GitHub account name.

* `branch_name` - (Required) Specifies the branch of the repository to get code from.

* `git_url` - (Required) Specifies the GitHub Enterprise host name. For example: https://github.mydomain.com. Use https://github.com for open source repositories.

* `repository_name` - (Required) Specifies the name of the git repository.

* `root_folder` - (Required) Specifies the root folder within the repository. Set to `/` for the top level.

-> **Note:** You must log in to the Data Factory management UI to complete the authentication to the GitHub repository.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the identity type of the Data Factory. At this time the only allowed value is `SystemAssigned`.

---

A `vsts_configuration` block supports the following:

* `account_name` - (Required) Specifies the VSTS account name.

* `branch_name` - (Required) Specifies the branch of the repository to get code from.

* `project_name` - (Required) Specifies the name of the VSTS project.

* `repository_name` - (Required) Specifies the name of the git repository.

* `root_folder` - (Required) Specifies the root folder within the repository. Set to `/` for the top level.

* `tenant_id` - (Required) Specifies the Tenant ID associated with the VSTS account.

## Attributes Reference

The following attributes are exported:

* `id` - The Data Factory ID.

* `identity` - An `identity` block as defined below.

---

The `identity` block exports the following:

* `principal_id` - The ID of the Principal (Client) in Azure Active Directory

* `tenant_id` - The ID of the Azure Active Directory Tenant.


## Import

Data Factory can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example
```
