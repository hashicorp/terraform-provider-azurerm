change "dependency" {
  body = "dependencies: `kusto` - update API version to `2025-02-14`"
}

change "new-property" {
  body = "`azurerm_kusto_cluster` - add support for the `customer_managed_key` block"
}

change "resource-fix" {
  body = "`azurerm_kusto_cluster` - fix updates removing an externally managed `customer_managed_key`"
}
