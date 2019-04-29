resource "azurerm_resource_group" "testrg" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_storage_account" "testsa" {
  name                     = "${var.prefix}sa"
  resource_group_name      = "${azurerm_resource_group.testrg.name}"
  location                 = "${azurerm_resource_group.testrg.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "${var.prefix}-sc"
  resource_group_name   = "${azurerm_resource_group.testrg.name}"
  storage_account_name  = "${azurerm_storage_account.testsa.name}"
  container_access_type = "private"
}

resource "azurerm_app_service_plan" "test" {
  name                = "${var.prefix}-splan"
  location            = "${azurerm_resource_group.testrg.location}"
  resource_group_name = "${azurerm_resource_group.testrg.name}"

  sku {
    tier = "Standard"
    size = "S1"
  }
}

data "azurerm_storage_account_sas" "test" {
  connection_string = "${azurerm_storage_account.testsa.primary_connection_string}"
  https_only        = true

  resource_types {
    service   = false
    container = false
    object    = true
  }

  services {
    blob  = true
    queue = false
    table = false
    file  = false
  }

  start  = "2019-03-21"
  expiry = "2020-03-21"

    permissions {
    read    = false
    write   = true
    delete  = false
    list    = false
    add     = false
    create  = false
    update  = false
    process = false
  }
}

resource "azurerm_app_service" "test" {
  name                    = "${var.prefix}-appservice"
  location                = "${azurerm_resource_group.testrg.location}"
  resource_group_name     = "${azurerm_resource_group.testrg.name}"
  app_service_plan_id     = "${azurerm_app_service_plan.test.id}"
  storage_account_url     = "https://${azurerm_storage_account.testsa.name}.blob.core.windows.net/${azurerm_storage_container.test.name}${data.azurerm_storage_account_sas.test.sas}&sr=b"

  backup_schedule {
    frequency_interval       = "30"
#    frequency_unit           = "Day"
#    keep_at_least_one_backup = false
#    retention_period_in_days = "9"
#    start_time               = "2019-04-29T09:40:00+02:00"
  }

  site_config {
    dotnet_framework_version = "v4.0"
    scm_type                 = "LocalGit"
  }

  app_settings = {
    "SOME_KEY" = "some-value"
  }

}
