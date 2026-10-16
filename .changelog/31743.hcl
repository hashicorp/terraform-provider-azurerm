
change "new-resource" {
  body = "**New Resource**: `azurerm_storage_mover_nfs_file_share_target_endpoint`"
}
change "new-list-resource" {
  body = "**New List Resource**: `azurerm_storage_mover_nfs_file_share_target_endpoint`"
}
change "resource-fix" {
  body = "`azurerm_storage_mover_source_endpoint` - reject endpoint types outside the NFS mount resource during import, read, update, and delete"
}
change "resource-fix" {
  body = "`azurerm_storage_mover_target_endpoint` - reject endpoint types outside the blob container resource during import, read, update, and delete"
}
