resource "azurerm_resource_group" "example" {
  name     = "locked-resource-group"
  location = "West Europe"
}

resource "azurerm_public_ip" "example" {
  name                         = "locked-publicip"
  location                     = "${azurerm_resource_group.example.location}"
  resource_group_name          = "${azurerm_resource_group.example.name}"
  public_ip_address_allocation = "static"
  idle_timeout_in_minutes      = 30
}

resource "azurerm_management_lock" "example" {
  name       = "resource-ip"
  scope      = "${azurerm_public_ip.example.id}"
  lock_level = "CanNotDelete"
  notes      = "Locked because it's needed by a third-party"
}
