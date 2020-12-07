package iothub

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMIotHubDPSCertificate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDPSCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubDPSCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSCertificateExists(data.ResourceName),
				),
			},
			data.ImportStep("certificate_content"),
		},
	})
}

func TestAccAzureRMIotHubDPSCertificate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDPSCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubDPSCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSCertificateExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMIotHubDPSCertificate_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_iothub_dps_certificate"),
			},
		},
	})
}

func TestAccAzureRMIotHubDPSCertificate_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_dps_certificate", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDPSCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubDPSCertificate_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSCertificateExists(data.ResourceName),
				),
			},
			data.ImportStep("certificate_content"),
			{
				Config: testAccAzureRMIotHubDPSCertificate_update(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSCertificateExists(data.ResourceName),
				),
			},
			data.ImportStep("certificate_content"),
		},
	})
}

func testCheckAzureRMIotHubDPSCertificateDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.DPSCertificateClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub_dps_certificate" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		iotDPSName := rs.Primary.Attributes["iot_dps_name"]

		resp, err := client.Get(ctx, name, resourceGroup, iotDPSName, "")
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("IoT Device Provisioning Service Certificate %s still exists in (device provisioning service %s / resource group %s)", name, iotDPSName, resourceGroup)
		}
	}
	return nil
}

func testCheckAzureRMIotHubDPSCertificateExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.DPSCertificateClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		iotDPSName := rs.Primary.Attributes["iot_dps_name"]

		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for IoT Device Provisioning Service Certificate: %s", name)
		}

		resp, err := client.Get(ctx, name, resourceGroup, iotDPSName, "")
		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return fmt.Errorf("Bad: IoT Device Provisioning Service Certificate %q (Device Provisioning Service %q / Resource Group %q) does not exist", name, iotDPSName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on iothubDPSCertificateClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMIotHubDPSCertificate_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }
}

resource "azurerm_iothub_dps_certificate" "test" {
  name                = "acctestIoTDPSCertificate-%d"
  resource_group_name = azurerm_resource_group.test.name
  iot_dps_name        = azurerm_iothub_dps.test.name

  certificate_content = filebase64("testdata/batch_certificate.cer")
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMIotHubDPSCertificate_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMIotHubDPSCertificate_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_dps_certificate" "import" {
  name                = azurerm_iothub_dps_certificate.test.name
  resource_group_name = azurerm_resource_group.test.name
  iot_dps_name        = azurerm_iothub_dps.test.name

  certificate_content = filebase64("testdata/batch_certificate.cer")
}
`, template)
}

func testAccAzureRMIotHubDPSCertificate_update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "S1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_dps_certificate" "test" {
  name                = "acctestIoTDPSCertificate-%d"
  resource_group_name = azurerm_resource_group.test.name
  iot_dps_name        = azurerm_iothub_dps.test.name

  certificate_content = filebase64("testdata/application_gateway_test.cer")
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
