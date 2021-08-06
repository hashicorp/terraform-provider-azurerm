package cdn_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func preCheck(t *testing.T) {
	variables := []string{
		"ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME",
		"ARM_TEST_DNS_ZONE_NAME",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Skipf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

type CdnEndpointCustomDomainResource struct {
	DNSZoneRG     string
	DNSZoneName   string
	SubDomainName string
}

func NewCdnEndpointCustomDomainResource(dnsZoneRg, dnsZoneName string) CdnEndpointCustomDomainResource {
	return CdnEndpointCustomDomainResource{
		DNSZoneRG:     dnsZoneRg,
		DNSZoneName:   dnsZoneName,
		SubDomainName: acceptance.RandString(3),
	}
}

func TestAccCdnEndpointCustomDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")

	preCheck(t)

	r := NewCdnEndpointCustomDomainResource(os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"))

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

func TestAccCdnEndpointCustomDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cdn_endpoint_custom_domain", "test")

	preCheck(t)

	r := NewCdnEndpointCustomDomainResource(os.Getenv("ARM_TEST_DNS_ZONE_RESOURCE_GROUP_NAME"), os.Getenv("ARM_TEST_DNS_ZONE_NAME"))

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

func (r CdnEndpointCustomDomainResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Cdn.CustomDomainsClient.Get(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %q: %+v", id, err)
	}
	return utils.Bool(true), nil
}

func (r CdnEndpointCustomDomainResource) Destroy(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.CustomDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	c := client.Cdn.CustomDomainsClient
	future, err := c.Delete(ctx, id.ResourceGroup, id.ProfileName, id.EndpointName, id.Name)
	if err != nil {
		return nil, fmt.Errorf("deleting %q: %+v", id, err)
	}
	if err := future.WaitForCompletionRef(ctx, c.Client); err != nil {
		return nil, fmt.Errorf("waiting for deletion of %q: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r CdnEndpointCustomDomainResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_endpoint_custom_domain" "test" {
  name            = "acceptance-customdomain"
  cdn_endpoint_id = azurerm_cdn_endpoint.test.id
  host_name       = "${azurerm_dns_cname_record.test.name}.${data.azurerm_dns_zone.test.name}"
}
`, template)
}

func (r CdnEndpointCustomDomainResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cdn_endpoint_custom_domain" "import" {
  name            = azurerm_cdn_endpoint_custom_domain.test.name
  cdn_endpoint_id = azurerm_cdn_endpoint_custom_domain.test.cdn_endpoint_id
  host_name       = azurerm_cdn_endpoint_custom_domain.test.host_name
}
`, template)
}

func (r CdnEndpointCustomDomainResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acceptanceRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acceptancesa%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acceptancecdnprof%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

resource "azurerm_cdn_endpoint" "test" {
  name                = "acceptancecdnend%[1]d"
  profile_name        = azurerm_cdn_profile.test.name
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  origin {
    name      = "test"
    host_name = azurerm_storage_account.test.primary_blob_host
  }
}

data "azurerm_dns_zone" "test" {
  name                = "%[3]s"
  resource_group_name = "%[4]s"
}

resource "azurerm_dns_cname_record" "test" {
  name                = "%[5]s"
  zone_name           = data.azurerm_dns_zone.test.name
  resource_group_name = data.azurerm_dns_zone.test.resource_group_name
  ttl                 = 3600
  target_resource_id  = azurerm_cdn_endpoint.test.id
}
`, data.RandomIntOfLength(8), data.Locations.Primary, r.DNSZoneName, r.DNSZoneRG, r.SubDomainName)
}
