---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_databricks_access_connector"
description: |-
  Manages a Databricks Access Connector
---

# azurerm_databricks_access_connector

~> **NOTE:** Databricks Access Connectors are in Private Preview and potentially subject to breaking change without notice. If you would like to use these features please contact your Microsoft support representative on how to opt-in to the Databricks Access Connector Private Preview feature program.

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

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Databricks Access Connector in the Azure management plane.

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
