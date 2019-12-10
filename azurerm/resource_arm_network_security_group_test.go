package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetworkSecurityGroup_basic(t *testing.T) {
	resourceName := "azurerm_network_security_group.test"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),
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

func TestAccAzureRMNetworkSecurityGroup_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_network_security_group.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkSecurityGroup_requiresImport(rInt, location),
				ExpectError: testRequiresImportError("azurerm_network_security_group"),
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_singleRule(t *testing.T) {
	resourceName := "azurerm_network_security_group.test"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_singleRule(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),
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

func TestAccAzureRMNetworkSecurityGroup_update(t *testing.T) {
	resourceName := "azurerm_network_security_group.test"
	rInt := tf.AccRandTimeInt()
	location := testLocation()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_singleRule(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),

					// The configuration for this step contains one security_rule
					// block, which should now be reflected in the state.
					resource.TestCheckResourceAttr(resourceName, "security_rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMNetworkSecurityGroup_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),

					// The configuration for this step contains no security_rule
					// blocks at all, which means "ignore any existing security groups"
					// and thus the one from the previous step is preserved.
					resource.TestCheckResourceAttr(resourceName, "security_rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMNetworkSecurityGroup_rulesExplicitZero(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),

					// The configuration for this step assigns security_rule = []
					// to state explicitly that no rules are desired, so the
					// rule from the first step should now be removed.
					resource.TestCheckResourceAttr(resourceName, "security_rule.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_disappears(t *testing.T) {
	resourceName := "azurerm_network_security_group.test"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),
					testCheckAzureRMNetworkSecurityGroupDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_withTags(t *testing.T) {
	resourceName := "azurerm_network_security_group.test"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_withTags(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(resourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAzureRMNetworkSecurityGroup_withTagsUpdate(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.environment", "staging"),
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

func TestAccAzureRMNetworkSecurityGroup_addingExtraRules(t *testing.T) {
	resourceName := "azurerm_network_security_group.test"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_singleRule(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "security_rule.#", "1"),
				),
			},

			{
				Config: testAccAzureRMNetworkSecurityGroup_anotherRule(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "security_rule.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_augmented(t *testing.T) {
	resourceName := "azurerm_network_security_group.test"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_augmented(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "security_rule.#", "1"),
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

func TestAccAzureRMNetworkSecurityGroup_applicationSecurityGroup(t *testing.T) {
	resourceName := "azurerm_network_security_group.test"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_applicationSecurityGroup(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "security_rule.#", "1"),
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

func testCheckAzureRMNetworkSecurityGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		sgName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for network security group: %q", sgName)
		}

		client := testAccProvider.Meta().(*ArmClient).Network.SecurityGroupClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, resourceGroup, sgName, "")
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Network Security Group %q (resource group: %q) does not exist", sgName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on secGroupClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkSecurityGroupDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		sgName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for network security group: %q", sgName)
		}

		client := testAccProvider.Meta().(*ArmClient).Network.SecurityGroupClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		future, err := client.Delete(ctx, resourceGroup, sgName)
		if err != nil {
			if !response.WasNotFound(future.Response()) {
				return fmt.Errorf("Error deleting NSG %q (Resource Group %q): %+v", sgName, resourceGroup, err)
			}
		}

		return nil
	}
}

func testCheckAzureRMNetworkSecurityGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).Network.SecurityGroupClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_network_security_group" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name, "")

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}
			return err
		}

		return fmt.Errorf("Network Security Group still exists:\n%#v", resp.SecurityGroupPropertiesFormat)
	}

	return nil
}

func testAccAzureRMNetworkSecurityGroup_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, rInt, location)
}

func testAccAzureRMNetworkSecurityGroup_requiresImport(rInt int, location string) string {
	template := testAccAzureRMNetworkSecurityGroup_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_security_group" "test" {
  name                = "${azurerm_network_security_group.test.name}"
  location            = "${azurerm_network_security_group.test.location}"
  resource_group_name = "${azurerm_network_security_group.test.resource_group_name}"
}
`, template)
}

func testAccAzureRMNetworkSecurityGroup_rulesExplicitZero(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule = []
}
`, rInt, location)
}

func testAccAzureRMNetworkSecurityGroup_singleRule(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "TCP"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}
`, rInt, location)
}

func testAccAzureRMNetworkSecurityGroup_anotherRule(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "testDeny"
    priority                   = 101
    direction                  = "Inbound"
    access                     = "Deny"
    protocol                   = "Udp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}
`, rInt, location)
}

func testAccAzureRMNetworkSecurityGroup_withTags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, rInt, location)
}

func testAccAzureRMNetworkSecurityGroup_withTagsUpdate(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  tags = {
    environment = "staging"
  }
}
`, rInt, location)
}

func testAccAzureRMNetworkSecurityGroup_augmented(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                         = "test123"
    priority                     = 100
    direction                    = "Inbound"
    access                       = "Allow"
    protocol                     = "Tcp"
    source_port_ranges           = ["10000-40000"]
    destination_port_ranges      = ["80", "443", "8080", "8190"]
    source_address_prefixes      = ["10.0.0.0/8", "192.168.0.0/16"]
    destination_address_prefixes = ["172.16.0.0/20", "8.8.8.8"]
  }
}
`, rInt, location)
}

func testAccAzureRMNetworkSecurityGroup_applicationSecurityGroup(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_security_group" "first" {
  name                = "acctest-first%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_application_security_group" "second" {
  name                = "acctest-second%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                                       = "test123"
    priority                                   = 100
    direction                                  = "Inbound"
    access                                     = "Allow"
    protocol                                   = "Tcp"
    source_application_security_group_ids      = ["${azurerm_application_security_group.first.id}"]
    destination_application_security_group_ids = ["${azurerm_application_security_group.second.id}"]
    source_port_ranges                         = ["10000-40000"]
    destination_port_ranges                    = ["80", "443", "8080", "8190"]
  }
}
`, rInt, location, rInt, rInt, rInt)
}
