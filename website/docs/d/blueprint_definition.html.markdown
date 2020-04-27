---
subcategory: "Blueprints"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_blueprint_definition"
description: |-
  Gets information about an existing Blueprint Definition
---

# Data Source: azurerm_blueprint_definition

Use this data source to access information about an existing Azure Blueprint Definition

## Example Usage

```hcl
data "azurerm_blueprint_definition" "example" {

}

```

## Argument Reference

* `name` - (Required) The name of the Blueprint

* `scope_type` - (Required) The scope at which the blueprint definition is stored. Possible values are `subscription` and `managementGroup`.  

* `scope_name` - (Required) The name of the scope. This is a subscription ID or Management Group name, depending on the `scope_type`.  

## Attribute Reference

* `id` - The Azure Resource ID of the Blueprint Definition.  

* `target_scope` - The target scope.  

* `display_name` - The display name of the Blueprint Definition.  

* `description` - The description of the Blueprint Definition.  

* `time_created` - The timestamp of when this Blueprint Definition was created.  

* `last_modified` - The timestamp of when this last modification was saved to the Blueprint Definition.  

* `versions` - A list of versions published for this Blueprint Definition.  


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Blueprint Published Version.  
