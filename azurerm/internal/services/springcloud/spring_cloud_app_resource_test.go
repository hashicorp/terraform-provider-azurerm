package springcloud_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMSpringCloudApp_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpringCloudAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudAppExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSpringCloudApp_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpringCloudAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudAppExists(data.ResourceName),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMSpringCloudApp_requiresImport),
		},
	})
}

func TestAccAzureRMSpringCloudApp_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpringCloudAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudApp_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudAppExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.principal_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "identity.0.tenant_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMSpringCloudApp_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_spring_cloud_app", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMSpringCloudAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMSpringCloudApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudAppExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSpringCloudApp_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudAppExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMSpringCloudApp_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMSpringCloudAppExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMSpringCloudAppExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Spring Cloud App not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["service_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).AppPlatform.AppsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		if resp, err := client.Get(ctx, resourceGroup, serviceName, name, ""); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Spring Cloud App %q (Spring Cloud Name %q / Resource Group %q) does not exist", name, serviceName, resourceGroup)
			}
			return fmt.Errorf("bad: Get on AppPlatform.AppsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMSpringCloudAppDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).AppPlatform.AppsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_spring_cloud_app" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resGroup := rs.Primary.Attributes["resource_group_name"]
		serviceName := rs.Primary.Attributes["service_name"]

		if resp, err := client.Get(ctx, resGroup, serviceName, name, ""); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("bad: Get on AppPlatform.AppsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMSpringCloudApp_basic(data acceptance.TestData) string {
	template := testAccAzureRMSpringCloudApp_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name
}
`, template, data.RandomInteger)
}

func testAccAzureRMSpringCloudApp_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMSpringCloudApp_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app" "import" {
  name                = azurerm_spring_cloud_app.test.name
  resource_group_name = azurerm_spring_cloud_app.test.resource_group_name
  service_name        = azurerm_spring_cloud_app.test.service_name
}
`, template)
}

func testAccAzureRMSpringCloudApp_complete(data acceptance.TestData) string {
	template := testAccAzureRMSpringCloudApp_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_spring_cloud_app" "test" {
  name                = "acctest-sca-%d"
  resource_group_name = azurerm_spring_cloud_service.test.resource_group_name
  service_name        = azurerm_spring_cloud_service.test.name

  identity {
    type = "SystemAssigned"
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMSpringCloudApp_template(data acceptance.TestData) string {
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
