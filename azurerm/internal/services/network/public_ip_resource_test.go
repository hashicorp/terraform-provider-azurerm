package network_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type PublicIPResource struct {
}

func TestAccPublicIpStatic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.static_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("allocation_method").HasValue("Static"),
				check.That(data.ResourceName).Key("ip_version").HasValue("IPv4"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.static_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("allocation_method").HasValue("Static"),
				check.That(data.ResourceName).Key("ip_version").HasValue("IPv4"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_public_ip"),
		},
	})
}

func TestAccPublicIpStatic_zones(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withZone(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("allocation_method").HasValue("Static"),
				check.That(data.ResourceName).Key("zones.#").HasValue("1"),
				check.That(data.ResourceName).Key("zones.0").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_basic_withDNSLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}
	dnl := fmt.Sprintf("acctestdnl-%d", data.RandomInteger)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic_withDNSLabel(data, dnl),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("allocation_method").HasValue("Static"),
				check.That(data.ResourceName).Key("domain_name_label").HasValue(dnl),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_standard_withIPv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standard_withIPVersion(data, "IPv6"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_version").HasValue("IPv6"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpDynamic_basic_withIPv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}
	ipVersion := "Ipv6"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dynamic_basic_withIPVersion(data, ipVersion),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_version").HasValue("IPv6"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_basic_defaultsToIPv4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.static_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_version").HasValue("IPv4"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_basic_withIPv4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}
	ipVersion := "IPv4"

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.static_basic_withIPVersion(data, ipVersion),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_version").HasValue("IPv4"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_standard(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standard(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.static_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				testCheckPublicIpDisappears(data.ResourceName),
			),
			ExpectNonEmptyPlan: true,
		},
	})
}

func TestAccPublicIpStatic_idleTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.idleTimeout(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("idle_timeout_in_minutes").HasValue("30"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("staging"),
			),
		},
	})
}

func TestAccPublicIpStatic_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.static_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("domain_name_label").HasValue(fmt.Sprintf("acctest-%d", data.RandomInteger)),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_standardPrefix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standardPrefix(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccPublicIpStatic_standardPrefixWithTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.standardPrefixWithTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		{
			Config: r.standardPrefixWithTagsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("staging"),
			),
		},
	})
}

func TestAccPublicIpDynamic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dynamic_basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_importIdError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.static_basic(data),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateId:     fmt.Sprintf("/subscriptions/%s/resourceGroups/acctestRG-%d/providers/Microsoft.Network/publicIPAdresses/acctestpublicip-%d", os.Getenv("ARM_SUBSCRIPTION_ID"), data.RandomInteger, data.RandomInteger),
			ExpectError:       regexp.MustCompile("Error parsing supplied resource id."),
		},
	})
}

func TestAccPublicIpStatic_canLabelBe63(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.canLabelBe63(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("allocation_method").HasValue("Static"),
			),
		},
		data.ImportStep(),
	})
}

func (t PublicIPResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resGroup := id.ResourceGroup
	name := id.Path["publicIPAddresses"]

	resp, err := clients.Network.PublicIPsClient.Get(ctx, resGroup, name, "")
	if err != nil {
		return nil, fmt.Errorf("reading Public IP (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func testCheckPublicIpDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Network.PublicIPsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		publicIpName := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for public ip: %s", publicIpName)
		}

		future, err := client.Delete(ctx, resourceGroup, publicIpName)
		if err != nil {
			return fmt.Errorf("Error deleting Public IP %q (Resource Group %q): %+v", publicIpName, resourceGroup, err)
		}

		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Error waiting for deletion of Public IP %q (Resource Group %q): %+v", publicIpName, resourceGroup, err)
		}

		return nil
	}
}

func (PublicIPResource) static_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r PublicIPResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip" "import" {
  name                = azurerm_public_ip.test.name
  location            = azurerm_public_ip.test.location
  resource_group_name = azurerm_public_ip.test.resource_group_name
  allocation_method   = azurerm_public_ip.test.allocation_method
}
`, r.static_basic(data))
}

func (PublicIPResource) withZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPResource) basic_withDNSLabel(data acceptance.TestData, dnsNameLabel string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, dnsNameLabel)
}

func (PublicIPResource) static_basic_withIPVersion(data acceptance.TestData, ipVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  ip_version          = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, ipVersion)
}

func (PublicIPResource) standard(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPResource) standardPrefix(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  public_ip_prefix_id = azurerm_public_ip_prefix.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PublicIPResource) standardPrefixWithTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  public_ip_prefix_id = azurerm_public_ip_prefix.test.id

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PublicIPResource) standardPrefixWithTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  public_ip_prefix_id = azurerm_public_ip_prefix.test.id

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PublicIPResource) standard_withIPVersion(data acceptance.TestData, ipVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  ip_version          = "%s"
  sku                 = "Standard"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, ipVersion)
}

func (PublicIPResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  domain_name_label   = "acctest-%d"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PublicIPResource) idleTimeout(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpublicip-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  idle_timeout_in_minutes = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPResource) dynamic_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPResource) dynamic_basic_withIPVersion(data acceptance.TestData, ipVersion string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Dynamic"

  ip_version = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, ipVersion)
}

func (PublicIPResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPResource) canLabelBe63(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  allocation_method = "Static"
  domain_name_label = "k2345678-1-2345678-2-2345678-3-2345678-4-2345678-5-2345678-6-23"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
