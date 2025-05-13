---
subcategory: "Arc Resource Bridge"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_arc_resource_bridge_appliance"
description: |-
  Manages an Arc Resource Bridge Appliance.
---

# azurerm_arc_resource_bridge_appliance

Manages an Arc Resource Bridge Appliance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_arc_resource_bridge_appliance" "example" {
  name                    = "example-appliance"
  location                = azurerm_resource_group.example.location
  resource_group_name     = azurerm_resource_group.example.name
  distro                  = "AKSEdge"
  infrastructure_provider = "VMWare"

  identity {
    type = "SystemAssigned"
  }

  tags = {
    "hello" = "world"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The Name which should be used for this Arc Resource Bridge Appliance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the Arc Resource Bridge Appliance exists. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Arc Resource Bridge Appliance should exist. Changing this forces a new resource to be created.

* `distro` - (Required) Specifies a supported Fabric/Infrastructure for this Arc Resource Bridge Appliance. The possible value is `AKSEdge`.

* `identity` - (Required) An `identity` block as defined below. Changing this forces a new resource to be created.

* `infrastructure_provider` - (Required) The infrastructure provider about the connected Arc Resource Bridge Appliance. Possible values are `HCI`,`SCVMM` and `VMWare`. Changing this forces a new resource to be created.

* `public_key_base64` - (Optional) The `public_key_base64` is an RSA public key in PKCS1 format encoded in base64. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Arc Resource Bridge Appliance.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Arc Resource Bridge Appliance. The only possible value is `SystemAssigned`. Changing this forces a new resource to be created.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Arc Resource Bridge Appliance.

* `identity` - An `identity` block as defined below.

---
An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Arc Resource Bridge Appliance.
* `read` - (Defaults to 5 minutes) Used when retrieving the Arc Resource Bridge Appliance.
* `update` - (Defaults to 30 minutes) Used when updating the Arc Resource Bridge Appliance.
* `delete` - (Defaults to 30 minutes) Used when deleting the Arc Resource Bridge Appliance.

## Import

Arc Resource Bridge Appliance can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_arc_resource_bridge_appliance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.ResourceConnector/appliances/appliancesExample
```
