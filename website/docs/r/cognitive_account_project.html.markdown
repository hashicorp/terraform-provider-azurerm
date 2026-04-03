---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_cognitive_account_project"
description: |-
  Manages a Cognitive Account Project.
---

# azurerm_cognitive_account_project

Manages a Cognitive Account Project.

~> **Note:** Cognitive Account Projects can only be created under a Cognitive Account that has `project_management_enabled = true`, `kind = "AIServices"`, a managed identity configured, and a `custom_subdomain_name` specified.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_cognitive_account" "example" {
  name                       = "example-account"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  kind                       = "AIServices"
  sku_name                   = "S0"
  project_management_enabled = true
  custom_subdomain_name      = "example-account-subdomain"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_cognitive_account_project" "example" {
  name                 = "example-project"
  cognitive_account_id = azurerm_cognitive_account.example.id
  location             = azurerm_resource_group.example.location
  description          = "Example cognitive services project"
  display_name         = "Example Project"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Environment = "test"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Cognitive Account Project. Changing this forces a new resource to be created.

* `cognitive_account_id` - (Required) The ID of the Cognitive Account where the Project should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Cognitive Account Project should exist. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below.

* `description` - (Optional) A description of the Cognitive Account Project.

* `display_name` - (Optional) The display name of the Cognitive Account Project.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Cognitive Account Project. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Cognitive Account Project.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Cognitive Account Project.

* `default` - Whether this project is the default project for the Cognitive Account.

* `endpoints` - A mapping of endpoint names to endpoint URLs for the project.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Cognitive Account Project.
* `read` - (Defaults to 5 minutes) Used when retrieving the Cognitive Account Project.
* `update` - (Defaults to 30 minutes) Used when updating the Cognitive Account Project.
* `delete` - (Defaults to 30 minutes) Used when deleting the Cognitive Account Project.

## Import

Cognitive Account Projects can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_cognitive_account_project.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.CognitiveServices/accounts/account1/projects/project1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CognitiveServices` - 2025-06-01
