// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SecurityCenterWorkspaceResource struct{}

func testAccSecurityCenterWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_workspace", "test")
	r := SecurityCenterWorkspaceResource{}

	scope := fmt.Sprintf("/subscriptions/%s", os.Getenv("ARM_SUBSCRIPTION_ID"))

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicCfg(data, scope),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue(scope),
			),
		},
		data.ImportStep(),
	})
}

func testAccSecurityCenterWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_workspace", "test")
	r := SecurityCenterWorkspaceResource{}
	scope := fmt.Sprintf("/subscriptions/%s", os.Getenv("ARM_SUBSCRIPTION_ID"))

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicCfg(data, scope),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue(scope),
			),
		},
		{
			Config: r.differentWorkspaceCfg(data, scope),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue(scope),
			),
		},
		data.ImportStep(),
	})
}

func testAccSecurityCenterWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_security_center_workspace", "test")
	r := SecurityCenterWorkspaceResource{}
	scope := fmt.Sprintf("/subscriptions/%s", os.Getenv("ARM_SUBSCRIPTION_ID"))

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicCfg(data, scope),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("scope").HasValue(scope),
			),
		},
		{
			Config:      r.requiresImportCfg(data, scope),
			ExpectError: acceptance.RequiresImportError("azurerm_security_center_workspace"),
		},
	})
}

func (SecurityCenterWorkspaceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.SecurityCenter.WorkspaceClient.Get(ctx, id.WorkspaceSettingName)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.WorkspaceSettingProperties != nil), nil
}

func (SecurityCenterWorkspaceResource) basicCfg(data acceptance.TestData, scope string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_security_center_subscription_pricing" "test" {
  tier          = "Free"
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
  tier = "Free"
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
