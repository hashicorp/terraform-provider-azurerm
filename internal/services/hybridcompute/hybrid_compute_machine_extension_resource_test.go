package hybridcompute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-03-10/machineextensions"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type HybridComputeMachineExtensionResource struct{}

func TestAccHybridComputeMachineExtension_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine_extension", "test")
	r := HybridComputeMachineExtensionResource{}
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

func TestAccHybridComputeMachineExtension_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine_extension", "test")
	r := HybridComputeMachineExtensionResource{}
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

func TestAccHybridComputeMachineExtension_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine_extension", "test")
	r := HybridComputeMachineExtensionResource{}
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

func TestAccHybridComputeMachineExtension_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_hybrid_compute_machine_extension", "test")
	r := HybridComputeMachineExtensionResource{}
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

func (r HybridComputeMachineExtensionResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := machineextensions.ParseExtensionID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.HybridCompute.MachineExtensionsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r HybridComputeMachineExtensionResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_hybrid_compute_machine" "test" {
  name                = "acctest-hcm-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r HybridComputeMachineExtensionResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_hybrid_compute_machine_extension" "test" {
  name                      = "acctest-hcme-%d"
  hybrid_compute_machine_id = azurerm_hybrid_compute_machine.test.id
  location                  = "%s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r HybridComputeMachineExtensionResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_hybrid_compute_machine_extension" "import" {
  name                      = azurerm_hybrid_compute_machine_extension.test.name
  hybrid_compute_machine_id = azurerm_hybrid_compute_machine.test.id
  location                  = "%s"
}
`, config, data.Locations.Primary)
}

func (r HybridComputeMachineExtensionResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_hybrid_compute_machine_extension" "test" {
  name                       = "acctest-hcme-%d"
  hybrid_compute_machine_id  = azurerm_hybrid_compute_machine.test.id
  location                   = "%s"
  auto_upgrade_minor_version_enabled = false
  automatic_upgrade_enabled   = false
  force_update_tag           = ""
  publisher                  = ""
  type_handler_version       = ""
  instance_view {
    name                 = ""
    type                 = ""
    type_handler_version = ""
    status {
      code           = ""
      display_status = ""
      level          = ""
      message        = ""
      time           = ""
    }
  }
  protected_settings = jsonencode({
    "key" : "value"
  })
  settings = jsonencode({
    "key" : "value"
  })
  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r HybridComputeMachineExtensionResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_hybrid_compute_machine_extension" "test" {
  name                       = "acctest-hcme-%d"
  hybrid_compute_machine_id  = azurerm_hybrid_compute_machine.test.id
  location                   = "%s"
  auto_upgrade_minor_version_enabled = false
  automatic_upgrade_enabled   = false
  force_update_tag           = ""
  publisher                  = ""
  type_handler_version       = ""
  instance_view {
    name                 = ""
    type                 = ""
    type_handler_version = ""
    status {
      code           = ""
      display_status = ""
      level          = ""
      message        = ""
      time           = ""
    }
  }
  protected_settings = jsonencode({
    "key" : "value"
  })
  settings = jsonencode({
    "key" : "value"
  })
  tags = {
    key = "value"
  }

}
`, template, data.RandomInteger, data.Locations.Primary)
}
