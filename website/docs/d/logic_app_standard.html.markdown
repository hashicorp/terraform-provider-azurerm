---
subcategory: "Logic App"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_logic_app_standard"
description: |-
  Gets information about an existing Logic App Standard instance.
---

# Data Source: azurerm_logic_app_standard

Use this data source to access information about an existing Logic App Standard instance.

## Example Usage

```hcl
data "azurerm_logic_app_standard" "example" {
  name                = "logicappstd"
  resource_group_name = "example-rg"
}

output "id" {
  value = data.azurerm_logic_app_standard.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - The name of this Logic App.

* `resource_group_name` - The name of the Resource Group where the Logic App exists.

## Attributes Reference

The following attributes are exported:

* `id` - The Logic App Standard ID.

* `location` - The Azure location where the Logic App Standard exists.

* `identity` - An `identity` block as defined below.

---

The `identity` block exports the following:

* `type` - The Type of Managed Identity assigned to this Logic App Workflow.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this Logic App Workflow.

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this Logic App Workflow.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Logic App Workflow.
