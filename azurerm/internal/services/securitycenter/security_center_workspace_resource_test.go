package securitycenter_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func testAccSecurityCenterWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_workspace", "test")

	scope := fmt.Sprintf("/subscriptions/%s", os.Getenv("ARM_SUBSCRIPTION_ID"))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSecurityCenterWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterWorkspace_basicCfg(data, scope),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityCenterWorkspaceExists(),
					resource.TestCheckResourceAttr(data.ResourceName, "scope", scope),
				),
			},
			data.ImportStep(),
			{
				// reset pricing to free
				Config: testAccSecurityCenterSubscriptionPricing_tier("Free", "VirtualMachines"),
			},
		},
	})
}

func testAccSecurityCenterWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_workspace", "test")
	scope := fmt.Sprintf("/subscriptions/%s", os.Getenv("ARM_SUBSCRIPTION_ID"))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSecurityCenterWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterWorkspace_basicCfg(data, scope),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityCenterWorkspaceExists(),
					resource.TestCheckResourceAttr(data.ResourceName, "scope", scope),
				),
			},
			{
				Config:      testAccSecurityCenterWorkspace_requiresImportCfg(data, scope),
				ExpectError: acceptance.RequiresImportError("azurerm_security_center_workspace"),
			},
			{
				// reset pricing to free
				Config: testAccSecurityCenterSubscriptionPricing_tier("Free", "VirtualMachines"),
			},
		},
	})
}

func testAccSecurityCenterWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_workspace", "test")
	scope := fmt.Sprintf("/subscriptions/%s", os.Getenv("ARM_SUBSCRIPTION_ID"))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckSecurityCenterWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityCenterWorkspace_basicCfg(data, scope),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityCenterWorkspaceExists(),
					resource.TestCheckResourceAttr(data.ResourceName, "scope", scope),
				),
			},
			{
				Config: testAccSecurityCenterWorkspace_differentWorkspaceCfg(data, scope),
				Check: resource.ComposeTestCheckFunc(
					testCheckSecurityCenterWorkspaceExists(),
					resource.TestCheckResourceAttr(data.ResourceName, "scope", scope),
				),
			},
			data.ImportStep(),
			{
				// reset pricing to free
				Config: testAccSecurityCenterSubscriptionPricing_tier("Free", "VirtualMachines"),
			},
		},
	})
}

func testCheckSecurityCenterWorkspaceExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.WorkspaceClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		workspaceSettingName := "default"
		resp, err := client.Get(ctx, workspaceSettingName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Security Center Subscription Workspace %q was not found: %+v", workspaceSettingName, err)
			}

			return fmt.Errorf("Bad: Get: %+v", err)
		}

		return nil
	}
}

func testCheckSecurityCenterWorkspaceDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).SecurityCenter.WorkspaceClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_security_center_workspace" {
			continue
		}

		workspaceSettingName := "default"
		resp, err := client.Get(ctx, workspaceSettingName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("security center workspace settings still exists")
	}

	return nil
}

func testAccSecurityCenterWorkspace_basicCfg(data acceptance.TestData, scope string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_subscription_pricing" "test" {
  tier          = "Standard"
  resource_type = "VirtualMachines"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-%d-1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_security_center_workspace" "test" {
  scope        = "%s"
  workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, scope)
}

func testAccSecurityCenterWorkspace_requiresImportCfg(data acceptance.TestData, scope string) string {
	template := testAccSecurityCenterWorkspace_basicCfg(data, scope)
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_workspace" "import" {
  scope        = azurerm_security_center_workspace.test.scope
  workspace_id = azurerm_security_center_workspace.test.workspace_id
}
`, template)
}

func testAccSecurityCenterWorkspace_differentWorkspaceCfg(data acceptance.TestData, scope string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_subscription_pricing" "test" {
  tier = "Standard"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test2" {
  name                = "acctest-%d-2"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_security_center_workspace" "test" {
  scope        = "%s"
  workspace_id = azurerm_log_analytics_workspace.test2.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, scope)
}
