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
  name                = "example-resource"
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

In addition to the Arguments listed above - the following Attributes are exported:

* `name` - (Required) Specifies the name of the Databricks Access Connector resource. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Databricks Access Connector should exist. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource has to be created. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on the Databricks Access Connector. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned`.

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to the Databricks Access Connector. Only one User Assigned Managed Identity ID is supported per Databricks Access Connector resource.

~> **Note:** `identity_ids` are required when `type` is set to `UserAssigned`.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Databricks Access Connector in the Azure management plane.

* `identity` - A list of `identity` blocks containing the system-assigned managed identities as defined below.

---

An `identity` block exports the following:

* `type` - (Required) The type of Managed Service Identity that is configured on this Access Connector.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Access Connector.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Access Connector.

* `identity_ids` - (Optional) The list of User Assigned Managed Identity IDs assigned to this Access Connector. 


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Databricks Access Connector.
* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Access Connector.
* `update` - (Defaults to 30 minutes) Used when updating the Databricks Access Connector.
* `delete` - (Defaults to 30 minutes) Used when deleting the Databricks Access Connector.

## Import

Databricks Access Connectors can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_databricks_access_connector.connector1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Databricks/accessConnectors/connector1
```
