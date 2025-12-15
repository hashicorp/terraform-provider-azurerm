---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_network_security_perimeter_association"
description: |-
  Manages a Network Security Perimeter Association.
---

# azurerm_network_security_perimeter_association

Manages a Network Security Perimeter Association.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example"
  location            = "West Europe"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_network_security_perimeter" "example" {
  name = "example"
  resource_group_name = azurerm_resource_group.example.name
  location = "West Europe"
}

resource "azurerm_network_security_perimeter_profile" "example" {
  name = "example"
  perimeter_id = azurerm_network_security_perimeter.example.id
}

resource "azurerm_network_security_perimeter_association" "example" {
  name = "example"
  access_mode = "Enforced"

  profile_id = azurerm_network_security_perimeter_profile.example.id
  resource_id = azurerm_log_analytics_workspace.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Network Security Perimeter Association. Changing this forces a new Network Security Perimeter Association to be created.

* `access_mode` - (Required) Access mode for the associated resource on the Network Security Perimeter. Possible values are `Audit`, `Enforced`, and `Learning`.

* `profile_id` - (Required) The ID of the Network Security Perimeter Profile. Changing this forces a new Network Security Perimeter Association to be created.

* `resource_id` - (Required) The ID of the associated resource. Changing this forces a new Network Security Perimeter Association to be created.

-> **Note:** A resource can only be associated with one Network Security Perimeter at a time.  
If the target resource is already associated with another Network Security Perimeter, the creation of this association may appear to succeed in Terraform but will not be reflected in Azure. In such cases, the association will not actually exist and subsequent Terraform operations may show unexpected behavior.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Network Security Perimeter Association.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Network Security Perimeter Association.
* `read` - (Defaults to 5 minutes) Used when retrieving the Network Security Perimeter Association.
* `update` - (Defaults to 30 minutes) Used when updating the Network Security Perimeter Association.
* `delete` - (Defaults to 30 minutes) Used when deleting the Network Security Perimeter Association.

## Import

Network Security Perimeter Associations can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_network_security_perimeter_association.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Network/networkSecurityPerimeters/example-nsp/resourceAssociations/example-assoc
```