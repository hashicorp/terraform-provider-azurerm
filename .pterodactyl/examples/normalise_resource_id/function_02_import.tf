#> ## Example - Import
import {
  id = provider::azurerm::normalise_resource_id("/Subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/import-example")
  to = azurerm_resource_group.test
}

resource "azurerm_resource_group" "test" {
  name     = "import-example"
  location = "westeurope"
}
