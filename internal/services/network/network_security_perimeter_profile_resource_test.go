package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkSecurityPerimeterProfileTestResource struct{}

func TestAccNetworkSecurityPerimeterProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_profile", "test")
	r := NetworkSecurityPerimeterProfileTestResource{}

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

func TestAccNetworkSecurityPerimeterProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_profile", "test")
	r := NetworkSecurityPerimeterProfileTestResource{}

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

func TestAccNetworkSecurityPerimeterProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_profile", "test")
	r := NetworkSecurityPerimeterProfileTestResource{}

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

func TestAccNetworkSecurityPerimeterProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_profile", "test")
	r := NetworkSecurityPerimeterProfileTestResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
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

func (NetworkSecurityPerimeterProfileTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networksecurityperimeterprofiles.ParseProfileID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.NetworkSecurityPerimeterProfilesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (NetworkSecurityPerimeterProfileTestResource) basic(data acceptance.TestData) string {
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

resource "azurerm_network_security_perimeter_profile" "test" {
	name = "acctestProfile-%d"
	perimeter_id = azurerm_network_security_perimeter.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r NetworkSecurityPerimeterProfileTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_security_perimeter_profile" "import" {
	name = azurerm_network_security_perimeter_profile.test.name
	perimeter_id = azurerm_network_security_perimeter_profile.test.perimeter_id
}
`, r.basic(data))
}

func (NetworkSecurityPerimeterProfileTestResource) complete(data acceptance.TestData) string {
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

resource "azurerm_network_security_perimeter_profile" "test" {
	name = "acctestProfile-%d"
	perimeter_id = azurerm_network_security_perimeter.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
