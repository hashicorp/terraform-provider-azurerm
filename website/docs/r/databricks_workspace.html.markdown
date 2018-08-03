---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_workspace"
sidebar_current: "docs-azurerm-resource-databricks-workspace"
description: |-
  Manages a new Databricks Workspace resource
---

# azurerm_databricks_workspace

Manages a new Databricks Workspace

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "resourceGroup1"
  location = "West US"
}

resource "azurerm_databricks_workspace" "test" {
  name               = "databricks-test"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku = "Standard"

  tags {
    pricing = "Premium"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Databricks Workspace resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group under which the Databricks Workspace resource has to be created. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource has to be created. Changing this forces a new resource to be created.

* `sku` - (Required) A `sku` for the Databricks Workspace. Possible values are `Standard` or `Premium`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The Databricks Workspace ID.

* `managed_resource_group_id` - The Managed Resource Group ID of the created Databricks Workspace.
