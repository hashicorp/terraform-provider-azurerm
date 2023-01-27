---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_disk_sas_token"
description: |-
  Manages a Disk SAS Token.
---

# azurerm_managed_disk_sas_token

Manages a Disk SAS Token.

Use this resource to obtain a Shared Access Signature (SAS Token) for an existing Managed Disk.

Shared access signatures allow fine-grained, ephemeral access control to various aspects of Managed Disk similar to blob/storage account container.

With the help of this resource, data from the disk can be copied from managed disk to a storage blob or to some other system without the need of azcopy.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "testrg"
  location = "West Europe"
}

resource "azurerm_managed_disk" "test" {
  name                 = "tst-disk-export"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"
}

resource "azurerm_managed_disk_sas_token" "test" {
  managed_disk_id     = azurerm_managed_disk.test.id
  duration_in_seconds = 300
  access_level        = "Read"
}
```

## Arguments Reference

The following arguments are supported:

* `managed_disk_id` - (Required) The ID of an existing Managed Disk which should be exported. Changing this forces a new resource to be created.

* `duration_in_seconds` - (Required) The duration for which the export should be allowed. Should be between 30 & 4294967295 seconds. Changing this forces a new resource to be created.

* `access_level` - (Required) The level of access required on the disk. Supported are Read, Write. Changing this forces a new resource to be created.

Refer to the [SAS creation reference from Azure](https://docs.microsoft.com/rest/api/compute/disks/grant-access)
for additional details on the fields above.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Disk Export resource.

* `sas_url` - The computed Shared Access Signature (SAS) of the Managed Disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Disk.
* `read` - (Defaults to 5 minutes) Used when retrieving the Disk.
* `delete` - (Defaults to 30 minutes) Used when deleting the Disk.

## Import

Disk SAS Token can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_managed_disk_sas_token.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/disks/manageddisk1
```
