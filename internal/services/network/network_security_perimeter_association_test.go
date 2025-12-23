package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networksecurityperimeterassociations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkSecurityPerimeterAssociationTestResource struct{}

func TestAccNetworkSecurityPerimeterAssociation_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_association", "test")
	r := NetworkSecurityPerimeterAssociationTestResource{}

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

func TestAccNetworkSecurityPerimeterAssociation_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_association", "test")
	r := NetworkSecurityPerimeterAssociationTestResource{}

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

func TestAccNetworkSecurityPerimeterAssociation_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_association", "test")
	r := NetworkSecurityPerimeterAssociationTestResource{}

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

func TestAccNetworkSecurityPerimeterAssociation_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_network_security_perimeter_association", "test")
	r := NetworkSecurityPerimeterAssociationTestResource{}

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

func (NetworkSecurityPerimeterAssociationTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := networksecurityperimeterassociations.ParseResourceAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.NetworkSecurityPerimeterAssociationsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (NetworkSecurityPerimeterAssociationTestResource) basic(data acceptance.TestData) string {
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

resource "azurerm_log_analytics_workspace" "test" {
  name                = "example"
  location            = azurerm_network_security_perimeter.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_perimeter_association" "test" {
	name = "acctestassoc-%d"
	profile_id = azurerm_network_security_perimeter_profile.test.id 
	resource_id = azurerm_log_analytics_workspace.test.id

	access_mode = "Learning"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r NetworkSecurityPerimeterAssociationTestResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_network_security_perimeter_association" "import" {
	name = azurerm_network_security_perimeter_association.test.name
	profile_id = azurerm_network_security_perimeter_association.test.profile_id
	resource_id = azurerm_network_security_perimeter_association.test.resource_id

	access_mode = azurerm_network_security_perimeter_association.test.access_mode
}
`, r.basic(data))
}

func (NetworkSecurityPerimeterAssociationTestResource) complete(data acceptance.TestData) string {
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

resource "azurerm_log_analytics_workspace" "test" {
  name                = "example"
  location            = azurerm_network_security_perimeter.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_network_security_perimeter_association" "test" {
	name = "acctestassoc-%d"
	profile_id = azurerm_network_security_perimeter_profile.test.id 
	resource_id = azurerm_log_analytics_workspace.test.id

	access_mode = "Learning"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
