
resource "azurerm_resource_group" "rg" {
    name     = "${var.resource_group_name}"
    location = "${var.resource_group_location}"
}

resource "azurerm_automation_account" "account" {
    name                = "tfex-automation-account"
    location            = "${azurerm_resource_group.rg.location}"
    resource_group_name = "${azurerm_resource_group.rg.name}"

    sku {
        name = "Basic"
    }
}

resource "azurerm_automation_schedule" "one-time" {
    name                    = "tfex-automation_schedule-one_time"
    resource_group_name     = "${azurerm_resource_group.rg.name}"
    automation_account_name = "${azurerm_automation_account.account.name}"
    frequency	            = "OneTime"
    //defaults start_time to now + 7 min
}

resource "azurerm_automation_schedule" "hour" {
    name                    = "tfex-automation_schedule-hour"
    resource_group_name     = "${azurerm_resource_group.rg.name}"
    automation_account_name = "${azurerm_automation_account.account.name}"
    frequency	            = "Hour"
    interval                = 2
    //timezone defaults to UTC
}
