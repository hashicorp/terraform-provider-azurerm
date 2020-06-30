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
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
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
