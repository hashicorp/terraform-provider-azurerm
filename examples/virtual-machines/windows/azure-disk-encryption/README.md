##Example: Windows Virtual Machine with Azure Disk Encryption

**NOTE:** Azure Disk Encryption uses Virtual Machine Extension, which creates unmanaged Secret in Key Vault, suggest using [Server-side Encryption](https://docs.microsoft.com/azure/virtual-machines/disk-encryption) instead which uses [disk_encryption_set_id](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/windows_virtual_machine#disk_encryption_set_id)

This example provisions a Windows Virtual Machine with disk encrypted using Azure Disk Encryption.

For more information, please refer to [Azure Disk Encryption Extension for Windows](https://docs.microsoft.com/azure/virtual-machines/extensions/azure-disk-enc-windows)
