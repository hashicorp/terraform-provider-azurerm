---
subcategory: "Workloads"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_workloads_sap_discovery_virtual_instance"
description: |-
  Manages an SAP Discovery Virtual Instance.
---

# azurerm_workloads_sap_discovery_virtual_instance

Manages an SAP Discovery Virtual Instance.

-> **Note:** Before using this resource, it's required to submit the request of registering the Resource Provider with Azure CLI `az provider register --namespace "Microsoft.Workloads"`. The Resource Provider can take a while to register, you can check the status by running `az provider show --namespace "Microsoft.Workloads" --query "registrationState"`. Once this outputs "Registered" the Resource Provider is available for use.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-sapvis"
  location = "West Europe"
}

resource "azurerm_workloads_sap_discovery_virtual_instance" "example" {
  name                              = "X01"
  resource_group_name               = azurerm_resource_group.example.name
  location                          = azurerm_resource_group.example.location
  environment                       = "NonProd"
  sap_product                       = "S4HANA"
  central_server_virtual_machine_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/exampleRG/providers/Microsoft.Compute/virtualMachines/csvm1"
  managed_storage_account_name      = "managedsa"

  identity {
    type = "UserAssigned"

    identity_ids = [
      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/exampleRG/providers/Microsoft.ManagedIdentity/userAssignedIdentities/uai1"
    ]
  }

  lifecycle {
    ignore_changes = [managed_resource_group_name]
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the SAP Discovery Virtual Instance. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the SAP Discovery Virtual Instance should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the SAP Discovery Virtual Instance should exist. Changing this forces a new resource to be created.

* `central_server_virtual_machine_id` - (Required) The ID of the Virtual Machine of the Central Server. Changing this forces a new resource to be created.

* `environment` - (Required) The environment type for the SAP Discovery Virtual Instance. Possible values are `NonProd` and `Prod`. Changing this forces a new resource to be created.

* `sap_product` - (Required) The SAP Product type for the SAP Discovery Virtual Instance. Possible values are `ECC`, `Other` and `S4HANA`. Changing this forces a new resource to be created.

* `identity` - (Optional) An `identity` block as defined below.

* `managed_resource_group_name` - (Optional) The name of the managed Resource Group for the SAP Discovery Virtual Instance. Changing this forces a new resource to be created.

* `managed_resources_network_access_type` - (Optional) The network access type for managed resources. Possible values are `Private` and `Public`. Defaults to `Public`.

* `managed_storage_account_name` - (Optional) The name of the custom Storage Account created by the service in the managed Resource Group. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the SAP Discovery Virtual Instance.

---

An `identity` block supports the following:

* `type` - (Required) The type of Managed Service Identity that should be configured on this SAP Discovery Virtual Instance. The only possible value is `UserAssigned`.

* `identity_ids` - (Required) A list of User Assigned Managed Identity IDs to be assigned to this SAP Discovery Virtual Instance.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the SAP Discovery Virtual Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the SAP Discovery Virtual Instance.
* `read` - (Defaults to 5 minutes) Used when retrieving the SAP Discovery Virtual Instance.
* `update` - (Defaults to 1 hour) Used when updating the SAP Discovery Virtual Instance.
* `delete` - (Defaults to 1 hour) Used when deleting the SAP Discovery Virtual Instance.

## Import

SAP Discovery Virtual Instances can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_workloads_sap_discovery_virtual_instance.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.Workloads/sapVirtualInstances/vis1
```
