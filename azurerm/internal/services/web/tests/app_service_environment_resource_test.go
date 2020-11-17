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

func TestAccAzureRMAppServiceEnvironment_basic(t *testing.T) {
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

func TestAccAzureRMAppServiceEnvironment_requiresImport(t *testing.T) {
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
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMAppServiceEnvironment_requiresImport),
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
				Config: testAccAzureRMAppServiceEnvironment_tierAndScaleFactor(data),
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
					resource.TestCheckResourceAttrPair(data.ResourceName, "id", aspData.ResourceName, "app_service_environment_id"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_dedicatedResourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_dedicatedResourceGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_withCertificatePfx(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")
	certData := acceptance.BuildTestData(t, "azurerm_app_service_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_withCertificatePfx(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCertHostingEnvProfileIdMatchesAseId(data.ResourceName, certData.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMAppServiceEnvironment_internalLoadBalancer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_app_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAppServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAppServiceEnvironment_internalLoadBalancerAndWhitelistedIpRanges(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAppServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "internal_load_balancing_mode", "Web, Publishing"),
				),
			},
			data.ImportStep(),
		},
	})
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

func testCheckAzureRMCertHostingEnvProfileIdMatchesAseId(aseResourceName, certResourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		certsConn := acceptance.AzureProvider.Meta().(*clients.Client).Web.CertificatesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		aseRs, ok := s.RootModule().Resources[aseResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", aseResourceName)
		}
		certRs, ok := s.RootModule().Resources[certResourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", certResourceName)
		}

		aseId := aseRs.Primary.Attributes["id"]

		certName := certRs.Primary.Attributes["name"]
		certResourceGroup, hasResourceGroup := certRs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Certificate: %s", certName)
		}

		certResp, err := certsConn.Get(ctx, certResourceGroup, certName)
		if err != nil {
			if utils.ResponseWasNotFound(certResp.Response) {
				return fmt.Errorf("Bad: Certificatet %q (resource group: %q) does not exist", certName, certResourceGroup)
			}

			return fmt.Errorf("Bad: Get on certificatesClient: %+v", err)
		}

		if *certResp.HostingEnvironmentProfile.ID != aseId {
			return fmt.Errorf("Bad: Certificate hostingEnvironmentProfile.ID (%s) not equal to ASE ID (%s)", *certResp.HostingEnvironmentProfile.ID, aseId)
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

func testAccAzureRMAppServiceEnvironment_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "import" {
  name      = azurerm_app_service_environment.test.name
  subnet_id = azurerm_app_service_environment.test.subnet_id
}
`, template)
}

func testAccAzureRMAppServiceEnvironment_tierAndScaleFactor(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "test" {
  name                   = "acctest-ase-%d"
  subnet_id              = azurerm_subnet.ase.id
  pricing_tier           = "I2"
  front_end_scale_factor = 10
}
`, template, data.RandomInteger)
}

func testAccAzureRMAppServiceEnvironment_withAppServicePlan(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_plan" "test" {
  name                       = "acctest-ASP-%d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  app_service_environment_id = azurerm_app_service_environment.test.id

  sku {
    tier     = "Isolated"
    size     = "I1"
    capacity = 1
  }
}
`, template, data.RandomInteger)
}

func testAccAzureRMAppServiceEnvironment_dedicatedResourceGroup(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group" "test2" {
  name     = "acctestRG2-%[2]d"
  location = "%s"
}

resource "azurerm_app_service_environment" "test" {
  name                = "acctest-ase-%[2]d"
  resource_group_name = azurerm_resource_group.test2.name
  subnet_id           = azurerm_subnet.ase.id
}
`, template, data.RandomInteger, data.Locations.Secondary)
}

func testAccAzureRMAppServiceEnvironment_withCertificatePfx(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_certificate" "test" {
  name                           = "acctest-cert-%d"
  resource_group_name            = azurerm_app_service_environment.test.resource_group_name
  location                       = azurerm_resource_group.test.location
  pfx_blob                       = filebase64("testdata/app_service_certificate.pfx")
  password                       = "terraform"
  hosting_environment_profile_id = azurerm_app_service_environment.test.id
}
`, template, data.RandomInteger)
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

func testAccAzureRMAppServiceEnvironment_internalLoadBalancerAndWhitelistedIpRanges(data acceptance.TestData) string {
	template := testAccAzureRMAppServiceEnvironment_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_app_service_environment" "test" {
  name                         = "acctest-ase-%d"
  subnet_id                    = azurerm_subnet.ase.id
  pricing_tier                 = "I1"
  front_end_scale_factor       = 5
  internal_load_balancing_mode = "Web, Publishing"
  allowed_user_ip_cidrs        = ["11.22.33.44/32", "55.66.77.0/24"]
}
`, template, data.RandomInteger)
}
