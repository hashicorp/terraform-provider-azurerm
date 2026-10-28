change "new-property" {
  body = "`azurerm_oracle_autonomous_database` - add support for the `data_storage_size_in_gbs` property"

  make .tools/bin/changeloggy
  .tools/bin/changeloggy add -p 33295 -t breaking '`azurerm_oracle_autonomous_database` - drop support for `data_storage_size_in_tbs` in favour of `data_storage_size_in_gbs`'
}