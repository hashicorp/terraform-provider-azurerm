---
subcategory: "Lab Service"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_lab_service_user"
description: |-
  Manages a Lab Service User.
---

# azurerm_lab_service_user

Manages a Lab Service User.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_lab_service_lab" "example" {
  name                = "example-lsl"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_lab_service_user" "example" {
  name                = "example-lsu"
  lab_services_lab_id = azurerm_lab_service_lab.test.id
  email               = "terraform-acctest@hashicorp.com"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Lab Service User. Changing this forces a new resource to be created.

* `lab_id` - (Required) The ID of the Lab Service Lab. Changing this forces a new resource to be created.

* `email` - (Required) The email address of the user. Changing this forces a new resource to be created.

* `additional_usage_quota` - (Optional) The amount of usage quota time the user gets in addition to the lab usage quota. Defaults to `PT0S`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Lab Services User.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Lab Service User.
* `read` - (Defaults to 5 minutes) Used when retrieving the Lab Service User.
* `update` - (Defaults to 30 minutes) Used when updating the Lab Service User.
* `delete` - (Defaults to 30 minutes) Used when deleting the Lab Service User.

## Import

Lab Service Users can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_lab_service_user.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.LabServices/labs/lab1/users/user1
```
