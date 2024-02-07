// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type EventGridDomainResource struct{}

func TestAccEventGridDomain_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridDomain_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("auto_create_topic_with_first_subscription").HasValue("true"),
				check.That(data.ResourceName).Key("auto_delete_topic_with_last_subscription").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("endpoint").Exists(),
				check.That(data.ResourceName).Key("primary_access_key").Exists(),
				check.That(data.ResourceName).Key("secondary_access_key").Exists(),
				check.That(data.ResourceName).Key("local_auth_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("auto_create_topic_with_first_subscription").HasValue("false"),
				check.That(data.ResourceName).Key("auto_delete_topic_with_last_subscription").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridDomain_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mapping(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithTags(data),
			Check: acceptance.ComposeTestCheckFunc(
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

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.inboundIPRules(data),
			Check: acceptance.ComposeTestCheckFunc(
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

func TestAccEventGridDomain_basicWithSystemManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithSystemManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("0"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccEventGridDomain_basicWithUserAssignedManagedIdentity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_eventgrid_domain", "test")
	r := EventGridDomainResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithUserAssignedManagedIdentity(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").IsEmpty(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").IsEmpty(),
			),
		},
		data.ImportStep(),
	})
}

func (EventGridDomainResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := domains.ParseDomainID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.EventGrid.Domains.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
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

func (EventGridDomainResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                                      = "acctesteg-%d"
  location                                  = azurerm_resource_group.test.location
  resource_group_name                       = azurerm_resource_group.test.name
  local_auth_enabled                        = false
  auto_create_topic_with_first_subscription = false
  auto_delete_topic_with_last_subscription  = false
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

func (EventGridDomainResource) basicWithSystemManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "SystemAssigned"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (EventGridDomainResource) basicWithUserAssignedManagedIdentity(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctesteg-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_eventgrid_domain" "test" {
  name                = "acctesteg-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  identity {
    type = "UserAssigned"
    identity_ids = [
      azurerm_user_assigned_identity.test.id
    ]
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (EventGridDomainResource) complete(data acceptance.TestData) string {
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

  input_schema = "CustomEventSchema"

  input_mapping_fields {
    topic      = "test"
    event_type = "test"
  }

  input_mapping_default_values {
    data_version = "1.0"
    subject      = "DefaultSubject"
  }

  tags = {
    "foo" = "bar"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
