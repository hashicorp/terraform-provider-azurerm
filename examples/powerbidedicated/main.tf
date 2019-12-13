resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_powerbidedicated_capacity" "example" {
  name                = "${var.prefix}powerbidedicatedcapacity"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  sku                 = "A1"
  administrators      = ["test2@microsoft.onmicrosoft.com"]
}
