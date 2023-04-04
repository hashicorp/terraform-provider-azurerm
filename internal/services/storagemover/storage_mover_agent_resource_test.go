package storagemover_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/storagemover/2023-03-01/agents"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type StorageMoverAgentTestResource struct{}

func TestAccStorageMoverAgent_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_agent", "test")
	r := StorageMoverAgentTestResource{}
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

func TestAccStorageMoverAgent_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_agent", "test")
	r := StorageMoverAgentTestResource{}
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

func TestAccStorageMoverAgent_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_agent", "test")
	r := StorageMoverAgentTestResource{}
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

func TestAccStorageMoverAgent_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_storage_mover_agent", "test")
	r := StorageMoverAgentTestResource{}
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

func (r StorageMoverAgentTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := agents.ParseAgentID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.StorageMover.AgentsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r StorageMoverAgentTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}

resource "azurerm_storage_mover" "test" {
  name                = "acctest-ssm-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r StorageMoverAgentTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`

%s

resource "azurerm_storage_mover_agent" "test" {
  name             = "acctest-sa-%d"
  storage_mover_id = azurerm_storage_mover.test.id
  arc_resource_id  = "${azurerm_resource_group.test.id}/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"
  arc_vm_uuid      = "3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"
}
`, template, data.RandomInteger)
}

func (r StorageMoverAgentTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`

%s

resource "azurerm_storage_mover_agent" "import" {
  name             = azurerm_storage_mover_agent.test.name
  storage_mover_id = azurerm_storage_mover_agent.test.storage_mover_id
  arc_resource_id  = azurerm_storage_mover_agent.test.arc_resource_id
  arc_vm_uuid      = azurerm_storage_mover_agent.test.arc_vm_uuid
}
`, config)
}

func (r StorageMoverAgentTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`

%s

resource "azurerm_storage_mover_agent" "test" {
  name             = "acctest-sa-%d"
  storage_mover_id = azurerm_storage_mover.test.id
  arc_resource_id  = "${azurerm_resource_group.test.id}/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"
  arc_vm_uuid      = "3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"
  description      = "Example Agent Description"
}
`, template, data.RandomInteger)
}

func (r StorageMoverAgentTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`

%s

resource "azurerm_storage_mover_agent" "test" {
  name             = "acctest-sa-%d"
  storage_mover_id = azurerm_storage_mover.test.id
  arc_resource_id  = "${azurerm_resource_group.test.id}/providers/Microsoft.HybridCompute/machines/examples-hybridComputeName"
  arc_vm_uuid      = "3bb2c024-eba9-4d18-9e7a-1d772fcc5fe9"
  description      = "Update Example Agent Description"

}
`, template, data.RandomInteger)
}
