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

type DigitalTwinsEndpointEventGridResource struct{}

func TestAccDigitalTwinsEndpointEventGrid_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventgrid", "test")
	r := DigitalTwinsEndpointEventGridResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
	})
}

func TestAccDigitalTwinsEndpointEventGrid_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventgrid", "test")
	r := DigitalTwinsEndpointEventGridResource{}

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

func TestAccDigitalTwinsEndpointEventGrid_updateEventGrid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventgrid", "test")
	r := DigitalTwinsEndpointEventGridResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
		{
			Config: r.updateEventGrid(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
		{
			Config: r.updateEventGridRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
	})
}

func TestAccDigitalTwinsEndpointEventGrid_updateDeadLetter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventgrid", "test")
	r := DigitalTwinsEndpointEventGridResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
		{
			Config: r.updateDeadLetter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key", "dead_letter_storage_secret"),
		{
			Config: r.updateDeadLetterRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key", "dead_letter_storage_secret"),
	})
}

func (r DigitalTwinsEndpointEventGridResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

	return utils.Bool(true), nil
}

func (r DigitalTwinsEndpointEventGridResource) template(data acceptance.TestData) string {
	iR := DigitalTwinsInstanceResource{}
	digitalTwinsInstance := iR.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_eventgrid_topic" "test" {
  name                = "acctesteg-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, digitalTwinsInstance, data.RandomInteger)
}

func (r DigitalTwinsEndpointEventGridResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_endpoint_eventgrid" "test" {
  name                                 = "acctest-EG-%d"
  digital_twins_id                     = azurerm_digital_twins_instance.test.id
  eventgrid_topic_endpoint             = azurerm_eventgrid_topic.test.endpoint
  eventgrid_topic_primary_access_key   = azurerm_eventgrid_topic.test.primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_eventgrid_topic.test.secondary_access_key
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointEventGridResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_endpoint_eventgrid" "import" {
  name                                 = azurerm_digital_twins_endpoint_eventgrid.test.name
  digital_twins_id                     = azurerm_digital_twins_endpoint_eventgrid.test.digital_twins_id
  eventgrid_topic_endpoint             = azurerm_digital_twins_endpoint_eventgrid.test.eventgrid_topic_endpoint
  eventgrid_topic_primary_access_key   = azurerm_digital_twins_endpoint_eventgrid.test.eventgrid_topic_primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_digital_twins_endpoint_eventgrid.test.eventgrid_topic_secondary_access_key
}
`, r.basic(data))
}

func (r DigitalTwinsEndpointEventGridResource) updateEventGrid(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_eventgrid_topic" "test_alt" {
  name                = "acctesteg-alt-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_digital_twins_endpoint_eventgrid" "test" {
  name                                 = "acctest-EG-%[2]d"
  digital_twins_id                     = azurerm_digital_twins_instance.test.id
  eventgrid_topic_endpoint             = azurerm_eventgrid_topic.test_alt.endpoint
  eventgrid_topic_primary_access_key   = azurerm_eventgrid_topic.test_alt.primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_eventgrid_topic.test_alt.secondary_access_key
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointEventGridResource) updateEventGridRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_eventgrid_topic" "test_alt" {
  name                = "acctesteg-alt-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_digital_twins_endpoint_eventgrid" "test" {
  name                                 = "acctest-EG-%[2]d"
  digital_twins_id                     = azurerm_digital_twins_instance.test.id
  eventgrid_topic_endpoint             = azurerm_eventgrid_topic.test.endpoint
  eventgrid_topic_primary_access_key   = azurerm_eventgrid_topic.test.primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_eventgrid_topic.test.secondary_access_key
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointEventGridResource) updateDeadLetter(data acceptance.TestData) string {
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

resource "azurerm_digital_twins_endpoint_eventgrid" "test" {
  name                                 = "acctest-EG-%[3]d"
  digital_twins_id                     = azurerm_digital_twins_instance.test.id
  eventgrid_topic_endpoint             = azurerm_eventgrid_topic.test.endpoint
  eventgrid_topic_primary_access_key   = azurerm_eventgrid_topic.test.primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_eventgrid_topic.test.secondary_access_key
  dead_letter_storage_secret           = "${azurerm_storage_container.test.id}?${azurerm_storage_account.test.primary_access_key}"

}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r DigitalTwinsEndpointEventGridResource) updateDeadLetterRestore(data acceptance.TestData) string {
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

resource "azurerm_digital_twins_endpoint_eventgrid" "test" {
  name                                 = "acctest-EG-%[3]d"
  digital_twins_id                     = azurerm_digital_twins_instance.test.id
  eventgrid_topic_endpoint             = azurerm_eventgrid_topic.test.endpoint
  eventgrid_topic_primary_access_key   = azurerm_eventgrid_topic.test.primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_eventgrid_topic.test.secondary_access_key

}
`, r.template(data), data.RandomString, data.RandomInteger)
}
