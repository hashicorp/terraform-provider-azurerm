// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package springcloud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SpringCloudServiceResource struct{}

func TestAccSpringCloudService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			// those field returned by api are "*"
			// import state verify ignore those fields
			"config_server_git_setting.0.ssh_auth.0.private_key",
			"config_server_git_setting.0.ssh_auth.0.host_key",
			"config_server_git_setting.0.ssh_auth.0.host_key_algorithm",
			"config_server_git_setting.0.repository.0.http_basic_auth.0.username",
			"config_server_git_setting.0.repository.0.http_basic_auth.0.password",
		),
		{
			Config: r.singleGitRepo(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			"config_server_git_setting.0.ssh_auth.0.private_key",
			"config_server_git_setting.0.ssh_auth.0.host_key",
			"config_server_git_setting.0.ssh_auth.0.host_key_algorithm"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(
			// those field returned by api are "*"
			// import state verify ignore those fields
			"config_server_git_setting.0.ssh_auth.0.private_key",
			"config_server_git_setting.0.ssh_auth.0.host_key",
			"config_server_git_setting.0.ssh_auth.0.host_key_algorithm",
			"config_server_git_setting.0.repository.0.http_basic_auth.0.username",
			"config_server_git_setting.0.repository.0.http_basic_auth.0.password",
		),
	})
}

func TestAccSpringCloudService_virtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.virtualNetwork(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("network.0.service_runtime_network_resource_group").Exists(),
				check.That(data.ResourceName).Key("network.0.app_network_resource_group").Exists(),
				check.That(data.ResourceName).Key("outbound_public_ip_addresses.0").Exists(),
				check.That(data.ResourceName).Key("required_network_traffic_rules.0.protocol").Exists(),
				check.That(data.ResourceName).Key("required_network_traffic_rules.0.port").Exists(),
				check.That(data.ResourceName).Key("required_network_traffic_rules.0.ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("required_network_traffic_rules.0.fqdns.#").Exists(),
				check.That(data.ResourceName).Key("required_network_traffic_rules.0.direction").Exists(),
			),
		},
		data.ImportStep(
			// those field returned by api are "*"
			// import state verify ignore those fields
			"config_server_git_setting.0.ssh_auth.0.private_key",
			"config_server_git_setting.0.ssh_auth.0.host_key",
			"config_server_git_setting.0.ssh_auth.0.host_key_algorithm",
		),
	})
}

func TestAccSpringCloudService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccSpringCloudService_serviceRegistry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.serviceRegistry(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceRegistry(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("service_registry_id").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.serviceRegistry(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudService_buildAgentPool(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.buildAgentPool(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.buildAgentPool(data, "S2"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.buildAgentPool(data, "S1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudService_zoneRedundant(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zoneRedundant(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSpringCloudService_containerRegistry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.containerRegistry(data, "first"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("container_registry.0.password", "container_registry.1.password"),
		{
			Config: r.containerRegistry(data, "second"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("container_registry.0.password", "container_registry.1.password"),
	})
}

func TestAccSpringCloudService_marketplace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")
	r := SpringCloudServiceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.marketplace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t SpringCloudServiceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SpringCloudServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.AppPlatform.ServicesClient.Get(ctx, id.ResourceGroup, id.SpringName)
	if err != nil {
		return nil, fmt.Errorf("unable to read Spring Cloud Service %q (Resource Group %q): %+v", id.SpringName, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

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

func (SpringCloudServiceResource) zoneRedundant(data acceptance.TestData) string {
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
  zone_redundant      = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (SpringCloudServiceResource) serviceRegistry(data acceptance.TestData, enabled bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                     = "acctest-sc-%d"
  location                 = azurerm_resource_group.test.location
  resource_group_name      = azurerm_resource_group.test.name
  sku_name                 = "E0"
  service_registry_enabled = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, enabled)
}

func (SpringCloudServiceResource) buildAgentPool(data acceptance.TestData, size string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%d"
  location = "%s"
}

resource "azurerm_spring_cloud_service" "test" {
  name                  = "acctest-sc-%d"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  sku_name              = "E0"
  build_agent_pool_size = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, size)
}

func (SpringCloudServiceResource) singleGitRepo(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

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

    repository {
      name         = "repo2"
      uri          = "https://github.com/Azure-Samples/piggymetrics"
      label        = "config"
      search_paths = ["dir1", "dir2"]
    }
  }

  trace {
    connection_string = azurerm_application_insights.test.connection_string
    sample_rate       = 20
  }

  tags = {
    Env     = "Test"
    version = "1"
  }
}
	  `, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (SpringCloudServiceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

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

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

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

    repository {
      name         = "repo1"
      uri          = "https://github.com/Azure-Samples/piggymetrics"
      label        = "config"
      search_paths = ["dir1", "dir2"]
      http_basic_auth {
        username = "username"
        password = "password"
      }
    }

    repository {
      name         = "repo2"
      uri          = "https://github.com/Azure-Samples/piggymetrics"
      label        = "config"
      search_paths = ["dir1", "dir2"]
    }
  }

  trace {
    connection_string = azurerm_application_insights.test.connection_string
    sample_rate       = 20
  }

  tags = {
    Env     = "Test"
    version = "1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
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

func (SpringCloudServiceResource) containerRegistry(data acceptance.TestData, containerRegistryName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-spring-%[2]d"
  location = "%[1]s"
}

resource "azurerm_container_registry" "first" {
  name                = "acctestacr%[2]d1"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = true
}

resource "azurerm_container_registry" "second" {
  name                = "acctestacr%[2]d2"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = true
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "E0"

  container_registry {
    name     = "first"
    server   = azurerm_container_registry.first.login_server
    username = azurerm_container_registry.first.admin_username
    password = azurerm_container_registry.first.admin_password
  }

  container_registry {
    name     = "second"
    server   = azurerm_container_registry.second.login_server
    username = azurerm_container_registry.second.admin_username
    password = azurerm_container_registry.second.admin_password
  }

  default_build_service {
    container_registry_name = "%[3]s"
  }
}
`, data.Locations.Primary, data.RandomInteger, containerRegistryName)
}

func (SpringCloudServiceResource) marketplace(data acceptance.TestData) string {
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
  sku_name            = "E0"

  marketplace {
    plan      = "asa-ent-hr-mtr"
    publisher = "vmware-inc"
    product   = "azure-spring-cloud-vmware-tanzu-2"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r SpringCloudServiceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_service" "import" {
  name                = azurerm_spring_cloud_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, r.basic(data))
}
