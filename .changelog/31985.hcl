change "resource-fix" {
  body = "`azurerm_data_protection_backup_vault` - fix creation when `datastore_type` is `ArchiveStore`"
}

change "breaking" {
  body = "`azurerm_data_protection_backup_vault` - replace `datastore_type` and `redundancy` with `storage_setting` blocks in the resource and `storage_settings` in the data source in 6.0"
}
