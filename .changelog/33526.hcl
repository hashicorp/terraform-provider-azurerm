change "resource-fix" {
  body = "`azurerm_subnet` - the `actions` property within the `delegation.service_delegation` block is now deprecated and computed, as this property is read-only in the Azure API and any configured value is ignored"
}
