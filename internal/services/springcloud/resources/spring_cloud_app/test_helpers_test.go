// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_app_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type SpringCloudServiceResource struct{}

func (SpringCloudServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (SpringCloudServiceResource) virtualNetwork(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

provider "azuread" {}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestai-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.1.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test1" {
  name                 = "internal1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.0.0/24"]
}

resource "azurerm_subnet" "test2" {
  name                 = "internal2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.1.1.0/24"]
}

data "azuread_service_principal" "test" {
  display_name = "Azure Spring Cloud Resource Provider"
}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_virtual_network.test.id
  role_definition_name = "Owner"
  principal_id         = data.azuread_service_principal.test.object_id
}

resource "azurerm_spring_cloud_service" "test" {
  name                               = "acctest-sc-%d"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  log_stream_public_endpoint_enabled = true

  network {
    app_subnet_id             = azurerm_subnet.test1.id
    service_runtime_subnet_id = azurerm_subnet.test2.id
    cidr_ranges               = ["10.4.0.0/16", "10.5.0.0/16", "10.3.0.1/16"]
    read_timeout_seconds      = 2
    outbound_type             = "loadBalancer"
  }

  config_server_git_setting {
    uri          = "git@bitbucket.org:Azure-Samples/piggymetrics.git"
    label        = "config"
    search_paths = ["dir1", "dir4"]

    ssh_auth {
      private_key                      = file("testdata/private_key")
      host_key                         = file("testdata/host_key")
      host_key_algorithm               = "ssh-rsa"
      strict_host_key_checking_enabled = false
    }
  }

  trace {
    connection_string = azurerm_application_insights.test.connection_string
  }

  tags = {
    Env = "Test"
  }

  depends_on = [azurerm_role_assignment.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
