change "resource-fix" {
  body = "`azurerm_windows_virtual_machine` - fix deployment from specialized images by omitting `OSProfile` when `admin_username` is not specified"
}