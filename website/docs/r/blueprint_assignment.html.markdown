---
subcategory: "Blueprint"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_blueprint_assignment"
description: |-
  Manages a Blueprint Assignment.
---

# azurerm_blueprint_assignment

Manages a Blueprint Assignment.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "current" {}

resource "azurerm_blueprint_assignment" "example" {
  name     = "example-blueprint-assignment"
  location = "West Europe"
  scope    = data.azurerm_subscription.current.id

  blueprint_definition_id = "${data.azurerm_subscription.current.id}/providers/Microsoft.Blueprint/blueprints/MyBlueprint"

  identity {
    type = "SystemAssigned"
  }

  resource_groups = <<GROUPS
        {
          "prodRG": {
            "name": "example-RG-from-blueprint",
            "location": "West Europe"
          }
        }
  GROUPS

  parameter_values = <<VALUES
        {
          "tagName": {
            "value": "ENV"
          },
          "tagValue": {
            "value": "Example"
          }
        }
  VALUES
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Blueprint Assignment. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Blueprint Assignment exists. Changing this forces a new resource to be created.

* `scope` - (Required) The Scope at which the Blueprint Assignment should be applied, which must be either a Subscription (e.g. `/subscriptions/00000000-0000-0000-0000-000000000000`). Changing this forces a new resource to be created.

* `blueprint_definition_id` - (Required) The resource ID of the published version of a blueprint definition. Changing this forces a new resource to be created.

* `identity` - (Required) An `identity` block as defined below.

* `description` - (Optional) A multi-lined string that describes this Blueprint Assignment.

* `display_name` - (Optional) A one-lined string that explains the purpose of this Blueprint Assignment.

* `parameter_values` - (Optional) The parameter values for the Blueprint Assignment. This field is a json object that allows you to fill the parameters defined in your blueprint definition.

~> **NOTE:** There is a [`file` function available](https://www.terraform.io/docs/configuration/functions/file.html) which allows you to read this from an external file, which helps makes this more resource more readable.

* `resource_groups` - (Optional) The resource group values defined for the Blueprint Assignment. This field is a json object, specifying the name and location values of resource group placeholders defined in your blueprint definition.

~> **NOTE:** There is a [`file` function available](https://www.terraform.io/docs/configuration/functions/file.html) which allows you to read this from an external file, which helps makes this more resource more readable.

* `lock_mode` - (Optional) Defines how the resources that are deployed by a blueprint assignment are locked. Possible values include: `None`, `AllResourcesReadOnly`, `AllResourcesDoNotDelete`. Defaults to `None`.

* `lock_exclude_principals` - (Optional) A set of Azure AD service principal IDs that are excluded from the [deny assignment](https://docs.microsoft.com/en-us/azure/role-based-access-control/deny-assignments) the blueprint assignment creates. A total of up to 5 service principals can be defined.

---

A `identity` block supports the following:

* `type` - (Required) The type of Managed Identity which should be assigned to the Blueprint Assignment. Possible values are `SystemAssigned`, `UserAssigned`.

* `identity_ids` - (Optional) A list of User Managed Identity ID's which should be assigned to the Blueprint Assignment.

~> **NOTE:** This is required when `type` is set to `UserAssigned`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Blueprint Assignment.

* `published_blueprint_id` - The ID of the published blueprint ID which the Blueprint Assignment is using.

* `identity` - An `identity` block as documented below.

---

An `identity` block exports the following:

* `principal_id` - The ID of the System Managed Service Principal.

* `tenant_id` - The ID of the Tenant the Service Principal is assigned in.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Blueprint Assignment.
* `update` - (Defaults to 30 minutes) Used when updating the Blueprint Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Blueprint Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Blueprint Assignment.

## Import

Blueprint Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_blueprint_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Blueprint/blueprintAssignments/assignment1
```
