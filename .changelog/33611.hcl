change "resource-fix" {
  body = "`azurerm_kubernetes_cluster` - tag-only updates now use the `UpdateTags` API (`PATCH`) instead of `CreateOrUpdate` (`PUT`), avoiding an unnecessary cluster reconcile"
}
