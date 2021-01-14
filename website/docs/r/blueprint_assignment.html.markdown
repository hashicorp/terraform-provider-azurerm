---
subcategory: "Blueprints"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_blueprint_assignment"
description: |-
  Manages a Blueprint Assignment resource
---

# azurerm_blueprint_assignment

Manages a Blueprint Assignment resource

~> **NOTE:** Azure Blueprints are in Preview and potentially subject to breaking change without notice.

~> **NOTE:** Azure Blueprint Assignments can only be applied to Subscriptions.  Assignments to Management Groups is not currently supported by the service or by Terraform.

## Example Usage
```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azurerm_subscription" "example" {}

data "azurerm_blueprint_definition" "example" {
  name     = "exampleBlueprint"
  scope_id = data.azurerm_subscription.example.id
}

data "azurerm_blueprint_published_version" "example" {
  scope_id       = data.azurerm_blueprint_definition.example.scope_id
  blueprint_name = data.azurerm_blueprint_definition.example.name
  version        = "v1.0.0"
}

resource "azurerm_resource_group" "example" {
  name     = "exampleRG-bp"
  location = "westeurope"

  tags = {
    Environment = "example"
  }
}

resource "azurerm_user_assigned_identity" "example" {
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  name                = "bp-user-example"
}

resource "azurerm_role_assignment" "operator" {
  scope                = data.azurerm_subscription.example.id
  role_definition_name = "Blueprint Operator"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_role_assignment" "owner" {
  scope                = data.azurerm_subscription.example.id
  role_definition_name = "Owner"
  principal_id         = azurerm_user_assigned_identity.example.principal_id
}

resource "azurerm_blueprint_assignment" "example" {
  name                   = "testAccBPAssignment"
  target_subscription_id = data.azurerm_subscription.example.id
  version_id             = data.azurerm_blueprint_published_version.example.id
  location               = azurerm_resource_group.example.location

  lock_mode = "AllResourcesDoNotDelete"

  lock_exclude_principals = [
    data.azurerm_client_config.current.object_id,
  ]

  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  resource_groups = <<GROUPS
    {
      "ResourceGroup": {
        "name": "exampleRG-bp"
      }
    }
  GROUPS

  parameter_values = <<VALUES
    {
      "allowedlocationsforresourcegroups_listOfAllowedLocations": {
        "value": ["westus", "westus2", "eastus", "centralus", "centraluseuap", "southcentralus", "northcentralus", "westcentralus", "eastus2", "eastus2euap", "brazilsouth", "brazilus", "northeurope", "westeurope", "eastasia", "southeastasia", "japanwest", "japaneast", "koreacentral", "koreasouth", "indiasouth", "indiawest", "indiacentral", "australiaeast", "australiasoutheast", "canadacentral", "canadaeast", "uknorth", "uksouth2", "uksouth", "ukwest", "francecentral", "francesouth", "australiacentral", "australiacentral2", "uaecentral", "uaenorth", "southafricanorth", "southafricawest", "switzerlandnorth", "switzerlandwest", "germanynorth", "germanywestcentral", "norwayeast", "norwaywest"]
      }
    }
  VALUES

  depends_on = [
    azurerm_role_assignment.operator,
    azurerm_role_assignment.owner
  ]
}

```


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

* `lock_mode` - (Optional) The locking mode of the Blueprint Assignment.  One of `None` (Default), `AllResourcesReadOnly`, or `AlResourcesDoNotDelete`.

* `lock_exclude_principals` - (Optional) a list of up to 5 Principal IDs that are permitted to bypass the locks applied by the Blueprint.

---

An `identity` block supports the following Arguments

* `type` - (Required) The Identity type for the Managed Service Identity. Currently only `UserAssigned` is supported.

* `user_assigned_identities` - (Required) a list of User Assigned Identity ID's. At least one ID is required.


## Attribute Reference

* `id` - The ID of the Blueprint Assignment

* `description` - The Description on the Blueprint

* `display_name` - The display name of the blueprint

* `blueprint_name` - The name of the blueprint assigned

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Blueprint Assignment.
* `update` - (Defaults to 30 minutes) Used when updating the Blueprint Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Blueprint Assignment.
* `delete` - (Defaults to 5 minutes) Used when deleting the Blueprint Assignment.

## Import

Azure Blueprint Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_blueprint_assignment.example "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprintAssignments/assignSimpleBlueprint"
```
