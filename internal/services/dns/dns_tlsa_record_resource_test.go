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

type DnsTLSARecordResource struct{}

func TestAccDnsTLSARecord_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_tlsa_record", "test")
	r := DnsTLSARecordResource{}

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

func TestAccDnsTLSARecord_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_tlsa_record", "test")
	r := DnsTLSARecordResource{}

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

func TestAccDnsTLSARecord_updateRecords(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_tlsa_record", "test")
	r := DnsTLSARecordResource{}

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

func TestAccDnsTLSARecord_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_tlsa_record", "test")
	r := DnsTLSARecordResource{}

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

func (DnsTLSARecordResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (DnsTLSARecordResource) basic(data acceptance.TestData) string {
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

resource "azurerm_dns_tlsa_record" "test" {
  name        = "myarecord%d"
  dns_zone_id = azurerm_dns_zone.test.id
  ttl         = 300

  record {
    matching_type                = 1
    selector                     = 1
    usage                        = 3
    certificate_association_data = "370C66FD4A0673CE1B62E76B819835DABB20702E4497CB10AFFE46E8135381E7"
  }

  record {
    matching_type                = 1
    selector                     = 0
    usage                        = 0
    certificate_association_data = "d2abde240d7cd3ee6b4b28c54df034b97983a1d16e8a410e4561cb106618e971"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r DnsTLSARecordResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dns_tlsa_record" "import" {
  name        = azurerm_dns_tlsa_record.test.name
  dns_zone_id = azurerm_dns_zone.test.id
  ttl         = 300

  record {
    matching_type                = 1
    selector                     = 1
    usage                        = 3
    certificate_association_data = "370C66FD4A0673CE1B62E76B819835DABB20702E4497CB10AFFE46E8135381E7"
  }

  record {
    matching_type                = 1
    selector                     = 0
    usage                        = 0
    certificate_association_data = "d2abde240d7cd3ee6b4b28c54df034b97983a1d16e8a410e4561cb106618e971"
  }
}
`, r.basic(data))
}

func (DnsTLSARecordResource) updateRecords(data acceptance.TestData) string {
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

resource "azurerm_dns_tlsa_record" "test" {
  name        = "myarecord%d"
  dns_zone_id = azurerm_dns_zone.test.id
  ttl         = 300

  record {
    matching_type                = 1
    selector                     = 1
    usage                        = 3
    certificate_association_data = "370C66FD4A0673CE1B62E76B819835DABB20702E4497CB10AFFE46E8135381E7"
  }

  record {
    matching_type                = 1
    selector                     = 0
    usage                        = 0
    certificate_association_data = "d2abde240d7cd3ee6b4b28c54df034b97983a1d16e8a410e4561cb106618e971"
  }

  record {
    matching_type                = 1
    selector                     = 1
    usage                        = 3
    certificate_association_data = "0C72AC70B745AC19998811B131D662C9AC69DBDBE7CB23E5B514B56664C5D3D6"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (DnsTLSARecordResource) withTags(data acceptance.TestData) string {
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

resource "azurerm_dns_tlsa_record" "test" {
  name        = "myarecord%d"
  dns_zone_id = azurerm_dns_zone.test.id
  ttl         = 300

  record {
    matching_type                = 1
    selector                     = 1
    usage                        = 3
    certificate_association_data = "370C66FD4A0673CE1B62E76B819835DABB20702E4497CB10AFFE46E8135381E7"
  }

  record {
    matching_type                = 1
    selector                     = 0
    usage                        = 0
    certificate_association_data = "d2abde240d7cd3ee6b4b28c54df034b97983a1d16e8a410e4561cb106618e971"
  }

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (DnsTLSARecordResource) withTagsUpdate(data acceptance.TestData) string {
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

resource "azurerm_dns_tlsa_record" "test" {
  name        = "myarecord%d"
  dns_zone_id = azurerm_dns_zone.test.id
  ttl         = 300

  record {
    matching_type                = 1
    selector                     = 1
    usage                        = 3
    certificate_association_data = "370C66FD4A0673CE1B62E76B819835DABB20702E4497CB10AFFE46E8135381E7"
  }

  record {
    matching_type                = 1
    selector                     = 0
    usage                        = 0
    certificate_association_data = "d2abde240d7cd3ee6b4b28c54df034b97983a1d16e8a410e4561cb106618e971"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
