---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_storage_path"
description: |-
  Manages an Azure Stack HCI Storage Path.
---

# azurerm_stack_hci_storage_path

Manages an Azure Stack HCI Storage Path.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "West Europe"
}
resource "azurerm_stack_hci_storage_path" "example" {
  name                = "example-sp"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  path                = "C:\\ClusterStorage\\UserStorage_2\\sp-example"
  tags = {
    foo = "bar"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Stack HCI Storage Path. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Stack HCI Storage Path should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure Stack HCI Storage Path should exist. Changing this forces a new resource to be created.

* `custom_location_id` - (Required) The ID of Custom Location where the Azure Stack HCI Storage Path should exist. Changing this forces a new resource to be created.

* `path` - (Required) The file path on the disk to create the Storage Path. Changing this forces a new resource to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Stack HCI Storage Path.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The resource ID of the Azure Stack HCI Storage Path.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Azure Stack HCI Storage Path.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Storage Path.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Stack HCI Storage Path.
* `delete` - (Defaults to 30 minutes) Used when deleting the Azure Stack HCI Storage Path.

## Import

Azure Stack HCI Storage Paths can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_storage_path.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.AzureStackHCI/storageContainers/storage1
```
