resource "azurerm_resource_group" "example" {
  name     = "${var.resource_group}"
  location = "${var.location}"
}

resource "random_integer" "ri" {
  min = 10000
  max = 99999
}

resource "azurerm_data_lake_store" "example" {
  name                = "tfexdlstore${random_integer.ri.result}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"

  tier = "Consumption"
}

resource "azurerm_data_lake_store_firewall_rule" "test" {
  name                = "tfex-datalakestore-fwrule"
  account_name        = "${azurerm_data_lake_store.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}

resource "azurerm_data_lake_analytics_account" "example" {
  name                = "tfexdlanalytics${random_integer.ri.result}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  location            = "${azurerm_resource_group.example.location}"
  tier                = "Consumption"

  default_store_account_name = "${azurerm_data_lake_store.example.name}"
}

resource "azurerm_data_lake_analytics_firewall_rule" "test" {
  name                = "tfex-datalakestore-fwrule"
  account_name        = "${azurerm_data_lake_analytics_account.example.name}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  start_ip_address    = "0.0.0.0"
  end_ip_address      = "0.0.0.0"
}
