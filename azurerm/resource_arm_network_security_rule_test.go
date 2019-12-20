package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetworkSecurityRule_basic(t *testing.T) {
	resourceName := "azurerm_network_security_rule.test"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_basic(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists(resourceName),
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

func TestAccAzureRMNetworkSecurityRule_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_network_security_rule.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkSecurityRule_requiresImport(rInt, location),
				ExpectError: acceptance.RequiresImportError("azurerm_network_security_rule"),
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityRule_disappears(t *testing.T) {
	resourceGroup := "azurerm_network_security_rule.test"
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_basic(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists(resourceGroup),
					testCheckAzureRMNetworkSecurityRuleDisappears(resourceGroup),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityRule_addingRules(t *testing.T) {
	rInt := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_updateBasic(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists("azurerm_network_security_rule.test1"),
				),
			},

			{
				Config: testAccAzureRMNetworkSecurityRule_updateExtraRule(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists("azurerm_network_security_rule.test2"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityRule_augmented(t *testing.T) {
	resourceName := "azurerm_network_security_rule.test1"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_augmented(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists(resourceName),
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

func TestAccAzureRMNetworkSecurityRule_applicationSecurityGroups(t *testing.T) {
	resourceName := "azurerm_network_security_rule.test1"
	rInt := tf.AccRandTimeInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityRule_applicationSecurityGroups(rInt, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityRuleExists(resourceName),
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

func testCheckAzureRMNetworkSecurityRuleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		sgName := rs.Primary.Attributes["network_security_group_name"]
		sgrName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for network security rule: %q", sgName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SecurityRuleClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, sgName, sgrName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Network Security Rule %q (resource group: %q) (network security group: %q) does not exist", sgrName, sgName, resourceGroup)
			}
			return fmt.Errorf("Error retrieving Network Security Rule %q (NSG %q / Resource Group %q): %+v", sgrName, sgName, resourceGroup, err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkSecurityRuleDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		sgName := rs.Primary.Attributes["network_security_group_name"]
		sgrName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for network security rule: %s", sgName)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SecurityRuleClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		future, err := client.Delete(ctx, resourceGroup, sgName, sgrName)
		if err != nil {
			if !response.WasNotFound(future.Response()) {
				return fmt.Errorf("Error deleting Network Security Rule %q (NSG %q / Resource Group %q): %+v", sgrName, sgName, resourceGroup, err)
			}
		}

		return nil
	}
}

func testCheckAzureRMNetworkSecurityRuleDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SecurityRuleClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_network_security_rule" {
			continue
		}

		sgName := rs.Primary.Attributes["network_security_group_name"]
		sgrName := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, sgName, sgrName)

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
  network_security_group_name = "${azurerm_network_security_group.test.name}"
  resource_group_name         = "${azurerm_resource_group.test.name}"
  priority                    = 100
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
}
`, rInt, location)
}

func testAccAzureRMNetworkSecurityRule_requiresImport(rInt int, location string) string {
	template := testAccAzureRMNetworkSecurityRule_basic(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_network_security_rule" "import" {
  name                        = "${azurerm_network_security_rule.test.name}"
  network_security_group_name = "${azurerm_network_security_rule.test.network_security_group_name}"
  resource_group_name         = "${azurerm_network_security_rule.test.resource_group_name}"
  priority                    = "${azurerm_network_security_rule.test.priority}"
  direction                   = "${azurerm_network_security_rule.test.direction}"
  access                      = "${azurerm_network_security_rule.test.access}"
  protocol                    = "${azurerm_network_security_rule.test.protocol}"
  source_port_range           = "${azurerm_network_security_rule.test.source_port_range}"
  destination_port_range      = "${azurerm_network_security_rule.test.destination_port_range}"
  source_address_prefix       = "${azurerm_network_security_rule.test.source_address_prefix}"
  destination_address_prefix  = "${azurerm_network_security_rule.test.destination_address_prefix}"
}
`, template)
}

func testAccAzureRMNetworkSecurityRule_updateBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test1" {
  name                = "acceptanceTestSecurityGroup2"
  location            = "${azurerm_resource_group.test1.location}"
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
  location            = "${azurerm_resource_group.test1.location}"
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
  protocol                    = "Icmp"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${azurerm_resource_group.test1.name}"
  network_security_group_name = "${azurerm_network_security_group.test1.name}"
}
`, rInt, location)
}

func testAccAzureRMNetworkSecurityRule_augmented(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test1" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test1" {
  name                = "acceptanceTestSecurityGroup2"
  location            = "${azurerm_resource_group.test1.location}"
  resource_group_name = "${azurerm_resource_group.test1.name}"
}

resource "azurerm_network_security_rule" "test1" {
  name                         = "test123"
  priority                     = 100
  direction                    = "Outbound"
  access                       = "Allow"
  protocol                     = "Tcp"
  source_port_ranges           = ["10000-40000"]
  destination_port_ranges      = ["80", "443", "8080", "8190"]
  source_address_prefixes      = ["10.0.0.0/8", "192.168.0.0/16"]
  destination_address_prefixes = ["172.16.0.0/20", "8.8.8.8"]
  resource_group_name          = "${azurerm_resource_group.test1.name}"
  network_security_group_name  = "${azurerm_network_security_group.test1.name}"
}
`, rInt, location)
}

func testAccAzureRMNetworkSecurityRule_applicationSecurityGroups(rInt int, location string) string {
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
}

resource "azurerm_network_security_rule" "test1" {
  name                                       = "test123"
  resource_group_name                        = "${azurerm_resource_group.test.name}"
  network_security_group_name                = "${azurerm_network_security_group.test.name}"
  priority                                   = 100
  direction                                  = "Outbound"
  access                                     = "Allow"
  protocol                                   = "Tcp"
  source_application_security_group_ids      = ["${azurerm_application_security_group.first.id}"]
  destination_application_security_group_ids = ["${azurerm_application_security_group.second.id}"]
  source_port_ranges                         = ["10000-40000"]
  destination_port_ranges                    = ["80", "443", "8080", "8190"]
}
`, rInt, location, rInt, rInt, rInt)
}
