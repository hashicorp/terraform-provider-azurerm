---
subcategory: "Resource Connector"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_resource_connector_appliance"
description: |-
  Manages a Resource Connector Appliance.
---

# azurerm_resource_connector_appliance

Manages a Resource Connector Appliance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_resource_connector_appliance" "example" {
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

* `name` - (Required) The Name which should be used for this Resource Connector Appliance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the resource group where the Resource Connector Appliance exists. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Resource Connector Appliance should exist. Changing this forces a new resource to be created.

* `distro` - (Required) Specifies a supported Fabric/Infrastructure for this Resource Connector Appliance. The possible value is `AKSEdge`.

* `identity` - (Required) An `identity` block as defined below. Changing this forces a new resource to be created.

* `infrastructure_provider`- (Required) The infrastructure provider about the connected Resource Connector Appliance. Possible values are `HCI`,`SCVMM` and `VMWare`. Changing this forces a new resource to be created.

* `public_key`- (Optional) The `public_key` is an RSA public key in PKCS1 format encoded in base64. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Resource Connector Appliance.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Resource Connector Appliance. The only possible value is `SystemAssigned`.


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Resource Connector Appliance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 minutes) Used when creating the Resource Connector Appliance.
* `read` - (Defaults to 5 minutes) Used when retrieving the Resource Connector Appliance.
* `update` - (Defaults to 30 minutes) Used when updating the Resource Connector Appliance.
* `delete` - (Defaults to 30 minutes) Used when deleting the Resource Connector Appliance.

## Import

Resource Connector Appliance can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_resource_connector_appliance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/providers/Microsoft.ResourceConnector/appliances/appliancesExample
```
