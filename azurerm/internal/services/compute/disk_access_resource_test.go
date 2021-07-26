package compute_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/compute/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DiskAccessResource struct {
}

func TestAccDiskAccess_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_access", "test")
	r := DiskAccessResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDiskAccess_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_access", "test")
	r := DiskAccessResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_disk_access"),
		},
	})
}

func TestAccDiskAccess_import(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_disk_access", "test")
	r := DiskAccessResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (t DiskAccessResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DiskAccessID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Compute.DiskAccessClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Compute Disk Access %q", id.String())
	}

	return utils.Bool(resp.ID != nil), nil
}

func (DiskAccessResource) empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_disk_access" "test" {
  name                = "acctestda-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}

	`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r DiskAccessResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
	%s

resource "azurerm_disk_access" "import" {
  name                = azurerm_disk_access.test.name
  location            = azurerm_disk_access.test.location
  resource_group_name = azurerm_disk_access.test.resource_group_name

  tags = {
    environment = "acctest"
    cost-center = "ops"
  }
}
`, r.empty(data))
}

func (DiskAccessResource) importConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_disk_access" "test" {
  name                = "accda%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location


  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
