package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type VirtualHubSecurityPartnerProviderResource struct {
}

func TestAccVirtualHubSecurityPartnerProvider_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_security_partner_provider", "test")
	r := VirtualHubSecurityPartnerProviderResource{}
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

func TestAccVirtualHubSecurityPartnerProvider_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_security_partner_provider", "test")
	r := VirtualHubSecurityPartnerProviderResource{}
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

func TestAccVirtualHubSecurityPartnerProvider_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_security_partner_provider", "test")
	r := VirtualHubSecurityPartnerProviderResource{}
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

func TestAccVirtualHubSecurityPartnerProvider_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_virtual_hub_security_partner_provider", "test")
	r := VirtualHubSecurityPartnerProviderResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t VirtualHubSecurityPartnerProviderResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.SecurityPartnerProviderID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.SecurityPartnerProviderClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("reading Security Partner Provider (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (VirtualHubSecurityPartnerProviderResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-vhub-%d"
  location = "%s"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest-vwan-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest-VHUB-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.2.0/24"
}

resource "azurerm_vpn_gateway" "test" {
  name                = "acctest-VPNG-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  virtual_hub_id      = azurerm_virtual_hub.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (r VirtualHubSecurityPartnerProviderResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_security_partner_provider" "test" {
  name                   = "acctest-SPP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  security_provider_name = "ZScaler"

  depends_on = [azurerm_vpn_gateway.test]
}
`, r.template(data), data.RandomInteger)
}

func (r VirtualHubSecurityPartnerProviderResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_security_partner_provider" "import" {
  name                   = azurerm_virtual_hub_security_partner_provider.test.name
  resource_group_name    = azurerm_virtual_hub_security_partner_provider.test.resource_group_name
  location               = azurerm_virtual_hub_security_partner_provider.test.location
  security_provider_name = azurerm_virtual_hub_security_partner_provider.test.security_provider_name
}
`, r.basic(data))
}

func (r VirtualHubSecurityPartnerProviderResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_hub_security_partner_provider" "test" {
  name                   = "acctest-SPP-%d"
  resource_group_name    = azurerm_resource_group.test.name
  location               = azurerm_resource_group.test.location
  virtual_hub_id         = azurerm_virtual_hub.test.id
  security_provider_name = "ZScaler"

  tags = {
    ENv = "Test"
  }

  depends_on = [azurerm_vpn_gateway.test]
}
`, r.template(data), data.RandomInteger)
}
