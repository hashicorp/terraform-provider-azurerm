---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_data_factory_v2"
sidebar_current: "docs-azurerm-resource-data-factory-v2-x"
description: |-
  Manage an Azure Data Factory (Version 2).
---

# azurerm_data_factory_v2

Manage an Azure Data Factory (Version 2).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example"
  location = "northeurope"
}

resource "azurerm_data_factory_v2" "example" {
  name                = "example"
  location            = "northeurope"
  resource_group_name = "${azurerm_resource_group.resource_group.name}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Data Factory. Changing this forces a new resource to be created. Must be globally unique. See the [Microsoft documentation](https://docs.microsoft.com/en-us/azure/data-factory/naming-rules) for all restrictions.

* `resource_group_name` - (Required) The name of the resource group in which to create the Data Factory.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `github_configuration` - (Optional) A `github_configuration` block as defined below.

* `vsts_configuration` - (Optional) A `vsts_configuration` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`identity` supports the following:

* `type` - (Required) Specifies the identity type of the Data Factory. At this time the only allowed value is `SystemAssigned`.

---

`github_configuration` supports the following:

* `account_name` - (Required)

* `collaboration_branch` - (Required)

* `host_name` - (Required)

* `repository_name` - (Required)

* `root_folder` - (Required)

---

`vsts_configuration` supports the following:

* `account_name` - (Required)

* `collaboration_branch` - (Required)

* `project_name` - (Required)

* `repository_name` - (Required)

* `root_folder` - (Required)

* `tenant_id` - (Required)

## Attributes Reference

The following attributes are exported:

* `id` - The Data Factory ID.

## Import

Data Factory can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_data_factory_v2.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example/providers/Microsoft.DataFactory/factories/example
```
