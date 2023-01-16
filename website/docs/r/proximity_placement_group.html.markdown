---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_proximity_placement_group"
description: |-
  Manages a proximity placement group for virtual machines, virtual machine scale sets and availability sets.

---

# azurerm_proximity_placement_group

Manages a proximity placement group for virtual machines, virtual machine scale sets and availability sets.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_proximity_placement_group" "example" {
  name                = "exampleProximityPlacementGroup"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  tags = {
    environment = "Production"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the proximity placement group. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the availability set. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `allowed_vm_sizes` - (Optional) Specifies the supported sizes of Virtual Machines that can be created in the Proximity Placement Group. Possible values are `Basic_A4`, `Basic_A1`, `Basic_A3`, `Basic_A2`, `Basic_A0`, `Standard_A8`, `Standard_A8_v2`, `Standard_A8m_v2`, `Standard_A5`, `Standard_A4`, `Standard_A4_v2`, `Standard_A4m_v2`, `Standard_A9`, `Standard_A1`, `Standard_A11`, `Standard_A1_v2`, `Standard_A10`, `Standard_A7`, `Standard_A6`, `Standard_A3`, `Standard_A2`, `Standard_A2_v2`, `Standard_A2m_v2`, `Standard_A0`, `Standard_B8ms`, `Standard_B4ms`, `Standard_B1ms`, `Standard_B1s`, `Standard_B2ms`, `Standard_B2s`, `Standard_D8_v3`, `Standard_D8s_v3`, `Standard_D5_v2`, `Standard_D4`, `Standard_D4_v3`, `Standard_D4_v2`, `Standard_D4s_v3`, `Standard_D1`, `Standard_D15_v2`, `Standard_D14`, `Standard_D14_v2`, `Standard_D11`, `Standard_D11_v2`, `Standard_D16_v3`, `Standard_D16s_v3`, `Standard_D13`, `Standard_D13_v2`, `Standard_D12`, `Standard_D12_v2`, `Standard_D1_v2`, `Standard_DS5_v2`, `Standard_DS4`, `Standard_DS4_v2`, `Standard_DS1`, `Standard_DS15_v2`, `Standard_DS14`, `Standard_DS14-8_v2`, `Standard_DS14-4_v2`, `Standard_DS14_v2`, `Standard_DS11`, `Standard_DS11_v2`, `Standard_DS13`, `Standard_DS13-4_v2`, `Standard_DS13-2_v2`, `Standard_DS13_v2`, `Standard_DS12`, `Standard_DS12_v2`, `Standard_DS1_v2`, `Standard_DS3`, `Standard_DS3_v2`, `Standard_DS2`, `Standard_DS2_v2`, `Standard_D64_v3`, `Standard_D64s_v3`, `Standard_D3`, `Standard_D32_v3`, `Standard_D32s_v3`, `Standard_D3_v2`, `Standard_D2`, `Standard_D2_v3`, `Standard_D2_v2`, `Standard_D2s_v3`, `Standard_E8_v3`, `Standard_E8s_v3`, `Standard_E4_v3`, `Standard_E4s_v3`, `Standard_E16_v3`, `Standard_E16s_v3`, `Standard_E64-16s_v3`, `Standard_E64-32s_v3`, `Standard_E64_v3`, `Standard_E64s_v3`, `Standard_E32-8s_v3`, `Standard_E32-16_v3`, `Standard_E32_v3`, `Standard_E32s_v3`, `Standard_E2_v3`, `Standard_E2s_v3`, `Standard_F8`, `Standard_F8s`, `Standard_F8s_v2`, `Standard_F4`, `Standard_F4s`, `Standard_F4s_v2`, `Standard_F1`, `Standard_F16`, `Standard_F16s`, `Standard_F16s_v2`, `Standard_F1s`, `Standard_F72s_v2`, `Standard_F64s_v2`, `Standard_F32s_v2`, `Standard_F2`, `Standard_F2s`, `Standard_F2s_v2`, `Standard_G5`, `Standard_G4`, `Standard_G1`, `Standard_GS5`, `Standard_GS5-8`, `Standard_GS5-16`, `Standard_GS4`, `Standard_GS4-8`, `Standard_GS4-4`, `Standard_GS1`, `Standard_GS3`, `Standard_GS2`, `Standard_G3`, `Standard_G2`, `Standard_H8`, `Standard_H8m`, `Standard_H16`, `Standard_H16m`, `Standard_H16mr`, `Standard_H16r`, `Standard_L8s`, `Standard_L4s`, `Standard_L16s`, `Standard_L32s`, `Standard_M128-64ms`, `Standard_M128-32ms`, `Standard_M128ms`, `Standard_M128s`, `Standard_M64-16ms`, `Standard_M64-32ms`, `Standard_M64ms`, `Standard_M64s`, `Standard_NC12`, `Standard_NC12s_v3`, `Standard_NC12s_v2`, `Standard_NC6`, `Standard_NC6s_v3`, `Standard_NC6s_v2`, `Standard_NC24`, `Standard_NC24r`, `Standard_NC24rs_v3`, `Standard_NC24rs_v2`, `Standard_NC24s_v3`, `Standard_NC24s_v2`, `Standard_ND12s`, `Standard_ND6s`, `Standard_ND24rs`, `Standard_ND24s`, `Standard_NV12`, `Standard_NV6` and `Standard_NV24`.

* `zone` - (Optional) Specifies the supported zone of the Proximity Placement Group. Changing this forces a new resource to be created.

~> **NOTE:** `allowed_vm_sizes` must be set when `zone` is specified.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Proximity Placement Group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Proximity Placement Group.
* `update` - (Defaults to 30 minutes) Used when updating the Proximity Placement Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Proximity Placement Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Proximity Placement Group.

## Import

Proximity Placement Groups can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_proximity_placement_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/example-rg/providers/Microsoft.Compute/proximityPlacementGroups/example-ppg
```
