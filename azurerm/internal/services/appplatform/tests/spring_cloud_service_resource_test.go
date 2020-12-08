package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSpringCloudService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpringCloudServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSpringCloudService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpringCloudServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSpringCloudService_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudServiceExists(data.ResourceName),
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
				Config: testAccAzureRMSpringCloudService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudServiceExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSpringCloudService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpringCloudServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudService_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudServiceExists(data.ResourceName),
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
		},
	})
}

func TestAccAzureRMSpringCloudService_virtualNetwork(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpringCloudServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudService_virtualNetwork(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudServiceExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network.0.service_runtime_network_resource_group"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network.0.app_network_resource_group"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "outbound_public_ip_addresses.0"),
				),
			},
			data.ImportStep(
				// those field returned by api are "*"
				// import state verify ignore those fields
				"config_server_git_setting.0.ssh_auth.0.private_key",
				"config_server_git_setting.0.ssh_auth.0.host_key",
				"config_server_git_setting.0.ssh_auth.0.host_key_algorithm",
			),
		},
	})
}

func TestAccAzureRMSpringCloudService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_service", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpringCloudServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudService_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudServiceExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSpringCloudService_requiresImport),
		},
	})
}

func testCheckAzureRMSpringCloudServiceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Spring Cloud not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).AppPlatform.ServicesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Spring Cloud Service %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on AppPlatform.ServicesClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSpringCloudServiceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).AppPlatform.ServicesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_spring_cloud_service" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on AppPlatform.ServicesClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMSpringCloudService_basic(data acceptance.TestData) string {
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

func testAccAzureRMSpringCloudService_complete(data acceptance.TestData) string {
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
    instrumentation_key = azurerm_application_insights.test.instrumentation_key
  }

  tags = {
    Env     = "Test"
    version = "1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSpringCloudService_virtualNetwork(data acceptance.TestData) string {
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
  address_prefix       = "10.1.0.0/24"
}

resource "azurerm_subnet" "test2" {
  name                 = "internal2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.1.1.0/24"
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
  name                = "acctest-sc-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  network {
    app_subnet_id             = azurerm_subnet.test1.id
    service_runtime_subnet_id = azurerm_subnet.test2.id
    cidr_ranges               = ["10.4.0.0/16", "10.5.0.0/16", "10.3.0.1/16"]
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
    instrumentation_key = azurerm_application_insights.test.instrumentation_key
  }

  tags = {
    Env = "Test"
  }

  depends_on = [azurerm_role_assignment.test]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMSpringCloudService_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSpringCloudService_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_service" "import" {
  name                = azurerm_spring_cloud_service.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
