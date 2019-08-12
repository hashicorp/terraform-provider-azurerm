package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func testCheckAzureRMPublicIPPrefixExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		publicIpPrefixName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for public ip prefix: %s", publicIpPrefixName)
		}

		client := testAccProvider.Meta().(*ArmClient).network.PublicIPPrefixesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := client.Get(ctx, resourceGroup, publicIpPrefixName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get on publicIPPrefixClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Public IP Prefix %q (resource group: %q) does not exist", publicIpPrefixName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMPublicIPPrefixDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		publicIpPrefixName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for public ip prefix: %s", publicIpPrefixName)
		}

		client := testAccProvider.Meta().(*ArmClient).network.PublicIPPrefixesClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		future, err := client.Delete(ctx, resourceGroup, publicIpPrefixName)
		if err != nil {
			return fmt.Errorf("Error deleting Public IP Prefix %q (Resource Group %q): %+v", publicIpPrefixName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Public IP Prefix %q (Resource Group %q): %+v", publicIpPrefixName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMPublicIPPrefixDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).network.PublicIPPrefixesClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_public_ip_prefix" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Public IP Prefix still exists:\n%#v", resp.PublicIPPrefixPropertiesFormat)
		}
	}

	return nil
}

func TestAccAzureRMPublicIpPrefix_basic(t *testing.T) {
	resourceName := "azurerm_public_ip_prefix.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPublicIPPrefix_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIPPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIPPrefixExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_prefix"),
					resource.TestCheckResourceAttr(resourceName, "prefix_length", "28"),
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

func TestAccAzureRMPublicIpPrefix_prefixLength(t *testing.T) {
	resourceName := "azurerm_public_ip_prefix.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPublicIPPrefix_prefixLength(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIPPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIPPrefixExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "ip_prefix"),
					resource.TestCheckResourceAttr(resourceName, "prefix_length", "31"),
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

func TestAccAzureRMPublicIpPrefix_update(t *testing.T) {
	resourceName := "azurerm_public_ip_prefix.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()
	preConfig := testAccAzureRMPublicIPPrefix_withTags(ri, location)
	postConfig := testAccAzureRMPublicIPPrefix_withTagsUpdate(ri, location)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIPPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIPPrefixExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIPPrefixExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
				),
			},
		},
	})
}

func TestAccAzureRMPublicIpPrefix_disappears(t *testing.T) {
	resourceName := "azurerm_public_ip_prefix.test"
	ri := tf.AccRandTimeInt()
	config := testAccAzureRMPublicIPPrefix_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMPublicIPPrefixDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMPublicIPPrefixExists(resourceName),
					testCheckAzureRMPublicIPPrefixDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAzureRMPublicIPPrefix_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPPrefix_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPPrefix_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "staging"
  }
}
`, rInt, location, rInt)
}

func testAccAzureRMPublicIPPrefix_prefixLength(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  prefix_length = 31
}
`, rInt, location, rInt)
}
