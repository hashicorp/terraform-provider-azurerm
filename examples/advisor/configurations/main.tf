resource "azurerm_resource_group" "test" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_advisor" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  exclude             = false
  low_cpu_threshold   = "5"
}

