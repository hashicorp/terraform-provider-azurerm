package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkSecurityPerimeterTestResource struct{}

func TestAccNetworkSecurityPerimeter_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter", "test")
	r := NetworkSecurityPerimeterTestResource{}

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

func TestAccNetworkSecurityPerimeter_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter", "test")
	r := NetworkSecurityPerimeterTestResource{}

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

func TestAccNetworkSecurityPerimeter_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter", "test")
	r := NetworkSecurityPerimeterTestResource{}

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

func TestAccNetworkSecurityPerimeter_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter", "test")
	r := NetworkSecurityPerimeterTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.env").HasValue("test"),
			),
		},
		data.ImportStep(),
	})
}

func (NetworkSecurityPerimeterTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networksecurityperimeters.ParseNetworkSecurityPerimeterID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.NetworkSecurityPerimetersClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (NetworkSecurityPerimeterTestResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_perimeter" "test" {
  name     = "acctestNsp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}

func (r NetworkSecurityPerimeterTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_security_perimeter" "import" {
  name     = azurerm_network_security_perimeter.test.name
  location = azurerm_network_security_perimeter.test.location
  resource_group_name = azurerm_network_security_perimeter.test.resource_group_name
}
`, r.basic(data))
}

func (NetworkSecurityPerimeterTestResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_perimeter" "test" {
  name     = "acctestNsp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = "%s"

  tags = {
    env = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary)
}
