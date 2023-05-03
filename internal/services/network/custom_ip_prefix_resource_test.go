package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CustomIpPrefixResource struct{}

const (
	ipv4TestCidr = "194.41.20.0/24"
	ipv6TestCidr = "2620:10c:5001::/48"
)

func TestAccCustomIpPrefix(t *testing.T) {
	acceptance.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		//"ipv4": {
		//	"basic":                testAccCustomIpPrefix_ipv4,
		//	"commissioned":         testAccCustomIpPrefix_ipv4Commissioned,
		//	"commissionedRegional": testAccCustomIpPrefix_ipv4CommissionedRegional,
		//	"update":               testAccCustomIpPrefix_ipv4Update,
		//	"requiresImport":       testAccCustomIpPrefix_requiresImport,
		//},
		"ipv6": {
			"basic": testAccCustomIpPrefix_ipv6,
			//"commissioned":         testAccCustomIpPrefix_ipv6Commissioned,
			//"commissionedRegional": testAccCustomIpPrefix_ipv6CommissionedRegional,
			//"update":               testAccCustomIpPrefix_ipv6Update,
		},
	})
}

func testAccCustomIpPrefix_ipv4(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_ipv4Commissioned(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4Commissioned(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_ipv4CommissionedRegional(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4CommissionedRegional(data),
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
	})
}

func testAccCustomIpPrefix_ipv4Update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
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
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_ipv6(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv6(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func testAccCustomIpPrefix_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_custom_ip_prefix", "test")
	r := CustomIpPrefixResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.ipv4(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (CustomIpPrefixResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomIpPrefixID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Network.CustomIPPrefixesClient.Get(ctx, id.ResourceGroup, id.CustomIpPrefixeName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r CustomIpPrefixResource) ipv4(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_custom_ip_prefix" "test" {
  name                          = "acctest-%[1]d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  cidr                          = "%[3]s"
  roa_validity_end_date         = "2099-12-12"
  wan_validation_signed_message = "signed message for WAN validation"
  zones                         = ["1"]
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

  commissioning_enabled         = true
  cidr                          = "%[3]s"
  internet_advertising_disabled = false
  roa_validity_end_date         = "2099-12-12"
  wan_validation_signed_message = "signed message for WAN validation"
  zones                         = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, ipv4TestCidr)
}

func (r CustomIpPrefixResource) ipv4CommissionedRegional(data acceptance.TestData) string {
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

  commissioning_enabled         = true
  cidr                          = "%[3]s"
  internet_advertising_disabled = true
  roa_validity_end_date         = "2099-12-12"
  wan_validation_signed_message = "signed message for WAN validation"
  zones                         = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, ipv4TestCidr)
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
  name                          = "acctest-%[1]d"
  location                      = azurerm_resource_group.test.location
  resource_group_name           = azurerm_resource_group.test.name
  cidr                          = "%[3]s"
  roa_validity_end_date         = "2199-12-12"
  wan_validation_signed_message = "signed message for WAN validation"
}

resource "azurerm_custom_ip_prefix" "regional" {
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  cidr                = cidrsubnet(azurerm_custom_ip_prefix.global.cidr, 16, 1)
  zones               = ["1"]
}
`, data.RandomInteger, data.Locations.Primary, ipv6TestCidr)
}

func (r CustomIpPrefixResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_custom_ip_prefix" "import" {
  name                = azurerm_custom_ip_prefix.test.name
  location            = azurerm_custom_ip_prefix.test.location
  resource_group_name = azurerm_custom_ip_prefix.test.resource_group_name

  cidr                          = azurerm_custom_ip_prefix.test.cidr
  roa_validity_end_date         = azurerm_custom_ip_prefix.test.roa_validity_end_date
  wan_validation_signed_message = azurerm_custom_ip_prefix.test.wan_validation_signed_message
}
`, r.ipv4(data))
}
