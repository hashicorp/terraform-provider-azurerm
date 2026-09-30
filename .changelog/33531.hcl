change "resource-enhancement" {
  body = "`azurerm_storage_account` - prevent replacement of resource when `account_kind` changes from `StorageV2` to `Storage` as `Storage` is deprecated and can no longer be used to create accounts"
}
