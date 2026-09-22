// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package springcloud_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SpringCloudAppDataSource struct{}

func TestAccDataSourceSpringCloudApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_spring_cloud_app", "test")
	r := SpringCloudAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("addon_json").HasValue(""),
				check.That(data.ResourceName).Key("custom_persistent_disk.#").HasValue("0"),
				check.That(data.ResourceName).Key("ingress_settings.0.backend_protocol").HasValue("Default"),
				check.That(data.ResourceName).Key("ingress_settings.0.read_timeout_in_seconds").HasValue("300"),
				check.That(data.ResourceName).Key("ingress_settings.0.send_timeout_in_seconds").HasValue("60"),
				check.That(data.ResourceName).Key("ingress_settings.0.session_affinity").HasValue("None"),
				check.That(data.ResourceName).Key("ingress_settings.0.session_cookie_max_age").HasValue("0"),
				check.That(data.ResourceName).Key("public_endpoint_enabled").HasValue("false"),
			),
		},
	})
}

func TestAccDataSourceSpringCloudApp_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_spring_cloud_app", "test")
	r := SpringCloudAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("addon_json").Exists(),
				check.That(data.ResourceName).Key("custom_persistent_disk.#").HasValue("1"),
				check.That(data.ResourceName).Key("custom_persistent_disk.0.mount_path").HasValue("/temp"),
				check.That(data.ResourceName).Key("custom_persistent_disk.0.share_name").HasValue("testname"),
				check.That(data.ResourceName).Key("custom_persistent_disk.0.mount_options.#").HasValue("4"),
				check.That(data.ResourceName).Key("custom_persistent_disk.0.read_only_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("ingress_settings.0.backend_protocol").HasValue("GRPC"),
				check.That(data.ResourceName).Key("ingress_settings.0.read_timeout_in_seconds").HasValue("700"),
				check.That(data.ResourceName).Key("ingress_settings.0.send_timeout_in_seconds").HasValue("700"),
				check.That(data.ResourceName).Key("ingress_settings.0.session_affinity").HasValue("Cookie"),
				check.That(data.ResourceName).Key("ingress_settings.0.session_cookie_max_age").HasValue("700"),
			),
		},
	})
}

func TestAccDataSourceSpringCloudApp_publicEndpoint(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_spring_cloud_app", "test")
	r := SpringCloudAppDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.publicEndpoint(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("public_endpoint_enabled").HasValue("true"),
			),
		},
	})
}

func (SpringCloudAppDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_spring_cloud_app" "test" {
  name                = azurerm_spring_cloud_app.test.name
  resource_group_name = azurerm_spring_cloud_app.test.resource_group_name
  service_name        = azurerm_spring_cloud_app.test.service_name
}
`, SpringCloudAppResource{}.complete(data))
}

func (SpringCloudAppDataSource) complete(data acceptance.TestData) string {
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

resource "azurerm_storage_account" "test1" {
  name                     = "acctest1%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_spring_cloud_storage" "test1" {
  name                    = "acctest-test1-%[2]d"
  spring_cloud_service_id = azurerm_spring_cloud_service.test.id
  storage_account_name    = azurerm_storage_account.test1.name
  storage_account_key     = azurerm_storage_account.test1.primary_access_key
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

  custom_persistent_disk {
    storage_name      = azurerm_spring_cloud_storage.test1.name
    mount_path        = "/temp"
    share_name        = "testname"
    mount_options     = ["uid=1000", "gid=1000", "file_mode=0755", "dir_mode=0755"]
    read_only_enabled = true
  }

  ingress_settings {
    session_affinity        = "Cookie"
    read_timeout_in_seconds = 700
    send_timeout_in_seconds = 700
    session_cookie_max_age  = 700
    backend_protocol        = "GRPC"
  }
}

data "azurerm_spring_cloud_app" "test" {
  name                = azurerm_spring_cloud_app.test.name
  resource_group_name = azurerm_spring_cloud_app.test.resource_group_name
  service_name        = azurerm_spring_cloud_app.test.service_name
}
`, data.Locations.Primary, data.RandomInteger, data.RandomStringOfLength(10))
}

func (SpringCloudAppDataSource) publicEndpoint(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app" "test" {
  name                    = "acctest-sca-%[2]d"
  resource_group_name     = azurerm_spring_cloud_service.test.resource_group_name
  service_name            = azurerm_spring_cloud_service.test.name
  public_endpoint_enabled = true
}

data "azurerm_spring_cloud_app" "test" {
  name                = azurerm_spring_cloud_app.test.name
  resource_group_name = azurerm_spring_cloud_app.test.resource_group_name
  service_name        = azurerm_spring_cloud_app.test.service_name
}
`, SpringCloudServiceResource{}.virtualNetwork(data), data.RandomInteger)
}
