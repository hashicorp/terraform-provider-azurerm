package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMVirtualWan_basic(t *testing.T) {
	resourceName := "azurerm_virtual_wan.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualWanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualWan_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualWanExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "disable_vpn_encryption"),
					resource.TestCheckResourceAttrSet(resourceName, "allow_branch_to_branch_traffic"),
					resource.TestCheckResourceAttrSet(resourceName, "allow_vnet_to_vnet_traffic"),
					resource.TestCheckResourceAttrSet(resourceName, "office365_local_breakout_category"),
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
func TestAccAzureRMVirtualWan_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}
	resourceName := "azurerm_virtual_wan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualWanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualWan_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualWanExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "disable_vpn_encryption"),
					resource.TestCheckResourceAttrSet(resourceName, "security_provider_name"),
					resource.TestCheckResourceAttrSet(resourceName, "allow_branch_to_branch_traffic"),
					resource.TestCheckResourceAttrSet(resourceName, "allow_vnet_to_vnet_traffic"),
					resource.TestCheckResourceAttrSet(resourceName, "office365_local_breakout_category"),
				),
			},
			{
				Config:      testAccAzureRMVirtualWan_requiresImport(ri, testLocation()),
				ExpectError: testRequiresImportError("azurerm_virtual_wan"),
			},
		},
	})
}

func TestAccAzureRMVirtualWan_complete(t *testing.T) {
	resourceName := "azurerm_virtual_wan.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMVirtualWanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMVirtualWan_complete(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMVirtualWanExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "disable_vpn_encryption", "false"),
					resource.TestCheckResourceAttr(resourceName, "security_provider_name", ""),
					resource.TestCheckResourceAttr(resourceName, "allow_branch_to_branch_traffic", "true"),
					resource.TestCheckResourceAttr(resourceName, "allow_vnet_to_vnet_traffic", "true"),
					resource.TestCheckResourceAttr(resourceName, "office365_local_breakout_category", "All"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.Hello", "There"),
					resource.TestCheckResourceAttr(resourceName, "tags.World", "Example"),
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

func testCheckAzureRMVirtualWanDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).network.VirtualWanClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_virtual_wan" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Virtual WAN still exists:\n%+v", resp)
		}
	}

	return nil
}

func testCheckAzureRMVirtualWanExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		virtualWanName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Virtual WAN: %s", virtualWanName)
		}

		client := testAccProvider.Meta().(*ArmClient).network.VirtualWanClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, virtualWanName)
		if err != nil {
			return fmt.Errorf("Bad: Get on virtualWanClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Virtual WAN %q (resource group: %q) does not exist", virtualWanName, resourceGroup)
		}

		return nil
	}
}

func testAccAzureRMVirtualWan_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}
`, rInt, location, rInt)
}

func testAccAzureRMVirtualWan_requiresImport(rInt int, location string) string {
	template := testAccAzureRMVirtualWan_basic(rInt, location)

	return fmt.Sprintf(`
%s

resource "azurerm_virtual_wan" "import" {
  name                = "${azurerm_virtual_wan.test.name}"
  resource_group_name = "${azurerm_virtual_wan.test.resource_group_name}"
  location            = "${azurerm_virtual_wan.test.location}"
}
`, template)
}

func testAccAzureRMVirtualWan_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctestvwan%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"

  disable_vpn_encryption = false

  allow_branch_to_branch_traffic = true
  allow_vnet_to_vnet_traffic     = true

  office365_local_breakout_category = "All"

  tags = {
    Hello = "There"
    World = "Example"
  }
}
`, rInt, location, rInt)
}
