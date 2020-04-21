package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMIoTCentralApplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotCentralApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotCentralApplication_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotCentralApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "ST1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIoTCentralApplication_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotCentralApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotCentralApplication_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotCentralApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "template", "iotc-default@1.0.0"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIoTCentralApplication_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotCentralApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotCentralApplication_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotCentralApplicationExists(data.ResourceName),
				),
			},
			{
				Config: testAccAzureRMIotCentralApplication_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotCentralApplicationExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "ST1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.ENV", "Test"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIoTCentralApplication_requiresImportErrorStep(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iotcentral_application", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotCentralApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotCentralApplication_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotCentralApplicationExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMIotCentralApplication_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_iotcentral_application"),
			},
		},
	})
}

func testCheckAzureRMIotCentralApplicationExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found:  %s", resourceName)
		}
		appName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for IoT Central Application:  %s", appName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTCentral.AppsClient
		resp, err := client.Get(ctx, resourceGroup, appName)
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: IoT Central Application %q (Resource Group  %q) does not exist", appName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get IoT Central Application:  %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMIotCentralApplicationDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTCentral.AppsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iotcentral_application" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if resp.StatusCode != http.StatusNotFound {
				return fmt.Errorf("Bad: Get IoT Central Application:  %+v", err)
			}
		}
	}
	return nil
}

func testAccAzureRMIotCentralApplication_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "acctest-iotcentralapp-%[1]d"
  sku                 = "ST1"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMIotCentralApplication_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sub_domain          = "acctest-iotcentralapp-%[2]d"
  display_name        = "acctest-iotcentralapp-%[2]d"
  sku                 = "ST1"
  template            = "iotc-default@1.0.0"
  tags = {
    ENV = "Test"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMIotCentralApplication_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[2]d"
  location = "%[1]s"
}

resource "azurerm_iotcentral_application" "test" {
  name                = "acctest-iotcentralapp-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sub_domain          = "acctest-iotcentralapp-%[2]d"
  display_name        = "acctest-iotcentralapp-%[2]d"
  sku                 = "ST1"
  tags = {
    ENV = "Test"
  }
}
`, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMIotCentralApplication_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMIotCentralApplication_basic(data)
	return fmt.Sprintf(`
%s
resource "azurerm_iotcentral_application" "import" {
  name                = azurerm_iotcentral_application.test.name
  resource_group_name = azurerm_iotcentral_application.test.resource_group_name
  location            = azurerm_iotcentral_application.test.location
  sub_domain          = azurerm_iotcentral_application.test.sub_domain
  display_name        = azurerm_iotcentral_application.test.display_name
  sku                 = "ST1"
}
`, template)
}
