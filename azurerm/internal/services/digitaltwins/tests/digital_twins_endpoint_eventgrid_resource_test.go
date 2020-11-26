package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/digitaltwins/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type DigitalTwinsEndpointEventGrid struct{}

func TestAccAzureRMDigitalTwinsEndpointEventGrid_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventgrid", "test")
	r := DigitalTwinsEndpointEventGrid{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
	})
}

func TestAccAzureRMDigitalTwinsEndpointEventGrid_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventgrid", "test")
	r := DigitalTwinsEndpointEventGrid{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAzureRMDigitalTwinsEndpointEventGrid_updateEventGrid(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventgrid", "test")
	r := DigitalTwinsEndpointEventGrid{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
		{
			Config: r.updateEventGrid(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
		{
			Config: r.updateEventGridRestore(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
	})
}

func TestAccAzureRMDigitalTwinsEndpointEventGrid_updateDeadLetter(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_digital_twins_endpoint_eventgrid", "test")
	r := DigitalTwinsEndpointEventGrid{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
		{
			Config: r.updateDeadLetter(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
		{
			Config: r.updateDeadLetterRestore(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("eventgrid_topic_endpoint", "eventgrid_topic_primary_access_key", "eventgrid_topic_secondary_access_key"),
	})
}

func (r DigitalTwinsEndpointEventGrid) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.DigitalTwinsEndpointID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.DigitalTwins.EndpointClient.Get(ctx, id.ResourceGroup, id.DigitalTwinsInstanceName, id.EndpointName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving Digital Twins EventGrid Endpoint %q (Resource Group %q / Digital Twins Instance Name %q): %+v", id.EndpointName, id.ResourceGroup, id.DigitalTwinsInstanceName, err)
	}

	return utils.Bool(true), nil
}

func (r DigitalTwinsEndpointEventGrid) template(data acceptance.TestData) string {
	digitalTwins := testAccAzureRMDigitalTwinsInstance_basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_eventgrid_topic" "test" {
  name                = "acctesteg-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, digitalTwins, data.RandomInteger)
}

func (r DigitalTwinsEndpointEventGrid) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_endpoint_eventgrid" "test" {
  name                                 = "acctest-EG-%d"
  digital_twins_instance_id            = azurerm_digital_twins_instance.test.id
  eventgrid_topic_endpoint             = azurerm_eventgrid_topic.test.endpoint
  eventgrid_topic_primary_access_key   = azurerm_eventgrid_topic.test.primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_eventgrid_topic.test.secondary_access_key
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointEventGrid) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_digital_twins_endpoint_eventgrid" "import" {
  name                                 = azurerm_digital_twins_endpoint_eventgrid.test.name
  digital_twins_instance_id            = azurerm_digital_twins_endpoint_eventgrid.test.digital_twins_instance_id
  eventgrid_topic_endpoint             = azurerm_digital_twins_endpoint_eventgrid.test.eventgrid_topic_endpoint
  eventgrid_topic_primary_access_key   = azurerm_digital_twins_endpoint_eventgrid.test.eventgrid_topic_primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_digital_twins_endpoint_eventgrid.test.eventgrid_topic_secondary_access_key
}
`, r.basic(data))
}

func (r DigitalTwinsEndpointEventGrid) updateEventGrid(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_eventgrid_topic" "test_alt" {
  name                = "acctesteg-alt-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_digital_twins_endpoint_eventgrid" "test" {
  name                                 = "acctest-EG-%[2]d"
  digital_twins_instance_id            = azurerm_digital_twins_instance.test.id
  eventgrid_topic_endpoint             = azurerm_eventgrid_topic.test_alt.endpoint
  eventgrid_topic_primary_access_key   = azurerm_eventgrid_topic.test_alt.primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_eventgrid_topic.test_alt.secondary_access_key
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointEventGrid) updateEventGridRestore(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_eventgrid_topic" "test_alt" {
  name                = "acctesteg-alt-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_digital_twins_endpoint_eventgrid" "test" {
  name                                 = "acctest-EG-%[2]d"
  digital_twins_instance_id            = azurerm_digital_twins_instance.test.id
  eventgrid_topic_endpoint             = azurerm_eventgrid_topic.test.endpoint
  eventgrid_topic_primary_access_key   = azurerm_eventgrid_topic.test.primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_eventgrid_topic.test.secondary_access_key
}
`, r.template(data), data.RandomInteger)
}

func (r DigitalTwinsEndpointEventGrid) updateDeadLetter(data acceptance.TestData) string {
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
  digital_twins_instance_id            = azurerm_digital_twins_instance.test.id
  eventgrid_topic_endpoint             = azurerm_eventgrid_topic.test.endpoint
  eventgrid_topic_primary_access_key   = azurerm_eventgrid_topic.test.primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_eventgrid_topic.test.secondary_access_key
  dead_letter_storage_secret           = "${azurerm_storage_container.test.id}?${azurerm_storage_account.test.primary_access_key}"

}
`, r.template(data), data.RandomString, data.RandomInteger)
}

func (r DigitalTwinsEndpointEventGrid) updateDeadLetterRestore(data acceptance.TestData) string {
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
  digital_twins_instance_id            = azurerm_digital_twins_instance.test.id
  eventgrid_topic_endpoint             = azurerm_eventgrid_topic.test.endpoint
  eventgrid_topic_primary_access_key   = azurerm_eventgrid_topic.test.primary_access_key
  eventgrid_topic_secondary_access_key = azurerm_eventgrid_topic.test.secondary_access_key

}
`, r.template(data), data.RandomString, data.RandomInteger)
}
