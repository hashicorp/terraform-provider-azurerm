---
subcategory: "Blueprints"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_blueprint_published_version"
description: |-
  Gets information about an existing Blueprint Published Version
---

# Data Source: azurerm_blueprint_published_version

Use this data source to access information about an existing Blueprint Published Version

~> **Note:** Azure Blueprints are in Preview and potentially subject to breaking change without notice.

## Example Usage

```hcl
data "azurerm_subscription" "current" {}

data "azurerm_blueprint_published_version" "test" {
  scope_id       = data.azurerm_subscription.current.id
  blueprint_name = "exampleBluePrint"
  version        = "dev_v2.3"
}
```

## Argument Reference

* `blueprint_name` - (Required) The name of the Blueprint Definition

* `scope_id` - (Required) The ID of the Management Group / Subscription where this Blueprint Definition is stored.

* `version` - (Required) The Version name of the Published Version of the Blueprint Definition

## Attributes Reference

* `id` - The ID of the Published Version

* `type` - The type of the Blueprint

* `target_scope` - The target scope

* `display_name` - The display name of the Blueprint Published Version

* `description` - The description of the Blueprint Published Version

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Blueprint Published Version.

## API Providers
<!-- This section is generated, changes will be overwritten -->
This data source uses the following Azure API Providers:

* `Microsoft.Blueprint`: 2018-11-01-preview
