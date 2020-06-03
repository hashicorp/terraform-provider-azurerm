---
subcategory: "Blueprints"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_blueprint_assignment"
description: |-
  Manages a Blueprint Assignment resource
---

# azurerm_blueprint_assignment

Manages a Blueprint Assignment resource

## Example Usage

## Argument Reference

* `name` - (Required) The name of the Blueprint Assignment

* `target_subscription_id` - (Required) The Subscription ID the Blueprint Published Version is to be applied to.

* `location` - (Required) The Azure location of the Assignment. 

* `identitiy` - (Required) an identity block, as detailed below.

* `version_id` - (Required) The ID of the Published Version of the blueprint to be assigned. 

* `parameter_values` - (Optional) a JSON string to supply Blueprint Assignment parameter values.

~> **NOTE:** Improperly formatted JSON, or missing values required by a Blueprint will cause the assignment to fail.

* `resource_groups` - (Optional) a JSON string to supply the Blueprint Resource Group information. 

~> **NOTE:** Improperly formatted JSON, or missing values required by a Blueprint will cause the assignment to fail.

* `lock_mode` - (Optional) The locking mode of the Blueprint Assignment.  One of `None`, `AllResourcesReadOnly`, or `AlResourcesDoNotDelete`.

* `lock_exclude_principals` - (Optional) a list of up to 5 Principal IDs that are permitted to bypass the locks applied by the Blueprint.

---

An `identity` block supports the following Arguments

* `type` - (Required) The Identity type for the Managed Service Identity. Currently only `UserAssigned` is supported.

* `user_assigned_identities` - (Required) a list of User Assigned Identity ID's. At least one ID is required.


## Attribute Reference

* `id` - the Azure Resource ID of the Blueprint Assignment

* `description` - The Description on the Blueprint

* `display_name` - The display name of the blueprint

* `blueprint_name` - The name of the blueprint assigned
