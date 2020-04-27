---
subcategory: "Blueprints"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_blueprint_published_version"
description: |-
  Gets information about an existing Blueprint Published Version
---

# Data Source: azurerm_blueprint_published_version

Use this data source to access information about an existing Azure Blueprint Published Version

## Example Usage
```hcl
data "azurerm_blueprint_published_version" "subscription_example" {
  subscription_id = "00000000-0000-0000-0000-0000000000000"
  blueprint_name  = "exampleBluePrint"
  version         = "v1.0"
}

data "azurerm_blueprint_published_version" "management_group_example" {
  management_group = "devManagementGroup"
  blueprint_name   = "exampleMGBluePrint"
  version          = "dev_v2.3"
}
```


## Argument Reference

* `subscription_id` - (Optional) The ID of the Subscription where the Blueprint is stored.  One of `subscription_id` or `management_group` is required.

* `management_group` - (Optional) The ID of the Management Group where the Blueprint is stored.  One of `management_group` or `subscription_id` is required.

* `blueprint_name` - (Required) The name of the Blueprint Definition

* `version` - (Required) The Version name of the Published Version of the Blueprint Definition


## Attribute Reference

* `id` - The Azure Resource ID of the Published Version  

* `type` - The type of the Blueprint  

* `target_scope` - The target scope  

* `display_name` - The display name of the Blueprint Published Version  

* `description` - The description of the Blueprint Published Version  


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Blueprint Published Version.
