// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/customipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type CustomIpPrefixResource struct{}

const (
	ipv4TestCidr = "194.41.20.0/24"
	ipv6TestCidr = "2620:10c:5001::/48"
)

func TestAccCustomIpPrefixV4(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"ipv4": {
			"commissioned":   testAccCustomIpPrefix_ipv4,
			"requiresImport": testAccCustomIpPrefix_ipv4RequiresImport,
		},
	})
}

func TestAccCustomIpPrefixV6(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"ipv6": {
			"commissioned": testAccCustomIpPrefix_ipv6,
		},
	})
}

func testAccCustomIpPrefix_ipv4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4Provisioned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4CommissionedUnadvertised(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4Commissioned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv4Provisioned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_ipv4RequiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4Provisioned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.ipv4RequiresImport),
	})
}

func testAccCustomIpPrefix_ipv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "global")
	data2 := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "regional")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv6(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv6CommissionedGlobalUnadvertised(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.ipv6Commissioned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data2.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (CustomIpPrefixResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := customipprefixes.ParseCustomIPPrefixID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.Client.CustomIPPrefixes.Get(ctx, *id, customipprefixes.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return pointer.To(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(true), nil
}

func (r CustomIpPrefixResource) ipv4Provisioned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidr  = "%[3]s"
  zones = ["1"]

  roa_validity_end_date         = "2099-12-12"
  wan_validation_signed_message = "signed message for WAN validation"
}
`, data.RandomInteger, data.Locations.Primary, ipv4TestCidr)
}

func (r CustomIpPrefixResource) ipv4Commissioned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidr  = "%[3]s"
  zones = ["1"]

  roa_validity_end_date         = "2099-12-12"
  wan_validation_signed_message = "signed message for WAN validation"

  commissioning_enabled         = true
  internet_advertising_disabled = false
}
`, data.RandomInteger, data.Locations.Primary, ipv4TestCidr)
}

func (r CustomIpPrefixResource) ipv4CommissionedUnadvertised(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "test" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidr  = "%[3]s"
  zones = ["1"]

  roa_validity_end_date         = "2099-12-12"
  wan_validation_signed_message = "signed message for WAN validation"

  commissioning_enabled         = true
  internet_advertising_disabled = true
}
`, data.RandomInteger, data.Locations.Primary, ipv4TestCidr)
}

func (r CustomIpPrefixResource) ipv4RequiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_custom_ip_prefix" "import" {
  name                = azurerm_custom_ip_prefix.test.name
  location            = azurerm_custom_ip_prefix.test.location
  resource_group_name = azurerm_custom_ip_prefix.test.resource_group_name

  cidr  = azurerm_custom_ip_prefix.test.cidr
  zones = azurerm_custom_ip_prefix.test.zones

  roa_validity_end_date         = azurerm_custom_ip_prefix.test.roa_validity_end_date
  wan_validation_signed_message = azurerm_custom_ip_prefix.test.wan_validation_signed_message
}
`, r.ipv4Provisioned(data))
}

func (r CustomIpPrefixResource) ipv6(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "global" {
  name                = "acctest-v6global-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidr = "%[3]s"

  roa_validity_end_date         = "2199-12-12"
  wan_validation_signed_message = "signed message for WAN validation"
}

resource "azurerm_custom_ip_prefix" "regional" {
  name                       = "acctest-v6regional-%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  parent_custom_ip_prefix_id = azurerm_custom_ip_prefix.global.id

  cidr  = cidrsubnet(azurerm_custom_ip_prefix.global.cidr, 16, 1)
  zones = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, ipv6TestCidr)
}

func (r CustomIpPrefixResource) ipv6CommissionedGlobalUnadvertised(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "global" {
  name                = "acctest-v6global-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidr = "%[3]s"

  roa_validity_end_date         = "2199-12-12"
  wan_validation_signed_message = "signed message for WAN validation"

  commissioning_enabled = false
}

resource "azurerm_custom_ip_prefix" "regional" {
  name                       = "acctest-v6regional-%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  parent_custom_ip_prefix_id = azurerm_custom_ip_prefix.global.id

  cidr  = cidrsubnet(azurerm_custom_ip_prefix.global.cidr, 16, 1)
  zones = ["1"]

  commissioning_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, ipv6TestCidr)
}

func (r CustomIpPrefixResource) ipv6Commissioned(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "global" {
  name                = "acctest-v6global-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  cidr = "%[3]s"

  roa_validity_end_date         = "2199-12-12"
  wan_validation_signed_message = "signed message for WAN validation"

  commissioning_enabled = true
}

resource "azurerm_custom_ip_prefix" "regional" {
  name                       = "acctest-v6regional-%[1]d"
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  parent_custom_ip_prefix_id = azurerm_custom_ip_prefix.global.id

  cidr  = cidrsubnet(azurerm_custom_ip_prefix.global.cidr, 16, 1)
  zones = ["1"]

  commissioning_enabled = true
}
`, data.RandomInteger, data.Locations.Primary, ipv6TestCidr)
}
