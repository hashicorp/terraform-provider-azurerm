change "resource-fix" {
  body = "`azurerm_oracle_autonomous_database` - fix Autonomous JSON Database creation by defaulting `character_set` and `national_character_set` and omitting both properties from AJD create requests"
}
