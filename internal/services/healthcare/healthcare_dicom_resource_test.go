package healthcare_test

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type HealthCareDicomResource struct{}

func TestAccHealthCareDicom_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_dicom_service", "test")
	r := HealthCareDicomResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r)),
		},
		data.ImportStep(),
	})
}

func TestAccHealthCareDicom_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_healthcare_dicom_service", "test")
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

func (HealthCareDicomResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DicomServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.HealthCare.HealthcareWorkspaceDicomServiceClient.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s, %+v", *id, err)
	}

	return utils.Bool(resp.DicomServiceProperties != nil), nil
}

func (HealthCareDicomResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-hcwk-%d"
  location = "%s"
}

resource "azurerm_healthcare_workspace" "test" {
  name                = "acctestwk%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_healthcare_dicom_service" "test" {
  name         = "acctest-dicom%d"
  workspace_id = azurerm_healthcare_workspace.test.id
  location     = "east us"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(8), data.RandomIntOfLength(8))
}

func (r HealthCareDicomResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_healthcare_dicom_service" "import" {
  name         = azurerm_healthcare_dicom_service.test.name
  workspace_id = azurerm_healthcare_dicom_service.test.workspace_id
  location     = azurerm_healthcare_dicom_service.test.location
}
`, r.basic(data))
}
