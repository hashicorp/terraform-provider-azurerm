---
subcategory: "Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_express_route_port"
description: |-
  Manages a Express Route Port.
---

# azurerm_express_route_port

Manages a Express Route Port.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West US"
}

resource "azurerm_express_route_port" "example" {
  name                = "port1"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  peering_location    = "Airtel-Chennai-CLS"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Express Route Port. Changing this forces a new Express Route Port to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Express Route Port should exist. Changing this forces a new Express Route Port to be created.
  
* `location` - (Required) The Azure Region where the Express Route Port should exist. Changing this forces a new Express Route Port to be created.
  
* `bandwidth_in_gbps` - (Required) Bandwidth of the Express Route Port in Gbps. Changing this forces a new Express Route Port to be created.

* `encapsulation` - (Required) The encapsulation method used for the Express Route Port. Changing this forces a new Express Route Port to be created. Possible values are: `Dot1Q`, `QinQ`.

* `peering_location` - (Required) The name of the peering location that this Express Route Port is physically mapped to. Changing this forces a new Express Route Port to be created.

* `link1` - (Optional) A list of `link` blocks as defined below.

* `link2` - (Optional) A list of `link` blocks as defined below.

* `billing_type` - (Optional) The billing type of the Express Route Port. Possible values are `MeteredData` and `UnlimitedData`. Defaults to `MeteredData`.

---

* `identity` - (Optional) An `identity` block as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Express Route Port.

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity that should be configured on this Express Route Port. Only possible value is `UserAssigned`.

* `identity_ids` - (Required) Specifies a list of User Assigned Managed Identity IDs to be assigned to this Express Route Port.

---

A `link` block supports the following:

* `admin_enabled` - (Optional) Whether enable administration state on the Express Route Port Link? Defaults to `false`.
  
* `macsec_cipher` - (Optional) The MACSec cipher used for this Express Route Port Link. Possible values are `GcmAes128` and `GcmAes256`. Defaults to `GcmAes128`.

* `macsec_ckn_keyvault_secret_id` - (Optional) The ID of the Key Vault Secret that contains the MACSec CKN key for this Express Route Port Link.

* `macsec_cak_keyvault_secret_id` - (Optional) The ID of the Key Vault Secret that contains the Mac security CAK key for this Express Route Port Link.

* `macsec_sci_enabled` - (Optional) Should Secure Channel Identifier on the Express Route Port Link be enabled? Defaults to `false`.

~> **Note:** `macsec_ckn_keyvault_secret_id` and `macsec_cak_keyvault_secret_id` should be used together with `identity`, so that the Express Route Port instance have the right permission to access the Key Vault.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Express Route Port.

* `identity` - A `identity` block as defined below.
  
* `link1` - A list of `link` blocks as defined below.

* `link2` - A list of `link` blocks as defined below.

* `guid` - The resource GUID of the Express Route Port.
  
* `ethertype` - The EtherType of the Express Route Port.
  
* `mtu` - The maximum transmission unit of the Express Route Port.

---

A `link` block exports the following:

* `id` - The ID of this Express Route Port Link.
  
* `router_name` - The name of the Azure router associated with the Express Route Port Link.

* `interface_name` - The interface name of the Azure router associated with the Express Route Port Link.

* `patch_panel_id` - The ID that maps from the Express Route Port Link to the patch panel port.
  
* `rack_id` - The ID that maps from the patch panel port to the rack.

* `connector_type` - The connector type of the Express Route Port Link.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Express Route Port.
* `read` - (Defaults to 5 minutes) Used when retrieving the Express Route Port.
* `update` - (Defaults to 30 minutes) Used when updating the Express Route Port.
* `delete` - (Defaults to 30 minutes) Used when deleting the Express Route Port.

## Import

Express Route Ports can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_express_route_port.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Network/expressRoutePorts/port1
```
