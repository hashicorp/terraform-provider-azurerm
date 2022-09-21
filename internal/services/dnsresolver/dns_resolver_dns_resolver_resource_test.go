package dnsresolver_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/dnsresolver/2022-07-01/dnsresolvers"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DNSResolverDnsResolverResource struct{}

func TestAccDNSResolverDnsResolver_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_resolver", "test")
	r := DNSResolverDnsResolverResource{}
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

func TestAccDNSResolverDnsResolver_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_resolver", "test")
	r := DNSResolverDnsResolverResource{}
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

func TestAccDNSResolverDnsResolver_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_resolver", "test")
	r := DNSResolverDnsResolverResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccDNSResolverDnsResolver_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_dns_resolver", "test")
	r := DNSResolverDnsResolverResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r DNSResolverDnsResolverResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := dnsresolvers.ParseDnsResolverID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.DNSResolver.DnsResolversClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r DNSResolverDnsResolverResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-rg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r DNSResolverDnsResolverResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_dns_resolver" "test" {
  name                = "acctest-dr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_network_id  = azurerm_virtual_network.test.id
}
`, template, data.RandomInteger)
}

func (r DNSResolverDnsResolverResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_dns_resolver" "import" {
  name                = azurerm_dns_resolver.test.name
  resource_group_name = azurerm_dns_resolver.test.resource_group_name
  location            = azurerm_dns_resolver.test.location
  virtual_network_id  = azurerm_dns_resolver.test.virtual_network_id
}
`, config)
}

func (r DNSResolverDnsResolverResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_dns_resolver" "test" {
  name                = "acctest-dr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_network_id  = azurerm_virtual_network.test.id
  tags = {
    key = "value"
  }
}
`, template, data.RandomInteger)
}

func (r DNSResolverDnsResolverResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_dns_resolver" "test" {
  name                = "acctest-dr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_network_id  = azurerm_virtual_network.test.id
  tags = {
    key = "updated value"
  }

}
`, template, data.RandomInteger)
}
