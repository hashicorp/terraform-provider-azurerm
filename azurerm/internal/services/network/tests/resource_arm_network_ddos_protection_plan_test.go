package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// NOTE: this is a test group to avoid each test case to run in parallel, since Azure only allows one DDoS Protection
// Plan per region.
func TestAccAzureRMNetworkDDoSProtectionPlan(t *testing.T) {
	testCases := map[string]map[string]func(t *testing.T){
		"normal": {
			"basic":          testAccAzureRMNetworkDDoSProtectionPlan_basic,
			"requiresImport": testAccAzureRMNetworkDDoSProtectionPlan_requiresImport,
			"withTags":       testAccAzureRMNetworkDDoSProtectionPlan_withTags,
			"disappears":     testAccAzureRMNetworkDDoSProtectionPlan_disappears,
		},
		"datasource": {
			"basic": testAccAzureRMNetworkDDoSProtectionPlanDataSource_basic,
		},
	}

	for group, steps := range testCases {
		t.Run(group, func(t *testing.T) {
			for name, tc := range steps {
				t.Run(name, func(t *testing.T) {
					tc(t)
				})
			}
		})
	}
}

func testAccAzureRMNetworkDDoSProtectionPlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_protection_plan", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkDDoSProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkDDoSProtectionPlan_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkDDoSProtectionPlanExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_ids.#"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkDDoSProtectionPlan_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_protection_plan", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkDDoSProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkDDoSProtectionPlan_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkDDoSProtectionPlanExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMNetworkDDoSProtectionPlan_requiresImportConfig(data),
				ExpectError: acceptance.RequiresImportError("azurerm_network_ddos_protection_plan"),
			},
		},
	})
}

func testAccAzureRMNetworkDDoSProtectionPlan_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_protection_plan", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkDDoSProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkDDoSProtectionPlan_withTagsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkDDoSProtectionPlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMNetworkDDoSProtectionPlan_withUpdatedTagsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkDDoSProtectionPlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Staging"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMNetworkDDoSProtectionPlan_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_ddos_protection_plan", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMNetworkDDoSProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMNetworkDDoSProtectionPlan_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMNetworkDDoSProtectionPlanExists(data.ResourceName),
					testCheckAzureRMNetworkDDoSProtectionPlanDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMNetworkDDoSProtectionPlanExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.DDOSProtectionPlansClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DDoS Protection Plan: %q", name)
		}

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: DDoS Protection Plan %q (Resource Group: %q) does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on NetworkDDoSProtectionPlanClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkDDoSProtectionPlanDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.DDOSProtectionPlansClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for DDoS Protection Plan: %q", name)
		}

		future, err := client.Delete(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Delete on NetworkDDoSProtectionPlanClient: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Bad: waiting for Deletion on NetworkDDoSProtectionPlanClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMNetworkDDoSProtectionPlanDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.DDOSProtectionPlansClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_network_ddos_protection_plan" {
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

		return fmt.Errorf("DDoS Protection Plan still exists:\n%#v", resp.DdosProtectionPlanPropertiesFormat)
	}

	return nil
}

func testAccAzureRMNetworkDDoSProtectionPlan_basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMNetworkDDoSProtectionPlan_requiresImportConfig(data acceptance.TestData) string {
	basicConfig := testAccAzureRMNetworkDDoSProtectionPlan_basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_network_ddos_protection_plan" "import" {
  name                = azurerm_network_ddos_protection_plan.test.name
  location            = azurerm_network_ddos_protection_plan.test.location
  resource_group_name = azurerm_network_ddos_protection_plan.test.resource_group_name
}
`, basicConfig)
}

func testAccAzureRMNetworkDDoSProtectionPlan_withTagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMNetworkDDoSProtectionPlan_withUpdatedTagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
