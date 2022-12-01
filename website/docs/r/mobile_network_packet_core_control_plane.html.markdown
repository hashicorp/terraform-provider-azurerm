---
subcategory: "Mobile Network"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mobile_network_packet_core_control_plane"
description: |-
  Manages a Mobile Network Packet Core Control Plane.
---

# azurerm_mobile_network_packet_core_control_plane

Manages a Mobile Network Packet Core Control Plane.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_mobile_network" "example" {
  name                = "example-mn"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  mobile_country_code = "001"
  mobile_network_code = "01"
}

resource "azurerm_mobile_network_packet_core_control_plane" "example" {
  name                = "example-mnpccp"
  resource_group_name = azurerm_resource_group.example.name
  location            = "West Europe"
  sku                 = "EvaluationPackage"
  mobile_network_id   = azurerm_mobile_network.example.id

  control_plane_access_interface {
    name         = "default-interface"
    ipv4_address = "192.168.1.199"
    ipv4_gateway = "192.168.1.1"
    ipv4_subnet  = "192.168.1.0/25"
  }

  platform {
    type = "BaseVM"
  }

  interop_settings = jsonencode({
    "key" : "value"
  })

  tags = {
    key = "value"
  }

}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Mobile Network Packet Core Control Plane. Changing this forces a new Mobile Network Packet Core Control Plane to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Mobile Network Packet Core Control Plane should exist. Changing this forces a new Mobile Network Packet Core Control Plane to be created.

* `location` - (Required) Specifies the Azure Region where the Mobile Network Packet Core Control Plane should exist. Changing this forces a new Mobile Network Packet Core Control Plane to be created.

* `control_plane_access_interface` - (Required) A `control_plane_access_interface` block as defined below. The control plane interface on the access network. For 5G networks, this is the N2 interface. For 4G networks, this is the S1-MME interface.

* `mobile_network_id` - (Required) The ID of Mobile Network in which this packet core control plane is deployed.

* `sku` - (Required) The SKU defining the throughput and SIM allowances for this packet core control plane deployment. Possible values are `EdgeSite4GBPS`, `EdgeSite3GBPS`, `EdgeSite2GBPS`, `EvaluationPackage`, `FlagshipStarterPackage`, `LargePackage` and `MediumPackage`.

* `core_network_technology` - (Optional) The core network technology generation. Possible values are `EPG` and `5GC`.

* `platform` - (Optional) A `platform` block as defined below.

* `identity` - (Optional) An `identity` block as defined below.

* `interop_settings` - (Optional) Settings to allow interoperability with third party components e.g. RANs and UEs.

* `local_diagnostics_access_certificate_url` - (Optional) A versionless certificate URL, which used to secure local access to packet core diagnostics over local APIs by the kubernetes ingress.

* `tags` - (Optional) A mapping of tags which should be assigned to the Mobile Network Packet Core Control Plane.

* `version` - (Optional) Specifies the version of the packet core software that is deployed.

---

A `control_plane_access_interface` block supports the following:

* `name` - (Optional) Specifies the logical name for this interface. This should match one of the interfaces configured on your Azure Stack Edge device.

* `ipv4_address` - (Optional) The IPv4 address.

* `ipv4_subnet` - (Optional) The IPv4 subnet.

* `ipv4_gateway` - (Optional) The default IPv4 gateway (router).

---

An `identity` block supports the following:

* `type` - (Required) Specifies the type of Managed Service Identity. Possible values are `SystemAssigned`, `UserAssigned`, `SystemAssigned, UserAssigned` (to enable both).

* `identity_ids` - (Optional) A list of IDs for User Assigned Managed Identity resources to be assigned.

---

A `platform` block supports the following:

* `type` - (Required) Specifies the platform type where packet core is deployed. Possible values are `BaseVM` and `AKS-HCI`.

* `edge_device_id` - (Optional) The ID of Azure Stack Edge device where the packet core is deployed. If the device is part of a fault tolerant pair, either device in the pair can be specified.

~> **NOTE:** `edge_device_id` is required when `type` is set to `AKS-HCI`, and it should not be specified when `type` is set to `BaseVM`.

* `connected_cluster_id` - (Optional) The ID of Azure Arc connected cluster where the packet core is deployed.

* `custom_location_id` - (Optional) The ID of Azure Arc custom location where the packet core is deployed.

---

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Mobile Network Packet Core Control Plane.

* `identity` - An `identity` block as defined below.

---

An `identity` block exports the following:

* `principal_id` - The Principal ID associated with this Managed Service Identity.

* `tenant_id` - The Tenant ID associated with this Managed Service Identity.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Mobile Network Packet Core Control Plane.
* `read` - (Defaults to 5 minutes) Used when retrieving the Mobile Network Packet Core Control Plane.
* `update` - (Defaults to 30 minutes) Used when updating the Mobile Network Packet Core Control Plane.
* `delete` - (Defaults to 30 minutes) Used when deleting the Mobile Network Packet Core Control Plane.

## Import

Mobile Network Packet Core Control Plane can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_mobile_network_packet_core_control_plane.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.MobileNetwork/packetCoreControlPlanes/packetCoreControlPlane1
```
