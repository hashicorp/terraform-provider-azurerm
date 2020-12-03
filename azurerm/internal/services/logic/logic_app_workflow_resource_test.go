package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type LogicAppWorkflowResource struct {
}

func TestAccLogicAppWorkflow_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("connector_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("connector_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_outbound_ip_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppWorkflow_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_logic_app_workflow"),
		},
	})
}

func TestAccLogicAppWorkflow_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("parameters.%").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.Source").HasValue("AcceptanceTests"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccLogicAppWorkflow_integrationAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.integrationAccount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("logic_app_integration_account_id"),
		{
			Config: r.integrationAccountUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("logic_app_integration_account_id"),
		{
			Config: r.empty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("logic_app_integration_account_id"),
	})
}

func TestAccLogicAppWorkflow_integrationServiceEnvironment(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_logic_app_workflow", "test")
	r := LogicAppWorkflowResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.integrationServiceEnvironment(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (LogicAppWorkflowResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	name := id.Path["workflows"]

	resp, err := clients.Logic.WorkflowClient.Get(ctx, id.ResourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Logic App Workflow %s (resource group: %s): %v", name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.WorkflowProperties != nil), nil
}

func (LogicAppWorkflowResource) empty(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r LogicAppWorkflowResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_workflow" "import" {
  name                = azurerm_logic_app_workflow.test.name
  location            = azurerm_logic_app_workflow.test.location
  resource_group_name = azurerm_logic_app_workflow.test.resource_group_name
}
`, r.empty(data))
}

func (LogicAppWorkflowResource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_logic_app_workflow" "test" {
  name                = "acctestlaw-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    "Source" = "AcceptanceTests"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (LogicAppWorkflowResource) integrationAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-IA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_logic_app_integration_account" "test2" {
  name                = "acctest-IA2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_logic_app_workflow" "test" {
  name                             = "acctestlaw-%[1]d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  logic_app_integration_account_id = azurerm_logic_app_integration_account.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (LogicAppWorkflowResource) integrationAccountUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%[1]d"
  location = "%[2]s"
}

resource "azurerm_logic_app_integration_account" "test" {
  name                = "acctest-IA-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_logic_app_integration_account" "test2" {
  name                = "acctest-IA2-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_logic_app_workflow" "test" {
  name                             = "acctestlaw-%[1]d"
  location                         = azurerm_resource_group.test.location
  resource_group_name              = azurerm_resource_group.test.name
  logic_app_integration_account_id = azurerm_logic_app_integration_account.test2.id
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r LogicAppWorkflowResource) integrationServiceEnvironment(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_logic_app_workflow" "test" {
  name                               = "acctestlaw-%d"
  location                           = azurerm_resource_group.test.location
  resource_group_name                = azurerm_resource_group.test.name
  integration_service_environment_id = azurerm_integration_service_environment.test.id
}
`, IntegrationServiceEnvironmentResource{}.basic(data), data.RandomInteger)
}
