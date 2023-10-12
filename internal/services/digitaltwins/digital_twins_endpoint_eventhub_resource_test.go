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

type DigitalTwinsEndpointEventHubResource struct{}

func TestAccDigitalTwinsEndpointEventHub_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventhub", "test")
	r := DigitalTwinsEndpointEventHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub_primary_connection_string", "eventhub_secondary_connection_string"),
	})
}

func TestAccDigitalTwinsEndpointEventHub_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventhub", "test")
	r := DigitalTwinsEndpointEventHubResource{}

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

func TestAccDigitalTwinsEndpointEventHub_updateEventHub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventhub", "test")
	r := DigitalTwinsEndpointEventHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub_primary_connection_string", "eventhub_secondary_connection_string"),
		{
			Config: r.updateEventHub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub_primary_connection_string", "eventhub_secondary_connection_string"),
		{
			Config: r.updateEventHubRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub_primary_connection_string", "eventhub_secondary_connection_string"),
	})
}

func TestAccDigitalTwinsEndpointEventHub_updateDeadLetter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventhub", "test")
	r := DigitalTwinsEndpointEventHubResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub_primary_connection_string", "eventhub_secondary_connection_string"),
		{
			Config: r.updateDeadLetter(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub_primary_connection_string", "eventhub_secondary_connection_string", "dead_letter_storage_secret"),
		{
			Config: r.updateDeadLetterRestore(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventhub_primary_connection_string", "eventhub_secondary_connection_string", "dead_letter_storage_secret"),
	})
}

func (r DigitalTwinsEndpointEventHubResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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

func (r DigitalTwinsEndpointEventHubResource) template(data acceptance.TestData) string {
	iR := DigitalTwinsInstanceResource{}
	digitalTwinsInstance := iR.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku = "Standard"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%[2]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  partition_count   = 2
  message_retention = 1
}

resource "azurerm_eventhub_authorization_rule" "test" {
  name                = "acctest-r%[2]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test.name
  resource_group_name = azurerm_resource_group.test.name

  listen = false
  send   = true
  manage = false
}
`, digitalTwinsInstance, data.RandomInteger)
}

func (r DigitalTwinsEndpointEventHubResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_endpoint_eventhub" "test" {
  name                                 = "acctest-EH-%d"
  digital_twins_id                     = azurerm_digital_twins_instance.test.id
  eventhub_primary_connection_string   = azurerm_eventhub_authorization_rule.test.primary_connection_string
  eventhub_secondary_connection_string = azurerm_eventhub_authorization_rule.test.secondary_connection_string
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointEventHubResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_endpoint_eventhub" "import" {
  name                                 = azurerm_digital_twins_endpoint_eventhub.test.name
  digital_twins_id                     = azurerm_digital_twins_endpoint_eventhub.test.digital_twins_id
  eventhub_primary_connection_string   = azurerm_digital_twins_endpoint_eventhub.test.eventhub_primary_connection_string
  eventhub_secondary_connection_string = azurerm_digital_twins_endpoint_eventhub.test.eventhub_secondary_connection_string
}
`, r.basic(data))
}

func (r DigitalTwinsEndpointEventHubResource) updateEventHub(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub" "test_alt" {
  name                = "acctesteventhub-alt-%[2]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  partition_count   = 2
  message_retention = 1
}

resource "azurerm_eventhub_authorization_rule" "test_alt" {
  name                = "acctest-r%[2]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test_alt.name
  resource_group_name = azurerm_resource_group.test.name

  listen = false
  send   = true
  manage = false
}

resource "azurerm_digital_twins_endpoint_eventhub" "test" {
  name                                 = "acctest-EH-%[2]d"
  digital_twins_id                     = azurerm_digital_twins_instance.test.id
  eventhub_primary_connection_string   = azurerm_eventhub_authorization_rule.test_alt.primary_connection_string
  eventhub_secondary_connection_string = azurerm_eventhub_authorization_rule.test_alt.secondary_connection_string
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointEventHubResource) updateEventHubRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_eventhub" "test_alt" {
  name                = "acctesteventhub-alt-%[2]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name

  partition_count   = 2
  message_retention = 1
}

resource "azurerm_eventhub_authorization_rule" "test_alt" {
  name                = "acctest-r%[2]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  eventhub_name       = azurerm_eventhub.test_alt.name
  resource_group_name = azurerm_resource_group.test.name

  listen = false
  send   = true
  manage = false
}

resource "azurerm_digital_twins_endpoint_eventhub" "test" {
  name                                 = "acctest-EH-%[2]d"
  digital_twins_id                     = azurerm_digital_twins_instance.test.id
  eventhub_primary_connection_string   = azurerm_eventhub_authorization_rule.test.primary_connection_string
  eventhub_secondary_connection_string = azurerm_eventhub_authorization_rule.test.secondary_connection_string
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointEventHubResource) updateDeadLetter(data acceptance.TestData) string {
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

resource "azurerm_digital_twins_endpoint_eventhub" "test" {
  name                                 = "acctest-EH-%[3]d"
  digital_twins_id                     = azurerm_digital_twins_instance.test.id
  eventhub_primary_connection_string   = azurerm_eventhub_authorization_rule.test.primary_connection_string
  eventhub_secondary_connection_string = azurerm_eventhub_authorization_rule.test.secondary_connection_string
  dead_letter_storage_secret           = "${azurerm_storage_container.test.id}?${azurerm_storage_account.test.primary_access_key}"
}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r DigitalTwinsEndpointEventHubResource) updateDeadLetterRestore(data acceptance.TestData) string {
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

resource "azurerm_digital_twins_endpoint_eventhub" "test" {
  name                                 = "acctest-EH-%[3]d"
  digital_twins_id                     = azurerm_digital_twins_instance.test.id
  eventhub_primary_connection_string   = azurerm_eventhub_authorization_rule.test.primary_connection_string
  eventhub_secondary_connection_string = azurerm_eventhub_authorization_rule.test.secondary_connection_string
}
`, r.template(data), data.RandomString, data.RandomInteger)
}
