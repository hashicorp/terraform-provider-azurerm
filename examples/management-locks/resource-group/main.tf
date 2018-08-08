resource "azurerm_resource_group" "example" {
  name     = "locked-resource-group"
  location = "West Europe"
}

resource "azurerm_management_lock" "example" {
  name       = "resource-group-level"
  scope      = "${azurerm_resource_group.example.id}"
  lock_level = "ReadOnly"
  notes      = "This Resource Group is Read-Only"
}
