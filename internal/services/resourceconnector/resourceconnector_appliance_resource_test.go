package resourceconnector_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resourceconnector/2022-10-27/appliances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ResourceConnectorApplianceResource struct{}

func TestAccResourceConnectorAppliance_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_connector_appliance", "test")
	l := ResourceConnectorApplianceResource{}

	data.ResourceTest(t, l, []acceptance.TestStep{
		{
			Config: l.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(l),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceConnectorAppliance_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_connector_appliance", "test")
	l := ResourceConnectorApplianceResource{}

	data.ResourceTest(t, l, []acceptance.TestStep{
		{
			Config: l.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(l),
			),
		},
		data.ImportStep(),
		{
			Config: l.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(l),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceConnectorAppliance_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_connector_appliance", "test")
	l := ResourceConnectorApplianceResource{}

	data.ResourceTest(t, l, []acceptance.TestStep{
		{
			Config: l.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(l),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceConnectorAppliance_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_connector_appliance", "test")
	l := ResourceConnectorApplianceResource{}
	data.ResourceTest(t, l, []acceptance.TestStep{
		{
			Config: l.basic(data),
		},
		{
			Config:      l.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_resource_connector_appliance"),
		},
	})
}

func (r ResourceConnectorApplianceResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := appliances.ParseApplianceID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.ResourceConnector.AppliancesClient
	resp, err := client.Get(ctx, *id)

	if err != nil {
		return nil, fmt.Errorf("retrieving %q: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (r ResourceConnectorApplianceResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-Appliances-%[1]d"
  location = "%[2]s"
}
resource "azurerm_resource_connector_appliance" "test" {
  name                    = "acctestAppliances-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  distro                  = "AKSEdge"
  infrastructure_provider = "VMWare"
  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceConnectorApplianceResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-Appliances-%[1]d"
  location = "%[2]s"
}
resource "azurerm_resource_connector_appliance" "test" {
  name                    = "acctestAppliances-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  distro                  = "AKSEdge"
  infrastructure_provider = "VMWare"
  public_key              = "MIGJAoGBALXhHAjXWcYsF5oMPrYPfYwA/Jim7ErxlRM7laOhvUqFuMkEGxOGf76W4NhMoouFty7SUeio+IWgHjwUmiXDBhVsNie2Pe5XSyuSmvhRIFOoULfKUgv3qEBIHq0ylZOoaNIFN/HFALRIqejEh2MF5URi3fBxJA4tDV4tgR+KdYJ9AgMBAAE="
  identity {
    type = "SystemAssigned"
  }
  tags = {
    "hello" = "world"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceConnectorApplianceResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-appliances-%[1]d"
  location = "%[2]s"
}
resource "azurerm_resource_connector_appliance" "test" {
  name                    = "acctestappliances-%[1]d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  distro                  = "AKSEdge"
  infrastructure_provider = "VMWare"
  identity {
    type = "SystemAssigned"
  }
  tags = {
    "hello" = "world"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ResourceConnectorApplianceResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s
resource "azurerm_resource_connector_appliance" "import" {
  name                    = azurerm_resource_connector_appliance.test.name
  location                = azurerm_resource_connector_appliance.test.location
  resource_group_name     = azurerm_resource_connector_appliance.test.resource_group_name
  distro                  = azurerm_resource_connector_appliance.test.distro
  infrastructure_provider = azurerm_resource_connector_appliance.test.infrastructure_provider
  identity {
    type = "SystemAssigned"
  }
}
`, r.basic(data), data.RandomInteger)
}
