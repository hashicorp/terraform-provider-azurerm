package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LocalNetworkGatewayResource struct {
}

func TestAccLocalNetworkGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLocalNetworkGateway_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_local_network_gateway"),
		},
	})
}

func TestAccLocalNetworkGateway_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				testCheckLocalNetworkGatewayDisappears(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccLocalNetworkGateway_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("acctest"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLocalNetworkGateway_bgpSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.bgpSettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLocalNetworkGateway_bgpSettingsDisable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.bgpSettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.asn").HasValue("2468"),
				check.That(data.ResourceName).Key("bgp_settings.0.bgp_peering_address").HasValue("10.104.1.1"),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("0"),
			),
		},
	})
}

func TestAccLocalNetworkGateway_bgpSettingsEnable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("0"),
			),
		},
		{
			Config: r.bgpSettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.asn").HasValue("2468"),
				check.That(data.ResourceName).Key("bgp_settings.0.bgp_peering_address").HasValue("10.104.1.1"),
			),
		},
	})
}

func TestAccLocalNetworkGateway_bgpSettingsComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.bgpSettingsComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.0.0/8"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.asn").HasValue("2468"),
				check.That(data.ResourceName).Key("bgp_settings.0.bgp_peering_address").HasValue("10.104.1.1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peer_weight").HasValue("15"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLocalNetworkGateway_updateAddressSpace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.multipleAddressSpace(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multipleAddressSpaceUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLocalNetworkGateway_fqdn(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.fqdn(data, "www.foo.com"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.fqdn(data, "www.bar.com"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t LocalNetworkGatewayResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	name := id.Path["localNetworkGateways"]
	resGroup := id.ResourceGroup

	resp, err := clients.Network.LocalNetworkGatewaysClient.Get(ctx, resGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading Local Network Gateway (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func testCheckLocalNetworkGatewayDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// first, check that it exists on Azure:
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.LocalNetworkGatewaysClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// then check within the schema for the local network gateway:
		res, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Local network gateway '%s' not found.", resourceName)
		}

		// and finally, extract the name and the resource group:
		id, err := azure.ParseAzureResourceID(res.Primary.ID)
		if err != nil {
			return err
		}
		localNetName := id.Path["localNetworkGateways"]
		resourceGroup := id.ResourceGroup

		future, err := client.Delete(ctx, resourceGroup, localNetName)
		if err != nil {
			if response.WasNotFound(future.Response()) {
				return fmt.Errorf("Local network gateway %q (resource group %q) does not exist on Azure.", localNetName, resourceGroup)
			}
			return fmt.Errorf("Error deleting the state of local network gateway %q: %+v", localNetName, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of the local network gateway %q to complete: %+v", localNetName, err)
		}

		return nil
	}
}

func (LocalNetworkGatewayResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LocalNetworkGatewayResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_local_network_gateway" "import" {
  name                = azurerm_local_network_gateway.test.name
  location            = azurerm_local_network_gateway.test.location
  resource_group_name = azurerm_local_network_gateway.test.resource_group_name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]
}
`, r.basic(data))
}

func (LocalNetworkGatewayResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  tags = {
    environment = "acctest"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LocalNetworkGatewayResource) bgpSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  bgp_settings {
    asn                 = 2468
    bgp_peering_address = "10.104.1.1"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LocalNetworkGatewayResource) bgpSettingsComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/8"]

  bgp_settings {
    asn                 = 2468
    bgp_peering_address = "10.104.1.1"
    peer_weight         = 15
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LocalNetworkGatewayResource) multipleAddressSpace(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.0.0/24", "127.0.1.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LocalNetworkGatewayResource) multipleAddressSpaceUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.1.0/24", "127.0.0.0/24"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LocalNetworkGatewayResource) fqdn(data acceptance.TestData, fqdn string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_fqdn        = %q
  address_space       = ["127.0.0.0/8"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, fqdn)
}
