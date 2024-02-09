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
  name = "existing"
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

* `identity_ids` - A `identity_ids` block as defined below.

* `principal_id` - The ID of the TODO.

* `tenant_id` - The ID of the TODO.

* `type` - TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Databricks Access Connector.