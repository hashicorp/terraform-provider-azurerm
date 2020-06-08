package dd

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceEnvironment_withCertificatePfx(t *testing.T) {
	subscription := os.Getenv("ARM_SUBSCRIPTION_ID")
	if len(subscription) < 1 {
		t.Fatal("error retrieving subscription ID from environment")
	}

	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	certData := acceptance.BuildTestData(t, "azurerm_app_service_certificate", "test")

	expectedHepi := fmt.Sprintf(
		"/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Web/hostingEnvironments/acctest-ase-%d",
		subscription, data.RandomInteger, data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_withCertificatePfx(data, expectedHepi),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(certData.ResourceName, "hosting_environment_profile_id", expectedHepi),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMAppServiceEnvironmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServiceEnvironmentsClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_app_service_environment" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, resourceGroup, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return nil
	}

	return nil
}

func testAccAzureRMAppServiceEnvironment_basic(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "test" {
  name      = "acctest-ase-%d"
  subnet_id = azurerm_subnet.ase.id
}
`, template, data.RandomInteger)
}

func testAccAzureRMAppServiceEnvironment_withCertificatePfx(data acceptance.TestData, hepi string) string {
	template := testAccAzureRMAppServiceEnvironment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_certificate" "test" {
  name                           = "acctest-cert-%d"
  resource_group_name            = azurerm_resource_group.test.name
  location                       = azurerm_resource_group.test.location
  pfx_blob                       = filebase64("testdata/app_service_certificate.pfx")
  password                       = "terraform"
  hosting_environment_profile_id = "%s"
}
`, template, data.RandomInteger, hepi)
}

func testAccAzureRMAppServiceEnvironment_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "ase" {
  name                 = "asesubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_subnet" "gateway" {
  name                 = "gatewaysubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
