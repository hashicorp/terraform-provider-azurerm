// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package digitaltwins_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2023-01-31/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type DigitalTwinsEndpointServiceBusResource struct{}

func TestAccDigitalTwinsEndpointServicebus_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_servicebus", "test")
	r := DigitalTwinsEndpointServiceBusResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("servicebus_primary_connection_string", "servicebus_secondary_connection_string"),
	})
}

func TestAccDigitalTwinsEndpointServicebus_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_servicebus", "test")
	r := DigitalTwinsEndpointServiceBusResource{}

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

func TestAccDigitalTwinsEndpointServicebus_updateServiceBus(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_servicebus", "test")
	r := DigitalTwinsEndpointServiceBusResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("servicebus_primary_connection_string", "servicebus_secondary_connection_string"),
		{
			Config: r.updateServiceBus(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("servicebus_primary_connection_string", "servicebus_secondary_connection_string"),
		{
			Config: r.updateServiceBusRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("servicebus_primary_connection_string", "servicebus_secondary_connection_string"),
	})
}

func TestAccDigitalTwinsEndpointServicebus_updateDeadLetter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_servicebus", "test")
	r := DigitalTwinsEndpointServiceBusResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("servicebus_primary_connection_string", "servicebus_secondary_connection_string"),
		{
			Config: r.updateDeadLetter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("servicebus_primary_connection_string", "servicebus_secondary_connection_string", "dead_letter_storage_secret"),
		{
			Config: r.updateDeadLetterRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("servicebus_primary_connection_string", "servicebus_secondary_connection_string", "dead_letter_storage_secret"),
	})
}

func (r DigitalTwinsEndpointServiceBusResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := endpoints.ParseEndpointID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DigitalTwins.EndpointClient.DigitalTwinsEndpointGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r DigitalTwinsEndpointServiceBusResource) template(data acceptance.TestData) string {
	iR := DigitalTwinsInstanceResource{}
	digitalTwinsInstance := iR.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestservicebusnamespace-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
}

resource "azurerm_servicebus_topic" "test" {
  name         = "acctestservicebustopic-%[2]d"
  namespace_id = azurerm_servicebus_namespace.test.id
}

resource "azurerm_servicebus_topic_authorization_rule" "test" {
  name     = "acctest-rule-%[2]d"
  topic_id = azurerm_servicebus_topic.test.id

  listen = false
  send   = true
  manage = false
}
`, digitalTwinsInstance, data.RandomInteger)
}

func (r DigitalTwinsEndpointServiceBusResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_endpoint_servicebus" "test" {
  name                                   = "acctest-EndpointSB-%d"
  digital_twins_id                       = azurerm_digital_twins_instance.test.id
  servicebus_primary_connection_string   = azurerm_servicebus_topic_authorization_rule.test.primary_connection_string
  servicebus_secondary_connection_string = azurerm_servicebus_topic_authorization_rule.test.secondary_connection_string
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointServiceBusResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_endpoint_servicebus" "import" {
  name                                   = azurerm_digital_twins_endpoint_servicebus.test.name
  digital_twins_id                       = azurerm_digital_twins_endpoint_servicebus.test.digital_twins_id
  servicebus_primary_connection_string   = azurerm_digital_twins_endpoint_servicebus.test.servicebus_primary_connection_string
  servicebus_secondary_connection_string = azurerm_digital_twins_endpoint_servicebus.test.servicebus_secondary_connection_string
}
`, r.basic(data))
}

func (r DigitalTwinsEndpointServiceBusResource) updateServiceBus(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_servicebus_namespace" "test_alt" {
  name                = "acctestservicebusnamespace-alt-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_servicebus_topic" "test_alt" {
  name         = "acctestservicebustopic-alt-%[2]d"
  namespace_id = azurerm_servicebus_namespace.test.id
}

resource "azurerm_servicebus_topic_authorization_rule" "test_alt" {
  name     = "acctest-rule-alt-%[2]d"
  topic_id = azurerm_servicebus_topic.test.id

  listen = false
  send   = true
  manage = false
}

resource "azurerm_digital_twins_endpoint_servicebus" "test" {
  name                                   = "acctest-EndpointSB-%[2]d"
  digital_twins_id                       = azurerm_digital_twins_instance.test.id
  servicebus_primary_connection_string   = azurerm_servicebus_topic_authorization_rule.test_alt.primary_connection_string
  servicebus_secondary_connection_string = azurerm_servicebus_topic_authorization_rule.test_alt.secondary_connection_string
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointServiceBusResource) updateServiceBusRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_servicebus_namespace" "test_alt" {
  name                = "acctestservicebusnamespace-alt-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_servicebus_topic" "test_alt" {
  name         = "acctestservicebustopic-alt-%[2]d"
  namespace_id = azurerm_servicebus_namespace.test.id
}

resource "azurerm_servicebus_topic_authorization_rule" "test_alt" {
  name     = "acctest-rule-alt-%[2]d"
  topic_id = azurerm_servicebus_topic.test.id

  listen = false
  send   = true
  manage = false
}

resource "azurerm_digital_twins_endpoint_servicebus" "test" {
  name                                   = "acctest-EndpointSB-%[2]d"
  digital_twins_id                       = azurerm_digital_twins_instance.test.id
  servicebus_primary_connection_string   = azurerm_servicebus_topic_authorization_rule.test.primary_connection_string
  servicebus_secondary_connection_string = azurerm_servicebus_topic_authorization_rule.test.secondary_connection_string
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointServiceBusResource) updateDeadLetter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_digital_twins_endpoint_servicebus" "test" {
  name                                   = "acctest-EndpointSB-%[3]d"
  digital_twins_id                       = azurerm_digital_twins_instance.test.id
  servicebus_primary_connection_string   = azurerm_servicebus_topic_authorization_rule.test.primary_connection_string
  servicebus_secondary_connection_string = azurerm_servicebus_topic_authorization_rule.test.secondary_connection_string
  dead_letter_storage_secret             = "${azurerm_storage_container.test.id}?${azurerm_storage_account.test.primary_access_key}"
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r DigitalTwinsEndpointServiceBusResource) updateDeadLetterRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctestacc%[2]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "test" {
  name                  = "vhds"
  storage_account_name  = azurerm_storage_account.test.name
  container_access_type = "private"
}

resource "azurerm_digital_twins_endpoint_servicebus" "test" {
  name                                   = "acctest-EndpointSB-%[3]d"
  digital_twins_id                       = azurerm_digital_twins_instance.test.id
  servicebus_primary_connection_string   = azurerm_servicebus_topic_authorization_rule.test.primary_connection_string
  servicebus_secondary_connection_string = azurerm_servicebus_topic_authorization_rule.test.secondary_connection_string
}
`, r.template(data), data.RandomString, data.RandomInteger)
}
