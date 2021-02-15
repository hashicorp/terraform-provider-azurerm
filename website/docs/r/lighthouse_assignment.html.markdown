---
subcategory: "Lighthouse"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lighthouse_assignment"
description: |-
    Manages a Lighthouse Assignment to a subscription or to a resource group.

---

# azurerm_lighthouse_assignment

Manages a [Lighthouse](https://docs.microsoft.com/en-us/azure/lighthouse) Assignment to a subscription, or to a resource group.

## Example Usage

```hcl

resource "azurerm_lighthouse_assignment" "example" {
  scope                    = "/subscription/00000000-0000-0000-0000-000000000000"
  lighthouse_definition_id = "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.ManagedServices/registrationDefinitions/00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) A unique UUID/GUID which identifies this lighthouse assignment- one will be generated if not specified. Changing this forces a new resource to be created.

* `scope` - (Required) The scope at which the Lighthouse Assignment applies too, such as `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333` or `/subscriptions/0b1f6471-1bf0-4dda-aec3-111122223333/resourceGroups/myGroup`. Changing this forces a new resource to be created.

* `lighthouse_definition_id` - (Required) A Fully qualified path of the lighthouse definition, such as `/subscriptions/0afefe50-734e-4610-8c82-a144aff49dea/providers/Microsoft.ManagedServices/registrationDefinitions/26c128c2-fefa-4340-9bb1-8e081c90ada2`. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

* `id` - the fully qualified ID of the Lighthouse Assignment.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Lighthouse Assignment.
* `read` - (Defaults to 5 minutes) Used when retrieving the Lighthouse Assignment.
* `delete` - (Defaults to 30 minutes) Used when deleting the Lighthouse Assignment.

## Import

Lighthouse Assignments can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lighthouse_assignment.example /subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.ManagedServices/registrationAssignments/00000000-0000-0000-0000-000000000000
```
