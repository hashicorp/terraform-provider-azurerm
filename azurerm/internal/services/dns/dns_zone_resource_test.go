package dns_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/dns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DnsZoneResource struct {
}

func TestAccDnsZone_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_zone", "test")
	r := DnsZoneResource{}

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

func TestAccDnsZone_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_zone", "test")
	r := DnsZoneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_dns_zone"),
		},
	})
}

func TestAccDnsZone_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_zone", "test")
	r := DnsZoneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDnsZone_withSOARecord(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_zone", "test")
	r := DnsZoneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.withBasicSOARecord(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCompletedSOARecord(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBasicSOARecord(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (DnsZoneResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DnsZoneID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Dns.ZonesClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving DNS zone %s (resource group: %s): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.ZoneProperties != nil), nil
}

func (DnsZoneResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DnsZoneResource) requiresImport(data acceptance.TestData) string {
	template := DnsZoneResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_dns_zone" "import" {
  name                = azurerm_dns_zone.test.name
  resource_group_name = azurerm_dns_zone.test.resource_group_name
}
`, template)
}

func (DnsZoneResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DnsZoneResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DnsZoneResource) withBasicSOARecord(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dns-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name

  soa_record {
    email     = "testemail.com"
    host_name = "testhost.contoso.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (DnsZoneResource) withCompletedSOARecord(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-dns-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name

  soa_record {
    email         = "testemail.com"
    host_name     = "testhost.contoso.com"
    expire_time   = 2419200
    minimum_ttl   = 200
    refresh_time  = 2600
    retry_time    = 200
    serial_number = 1
    ttl           = 100

    tags = {
      ENv = "Test"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
