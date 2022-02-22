---
subcategory: "Compute"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_managed_disk_export"
description: |-
  Gets a Shared Access Signature (SAS Token) for an existing Managed Disk.

---

# Data Source: azurerm_managed_disk_export

Use this data source to obtain a Shared Access Signature (SAS Token) for an existing Managed Disk.

Shared access signatures allow fine-grained, ephemeral access control to various aspects of Managed Disk similar to blob/storage account container.

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

data "azurerm_managed_disk_export" "disk" {
  managed_disk_id     = azurerm_managed_disk.disk.id
  duration_in_seconds = 300
  access              = "Read"
}

output "sas_url_query_string" {
  value = data.azurerm_managed_disk_export.disk.sas
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
