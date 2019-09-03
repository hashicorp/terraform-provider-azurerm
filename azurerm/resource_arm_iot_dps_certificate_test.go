package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMIotDPSCertificate_basic(t *testing.T) {
	resourceName := "azurerm_iot_dps_certificate.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotDPSCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotDPSCertificate_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotDPSCertificateExists(resourceName),
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

func TestAccAzureRMIotDPSCertificate_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_iot_dps_certificate.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotDPSCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotDPSCertificate_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotDPSCertificateExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMIotDPSCertificate_requiresImport(rInt, location),
				ExpectError: testRequiresImportError("azurerm_iotdps"),
			},
		},
	})
}

func TestAccAzureRMIotDPSCertificate_update(t *testing.T) {
	resourceName := "azurerm_iot_dps_certificate.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotDPSCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotDPSCertificate_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotDPSCertificateExists(resourceName),
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
				Config: testAccAzureRMIotDPSCertificate_update(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotDPSCertificateExists(resourceName),
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

func testCheckAzureRMIotDPSCertificateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).iothub.DPSCertificateClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iot_dps_certificate" {
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

func testCheckAzureRMIotDPSCertificateExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

		client := testAccProvider.Meta().(*ArmClient).iothub.DPSCertificateClient
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

func testAccAzureRMIotDPSCertificate_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iot_dps" "test" {
  name                = "acctestIoTDPS-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "S1"
    tier     = "Standard"
    capacity = "1"
  }
}

resource "azurerm_iot_dps_certificate" "test" {
  name                = "acctestIoTDPSCertificate-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  iot_dps_name        = "${azurerm_iot_dps.test.name}"

  certificate_content = "${filebase64("testdata/batch_certificate.cer")}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMIotDPSCertificate_requiresImport(rInt int, location string) string {
	template := testAccAzureRMIotDPS_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_iot_dps_certificate" "test" {
  name                = "${azurerm_iot_dps_certificate.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  iot_dps_name        = "${azurerm_iot_dps.test.name}"

  certificate_content = "${filebase64("testdata/batch_certificate.cer")}"
}
`, template)
}

func testAccAzureRMIotDPSCertificate_update(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iot_dps" "test" {
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

resource "azurerm_iot_dps_certificate" "test" {
  name                = "acctestIoTDPSCertificate-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  iot_dps_name        = "${azurerm_iot_dps.test.name}"

  certificate_content = "${filebase64("testdata/application_gateway_test.cer")}"
}
`, rInt, location, rInt, rInt)
}
