---
subcategory: "Azure Stack HCI"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_stack_hci_marketplace_gallery_image"
description: |-
  Manages an Azure Stack HCI Marketplace Gallery Image.
---

# azurerm_stack_hci_marketplace_gallery_image

Manages an Azure Stack HCI Marketplace Gallery Image.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "examples"
  location = "West Europe"
}

data "azurerm_client_config" "example" {}

// service principal of 'Microsoft.AzureStackHCI Resource Provider'
data "azuread_service_principal" "hciRp" {
  client_id = "1412d89f-b8a8-4111-b4fd-e82905cbd85d"
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_resource_group.example.id
  role_definition_name = "Azure Connected Machine Resource Manager"
  principal_id         = data.azuread_service_principal.hciRp.object_id
}

resource "azurerm_stack_hci_marketplace_gallery_image" "example" {
  name                = "example-mgi"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.ExtendedLocation/customLocations/cl1"
  hyperv_generation   = "V2"
  os_type             = "Windows"
  version             = "20348.2655.240905"
  identifier {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2022-datacenter-azure-edition-core"
  }
  tags = {
    foo = "bar"
    env = "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Stack HCI Marketplace Gallery Image. Changing this forces a new Azure Stack HCI Marketplace Gallery Image to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Stack HCI Marketplace Gallery Image should exist. Changing this forces a new Azure Stack HCI Marketplace Gallery Image to be created.

* `location` - (Required) The Azure Region where the Azure Stack HCI Marketplace Gallery Image should exist. Changing this forces a new Azure Stack HCI Marketplace Gallery Image to be created.

* `custom_location_id` - (Required) The ID of the Custom Location where the Azure Stack HCI Marketplace Gallery Image should exist. Changing this forces a new resource to be created.

* `hyperv_generation` - (Required) The hypervisor generation of the Azure Stack HCI Marketplace Gallery Image. Possible values are `V1` and `V2`. Changing this forces a new Azure Stack HCI Marketplace Gallery Image to be created.

* `identifier` - (Required) An `identifier` block as defined below. Changing this forces a new Azure Stack HCI Marketplace Gallery Image to be created.

* `os_type` - (Required) The Operating System type of the Azure Stack HCI Marketplace Gallery Image. Possible values are `Windows` and `Linux`. Changing this forces a new Azure Stack HCI Marketplace Gallery Image to be created.

* `version` - (Required) The version of the Azure Stack HCI Marketplace Gallery Image. Changing this forces a new Azure Stack HCI Marketplace Gallery Image to be created.

---

* `storage_path_id` - (Optional) The ID of the Azure Stack HCI Storage Path used for this Marketplace Gallery Image. Changing this forces a new Azure Stack HCI Virtual Hard Disk to be created.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Stack HCI Marketplace Gallery Image.

---

An `identifier` block supports the following:

* `offer` - (Required) The offer of the Azure Stack HCI Marketplace Gallery Image. Changing this forces a new Azure Stack HCI Marketplace Gallery Image to be created.

* `publisher` - (Required) The publisher of the Azure Stack HCI Marketplace Gallery Image. Changing this forces a new Azure Stack HCI Marketplace Gallery Image to be created.

* `sku` - (Required) The sku of the Azure Stack HCI Marketplace Gallery Image. Changing this forces a new Azure Stack HCI Marketplace Gallery Image to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Azure Stack HCI Marketplace Gallery Image.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 hours) Used when creating the Azure Stack HCI Marketplace Gallery Image.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Stack HCI Marketplace Gallery Image.
* `update` - (Defaults to 30 minutes) Used when updating the Azure Stack HCI Marketplace Gallery Image.
* `delete` - (Defaults to 1 hour) Used when deleting the Azure Stack HCI Marketplace Gallery Image.

## Import

Azure Stack HCI Marketplace Gallery Images can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_stack_hci_marketplace_gallery_image.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.AzureStackHCI/marketplaceGalleryImages/image1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.AzureStackHCI`: 2024-01-01
