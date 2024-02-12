---
subcategory: "Databricks"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_databricks_access_connector"
description: |-
  Gets information about an existing Databricks Access Connector.
---

# Data Source: azurerm_databricks_access_connector

Use this data source to access information about an existing Databricks Access Connector.

## Example Usage

```hcl
data "azurerm_databricks_access_connector" "example" {
  name                = "existing"
  resource_group_name = "existing"
}

output "id" {
  value = data.azurerm_databricks_access_connector.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Databricks Access Connector.

* `resource_group_name` - (Required) The name of the Resource Group where the Databricks Access Connector exists. Changing this forces a new Databricks Access Connector to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Databricks Access Connector.

* `identity` - A `identity` block as defined below.

* `location` - The Azure Region where the Databricks Access Connector exists.

* `tags` - A mapping of tags assigned to the Databricks Access Connector.

---

A `identity` block exports the following:

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Access Connector.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Access Connector.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Access Connector.

* `type` - The type of Managed Service Identity that is configured on this Access Connector.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Access Connector.
