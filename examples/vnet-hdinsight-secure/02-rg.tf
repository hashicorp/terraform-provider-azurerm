# Create a resource group
resource "azurerm_resource_group" "rg" {
  name     = "${var.resource_prefix}-rg"
  location = "${var.location}"
  tags     = "${var.tags}"
}
