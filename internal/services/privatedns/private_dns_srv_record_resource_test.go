// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package privatedns_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type PrivateDnsSrvRecordResource struct{}

func TestAccPrivateDnsSrvRecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_srv_record", "test")
	r := PrivateDnsSrvRecordResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("fqdn").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccPrivateDnsSrvRecord_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_srv_record", "test")
	r := PrivateDnsSrvRecordResource{}
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

func TestAccPrivateDnsSrvRecord_updateRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_srv_record", "test")
	r := PrivateDnsSrvRecordResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("record.#").HasValue("2"),
			),
		},
		{
			Config: r.updateRecords(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("record.#").HasValue("3"),
			),
		},
	})
}

func TestAccPrivateDnsSrvRecord_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_private_dns_srv_record", "test")
	r := PrivateDnsSrvRecordResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		{
			Config: r.withTagsUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func (t PrivateDnsSrvRecordResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := recordsets.ParseRecordTypeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.PrivateDns.RecordSetsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading Private DNS SRV Record (%s): %+v", id.String(), err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (PrivateDnsSrvRecordResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-prvdns-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_srv_record" "test" {
  name                = "testaccsrv%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }

  record {
    priority = 10
    weight   = 10
    port     = 8080
    target   = "target2.contoso.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r PrivateDnsSrvRecordResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_srv_record" "import" {
  name                = azurerm_private_dns_srv_record.test.name
  resource_group_name = azurerm_private_dns_srv_record.test.resource_group_name
  zone_name           = azurerm_private_dns_srv_record.test.zone_name
  ttl                 = 300
  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }
  record {
    priority = 10
    weight   = 10
    port     = 8080
    target   = "target2.contoso.com"
  }
}
`, r.basic(data))
}

func (PrivateDnsSrvRecordResource) updateRecords(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_srv_record" "test" {
  name                = "test%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }
  record {
    priority = 10
    weight   = 10
    port     = 8080
    target   = "target2.contoso.com"
  }
  record {
    priority = 20
    weight   = 100
    port     = 8080
    target   = "target3.contoso.com"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PrivateDnsSrvRecordResource) withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_srv_record" "test" {
  name                = "test%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }
  record {
    priority = 10
    weight   = 10
    port     = 8080
    target   = "target2.contoso.com"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (PrivateDnsSrvRecordResource) withTagsUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "testzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_private_dns_srv_record" "test" {
  name                = "test%d"
  resource_group_name = azurerm_resource_group.test.name
  zone_name           = azurerm_private_dns_zone.test.name
  ttl                 = 300
  record {
    priority = 1
    weight   = 5
    port     = 8080
    target   = "target1.contoso.com"
  }
  record {
    priority = 10
    weight   = 10
    port     = 8080
    target   = "target2.contoso.com"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
