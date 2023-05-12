package privatedns_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2018-09-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PrivateDnsZoneResource struct{}

func TestAccPrivateDnsZone_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone", "test")
	r := PrivateDnsZoneResource{}
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

func TestAccPrivateDnsZone_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone", "test")
	r := PrivateDnsZoneResource{}
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

func TestAccPrivateDnsZone_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone", "test")
	r := PrivateDnsZoneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateDnsZone_withSOARecord(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_zone", "test")
	r := PrivateDnsZoneResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withBasicSOARecord(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withCompletedSOARecord(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.withBasicSOARecord(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (t PrivateDnsZoneResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := privatezones.ParsePrivateDnsZoneID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.PrivateDns.PrivateZonesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Private DNS Zone (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (PrivateDnsZoneResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    Env1 = "Test1"
    Env2 = "Test2"
    Env3 = "Test3"
    Env4 = "Test4"
    Env5 = "Test5"
    Env6 = "Test6"
    Env7 = "Test7"
    Env8 = "Test8"
    Env9 = "Test9"
    Env10 = "Test0"
    Env11 = "Test11"
    Env12 = "Test12"
    Env13 = "Test13"
    Env14 = "Test14"
    Env15 = "Test15"
    Env16 = "Test16"
    Env17 = "Test17"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r PrivateDnsZoneResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_zone" "import" {
  name                = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_private_dns_zone.test.resource_group_name
}
`, r.basic(data))
}

func (PrivateDnsZoneResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PrivateDnsZoneResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PrivateDnsZoneResource) withBasicSOARecord(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatedns-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name

  soa_record {
    email = "testemail.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PrivateDnsZoneResource) withCompletedSOARecord(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-privatedns-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name

  soa_record {
    email        = "testemail.com"
    expire_time  = 2419200
    minimum_ttl  = 200
    refresh_time = 2600
    retry_time   = 200
    ttl          = 100

    tags = {
      ENv = "Test"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
