// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type PublicIPResource struct{}

func TestAccPublicIpStatic_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.static_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("allocation_method").HasValue("Static"),
				check.That(data.ResourceName).Key("ip_version").HasValue("IPv4"),
				check.That(data.ResourceName).Key("ddos_protection_mode").HasValue("VirtualNetworkInherited"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.static_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func TestAccPublicIp_zonesSingle(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zonesSingle(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("allocation_method").HasValue("Static"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIp_zonesMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zonesMultiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("allocation_method").HasValue("Static"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_basic_withDNSLabel(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}
	dnl := fmt.Sprintf("acctestdnl-%d", data.RandomInteger)

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic_withDNSLabel(data, dnl),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standard_withIPVersion(data, "IPv6"),
			Check: acceptance.ComposeTestCheckFunc(
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
	ipVersion := "IPv6"

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dynamic_basic_withIPVersion(data, ipVersion),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.static_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.static_basic_withIPVersion(data, ipVersion),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standard(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_standard_withDDoS(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardDDoSDisabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ddos_protection_mode").HasValue("Disabled"),
			),
		},
		data.ImportStep(),
		{
			Config: r.standardDDoSEnabled(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ddos_protection_mode").HasValue("Enabled"),
				check.That(data.ResourceName).Key("ddos_protection_plan_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.static_basic,
			TestResource: r,
		}),
	})
}

func TestAccPublicIpStatic_idleTimeout(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.idleTimeout(data, 30),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.idleTimeout(data, 15),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.static_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardPrefix(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccPublicIpStatic_standardPrefixWithTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standardPrefixWithTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("Production"),
				check.That(data.ResourceName).Key("tags.cost_center").HasValue("MSFT"),
			),
		},
		{
			Config: r.standardPrefixWithTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.dynamic_basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_canLabelBe63(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.canLabelBe63(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("allocation_method").HasValue("Static"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_ipTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.standard_IpTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_tags.RoutingPreference").HasValue("Internet"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_globalTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.globalTier(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("sku").HasValue("Standard"),
				check.That(data.ResourceName).Key("sku_tier").HasValue("Global"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_regionalTier(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.regionalTier(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_address").Exists(),
				check.That(data.ResourceName).Key("sku").HasValue("Basic"),
				check.That(data.ResourceName).Key("sku_tier").HasValue("Regional"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpStatic_edgeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip", "test")
	r := PublicIPResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.edgeZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t PublicIPResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParsePublicIPAddressID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.PublicIPAddresses.Get(ctx, *id, publicipaddresses.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (PublicIPResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParsePublicIPAddressID(state.ID)
	if err != nil {
		return nil, err
	}

	ctx2, cancel := context.WithTimeout(ctx, 30*time.Minute)
	defer cancel()
	if err := client.Network.PublicIPAddresses.DeleteThenPoll(ctx2, *id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return pointer.To(true), nil
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
  sku                 = "Basic"
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
  sku                 = "Basic"
}
`, r.static_basic(data))
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
  sku                 = "Basic"
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
  sku                 = "Basic"
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

func (PublicIPResource) standardDDoSDisabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                 = "acctestpublicip-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  allocation_method    = "Static"
  sku                  = "Standard"
  ddos_protection_mode = "Disabled"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPResource) standardDDoSEnabled(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_ddos_protection_plan" "test" {
  name                = "acctestddospplan-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_public_ip" "test" {
  name                    = "acctestpublicip-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  sku                     = "Standard"
  ddos_protection_mode    = "Enabled"
  ddos_protection_plan_id = azurerm_network_ddos_protection_plan.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
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
  name                    = "acctestpublicip-%d"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  allocation_method       = "Static"
  sku                     = "Basic"
  domain_name_label       = "acctest-%d"
  idle_timeout_in_minutes = 30
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PublicIPResource) idleTimeout(data acceptance.TestData, idleTimeout int) string {
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
  sku                     = "Basic"
  idle_timeout_in_minutes = %d
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, idleTimeout)
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
  sku                 = "Basic"
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
  sku                 = "Basic"

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
  sku                 = "Basic"

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
  sku                 = "Basic"

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
  sku               = "Basic"
  domain_name_label = "%s"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, strings.ToLower(data.RandomStringOfLength(63)))
}

func (PublicIPResource) standard_IpTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1", "2", "3"]

  ip_tags = {
    RoutingPreference = "Internet"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPResource) globalTier(data acceptance.TestData) string {
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
  sku_tier            = "Global"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPResource) regionalTier(data acceptance.TestData) string {
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
  sku                 = "Basic"
  sku_tier            = "Regional"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPResource) zonesSingle(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (PublicIPResource) zonesMultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  zones               = ["1", "2", "3"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (PublicIPResource) edgeZone(data acceptance.TestData) string {
	// @tombuildsstuff: WestUS has an edge zone available - so hard-code to that for now
	data.Locations.Primary = "westus"

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

data "azurerm_extended_locations" "test" {
  location = azurerm_resource_group.test.location
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpublicip-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  edge_zone           = data.azurerm_extended_locations.test.extended_locations[0]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
