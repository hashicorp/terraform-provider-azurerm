package azurerm

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMIotHubSharedAccessPolicy_basic(t *testing.T) {
	resourceName := "azurerm_iothub_shared_access_policy.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubSharedAccessPolicy_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubSharedAccessPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "registry_read", "true"),
					resource.TestCheckResourceAttr(resourceName, "registry_write", "true"),
					resource.TestCheckResourceAttr(resourceName, "service_connect", "false"),
					resource.TestCheckResourceAttr(resourceName, "device_connect", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", "acctest"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMIotHubSharedAccessPolicy_writeWithoutRead(t *testing.T) {
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMIotHubSharedAccessPolicy_writeWithoutRead(rInt, testLocation()),
				ExpectError: regexp.MustCompile("If `registry_write` is set to true, `registry_read` must also be set to true"),
			},
		},
	})
}

func TestAccAzureRMIotHubSharedAccessPolicy_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_iothub_shared_access_policy.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMIotHubSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubSharedAccessPolicy_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubSharedAccessPolicyExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMIotHubSharedAccessPolicy_requiresImport(rInt, location),
				ExpectError: testRequiresImportError("azurerm_iothub_shared_access_policy"),
			},
		},
	})
}

func testAccAzureRMIotHubSharedAccessPolicy_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "B1"
    tier     = "Basic"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_shared_access_policy" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_name         = "${azurerm_iothub.test.name}"
  name                = "acctest"

  registry_read  = true
  registry_write = true
}
`, rInt, location, rInt)
}

func testAccAzureRMIotHubSharedAccessPolicy_requiresImport(rInt int, location string) string {
	template := testAccAzureRMIotHubSharedAccessPolicy_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_shared_access_policy" "import" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_name         = "${azurerm_iothub.test.name}"
  name                = "acctest"

  registry_read  = true
  registry_write = true
}
`, template)
}

func testAccAzureRMIotHubSharedAccessPolicy_writeWithoutRead(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  sku {
    name     = "B1"
    tier     = "Basic"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_shared_access_policy" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  iothub_name         = "${azurerm_iothub.test.name}"
  name                = "acctest"

  registry_write = true
}
`, rInt, location, rInt)
}

func testCheckAzureRMIotHubSharedAccessPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		parsedIothubId, err := azure.ParseAzureResourceID(rs.Primary.ID)
		if err != nil {
			return err
		}
		iothubName := parsedIothubId.Path["IotHubs"]
		keyName := parsedIothubId.Path["IotHubKeys"]
		resourceGroup := parsedIothubId.ResourceGroup

		client := testAccProvider.Meta().(*ArmClient).iothub.ResourceClient

		for accessPolicyIterator, err := client.ListKeysComplete(ctx, resourceGroup, iothubName); accessPolicyIterator.NotDone(); err = accessPolicyIterator.NextWithContext(ctx) {
			if err != nil {
				return fmt.Errorf("Error loading Shared Access Profiles of IotHub %q (Resource Group %q): %+v", iothubName, resourceGroup, err)
			}

			if strings.EqualFold(*accessPolicyIterator.Value().KeyName, keyName) {
				return nil
			}
		}

		return fmt.Errorf("Bad: No shared access policy %s defined for IotHub %s", keyName, iothubName)

	}
}

func testCheckAzureRMIotHubSharedAccessPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).iothub.ResourceClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_iothub_shared_access_policy" {
			continue
		}

		keyName := rs.Primary.Attributes["name"]
		iothubName := rs.Primary.Attributes["iothub_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, iothubName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return fmt.Errorf("Bad: Get on iothubResourceClient: %+v", err)
		}

		for _, sharedAccessPolicy := range *resp.Properties.AuthorizationPolicies {
			if *sharedAccessPolicy.KeyName == keyName {
				return fmt.Errorf("Bad: Shared Access Policy %s still exists on IoTHb %s", keyName, iothubName)
			}
		}
	}
	return nil
}
