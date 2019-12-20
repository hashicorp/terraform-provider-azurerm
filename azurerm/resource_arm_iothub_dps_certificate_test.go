package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMIotHubDPSCertificate_basic(t *testing.T) {
	resourceName := "azurerm_iothub_dps_certificate.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDPSCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubDPSCertificate_basic(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSCertificateExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate_content",
				},
			},
		},
	})
}

func TestAccAzureRMIotHubDPSCertificate_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_iothub_dps_certificate.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDPSCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubDPSCertificate_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSCertificateExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMIotHubDPSCertificate_requiresImport(rInt, location),
				ExpectError: acceptance.RequiresImportError("azurerm_iothubdps"),
			},
		},
	})
}

func TestAccAzureRMIotHubDPSCertificate_update(t *testing.T) {
	resourceName := "azurerm_iothub_dps_certificate.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDPSCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubDPSCertificate_basic(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSCertificateExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate_content",
				},
			},
			{
				Config: testAccAzureRMIotHubDPSCertificate_update(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDPSCertificateExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate_content",
				},
			},
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

		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.DPSCertificateClient
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

func testAccAzureRMIotHubDPSCertificate_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }
}

resource "azurerm_iothub_dps_certificate" "test" {
  name                = "acctestIoTDPSCertificate-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  iot_dps_name        = "${azurerm_iothub_dps.test.name}"

  certificate_content = "${filebase64("testdata/batch_certificate.cer")}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMIotHubDPSCertificate_requiresImport(rInt int, location string) string {
	template := testAccAzureRMIotHubDPS_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_dps_certificate" "test" {
  name                = "${azurerm_iothub_dps_certificate.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  iot_dps_name        = "${azurerm_iothub_dps.test.name}"

  certificate_content = "${filebase64("testdata/batch_certificate.cer")}"
}
`, template)
}

func testAccAzureRMIotHubDPSCertificate_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_dps_certificate" "test" {
  name                = "acctestIoTDPSCertificate-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  iot_dps_name        = "${azurerm_iothub_dps.test.name}"

  certificate_content = "${filebase64("testdata/application_gateway_test.cer")}"
}
`, rInt, location, rInt, rInt)
}
