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

func TestValidateDevTestVirtualNetworkName(t *testing.T) {
	validNames := []string{
		"valid-name",
		"valid02-name",
		"validName1",
		"-validname1",
		"valid_name",
		"double-hyphen--valid",
	}
	for _, v := range validNames {
		_, errors := validateDevTestVirtualNetworkName()(v, "example")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid Dev Test Virtual Network Name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"invalid!",
		"!@£",
	}
	for _, v := range invalidNames {
		_, errors := validateDevTestVirtualNetworkName()(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid Dev Test Virtual Network Name", v)
		}
	}
}

func TestAccAzureRMDevTestVirtualNetwork_basic(t *testing.T) {
	resourceName := "azurerm_dev_test_virtual_network.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDevTestVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestVirtualNetwork_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestVirtualNetworkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccAzureRMDevTestVirtualNetwork_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_dev_test_virtual_network.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDevTestVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestVirtualNetwork_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestVirtualNetworkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				Config:      testAccAzureRMDevTestVirtualNetwork_requiresImport(rInt, location),
				ExpectError: testRequiresImportError("azurerm_dev_test_virtual_network"),
			},
		},
	})
}

func TestAccAzureRMDevTestVirtualNetwork_subnet(t *testing.T) {
	resourceName := "azurerm_dev_test_virtual_network.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDevTestVirtualNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDevTestVirtualNetwork_subnets(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDevTestVirtualNetworkExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "subnet.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "subnet.0.use_public_ip_address", "Deny"),
					resource.TestCheckResourceAttr(resourceName, "subnet.0.use_in_virtual_machine_creation", "Allow"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func testCheckAzureRMDevTestVirtualNetworkExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		virtualNetworkName := rs.Primary.Attributes["name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		conn := testAccProvider.Meta().(*ArmClient).devTestLabs.VirtualNetworksClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, labName, virtualNetworkName, "")
		if err != nil {
			return fmt.Errorf("Bad: Get devTestVirtualNetworksClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: DevTest Virtual Network %q (Lab %q / Resource Group: %q) does not exist", virtualNetworkName, labName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMDevTestVirtualNetworkDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).devTestLabs.VirtualNetworksClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_dev_test_virtual_network" {
			continue
		}

		virtualNetworkName := rs.Primary.Attributes["name"]
		labName := rs.Primary.Attributes["lab_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, labName, virtualNetworkName, "")

		if err != nil {
			if resp.StatusCode == http.StatusNotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("DevTest Virtual Network still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMDevTestVirtualNetwork_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dev_test_virtual_network" "test" {
  name                = "acctestdtvn%d"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMDevTestVirtualNetwork_requiresImport(rInt int, location string) string {
	template := testAccAzureRMDevTestVirtualNetwork_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_dev_test_virtual_network" "import" {
  name                = "${azurerm_dev_test_virtual_network.test.name}"
  lab_name            = "${azurerm_dev_test_virtual_network.test.lab_name}"
  resource_group_name = "${azurerm_dev_test_virtual_network.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMDevTestVirtualNetwork_subnets(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dev_test_lab" "test" {
  name                = "acctestdtl%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dev_test_virtual_network" "test" {
  name                = "acctestdtvn%d"
  lab_name            = "${azurerm_dev_test_lab.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  subnet {
    use_public_ip_address           = "Deny"
    use_in_virtual_machine_creation = "Allow"
  }
}
`, rInt, location, rInt, rInt)
}
