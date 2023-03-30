---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_virtual_instance"
description: |-
  Manages a SAP Virtual Instance.
---

# azurerm_workloads_sap_virtual_instance

Manages a SAP Virtual Instance.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_workloads_sap_virtual_instance" "example" {
  name                = "X00"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  environment         = "NonProd"
  sap_product         = "S4HANA"

  configuration {
    central_server_vm_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this SAP Virtual Instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SAP Virtual Instance should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the SAP Virtual Instance should exist. Changing this forces a new resource to be created.

* `environment` - (Required) The environment type for the SAP Virtual Instance. Possible values are `NonProd` and `Prod`. Changing this forces a new resource to be created.

* `sap_product` - (Required) The SAP Product type for the SAP Virtual Instance. Possible values are `ECC`, `Other` and `S4HANA`. Changing this forces a new resource to be created.

* `discovery_configuration` - (Optional) A `discovery_configuration` block as defined below. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `managed_resource_group_name` - (Optional) The name of the managed Resource Group for the SAP Virtual Instance.

* `tags` - (Optional) A mapping of tags which should be assigned to the SAP Virtual Instance.

---

A `discovery_configuration` block supports the following:

* `central_server_vm_id` - (Optional) The resource ID of the Virtual Machine of the Central Server. Changing this forces a new resource to be created.

* `managed_storage_account_name` - (Optional) The name of the custom Storage Account created by the service in the managed Resource Group. Changing this forces a new resource to be created.

~> **Note:** If not provided, the service will create the Storage Account with a random name.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity that should be configured on this SAP Virtual Instance. Only possible value is `UserAssigned`.

* `identity_ids` - (Required) A list of User Assigned Managed Identity IDs to be assigned to this SAP Virtual Instance.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SAP Virtual Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the SAP Virtual Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the SAP Virtual Instance.
* `update` - (Defaults to 30 minutes) Used when updating the SAP Virtual Instance.
* `delete` - (Defaults to 30 minutes) Used when deleting the SAP Virtual Instance.

## Import

SAP Virtual Instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_workloads_sap_virtual_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Workloads/sapVirtualInstances/vis1
```
