// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_center_workspace_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apicenter/2024-03-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ApiCenterWorkspaceResource struct{}

func TestAccApiCenterWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_workspace", "test")
	r := ApiCenterWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("title").MatchesOtherKey(check.That(data.ResourceName).Key("name")),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiCenterWorkspace_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_workspace", "test")
	r := ApiCenterWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiCenterWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_center_workspace", "test")
	r := ApiCenterWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (ApiCenterWorkspaceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workspaces.ParseWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiCenter.WorkspacesClient.Get(ctx, *id)
	if response.WasNotFound(resp.HttpResponse) {
		return pointer.To(false), nil
	}
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (ApiCenterWorkspaceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apicenter-%d"
  location = "%s"
}

resource "azurerm_api_center_service" "test" {
  name                = "acctestApiCSvc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ApiCenterWorkspaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_center_workspace" "test" {
  name                  = "acctest%d"
  api_center_service_id = azurerm_api_center_service.test.id
}
`, r.template(data), data.RandomInteger)
}

func (r ApiCenterWorkspaceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%s

resource "azurerm_api_center_workspace" "test" {
  name                  = "acctest%d"
  api_center_service_id = azurerm_api_center_service.test.id
  title                 = "testtitle"
  description           = "testing workspace"
}
`, r.template(data), data.RandomInteger)
}

func (r ApiCenterWorkspaceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_center_workspace" "import" {
  name                  = azurerm_api_center_workspace.test.name
  api_center_service_id = azurerm_api_center_workspace.test.api_center_service_id
  title                 = azurerm_api_center_workspace.test.title
}
`, r.basic(data))
}
