---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_access_connector"
description: |-
  Manages a Databricks Access Connector
---

# azurerm_databricks_access_connector

Manages a Databricks Access Connector

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_databricks_access_connector" "example" {
  name                = "databrickstest"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Databricks Access Connector resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Databricks Access Connector should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource has to be created. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

An `identity` block supports the following:
* `type` - (Required) The type of identity to use for this Access Connector. `SystemAssigned` is the only possible value.
* `principal_id` - (Optional) The object id of an existing principal. If not specified, a new system-assigned managed identity is created.
* `tenant_id` - (Optional) The tenant id in which the principal resides.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Databricks Access Connector in the Azure management plane.
* `identity`  - A list of `identity` blocks containing the system-assigned managed identities as defined below.

An `identity` block exports the following:
* `type` - The type of identity.
* `principal_id` - The Principal Id associated with this system-assigned managed identity.
* `tenant_id` - The Tenant Id associated with this system-assigned managed identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Databricks Access Connector.
* `update` - (Defaults to 5 minutes) Used when updating the Databricks Access Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Access Connector.
* `delete` - (Defaults to 5 minutes) Used when deleting the Databricks Access Connector.

## Import

Databricks Access Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_databricks_access_connector.connector1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Databricks/accessConnectors/connector1
```
