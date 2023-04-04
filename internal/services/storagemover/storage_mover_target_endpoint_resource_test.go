package storagemover_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/endpoints"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageMoverTargetEndpointTestResource struct{}

func TestAccstoragemoverEndpoint_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storagemover_endpoint", "test")
	r := StorageMoverTargetEndpointTestResource{}
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

func TestAccstoragemoverEndpoint_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storagemover_endpoint", "test")
	r := StorageMoverTargetEndpointTestResource{}
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

func TestAccstoragemoverEndpoint_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storagemover_endpoint", "test")
	r := StorageMoverTargetEndpointTestResource{}
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

func TestAccstoragemoverEndpoint_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storagemover_endpoint", "test")
	r := StorageMoverTargetEndpointTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
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
	})
}

func (r StorageMoverTargetEndpointTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := endpoints.ParseEndpointID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.storagemover.EndpointsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r StorageMoverTargetEndpointTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_storagemover_storage_mover" "test" {
  name                = "acctest-ssm-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StorageMoverTargetEndpointTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_storagemover_endpoint" "test" {
  name                          = "acctest-se-%d"
  StorageMoverId = azurerm_storagemover_storage_mover.test.id
  endpoint_type                 = ""
}
`, template, data.RandomInteger)
}

func (r StorageMoverTargetEndpointTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_storagemover_endpoint" "import" {
  name                          = azurerm_storagemover_endpoint.test.name
  StorageMoverId = azurerm_storagemover_storage_mover.test.id
  endpoint_type                 = ""
}
`, config)
}

func (r StorageMoverTargetEndpointTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_storagemover_endpoint" "test" {
  name                          = "acctest-se-%d"
  StorageMoverId = azurerm_storagemover_storage_mover.test.id
  description                   = ""
  endpoint_type                 = ""

}
`, template, data.RandomInteger)
}

func (r StorageMoverTargetEndpointTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_storagemover_endpoint" "test" {
  name                          = "acctest-se-%d"
  StorageMoverId = azurerm_storagemover_storage_mover.test.id
  description                   = ""
  endpoint_type                 = ""

}
`, template, data.RandomInteger)
}
