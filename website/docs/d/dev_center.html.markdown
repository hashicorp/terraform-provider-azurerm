---
subcategory: "Dev Center"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_dev_center"
description: |-
  Gets information about an existing Dev Center.
---

# Data Source: azurerm_dev_center

Use this data source to access information about an existing Dev Center.

## Example Usage

```hcl
data "azurerm_dev_center" "example" {
  name                = "example"
  resource_group_name = "example-resource-group"
}

output "id" {
  value = data.azurerm_dev_center.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of this Dev Center.

* `resource_group_name` - (Required) The name of the Resource Group where the Dev Center exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dev Center.

* `dev_center_uri` - The URI of the Dev Center.

* `identity` - An `identity` block as defined below.

* `location` - The Azure Region where the Dev Center exists.

* `tags` - A mapping of tags assigned to the Dev Center.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Dev Center.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Dev Center.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Dev Center.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Dev Center.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Dev Center.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.DevCenter`: 2025-02-01
