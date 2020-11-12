---
subcategory: "Blueprints"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_blueprint_definition"
description: |-
  Gets information about an existing Blueprint Definition
---

# Data Source: azurerm_blueprint_definition

Use this data source to access information about an existing Azure Blueprint Definition

~> **NOTE:** Azure Blueprints are in Preview and potentially subject to breaking change without notice.

## Example Usage

```hcl
data "azurerm_client_config" "current" {}

data "azurerm_management_group" "root" {
  name = data.azurerm_client_config.current.tenant_id
}

data "azurerm_blueprint_definition" "example" {
  name     = "exampleManagementGroupBP"
  scope_id = data.azurerm_management_group.root.id
}

```

## Argument Reference

* `name` - (Required) The name of the Blueprint.

* `scope_id` - (Required) The ID of the Subscription or Management Group, as the scope at which the blueprint definition is stored.

## Attribute Reference

* `id` - The ID of the Blueprint Definition.  

* `description` - The description of the Blueprint Definition.  

* `display_name` - The display name of the Blueprint Definition.  

* `last_modified` - The timestamp of when this last modification was saved to the Blueprint Definition.  

* `target_scope` - The target scope.  

* `time_created` - The timestamp of when this Blueprint Definition was created.  

* `versions` - A list of versions published for this Blueprint Definition.  


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Blueprint Published Version.  
