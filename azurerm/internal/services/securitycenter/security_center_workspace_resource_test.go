package securitycenter_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type SecurityCenterWorkspaceResource struct {
}

func testAccSecurityCenterWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_workspace", "test")
	r := SecurityCenterWorkspaceResource{}

	scope := fmt.Sprintf("/subscriptions/%s", os.Getenv("ARM_SUBSCRIPTION_ID"))

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicCfg(data, scope),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue(scope),
			),
		},
		data.ImportStep(),
		{
			// reset pricing to free
			Config: SecurityCenterSubscriptionPricingResource{}.tier("Free", "VirtualMachines"),
		},
	})
}

func testAccSecurityCenterWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_workspace", "test")
	r := SecurityCenterWorkspaceResource{}
	scope := fmt.Sprintf("/subscriptions/%s", os.Getenv("ARM_SUBSCRIPTION_ID"))

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicCfg(data, scope),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue(scope),
			),
		},
		{
			Config:      r.requiresImportCfg(data, scope),
			ExpectError: acceptance.RequiresImportError("azurerm_security_center_workspace"),
		},
		{
			// reset pricing to free
			Config: SecurityCenterSubscriptionPricingResource{}.tier("Free", "VirtualMachines"),
		},
	})
}

func testAccSecurityCenterWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_workspace", "test")
	r := SecurityCenterWorkspaceResource{}
	scope := fmt.Sprintf("/subscriptions/%s", os.Getenv("ARM_SUBSCRIPTION_ID"))

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicCfg(data, scope),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue(scope),
			),
		},
		{
			Config: r.differentWorkspaceCfg(data, scope),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue(scope),
			),
		},
		data.ImportStep(),
		{
			// reset pricing to free
			Config: SecurityCenterSubscriptionPricingResource{}.tier("Free", "VirtualMachines"),
		},
	})
}

func (t SecurityCenterWorkspaceResource) Exists(ctx context.Context, clients *clients.Client, _ *terraform.InstanceState) (*bool, error) {
	workspaceSettingName := "default"

	resp, err := clients.SecurityCenter.WorkspaceClient.Get(ctx, workspaceSettingName)
	if err != nil {
		return nil, fmt.Errorf("reading Security Center Subscription Workspace Rule Set (%s): %+v", workspaceSettingName, err)
	}

	return utils.Bool(resp.WorkspaceSettingProperties != nil), nil
}

func (SecurityCenterWorkspaceResource) basicCfg(data acceptance.TestData, scope string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_subscription_pricing" "test" {
  tier          = "Standard"
  resource_type = "VirtualMachines"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sc-%d"
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

func (r SecurityCenterWorkspaceResource) requiresImportCfg(data acceptance.TestData, scope string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_security_center_workspace" "import" {
  scope        = azurerm_security_center_workspace.test.scope
  workspace_id = azurerm_security_center_workspace.test.workspace_id
}
`, r.basicCfg(data, scope))
}

func (SecurityCenterWorkspaceResource) differentWorkspaceCfg(data acceptance.TestData, scope string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_subscription_pricing" "test" {
  tier = "Standard"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sc-%d"
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
