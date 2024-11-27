package fabric_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/fabric/2023-11-01/fabriccapacities"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type FabricFabricCapacityResource struct{}

func TestAccFabricFabricCapacity_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fabric_capacity", "test")
	r := FabricFabricCapacityResource{}
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

func TestAccFabricFabricCapacity_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fabric_capacity", "test")
	r := FabricFabricCapacityResource{}
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

func TestAccFabricFabricCapacity_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fabric_capacity", "test")
	r := FabricFabricCapacityResource{}
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

func TestAccFabricFabricCapacity_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_fabric_capacity", "test")
	r := FabricFabricCapacityResource{}
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
	})
}

func (r FabricFabricCapacityResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := fabriccapacities.ParseCapacityID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Fabric.FabricCapacitiesClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r FabricFabricCapacityResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r FabricFabricCapacityResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_fabric_capacity" "test" {
  name                   = "acctestffc%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = "%s"
  administration_members = [data.azurerm_client_config.current.object_id]

  sku {
    name = "F32"
    tier = "Fabric"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r FabricFabricCapacityResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_fabric_capacity" "import" {
  name                   = azurerm_fabric_capacity.test.name
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  administration_members = azurerm_fabric_capacity.test.administration_members

  sku {
    name = azurerm_fabric_capacity.test.sku[0].name
    tier = azurerm_fabric_capacity.test.sku[0].tier
  }
}
`, r.basic(data))
}

func (r FabricFabricCapacityResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_fabric_capacity" "test" {
  name                = "acctestffc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"

  administration_members = [data.azurerm_client_config.current.object_id]

  sku {
    name = "F32"
    tier = "Fabric"
  }

  tags = {
    environment = "test"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r FabricFabricCapacityResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_fabric_capacity" "test" {
  name                = "acctestffc%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"

  administration_members = []

  sku {
    name = "F64"
    tier = "Fabric"
  }

  tags = {
    environment = "test1"
    environment = "test2"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}
