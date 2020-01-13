package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccAzureRMDDoSProtectionPlan_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ddos_protection_plan", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDDoSProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDDoSProtectionPlan_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDDoSProtectionPlanExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_ids.#"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMDDoSProtectionPlan_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_ddos_protection_plan", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDDoSProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDDoSProtectionPlan_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDDoSProtectionPlanExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMDDoSProtectionPlan_requiresImportConfig(data),
				ExpectError: acceptance.RequiresImportError("azurerm_ddos_protection_plan"),
			},
		},
	})
}

func testAccAzureRMDDoSProtectionPlan_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ddos_protection_plan", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDDoSProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDDoSProtectionPlan_withTagsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDDoSProtectionPlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
			{
				Config: testAccAzureRMDDoSProtectionPlan_withUpdatedTagsConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDDoSProtectionPlanExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Staging"),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMDDoSProtectionPlan_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ddos_protection_plan", "test")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDDoSProtectionPlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMDDoSProtectionPlan_basicConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMDDoSProtectionPlanExists(data.ResourceName),
					testCheckAzureRMDDoSProtectionPlanDisappears(data.ResourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMDDoSProtectionPlanExists(resourceName string) resource.TestCheckFunc {
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

			return fmt.Errorf("Bad: Get on ddosProtectionPlanClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDDoSProtectionPlanDisappears(resourceName string) resource.TestCheckFunc {
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
			return fmt.Errorf("Bad: Delete on ddosProtectionPlanClient: %+v", err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Bad: waiting for Deletion on ddosProtectionPlanClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMDDoSProtectionPlanDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Network.DDOSProtectionPlansClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_ddos_protection_plan" {
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

func testAccAzureRMDDoSProtectionPlan_basicConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDDoSProtectionPlan_requiresImportConfig(data acceptance.TestData) string {
	basicConfig := testAccAzureRMDDoSProtectionPlan_basicConfig(data)
	return fmt.Sprintf(`
%s

resource "azurerm_ddos_protection_plan" "import" {
  name                = "${azurerm_ddos_protection_plan.test.name}"
  location            = "${azurerm_ddos_protection_plan.test.location}"
  resource_group_name = "${azurerm_ddos_protection_plan.test.resource_group_name}"
}
`, basicConfig)
}

func testAccAzureRMDDoSProtectionPlan_withTagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccAzureRMDDoSProtectionPlan_withUpdatedTagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  tags = {
    environment = "Staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
