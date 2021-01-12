package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMNetworkSecurityGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkSecurityGroup_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_network_security_group"),
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_singleRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_singleRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_singleRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),

					// The configuration for this step contains one security_rule
					// block, which should now be reflected in the state.
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMNetworkSecurityGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),

					// The configuration for this step contains no security_rule
					// blocks at all, which means "ignore any existing security groups"
					// and thus the one from the previous step is preserved.
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.#", "1"),
				),
			},
			{
				Config: testAccAzureRMNetworkSecurityGroup_rulesExplicitZero(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),

					// The configuration for this step assigns security_rule = []
					// to state explicitly that no rules are desired, so the
					// rule from the first step should now be removed.
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.#", "0"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
					testCheckAzureRMNetworkSecurityGroupDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkSecurityGroup_withTagsUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "staging"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_addingExtraRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_singleRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.#", "1"),
				),
			},

			{
				Config: testAccAzureRMNetworkSecurityGroup_anotherRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.#", "2"),
				),
			},
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_augmented(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_augmented(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_applicationSecurityGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_applicationSecurityGroup(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.#", "1"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMNetworkSecurityGroup_deleteRule(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_group", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkSecurityGroup_singleRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMNetworkSecurityGroup_deleteRule(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkSecurityGroupExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "security_rule.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMNetworkSecurityGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SecurityGroupClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		sgName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for network security group: %q", sgName)
		}

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
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SecurityGroupClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		sgName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for network security group: %q", sgName)
		}

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
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.SecurityGroupClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

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

func testAccAzureRMNetworkSecurityGroup_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMNetworkSecurityGroup_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMNetworkSecurityGroup_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_security_group" "import" {
  name                = azurerm_network_security_group.test.name
  location            = azurerm_network_security_group.test.location
  resource_group_name = azurerm_network_security_group.test.resource_group_name
}
`, template)
}

func testAccAzureRMNetworkSecurityGroup_rulesExplicitZero(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule = []
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMNetworkSecurityGroup_singleRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMNetworkSecurityGroup_anotherRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMNetworkSecurityGroup_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMNetworkSecurityGroup_withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMNetworkSecurityGroup_augmented(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

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
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMNetworkSecurityGroup_applicationSecurityGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_application_security_group" "first" {
  name                = "acctest-first%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_application_security_group" "second" {
  name                = "acctest-second%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  security_rule {
    name                                       = "test123"
    priority                                   = 100
    direction                                  = "Inbound"
    access                                     = "Allow"
    protocol                                   = "Tcp"
    source_application_security_group_ids      = [azurerm_application_security_group.first.id]
    destination_application_security_group_ids = [azurerm_application_security_group.second.id]
    source_port_ranges                         = ["10000-40000"]
    destination_port_ranges                    = ["80", "443", "8080", "8190"]
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccAzureRMNetworkSecurityGroup_deleteRule(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acceptanceTestSecurityGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  security_rule       = []
}
`, data.RandomInteger, data.Locations.Primary)
}
