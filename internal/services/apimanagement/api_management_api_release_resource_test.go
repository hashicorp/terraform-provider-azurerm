package apimanagement_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApiManagementApiReleaseResource struct {
}

func TestAccApiManagementApiRelease_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_release", "test")
	r := ApiManagementApiReleaseResource{}

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

func TestAccApiManagementApiRelease_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_release", "test")
	r := ApiManagementApiReleaseResource{}

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

func TestAccApiManagementApiRelease_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_release", "test")
	r := ApiManagementApiReleaseResource{}

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

func TestAccApiManagementApiRelease_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_management_api_release", "test")
	r := ApiManagementApiReleaseResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ApiManagementApiReleaseResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ApiReleaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.ApiManagement.ApiReleasesClient.Get(ctx, id.ResourceGroup, id.ServiceName, id.ApiName, id.ReleaseName)
	if err != nil {
		return nil, fmt.Errorf("reading ApiManagement Api Release (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (r ApiManagementApiReleaseResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_release" "test" {
  name   = "acctest-ApiRelease-%d"
  api_id = azurerm_api_management_api.test.id
}
`, ApiManagementApiResource{}.basic(data), data.RandomInteger)
}

func (r ApiManagementApiReleaseResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_release" "test" {
  name   = "acctest-ApiRelease-%d"
  api_id = azurerm_api_management_api.test.id
  notes  = "Release 1.0"
}
`, ApiManagementApiResource{}.basic(data), data.RandomInteger)
}
func (r ApiManagementApiReleaseResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_release" "test" {
  name   = "acctest-ApiRelease-%d"
  api_id = azurerm_api_management_api.test.id
  notes  = "Release 2.0"
}
`, ApiManagementApiResource{}.basic(data), data.RandomInteger)
}

func (r ApiManagementApiReleaseResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_api_management_api_release" "import" {
  name   = azurerm_api_management_api_release.test.name
  api_id = azurerm_api_management_api_release.test.api_id
}
`, r.basic(data))
}
