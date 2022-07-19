provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "xiaxintestRG-WAlogs"
  location = "east us"
}

resource "azurerm_service_plan" "test" {
  name                = "xiaxintestASP-logging"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "B1"
}

resource "azurerm_linux_web_app" "test" {
  name                = "LinuxWA-logging"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id

  site_config {}

  /* logs {
    application_logs {
      file_system_level = "Off"
      azure_blob_storage {
        level             = "Information"
        sas_url           = "http://x.com/"
        retention_in_days = 2
      }
    }

    http_logs {
      file_system {
        retention_in_days = 4
        retention_in_mb   = 25
      }
    }
    detailed_error_messages = true
  } */
}

resource "azurerm_linux_web_app" "test1" {
  name                = "LinuxWA-loggingDefault"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id
  site_config {}
  logs {
     application_logs {
      file_system_level = "Information"
      azure_blob_storage {
        level             = "Warning"
        sas_url           = "http://x.com/"
        retention_in_days = 2
      }
    }

     http_logs {
      file_system {
        retention_in_days = 7
        retention_in_mb   = 25
      }
    }
    detailed_error_messages = true
  }
}

resource "azurerm_linux_web_app" "test2" {
  name                = "LinuxWA-logging1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  service_plan_id     = azurerm_service_plan.test.id
  site_config {}
  logs {
#    application_logs {
#      file_system_level = "Off"
#      azure_blob_storage {
#        level             = "Warning"
#        sas_url           = "http://x.com/"
#        retention_in_days = 2
#      }
#    }

#    http_logs {
#      file_system {
#        retention_in_days = 7
#        retention_in_mb   = 25
#      }
#    }
    detailed_error_messages = true
  }
}
