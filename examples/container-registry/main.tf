provider "azurerm" {
   # currently using a version build from GitHub while PR on ACR has not been approved.
   # https://github.com/terraform-providers/terraform-provider-azurerm/pull/2055 
   # While done, uncomment and set new version 
   # version = "~>1.5"
}
resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_group_name}"
  location = "${var.resource_group_location}"
}

resource "random_integer" "ri" {
  min = 10000
  max = 99999
}

resource "azurerm_container_registry" "acr" {
  name                   = "acr${random_integer.ri.result}"
  resource_group_name    = "${azurerm_resource_group.rg.name}"
  location               = "${azurerm_resource_group.rg.location}"
  sku                    = "${var.sku}"
  admin_enabled          = "${var.admin_enabled}"
  georeplication_locations = "${var.georeplication_locations}"
}
