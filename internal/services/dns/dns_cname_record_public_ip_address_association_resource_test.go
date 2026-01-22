// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dns_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/publicipaddresses"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type DnsCnameRecordPublicIpAddressAssociationResource struct{}

func TestAccDnsCnameRecordPublicIpAddressAssociation_basic(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE or ARM_TEST_DATA_RESOURCE_GROUP is not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_dns_cname_record_public_ip_address_association", "test")
	r := DnsCnameRecordPublicIpAddressAssociationResource{}

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

func TestAccDnsCnameRecordPublicIpAddressAssociation_requiresImport(t *testing.T) {
	if os.Getenv("ARM_TEST_DNS_ZONE") == "" || os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP") == "" {
		t.Skip("Skipping as ARM_TEST_DNS_ZONE or ARM_TEST_DATA_RESOURCE_GROUP is not set")
	}

	data := acceptance.BuildTestData(t, "azurerm_dns_cname_record_public_ip_address_association", "test")
	r := DnsCnameRecordPublicIpAddressAssociationResource{}

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

func (r DnsCnameRecordPublicIpAddressAssociationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.DnsCnameRecordPublicIpAddressAssociationID(state.ID)
	if err != nil {
		return nil, err
	}

	publicIp, err := client.Network.PublicIPAddresses.Get(ctx, id.PublicIpAddressId, publicipaddresses.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(publicIp.HttpResponse) {
			return pointer.To(false), nil
		}
		return nil, fmt.Errorf("retrieving Public IP Address %s: %+v", id.PublicIpAddressId, err)
	}

	if publicIp.Model == nil || publicIp.Model.Properties == nil {
		return pointer.To(false), nil
	}

	if publicIp.Model.Properties.DnsSettings == nil || publicIp.Model.Properties.DnsSettings.ReverseFqdn == nil {
		return pointer.To(false), nil
	}

	return pointer.To(*publicIp.Model.Properties.DnsSettings.ReverseFqdn != ""), nil
}

func (r DnsCnameRecordPublicIpAddressAssociationResource) basic(data acceptance.TestData) string {
	dnsZone := os.Getenv("ARM_TEST_DNS_ZONE")
	dataResourceGroup := os.Getenv("ARM_TEST_DATA_RESOURCE_GROUP")

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

data "azurerm_dns_zone" "test" {
  name                = "%[3]s"
  resource_group_name = "%[4]s"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
  domain_name_label   = "acctestpip%[1]d"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "myarecord%[1]d"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 300
  record              = azurerm_public_ip.test.fqdn
}

resource "azurerm_dns_cname_record_public_ip_address_association" "test" {
  dns_cname_record_id  = azurerm_dns_cname_record.test.id
  public_ip_address_id = azurerm_public_ip.test.id
}
`, data.RandomInteger, data.Locations.Primary, dnsZone, dataResourceGroup)
}

func (r DnsCnameRecordPublicIpAddressAssociationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_dns_cname_record_public_ip_address_association" "import" {
  dns_cname_record_id  = azurerm_dns_cname_record_public_ip_address_association.test.dns_cname_record_id
  public_ip_address_id = azurerm_dns_cname_record_public_ip_address_association.test.public_ip_address_id
}
`, r.basic(data))
}
