package communication_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/communication/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type CommunicationServiceResource struct{}

func TestAccCommunicationService_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_communication_service", "test")
	r := CommunicationServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCommunicationService_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_communication_service", "test")
	r := CommunicationServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCommunicationService_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_communication_service", "test")
	r := CommunicationServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCommunicationService_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_communication_service", "test")
	r := CommunicationServiceResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CommunicationServiceResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	clusterClient := client.Communication.ServiceClient
	id, err := parse.CommunicationServiceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clusterClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Communication Service %q: %+v", state.ID, err)
	}

	return utils.Bool(resp.ServiceProperties != nil), nil
}

func (r CommunicationServiceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_communication_service" "test" {
  name                = "acctest-CommunicationService-%d"
  resource_group_name = azurerm_resource_group.test.name
}
`, r.template(data), data.RandomInteger)
}

func (r CommunicationServiceResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_communication_service" "import" {
  name                = azurerm_communication_service.test.name
  resource_group_name = azurerm_communication_service.test.resource_group_name
  data_location       = azurerm_communication_service.test.data_location
}
`, config)
}

func (r CommunicationServiceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_communication_service" "test" {
  name                = "acctest-CommunicationService-%d"
  resource_group_name = azurerm_resource_group.test.name
  data_location       = "United States"

  tags = {
    env = "Test"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CommunicationServiceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_communication_service" "test" {
  name                = "acctest-CommunicationService-%d"
  resource_group_name = azurerm_resource_group.test.name
  data_location       = "Australia"

  tags = {
    env = "Test2"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r CommunicationServiceResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-communicationservice-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
