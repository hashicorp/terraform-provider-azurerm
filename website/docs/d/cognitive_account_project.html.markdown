---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: Data Source: azurerm_cognitive_account_project"
description: |-
  Gets information about an existing Cognitive Services Account Project.
---

# Data Source: azurerm_cognitive_account_project

Use this data source to access information about an existing Cognitive Services Account Project.

## Example Usage

```hcl
data "azurerm_cognitive_account_project" "example" {
  name                   = "example-project"
  cognitive_account_name = "example-account"
  resource_group_name    = "example-resources"
}

output "id" {
  value = data.azurerm_cognitive_account_project.example.id
}

output "location" {
  value = data.azurerm_cognitive_account_project.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Services Account Project.

* `cognitive_account_name` - (Required) The name of the Cognitive Services Account in which the Project exists.

* `resource_group_name` - (Required) The name of the Resource Group where the Cognitive Services Account exists.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Services Account Project.

* `location` - The Azure Region where the Cognitive Services Account Project exists.

* `default` - Whether this is the default project for the Cognitive Services Account.

* `description` - The description of the Cognitive Services Account Project.

* `display_name` - The display name of the Cognitive Services Account Project.

* `endpoints` - A mapping of endpoint names to endpoint URLs for the Cognitive Services Account Project.

* `identity` - An `identity` block as defined below.

* `tags` - A mapping of tags assigned to the Cognitive Services Account Project.

---

An `identity` block exports the following:

* `type` - The type of Managed Service Identity that is configured on this Cognitive Services Account Project.

* `principal_id` - The Principal ID of the System Assigned Managed Service Identity that is configured on this Cognitive Services Account Project.

* `tenant_id` - The Tenant ID of the System Assigned Managed Service Identity that is configured on this Cognitive Services Account Project.

* `identity_ids` - The list of User Assigned Managed Identity IDs assigned to this Cognitive Services Account Project.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Services Account Project.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.CognitiveServices` - 2025-06-01
