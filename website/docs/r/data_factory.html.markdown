---
subcategory: "Data Factory"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory"
description: |-
  Manages an Azure Data Factory (Version 2).
---

# azurerm_data_factory

Manages an Azure Data Factory (Version 2).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_data_factory" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `github_configuration` - (Optional) A `github_configuration` block as defined below.

* `global_parameter` - (Optional) A list of `global_parameter` blocks as defined above.

* `identity` - (Optional) An `identity` block as defined below.

* `vsts_configuration` - (Optional) A `vsts_configuration` block as defined below.

* `managed_virtual_network_enabled` - (Optional) Is Managed Virtual Network enabled?

* `public_network_enabled` - (Optional) Is the Data Factory visible to the public network? Defaults to `true`.

* `customer_managed_key_id` - (Optional) Specifies the Azure Key Vault Key ID to be used as the Customer Managed Key (CMK) for double encryption. Required with user assigned identity.

* `customer_managed_key_identity_id` - (Optional) Specifies the ID of the user assigned identity associated with the Customer Managed Key. Must be supplied if `customer_managed_key_id` is set.

* `purview_id` - (Optional) Specifies the ID of the purview account resource associated with the Data Factory.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `github_configuration` block supports the following:

* `account_name` - (Required) Specifies the GitHub account name.

* `branch_name` - (Required) Specifies the branch of the repository to get code from.

* `git_url` - (Optional) Specifies the GitHub Enterprise host name. For example: <https://github.mydomain.com>. Use <https://github.com> for open source repositories.

* `repository_name` - (Required) Specifies the name of the git repository.

* `root_folder` - (Required) Specifies the root folder within the repository. Set to `/` for the top level.

* `publishing_enabled` - (Optional) Is automated publishing enabled? Defaults to `true`.

-> **Note:** You must log in to the Data Factory management UI to complete the authentication to the GitHub repository.

---

A `global_parameter` block supports the following:

* `name` - (Required) Specifies the global parameter name.

* `type` - (Required) Specifies the global parameter type. Possible Values are `Array`, `Bool`, `Float`, `Int`, `Object` or `String`.

* `value` - (Required) Specifies the global parameter value.

-> **Note:** For type `Array` and `Object` it is recommended to use `jsonencode()` for the value

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Data Factory. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Data Factory.

~> **NOTE:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

---

A `vsts_configuration` block supports the following:

* `account_name` - (Required) Specifies the VSTS account name.

* `branch_name` - (Required) Specifies the branch of the repository to get code from.

* `project_name` - (Required) Specifies the name of the VSTS project.

* `repository_name` - (Required) Specifies the name of the git repository.

* `root_folder` - (Required) Specifies the root folder within the repository. Set to `/` for the top level.

* `tenant_id` - (Required) Specifies the Tenant ID associated with the VSTS account.

* `publishing_enabled` - (Optional) Is automated publishing enabled? Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Data Factory.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Data Factory.
* `update` - (Defaults to 30 minutes) Used when updating the Data Factory.
* `read` - (Defaults to 5 minutes) Used when retrieving the Data Factory.
* `delete` - (Defaults to 30 minutes) Used when deleting the Data Factory.

## Import

Data Factory can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example
```
