change "resource-fix" {
  body = "`azurerm_netapp_volume` - remove the default for `unix_permissions` as this prevented provisioning SMB or dual protocol volumes with NTFS security style"
}
