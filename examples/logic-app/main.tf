resource "azurerm_resource_group" "main" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_logic_app_workflow" "main" {
  name                = "${var.prefix}-logicapp"
  location            = "${azurerm_resource_group.main.location}"
  resource_group_name = "${azurerm_resource_group.main.name}"
}

resource "azurerm_logic_app_trigger_recurrence" "hourly" {
  name         = "run-every-hour"
  logic_app_id = "${azurerm_logic_app_workflow.main.id}"
  frequency    = "Hour"
  interval     = 1
}

resource "azurerm_logic_app_action_http" "main" {
  name         = "clear-stale-objects"
  logic_app_id = "${azurerm_logic_app_workflow.main.id}"
  method       = "DELETE"
  uri          = "http://example.com/clear-stable-objects"
}
