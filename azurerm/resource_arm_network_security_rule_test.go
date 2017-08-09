package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMNetworkSecurityRule_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists("azurerm_network_security_rule.test"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityRule_disappears(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_basic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists("azurerm_network_security_rule.test"),
					testCheckAzureRMNetworkSecurityRuleDisappears("azurerm_network_security_rule.test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityRule_addingRules(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_updateBasic(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists("azurerm_network_security_rule.test1"),
				),
			},

			{
				Config: testAccAzureRMNetworkSecurityRule_updateExtraRule(rInt, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists("azurerm_network_security_rule.test2"),
				),
			},
		},
	})
}

func testCheckAzureRMNetworkSecurityRuleExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		sgName := rs.Primary.Attributes["network_security_group_name"]
		sgrName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for network security rule: %s", sgName)
		}

		conn := testAccProvider.Meta().(*ArmClient).secRuleClient

		resp, err := conn.Get(resourceGroup, sgName, sgrName)
		if err != nil {
			return fmt.Errorf("Bad: Get on secRuleClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Network Security Rule %q (resource group: %q) (network security group: %q) does not exist", sgrName, sgName, resourceGroup)
		}

		return nil
	}
}

func testCheckAzureRMNetworkSecurityRuleDisappears(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		sgName := rs.Primary.Attributes["network_security_group_name"]
		sgrName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for network security rule: %s", sgName)
		}

		conn := testAccProvider.Meta().(*ArmClient).secRuleClient

		_, error := conn.Delete(resourceGroup, sgName, sgrName, make(chan struct{}))
		err := <-error
		if err != nil {
			return fmt.Errorf("Bad: Delete on secRuleClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkSecurityRuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).secRuleClient

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "azurerm_network_security_rule" {
			continue
		}

		sgName := rs.Primary.Attributes["network_security_group_name"]
		sgrName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, sgName, sgrName)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Network Security Rule still exists:\n%#v", resp.SecurityRulePropertiesFormat)
		}
	}

	return nil
}

func testAccAzureRMNetworkSecurityRule_basic(rInt int, location string) string {
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

resource "azurerm_network_security_rule" "test" {
  name                        = "test123"
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  network_security_group_name = "${azurerm_network_security_group.test.name}"
}

`, rInt, location)
}

func testAccAzureRMNetworkSecurityRule_updateBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test1" {
  name                = "acceptanceTestSecurityGroup2"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test1.name}"
}

resource "azurerm_network_security_rule" "test1" {
  name                        = "test123"
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${azurerm_resource_group.test1.name}"
  network_security_group_name = "${azurerm_network_security_group.test1.name}"
}
`, rInt, location)
}

func testAccAzureRMNetworkSecurityRule_updateExtraRule(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test1" {
  name                = "acceptanceTestSecurityGroup2"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test1.name}"
}

resource "azurerm_network_security_rule" "test1" {
  name                        = "test123"
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${azurerm_resource_group.test1.name}"
  network_security_group_name = "${azurerm_network_security_group.test1.name}"
}

resource "azurerm_network_security_rule" "test2" {
  name                        = "testing456"
  priority                    = 101
  direction                   = "Inbound"
  access                      = "Deny"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${azurerm_resource_group.test1.name}"
  network_security_group_name = "${azurerm_network_security_group.test1.name}"
}
`, rInt, location)
}
