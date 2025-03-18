package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/workspace"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApiManagementWorkspaceResource struct{}

func TestAccApiManagementWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace", "test")
	r := ApiManagementWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("workspace_name").HasValue("acctest-workspace"),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("service_name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspace_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace", "test")
	r := ApiManagementWorkspaceResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("workspace_name").HasValue("acctest-workspace-updated"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccApiManagementWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_workspace", "test")
	r := ApiManagementWorkspaceResource{}

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

func (r ApiManagementWorkspaceResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := workspace.ParseWorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.ApiManagement.WorkspaceClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r ApiManagementWorkspaceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestapim%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name           = "Developer_1"
}

resource "azurerm_api_management_workspace" "test" {
  name                = "acctest-ws-%d"
  workspace_name      = "acctest-workspace"
  service_name        = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementWorkspaceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-apim-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestapim%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name           = "Developer_1"
}

resource "azurerm_api_management_workspace" "test" {
  name                = "acctest-ws-%d"
  workspace_name      = "acctest-workspace-updated"
  service_name        = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ApiManagementWorkspaceResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_workspace" "import" {
  name                = azurerm_api_management_workspace.test.name
  workspace_name      = azurerm_api_management_workspace.test.workspace_name
  service_name        = azurerm_api_management_workspace.test.service_name
  resource_group_name = azurerm_api_management_workspace.test.resource_group_name
}
`, template)
}
