// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package eventgrid_system_topic_event_subscription_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func azureFunctionTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_storage_account" "test" {
  name                     = substr("acctestsa%[1]d", 0, 24)
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_service_plan" "test" {
  name                = "acctestASP-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  os_type             = "Linux"
  sku_name            = "S1"
}

resource "azurerm_linux_function_app" "test" {
  name                       = "acctestLA-%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  service_plan_id            = azurerm_service_plan.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key
  app_settings = {
    WEBSITE_CONTENTSHARE                     = "testacc-content-app"
    WEBSITE_CONTENTAZUREFILECONNECTIONSTRING = azurerm_storage_account.test.primary_blob_connection_string
    AzureWebJobsSecretStorageType            = "Blob"
  }
  site_config {
    application_stack {
      python_version = "3.9"
    }
    cors {
      allowed_origins     = ["https://portal.azure.com"]
      support_credentials = false
    }
  }
  lifecycle {
    ignore_changes = [
      app_settings["WEBSITE_CONTENTSHARE"],
    ]
  }
}

resource "azurerm_function_app_function" "test" {
  name            = "acctest-FnAppFn-%[1]d"
  function_app_id = azurerm_linux_function_app.test.id
  language        = "Python"
  test_data = jsonencode({
    "name" = "Azure"
  })
  config_json = jsonencode({
    "bindings" = [
      {
        "authLevel" = "function"
        "direction" = "in"
        "methods" = [
          "get",
          "post",
        ]
        "name" = "req"
        "type" = "EventGridTrigger"
      },
      {
        "direction" = "out"
        "name"      = "$return"
        "type"      = "http"
      },
    ]
  })
}

resource "azurerm_linux_function_app_slot" "test" {
  name                       = "acctest-LFAS-%[1]d"
  function_app_id            = azurerm_linux_function_app.test.id
  storage_account_name       = azurerm_storage_account.test.name
  storage_account_access_key = azurerm_storage_account.test.primary_access_key

  app_settings = {
    WEBSITE_CONTENTSHARE                     = "testacc-content-appslot"
    WEBSITE_CONTENTAZUREFILECONNECTIONSTRING = azurerm_storage_account.test.primary_blob_connection_string
    AzureWebJobsSecretStorageType            = "Blob"
  }

  site_config {
    application_stack {
      python_version = "3.9"
    }

    cors {
      support_credentials = false
      allowed_origins     = ["https://portal.azure.com"]
    }
  }

  lifecycle {
    ignore_changes = [
      app_settings["WEBSITE_CONTENTSHARE"],
    ]
  }
}

resource "azurerm_function_app_active_slot" "test" {
  slot_id = azurerm_linux_function_app_slot.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}
