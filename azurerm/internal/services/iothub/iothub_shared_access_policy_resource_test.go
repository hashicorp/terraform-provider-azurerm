package iothub

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMIotHubSharedAccessPolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_shared_access_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubSharedAccessPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubSharedAccessPolicyExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "registry_read", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "registry_write", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "service_connect", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "device_connect", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "name", "acctest"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIotHubSharedAccessPolicy_writeWithoutRead(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_shared_access_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccAzureRMIotHubSharedAccessPolicy_writeWithoutRead(data),
				ExpectError: regexp.MustCompile("If `registry_write` is set to true, `registry_read` must also be set to true"),
			},
		},
	})
}

func TestAccAzureRMIotHubSharedAccessPolicy_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_iothub_shared_access_policy", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIotHubSharedAccessPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIotHubSharedAccessPolicy_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIotHubSharedAccessPolicyExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMIotHubSharedAccessPolicy_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_iothub_shared_access_policy"),
			},
		},
	})
}

func testAccAzureRMIotHubSharedAccessPolicy_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_shared_access_policy" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  registry_read  = true
  registry_write = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMIotHubSharedAccessPolicy_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMIotHubSharedAccessPolicy_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_iothub_shared_access_policy" "import" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  registry_read  = true
  registry_write = true
}
`, template)
}

func testAccAzureRMIotHubSharedAccessPolicy_writeWithoutRead(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_iothub" "test" {
  name                = "acctestIoTHub-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku {
    name     = "B1"
    capacity = "1"
  }

  tags = {
    purpose = "testing"
  }
}

resource "azurerm_iothub_shared_access_policy" "test" {
  resource_group_name = azurerm_resource_group.test.name
  iothub_name         = azurerm_iothub.test.name
  name                = "acctest"

  registry_write = true
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testCheckAzureRMIotHubSharedAccessPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
	client := acceptance.AzureProvider.Meta().(*clients.Client).IoTHub.ResourceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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
