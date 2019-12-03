package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMDedicatedHostGroup_basic(t *testing.T) {
	resourceName := "azurerm_dedicated_host_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	rName := acctest.RandStringFromCharSet(4, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHostGroup_basic(ri, location, rName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostGroupExists(resourceName),
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

func TestAccAzureRMDedicatedHostGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_dedicated_host_group.test"
	ri := tf.AccRandTimeInt()
	rName := acctest.RandStringFromCharSet(4, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHostGroup_basic(ri, testLocation(), rName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostGroupExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMDedicatedHostGroup_requiresImport(ri, testLocation(), rName),
				ExpectError: testRequiresImportError("azurerm_dedicated_host_group"),
			},
		},
	})
}

func TestAccAzureRMDedicatedHostGroup_complete(t *testing.T) {
	resourceName := "azurerm_dedicated_host_group.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	rName := acctest.RandStringFromCharSet(4, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDedicatedHostGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDedicatedHostGroup_complete(ri, location, rName),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDedicatedHostGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "zones.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "zones.0", "1"),
					resource.TestCheckResourceAttr(resourceName, "platform_fault_domain_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.ENV", "prod"),
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

func testCheckAzureRMDedicatedHostGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Dedicated Host Group not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		client := testAccProvider.Meta().(*ArmClient).Compute.DedicatedHostGroupsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Dedicated Host Group %q (Resource Group %q) does not exist", name, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on Compute.DedicatedHostGroupsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDedicatedHostGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).Compute.DedicatedHostGroupsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dedicated_host_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		if resp, err := client.Get(ctx, resourceGroup, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on Compute.DedicatedHostGroupsClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMDedicatedHostGroup_basic(rInt int, location string, rName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dedicated_host_group" "test" {
  resource_group_name 			= "${azurerm_resource_group.test.name}"
  name 							= "%s"
  location            			= "${azurerm_resource_group.test.location}"
  platform_fault_domain_count 	= 2
}
`, rInt, location, rName)
}

func testAccAzureRMDedicatedHostGroup_requiresImport(rInt int, location string, rName string) string {
	return fmt.Sprintf(`
%s
resource "azurerm_dedicated_host_group" "import" {
  resource_group_name = "${azurerm_dedicated_host_group.test.resource_group_name}"
  name                = "${azurerm_dedicated_host_group.test.name}"
  location            = "${azurerm_dedicated_host_group.test.location}"
  platform_fault_domain_count 	= 2
}
`, testAccAzureRMDedicatedHostGroup_basic(rInt, location, rName))
}

func testAccAzureRMDedicatedHostGroup_complete(rInt int, location string, rName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dedicated_host_group" "test" {
  resource_group_name 			= "${azurerm_resource_group.test.name}"
  name 							= "%s"
  location            			= "${azurerm_resource_group.test.location}"
  platform_fault_domain_count 	= 2
  zones = ["1"]
  tags = {
    ENV = "prod"
  }
}
`, rInt, location, rName)
}
