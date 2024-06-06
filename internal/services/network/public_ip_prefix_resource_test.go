// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PublicIPPrefixResource struct{}

func (t PublicIPPrefixResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := publicipprefixes.ParsePublicIPPrefixID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.PublicIPPrefixes.Get(ctx, *id, publicipprefixes.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (PublicIPPrefixResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := publicipprefixes.ParsePublicIPPrefixID(state.ID)
	if err != nil {
		return nil, err
	}

	if err := client.Network.PublicIPPrefixes.DeleteThenPoll(ctx, *id); err != nil {
		return nil, fmt.Errorf("deleting %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func TestAccPublicIpPrefix_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")
	r := PublicIPPrefixResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_prefix").Exists(),
				check.That(data.ResourceName).Key("prefix_length").HasValue("28"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpPrefix_ipv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")
	r := PublicIPPrefixResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv6(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_prefix").Exists(),
				check.That(data.ResourceName).Key("prefix_length").HasValue("126"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpPrefix_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")
	r := PublicIPPrefixResource{}

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

func TestAccPublicIpPrefix_prefixLength31(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")
	r := PublicIPPrefixResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.prefixLength31(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_prefix").Exists(),
				check.That(data.ResourceName).Key("prefix_length").HasValue("31"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpPrefix_prefixLength24(t *testing.T) {
	// NOTE: This test will fail unless the subscription is updated
	//        to accept a minimum PrefixLength of 24
	// more detail about [public ip limits](https://learn.microsoft.com/en-us/azure/virtual-network/ip-services/public-ip-addresses#limits)
	// you can submit a support request to increase the limit in [Azure Portal](https://learn.microsoft.com/en-us/azure/networking/check-usage-against-limits#azure-portal)
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")
	r := PublicIPPrefixResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.prefixLength24(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("ip_prefix").Exists(),
				check.That(data.ResourceName).Key("prefix_length").HasValue("24"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpPrefix_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")
	r := PublicIPPrefixResource{}

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

func TestAccPublicIpPrefix_disappears(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")
	r := PublicIPPrefixResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		data.DisappearsStep(acceptance.DisappearsStepData{
			Config:       r.basic,
			TestResource: r,
		}),
	})
}

func TestAccPublicIpPrefix_zonesSingle(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")
	r := PublicIPPrefixResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zonesSingle(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPublicIpPrefix_zonesMultiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_public_ip_prefix", "test")
	r := PublicIPPrefixResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.zonesMultiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (PublicIPPrefixResource) basic(data acceptance.TestData) string {
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
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPPrefixResource) ipv6(data acceptance.TestData) string {
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
  ip_version          = "IPv6"

  prefix_length = 126
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r PublicIPPrefixResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_public_ip_prefix" "import" {
  name                = azurerm_public_ip_prefix.test.name
  location            = azurerm_public_ip_prefix.test.location
  resource_group_name = azurerm_public_ip_prefix.test.resource_group_name
}
`, r.basic(data))
}

func (PublicIPPrefixResource) withTags(data acceptance.TestData) string {
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

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPPrefixResource) withTagsUpdate(data acceptance.TestData) string {
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

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPPrefixResource) prefixLength31(data acceptance.TestData) string {
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

  prefix_length = 31
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPPrefixResource) prefixLength24(data acceptance.TestData) string {
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

  prefix_length = 24
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PublicIPPrefixResource) zonesSingle(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  zones               = ["1"]
}
`, data.RandomInteger, data.Locations.Primary)
}

func (PublicIPPrefixResource) zonesMultiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "acctestpublicipprefix-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  zones               = ["1", "2", "3"]
}
`, data.RandomInteger, data.Locations.Primary)
}
