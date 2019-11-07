---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_blueprint"
sidebar_current: "docs-azurerm-resource-blueprint"
description: |-
  Manages an Azure Blueprint definition.

---

# azurerm_blueprint

Manages an Azure Blueprint definition.

## Example Usage

```hcl
data "azurerm_subscription" "current" {}

resource "azurerm_blueprint" "example" {
  name  = "example-blueprint"
  scope = data.azurerm_subscription.current.id
  type  = "Microsoft.Blueprint/blueprints"
  properties {
    description  = "example blueprint description"
    display_name = "example blueprint"
    target_scope = "subscription"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Blueprint

* `scope` - (Required) Specifies the scope at which the blueprint is stored. This is either a Subscription ID or a Management Group ID.

~> **NOTE:** For Subscriptions use standard `/subscriptions/00000000-0000-0000-0000-000000000000` format, for Management Groups use `/providers/Microsoft.Management/managementGroups/000000000-0000-0000-0000-000000000000`.  The GUID in this case is the Management Group GUID, or in the case of the root management group, this is the Tenant ID.

* `type` - (Required) This is (currently) always `Microsoft.Blueprint/blueprints`

* `properties` - (Required) A `properties` block as defined below

`properties` supports the following:

* `description` - (Optional) A description of this Blueprint, supports multi-line strings

* `display_name` - (Optional) A display name as viewed in the Portal

* `target_scope` - (Required) Currently only `subscription` is supported, `managementGroup` is reserved for future use and may be supported later.

* `parameters` - (Optional) JSON string of parameters for the blueprint configuration items (Policies, Templates etc)

* `resource_groups` - (Optional) JSON String of resource group definitions

* `versions` - (Optional) JSON object describing the versions of this blueprint


## Attribute Reference

* `id` - The resource ID of the Blueprint

~> **NOTE:** This is of the same form as `scope` above

* `status` - provides details on `created` and `last_modified` time and date
