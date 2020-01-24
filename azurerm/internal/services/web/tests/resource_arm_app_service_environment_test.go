package tests

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAppServiceEnvironment_basicWindows(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "pricing_tier", "I1"),
					resource.TestCheckResourceAttr(data.ResourceName, "front_end_scale_factor", "15"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "pricing_tier", "I1"),
					resource.TestCheckResourceAttr(data.ResourceName, "front_end_scale_factor", "15"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMAppServiceEnvironment_update(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "pricing_tier", "I2"),
					resource.TestCheckResourceAttr(data.ResourceName, "front_end_scale_factor", "10"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_tierAndScaleFactor(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_tierAndScaleFactor(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "pricing_tier", "I2"),
					resource.TestCheckResourceAttr(data.ResourceName, "front_end_scale_factor", "10"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_withAppServicePlan(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	aspData := acceptance.BuildTestData(t, "azurerm_app_service_plan", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_withAppServicePlan(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(data.ResourceName),
					testCheckAppServicePlanMemberOfAppServiceEnvironment(data.ResourceName, aspData.ResourceName),
				),
			},
		},
	})
}

func testAccAzureRMAppServiceEnvironment_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
	  
resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]

  subnet {
    name           = "asesubnet"
    address_prefix = "10.0.1.0/24"
	}
	  
  subnet {
    name           = "gatewaysubnet"
    address_prefix = "10.0.2.0/24"
  }
}

data "azurerm_subnet" "test" {
  name                 = "asesubnet"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
}

resource "azurerm_app_service_environment" "test" {
  name      = "acctest-ase-%[1]d"
  subnet_id = "${data.azurerm_subnet.test.id}"
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMAppServiceEnvironment_tierAndScaleFactor(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
	  
resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%[1]d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  address_space       = ["10.0.0.0/16"]

  subnet {
    name           = "asesubnet"
    address_prefix = "10.0.1.0/24"
	}
	  
  subnet {
    name           = "gatewaysubnet"
    address_prefix = "10.0.2.0/24"
  }
}

data "azurerm_subnet" "test" {
  name                 = "asesubnet"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
}

resource "azurerm_app_service_environment" "test" {
  name                   = "acctest-ase-%[1]d"
  subnet_id              = "${data.azurerm_subnet.test.id}"
  pricing_tier           = "I2"
  front_end_scale_factor = 10
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMAppServiceEnvironment_update(data acceptance.TestData) string {
	return testAccAzureRMAppServiceEnvironment_tierAndScaleFactor(data)
}

func testAccAzureRMAppServiceEnvironment_withAppServicePlan(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_plan" "test"{
  name                       = "acctest-ASP-%d"
  location                   = "${azurerm_resource_group.test.location}"
  resource_group_name        = "${azurerm_resource_group.test.name}"
  app_service_environment_id = "${azurerm_app_service_environment.test.id}"

  sku {
    tier = "Basic"
    size = "B1"
  }

}
`, template, data.RandomInteger)
}

func testCheckAzureRMAppServiceEnvironmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServiceEnvironmentsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		appServiceEnvironmentName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for App Service Environment: %s", appServiceEnvironmentName)
		}

		resp, err := client.Get(ctx, resourceGroup, appServiceEnvironmentName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: App Service Environment %q (resource group %q) does not exist", appServiceEnvironmentName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on appServiceEnvironmentClient: %+v", err)
		}

		return nil
	}
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

func testCheckAppServicePlanMemberOfAppServiceEnvironment(ase string, asp string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		aseClient := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServiceEnvironmentsClient
		aspClient := acceptance.AzureProvider.Meta().(*clients.Client).Web.AppServicePlansClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		aseResource, ok := s.RootModule().Resources[ase]
		if !ok {
			return fmt.Errorf("Not found: %s", ase)
		}

		appServiceEnvironmentName := aseResource.Primary.Attributes["name"]
		appServiceEnvironmentResourceGroup := aseResource.Primary.Attributes["resource_group_name"]

		aseResp, err := aseClient.Get(ctx, appServiceEnvironmentResourceGroup, appServiceEnvironmentName)
		if err != nil {
			if utils.ResponseWasNotFound(aseResp.Response) {
				return fmt.Errorf("Bad: App Service Environment %q (resource group %q) does not exist: %+v", appServiceEnvironmentName, appServiceEnvironmentResourceGroup, err)
			}
		}

		aspResource, ok := s.RootModule().Resources[asp]
		if !ok {
			return fmt.Errorf("Not found: %s", ase)
		}

		appServicePlanName := aspResource.Primary.Attributes["name"]
		appServicePlanResourceGroup := aspResource.Primary.Attributes["resource_group_name"]

		aspResp, err := aspClient.Get(ctx, appServicePlanResourceGroup, appServicePlanName)
		if err != nil {
			if utils.ResponseWasNotFound(aseResp.Response) {
				return fmt.Errorf("Bad: App Service Plan %q (resource group %q) does not exist: %+v", appServicePlanName, appServicePlanResourceGroup, err)
			}
		}
		if aspResp.HostingEnvironmentProfile.ID != aseResp.ID {
			return fmt.Errorf("Bad: App Service Plan %s not a member of App Service Environment %s", appServicePlanName, appServiceEnvironmentName)
		}

		return nil
	}
}
