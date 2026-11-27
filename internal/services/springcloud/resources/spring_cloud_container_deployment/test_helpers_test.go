// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_container_deployment_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type SpringCloudAppResource struct{}

func (r SpringCloudAppResource) addon(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%[2]d"
  location = "%[1]s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                     = "acctest-sc-%[2]d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  sku_name                 = "E0"
  service_registry_enabled = true
}

resource "azurerm_spring_cloud_configuration_service" "test" {
  name                    = "default"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  repository {
    name                     = "fake"
    label                    = "master"
    patterns                 = ["auth-service", "gateway"]
    uri                      = "https://github.com/Azure-Samples/piggymetrics-config"
    search_paths             = ["/"]
    strict_host_key_checking = false
    username                 = "adminuser"
    password                 = "H@Sh1CoR3!"
  }
}

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%[2]d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
  addon_json = jsonencode({
    applicationConfigurationService = {
      resourceId = azurerm_spring_cloud_configuration_service.test.id
    }
    serviceRegistry = {
      resourceId = azurerm_spring_cloud_service.test.service_registry_id
    }
  })
}
`, data.Locations.Primary, data.RandomInteger)
}
