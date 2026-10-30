change "new-property" {
  body = "`azurerm_oracle_autonomous_database` - add support for the `data_storage_size_in_gbs` property"
}

change "breaking" {
  body = "`azurerm_oracle_autonomous_database` - drop support for `data_storage_size_in_tbs` in favour of `data_storage_size_in_gbs`"
}
