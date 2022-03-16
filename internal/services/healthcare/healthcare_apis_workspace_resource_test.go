package healthcare_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type HealthCareWorkspaceResouce struct{}

func TestAccHealthCareWorkspace_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_workspace", "test")
	r := HealthCareWorkspaceResouce{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareWorkspace_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_workspace", "test")
	r := HealthCareWorkspaceResouce{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (HealthCareWorkspaceResouce) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.WorkspaceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HealthCare.HealthcareWorkspaceClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s, %+v", *id, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}
func (HealthCareWorkspaceResouce) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-health-%d"
  location = "%s"
}

resource "azurerm_healthcare_workspace" "test" {
  name                          = "acctestwk%d"
  resource_group_name           = azurerm_resource_group.test.name
  location                      = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(8))
}

func (r HealthCareWorkspaceResouce) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_workspace" "import" {
  name                = azurerm_healthcare_workspace.test.name
  resource_group_name = azurerm_healthcare_workspace.test.resource_group_name
  location            = azurerm_healthcare_workspace.test.location
}


`, r.basic(data))
}
