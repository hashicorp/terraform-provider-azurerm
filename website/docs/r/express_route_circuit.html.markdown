---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_circuit"
sidebar_current: "docs-azurerm-resource-network-express-route-circuit-x"
description: |-
  Manages an ExpressRoute circuit.
---

# azurerm_express_route_circuit

Manages an ExpressRoute circuit.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "exprtTest"
  location = "West US"
}

resource "azurerm_express_route_circuit" "test" {
  name                     = "expressRoute1"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  service_provider_name    = "Equinix"
  peering_location         = "Silicon Valley"
  bandwidth_in_mbps        = 50
  sku {
    tier   = "Standard"
    family = "MeteredData"
  }

  tags {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ExpressRoute circuit. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the ExpressRoute circuit. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `service_provider_name` - (Required) The name of the ExpressRoute Service Provider.

* `peering_location` - (Required) The name of the peering location and **not** the Azure resource location.

* `bandwidth_in_mbps` - (Required) The bandwidth in Mbps of the circuit being created.

~> **NOTE:** Once you increase your bandwidth, you will not be able to decrease it to it's previous value.

* `sku` - (Required) A `sku` block for the ExpressRoute circuit as documented below.

* `allow_classic_operations` - (Optional) Allow the circuit to interact with classic (RDFE) resources. The default value is `false`.

* `tags` - (Optional) A mapping of tags to assign to the resource.

---

`sku` supports the following:

* `tier` - (Required) The service tier. Possible values are `Standard` or `Premium`.

* `family` - (Required) The billing mode for bandwidth. Possible values are `MeteredData` or `UnlimitedData`.

~> **NOTE:** You can migrate from `MeteredData` to `UnlimitedData`, but not the other way around.

## Attributes Reference

The following attributes are exported:

* `id` - The Resource ID of the ExpressRoute circuit.
* `service_provider_provisioning_state` - The ExpressRoute circuit provisioning state from your chosen service provider. Possible values are "NotProvisioned", "Provisioning", "Provisioned", and "Deprovisioning".
* `service_key` - The string needed by the service provider to provision the ExpressRoute circuit.

## Import

ExpressRoute circuits can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_express_route_circuit.myExpressRoute /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/expressRouteCircuits/myExpressRoute
```
