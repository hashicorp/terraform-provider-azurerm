package connections_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/web/2016-06-01/connections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApiConnectionTestResource struct {
}

func TestAccApiConnection_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_connection", "test")
	r := ApiConnectionTestResource{}
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

func TestAccApiConnection_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_connection", "test")
	r := ApiConnectionTestResource{}
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

func TestAccApiConnection_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_api_connection", "test")
	r := ApiConnectionTestResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parameter_values"),
		{
			// re-add
			Config: r.completeUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("parameter_values"),
	})
}

func (t ApiConnectionTestResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := connections.ParseConnectionID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Connections.ConnectionsClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (t ApiConnectionTestResource) basic(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_connection" "test" {
  name                = "acctestconn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  managed_api_id      = data.azurerm_managed_api.test.id
}
`, template, data.RandomInteger)
}

func (t ApiConnectionTestResource) requiresImport(data acceptance.TestData) string {
	template := t.basic(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_api_connection" "import" {
  name                = azurerm_api_connection.test.name
  resource_group_name = azurerm_api_connection.test.resource_group_name
  managed_api_id      = azurerm_api_connection.test.managed_api_id
}
`, template)
}

func (t ApiConnectionTestResource) complete(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_connection" "test" {
  name                = "acctestconn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  managed_api_id      = data.azurerm_managed_api.test.id
  display_name        = "Example 1"

  parameter_values = {
    connectionString = azurerm_servicebus_namespace.test.default_primary_connection_string
  }

  tags = {
    Hello = "World"
  }

  lifecycle {
    ignore_changes = ["parameter_values"] # not returned from the API
  }
}
`, template, data.RandomInteger)
}

func (t ApiConnectionTestResource) completeUpdated(data acceptance.TestData) string {
	template := t.template(data)
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

%[1]s

resource "azurerm_api_connection" "test" {
  name                = "acctestconn-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  managed_api_id      = data.azurerm_managed_api.test.id
  display_name        = "Example 1"

  parameter_values = {
    connectionString = azurerm_servicebus_namespace.test.default_primary_connection_string
  }

  tags = {
    Another = "Tag"
    Hello   = "World"
  }

  lifecycle {
    ignore_changes = ["parameter_values"] # not returned from the API
  }
}
`, template, data.RandomInteger)
}

func (ApiConnectionTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-conn-%[1]d"
  location = %[2]q
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-conn-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_servicebus_namespace" "test" {
  name                = "acctestsbn-conn-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

data "azurerm_managed_api" "test" {
  name     = "servicebus"
  location = azurerm_resource_group.test.location
}
`, data.RandomInteger, data.Locations.Primary)
}
