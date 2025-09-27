// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspacepolicy"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiManagementWorkspacePolicyTestResource struct{}

func TestAccApiManagementWorkspacePolicy_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_policy", "test")
	r := ApiManagementWorkspacePolicyTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspacePolicy_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_policy", "test")
	r := ApiManagementWorkspacePolicyTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("xml_link"),
	})
}

func TestAccApiManagementWorkspacePolicy_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace_policy", "test")
	r := ApiManagementWorkspacePolicyTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("xml_link"),
		{
			Config: r.updated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("xml_link"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("xml_link"),
	})
}

func (ApiManagementWorkspacePolicyTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workspacepolicy.ParseWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.WorkspacePolicyClient.Get(ctx, *id, workspacepolicy.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ApiManagementWorkspacePolicyTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_policy" "test" {
  api_management_workspace_id = azurerm_api_management_workspace.test.id

  xml_content = <<XML
<policies>
  <inbound>
    <find-and-replace from="xyz" to="abc" />
  </inbound>
</policies>
XML
}
`, r.template(data))
}

func (r ApiManagementWorkspacePolicyTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_policy" "test" {
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  xml_link                    = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/main/internal/services/apimanagement/testdata/api_management_policy_test.xml"
}
`, r.template(data))
}

func (r ApiManagementWorkspacePolicyTestResource) updated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_management_workspace_policy" "test" {
  api_management_workspace_id = azurerm_api_management_workspace.test.id
  xml_link                    = "https://raw.githubusercontent.com/hashicorp/terraform-provider-azurerm/main/internal/services/apimanagement/testdata/api_management_policy_update_test.xml"
}
`, r.template(data))
}

func (ApiManagementWorkspacePolicyTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Premium_1"
}

resource "azurerm_api_management_workspace" "test" {
  name              = "acctestws-%d"
  api_management_id = azurerm_api_management.test.id
  display_name      = "Test Workspace"
  description       = "Test workspace for policy"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
