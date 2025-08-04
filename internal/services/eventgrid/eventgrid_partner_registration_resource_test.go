package eventgrid_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/partnerregistrations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type EventGridPartnerRegistrationTestResource struct{}

func TestAccEventGridPartnerRegistration_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_registration", "test")
	r := EventGridPartnerRegistrationTestResource{}

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

func TestAccEventGridPartnerRegistration_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_registration", "test")
	r := EventGridPartnerRegistrationTestResource{}

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

func TestAccEventGridPartnerRegistration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_partner_registration", "test")
	r := EventGridPartnerRegistrationTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
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

func (EventGridPartnerRegistrationTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := partnerregistrations.ParsePartnerRegistrationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.EventGrid.PartnerRegistrations.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r EventGridPartnerRegistrationTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventgrid_partner_registration" "test" {
  name                = "acctestPartnerReg-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  tags                = { "environment" = "test" }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r EventGridPartnerRegistrationTestResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventgrid_partner_registration" "test" {
  name                = "acctestPartnerReg-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  tags                = { "environment" = "updated", "foo" = "bar" }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r EventGridPartnerRegistrationTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_eventgrid_partner_registration" "import" {
  name                = azurerm_eventgrid_partner_registration.test.name
  resource_group_name = azurerm_eventgrid_partner_registration.test.resource_group_name
  tags                = { "environment" = "test" }
}
`, r.basic(data))
}
