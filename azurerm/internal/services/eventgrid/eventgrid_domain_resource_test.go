package eventgrid_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/eventgrid/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type EventGridDomainResource struct {
}

func TestAccEventGridDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_eventgrid_domain"),
		},
	})
}

func TestAccEventGridDomain_mapping(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.mapping(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("input_mapping_fields.0.topic").HasValue("test"),
				check.That(data.ResourceName).Key("input_mapping_fields.0.topic").HasValue("test"),
				check.That(data.ResourceName).Key("input_mapping_default_values.0.data_version").HasValue("1.0"),
				check.That(data.ResourceName).Key("input_mapping_default_values.0.subject").HasValue("DefaultSubject"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridDomain_basicWithTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basicWithTags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.foo").HasValue("bar"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridDomain_inboundIPRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.inboundIPRules(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("inbound_ip_rule.#").HasValue("2"),
				check.That(data.ResourceName).Key("inbound_ip_rule.0.ip_mask").HasValue("10.0.0.0/16"),
				check.That(data.ResourceName).Key("inbound_ip_rule.1.ip_mask").HasValue("10.1.0.0/16"),
				check.That(data.ResourceName).Key("inbound_ip_rule.0.action").HasValue("Allow"),
				check.That(data.ResourceName).Key("inbound_ip_rule.1.action").HasValue("Allow"),
			),
		},
		data.ImportStep(),
	})
}

func (EventGridDomainResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DomainID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.EventGrid.DomainsClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving EventGrid Domain %q (resource group: %q): %+v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.DomainProperties != nil), nil
}

func (EventGridDomainResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventGridDomainResource) requiresImport(data acceptance.TestData) string {
	template := EventGridDomainResource{}.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_eventgrid_domain" "import" {
  name                = azurerm_eventgrid_domain.test.name
  location            = azurerm_eventgrid_domain.test.location
  resource_group_name = azurerm_eventgrid_domain.test.resource_group_name
}
`, template)
}

func (EventGridDomainResource) mapping(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  input_schema = "CustomEventSchema"

  input_mapping_fields {
    topic      = "test"
    event_type = "test"
  }

  input_mapping_default_values {
    data_version = "1.0"
    subject      = "DefaultSubject"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventGridDomainResource) basicWithTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "foo" = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (EventGridDomainResource) inboundIPRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  public_network_access_enabled = true

  inbound_ip_rule {
    ip_mask = "10.0.0.0/16"
    action  = "Allow"
  }

  inbound_ip_rule {
    ip_mask = "10.1.0.0/16"
    action  = "Allow"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
