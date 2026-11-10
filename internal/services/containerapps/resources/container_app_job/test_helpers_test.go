// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_job_test

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

type ContainerAppEnvironmentDaprComponentResource struct{}

func (r ContainerAppEnvironmentDaprComponentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_storage_account" "test" {
  name                = "unlikely23exst2acct%[3]s"
  resource_group_name = azurerm_resource_group.test.name

  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"

  tags = {
    environment = "production"
  }
}

resource "azurerm_storage_container" "test" {
  name                  = "container-app-storage"
  storage_account_id    = azurerm_storage_account.test.id
  container_access_type = "private"
}

resource "azurerm_container_app_environment_dapr_component" "test" {
  name                         = "acctest-dapr-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  component_type               = "state.azure.blobstorage"
  version                      = "v1"

  init_timeout  = "10s"
  ignore_errors = true

  secret {
    name  = "secret"
    value = "sauce"
  }

  secret {
    name  = "storage-account-access-key"
    value = azurerm_storage_account.test.primary_access_key
  }

  metadata {
    name        = "storage-account-key"
    secret_name = "storage-account-access-key"
  }

  metadata {
    name  = "storage-container-name"
    value = azurerm_storage_container.test.name
  }

  metadata {
    name  = "SOME_APP_SETTING"
    value = "scwiffy"
  }

  scopes = ["testapp"]
}
`, r.template(data), data.RandomInteger, data.RandomString)
}

func (r ContainerAppEnvironmentDaprComponentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-CAEnv-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestCAEnv-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[1]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  logs_destination           = "log-analytics"
}
`, data.RandomInteger, data.Locations.Primary)
}

type ContainerAppEnvironmentResource struct{}

func (r ContainerAppEnvironmentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_container_app_environment" "test" {
  name                       = "acctest-CAEnv%[2]d"
  resource_group_name        = azurerm_resource_group.test.name
  location                   = azurerm_resource_group.test.location
  logs_destination           = "log-analytics"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  infrastructure_subnet_id   = azurerm_subnet.control.id

  internal_load_balancer_enabled = true
  zone_redundancy_enabled        = true
  mutual_tls_enabled             = true

  workload_profile {
    maximum_count         = 3
    minimum_count         = 0
    name                  = "D4-01"
    workload_profile_type = "D4"
  }

  tags = {
    Foo    = "Bar"
    secret = "sauce"
  }
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "diagnostics"
  target_resource_id         = azurerm_container_app_environment.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  enabled_log {
    category_group = "allLogs"
  }

  enabled_metric {
    category = "AllMetrics"
  }
}
`, r.templateVNet(data), data.RandomInteger)
}

func (r ContainerAppEnvironmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-CAE-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ContainerAppEnvironmentResource) templateVNet(data acceptance.TestData) string {
	return fmt.Sprintf(`


%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "control" {
  name                 = "control-plane"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/23"]

  delegation {
    name = "acctestdelegation%[2]d"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.App/environments"
    }
  }
}


`, r.template(data), data.RandomInteger)
}

type ContainerAppResource struct{}

func (ContainerAppResource) templatePlusExtras(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_container_registry" "test" {
  name                = "testacccr%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Basic"
  admin_enabled       = true
}

resource "azurerm_storage_share" "test" {
  name               = "testshare%[3]s"
  storage_account_id = azurerm_storage_account.test.id
  quota              = 1
}

resource "azurerm_container_app_environment_storage" "test" {
  name                         = "testacc-caes-%[2]d"
  container_app_environment_id = azurerm_container_app_environment.test.id
  account_name                 = azurerm_storage_account.test.name
  access_key                   = azurerm_storage_account.test.primary_access_key
  share_name                   = azurerm_storage_share.test.name
  access_mode                  = "ReadWrite"
}
	`, ContainerAppEnvironmentDaprComponentResource{}.complete(data), data.RandomInteger, data.RandomString)
}
