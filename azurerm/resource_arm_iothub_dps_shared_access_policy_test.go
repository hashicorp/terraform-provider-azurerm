package azurerm

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMIotHubDpsSharedAccessPolicy_basic(t *testing.T) {
	resourceName := "azurerm_iothub_dps_shared_access_policy.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDpsSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubDpsSharedAccessPolicy_basic(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubDpsSharedAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "acctest"),
					resource.TestCheckResourceAttr(resourceName, "enrollment_read", "false"),
					resource.TestCheckResourceAttr(resourceName, "enrollment_write", "false"),
					resource.TestCheckResourceAttr(resourceName, "registration_read", "false"),
					resource.TestCheckResourceAttr(resourceName, "registration_write", "false"),
					resource.TestCheckResourceAttr(resourceName, "service_config", "true"),
				),
			},
		},
	})
}

func TestAccAzureRMIotHubDpsSharedAccessPolicy_writeWithoutRead(t *testing.T) {
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDpsSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMIotHubDpsSharedAccessPolicy_writeWithoutRead(rInt, acceptance.Location()),
				ExpectError: regexp.MustCompile("If `registration_write` is set to true, `registration_read` must also be set to true"),
			},
		},
	})
}

func TestAccAzureRMIotHubDpsSharedAccessPolicy_enrollmentReadWithoutRegistration(t *testing.T) {
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDpsSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMIotHubDpsSharedAccessPolicy_enrollmentReadWithoutRegistration(rInt, acceptance.Location()),
				ExpectError: regexp.MustCompile("If `enrollment_read` is set to true, `registration_read` must also be set to true"),
			},
		},
	})
}

func TestAccAzureRMIotHubDpsSharedAccessPolicy_enrollmentWriteWithoutOthers(t *testing.T) {
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubDpsSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMIotHubDpsSharedAccessPolicy_enrollmentWriteWithoutOthers(rInt, acceptance.Location()),
				ExpectError: regexp.MustCompile("If `enrollment_write` is set to true, `enrollment_read`, `registration_read`, and `registration_write` must also be set to true"),
			},
		},
	})
}

func testAccAzureRMIotHubDpsSharedAccessPolicy_basic(rInt int, location string) string {
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
resource "azurerm_iothub_dps_shared_access_policy" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_dps_name     = "${azurerm_iothub_dps.test.name}"
  name                = "acctest"
  service_config  = true
}
`, rInt, location, rInt)
}

func testAccAzureRMIotHubDpsSharedAccessPolicy_writeWithoutRead(rInt int, location string) string {
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
resource "azurerm_iothub_dps_shared_access_policy" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_dps_name     = "${azurerm_iothub_dps.test.name}"
  name                = "acctest"
  registration_write = true
}
`, rInt, location, rInt)
}

func testAccAzureRMIotHubDpsSharedAccessPolicy_enrollmentReadWithoutRegistration(rInt int, location string) string {
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
resource "azurerm_iothub_dps_shared_access_policy" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_dps_name     = "${azurerm_iothub_dps.test.name}"
  name                = "acctest"
  enrollment_read = true
}
`, rInt, location, rInt)
}

func testAccAzureRMIotHubDpsSharedAccessPolicy_enrollmentWriteWithoutOthers(rInt int, location string) string {
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
resource "azurerm_iothub_dps_shared_access_policy" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_dps_name     = "${azurerm_iothub_dps.test.name}"
  name                = "acctest"
  enrollment_write = true
}
`, rInt, location, rInt)
}

func testCheckAzureRMIotHubDpsSharedAccessPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		_, err := azure.ParseAzureResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}

		keyName := rs.Primary.Attributes["name"]
		iothubDpsName := rs.Primary.Attributes["iothub_dps_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.DPSResourceClient

		_, err = client.ListKeysForKeyName(ctx, iothubDpsName, keyName, resourceGroup)
		if err != nil {
			return fmt.Errorf("Bad: No shared access policy %s defined for IotHub DPS %s", keyName, iothubDpsName)
		}

		return nil
	}
}

func testCheckAzureRMIotHubDpsSharedAccessPolicyDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.DPSResourceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub_dps_shared_access_policy" {
			continue
		}

		keyName := rs.Primary.Attributes["name"]
		iothubDpsName := rs.Primary.Attributes["iothub_dps_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, iothubDpsName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Get on iothubDPSResourceClient: %+v", err)
		}

		for _, sharedAccessPolicy := range *resp.Properties.AuthorizationPolicies {
			if *sharedAccessPolicy.KeyName == keyName {
				return fmt.Errorf("Bad: Shared Access Policy %s still exists on IoTHb DPS %s", keyName, iothubDpsName)
			}
		}
	}
	return nil
}
