---
subcategory: "Management"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_management_private_link"
description: |-
  Manages a Resource Management Private Link to restrict access for managing resources in the tenant.
---

# azurerm_resource_management_private_link

Manages a Resource Management Private Link to restrict access for managing resources in the tenant.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}
resource "azurerm_resource_management_private_link" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of this Resource Management Private Link. Changing this forces a new Resource Management Private Link to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group within which this Resource Management Private Link should exist. Changing this forces a new Resource Management Private Link to be created.
 
* `location` - (Required) The Azure Region where the Resource Management Private Link should exist. Changing this forces a new Resource Management Private Link to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Resource Management Private Link.

---

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Resource Management Private Link.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Management Private Link.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Management Private Link.

## Import

An existing Resource Management Private Link can be imported into Terraform using the `resource id`, e.g.

```shell
terraform import azurerm_resource_management_private_link.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Authorization/resourceManagementPrivateLinks/link1
```
