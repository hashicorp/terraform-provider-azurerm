// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package dns_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/dns/2023-07-01-preview/recordsets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DnsDSRecordResource struct{}

func TestAccDnsDSRecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_ds_record", "test")
	r := DnsDSRecordResource{}

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

func TestAccDnsDSRecord_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_ds_record", "test")
	r := DnsDSRecordResource{}

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

func TestAccDnsDSRecord_updateRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_ds_record", "test")
	r := DnsDSRecordResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("record.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateRecords(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("record.#").HasValue("3"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDnsDSRecord_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_ds_record", "test")
	r := DnsDSRecordResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
			),
		},
		data.ImportStep(),
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

func (DnsDSRecordResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := recordsets.ParseRecordTypeID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Dns.RecordSets.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func (DnsDSRecordResource) basic(data acceptance.TestData) string {
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

resource "azurerm_dns_ds_record" "test" {
  name        = "myarecord%d"
  dns_zone_id = azurerm_dns_zone.test.id
  ttl         = 300

  record {
    algorithm    = 13
    key_tag      = 28237
    digest_type  = 2
    digest_value = "40F628643831D5EAF7D005D3237DE32F3F37AE6025C7891D202B0BAFA9924778"
  }

  record {
    algorithm    = 13
    key_tag      = 46872
    digest_type  = 2
    digest_value = "2C0BAC20EB5C8C315694CBEB62E56C71CDC0069D058A8B80992E6499D91DD247"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r DnsDSRecordResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dns_ds_record" "import" {
  name        = azurerm_dns_ds_record.test.name
  dns_zone_id = azurerm_dns_zone.test.id
  ttl         = 300

  record {
    algorithm    = 13
    key_tag      = 28237
    digest_type  = 2
    digest_value = "40F628643831D5EAF7D005D3237DE32F3F37AE6025C7891D202B0BAFA9924778"
  }

  record {
    algorithm    = 13
    key_tag      = 46872
    digest_type  = 2
    digest_value = "2C0BAC20EB5C8C315694CBEB62E56C71CDC0069D058A8B80992E6499D91DD247"
  }
}
`, r.basic(data))
}

func (DnsDSRecordResource) updateRecords(data acceptance.TestData) string {
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

resource "azurerm_dns_ds_record" "test" {
  name        = "myarecord%d"
  dns_zone_id = azurerm_dns_zone.test.id
  ttl         = 300

  record {
    algorithm    = 13
    key_tag      = 28237
    digest_type  = 2
    digest_value = "40F628643831D5EAF7D005D3237DE32F3F37AE6025C7891D202B0BAFA9924778"
  }

  record {
    algorithm    = 13
    key_tag      = 46872
    digest_type  = 2
    digest_value = "2C0BAC20EB5C8C315694CBEB62E56C71CDC0069D058A8B80992E6499D91DD247"
  }

  record {
    algorithm    = 13
    key_tag      = 20795
    digest_type  = 2
    digest_value = "55E20DB8044B0C6190A925598F08F8146C9A0D4F668F8CA5A7276EB54064C5E3"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (DnsDSRecordResource) withTags(data acceptance.TestData) string {
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

resource "azurerm_dns_ds_record" "test" {
  name        = "myarecord%d"
  dns_zone_id = azurerm_dns_zone.test.id
  ttl         = 300

  record {
    algorithm    = 13
    key_tag      = 28237
    digest_type  = 2
    digest_value = "40F628643831D5EAF7D005D3237DE32F3F37AE6025C7891D202B0BAFA9924778"
  }

  record {
    algorithm    = 13
    key_tag      = 46872
    digest_type  = 2
    digest_value = "2C0BAC20EB5C8C315694CBEB62E56C71CDC0069D058A8B80992E6499D91DD247"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (DnsDSRecordResource) withTagsUpdate(data acceptance.TestData) string {
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

resource "azurerm_dns_ds_record" "test" {
  name        = "myarecord%d"
  dns_zone_id = azurerm_dns_zone.test.id
  ttl         = 300

  record {
    algorithm    = 13
    key_tag      = 28237
    digest_type  = 2
    digest_value = "40F628643831D5EAF7D005D3237DE32F3F37AE6025C7891D202B0BAFA9924778"
  }

  record {
    algorithm    = 13
    key_tag      = 46872
    digest_type  = 2
    digest_value = "2C0BAC20EB5C8C315694CBEB62E56C71CDC0069D058A8B80992E6499D91DD247"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
