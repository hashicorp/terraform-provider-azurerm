package labservice_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/lab"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LabServiceLabResource struct{}

func TestAccLabServicesLab_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_lab", "test")
	r := LabServiceLabResource{}

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

func TestAccLabServiceLab_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_lab", "test")
	r := LabServiceLabResource{}

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

func TestAccLabServiceLab_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_lab", "test")
	r := LabServiceLabResource{}

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

func TestAccLabServiceLab_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_lab_service_lab", "test")
	r := LabServiceLabResource{}

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
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r LabServiceLabResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := lab.ParseLabID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.LabService.LabClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r LabServiceLabResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lab-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LabServiceLabResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_lab" "test" {
  name                = "acctest-lab-%d"
  resource_group_name = azurerm_resource_group.lab.name
  location            = azurerm_resource_group.lab.location
  description         = "Testing description"
  title               = "Testing title"
}
`, r.template(data), data.RandomInteger)
}

func (r LabServiceLabResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_lab" "import" {
  name                = azurerm_lab_service_lab.test.name
  resource_group_name = azurerm_lab_service_lab.test.resource_group_name
  location            = azurerm_lab_service_lab.test.location
  description         = azurerm_lab_service_lab.test.description
  title               = azurerm_lab_service_lab.test.title
}
`, r.basic(data))
}

func (r LabServiceLabResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_lab" "test" {
  name                = "acctest-lab-%d"
  resource_group_name = azurerm_resource_group.lab.name
  location            = azurerm_resource_group.lab.location
  description         = "Testing description"
  title               = "Testing title"
}
`, r.template(data), data.RandomInteger)
}

func (r LabServiceLabResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lab_service_lab" "test" {
  name                = "acctest-lab-%d"
  resource_group_name = azurerm_resource_group.lab.name
  location            = azurerm_resource_group.lab.location
  description         = "Testing description"
  title               = "Testing title"
}
`, r.template(data), data.RandomInteger)
}
