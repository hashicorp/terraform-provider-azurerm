---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_disk_export_revoke"
description: |-
  Revokes a Shared Access Signature (SAS Token) for an existing Managed Disk.

---

# Data Source: azurerm_managed_disk_export_revoke

Use this data source to revoke Shared Access Signature (SAS Token) for an existing Managed Disk obtained after grant.

Shared access signatures allow fine-grained, ephemeral access control to various aspects of Managed Disk similar to blob/storage account container.

With the help of this resource, data from the disk can be copied from managed disk to a storage blob or to some other system without the need of azcopy.

## Example Usage

```hcl
resource "azurerm_resource_group" "rg" {
  name     = "resourceGroupName"
  location = "West Europe"
}

resource "azurerm_managed_disk" "disk" {
  name                 = "azureManagedDisk"
  location             = azurerm_resource_group.rg.location
  resource_group_name  = azurerm_resource_group.rg.name
  storage_account_type = "Standard_LRS"
  create_option        = "Empty"
  disk_size_gb         = "1"

  tags = {
    environment = "staging"
  }
}

resource "azurerm_storage_account" "example" {
  name                     = "examplestoracc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "example" {
  name                  = "content"
  storage_account_name  = azurerm_storage_account.example.name
  container_access_type = "private"
}

resource "azurerm_storage_blob" "example" {
  name                   = "${azurerm_managed_disk.disk.name}.vhd"
  storage_account_name   = azurerm_storage_account.example.name
  storage_container_name = azurerm_storage_container.example.name
  type                   = "Page"
  source_uri             = data.azurerm_managed_disk_export.export.sas
}

data "azurerm_managed_disk_export" "export" {
  managed_disk_id     = azurerm_managed_disk.disk.id
  duration_in_seconds = 300
  access              = "Read"
}

## revoke the sas token post export is complete
data "azurerm_managed_disk_export_revoke" "revoke" {
  depends_on      = [azurerm_storage_blob.example]
  managed_disk_id = "/subscriptions/14b86a40-8d8f-4e69-abaf-42cbb0b8a331/resourceGroups/FACTORY/providers/Microsoft.Compute/disks/ash-rhel7-image"
}
```

## Argument Reference

* `managed_disk_id` - The ID of an existing Managed Disk which should be exported. Changing this forces a new resource to be created.

* `duration_in_seconds` - The duration for which the export should be allowed. Should be greater than 30 seconds.

* `access` - (Optional) The level of access required on the disk. Supported are Read, Write. Defaults to Read. 

Refer to the [SAS creation reference from Azure](https://docs.microsoft.com/en-us/rest/api/compute/disks/grant-access)
for additional details on the fields above.

## Attributes Reference

* `sas` - The computed Shared Access Signature (SAS) of the Managed Disk.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Blob Container.
