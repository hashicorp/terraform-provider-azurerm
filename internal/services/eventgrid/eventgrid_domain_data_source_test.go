package eventgrid_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type EventGridDomainDataSource struct {
}

func TestAccEventGridDomainDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_domain", "test")
	r := EventGridDomainDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridDomainDataSource_mapping(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_domain", "test")
	r := EventGridDomainDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.mapping(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("input_mapping_fields.0.topic").Exists(),
				check.That(data.ResourceName).Key("input_mapping_fields.0.topic").Exists(),
				check.That(data.ResourceName).Key("input_mapping_default_values.0.data_version").Exists(),
				check.That(data.ResourceName).Key("input_mapping_default_values.0.subject").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridDomainDataSource_basicWithTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithTags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").Exists(),
				check.That(data.ResourceName).Key("tags.foo").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridDomainDataSource_inboundIPRules(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.inboundIPRules(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("inbound_ip_rule.#").Exists(),
				check.That(data.ResourceName).Key("inbound_ip_rule.0.ip_mask").Exists(),
				check.That(data.ResourceName).Key("inbound_ip_rule.1.ip_mask").Exists(),
				check.That(data.ResourceName).Key("inbound_ip_rule.0.action").Exists(),
				check.That(data.ResourceName).Key("inbound_ip_rule.1.action").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func (EventGridDomainDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_domain" "test" {
  name                = azurerm_eventgrid_domain.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventGridDomainResource{}.basic(data))
}

func (EventGridDomainDataSource) mapping(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_domain" "test" {
  name                = azurerm_eventgrid_domain.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventGridDomainResource{}.mapping(data))
}

func (EventGridDomainDataSource) basicWithTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_domain" "test" {
  name                = azurerm_eventgrid_domain.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventGridDomainResource{}.basicWithTags(data))
}

func (EventGridDomainDataSource) inboundIPRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_eventgrid_domain" "test" {
  name                = azurerm_eventgrid_domain.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, EventGridDomainResource{}.inboundIPRules(data))
}
