---
subcategory: "Cognitive Services"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_ai_foundry_project"
description: |-
  Manages an AI Foundry Project.
---

# azurerm_ai_foundry_project

Manages an AI Foundry Project.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_ai_foundry" "example" {
  name                       = "example-account"
  location                   = azurerm_resource_group.example.location
  resource_group_name        = azurerm_resource_group.example.name
  project_management_enabled = true
  custom_subdomain_name      = "example-account"
  sku_name                   = "S0"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Acceptance = "Test"
  }
}

resource "azurerm_ai_foundry_project" "example" {
  name          = "example-project"
  location      = azurerm_resource_group.example.location
  ai_foundry_id = azurerm_ai_foundry.example.id
  description   = "Project description for example-project"
  display_name  = "Example Project"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    Acceptance = "Test"
  }
}
```

-> **Note:** You must have a `azurerm_ai_foundry` resource where `project_management_enabled` is set to `true`, one `identity` block at least is set, and `custom_subdomain_name` is set when you will create `azurerm_ai_foundry_project` resources.

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the AI Services Account. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `description` - (Optional) The description of this AI Foundry Project.

* `display_name` - (Optional) The display name of this AI Foundry Project.

* `identity` - (Required) An `identity` block as defined below.

~> **Note:** You must have one managed identity block at least on your `azurerm_ai_foundry_project` resource.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

A `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this AI Services Account. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned`

* `identity_ids` - (Optional) Specifies a list of User Assigned Managed Identity IDs to be assigned to this AI Services Account.

~> **Note:** This is required when `type` is set to `UserAssigned` or `SystemAssigned, UserAssigned`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the AI Foundry Project.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the AI Foundry Project.
* `read` - (Defaults to 5 minutes) Used when retrieving the AI Foundry Project.
* `update` - (Defaults to 30 minutes) Used when updating the AI Foundry Project.
* `delete` - (Defaults to 30 minutes) Used when deleting the AI Foundry Project.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.CognitiveServices` - 2025-06-01
